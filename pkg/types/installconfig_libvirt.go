// +build libvirt

package types

import (
	"sort"

	"github.com/openshift/installer/pkg/types/libvirt"
)

func init() {
	PlatformNames = append(PlatformNames, libvirt.Name)
	sort.Strings(PlatformNames)
}
