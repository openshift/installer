package rhcos

import (
	"os"
	"fmt"
)

// QEMU fetches the URL of the latest Red Hat CoreOS release.
func QEMU(build Build) string {
	qcowImage, ok := os.LookupEnv("OPENSHIFT_INSTALL_LIBVIRT_IMAGE")
	if ok {
		return qcowImage
	}
	return fmt.Sprintf("%s/%s/%s/%s", build.BaseURL, build.Channel, build.Meta.OSTreeVersion, build.Meta.Images.QEMU.Path)
}
