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
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck //CORS-3563

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	gcpic "github.com/openshift/installer/pkg/asset/installconfig/gcp"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	gcpconsts "github.com/openshift/installer/pkg/constants/gcp"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
)

const (
	// InstanceGroupRoleTag is the tag used in the instance
	// group name to maintain compatibility between MAPI & CAPI.
	InstanceGroupRoleTag = "master"

	resourceDescription = "Created By OpenShift Installer"
)

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

	firewallRules := []capg.FirewallRule{}
	managementPolicy := capg.RulesManagementUnmanaged
	if installConfig.Config.GCP.FirewallRulesManagement != gcp.UnmanagedFirewallRules {
		managementPolicy = capg.RulesManagementManaged

		firewallRules = append(firewallRules, createBootstrapFirewallRuleForCAPG(installConfig, clusterID)...)
		firewallRules = append(firewallRules, createFirewallRulesForCAPG(installConfig, clusterID)...)
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
				Name:                  ptr.To(networkName),
				Subnets:               subnets,
				AutoCreateSubnetworks: ptr.To(autoCreateSubnets),
				Firewall: capg.FirewallSpec{
					DefaultRulesManagement: managementPolicy,
					FirewallRules:          firewallRules,
				},
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

func createFirewallRulesForCAPG(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) []capg.FirewallRule {
	firewallRules := []capg.FirewallRule{}
	workerTag := fmt.Sprintf("%s-worker", clusterID.InfraID)
	masterTag := fmt.Sprintf("%s-control-plane", clusterID.InfraID)

	// control-plane rules are needed for worker<->master communication for worker provisioning
	firewallRules = append(firewallRules, capg.FirewallRule{
		Name:         fmt.Sprintf("%s-control-plane", clusterID.InfraID),
		SourceRanges: []string{},
		TargetTags:   []string{masterTag},
		SourceTags:   []string{workerTag, masterTag},
		Direction:    capg.FirewallRuleDirectionIngress,
		Description:  resourceDescription,
		Priority:     1000,
		Allowed: []capg.FirewallDescriptor{
			{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"22623"}, // Ignition
			}, {
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"10257"}, // Kube manager
			}, {
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"10259"}, // Kube scheduler
			},
		},
	})

	// etcd are needed for master communication for etcd nodes
	firewallRules = append(firewallRules, capg.FirewallRule{
		Name:         fmt.Sprintf("%s-etcd", clusterID.InfraID),
		SourceRanges: []string{},
		TargetTags:   []string{masterTag},
		SourceTags:   []string{masterTag},
		Direction:    capg.FirewallRuleDirectionIngress,
		Description:  resourceDescription,
		Priority:     1000,
		Allowed: []capg.FirewallDescriptor{
			{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"2379-2380"},
			},
		},
	})

	// Add a single firewall rule to allow the Google Cloud Engine health checks to access all of the services.
	// This rule enables the ingress load balancers to determine the health status of their instances.
	healthCheckSrcRanges := []string{"35.191.0.0/16", "130.211.0.0/22"}
	if installConfig.Config.Publish == types.InternalPublishingStrategy {
		healthCheckSrcRanges = append(healthCheckSrcRanges, []string{"209.85.152.0/22", "209.85.204.0/22"}...)
	}
	firewallRules = append(firewallRules, capg.FirewallRule{
		Name:         fmt.Sprintf("%s-health-checks", clusterID.InfraID),
		SourceRanges: healthCheckSrcRanges,
		TargetTags:   []string{masterTag},
		SourceTags:   []string{},
		Direction:    capg.FirewallRuleDirectionIngress,
		Description:  resourceDescription,
		Priority:     1000,
		Allowed: []capg.FirewallDescriptor{
			{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"6080", "6443", "22624"},
			},
		},
	})

	// internal-cluster rules are needed for worker<->master communication for k8s nodes
	firewallRules = append(firewallRules, capg.FirewallRule{
		Name:         fmt.Sprintf("%s-internal-cluster", clusterID.InfraID),
		SourceRanges: []string{},
		TargetTags:   []string{workerTag, masterTag},
		SourceTags:   []string{workerTag, masterTag},
		Direction:    capg.FirewallRuleDirectionIngress,
		Description:  resourceDescription,
		Priority:     1000,
		Allowed: []capg.FirewallDescriptor{
			{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"30000-32767"}, // k8s NodePorts
			}, {
				IPProtocol: capg.FirewallProtocolUDP,
				Ports:      []string{"30000-32767"}, // k8s NodePorts
			}, {
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"9000-9999"}, // host-level services
			}, {
				IPProtocol: capg.FirewallProtocolUDP,
				Ports:      []string{"9000-9999"}, // host-level services
			}, {
				IPProtocol: capg.FirewallProtocolUDP,
				Ports:      []string{"4789", "6081"}, // VXLAN and GENEVE
			}, {
				IPProtocol: capg.FirewallProtocolUDP,
				Ports:      []string{"500", "4500"}, // IKE and IKE(NAT-T)
			}, {
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"10250"}, // kubelet secure
			}, {
				IPProtocol: capg.FirewallProtocolESP,
			},
		},
	})

	// api rules are needed to access the kube-apiserver on master nodes
	machineCIDR := installConfig.Config.Networking.MachineNetwork[0].CIDR.String()
	apiSrcRanges := []string{}
	if !installConfig.Config.PublicAPI() {
		// For Internal, limit the source to the machineCIDR
		apiSrcRanges = append(apiSrcRanges, machineCIDR)
	}
	firewallRules = append(firewallRules, capg.FirewallRule{
		Name:        fmt.Sprintf("%s-api", clusterID.InfraID),
		Direction:   capg.FirewallRuleDirectionIngress,
		Description: resourceDescription,
		Allowed: []capg.FirewallDescriptor{
			{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"6443"}, // kube-apiserver
			},
		},
		TargetTags:   []string{masterTag},
		SourceRanges: apiSrcRanges,
	})

	// internal-network rules are used to access ssh and icmp over the machine network
	firewallRules = append(firewallRules, capg.FirewallRule{
		Name:         fmt.Sprintf("%s-internal-network", clusterID.InfraID),
		SourceRanges: []string{machineCIDR},
		TargetTags:   []string{workerTag, masterTag},
		SourceTags:   []string{},
		Description:  resourceDescription,
		Direction:    capg.FirewallRuleDirectionIngress,
		Priority:     1000,
		Allowed: []capg.FirewallDescriptor{
			{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"22"}, // SSH
			}, {
				IPProtocol: capg.FirewallProtocolICMP,
			},
		},
	})

	return firewallRules
}

func createBootstrapFirewallRuleForCAPG(installConfig *installconfig.InstallConfig, clusterID *installconfig.ClusterID) []capg.FirewallRule {
	bootstrapTag := fmt.Sprintf("%s-control-plane", clusterID.InfraID)
	sourceRanges := []string{}
	if installConfig.Config.Publish == types.ExternalPublishingStrategy {
		sourceRanges = append(sourceRanges, []string{"0.0.0.0/0"}...)
	} else {
		machineCIDR := installConfig.Config.Networking.MachineNetwork[0].CIDR.String()
		sourceRanges = append(sourceRanges, []string{machineCIDR}...)
	}

	return []capg.FirewallRule{{
		Name:         fmt.Sprintf("%s-bootstrap-in-ssh", clusterID.InfraID),
		SourceTags:   []string{},
		TargetTags:   []string{bootstrapTag},
		SourceRanges: sourceRanges,
		Description:  resourceDescription,
		Direction:    capg.FirewallRuleDirectionIngress,
		Priority:     1000,
		Allowed: []capg.FirewallDescriptor{
			{
				IPProtocol: capg.FirewallProtocolTCP,
				Ports:      []string{"22"}, // SSH
			}, {
				IPProtocol: capg.FirewallProtocolICMP,
			},
		},
	}}
}
