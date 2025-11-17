/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scope

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/api/compute/v1"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterScopeParams defines the input parameters used to create a new Scope.
type ClusterScopeParams struct {
	GCPServices
	Client     client.Client
	Cluster    *clusterv1.Cluster
	GCPCluster *infrav1.GCPCluster
}

// NewClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewClusterScope(ctx context.Context, params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.GCPCluster == nil {
		return nil, errors.New("failed to generate new scope from nil GCPCluster")
	}

	if params.Compute == nil {
		computeSvc, err := newComputeService(ctx, params.GCPCluster.Spec.CredentialsRef, params.Client, params.GCPCluster.Spec.ServiceEndpoints)
		if err != nil {
			return nil, errors.Errorf("failed to create gcp compute client: %v", err)
		}

		params.Compute = computeSvc
	}

	helper, err := patch.NewHelper(params.GCPCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &ClusterScope{
		client:      params.Client,
		Cluster:     params.Cluster,
		GCPCluster:  params.GCPCluster,
		GCPServices: params.GCPServices,
		patchHelper: helper,
	}, nil
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	client      client.Client
	patchHelper *patch.Helper

	Cluster    *clusterv1.Cluster
	GCPCluster *infrav1.GCPCluster
	GCPServices
}

// ANCHOR: ClusterGetter

// Cloud returns initialized cloud.
func (s *ClusterScope) Cloud() cloud.Cloud {
	return newCloud(s.Project(), s.GCPServices)
}

// NetworkCloud returns initialized cloud.
func (s *ClusterScope) NetworkCloud() cloud.Cloud {
	return newCloud(s.NetworkProject(), s.GCPServices)
}

// Project returns the current project name.
func (s *ClusterScope) Project() string {
	return s.GCPCluster.Spec.Project
}

// NetworkProject returns the project name where network resources should exist.
// The network project defaults to the Project when one is not supplied.
func (s *ClusterScope) NetworkProject() string {
	return ptr.Deref(s.GCPCluster.Spec.Network.HostProject, s.Project())
}

// SkipFirewallRuleCreation returns whether the spec indicates that firewall rules
// should be created or not. If the RulesManagement for the default firewall rules is
// set to unmanaged or when the cluster will include a shared VPC, the default firewall
// rule creation will be skipped.
func (s *ClusterScope) SkipFirewallRuleCreation() bool {
	return (s.GCPCluster.Spec.Network.Firewall.DefaultRulesManagement == infrav1.RulesManagementUnmanaged) || s.IsSharedVpc()
}

// IsSharedVpc returns true If sharedVPC used else , returns false.
func (s *ClusterScope) IsSharedVpc() bool {
	return s.NetworkProject() != s.Project()
}

// Region returns the cluster region.
func (s *ClusterScope) Region() string {
	return s.GCPCluster.Spec.Region
}

// Name returns the cluster name.
func (s *ClusterScope) Name() string {
	return s.Cluster.Name
}

// Namespace returns the cluster namespace.
func (s *ClusterScope) Namespace() string {
	return s.Cluster.Namespace
}

// NetworkName returns the cluster network unique identifier.
func (s *ClusterScope) NetworkName() string {
	return ptr.Deref(s.GCPCluster.Spec.Network.Name, "default")
}

// NetworkMtu returns the Network MTU of 1440 which is the default, otherwise returns back what is being set.
// Mtu: Maximum Transmission Unit in bytes. The minimum value for this field is
// 1300 and the maximum value is 8896. The suggested value is 1500, which is
// the default MTU used on the Internet, or 8896 if you want to use Jumbo
// frames. If unspecified, the value defaults to 1460.
// More info
// - https://pkg.go.dev/google.golang.org/api/compute/v1#Network
// - https://cloud.google.com/vpc/docs/mtu
func (s *ClusterScope) NetworkMtu() int64 {
	if s.GCPCluster.Spec.Network.Mtu == 0 {
		return int64(1460)
	}
	return s.GCPCluster.Spec.Network.Mtu
}

// NetworkLink returns the partial URL for the network.
func (s *ClusterScope) NetworkLink() string {
	return fmt.Sprintf("projects/%s/global/networks/%s", s.NetworkProject(), s.NetworkName())
}

// Network returns the cluster network object.
func (s *ClusterScope) Network() *infrav1.Network {
	return &s.GCPCluster.Status.Network
}

// AdditionalLabels returns the cluster additional labels.
func (s *ClusterScope) AdditionalLabels() infrav1.Labels {
	return s.GCPCluster.Spec.AdditionalLabels
}

// LoadBalancer returns the LoadBalancer configuration.
func (s *ClusterScope) LoadBalancer() infrav1.LoadBalancerSpec {
	return s.GCPCluster.Spec.LoadBalancer
}

// ResourceManagerTags returns ResourceManagerTags from the scope's GCPCluster. The returned value will never be nil.
func (s *ClusterScope) ResourceManagerTags() infrav1.ResourceManagerTags {
	if len(s.GCPCluster.Spec.ResourceManagerTags) == 0 {
		s.GCPCluster.Spec.ResourceManagerTags = infrav1.ResourceManagerTags{}
	}

	return s.GCPCluster.Spec.ResourceManagerTags.DeepCopy()
}

// ControlPlaneEndpoint returns the cluster control-plane endpoint.
func (s *ClusterScope) ControlPlaneEndpoint() clusterv1.APIEndpoint {
	endpoint := s.GCPCluster.Spec.ControlPlaneEndpoint
	endpoint.Port = 443
	if c := s.Cluster.Spec.ClusterNetwork; c != nil {
		endpoint.Port = ptr.Deref(c.APIServerPort, 443)
	}
	return endpoint
}

// FailureDomains returns the cluster failure domains.
func (s *ClusterScope) FailureDomains() clusterv1.FailureDomains {
	return s.GCPCluster.Status.FailureDomains
}

// ANCHOR_END: ClusterGetter

// ANCHOR: ClusterSetter

// SetReady sets cluster ready status.
func (s *ClusterScope) SetReady() {
	s.GCPCluster.Status.Ready = true
}

// SetFailureDomains sets cluster failure domains.
func (s *ClusterScope) SetFailureDomains(fd clusterv1.FailureDomains) {
	s.GCPCluster.Status.FailureDomains = fd
}

// SetControlPlaneEndpoint sets cluster control-plane endpoint.
func (s *ClusterScope) SetControlPlaneEndpoint(endpoint clusterv1.APIEndpoint) {
	s.GCPCluster.Spec.ControlPlaneEndpoint = endpoint
}

// ANCHOR_END: ClusterSetter

// ANCHOR: ClusterNetworkSpec

// NetworkSpec returns google compute network spec.
func (s *ClusterScope) NetworkSpec() *compute.Network {
	createSubnet := ptr.Deref(s.GCPCluster.Spec.Network.AutoCreateSubnetworks, true)
	network := &compute.Network{
		Name:                  s.NetworkName(),
		Description:           infrav1.ClusterTagKey(s.Name()),
		AutoCreateSubnetworks: createSubnet,
		ForceSendFields:       []string{"AutoCreateSubnetworks"},
		Mtu:                   s.NetworkMtu(),
	}

	return network
}

// NatRouterSpec returns google compute nat router spec.
func (s *ClusterScope) NatRouterSpec() *compute.Router {
	networkSpec := s.NetworkSpec()
	return &compute.Router{
		Name: fmt.Sprintf("%s-%s", networkSpec.Name, "router"),
		Nats: []*compute.RouterNat{
			{
				Name:                          fmt.Sprintf("%s-%s", networkSpec.Name, "nat"),
				NatIpAllocateOption:           "AUTO_ONLY",
				SourceSubnetworkIpRangesToNat: "ALL_SUBNETWORKS_ALL_IP_RANGES",
				MinPortsPerVm:                 s.GCPCluster.Spec.Network.MinPortsPerVM,
			},
		},
	}
}

// ANCHOR_END: ClusterNetworkSpec

// SubnetSpecs returns google compute subnets spec.
func (s *ClusterScope) SubnetSpecs() []*compute.Subnetwork {
	subnets := []*compute.Subnetwork{}
	for _, subnetwork := range s.GCPCluster.Spec.Network.Subnets {
		secondaryIPRanges := []*compute.SubnetworkSecondaryRange{}
		for rangeName, secondaryCidrBlock := range subnetwork.SecondaryCidrBlocks {
			secondaryIPRanges = append(secondaryIPRanges, &compute.SubnetworkSecondaryRange{RangeName: rangeName, IpCidrRange: secondaryCidrBlock})
		}
		subnets = append(subnets, &compute.Subnetwork{
			Name:                  subnetwork.Name,
			Region:                subnetwork.Region,
			EnableFlowLogs:        ptr.Deref(subnetwork.EnableFlowLogs, false),
			PrivateIpGoogleAccess: ptr.Deref(subnetwork.PrivateGoogleAccess, false),
			IpCidrRange:           subnetwork.CidrBlock,
			SecondaryIpRanges:     secondaryIPRanges,
			Description:           ptr.Deref(subnetwork.Description, infrav1.ClusterTagKey(s.Name())),
			Network:               s.NetworkLink(),
			Purpose:               ptr.Deref(subnetwork.Purpose, "PRIVATE_RFC_1918"),
			Role:                  "ACTIVE",
			StackType:             subnetwork.StackType,
		})
	}

	return subnets
}

// ANCHOR: ClusterFirewallSpec

// FirewallRulesSpec returns google compute firewall spec.
func (s *ClusterScope) FirewallRulesSpec() []*compute.Firewall {
	firewallRules := []*compute.Firewall{
		{
			Name:    fmt.Sprintf("allow-%s-healthchecks", s.Name()),
			Network: s.NetworkLink(),
			Allowed: []*compute.FirewallAllowed{
				{
					IPProtocol: "TCP",
					Ports: []string{
						strconv.FormatInt(6443, 10),
					},
				},
			},
			Direction: "INGRESS",
			SourceRanges: []string{
				"35.191.0.0/16",
				"130.211.0.0/22",
			},
			TargetTags: []string{
				s.Name() + "-control-plane",
			},
		},
		{
			Name:    fmt.Sprintf("allow-%s-cluster", s.Name()),
			Network: s.NetworkLink(),
			Allowed: []*compute.FirewallAllowed{
				{
					IPProtocol: "all",
				},
			},
			Direction: "INGRESS",
			SourceTags: []string{
				s.Name() + "-control-plane",
				s.Name() + "-node",
			},
			TargetTags: []string{
				s.Name() + "-control-plane",
				s.Name() + "-node",
			},
		},
	}

	return firewallRules
}

// ANCHOR_END: ClusterFirewallSpec

// ANCHOR: ClusterControlPlaneSpec

// AddressSpec returns google compute address spec.
func (s *ClusterScope) AddressSpec(lbname string) *compute.Address {
	return &compute.Address{
		Name:        fmt.Sprintf("%s-%s", s.Name(), lbname),
		AddressType: "EXTERNAL",
		IpVersion:   "IPV4",
	}
}

// BackendServiceSpec returns google compute backend-service spec.
func (s *ClusterScope) BackendServiceSpec(lbname string) *compute.BackendService {
	return &compute.BackendService{
		Name:                fmt.Sprintf("%s-%s", s.Name(), lbname),
		LoadBalancingScheme: "EXTERNAL",
		PortName:            "apiserver",
		Protocol:            "TCP",
		TimeoutSec:          int64((10 * time.Minute).Seconds()),
	}
}

// ForwardingRuleSpec returns google compute forwarding-rule spec.
func (s *ClusterScope) ForwardingRuleSpec(lbname string) *compute.ForwardingRule {
	port := int32(443)
	if c := s.Cluster.Spec.ClusterNetwork; c != nil {
		port = ptr.Deref(c.APIServerPort, 443)
	}
	portRange := fmt.Sprintf("%d-%d", port, port)
	return &compute.ForwardingRule{
		Name:                fmt.Sprintf("%s-%s", s.Name(), lbname),
		IPProtocol:          "TCP",
		LoadBalancingScheme: "EXTERNAL",
		PortRange:           portRange,
		Labels:              s.AdditionalLabels(),
	}
}

// HealthCheckSpec returns google compute health-check spec.
func (s *ClusterScope) HealthCheckSpec(lbname string) *compute.HealthCheck {
	return &compute.HealthCheck{
		Name: fmt.Sprintf("%s-%s", s.Name(), lbname),
		Type: "HTTPS",
		HttpsHealthCheck: &compute.HTTPSHealthCheck{
			Port:              6443,
			PortSpecification: "USE_FIXED_PORT",
			RequestPath:       "/readyz",
		},
		CheckIntervalSec:   10,
		TimeoutSec:         5,
		HealthyThreshold:   5,
		UnhealthyThreshold: 3,
	}
}

// InstanceGroupSpec returns google compute instance-group spec.
func (s *ClusterScope) InstanceGroupSpec(zone string) *compute.InstanceGroup {
	port := ptr.Deref(s.GCPCluster.Spec.Network.LoadBalancerBackendPort, 6443)
	tag := ptr.Deref(s.GCPCluster.Spec.LoadBalancer.APIServerInstanceGroupTagOverride, infrav1.APIServerRoleTagValue)
	return &compute.InstanceGroup{
		Name: fmt.Sprintf("%s-%s-%s", s.Name(), tag, zone),
		NamedPorts: []*compute.NamedPort{
			{
				Name: "apiserver",
				Port: int64(port),
			},
		},
	}
}

// TargetTCPProxySpec returns google compute target-tcp-proxy spec.
func (s *ClusterScope) TargetTCPProxySpec() *compute.TargetTcpProxy {
	return &compute.TargetTcpProxy{
		Name:        fmt.Sprintf("%s-%s", s.Name(), infrav1.APIServerRoleTagValue),
		ProxyHeader: "NONE",
	}
}

// ANCHOR_END: ClusterControlPlaneSpec

// PatchObject persists the cluster configuration and status.
func (s *ClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.GCPCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *ClusterScope) Close() error {
	return s.PatchObject()
}
