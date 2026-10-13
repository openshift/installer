package aws

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["aws"] = New
}
