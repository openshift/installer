package infraready

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// FloatingIPs creates or gets the API and Ingress ports and attaches the Floating IPs to them.
func FloatingIPs(ctx context.Context, cluster *capo.OpenStackCluster, installConfig *installconfig.InstallConfig, infraID string) error {
	platformOpenstack := installConfig.Config.OpenStack
	if lb := platformOpenstack.LoadBalancer; lb != nil && lb.Type == configv1.LoadBalancerTypeUserManaged {
		return nil
	}
	networkClient, err := openstackdefaults.NewServiceClient(ctx, "network", openstackdefaults.DefaultClientOpts(installConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return err
	}
	var apiPort, ingressPort *ports.Port
	if platformOpenstack.ControlPlanePort != nil && len(platformOpenstack.ControlPlanePort.FixedIPs) == 2 {
		// To avoid unnecessary calls to Neutron, let's fetch the Ports in case there is a need to attach FIPs
		if platformOpenstack.APIFloatingIP != "" {
			// Using the first VIP as both API VIPs must be allocated on the same Port
			apiPort, err = getPort(ctx, networkClient, cluster.Status.Network.ID, platformOpenstack.APIVIPs[0])
			if err != nil {
				return err
			}
		}
		if platformOpenstack.IngressFloatingIP != "" {
			// Using the first VIP as both Ingress VIPs must be allocated on the same Port
			ingressPort, err = getPort(ctx, networkClient, cluster.Status.Network.ID, platformOpenstack.IngressVIPs[0])
			if err != nil {
				return err
			}
		}
	} else {
		apiPort, err = createPort(ctx, networkClient, "api", infraID, cluster.Status.Network.ID, cluster.Status.Network.Subnets[0].ID, platformOpenstack.APIVIPs[0])
		if err != nil {
			return err
		}
		ingressPort, err = createPort(ctx, networkClient, "ingress", infraID, cluster.Status.Network.ID, cluster.Status.Network.Subnets[0].ID, platformOpenstack.IngressVIPs[0])
		if err != nil {
			return err
		}
	}

	if platformOpenstack.APIFloatingIP != "" {
		if err := assignFIP(ctx, networkClient, platformOpenstack.APIFloatingIP, apiPort); err != nil {
			return err
		}
	}

	if platformOpenstack.IngressFloatingIP != "" {
		if err := assignFIP(ctx, networkClient, platformOpenstack.IngressFloatingIP, ingressPort); err != nil {
			return err
		}
	}

	return nil
}

// getPort retrieves a Neutron Port based on a network ID and the Fixed IP.
func getPort(ctx context.Context, client *gophercloud.ServiceClient, networkID, fixedIP string) (*ports.Port, error) {
	listOpts := ports.ListOpts{
		NetworkID: networkID,
		FixedIPs: []ports.FixedIPOpts{
			{
				IPAddress: fixedIP,
			}},
	}
	allPages, err := ports.List(client, listOpts).AllPages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list Ports: %w", err)
	}
	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return nil, fmt.Errorf("failed to extract Ports: %w", err)
	}
	if len(allPorts) != 1 {
		return nil, fmt.Errorf("could not find Port with IP: %s", fixedIP)
	}
	return &allPorts[0], nil
}

func createPort(ctx context.Context, client *gophercloud.ServiceClient, role, infraID, networkID, subnetID, fixedIP string) (*ports.Port, error) {
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

	port, err := ports.Create(ctx, client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	tag := fmt.Sprintf("openshiftClusterID=%s", infraID)
	err = attributestags.Add(ctx, client, "ports", port.ID, tag).ExtractErr()
	if err != nil {
		return nil, err
	}
	return port, err
}

func assignFIP(ctx context.Context, client *gophercloud.ServiceClient, address string, port *ports.Port) error {
	listOpts := floatingips.ListOpts{
		FloatingIP: address,
	}
	allPages, err := floatingips.List(client, listOpts).AllPages(ctx)
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

	_, err = floatingips.Update(ctx, client, fip.ID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("failed to attach floating IP to port: %w", err)
	}
	return nil
}
