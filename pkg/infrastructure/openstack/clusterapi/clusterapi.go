package clusterapi

import (
	"context"
	"fmt"

	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/infrastructure/clusterapi"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// Provider defines the InfraProvider.
type Provider struct {
	clusterapi.InfraProvider
}

// func (p Provider) PreProvision(in clusterapi.PreProvisionInput) error {
// 	return nil
// }

func (p Provider) ControlPlaneAvailable(in clusterapi.ControlPlaneAvailableInput) error {
	ospCluster := &capo.OpenStackCluster{}
	key := client.ObjectKey{
		Name:      in.InfraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.Client.Get(context.Background(), key, ospCluster); err != nil {
		return fmt.Errorf("failed to get OSPCluster: %w", err)
	}

	networkClient, err := openstackdefaults.NewServiceClient("network", openstackdefaults.DefaultClientOpts(in.InstallConfig.Config.Platform.OpenStack.Cloud))
	if err != nil {
		return err
	}

	apiPort, err := createPort(networkClient, "api", in.InfraID, ospCluster.Status.Network.ID)
	if err != nil {
		return err
	}
	if in.InstallConfig.Config.OpenStack.APIFloatingIP != "" {
		err = assignFIP(networkClient, in.InstallConfig.Config.OpenStack.APIFloatingIP, apiPort)
		if err != nil {
			return err
		}
	}

	ingressPort, err := createPort(networkClient, "ingress", in.InfraID, ospCluster.Status.Network.ID)
	if err != nil {
		return err
	}
	if in.InstallConfig.Config.OpenStack.IngressFloatingIP != "" {
		err = assignFIP(networkClient, in.InstallConfig.Config.OpenStack.IngressFloatingIP, ingressPort)
		if err != nil {
			return err
		}
	}

	return nil
}

func createPort(client *gophercloud.ServiceClient, role, infraID, networkID string) (*ports.Port, error) {
	createOtps := ports.CreateOpts{
		Name:        fmt.Sprintf("%s-%s-port", infraID, role),
		NetworkID:   networkID,
		Description: "Created By OpenShift Installer",
	}

	port, err := ports.Create(client, createOtps).Extract()
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
