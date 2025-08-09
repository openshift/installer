package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	azic "github.com/openshift/installer/pkg/asset/installconfig/azure"
	capzash "github.com/openshift/installer/pkg/asset/manifests/azure/stack/v1beta1"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils/cidr"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	logrus.Debugf("XXX: GenerateClusterAssets")
	resourceGroup := installConfig.Config.Platform.Azure.ClusterResourceGroupName(clusterID.InfraID)
	controlPlaneSubnet := installConfig.Config.Platform.Azure.ControlPlaneSubnetName(clusterID.InfraID)
	computeSubnet := installConfig.Config.Platform.Azure.ComputeSubnetName(clusterID.InfraID)
	networkSecurityGroup := installConfig.Config.Platform.Azure.NetworkSecurityGroupName(clusterID.InfraID)
	securityGroup := capz.SecurityGroup{Name: networkSecurityGroup}

	session, err := installConfig.Azure.Session()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Azure session")
	}

	// CAPZ expects the capz-system to be created.
	azureNamespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "capz-system"}}
	azureNamespace.SetGroupVersionKind(corev1.SchemeGroupVersion.WithKind("Namespace"))

	manifests := []*asset.RuntimeFile{}
	manifests = append(manifests, &asset.RuntimeFile{
		Object: azureNamespace,
		File:   asset.File{Filename: "00_azure-namespace.yaml"},
	})

	// Setting ID on the Subnet disables natgw creation. See:
	// https://github.com/kubernetes-sigs/cluster-api-provider-azure/blob/21479a9a4c640b43e0bef028487c522c55605d06/api/v1beta1/azurecluster_default.go#L160
	// CAPZ enables NAT Gateways by default, so we are using this hack to disable
	// nat gateways when we prefer to use load balancers for node egress.
	nodeSubnetID := ""
	switch installConfig.Config.Platform.Azure.OutboundType {
	// Because the node subnet does not already exist, we are using an arbitrary value.
	// We could populate this with the proper subnet ID in the case of BYO VNET, but
	// the value currently has no practical effect.
	case azure.LoadbalancerOutboundType:
		fallthrough
	case azure.UserDefinedRoutingOutboundType:
		nodeSubnetID = "UNKNOWN"
	}

	//apiServerLBFrontendConfigName := azure.InternalFrontendIPv4ConfigName
	apiServerLBFrontendConfigName := fmt.Sprintf("%s-internal-frontEnd", clusterID.InfraID)
	//apiServerBackendPoolName := azure.InternalBackendAddressPoolIPv4Name
	apiServerBackendPoolName := fmt.Sprintf("%s-internal", clusterID.InfraID)
	if installConfig.Config.Networking.IsIPv6() && !installConfig.Config.Networking.IsDualStack() {
		//apiServerLBFrontendConfigName = azure.InternalFrontendIPv6ConfigName
		apiServerLBFrontendConfigName = fmt.Sprintf("%s-internal-frontEnd-ipv6", clusterID.InfraID)
		//apiServerBackendPoolName = azure.InternalBackendAddressPoolIPv6Name
		apiServerBackendPoolName = fmt.Sprintf("%s-internal-ipv6", clusterID.InfraID)
	}
	apiServerLB := capz.LoadBalancerSpec{
		Name: fmt.Sprintf("%s-internal", clusterID.InfraID),
		BackendPool: capz.BackendPool{
			Name: apiServerBackendPoolName,
		},
		LoadBalancerClassSpec: capz.LoadBalancerClassSpec{
			Type: capz.Internal,
		},
	}

	//controlPlaneBackendPoolName := azure.OutboundBackendAddressPoolIPv4Name
	controlPlaneBackendPoolName := fmt.Sprintf("%s-outbound-lb-outboundBackendPool", clusterID.InfraID)
	if installConfig.Config.Networking.IsIPv6() && !installConfig.Config.Networking.IsDualStack() {
		//controlPlaneBackendPoolName = azure.OutboundBackendAddressPoolIPv6Name
		controlPlaneBackendPoolName = fmt.Sprintf("%s-outbound-lb-outboundBackendPool-ipv6", clusterID.InfraID)
	}
	controlPlaneOutboundLB := &capz.LoadBalancerSpec{
		Name:             clusterID.InfraID,
		FrontendIPsCount: to.Ptr(int32(1)),
		BackendPool: capz.BackendPool{
			Name: controlPlaneBackendPoolName,
		},
	}

	if installConfig.Config.Platform.Azure.OutboundType == azure.UserDefinedRoutingOutboundType {
		controlPlaneOutboundLB = nil
	}

	// XXX: keep this in for now to keep the LB code happy, but rewrite
	// later
	subnets, err := cidr.SplitIntoSubnetsIPv4(capiutils.CIDRFromInstallConfig(installConfig).String(), 2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to split CIDR into subnets")
	}
	virtualNetworkID := ""
	lbip := capz.DefaultInternalLBIPAddress
	lbip = getIPWithinCIDR(subnets, lbip)

	if controlPlaneSub := installConfig.Config.Azure.ControlPlaneSubnet; controlPlaneSub != "" {
		client, err := installConfig.Azure.Client()
		if err != nil {
			return nil, fmt.Errorf("failed to get azure client: %w", err)
		}
		ctx := context.TODO()
		controlPlaneSubnet, err := client.GetControlPlaneSubnet(ctx, installConfig.Config.Azure.NetworkResourceGroupName, installConfig.Config.Azure.VirtualNetwork, controlPlaneSub)
		if err != nil || controlPlaneSubnet == nil {
			return nil, fmt.Errorf("failed to get azure control plane subnet: %w", err)
		} else if controlPlaneSubnet.AddressPrefixes == nil && controlPlaneSubnet.AddressPrefix == nil {
			return nil, fmt.Errorf("failed to get azure control plane subnet addresses: %w", err)
		}
		subnetList := []*net.IPNet{}
		if controlPlaneSubnet.AddressPrefixes != nil {
			for _, sub := range *controlPlaneSubnet.AddressPrefixes {
				_, ipnet, err := net.ParseCIDR(sub)
				if err != nil {
					return nil, fmt.Errorf("failed to get translate azure control plane subnet addresses: %w", err)
				}
				subnetList = append(subnetList, ipnet)
			}
		}

		if controlPlaneSubnet.AddressPrefix != nil {
			_, ipnet, err := net.ParseCIDR(*controlPlaneSubnet.AddressPrefix)
			if err != nil {
				return nil, fmt.Errorf("failed to get translate azure control plane subnet address prefix: %w", err)
			}
			subnetList = append(subnetList, ipnet)
		}
		lbip = getIPWithinCIDR(subnetList, lbip)
	}

	apiServerLB.FrontendIPs = []capz.FrontendIP{{
		Name: apiServerLBFrontendConfigName,
		FrontendIPClass: capz.FrontendIPClass{
			PrivateIPAddress: lbip,
		},
	}}

	if installConfig.Config.Azure.VirtualNetwork != "" {
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
	}

	azEnv := string(installConfig.Azure.CloudName)

	computeSubnetSpec := capz.SubnetSpec{
		ID: nodeSubnetID,
		SubnetClassSpec: capz.SubnetClassSpec{
			Name: computeSubnet,
			Role: capz.SubnetNode,
			/*
				CIDRBlocks: []string{
					subnets[1].String(),
				},
			*/
		},
		//SecurityGroup: securityGroup,
	}

	if installConfig.Config.Azure.OutboundType == azure.NATGatewaySingleZoneOutboundType {
		computeSubnetSpec.NatGateway = capz.NatGateway{
			NatGatewayClassSpec: capz.NatGatewayClassSpec{Name: fmt.Sprintf("%s-natgw", clusterID.InfraID)},
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
				AzureEnvironment: azEnv,
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
				},
				APIServerLB:            &apiServerLB,
				ControlPlaneOutboundLB: controlPlaneOutboundLB,
				Subnets: capz.Subnets{
					{
						SubnetClassSpec: capz.SubnetClassSpec{
							Name: controlPlaneSubnet,
							Role: capz.SubnetControlPlane,
						},
						//SecurityGroup: securityGroup,
					},
					computeSubnetSpec,
				},
			},
		},
	}

	///////////////////////////////////////////////////////////////////////////////////////////
	// If a single CIDR is in the install-config:
	//   - Split it:
	//     - Use first half for control plane subnet
	//     - Use second half for compute subnet
	//
	// If more than a single CIDR is specified:
	//   - If multiple of same protocol:
	//     - Use first for control plane subnet
	//     - Use second for compute subnet
	//     - Create additional subnets for any additional
	//
	// If dual stack:
	//   - Same as above, but bind both protocols to each subnet
	//
	// Example #1:
	//
	// machineNetwork:
	// - cidr: 10.0.0.0/17     <- control plane subnet
	// - cidr: 10.0.128.0/17   <- compute subnet
	// - cidr: 192.168.0.0/16  <- subnet #x
	// - cidr: 172.16.0.0/16   <- subnet #y
	// - cidr: ffd0::/48       <- splits into ffd0:0:0:1::/64 and ffd0:0:0:2::/64
	//                            where ffd0:0:0:1::/64 is control plane subnet
	//                            and ffd0:0:0:2::/64 is compute subnet
	//                            ffd0:0:0:1::/64 is added to 10.0.0.0/17 subnet cidr blocks
	//                            ffd0:0:0:2::/64 is added to 10.0.128.0/17 subnet cidr blocks
	// Example #2:
	//
	// machineNetwork:
	// - cidr: 10.0.0.0/16
	//
	// ^ Splits into 10.0.0.0/17 and 10.0.128.0/17
	//   10.0.0.0/17 for the control plane subnet
	//   10.0.128.0/17 for the compute subnet
	//
	// Example #3:
	//
	// machineNetwork:
	// - cidr: ffd0::/48
	//
	// ^ Splits into ffd0:0:0:1::/64 and ffd0:0:0:2::/64
	//   ffd0:0:0:1::/64 is for the control plane subnet
	//   ffd0:0:0:2::/64 is for the compute subnet
	//
	// Example #4:
	//
	// machineNetwork:
	// - cidr: 10.0.0.0/16
	// - cidr: ffd0::/48
	//
	// ^ 10.0.0.0/16 splits into 10.0.0.0/17 and 10.0.128.0/17
	//   ffd0::/48 splits into ffd0:0:0:1::/64 and ffd0:0:0:2::/64
	//   10.0.0.0/17 and ffd0:0:0:1::/64 and the control plane cidr blocks
	//   10.0.128.0/17 and ffd0:0:0:1::/64 are the compute cidr blocks
	///////////////////////////////////////////////////////////////////////////////////////////
	CIDRs := capiutils.CIDRsFromInstallConfig(installConfig)

	// XXXX
	// XXX: Get back to this and rewrite using install config networking
	// IsIP* functions
	// XXXX

	/*
	   if installConfig.Networking.IsDualStack() {
	           p.DualStack = true
	   } else if installConfig.Networking.IsIPv4() {
	           p.FrontendIPConfigName = frontendIPv4ConfigName
	   } else if installConfig.Networking.IsIPv6() {
	           p.FrontendIPConfigName = frontendIPv6ConfigName
	   }
	*/
	//installConfig.Config.Networking.Is

	// Split CIDRs into IPv4 and IPv6 CIDRs
	var ipv4CIDRs, ipv6CIDRs []ipnet.IPNet
	for _, CIDR := range CIDRs {
		switch len(CIDR.IP) {
		case net.IPv4len:
			ipv4CIDRs = append(ipv4CIDRs, CIDR)
		case net.IPv6len:
			ipv6CIDRs = append(ipv6CIDRs, CIDR)
		}
	}

	// Split IPv4 CIDRs into IPv4 subnets
	var ipv4Subnets []*net.IPNet
	logrus.Debugf("XXX: len(ipv4CIDRs)=%d", len(ipv4CIDRs))
	switch len(ipv4CIDRs) {
	case 1:
		ipv4Subnets, err = cidr.SplitIntoSubnetsIPv4(ipv4CIDRs[0].String(), 2)
		if err != nil {
			return nil, errors.Wrap(err, "failed to split IPv4 CIDR into subnets")
		}
	default:
		for _, ipv4Cidr := range ipv4CIDRs {
			ipv4Subnets = append(ipv4Subnets, &net.IPNet{
				IP:   ipv4Cidr.IP,
				Mask: ipv4Cidr.Mask,
			})
		}
	}

	// Split IPv6 CIDRs into IPv6 subnets
	var ipv6Subnets []*net.IPNet
	logrus.Debugf("XXX: len(ipv6CIDRs)=%d", len(ipv6CIDRs))
	switch len(ipv6CIDRs) {
	case 1:
		ipv6Subnets, err = cidr.SplitIntoSubnetsIPv6(ipv6CIDRs[0].String(), 2)
		if err != nil {
			return nil, errors.Wrap(err, "failed to split IPv6 CIDR into subnets")
		}
	default:
		for _, ipv6Cidr := range ipv6CIDRs {
			ipv6Subnets = append(ipv6Subnets, &net.IPNet{
				IP:   ipv6Cidr.IP,
				Mask: ipv6Cidr.Mask,
			})
		}
	}

	securityRule := capz.SecurityRule{
		Name:             "apiserver_in",
		Protocol:         capz.SecurityGroupProtocolTCP,
		Direction:        capz.SecurityRuleDirectionInbound,
		Priority:         100,
		SourcePorts:      ptr.To("*"),
		DestinationPorts: ptr.To("6443"),
		Source:           ptr.To("*"),
		Destination:      ptr.To("*"),
		Action:           capz.SecurityRuleActionAllow,
	}

	// If we are using Internal publishing, we need a security rule for each CIDR
	var securityRules []capz.SecurityRule
	var securityRulePriority int32 = 100
	if len(ipv4CIDRs) > 0 && installConfig.Config.Publish == types.InternalPublishingStrategy {
		for i, ipv4CIDR := range ipv4CIDRs {
			ipv4SecurityRule := securityRule
			ipv4SecurityRule.Name = fmt.Sprintf("apiserver_in_ipv4_%02d", i)
			ipv4SecurityRule.Source = to.Ptr(ipv4CIDR.String())
			ipv4SecurityRule.Priority = securityRulePriority
			securityRules = append(securityRules, ipv4SecurityRule)
			securityRulePriority += 10
		}
	}
	if len(ipv6CIDRs) > 0 && installConfig.Config.Publish == types.InternalPublishingStrategy {
		for i, ipv6CIDR := range ipv6CIDRs {
			ipv6SecurityRule := securityRule
			ipv6SecurityRule.Name = fmt.Sprintf("apiserver_in_ipv6_%02d", i)
			ipv6SecurityRule.Source = to.Ptr(ipv6CIDR.String())
			ipv6SecurityRule.Priority = securityRulePriority
			securityRules = append(securityRules, ipv6SecurityRule)
			securityRulePriority += 10
		}
	}
	if len(securityRules) == 0 {
		securityRules = append(securityRules, securityRule)
	}
	securityGroup.SecurityGroupClass.SecurityRules = securityRules
	logrus.Debugf("XXX: len(securityRules)=%d", len(securityRules))

	subnetsPtr := &azureCluster.Spec.NetworkSpec.Subnets

	// Create additional subnets if necessary
	subnetsCount := max(2, len(ipv4Subnets), len(ipv6Subnets))
	logrus.Debugf("XXX: subnetsCount=%d", subnetsCount)
	if subnetsCount > 2 {
		for i := 2; i < subnetsCount; i++ {
			newSubnet := capz.SubnetSpec{
				SubnetClassSpec: capz.SubnetClassSpec{
					Name: fmt.Sprintf("subnet_%02d", i),
					Role: capz.SubnetNode,
				},
			}
			*subnetsPtr = append(*subnetsPtr, newSubnet)
		}
	}

	var vnetCIDRblocks, controlPlaneCIDRblocks, computeCIDRblocks []string

	// Set the VNet CIDR blockCIDs
	for _, CIDR := range CIDRs {
		logrus.Debugf("XXX: CIDR=%s", CIDR.String())
		vnetCIDRblocks = append(vnetCIDRblocks, CIDR.String())
	}

	// XXX: 10.1.0.0 is node CIDR, where does 10.2.0.0 come from?
	vnetCIDRblocks = append(vnetCIDRblocks, "10.1.0.0/16")
	vnetCIDRblocks = append(vnetCIDRblocks, "10.2.0.0/16")

	logrus.Debugf("XXX: len(ipv4Subnets)=%d", len(ipv4Subnets))
	logrus.Debugf("XXX: len(ipv6Subnets)=%d", len(ipv6Subnets))

	// Set control plane and compute network
	for i, ipv4Subnet := range ipv4Subnets {
		switch i {
		case 0:
			controlPlaneCIDRblocks = append(controlPlaneCIDRblocks, ipv4Subnet.String())
		case 1:
			computeCIDRblocks = append(computeCIDRblocks, ipv4Subnet.String())
		}
	}

	// Set control plane and compute network IPv6 CIDR blocks
	for i, ipv6Subnet := range ipv6Subnets {
		switch i {
		case 0:
			controlPlaneCIDRblocks = append(controlPlaneCIDRblocks, ipv6Subnet.String())
		case 1:
			computeCIDRblocks = append(computeCIDRblocks, ipv6Subnet.String())
		}
	}

	azureCluster.Spec.NetworkSpec.Vnet.VnetClassSpec.CIDRBlocks = vnetCIDRblocks
	azureCluster.Spec.NetworkSpec.Subnets[0].SubnetClassSpec.CIDRBlocks = controlPlaneCIDRblocks
	azureCluster.Spec.NetworkSpec.Subnets[1].SubnetClassSpec.CIDRBlocks = computeCIDRblocks

	logrus.Debugf("XXX: control plane subnet CIDR blocks = %+v", azureCluster.Spec.NetworkSpec.Subnets[0].SubnetClassSpec.CIDRBlocks)

	// Set the security group for each subnet
	for i := 0; i < subnetsCount; i++ {
		azureCluster.Spec.NetworkSpec.Subnets[i].SecurityGroup = securityGroup
	}

	// We are maintaining a fork of CAPZ for azurestack. The only API difference
	// is the ARMEndpoint field, so we can use the CAPZ cluster object, and if
	// running on ASH convert to the fork API, and add the field.
	var cluster client.Object
	if !strings.EqualFold(azEnv, string(azure.StackCloud)) {
		azureCluster.SetGroupVersionKind(capz.GroupVersion.WithKind("AzureCluster"))
		cluster = azureCluster
	} else {
		var ashCluster capzash.AzureCluster
		if err := deepCopy(azureCluster, &ashCluster); err != nil {
			return nil, fmt.Errorf("failed to convert azureCluster to azure-stack cluster: %w", err)
		}
		ashCluster.Spec.ARMEndpoint = session.Environment.ServiceManagementEndpoint
		ashCluster.SetGroupVersionKind(capzash.GroupVersion.WithKind("AzureCluster"))
		cluster = &ashCluster
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: cluster,
		File:   asset.File{Filename: "02_azure-cluster.yaml"},
	})

	id := &capz.AzureClusterIdentity{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capz.AzureClusterIdentitySpec{
			AllowedNamespaces: &capz.AllowedNamespaces{}, // Allow all namespaces.
			ClientID:          session.Credentials.ClientID,
			TenantID:          session.Credentials.TenantID,
		},
	}

	switch session.AuthType {
	case azic.ManagedIdentityAuth:
		id.Spec.Type = capz.UserAssignedMSI
	case azic.ClientSecretAuth:
		id.Spec.Type = capz.ServicePrincipal
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

		id.Spec.ClientSecret = corev1.SecretReference{
			Name:      azureClientSecret.Name,
			Namespace: azureClientSecret.Namespace,
		}
	case azic.ClientCertificateAuth:
		id.Spec.Type = capz.ServicePrincipalCertificate
		id.Spec.CertPath = session.Credentials.ClientCertificatePath
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
	// AzureStack does not support the call to CheckIPAddressAvailability.
	if installConfig.Azure.CloudName == azure.StackCloud {
		cidr := machineCidr[0]
		if cidr.CIDR.Contains(net.IP(lbip)) {
			return lbip, nil
		}
		ipSubnets := cidr.CIDR.IP
		ipSubnets[len(ipSubnets)-1] += 4
		return ipSubnets.String(), nil
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
		return "", fmt.Errorf("failed to get an available IP in given virtual network for LB: this error may be caused by lack of necessary permissions")
	}
	for _, ip := range *ipAvail.AvailableIPAddresses {
		for _, cidrRange := range machineCidr {
			if cidrRange.CIDR.Contains(net.ParseIP(lbip)) {
				return ip, nil
			}
		}
	}
	return "", fmt.Errorf("failed to get an IP that's available and in the given machine network: this error may be caused by lack of necessary permissions")
}

func deepCopy(src, dst interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}

func apiserverSecurityRule(name string, priority int32) capz.SecurityRule {
	return capz.SecurityRule{
		Name:             name,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Direction:        capz.SecurityRuleDirectionInbound,
		Priority:         priority,
		SourcePorts:      ptr.To("*"),
		DestinationPorts: ptr.To("6443"),
		Source:           ptr.To("*"),
		Destination:      ptr.To("*"),
		Action:           capz.SecurityRuleActionAllow,
	}
}
