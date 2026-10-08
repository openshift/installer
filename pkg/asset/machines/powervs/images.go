package powervs

import (
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types"
)

// OSImageNameFromStream returns the PowerVS catalog image name that corresponds
// to the given osImageStream value from the install-config.
// The stream value (e.g. "rhel-9", "rhel-10") is mapped to a catalog image
// name of the form "RHEL-CoreOS-<major>" (e.g. "RHEL-CoreOS-9").
// If the stream value is empty or does not contain a "-", it defaults to
// "RHEL-CoreOS-9".
func OSImageNameFromStream(stream types.OSImageStream) string {
	parts := strings.SplitN(string(stream), "-", 2)
	if len(parts) == 2 && parts[1] != "" {
		return fmt.Sprintf("RHEL-CoreOS-%s", parts[1])
	}
	return "RHEL-CoreOS-9"
}
