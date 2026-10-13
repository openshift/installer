// package arch contains mappings between the Golang architecture and
// the RPM architecture used by Fedora CoreOS and derivatives.
package arch

import "runtime"

type mapping struct {
	rpmArch string
	goArch  string
}

// If an architecture isn't defined here, we assume it's
// pass through.
var translations = []mapping{
	{
		rpmArch: "x86_64",
		goArch:  "amd64",
	},
	{
		rpmArch: "aarch64",
		goArch:  "arm64",
	},
}

// CurrentRpmArch returns the current architecture in RPM terms.
func CurrentRpmArch() string {
	return RpmArch(runtime.GOARCH)
}

// RpmArch translates a Go architecture to RPM.
func RpmArch(goarch string) string {
	for _, m := range translations {
		if m.goArch == goarch {
			return m.rpmArch
		}
	}
	return goarch
}

// GoArch translates an RPM architecture to Go.
func GoArch(rpmarch string) string {
	for _, m := range translations {
		if m.rpmArch == rpmarch {
			return m.goArch
		}
	}
	return rpmarch
}
