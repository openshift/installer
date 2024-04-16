// Package baremetal provides a cluster-destroyer for bare metal clusters.
package baremetal

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["baremetal"] = New
}
