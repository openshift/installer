// Package packet supply utilities to extract information from terraform state
package packet

import (
	"github.com/openshift/installer/pkg/terraform"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// BootstrapIP returns the ip address for bootstrap host.
// TODO(displague) implement
func BootstrapIP(tfs *terraform.State) (string, error) {
	br, err := terraform.LookupResource(tfs, "module.bootstrap", "packet_device", "lb")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap")
	}
	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}
	bootstrap, _, err := unstructured.NestedString(br.Instances[0].Attributes, "access_public_ipv4")
	return bootstrap, err
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
// TODO(displague) implement
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	return []string{""}, nil
}
