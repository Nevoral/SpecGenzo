package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"

	"github.com/Nevoral/SpecGenzo/data/mdn"
	"github.com/Nevoral/SpecGenzo/model"
	"github.com/Nevoral/SpecGenzo/parser"
)

// downloadFile downloads a URL and writes the content to a file.
func downloadFile(url string, filePath string) error {
	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create a file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}

func main() {
	m := mdn.ScrapeMDNSource()
	spec, err := m.ExtractWebSpecification()
	if err != nil {
		panic(err)
	}
	if err = parser.CreateJsonSpec("./spec/mdnCompatDataSource.json", spec); err != nil {
		panic(err)
	}

	htmlElements := mdn.GetAllHtmlTags("https://developer.mozilla.org/en-US/docs/Web/HTML/Element", nil)
	compareNodeConfigs(spec.Spec[model.HTML].Nodes, htmlElements)
}

func compareNodeConfigs(slice1, slice2 []*model.NodeConfig) {
	map1 := make(map[string]*model.NodeConfig)
	map2 := make(map[string]*model.NodeConfig)

	for _, config := range slice1 {
		map1[config.Name] = config
	}
	for _, config := range slice2 {
		map2[config.Name] = config
	}

	for name, value := range map1 {
		if _, exists := map2[name]; !exists {
			fmt.Printf("%s - ?\n", name)
			continue
		}
		if checkStatusTag(value.Tags, map2[name].Tags) {
			fmt.Println(name)
		}
	}

	for name, value := range map2 {
		if _, exists := map1[name]; !exists {
			fmt.Printf("? - %s\n", name)
			continue
		}
		if checkStatusTag(value.Tags, map1[name].Tags) {
			fmt.Println(name)
		}
	}
}

func checkStatusTag(slice1, slice2 []model.Tag) bool {
	var r bool = false
	for _, value := range slice1 {
		if !slices.Contains(slice2, value) {
			fmt.Printf("\t%s - ?\n", value.String())
			r = true
		}
	}
	for _, value := range slice2 {
		if !slices.Contains(slice1, value) {
			fmt.Printf("\t? - %s\n", value.String())
			r = true
		}
	}
	return r
}
