package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utilkubeconfig "sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	capimanifests "github.com/openshift/installer/pkg/asset/manifests/clusterapi"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/clusterapi"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/types"
)

// Ensure that clusterapi.InfraProvider implements
// the infrastructure.Provider interface, which is the
// interface the installer uses to call this provider.
var _ infrastructure.Provider = (*InfraProvider)(nil)

// InfraProvider implements common Cluster API logic and
// contains the platform CAPI provider, which is called
// in the lifecycle defined by the Provider interface.
type InfraProvider struct {
	impl Provider
}

// InitializeProvider returns a ClusterAPI provider implementation
// for a specific cloud platform.
func InitializeProvider(platform Provider) infrastructure.Provider {
	return &InfraProvider{impl: platform}
}

// Provision creates cluster resources by applying CAPI manifests to a locally running control plane.
//
//nolint:gocyclo
func (i *InfraProvider) Provision(dir string, parents asset.Parents) ([]*asset.File, error) {
	manifestsAsset := &manifests.Manifests{}
	capiManifestsAsset := &capimanifests.Cluster{}
	capiMachinesAsset := &machines.ClusterAPI{}
	clusterKubeconfigAsset := &kubeconfig.AdminClient{}
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	bootstrapIgnAsset := &bootstrap.Bootstrap{}
	masterIgnAsset := &machine.Master{}
	parents.Get(
		manifestsAsset,
		capiManifestsAsset,
		clusterKubeconfigAsset,
		clusterID,
		installConfig,
		rhcosImage,
		bootstrapIgnAsset,
		masterIgnAsset,
		capiMachinesAsset,
	)

	fileList := []*asset.File{}

	// Collect cluster and non-machine-related infra manifests
	// to be applied during the initial stage.
	infraManifests := []client.Object{}
	for _, m := range capiManifestsAsset.RuntimeFiles() {
		infraManifests = append(infraManifests, m.Object)
	}

	// Machine manifests will be applied after the infra
	// manifests and subsequent hooks.
	machineManifests := []client.Object{}
	for _, m := range capiMachinesAsset.RuntimeFiles() {
		machineManifests = append(machineManifests, m.Object)
	}

	// TODO(vincepri): The context should be passed down from the caller,
	// although today the Asset interface doesn't allow it, refactor once it does.
	ctx, cancel := context.WithCancel(signals.SetupSignalHandler())
	go func() {
		<-ctx.Done()
		cancel()
		clusterapi.System().Teardown()
	}()

	if p, ok := i.impl.(PreProvider); ok {
		preProvisionInput := PreProvisionInput{
			InfraID:          clusterID.InfraID,
			InstallConfig:    installConfig,
			RhcosImage:       rhcosImage,
			ManifestsAsset:   manifestsAsset,
			MachineManifests: machineManifests,
		}

		if err := p.PreProvision(ctx, preProvisionInput); err != nil {
			return fileList, fmt.Errorf("failed during pre-provisioning: %w", err)
		}
	} else {
		logrus.Debugf("No pre-provisioning requirements for the %s provider", i.impl.Name())
	}

	// Run the CAPI system.
	capiSystem := clusterapi.System()
	if err := capiSystem.Run(ctx, installConfig); err != nil {
		return fileList, fmt.Errorf("failed to run cluster api system: %w", err)
	}

	// Grab the client.
	cl := capiSystem.Client()

	// Create the infra manifests.
	for _, m := range infraManifests {
		m.SetNamespace(capiutils.Namespace)
		if err := cl.Create(ctx, m); err != nil {
			return fileList, fmt.Errorf("failed to create infrastructure manifest: %w", err)
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
		if err := cl.Get(ctx, key, cluster); err != nil {
			return fileList, err
		}
		// Create the secret.
		clusterKubeconfig := clusterKubeconfigAsset.Files()[0].Data
		secret := utilkubeconfig.GenerateSecret(cluster, clusterKubeconfig)
		if err := cl.Create(ctx, secret); err != nil {
			return fileList, err
		}
	}

	// Wait for successful provisioning by checking the InfrastructureReady
	// status on the cluster object.
	var cluster *clusterv1.Cluster
	{
		if err := wait.ExponentialBackoffWithContext(ctx, wait.Backoff{
			Duration: time.Second * 10,
			Factor:   float64(1.5),
			Steps:    32,
		}, func(ctx context.Context) (bool, error) {
			c := &clusterv1.Cluster{}
			if err := cl.Get(ctx, client.ObjectKey{
				Name:      clusterID.InfraID,
				Namespace: capiutils.Namespace,
			}, c); err != nil {
				if apierrors.IsNotFound(err) {
					return false, nil
				}
				return false, err
			}
			cluster = c
			return cluster.Status.InfrastructureReady, nil
		}); err != nil {
			return fileList, err
		}
		if cluster == nil {
			return fileList, fmt.Errorf("error occurred during load balancer ready check")
		}
		if cluster.Spec.ControlPlaneEndpoint.Host == "" {
			return fileList, fmt.Errorf("control plane endpoint is not set")
		}
	}

	if p, ok := i.impl.(InfraReadyProvider); ok {
		infraReadyInput := InfraReadyInput{
			Client:        cl,
			InstallConfig: installConfig,
			InfraID:       clusterID.InfraID,
		}

		if err := p.InfraReady(ctx, infraReadyInput); err != nil {
			return fileList, fmt.Errorf("failed provisioning resources after infrastructure ready: %w", err)
		}
	} else {
		logrus.Debugf("No infrastructure ready requirements for the %s provider", i.impl.Name())
	}

	bootstrapIgnData, err := injectInstallInfo(bootstrapIgnAsset.Files()[0].Data)
	if err != nil {
		return nil, fmt.Errorf("unable to inject installation info: %w", err)
	}

	// The cloud-platform may need to override the default
	// bootstrap ignition behavior.
	if p, ok := i.impl.(IgnitionProvider); ok {
		ignInput := IgnitionInput{
			Client:           cl,
			BootstrapIgnData: bootstrapIgnData,
			InfraID:          clusterID.InfraID,
			InstallConfig:    installConfig,
		}

		if bootstrapIgnData, err = p.Ignition(ctx, ignInput); err != nil {
			return fileList, fmt.Errorf("failed preparing ignition data: %w", err)
		}
	} else {
		logrus.Debugf("No Ignition requirements for the %s provider", i.impl.Name())
	}
	bootstrapIgnSecret := IgnitionSecret(bootstrapIgnData, clusterID.InfraID, "bootstrap")
	masterIgnSecret := IgnitionSecret(masterIgnAsset.Files()[0].Data, clusterID.InfraID, "master")
	machineManifests = append(machineManifests, bootstrapIgnSecret, masterIgnSecret)

	// Create the machine manifests.
	for _, m := range machineManifests {
		m.SetNamespace(capiutils.Namespace)
		if err := cl.Create(ctx, m); err != nil {
			return fileList, fmt.Errorf("failed to create control-plane manifest: %w", err)
		}
		logrus.Infof("Created manifest %+T, namespace=%s name=%s", m, m.GetNamespace(), m.GetName())
	}

	{
		masterCount := int64(1)
		if reps := installConfig.Config.ControlPlane.Replicas; reps != nil {
			masterCount = *reps
		}

		logrus.Debugf("Waiting for machines to provision")
		if err := wait.ExponentialBackoffWithContext(ctx, wait.Backoff{
			Duration: time.Second * 10,
			Factor:   float64(1.5),
			Steps:    32,
		}, func(ctx context.Context) (bool, error) {
			for i := int64(0); i < masterCount; i++ {
				machine := &clusterv1.Machine{}
				if err := cl.Get(ctx, client.ObjectKey{
					Name:      fmt.Sprintf("%s-%s-%d", clusterID.InfraID, "master", i),
					Namespace: capiutils.Namespace,
				}, machine); err != nil {
					if apierrors.IsNotFound(err) {
						logrus.Debugf("Not found")
						return false, nil
					}
					return false, err
				}
				if machine.Status.Phase != string(clusterv1.MachinePhaseProvisioned) &&
					machine.Status.Phase != string(clusterv1.MachinePhaseRunning) {
					return false, nil
				} else if machine.Status.Phase == string(clusterv1.MachinePhaseFailed) {
					return false, fmt.Errorf("machine %s failed to provision: %q", machine.Name, *machine.Status.FailureMessage)
				}
				logrus.Debugf("Machine %s is ready. Phase: %s", machine.Name, machine.Status.Phase)
			}
			return true, nil
		}); err != nil {
			return fileList, err
		}
	}

	if p, ok := i.impl.(PostProvider); ok {
		postMachineInput := PostProvisionInput{
			Client:        cl,
			InstallConfig: installConfig,
			InfraID:       clusterID.InfraID,
		}

		if err = p.PostProvision(ctx, postMachineInput); err != nil {
			return fileList, fmt.Errorf("failed during post-machine creation hook: %w", err)
		}
	}

	// For each manifest we created, retrieve it and store it in the asset.
	manifests := []client.Object{}
	manifests = append(manifests, infraManifests...)
	manifests = append(manifests, machineManifests...)
	for _, m := range manifests {
		key := client.ObjectKey{
			Name:      m.GetName(),
			Namespace: m.GetNamespace(),
		}
		if err := cl.Get(ctx, key, m); err != nil {
			return fileList, fmt.Errorf("failed to get manifest: %w", err)
		}

		gvk, err := cl.GroupVersionKindFor(m)
		if err != nil {
			return fileList, fmt.Errorf("failed to get GVK for manifest: %w", err)
		}
		fileName := fmt.Sprintf("%s-%s-%s.yaml", gvk.Kind, m.GetNamespace(), m.GetName())
		objData, err := yaml.Marshal(m)
		if err != nil {
			return fileList, fmt.Errorf("failed to create infrastructure manifest %s from InstallConfig: %w", fileName, err)
		}
		fileList = append(fileList, &asset.File{
			Filename: fileName,
			Data:     objData,
		})
	}

	logrus.Infof("Cluster API resources have been created. Waiting for cluster to become ready...")
	return fileList, nil
}

// DestroyBootstrap destroys the temporary bootstrap resources.
func (i *InfraProvider) DestroyBootstrap(dir string) error {
	metadata, err := metadata.Load(dir)
	if err != nil {
		return err
	}

	// TODO(padillon): start system if not running
	if sys := clusterapi.System(); sys.State() == clusterapi.SystemStateRunning {
		machineName := capiutils.GenerateBoostrapMachineName(metadata.InfraID)
		machineNamespace := capiutils.Namespace
		if err := sys.Client().Delete(context.TODO(), &clusterv1.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      machineName,
				Namespace: machineNamespace,
			},
		}); err != nil {
			return fmt.Errorf("failed to delete bootstrap machine: %w", err)
		}

		machineDeletionTimeout := 2 * time.Minute
		logrus.Infof("Waiting up to %v for bootstrap machine deletion %s/%s...", machineDeletionTimeout, machineNamespace, machineName)
		machineContext, cancel := context.WithTimeout(context.TODO(), machineDeletionTimeout)
		wait.Until(func() {
			err := sys.Client().Get(context.TODO(), client.ObjectKey{
				Name:      machineName,
				Namespace: machineNamespace,
			}, &clusterv1.Machine{})
			if err != nil {
				if apierrors.IsNotFound(err) {
					logrus.Debugf("Machine deleted: %s", machineName)
					cancel()
				} else {
					logrus.Debugf("Error when deleting bootstrap machine: %s", err)
				}
			}
		}, 2*time.Second, machineContext.Done())

		err = machineContext.Err()
		if err != nil && !errors.Is(err, context.Canceled) {
			logrus.Infof("Timeout deleting bootstrap machine: %s", err)
		}
	}
	logrus.Infof("Finished destroying bootstrap resources")
	clusterapi.System().Teardown()

	return nil
}

// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
func (i *InfraProvider) ExtractHostAddresses(dir string, config *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	return nil
}

// IgnitionSecret provides the basic formatting for creating the
// ignition secret.
func IgnitionSecret(ign []byte, infraID, role string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", infraID, role),
			Namespace: capiutils.Namespace,
			Labels: map[string]string{
				"cluster.x-k8s.io/cluster-name": infraID,
			},
		},
		Data: map[string][]byte{
			"format": []byte("ignition"),
			"value":  ign,
		},
	}
	secret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
	return secret
}
