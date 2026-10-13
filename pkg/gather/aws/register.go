package aws

import "github.com/openshift/installer/pkg/gather/providers"

func init() {
	providers.Registry["aws"] = New
}
