// Package openstack contains utilities that help gather Openstack specific
// information from terraform state.
package openstack

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/terraform"
)

// BootstrapIP returns the ip address for bootstrap host.
func BootstrapIP(tfs *terraform.State) (string, error) {
	br, err := terraform.LookupResource(tfs, "module.bootstrap", "openstack_compute_instance_v2", "bootstrap")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap")
	}
	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}
	bootstrap, _, err := unstructured.NestedString(br.Instances[0].Attributes, "access_ip_v4")
	if err != nil {
		return "", errors.New("no public_ip found for bootstrap")
	}
	return bootstrap, nil
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	mrs, err := terraform.LookupResource(tfs, "module.masters", "openstack_compute_instance_v2", "master_conf")
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup masters")
	}
	var errs []error
	var masters []string
	for idx, inst := range mrs.Instances {
		master, _, err := unstructured.NestedString(inst.Attributes, "access_ip_v4")
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "no access_ip_v4 for master_conf.%d", idx))
		}
		masters = append(masters, master)
	}
	return masters, utilerrors.NewAggregate(errs)
}
