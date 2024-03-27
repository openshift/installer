package powervs

import (
	"fmt"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capibm "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	powervstypes "github.com/openshift/installer/pkg/types/powervs"
)

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, bucket string, object string) (*capiutils.GenerateClusterAssetsOutput, error) {
	var (
		manifests          []*asset.RuntimeFile
		network            string
		service            capibm.IBMPowerVSResourceReference
		vpcName            string
		vpcRegion          string
		transitGatewayName string
		cosName            string
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
		logrus.Debugf("GenerateClusterAssets: imageName = %v", imageName)
		logrus.Debugf("GenerateClusterAssets: bucketName = %v", bucketName)
		logrus.Debugf("GenerateClusterAssets: powerVSCluster.Spec.ControlPlaneEndpoint.Host = %v", powerVSCluster.Spec.ControlPlaneEndpoint.Host)
	}()

	manifests = []*asset.RuntimeFile{}

	network = fmt.Sprintf("%s-network", clusterID.InfraID)

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

	transitGatewayName = fmt.Sprintf("%s-tg", clusterID.InfraID)

	cosName = fmt.Sprintf("%s-cos", clusterID.InfraID)

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
				BucketRegion: vpcRegion,
			},
			Ignition: &capibm.Ignition{
				Version: "3.4",
			},
		},
	}

	// Avoid empty manifest elements
	logrus.Debugf("GenerateClusterAssets: len(VPCSubnets) = %d", len(installConfig.Config.Platform.PowerVS.VPCSubnets))
	if len(installConfig.Config.Platform.PowerVS.VPCSubnets) > 0 {
		powerVSCluster.Spec.VPCSubnets = make([]capibm.Subnet, 0)
		for i := range installConfig.Config.Platform.PowerVS.VPCSubnets {
			// We cannot get the string in the loop and add it as it appears duplicated in the generated file
			powerVSCluster.Spec.VPCSubnets = append(powerVSCluster.Spec.VPCSubnets,
				capibm.Subnet{
					ID: &installConfig.Config.Platform.PowerVS.VPCSubnets[i],
				})
		}
		logrus.Debugf("GenerateClusterAssets: subnets = %+v", powerVSCluster.Spec.VPCSubnets)
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: powerVSCluster,
		File:   asset.File{Filename: "02_powervs-cluster.yaml"},
	})

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
			Region:          &vpcRegion,
		},
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: powerVSImage,
		File:   asset.File{Filename: "03_powervs-image.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRef: &corev1.ObjectReference{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
			Kind:       "IBMPowerVSCluster",
			Name:       powerVSCluster.Name,
			Namespace:  powerVSCluster.Namespace,
		},
	}, nil
}
