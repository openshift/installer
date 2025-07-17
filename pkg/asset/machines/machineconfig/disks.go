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

// DiskMountUnit is used to supply the template the proper fields to produce the unit string.
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

// ForDiskSetup generates a machine config for the three disk setup types, etcd, swap or user-defined.
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
