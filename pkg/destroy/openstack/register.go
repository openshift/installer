package openstack

import (
	"github.com/openshift/installer/pkg/destroy/providers"
)

func init() {
	providers.Registry["openstack"] = New
}
