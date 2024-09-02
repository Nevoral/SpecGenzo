package model

import (
	"fmt"
	"strings"
)

type Comment string

// BuildComment - this method will create comment for attribute or tag.
// starts after "[TagName/AttributeName] - %s", BuildComment
func (c Comment) BuildComment() string {
	parts := strings.Split(strings.TrimSpace(string(c)), "\n")

	var result string
	for index, val := range parts {
		if index == 0 {
			result += fmt.Sprintf("%s", val)
		}
		result += fmt.Sprintf("\n%s", val)
	}
	return result
}

type WebSpecification struct {
	Version string
	Spec    map[Namespace]*NamespaceConfig
}

func (w *WebSpecification) GetConfig(tagType Namespace) *NamespaceConfig {
	return w.Spec[tagType]
}
