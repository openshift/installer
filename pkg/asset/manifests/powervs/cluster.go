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
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	powervsconfig "github.com/openshift/installer/pkg/asset/installconfig/powervs"
	machinespowervs "github.com/openshift/installer/pkg/asset/machines/powervs"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	var (
		manifests          []*asset.RuntimeFile
		network            string
		dhcpSubnetCIDR     string
		service            capibm.ResourceIdentifier
		workspaceSource    capibm.WorkspaceSource
		vpcNameOrID        string
		vpcStruct          *vpcv1.VPC
		vpcRegion          string
		cosName            string
		cosRegion          string
		bucketName         string
		transitGatewayName string
		client             *powervsconfig.Client
		err                error
		powerVSCluster     *capibm.IBMPowerVSCluster
	)

	defer func() {
		logrus.Debugf("GenerateClusterAssets: installConfig = %+v, clusterID = %v", installConfig, *clusterID)
		logrus.Debugf("GenerateClusterAssets: ic.ObjectMeta = %+v", installConfig.Config.ObjectMeta.Name)
		logrus.Debugf("GenerateClusterAssets: installConfig.Config.PowerVS = %+v", *installConfig.Config.PowerVS)
		logrus.Debugf("GenerateClusterAssets: installConfig.Config.Platform.PowerVS.VPC = %v", installConfig.Config.Platform.PowerVS.VPC)
		logrus.Debugf("GenerateClusterAssets: vpcNameOrID = %v", vpcNameOrID)
		logrus.Debugf("GenerateClusterAssets: vpcRegion = %v", vpcRegion)
		logrus.Debugf("GenerateClusterAssets: installConfig.Config.Platform.PowerVS.TransitGateway = %v", installConfig.Config.Platform.PowerVS.TransitGateway)
		logrus.Debugf("GenerateClusterAssets: transitGatewayName = %v", transitGatewayName)
		logrus.Debugf("GenerateClusterAssets: cosName = %v", cosName)
		logrus.Debugf("GenerateClusterAssets: cosRegion = %v", cosRegion)
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
	if installConfig.Config.PowerVS.ServiceInstanceGUID == "" {
		serviceName := fmt.Sprintf("%s-power-iaas", clusterID.InfraID)
		service = capibm.ResourceIdentifier{Name: serviceName}
		workspaceSource = capibm.WorkspaceSource{
			Type:      capibm.SourceTypeProvision,
			Provision: capibm.WorkspaceProvisionConfig{Name: serviceName},
		}
	} else {
		service = capibm.ResourceIdentifier{
			ID: installConfig.Config.PowerVS.ServiceInstanceGUID,
		}
		workspaceSource = capibm.WorkspaceSource{
			Type:      capibm.SourceTypeReference,
			Reference: service,
		}
	}

	client, err = powervsconfig.NewClient()
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
			VPCConnection: capibm.TransitGatewayConnectionSource{
				Type: capibm.SourceTypeProvision,
			},
			PowerVSConnection: capibm.TransitGatewayConnectionSource{
				Type: capibm.SourceTypeProvision,
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
			VPCConnection: capibm.TransitGatewayConnectionSource{
				Type: capibm.SourceTypeProvision,
			},
			PowerVSConnection: capibm.TransitGatewayConnectionSource{
				Type: capibm.SourceTypeProvision,
			},
		}
	}

	cosName = fmt.Sprintf("%s-cos", clusterID.InfraID)

	if cosRegion, err = powervstypes.COSRegionForPowerVSRegion(installConfig.Config.PowerVS.Region); err != nil {
		return nil, fmt.Errorf("unable to derive cosRegion from region: %s %w", installConfig.Config.PowerVS.Region, err)
	}

	bucketName = fmt.Sprintf("%s-bootstrap-ign", clusterID.InfraID)

	cosSource := capibm.COSInstanceSource{
		Type:         capibm.SourceTypeProvision,
		BucketName:   bucketName,
		BucketRegion: cosRegion,
		Provision: capibm.COSInstanceProvision{
			Name: cosName,
		},
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

	loadBalancers := buildLoadBalancers(clusterID.InfraID, useDHCPServer, installConfig.Config.Publish)

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
		pgwName        string
		pgwVPCZone     string
	)
	if !useDHCPServer {
		// Derive the VPC zone from the PowerVS zone for the Public Gateway and Routing Table.
		var vpcZoneErr error
		pgwVPCZone, vpcZoneErr = powervstypes.VPCZoneForPowerVSZone(installConfig.Config.Platform.PowerVS.Zone)
		if vpcZoneErr != nil {
			return nil, fmt.Errorf("unable to derive VPC zone for Public Gateway from PowerVS zone %s: %w",
				installConfig.Config.Platform.PowerVS.Zone, vpcZoneErr)
		}

		// Public Gateway — outbound internet access for VPC subnets.
		pgwName = fmt.Sprintf("%s-pgw", clusterID.InfraID)
		publicGateways = []capibm.VPCPublicGateway{
			{
				Type: capibm.SourceTypeProvision,
				Zone: pgwVPCZone,
				Provision: capibm.VPCPublicGatewayProvision{
					Name: pgwName,
				},
			},
		}

		// Routing Table — advertises 0.0.0.0/0 via the route-mode NLB to the Transit Gateway.
		// The nextHop is intentionally left empty: CAPIBM auto-populates it from
		// LoadBalancerStatus.PrivateIPs[0] of the private route-mode load balancer.
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
							Zone:        pgwVPCZone,
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
			// Topology must be set so the CAPIBM webhook accepts the object.
			// Both DHCPSubnet and DHCPServer paths use a VPC and Load Balancer,
			// so PowerVSLoadBalancerTopology is always correct here.
			Topology:          capibm.PowerVSLoadBalancerTopology,
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

	// DHCPSubnet path: explicitly provision all three VPC subnets (one per zone) so that
	// CAPIBM does not auto-expand them without a PublicGateway reference. Only the subnet
	// in the same zone as the PGW gets the PublicGateway field set — IBM Cloud only supports
	// one PGW per zone, and the route-mode NLB also selects this zone.
	if !useDHCPServer && vpcStruct == nil {
		allVPCZones, vpcZonesErr := powervstypes.AvailableVPCZones(installConfig.Config.Platform.PowerVS.Region)
		if vpcZonesErr != nil {
			return nil, fmt.Errorf("unable to derive VPC zones for PowerVS region %s: %w",
				installConfig.Config.Platform.PowerVS.Region, vpcZonesErr)
		}
		for _, zone := range allVPCZones {
			subnetName := fmt.Sprintf("%s-subnet-%s", clusterID.InfraID, zone)
			subnetSource := capibm.VPCSubnetSource{
				Type: capibm.SourceTypeProvision,
				Zone: zone,
				Provision: capibm.VPCSubnetProvision{
					Name: subnetName,
				},
			}
			// Attach the PGW only to the subnet that shares its zone.
			if zone == pgwVPCZone {
				subnetSource.PublicGateway = capibm.ResourceIdentifier{
					Name: pgwName,
				}
			}
			powerVSCluster.Spec.VPCSubnets = append(powerVSCluster.Spec.VPCSubnets, subnetSource)
		}
		logrus.Debugf("GenerateClusterAssets: VPCSubnets with PGW attachment = %+v", powerVSCluster.Spec.VPCSubnets)
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
			// Disable SNAT for disconnected scenario.
			powerVSCluster.Spec.Network.Provision.DHCPServer.Snat = dhcpServerSnatPolicy(installConfig.Config)
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
			subnetSource := capibm.VPCSubnetSource{
				Type: capibm.SourceTypeReference,
				Reference: capibm.ResourceIdentifier{
					ID:   ptr.Deref(subnet.ID, ""),
					Name: ptr.Deref(subnet.Name, ""),
				},
			}
			// For DHCPSubnet mode, attach the PGW to the subnet in the matching zone.
			if !useDHCPServer && subnet.Zone != nil && subnet.Zone.Name != nil && *subnet.Zone.Name == pgwVPCZone {
				subnetSource.PublicGateway = capibm.ResourceIdentifier{
					Name: pgwName,
				}
			}
			powerVSCluster.Spec.VPCSubnets = append(powerVSCluster.Spec.VPCSubnets, subnetSource)
		}
		logrus.Debugf("GenerateClusterAssets: subnets = %+v", powerVSCluster.Spec.VPCSubnets)
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: powerVSCluster,
		File:   asset.File{Filename: "02_powervs-cluster.yaml"},
	})

	// When osImage is set, import the image from COS into the workspace; otherwise
	// the machines use a stock catalog image and no image manifest is needed.
	if osImage := installConfig.Config.Platform.PowerVS.OSImage; osImage != "" {
		imageBucket, imageObject, imageRegion, err := powervstypes.ParseOSImageURL(osImage)
		if err != nil {
			return nil, fmt.Errorf("invalid platform.powervs.osImage %q: %w", osImage, err)
		}
		logrus.Infof("Importing PowerVS image from bucket %s, object %s, region %s", imageBucket, imageObject, imageRegion)

		manifests = append(manifests, &asset.RuntimeFile{
			Object: &capibm.IBMPowerVSImage{
				TypeMeta: metav1.TypeMeta{
					APIVersion: capibm.GroupVersion.String(),
					Kind:       "IBMPowerVSImage",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      machinespowervs.ImportedImageName(clusterID.InfraID),
					Namespace: capiutils.Namespace,
				},
				Spec: capibm.IBMPowerVSImageSpec{
					ClusterName: clusterID.InfraID,
					Workspace:   service,
					Bucket:      imageBucket,
					Object:      imageObject,
					Region:      imageRegion,
				},
			},
			File: asset.File{Filename: "03_powervs-image.yaml"},
		})
	} else {
		logrus.Infof("Using PowerVS stock image: %s", machinespowervs.OSImageNameFromStream(installConfig.Config.OSImageStream))
	}

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

// dhcpServerSnatPolicy returns the SNAT policy for the legacy DHCPServer network.
// SNAT is disabled for disconnected clusters (Internal publishing with image mirror
// sources configured) so that egress is not NATed through the DHCP server.
func dhcpServerSnatPolicy(ic *types.InstallConfig) capibm.DHCPSnatPolicy {
	if ic.Publish == types.InternalPublishingStrategy &&
		(len(ic.DeprecatedImageContentSources) > 0 || len(ic.ImageDigestSources) > 0) {
		return capibm.DHCPSnatPolicyDisabled
	}
	return capibm.DHCPSnatPolicyEnabled
}

// buildLoadBalancers builds the v1beta3 load balancers for the cluster.
//
// DHCPSubnet mode: uses a private "network-fixed" NLB with RouteMode enabled so it acts
// as a private L3 next-hop router (VNF). IBM Cloud rejects route_mode=true on public NLBs,
// so the LB visibility is explicitly set to Private. The CAPIBM webhook enforces that
// RouteMode=Enabled requires Profile=network-fixed. This NLB is only the egress next hop;
// it is never an API endpoint and gets no DNS record.
//
// DHCPSubnet + publish External: a separate public "application" LB named
// "<infraID>-loadbalancer" serves the external API (api.*), because the route-mode NLB
// is private and unreachable from outside the VPC. The route-mode NLB is named
// "<infraID>-loadbalancer-egress" so it does not collide with the public LB's name,
// which the DNS record creation and control plane machine sets key off.
//
// DHCPServer mode: uses a public LB with the default "application" profile (no RouteMode)
// because the DHCP Server's own SNAT provides egress — a route-mode NLB is not needed.
//
// The private/internal LB always uses the "application" profile.
func buildLoadBalancers(infraID string, useDHCPServer bool, publish types.PublishingStrategy) []capibm.LoadBalancerSource {
	lbName := fmt.Sprintf("%s-loadbalancer", infraID)
	lbIntName := fmt.Sprintf("%s-loadbalancer-int", infraID)
	lbEgressName := fmt.Sprintf("%s-loadbalancer-egress", infraID)

	sshListener := []capibm.AdditionalListener{
		{
			Port:     22,
			Protocol: capibm.LoadBalancerListenerProtocolTCP,
		},
	}

	var loadBalancers []capibm.LoadBalancerSource
	if useDHCPServer {
		loadBalancers = append(loadBalancers, capibm.LoadBalancerSource{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.LoadBalancerProvision{
				Name:                lbName,
				Type:                capibm.LoadBalancerTypePublic,
				AdditionalListeners: sshListener,
			},
		})
	} else {
		// DHCPSubnet path: route-mode NLB required for VPC Routing Table egress.
		// IBM Cloud requires route_mode=true only on private NLBs; public NLBs do not
		// support route mode. The route-mode LB acts as a private L3 next-hop router
		// (VNF), not as the external API server endpoint.
		routeModeName := lbName
		if publish == types.ExternalPublishingStrategy {
			routeModeName = lbEgressName
			// The public LB is listed first so it is preferred as the API endpoint.
			loadBalancers = append(loadBalancers, capibm.LoadBalancerSource{
				Type: capibm.SourceTypeProvision,
				Provision: capibm.LoadBalancerProvision{
					Name:                lbName,
					Type:                capibm.LoadBalancerTypePublic,
					AdditionalListeners: sshListener,
				},
			})
		}
		loadBalancers = append(loadBalancers, capibm.LoadBalancerSource{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.LoadBalancerProvision{
				Name:                routeModeName,
				Type:                capibm.LoadBalancerTypePrivate,
				AdditionalListeners: sshListener,
				Profile:             capibm.PowerVSLoadBalancerProfileNetworkFixed,
				RouteMode:           capibm.PowerVSLoadBalancerRouteModeEnabled,
			},
		})
	}

	return append(loadBalancers, capibm.LoadBalancerSource{
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
	})
}
