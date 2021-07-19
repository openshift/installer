package ibmcloud

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["ibmcloud"] = New
}
