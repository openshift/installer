package aws

import "github.com/openshift/installer/pkg/destroy"

func init() {
	destroy.Registry["aws"] = New
}
