package terraform

import (
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/configs"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform-plugin-sdk/tfdiags"
)

// ApplyGraphBuilder implements GraphBuilder and is responsible for building
// a graph for applying a Terraform diff.
//
// Because the graph is built from the diff (vs. the config or state),
// this helps ensure that the apply-time graph doesn't modify any resources
// that aren't explicitly in the diff. There are other scenarios where the
// diff can be deviated, so this is just one layer of protection.
type ApplyGraphBuilder struct {
	// Config is the configuration tree that the diff was built from.
	Config *configs.Config

	// Changes describes the changes that we need apply.
	Changes *plans.Changes

	// State is the current state
	State *states.State

	// Components is a factory for the plug-in components (providers and
	// provisioners) available for use.
	Components contextComponentFactory

	// Schemas is the repository of schemas we will draw from to analyse
	// the configuration.
	Schemas *Schemas

	// Targets are resources to target. This is only required to make sure
	// unnecessary outputs aren't included in the apply graph. The plan
	// builder successfully handles targeting resources. In the future,
	// outputs should go into the diff so that this is unnecessary.
	Targets []addrs.Targetable

	// Validate will do structural validation of the graph.
	Validate bool
}

// See GraphBuilder
func (b *ApplyGraphBuilder) Build(path addrs.ModuleInstance) (*Graph, tfdiags.Diagnostics) {
	return (&BasicGraphBuilder{
		Steps:    b.Steps(),
		Validate: b.Validate,
		Name:     "ApplyGraphBuilder",
	}).Build(path)
}

// See GraphBuilder
func (b *ApplyGraphBuilder) Steps() []GraphTransformer {
	// Custom factory for creating providers.
	concreteProvider := func(a *NodeAbstractProvider) dag.Vertex {
		return &NodeApplyableProvider{
			NodeAbstractProvider: a,
		}
	}

	concreteResource := func(a *NodeAbstractResource) dag.Vertex {
		return &nodeExpandApplyableResource{
			NodeAbstractResource: a,
		}
	}

	concreteResourceInstance := func(a *NodeAbstractResourceInstance) dag.Vertex {
		return &NodeApplyableResourceInstance{
			NodeAbstractResourceInstance: a,
		}
	}

	steps := []GraphTransformer{
		// Creates all the resources represented in the config. During apply,
		// we use this just to ensure that the whole-resource metadata is
		// updated to reflect things such as whether the count argument is
		// set in config, or which provider configuration manages each resource.
		&ConfigTransformer{
			Concrete: concreteResource,
			Config:   b.Config,
		},

		// Add dynamic values
		&RootVariableTransformer{Config: b.Config},
		&ModuleVariableTransformer{Config: b.Config},
		&LocalTransformer{Config: b.Config},
		&OutputTransformer{Config: b.Config, Changes: b.Changes},

		// Creates all the resource instances represented in the diff, along
		// with dependency edges against the whole-resource nodes added by
		// ConfigTransformer above.
		&DiffTransformer{
			Concrete: concreteResourceInstance,
			State:    b.State,
			Changes:  b.Changes,
		},

		// Attach the state
		&AttachStateTransformer{State: b.State},

		// Create orphan output nodes
		&OrphanOutputTransformer{Config: b.Config, State: b.State},

		// Attach the configuration to any resources
		&AttachResourceConfigTransformer{Config: b.Config},

		// Provisioner-related transformations
		&MissingProvisionerTransformer{Provisioners: b.Components.ResourceProvisioners()},
		&ProvisionerTransformer{},

		// add providers
		TransformProviders(b.Components.ResourceProviders(), concreteProvider, b.Config),

		// Remove modules no longer present in the config
		&RemovedModuleTransformer{Config: b.Config, State: b.State},

		// Must attach schemas before ReferenceTransformer so that we can
		// analyze the configuration to find references.
		&AttachSchemaTransformer{Schemas: b.Schemas, Config: b.Config},

		// Create expansion nodes for all of the module calls. This must
		// come after all other transformers that create nodes representing
		// objects that can belong to modules.
		&ModuleExpansionTransformer{Config: b.Config},

		// Connect references so ordering is correct
		&ReferenceTransformer{},
		&AttachDependenciesTransformer{},

		// Detect when create_before_destroy must be forced on for a particular
		// node due to dependency edges, to avoid graph cycles during apply.
		&ForcedCBDTransformer{},

		// Destruction ordering
		&DestroyEdgeTransformer{
			Config:  b.Config,
			State:   b.State,
			Schemas: b.Schemas,
		},
		&CBDEdgeTransformer{
			Config:  b.Config,
			State:   b.State,
			Schemas: b.Schemas,
		},

		// We need to remove configuration nodes that are not used at all, as
		// they may not be able to evaluate, especially during destroy.
		// These include variables, locals, and instance expanders.
		&pruneUnusedNodesTransformer{},

		// Target
		&TargetsTransformer{Targets: b.Targets},

		// Add the node to fix the state count boundaries
		&CountBoundaryTransformer{
			Config: b.Config,
		},

		// Close opened plugin connections
		&CloseProviderTransformer{},
		&CloseProvisionerTransformer{},

		// close the root module
		&CloseRootModuleTransformer{},

		// Perform the transitive reduction to make our graph a bit
		// more sane if possible (it usually is possible).
		&TransitiveReductionTransformer{},
	}

	return steps
}
