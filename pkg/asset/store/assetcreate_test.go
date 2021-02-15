package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/targets"
)

const userProvidedAssets = `{
  "*installconfig.baseDomain": {
    "BaseDomain": "test-domain"
  },
  "*installconfig.clusterID": {
    "ClusterID": "test-cluster-id"
  },
  "*installconfig.clusterName": {
    "ClusterName": "test-cluster"
  },
  "*installconfig.platform": {
    "none": {}
  },
  "*installconfig.pullSecret": {
    "PullSecret": "{\"auths\": {\"example.com\": {\"auth\": \"test-auth\"}}}\n"
  },
  "*installconfig.sshPublicKey": {}
}`

const singleNodeBootstrapInPlaceInstallConfig = `
apiVersion: v1
baseDomain: test-domain
metadata:
  name: test-cluster
platform:
  none: {}
controlPlane:
  replicas: 1
bootstrapInPlace:
  installationDisk: /dev/sda
pullSecret: |
  {
    "auths": {
      "example.com": {
        "auth": "test-auth"
       }
    }
  }
`

func TestCreatedAssetsAreNotDirty(t *testing.T) {
	cases := []struct {
		name    string
		targets []asset.WritableAsset
		files   map[string]string
	}{
		{
			name:    "install config",
			targets: targets.InstallConfig,
			files:   map[string]string{stateFileName: userProvidedAssets},
		},
		{
			name:    "manifest templates",
			targets: targets.ManifestTemplates,
			files:   map[string]string{stateFileName: userProvidedAssets},
		},
		{
			name:    "manifests",
			targets: targets.Manifests,
			files:   map[string]string{stateFileName: userProvidedAssets},
		},
		{
			name:    "ignition configs",
			targets: targets.IgnitionConfigs,
			files:   map[string]string{stateFileName: userProvidedAssets},
		},
		{
			name:    "single node bootstrap-in-place ignition config",
			targets: targets.SingleNodeIgnitionConfig,
			files:   map[string]string{"install-config.yaml": singleNodeBootstrapInPlaceInstallConfig},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := ioutil.TempDir("", "TestCreatedAssetsAreNotDirty")
			if err != nil {
				t.Fatalf("could not create the temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			for name, contents := range tc.files {
				if err := ioutil.WriteFile(filepath.Join(tempDir, name), []byte(contents), 0666); err != nil {
					t.Fatalf("could not write the %s file: %v", name, err)
				}
			}

			assetStore, err := newStore(tempDir)
			if err != nil {
				t.Fatalf("failed to create asset store: %v", err)
			}

			for _, a := range tc.targets {
				if err := assetStore.Fetch(a, tc.targets...); err != nil {
					t.Fatalf("failed to fetch %q: %v", a.Name(), err)
				}

				if err := asset.PersistToFile(a, tempDir); err != nil {
					t.Fatalf("failed to write asset %q to disk: %v", a.Name(), err)
				}
			}

			newAssetStore, err := newStore(tempDir)
			if err != nil {
				t.Fatalf("failed to create new asset store: %v", err)
			}

			emptyAssets := map[string]bool{
				"Master Machines":    true, // no files for the 'none' platform
				"Worker Machines":    true, // no files for the 'none' platform
				"Metadata":           true, // read-only
				"Kubeadmin Password": true, // read-only
			}
			for _, a := range tc.targets {
				name := a.Name()
				newAsset := reflect.New(reflect.TypeOf(a).Elem()).Interface().(asset.WritableAsset)
				if err := newAssetStore.Fetch(newAsset, tc.targets...); err != nil {
					t.Fatalf("failed to fetch %q in new store: %v", a.Name(), err)
				}
				assetState := newAssetStore.assets[reflect.TypeOf(a)]
				if !emptyAssets[name] {
					assert.Truef(t, assetState.presentOnDisk, "asset %q was not found on disk", a.Name())
				}
			}

			assert.Equal(t, len(assetStore.assets), len(newAssetStore.assets), "new asset store does not have the same number of assets as original")

			for _, a := range newAssetStore.assets {
				if a.source == unfetched {
					continue
				}
				if emptyAssets[a.asset.Name()] {
					continue
				}
				originalAssetState, ok := assetStore.assets[reflect.TypeOf(a.asset)]
				if !ok {
					t.Fatalf("asset %q not found in original store", a.asset.Name())
				}
				assert.Equalf(t, originalAssetState.asset, a.asset, "fetched and generated asset %q are not equal", a.asset.Name())
				assert.Equalf(t, stateFileSource, a.source, "asset %q was not fetched from the state file", a.asset.Name())
			}
		})
	}
}
