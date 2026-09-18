package powervs

import "github.com/openshift/installer/pkg/types"

const (
	// OSImageNameRHCOS9 is the catalog image name for RHEL CoreOS 9.
	OSImageNameRHCOS9 = "RHEL-CoreOS-9"

	// OSImageNameRHCOS10 is the catalog image name for RHEL CoreOS 10.
	OSImageNameRHCOS10 = "RHEL-CoreOS-10"
)

// OSImageNameFromStream returns the PowerVS catalog image name that corresponds
// to the given osImageStream value from the install-config.
// types.OSImageStreamRHCOS10 ("rhel-10") maps to RHEL-CoreOS-10;
// anything else (including empty) defaults to RHEL-CoreOS-9.
func OSImageNameFromStream(stream types.OSImageStream) string {
	if stream == types.OSImageStreamRHCOS10 {
		return OSImageNameRHCOS10
	}
	return OSImageNameRHCOS9
}
