package model

import "strings"

type Tag int

const (
	Experimental Tag = iota
	Deprecated
	NonStandard
	Standard
)

func RegisterTag(tag string) Tag {
	switch strings.TrimSpace(tag) {
	case "Experimental":
		return Experimental
	case "Deprecated":
		return Deprecated
	case "NonStandard":
		return NonStandard
	default:
		return Standard
	}
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

type NodeType int

const (
	DoctypeType NodeType = iota
	SelfClosingType
	CommentType
	TextContentType
	FullTagType
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
