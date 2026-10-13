package gcp

import "github.com/openshift/installer/pkg/gather/providers"

func init() {
	providers.Registry["gcp"] = New
}
