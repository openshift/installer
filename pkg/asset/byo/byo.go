package byo

import (
	"github.com/openshift/installer/pkg/asset"
	"path/filepath"
)

const (
	PluginsDir = "plugins"
)

type Deployment struct {
	FileList []*asset.File
}

var _ asset.Asset = (*Deployment)(nil)

func (d *Deployment) Name() string {
	return "BYO Deployment"
}

func (d *Deployment) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

func (d *Deployment) Generate(asset.Parents) error {
	return nil
}

func (d *Deployment) Files() []*asset.File {
	return d.FileList
}

// Load returns the deployment asset from disk.
func (d *Deployment) Load(f asset.FileFetcher) (bool, error) {
	var fileList []*asset.File

	// get the main.tf
	mainFile, err := f.FetchByName("main.tf")
	if err != nil {
		return false, err
	}
	fileList = append(fileList, mainFile)

	// get the variables-*.tf
	variablesFile, err := f.FetchByPattern("variables-*.tf")
	if err != nil {
		return false, err
	}
	fileList = append(fileList, variablesFile...)

	// get all modules
	modulesFileList, err := f.FetchByPattern(filepath.Join("**/*.tf"))
	if err != nil {
		return false, err
	}
	fileList = append(fileList, modulesFileList...)

	// get plugins
	pluginsFileList, _ := f.FetchByPattern(filepath.Join(PluginsDir, "**/*"))
	fileList = append(fileList, pluginsFileList...)

	d.FileList = fileList

	asset.SortFiles(d.FileList)

	return true, nil
}
