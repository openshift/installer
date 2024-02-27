package infraready

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"

	"github.com/openshift/installer/pkg/asset/installconfig"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// FloatingIPs creates the API and Ingress ports and attaches the Floating IPs to them.
func FloatingIPs(cluster *capo.OpenStackCluster, installConfig *installconfig.InstallConfig, infraID string) error {
	networkClient, err := openstackdefaults.NewServiceClient("network", openstackdefaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return err
	}

	apiPort, err := createPort(networkClient, "api", infraID, cluster.Status.Network.ID, cluster.Status.Network.Subnets[0].ID, installConfig.Config.Platform.OpenStack.APIVIPs[0])
	if err != nil {
		return err
	}
	if installConfig.Config.OpenStack.APIFloatingIP != "" {
		if err := assignFIP(networkClient, installConfig.Config.OpenStack.APIFloatingIP, apiPort); err != nil {
			return err
		}
	}

	ingressPort, err := createPort(networkClient, "ingress", infraID, cluster.Status.Network.ID, cluster.Status.Network.Subnets[0].ID, installConfig.Config.Platform.OpenStack.IngressVIPs[0])
	if err != nil {
		return err
	}
	if installConfig.Config.OpenStack.IngressFloatingIP != "" {
		if err := assignFIP(networkClient, installConfig.Config.OpenStack.IngressFloatingIP, ingressPort); err != nil {
			return err
		}
	}

	return nil
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
