// Package kubevirt supply utilities to extract information from terraform state
package kubevirt

import (
	"github.com/openshift/installer/pkg/terraform"
)

// BootstrapIP returns the ip address for bootstrap host.
// still unsupported, because qemu-ga is not available - see https://bugzilla.redhat.com/show_bug.cgi?id=1764804
func BootstrapIP(tfs *terraform.State) (string, error) {
	return "", nil
}

// ControlPlaneIPs returns the ip addresses for control plane hosts.
// still unsupported, because qemu-ga is not available  - see https://bugzilla.redhat.com/show_bug.cgi?id=1764804
func ControlPlaneIPs(tfs *terraform.State) ([]string, error) {
	return []string{""}, nil
}
