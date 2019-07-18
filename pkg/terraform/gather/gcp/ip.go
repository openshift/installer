package gcp

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/terraform"
)

// BootstrapIP returns the ip address for bootstrap host.
func BootstrapIP(tfs *terraform.State) (string, error) {
	br, err := terraform.LookupResource(tfs, "module.bootstrap", "google_compute_instance", "bootstrap")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap")
	}
	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}

	networkInterfaces, found, err := unstructured.NestedSlice(br.Instances[0].Attributes, "network_interface")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup network interface")
	}
	if !found {
		return "", errors.New("bootstrap does not contain network_interface")
	}

	accessConfigs, found, err := unstructured.NestedSlice(networkInterfaces[0].(map[string]interface{}), "access_config")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup access config")
	}
	if !found {
		return "", errors.New("bootstrap's network interface does not contain access_config")
	}

	bootstrap, found, err := unstructured.NestedString(accessConfigs[0].(map[string]interface{}), "nat_ip")
	if err != nil {
		return "", errors.New("failed to lookup public ip address")
	}
	if !found {
		return "", errors.New("access config does not contain nat_ip")
	}
	return bootstrap, nil
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	mrs, err := terraform.LookupResource(tfs, "module.master", "google_compute_instance", "master")
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup masters")
	}
	var errs []error
	var masters []string
	for idx, inst := range mrs.Instances {
		networkInterfaces, found, err := unstructured.NestedSlice(inst.Attributes, "network_interface")
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "failed to lookup network interface for master.%d", idx))
		}
		if !found {
			errs = append(errs, errors.Errorf("no network_interface found for master.%d", idx))
		}

		master, found, err := unstructured.NestedString(networkInterfaces[0].(map[string]interface{}), "network_ip")
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "no network_ip for master.%d", idx))
		}
		if !found {
			errs = append(errs, errors.Errorf("no network_ip found for master.%d", idx))
		}
		masters = append(masters, master)
	}
	return masters, utilerrors.NewAggregate(errs)
}
