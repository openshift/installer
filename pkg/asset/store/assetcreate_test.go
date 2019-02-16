package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/machines"
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

func TestCreatedAssetsAreNotDirty(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	cases := []struct {
		name    string
		targets []asset.WritableAsset
	}{
		{
			name:    "install config",
			targets: targets.InstallConfig,
		},
		{
			name:    "manifest templates",
			targets: targets.ManifestTemplates,
		},
		{
			name:    "manifests",
			targets: targets.Manifests,
		},
		{
			name:    "ignition configs",
			targets: targets.IgnitionConfigs,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir, err := ioutil.TempDir("", "TestCreatedAssetsAreNotDirty")
			if err != nil {
				t.Fatalf("could not create the temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			if err := ioutil.WriteFile(filepath.Join(tempDir, stateFileName), []byte(userProvidedAssets), 0666); err != nil {
				t.Fatalf("could not write the state file: %v", err)
			}

			assetStore, err := newStore(tempDir)
			if err != nil {
				t.Fatalf("failed to create asset store: %v", err)
			}

			for _, a := range tc.targets {
				if err := assetStore.Fetch(a); err != nil {
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

			emptyAssets := map[reflect.Type]bool{
				reflect.TypeOf(&machines.ControlPlane{}): true, // no files for the 'none' platform
				reflect.TypeOf(&cluster.Metadata{}):      true, // read-only
			}
			for _, a := range tc.targets {
				newAsset := reflect.New(reflect.TypeOf(a).Elem()).Interface().(asset.WritableAsset)
				if err := newAssetStore.Fetch(newAsset); err != nil {
					t.Fatalf("failed to fetch %q in new store: %v", a.Name(), err)
				}
				if _, ok := emptyAssets[reflect.TypeOf(a)]; !ok {
					assetState := newAssetStore.assets[reflect.TypeOf(a)]
					assert.Truef(t, assetState.presentOnDisk, "asset %q was not found on disk", a.Name())
				}
			}

			assert.Equal(t, len(assetStore.assets), len(newAssetStore.assets), "new asset store does not have the same number of assets as original")

			for _, a := range newAssetStore.assets {
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
