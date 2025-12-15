package fwxschema

import (
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// NestedBlockObjectWithValidators is an optional interface on
// NestedBlockObject which enables Object validation support.
type NestedBlockObjectWithValidators interface {
	fwschema.NestedBlockObject

	// ObjectValidators should return a list of Object validators.
	ObjectValidators() []validator.Object
}
