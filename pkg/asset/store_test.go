package asset

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type generationLog struct {
	log []string
}

func (l *generationLog) logGeneration(a Asset) {
	l.log = append(l.log, a.Name())
}

type testStoreAsset interface {
	Asset
	SetDependencies([]Asset)
	SetDirty(bool)
}

type testStoreAssetImpl struct {
	name          string
	dependencies  []Asset
	generationLog *generationLog
	dirty         bool
}

func (a *testStoreAssetImpl) Dependencies() []Asset {
	return a.dependencies
}

func (a *testStoreAssetImpl) Generate(Parents) error {
	a.generationLog.logGeneration(a)
	return nil
}

func (a *testStoreAssetImpl) Name() string {
	return a.name
}

func (a *testStoreAssetImpl) SetDependencies(dependencies []Asset) {
	a.dependencies = dependencies
}

func (a *testStoreAssetImpl) SetDirty(dirty bool) {
	a.dirty = dirty
}
func (a *testStoreAssetImpl) Files() []*File {
	return []*File{{Filename: a.name}}
}

func (a *testStoreAssetImpl) Load(FileFetcher) (bool, error) {
	return a.dirty, nil
}

type testStoreAssetA struct {
	testStoreAssetImpl
}

type testStoreAssetB struct {
	testStoreAssetImpl
}

type testStoreAssetC struct {
	testStoreAssetImpl
}

type testStoreAssetD struct {
	testStoreAssetImpl
}

func newTestStoreAsset(gl *generationLog, name string) testStoreAsset {
	ta := testStoreAssetImpl{
		name:          name,
		generationLog: gl,
	}
	switch name {
	case "a":
		return &testStoreAssetA{ta}
	case "b":
		return &testStoreAssetB{ta}
	case "c":
		return &testStoreAssetC{ta}
	case "d":
		return &testStoreAssetD{ta}
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
			gl := &generationLog{
				log: []string{},
			}
			dir, err := ioutil.TempDir("", "TestStoreFetch")
			if err != nil {
				t.Fatalf("failed to create temporary directory: %v", err)
			}
			defer os.RemoveAll(dir)
			store := &StoreImpl{
				directory: dir,
				assets:    make(map[reflect.Type]assetState),
			}
			assets := make(map[string]testStoreAsset, len(tc.assets))
			for name := range tc.assets {
				assets[name] = newTestStoreAsset(gl, name)
			}
			for name, deps := range tc.assets {
				dependencies := make([]Asset, len(deps))
				for i, d := range deps {
					dependencies[i] = assets[d]
				}
				assets[name].SetDependencies(dependencies)
			}
			for _, assetName := range tc.existingAssets {
				asset := assets[assetName]
				store.assets[reflect.TypeOf(asset)] = assetState{asset: asset}
			}
			err = store.Fetch(assets[tc.target])
			assert.NoError(t, err, "error fetching asset")
			assert.EqualValues(t, tc.expectedGenerationLog, gl.log)
		})
	}
}

func TestStoreFetchDirty(t *testing.T) {
	cases := []struct {
		name                  string
		assets                map[string][]string
		dirtyAssets           []string
		target                string
		expectedGenerationLog []string
		expectedDirty         bool
	}{
		{
			name: "no dirty assets",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			dirtyAssets:           nil,
			target:                "a",
			expectedGenerationLog: []string{"b", "a"},
			expectedDirty:         false,
		},
		{
			name: "dirty asset causes re-generation",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			dirtyAssets:           []string{"a"},
			target:                "a",
			expectedGenerationLog: []string{"b"},
			expectedDirty:         true,
		},
		{
			name: "dirty dependent asset causes re-generation",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			dirtyAssets:           []string{"b"},
			target:                "a",
			expectedGenerationLog: []string{"a"},
			expectedDirty:         true,
		},
		{
			name: "dirty dependents invalidate all its children",
			assets: map[string][]string{
				"a": {"b", "c"},
				"b": {"d"},
				"c": {"d"},
				"d": {},
			},
			dirtyAssets:           []string{"d"},
			target:                "a",
			expectedGenerationLog: []string{"b", "c", "a"},
			expectedDirty:         true,
		},
		{
			name: "re-generate when both parents and childre are dirty",
			assets: map[string][]string{
				"a": {"b"},
				"b": {},
			},
			dirtyAssets:           []string{"a", "b"},
			target:                "a",
			expectedGenerationLog: []string{"a"},
			expectedDirty:         true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gl := &generationLog{
				log: []string{},
			}
			store := &StoreImpl{
				assets: make(map[reflect.Type]assetState),
			}
			assets := make(map[string]testStoreAsset, len(tc.assets))
			for name := range tc.assets {
				assets[name] = newTestStoreAsset(gl, name)
			}
			for name, deps := range tc.assets {
				dependencies := make([]Asset, len(deps))
				for i, d := range deps {
					dependencies[i] = assets[d]
				}
				assets[name].SetDependencies(dependencies)
			}
			for _, name := range tc.dirtyAssets {
				assets[name].SetDirty(true)
			}
			dirty, err := store.fetch(assets[tc.target], "")
			assert.NoError(t, err, "unexpected error")
			assert.EqualValues(t, tc.expectedGenerationLog, gl.log)
			assert.Equal(t, tc.expectedDirty, dirty)
		})
	}
}
