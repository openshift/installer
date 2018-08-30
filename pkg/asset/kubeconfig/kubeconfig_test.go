package kubeconfig

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestKubeconfigGenerate(t *testing.T) {
	testDir, err := ioutil.TempDir(os.TempDir(), "kubeconfig_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

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
		userName     string
		certKey      asset.Asset
		parents      map[asset.Asset]*asset.State
		errString    string
		expectedData []byte
	}{
		{
			userName: "admin",
			certKey:  adminCertKey,
			parents: map[asset.Asset]*asset.State{
				rootCA:        rootCAState,
				adminCertKey:  adminCertState,
				installConfig: installConfigState,
			},
			errString: "",
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
			userName: "kubelet",
			certKey:  kubeletCertKey,
			parents: map[asset.Asset]*asset.State{
				rootCA:         rootCAState,
				kubeletCertKey: kubeletCertState,
				installConfig:  installConfigState,
			},
			errString: "",
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

	for i, tt := range tests {
		kubeconfig := &Kubeconfig{
			rootDir:       testDir,
			userName:      tt.userName,
			rootCA:        rootCA,
			certKey:       tt.certKey,
			installConfig: installConfig,
		}
		st, err := kubeconfig.Generate(tt.parents)
		if err != nil {
			assert.EqualErrorf(t, err, tt.errString, fmt.Sprintf("test #%d expected %v, saw %v", i, tt.errString, err))
			continue
		} else {
			if tt.errString != "" {
				t.Errorf("test #%d expect error %v, saw nil", i, tt.errString)
			}
		}

		filename := filepath.Join(testDir, "auth", fmt.Sprintf("kubeconfig-%s", tt.userName))
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Errorf("test #%d failed to read kubeconfig %q: %v", i, filename, err)
		}

		if !bytes.Equal(st.Contents[0].Data, data) {
			t.Errorf("test #%d expect kubeconfig data: %q, saw %q", i, string(st.Contents[0].Data), string(data))
		}

		if !bytes.Equal(tt.expectedData, data) {
			t.Errorf("test #%d expect kubeconfig data: %q, saw %q", i, string(tt.expectedData), string(data))
		}
	}
}
