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
}

type global struct {
	SecretName      string `ini:"secret-name"`
	SecretNamespace string `ini:"secret-namespace"`
	InsecureFlag    int    `ini:"insecure-flag"`
}

type workspace struct {
	Server           string `ini:"server"`
	Datacenter       string `ini:"datacenter"`
	DefaultDatastore string `ini:"default-datastore"`
	Folder           string `ini:"folder"`
}

type virtualCenter struct {
	Datacenters string `ini:"datacenters"`
}

// CloudProviderConfig generates the cloud provider config for the vSphere platform.
func CloudProviderConfig(clusterName string, p *vspheretypes.Platform) (string, error) {
	file := ini.Empty()
	config := &config{
		Global: global{
			SecretName:      "vsphere-creds",
			SecretNamespace: "kube-system",
			InsecureFlag:    1,
		},
		Workspace: workspace{
			Server:           p.VCenter,
			Datacenter:       p.Datacenter,
			DefaultDatastore: p.DefaultDatastore,
			Folder:           clusterName,
		},
	}
	if err := file.ReflectFrom(config); err != nil {
		return "", errors.Wrap(err, "failed to reflect from config")
	}
	s, err := file.NewSection(fmt.Sprintf("VirtualCenter %q", p.VCenter))
	if err != nil {
		return "", errors.Wrapf(err, "failed to create section for virtual center")
	}
	if err := s.ReflectFrom(
		&virtualCenter{
			Datacenters: p.Datacenter,
		}); err != nil {
		return "", errors.Wrapf(err, "failed to reflect from virtual center")
	}
	buf := &bytes.Buffer{}
	if _, err := file.WriteTo(buf); err != nil {
		return "", errors.Wrap(err, "failed to write out cloud provider config")
	}
	return buf.String(), nil
}
