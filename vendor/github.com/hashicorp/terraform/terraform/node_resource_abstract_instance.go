package terraform

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform-plugin-sdk/tfdiags"
)

// NodeAbstractResourceInstance represents a resource instance with no
// associated operations. It embeds NodeAbstractResource but additionally
// contains an instance key, used to identify one of potentially many
// instances that were created from a resource in configuration, e.g. using
// the "count" or "for_each" arguments.
type NodeAbstractResourceInstance struct {
	NodeAbstractResource
	Addr addrs.AbsResourceInstance

	// These are set via the AttachState method.
	instanceState *states.ResourceInstance
	// storedProviderConfig is the provider address retrieved from the
	// state, but since it is only stored in the whole Resource rather than the
	// ResourceInstance, we extract it out here.
	storedProviderConfig addrs.AbsProviderConfig

	Dependencies []addrs.ConfigResource
}

// NewNodeAbstractResourceInstance creates an abstract resource instance graph
// node for the given absolute resource instance address.
func NewNodeAbstractResourceInstance(addr addrs.AbsResourceInstance) *NodeAbstractResourceInstance {
	// Due to the fact that we embed NodeAbstractResource, the given address
	// actually ends up split between the resource address in the embedded
	// object and the InstanceKey field in our own struct. The
	// ResourceInstanceAddr method will stick these back together again on
	// request.
	r := NewNodeAbstractResource(addr.ContainingResource().Config())
	return &NodeAbstractResourceInstance{
		NodeAbstractResource: *r,
		Addr:                 addr,
	}
}

func (n *NodeAbstractResourceInstance) Name() string {
	return n.ResourceInstanceAddr().String()
}

func (n *NodeAbstractResourceInstance) Path() addrs.ModuleInstance {
	return n.Addr.Module
}

// GraphNodeReferenceable
func (n *NodeAbstractResourceInstance) ReferenceableAddrs() []addrs.Referenceable {
	addr := n.ResourceInstanceAddr()
	return []addrs.Referenceable{
		addr.Resource,

		// A resource instance can also be referenced by the address of its
		// containing resource, so that e.g. a reference to aws_instance.foo
		// would match both aws_instance.foo[0] and aws_instance.foo[1].
		addr.ContainingResource().Resource,
	}
}

// GraphNodeReferencer
func (n *NodeAbstractResourceInstance) References() []*addrs.Reference {
	// If we have a configuration attached then we'll delegate to our
	// embedded abstract resource, which knows how to extract dependencies
	// from configuration. If there is no config, then the dependencies will
	// be connected during destroy from those stored in the state.
	if n.Config != nil {
		if n.Schema == nil {
			// We'll produce a log message about this out here so that
			// we can include the full instance address, since the equivalent
			// message in NodeAbstractResource.References cannot see it.
			log.Printf("[WARN] no schema is attached to %s, so config references cannot be detected", n.Name())
			return nil
		}
		return n.NodeAbstractResource.References()
	}

	// If we have neither config nor state then we have no references.
	return nil
}

// StateDependencies returns the dependencies saved in the state.
func (n *NodeAbstractResourceInstance) StateDependencies() []addrs.ConfigResource {
	if s := n.instanceState; s != nil {
		if s.Current != nil {
			return s.Current.Dependencies
		}
	}

	return nil
}

// GraphNodeProviderConsumer
func (n *NodeAbstractResourceInstance) ProvidedBy() (addrs.ProviderConfig, bool) {
	// If we have a config we prefer that above all else
	if n.Config != nil {
		relAddr := n.Config.ProviderConfigAddr()
		return addrs.LocalProviderConfig{
			LocalName: relAddr.LocalName,
			Alias:     relAddr.Alias,
		}, false
	}

	// See if we have a valid provider config from the state.
	if n.storedProviderConfig.Provider.Type != "" {
		// An address from the state must match exactly, since we must ensure
		// we refresh/destroy a resource with the same provider configuration
		// that created it.
		return n.storedProviderConfig, true
	}

	// No provider configuration found; return a default address
	return addrs.AbsProviderConfig{
		Provider: n.Provider(),
		Module:   n.ModulePath(),
	}, false
}

// GraphNodeProviderConsumer
func (n *NodeAbstractResourceInstance) Provider() addrs.Provider {
	if n.Config != nil {
		return n.Config.Provider
	}
	return addrs.ImpliedProviderForUnqualifiedType(n.Addr.Resource.ContainingResource().ImpliedProvider())
}

// GraphNodeResourceInstance
func (n *NodeAbstractResourceInstance) ResourceInstanceAddr() addrs.AbsResourceInstance {
	return n.Addr
}

// GraphNodeAttachResourceState
func (n *NodeAbstractResourceInstance) AttachResourceState(s *states.Resource) {
	if s == nil {
		log.Printf("[WARN] attaching nil state to %s", n.Addr)
		return
	}
	n.instanceState = s.Instance(n.Addr.Resource.Key)
	n.storedProviderConfig = s.ProviderConfig
}

// readDiff returns the planned change for a particular resource instance
// object.
func (n *NodeAbstractResourceInstance) readDiff(ctx EvalContext, providerSchema *ProviderSchema) (*plans.ResourceInstanceChange, error) {
	changes := ctx.Changes()
	addr := n.ResourceInstanceAddr()

	schema, _ := providerSchema.SchemaForResourceAddr(addr.Resource.Resource)
	if schema == nil {
		// Should be caught during validation, so we don't bother with a pretty error here
		return nil, fmt.Errorf("provider does not support resource type %q", addr.Resource.Resource.Type)
	}

	gen := states.CurrentGen
	csrc := changes.GetResourceInstanceChange(addr, gen)
	if csrc == nil {
		log.Printf("[TRACE] EvalReadDiff: No planned change recorded for %s", n.Addr)
		return nil, nil
	}

	change, err := csrc.Decode(schema.ImpliedType())
	if err != nil {
		return nil, fmt.Errorf("failed to decode planned changes for %s: %s", n.Addr, err)
	}

	log.Printf("[TRACE] EvalReadDiff: Read %s change from plan for %s", change.Action, n.Addr)

	return change, nil
}

func (n *NodeAbstractResourceInstance) checkPreventDestroy(change *plans.ResourceInstanceChange) error {
	if change == nil || n.Config == nil || n.Config.Managed == nil {
		return nil
	}

	preventDestroy := n.Config.Managed.PreventDestroy

	if (change.Action == plans.Delete || change.Action.IsReplace()) && preventDestroy {
		var diags tfdiags.Diagnostics
		diags = diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Instance cannot be destroyed",
			Detail: fmt.Sprintf(
				"Resource %s has lifecycle.prevent_destroy set, but the plan calls for this resource to be destroyed. To avoid this error and continue with the plan, either disable lifecycle.prevent_destroy or reduce the scope of the plan using the -target flag.",
				n.Addr.String(),
			),
			Subject: &n.Config.DeclRange,
		})
		return diags.Err()
	}

	return nil
}
