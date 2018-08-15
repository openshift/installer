package tls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/openshift/installer/pkg/installerassets"
)

func init() {
	installerassets.Defaults["tls/kubelet-client.key"] = privateKey
	installerassets.Rebuilders["tls/kubelet-client.crt"] = certificateRebuilder(
		"tls/kubelet-client.crt",
		"tls/kubelet-client.key",
		"tls/kube-ca.crt",
		"tls/kube-ca.key",
		&x509.Certificate{
			BasicConstraintsValid: true,
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			NotAfter:              time.Now().Add(validityThirtyMinutes),
			NotBefore:             time.Now(),
			SerialNumber:          new(big.Int).SetInt64(0),
			// system:masters is a hack to get the kubelet up without kube-core
			// TODO(node): make kubelet bootstrapping secure with minimal permissions eventually switching to system:node:* CommonName
			Subject: pkix.Name{CommonName: "system:serviceaccount:kube-system:default", Organization: []string{"system:serviceaccounts:kube-system", "system:masters"}},
		},
		nil,
	)
}
