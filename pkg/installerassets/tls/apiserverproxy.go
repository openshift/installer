package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["tls/api-server-proxy.key"] = privateKey
	installerassets.Rebuilders["tls/api-server-proxy.crt"] = certificateRebuilder(
		"tls/api-server-proxy.crt",
		"tls/api-server-proxy.key",
		"tls/kube-ca.crt",
		"tls/kube-ca.key",
		&x509.Certificate{
			BasicConstraintsValid: true,
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			NotAfter:              time.Now().Add(validityTenYears),
			NotBefore:             time.Now(),
			SerialNumber:          new(big.Int).SetInt64(0),
			Subject:               pkix.Name{CommonName: "system:kube-apiserver-proxy", OrganizationalUnit: []string{"kube-master"}},
		},
		nil,
	)
}
