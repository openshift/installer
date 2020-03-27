// Package azure contains utilities that help gather Azure specific
// information from terraform state.
package azure

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"

	"github.com/openshift/installer/pkg/terraform"
)

// BootstrapIP returns the ip address for bootstrap host.
func BootstrapIP(tfs *terraform.State) (string, error) {
	var bootstrap string

	publicIP, err := terraform.LookupResource(tfs, "module.bootstrap", "azurerm_public_ip", "bootstrap_public_ip")
	if err != nil {
		publicIP, err = terraform.LookupResource(tfs, "module.bootstrap", "azurerm_public_ip", "bootstrap_public_ip_v4")
	}
	if err != nil {
		publicIP, err = terraform.LookupResource(tfs, "module.bootstrap", "azurerm_public_ip", "bootstrap_public_ip_v6")
	}
	if err == nil && len(publicIP.Instances) > 0 {
		bootstrap, _, err = unstructured.NestedString(publicIP.Instances[0].Attributes, "ip_address")
		if err != nil {
			return "", errors.New("no public_ip found for bootstrap")
		}
		return bootstrap, nil
	}

	br, err := terraform.LookupResource(tfs, "module.bootstrap", "azurerm_network_interface", "bootstrap")
	if err != nil {
		return "", errors.Wrap(err, "failed to lookup bootstrap network interface")
	}
	if len(br.Instances) == 0 {
		return "", errors.New("no bootstrap instance found")
	}
	bootstrap, _, err = unstructured.NestedString(br.Instances[0].Attributes, "private_ip_address")
	if err != nil {
		return "", errors.New("no private_ip_address found for bootstrap")
	}

	return bootstrap, nil
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	mrs, err := terraform.LookupResource(tfs, "module.master", "azurerm_network_interface", "master")
	if err != nil {
		return nil, errors.Wrap(err, "failed to lookup masters")
	}
	var errs []error
	var masters []string
	for idx, inst := range mrs.Instances {
		master, _, err := unstructured.NestedString(inst.Attributes, "private_ip_address")
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "no private_ip for master.%d", idx))
		}
		masters = append(masters, master)
	}
	return masters, utilerrors.NewAggregate(errs)
}
