package kubeconfig

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/openshift/installer/pkg/asset"
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
				Data: []byte(`admin:
  email: foo
  password: foo
baseDomain: test.base.domain
clusterID: e93844a3-a727-11e8-ab80-28d244b00276
license: foo
machines:
metadata:
  creationTimestamp: null
  name: test-cluster-name
networking:
  podCIDR:
    IP: ""
    Mask: null
  serviceCIDR:
    IP: ""
    Mask: null
  type: ""
platform:
  aws:
    keyPairName: foo
    region: us-west-2
    vpcCIDRBlock: ""
    vpcID: ""
pullSecret: foo
`),
			},
		},
	}

	tests := []struct {
		userName     string
		certKey      asset.Asset
		parents      map[asset.Asset]*asset.State
		err          bool
		expectedData []byte
	}{
		{
			"admin",
			adminCertKey,
			map[asset.Asset]*asset.State{
				rootCA:        rootCAState,
				adminCertKey:  adminCertState,
				installConfig: installConfigState,
			},
			false,
			[]byte(`clusters:
- cluster:
    certificate-authority-data: VkVoSlV5QkpVeUJTVDA5VUlFTkJJRU5GVWxRZ1JFRlVRUT09
    server: https://test-cluster-name-api.test.base.domain:6443
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
    client-certificate-data: VkVoSlV5QkpVeUJCUkUxSlRpQkRSVkpVSUVSQlZFRT0=
    client-key-data: VkVoSlV5QkpVeUJCUkUxSlRpQkxSVmtnUkVGVVFRPT0=
`),
		},
		{
			"kubelet",
			kubeletCertKey,
			map[asset.Asset]*asset.State{
				rootCA:         rootCAState,
				kubeletCertKey: kubeletCertState,
				installConfig:  installConfigState,
			},
			false,
			[]byte(`clusters:
- cluster:
    certificate-authority-data: VkVoSlV5QkpVeUJTVDA5VUlFTkJJRU5GVWxRZ1JFRlVRUT09
    server: https://test-cluster-name-api.test.base.domain:6443
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
    client-certificate-data: VkVoSlV5QkpVeUJMVlVKRlRFVlVJRU5GVWxRZ1JFRlVRUT09
    client-key-data: VkVoSlV5QkpVeUJMVlVKRlRFVlVJRXRGV1NCRVFWUkI=
`),
		},
		{
			"admin", // No root ca in parents.
			adminCertKey,
			map[asset.Asset]*asset.State{
				adminCertKey:  adminCertState,
				installConfig: installConfigState,
			},
			true,
			nil,
		},
		{
			"admin", // No admin cert in parents.
			adminCertKey,
			map[asset.Asset]*asset.State{
				rootCA:         rootCAState,
				kubeletCertKey: kubeletCertState,
				installConfig:  installConfigState,
			},
			true,
			nil,
		},
		{
			"kubelet", // No kubelet cert in parents.
			kubeletCertKey,
			map[asset.Asset]*asset.State{
				rootCA:        rootCAState,
				adminCertKey:  adminCertState,
				installConfig: installConfigState,
			},
			true,
			nil,
		},
		{
			"admin", // No install config in parents.
			adminCertKey,
			map[asset.Asset]*asset.State{
				rootCA:       rootCAState,
				adminCertKey: adminCertState,
			},
			true,
			nil,
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
		if tt.err != (err != nil) {
			t.Errorf("test #%d error is not expected, expect %v, saw %v, err: %v", i, tt.err, err != nil, err)
		}
		if err != nil {
			continue
		}

		filename := assetFilePath(testDir, tt.userName)
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
