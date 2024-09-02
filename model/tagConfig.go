package model

type TagType int

const (
	DoctypeType TagType = iota
	SelfClosingType
	CommentType
	TextContentType
	FullTagType
)

func (t TagType) String() string {
	switch t {
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

type TagConfig struct {
	Name                       string
	TagType                    TagType
	Comment                    Comment
	DocumentationURL           string
	AttributesCategorySupports map[AttributeCategories][]string
	SpecificAttributes         []*AttributeConfig
	SupportedChildrenTags      []string
}

func (t *TagConfig) IsSelfClosing() bool {
	if t.TagType == SelfClosingType {
		return true
	}
	return false
}
