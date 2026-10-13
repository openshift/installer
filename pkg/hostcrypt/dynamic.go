//go:build fipscapable
// +build fipscapable

package hostcrypt

import "fmt"

func allowFIPSCluster() error {
	fipsEnabled, err := hostFIPSEnabled()
	if err != nil {
		return err
	}
	if fipsEnabled {
		return nil
	}
	return fmt.Errorf("enable FIPS mode on the host")
}
