package asset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type generationLog struct {
	log []string
}

func (l *generationLog) logGeneration(a *testAsset) {
	l.log = append(l.log, a.name)
}

type testAsset struct {
	name          string
	dependencies  []Asset
	generationLog *generationLog
}

func (a *testAsset) Dependencies() []Asset {
	return a.dependencies
}

func (a *testAsset) Generate(map[Asset]*State) (*State, error) {
	a.generationLog.logGeneration(a)
	return nil, nil
}

func (a *testAsset) Name() string {
	return "Test Asset"
}

func newTestAsset(gl *generationLog, name string) *testAsset {
	return &testAsset{
		name:          name,
		generationLog: gl,
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
			store := &StoreImpl{
				assets: map[Asset]*State{},
			}
			assets := make(map[string]*testAsset, len(tc.assets))
			for name := range tc.assets {
				assets[name] = newTestAsset(gl, name)
			}
			for name, deps := range tc.assets {
				dependencies := make([]Asset, len(deps))
				for i, d := range deps {
					dependencies[i] = assets[d]
				}
				assets[name].dependencies = dependencies
			}
			for _, assetName := range tc.existingAssets {
				asset := assets[assetName]
				store.assets[asset] = nil
			}
			_, err := store.Fetch(assets[tc.target])
			assert.NoError(t, err, "error fetching asset")
			assert.EqualValues(t, tc.expectedGenerationLog, gl.log)
		})
	}
}
