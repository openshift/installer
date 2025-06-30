package machineconfig

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset/ignition"
)

type DiskMountUnit struct {
	MountPath string
	Label     string
	UnitName  string
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

func ForDiskSetup(role, device, label, path string) (*mcfgv1.MachineConfig, error) {
	ignConfig := igntypes.Config{
		Ignition: igntypes.Ignition{
			Version: igntypes.MaxVersion.String(),
		},
	}

	unitName := strings.TrimPrefix(path, "/")
	unitName = strings.Replace(unitName, "/", "-", -1)

	mountUnit := DiskMountUnit{
		MountPath: path,
		Label:     label,
		UnitName:  unitName,
	}

	diskMountUnitTemplate, err := template.New("mountUnit").Parse(diskMountUnit)
	if err != nil {
		return nil, err
	}

	var dmu bytes.Buffer
	err = diskMountUnitTemplate.Execute(&dmu, mountUnit)
	if err != nil {
		return nil, err
	}

	units := dmu.String()

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
	rawExt, err := ignition.ConvertToRawExtension(ignConfig)
	if err != nil {
		return nil, err
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
