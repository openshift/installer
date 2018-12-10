package rhcos

import (
	"fmt"
)

// QEMU fetches the URL of the latest Red Hat CoreOS release.
func QEMU(build Build) string {
	return fmt.Sprintf("%s/%s/%s/%s", build.BaseURL, build.Channel, build.Meta.OSTreeVersion, build.Meta.Images.QEMU.Path)
}
