package tls

import (
	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["tls/service-account.key"] = privateKey
}
