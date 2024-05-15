package gcp

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	capg "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
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

	masterSubnet := gcp.DefaultSubnetName(clusterID.InfraID, "master")
	if installConfig.Config.GCP.ControlPlaneSubnet != "" {
		masterSubnet = installConfig.Config.GCP.ControlPlaneSubnet
	}

	master := capg.SubnetSpec{
		Name:        masterSubnet,
		CidrBlock:   "",
		Description: ptr.To(description),
		Region:      installConfig.Config.GCP.Region,
	}

	workerSubnet := gcp.DefaultSubnetName(clusterID.InfraID, "worker")
	if installConfig.Config.GCP.ComputeSubnet != "" {
		workerSubnet = installConfig.Config.GCP.ComputeSubnet
	}

	worker := capg.SubnetSpec{
		Name:        workerSubnet,
		CidrBlock:   "",
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
		master.CidrBlock = masterCIDR.String()
	}

	if installConfig.Config.GCP.ComputeSubnet == "" {
		workerCIDR, err := cidr.Subnet(ipv4Net, 1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to create the worker subnet %w", err)
		}
		worker.CidrBlock = workerCIDR.String()
	}

	subnets := []capg.SubnetSpec{master, worker}
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
		},
	}
	gcpCluster.SetGroupVersionKind(capg.GroupVersion.WithKind("GCPCluster"))

	manifests = append(manifests, &asset.RuntimeFile{
		Object: gcpCluster,
		File:   asset.File{Filename: "02_gcp-cluster.yaml"},
	})

	return &capiutils.GenerateClusterAssetsOutput{
		Manifests: manifests,
		InfrastructureRef: &corev1.ObjectReference{
			APIVersion: capg.GroupVersion.String(),
			Kind:       "GCPCluster",
			Name:       gcpCluster.Name,
			Namespace:  gcpCluster.Namespace,
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
