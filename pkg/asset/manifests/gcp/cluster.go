package gcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/apparentlymart/go-cidr/cidr"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

// InstanceGroupRoleTag is the tag used in the instance
// group name to maintain compatibility between MAPI & CAPI.
const InstanceGroupRoleTag = "master"

// GenerateClusterAssets generates the manifests for the cluster-api.
func GenerateClusterAssets(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) (*capiutils.GenerateClusterAssetsOutput, error) {
	manifests := []*asset.RuntimeFile{}
	const description = "Created By OpenShift Installer"

	networkName := fmt.Sprintf("%s-network", clusterID.InfraID)
	if installConfig.Config.GCP.Network != "" {
		networkName = installConfig.Config.GCP.Network
	}

	networkProject := installConfig.Config.GCP.ProjectID
	if installConfig.Config.GCP.NetworkProjectID != "" {
		networkProject = installConfig.Config.GCP.NetworkProjectID
	}

	controlPlaneSubnetName := gcp.DefaultSubnetName(clusterID.InfraID, "master")
	controlPlaneSubnetCidr := ""
	if installConfig.Config.GCP.ControlPlaneSubnet != "" {
		controlPlaneSubnetName = installConfig.Config.GCP.ControlPlaneSubnet

		controlPlaneSubnet, err := getSubnet(context.TODO(), networkProject, installConfig.Config.GCP.Region, controlPlaneSubnetName, installConfig.Config.GCP.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get control plane subnet: %w", err)
		}
		// IpCidr is the IPv4 version, the IPv6 version can be accessed as well
		controlPlaneSubnetCidr = controlPlaneSubnet.IpCidrRange
	}

	controlPlane := capg.SubnetSpec{
		Name:        controlPlaneSubnetName,
		CidrBlock:   controlPlaneSubnetCidr,
		Description: ptr.To(description),
		Region:      installConfig.Config.GCP.Region,
	}

	computeSubnetName := gcp.DefaultSubnetName(clusterID.InfraID, "worker")
	computeSubnetCidr := ""
	if installConfig.Config.GCP.ComputeSubnet != "" {
		computeSubnetName = installConfig.Config.GCP.ComputeSubnet

		computeSubnet, err := getSubnet(context.TODO(), networkProject, installConfig.Config.GCP.Region, computeSubnetName, installConfig.Config.GCP.Endpoint)
		if err != nil {
			return nil, fmt.Errorf("failed to get compute subnet: %w", err)
		}
		// IpCidr is the IPv4 version, the IPv6 version can be accessed as well
		computeSubnetCidr = computeSubnet.IpCidrRange
	}

	compute := capg.SubnetSpec{
		Name:        computeSubnetName,
		CidrBlock:   computeSubnetCidr,
		Description: ptr.To(description),
		Region:      installConfig.Config.GCP.Region,
	}

	// Add the CIDR information.
	machineV4CIDRs := []string{}
	for _, network := range installConfig.Config.Networking.MachineNetwork {
		if network.CIDR.IPNet.IP.To4() != nil {
			machineV4CIDRs = append(machineV4CIDRs, network.CIDR.IPNet.String())
		}
	}

	if len(machineV4CIDRs) == 0 {
		return nil, fmt.Errorf("failed to parse machine CIDRs")
	}

	_, ipv4Net, err := net.ParseCIDR(machineV4CIDRs[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse machine network CIDR: %w", err)
	}

	if installConfig.Config.GCP.ControlPlaneSubnet == "" {
		masterCIDR, err := cidr.Subnet(ipv4Net, 1, 0)
		if err != nil {
			return nil, fmt.Errorf("failed to create the master subnet %w", err)
		}
		controlPlane.CidrBlock = masterCIDR.String()
	}

	if installConfig.Config.GCP.ComputeSubnet == "" {
		computeCIDR, err := cidr.Subnet(ipv4Net, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to create the compute subnet %w", err)
		}
		compute.CidrBlock = computeCIDR.String()
	}

	subnets := []capg.SubnetSpec{controlPlane, compute}
	// Subnets should never be auto created, even in shared VPC installs
	autoCreateSubnets := false

	labels := map[string]string{}
	labels[fmt.Sprintf(gcpconsts.ClusterIDLabelFmt, clusterID.InfraID)] = "owned"
	labels[fmt.Sprintf("capg-cluster-%s", clusterID.InfraID)] = "owned"
	for _, label := range installConfig.Config.GCP.UserLabels {
		labels[label.Key] = label.Value
	}

	capgLoadBalancerType := capg.InternalExternal
	if installConfig.Config.Publish == types.InternalPublishingStrategy {
		capgLoadBalancerType = capg.Internal
	}

	gcpCluster := &capg.GCPCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      clusterID.InfraID,
			Namespace: capiutils.Namespace,
		},
		Spec: capg.GCPClusterSpec{
			Project:              installConfig.Config.GCP.ProjectID,
			Region:               installConfig.Config.GCP.Region,
			ControlPlaneEndpoint: clusterv1.APIEndpoint{Port: 6443},
			Network: capg.NetworkSpec{
				// TODO: Need a network project for installs where the network resources will exist in another
				// project such as shared vpc installs
				Name:                  ptr.To(networkName),
				Subnets:               subnets,
				AutoCreateSubnetworks: ptr.To(autoCreateSubnets),
			},
			AdditionalLabels: labels,
			FailureDomains:   findFailureDomains(installConfig),
			LoadBalancer: capg.LoadBalancerSpec{
				APIServerInstanceGroupTagOverride: ptr.To(InstanceGroupRoleTag),
				LoadBalancerType:                  ptr.To(capgLoadBalancerType),
			},
			ResourceManagerTags: GetTagsFromInstallConfig(installConfig),
		},
	}
	gcpCluster.SetGroupVersionKind(capg.GroupVersion.WithKind("GCPCluster"))

	if endpoint := installConfig.Config.GCP.Endpoint; gcp.ShouldUseEndpointForInstaller(endpoint) {
		gcpCluster.Spec.ServiceEndpoints = &capg.ServiceEndpoints{
			ComputeServiceEndpoint:         fmt.Sprintf("https://compute-%s.p.googleapis.com/compute/v1/", endpoint.Name),
			ContainerServiceEndpoint:       fmt.Sprintf("https://container-%s.p.googleapis.com/container/v1/", endpoint.Name),
			IAMServiceEndpoint:             fmt.Sprintf("https://iam-%s.p.googleapis.com/", endpoint.Name),
			ResourceManagerServiceEndpoint: fmt.Sprintf("https://cloudresourcemanager-%s.p.googleapis.com/", endpoint.Name),
		}
	}

	// Set the network project during shared vpc installs
	if installConfig.Config.GCP.NetworkProjectID != "" {
		gcpCluster.Spec.Network.HostProject = ptr.To(installConfig.Config.GCP.NetworkProjectID)
	}

	manifests = append(manifests, &asset.RuntimeFile{
		Object: gcpCluster,
		File:   asset.File{Filename: "02_gcp-cluster.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRefs: []*corev1.ObjectReference{
			{
				APIVersion: capg.GroupVersion.String(),
				Kind:       "GCPCluster",
				Name:       gcpCluster.Name,
				Namespace:  gcpCluster.Namespace,
			},
		},
	}, nil
}

// findFailureDomains will find the failure domains or availability zones for the GCP platform.
// When the default machine platform is defined, take any zone from the compute node(s) and
// any defined in the control plane node(s). When the default machine platform is not defined,
// only use zones if both the compute and control plane node availability zones exist.
func findFailureDomains(installConfig *installconfig.InstallConfig) []string {
	zones := sets.New[string]()

	var controlPlaneZones, computeZones []string
	if installConfig.Config.ControlPlane.Platform.GCP != nil {
		controlPlaneZones = installConfig.Config.ControlPlane.Platform.GCP.Zones
	}

	if installConfig.Config.Compute[0].Platform.GCP != nil {
		computeZones = installConfig.Config.Compute[0].Platform.GCP.Zones
	}

	def := installConfig.Config.GCP.DefaultMachinePlatform
	if def != nil && len(def.Zones) > 0 {
		for _, zone := range def.Zones {
			zones.Insert(zone)
		}

		for _, zone := range controlPlaneZones {
			zones.Insert(zone)
		}

		for _, zone := range computeZones {
			zones.Insert(zone)
		}
	} else if len(controlPlaneZones) > 0 && len(computeZones) > 0 {
		for _, zone := range controlPlaneZones {
			zones.Insert(zone)
		}
		for _, zone := range computeZones {
			zones.Insert(zone)
		}
	}

	return zones.UnsortedList()
}

// getSubnet will find a subnet in a project by the name. The matching subnet structure will be returned if
// one is found.
func getSubnet(ctx context.Context, project, region, subnetName string, endpoint *gcp.PSCEndpoint) (*compute.Subnetwork, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()

	opts := []option.ClientOption{}
	if gcp.ShouldUseEndpointForInstaller(endpoint) {
		opts = append(opts, gcpic.CreateEndpointOption(endpoint.Name, gcpic.ServiceNameGCPCompute))
	}
	computeService, err := gcpic.GetComputeService(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create compute service: %w", err)
	}

	subnetService := compute.NewSubnetworksService(computeService)
	subnet, err := subnetService.Get(project, region, subnetName).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to find subnet %s: %w", subnetName, err)
	} else if subnet == nil {
		return nil, fmt.Errorf("subnet %s is empty", subnetName)
	}

	return subnet, nil
}

// GetTagsFromInstallConfig will return a slice of ResourceManagerTags from UserTags in install-config.
func GetTagsFromInstallConfig(installConfig *installconfig.InstallConfig) []capg.ResourceManagerTag {
	tags := make([]capg.ResourceManagerTag, len(installConfig.Config.Platform.GCP.UserTags))
	for i, tag := range installConfig.Config.Platform.GCP.UserTags {
		tags[i] = capg.ResourceManagerTag{
			ParentID: tag.ParentID,
			Key:      tag.Key,
			Value:    tag.Value,
		}
	}

	return tags
}
