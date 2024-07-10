package manifests

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/library-go/pkg/crypto"
)

type configurationObject struct {
	metav1.TypeMeta

	Metadata metadata    `json:"metadata,omitempty"`
	Data     genericData `json:"data,omitempty"`
}

type metadata struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
}

// TLSCertificate contains the private key and certificate contents in []byte format.
type TLSCertificate struct {
	privateKey  []byte
	certificate []byte
}

const (
	tlsExpirationDays = 1
)

func configMap(namespace, name string, data genericData) *configurationObject {
	return &configurationObject{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		Metadata: metadata{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

func getAPIServerURL(ic *types.InstallConfig) string {
	return fmt.Sprintf("https://api.%s:6443", ic.ClusterDomain())
}

func getInternalAPIServerURL(ic *types.InstallConfig) string {
	return fmt.Sprintf("https://api-int.%s:6443", ic.ClusterDomain())
}

func generateTLSCertificate(provisioningIP string) (TLSCertificate, error) {
	caConfig, err := crypto.MakeSelfSignedCAConfig("metal3-ironic", tlsExpirationDays)
	if err != nil {
		return TLSCertificate{}, err
	}

	ca := crypto.CA{
		Config:          caConfig,
		SerialGenerator: &crypto.RandomSerialGenerator{},
	}

	var host string
	if provisioningIP == "" {
		host = "localhost"
	} else {
		host = provisioningIP
	}

	config, err := ca.MakeServerCert(sets.New(host), tlsExpirationDays)

	if err != nil {
		return TLSCertificate{}, err
	}

	certBytes, keyBytes, err := config.GetPEMBytes()
	if err != nil {
		return TLSCertificate{}, err
	}

	return TLSCertificate{
		privateKey:  keyBytes,
		certificate: certBytes,
	}, nil
}
