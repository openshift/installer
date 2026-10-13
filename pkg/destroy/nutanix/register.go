// Package nutanix provides a cluster-destroyer for nutanix clusters.
package nutanix

import (
	"github.com/openshift/installer/pkg/destroy/providers"
)

func init() {
	providers.Registry["nutanix"] = New
}
