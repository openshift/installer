package openshift

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/templates/content"
)

const (
	netopCRDfilename = "cluster-networkconfig-crd.yaml"
)

var _ asset.WritableAsset = (*NetworkCRDs)(nil)

// NetworkCRDs is the custom resource definitions for the network operator types:
// - NetworkConfig.networkoperator.openshift.io
type NetworkCRDs struct {
	asset.DefaultFileListWriter
}

// Dependencies returns all of the dependencies directly needed by the asset
func (t *NetworkCRDs) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Name returns the human-friendly name of the asset.
func (t *NetworkCRDs) Name() string {
	return "Network CRDs"
}

// Generate generates the actual files by this asset
func (t *NetworkCRDs) Generate(parents asset.Parents) error {
	data, err := content.GetOpenshiftTemplate(netopCRDfilename)
	if err != nil {
		return err
	}
	t.FileList = append(t.FileList, &asset.File{
		Filename: filepath.Join(content.TemplateDir, netopCRDfilename),
		Data:     []byte(data),
	})
	return nil
}

// Load returns the asset from disk.
func (t *NetworkCRDs) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(filepath.Join(content.TemplateDir, netopCRDfilename))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	t.FileList = append(t.FileList, file)

	return true, nil
}
