package azure

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	aztypes "github.com/openshift/installer/pkg/types/azure"
	"github.com/sirupsen/logrus"
	"k8s.io/utils/ptr"
)

type lbInput struct {
	subscriptionID          string
	infraID                 string
	region                  string
	resourceGroupName       string
	loadBalancerName        string
	frontendIPConfigName    string
	backendAddressPoolName  string
	outboundAddressPoolName string
	idPrefix                string
	stackType               aztypes.StackType
	networkClientFactory    *armnetwork.ClientFactory
	lbClient                *armnetwork.LoadBalancersClient
	frontendIPConfiguration *armnetwork.FrontendIPConfiguration
	backendAddressPool      *armnetwork.BackendAddressPool
	pip                     *armnetwork.PublicIPAddress
	tags                    map[string]*string
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
	loadBalancerName     string
	vmIDs                []string
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
	backendIPConfigID    string
	frontendIPConfigName string
	inboundNatRuleID     string
	inboundNatRuleName   string
	inboundNatRulePort   int32
	networkClientFactory *armnetwork.ClientFactory
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
					//DomainNameLabel: to.Ptr(in.infraID),
					DomainNameLabel: to.Ptr(in.name),
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

func createAPILoadBalancer(ctx context.Context, pipv4 *armnetwork.PublicIPAddress, pipv6 *armnetwork.PublicIPAddress, in *lbInput) (*armnetwork.LoadBalancer, error) {
	probeName := "api-probe"
	lbRuleName := "api-ipv4"
	if in.stackType == aztypes.StackTypeIPv6 {
		lbRuleName = "api-ipv6"
	}

	frontendIPConfigurations := []*armnetwork.FrontendIPConfiguration{
		{
			Name: &in.frontendIPConfigName,
			Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
				PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
				PublicIPAddress:           pipv4,
			},
		},
	}
	if pipv6 != nil {
		frontendIPConfigurations = append(frontendIPConfigurations, &armnetwork.FrontendIPConfiguration{
			Name: &in.frontendIPConfigName,
			Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
				PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
				PublicIPAddress:           pipv4,
			},
		})
	}

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
				FrontendIPConfigurations: frontendIPConfigurations,
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
						Name: to.Ptr(lbRuleName),
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
		Name: to.Ptr("api-probe"),
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
			AllocatedOutboundPorts: to.Ptr[int32](1024),
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

func getLoadBalancer(ctx context.Context, in *lbInput) (*armnetwork.LoadBalancer, error) {
	loadBalancersClient := in.networkClientFactory.NewLoadBalancersClient()
	loadBalancersResp, err := loadBalancersClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	return &loadBalancersResp.LoadBalancer, nil
}

func newFrontendIPConfigurationIPv6(name string, subnet *armnetwork.Subnet) *armnetwork.FrontendIPConfiguration {
	return &armnetwork.FrontendIPConfiguration{
		Name: to.Ptr(name),
		Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
			PrivateIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv6),
			PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
			Subnet:                    subnet,
		},
	}
}

func addFrontendIPConfigurationToLoadBalancer(ctx context.Context, in *lbInput) (*armnetwork.LoadBalancer, error) {
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	lb := lbResp.LoadBalancer

	// We need a vnet & subnet
	lb.Properties.FrontendIPConfigurations = append(lb.Properties.FrontendIPConfigurations,
		in.frontendIPConfiguration,
	)

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

func newBackendAddressPool(name string, virtualNetwork *armnetwork.VirtualNetwork) *armnetwork.BackendAddressPool {
	return &armnetwork.BackendAddressPool{
		Name: to.Ptr(name),
		/*
			Properties: &armnetwork.BackendAddressPoolPropertiesFormat{
				VirtualNetwork: &armnetwork.SubResource{
					ID: virtualNetwork.ID,
				},
			},
		*/
	}
}

func addBackendAddressPoolToLoadBalancer(ctx context.Context, in *lbInput) (*armnetwork.LoadBalancer, error) {
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get load balancer: %w", err)
	}
	lb := lbResp.LoadBalancer

	lb.Properties.BackendAddressPools = append(lb.Properties.BackendAddressPools, in.backendAddressPool)

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

func updateOutboundLoadBalancerToAPILoadBalancerOld(ctx context.Context, in *lbInput) (*armnetwork.LoadBalancer, error) {
	probeName := "api-probe"
	lbRuleName := "api-ipv4"
/*
	if in.stackType == aztypes.StackTypeIPv6 {
		lbRuleName = "api-ipv6"
	}
*/
	logrus.Debugf("XXX: updateOutboundLoadBalancerToAPILoadBalancer: lbRuleName=%s", lbRuleName)

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
	if in.stackType == aztypes.StackTypeIPv6 {
		extLB.Properties.FrontendIPConfigurations = append(extLB.Properties.FrontendIPConfigurations,
			&armnetwork.FrontendIPConfiguration{
				Name: &in.frontendIPConfigName,
				Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
					PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
					PublicIPAddress:           in.pip,
				},
			})
		extLB.Properties.BackendAddressPools = append(extLB.Properties.BackendAddressPools,
			&armnetwork.BackendAddressPool{
				Name: &in.backendAddressPoolName,
			})
	}

	if len(extLB.Properties.Probes) == 0 {
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
	}

	extLB.Properties.LoadBalancingRules = append(extLB.Properties.LoadBalancingRules,
		&armnetwork.LoadBalancingRule{
			Name: to.Ptr(lbRuleName),
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
	/*
		outboundRuleName := "OutboundNATAllProtocolsIPv6"
		if len(extLB.Properties.OutboundRules) == 1 && in.stackType == aztypes.StackTypeIPv6 {
			extLB.Properties.BackendAddressPools = append(extLB.Properties.BackendAddressPools,
				&armnetwork.BackendAddressPool{
					Name: to.Ptr(in.outboundAddressPoolName),
				})
			extLB.Properties.OutboundRules = append(extLB.Properties.OutboundRules,
				&armnetwork.OutboundRule{
					Name: to.Ptr(outboundRuleName),
					Properties: &armnetwork.OutboundRulePropertiesFormat{
						Protocol:               to.Ptr(armnetwork.LoadBalancerOutboundRuleProtocolAll),
						EnableTCPReset:         to.Ptr(true),
						AllocatedOutboundPorts: to.Ptr(1024),
						FrontendIPConfigurations: []*armnetwork.SubResource{
							{
								ID: to.Ptr(fmt.Sprintf("/%s/%s/frontendIPConfigurations/%s", in.idPrefix, in.loadBalancerName, in.frontendIPConfigName)),
							},
						},
						BackendAddressPool: &armnetwork.SubResource{
							ID: to.Ptr(fmt.Sprintf("/%s/%s/backendAddressPools/%s", in.idPrefix, in.loadBalancerName, in.outboundAddressPoolName)),
						},
					},
				})
		}
	*/

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

func updateOutboundLoadBalancerToAPILoadBalancer(ctx context.Context, pip *armnetwork.PublicIPAddress, in *lbInput) (*armnetwork.LoadBalancer, error) {
        probeName := "api-probe"

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
                                PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
                                PublicIPAddress:           pip,
                        },
                })
        extLB.Properties.BackendAddressPools = append(extLB.Properties.BackendAddressPools,
                &armnetwork.BackendAddressPool{
                        Name: &in.backendAddressPoolName,
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


func updateInternalLoadBalancerIPv6(ctx context.Context, in *lbInput) (*armnetwork.LoadBalancer, error) {
	//probeName := "sint-probe"
	lbRuleName := "sint-ipv6"

	// Get the CAPI-created internal load balancer so we can modify it.
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get internal load balancer: %w", err)
	}
	intLB := lbResp.LoadBalancer

	// XXX fix this
	//mcsProbe := in.probe
	mcsProbe := mcsProbe()
	mcsProbeName := *mcsProbe.Name

	existingFrontEndIPConfig := intLB.Properties.FrontendIPConfigurations
	if len(existingFrontEndIPConfig) == 0 {
		return nil, fmt.Errorf("could not get frontEndIPConfig for internal LB %s", *intLB.Name)
	}
	controlPlaneSubnet := existingFrontEndIPConfig[0].Properties.Subnet

	intLB.Properties.FrontendIPConfigurations = append(intLB.Properties.FrontendIPConfigurations,
		&armnetwork.FrontendIPConfiguration{
			Name: &in.frontendIPConfigName,
			Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
				PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
				PrivateIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv6),
				Subnet:                    controlPlaneSubnet,
			},
		})

	intLB.Properties.BackendAddressPools = append(intLB.Properties.BackendAddressPools,
		&armnetwork.BackendAddressPool{
			Name: &in.backendAddressPoolName,
		})

	mcsRule := &armnetwork.LoadBalancingRule{
		Name: to.Ptr(lbRuleName),
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
				ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", in.idPrefix, in.loadBalancerName, mcsProbeName)),
			},
		},
	}

	// XXX: This is so hacky, do this for now for testing this code
	var found bool
	for _, probe := range intLB.Properties.Probes {
		if *probe.Name == *mcsProbe.Name {
			found = true
			break
		}
	}
	if found == false {
		intLB.Properties.Probes = append(intLB.Properties.Probes, mcsProbe)
	}

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


func updateInternalLoadBalancerOld(ctx context.Context, in *lbInput) (*armnetwork.LoadBalancer, error) {
	probeName := "sint-probe"
	lbRuleName := "sint-ipv4"
	if in.stackType == aztypes.StackTypeIPv6 {
		lbRuleName = "sint-ipv6"
	}

	logrus.Debugf("XXX: updateInternalLoadBalancer: frontendIPConfigName=%s backendAddressPoolName=%s lbRuleName=%s",
		in.frontendIPConfigName, in.backendAddressPoolName, lbRuleName)

	// Get the CAPI-created internal load balancer so we can modify it.
	lbResp, err := in.lbClient.Get(ctx, in.resourceGroupName, in.loadBalancerName, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get internal load balancer: %w", err)
	}
	intLB := lbResp.LoadBalancer

	var subnet *armnetwork.Subnet

	var found bool = false
	for _, frontend := range intLB.Properties.FrontendIPConfigurations {
		subnet = frontend.Properties.Subnet
		if *frontend.Name == in.frontendIPConfigName {
			found = true
			break
		}
	}

	// Create frontend and backend
	if found == false {
		intLB.Properties.FrontendIPConfigurations = append(intLB.Properties.FrontendIPConfigurations,
			&armnetwork.FrontendIPConfiguration{
				Name: &in.frontendIPConfigName,
				Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
					// XXX: This should probably be static...
					// for DNS for internal LB - handle
					// later
					PrivateIPAddressVersion:   to.Ptr(armnetwork.IPVersionIPv6),
					PrivateIPAllocationMethod: to.Ptr(armnetwork.IPAllocationMethodDynamic),
					//PrivateIPAddress:          to.Ptr("fc00:0:0:100"),
					Subnet: subnet,
					//PublicIPAddress:           in.pip,
				},
			})
		intLB.Properties.BackendAddressPools = append(intLB.Properties.BackendAddressPools,
			&armnetwork.BackendAddressPool{
				Name: &in.backendAddressPoolName,
			})
	}

	/*
		existingFrontEndIPConfig := intLB.Properties.FrontendIPConfigurations
		if len(existingFrontEndIPConfig) == 0 {
			return nil, fmt.Errorf("could not get frontEndIPConfig for internal LB %s", *intLB.Name)
		}
		existingFrontEndIPConfigName := *(existingFrontEndIPConfig[0].Name)
	*/

	mcsProbe := &armnetwork.Probe{
		Name: to.Ptr(probeName),
		Properties: &armnetwork.ProbePropertiesFormat{
			Protocol:          to.Ptr(armnetwork.ProbeProtocolHTTPS),
			Port:              to.Ptr[int32](22623),
			IntervalInSeconds: to.Ptr[int32](5),
			NumberOfProbes:    to.Ptr[int32](2),
			RequestPath:       to.Ptr("/healthz"),
		},
	}

	found = false
	for _, probe := range intLB.Properties.Probes {
		if *probe.Name == *mcsProbe.Name {
			found = true
			break
		}
	}
	if found == false {
		intLB.Properties.Probes = append(intLB.Properties.Probes, mcsProbe)
	}

	mcsRule := &armnetwork.LoadBalancingRule{
		Name: to.Ptr(lbRuleName),
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
				ID: to.Ptr(fmt.Sprintf("/%s/%s/probes/%s", in.idPrefix, in.loadBalancerName, *mcsProbe.Name)),
			},
		},
	}

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

// sony fw900
func dumpBackendAddressPool(pool *armnetwork.BackendAddressPool) {
	if pool == nil {
		return
	}
	if pool.Properties == nil {
		return
	}

	if pool.Name != nil {
		logrus.Debugf("XXX: Dumping backend pool %s", *pool.Name)
	}
	for _, ipConfig := range pool.Properties.BackendIPConfigurations {
		if ipConfig == nil {
			logrus.Debugf("XXX: ipConfig is NULL")
			continue
		}

		if ipConfig.Name != nil {
			logrus.Debugf("XXX: ipConfig Name=%s", *ipConfig.Name)
		}
		if ipConfig.ID != nil {
			logrus.Debugf("XXX: ipConfig ID=%s", *ipConfig.ID)
		}

		//logrus.Debugf("XXX: ipConfig Name=%s ID=%s", *ipConfig.Name, *ipConfig.ID)
		if ipConfig.Properties != nil && ipConfig.Properties.PrivateIPAddress != nil {
			logrus.Debugf("XXX: ipConfig Name=%s PrivateIPAddress=%s", *ipConfig.Name, *ipConfig.Properties.PrivateIPAddress)
		}
	}
}

func associateVMToIPv4BackendPool(ctx context.Context, in vmInput) error {
	for _, id := range in.vmIDs {
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
			for _, ipconfig := range nic.Properties.IPConfigurations {
				if *ipconfig.Name == "ipConfigv6" {
					continue		
				}
				ipconfig.Properties.LoadBalancerBackendAddressPools = append(ipconfig.Properties.LoadBalancerBackendAddressPools, in.backendAddressPools...)
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

func associateVMToBackendPool(ctx context.Context, in vmInput) error {
	logrus.Debugf("XXX: associateVMToBackendPool: in.backendAddressPools=%+v", in.backendAddressPools)

	networkInterfacesClient := in.networkClientFactory.NewInterfacesClient()
	for _, vmID := range in.vmIDs {
		vmName := path.Base(vmID)
		vm, err := in.vmClient.Get(ctx, in.resourceGroupName, vmName, nil)
		if err != nil {
			return fmt.Errorf("failed to get vm %s: %w", vmName, err)
		}

		if nics := vm.Properties.NetworkProfile.NetworkInterfaces; len(nics) == 1 {
			nicRef := nics[0]
			nicName := path.Base(*nicRef.ID)

			nic, err := networkInterfacesClient.Get(ctx, in.resourceGroupName, nicName, nil)
			if err != nil {
				return fmt.Errorf("failed to get nic for vm %s: %w", vmName, err)
			}
			logrus.Debugf("XXX: *nic.ID=%s", *nic.ID)

			// We need to CREATE the ipv6 backend and associate here...
			for _, ipConfig := range nic.Properties.IPConfigurations {
				if *ipConfig.Name == "pipConfig" {
					for _, pool := range in.backendAddressPools {
						if !strings.HasSuffix(*pool.Name, "-ipv6") {
							ipConfig.Properties.LoadBalancerBackendAddressPools = append(
								ipConfig.Properties.LoadBalancerBackendAddressPools,
								[]*armnetwork.BackendAddressPool{
									pool,
								}...,
							)
						}
					}

				} else if *ipConfig.Name == "ipConfigv6" {
					for _, pool := range in.backendAddressPools {
						if strings.HasSuffix(*pool.Name, "-ipv6") {
							ipConfig.Properties.LoadBalancerBackendAddressPools = append(
								ipConfig.Properties.LoadBalancerBackendAddressPools,
								[]*armnetwork.BackendAddressPool{
									pool,
								}...,
							)
						}
					}
				}
			}

			//, err := json.Marshal(user)

			pollerResp, err := networkInterfacesClient.BeginCreateOrUpdate(ctx, in.resourceGroupName, nicName, nic.Interface, nil)
			if err != nil {
				return fmt.Errorf("failed to update nic for %s: %w", vmName, err)
			}
			_, err = pollerResp.PollUntilDone(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to poll update nic for vm %s: %w", vmName, err)
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
				/*
					BackendIPConfiguration: &armnetwork.InterfaceIPConfiguration{
						ID: to.Ptr(in.backendIPConfigID),
					},
				*/
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

	frontendIPConfigurationsClient := in.networkClientFactory.NewLoadBalancerFrontendIPConfigurationsClient()
	frontendIPConfigurationsResp, err := frontendIPConfigurationsClient.Get(ctx,
		in.resourceGroupName,
		in.loadBalancerName,
		in.frontendIPConfigName,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get frontend IP configuration: %w", err)
	}

	frontendIPConfiguration := frontendIPConfigurationsResp.FrontendIPConfiguration
	logrus.Debugf("XXX: NAT frontendIPConfiguration.ID=%s Name=%s", *frontendIPConfiguration.ID, *frontendIPConfiguration.Name)

	for i, ipConfig := range bootstrapInterface.Properties.IPConfigurations {
		logrus.Debugf("XXX: NAT ipConfig.ID=%s Name=%s=", *ipConfig.ID, *ipConfig.Name)
		if *ipConfig.ID == *frontendIPConfiguration.ID {
			logrus.Debugf("XXX: MATCH! ipConfig.ID=%s Name=%s", *ipConfig.ID, *ipConfig.Name)
			ipConfig.Properties.LoadBalancerInboundNatRules = append(ipConfig.Properties.LoadBalancerInboundNatRules,
				&inboundNatRule,
			)
			bootstrapInterface.Properties.IPConfigurations[i] = ipConfig
			break
		}
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
