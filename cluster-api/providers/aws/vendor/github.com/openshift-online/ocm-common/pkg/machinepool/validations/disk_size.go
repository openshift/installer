package validations

import (
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"
)

const (
	machinePoolRootAWSVolumeSizeMin = 128
	// The following constants are in the helper file because putting them in the models creates
	// a circular dependency between the models and the helpers
	// machinePoolRootVolumeSizeMaxBefore414 is the maximum size of the root volume before 4.14
	// 1 TiB - limit before 4.14 due to some filesystem growing issues
	machinePoolRootVolumeSizeMaxBefore414 = 1024
	// machinePoolRootVolumeSizeMaxAsOf414 is the maximum size of the root volume as of 4.14
	// 16 TiB - limit as of 4.14
	machinePoolRootVolumeSizeMaxAsOf414 = 16384
	// constants for node pool root size validation
	nodePoolRootAWSVolumeSizeMin = 75
	nodePoolRootAWSVolumeSizeMax = 16384
)

// ValidateMachinePoolRootDiskSize validates the root volume size for a machine pool in AWS.
func ValidateMachinePoolRootDiskSize(version string, machinePoolRootVolumeSize int) error {
	machinePoolRootVolumeSizeMax, err := getAWSVolumeMaxSize(version)
	if err != nil {
		return err
	}

	if machinePoolRootVolumeSize < machinePoolRootAWSVolumeSizeMin ||
		machinePoolRootVolumeSize > machinePoolRootVolumeSizeMax {
		return fmt.Errorf("Invalid root disk size: %d GiB. Must be between %d GiB and %d GiB.",
			machinePoolRootVolumeSize,
			machinePoolRootAWSVolumeSizeMin,
			machinePoolRootVolumeSizeMax)
	}

	return nil
}

// ValidateNodePoolRootDiskSize validates the root volume size for a node pool in AWS.
func ValidateNodePoolRootDiskSize(nodePoolRootVolumeSize int) error {
	if nodePoolRootVolumeSize < nodePoolRootAWSVolumeSizeMin ||
		nodePoolRootVolumeSize > nodePoolRootAWSVolumeSizeMax {
		return fmt.Errorf("Invalid root disk size: %d GiB. Must be between %d GiB and %d GiB.",
			nodePoolRootVolumeSize,
			nodePoolRootAWSVolumeSizeMin,
			nodePoolRootAWSVolumeSizeMax)
	}

	return nil
}

// getAWSVolumeMaxSize returns the maximum size of the root volume for a machine pool in AWS.
func getAWSVolumeMaxSize(version string) (int, error) {
	version414, _ := semver.NewVersion("4.14.0")
	currentVersion, err := semver.NewVersion(strings.Replace(version, "openshift-v", "", 1))
	if err != nil {
		return 0, err
	}

	if currentVersion.GreaterThanOrEqual(version414) {
		return machinePoolRootVolumeSizeMaxAsOf414, nil
	}

	return machinePoolRootVolumeSizeMaxBefore414, nil
}
