package kubeconfig

import (
	"fmt"
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/stretchr/testify/assert"
)

type fakeAsset int

var _ asset.Asset = fakeAsset(0)

func (f fakeAsset) Dependencies() []asset.Asset {
	return nil
}

func (f fakeAsset) Generate(map[asset.Asset]*asset.State) (*asset.State, error) {
	return nil, nil
}

func (f fakeAsset) Name() string {
	return "Fake Asset"
}

func TestKubeconfigGenerate(t *testing.T) {
	rootCA := fakeAsset(0)
	adminCertKey := fakeAsset(1)
	kubeletCertKey := fakeAsset(2)
	installConfig := fakeAsset(3)

	rootCAState := &asset.State{
		Contents: []asset.Content{
			{
				Name: "root-ca.key",
				Data: []byte("THIS IS ROOT CA KEY DATA"),
			},
			{
				Name: "root-ca.crt",
				Data: []byte("THIS IS ROOT CA CERT DATA"),
			},
		},
	}

	adminCertState := &asset.State{
		Contents: []asset.Content{
			{
				Name: "admin.key",
				Data: []byte("THIS IS ADMIN KEY DATA"),
			},
			{
				Name: "admin.crt",
				Data: []byte("THIS IS ADMIN CERT DATA"),
			},
		},
	}

	kubeletCertState := &asset.State{
		Contents: []asset.Content{
			{
				Name: "kubelet.key",
				Data: []byte("THIS IS KUBELET KEY DATA"),
			},
			{
				Name: "kubelet.crt",
				Data: []byte("THIS IS KUBELET CERT DATA"),
			},
		},
	}

	installConfigState := &asset.State{
		Contents: []asset.Content{
			{
				Name: "installconfig.yaml",
				Data: []byte(`baseDomain: test.example.com
metadata:
  name: test-cluster-name
`),
			},
		},
	}

	tests := []struct {
		name         string
		userName     string
		certKey      asset.Asset
		parents      map[asset.Asset]*asset.State
		errString    string
		expectedData []byte
	}{
		{
			name:     "admin kubeconfig",
			userName: "admin",
			certKey:  adminCertKey,
			parents: map[asset.Asset]*asset.State{
				rootCA:        rootCAState,
				adminCertKey:  adminCertState,
				installConfig: installConfigState,
			},
			expectedData: []byte(`clusters:
- cluster:
    certificate-authority-data: VEhJUyBJUyBST09UIENBIENFUlQgREFUQQ==
    server: https://test-cluster-name-api.test.example.com:6443
  name: test-cluster-name
contexts:
- context:
    cluster: test-cluster-name
    user: admin
  name: admin
current-context: admin
preferences: {}
users:
- name: admin
  user:
    client-certificate-data: VEhJUyBJUyBBRE1JTiBDRVJUIERBVEE=
    client-key-data: VEhJUyBJUyBBRE1JTiBLRVkgREFUQQ==
`),
		},
		{
			name:     "kubelet kubeconfig",
			userName: "kubelet",
			certKey:  kubeletCertKey,
			parents: map[asset.Asset]*asset.State{
				rootCA:         rootCAState,
				kubeletCertKey: kubeletCertState,
				installConfig:  installConfigState,
			},
			expectedData: []byte(`clusters:
- cluster:
    certificate-authority-data: VEhJUyBJUyBST09UIENBIENFUlQgREFUQQ==
    server: https://test-cluster-name-api.test.example.com:6443
  name: test-cluster-name
contexts:
- context:
    cluster: test-cluster-name
    user: kubelet
  name: kubelet
current-context: kubelet
preferences: {}
users:
- name: kubelet
  user:
    client-certificate-data: VEhJUyBJUyBLVUJFTEVUIENFUlQgREFUQQ==
    client-key-data: VEhJUyBJUyBLVUJFTEVUIEtFWSBEQVRB
`),
		},
		{
			name:     "no root ca",
			userName: "admin", // No root ca in parents.
			certKey:  adminCertKey,
			parents: map[asset.Asset]*asset.State{
				adminCertKey:  adminCertState,
				installConfig: installConfigState,
			},
			errString:    "failed to find kubeconfig.fakeAsset in parents",
			expectedData: nil,
		},
		{
			name:     "no admin certs",
			userName: "admin", // No admin cert in parents.
			certKey:  adminCertKey,
			parents: map[asset.Asset]*asset.State{
				rootCA:         rootCAState,
				kubeletCertKey: kubeletCertState,
				installConfig:  installConfigState,
			},
			errString:    "failed to find kubeconfig.fakeAsset in parents",
			expectedData: nil,
		},
		{
			name:     "no kubelet certs",
			userName: "kubelet", // No kubelet cert in parents.
			certKey:  kubeletCertKey,
			parents: map[asset.Asset]*asset.State{
				rootCA:        rootCAState,
				adminCertKey:  adminCertState,
				installConfig: installConfigState,
			},
			errString:    "failed to find kubeconfig.fakeAsset in parents",
			expectedData: nil,
		},
		{
			name:     "no install config",
			userName: "admin", // No install config in parents.
			certKey:  adminCertKey,
			parents: map[asset.Asset]*asset.State{
				rootCA:       rootCAState,
				adminCertKey: adminCertState,
			},
			errString:    "failed to find kubeconfig.fakeAsset in parents",
			expectedData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kubeconfig := &Kubeconfig{
				userName:      tt.userName,
				rootCA:        rootCA,
				certKey:       tt.certKey,
				installConfig: installConfig,
			}
			st, err := kubeconfig.Generate(tt.parents)
			if err != nil {
				assert.EqualErrorf(t, err, tt.errString, fmt.Sprintf("expected %v, saw %v", tt.errString, err))
				return
			} else if tt.errString != "" {
				t.Errorf("expect error %v, saw nil", tt.errString)
			}

			assert.Equal(t, tt.expectedData, st.Contents[0].Data, "unexpected data in kubeconfig")
		})
	}

}
