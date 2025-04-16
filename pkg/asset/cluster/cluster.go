package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/aws"
	"github.com/openshift/installer/pkg/asset/cluster/azure"
	"github.com/openshift/installer/pkg/asset/cluster/openstack"
	"github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	capimanifests "github.com/openshift/installer/pkg/asset/manifests/clusterapi"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/quota"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/clusterapi"
	infra "github.com/openshift/installer/pkg/infrastructure/platform"
	typesaws "github.com/openshift/installer/pkg/types/aws"
	typesazure "github.com/openshift/installer/pkg/types/azure"
	typesopenstack "github.com/openshift/installer/pkg/types/openstack"
)

var (
	// InstallDir is the directory containing install assets.
	InstallDir string
)

// Cluster uses the terraform executable to launch a cluster
// with the given terraform tfvar and generated templates.
type Cluster struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*Cluster)(nil)

// Name returns the human-friendly name of the asset.
func (c *Cluster) Name() string {
	return "Cluster"
}

// Dependencies returns the direct dependency for launching
// the cluster.
func (c *Cluster) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		// PlatformCredsCheck, PlatformPermsCheck, PlatformProvisionCheck, and VCenterContexts.
		// perform validations & check perms required to provision infrastructure.
		// We do not actually use them in this asset directly, hence
		// they are put in the dependencies but not fetched in Generate.
		&installconfig.PlatformCredsCheck{},
		&installconfig.PlatformPermsCheck{},
		&installconfig.PlatformProvisionCheck{},
		new(rhcos.Image),
		&quota.PlatformQuotaCheck{},
		&tfvars.TerraformVariables{},
		&password.KubeadminPassword{},
		&manifests.Manifests{},
		&capimanifests.Cluster{},
		&kubeconfig.AdminClient{},
		&bootstrap.Bootstrap{},
		&machine.Master{},
		&machine.Worker{},
		&machines.Worker{},
		&machines.ClusterAPI{},
		new(rhcos.Image),
		&manifests.Manifests{},
		&tls.RootCA{},
	}
}

// Generate launches the cluster and generates the terraform state file on disk.
func (c *Cluster) Generate(ctx context.Context, parents asset.Parents) (err error) {
	if InstallDir == "" {
		logrus.Fatalf("InstallDir has not been set for the %q asset", c.Name())
	}

	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	terraformVariables := &tfvars.TerraformVariables{}
	parents.Get(clusterID, installConfig, terraformVariables, rhcosImage)

	if fs := installConfig.Config.FeatureSet; strings.HasSuffix(string(fs), "NoUpgrade") {
		logrus.Warnf("FeatureSet %q is enabled. This FeatureSet does not allow upgrades and may affect the supportability of the cluster.", fs)
	}

	if installConfig.Config.Platform.None != nil {
		return errors.New("cluster cannot be created with platform set to 'none'")
	}

	if installConfig.Config.BootstrapInPlace != nil {
		return errors.New("cluster cannot be created with bootstrapInPlace set")
	}

	platform := installConfig.Config.Platform.Name()

	if azure := installConfig.Config.Platform.Azure; azure != nil && azure.CloudName == typesazure.StackCloud {
		platform = typesazure.StackTerraformName
	}

	// TODO(padillon): determine whether CAPI handles tagging shared subnets, in which case we should be able
	// to encapsulate these into the terraform package.
	logrus.Infof("Creating infrastructure resources...")
	switch platform {
	case typesaws.Name:
		if err := aws.PreTerraform(ctx, clusterID.InfraID, installConfig); err != nil {
			return err
		}
	case typesazure.Name, typesazure.StackTerraformName:
		if err := azure.PreTerraform(ctx, clusterID.InfraID, installConfig); err != nil {
			return err
		}
	case typesopenstack.Name:
		var tfvarsFile *asset.File
		for _, f := range terraformVariables.Files() {
			if f.Filename == tfvars.TfPlatformVarsFileName {
				tfvarsFile = f
				break
			}
		}
		if err := openstack.PreTerraform(ctx, tfvarsFile, installConfig, clusterID, rhcosImage); err != nil {
			return err
		}
	}

	provider, err := infra.ProviderForPlatform(platform, installConfig.Config.EnabledFeatureGates())
	if err != nil {
		return fmt.Errorf("error getting infrastructure provider: %w", err)
	}
	files, err := provider.Provision(ctx, InstallDir, parents)
	if files != nil {
		c.FileList = append(c.FileList, files...) // append state files even in case of failure
	}
	if err != nil {
		return fmt.Errorf("%s: %w", asset.ClusterCreationError, err)
	}

	return nil
}

// Files returns the FileList generated by the asset.
func (c *Cluster) Files() []*asset.File {
	return c.FileList
}

// Load returns error if the tfstate file is already on-disk, because we want to
// prevent user from accidentally re-launching the cluster.
func (c *Cluster) Load(f asset.FileFetcher) (found bool, err error) {
	matches, err := filepath.Glob("terraform(.*)?.tfstate")
	if err != nil {
		return true, err
	}
	if len(matches) != 0 {
		return true, fmt.Errorf("terraform state files already exist.  There may already be a running cluster")
	}

	matches, err = filepath.Glob(filepath.Join(InstallDir, clusterapi.ArtifactsDir, "envtest.kubeconfig"))
	if err != nil {
		return true, fmt.Errorf("error checking for existence of envtest.kubeconfig: %w", err)
	}

	// Cluster-API based installs can be re-entered, but this is an experimental feature
	// that should be opted into and only used for testing and development.
	reentrant := strings.EqualFold(os.Getenv("OPENSHIFT_INSTALL_REENTRANT"), "true")

	if !reentrant && len(matches) != 0 {
		return true, fmt.Errorf("local infrastructure provisioning artifacts already exist. There may already be a running cluster")
	}
	return false, nil
}
