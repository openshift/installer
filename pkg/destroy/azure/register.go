package azure

import (
	"github.com/openshift/installer/pkg/destroy"
)

func init() {
	destroy.Registry["azure"] = New
}
