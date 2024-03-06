package azure

import (
	"context"
	"fmt"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
)

func createPublicIP(ctx context.Context, infraID, region string, pipClient *armnetwork.PublicIPAddressesClient) (*armnetwork.PublicIPAddress, error) {
	publicIPAddressName := fmt.Sprintf("%s-pip-v4", infraID)
	resourceGroupName := fmt.Sprintf("%s-rg", infraID)

	pollerResp, err := pipClient.BeginCreateOrUpdate(
		ctx,
		resourceGroupName,
		publicIPAddressName,
		armnetwork.PublicIPAddress{
			Name:     to.Ptr(publicIPAddressName),
			Location: to.Ptr(region),
			SKU: &armnetwork.PublicIPAddressSKU{
				Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
				Tier: to.Ptr(armnetwork.PublicIPAddressSKUTierRegional),
			},
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv4),
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
			},
		},
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.PublicIPAddress, nil
}

func createExternalLoadBalancer(ctx context.Context, infraID, region, subscriptionID string, pip *armnetwork.PublicIPAddress, lbClient *armnetwork.LoadBalancersClient) (*armnetwork.LoadBalancer, error) {
	resourceGroupName := fmt.Sprintf("%s-rg", infraID)
	loadBalancerName := infraID
	probeName := "api-probe"
	frontEndIPConfigName := "public-lb-ip-v4"
	backEndAddressPoolName := infraID
	idPrefix := fmt.Sprintf("subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers", subscriptionID, resourceGroupName)

	pollerResp, err := lbClient.BeginCreateOrUpdate(ctx,
		resourceGroupName,
		loadBalancerName,
		armnetwork.LoadBalancer{
			Location: to.Ptr(region),
			SKU: &armnetwork.LoadBalancerSKU{
				Name: to.Ptr(armnetwork.LoadBalancerSKUNameStandard),
				Tier: to.Ptr(armnetwork.LoadBalancerSKUTierRegional),
			},
			Properties: &armnetwork.LoadBalancerPropertiesFormat{
				FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
					{
						Name: &frontEndIPConfigName,
						Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
							PublicIPAddress:           pip,
						},
					},
				},
				BackendAddressPools: []*armnetwork.BackendAddressPool{
					{
						Name: &backEndAddressPoolName,
					},
				},
				Probes: []*armnetwork.Probe{
					{
						Name: &probeName,
						Properties: &armnetwork.ProbePropertiesFormat{
							Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
							Port:              to.Ptr[int32](6443),
							IntervalInSeconds: to.Ptr[int32](5),
							NumberOfProbes:    to.Ptr[int32](2),
							RequestPath:       to.Ptr("/readyz"),
						},
					},
				},
				LoadBalancingRules: []*armnetwork.LoadBalancingRule{
					{
						Name: to.Ptr("api-v4"),
						Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
							Protocol:             to.Ptr(armnetwork.TransportProtocolTCP),
							FrontendPort:         to.Ptr[int32](6443),
							BackendPort:          to.Ptr[int32](6443),
							IdleTimeoutInMinutes: to.Ptr[int32](30),
							EnableFloatingIP:     to.Ptr(false),
							LoadDistribution:     to.Ptr(armnetwork.LoadDistributionDefault),
							FrontendIPConfiguration: &armnetwork.SubResource{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", idPrefix, loadBalancerName, frontEndIPConfigName)),
							},
							BackendAddressPool: &armnetwork.SubResource{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", idPrefix, loadBalancerName, backEndAddressPoolName)),
							},
							Probe: &armnetwork.SubResource{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", idPrefix, loadBalancerName, probeName)),
							},
						},
					},
				},
			},
		}, nil)

	if err != nil {
		return nil, fmt.Errorf("cannot create load balancer: %v", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func updateInternalLoadBalancer(ctx context.Context, infraID, region, subscriptionID string, lbClient *armnetwork.LoadBalancersClient) (*armnetwork.LoadBalancer, error) {
	resourceGroupName := fmt.Sprintf("%s-rg", infraID)
	loadBalancerName := fmt.Sprintf("%s-internal", infraID)
	mcsProbeName := "sint-probe"
	backEndAddressPoolName := fmt.Sprintf("%s-internal", infraID)
	idPrefix := fmt.Sprintf("subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers", subscriptionID, resourceGroupName)

	// Get the CAPI-created internal load balancer so we can modify it.
	lbResp, err := lbClient.Get(ctx, resourceGroupName, loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get internal load balancer: %w", err)
	}
	intLB := lbResp.LoadBalancer

	mcsProbe := &armnetwork.Probe{
		Name: to.Ptr(mcsProbeName),
		Properties: &armnetwork.ProbePropertiesFormat{
			Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
			Port:              to.Ptr[int32](22623),
			IntervalInSeconds: to.Ptr[int32](5),
			NumberOfProbes:    to.Ptr[int32](2),
			RequestPath:       to.Ptr("/healthz"),
		},
	}

	existingFrontEndIPConfig := intLB.Properties.FrontendIPConfigurations
	if len(existingFrontEndIPConfig) == 0 {
		return nil, fmt.Errorf("could not get frontEndIPConfig for internal LB %s", *intLB.Name)
	}
	existingFrontEndIPConfigName := *(existingFrontEndIPConfig[0].Name)

	mcsRule := &armnetwork.LoadBalancingRule{
		Name: to.Ptr("sint-v4"),
		Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
			Protocol:             to.Ptr(armnetwork.TransportProtocolTCP),
			FrontendPort:         to.Ptr[int32](22623),
			BackendPort:          to.Ptr[int32](22623),
			IdleTimeoutInMinutes: to.Ptr[int32](30),
			EnableFloatingIP:     to.Ptr(false),
			LoadDistribution:     to.Ptr(armnetwork.LoadDistributionDefault),
			FrontendIPConfiguration: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", idPrefix, loadBalancerName, existingFrontEndIPConfigName)),
			},
			BackendAddressPool: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", idPrefix, loadBalancerName, backEndAddressPoolName)),
			},
			Probe: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", idPrefix, loadBalancerName, mcsProbeName)),
			},
		},
	}

	intLB.Properties.Probes = append(intLB.Properties.Probes, mcsProbe)
	intLB.Properties.LoadBalancingRules = append(intLB.Properties.LoadBalancingRules, mcsRule)
	pollerResp, err := lbClient.BeginCreateOrUpdate(ctx,
		resourceGroupName,
		loadBalancerName,
		intLB,
		nil)
	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %v", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func associateVMToBackendPool(
	ctx context.Context,
	infraID string,
	ids []string,
	bap *armnetwork.BackendAddressPool,
	vmClient *armcompute.VirtualMachinesClient,
	nicClient *armnetwork.InterfacesClient,
) error {
	resourceGroupName := fmt.Sprintf("%s-rg", infraID)

	for _, id := range ids {
		vmName := path.Base(id)
		vm, err := vmClient.Get(ctx, resourceGroupName, vmName, nil)
		if err != nil {
			return fmt.Errorf("failed to get vm %s: %w", vmName, err)
		}

		if nics := vm.Properties.NetworkProfile.NetworkInterfaces; len(nics) == 1 {
			nicRef := nics[0]

			nicName := path.Base(*nicRef.ID)
			nic, err := nicClient.Get(ctx, resourceGroupName, nicName, nil)
			if err != nil {
				return fmt.Errorf("failed to get nic for vm %s: %w", vmName, err)
			}
			for _, ipconfig := range nic.Properties.IPConfigurations {
				baps := ipconfig.Properties.LoadBalancerBackendAddressPools
				baps = append(baps, bap)
			}
			pollerResp, err := nicClient.BeginCreateOrUpdate(ctx, resourceGroupName, nicName, nic.Interface, nil)
			if err != nil {
				return fmt.Errorf("failed to update nic for %s: %w", vmName, err)
			}
			_, err = pollerResp.PollUntilDone(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to update nic for vm %s: %w", vmName, err)
			}
		} else {
			return fmt.Errorf("vm %s does not have a single VM: %v", vmName, nics)
		}
	}
	return nil
}
