package updater

import (
	"bytes"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func GetAllHtmlTags(url string, selector func(*html.Node) bool) []string {
	page, err := fetchPage(url)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(page))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}
	var tags []string
	doc.Find("li > details").
		Each(func(index int, item *goquery.Selection) {
			item.Find("summary").Each(func(ind int, summaryText *goquery.Selection) {
				if strings.Contains("HTML elements, Global attributes, Attributes, <input> types", summaryText.Text()) {
					item.Find("ol > li").Each(func(ind int, line *goquery.Selection) {
						tags = append(tags, line.Text())
					})
				}
			})
		})
	return tags
}
