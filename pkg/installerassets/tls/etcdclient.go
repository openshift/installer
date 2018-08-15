package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["tls/etcd-client.key"] = privateKey
	installerassets.Rebuilders["tls/etcd-client.crt"] = certificateRebuilder(
		"tls/etcd-client.crt",
		"tls/etcd-client.key",
		"tls/etcd-ca.crt",
		"tls/etcd-ca.key",
		&x509.Certificate{
			BasicConstraintsValid: true,
			KeyUsage:              x509.KeyUsageKeyEncipherment,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			NotAfter:              time.Now().Add(validityTenYears),
			NotBefore:             time.Now(),
			SerialNumber:          new(big.Int).SetInt64(0),
			Subject:               pkix.Name{CommonName: "etcd", OrganizationalUnit: []string{"etcd"}},
		},
		nil,
	)
}
