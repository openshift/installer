package tfsdk

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Attribute must satify the fwschema.Attribute interface. It must also satisfy
// fwxschema.AttributeWithPlanModifiers and fwxschema.AttributeWithValidators
// interfaces, however we cannot check that here or it would introduce an
// import cycle.
var _ fwschema.Attribute = Attribute{}

// Attribute defines the constraints and behaviors of a single value field in a
// schema. Attributes are the fields that show up in Terraform state files and
// can be used in configuration files.
type Attribute struct {
	// Type indicates what kind of attribute this is. You'll most likely
	// want to use one of the types in the types package.
	//
	// If Type is set, Attributes cannot be.
	Type attr.Type

	// Attributes can have their own, nested attributes. This nested map of
	// attributes behaves exactly like the map of attributes on the Schema
	// type.
	//
	// If Attributes is set, Type cannot be.
	Attributes fwschema.NestedAttributes

	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this attribute is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this attribute is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription string

	// Required indicates whether the practitioner must enter a value for
	// this attribute or not. Required and Optional cannot both be true,
	// and Required and Computed cannot both be true.
	Required bool

	// Optional indicates whether the practitioner can choose not to enter
	// a value for this attribute or not. Optional and Required cannot both
	// be true.
	//
	// When defining an attribute that has Optional set to true,
	// and uses PlanModifiers to set a "default value" when none is provided,
	// Computed must also be set to true. This is necessary because default
	// values are, in effect, set by the provider (i.e. computed).
	Optional bool

	// Computed indicates whether the provider may return its own value for
	// this Attribute or not. Required and Computed cannot both be true. If
	// Required and Optional are both false, Computed must be true, and the
	// attribute will be considered "read only" for the practitioner, with
	// only the provider able to set its value.
	//
	// When defining an Optional Attribute that has a "default value"
	// plan modifier, Computed must also be set to true. Otherwise,
	// Terraform will return an error like:
	//
	//      planned value ... for a non-computed attribute
	//
	Computed bool

	// Sensitive indicates whether the value of this attribute should be
	// considered sensitive data. Setting it to true will obscure the value
	// in CLI output. Sensitive does not impact how values are stored, and
	// practitioners are encouraged to store their state as if the entire
	// file is sensitive.
	Sensitive bool

	// DeprecationMessage defines warning diagnostic details to display when
	// practitioner configurations use this Attribute. The warning diagnostic
	// summary is automatically set to "Attribute Deprecated" along with
	// configuration source file and line information.
	//
	// Set this field to a practitioner actionable message such as:
	//
	//  - "Configure other_attribute instead. This attribute will be removed
	//    in the next major version of the provider."
	//  - "Remove this attribute's configuration as it no longer is used and
	//    the attribute will be removed in the next major version of the
	//    provider."
	//
	// In Terraform 1.2.7 and later, this warning diagnostic is displayed any
	// time a practitioner attempts to configure a value for this attribute and
	// certain scenarios where this attribute is referenced.
	//
	// In Terraform 1.2.6 and earlier, this warning diagnostic is only
	// displayed when the Attribute is Required or Optional, and if the
	// practitioner configuration sets the value to a known or unknown value
	// (which may eventually be null). It has no effect when the Attribute is
	// Computed-only (read-only; not Required or Optional).
	//
	// Across any Terraform version, there are no warnings raised for
	// practitioner configuration values set directly to null, as there is no
	// way for the framework to differentiate between an unset and null
	// configuration due to how Terraform sends configuration information
	// across the protocol.
	//
	// Additional information about deprecation enhancements for read-only
	// attributes can be found in:
	//
	//  - https://github.com/hashicorp/terraform/issues/7569
	//
	DeprecationMessage string

	// Validators define value validation functionality for the attribute. All
	// elements of the slice of AttributeValidator are run, regardless of any
	// previous error diagnostics.
	//
	// Many common use case validators can be found in the
	// github.com/hashicorp/terraform-plugin-framework-validators Go module.
	//
	// If the Type field points to a custom type that implements the
	// xattr.TypeWithValidate interface, the validators defined in this field
	// are run in addition to the validation defined by the type.
	Validators []AttributeValidator

	// PlanModifiers defines a sequence of modifiers for this attribute at
	// plan time. Attribute-level plan modifications occur before any
	// resource-level plan modifications.
	//
	// Any errors will prevent further execution of this sequence
	// of modifiers and modifiers associated with any nested Attribute, but
	// will not prevent execution of PlanModifiers on any other Attribute or
	// Block in the Schema.
	//
	// Plan modification only applies to resources, not data sources or
	// providers. Setting PlanModifiers on a data source or provider attribute
	// will have no effect.
	//
	// When providing PlanModifiers, it's necessary to set Computed to true.
	PlanModifiers AttributePlanModifiers
}

// ApplyTerraform5AttributePathStep transparently calls
// ApplyTerraform5AttributePathStep on a.Type or a.Attributes, whichever is
// non-nil. It allows Attributes to be walked using tftypes.Walk and
// tftypes.Transform.
func (a Attribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	if a.Type != nil {
		return a.Type.ApplyTerraform5AttributePathStep(step)
	}
	if a.Attributes != nil {
		return a.Attributes.ApplyTerraform5AttributePathStep(step)
	}
	return nil, errors.New("Attribute has no type or nested attributes")
}

// Equal returns true if `a` and `o` should be considered Equal.
func (a Attribute) Equal(o fwschema.Attribute) bool {
	if _, ok := o.(Attribute); !ok {
		return false
	}
	if a.GetType() == nil && o.GetType() != nil {
		return false
	} else if a.GetType() != nil && o.GetType() == nil {
		return false
	} else if a.GetType() != nil && o.GetType() != nil && !a.GetType().Equal(o.GetType()) {
		return false
	}
	if a.GetAttributes() == nil && o.GetAttributes() != nil {
		return false
	} else if a.GetAttributes() != nil && o.GetAttributes() == nil {
		return false
	} else if a.GetAttributes() != nil && o.GetAttributes() != nil && !a.GetAttributes().Equal(o.GetAttributes()) {
		return false
	}
	if a.GetDescription() != o.GetDescription() {
		return false
	}
	if a.GetMarkdownDescription() != o.GetMarkdownDescription() {
		return false
	}
	if a.IsRequired() != o.IsRequired() {
		return false
	}
	if a.IsOptional() != o.IsOptional() {
		return false
	}
	if a.IsComputed() != o.IsComputed() {
		return false
	}
	if a.IsSensitive() != o.IsSensitive() {
		return false
	}
	if a.GetDeprecationMessage() != o.GetDeprecationMessage() {
		return false
	}
	return true
}

// FrameworkType returns the framework type, whether the direct type or nested
// attributes type, of the attribute.
func (a Attribute) FrameworkType() attr.Type {
	if a.Attributes != nil {
		return a.Attributes.Type()
	}

	return a.Type
}

// GetAttributes satisfies the fwschema.Attribute interface.
func (a Attribute) GetAttributes() fwschema.NestedAttributes {
	return a.Attributes
}

// GetDeprecationMessage satisfies the fwschema.Attribute interface.
func (a Attribute) GetDeprecationMessage() string {
	return a.DeprecationMessage
}

// GetDescription satisfies the fwschema.Attribute interface.
func (a Attribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription satisfies the fwschema.Attribute interface.
func (a Attribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetPlanModifiers satisfies the fwxschema.AttributeWithPlanModifiers
// interface.
func (a Attribute) GetPlanModifiers() AttributePlanModifiers {
	return a.PlanModifiers
}

// GetValidators satisfies the fwxschema.AttributeWithValidators interface.
func (a Attribute) GetValidators() []AttributeValidator {
	return a.Validators
}

// GetType satisfies the fwschema.Attribute interface.
func (a Attribute) GetType() attr.Type {
	return a.Type
}

// IsComputed satisfies the fwschema.Attribute interface.
func (a Attribute) IsComputed() bool {
	return a.Computed
}

// IsOptional satisfies the fwschema.Attribute interface.
func (a Attribute) IsOptional() bool {
	return a.Optional
}

// IsRequired satisfies the fwschema.Attribute interface.
func (a Attribute) IsRequired() bool {
	return a.Required
}

// IsSensitive satisfies the fwschema.Attribute interface.
func (a Attribute) IsSensitive() bool {
	return a.Sensitive
}
