package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/capi"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/password"
)

// CAPICluster uses a CAPI Control Plane to generate
// cluster infrastructure.
type CAPICluster struct {
	FileList []*asset.File
}

var _ asset.WritableAsset = (*CAPICluster)(nil)

// Name returns the human-friendly name of the asset.
func (c *CAPICluster) Name() string {
	return "CAPI Cluster"
}

// Dependencies returns the direct dependency for launching
// the capi cluster.
func (c *CAPICluster) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.ClusterID{},
		&installconfig.InstallConfig{},
		// PlatformCredsCheck, PlatformPermsCheck and PlatformProvisionCheck
		// perform validations & check perms required to provision infrastructure.
		// We do not actually use them in this asset directly, hence
		// they are put in the dependencies but not fetched in Generate.

		// We probably want these but commenting them out for now
		// just to speed up development iterations.
		// &installconfig.PlatformCredsCheck{},
		// &installconfig.PlatformPermsCheck{},
		// &installconfig.PlatformProvisionCheck{},
		// &quota.PlatformQuotaCheck{},
		&password.KubeadminPassword{},
		&capi.CAPIControlPlane{},
		&manifests.ClusterAPI{},
		&machines.CAPIMachine{},
		&bootstrap.Bootstrap{},
		&machine.Master{},
	}
}

// Generate launches the cluster.
func (c *CAPICluster) Generate(parents asset.Parents) (err error) {
	if InstallDir == "" {
		logrus.Fatalf("InstallDir has not been set for the %q asset", c.Name())
	}

	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	capiControlPlane := &capi.CAPIControlPlane{}
	capiManifests := &manifests.ClusterAPI{}
	capiMachines := &machines.CAPIMachine{}
	bootstrapIgnAsset := &bootstrap.Bootstrap{}
	masterIgnAsset := &machine.Master{}

	parents.Get(
		clusterID,
		installConfig,
		capiControlPlane,
		capiManifests,
		capiMachines,
		bootstrapIgnAsset,
		masterIgnAsset,
	)

	// Only need the objects--not the files.
	manifests := []client.Object{}
	for _, m := range capiManifests.Manifests {
		manifests = append(manifests, m.Object)
	}
	manifests = append(manifests, capiMachines.Machines...)

	if fs := installConfig.Config.FeatureSet; strings.HasSuffix(string(fs), "NoUpgrade") {
		logrus.Warnf("FeatureSet %q is enabled. This FeatureSet does not allow upgrades and may affect the supportability of the cluster.", fs)
	}

	if installConfig.Config.Platform.None != nil {
		return errors.New("cluster cannot be created with platform set to 'none'")
	}

	if installConfig.Config.BootstrapInPlace != nil {
		return errors.New("cluster cannot be created with bootstrapInPlace set")
	}

	// Create a new client to interact with the cluster.
	cl, err := client.New(capiControlPlane.LocalCP.Cfg, client.Options{
		Scheme: capiControlPlane.LocalCP.Env.Scheme,
	})
	if err != nil {
		return err
	}

	// Create the namespace for the cluster.
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "openshift-cluster-api-guests",
		},
	}
	if err := cl.Create(context.Background(), ns); err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	// Gather the ignition files, and store them in a secret.
	{
		masterIgn := string(masterIgnAsset.Files()[0].Data)
		bootstrapIgn, err := injectInstallInfo(bootstrapIgnAsset.Files()[0].Data)
		if err != nil {
			return errors.Wrap(err, "unable to inject installation info")
		}
		manifests = append(manifests,
			&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%s", clusterID.InfraID, "master"),
					Namespace: ns.Name,
				},
				Data: map[string][]byte{
					"format": []byte("ignition"),
					"value":  []byte(masterIgn),
				},
			},
			&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%s", clusterID.InfraID, "bootstrap"),
					Namespace: ns.Name,
				},
				Data: map[string][]byte{
					"format": []byte("ignition"),
					"value":  []byte(bootstrapIgn),
				},
			},
		)
	}

	for _, m := range manifests {
		m.SetNamespace(ns.Name)
		if err := cl.Create(context.Background(), m); err != nil {
			return fmt.Errorf("failed to create manifest: %w", err)
		}
		logrus.Infof("Created manifest %+T, namespace=%s name=%s", m, m.GetNamespace(), m.GetName())
	}

	// List all namespaces in the cluster.
	namespaceList := &corev1.NamespaceList{}
	if err := cl.List(context.Background(), namespaceList); err != nil {
		return err
	}
	for _, n := range namespaceList.Items {
		spew.Dump(n.Name)
	}

	time.Sleep(20 * time.Minute)

	capiControlPlane.LocalCP.Stop()
	panic("not implemented")

	return nil
}

// Files returns the FileList generated by the asset.
func (c *CAPICluster) Files() []*asset.File {
	return c.FileList
}

// Load returns error if the tfstate file is already on-disk, because we want to
// prevent user from accidentally re-launching the cluster.
func (c *CAPICluster) Load(f asset.FileFetcher) (found bool, err error) {

	return false, nil
}
