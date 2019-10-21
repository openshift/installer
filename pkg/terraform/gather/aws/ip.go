// Package aws contains utilities that help gather AWS specific
// information from terraform state.
package aws

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/terraform"
)

// BootstrapIP returns the ip address for bootstrap host.
func BootstrapIP(tfs *terraform.State) (string, error) {
	br, err := terraform.LookupResource(tfs, "module.bootstrap", "aws_instance", "bootstrap")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap")
	}
	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}

	for _, att := range []string{"public_ip", "private_ip"} {
		ip, _, _ := unstructured.NestedString(br.Instances[0].Attributes, att)
		if ip != "" {
			return ip, nil
		}
	}
	return "", errors.New("no usable IP found for bootstrap instance")
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	mrs, err := terraform.LookupResource(tfs, "module.masters", "aws_instance", "master")
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup masters")
	}
	var errs []error
	var masters []string
	for idx, inst := range mrs.Instances {
		master, _, err := unstructured.NestedString(inst.Attributes, "private_ip")
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "no private_ip for master.%d", idx))
		}
		masters = append(masters, master)
	}
	return masters, utilerrors.NewAggregate(errs)
}
