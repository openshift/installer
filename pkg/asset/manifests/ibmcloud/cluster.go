package ibmcloud

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	ibmcloudic "github.com/openshift/installer/pkg/asset/installconfig/ibmcloud"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
)

const operatingSystem = "rhel-coreos-stable-amd64"

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID, imageName string) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	platform := installConfig.Config.Platform.IBMCloud

	// TODO(cjschaef): Add support for creating VPC Subnet Address Pools (CIDRs) during Infrastructure bring up.
	// mainCIDR := capiutils.CIDRFromInstallConfig(installConfig)
	// Make sure we have a fresh instance of Metadata, in case of any service endpoint overrides.
	metadata := ibmcloudic.NewMetadata(installConfig.Config)
	client, err := metadata.Client()
	if err != nil {
		return nil, fmt.Errorf("failed creating IBM Cloud client %w", err)
	}

	// Collect and build information for Cluster manifest.
	resourceGroup := clusterID.InfraID
	// Override Resource Group if provided in Platform.
	if platform.ResourceGroupName != "" {
		resourceGroup = platform.ResourceGroupName
	}
	networkResourceGroup := resourceGroup
	// Override Network Resource Group if provided in Platform.
	if platform.NetworkResourceGroupName != "" {
		networkResourceGroup = platform.NetworkResourceGroupName
	}
	vpcName := platform.GetVPCName()
	if vpcName == "" {
		vpcName = fmt.Sprintf("%s-vpc", clusterID.InfraID)
	}

	// NOTE(cjschaef): Add support to use existing Image details, rather than always creating a new VPC Custom Image.
	imageSpec := &capibmcloud.ImageSpec{
		Name:            ptr.To(ibmcloudic.VSIImageName(clusterID.InfraID)),
		COSBucket:       ptr.To(ibmcloudic.VSIImageCOSBucketName(clusterID.InfraID)),
		COSBucketRegion: ptr.To(platform.Region),
		COSInstance:     ptr.To(ibmcloudic.COSInstanceName(clusterID.InfraID)),
		COSObject:       ptr.To(imageName),
		OperatingSystem: ptr.To(operatingSystem),
		ResourceGroup: &capibmcloud.IBMCloudResourceReference{
			Name: ptr.To(resourceGroup),
		},
	}

	// Build and transform Subnets into CAPI.Subnets.
	controlPlaneSubnets, err := metadata.ControlPlaneSubnets(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed collecting control plane subnets %w", err)
	}
	// If no Control Plane subnets were provided in InstallConfig, we build a default set to cover all provided zones, or zones in the region.
	if len(controlPlaneSubnets) == 0 {
		var zones []string
		// Use provided Control Plane zones, or default to all zones in the Region.
		if installConfig.Config.ControlPlane.Platform.IBMCloud != nil && len(installConfig.Config.ControlPlane.Platform.IBMCloud.Zones) != 0 {
			zones = installConfig.Config.ControlPlane.Platform.IBMCloud.Zones
		} else {
			var err error
			zones, err = client.GetVPCZonesForRegion(context.TODO(), platform.Region)
			if err != nil {
				return nil, fmt.Errorf("failed collecting zones in region: %w", err)
			}
		}
		if controlPlaneSubnets == nil {
			controlPlaneSubnets = make(map[string]ibmcloudic.Subnet, 0)
		}
		for _, zone := range zones {
			subnetName, err := ibmcloudic.CreateSubnetName(clusterID.InfraID, "master", zone)
			if err != nil {
				return nil, fmt.Errorf("failed creating subnet name: %w", err)
			}
			// Typically, the map is keyed by the Subnet ID, but we don't have that if we are generating new subnet names. Since the ID's don't get used in Cluster manifest generation, we should be okay, as the key is ignored during ibmcloudic.Subnet to capibmcloud.Subnet transition.
			controlPlaneSubnets[subnetName] = ibmcloudic.Subnet{
				Name: subnetName,
				Zone: zone,
			}
		}
	}
	capiControlPlaneSubnets := getCAPISubnets(controlPlaneSubnets)
	// Build the Subnets for the Load Balancers.
	loadBalancerControlPlaneSubnets := make([]capibmcloud.VPCResource, len(capiControlPlaneSubnets))
	for index, subnet := range capiControlPlaneSubnets {
		loadBalancerControlPlaneSubnets[index] = capibmcloud.VPCResource{
			Name: subnet.Name,
		}
	}

	// Build and transform Compute Subnets into CAPI.Subnets.
	computeSubnets, err := metadata.ComputeSubnets(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed collecting compute subnets %w", err)
	}
	// If no Compute subnets were provided in InstallConfig, we build a default set to cover all specified zones, or zones in the region.
	if len(computeSubnets) == 0 {
		var zones []string
		// Use provided Compute zones, or default to all zones in the Region.
		// NOTE(cjschaef): We only process the first Compute definition, which may result in complications if additional Compute definitions request different Zones.
		if installConfig.Config.Compute[0].Platform.IBMCloud != nil && len(installConfig.Config.Compute[0].Platform.IBMCloud.Zones) != 0 {
			zones = installConfig.Config.Compute[0].Platform.IBMCloud.Zones
		} else {
			var err error
			zones, err = client.GetVPCZonesForRegion(context.TODO(), platform.Region)
			if err != nil {
				return nil, fmt.Errorf("failed collecting zones in region: %w", err)
			}
		}
		if computeSubnets == nil {
			computeSubnets = make(map[string]ibmcloudic.Subnet, 0)
		}
		for _, zone := range zones {
			subnetName, err := ibmcloudic.CreateSubnetName(clusterID.InfraID, "worker", zone)
			if err != nil {
				return nil, fmt.Errorf("failed creating subnet name: %w", err)
			}
			// Typically, the map is keyed by the Subnet ID, but we don't have that if we are generating new subnet names. Since the ID's don't get used in Cluster manifest generation, we should be okay, as the key is ignored during ibmcloudic.Subnet to capibmcloud.Subnet transition.
			computeSubnets[subnetName] = ibmcloudic.Subnet{
				Name: subnetName,
				Zone: zone,
			}
		}
	}
	capiComputeSubnets := getCAPISubnets(computeSubnets)

	// Create a consolidated set of all subnets, to use when generating SecurityGroups (this should prevent duplicates that appear in both subnet slices, resulting in duplicate SecurityGroupRules for subnet CIDR's). We may not have CIDR's until Infrastructure creation, so rely on Subnet names, to lookup CIDR's at runtime.
	capiConsolidatedSubnets := consolidateCAPISubnets(capiControlPlaneSubnets, capiComputeSubnets)

	// Build the necessary Security Groups and Rules for the Cluster.
	vpcSecurityGroups := getVPCSecurityGroups(clusterID.InfraID, capiConsolidatedSubnets, installConfig.Config.Publish)
	// Build the Security Groups for the Load Balancers.
	loadBalancerSecurityGroups := []capibmcloud.VPCResource{
		{
			Name: ptr.To(fmt.Sprintf("%s-%s", clusterID.InfraID, kubeAPILBSGNameSuffix)),
		},
	}

	// Build the necessary Load Balancers.
	loadBalancers := getLoadBalancers(clusterID.InfraID, loadBalancerSecurityGroups, loadBalancerControlPlaneSubnets, installConfig.Config.Publish)

	// Create the IBMVPCCluster manifest.
	ibmcloudCluster := &capibmcloud.IBMVPCCluster{
		TypeMeta: metav1.TypeMeta{
			APIVersion: capibmcloud.GroupVersion.String(),
			Kind:       "IBMVPCCluster",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capibmcloud.IBMVPCClusterSpec{
			ControlPlaneEndpoint: capi.APIEndpoint{
				Host: fmt.Sprintf("api.%s.%s", installConfig.Config.ObjectMeta.Name, installConfig.Config.BaseDomain),
				Port: 6443,
			},
			Image: imageSpec,
			Network: &capibmcloud.VPCNetworkSpec{
				ControlPlaneSubnets: capiControlPlaneSubnets,
				LoadBalancers:       loadBalancers,
				ResourceGroup: &capibmcloud.IBMCloudResourceReference{
					Name: ptr.To(networkResourceGroup),
				},
				SecurityGroups: vpcSecurityGroups,
				VPC: &capibmcloud.VPCResource{
					Name: ptr.To(vpcName),
				},
				WorkerSubnets: capiComputeSubnets,
			},
			Region:        platform.Region,
			ResourceGroup: resourceGroup,
		},
	}

	ibmcloudCluster.SetGroupVersionKind(capibmcloud.GroupVersion.WithKind("IBMVPCCluster"))

	manifests = append(manifests, &asset.RuntimeFile{
		Object: ibmcloudCluster,
		File:   asset.File{Filename: "01_ibmcloud-cluster.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: "infrastructure.cluster.x-k8s.io/v1beta2",
				Kind:       "IBMVPCCluster",
				Name:       ibmcloudCluster.Name,
				Namespace:  ibmcloudCluster.Namespace,
			},
		},
	}, nil
}

// consolidateCAPISubnets will attempt to consolidate two Subnet slices, and attempt to remove any duplicated Subnets (appear in both slices).
// This does not attempt to remove duplicate Subnets that exist in a single slice however.
func consolidateCAPISubnets(subnetsA []capibmcloud.Subnet, subnetsB []capibmcloud.Subnet) []capibmcloud.Subnet {
	consolidatedSubnets := make([]capibmcloud.Subnet, len(subnetsA))
	copiedSubnetNames := make(map[string]bool, 0)

	for index, subnet := range subnetsA {
		consolidatedSubnets[index] = subnet
		copiedSubnetNames[*subnet.Name] = true
	}

	for _, subnet := range subnetsB {
		// If we don't already have the Subnet from subnetsA, append it to the consolidated list.
		if _, okay := copiedSubnetNames[*subnet.Name]; !okay {
			consolidatedSubnets = append(consolidatedSubnets, subnet)
		}
	}
	return consolidatedSubnets
}

// getCAPISubnets converts InstallConfig based Subnets to CAPI based Subnets for Cluster manifest generation.
func getCAPISubnets(subnets map[string]ibmcloudic.Subnet) []capibmcloud.Subnet {
	subnetList := make([]capibmcloud.Subnet, 0)
	for _, subnet := range subnets {
		subnetList = append(subnetList, capibmcloud.Subnet{
			Name: ptr.To(subnet.Name),
			Zone: ptr.To(subnet.Zone),
		})
	}
	return subnetList
}
