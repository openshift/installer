package alibabacloud

import "github.com/openshift/installer/pkg/destroy/providers"

func init() {
	providers.Registry["alibabacloud"] = New
}
