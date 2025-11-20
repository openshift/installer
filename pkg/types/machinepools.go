package types

import (
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/openshift/installer/pkg/types/nutanix"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/powervc"
	"github.com/openshift/installer/pkg/types/powervs"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const (
	// MachinePoolComputeRoleName name associated with the compute machinepool.
	MachinePoolComputeRoleName = "worker"
	// MachinePoolEdgeRoleName name associated with the compute edge machinepool.
	MachinePoolEdgeRoleName = "edge"
	// MachinePoolControlPlaneRoleName name associated with the control plane machinepool.
	MachinePoolControlPlaneRoleName = "master"
	// MachinePoolArbiterRoleName name associated with the control plane machinepool for smaller sized limited nodes.
	MachinePoolArbiterRoleName = "arbiter"
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
	// ArchitectureARM64 indicates arm (aarch64) systems
	ArchitectureARM64 = "arm64"
)

// DiskType is the string representation of the three types disk setups
// +kubebuilder:validation:Enum=etcd;swap;user-defined
type DiskType string

const (
	// Etcd indicates etcd disk setup.
	Etcd DiskType = "etcd"
	// Swap indicates swap disk setup.
	Swap DiskType = "swap"
	// UserDefined indicates user-defined disk setup.
	UserDefined DiskType = "user-defined"
)

// Disk defines the type of disk (etcd, swap or user-defined) and the configuration
// of each disk type.
type Disk struct {
	Type DiskType `json:"type,omitempty"`

	UserDefined *DiskUserDefined `json:"userDefined,omitempty"`
	Etcd        *DiskEtcd        `json:"etcd,omitempty"`
	Swap        *DiskSwap        `json:"swap,omitempty"`
}

// DiskUserDefined defines a disk type of user-defined.
type DiskUserDefined struct {
	PlatformDiskID string `json:"platformDiskID,omitempty"`
	MountPath      string `json:"mountPath,omitempty"`
}

// DiskSwap defines a disk type of swap.
type DiskSwap struct {
	PlatformDiskID string `json:"platformDiskID,omitempty"`
}

// DiskEtcd defines a disk type of etcd.
type DiskEtcd struct {
	PlatformDiskID string `json:"platformDiskID,omitempty"`
}

// MachinePool is a pool of machines to be installed.
type MachinePool struct {
	// Name is the name of the machine pool.
	// For the control plane machine pool, the name will always be "master".
	// For the compute machine pools, the only valid name is "worker".
	// For the arbiter machine pools, the only valid name is "arbiter".
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

	// Fencing stores the information about a baremetal host's management controller.
	// Fencing may only be set for control plane nodes.
	// +optional
	Fencing *Fencing `json:"fencing,omitempty"`

	// DiskSetup stores the type of disks that will be setup with MachineConfigs.
	// The available types are etcd, swap and user-defined.
	// +optional
	DiskSetup []Disk `json:"diskSetup,omitempty"`
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

	// IBMCloud is the configuration used when installing on IBM Cloud.
	IBMCloud *ibmcloud.MachinePool `json:"ibmcloud,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
	OpenStack *openstack.MachinePool `json:"openstack,omitempty"`

	// VSphere is the configuration used when installing on vSphere.
	VSphere *vsphere.MachinePool `json:"vsphere,omitempty"`

	// Ovirt is the configuration used when installing on oVirt.
	Ovirt *ovirt.MachinePool `json:"ovirt,omitempty"`

	// PowerVC is the configuration used when installing on IBM Power VC.
	PowerVC *powervc.MachinePool `json:"powervc,omitempty"`

	// PowerVS is the configuration used when installing on IBM Power VS.
	PowerVS *powervs.MachinePool `json:"powervs,omitempty"`

	// Nutanix is the configuration used when installing on Nutanix.
	Nutanix *nutanix.MachinePool `json:"nutanix,omitempty"`
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
	case p.IBMCloud != nil:
		return ibmcloud.Name
	case p.OpenStack != nil:
		return openstack.Name
	case p.VSphere != nil:
		return vsphere.Name
	case p.Ovirt != nil:
		return ovirt.Name
	case p.PowerVS != nil:
		return powervs.Name
	case p.Nutanix != nil:
		return nutanix.Name
	default:
		return ""
	}
}

// Fencing stores the information about a baremetal host's management controller.
type Fencing struct {
	// Credentials stores the information about a baremetal host's management controller.
	// +optional
	Credentials []*Credential `json:"credentials,omitempty"`
}

// CertificateVerificationPolicy represents the options for CertificateVerification .
type CertificateVerificationPolicy string

const (
	// CertificateVerificationEnabled enables ssl certificate verification.
	CertificateVerificationEnabled CertificateVerificationPolicy = "Enabled"
	// CertificateVerificationDisabled disables ssl certificate verification.
	CertificateVerificationDisabled CertificateVerificationPolicy = "Disabled"
)

// Credential stores the information about a baremetal host's management controller.
type Credential struct {
	HostName string `json:"hostName,omitempty" validate:"required,uniqueField"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Address  string `json:"address" validate:"required,uniqueField"`
	// CertificateVerification Defines whether ssl certificate verification is required or not.
	// If omitted, the platform chooses a default, that default is enabled.
	// +kubebuilder:default:="Enabled"
	// +kubebuilder:validation:Enum=Enabled;Disabled
	// +optional
	CertificateVerification CertificateVerificationPolicy `json:"certificateVerification,omitempty"`
}
