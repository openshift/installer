package mirror

import (
	"github.com/openshift/installer/pkg/asset"
)

const (
	mirrorConfigDir = "mirror"
)

var (
	_ asset.WritableAsset = (*AgentMirror)(nil)
)

// AgentMirror generates all the files required for disconnected installations.
type AgentMirror struct {
	asset.DefaultFileListWriter
}

// Name returns a human friendly name.
func (m *AgentMirror) Name() string {
	return "Agent Mirror Files"
}

// Dependencies returns all of the dependencies directly needed the asset.
func (m *AgentMirror) Dependencies() []asset.Asset {
	return []asset.Asset{
		&RegistriesConf{},
		&CaBundle{},
	}
}

// Generate generates the respective mirror files.
func (m *AgentMirror) Generate(dependencies asset.Parents) error {
	for _, a := range []asset.WritableAsset{
		&RegistriesConf{},
		&CaBundle{},
	} {
		dependencies.Get(a)
		m.FileList = append(m.FileList, a.Files()...)
	}

	asset.SortFiles(m.FileList)

	return nil
}

// Load returns the mirror files from disk.
func (m *AgentMirror) Load(f asset.FileFetcher) (bool, error) {

	return false, nil
}
