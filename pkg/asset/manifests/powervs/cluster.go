package powervs

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, bucket string, object string) (*capiutils.GenerateClusterAssetsOutput, error) {
	var (
		manifests          []*asset.RuntimeFile
		network            string
		dhcpSubnet         string
		service            capibm.IBMPowerVSResourceReference
		vpcName            string
		vpcRegion          string
		transitGatewayName string
		cosName            string
		cosRegion          string
		imageName          string
		bucketName         string
		err                error
		powerVSCluster     *capibm.IBMPowerVSCluster
		powerVSImage       *capibm.IBMPowerVSImage
	)

	defer func() {
		logrus.Debugf("GenerateClusterAssets: installConfig = %+v, clusterID = %v, bucket = %v, object = %v", installConfig, *clusterID, bucket, object)
		logrus.Debugf("GenerateClusterAssets: ic.ObjectMeta = %+v", installConfig.Config.ObjectMeta.Name)
		logrus.Debugf("GenerateClusterAssets: installConfig.Config.PowerVS = %+v", *installConfig.Config.PowerVS)
		logrus.Debugf("GenerateClusterAssets: vpcName = %v", vpcName)
		logrus.Debugf("GenerateClusterAssets: vpcRegion = %v", vpcRegion)
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
	dhcpSubnet = installConfig.Config.Networking.MachineNetwork[0].CIDR.String()
	if numNetworks := len(installConfig.Config.Networking.MachineNetwork); numNetworks > 1 {
		logrus.Infof("Warning: %d machineNetwork found! Ignoring all except the first.", numNetworks)
	}
	logrus.Debugf("GenerateClusterAssets: dhcpSubnet = %s", dhcpSubnet)

	if installConfig.Config.PowerVS.ServiceInstanceGUID == "" {
		serviceName := fmt.Sprintf("%s-power-iaas", clusterID.InfraID)

		service = capibm.IBMPowerVSResourceReference{
			Name: &serviceName,
		}
	} else {
		service = capibm.IBMPowerVSResourceReference{
			ID: &installConfig.Config.PowerVS.ServiceInstanceGUID,
		}
	}

	vpcName = installConfig.Config.Platform.PowerVS.VPCName
	if vpcName == "" {
		vpcName = fmt.Sprintf("vpc-%s", clusterID.InfraID)
	}

	vpcRegion = installConfig.Config.Platform.PowerVS.VPCRegion
	if vpcRegion == "" {
		if vpcRegion, err = powervstypes.VPCRegionForPowerVSRegion(installConfig.Config.PowerVS.Region); err != nil {
			return nil, fmt.Errorf("unable to derive vpcRegion from region: %s %w", installConfig.Config.PowerVS.Region, err)
		}
	}

	transitGatewayName = installConfig.Config.Platform.PowerVS.TGName
	if transitGatewayName == "" {
		transitGatewayName = fmt.Sprintf("%s-tg", clusterID.InfraID)
	}

	cosName = fmt.Sprintf("%s-cos", clusterID.InfraID)

	if cosRegion, err = powervstypes.COSRegionForPowerVSRegion(installConfig.Config.PowerVS.Region); err != nil {
		return nil, fmt.Errorf("unable to derive cosRegion from region: %s %w", installConfig.Config.PowerVS.Region, err)
	}

	imageName = fmt.Sprintf("rhcos-%s", clusterID.InfraID)

	bucketName = fmt.Sprintf("%s-bootstrap-ign", clusterID.InfraID)

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
			Network: capibm.IBMPowerVSResourceReference{
				Name: &network,
			},
			DHCPServer: &capibm.DHCPServer{
				Cidr: &dhcpSubnet,
			},
			ServiceInstance: &service,
			Zone:            &installConfig.Config.Platform.PowerVS.Zone,
			ResourceGroup: &capibm.IBMPowerVSResourceReference{
				Name: &installConfig.Config.Platform.PowerVS.PowerVSResourceGroup,
			},
			VPC: &capibm.VPCResourceReference{
				Name:   &vpcName,
				Region: &vpcRegion,
			},
			TransitGateway: &capibm.TransitGateway{
				Name: &transitGatewayName,
			},
			LoadBalancers: []capibm.VPCLoadBalancerSpec{
				{
					Name:   fmt.Sprintf("%s-loadbalancer", clusterID.InfraID),
					Public: ptr.To(true),
					AdditionalListeners: []capibm.AdditionalListenerSpec{
						{
							Port: 22,
						},
						// @BUG We should be able to specify this:
						// capibm.AdditionalListenerSpec{
						//	Port: 6443,
						// },
					},
				},
				{
					Name:   fmt.Sprintf("%s-loadbalancer-int", clusterID.InfraID),
					Public: ptr.To(false),
					AdditionalListeners: []capibm.AdditionalListenerSpec{
						// @BUG We should be able to specify this:
						// capibm.AdditionalListenerSpec{
						//	Port: 6443,
						// },
						{
							Port: 22623,
						},
					},
				},
			},
			CosInstance: &capibm.CosInstance{
				Name:         cosName,
				BucketName:   bucketName,
				BucketRegion: cosRegion,
			},
			Ignition: &capibm.Ignition{
				Version: "3.4",
			},
		},
	}

	// Use a custom resolver if using an Internal publishing strategy
	if installConfig.Config.Publish == types.InternalPublishingStrategy {
		var dnsServerIP string
		err := installConfig.PowerVS.EnsureVPCNameIsSpecifiedForInternal(installConfig.Config.PowerVS.VPCName)
		if err != nil {
			return nil, err
		}
		dnsServerIP, err = installConfig.PowerVS.GetDNSServerIP(context.TODO(), installConfig.Config.PowerVS.VPCName)
		if err != nil {
			return nil, fmt.Errorf("unable to find a DNS server for specified VPC: %s %w", installConfig.Config.PowerVS.VPCName, err)
		}

		powerVSCluster.Spec.DHCPServer.DNSServer = &dnsServerIP
		// TODO(mjturek): Restore once work is finished in 4.18 for disconnected scenario.
		if !(len(installConfig.Config.DeprecatedImageContentSources) == 0 && len(installConfig.Config.ImageDigestSources) == 0) {
			return nil, fmt.Errorf("deploying a disconnected cluster directly in 4.17 is not supported for Power VS. Please deploy disconnected in 4.16 and upgrade to 4.17")
		}
	}

	// If a VPC was specified, pass all subnets in it to cluster API
	if installConfig.Config.Platform.PowerVS.VPCName != "" {
		logrus.Debugf("GenerateClusterAssets: VPCName = %s", installConfig.Config.Platform.PowerVS.VPCName)
		if installConfig.Config.Publish == types.InternalPublishingStrategy {
			err = installConfig.PowerVS.EnsureVPCIsPermittedNetwork(context.TODO(), installConfig.Config.PowerVS.VPCName)
		}
		if err != nil {
			return nil, fmt.Errorf("error ensuring VPC is permitted: %s %w", installConfig.Config.PowerVS.VPCName, err)
		}
		subnets, err := installConfig.PowerVS.GetVPCSubnets(context.TODO(), vpcName)
		if err != nil {
			return nil, fmt.Errorf("error getting subnets in specified VPC: %s %w", installConfig.Config.PowerVS.VPCName, err)
		}
		for _, subnet := range subnets {
			powerVSCluster.Spec.VPCSubnets = append(powerVSCluster.Spec.VPCSubnets,
				capibm.Subnet{
					ID:   subnet.ID,
					Name: subnet.Name,
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
			ClusterName:     clusterID.InfraID,
			ServiceInstance: &service,
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
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
				Kind:       "IBMPowerVSCluster",
				Name:       powerVSCluster.Name,
				Namespace:  powerVSCluster.Namespace,
			},
		},
	}, nil
}
