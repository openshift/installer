//go:build !libvirt
// +build !libvirt

package hostcrypt

import "fmt"

const binaryInstructions = "To obtain a suitable binary, download the openshift-install-rhel8 archive from the client mirror, or extract the openshift-install-fips command from the release payload."

func allowFIPSCluster() error {
	hostMsg := ""
	if fipsEnabled, err := hostFIPSEnabled(); err != nil || !fipsEnabled {
		hostMsg = " on a host with FIPS enabled"
	}
	return fmt.Errorf("use the FIPS-capable installer binary for RHEL 8%s.\n%s",
		hostMsg, binaryInstructions)
}
