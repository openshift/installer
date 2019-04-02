package vsphere

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	ini "gopkg.in/ini.v1"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

type config struct {
	Global    global
	Workspace workspace
	Disk      disk
	Network   network
}

type global struct {
	SecretName      string `ini:"secret-name"`
	SecretNamespace string `ini:"secret-namespace"`
}

type workspace struct {
	Server           string `ini:"server"`
	Datacenter       string `ini:"datacenter"`
	DefaultDatastore string `ini:"default-datastore"`
	ResourcePoolPath string `ini:"resourcepool-path,omitempty"`
	Folder           string `ini:"folder"`
}

type disk struct {
	SCSIControllerType string `ini:"scsicontrollertype"`
}

type network struct {
	PublicNetwork string `ini:"public-network"`
}

type virtualCenter struct {
	Datacenters []string `ini:"datacenters"`
}

// CloudProviderConfig generates the cloud provider config for the vSphere platform.
func CloudProviderConfig(p *vspheretypes.Platform) (string, error) {
	file := ini.Empty()
	config := &config{
		Global: global{
			SecretName:      "vsphere-creds",
			SecretNamespace: "kube-system",
		},
		Workspace: workspace{
			Server:           p.Workspace.Server,
			Datacenter:       p.Workspace.Datacenter,
			DefaultDatastore: p.Workspace.DefaultDatastore,
			ResourcePoolPath: p.Workspace.ResourcePoolPath,
			Folder:           p.Workspace.Folder,
		},
		Disk: disk{
			SCSIControllerType: p.SCSIControllerType,
		},
		Network: network{
			PublicNetwork: p.PublicNetwork,
		},
	}
	if err := file.ReflectFrom(config); err != nil {
		return "", errors.Wrap(err, "failed to reflect from config")
	}
	for _, vc := range p.VirtualCenters {
		s, err := file.NewSection(fmt.Sprintf("VirtualCenter %q", vc.Name))
		if err != nil {
			return "", errors.Wrapf(err, "failed to create section for virtual center %q", vc.Name)
		}
		if err := s.ReflectFrom(
			&virtualCenter{
				Datacenters: vc.Datacenters,
			}); err != nil {
			return "", errors.Wrapf(err, "failed to reflect from virtual center %q", vc.Name)
		}
	}
	buf := &bytes.Buffer{}
	if _, err := file.WriteTo(buf); err != nil {
		return "", errors.Wrap(err, "failed to write out cloud provider config")
	}
	return buf.String(), nil
}
