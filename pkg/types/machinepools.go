package types

import (
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// MachinePool is a pool of machines to be installed.
type MachinePool struct {
	// Name is the name of the machine pool.
	// For the control plane machine pool, the name will always be "master".
	// For the compute machine pools, the only valid name is "worker".
	Name string `json:"name"`

	// Replicas is the count of machines for this machine pool.
	Replicas *int64 `json:"replicas,omitempty"`

	// Platform is configuration for machine pool specific to the platfrom.
	Platform MachinePoolPlatform `json:"platform"`
}

// MachinePoolPlatform is the platform-specific configuration for a machine
// pool. Only one of the platforms should be set.
type MachinePoolPlatform struct {
	// AWS is the configuration used when installing on AWS.
	AWS *aws.MachinePool `json:"aws,omitempty"`

	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *libvirt.MachinePool `json:"libvirt,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
	OpenStack *openstack.MachinePool `json:"openstack,omitempty"`

	// VSphere is the configuration used when installing on vSphere.
	VSphere *vsphere.MachinePool `json:"vsphere,omitempty"`
}

// Name returns a string representation of the platform (e.g. "aws" if
// AWS is non-nil).  It returns an empty string if no platform is
// configured.
func (p *MachinePoolPlatform) Name() string {
	switch {
	case p == nil:
		return ""
	case p.AWS != nil:
		return aws.Name
	case p.Libvirt != nil:
		return libvirt.Name
	case p.OpenStack != nil:
		return openstack.Name
	case p.VSphere != nil:
		return vsphere.Name
	default:
		return ""
	}
}
