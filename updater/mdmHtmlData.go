package updater

import (
	"bytes"
	"log"
	"slices"
	"strings"

	"github.com/Nevoral/SpecGenzo/model"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

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
							Tags:                       []model.Tag{model.RegisterTag(abbrNode)},
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

	eventHandler := []string{}
	globalAttributes := []*model.AttributeConfig{}
	page, err = fetchPage("https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes")
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
			Comment:          "",
			DocumentationURL: url,
			InitialValue:     "",
			SupportedValues:  map[string]model.Comment{},
		}
		attr.Find("dt").Each(func(i int, conf *goquery.Selection) {
			aNode := conf.Find("a")
			gl.Name = aNode.Text()
			gl.DocumentationURL = "https://developer.mozilla.org" + aNode.AttrOr("href", "")
			gl.Tags = append(gl.Tags, model.RegisterTag(conf.Find("abbr").Text()))
		})
		gl.Comment = model.Comment(attr.Find("dd > p").Text())
		globalAttributes = append(globalAttributes, &gl)
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
