package hostcrypt

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	fipsFile = "/proc/sys/crypto/fips_enabled"
)

// VerifyHostTargetState checks that the current binary matches the expected cryptographic state
// for the target cluster.
func VerifyHostTargetState(fips bool) error {
	if !fips {
		return nil
	}
	if runtime.GOOS != "linux" {
		return fmt.Errorf("target cluster is in FIPS mode, operation requires a Linux client")
	}

	hostFIPSData, err := os.ReadFile(fipsFile)
	if err != nil {
		return fmt.Errorf("target cluster is in FIPS mode, but failed to read client FIPS state %s: %w", fipsFile, err)
	}

	hostFIPS, err := strconv.ParseBool(strings.TrimSuffix(string(hostFIPSData), "\n"))
	if err != nil {
		return fmt.Errorf("target cluster is in FIPS mode, but failed to parse client FIPS state %s: %w", fipsFile, err)
	}

	if !hostFIPS {
		return fmt.Errorf("target cluster is in FIPS mode, operation requires a FIPS enabled client")
	}

	return nil
}
