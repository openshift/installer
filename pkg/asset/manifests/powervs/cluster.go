package powervs

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capibmv2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta2"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, bucket string, object string) (*capiutils.GenerateClusterAssetsOutput, error) {
	var (
		manifests          []*asset.RuntimeFile
		network            string
		dhcpSubnetCIDR     string
		vpcNameOrID        string
		vpcStruct          *vpcv1.VPC
		vpcRegion          string
		cosName            string
		cosRegion          string
		imageName          string
		bucketName         string
		transitGatewayName string
		client             *powervsconfig.Client
		err                error
		powerVSCluster     *capibm.IBMPowerVSCluster
		powerVSImage       *capibmv2.IBMPowerVSImage
	)

	defer func() {
		logrus.Debugf("GenerateClusterAssets: installConfig = %+v, clusterID = %v, bucket = %v, object = %v", installConfig, *clusterID, bucket, object)
		logrus.Debugf("GenerateClusterAssets: ic.ObjectMeta = %+v", installConfig.Config.ObjectMeta.Name)
		logrus.Debugf("GenerateClusterAssets: installConfig.Config.PowerVS = %+v", *installConfig.Config.PowerVS)
		logrus.Debugf("GenerateClusterAssets: installConfig.Config.Platform.PowerVS.VPC = %v", installConfig.Config.Platform.PowerVS.VPC)
		logrus.Debugf("GenerateClusterAssets: vpcNameOrID = %v", vpcNameOrID)
		logrus.Debugf("GenerateClusterAssets: vpcRegion = %v", vpcRegion)
		logrus.Debugf("GenerateClusterAssets: installConfig.Config.Platform.PowerVS.TransitGateway = %v", installConfig.Config.Platform.PowerVS.TransitGateway)
		logrus.Debugf("GenerateClusterAssets: transitGatewayName = %v", transitGatewayName)
		logrus.Debugf("GenerateClusterAssets: cosName = %v", cosName)
		logrus.Debugf("GenerateClusterAssets: cosRegion = %v", cosRegion)
		logrus.Debugf("GenerateClusterAssets: imageName = %v", imageName)
		logrus.Debugf("GenerateClusterAssets: bucketName = %v", bucketName)
		logrus.Debugf("GenerateClusterAssets: powerVSCluster.Spec.ControlPlaneEndpoint.Host = %v", powerVSCluster.Spec.ControlPlaneEndpoint.Host)
	}()

	manifests = []*asset.RuntimeFile{}

	network = fmt.Sprintf("%s-network", clusterID.InfraID)

	logrus.Debugf("GenerateClusterAssets: len MachineNetwork = %d", len(installConfig.Config.Networking.MachineNetwork))
	dhcpSubnetCIDR = installConfig.Config.Networking.MachineNetwork[0].CIDR.String()
	if numNetworks := len(installConfig.Config.Networking.MachineNetwork); numNetworks > 1 {
		logrus.Infof("Warning: %d machineNetwork found! Ignoring all except the first.", numNetworks)
	}
	logrus.Debugf("GenerateClusterAssets: dhcpSubnetCIDR = %s", dhcpSubnetCIDR)

	// Build v1beta3 WorkspaceSource
	var workspaceSource capibm.WorkspaceSource
	if installConfig.Config.PowerVS.ServiceInstanceGUID == "" {
		workspaceName := fmt.Sprintf("%s-power-iaas", clusterID.InfraID)
		workspaceSource = capibm.WorkspaceSource{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.WorkspaceProvisionConfig{
				Name: workspaceName,
			},
		}
	} else {
		workspaceSource = capibm.WorkspaceSource{
			Type: capibm.SourceTypeReference,
			Reference: capibm.ResourceIdentifier{
				ID: installConfig.Config.PowerVS.ServiceInstanceGUID,
			},
		}
	}

	client, err = powervsconfig.NewClientWithEndpoints(installConfig.Config.Platform.PowerVS.ServiceEndpoints)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	vpcRegion = installConfig.Config.Platform.PowerVS.VPCRegion
	if vpcRegion == "" {
		if vpcRegion, err = powervstypes.VPCRegionForPowerVSRegion(installConfig.Config.PowerVS.Region); err != nil {
			return nil, fmt.Errorf("unable to derive vpcRegion from region: %s %w", installConfig.Config.PowerVS.Region, err)
		}
	}

	// Build v1beta3 VPCSource
	// The VPC specified can be either:
	// 1) blank - CAPI will create one for us.
	// 2) an id of an existing VPC.
	// 3) a name of an existing VPC.
	vpcNameOrID = installConfig.Config.Platform.PowerVS.VPC
	var vpcSource capibm.VPCSource
	if vpcStruct, err = client.GetVPCByID(ctx, vpcNameOrID, vpcRegion); err == nil {
		// #2 — ID
		logrus.Debugf("GenerateClusterAssets: PowerVS.VPC ID is valid")
		vpcSource = capibm.VPCSource{
			Type:   capibm.SourceTypeReference,
			Region: vpcRegion,
			Reference: capibm.ResourceIdentifier{
				ID: installConfig.Config.Platform.PowerVS.VPC,
			},
		}
	} else if vpcStruct, err = client.GetVPCByName(ctx, vpcNameOrID); err == nil {
		// #3 — Name
		logrus.Debugf("GenerateClusterAssets: PowerVS.VPC Name is valid")
		vpcSource = capibm.VPCSource{
			Type:   capibm.SourceTypeReference,
			Region: vpcRegion,
			Reference: capibm.ResourceIdentifier{
				Name: vpcNameOrID,
			},
		}
	} else {
		if vpcNameOrID == "" {
			// #1 — provision
			logrus.Debugf("GenerateClusterAssets: PowerVS.VPC is empty")
			vpcNameOrID = fmt.Sprintf("vpc-%s", clusterID.InfraID)
			vpcStruct = nil
			vpcSource = capibm.VPCSource{
				Type:   capibm.SourceTypeProvision,
				Region: vpcRegion,
				Provision: capibm.VPCProvision{
					Name: vpcNameOrID,
				},
			}
		} else {
			return nil, fmt.Errorf("generateClusterAssets could not handle vpc")
		}
	}

	// Build VPC Security Groups (unchanged — uses powervs/v1beta2 compatible type via getVPCSecurityGroups)
	vpcSecurityGroups := getVPCSecurityGroups(clusterID.InfraID)

	// Build v1beta3 TransitGatewaySource
	// The Transit Gateway can be either:
	// 1) blank - CAPI will create one for us.
	// 2) an id of an existing TG.
	// 3) a name of an existing TG.
	transitGatewayName = installConfig.Config.Platform.PowerVS.TransitGateway
	var tgSource capibm.TransitGatewaySource
	if err = client.TransitGatewayIDValid(ctx, transitGatewayName); err == nil {
		logrus.Debugf("GenerateClusterAssets: TG ID is valid")
		tgSource = capibm.TransitGatewaySource{
			Type: capibm.SourceTypeReference,
			Reference: capibm.ResourceIdentifier{
				ID: installConfig.Config.Platform.PowerVS.TransitGateway,
			},
		}
	} else {
		if transitGatewayName == "" {
			logrus.Debugf("GenerateClusterAssets: PowerVS.TransitGateway is empty")
			transitGatewayName = fmt.Sprintf("%s-tg", clusterID.InfraID)
		}
		tgSource = capibm.TransitGatewaySource{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.TransitGatewayProvision{
				Name: transitGatewayName,
			},
		}
	}

	// Build v1beta3 COSInstanceSource
	if installConfig.Config.PowerVS.COSInstanceCRN != "" {
		logrus.Debugf("Using existing COS instance from CRN: %s", installConfig.Config.PowerVS.COSInstanceCRN)
		cosInstance, cosErr := client.GetResourceInstanceByCRN(ctx, installConfig.Config.PowerVS.COSInstanceCRN)
		if cosErr != nil {
			return nil, fmt.Errorf("failed to get COS instance details from CRN %s: %w", installConfig.Config.PowerVS.COSInstanceCRN, cosErr)
		}
		if cosInstance == nil || cosInstance.Name == nil {
			return nil, fmt.Errorf("COS instance not found or has no name for CRN: %s", installConfig.Config.PowerVS.COSInstanceCRN)
		}
		cosName = *cosInstance.Name
		logrus.Infof("Using existing COS instance: %s", cosName)
	} else {
		cosName = fmt.Sprintf("%s-cos", clusterID.InfraID)
		logrus.Debugf("Will create new COS instance: %s", cosName)
	}

	if cosRegion, err = powervstypes.COSRegionForPowerVSRegion(installConfig.Config.PowerVS.Region); err != nil {
		return nil, fmt.Errorf("unable to derive cosRegion from region: %s %w", installConfig.Config.PowerVS.Region, err)
	}

	if vpcRegion != cosRegion {
		logrus.Debugf("GenerateClusterAssets: vpcRegion(%s) is different than cosRegion(%s), overriding bucket name", vpcRegion, cosRegion)
		bucket = fmt.Sprintf("rhcos-powervs-images-%s", cosRegion)
	}

	bucketName = fmt.Sprintf("%s-bootstrap-ign", clusterID.InfraID)
	imageName = fmt.Sprintf("rhcos-%s", clusterID.InfraID)

	var cosSource capibm.COSInstanceSource
	if installConfig.Config.PowerVS.COSInstanceCRN != "" {
		cosSource = capibm.COSInstanceSource{
			Type:         capibm.SourceTypeReference,
			BucketName:   bucketName,
			BucketRegion: cosRegion,
			Reference: capibm.ResourceIdentifier{
				Name: cosName,
			},
		}
	} else {
		cosSource = capibm.COSInstanceSource{
			Type:         capibm.SourceTypeProvision,
			BucketName:   bucketName,
			BucketRegion: cosRegion,
			Provision: capibm.COSInstanceProvision{
				Name: cosName,
			},
		}
	}

	// useDHCPServer is true when the user has explicitly requested the legacy DHCPServer mode.
	// Default (empty string) means DHCPSubnet.
	useDHCPServer := installConfig.Config.Platform.PowerVS.NetworkType == powervstypes.NetworkTypeDHCPServer

	// Build v1beta3 NetworkSource.
	//
	// DHCPSubnet (default): creates a plain PowerVS network with DHCP natively enabled.
	// No separate DHCP server VM is created. Enables the full new networking model.
	//
	// DHCPServer (legacy): creates a dedicated DHCP Service VM. The DHCP Server handles its
	// own SNAT/egress so the route-mode NLB, Routing Table, and Public Gateway are omitted.
	var networkSource capibm.NetworkSource
	if useDHCPServer {
		networkSource = capibm.NetworkSource{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.NetworkProvisionConfig{
				Type: capibm.NetworkProvisionTypeDHCPServer,
				DHCPServer: capibm.DHCPServer{
					Name: network,
					CIDR: dhcpSubnetCIDR,
					Snat: capibm.DHCPSnatPolicyEnabled,
				},
			},
		}
	} else {
		networkSource = capibm.NetworkSource{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.NetworkProvisionConfig{
				Type: capibm.NetworkProvisionTypeDHCPSubnet,
				DHCPSubnet: capibm.DHCPSubnet{
					Name: network,
					CIDR: dhcpSubnetCIDR,
				},
			},
		}
	}

	// Build v1beta3 Load Balancers.
	//
	// DHCPSubnet mode: the public LB uses the "network-fixed" profile with RouteMode enabled
	// so it acts as a Layer-3 next-hop router. The CAPIBM webhook enforces that
	// RouteMode=Enabled requires Profile=network-fixed.
	//
	// DHCPServer mode: the public LB uses the default "application" profile (no RouteMode)
	// because the DHCP Server's own SNAT provides egress — a route-mode NLB is not needed.
	//
	// The private/internal LB always uses the "application" profile.
	lbName := fmt.Sprintf("%s-loadbalancer", clusterID.InfraID)
	lbIntName := fmt.Sprintf("%s-loadbalancer-int", clusterID.InfraID)

	publicLBProvision := capibm.LoadBalancerProvision{
		Name: lbName,
		Type: capibm.LoadBalancerTypePublic,
		AdditionalListeners: []capibm.AdditionalListener{
			{
				Port:     22,
				Protocol: capibm.LoadBalancerListenerProtocolTCP,
			},
		},
	}
	if !useDHCPServer {
		// DHCPSubnet path: route-mode NLB required for VPC Routing Table egress.
		publicLBProvision.Profile = capibm.PowerVSLoadBalancerProfileNetworkFixed
		publicLBProvision.RouteMode = capibm.PowerVSLoadBalancerRouteModeEnabled
	}

	loadBalancers := []capibm.LoadBalancerSource{
		{
			Type:      capibm.SourceTypeProvision,
			Provision: publicLBProvision,
		},
		{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.LoadBalancerProvision{
				Name: lbIntName,
				Type: capibm.LoadBalancerTypePrivate,
				AdditionalListeners: []capibm.AdditionalListener{
					{
						Port:     22623,
						Protocol: capibm.LoadBalancerListenerProtocolTCP,
					},
				},
			},
		},
	}

	// ResourceGroup
	rgSource := capibm.ResourceGroupSource{
		Type: capibm.SourceTypeReference,
		Reference: capibm.ResourceIdentifier{
			Name: installConfig.Config.Platform.PowerVS.PowerVSResourceGroup,
		},
	}

	// Build v1beta3 VPCPublicGateways and VPCRoutingTables.
	//
	// These are only needed for the DHCPSubnet networking model. When using DHCPServer,
	// the DHCP server VM handles SNAT and provides egress — no VPC routing is required.
	var (
		publicGateways []capibm.VPCPublicGateway
		routingTables  []capibm.VPCRoutingTableSource
	)
	if !useDHCPServer {
		// Derive the VPC zone from the PowerVS zone for the Public Gateway and Routing Table.
		vpcZone, vpcZoneErr := powervstypes.VPCZoneForPowerVSZone(installConfig.Config.Platform.PowerVS.Zone)
		if vpcZoneErr != nil {
			return nil, fmt.Errorf("unable to derive VPC zone for Public Gateway from PowerVS zone %s: %w",
				installConfig.Config.Platform.PowerVS.Zone, vpcZoneErr)
		}

		// Public Gateway — outbound internet access for VPC subnets.
		pgwName := fmt.Sprintf("%s-pgw", clusterID.InfraID)
		publicGateways = []capibm.VPCPublicGateway{
			{
				Type: capibm.SourceTypeProvision,
				Zone: vpcZone,
				Provision: capibm.VPCPublicGatewayProvision{
					Name: pgwName,
				},
			},
		}

		// Routing Table — advertises 0.0.0.0/0 via the route-mode NLB to the Transit Gateway.
		// The nextHop is intentionally left empty: CAPIBM auto-populates it from
		// LoadBalancerStatus.PrivateIPs[0] of the public route-mode load balancer.
		rtName := fmt.Sprintf("%s-rt", clusterID.InfraID)
		routeName := fmt.Sprintf("%s-rt-default", clusterID.InfraID)
		routingTables = []capibm.VPCRoutingTableSource{
			{
				Type: capibm.SourceTypeProvision,
				Provision: capibm.VPCRoutingTableProvision{
					Name:                       rtName,
					RouteTransitGatewayIngress: capibm.RoutingTableIngressPolicyEnabled,
					AdvertiseRoutesTo: []capibm.AdvertiseRouteTarget{
						capibm.AdvertiseRouteTargetTransitGateway,
					},
					Routes: []capibm.VPCRoutingTableRoute{
						{
							Name:        routeName,
							Action:      capibm.VPCRoutingTableRouteActionDeliver,
							Destination: "0.0.0.0/0",
							Zone:        vpcZone,
							Advertise:   capibm.RouteAdvertisePolicyEnabled,
						},
					},
				},
			},
		}
	}

	powerVSCluster = &capibm.IBMPowerVSCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: capibm.GroupVersion.String(),
			Kind:       "IBMPowerVSCluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
			Annotations: map[string]string{
				"powervs.cluster.x-k8s.io/create-infra": "true",
			},
		},
		Spec: capibm.IBMPowerVSClusterSpec{
			Workspace:         workspaceSource,
			Network:           networkSource,
			Zone:              installConfig.Config.Platform.PowerVS.Zone,
			ResourceGroup:     rgSource,
			VPC:               vpcSource,
			VPCSecurityGroups: vpcSecurityGroups,
			TransitGateway:    tgSource,
			LoadBalancers:     loadBalancers,
			VPCPublicGateways: publicGateways,
			VPCRoutingTables:  routingTables,
			COSInstance:       cosSource,
			Ignition: capibm.Ignition{
				Version: "3.4",
			},
		},
	}

	// Use a custom DNS resolver for Internal publishing strategy
	if installConfig.Config.Publish == types.InternalPublishingStrategy {
		var dnsServerIP string
		err = installConfig.PowerVS.EnsureVPCNameIsSpecifiedForInternal(installConfig.Config.PowerVS.VPC)
		if err != nil {
			return nil, err
		}
		dnsServerIP, err = installConfig.PowerVS.GetDNSServerIP(ctx, installConfig.Config.PowerVS.VPC)
		if err != nil {
			return nil, fmt.Errorf("unable to find a DNS server for specified VPC: %s %w", installConfig.Config.PowerVS.VPC, err)
		}
		// Inject the DNS server IP into the network provision config.
		// DHCPSubnet uses a []string slice; DHCPServer uses a single string field.
		if useDHCPServer {
			powerVSCluster.Spec.Network.Provision.DHCPServer.DNSServer = dnsServerIP
		} else {
			powerVSCluster.Spec.Network.Provision.DHCPSubnet.DNSServers = []string{dnsServerIP}
		}
	}

	// If a VPC was specified, pass all subnets in it to cluster API
	if vpcStruct != nil {
		if installConfig.Config.Publish == types.InternalPublishingStrategy {
			err = installConfig.PowerVS.EnsureVPCIsPermittedNetwork(ctx, vpcStruct)
		}
		if err != nil {
			return nil, fmt.Errorf("error ensuring VPC is permitted: %s %w", *vpcStruct.Name, err)
		}
		subnets, err := installConfig.PowerVS.GetVPCSubnets(ctx, vpcStruct)
		if err != nil {
			return nil, fmt.Errorf("error getting subnets in specified VPC: %s %w", *vpcStruct.Name, err)
		}
		for _, subnet := range subnets {
			powerVSCluster.Spec.VPCSubnets = append(powerVSCluster.Spec.VPCSubnets,
				capibm.VPCSubnetSource{
					Type: capibm.SourceTypeReference,
					Reference: capibm.ResourceIdentifier{
						ID:   ptr.Deref(subnet.ID, ""),
						Name: ptr.Deref(subnet.Name, ""),
					},
				})
		}
		logrus.Debugf("GenerateClusterAssets: subnets = %+v", powerVSCluster.Spec.VPCSubnets)
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: powerVSCluster,
		File:   asset.File{Filename: "02_powervs-cluster.yaml"},
	})

	// IBMPowerVSImage still uses powervs/v1beta2
	// Build v1beta2 service reference for the image
	var imageService capibmv2.IBMPowerVSResourceReference
	if installConfig.Config.PowerVS.ServiceInstanceGUID == "" {
		serviceName := fmt.Sprintf("%s-power-iaas", clusterID.InfraID)
		imageService = capibmv2.IBMPowerVSResourceReference{
			Name: &serviceName,
		}
	} else {
		imageService = capibmv2.IBMPowerVSResourceReference{
			ID: &installConfig.Config.PowerVS.ServiceInstanceGUID,
		}
	}

	powerVSImage = &capibmv2.IBMPowerVSImage{
		TypeMeta: metav1.TypeMeta{
			APIVersion: capibmv2.GroupVersion.String(),
			Kind:       "IBMPowerVSImage",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      imageName,
			Namespace: capiutils.Namespace,
		},
		Spec: capibmv2.IBMPowerVSImageSpec{
			ClusterName:     clusterID.InfraID,
			ServiceInstance: &imageService,
			Bucket:          &bucket,
			Object:          &object,
			Region:          &cosRegion,
		},
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: powerVSImage,
		File:   asset.File{Filename: "03_powervs-image.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: capibm.GroupVersion.String(),
				Kind:       "IBMPowerVSCluster",
				Name:       powerVSCluster.Name,
				Namespace:  powerVSCluster.Namespace,
			},
		},
	}, nil
}
