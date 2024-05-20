package cfschema

type Tagging struct {
	Taggable                 *bool                `json:"taggable,omitempty"`
	TagOnCreate              *bool                `json:"tagOnCreate,omitempty"`
	TagUpdatable             *bool                `json:"tagUpdatable,omitempty"`
	CloudFormationSystemTags *bool                `json:"cloudFormationSystemTags,omitempty"`
	TagProperty              *PropertyJsonPointer `json:"tagProperty,omitempty"`
}
