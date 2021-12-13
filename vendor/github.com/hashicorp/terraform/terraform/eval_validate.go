package terraform

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/provisioners"
	"github.com/hashicorp/terraform-plugin-sdk/tfdiags"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// EvalValidateCount is an EvalNode implementation that validates
// the count of a resource.
type EvalValidateCount struct {
	Resource *configs.Resource
}

// TODO: test
func (n *EvalValidateCount) Eval(ctx EvalContext) (interface{}, error) {
	var diags tfdiags.Diagnostics
	var count int
	var err error

	val, valDiags := ctx.EvaluateExpr(n.Resource.Count, cty.Number, nil)
	diags = diags.Append(valDiags)
	if valDiags.HasErrors() {
		goto RETURN
	}
	if val.IsNull() || !val.IsKnown() {
		goto RETURN
	}

	err = gocty.FromCtyValue(val, &count)
	if err != nil {
		// The EvaluateExpr call above already guaranteed us a number value,
		// so if we end up here then we have something that is out of range
		// for an int, and the error message will include a description of
		// the valid range.
		rawVal := val.AsBigFloat()
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid count value",
			Detail:   fmt.Sprintf("The number %s is not a valid count value: %s.", rawVal, err),
			Subject:  n.Resource.Count.Range().Ptr(),
		})
	} else if count < 0 {
		rawVal := val.AsBigFloat()
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Invalid count value",
			Detail:   fmt.Sprintf("The number %s is not a valid count value: count must not be negative.", rawVal),
			Subject:  n.Resource.Count.Range().Ptr(),
		})
	}

RETURN:
	return nil, diags.NonFatalErr()
}

// EvalValidateProvisioner validates the configuration of a provisioner
// belonging to a resource. The provisioner config is expected to contain the
// merged connection configurations.
type EvalValidateProvisioner struct {
	ResourceAddr       addrs.Resource
	Provisioner        *provisioners.Interface
	Schema             **configschema.Block
	Config             *configs.Provisioner
	ResourceHasCount   bool
	ResourceHasForEach bool
}

func (n *EvalValidateProvisioner) Validate(ctx EvalContext) error {
	provisioner := *n.Provisioner
	config := *n.Config
	schema := *n.Schema

	var diags tfdiags.Diagnostics

	// Validate the provisioner's own config first
	configVal, _, configDiags := n.evaluateBlock(ctx, config.Config, schema)
	diags = diags.Append(configDiags)
	if configDiags.HasErrors() {
		return diags.Err()
	}

	if configVal == cty.NilVal {
		// Should never happen for a well-behaved EvaluateBlock implementation
		return fmt.Errorf("EvaluateBlock returned nil value")
	}

	req := provisioners.ValidateProvisionerConfigRequest{
		Config: configVal,
	}

	resp := provisioner.ValidateProvisionerConfig(req)
	diags = diags.Append(resp.Diagnostics)

	// Now validate the connection config, which contains the merged bodies
	// of the resource and provisioner connection blocks.
	connDiags := n.validateConnConfig(ctx, config.Connection, n.ResourceAddr)
	diags = diags.Append(connDiags)

	return diags.NonFatalErr()
}

func (n *EvalValidateProvisioner) validateConnConfig(ctx EvalContext, config *configs.Connection, self addrs.Referenceable) tfdiags.Diagnostics {
	// We can't comprehensively validate the connection config since its
	// final structure is decided by the communicator and we can't instantiate
	// that until we have a complete instance state. However, we *can* catch
	// configuration keys that are not valid for *any* communicator, catching
	// typos early rather than waiting until we actually try to run one of
	// the resource's provisioners.

	var diags tfdiags.Diagnostics

	if config == nil || config.Config == nil {
		// No block to validate
		return diags
	}

	// We evaluate here just by evaluating the block and returning any
	// diagnostics we get, since evaluation alone is enough to check for
	// extraneous arguments and incorrectly-typed arguments.
	_, _, configDiags := n.evaluateBlock(ctx, config.Config, connectionBlockSupersetSchema)
	diags = diags.Append(configDiags)

	return diags
}

func (n *EvalValidateProvisioner) evaluateBlock(ctx EvalContext, body hcl.Body, schema *configschema.Block) (cty.Value, hcl.Body, tfdiags.Diagnostics) {
	keyData := EvalDataForNoInstanceKey
	selfAddr := n.ResourceAddr.Instance(addrs.NoKey)

	if n.ResourceHasCount {
		// For a resource that has count, we allow count.index but don't
		// know at this stage what it will return.
		keyData = InstanceKeyEvalData{
			CountIndex: cty.UnknownVal(cty.Number),
		}

		// "self" can't point to an unknown key, but we'll force it to be
		// key 0 here, which should return an unknown value of the
		// expected type since none of these elements are known at this
		// point anyway.
		selfAddr = n.ResourceAddr.Instance(addrs.IntKey(0))
	} else if n.ResourceHasForEach {
		// For a resource that has for_each, we allow each.value and each.key
		// but don't know at this stage what it will return.
		keyData = InstanceKeyEvalData{
			EachKey:   cty.UnknownVal(cty.String),
			EachValue: cty.DynamicVal,
		}

		// "self" can't point to an unknown key, but we'll force it to be
		// key "" here, which should return an unknown value of the
		// expected type since none of these elements are known at
		// this point anyway.
		selfAddr = n.ResourceAddr.Instance(addrs.StringKey(""))
	}

	return ctx.EvaluateBlock(body, schema, selfAddr, keyData)
}

// connectionBlockSupersetSchema is a schema representing the superset of all
// possible arguments for "connection" blocks across all supported connection
// types.
//
// This currently lives here because we've not yet updated our communicator
// subsystem to be aware of schema itself. Once that is done, we can remove
// this and use a type-specific schema from the communicator to validate
// exactly what is expected for a given connection type.
var connectionBlockSupersetSchema = &configschema.Block{
	Attributes: map[string]*configschema.Attribute{
		// NOTE: "type" is not included here because it's treated special
		// by the config loader and stored away in a separate field.

		// Common attributes for both connection types
		"host": {
			Type:     cty.String,
			Required: true,
		},
		"type": {
			Type:     cty.String,
			Optional: true,
		},
		"user": {
			Type:     cty.String,
			Optional: true,
		},
		"password": {
			Type:     cty.String,
			Optional: true,
		},
		"port": {
			Type:     cty.String,
			Optional: true,
		},
		"timeout": {
			Type:     cty.String,
			Optional: true,
		},
		"script_path": {
			Type:     cty.String,
			Optional: true,
		},

		// For type=ssh only (enforced in ssh communicator)
		"private_key": {
			Type:     cty.String,
			Optional: true,
		},
		"certificate": {
			Type:     cty.String,
			Optional: true,
		},
		"host_key": {
			Type:     cty.String,
			Optional: true,
		},
		"agent": {
			Type:     cty.Bool,
			Optional: true,
		},
		"agent_identity": {
			Type:     cty.String,
			Optional: true,
		},
		"bastion_host": {
			Type:     cty.String,
			Optional: true,
		},
		"bastion_host_key": {
			Type:     cty.String,
			Optional: true,
		},
		"bastion_port": {
			Type:     cty.Number,
			Optional: true,
		},
		"bastion_user": {
			Type:     cty.String,
			Optional: true,
		},
		"bastion_password": {
			Type:     cty.String,
			Optional: true,
		},
		"bastion_private_key": {
			Type:     cty.String,
			Optional: true,
		},
		"bastion_certificate": {
			Type:     cty.String,
			Optional: true,
		},

		// For type=winrm only (enforced in winrm communicator)
		"https": {
			Type:     cty.Bool,
			Optional: true,
		},
		"insecure": {
			Type:     cty.Bool,
			Optional: true,
		},
		"cacert": {
			Type:     cty.String,
			Optional: true,
		},
		"use_ntlm": {
			Type:     cty.Bool,
			Optional: true,
		},
	},
}

// connectionBlockSupersetSchema is a schema representing the superset of all
// possible arguments for "connection" blocks across all supported connection
// types.
//
// This currently lives here because we've not yet updated our communicator
// subsystem to be aware of schema itself. It's exported only for use in the
// configs/configupgrade package and should not be used from anywhere else.
// The caller may not modify any part of the returned schema data structure.
func ConnectionBlockSupersetSchema() *configschema.Block {
	return connectionBlockSupersetSchema
}

// EvalValidateResource validates the configuration of a resource.
type EvalValidateResource struct {
	Addr           addrs.Resource
	Provider       *providers.Interface
	ProviderSchema **ProviderSchema
	Config         *configs.Resource
	ProviderMetas  map[addrs.Provider]*configs.ProviderMeta

	// IgnoreWarnings means that warnings will not be passed through. This allows
	// "just-in-time" passes of validation to continue execution through warnings.
	IgnoreWarnings bool

	// ConfigVal, if non-nil, will be updated with the value resulting from
	// evaluating the given configuration body. Since validation is performed
	// very early, this value is likely to contain lots of unknown values,
	// but its type will conform to the schema of the resource type associated
	// with the resource instance being validated.
	ConfigVal *cty.Value
}

func (n *EvalValidateResource) Validate(ctx EvalContext) error {
	if n.ProviderSchema == nil || *n.ProviderSchema == nil {
		return fmt.Errorf("EvalValidateResource has nil schema for %s", n.Addr)
	}

	var diags tfdiags.Diagnostics
	provider := *n.Provider
	cfg := *n.Config
	schema := *n.ProviderSchema
	mode := cfg.Mode

	keyData := EvalDataForNoInstanceKey

	switch {
	case n.Config.Count != nil:
		// If the config block has count, we'll evaluate with an unknown
		// number as count.index so we can still type check even though
		// we won't expand count until the plan phase.
		keyData = InstanceKeyEvalData{
			CountIndex: cty.UnknownVal(cty.Number),
		}

		// Basic type-checking of the count argument. More complete validation
		// of this will happen when we DynamicExpand during the plan walk.
		countDiags := n.validateCount(ctx, n.Config.Count)
		diags = diags.Append(countDiags)

	case n.Config.ForEach != nil:
		keyData = InstanceKeyEvalData{
			EachKey:   cty.UnknownVal(cty.String),
			EachValue: cty.UnknownVal(cty.DynamicPseudoType),
		}

		// Evaluate the for_each expression here so we can expose the diagnostics
		forEachDiags := n.validateForEach(ctx, n.Config.ForEach)
		diags = diags.Append(forEachDiags)
	}

	diags = diags.Append(validateDependsOn(ctx, n.Config.DependsOn))

	// Validate the provider_meta block for the provider this resource
	// belongs to, if there is one.
	//
	// Note: this will return an error for every resource a provider
	// uses in a module, if the provider_meta for that module is
	// incorrect. The only way to solve this that we've foudn is to
	// insert a new ProviderMeta graph node in the graph, and make all
	// that provider's resources in the module depend on the node. That's
	// an awful heavy hammer to swing for this feature, which should be
	// used only in limited cases with heavy coordination with the
	// Terraform team, so we're going to defer that solution for a future
	// enhancement to this functionality.
	/*
		if n.ProviderMetas != nil {
			if m, ok := n.ProviderMetas[n.ProviderAddr.ProviderConfig.Type]; ok && m != nil {
				// if the provider doesn't support this feature, throw an error
				if (*n.ProviderSchema).ProviderMeta == nil {
					diags = diags.Append(&hcl.Diagnostic{
						Severity: hcl.DiagError,
						Summary:  fmt.Sprintf("Provider %s doesn't support provider_meta", cfg.ProviderConfigAddr()),
						Detail:   fmt.Sprintf("The resource %s belongs to a provider that doesn't support provider_meta blocks", n.Addr),
						Subject:  &m.ProviderRange,
					})
				} else {
					_, _, metaDiags := ctx.EvaluateBlock(m.Config, (*n.ProviderSchema).ProviderMeta, nil, EvalDataForNoInstanceKey)
					diags = diags.Append(metaDiags)
				}
			}
		}
	*/
	// BUG(paddy): we're not validating provider_meta blocks on EvalValidate right now
	// because the ProviderAddr for the resource isn't available on the EvalValidate
	// struct.

	// Provider entry point varies depending on resource mode, because
	// managed resources and data resources are two distinct concepts
	// in the provider abstraction.
	switch mode {
	case addrs.ManagedResourceMode:
		schema, _ := schema.SchemaForResourceType(mode, cfg.Type)
		if schema == nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid resource type",
				Detail:   fmt.Sprintf("The provider %s does not support resource type %q.", cfg.ProviderConfigAddr(), cfg.Type),
				Subject:  &cfg.TypeRange,
			})
			return diags.Err()
		}

		configVal, _, valDiags := ctx.EvaluateBlock(cfg.Config, schema, nil, keyData)
		diags = diags.Append(valDiags)
		if valDiags.HasErrors() {
			return diags.Err()
		}

		if cfg.Managed != nil { // can be nil only in tests with poorly-configured mocks
			for _, traversal := range cfg.Managed.IgnoreChanges {
				// validate the ignore_changes traversals apply.
				moreDiags := schema.StaticValidateTraversal(traversal)
				diags = diags.Append(moreDiags)

				// TODO: we want to notify users that they can't use
				// ignore_changes for computed attributes, but we don't have an
				// easy way to correlate the config value, schema and
				// traversal together.
			}
		}

		// Use unmarked value for validate request
		unmarkedConfigVal, _ := configVal.UnmarkDeep()
		req := providers.ValidateResourceTypeConfigRequest{
			TypeName: cfg.Type,
			Config:   unmarkedConfigVal,
		}

		resp := provider.ValidateResourceTypeConfig(req)
		diags = diags.Append(resp.Diagnostics.InConfigBody(cfg.Config))

		if n.ConfigVal != nil {
			*n.ConfigVal = configVal
		}

	case addrs.DataResourceMode:
		schema, _ := schema.SchemaForResourceType(mode, cfg.Type)
		if schema == nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid data source",
				Detail:   fmt.Sprintf("The provider %s does not support data source %q.", cfg.ProviderConfigAddr(), cfg.Type),
				Subject:  &cfg.TypeRange,
			})
			return diags.Err()
		}

		configVal, _, valDiags := ctx.EvaluateBlock(cfg.Config, schema, nil, keyData)
		diags = diags.Append(valDiags)
		if valDiags.HasErrors() {
			return diags.Err()
		}

		// Use unmarked value for validate request
		unmarkedConfigVal, _ := configVal.UnmarkDeep()
		req := providers.ValidateDataSourceConfigRequest{
			TypeName: cfg.Type,
			Config:   unmarkedConfigVal,
		}

		resp := provider.ValidateDataSourceConfig(req)
		diags = diags.Append(resp.Diagnostics.InConfigBody(cfg.Config))
	}

	if n.IgnoreWarnings {
		// If we _only_ have warnings then we'll return nil.
		if diags.HasErrors() {
			return diags.NonFatalErr()
		}
		return nil
	} else {
		// We'll return an error if there are any diagnostics at all, even if
		// some of them are warnings.
		return diags.NonFatalErr()
	}
}

func (n *EvalValidateResource) validateCount(ctx EvalContext, expr hcl.Expression) (diags tfdiags.Diagnostics) {
	val, countDiags := evaluateCountExpressionValue(expr, ctx)
	// If the value isn't known then that's the best we can do for now, but
	// we'll check more thoroughly during the plan walk
	if !val.IsKnown() {
		return diags
	}

	if countDiags.HasErrors() {
		diags = diags.Append(countDiags)
	}

	return diags
}

func (n *EvalValidateResource) validateForEach(ctx EvalContext, expr hcl.Expression) (diags tfdiags.Diagnostics) {
	val, forEachDiags := evaluateForEachExpressionValue(expr, ctx)
	// If the value isn't known then that's the best we can do for now, but
	// we'll check more thoroughly during the plan walk
	if !val.IsKnown() {
		return diags
	}

	if forEachDiags.HasErrors() {
		diags = diags.Append(forEachDiags)
	}

	return diags
}

func validateDependsOn(ctx EvalContext, dependsOn []hcl.Traversal) (diags tfdiags.Diagnostics) {
	for _, traversal := range dependsOn {
		ref, refDiags := addrs.ParseRef(traversal)
		diags = diags.Append(refDiags)
		if !refDiags.HasErrors() && len(ref.Remaining) != 0 {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Invalid depends_on reference",
				Detail:   "References in depends_on must be to a whole object (resource, etc), not to an attribute of an object.",
				Subject:  ref.Remaining.SourceRange().Ptr(),
			})
		}

		// The ref must also refer to something that exists. To test that,
		// we'll just eval it and count on the fact that our evaluator will
		// detect references to non-existent objects.
		if !diags.HasErrors() {
			scope := ctx.EvaluationScope(nil, EvalDataForNoInstanceKey)
			if scope != nil { // sometimes nil in tests, due to incomplete mocks
				_, refDiags = scope.EvalReference(ref, cty.DynamicPseudoType)
				diags = diags.Append(refDiags)
			}
		}
	}
	return diags
}
