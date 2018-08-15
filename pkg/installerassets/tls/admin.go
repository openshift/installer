package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["tls/admin-client.key"] = privateKey
	installerassets.Rebuilders["tls/admin-client.crt"] = certificateRebuilder(
		"tls/admin-client.crt",
		"tls/admin-client.key",
		"tls/kube-ca.crt",
		"tls/kube-ca.key",
		&x509.Certificate{
			BasicConstraintsValid: true,
			IsCA:         true,
			KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth /* FIXME: why? */, x509.ExtKeyUsageClientAuth},
			NotAfter:     time.Now().Add(validityTenYears),
			NotBefore:    time.Now(),
			SerialNumber: new(big.Int).SetInt64(0),
			Subject:      pkix.Name{CommonName: "system:admin", OrganizationalUnit: []string{"system:masters"}},
		},
		nil,
	)
}
