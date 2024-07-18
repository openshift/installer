package azure

import (
	"context"
	"fmt"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"k8s.io/utils/ptr"
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
	infraID         string
	name            string
	region          string
	resourceGroup   string
	createDNSRecord bool
	ipVersion       armnetwork.IPVersion
	pipClient       *armnetwork.PublicIPAddressesClient
	tags            map[string]*string
}

type vmInput struct {
	infraID       string
	resourceGroup string
	ids           []string
	bap           *armnetwork.BackendAddressPool
	vmClient      *armcompute.VirtualMachinesClient
	nicClient     *armnetwork.InterfacesClient
}

type securityGroupInput struct {
	resourceGroupName    string
	securityGroupName    string
	securityRuleName     string
	securityRulePort     string
	securityRulePriority int32
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
	var dnsSettings *armnetwork.PublicIPAddressDNSSettings
	if in.createDNSRecord {
		dnsSettings = &armnetwork.PublicIPAddressDNSSettings{
			DomainNameLabel: to.Ptr(in.infraID),
		}
	}

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
				PublicIPAddressVersion:   to.Ptr(in.ipVersion),
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
				DNSSettings:              dnsSettings,
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
				ipconfig.Properties.LoadBalancerBackendAddressPools = append(ipconfig.Properties.LoadBalancerBackendAddressPools, in.bap)
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

func addSecurityGroupRule(ctx context.Context, in *securityGroupInput) error {
	securityRulesClient := in.networkClientFactory.NewSecurityRulesClient()

	// Assume inbound tcp connections from any port to destination port for now
	securityRuleResp, err := securityRulesClient.BeginCreateOrUpdate(ctx,
		in.resourceGroupName,
		in.securityGroupName,
		in.securityRuleName,
		armnetwork.SecurityRule{
			Name: ptr.To(in.securityRuleName),
			Properties: &armnetwork.SecurityRulePropertiesFormat{
				Access:                   ptr.To(armnetwork.SecurityRuleAccessAllow),
				Direction:                ptr.To(armnetwork.SecurityRuleDirectionInbound),
				Protocol:                 ptr.To(armnetwork.SecurityRuleProtocolTCP),
				DestinationAddressPrefix: ptr.To("*"),
				DestinationPortRange:     ptr.To(in.securityRulePort),
				Priority:                 ptr.To[int32](in.securityRulePriority),
				SourceAddressPrefix:      ptr.To("*"),
				SourcePortRange:          ptr.To("*"),
			},
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to add security rule: %w", err)
	}
	_, err = securityRuleResp.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to add security rule: %w", err)
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

// CreateNatGatewayInput contains input parameters for NAT gateway creation.
type CreateNatGatewayInput struct {
	InfraID                string
	SubscriptionID         string
	ResourceGroupName      string
	VirtualNetworkName     string
	ControlPlaneSubnetName string
	WorkerSubnetName       string
	Region                 string
	Tags                   map[string]*string
	NetworkClientFactory   *armnetwork.ClientFactory
}

// CreateNatGateway creates a NAT gateway and associated resources.
func CreateNatGateway(ctx context.Context, in *CreateNatGatewayInput) (*armnetwork.NatGateway, error) {
	// Create a public IPv4 address.
	publicIPv4Address, err := createPublicIP(ctx, &pipInput{
		name:            fmt.Sprintf("%s-natgw-pip-v4", in.InfraID),
		infraID:         in.InfraID,
		region:          in.Region,
		resourceGroup:   in.ResourceGroupName,
		createDNSRecord: false,
		ipVersion:       armnetwork.IPVersionIPv4,
		pipClient:       in.NetworkClientFactory.NewPublicIPAddressesClient(),
		tags:            in.Tags,
	})
	if err != nil {
		return nil, err
	}

	// Create the NAT gateway.
	natGatewaysClient := in.NetworkClientFactory.NewNatGatewaysClient()
	natGatewayPollerResp, err := natGatewaysClient.BeginCreateOrUpdate(
		ctx,
		in.ResourceGroupName,
		fmt.Sprintf("%s-natgw", in.InfraID),
		armnetwork.NatGateway{
			Location: to.Ptr(in.Region),
			Properties: &armnetwork.NatGatewayPropertiesFormat{
				IdleTimeoutInMinutes: to.Ptr[int32](10),
				PublicIPAddresses: []*armnetwork.SubResource{
					{
						ID: publicIPv4Address.ID,
					},
				},
			},
			SKU: &armnetwork.NatGatewaySKU{
				Name: ptr.To(armnetwork.NatGatewaySKUNameStandard),
			},
			Tags: in.Tags,
		},
		nil,
	)
	if err != nil {
		return nil, err
	}
	natGatewaysResp, err := natGatewayPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	natGateway := natGatewaysResp.NatGateway

	// Update the control plane subnet to use the NAT gateway.
	subnetsClient := in.NetworkClientFactory.NewSubnetsClient()

	// Get the control plane subnet.
	subnetResp, err := subnetsClient.Get(ctx, in.ResourceGroupName, in.VirtualNetworkName, in.ControlPlaneSubnetName, nil)
	if err != nil {
		return nil, err
	}
	controlPlaneSubnet := subnetResp.Subnet

	// Update the control plane subnet to use the NAT gateway.
	controlPlaneSubnet.Properties.NatGateway = &armnetwork.SubResource{
		ID: natGateway.ID,
	}
	subnetsPollerResp, err := subnetsClient.BeginCreateOrUpdate(ctx,
		in.ResourceGroupName,
		in.VirtualNetworkName,
		in.ControlPlaneSubnetName,
		controlPlaneSubnet,
		nil,
	)
	if err != nil {
		return nil, err
	}
	_, err = subnetsPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Get the worker subnet.
	subnetResp, err = subnetsClient.Get(ctx, in.ResourceGroupName, in.VirtualNetworkName, in.WorkerSubnetName, nil)
	if err != nil {
		return nil, err
	}
	workerSubnet := subnetResp.Subnet

	// Update the worker subnet to use the NAT gateway.
	workerSubnet.Properties.NatGateway = &armnetwork.SubResource{
		ID: natGateway.ID,
	}
	subnetsPollerResp, err = subnetsClient.BeginCreateOrUpdate(ctx,
		in.ResourceGroupName,
		in.VirtualNetworkName,
		in.WorkerSubnetName,
		workerSubnet,
		nil,
	)
	if err != nil {
		return nil, err
	}
	_, err = subnetsPollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &natGateway, nil
}
