package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["tls/aggregator-ca.key"] = privateKey
	installerassets.Rebuilders["tls/aggregator-ca.crt"] = certificateRebuilder(
		"tls/aggregator-ca.crt",
		"tls/aggregator-ca.key",
		"tls/root-ca.crt",
		"tls/root-ca.key",
		&x509.Certificate{
			BasicConstraintsValid: true,
			IsCA:         true,
			KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			NotAfter:     time.Now().Add(validityTenYears),
			NotBefore:    time.Now(),
			SerialNumber: new(big.Int).SetInt64(0),
			Subject:      pkix.Name{CommonName: "aggregator", OrganizationalUnit: []string{"bootkube"}},
		},
		nil,
	)
}
