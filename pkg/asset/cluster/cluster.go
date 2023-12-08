package cluster

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utilkubeconfig "sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/aws"
	"github.com/openshift/installer/pkg/asset/cluster/azure"
	"github.com/openshift/installer/pkg/asset/cluster/openstack"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	capimanifests "github.com/openshift/installer/pkg/asset/manifests/clusterapi"
	"github.com/openshift/installer/pkg/asset/password"
	"github.com/openshift/installer/pkg/asset/quota"
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

// ResourceProvisioner provides hooks for creating additional resources during the
// provisioning lifecycle.
type ResourceProvisioner interface {
	// PreProvision is called before provisioning using CAPI controllers has begun.
	// and should be used to create dependencies needed for CAPI provisioning,
	// such as IAM roles or policies.
	PreProvision(clusterID string) error

	// ValidControlPlaneEndpoint is called once cluster.Spec.ControlPlaneEndpoint.IsValid()
	// returns true, typically after load balancers have been provisioned. It can be used
	// to create DNS records.
	ValidControlPlaneEndpoint(*clusterv1.Cluster) error
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
		// PlatformCredsCheck, PlatformPermsCheck and PlatformProvisionCheck
		// perform validations & check perms required to provision infrastructure.
		// We do not actually use them in this asset directly, hence
		// they are put in the dependencies but not fetched in Generate.
		&installconfig.PlatformCredsCheck{},
		&installconfig.PlatformPermsCheck{},
		&installconfig.PlatformProvisionCheck{},
		&quota.PlatformQuotaCheck{},
		&TerraformVariables{},
		&password.KubeadminPassword{},
		&capimanifests.Cluster{},
		&kubeconfig.AdminClient{},
	}
}

// Generate launches the cluster and generates the terraform state file on disk.
func (c *Cluster) Generate(parents asset.Parents) (err error) {
	if InstallDir == "" {
		logrus.Fatalf("InstallDir has not been set for the %q asset", c.Name())
	}

	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	terraformVariables := &TerraformVariables{}
	parents.Get(clusterID, installConfig, terraformVariables)

	if fs := installConfig.Config.FeatureSet; strings.HasSuffix(string(fs), "NoUpgrade") {
		logrus.Warnf("FeatureSet %q is enabled. This FeatureSet does not allow upgrades and may affect the supportability of the cluster.", fs)
	}

	if installConfig.Config.Platform.None != nil {
		return errors.New("cluster cannot be created with platform set to 'none'")
	}

	if installConfig.Config.BootstrapInPlace != nil {
		return errors.New("cluster cannot be created with bootstrapInPlace set")
	}

	// Check if we're using Cluster API.
	if capiutils.IsEnabled(installConfig) {
		return c.provisionWithClusterAPI(context.TODO(), parents, installConfig, clusterID)
	}

	// Otherwise, use the normal path.
	return c.provision(installConfig, clusterID, terraformVariables)
}

func (c *Cluster) provision(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, terraformVariables *TerraformVariables) error {
	platform := installConfig.Config.Platform.Name()

	if azure := installConfig.Config.Platform.Azure; azure != nil && azure.CloudName == typesazure.StackCloud {
		platform = typesazure.StackTerraformName
	}

	logrus.Infof("Creating infrastructure resources...")
	switch platform {
	case typesaws.Name:
		if err := aws.PreTerraform(context.TODO(), clusterID.InfraID, installConfig); err != nil {
			return err
		}
	case typesazure.Name, typesazure.StackTerraformName:
		if err := azure.PreTerraform(context.TODO(), clusterID.InfraID, installConfig); err != nil {
			return err
		}
	case typesopenstack.Name:
		if err := openstack.PreTerraform(); err != nil {
			return err
		}
	}

	tfvarsFiles := []*asset.File{}
	for _, file := range terraformVariables.Files() {
		tfvarsFiles = append(tfvarsFiles, file)
	}

	provider, err := infra.ProviderForPlatform(platform, installConfig.Config.EnabledFeatureGates())
	if err != nil {
		return fmt.Errorf("error getting infrastructure provider: %w", err)
	}
	files, err := provider.Provision(InstallDir, tfvarsFiles)
	if files != nil {
		c.FileList = append(c.FileList, files...) // append state files even in case of failure
	}
	if err != nil {
		return fmt.Errorf("%s: %w", asset.ClusterCreationError, err)
	}

	return nil
}

func (c *Cluster) provisionWithClusterAPI(ctx context.Context, parents asset.Parents, installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) error {
	capiManifests := &capimanifests.Cluster{}
	clusterKubeconfigAsset := &kubeconfig.AdminClient{}
	parents.Get(
		capiManifests,
		clusterKubeconfigAsset,
	)

	// supplementalProvisioner creates resources not provided in CAPI provisioning.
	supplementalProvisioner := initProvisioner(installConfig)

	if err := supplementalProvisioner.PreProvision(clusterID.InfraID); err != nil {
		return fmt.Errorf("failed to pre-provision resources: %w", err)
	}

	// Only need the objects--not the files.
	manifests := []client.Object{}
	for _, m := range capiManifests.RuntimeFiles() {
		manifests = append(manifests, m.Object)
	}

	// Run the CAPI system.
	capiSystem := clusterapi.System()
	if err := capiSystem.Run(ctx, installConfig); err != nil {
		return fmt.Errorf("failed to run cluster api system: %w", err)
	}

	// Grab the client.
	cl := capiSystem.Client()

	// Create all the manifests and store them.
	for _, m := range manifests {
		m.SetNamespace(capiutils.Namespace)
		if err := cl.Create(context.Background(), m); err != nil {
			return fmt.Errorf("failed to create manifest: %w", err)
		}
		logrus.Infof("Created manifest %+T, namespace=%s name=%s", m, m.GetNamespace(), m.GetName())
	}

	// Pass cluster kubeconfig and store it in; this is usually the role of a bootstrap provider.
	{
		key := client.ObjectKey{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		}
		cluster := &clusterv1.Cluster{}
		if err := cl.Get(context.Background(), key, cluster); err != nil {
			return err
		}
		// Create the secret.
		clusterKubeconfig := clusterKubeconfigAsset.Files()[0].Data
		secret := utilkubeconfig.GenerateSecret(cluster, clusterKubeconfig)
		if err := cl.Create(context.Background(), secret); err != nil {
			return err
		}
	}

	// Wait for the load balancer to be ready by checking the control plane endpoint
	// on the cluster object.
	var cluster *clusterv1.Cluster
	{
		if err := wait.ExponentialBackoff(wait.Backoff{
			Duration: time.Second * 10,
			Factor:   float64(1.5),
			Steps:    32,
		}, func() (bool, error) {
			c := &clusterv1.Cluster{}
			if err := cl.Get(context.Background(), client.ObjectKey{
				Name:      clusterID.InfraID,
				Namespace: capiutils.Namespace,
			}, c); err != nil {
				if apierrors.IsNotFound(err) {
					return false, nil
				}
				return false, err
			}
			cluster = c
			return cluster.Spec.ControlPlaneEndpoint.IsValid(), nil
		}); err != nil {
			return err
		}
		if cluster == nil {
			return errors.New("error occurred during load balancer ready check")
		}
		if cluster.Spec.ControlPlaneEndpoint.Host == "" {
			return errors.New("control plane endpoint is not set")
		}
	}

	if err := supplementalProvisioner.ValidControlPlaneEndpoint(cluster); err != nil {
		return fmt.Errorf("failed to create supplemental resources for valid control plane endpoint: %w", err)
	}

	// For each manifest we created, retrieve it and store it in the asset.
	for _, m := range manifests {
		key := client.ObjectKey{
			Name:      m.GetName(),
			Namespace: m.GetNamespace(),
		}
		if err := cl.Get(context.Background(), key, m); err != nil {
			return fmt.Errorf("failed to get manifest: %w", err)
		}

		gvk, err := cl.GroupVersionKindFor(m)
		if err != nil {
			return fmt.Errorf("failed to get GVK for manifest: %w", err)
		}
		fileName := fmt.Sprintf("%s-%s-%s.yaml", gvk.Kind, m.GetNamespace(), m.GetName())
		objData, err := yaml.Marshal(m)
		if err != nil {
			errMsg := fmt.Sprintf("failed to create infrastructure manifest %s from InstallConfig", fileName)
			return errors.Wrapf(err, errMsg)
		}
		c.FileList = append(c.FileList, &asset.File{
			Filename: fileName,
			Data:     objData,
		})
	}

	logrus.Infof("Cluster API resources have been created. Waiting for cluster to become ready...")
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
		return true, errors.Errorf("terraform state files alread exist.  There may already be a running cluster")
	}

	return false, nil
}

// defaultProvisioner does nothing and can be used if a cloud platform does not need
// to provision additional resources.
type defaultProvisioner struct{}

func (d defaultProvisioner) PreProvision(clusterID string) error                        { return nil }
func (d defaultProvisioner) ValidControlPlaneEndpoint(cluster *clusterv1.Cluster) error { return nil }

func initProvisioner(ic *installconfig.InstallConfig) ResourceProvisioner {
	switch ic.Config.Platform.Name() {
	case typesaws.Name:
		return aws.InitAWSProvisioner(ic)
	default:
		return defaultProvisioner{}
	}
}
