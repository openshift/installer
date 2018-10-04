package machine

import (
	"fmt"
	"github.com/openshift/installer/pkg/asset"
)

type testAsset struct {
	name string
}

func (a *testAsset) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

func (a *testAsset) Generate(map[string]*asset.State) (*asset.State, error) {
	return nil, nil
}

func (a *testAsset) Name() string {
	return fmt.Sprintf("Test Asset (%s)", a.name)
}

func stateWithContentsData(contentsData ...string) *asset.State {
	state := &asset.State{
		Contents: make([]asset.Content, len(contentsData)),
	}
	for i, d := range contentsData {
		state.Contents[i].Data = []byte(d)
	}
	return state
}
