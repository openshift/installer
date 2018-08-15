package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["tls/machine-config-server.key"] = privateKey
	installerassets.Rebuilders["tls/machine-config-server.crt"] = certificateRebuilder(
		"tls/machine-config-server.crt",
		"tls/machine-config-server.key",
		"tls/root-ca.crt",
		"tls/root-ca.key",
		&x509.Certificate{
			BasicConstraintsValid: true,
			KeyUsage:              x509.KeyUsageKeyEncipherment,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			NotAfter:              time.Now().Add(validityTenYears),
			NotBefore:             time.Now(),
			SerialNumber:          new(big.Int).SetInt64(0),
		},
		func(ctx context.Context, asset *assets.Asset, getByName assets.GetByString, template *x509.Certificate) (err error) {
			parents, err := asset.GetParents(
				ctx,
				getByName,
				"base-domain",
				"cluster-name",
			)
			if err != nil {
				return err
			}

			hostname := fmt.Sprintf("%s-api.%s", string(parents["cluster-name"].Data), string(parents["base-domain"].Data))
			template.Subject = pkix.Name{CommonName: hostname}
			template.DNSNames = []string{hostname}

			return nil
		},
	)
}
