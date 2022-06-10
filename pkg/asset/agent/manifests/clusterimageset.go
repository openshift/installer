package manifests

import (
	"fmt"
	"os"
	"path/filepath"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

var (
	clusterImageSetFilename = filepath.Join(clusterManifestDir, "cluster-image-set.yaml")
)

// ClusterImageSet generates the cluster-image-set.yaml file.
type ClusterImageSet struct {
	asset.DefaultFileWriter

	Config *hivev1.ClusterImageSet
}

var _ asset.WritableAsset = (*ClusterImageSet)(nil)

// Name returns a human friendly name for the asset.
func (*ClusterImageSet) Name() string {
	return "ClusterImageSet Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*ClusterImageSet) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate generates the ClusterImageSet manifest.
func (a *ClusterImageSet) Generate(dependencies asset.Parents) error {
	return nil
}

// Load returns ClusterImageSet asset from the disk.
func (a *ClusterImageSet) Load(f asset.FileFetcher) (bool, error) {

	clusterImageSetFile, err := f.FetchByName(clusterImageSetFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", clusterImageSetFilename))
	}

	a.File = clusterImageSetFile

	clusterImageSet := &hivev1.ClusterImageSet{}
	if err := yaml.UnmarshalStrict(clusterImageSetFile.Data, clusterImageSet); err != nil {
		err = errors.Wrapf(err, "failed to unmarshal %s", clusterImageSetFilename)
		return false, err
	}
	a.Config = clusterImageSet

	return true, nil
}
