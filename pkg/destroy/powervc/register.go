package powervc

import (
	"github.com/openshift/installer/pkg/destroy/providers"
)

func init() {
	providers.Registry["powervc"] = New
}
