package aws

import (
	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["aws/external-vpc-id"] = installerassets.ConstantDefault(nil)
}
