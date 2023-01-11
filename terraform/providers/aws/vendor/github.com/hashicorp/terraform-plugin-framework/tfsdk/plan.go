package tfsdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Plan represents a Terraform plan.
type Plan struct {
	Raw    tftypes.Value
	Schema Schema
}

// Get populates the struct passed as `target` with the entire plan.
func (p Plan) Get(ctx context.Context, target interface{}) diag.Diagnostics {
	return p.data().Get(ctx, target)
}

// GetAttribute retrieves the attribute found at `path` and populates the
// `target` with the value.
func (p Plan) GetAttribute(ctx context.Context, path path.Path, target interface{}) diag.Diagnostics {
	return p.data().GetAtPath(ctx, path, target)
}

// PathMatches returns all matching path.Paths from the given path.Expression.
//
// If a parent path is null or unknown, which would prevent a full expression
// from matching, the parent path is returned rather than no match to prevent
// false positives.
func (p Plan) PathMatches(ctx context.Context, pathExpr path.Expression) (path.Paths, diag.Diagnostics) {
	return p.data().PathMatches(ctx, pathExpr)
}

// Set populates the entire plan using the supplied Go value. The value `val`
// should be a struct whose values have one of the attr.Value types. Each field
// must be tagged with the corresponding schema field.
func (p *Plan) Set(ctx context.Context, val interface{}) diag.Diagnostics {
	data := p.data()
	diags := data.Set(ctx, val)

	if diags.HasError() {
		return diags
	}

	p.Raw = data.TerraformValue

	return diags
}

// SetAttribute sets the attribute at `path` using the supplied Go value.
//
// The attribute path and value must be valid with the current schema. If the
// attribute path already has a value, it will be overwritten. If the attribute
// path does not have a value, it will be added, including any parent attribute
// paths as necessary.
//
// Lists can only have the next element added according to the current length.
func (p *Plan) SetAttribute(ctx context.Context, path path.Path, val interface{}) diag.Diagnostics {
	data := p.data()
	diags := data.SetAtPath(ctx, path, val)

	if diags.HasError() {
		return diags
	}

	p.Raw = data.TerraformValue

	return diags
}

func (p Plan) data() *fwschemadata.Data {
	return &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionPlan,
		Schema:         p.Schema,
		TerraformValue: p.Raw,
	}
}
