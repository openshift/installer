package clusterapi

import (
	"context"
	"fmt"

	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/controller-runtime/pkg/client"

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

	createOtps := ports.CreateOpts{
		Name:      "CAPO test",
		NetworkID: ospCluster.Status.Network.ID,
	}

	_, err = ports.Create(networkClient, createOtps).Extract()
	if err != nil {
		return err
	}

	_, err = ports.Create(networkClient, createOtps).Extract()
	if err != nil {
		return err
	}

	return nil
}
