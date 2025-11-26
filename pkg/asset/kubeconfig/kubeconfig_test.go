package kubeconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/types"
)

type testCertKey struct {
	key  string
	cert string
}

func (t *testCertKey) Key() []byte {
	return []byte(t.key)
}

func (t *testCertKey) Cert() []byte {
	return []byte(t.cert)
}

func TestKubeconfigGenerate(t *testing.T) {
	rootCA := &testCertKey{
		key:  "THIS IS ROOT CA KEY DATA",
		cert: "THIS IS ROOT CA CERT DATA",
	}

	adminCert := &testCertKey{
		key:  "THIS IS ADMIN KEY DATA",
		cert: "THIS IS ADMIN CERT DATA",
	}

	kubeletCert := &testCertKey{
		key:  "THIS IS KUBELET KEY DATA",
		cert: "THIS IS KUBELET CERT DATA",
	}

	installConfig := &types.InstallConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster-name",
		},
		BaseDomain: "test.example.com",
	}

	tests := []struct {
		name         string
		userName     string
		filename     string
		clientCert   tls.CertKeyInterface
		apiURL       string
		expectedData []byte
	}{
		{
			name:       "admin kubeconfig",
			userName:   "admin",
			filename:   "auth/kubeconfig",
			clientCert: adminCert,
			apiURL:     "https://api-int.test-cluster-name.test.example.com:6443",
			expectedData: []byte(`clusters:
- cluster:
    certificate-authority-data: VEhJUyBJUyBST09UIENBIENFUlQgREFUQQ==
    server: https://api-int.test-cluster-name.test.example.com:6443
  name: test-cluster-name
contexts:
- context:
    cluster: test-cluster-name
    user: admin
  name: admin
current-context: admin
users:
- name: admin
  user:
    client-certificate-data: VEhJUyBJUyBBRE1JTiBDRVJUIERBVEE=
    client-key-data: VEhJUyBJUyBBRE1JTiBLRVkgREFUQQ==
`),
		},
		{
			name:       "kubelet kubeconfig",
			userName:   "kubelet",
			filename:   "auth/kubeconfig-kubelet",
			clientCert: kubeletCert,
			apiURL:     "https://api-int.test-cluster-name.test.example.com:6443",
			expectedData: []byte(`clusters:
- cluster:
    certificate-authority-data: VEhJUyBJUyBST09UIENBIENFUlQgREFUQQ==
    server: https://api-int.test-cluster-name.test.example.com:6443
  name: test-cluster-name
contexts:
- context:
    cluster: test-cluster-name
    user: kubelet
  name: kubelet
current-context: kubelet
users:
- name: kubelet
  user:
    client-certificate-data: VEhJUyBJUyBLVUJFTEVUIENFUlQgREFUQQ==
    client-key-data: VEhJUyBJUyBLVUJFTEVUIEtFWSBEQVRB
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kc := &kubeconfig{}
			err := kc.generate(rootCA, tt.clientCert, tt.apiURL, installConfig.GetName(), tt.userName, tt.filename)
			assert.NoError(t, err, "unexpected error generating config")
			actualFiles := kc.Files()
			assert.Equal(t, 1, len(actualFiles), "unexpected number of files generated")
			assert.Equal(t, tt.filename, actualFiles[0].Filename, "unexpected file name generated")
			assert.Equal(t, tt.expectedData, actualFiles[0].Data, "unexpected config")
		})
	}

}
