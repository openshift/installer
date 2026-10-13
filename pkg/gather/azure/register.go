package azure

import "github.com/openshift/installer/pkg/gather/providers"

func init() {
	providers.Registry["azure"] = New
}
