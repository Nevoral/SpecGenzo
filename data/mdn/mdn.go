package mdn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Nevoral/SpecGenzo/model"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type dataOption int

const (
	htmlConfig dataOption = iota
	svgConfig
	mathConfig
	htmlNodesConfig
	htmlGlobalAttrConfig
	htmlVoidNodesConfig
	htmlEventAttrConfig
	svgNodesConfig
	svgGlobalAttrConfig
	svgVoidNodesConfig
	svgEventAttrConfig
	mathNodesConfig
	mathGlobalAttrConfig
	mathVoidNodesConfig
	mathEventAttrConfig
	specificRequest
)

type MDNSource struct {
	baseURL  string
	path     string
	filename string
	specData []byte
	mdnData  *MDNData
}

func ScrapeMDNSource(url ...string) *MDNSource {
	m := MDNSource{
		baseURL:  "https://unpkg.com/@mdn/browser-compat-data/data.json",
		path:     "./data",
		filename: "mdn_compact_data.json",
	}
	if url != nil && len(url) > 0 {
		m.baseURL = url[0]
	}
	return &m
}

func (m *MDNSource) DownloadData() error {
	err := downloadData("./"+filepath.Join(m.path, m.filename), m.baseURL)
	if err != nil {
		if !(strings.Contains(err.Error(), "Failed to fetch the page:") || strings.Contains(err.Error(), "Error fetching page: status code")) {
			if err = downloadData("./"+filepath.Join(m.path, m.filename), m.baseURL); err != nil {
				return err
			}
		}
	}

	if m.specData, err = readJsonFile("./" + filepath.Join(m.path, m.filename)); err != nil {
		panic(err)
	}
	return nil
}

func (m *MDNSource) LoadData() (err error) {
	m.specData, err = readJsonFile("./" + filepath.Join(m.path, m.filename))
	return
}

func (m *MDNSource) ParseData() error {
	var v MDNData
	if err := json.Unmarshal(m.specData, &v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	m.mdnData = &v
	return nil
}

func (m *MDNSource) ExtractWebSpecification() (*model.WebSpecification, error) {
	if err := m.LoadData(); err != nil {
		return nil, err
	}
	if err := m.ParseData(); err != nil {
		return nil, err
	}
	html, err := extratactListOfElements(m.mdnData.Html.Elements)
	if err != nil {
		return nil, err
	}
	html.SortAllSlicesAscending()
	svg, err := extratactListOfElements(m.mdnData.Svg.Elements)
	if err != nil {
		return nil, err
	}
	svg.SortAllSlicesAscending()
	math, err := extratactListOfElements(m.mdnData.MathMl.Elements)
	if err != nil {
		return nil, err
	}
	math.SortAllSlicesAscending()
	var webSpec = model.WebSpecification{
		Version: "0.5.1",
		Spec: map[model.Namespace]*model.NamespaceConfig{
			model.HTML: html,
			model.SVG:  svg,
			model.MATH: math,
		},
	}
	return &webSpec, nil
}

func extratactListOfElements(data map[string]any) (*model.NamespaceConfig, error) {
	var spec = model.NamespaceConfig{
		Nodes: []*model.NodeConfig{},
		AttributesCategories: map[model.AttributeCategories][]*model.AttributeConfig{
			model.GlobalAttributes: {},
			model.WindowActions:    {},
			model.DocumentActions:  {},
		},
	}
	var err error
	for element, value := range data {
		var specAttr []*model.AttributeConfig
		var elementCompat *CompatData
		switch q := value.(type) {
		case map[string]any:
			for attr, val := range q {
				if attr == "__compat" {
					elementCompat, err = extractCompat(val)
					if err != nil {
						return nil, err
					}
					continue
				}
				var attrCompat *CompatData
				switch w := val.(type) {
				case map[string]any:
					for attrSpec, attrValue := range w {
						if attrSpec == "__compat" {
							attrCompat, err = extractCompat(attrValue)
							if err != nil {
								return nil, err
							}
						}
					}
				}
				var docUrl string
				switch s := attrCompat.SpecURL.(type) {
				case string:
					docUrl = s
				case []string:
					if len(s) > 0 {
						docUrl = s[0]
					}
				}
				specAttr = append(specAttr, &model.AttributeConfig{
					Name:             attr,
					Boolean:          false,
					Tags:             model.RegisterTag(attrCompat.Status),
					Comment:          "",
					DocumentationURL: docUrl,
					InitialValue:     "",
					SupportedValues:  map[string]model.Comment{},
				})
			}
		}
		spec.Nodes = append(spec.Nodes, &model.NodeConfig{
			Name:                       element,
			NodeType:                   model.FullTagType,
			Tags:                       model.RegisterTag(elementCompat.Status),
			Comment:                    "",
			DocumentationURL:           elementCompat.MDNURL,
			AttributesCategorySupports: map[model.AttributeCategories][]string{},
			SpecificAttributes:         specAttr,
			SupportedChildrenTags:      []string{},
		})
	}
	return &spec, nil
}

// func extractGlobalAttributes(data map[string]any) []*model.AttributeConfig {
//
// }

func extractCompat(data any) (*CompatData, error) {
	var compatData CompatData
	switch compat := data.(type) {
	case map[string]any:
		jsonData, err := json.Marshal(compat)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonData, &compatData)
		if err != nil {
			return nil, err
		}
	}
	return &compatData, nil
}

type MDNData struct {
	Meta          MetaData              `json:"__meta"`
	Api           map[string]any        `json:"api"`
	Browsers      map[string]any        `json:"browsers"`
	Css           CssData               `json:"css"`
	JavaScript    JavaScriptData        `json:"javascript"`
	MathMl        NamespaceElementsData `json:"mathml"`
	WebDriver     WebDriverData         `json:"webdriver"`
	WebExtensions WebExtensionsData     `json:"webextensions"`
	Html          NamespaceElementsData `json:"html"`
	Svg           NamespaceElementsData `json:"svg"`
	WebAssembly   map[string]any        `json:"webassembly"`
}

type NamespaceElementsData struct {
	Elements         map[string]any `json:"elements"`
	GlobalAttributes map[string]any `json:"global_attributes"`
	Manifest         map[string]any `json:"manifest,omitempty"`
	AttributeValues  map[string]any `json:"attribute_values,omitempty"`
}

func (n *NamespaceElementsData) ExtractNamespaceConfig() (*model.NamespaceConfig, error) {
	globalAttribute, err := n.ExtractGlobalAttr()
	nodes, err := n.ExtractElements()
	var spec = model.NamespaceConfig{
		Nodes: nodes,
		AttributesCategories: map[model.AttributeCategories][]*model.AttributeConfig{
			model.GlobalAttributes: globalAttribute,
			model.WindowActions:    {},
			model.DocumentActions:  {},
		},
	}
	return &spec, err
}

func (n *NamespaceElementsData) ExtractElements() ([]*model.NodeConfig, error) {
	var err error
	var nodes []*model.NodeConfig
	for element, value := range n.Elements {
		var elementCompat *CompatData
		var specAttr []*model.AttributeConfig
		switch q := value.(type) {
		case map[string]any:
			for attr, val := range q {
				if attr == "__compat" {
					elementCompat, err = extractCompat(val)
					if err != nil {
						return nil, err
					}
					continue
				}
				var attrCompat *CompatData
				switch w := val.(type) {
				case map[string]any:
					for attrSpec, attrValue := range w {
						if attrSpec == "__compat" {
							attrCompat, err = extractCompat(attrValue)
							if err != nil {
								return nil, err
							}
						}
					}
				}
				var docURL string
				if attrCompat.MDNURL != "" {
					docURL = attrCompat.MDNURL
				} else {
					switch s := attrCompat.SpecURL.(type) {
					case string:
						docURL = s
					case []string:
						if len(s) > 0 {
							docURL = s[0]
						}
					}
				}
				specAttr = append(specAttr, &model.AttributeConfig{
					Name:             attr,
					Boolean:          false,
					Tags:             model.RegisterTag(attrCompat.Status),
					Comment:          "",
					DocumentationURL: docURL,
					InitialValue:     "",
					SupportedValues:  map[string]model.Comment{},
				})
			}
		}
		nodes = append(nodes, &model.NodeConfig{
			Name:                       element,
			NodeType:                   model.FullTagType,
			Tags:                       model.RegisterTag(elementCompat.Status),
			Comment:                    "",
			DocumentationURL:           elementCompat.MDNURL,
			AttributesCategorySupports: map[model.AttributeCategories][]string{},
			SpecificAttributes:         specAttr,
			SupportedChildrenTags:      []string{},
		})
	}
	return nodes, nil
}

func (n *NamespaceElementsData) ExtractGlobalAttr() ([]*model.AttributeConfig, error) {
	var err error
	var globalAttr []*model.AttributeConfig
	for attr, val := range n.GlobalAttributes {
		var attrCompat *CompatData
		switch w := val.(type) {
		case map[string]any:
			for attrSpec, attrValue := range w {
				if attrSpec == "__compat" {
					attrCompat, err = extractCompat(attrValue)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		var docURL string
		if attrCompat.MDNURL != "" {
			docURL = attrCompat.MDNURL
		} else {
			switch s := attrCompat.SpecURL.(type) {
			case string:
				docURL = s
			case []string:
				if len(s) > 0 {
					docURL = s[0]
				}
			}
		}
		globalAttr = append(globalAttr, &model.AttributeConfig{
			Name:             attr,
			Boolean:          false,
			Tags:             model.RegisterTag(attrCompat.Status),
			Comment:          "",
			DocumentationURL: docURL,
			InitialValue:     "",
			SupportedValues:  map[string]model.Comment{},
		})
	}
	return globalAttr, nil
}

type MetaData struct {
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

type CssData struct {
	AtRules    map[string]any `json:"at-rules"`
	Properties map[string]any `json:"properties"`
	Selectors  map[string]any `json:"selectors"`
	Types      map[string]any `json:"types"`
}

type JavaScriptData struct {
	Functions          map[string]any `json:"functions"`
	Grammar            map[string]any `json:"grammar"`
	Operators          map[string]any `json:"operators"`
	Statements         map[string]any `json:"statements"`
	RegularExpressions map[string]any `json:"regular_expressions"`
	Builtins           map[string]any `json:"builtins"`
	Classes            map[string]any `json:"classes"`
}

type WebDriverData struct {
	Commands map[string]any `json:"commands"`
}

type WebExtensionsData struct {
	Api           map[string]any `json:"api"`
	Manifest      map[string]any `json:"manifest"`
	MatchPatterns map[string]any `json:"match_patterns"`
}

type AttributeData struct {
	Compat CompatData `json:"__compat,omitempty"`
}

type CompatData struct {
	Description string          `json:"description,omitempty"`
	MDNURL      string          `json:"mdn_url,omitempty"`
	SpecURL     any             `json:"spec_url,omitempty"` // Can be string or array of strings
	Tags        []string        `json:"tags,omitempty"`
	SourceFile  string          `json:"source_file,omitempty"`
	Support     map[string]any  `json:"support"`
	Status      map[string]bool `json:"status"`
}

type SupportData struct {
	VersionAdded          any        `json:"version_added,omitempty"`
	VersionRemoved        any        `json:"version_removed,omitempty"`
	VersionLast           any        `json:"version_last,omitempty"`
	Prefix                string     `json:"prefix,omitempty"`
	AlternativeName       string     `json:"alternative_name,omitempty"`
	Flags                 []FlagData `json:"flags,omitempty"`
	ImplURL               string     `json:"impl_url,omitempty"` // Can be string or array of strings
	PartialImplementation bool       `json:"partial_implementation,omitempty"`
	Notes                 any        `json:"notes,omitempty"` // Can be string or array of strings
}

type FlagData struct {
	Type       string `json:"type"`
	Name       string `json:"name"`
	ValueToSet string `json:"value_to_set"`
}

func GetAllHtmlTags(url string, selector func(*html.Node) bool) []*model.NodeConfig {
	page, err := fetchPage(url)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}
	var nodes = []*model.NodeConfig{}
	doc.Find("li > details").
		Each(func(index int, item *goquery.Selection) {
			item.Find("summary").Each(func(ind int, summaryText *goquery.Selection) {
				if strings.Contains("HTML elements" /*, Global attributes, Attributes, <input> types"*/, summaryText.Text()) {
					item.Find("ol > li").Each(func(ind int, line *goquery.Selection) {
						codeNode := line.Find("code")
						abbrNode := line.Find("abbr").Text()
						aNode := line.Find("a")
						nodes = append(nodes, &model.NodeConfig{
							Name:                       strings.Trim(codeNode.Text(), "<>"),
							NodeType:                   model.FullTagType,
							Tags:                       model.RegisterTag(map[string]bool{abbrNode: true}),
							Comment:                    "",
							DocumentationURL:           "https://developer.mozilla.org" + aNode.AttrOr("href", ""),
							AttributesCategorySupports: map[model.AttributeCategories][]string{},
							SpecificAttributes:         []*model.AttributeConfig{},
							SupportedChildrenTags:      []string{},
						})
					})
				}
			})
		})
	voidNode := []string{}
	page, err = fetchPage("https://developer.mozilla.org/en-US/docs/Glossary/Void_element")
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}
	voidDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}
	voidDoc.Find("article.main-page-content > div.section-content > ul > li > a > code").Each(func(i int, row *goquery.Selection) {
		voidNode = append(voidNode, strings.Trim(row.Text(), "<>"))
	})

	for _, node := range nodes {
		if slices.Contains(voidNode, node.Name) {
			node.NodeType = model.SelfClosingType
		}
		page, err = fetchPage(node.DocumentationURL)
		if err != nil {
			log.Fatalf("Failed to read body: %v", err)
		}
		nodeDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
		if err != nil {
			log.Fatalf("Failed to parse HTML: %v", err)
		}
		nodeDoc.Find("article.main-page-content > div.section-content").Each(func(i int, documentation *goquery.Selection) {
			node.Comment = model.Comment(documentation.Text())
		})
		nodeDoc.Find("article.main-page-content > section > div.section-content > p > a").Each(func(i int, attr *goquery.Selection) {
			if attr.Text() == "global attributes" || attr.Text() == "global HTML attributes" {
				node.AttributesCategorySupports[model.GlobalAttributes] = []string{}
			}
		})
	}
	return nodes
}

func GetAllGLobalAttributes() []*model.AttributeConfig {
	eventHandler := []string{}
	globalAttributes := []*model.AttributeConfig{}
	page, err := fetchPage("https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes")
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}
	globalDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}
	globalDoc.Find("article.main-page-content > div.section-content > ul > li > code").Each(func(i int, row *goquery.Selection) {
		eventHandler = append(eventHandler, row.Text())
	})

	globalDoc.Find("article.main-page-content > section > div.section-content > dl ").Each(func(i int, attr *goquery.Selection) {
		var gl = model.AttributeConfig{
			Name:             "",
			Boolean:          false,
			Tags:             []model.Tag{},
			Comment:          "",
			DocumentationURL: "",
			InitialValue:     "",
			SupportedValues:  map[string]model.Comment{},
		}
		attr.Find("dt").Each(func(i int, conf *goquery.Selection) {
			aNode := conf.Find("a")
			gl.Name = aNode.Text()
			gl.DocumentationURL = "https://developer.mozilla.org" + aNode.AttrOr("href", "")
			gl.Tags = model.RegisterTag(map[string]bool{conf.Find("abbr").Text(): true})
		})
		gl.Comment = model.Comment(attr.Find("dd > p").Text())
		globalAttributes = append(globalAttributes, &gl)
	})
	return globalAttributes
}
