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
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/aws"
	"github.com/openshift/installer/pkg/asset/cluster/azure"
	"github.com/openshift/installer/pkg/asset/cluster/openstack"
	"github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/asset/installconfig"
	awsconfig "github.com/openshift/installer/pkg/asset/installconfig/aws"
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
		&tfvars.TerraformVariables{},
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
	terraformVariables := &tfvars.TerraformVariables{}
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
		// TODO(vincepri): The context should be passed down from the caller,
		// although today the Asset interface doesn't allow it, refactor once it does.
		ctx, cancel := context.WithCancel(signals.SetupSignalHandler())
		go func() {
			<-ctx.Done()
			cancel()
			clusterapi.System().Teardown()
		}()

		return c.provisionWithClusterAPI(ctx, parents, installConfig, clusterID)
	}

	// Otherwise, use the normal path.
	return c.provision(installConfig, clusterID, terraformVariables, parents)
}

func (c *Cluster) provision(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, terraformVariables *tfvars.TerraformVariables, parents asset.Parents) error {
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

	provider, err := infra.ProviderForPlatform(platform, installConfig.Config.EnabledFeatureGates())
	if err != nil {
		return fmt.Errorf("error getting infrastructure provider: %w", err)
	}
	files, err := provider.Provision(InstallDir, parents)
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

	// Run the post-provisioning steps for the platform we're on.
	// TODO(vincepri): The following should probably be in a separate package with a clear
	// interface and multiple hooks at different stages of the cluster lifecycle.
	switch installConfig.Config.Platform.Name() {
	case typesaws.Name:
		ssn, err := installConfig.AWS.Session(context.TODO())
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}
		client := awsconfig.NewClient(ssn)
		r53cfg := awsconfig.GetR53ClientCfg(ssn, "")
		err = client.CreateOrUpdateRecord(installConfig.Config, cluster.Spec.ControlPlaneEndpoint.Host, r53cfg)
		if err != nil {
			return fmt.Errorf("failed to create route53 records: %w", err)
		}
		logrus.Infof("Created Route53 records to control plane load balancer.")
	default:
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
