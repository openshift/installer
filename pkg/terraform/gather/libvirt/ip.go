// Package libvirt contains utilities that help gather Libvirt specific
// information from terraform state.
package libvirt

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/terraform"
)

// BootstrapIP returns the ip address for bootstrap host.
func BootstrapIP(tfs *terraform.State) (string, error) {
	br, err := terraform.LookupResource(tfs, "module.bootstrap", "libvirt_domain", "bootstrap")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap")
	}
	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}
	bootstrap, err := hostnameForDomain(br.Instances[0].Attributes)
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup hostname")
	}
	return bootstrap, nil
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	mrs, err := terraform.LookupResource(tfs, "", "libvirt_domain", "master")
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup masters")
	}
	var errs []error
	var masters []string
	for idx, inst := range mrs.Instances {
		master, err := hostnameForDomain(inst.Attributes)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to lookup hostname for master.%d", idx))
		}
		masters = append(masters, master)
	}
	return masters, utilerrors.NewAggregate(errs)
}

func hostnameForDomain(attr map[string]interface{}) (string, error) {
	nics, _, err := unstructured.NestedSlice(attr, "network_interface")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup network_interface")
	}
	if len(nics) == 0 {
		return "", errors.New("no network_interface found")
	}
	hostname, _, err := unstructured.NestedString(nics[0].(map[string]interface{}), "hostname")
	if err != nil {
		return "", errors.New("no hostname found")
	}
	return hostname, nil
}
