package postprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// FloatingIPs creates and attaches a Floating IP to the Bootstrap Machine.
func FloatingIPs(ctx context.Context, c client.Client, cluster *capo.OpenStackCluster, installConfig *installconfig.InstallConfig, infraID string) error {
	if cluster.Status.ExternalNetwork == nil {
		return nil
	}
	bootstrapMachine := &capo.OpenStackMachine{}
	key := client.ObjectKey{
		Name:      capiutils.GenerateBoostrapMachineName(infraID),
		Namespace: capiutils.Namespace,
	}
	if err := c.Get(ctx, key, bootstrapMachine); err != nil {
		return fmt.Errorf("failed to get bootstrap Machine: %w", err)
	}

	networkClient, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return err
	}

	bootstrapPort, err := getPortForInstance(ctx, networkClient, *bootstrapMachine.Status.InstanceID, cluster.Status.Network.ID)
	if err != nil {
		return err
	}

	_, err = createAndAttachFIP(ctx, networkClient, "bootstrap", infraID, cluster.Status.ExternalNetwork.ID, bootstrapPort.ID)
	if err != nil {
		return err
	}

	return nil
}

// Return the first port associated with the given instance and existing on the
// given network.
func getPortForInstance(ctx context.Context, client *gophercloud.ServiceClient, instanceID, networkID string) (*ports.Port, error) {
	var port *ports.Port
	if err := ports.List(client, ports.ListOpts{DeviceID: instanceID, NetworkID: networkID}).EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
		ports, err := ports.ExtractPorts(page)
		if err != nil {
			return false, err
		}
		if len(ports) > 0 {
			port = &ports[0]
			return false, nil
		}
		return true, nil
	}); err != nil {
		return nil, fmt.Errorf("failed to list the ports of the bootstrap server: %w", err)
	}
	if port == nil {
		return nil, fmt.Errorf("the bootstrap server has no associated ports")
	}
	return port, nil
}

// Create a floating IP.
func createAndAttachFIP(ctx context.Context, client *gophercloud.ServiceClient, role, infraID, networkID, portID string) (*floatingips.FloatingIP, error) {
	createOpts := floatingips.CreateOpts{
		Description:       fmt.Sprintf("%s-%s-fip", infraID, role),
		FloatingNetworkID: networkID,
		PortID:            portID,
	}

	floatingIP, err := floatingips.Create(ctx, client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tag := fmt.Sprintf("openshiftClusterID=%s", infraID)
	err = attributestags.Add(ctx, client, "floatingips", floatingIP.ID, tag).ExtractErr()
	if err != nil {
		return nil, err
	}

	return floatingIP, err
}
