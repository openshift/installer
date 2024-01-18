package clusterapi

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/manifests"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	"github.com/openshift/installer/pkg/infrastructure/openstack/preprovision"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types/openstack"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// Provider defines the InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

// Name contains the name of the openstack provider.
func (p Provider) Name() string {
	return openstack.Name
}

var _ clusterapi.PreProvider = Provider{}

// PreProvision tags the VIP ports and creates the security groups that are needed during CAPI provisioning.
func (p Provider) PreProvision(ctx context.Context, in clusterapi.PreProvisionInput) error {
	var (
		infraID            = in.InfraID
		installConfig      = in.InstallConfig
		rhcosImage         = string(*in.RhcosImage)
		mastersSchedulable bool
	)

	for _, f := range in.ManifestsAsset.Files() {
		if f.Filename == manifests.SchedulerCfgFilename {
			schedulerConfig := configv1.Scheduler{}
			if err := yaml.Unmarshal(f.Data, &schedulerConfig); err != nil {
				return fmt.Errorf("unable to decode the scheduler manifest: %w", err)
			}
			mastersSchedulable = schedulerConfig.Spec.MastersSchedulable
			break
		}
	}

	if err := preprovision.TagVIPPorts(ctx, installConfig, infraID); err != nil {
		return err
	}

	// upload the corresponding image to Glance if rhcosImage contains a
	// URL. If rhcosImage contains a name, then that points to an existing
	// Glance image.
	if imageName, isURL := rhcos.GenerateOpenStackImageName(rhcosImage, infraID); isURL {
		if err := preprovision.UploadBaseImage(ctx, installConfig.Config.Platform.OpenStack.Cloud, rhcosImage, imageName, infraID, installConfig.Config.Platform.OpenStack.ClusterOSImageProperties); err != nil {
			return err
		}
	}

	return preprovision.SecurityGroups(ctx, installConfig, infraID, mastersSchedulable)
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

	networkClient, err := openstackdefaults.NewServiceClient("network", openstackdefaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return err
	}

	apiPort, err := createPort(networkClient, "api", infraID, ospCluster.Status.Network.ID, ospCluster.Status.Network.Subnets[0].ID, installConfig.Config.Platform.OpenStack.APIVIPs[0])
	if err != nil {
		return err
	}
	if installConfig.Config.OpenStack.APIFloatingIP != "" {
		err = assignFIP(networkClient, installConfig.Config.OpenStack.APIFloatingIP, apiPort)
		if err != nil {
			return err
		}
	}

	ingressPort, err := createPort(networkClient, "ingress", infraID, ospCluster.Status.Network.ID, ospCluster.Status.Network.Subnets[0].ID, installConfig.Config.Platform.OpenStack.IngressVIPs[0])
	if err != nil {
		return err
	}
	if installConfig.Config.OpenStack.IngressFloatingIP != "" {
		err = assignFIP(networkClient, installConfig.Config.OpenStack.IngressFloatingIP, ingressPort)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ clusterapi.IgnitionProvider = Provider{}

// Ignition uploads the bootstrap machine's Ignition file to OpenStack.
func (p Provider) Ignition(ctx context.Context, in clusterapi.IgnitionInput) ([]byte, error) {
	logrus.Debugf("Uploading the bootstrap machine's Ignition file to OpenStack")
	var (
		bootstrapIgnData = in.BootstrapIgnData
		infraID          = in.InfraID
		installConfig    = in.InstallConfig
	)

	return preprovision.UploadIgnitionAndBuildShim(ctx, installConfig.Config.Platform.OpenStack.Cloud, infraID, installConfig.Config.Proxy, bootstrapIgnData)
}

func createPort(client *gophercloud.ServiceClient, role, infraID, networkID, subnetID, fixedIP string) (*ports.Port, error) {
	createOpts := ports.CreateOpts{
		Name:        fmt.Sprintf("%s-%s-port", infraID, role),
		NetworkID:   networkID,
		Description: "Created By OpenShift Installer",
		FixedIPs: []ports.IP{
			{
				IPAddress: fixedIP,
				SubnetID:  subnetID,
			}},
	}

	port, err := ports.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tag := fmt.Sprintf("openshiftClusterID=%s", infraID)
	err = attributestags.Add(client, "ports", port.ID, tag).ExtractErr()
	if err != nil {
		return nil, err
	}
	return port, err
}

func assignFIP(client *gophercloud.ServiceClient, address string, port *ports.Port) error {
	listOpts := floatingips.ListOpts{
		FloatingIP: address,
	}
	allPages, err := floatingips.List(client, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("failed to list floating IPs: %w", err)
	}
	allFIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return fmt.Errorf("failed to extract floating IPs: %w", err)
	}

	if len(allFIPs) != 1 {
		return fmt.Errorf("could not find FIP: %s", address)
	}

	fip := allFIPs[0]

	updateOpts := floatingips.UpdateOpts{
		PortID: &port.ID,
	}

	_, err = floatingips.Update(client, fip.ID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("failed to attach floating IP to port: %w", err)
	}
	return nil
}
