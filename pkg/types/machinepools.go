package types

import (
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/packet"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// HyperthreadingMode is the mode of hyperthreading for a machine.
// +kubebuilder:validation:Enum="";Enabled;Disabled
type HyperthreadingMode string

const (
	// HyperthreadingEnabled indicates that hyperthreading is enabled.
	HyperthreadingEnabled HyperthreadingMode = "Enabled"
	// HyperthreadingDisabled indicates that hyperthreading is disabled.
	HyperthreadingDisabled HyperthreadingMode = "Disabled"
)

// Architecture is the instruction set architecture for the machines in a pool.
// +kubebuilder:validation:Enum="";amd64
type Architecture string

const (
	// ArchitectureAMD64 indicates AMD64 (x86_64).
	ArchitectureAMD64 = "amd64"
	// ArchitectureS390X indicates s390x (IBM System Z).
	ArchitectureS390X = "s390x"
	// ArchitecturePPC64LE indicates ppc64 little endian (Power PC)
	ArchitecturePPC64LE = "ppc64le"
)

// MachinePool is a pool of machines to be installed.
type MachinePool struct {
	// Name is the name of the machine pool.
	// For the control plane machine pool, the name will always be "master".
	// For the compute machine pools, the only valid name is "worker".
	Name string `json:"name"`

	// Replicas is the machine count for the machine pool.
	Replicas *int64 `json:"replicas,omitempty"`

	// Platform is configuration for machine pool specific to the platform.
	Platform MachinePoolPlatform `json:"platform"`

	// Hyperthreading determines the mode of hyperthreading that machines in the
	// pool will utilize.
	// Default is for hyperthreading to be enabled.
	//
	// +kubebuilder:default=Enabled
	// +optional
	Hyperthreading HyperthreadingMode `json:"hyperthreading,omitempty"`

	// Architecture is the instruction set architecture of the machine pool.
	// Defaults to amd64.
	//
	// +kubebuilder:default=amd64
	// +optional
	Architecture Architecture `json:"architecture,omitempty"`
}

// MachinePoolPlatform is the platform-specific configuration for a machine
// pool. Only one of the platforms should be set.
type MachinePoolPlatform struct {
	// AWS is the configuration used when installing on AWS.
	AWS *aws.MachinePool `json:"aws,omitempty"`

	// Azure is the configuration used when installing on Azure.
	Azure *azure.MachinePool `json:"azure,omitempty"`

	// BareMetal is the configuration used when installing on bare metal.
	BareMetal *baremetal.MachinePool `json:"baremetal,omitempty"`

	// GCP is the configuration used when installing on GCP
	GCP *gcp.MachinePool `json:"gcp,omitempty"`

	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *libvirt.MachinePool `json:"libvirt,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
	OpenStack *openstack.MachinePool `json:"openstack,omitempty"`

	// VSphere is the configuration used when installing on vSphere.
	VSphere *vsphere.MachinePool `json:"vsphere,omitempty"`

	// Ovirt is the configuration used when installing on oVirt.
	Ovirt *ovirt.MachinePool `json:"ovirt,omitempty"`

	// Packet is the configuration used when installing on Packet.
	Packet *packet.MachinePool `json:"packet,omitempty"`
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
	case p.Azure != nil:
		return azure.Name
	case p.BareMetal != nil:
		return baremetal.Name
	case p.GCP != nil:
		return gcp.Name
	case p.Libvirt != nil:
		return libvirt.Name
	case p.OpenStack != nil:
		return openstack.Name
	case p.Ovirt != nil:
		return ovirt.Name
	case p.VSphere != nil:
		return vsphere.Name
	case p.Packet != nil:
		return packet.Name
	default:
		return ""
	}
}
