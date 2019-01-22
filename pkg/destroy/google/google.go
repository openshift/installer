package google

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/destroy"
	"github.com/openshift/installer/pkg/terraform"
	"github.com/openshift/installer/pkg/types"
)

// ClusterUninstaller holds the various options for the cluster we want to delete.
type ClusterUninstaller struct {
	Metadata *types.ClusterMetadata
	RootDir  string
	Logger   logrus.FieldLogger
}

// New returns an OpenStack destroyer from ClusterMetadata.
func New(logger logrus.FieldLogger, metadata *types.ClusterMetadata, rootDir string) (destroy.Destroyer, error) {
	return &ClusterUninstaller{
		Metadata: metadata,
		RootDir:  rootDir,
		Logger:   logger,
	}, nil
}

// Destroy uses Terraform to remove bootstrap resources.
func (u *ClusterUninstaller) Run() (err error) {
	metadata := u.Metadata
	dir := u.RootDir

	platform := metadata.Platform()
	if platform == "" {
		return errors.New("no platform configured in metadata")
	}

	copyNames := []string{terraform.StateFileName, terraform.VarFileName}

	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory for Terraform execution")
	}
	defer os.RemoveAll(tempDir)

	for _, filename := range copyNames {
		err = copy(filepath.Join(dir, filename), filepath.Join(tempDir, filename))
		if err != nil {
			return errors.Wrapf(err, "failed to copy %s to the temporary directory", filename)
		}
	}

	err = terraform.Destroy(tempDir, platform)
	if err != nil {
		return errors.Wrap(err, "Terraform destroy")
	}

	tempStateFilePath := filepath.Join(dir, terraform.StateFileName+".new")
	err = copy(filepath.Join(tempDir, terraform.StateFileName), tempStateFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to copy %s from the temporary directory", terraform.StateFileName)
	}
	return os.Rename(tempStateFilePath, filepath.Join(dir, terraform.StateFileName))
}

func copy(from string, to string) error {
	data, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(to, data, 0666)
}
