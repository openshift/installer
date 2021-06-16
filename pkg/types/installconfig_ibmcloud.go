// +build ibmcloud

package types

import (
	"sort"

	"github.com/openshift/installer/pkg/types/ibmcloud"
)

func init() {
	PlatformNames = append(PlatformNames, ibmcloud.Name)
	sort.Strings(PlatformNames)
}
