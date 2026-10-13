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

	if err := allowFIPSCluster(); err != nil {
		return fmt.Errorf("target cluster is in FIPS mode, %w", err)
	}
	return nil
}

func hostFIPSEnabled() (bool, error) {
	if runtime.GOOS != "linux" {
		return false, fmt.Errorf("operation requires a Linux client")
	}

	hostFIPSData, err := os.ReadFile(fipsFile)
	if err != nil {
		return false, fmt.Errorf("failed to read client FIPS state %s: %w", fipsFile, err)
	}

	hostFIPS, err := strconv.ParseBool(strings.TrimSuffix(string(hostFIPSData), "\n"))
	if err != nil {
		return false, fmt.Errorf("failed to parse client FIPS state %s: %w", fipsFile, err)
	}

	return hostFIPS, nil
}
