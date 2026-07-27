package ibmcloud

// Boot volume storage profiles supported by IBM Cloud VPC. The first
// generation profiles scale IOPS with capacity; sdp is the second generation
// profile with independently adjustable IOPS and bandwidth.
const (
	BootVolumeProfileGeneralPurpose = "general-purpose"
	BootVolumeProfile5IOPSTier      = "5iops-tier"
	BootVolumeProfile10IOPSTier     = "10iops-tier"
	BootVolumeProfileCustom         = "custom"
	BootVolumeProfileSDP            = "sdp"

	// DefaultBootVolumeProfile is the profile used when a boot volume does not
	// specify one, matching the IBM Cloud default.
	DefaultBootVolumeProfile = BootVolumeProfileGeneralPurpose
)

// BootVolumeProfiles is the set of supported boot volume storage profiles.
var BootVolumeProfiles = []string{
	BootVolumeProfileGeneralPurpose,
	BootVolumeProfile5IOPSTier,
	BootVolumeProfile10IOPSTier,
	BootVolumeProfileCustom,
	BootVolumeProfileSDP,
}

// MachinePool stores the configuration for a machine pool installed on IBM Cloud.
type MachinePool struct {
	// InstanceType is the VSI machine profile.
	InstanceType string `json:"type,omitempty"`

	// Zones is the list of availability zones used for machines in the pool.
	// +optional
	Zones []string `json:"zones,omitempty"`

	// BootVolume is the configuration for the machine's boot volume.
	// +optional
	BootVolume *BootVolume `json:"bootVolume,omitempty"`

	// DedicatedHosts is the configuration for the machine's dedicated host and profile.
	// +optional
	DedicatedHosts []DedicatedHost `json:"dedicatedHosts,omitempty"`
}

// BootVolume stores the configuration for an individual machine's boot volume.
type BootVolume struct {
	// EncryptionKey is the CRN referencing a Key Protect or Hyper Protect
	// Crypto Services key to use for volume encryption. If not specified, a
	// provider managed encryption key will be used.
	// +optional
	EncryptionKey string `json:"encryptionKey,omitempty"`

	// Profile is the block storage profile for the boot volume. First
	// generation profiles are `general-purpose`, `5iops-tier`, `10iops-tier`
	// and `custom`. The second generation profile is `sdp`, which supports
	// independently adjustable IOPS and bandwidth as well as capacities
	// beyond the first generation limit.
	// If not specified, `general-purpose` is used.
	// +optional
	Profile string `json:"profile,omitempty"`

	// SizeGiB is the size of the boot volume in GiB. The value must be at
	// least as large as the image's minimum provisioned size. First
	// generation profiles support up to 250 GiB, the `sdp` profile supports
	// up to 32000 GiB.
	// If not specified, the IBM Cloud default boot volume size is used.
	// +optional
	SizeGiB int64 `json:"sizeGiB,omitempty"`

	// IOPS is the maximum I/O operations per second for the boot volume.
	// Only configurable for the `custom` and `sdp` profiles; for the other
	// profiles the IOPS is determined by the profile itself.
	// +optional
	IOPS int64 `json:"iops,omitempty"`

	// Bandwidth is the maximum bandwidth in megabits per second for the boot
	// volume. Only configurable for the `sdp` profile. If not specified, it
	// is computed by IBM Cloud from the volume's IOPS and size.
	// +optional
	Bandwidth int64 `json:"bandwidth,omitempty"`
}

// DedicatedHost stores the configuration for the machine's dedicated host platform.
type DedicatedHost struct {
	// Name is the name of the dedicated host to provision the machine on. If
	// specified, machines will be created on pre-existing dedicated host.
	// +optional
	Name string `json:"name,omitempty"`

	// Profile is the profile ID for the dedicated host. If specified, new
	// dedicated host will be created for machines.
	// +optional
	Profile string `json:"profile,omitempty"`
}

// Set sets the values from `required` to `a`.
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}

	if required.InstanceType != "" {
		a.InstanceType = required.InstanceType
	}

	if len(required.Zones) > 0 {
		a.Zones = required.Zones
	}

	if required.BootVolume != nil {
		if a.BootVolume == nil {
			a.BootVolume = &BootVolume{}
		}
		if required.BootVolume.EncryptionKey != "" {
			a.BootVolume.EncryptionKey = required.BootVolume.EncryptionKey
		}
		if required.BootVolume.Profile != "" {
			a.BootVolume.Profile = required.BootVolume.Profile
		}
		if required.BootVolume.SizeGiB != 0 {
			a.BootVolume.SizeGiB = required.BootVolume.SizeGiB
		}
		if required.BootVolume.IOPS != 0 {
			a.BootVolume.IOPS = required.BootVolume.IOPS
		}
		if required.BootVolume.Bandwidth != 0 {
			a.BootVolume.Bandwidth = required.BootVolume.Bandwidth
		}
	}

	if len(required.DedicatedHosts) > 0 {
		a.DedicatedHosts = required.DedicatedHosts
	}
}
