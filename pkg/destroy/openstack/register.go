// Package openstack provides a cluster-destroyer for openstack clusters.
package openstack

import (
	"github.com/openshift/installer/pkg/destroy"
)

func init() {
	destroy.Registry["openstack"] = New
}
