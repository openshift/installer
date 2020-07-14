// +build baremetal

package types

import (
	"sort"

	"github.com/openshift/installer/pkg/types/baremetal"
)

func init() {
	PlatformNames = append(PlatformNames, baremetal.Name)
	sort.Strings(PlatformNames)
}
