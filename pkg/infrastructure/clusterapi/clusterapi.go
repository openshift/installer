package clusterapi

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563
	"sigs.k8s.io/cluster-api/util"
	utilkubeconfig "sigs.k8s.io/cluster-api/util/kubeconfig"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	"github.com/openshift/installer/pkg/asset/cluster/tfvars"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/machines"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	capimanifests "github.com/openshift/installer/pkg/asset/manifests/clusterapi"
	"github.com/openshift/installer/pkg/asset/rhcos"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/clusterapi"
	"github.com/openshift/installer/pkg/infrastructure"
	"github.com/openshift/installer/pkg/metrics/timer"
	"github.com/openshift/installer/pkg/types"
)

// Ensure that clusterapi.InfraProvider implements
// the infrastructure.Provider interface, which is the
// interface the installer uses to call this provider.
var _ infrastructure.Provider = (*InfraProvider)(nil)

const (
	preProvisionStage        = "Infrastructure Pre-provisioning"
	infrastructureStage      = "Network-infrastructure Provisioning"
	infrastructureReadyStage = "Post-network, pre-machine Provisioning"
	ignitionStage            = "Bootstrap Ignition Provisioning"
	machineStage             = "Machine Provisioning"
	postProvisionStage       = "Infrastructure Post-provisioning"
)

// InfraProvider implements common Cluster API logic and
// contains the platform CAPI provider, which is called
// in the lifecycle defined by the Provider interface.
type InfraProvider struct {
	impl Provider

	appliedManifests []client.Object
}

// InitializeProvider returns a ClusterAPI provider implementation
// for a specific cloud platform.
func InitializeProvider(platform Provider) infrastructure.Provider {
	return &InfraProvider{impl: platform}
}

// Provision creates cluster resources by applying CAPI manifests to a locally running control plane.
//
//nolint:gocyclo
func (i *InfraProvider) Provision(ctx context.Context, dir string, parents asset.Parents) (fileList []*asset.File, err error) {
	manifestsAsset := &manifests.Manifests{}
	workersAsset := &machines.Worker{}
	capiManifestsAsset := &capimanifests.Cluster{}
	capiMachinesAsset := &machines.ClusterAPI{}
	clusterKubeconfigAsset := &kubeconfig.AdminClient{}
	clusterID := &installconfig.ClusterID{}
	installConfig := &installconfig.InstallConfig{}
	rhcosImage := new(rhcos.Image)
	bootstrapIgnAsset := &bootstrap.Bootstrap{}
	masterIgnAsset := &machine.Master{}
	workerIgnAsset := &machine.Worker{}
	tfvarsAsset := &tfvars.TerraformVariables{}
	rootCA := &tls.RootCA{}
	parents.Get(
		manifestsAsset,
		workersAsset,
		capiManifestsAsset,
		clusterKubeconfigAsset,
		clusterID,
		installConfig,
		rhcosImage,
		bootstrapIgnAsset,
		masterIgnAsset,
		workerIgnAsset,
		capiMachinesAsset,
		tfvarsAsset,
		rootCA,
	)

	var capiClusters []*clusterv1.Cluster

	// Collect cluster and non-machine-related infra manifests
	// to be applied during the initial stage.
	infraManifests := []client.Object{}
	for _, m := range capiManifestsAsset.RuntimeFiles() {
		// Check for cluster definition so that we can collect the names.
		if cluster, ok := m.Object.(*clusterv1.Cluster); ok {
			capiClusters = append(capiClusters, cluster)
		}

		infraManifests = append(infraManifests, m.Object)
	}

	// Machine manifests will be applied after the infra
	// manifests and subsequent hooks.
	machineManifests := []client.Object{}
	for _, m := range capiMachinesAsset.RuntimeFiles() {
		machineManifests = append(machineManifests, m.Object)
	}

	if p, ok := i.impl.(PreProvider); ok {
		preProvisionInput := PreProvisionInput{
			InfraID:          clusterID.InfraID,
			InstallConfig:    installConfig,
			RhcosImage:       rhcosImage,
			ManifestsAsset:   manifestsAsset,
			MachineManifests: machineManifests,
			WorkersAsset:     workersAsset,
		}
		timer.StartTimer(preProvisionStage)
		if err := p.PreProvision(ctx, preProvisionInput); err != nil {
			return fileList, fmt.Errorf("failed during pre-provisioning: %w", err)
		}
		timer.StopTimer(preProvisionStage)
	} else {
		logrus.Debugf("No pre-provisioning requirements for the %s provider", i.impl.Name())
	}

	// If we're skipping bootstrap destroy, shutdown the local control plane.
	// Otherwise, we will shut it down after bootstrap destroy.
	// This has to execute as the last defer in the stack since previous defers might still need the local controlplane.
	if oi, ok := os.LookupEnv("OPENSHIFT_INSTALL_PRESERVE_BOOTSTRAP"); ok && oi != "" {
		defer func() {
			logrus.Warn("OPENSHIFT_INSTALL_PRESERVE_BOOTSTRAP is set, shutting down local control plane.")
			clusterapi.System().Teardown()
		}()
	}

	// Make sure to always return generated manifests, even if errors happened
	defer func(ctx context.Context) {
		var errs []error
		// Overriding the named return with the generated list
		fileList, errs = i.collectManifests(ctx, clusterapi.System().Client())
		// If Provision returned an error, add it to the list
		if err != nil {
			clusterapi.System().CleanEtcd()
			errs = append(errs, err)
		}
		err = utilerrors.NewAggregate(errs)
	}(ctx)

	// Run the CAPI system.
	timer.StartTimer(infrastructureStage)
	capiSystem := clusterapi.System()
	if err := capiSystem.Run(ctx); err != nil {
		return fileList, fmt.Errorf("failed to run cluster api system: %w", err)
	}

	// Grab the client.
	cl := capiSystem.Client()

	i.appliedManifests = []client.Object{}

	// Create the infra manifests.
	logrus.Info("Creating infra manifests...")
	for _, m := range infraManifests {
		m.SetNamespace(capiutils.Namespace)
		if err := cl.Create(ctx, m); err != nil {
			return fileList, fmt.Errorf("failed to create infrastructure manifest: %w", err)
		}
		i.appliedManifests = append(i.appliedManifests, m)
		logrus.Infof("Created manifest %+T, namespace=%s name=%s", m, m.GetNamespace(), m.GetName())
	}
	logrus.Info("Done creating infra manifests")

	// Pass cluster kubeconfig and store it in; this is usually the role of a bootstrap provider.
	for _, capiCluster := range capiClusters {
		logrus.Infof("Creating kubeconfig entry for capi cluster %v", capiCluster.GetName())
		key := client.ObjectKey{
			Name:      capiCluster.GetName(),
			Namespace: capiutils.Namespace,
		}
		cluster := &clusterv1.Cluster{}
		if err := cl.Get(ctx, key, cluster); err != nil {
			return fileList, err
		}
		// Create the secret.
		clusterKubeconfig := clusterKubeconfigAsset.Files()[0].Data
		secret := generateSecret(cluster, clusterKubeconfig)
		if err := cl.Create(ctx, secret); err != nil {
			return fileList, err
		}
	}

	var networkTimeout = 15 * time.Minute

	if p, ok := i.impl.(Timeouts); ok {
		networkTimeout = p.NetworkTimeout()
	}

	// Wait for successful provisioning by checking the InfrastructureReady
	// status on the cluster object.
	untilTime := time.Now().Add(networkTimeout)
	timezone, _ := untilTime.Zone()
	logrus.Infof("Waiting up to %v (until %v %s) for network infrastructure to become ready...", networkTimeout, untilTime.Format(time.Kitchen), timezone)
	{
		if err := wait.PollUntilContextTimeout(ctx, 15*time.Second, networkTimeout, true,
			func(ctx context.Context) (bool, error) {
				for _, capiCluster := range capiClusters {
					if err := cl.Get(ctx, client.ObjectKey{
						Name:      capiCluster.GetName(),
						Namespace: capiutils.Namespace,
					}, capiCluster); err != nil {
						if apierrors.IsNotFound(err) {
							return false, nil
						}
						return false, err
					}
				}
				for _, capiCluster := range capiClusters {
					if !capiCluster.Status.InfrastructureReady {
						return false, nil
					}
				}
				return true, nil
			}); err != nil {
			// Attempt to find and report falsy conditions in infra cluster if any.
			if len(capiClusters) > 0 {
				warnIfFalsyInfraConditions(ctx, capiClusters[0].Spec.InfrastructureRef, cl)
			}
			if wait.Interrupted(err) {
				return fileList, fmt.Errorf("infrastructure was not ready within %v", networkTimeout)
			}
			return fileList, fmt.Errorf("infrastructure is not ready: %w", err)
		}

		if len(capiClusters) == 0 {
			return fileList, fmt.Errorf("error occurred during load balancer ready check")
		}
		if capiClusters[0].Spec.ControlPlaneEndpoint.Host == "" {
			return fileList, fmt.Errorf("control plane endpoint is not set")
		}
	}
	timer.StopTimer(infrastructureStage)
	logrus.Info("Network infrastructure is ready")

	if p, ok := i.impl.(InfraReadyProvider); ok {
		infraReadyInput := InfraReadyInput{
			Client:        cl,
			InstallConfig: installConfig,
			InfraID:       clusterID.InfraID,
		}

		timer.StartTimer(infrastructureReadyStage)
		if err := p.InfraReady(ctx, infraReadyInput); err != nil {
			return fileList, fmt.Errorf("failed provisioning resources after infrastructure ready: %w", err)
		}
		timer.StopTimer(infrastructureReadyStage)
	} else {
		logrus.Debugf("No infrastructure ready requirements for the %s provider", i.impl.Name())
	}

	masterIgnData := masterIgnAsset.Files()[0].Data
	bootstrapIgnData, err := injectInstallInfo(bootstrapIgnAsset.Files()[0].Data)
	if err != nil {
		return fileList, fmt.Errorf("unable to inject installation info: %w", err)
	}
	workerIgnData := workerIgnAsset.Files()[0].Data
	ignitionSecrets := []*corev1.Secret{}

	// The cloud-platform may need to override the default
	// bootstrap ignition behavior.
	if p, ok := i.impl.(IgnitionProvider); ok {
		ignInput := IgnitionInput{
			Client:           cl,
			BootstrapIgnData: bootstrapIgnData,
			MasterIgnData:    masterIgnData,
			WorkerIgnData:    workerIgnData,
			InfraID:          clusterID.InfraID,
			InstallConfig:    installConfig,
			TFVarsAsset:      tfvarsAsset,
			RootCA:           rootCA,
		}

		timer.StartTimer(ignitionStage)
		if ignitionSecrets, err = p.Ignition(ctx, ignInput); err != nil {
			return fileList, fmt.Errorf("failed preparing ignition data: %w", err)
		}
		timer.StopTimer(ignitionStage)
	} else {
		logrus.Debugf("Using default ignition for the %s provider", i.impl.Name())
		bootstrapIgnSecret := IgnitionSecret(bootstrapIgnData, clusterID.InfraID, "bootstrap")
		masterIgnSecret := IgnitionSecret(masterIgnData, clusterID.InfraID, "master")
		ignitionSecrets = append(ignitionSecrets, bootstrapIgnSecret, masterIgnSecret)
	}

	for _, secret := range ignitionSecrets {
		machineManifests = append(machineManifests, secret)
	}

	// Create the machine manifests.
	timer.StartTimer(machineStage)
	machineNames := []string{}
	capiMachines := []*clusterv1.Machine{}

	for _, m := range machineManifests {
		m.SetNamespace(capiutils.Namespace)
		if err := cl.Create(ctx, m); err != nil {
			return fileList, fmt.Errorf("failed to create control-plane manifest: %w", err)
		}
		i.appliedManifests = append(i.appliedManifests, m)

		if machine, ok := m.(*clusterv1.Machine); ok {
			machineNames = append(machineNames, machine.Name)
			capiMachines = append(capiMachines, machine)
		}
		logrus.Infof("Created manifest %+T, namespace=%s name=%s", m, m.GetNamespace(), m.GetName())
	}

	var provisionTimeout = 15 * time.Minute

	if p, ok := i.impl.(Timeouts); ok {
		provisionTimeout = p.ProvisionTimeout()
	}

	{
		untilTime := time.Now().Add(provisionTimeout)
		timezone, _ := untilTime.Zone()
		reqBootstrapPubIP := installConfig.Config.Publish == types.ExternalPublishingStrategy && i.impl.PublicGatherEndpoint() == ExternalIP
		logrus.Infof("Waiting up to %v (until %v %s) for machines %v to provision...", provisionTimeout, untilTime.Format(time.Kitchen), timezone, machineNames)
		if err := wait.PollUntilContextTimeout(ctx, 15*time.Second, provisionTimeout, true,
			func(ctx context.Context) (bool, error) {
				allReady := true
				for _, machine := range capiMachines {
					if err := cl.Get(ctx, client.ObjectKey{
						Name:      machine.Name,
						Namespace: capiutils.Namespace,
					}, machine); err != nil {
						if apierrors.IsNotFound(err) {
							logrus.Debugf("Not found")
							return false, nil
						}
						return false, err
					}

					reqPubIP := reqBootstrapPubIP && machine.Name == capiutils.GenerateBoostrapMachineName(clusterID.InfraID)
					ready, err := checkMachineReady(machine, reqPubIP)
					if err != nil {
						return false, fmt.Errorf("failed waiting for machines: %w", err)
					}
					if !ready {
						allReady = false
					} else {
						logrus.Debugf("Machine %s is ready. Phase: %s", machine.Name, machine.Status.Phase)
					}
				}
				return allReady, nil
			}); err != nil {
			// Attempt to find and report falsy conditions in infra machines if any.
			for _, machine := range capiMachines {
				if machine != nil {
					warnIfFalsyInfraConditions(ctx, &machine.Spec.InfrastructureRef, cl)
				}
			}
			if wait.Interrupted(err) {
				return fileList, fmt.Errorf("%s within %v", asset.ControlPlaneCreationError, provisionTimeout)
			}
			return fileList, fmt.Errorf("%s: machines are not ready: %w", asset.ControlPlaneCreationError, err)
		}
	}
	timer.StopTimer(machineStage)
	logrus.Info("Control-plane machines are ready")

	if p, ok := i.impl.(PostProvider); ok {
		postMachineInput := PostProvisionInput{
			Client:        cl,
			InstallConfig: installConfig,
			InfraID:       clusterID.InfraID,
		}

		timer.StartTimer(postProvisionStage)
		if err = p.PostProvision(ctx, postMachineInput); err != nil {
			return fileList, fmt.Errorf("failed during post-machine creation hook: %w", err)
		}
		timer.StopTimer(postProvisionStage)
	}

	logrus.Infof("Cluster API resources have been created. Waiting for cluster to become ready...")

	return fileList, nil
}

// DestroyBootstrap destroys the temporary bootstrap resources.
func (i *InfraProvider) DestroyBootstrap(ctx context.Context, dir string) error {
	defer clusterapi.System().CleanEtcd()

	metadata, err := metadata.Load(dir)
	if err != nil {
		return err
	}

	sys := clusterapi.System()
	if sys.State() != clusterapi.SystemStateRunning {
		if err := sys.Run(ctx); err != nil {
			return fmt.Errorf("failed to run capi system: %w", err)
		}
	}

	if p, ok := i.impl.(BootstrapDestroyer); ok {
		bootstrapDestoryInput := BootstrapDestroyInput{
			Client:   sys.Client(),
			Metadata: *metadata,
		}

		if err = p.DestroyBootstrap(ctx, bootstrapDestoryInput); err != nil {
			return fmt.Errorf("failed during the destroy bootstrap hook: %w", err)
		}
	}

	machineName := capiutils.GenerateBoostrapMachineName(metadata.InfraID)
	machineNamespace := capiutils.Namespace
	if err := sys.Client().Delete(ctx, &clusterv1.Machine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineName,
			Namespace: machineNamespace,
		},
	}); err != nil {
		return fmt.Errorf("failed to delete bootstrap machine: %w", err)
	}

	machineDeletionTimeout := 5 * time.Minute
	logrus.Infof("Waiting up to %v for bootstrap machine deletion %s/%s...", machineDeletionTimeout, machineNamespace, machineName)
	cctx, cancel := context.WithTimeout(ctx, machineDeletionTimeout)
	wait.UntilWithContext(cctx, func(context.Context) {
		err := sys.Client().Get(cctx, client.ObjectKey{
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
	}, 2*time.Second)

	err = cctx.Err()
	if err != nil && !errors.Is(err, context.Canceled) {
		logrus.Warnf("Timeout deleting bootstrap machine: %s", err)
	}
	clusterapi.System().Teardown()

	if p, ok := i.impl.(PostDestroyer); ok {
		postDestroyInput := PostDestroyerInput{
			Metadata: *metadata,
		}
		if err := p.PostDestroy(ctx, postDestroyInput); err != nil {
			return fmt.Errorf("failed during post-destroy hook: %w", err)
		}
		logrus.Debugf("Finished running post-destroy hook")
	} else {
		logrus.Infof("no post-destroy requirements for the %s provider", i.impl.Name())
	}

	logrus.Infof("Finished destroying bootstrap resources")
	return nil
}

type machineManifest struct {
	Status struct {
		Addresses []clusterv1.MachineAddress `yaml:"addresses"`
	} `yaml:"status"`
}

// extractIPAddress extracts the IP address from a machine manifest file in a
// provider-agnostic way by reading only the "status" stanza, which should be
// present in all providers.
func extractIPAddress(manifestPath string) ([]string, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return []string{}, fmt.Errorf("failed to read machine manifest %s: %w", manifestPath, err)
	}
	var manifest machineManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return []string{}, fmt.Errorf("failed to unmarshal manifest %s: %w", manifestPath, err)
	}

	var externalIPAddrs []string
	var internalIPAddrs []string
	for _, addr := range manifest.Status.Addresses {
		switch addr.Type {
		case clusterv1.MachineExternalIP:
			externalIPAddrs = append(externalIPAddrs, addr.Address)
		case clusterv1.MachineInternalIP:
			internalIPAddrs = append(internalIPAddrs, addr.Address)
		default:
			continue
		}
	}

	// prioritize the external address in the front of the list
	externalIPAddrs = append(externalIPAddrs, internalIPAddrs...)

	return externalIPAddrs, nil
}

// ExtractHostAddresses extracts the IPs of the bootstrap and control plane machines.
func (i *InfraProvider) ExtractHostAddresses(dir string, config *types.InstallConfig, ha *infrastructure.HostAddresses) error {
	manifestsDir := filepath.Join(dir, clusterapi.ArtifactsDir)
	logrus.Debugf("Looking for machine manifests in %s", manifestsDir)

	addr, err := i.getBootstrapAddress(config, manifestsDir)
	if err != nil {
		return fmt.Errorf("failed to get bootstrap address: %w", err)
	}
	ha.Bootstrap = addr

	masterFiles, err := filepath.Glob(filepath.Join(manifestsDir, "Machine\\-openshift\\-cluster\\-api\\-guests\\-*\\-master\\-?.yaml"))
	if err != nil {
		return fmt.Errorf("failed to list master machine manifests: %w", err)
	}
	logrus.Debugf("master machine manifests found: %v", masterFiles)

	if replicas := int(*config.ControlPlane.Replicas); replicas != len(masterFiles) {
		logrus.Warnf("not all master manifests found: %d. Expected %d.", len(masterFiles), replicas)
	}
	for _, manifest := range masterFiles {
		addrs, err := extractIPAddress(manifest)
		if err != nil {
			// Log the error but keep parsing the remaining files
			logrus.Warnf("failed to extract IP address for %s: %v", manifest, err)
			continue
		}
		logrus.Debugf("found master address: %s", addrs)

		ha.Masters = append(ha.Masters, prioritizeIPv4(config, addrs))
	}

	return nil
}

func (i *InfraProvider) getBootstrapAddress(config *types.InstallConfig, manifestsDir string) (string, error) {
	// If the bootstrap node cannot have a public IP address, we
	// SSH through the load balancer, as is this case on Azure.
	if i.impl.PublicGatherEndpoint() == APILoadBalancer && config.Publish != types.InternalPublishingStrategy {
		return fmt.Sprintf("api.%s", config.ClusterDomain()), nil
	}

	bootstrapFiles, err := filepath.Glob(filepath.Join(manifestsDir, "Machine\\-openshift\\-cluster\\-api\\-guests\\-*\\-bootstrap.yaml"))
	if err != nil {
		return "", fmt.Errorf("failed to list bootstrap manifests: %w", err)
	}
	logrus.Debugf("bootstrap manifests found: %v", bootstrapFiles)

	if len(bootstrapFiles) != 1 {
		return "", fmt.Errorf("wrong number of bootstrap manifests found: %v. Expected exactly one", bootstrapFiles)
	}
	addrs, err := extractIPAddress(bootstrapFiles[0])
	if err != nil {
		return "", fmt.Errorf("failed to extract IP address for bootstrap: %w", err)
	}
	logrus.Debugf("found bootstrap address: %s", addrs)
	return prioritizeIPv4(config, addrs), nil
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

func (i *InfraProvider) collectManifests(ctx context.Context, cl client.Client) ([]*asset.File, []error) {
	logrus.Debug("Collecting applied cluster api manifests...")
	errorList := []error{}
	fileList := []*asset.File{}
	for _, m := range i.appliedManifests {
		key := client.ObjectKey{
			Name:      m.GetName(),
			Namespace: m.GetNamespace(),
		}
		if err := cl.Get(ctx, key, m); err != nil {
			errorList = append(errorList, fmt.Errorf("failed to get manifest %s: %w", m.GetName(), err))
			continue
		}

		gvk, err := cl.GroupVersionKindFor(m)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to get GVK for manifest %s: %w", m.GetName(), err))
			continue
		}
		fileName := filepath.Join(clusterapi.ArtifactsDir, fmt.Sprintf("%s-%s-%s.yaml", gvk.Kind, m.GetNamespace(), m.GetName()))
		objData, err := yaml.Marshal(m)
		if err != nil {
			errorList = append(errorList, fmt.Errorf("failed to marshal manifest %s: %w", fileName, err))
			continue
		}
		fileList = append(fileList, &asset.File{
			Filename: fileName,
			Data:     objData,
		})
	}
	return fileList, errorList
}

func checkMachineReady(machine *clusterv1.Machine, requirePublicIP bool) (bool, error) {
	logrus.Debugf("Checking that machine %s has provisioned...", machine.Name)
	if machine.Status.Phase != string(clusterv1.MachinePhaseProvisioned) &&
		machine.Status.Phase != string(clusterv1.MachinePhaseRunning) {
		logrus.Debugf("Machine %s has not yet provisioned: %s", machine.Name, machine.Status.Phase)
		return false, nil
	} else if machine.Status.Phase == string(clusterv1.MachinePhaseFailed) {
		//TODO: We need to update this to use non deprecated field
		msg := ptr.Deref(machine.Status.FailureMessage, "machine.Status.FailureMessage was not set") //nolint:staticcheck
		return false, fmt.Errorf("machine %s failed to provision: %s", machine.Name, msg)
	}
	logrus.Debugf("Machine %s has status: %s", machine.Name, machine.Status.Phase)
	return hasRequiredIP(machine, requirePublicIP), nil
}

func hasRequiredIP(machine *clusterv1.Machine, requirePublicIP bool) bool {
	logrus.Debugf("Checking that IP addresses are populated in the status of machine %s...", machine.Name)

	for _, addr := range machine.Status.Addresses {
		switch {
		case len(addr.Address) == 0:
			continue
		case addr.Type == clusterv1.MachineExternalIP:
			logrus.Debugf("Found external IP address: %s", addr.Address)
			return true
		case addr.Type == clusterv1.MachineInternalIP && !requirePublicIP:
			logrus.Debugf("Found internal IP address: %s", addr.Address)
			return true
		}
		logrus.Debugf("Checked IP %s: %s", addr.Type, addr.Address)
	}
	logrus.Debugf("Still waiting for machine %s to get required IPs", machine.Name)
	return false
}

// gatherInfraConditions gather conditions from CAPI infra cluster or machine
// in a provider-agnostic way from the "status.condition" field, which should be present in all providers.
// https://cluster-api.sigs.k8s.io/developer/providers/contracts/infra-cluster#infracluster-conditions
// https://cluster-api.sigs.k8s.io/developer/providers/contracts/infra-machine#inframachine-conditions
func gatherInfraConditions(ctx context.Context, objRef *corev1.ObjectReference, cl client.Client) (clusterv1.Conditions, error) {
	unstructuredObj := &unstructured.Unstructured{}
	unstructuredObj.SetGroupVersionKind(objRef.GroupVersionKind())

	if err := cl.Get(ctx, client.ObjectKey{
		Namespace: objRef.Namespace,
		Name:      objRef.Name,
	}, unstructuredObj); err != nil {
		return nil, err
	}

	// Field .status.conditions should be implemented by all providers
	// and has type clusterv1.Conditions
	infraObj := &struct {
		Status struct {
			Conditions clusterv1.Conditions `json:"conditions,omitempty"`
		} `json:"status,omitempty"`
	}{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), infraObj); err != nil {
		return nil, err
	}

	return infraObj.Status.Conditions, nil
}

// warnIfFalsyInfraConditions logs warning messages for any conditions that are not "True"
// in the infra cluster or machine status.
func warnIfFalsyInfraConditions(ctx context.Context, objRef *corev1.ObjectReference, cl client.Client) {
	apiVersion, kind := objRef.GroupVersionKind().ToAPIVersionAndKind()
	objInfo := fmt.Sprintf("apiVersion=%s, kind=%s, namespace=%s, name=%s", apiVersion, kind, objRef.Namespace, objRef.Name)

	logrus.Infof("Gathering conditions for %s", objInfo)
	conditions, err := gatherInfraConditions(ctx, objRef, cl)
	if err != nil {
		logrus.Warnf("Failed to gather conditions: %s", err.Error())
		return
	}

	logrus.Infof("Checking conditions for %s", objInfo)
	if len(conditions) > 0 {
		var falsyConditions clusterv1.Conditions
		for _, condition := range conditions {
			if condition.Status != corev1.ConditionTrue {
				falsyConditions = append(falsyConditions, condition)
			}
		}

		if len(falsyConditions) == 0 {
			logrus.Infof("All conditions are satisfied")
		}
		for _, condition := range falsyConditions {
			logrus.Warnf("Condition %s has status: %q, reason: %q, message: %q", condition.Type, condition.Status, condition.Reason, condition.Message)
		}
	} else {
		logrus.Infof("No conditions found")
	}
	logrus.Infof("Done checking conditions for %s", objInfo)
}

// generateSecret returns a Kubernetes secret for the given Cluster and kubeconfig data.
// TODO: upgrade our usage from the deprecated capi cluster v1beta1 -> v1beta2
// and remove this function in favor of Generate secret in capi util/kubeconfig.
func generateSecret(cluster *clusterv1.Cluster, data []byte) *corev1.Secret {
	name := util.ObjectKey(cluster)
	return utilkubeconfig.GenerateSecretWithOwner(name, data, metav1.OwnerReference{
		APIVersion: clusterv1.GroupVersion.String(),
		Kind:       "Cluster",
		Name:       cluster.Name,
		UID:        cluster.UID,
	})
}
