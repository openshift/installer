package machineconfig

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/types"
)

// diskMount is used to supply the template the proper fields to produce the unit string.
// It contains the mount path and partition label for generating systemd mount units.
type diskMount struct {
	MountPath string
	Label     string
}

const diskMountUnit = `
[Unit]
Requires=systemd-fsck@dev-disk-by\x2dpartlabel-{{.Label}}.service
After=systemd-fsck@dev-disk-by\x2dpartlabel-{{.Label}}.service

[Mount]
Where={{.MountPath}}
What=/dev/disk/by-partlabel/{{.Label}}
Type=xfs
Options=defaults,prjquota

[Install]
RequiredBy=local-fs.target
`

const swapMountUnit = `
[Swap]
What=/dev/disk/by-partlabel/{{.Label}}

[Install]
WantedBy=swap.target
`

const gptSwap = "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F"

// ForDiskSetup generates a MachineConfig for disk setup based on the specified disk type.
// It supports three disk setup types: etcd, swap, and user-defined disks.
//
// Parameters:
//   - role: the machine role (e.g., "master", "worker") that this config applies to
//   - device: the device path (e.g., "/dev/sdb") where the disk will be configured
//   - label: the partition label to assign to the disk partition
//   - path: the mount path for the disk (not used for swap disks)
//   - diskType: the type of disk configuration (types.Etcd, types.Swap, or types.UserDefined)
//
// Returns a MachineConfig that includes the necessary ignition configuration for:
//   - Creating and partitioning the disk
//   - Formatting the filesystem (XFS for data disks, swap for swap disks)
//   - Setting up systemd mount units
//   - Configuring mount options and dependencies
//
// The function sanitizes the label by removing non-alphanumeric characters and generates
// appropriate systemd units based on the disk type.
func ForDiskSetup(role, device, label, path string, diskType types.DiskType) (*mcfgv1.MachineConfig, error) {
	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	// Remove all non-alphanumeric characters from the label
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	label = reg.ReplaceAllString(label, "")

	mountUnit := diskMount{
		MountPath: path,
		Label:     label,
	}

	var templateStringToParse string
	switch diskType {
	case types.Etcd, types.UserDefined:
		templateStringToParse = diskMountUnit
	case types.Swap:
		templateStringToParse = swapMountUnit
	}

	diskMountUnitTemplate, err := template.New("mountUnit").Parse(templateStringToParse)
	if err != nil {
		return nil, err
	}

	var dmu bytes.Buffer
	err = diskMountUnitTemplate.Execute(&dmu, mountUnit)
	if err != nil {
		return nil, err
	}

	units := dmu.String()

	var rawExt runtime.RawExtension
	switch diskType {
	case types.Etcd, types.UserDefined:
		rawExt, err = getDiskIgnition(ignConfig, device, label, path, units)
		if err != nil {
			return nil, err
		}
	case types.Swap:
		rawExt, err = getSwapIgnition(ignConfig, device, label, units)
		if err != nil {
			return nil, err
		}
	}

	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("01-disk-setup-%s-%s", label, role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: rawExt,
		},
	}, nil
}

// getDiskIgnition creates an ignition configuration for data disks (etcd and user-defined).
// It configures the disk with a single partition, formats it as XFS, and creates a systemd
// mount unit to mount the disk at the specified path.
//
// Parameters:
//   - ignConfig: the base ignition configuration to modify
//   - device: the device path for the disk
//   - label: the partition label to assign
//   - path: the mount path for the disk
//   - units: the systemd unit content for mounting
//
// The function:
//   - Creates a single partition that uses the entire disk
//   - Formats the partition as XFS with project quota support
//   - Configures a systemd mount unit with proper dependencies
//   - Sets mount options to "defaults,prjquota" for project quota support
//
// Returns the ignition configuration as a runtime.RawExtension for use in MachineConfig.
func getDiskIgnition(ignConfig igntypes.Config, device, label, path, units string) (runtime.RawExtension, error) {
	unitName := strings.Trim(path, "/")
	unitName = strings.ReplaceAll(unitName, "/", "-")

	ignConfig.Storage.Disks = append(ignConfig.Storage.Disks, igntypes.Disk{
		Device: device,
		Partitions: []igntypes.Partition{{
			Label:    ptr.To(label),
			StartMiB: ptr.To(0),
			SizeMiB:  ptr.To(0),
		}},
		WipeTable: ptr.To(true),
	})

	ignConfig.Storage.Filesystems = append(ignConfig.Storage.Filesystems, igntypes.Filesystem{
		Device:         fmt.Sprintf("/dev/disk/by-partlabel/%s", label),
		Format:         ptr.To("xfs"),
		Label:          ptr.To(label),
		MountOptions:   []igntypes.MountOption{"defaults", "prjquota"},
		Path:           ptr.To(path),
		WipeFilesystem: ptr.To(true),
	})
	ignConfig.Systemd.Units = append(ignConfig.Systemd.Units, igntypes.Unit{
		Name:     fmt.Sprintf("%s.mount", unitName),
		Enabled:  ptr.To(true),
		Contents: &units,
	})
	return ignition.ConvertToRawExtension(ignConfig)
}

// getSwapIgnition creates an ignition configuration for swap disks.
// It configures the disk with a swap partition using the appropriate GPT partition type
// and creates a systemd swap unit to activate the swap space.
//
// Parameters:
//   - ignConfig: the base ignition configuration to modify
//   - device: the device path for the disk
//   - label: the partition label to assign (typically "swap")
//   - units: the systemd unit content for swap activation
//
// The function:
//   - Creates a single partition with the GPT swap partition type (0657FD6D-A4AB-43C4-84E5-0933C84B4F4F)
//   - Formats the partition as swap space
//   - Configures a systemd swap unit for automatic activation
//   - Uses a fixed unit name "dev-disk-by\x2dpartlabel-swap.swap"
//
// Returns the ignition configuration as a runtime.RawExtension for use in MachineConfig.
func getSwapIgnition(ignConfig igntypes.Config, device, label, units string) (runtime.RawExtension, error) {
	unitName := "dev-disk-by\\x2dpartlabel-swap.swap"
	ignConfig.Storage.Disks = append(ignConfig.Storage.Disks, igntypes.Disk{
		Device: device,
		Partitions: []igntypes.Partition{{
			Label:    ptr.To(label),
			StartMiB: ptr.To(0),
			SizeMiB:  ptr.To(0),
			GUID:     ptr.To(gptSwap),
		}},
		WipeTable: ptr.To(true),
	})

	ignConfig.Storage.Filesystems = append(ignConfig.Storage.Filesystems, igntypes.Filesystem{
		Device: fmt.Sprintf("/dev/disk/by-partlabel/%s", label),
		Format: ptr.To("swap"),
		Label:  ptr.To(label),
	})
	ignConfig.Systemd.Units = append(ignConfig.Systemd.Units, igntypes.Unit{
		Name:     unitName,
		Enabled:  ptr.To(true),
		Contents: &units,
	})
	return ignition.ConvertToRawExtension(ignConfig)
}
