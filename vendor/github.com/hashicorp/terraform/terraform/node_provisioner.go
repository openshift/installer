package terraform

import (
	"fmt"

	"github.com/hashicorp/terraform/addrs"
)

// NodeProvisioner represents a provider that has no associated operations.
// It registers all the common interfaces across operations for providers.
type NodeProvisioner struct {
	NameValue string
	PathValue addrs.ModuleInstance
}

var (
	_ GraphNodeModuleInstance = (*NodeProvisioner)(nil)
	_ GraphNodeProvisioner    = (*NodeProvisioner)(nil)
	_ GraphNodeExecutable     = (*NodeProvisioner)(nil)
)

func (n *NodeProvisioner) Name() string {
	result := fmt.Sprintf("provisioner.%s", n.NameValue)
	if len(n.PathValue) > 0 {
		result = fmt.Sprintf("%s.%s", n.PathValue.String(), result)
	}

	return result
}

// GraphNodeModuleInstance
func (n *NodeProvisioner) Path() addrs.ModuleInstance {
	return n.PathValue
}

// GraphNodeProvisioner
func (n *NodeProvisioner) ProvisionerName() string {
	return n.NameValue
}

// GraphNodeExecutable impl.
func (n *NodeProvisioner) Execute(ctx EvalContext, op walkOperation) error {
	return ctx.InitProvisioner(n.NameValue)
}
