package clusterapi

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	mapov1alpha1 "github.com/openshift/api/machine/v1alpha1"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/infrastructure/openstack/infraready"
	"github.com/openshift/installer/pkg/infrastructure/openstack/postprovision"
	"github.com/openshift/installer/pkg/infrastructure/openstack/preprovision"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types/openstack"
)

// Provider defines the InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

// Name contains the name of the openstack provider.
func (p Provider) Name() string {
	return openstack.Name
}

// PublicGatherEndpoint indicates that machine ready checks should NOT wait for an ExternalIP
// in the status when declaring machines ready.
func (Provider) PublicGatherEndpoint() clusterapi.GatherEndpoint { return clusterapi.InternalIP }

var _ clusterapi.PreProvider = Provider{}

// PreProvision tags the VIP ports, and creates the security groups and the
// server groups defined in the Machine manifests.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	var (
		infraID          = in.InfraID
		installConfig    = in.InstallConfig
		rhcosImage       = in.RhcosImage.ControlPlane
		manifestsAsset   = in.ManifestsAsset
		machineManifests = in.MachineManifests
		workersAsset     = in.WorkersAsset
	)

	if err := preprovision.TagVIPPorts(ctx, installConfig, infraID); err != nil {
		return fmt.Errorf("failed to tag VIP ports: %w", err)
	}

	// upload the corresponding image to Glance if rhcosImage contains a
	// URL. If rhcosImage contains a name, then that points to an existing
	// Glance image.
	if imageName, isURL := rhcos.GenerateOpenStackImageName(rhcosImage, infraID); isURL {
		if err := preprovision.UploadBaseImage(ctx, installConfig.Config.Platform.OpenStack.Cloud, rhcosImage, imageName, infraID, installConfig.Config.Platform.OpenStack.ClusterOSImageProperties); err != nil {
			return fmt.Errorf("failed to upload the RHCOS base image: %w", err)
		}
	}

	{
		var mastersSchedulable bool
		for _, f := range manifestsAsset.Files() {
			if f.Filename == manifests.SchedulerCfgFilename {
				schedulerConfig := configv1.Scheduler{}
				if err := yaml.Unmarshal(f.Data, &schedulerConfig); err != nil {
					return fmt.Errorf("unable to decode the scheduler manifest: %w", err)
				}
				mastersSchedulable = schedulerConfig.Spec.MastersSchedulable
				break
			}
		}
		if err := preprovision.SecurityGroups(ctx, installConfig, infraID, mastersSchedulable); err != nil {
			return fmt.Errorf("failed to create security groups: %w", err)
		}
	}

	{
		capiMachines := make([]capo.OpenStackMachine, 0, len(machineManifests))
		for i := range machineManifests {
			if m, ok := machineManifests[i].(*capo.OpenStackMachine); ok {
				capiMachines = append(capiMachines, *m)
			}
		}

		var workerSpecs []mapov1alpha1.OpenstackProviderSpec
		{
			machineSets, err := workersAsset.MachineSets()
			if err != nil {
				return fmt.Errorf("failed to extract MachineSets from the worker assets: %w", err)
			}
			workerSpecs = make([]mapov1alpha1.OpenstackProviderSpec, 0, len(machineSets))
			for _, machineSet := range machineSets {
				workerSpecs = append(workerSpecs, *machineSet.Spec.Template.Spec.ProviderSpec.Value.Object.(*mapov1alpha1.OpenstackProviderSpec))
			}
		}
		if err := preprovision.ServerGroups(ctx, installConfig, capiMachines, workerSpecs); err != nil {
			return fmt.Errorf("failed to create server groups: %w", err)
		}
	}

	return nil
}

var _ clusterapi.InfraReadyProvider = Provider{}

// InfraReady creates the API and Ingress ports and attaches the Floating IPs to them.
func (p Provider) InfraReady(ctx context.Context, in clusterapi.InfraReadyInput) error {
	var (
		k8sClient     = in.Client
		infraID       = in.InfraID
		installConfig = in.InstallConfig
	)

	ospCluster := &capo.OpenStackCluster{}
	key := client.ObjectKey{
		Name:      infraID,
		Namespace: capiutils.Namespace,
	}
	if err := k8sClient.Get(ctx, key, ospCluster); err != nil {
		return fmt.Errorf("failed to get OSPCluster: %w", err)
	}

	return infraready.FloatingIPs(ctx, ospCluster, installConfig, infraID)
}

var _ clusterapi.IgnitionProvider = Provider{}

// Ignition uploads the bootstrap machine's Ignition file to OpenStack.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]*corev1.Secret, error) {
	logrus.Debugf("Uploading the bootstrap machine's Ignition file to OpenStack")
	var (
		bootstrapIgnData = in.BootstrapIgnData
		infraID          = in.InfraID
		installConfig    = in.InstallConfig
	)

	ignShim, err := preprovision.UploadIgnitionAndBuildShim(ctx, installConfig.Config.Platform.OpenStack.Cloud, infraID, installConfig.Config.Proxy, bootstrapIgnData)
	if err != nil {
		return nil, fmt.Errorf("failed to upload and build ignition shim: %w", err)
	}

	ignSecrets := []*corev1.Secret{
		clusterapi.IgnitionSecret(ignShim, in.InfraID, "bootstrap"),
		clusterapi.IgnitionSecret(in.MasterIgnData, in.InfraID, "master"),
	}
	return ignSecrets, nil
}

var _ clusterapi.PostProvider = Provider{}

// PostProvision creates and attaches a Floating IP to the Bootstrap Machine.
func (p Provider) PostProvision(ctx context.Context, in clusterapi.PostProvisionInput) error {
	var (
		k8sClient     = in.Client
		infraID       = in.InfraID
		installConfig = in.InstallConfig
	)

	ospCluster := &capo.OpenStackCluster{}
	key := client.ObjectKey{
		Name:      infraID,
		Namespace: capiutils.Namespace,
	}
	if err := k8sClient.Get(ctx, key, ospCluster); err != nil {
		return fmt.Errorf("failed to get OSPCluster: %w", err)
	}

	return postprovision.FloatingIPs(ctx, k8sClient, ospCluster, installConfig, infraID)
}
