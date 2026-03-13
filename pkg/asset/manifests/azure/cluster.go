package azure

import (
	"context"
	"fmt"
	"net"

	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	azic "github.com/openshift/installer/pkg/asset/installconfig/azure"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils/cidr"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	mainCIDR := capiutils.CIDRFromInstallConfig(installConfig).String()

	session, err := installConfig.Azure.Session()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Azure session")
	}

	subnets, err := cidr.SplitIntoSubnetsIPv4(mainCIDR, 2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to split CIDR into subnets")
	}

	virtualNetworkAddressPrefixes := []string{mainCIDR}
	controlPlaneAddressPrefixes := []string{subnets[0].String()}
	computeAddressPrefixes := []string{subnets[1].String()}

	// CAPZ expects the capz-system to be created.
	azureNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "capz-system"}}
	azureNamespace.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Namespace"))
	manifests = append(manifests, &asset.RuntimeFile{
		Object: azureNamespace,
		File:   asset.File{Filename: "00_azure-namespace.yaml"},
	})

	resourceGroup := installConfig.Config.Platform.Azure.ClusterResourceGroupName(clusterID.InfraID)
	controlPlaneSubnet := installConfig.Config.Platform.Azure.ControlPlaneSubnetName(clusterID.InfraID)
	computeSubnet := installConfig.Config.Platform.Azure.ComputeSubnetName(clusterID.InfraID)
	networkSecurityGroup := installConfig.Config.Platform.Azure.NetworkSecurityGroupName(clusterID.InfraID)

	source := "*"
	if installConfig.Config.Publish == types.InternalPublishingStrategy {
		source = mainCIDR
	}

	securityGroup := capz.SecurityGroup{
		Name: networkSecurityGroup,
		SecurityGroupClass: capz.SecurityGroupClass{
			SecurityRules: []capz.SecurityRule{
				{
					Name:             "apiserver_in",
					Protocol:         capz.SecurityGroupProtocolTCP,
					Direction:        capz.SecurityRuleDirectionInbound,
					Priority:         101,
					SourcePorts:      ptr.To("*"),
					DestinationPorts: ptr.To("6443"),
					Source:           ptr.To(source),
					Destination:      ptr.To("*"),
					Action:           capz.SecurityRuleActionAllow,
				},
			},
		},
	}

	// Setting ID on the Subnet disables natgw creation. See:
	// https://github.com/kubernetes-sigs/cluster-api-provider-azure/blob/21479a9a4c640b43e0bef028487c522c55605d06/api/v1beta1/azurecluster_default.go#L160
	// CAPZ enables NAT Gateways by default, so we are using this hack to disable
	// nat gateways when we prefer to use load balancers for node egress.
	nodeSubnetID := ""
	if installConfig.Config.Platform.Azure.OutboundType != azure.NatGatewayOutboundType {
		// Because the node subnet does not already exist, we are using an arbitrary value.
		// We could populate this with the proper subnet ID in the case of BYO VNET, but
		// the value currently has no practical effect.
		nodeSubnetID = "UNKNOWN"
	}

	apiServerLB := capz.LoadBalancerSpec{
		Name: fmt.Sprintf("%s-internal", clusterID.InfraID),
		BackendPool: capz.BackendPool{
			Name: fmt.Sprintf("%s-internal", clusterID.InfraID),
		},
		LoadBalancerClassSpec: capz.LoadBalancerClassSpec{
			Type: capz.Internal,
		},
	}

	controlPlaneOutboundLB := &capz.LoadBalancerSpec{
		Name:             clusterID.InfraID,
		FrontendIPsCount: to.Ptr(int32(1)),
	}
	if installConfig.Config.Platform.Azure.OutboundType == azure.UserDefinedRoutingOutboundType {
		controlPlaneOutboundLB = nil
	}

	virtualNetworkID := ""
	lbip := capz.DefaultInternalLBIPAddress
	lbip = getIPWithinCIDR(subnets, lbip)

	if controlPlaneSubnetName := installConfig.Config.Azure.ControlPlaneSubnet; controlPlaneSubnetName != "" {
		controlPlaneSubnet, err := getSubnet(installConfig, clusterID, "controlPlane", controlPlaneSubnetName)
		if err != nil {
			return nil, fmt.Errorf("failed to get control plane subnet: %w", err)
		}
		subnetList, err := getSubnetAddressPrefixes(controlPlaneSubnet)
		if err != nil {
			return nil, fmt.Errorf("failed to get control plane subnet address prefixes: %w", err)
		}
		controlPlaneAddressPrefixes = stringifyAddressPrefixes(subnetList)
		lbip = getIPWithinCIDR(subnetList, lbip)
	}

	if computeSubnetName := installConfig.Config.Azure.ComputeSubnet; computeSubnetName != "" {
		computeSubnet, err := getSubnet(installConfig, clusterID, "compute", computeSubnetName)
		if err != nil {
			return nil, fmt.Errorf("failed to get compute subnet: %w", err)
		}
		subnetList, err := getSubnetAddressPrefixes(computeSubnet)
		if err != nil {
			return nil, fmt.Errorf("failed to get compute subnet address prefixes: %w", err)
		}
		computeAddressPrefixes = stringifyAddressPrefixes(subnetList)
	}

	apiServerLB.FrontendIPs = []capz.FrontendIP{{
		Name: fmt.Sprintf("%s-internal-frontEnd", clusterID.InfraID),
		FrontendIPClass: capz.FrontendIPClass{
			PrivateIPAddress: lbip,
		},
	}}
	if installConfig.Config.Azure.VirtualNetwork != "" {
		virtualNetworkAddressPrefixes = make([]string, 0)

		client, err := installConfig.Azure.Client()
		if err != nil {
			return nil, fmt.Errorf("failed to get azure client: %w", err)
		}
		ctx := context.TODO()
		virtualNetwork, err := client.GetVirtualNetwork(ctx, installConfig.Config.Azure.NetworkResourceGroupName, installConfig.Config.Azure.VirtualNetwork)
		if err != nil {
			return nil, fmt.Errorf("failed to get azure virtual network: %w", err)
		}
		if virtualNetwork != nil {
			virtualNetworkID = *virtualNetwork.ID
		}
		lbip, err := getNextAvailableIPForLoadBalancer(ctx, installConfig, lbip)
		if err != nil {
			return nil, err
		}
		apiServerLB.FrontendIPs[0].FrontendIPClass = capz.FrontendIPClass{
			PrivateIPAddress: lbip,
		}
		if virtualNetwork.AddressSpace != nil && virtualNetwork.AddressSpace.AddressPrefixes != nil {
			virtualNetworkAddressPrefixes = append(virtualNetworkAddressPrefixes, *virtualNetwork.AddressSpace.AddressPrefixes...)
		}
	}

	azureCluster := &capz.AzureCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capz.AzureClusterSpec{
			ResourceGroup: resourceGroup,
			AzureClusterClassSpec: capz.AzureClusterClassSpec{
				SubscriptionID:   session.Credentials.SubscriptionID,
				Location:         installConfig.Config.Azure.Region,
				AdditionalTags:   installConfig.Config.Platform.Azure.UserTags,
				AzureEnvironment: string(installConfig.Azure.CloudName),
				IdentityRef: &corev1.ObjectReference{
					APIVersion: capz.GroupVersion.String(),
					Kind:       "AzureClusterIdentity",
					Name:       clusterID.InfraID,
				},
			},
			NetworkSpec: capz.NetworkSpec{
				NetworkClassSpec: capz.NetworkClassSpec{
					PrivateDNSZoneName: installConfig.Config.ClusterDomain(),
				},
				Vnet: capz.VnetSpec{
					ResourceGroup: installConfig.Config.Azure.NetworkResourceGroupName,
					Name:          installConfig.Config.Azure.VirtualNetwork,
					// The ID is set to virtual network here for existing vnets here. This is to force CAPZ to consider this resource as
					// "not managed" which would prevent the creation of an additional nsg and route table in the network resource group.
					// The ID field is not used for any other purpose in CAPZ except to set the "managed" status.
					// See https://github.com/kubernetes-sigs/cluster-api-provider-azure/blob/main/azure/scope/cluster.go#L585
					// https://github.com/kubernetes-sigs/cluster-api-provider-azure/commit/0f321e4089a3f4dc37f8420bf2ef6762c398c400
					ID: virtualNetworkID,
					VnetClassSpec: capz.VnetClassSpec{
						CIDRBlocks: virtualNetworkAddressPrefixes,
					},
				},
				APIServerLB:            apiServerLB,
				ControlPlaneOutboundLB: controlPlaneOutboundLB,
				Subnets: capz.Subnets{
					{
						SubnetClassSpec: capz.SubnetClassSpec{
							Name:       controlPlaneSubnet,
							Role:       capz.SubnetControlPlane,
							CIDRBlocks: controlPlaneAddressPrefixes,
						},
						SecurityGroup: securityGroup,
					},
					{
						ID: nodeSubnetID,
						SubnetClassSpec: capz.SubnetClassSpec{
							Name:       computeSubnet,
							Role:       capz.SubnetNode,
							CIDRBlocks: computeAddressPrefixes,
						},
						SecurityGroup: securityGroup,
					},
				},
			},
		},
	}
	azureCluster.SetGroupVersionKind(capz.GroupVersion.WithKind("AzureCluster"))
	manifests = append(manifests, &asset.RuntimeFile{
		Object: azureCluster,
		File:   asset.File{Filename: "02_azure-cluster.yaml"},
	})

	azureClientSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID + "-azure-client-secret",
			Namespace: capiutils.Namespace,
		},
		StringData: map[string]string{
			"clientSecret": session.Credentials.ClientSecret,
		},
	}
	azureClientSecret.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Secret"))
	manifests = append(manifests, &asset.RuntimeFile{
		Object: azureClientSecret,
		File:   asset.File{Filename: "01_azure-client-secret.yaml"},
	})

	id := &capz.AzureClusterIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capz.AzureClusterIdentitySpec{
			Type:              capz.ServicePrincipal,
			AllowedNamespaces: &capz.AllowedNamespaces{}, // Allow all namespaces.
			ClientID:          session.Credentials.ClientID,
			ClientSecret: corev1.SecretReference{
				Name:      azureClientSecret.Name,
				Namespace: azureClientSecret.Namespace,
			},
			TenantID: session.Credentials.TenantID,
		},
	}
	if session.AuthType == azic.ManagedIdentityAuth {
		id.Spec.Type = capz.UserAssignedMSI
		id.Spec.ClientSecret = corev1.SecretReference{}
	}
	id.SetGroupVersionKind(capz.GroupVersion.WithKind("AzureClusterIdentity"))
	manifests = append(manifests, &asset.RuntimeFile{
		Object: id,
		File:   asset.File{Filename: "01_azure-cluster-controller-identity-default.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: capz.GroupVersion.String(),
				Kind:       "AzureCluster",
				Name:       azureCluster.Name,
				Namespace:  azureCluster.Namespace,
			},
		},
	}, nil
}

func getSubnet(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, subnetType, subnetName string) (*aznetwork.Subnet, error) {
	var subnet *aznetwork.Subnet

	azClient, err := installConfig.Azure.Client()
	if err != nil {
		return nil, fmt.Errorf("failed to get azure client: %w", err)
	}
	ctx := context.TODO()

	if subnetType == "controlPlane" {
		subnet, err = azClient.GetControlPlaneSubnet(ctx,
			installConfig.Config.Azure.NetworkResourceGroupName,
			installConfig.Config.Azure.VirtualNetwork,
			subnetName,
		)
	} else if subnetType == "compute" {
		subnet, err = azClient.GetComputeSubnet(ctx,
			installConfig.Config.Azure.NetworkResourceGroupName,
			installConfig.Config.Azure.VirtualNetwork,
			subnetName,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get subnet: %w", err)
	}
	if subnet == nil {
		return nil, fmt.Errorf("failed to get subnet")
	}
	if subnet.AddressPrefixes == nil && subnet.AddressPrefix == nil {
		return nil, fmt.Errorf("failed to get subnet addresses: %w", err)
	}
	return subnet, nil
}

func getSubnetAddressPrefixes(subnet *aznetwork.Subnet) ([]*net.IPNet, error) {
	subnetList := []*net.IPNet{}
	if subnet.AddressPrefixes != nil {
		for _, sub := range *subnet.AddressPrefixes {
			_, ipnet, err := net.ParseCIDR(sub)
			if err != nil {
				return subnetList, fmt.Errorf("failed to get translate azure subnet addresses: %w", err)
			}
			subnetList = append(subnetList, ipnet)
		}
	}
	if subnet.AddressPrefix != nil {
		_, ipnet, err := net.ParseCIDR(*subnet.AddressPrefix)
		if err != nil {
			return subnetList, fmt.Errorf("failed to get translate azure subnet address prefix: %w", err)
		}
		subnetList = append(subnetList, ipnet)
	}

	return subnetList, nil
}

func stringifyAddressPrefixes(addressPrefixes []*net.IPNet) []string {
	strAddressPrefixes := []string{}
	for _, addressPrefix := range addressPrefixes {
		strAddressPrefixes = append(strAddressPrefixes, addressPrefix.String())
	}
	return strAddressPrefixes
}

func getIPWithinCIDR(subnets []*net.IPNet, ip string) string {
	if subnets == nil || ip == "" {
		return ""
	}
	// Check if default lbip is within control plane network.
	// If not in control plane network, assign the first non-reserved IP in the CIDR to lbip.
	for _, subnet := range subnets {
		if subnet == nil {
			continue
		}
		if subnet.Contains(net.ParseIP(ip)) {
			return ip
		}
	}
	ipSubnets := make(net.IP, len(subnets[0].IP))
	copy(ipSubnets, subnets[0].IP)
	// Since the first 4 IP of the subnets are usually reserved[1], pick the next one that's available in the CIDR.
	// [1] - https://learn.microsoft.com/en-us/azure/virtual-network/ip-services/private-ip-addresses#allocation-method
	ipSubnets[len(ipSubnets)-1] += 4
	return ipSubnets.String()
}

func getNextAvailableIPForLoadBalancer(ctx context.Context, installConfig *installconfig.InstallConfig, lbip string) (string, error) {
	client, err := installConfig.Azure.Client()
	if err != nil {
		return "", fmt.Errorf("failed to get azure client: %w", err)
	}
	networkResourceGroupName := installConfig.Config.Azure.NetworkResourceGroupName
	virtualNetworkName := installConfig.Config.Azure.VirtualNetwork
	machineCidr := installConfig.Config.MachineNetwork
	if cpSubnet := installConfig.Config.Azure.ControlPlaneSubnet; cpSubnet != "" {
		controlPlane, err := client.GetControlPlaneSubnet(ctx, networkResourceGroupName, virtualNetworkName, cpSubnet)
		if err != nil {
			return "", fmt.Errorf("failed to get control plane subnet: %w", err)
		}
		if controlPlane.AddressPrefix == nil && controlPlane.AddressPrefixes == nil {
			return "", fmt.Errorf("failed to get control plane subnet addresses: %w", err)
		}
		prefixes := []*ipnet.IPNet{}
		if controlPlane.AddressPrefixes != nil {
			for _, sub := range *controlPlane.AddressPrefixes {
				ipnet, err := ipnet.ParseCIDR(sub)
				if err != nil {
					return "", fmt.Errorf("failed to get translate azure control plane subnet addresses: %w", err)
				}
				prefixes = append(prefixes, ipnet)
			}
		}

		if controlPlane.AddressPrefix != nil {
			ipnet, err := ipnet.ParseCIDR(*controlPlane.AddressPrefix)
			if err != nil {
				return "", fmt.Errorf("failed to get translate azure control plane subnet address prefix: %w", err)
			}
			prefixes = append(prefixes, ipnet)
		}
		cidrRange := []types.MachineNetworkEntry{}
		for _, prefix := range prefixes {
			if prefix != nil {
				cidrRange = append(cidrRange, types.MachineNetworkEntry{CIDR: *prefix})
			}
		}
		machineCidr = cidrRange
	}
	availableIP, err := client.CheckIPAddressAvailability(ctx, networkResourceGroupName, virtualNetworkName, lbip)
	if err != nil {
		return "", fmt.Errorf("failed to get azure ip availability: %w", err)
	}
	if availableIP == nil {
		return "", errors.New("failed to get available IP in given machine network: this error may be caused by lack of necessary permissions")
	}
	ipAvail := *availableIP
	if ipAvail.Available != nil && *ipAvail.Available {
		for _, cidrRange := range machineCidr {
			if cidrRange.CIDR.Contains(net.ParseIP(lbip)) {
				return lbip, nil
			}
		}
	}
	if ipAvail.AvailableIPAddresses == nil || len(*ipAvail.AvailableIPAddresses) == 0 {
		return "", fmt.Errorf("failed to get an available IP in given virtual network for LB: %w", err)
	}
	for _, ip := range *ipAvail.AvailableIPAddresses {
		for _, cidrRange := range machineCidr {
			if cidrRange.CIDR.Contains(net.ParseIP(lbip)) {
				return ip, nil
			}
		}
	}
	return "", fmt.Errorf("failed to get available IP in given machine network: %w", err)
}
