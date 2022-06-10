package manifests

import (
	"fmt"
	"os"
	"path/filepath"

	hivev1 "github.com/openshift/hive/apis/hive/v1"
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"

	"github.com/openshift/installer/pkg/asset"
)

var (
	clusterDeploymentFilename = filepath.Join(clusterManifestDir, "cluster-deployment.yaml")
)

// ClusterDeployment generates the cluster-deployment.yaml file.
type ClusterDeployment struct {
	asset.DefaultFileWriter

	Config *hivev1.ClusterDeployment
}

var _ asset.WritableAsset = (*ClusterDeployment)(nil)

// Name returns a human friendly name for the asset.
func (*ClusterDeployment) Name() string {
	return "ClusterDeployment Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*ClusterDeployment) Dependencies() []asset.Asset {
	return []asset.Asset{
		// TODO: add dependency to InstallConfig when generation is implemented
	}
}

// Generate generates the ClusterDeployment manifest.
func (i *ClusterDeployment) Generate(dependencies asset.Parents) error {
	// TODO: generate from install-config
	return nil
}

// Load returns ClusterDeployment asset from the disk.
func (i *ClusterDeployment) Load(f asset.FileFetcher) (bool, error) {

	file, err := f.FetchByName(clusterDeploymentFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrap(err, fmt.Sprintf("failed to load %s file", clusterDeploymentFilename))
	}

	config := &hivev1.ClusterDeployment{}
	if err := yaml.UnmarshalStrict(file.Data, config); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", clusterDeploymentFilename)
	}

	i.File, i.Config = file, config
	return true, nil
}
