package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
)

var (
	// It is unfortunate that these need to be global variables. However, the
	// asset store creates new assets by type, so the tests cannot store behavior
	// state in the assets themselves.
	generationLog []string
	dependencies  map[reflect.Type][]asset.Asset
	onDiskAssets  map[reflect.Type]bool
)

func clearAssetBehaviors() {
	generationLog = []string{}
	dependencies = map[reflect.Type][]asset.Asset{}
	onDiskAssets = map[reflect.Type]bool{}
}

func dependenciesTestStoreAsset(a asset.Asset) []asset.Asset {
	return dependencies[reflect.TypeOf(a)]
}

func generateTestStoreAsset(a asset.Asset) error {
	generationLog = append(generationLog, a.Name())
	return nil
}

func fileTestStoreAsset(a asset.Asset) []*asset.File {
	return []*asset.File{{Filename: a.Name()}}
}

func loadTestStoreAsset(a asset.Asset) (bool, error) {
	return onDiskAssets[reflect.TypeOf(a)], nil
}

type testStoreAssetA struct{}

func (a *testStoreAssetA) Name() string {
	return "a"
}

func (a *testStoreAssetA) Dependencies() []asset.Asset {
	return dependenciesTestStoreAsset(a)
}

func (a *testStoreAssetA) Generate(asset.Parents) error {
	return generateTestStoreAsset(a)
}

func (a *testStoreAssetA) Files() []*asset.File {
	return fileTestStoreAsset(a)
}

func (a *testStoreAssetA) Load(asset.FileFetcher) (bool, error) {
	return loadTestStoreAsset(a)
}

type testStoreAssetB struct{}

func (a *testStoreAssetB) Name() string {
	return "b"
}

func (a *testStoreAssetB) Dependencies() []asset.Asset {
	return dependenciesTestStoreAsset(a)
}

func (a *testStoreAssetB) Generate(asset.Parents) error {
	return generateTestStoreAsset(a)
}

func (a *testStoreAssetB) Files() []*asset.File {
	return fileTestStoreAsset(a)
}

func (a *testStoreAssetB) Load(asset.FileFetcher) (bool, error) {
	return loadTestStoreAsset(a)
}

type testStoreAssetC struct{}

func (a *testStoreAssetC) Name() string {
	return "c"
}

func (a *testStoreAssetC) Dependencies() []asset.Asset {
	return dependenciesTestStoreAsset(a)
}

func (a *testStoreAssetC) Generate(asset.Parents) error {
	return generateTestStoreAsset(a)
}

func (a *testStoreAssetC) Files() []*asset.File {
	return fileTestStoreAsset(a)
}

func (a *testStoreAssetC) Load(asset.FileFetcher) (bool, error) {
	return loadTestStoreAsset(a)
}

type testStoreAssetD struct{}

func (a *testStoreAssetD) Name() string {
	return "d"
}

func (a *testStoreAssetD) Dependencies() []asset.Asset {
	return dependenciesTestStoreAsset(a)
}

func (a *testStoreAssetD) Generate(asset.Parents) error {
	return generateTestStoreAsset(a)
}

func (a *testStoreAssetD) Files() []*asset.File {
	return fileTestStoreAsset(a)
}

func (a *testStoreAssetD) Load(asset.FileFetcher) (bool, error) {
	return loadTestStoreAsset(a)
}

func newTestStoreAsset(name string) asset.Asset {
	switch name {
	case "a":
		return &testStoreAssetA{}
	case "b":
		return &testStoreAssetB{}
	case "c":
		return &testStoreAssetC{}
	case "d":
		return &testStoreAssetD{}
	default:
		return nil
	}
}

// TestStoreFetch tests the Fetch method of StoreImpl.
func TestStoreFetch(t *testing.T) {
	cases := []struct {
		name                  string
		assets                map[string][]string
		existingAssets        []string
		target                string
		expectedGenerationLog []string
	}{
		{
			name: "no dependencies",
			assets: map[string][]string{
				"a": {},
			},
			target:                "a",
			expectedGenerationLog: []string{"a"},
		},
		{
			name: "single dependency",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			target:                "a",
			expectedGenerationLog: []string{"b", "a"},
		},
		{
			name: "multiple dependencies",
			assets: map[string][]string{
				"a": {"b", "c"},
				"b": {},
				"c": {},
			},
			target:                "a",
			expectedGenerationLog: []string{"b", "c", "a"},
		},
		{
			name: "grandchild dependency",
			assets: map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {},
			},
			target:                "a",
			expectedGenerationLog: []string{"c", "b", "a"},
		},
		{
			name: "intragenerational shared dependency",
			assets: map[string][]string{
				"a": {"b", "c"},
				"b": {"d"},
				"c": {"d"},
				"d": {},
			},
			target:                "a",
			expectedGenerationLog: []string{"d", "b", "c", "a"},
		},
		{
			name: "intergenerational shared dependency",
			assets: map[string][]string{
				"a": {"b", "c"},
				"b": {"c"},
				"c": {},
			},
			target:                "a",
			expectedGenerationLog: []string{"c", "b", "a"},
		},
		{
			name: "existing asset",
			assets: map[string][]string{
				"a": {},
			},
			existingAssets:        []string{"a"},
			target:                "a",
			expectedGenerationLog: []string{},
		},
		{
			name: "existing child asset",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			existingAssets:        []string{"b"},
			target:                "a",
			expectedGenerationLog: []string{"a"},
		},
		{
			name: "absent grandchild asset",
			assets: map[string][]string{
				"a": {"b"},
				"b": {"c"},
				"c": {},
			},
			existingAssets:        []string{"b"},
			target:                "a",
			expectedGenerationLog: []string{"a"},
		},
		{
			name: "absent grandchild with absent parent",
			assets: map[string][]string{
				"a": {"b", "c"},
				"b": {"d"},
				"c": {"d"},
				"d": {},
			},
			existingAssets:        []string{"b"},
			target:                "a",
			expectedGenerationLog: []string{"d", "c", "a"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearAssetBehaviors()
			dir, err := ioutil.TempDir("", "TestStoreFetch")
			if err != nil {
				t.Fatalf("failed to create temporary directory: %v", err)
			}
			defer os.RemoveAll(dir)
			store := &storeImpl{
				directory: dir,
				assets:    map[reflect.Type]*assetState{},
			}
			assets := make(map[string]asset.Asset, len(tc.assets))
			for name := range tc.assets {
				assets[name] = newTestStoreAsset(name)
			}
			for name, deps := range tc.assets {
				dependenciesOfAsset := make([]asset.Asset, len(deps))
				for i, d := range deps {
					dependenciesOfAsset[i] = assets[d]
				}
				dependencies[reflect.TypeOf(assets[name])] = dependenciesOfAsset
			}
			for _, assetName := range tc.existingAssets {
				asset := assets[assetName]
				store.assets[reflect.TypeOf(asset)] = &assetState{
					asset:  asset,
					source: generatedSource,
				}
			}
			err = store.Fetch(assets[tc.target])
			assert.NoError(t, err, "error fetching asset")
			assert.EqualValues(t, tc.expectedGenerationLog, generationLog)
		})
	}
}

func TestStoreFetchOnDiskAssets(t *testing.T) {
	cases := []struct {
		name                  string
		assets                map[string][]string
		onDiskAssets          []string
		target                string
		expectedGenerationLog []string
		expectedDirty         bool
	}{
		{
			name: "no on-disk assets",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			onDiskAssets:          nil,
			target:                "a",
			expectedGenerationLog: []string{"b", "a"},
			expectedDirty:         false,
		},
		{
			name: "on-disk asset does not need dependent generation",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			onDiskAssets:          []string{"a"},
			target:                "a",
			expectedGenerationLog: []string{},
			expectedDirty:         false,
		},
		{
			name: "on-disk dependent asset causes re-generation",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			onDiskAssets:          []string{"b"},
			target:                "a",
			expectedGenerationLog: []string{"a"},
			expectedDirty:         true,
		},
		{
			name: "on-disk dependents invalidate all its children",
			assets: map[string][]string{
				"a": {"b", "c"},
				"b": {"d"},
				"c": {"d"},
				"d": {},
			},
			onDiskAssets:          []string{"d"},
			target:                "a",
			expectedGenerationLog: []string{"b", "c", "a"},
			expectedDirty:         true,
		},
		{
			name: "re-generate when both parents and children are on-disk",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			onDiskAssets:          []string{"a", "b"},
			target:                "a",
			expectedGenerationLog: []string{"a"},
			expectedDirty:         true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearAssetBehaviors()
			store := &storeImpl{
				assets: map[reflect.Type]*assetState{},
			}
			assets := make(map[string]asset.Asset, len(tc.assets))
			for name := range tc.assets {
				assets[name] = newTestStoreAsset(name)
			}
			for name, deps := range tc.assets {
				dependenciesOfAsset := make([]asset.Asset, len(deps))
				for i, d := range deps {
					dependenciesOfAsset[i] = assets[d]
				}
				dependencies[reflect.TypeOf(assets[name])] = dependenciesOfAsset
			}
			for _, name := range tc.onDiskAssets {
				onDiskAssets[reflect.TypeOf(assets[name])] = true
			}
			err := store.fetch(assets[tc.target], "")
			assert.NoError(t, err, "unexpected error")
			assert.EqualValues(t, tc.expectedGenerationLog, generationLog)
			assert.Equal(t, tc.expectedDirty, store.assets[reflect.TypeOf(assets[tc.target])].anyParentsDirty)
		})
	}
}

func TestStoreFetchIdempotency(t *testing.T) {
	clearAssetBehaviors()

	tempDir, err := ioutil.TempDir("", "TestStoreFetchIdempotency")
	if err != nil {
		t.Fatalf("could not create the temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	for i := 0; i < 2; i++ {
		store, err := newStore(tempDir)
		if !assert.NoError(t, err, "(loop %d) unexpected error creating store", i) {
			t.Fatal()
		}
		assets := []asset.WritableAsset{&testStoreAssetA{}, &testStoreAssetB{}}
		for _, a := range assets {
			err = store.Fetch(a, assets...)
			if !assert.NoError(t, err, "(loop %d) unexpected error fetching asset %q", a.Name()) {
				t.Fatal()
			}
			err = asset.PersistToFile(a, tempDir)
			if !assert.NoError(t, err, "(loop %d) unexpected error persisting asset %q", a.Name()) {
				t.Fatal()
			}
			onDiskAssets[reflect.TypeOf(a)] = true
		}
	}

	expectedFiles := []string{"a", "b"}
	actualFiles := []string{}
	walkFunc := func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() || fi.Name() == stateFileName {
			return nil
		}
		actualFiles = append(actualFiles, fi.Name())
		return nil
	}
	filepath.Walk(tempDir, walkFunc)
	assert.Equal(t, expectedFiles, actualFiles, "unexpected files on disk")
}

func TestStoreLoadOnDiskAssets(t *testing.T) {
	cases := []struct {
		name               string
		assets             map[string][]string
		onDiskAssets       []string
		target             string
		expectedFoundValue bool
	}{
		{
			name: "on-disk assets",
			assets: map[string][]string{
				"a": {},
			},
			onDiskAssets:       []string{"a"},
			target:             "a",
			expectedFoundValue: true,
		},
		{
			name: "no on-disk assets",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			onDiskAssets:       nil,
			target:             "a",
			expectedFoundValue: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clearAssetBehaviors()
			store := &storeImpl{
				assets: map[reflect.Type]*assetState{},
			}
			assets := make(map[string]asset.Asset, len(tc.assets))
			for name := range tc.assets {
				assets[name] = newTestStoreAsset(name)
			}
			for name, deps := range tc.assets {
				dependenciesOfAsset := make([]asset.Asset, len(deps))
				for i, d := range deps {
					dependenciesOfAsset[i] = assets[d]
				}
				dependencies[reflect.TypeOf(assets[name])] = dependenciesOfAsset
			}
			for _, name := range tc.onDiskAssets {
				onDiskAssets[reflect.TypeOf(assets[name])] = true
			}
			found, err := store.Load(assets[tc.target])
			assert.NoError(t, err, "unexpected error")
			assert.EqualValues(t, tc.expectedFoundValue, found != nil)
		})
	}
}
