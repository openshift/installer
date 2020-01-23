// Package vsphere provides a cluster-destroyer for vsphere clusters.
package vsphere

import (
	"github.com/openshift/installer/pkg/destroy/providers"
)

func init() {
	providers.Registry["vsphere"] = New
}
