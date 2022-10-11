package fwserver

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// schema is a temporary function that will return a given fwschema.Schema as
// tfsdk.Schema or panic. This will be removed once tfsdk.Config, tfsdk.Plan,
// and tfsdk.State no longer use tfsdk.Schema directly or are replaced.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/366
func schema(s fwschema.Schema) tfsdk.Schema {
	switch s := s.(type) {
	case tfsdk.Schema:
		return s
	case *tfsdk.Schema:
		return *s
	default:
		panic(fmt.Sprintf("unknown fwserver fwschema.Schema type: %T", s))
	}
}
