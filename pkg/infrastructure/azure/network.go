package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	aztypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
)

type lbInput struct {
	loadBalancerName       string
	infraID                string
	region                 string
	resourceGroupName      string
	subscriptionID         string
	frontendIPConfigName   string
	backendAddressPoolName string
	idPrefix               string
	stackType              aztypes.StackType
	lbClient               *armnetwork.LoadBalancersClient
	tags                   map[string]*string
}

type lbRuleInput struct {
	idPrefix               string
	loadBalancerName       string
	probeName              string
	ruleName               string
	frontendIPConfigName   string
	backendAddressPoolName string
}

type vnetInput struct {
	resourceGroupName    string
	virtualNetworkName   string
	networkClientFactory *armnetwork.ClientFactory
}

type subnetInput struct {
	resourceGroupName    string
	virtualNetworkName   string
	subnetName           string
	networkClientFactory *armnetwork.ClientFactory
}

type pipInput struct {
	infraID           string
	name              string
	region            string
	resourceGroupName string
	stackType         aztypes.StackType
	pipClient         *armnetwork.PublicIPAddressesClient
	tags              map[string]*string
}

type vmInput struct {
	infraID              string
	resourceGroupName    string
	ids                  []string
	backendAddressPools  []*armnetwork.BackendAddressPool
	vmClient             *armcompute.VirtualMachinesClient
	nicClient            *armnetwork.InterfacesClient
	networkClientFactory *armnetwork.ClientFactory
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

func mcsProbe() *armnetwork.Probe {
	return &armnetwork.Probe{
		Name: to.Ptr("sint-probe"),
		Properties: &armnetwork.ProbePropertiesFormat{
			Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
			Port:              to.Ptr[int32](22623),
			IntervalInSeconds: to.Ptr[int32](5),
			NumberOfProbes:    to.Ptr[int32](2),
			RequestPath:       to.Ptr("/healthz"),
		},
	}
}

func mcsRule(in *lbRuleInput) *armnetwork.LoadBalancingRule {
	return &armnetwork.LoadBalancingRule{
		Name: to.Ptr(in.ruleName),
		Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
			Protocol:             to.Ptr(armnetwork.TransportProtocolTCP),
			FrontendPort:         to.Ptr[int32](22623),
			BackendPort:          to.Ptr[int32](22623),
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
				ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", in.idPrefix, in.loadBalancerName, in.probeName)),
			},
		},
	}
}

func apiProbe() *armnetwork.Probe {
	return &armnetwork.Probe{
		//Name: to.Ptr("api-probe"),
		Name: to.Ptr("HTTPSProbe"),
		Properties: &armnetwork.ProbePropertiesFormat{
			Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
			Port:              to.Ptr[int32](6443),
			IntervalInSeconds: to.Ptr[int32](5),
			NumberOfProbes:    to.Ptr[int32](2),
			RequestPath:       to.Ptr("/readyz"),
		},
	}
}

func apiRule(in *lbRuleInput) *armnetwork.LoadBalancingRule {
	return &armnetwork.LoadBalancingRule{
		Name: to.Ptr(in.ruleName),
		Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
			Protocol:             to.Ptr(armnetwork.TransportProtocolTCP),
			FrontendPort:         to.Ptr[int32](6443),
			BackendPort:          to.Ptr[int32](6443),
			IdleTimeoutInMinutes: to.Ptr[int32](30),
			DisableOutboundSnat:  to.Ptr(true),
			EnableFloatingIP:     to.Ptr(false),
			LoadDistribution:     to.Ptr(armnetwork.LoadDistributionDefault),
			FrontendIPConfiguration: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", in.idPrefix, in.loadBalancerName, in.frontendIPConfigName)),
			},
			BackendAddressPool: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", in.idPrefix, in.loadBalancerName, in.backendAddressPoolName)),
			},
			Probe: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", in.idPrefix, in.loadBalancerName, in.probeName)),
			},
		},
	}
}

func outboundRule(in *lbRuleInput) *armnetwork.OutboundRule {
	return &armnetwork.OutboundRule{
		Name: to.Ptr(in.ruleName),
		Properties: &armnetwork.OutboundRulePropertiesFormat{
			Protocol:               to.Ptr(armnetwork.LoadBalancerOutboundRuleProtocolAll),
			AllocatedOutboundPorts: to.Ptr[int32](1024), // configurable ??
			IdleTimeoutInMinutes:   nil,
			FrontendIPConfigurations: []*armnetwork.SubResource{
				&armnetwork.SubResource{
					ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", in.idPrefix, in.loadBalancerName, in.frontendIPConfigName)),
				},
			},
			BackendAddressPool: &armnetwork.SubResource{
				ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", in.idPrefix, in.loadBalancerName, in.backendAddressPoolName)),
			},
		},
	}
}

func addProbeToLoadBalancer(ctx context.Context, probe *armnetwork.Probe, in *lbInput) (*armnetwork.LoadBalancer, error) {
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	lb := lbResp.LoadBalancer

	lb.Properties.Probes = append(lb.Properties.Probes, probe)

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx, in.resourceGroupName, in.loadBalancerName, lb, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func addLoadBalancingRuleToLoadBalancer(ctx context.Context, rule *armnetwork.LoadBalancingRule, in *lbInput) (*armnetwork.LoadBalancer, error) {
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	lb := lbResp.LoadBalancer

	lb.Properties.LoadBalancingRules = append(lb.Properties.LoadBalancingRules, rule)

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx, in.resourceGroupName, in.loadBalancerName, lb, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func addOutboundRuleToLoadBalancer(ctx context.Context, rule *armnetwork.OutboundRule, in *lbInput) (*armnetwork.LoadBalancer, error) {
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	lb := lbResp.LoadBalancer

	lb.Properties.OutboundRules = append(lb.Properties.OutboundRules, rule)

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx, in.resourceGroupName, in.loadBalancerName, lb, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func getSubnet(ctx context.Context, in *subnetInput) (*armnetwork.Subnet, error) {
	subnetsClient := in.networkClientFactory.NewSubnetsClient()
	subnetResp, err := subnetsClient.Get(ctx, in.resourceGroupName, in.virtualNetworkName, in.subnetName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get subnet: %w", err)
	}
	return &subnetResp.Subnet, nil
}

func getVirtualNetwork(ctx context.Context, in *vnetInput) (*armnetwork.VirtualNetwork, error) {
	virtualNetworksClient := in.networkClientFactory.NewVirtualNetworksClient()
	virtualNetworksResp, err := virtualNetworksClient.Get(ctx, in.resourceGroupName, in.virtualNetworkName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual network: %w", err)
	}
	return &virtualNetworksResp.VirtualNetwork, nil
}

func addFrontendIPConfigurationToLoadBalancer(ctx context.Context, frontendIPConfiguration *armnetwork.FrontendIPConfiguration, in *lbInput) (*armnetwork.LoadBalancer, error) {
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	lb := lbResp.LoadBalancer

	// We need a vnet & subnet
	lb.Properties.FrontendIPConfigurations = append(lb.Properties.FrontendIPConfigurations, frontendIPConfiguration)

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx, in.resourceGroupName, in.loadBalancerName, lb, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func addBackendAddressPoolToLoadBalancer(ctx context.Context, backendAddressPool *armnetwork.BackendAddressPool, in *lbInput) (*armnetwork.LoadBalancer, error) {
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	lb := lbResp.LoadBalancer

	lb.Properties.BackendAddressPools = append(lb.Properties.BackendAddressPools, backendAddressPool)

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx, in.resourceGroupName, in.loadBalancerName, lb, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot update load balancer: %w", err)
	}

	resp, err := pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &resp.LoadBalancer, nil
}

func createPublicIP(ctx context.Context, in *pipInput) (*armnetwork.PublicIPAddress, error) {
	publicIPAddressVersion := armnetwork.IPVersionIPv4
	if in.stackType == aztypes.StackTypeIPv6 {
		publicIPAddressVersion = armnetwork.IPVersionIPv6
	}

	pollerResp, err := in.pipClient.BeginCreateOrUpdate(
		ctx,
		in.resourceGroupName,
		in.name,
		armnetwork.PublicIPAddress{
			Name:     to.Ptr(in.name),
			Location: to.Ptr(in.region),
			SKU: &armnetwork.PublicIPAddressSKU{
				Name: to.Ptr(armnetwork.PublicIPAddressSKUNameStandard),
				Tier: to.Ptr(armnetwork.PublicIPAddressSKUTierRegional),
			},
			Properties: &armnetwork.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   to.Ptr(publicIPAddressVersion),
				PublicIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodStatic),
				DNSSettings: &armnetwork.PublicIPAddressDNSSettings{
					DomainNameLabel: to.Ptr(in.name),
					//DomainNameLabel: to.Ptr(in.infraID),
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
		in.resourceGroupName,
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

	// XXX: Fix this, pass these in
	probeName := "api-probe"
	ruleName := "api-v4"
	privateIPAddressVersion := armnetwork.IPVersionIPv4

	if in.stackType == aztypes.StackTypeIPv6 {
		privateIPAddressVersion = armnetwork.IPVersionIPv6
		probeName = "api-probe-v6"
		ruleName = "api-v6"
	}

	// Get the CAPI-created outbound load balancer so we can modify it.
	extLB, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
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
				PrivateIPAddressVersion:   to.Ptr(privateIPAddressVersion),
				PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
				PublicIPAddress:           pip,
			},
		})
	extLB.Properties.BackendAddressPools = append(extLB.Properties.BackendAddressPools,
		&armnetwork.BackendAddressPool{
			Name: &in.backendAddressPoolName,
		})

	extLB.Properties.Probes = append(extLB.Properties.Probes,
		&armnetwork.Probe{
			Name: &probeName,
			Properties: &armnetwork.ProbePropertiesFormat{
				Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
				Port:              to.Ptr[int32](6443),
				IntervalInSeconds: to.Ptr[int32](5),
				NumberOfProbes:    to.Ptr[int32](2),
				RequestPath:       to.Ptr("/readyz"),
			},
		})
	extLB.Properties.LoadBalancingRules = append(extLB.Properties.LoadBalancingRules,
		&armnetwork.LoadBalancingRule{
			Name: to.Ptr(ruleName),
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
		})

	pollerResp, err := in.lbClient.BeginCreateOrUpdate(ctx,
		in.resourceGroupName,
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
				Probes:                   extLB.Properties.Probes,
				LoadBalancingRules:       extLB.Properties.LoadBalancingRules,
				OutboundRules:            extLB.Properties.OutboundRules,
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
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
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
		in.resourceGroupName,
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

func deepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

func associateVMToBackendPools(ctx context.Context, in vmInput) error {
	var ipv4BackendAddressPools, ipv6BackendAddressPools []*armnetwork.BackendAddressPool
	var backendAddressPools [4]armnetwork.BackendAddressPool
	var backendAddressPoolName string

	lbBackendAddressPoolsClient := in.networkClientFactory.NewLoadBalancerBackendAddressPoolsClient()
	loadBalancerName := in.infraID

	// Get the IPv4 backend address pools
	backendAddressPoolName = in.infraID
	resp, err := lbBackendAddressPoolsClient.Get(ctx,
		in.resourceGroupName,
		loadBalancerName,
		backendAddressPoolName,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to get backend address pool %s: %", backendAddressPoolName, err)
	}
	deepCopy(&resp.BackendAddressPool, &backendAddressPools[0])
	ipv4BackendAddressPools = append(ipv4BackendAddressPools, &backendAddressPools[0])

	backendAddressPoolName = fmt.Sprintf("%s-outbound-lb-outboundBackendPool", in.infraID)
	resp, err = lbBackendAddressPoolsClient.Get(ctx,
		in.resourceGroupName,
		loadBalancerName,
		backendAddressPoolName,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to get backend address pool %s: %", backendAddressPoolName, err)
	}
	deepCopy(&resp.BackendAddressPool, &backendAddressPools[1])
	ipv4BackendAddressPools = append(ipv4BackendAddressPools, &backendAddressPools[1])

	// Get the IPv6 backend address pools
	backendAddressPoolName = fmt.Sprintf("%s-ipv6", in.infraID)
	resp, err = lbBackendAddressPoolsClient.Get(ctx,
		in.resourceGroupName,
		loadBalancerName,
		backendAddressPoolName,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to get backend address pool %s: %", backendAddressPoolName, err)
	}
	deepCopy(&resp.BackendAddressPool, &backendAddressPools[2])
	ipv6BackendAddressPools = append(ipv6BackendAddressPools, &backendAddressPools[2])

	loadBalancerName = fmt.Sprintf("%s-internal", in.infraID)
	backendAddressPoolName = fmt.Sprintf("%s-internal-ipv6", in.infraID)
	resp, err = lbBackendAddressPoolsClient.Get(ctx,
		in.resourceGroupName,
		loadBalancerName,
		backendAddressPoolName,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to get backend address pool %s: %", backendAddressPoolName, err)
	}
	deepCopy(&resp.BackendAddressPool, &backendAddressPools[3])
	ipv6BackendAddressPools = append(ipv6BackendAddressPools, &backendAddressPools[3])

	for _, id := range in.ids {
		vmName := path.Base(id)
		vm, err := in.vmClient.Get(ctx, in.resourceGroupName, vmName, nil)
		if err != nil {
			return fmt.Errorf("failed to get vm %s: %w", vmName, err)
		}

		if nics := vm.Properties.NetworkProfile.NetworkInterfaces; len(nics) == 1 {
			nicRef := nics[0]
			nicName := path.Base(*nicRef.ID)
			nic, err := in.nicClient.Get(ctx, in.resourceGroupName, nicName, nil)
			if err != nil {
				return fmt.Errorf("failed to get nic for vm %s: %w", vmName, err)
			}
			for _, ipConfig := range nic.Properties.IPConfigurations {
				logrus.Debugf("XXX: nicName=%s ipConfig.Name=%s vmName=%s", nicName, *ipConfig.Name, vmName)
				if *ipConfig.Name == "pipConfig" {
					//ipConfig.Properties.LoadBalancerBackendAddressPools = []*armnetwork.BackendAddressPool{}
					for _, pool := range ipv4BackendAddressPools {
						ipConfig.Properties.LoadBalancerBackendAddressPools = append(
							ipConfig.Properties.LoadBalancerBackendAddressPools,
							[]*armnetwork.BackendAddressPool{{
								ID: pool.ID,
							}}...,
						)
					}

				} else if *ipConfig.Name == "ipConfigv6" {
					//ipConfig.Properties.LoadBalancerBackendAddressPools = []*armnetwork.BackendAddressPool{}
					for _, pool := range ipv6BackendAddressPools {
						ipConfig.Properties.LoadBalancerBackendAddressPools = append(
							ipConfig.Properties.LoadBalancerBackendAddressPools,
							[]*armnetwork.BackendAddressPool{{
								ID: pool.ID,
							}}...,
						)
					}
				}
			}
			pollerResp, err := in.nicClient.BeginCreateOrUpdate(ctx, in.resourceGroupName, nicName, nic.Interface, nil)
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
		if *ipConfig.Name == "ipConfigv6" {
			continue
		}
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
