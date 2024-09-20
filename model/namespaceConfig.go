package model

import (
	"fmt"
	"slices"
	"strings"
)

type Namespace string

const (
	HTML  Namespace = "HTML"
	SVG             = "SVG"
	MATH            = "MATH"
	XHTML           = "XHTML"
)

func (n Namespace) String() string {
	switch n {
	case HTML:
		return "HTML"
	case SVG:
		return "SVG"
	case MATH:
		return "MATH"
	case XHTML:
		return "XHTML"
	default:
		return "Unknown"
	}
}

type NamespaceConfig struct {
	Nodes                []*NodeConfig
	AttributesCategories map[AttributeCategories][]*AttributeConfig
}

func (n *NamespaceConfig) GetAttributeDefaultValue(name string, category AttributeCategories) string {
	atrIndex := slices.IndexFunc(n.AttributesCategories[category], func(e *AttributeConfig) bool {
		return e.Name == name
	})
	return n.AttributesCategories[category][atrIndex].InitialValue
}

func (n *NamespaceConfig) GetAttributeBoolean(name string, category AttributeCategories) bool {
	atrIndex := slices.IndexFunc(n.AttributesCategories[category], func(e *AttributeConfig) bool {
		return e.Name == name
	})
	return n.AttributesCategories[category][atrIndex].Boolean
}

func (n *NamespaceConfig) GetTagConfig(name string) (*NodeConfig, error) {
	tagIndex := slices.IndexFunc(n.Nodes, func(e *NodeConfig) bool {
		return strings.ToLower(e.Name) == name
	})
	if tagIndex < 0 {
		msg := fmt.Errorf("Error: in specification isn't any tag called %s.", name)
		return nil, msg
	}
	return n.Nodes[tagIndex], nil
}

func (n *NamespaceConfig) IsTagSelfClosing(name string) bool {
	for _, tag := range n.Nodes {
		if tag.Name == name && tag.IsSelfClosing() {
			return true
		}
	}
	return false
}

func (n *NamespaceConfig) SortAllSlicesAscending() {
	slices.SortFunc(n.Nodes, func(a, b *NodeConfig) int {
		if a.Name > b.Name {
			return 1
		} else if a.Name == b.Name {
			return 0
		}
		return -1
	})
	for _, value := range n.Nodes {
		slices.SortFunc(value.SpecificAttributes, func(a, b *AttributeConfig) int {
			if a.Name > b.Name {
				return 1
			} else if a.Name == b.Name {
				return 0
			}
			return -1
		})
	}
	for _, value := range n.AttributesCategories {
		slices.SortFunc(value, func(a, b *AttributeConfig) int {
			if a.Name > b.Name {
				return 1
			} else if a.Name == b.Name {
				return 0
			}
			return -1
		})
	}
}

func (n *NamespaceConfig) CheckValueValidity(name, value string, category AttributeCategories) bool {
	atrIndex := slices.IndexFunc(n.AttributesCategories[category], func(e *AttributeConfig) bool {
		return e.Name == name
	})
	if len(n.AttributesCategories[category][atrIndex].SupportedValues) == 0 {
		return true
	}
	if _, ok := n.AttributesCategories[category][atrIndex].SupportedValues[value]; !ok {
		return false
	}
	return true
}
