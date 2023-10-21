package timeouts

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/internal/validators"
)

const (
	attributeNameCreate = "create"
	attributeNameRead   = "read"
	attributeNameUpdate = "update"
	attributeNameDelete = "delete"
)

// Opts is used as an argument to Block and Attributes to indicate which attributes
// should be created.
type Opts struct {
	Create bool
	Read   bool
	Update bool
	Delete bool
}

// Block returns a schema.Block containing attributes for each of the fields
// in Opts which are set to true. Each attribute is defined as types.StringType
// and optional. A validator is used to verify that the value assigned to an
// attribute can be parsed as time.Duration.
func Block(ctx context.Context, opts Opts) schema.Block {
	return schema.SingleNestedBlock{
		Attributes: attributesMap(opts),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(opts),
			},
		},
	}
}

// BlockAll returns a schema.Block containing attributes for each of create, read,
// update and delete. Each attribute is defined as types.StringType and optional.
// A validator is used to verify that the value assigned to an attribute can be
// parsed as time.Duration.
func BlockAll(ctx context.Context) schema.Block {
	return Block(ctx, Opts{
		Create: true,
		Read:   true,
		Update: true,
		Delete: true,
	})
}

// Attributes returns a schema.SingleNestedAttribute which contains attributes for
// each of the fields in Opts which are set to true. Each attribute is defined as
// types.StringType and optional. A validator is used to verify that the value
// assigned to an attribute can be parsed as time.Duration.
func Attributes(ctx context.Context, opts Opts) schema.Attribute {
	return schema.SingleNestedAttribute{
		Attributes: attributesMap(opts),
		CustomType: Type{
			ObjectType: types.ObjectType{
				AttrTypes: attrTypesMap(opts),
			},
		},
		Optional: true,
	}
}

// AttributesAll returns a schema.SingleNestedAttribute which contains attributes
// for each of create, read, update and delete. Each attribute is defined as
// types.StringType and optional. A validator is used to verify that the value
// assigned to an attribute can be parsed as time.Duration.
func AttributesAll(ctx context.Context) schema.Attribute {
	return Attributes(ctx, Opts{
		Create: true,
		Read:   true,
		Update: true,
		Delete: true,
	})
}

func attributesMap(opts Opts) map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{}
	attribute := schema.StringAttribute{
		Optional: true,
		Validators: []validator.String{
			validators.TimeDuration(),
		},
	}

	if opts.Create {
		attributes[attributeNameCreate] = attribute
	}

	if opts.Read {
		attributes[attributeNameRead] = attribute
	}

	if opts.Update {
		attributes[attributeNameUpdate] = attribute
	}

	if opts.Delete {
		attributes[attributeNameDelete] = attribute
	}

	return attributes
}

func attrTypesMap(opts Opts) map[string]attr.Type {
	attrTypes := map[string]attr.Type{}

	if opts.Create {
		attrTypes[attributeNameCreate] = types.StringType
	}

	if opts.Read {
		attrTypes[attributeNameRead] = types.StringType
	}

	if opts.Update {
		attrTypes[attributeNameUpdate] = types.StringType
	}

	if opts.Delete {
		attrTypes[attributeNameDelete] = types.StringType
	}

	return attrTypes
}
