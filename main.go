package main

import (
	"fmt"

	"github.com/Nevoral/SpecGenzo/model"
	"github.com/Nevoral/SpecGenzo/updater"
	"golang.org/x/net/html"
)

func main() {
	for _, val := range updater.GetAllHtmlTags("https://developer.mozilla.org/en-US/docs/Web/HTML/Element", func(n *html.Node) bool {
		return n.Data == "a"
	}) {
		if _, exist := val.AttributesCategorySupports[model.GlobalAttributes]; exist {
			fmt.Println(len(val.AttributesCategorySupports))
		} else {
			fmt.Println(val.Name)
		}
	}
}
