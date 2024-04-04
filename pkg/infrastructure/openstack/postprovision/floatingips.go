package postprovision

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// FloatingIPs creates and attaches a Floating IP to the Bootstrap Machine.
func FloatingIPs(ctx context.Context, c client.Client, cluster *capo.OpenStackCluster, installConfig *installconfig.InstallConfig, infraID string) error {
	bootstrapMachine := &capo.OpenStackMachine{}
	key := client.ObjectKey{
		Name:      capiutils.GenerateBoostrapMachineName(infraID),
		Namespace: capiutils.Namespace,
	}
	if err := c.Get(ctx, key, bootstrapMachine); err != nil {
		return fmt.Errorf("failed to get bootstrap Machine: %w", err)
	}

	networkClient, err := openstackdefaults.NewServiceClient("network", openstackdefaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return err
	}

	bootstrapPort, err := getPortForInstance(networkClient, *bootstrapMachine.Status.InstanceID)
	if err != nil {
		return err
	}

	_, err = createAndAttachFIP(networkClient, "bootstrap", infraID, cluster.Status.ExternalNetwork.ID, bootstrapPort.ID)
	if err != nil {
		return err
	}

	return nil
}

// Get the first port associated with an instance.
//
// This is used to attach a FIP to the instance.
func getPortForInstance(client *gophercloud.ServiceClient, instanceID string) (*ports.Port, error) {
	listOpts := ports.ListOpts{
		DeviceID: instanceID,
	}
	allPages, err := ports.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("failed to list ports: %w", err)
	}
	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return nil, fmt.Errorf("failed to extract ports: %w", err)
	}

	if len(allPorts) < 1 {
		return nil, fmt.Errorf("bootstrap machine has no associated ports")
	}

	return &allPorts[0], nil
}

// Create a floating IP.
func createAndAttachFIP(client *gophercloud.ServiceClient, role, infraID, networkID, portID string) (*floatingips.FloatingIP, error) {
	createOpts := floatingips.CreateOpts{
		Description:       fmt.Sprintf("%s-%s-fip", infraID, role),
		FloatingNetworkID: networkID,
		PortID:            portID,
	}

	floatingIP, err := floatingips.Create(client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tag := fmt.Sprintf("openshiftClusterID=%s", infraID)
	err = attributestags.Add(client, "floatingips", floatingIP.ID, tag).ExtractErr()
	if err != nil {
		return nil, err
	}

	return floatingIP, err
}
