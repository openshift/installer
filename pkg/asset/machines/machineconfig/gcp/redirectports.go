package gcp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/pkg/errors"

	"github.com/coreos/ignition/config/util"
	igntypes "github.com/coreos/ignition/config/v2_2/types"
	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/asset/ignition"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	scriptTemplate = "bootstrap/gcp/files/usr/local/bin/redirect-ports.sh.template"
	scriptFile     = "/usr/local/bin/redirect-ports.sh"
	unitFile       = "bootstrap/gcp/systemd/units/redirect-ports.service"
)

// ForPortRedirection creates the MachineConfig that adds a port redirection systemd unit
// to machines on GCP. For masters, it redirects incoming traffic. For all nodes it redirects
// outgoing traffic to the internal control plane load balancer.
func ForPortRedirection(role string, clusterDomain string) (*mcfgv1.MachineConfig, error) {
	templateData := map[string]interface{}{
		"ClusterDomain": clusterDomain,
		"ControlPlane":  (role == "master"),
	}
	script, err := readTemplate(scriptTemplate, templateData)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read %s", scriptFile)
	}
	unit, err := readFile(unitFile)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read %s", unitFile)
	}
	return &mcfgv1.MachineConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machineconfiguration.openshift.io/v1",
			Kind:       "MachineConfig",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("99-%s-redirect-lb-ports", role),
			Labels: map[string]string{
				"machineconfiguration.openshift.io/role": role,
			},
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: igntypes.Config{
				Ignition: igntypes.Ignition{
					Version: igntypes.MaxVersion.String(),
				},
				Storage: igntypes.Storage{
					Files: []igntypes.File{
						ignition.FileFromBytes(scriptFile, "root", 0555, script),
					},
				},
				Systemd: igntypes.Systemd{
					Units: []igntypes.Unit{
						{
							Name:     "redirect-ports.service",
							Enabled:  util.BoolToPtr(true),
							Contents: string(unit),
						},
					},
				},
			},
		},
	}, nil
}

func readFile(name string) ([]byte, error) {
	file, err := data.Assets.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func readTemplate(name string, data map[string]interface{}) ([]byte, error) {
	content, err := readFile(name)
	if err != nil {
		return nil, err
	}
	t, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
