package powervs

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		manifests            []*asset.RuntimeFile
		dhcpSubnet           string
		service              capibm.ResourceIdentifier
		workspaceSource      capibm.WorkspaceSource
		vpcNameOrID          string
		vpcStruct            *vpcv1.VPC
		vpcRegion            string
		cosName              string
		cosRegion            string
		imageName            string
		bucketName           string
		transitGatewayName   string
		client               *powervsconfig.Client
		vpcSource            capibm.VPCSource
		transitGatewaySource capibm.TransitGatewaySource
		err                  error
		powerVSCluster       *capibm.IBMPowerVSCluster
		powerVSImage         *capibm.IBMPowerVSImage
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

	logrus.Debugf("GenerateClusterAssets: len MachineNetwork = %d", len(installConfig.Config.Networking.MachineNetwork))
	dhcpSubnet = installConfig.Config.Networking.MachineNetwork[0].CIDR.String()
	if numNetworks := len(installConfig.Config.Networking.MachineNetwork); numNetworks > 1 {
		logrus.Infof("Warning: %d machineNetwork found! Ignoring all except the first.", numNetworks)
	}
	logrus.Debugf("GenerateClusterAssets: dhcpSubnet = %s", dhcpSubnet)

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

	// The VPC specified can be either:
	// 1) blank - CAPI will create one for us.
	// 2) an id of an existing VPC.
	// 3) a name of an existing VPC.
	vpcNameOrID = installConfig.Config.Platform.PowerVS.VPC
	vpcSource = capibm.VPCSource{Region: vpcRegion}

	if vpcStruct, err = client.GetVPCByID(ctx, vpcNameOrID, vpcRegion); err == nil {
		logrus.Debugf("GenerateClusterAssets: PowerVS.VPC ID is valid")
		vpcSource.Type = capibm.SourceTypeReference
		vpcSource.Reference = capibm.ResourceIdentifier{ID: vpcNameOrID}
	} else if vpcStruct, err = client.GetVPCByName(ctx, vpcNameOrID); err == nil {
		logrus.Debugf("GenerateClusterAssets: PowerVS.VPC Name is valid")
		vpcSource.Type = capibm.SourceTypeReference
		vpcSource.Reference = capibm.ResourceIdentifier{Name: vpcNameOrID}
	} else {
		if vpcNameOrID == "" {
			logrus.Debugf("GenerateClusterAssets: PowerVS.VPC is empty")
			vpcNameOrID = fmt.Sprintf("vpc-%s", clusterID.InfraID)
			vpcStruct = nil
			vpcSource.Type = capibm.SourceTypeProvision
			vpcSource.Provision = capibm.VPCProvision{Name: vpcNameOrID}
		} else {
			return nil, fmt.Errorf("generateClusterAssets could not handle vpc")
		}
	}
	vpcSecurityGroups := getVPCSecurityGroups(clusterID.InfraID)
	// The Transit Gateway can be either:
	// 1) blank - CAPI will create one for us.
	// 2) an id of an existing TG.
	// 3) a name of an existing TG.
	transitGatewayName = installConfig.Config.Platform.PowerVS.TransitGateway
	transitGatewaySource = capibm.TransitGatewaySource{}

	if err = client.TransitGatewayIDValid(ctx, transitGatewayName); err == nil {
		logrus.Debugf("GenerateClusterAssets: TG ID is valid")
		transitGatewaySource.Type = capibm.SourceTypeReference
		transitGatewaySource.Reference = capibm.ResourceIdentifier{
			ID: transitGatewayName,
		}
	} else if transitGatewayName == "" {
		logrus.Debugf("GenerateClusterAssets: PowerVS.TransitGateway is empty")
		transitGatewayName = fmt.Sprintf("%s-tg", clusterID.InfraID)
		transitGatewaySource.Type = capibm.SourceTypeProvision
		transitGatewaySource.Provision = capibm.TransitGatewayProvision{
			Name: transitGatewayName,
		}
		transitGatewaySource.VPCConnection = capibm.TransitGatewayConnectionSource{
			Type: capibm.SourceTypeProvision,
		}
		transitGatewaySource.PowerVSConnection = capibm.TransitGatewayConnectionSource{
			Type: capibm.SourceTypeProvision,
		}
	} else {
		logrus.Debugf("GenerateClusterAssets: TG name is being used")
		transitGatewaySource.Type = capibm.SourceTypeReference
		transitGatewaySource.Reference = capibm.ResourceIdentifier{
			Name: transitGatewayName,
		}
	}
	cosName = fmt.Sprintf("%s-cos", clusterID.InfraID)

	if cosRegion, err = powervstypes.COSRegionForPowerVSRegion(installConfig.Config.PowerVS.Region); err != nil {
		return nil, fmt.Errorf("unable to derive cosRegion from region: %s %w", installConfig.Config.PowerVS.Region, err)
	}

	imageName = fmt.Sprintf("rhcos-%s", clusterID.InfraID)

	bucketName = fmt.Sprintf("%s-bootstrap-ign", clusterID.InfraID)

	vpcSubnets := make([]capibm.VPCSubnetSource, 0)
	if vpcStruct == nil {
		vpcSubnets = append(vpcSubnets, capibm.VPCSubnetSource{
			Type: capibm.SourceTypeProvision,
			Provision: capibm.VPCSubnetProvision{
				Name: fmt.Sprintf("%s-vpcsubnet", clusterID.InfraID),
			},
		})
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
			Topology: capibm.PowerVSLoadBalancerTopology,
			ResourceGroup: capibm.ResourceGroupSource{
				Type: capibm.SourceTypeReference,
				Reference: capibm.ResourceIdentifier{
					Name: installConfig.Config.Platform.PowerVS.PowerVSResourceGroup,
				},
			},
			Zone:      installConfig.Config.Platform.PowerVS.Zone,
			Workspace: workspaceSource,
			Network: capibm.NetworkSource{
				Type: capibm.SourceTypeProvision,
				Provision: capibm.NetworkProvisionConfig{
					DHCPServer: capibm.DHCPServer{
						Name: fmt.Sprintf("%s-dhcp", clusterID.InfraID),
						CIDR: dhcpSubnet,
					},
				},
			},
			VPC:               vpcSource,
			VPCSubnets:        vpcSubnets,
			VPCSecurityGroups: vpcSecurityGroups,
			TransitGateway:    transitGatewaySource,
			LoadBalancers: []capibm.LoadBalancerSource{
				{
					Type: capibm.SourceTypeProvision,
					Provision: capibm.LoadBalancerProvision{
						Name: fmt.Sprintf("%s-loadbalancer", clusterID.InfraID),
						Type: capibm.LoadBalancerTypePublic,
						AdditionalListeners: []capibm.AdditionalListener{
							{
								Port:     22,
								Protocol: capibm.LoadBalancerListenerProtocolTCP,
							},
							// @BUG We should be able to specify this:
							// capibm.AdditionalListener{
							//	Port: 6443,
							// },
						},
					},
				},
				{
					Type: capibm.SourceTypeProvision,
					Provision: capibm.LoadBalancerProvision{
						Name: fmt.Sprintf("%s-loadbalancer-int", clusterID.InfraID),
						Type: capibm.LoadBalancerTypePrivate,
						AdditionalListeners: []capibm.AdditionalListener{
							// @BUG We should be able to specify this:
							// capibm.AdditionalListener{
							//	Port: 6443,
							// },
							{
								Port:     22623,
								Protocol: capibm.LoadBalancerListenerProtocolTCP,
							},
						},
					},
				},
			},
			COSInstance: capibm.COSInstanceSource{
				Type:         capibm.SourceTypeProvision,
				BucketName:   bucketName,
				BucketRegion: cosRegion,
				Provision: capibm.COSInstanceProvision{
					Name: cosName,
				},
			},
			Ignition: capibm.Ignition{
				Version: "3.4",
			},
		},
	}

	// Use a custom resolver if using an Internal publishing strategy
	if installConfig.Config.Publish == types.InternalPublishingStrategy {
		var dnsServerIP string
		err := installConfig.PowerVS.EnsureVPCNameIsSpecifiedForInternal(installConfig.Config.PowerVS.VPC)
		if err != nil {
			return nil, err
		}
		dnsServerIP, err = installConfig.PowerVS.GetDNSServerIP(ctx, installConfig.Config.PowerVS.VPC)
		if err != nil {
			return nil, fmt.Errorf("unable to find a DNS server for specified VPC: %s %w", installConfig.Config.PowerVS.VPC, err)
		}

		powerVSCluster.Spec.Network.Provision.DHCPServer.DNSServer = dnsServerIP
		// Disable SNAT for disconnected scenario.
		if len(installConfig.Config.DeprecatedImageContentSources) == 0 && len(installConfig.Config.ImageDigestSources) == 0 {
			powerVSCluster.Spec.Network.Provision.DHCPServer.Snat = capibm.DHCPSnatPolicyEnabled
		} else {
			powerVSCluster.Spec.Network.Provision.DHCPServer.Snat = capibm.DHCPSnatPolicyDisabled
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
		powerVSCluster.Spec.VPCSubnets = nil
		for _, subnet := range subnets {
			powerVSCluster.Spec.VPCSubnets = append(powerVSCluster.Spec.VPCSubnets,
				capibm.VPCSubnetSource{
					Type: capibm.SourceTypeReference,
					Reference: capibm.ResourceIdentifier{
						ID: *subnet.ID,
					},
				})
		}
		logrus.Debugf("GenerateClusterAssets: subnets = %+v", powerVSCluster.Spec.VPCSubnets)
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: powerVSCluster,
		File:   asset.File{Filename: "02_powervs-cluster.yaml"},
	})

	if vpcRegion != cosRegion {
		logrus.Debugf("GenerateClusterAssets: vpcRegion(%s) is different than cosRegion(%s), cosRegion. Overriding bucket name", vpcRegion, cosRegion)
		bucket = fmt.Sprintf("rhcos-powervs-images-%s", cosRegion)
	}

	powerVSImage = &capibm.IBMPowerVSImage{
		TypeMeta: metav1.TypeMeta{
			APIVersion: capibm.GroupVersion.String(),
			Kind:       "IBMPowerVSImage",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      imageName,
			Namespace: capiutils.Namespace,
		},
		Spec: capibm.IBMPowerVSImageSpec{
			ClusterName: clusterID.InfraID,
			Workspace:   service,
			Bucket:      bucket,
			Object:      object,
			Region:      cosRegion,
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
