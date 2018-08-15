package tls

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
)

func init() {
	installerassets.Defaults["tls/api-server.key"] = privateKey
	installerassets.Rebuilders["tls/api-server.crt"] = certificateRebuilder(
		"tls/api-server.crt",
		"tls/api-server.key",
		"tls/kube-ca.crt",
		"tls/kube-ca.key",
		&x509.Certificate{
			BasicConstraintsValid: true,
			IsCA:         true,
			KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			NotAfter:     time.Now().Add(validityTenYears),
			NotBefore:    time.Now(),
			SerialNumber: new(big.Int).SetInt64(0),
			Subject:      pkix.Name{CommonName: "system:kube-apiserver", OrganizationalUnit: []string{"kube-master"}},
		},
		func(ctx context.Context, asset *assets.Asset, getByName assets.GetByString, template *x509.Certificate) (err error) {
			parents, err := asset.GetParents(
				ctx,
				getByName,
				"base-domain",
				"cluster-name",
				"network/service-cidr",
			)
			if err != nil {
				return err
			}

			ip, ipnet, err := net.ParseCIDR(string(parents["network/service-cidr"].Data))
			if err != nil {
				return errors.Wrap(err, "parse service CIDR")
			}
			ipnet.IP = ip

			apiServerAddress, err := cidr.Host(ipnet, 1)
			if err != nil {
				return errors.Wrap(err, "calculate API-server address")
			}

			template.DNSNames = []string{
				fmt.Sprintf("%s-api.%s", string(parents["cluster-name"].Data), string(parents["base-domain"].Data)),
				"kubernetes",
				"kubernetes.default",
				"kubernetes.default.svc",
				"kubernetes.default.svc.cluster.local",
				"localhost",
			}
			template.IPAddresses = []net.IP{
				apiServerAddress,
				net.ParseIP("127.0.0.1"),
			}

			return nil
		},
	)

	installerassets.Rebuilders["tls/api-server-chain.crt"] = certificateChainRebuilder(
		"tls/api-server-chain.crt",
		"tls/api-server.crt",
		"tls/kube-ca.crt",
	)
}
