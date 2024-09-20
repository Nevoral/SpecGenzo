package model

import (
	"slices"
	"strings"
)

type Tag string

const (
	Experimental Tag = "Experimental"
	Deprecated       = "Deprecated"
	NonStandard      = "NonStandard"
	Standard         = "Standard"
)

func RegisterTag(state map[string]bool) []Tag {
	var list []Tag
	for key, val := range state {
		if val {
			switch strings.TrimSpace(key) {
			case "Experimental":
				list = append(list, Experimental)
			case "experimental":
				list = append(list, Experimental)
			case "Deprecated":
				list = append(list, Deprecated)
			case "deprecated":
				list = append(list, Deprecated)
			case "NonStandard":
				list = append(list, NonStandard)
			default:
				list = append(list, Standard)
			}
		}
	}
	if slices.Contains(list, Standard) && (slices.Contains(list, Experimental) || slices.Contains(list, Deprecated)) {
		list = slices.DeleteFunc(list, func(t Tag) bool {
			return t == Standard
		})
	}
	return list
}

func (t Tag) String() string {
	switch t {
	case Experimental:
		return "Experimental"
	case Deprecated:
		return "Deprecated"
	case NonStandard:
		return "NonStandard"
	default:
		return "Standard"
	}
}

type NodeType string

const (
	DoctypeType     NodeType = "DoctypeType"
	SelfClosingType          = "SelfClosingType"
	CommentType              = "CommentType"
	TextContentType          = "TextContentType"
	FullTagType              = "FullTagType"
)

func (n NodeType) String() string {
	switch n {
	case DoctypeType:
		return "DoctypeType"
	case SelfClosingType:
		return "SelfClosingType"
	case CommentType:
		return "CommentType"
	case TextContentType:
		return "TextContentType"
	case FullTagType:
		return "FullTagType"
	default:
		return "Unknown"
	}
}

type NodeConfig struct {
	Name                       string
	NodeType                   NodeType
	Tags                       []Tag
	Comment                    Comment
	DocumentationURL           string
	AttributesCategorySupports map[AttributeCategories][]string
	SpecificAttributes         []*AttributeConfig
	SupportedChildrenTags      []string
}

func (t *NodeConfig) IsSelfClosing() bool {
	if t.NodeType == SelfClosingType {
		return true
	}
	return false
}
