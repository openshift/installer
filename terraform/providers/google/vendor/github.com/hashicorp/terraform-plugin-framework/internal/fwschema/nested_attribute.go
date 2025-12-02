package fwschema

// NestedAttribute defines a schema attribute that contains nested attributes.
type NestedAttribute interface {
	Attribute

	// GetNestedObject should return the object underneath the nested
	// attribute. For single nesting mode, the NestedAttributeObject can be
	// generated from the Attribute.
	GetNestedObject() NestedAttributeObject

	// GetNestingMode should return the nesting mode (list, map, set, or
	// single) of the nested attributes or left unset if this Attribute
	// does not represent nested attributes.
	GetNestingMode() NestingMode
}
