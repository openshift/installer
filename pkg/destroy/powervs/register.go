package powervs

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["powervs"] = New
}
