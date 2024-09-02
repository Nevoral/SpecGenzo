package main

import (
	"fmt"

	"github.com/Nevoral/SpecGenzo/updater"
	"golang.org/x/net/html"
)

func main() {
	fmt.Println(updater.GetAllHtmlTags("https://developer.mozilla.org/en-US/docs/Web/HTML/Element", func(n *html.Node) bool {
		return n.Data == "a"
	}))
}
