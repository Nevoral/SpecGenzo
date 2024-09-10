package model

type AttributeCategories string

const (
	SpecificAttributes                 AttributeCategories = "Specific Attributes"
	GlobalAttributes                                       = "Global Attributes"
	AriaAttributes                                         = "Aria Attributes"
	DocumentActions                                        = "Document Actions"
	WindowActions                                          = "Window Actions"
	AnimationAdditionAttributes                            = "Animation Addition Attributes"
	AnimationEventAttributes                               = "Animation Event Attributes"
	AnimationTargetElementAttributes                       = "Animation Target Element Attributes"
	AnimationAttributeTargetAttributes                     = "Animation Attribute Target Attributes"
	AnimationTimingAttributes                              = "Animation Timing Attributes"
	AnimationValueAttributes                               = "Animation Value Attributes"
	ConditionalProcessingAttributes                        = "Conditional Processing Attributes"
	CoreAttributes                                         = "Core Attributes"
	PresentationAttributes                                 = "Presentation Attributes"
	FilterPrimitiveAttributes                              = "Filter Primitive Attributes"
	TransferFunctionAttributes                             = "Transfer Function Attributes"
	GlobalEventAttributes                                  = "Global Event Attributes"
	DocumentElementEventAttributes                         = "Document Element Event Attributes"
)

func (a AttributeCategories) String() string {
	return string(a)
}

type AttributeConfig struct {
	Name             string
	Boolean          bool
	Tags             []Tag
	Comment          Comment
	DocumentationURL string
	InitialValue     string             // if exist value will be there if not it will be empty string ""
	SupportedValues  map[string]Comment // if map has 0 len its
}

func (a *AttributeConfig) IsBoolean() bool {
	return a.Boolean
}
