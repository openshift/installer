package azure

import (
	"context"
	"fmt"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

type lbInput struct {
	loadBalancerName       string
	infraID                string
	region                 string
	resourceGroup          string
	subscriptionID         string
	frontendIPConfigName   string
	backendAddressPoolName string
	idPrefix               string
	lbClient               *armnetwork.LoadBalancersClient
	tags                   map[string]*string
}

type pipInput struct {
	infraID       string
	name          string
	region        string
	resourceGroup string
	pipClient     *armnetwork.PublicIPAddressesClient
	tags          map[string]*string
}

type vmInput struct {
	infraID             string
	resourceGroup       string
	ids                 []string
	backendAddressPools []*armnetwork.BackendAddressPool
	vmClient            *armcompute.VirtualMachinesClient
	nicClient           *armnetwork.InterfacesClient
}

type securityGroupInput struct {
	resourceGroupName    string
	securityGroupName    string
	securityRuleName     string
	securityRulePort     string
	networkClientFactory *armnetwork.ClientFactory
}

type inboundNatRuleInput struct {
	resourceGroupName    string
	loadBalancerName     string
	bootstrapNicName     string
	frontendIPConfigID   string
	inboundNatRuleID     string
	inboundNatRuleName   string
	inboundNatRulePort   int32
	networkClientFactory *armnetwork.ClientFactory
}

func createPublicIP(ctx context.Context, in *pipInput) (*armnetwork.PublicIPAddress, error) {
	pollerResp, err := in.pipClient.BeginCreateOrUpdate(
		ctx,
		in.resourceGroup,
		in.name,
		armnetwork.PublicIPAddress{
			Name:     to.Ptr(in.name),
			Location: to.Ptr(in.region),
			SKU: &armnetwork.PublicIPAddressSKU{
				Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
				Tier: to.Ptr(armnetwork.PublicIPAddressSKUTierRegional),
			},
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv4),
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
				DNSSettings: &armnetwork.PublicIPAddressDNSSettings{
					DomainNameLabel: to.Ptr(in.infraID),
				},
			},
			Tags: in.tags,
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

func createAPILoadBalancer(ctx context.Context, pip *armnetwork.PublicIPAddress, in *lbInput) (*armnetwork.LoadBalancer, error) {
	probeName := "api-probe"

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx,
		in.resourceGroup,
		in.loadBalancerName,
		armnetwork.LoadBalancer{
			Location: to.Ptr(in.region),
			SKU: &armnetwork.LoadBalancerSKU{
				Name: to.Ptr(armnetwork.LoadBalancerSKUNameStandard),
				Tier: to.Ptr(armnetwork.LoadBalancerSKUTierRegional),
			},
			Properties: &armnetwork.LoadBalancerPropertiesFormat{
				FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
					{
						Name: &in.frontendIPConfigName,
						Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
							PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
							PublicIPAddress:           pip,
						},
					},
				},
				BackendAddressPools: []*armnetwork.BackendAddressPool{
					{
						Name: &in.backendAddressPoolName,
					},
				},
				Probes: []*armnetwork.Probe{
					{
						Name: &probeName,
						Properties: &armnetwork.ProbePropertiesFormat{
							Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
							Port:              to.Ptr[int32](6443),
							IntervalInSeconds: to.Ptr[int32](5),
							ProbeThreshold:    to.Ptr[int32](2),
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
								ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", in.idPrefix, in.loadBalancerName, in.frontendIPConfigName)),
							},
							BackendAddressPool: &armnetwork.SubResource{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", in.idPrefix, in.loadBalancerName, in.backendAddressPoolName)),
							},
							Probe: &armnetwork.SubResource{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", in.idPrefix, in.loadBalancerName, probeName)),
							},
						},
					},
				},
			},
			Tags: in.tags,
		}, nil)

	if err != nil {
		return nil, fmt.Errorf("cannot create load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func updateOutboundLoadBalancerToAPILoadBalancer(ctx context.Context, pip *armnetwork.PublicIPAddress, in *lbInput) (*armnetwork.LoadBalancer, error) {
	probeName := "api-probe"

	// Get the CAPI-created outbound load balancer so we can modify it.
	extLB, err := in.lbClient.Get(ctx, in.resourceGroup, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get external load balancer: %w", err)
	}

	// Get the existing frontend configuration and backend address pool and
	// create an additional frontend configuration and backend address
	// pool. Use the newly created public IP address with the additional
	// configuration so we can setup load balancing rules for the external
	// API server.
	extLB.Properties.FrontendIPConfigurations = append(extLB.Properties.FrontendIPConfigurations,
		&armnetwork.FrontendIPConfiguration{
			Name: &in.frontendIPConfigName,
			Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
				PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
				PublicIPAddress:           pip,
			},
		})
	extLB.Properties.BackendAddressPools = append(extLB.Properties.BackendAddressPools,
		&armnetwork.BackendAddressPool{
			Name: &in.backendAddressPoolName,
		})

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx,
		in.resourceGroup,
		in.loadBalancerName,
		armnetwork.LoadBalancer{
			Location: to.Ptr(in.region),
			SKU: &armnetwork.LoadBalancerSKU{
				Name: to.Ptr(armnetwork.LoadBalancerSKUNameStandard),
				Tier: to.Ptr(armnetwork.LoadBalancerSKUTierRegional),
			},
			Properties: &armnetwork.LoadBalancerPropertiesFormat{
				FrontendIPConfigurations: extLB.Properties.FrontendIPConfigurations,
				BackendAddressPools:      extLB.Properties.BackendAddressPools,
				Probes: []*armnetwork.Probe{
					{
						Name: &probeName,
						Properties: &armnetwork.ProbePropertiesFormat{
							Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
							Port:              to.Ptr[int32](6443),
							IntervalInSeconds: to.Ptr[int32](5),
							ProbeThreshold:    to.Ptr[int32](2),
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
								ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", in.idPrefix, in.loadBalancerName, in.frontendIPConfigName)),
							},
							BackendAddressPool: &armnetwork.SubResource{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", in.idPrefix, in.loadBalancerName, in.backendAddressPoolName)),
							},
							Probe: &armnetwork.SubResource{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", in.idPrefix, in.loadBalancerName, probeName)),
							},
						},
					},
				},
				OutboundRules: extLB.Properties.OutboundRules,
			},
			Tags: in.tags,
		}, nil)

	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func updateInternalLoadBalancer(ctx context.Context, in *lbInput) (*armnetwork.LoadBalancer, error) {
	mcsProbeName := "sint-probe"

	// Get the CAPI-created internal load balancer so we can modify it.
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroup, in.loadBalancerName, nil)
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
			ProbeThreshold:    to.Ptr[int32](2),
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
				ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", in.idPrefix, in.loadBalancerName, existingFrontEndIPConfigName)),
			},
			BackendAddressPool: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", in.idPrefix, in.loadBalancerName, in.backendAddressPoolName)),
			},
			Probe: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", in.idPrefix, in.loadBalancerName, mcsProbeName)),
			},
		},
	}

	intLB.Properties.Probes = append(intLB.Properties.Probes, mcsProbe)
	intLB.Properties.LoadBalancingRules = append(intLB.Properties.LoadBalancingRules, mcsRule)
	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx,
		in.resourceGroup,
		in.loadBalancerName,
		intLB,
		nil)
	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func associateVMToBackendPool(ctx context.Context, in vmInput) error {
	for _, id := range in.ids {
		vmName := path.Base(id)
		vm, err := in.vmClient.Get(ctx, in.resourceGroup, vmName, nil)
		if err != nil {
			return fmt.Errorf("failed to get vm %s: %w", vmName, err)
		}

		if nics := vm.Properties.NetworkProfile.NetworkInterfaces; len(nics) == 1 {
			nicRef := nics[0]

			nicName := path.Base(*nicRef.ID)
			nic, err := in.nicClient.Get(ctx, in.resourceGroup, nicName, nil)
			if err != nil {
				return fmt.Errorf("failed to get nic for vm %s: %w", vmName, err)
			}
			for _, ipconfig := range nic.Properties.IPConfigurations {
				ipconfig.Properties.LoadBalancerBackendAddressPools = append(ipconfig.Properties.LoadBalancerBackendAddressPools, in.backendAddressPools...)
			}
			pollerResp, err := in.nicClient.BeginCreateOrUpdate(ctx, in.resourceGroup, nicName, nic.Interface, nil)
			if err != nil {
				return fmt.Errorf("failed to update nic for %s: %w", vmName, err)
			}
			_, err = pollerResp.PollUntilDone(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to update nic for vm %s: %w", vmName, err)
			}
		} else {
			return fmt.Errorf("vm %s does not have a single nic: %w", vmName, err)
		}
	}
	return nil
}

func deleteSecurityGroupRule(ctx context.Context, in *securityGroupInput) error {
	securityRulesClient := in.networkClientFactory.NewSecurityRulesClient()
	securityRulesPoller, err := securityRulesClient.BeginDelete(ctx, in.resourceGroupName, in.securityGroupName, in.securityRuleName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete security rule: %w", err)
	}
	_, err = securityRulesPoller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete security rule: %w", err)
	}
	return nil
}

func addInboundNatRuleToLoadBalancer(ctx context.Context, in *inboundNatRuleInput) (*armnetwork.InboundNatRule, error) {
	inboundNatRulesClient := in.networkClientFactory.NewInboundNatRulesClient()
	inboundNatRulesPoller, err := inboundNatRulesClient.BeginCreateOrUpdate(ctx,
		in.resourceGroupName,
		in.loadBalancerName,
		in.inboundNatRuleName,
		armnetwork.InboundNatRule{
			Properties: &armnetwork.InboundNatRulePropertiesFormat{
				BackendPort: to.Ptr[int32](in.inboundNatRulePort),
				FrontendIPConfiguration: &armnetwork.SubResource{
					ID: to.Ptr(in.frontendIPConfigID),
				},
				FrontendPort: to.Ptr[int32](in.inboundNatRulePort),
				Protocol:     to.Ptr(armnetwork.TransportProtocolTCP), // assume TCP for now
			},
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add inbound nat rule to load balancer: %w", err)
	}
	inboundNatRuleResp, err := inboundNatRulesPoller.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to add inbound nat rule to load balancer: %w", err)
	}

	return &inboundNatRuleResp.InboundNatRule, nil
}

func deleteInboundNatRule(ctx context.Context, in *inboundNatRuleInput) error {
	inboundNatRulesClient := in.networkClientFactory.NewInboundNatRulesClient()
	inboundNatRulesPoller, err := inboundNatRulesClient.BeginDelete(ctx,
		in.resourceGroupName,
		in.loadBalancerName,
		in.inboundNatRuleName,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to delete inbound nat rule: %w", err)
	}
	_, err = inboundNatRulesPoller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete inbound nat rule: %w", err)
	}
	return nil
}

func associateInboundNatRuleToInterface(ctx context.Context, in *inboundNatRuleInput) (*armnetwork.Interface, error) {
	interfacesClient := in.networkClientFactory.NewInterfacesClient()
	interfaceResp, err := interfacesClient.Get(ctx,
		in.resourceGroupName,
		in.bootstrapNicName,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get interface: %w", err)
	}
	bootstrapInterface := interfaceResp.Interface

	inboundNatRulesClient := in.networkClientFactory.NewInboundNatRulesClient()
	inboundNatRulesResp, err := inboundNatRulesClient.Get(ctx,
		in.resourceGroupName,
		in.loadBalancerName,
		in.inboundNatRuleName,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get inbound nat rule: %w", err)
	}
	inboundNatRule := inboundNatRulesResp.InboundNatRule

	for i, ipConfig := range bootstrapInterface.Properties.IPConfigurations {
		ipConfig.Properties.LoadBalancerInboundNatRules = append(ipConfig.Properties.LoadBalancerInboundNatRules,
			&inboundNatRule,
		)
		bootstrapInterface.Properties.IPConfigurations[i] = ipConfig
	}

	interfacesPollerResp, err := interfacesClient.BeginCreateOrUpdate(ctx,
		in.resourceGroupName,
		in.bootstrapNicName,
		bootstrapInterface,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add inbound nat rule to interface: %w", err)
	}

	interfacesResp, err := interfacesPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to add inbound nat rule to interface: %w", err)
	}
	return &interfacesResp.Interface, nil
}

type natGatewayInput struct {
	infraID        string
	cl             client.Client
	subscriptionID string
	creds          azcore.TokenCredential
	cloudConfig    cloud.Configuration
}

func associateNatGatewayToSubnet(ctx context.Context, in natGatewayInput) error {
	clientOpts := &arm.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: in.cloudConfig,
		},
	}
	subnetsClient, err := armnetwork.NewSubnetsClient(in.subscriptionID, in.creds, clientOpts)
	if err != nil {
		return fmt.Errorf("failed to get subnet client: %w", err)
	}

	azureCluster := &capz.AzureCluster{}
	key := client.ObjectKey{
		Name:      in.infraID,
		Namespace: capiutils.Namespace,
	}
	if err := in.cl.Get(context.Background(), key, azureCluster); err != nil {
		return fmt.Errorf("failed to get AzureCluster: %w", err)
	}

	subnets := azureCluster.Spec.NetworkSpec.Subnets
	for _, existingSubnet := range subnets {
		if existingSubnet.Role == capz.SubnetControlPlane {
			continue
		}
		if existingSubnet.NatGateway.Name == "" {
			continue
		}
		natGatewayID := fmt.Sprintf(
			"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/natGateways/%s",
			in.subscriptionID,
			azureCluster.Spec.ResourceGroup,
			existingSubnet.NatGateway.Name,
		)

		subnet, err := subnetsClient.Get(ctx,
			azureCluster.Spec.NetworkSpec.Vnet.ResourceGroup,
			azureCluster.Spec.NetworkSpec.Vnet.Name,
			existingSubnet.Name,
			nil)
		if err != nil {
			return fmt.Errorf("failed to get subnet: %w", err)
		}

		subnet.Properties.NatGateway = &armnetwork.SubResource{
			ID: &natGatewayID,
		}

		poller, err := subnetsClient.BeginCreateOrUpdate(ctx,
			azureCluster.Spec.NetworkSpec.Vnet.ResourceGroup,
			azureCluster.Spec.NetworkSpec.Vnet.Name,
			*subnet.Name,
			subnet.Subnet,
			nil)
		if err != nil {
			return fmt.Errorf("failed to begin subnet update: %w", err)
		}

		_, err = poller.PollUntilDone(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to update subnet: %w", err)
		}
	}
	return nil
}
