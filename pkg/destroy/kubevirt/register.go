package kubevirt

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["kubevirt"] = New
}
