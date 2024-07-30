/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"k8s.io/utils/ptr"
)

const (
	// DefaultAPIServerPort defines the API server port when defining a Load Balancer.
	DefaultAPIServerPort = 6443
	// DefaultAPIServerPortString defines the API server port as a string for convenience.
	DefaultAPIServerPortString = "6443"
	// DefaultAPIServerHealthCheckPath the API server health check path.
	DefaultAPIServerHealthCheckPath = "/readyz"
	// DefaultAPIServerHealthCheckIntervalSec the API server health check interval in seconds.
	DefaultAPIServerHealthCheckIntervalSec = 10
	// DefaultAPIServerHealthCheckTimeoutSec the API server health check timeout in seconds.
	DefaultAPIServerHealthCheckTimeoutSec = 5
	// DefaultAPIServerHealthThresholdCount the API server health check threshold count.
	DefaultAPIServerHealthThresholdCount = 5
	// DefaultAPIServerUnhealthThresholdCount the API server unhealthy check threshold count.
	DefaultAPIServerUnhealthThresholdCount = 3

	// ZoneTypeAvailabilityZone defines the regular AWS zones in the Region.
	ZoneTypeAvailabilityZone ZoneType = "availability-zone"
	// ZoneTypeLocalZone defines the AWS zone type in Local Zone infrastructure.
	ZoneTypeLocalZone ZoneType = "local-zone"
	// ZoneTypeWavelengthZone defines the AWS zone type in Wavelength infrastructure.
	ZoneTypeWavelengthZone ZoneType = "wavelength-zone"
)

// NetworkStatus encapsulates AWS networking resources.
type NetworkStatus struct {
	// SecurityGroups is a map from the role/kind of the security group to its unique name, if any.
	SecurityGroups map[SecurityGroupRole]SecurityGroup `json:"securityGroups,omitempty"`

	// APIServerELB is the Kubernetes api server load balancer.
	APIServerELB LoadBalancer `json:"apiServerElb,omitempty"`

	// SecondaryAPIServerELB is the secondary Kubernetes api server load balancer.
	SecondaryAPIServerELB LoadBalancer `json:"secondaryAPIServerELB,omitempty"`

	// NatGatewaysIPs contains the public IPs of the NAT Gateways
	NatGatewaysIPs []string `json:"natGatewaysIPs,omitempty"`
}

// ELBScheme defines the scheme of a load balancer.
type ELBScheme string

var (
	// ELBSchemeInternetFacing defines an internet-facing, publicly
	// accessible AWS ELB scheme.
	ELBSchemeInternetFacing = ELBScheme("internet-facing")

	// ELBSchemeInternal defines an internal-only facing
	// load balancer internal to an ELB.
	ELBSchemeInternal = ELBScheme("internal")
)

func (e ELBScheme) String() string {
	return string(e)
}

// Equals returns true if two ELBScheme are equal.
func (e ELBScheme) Equals(other *ELBScheme) bool {
	if other == nil {
		return false
	}

	return e == *other
}

// ELBProtocol defines listener protocols for a load balancer.
type ELBProtocol string

func (e ELBProtocol) String() string {
	return string(e)
}

var (
	// ELBProtocolTCP defines the ELB API string representing the TCP protocol.
	ELBProtocolTCP = ELBProtocol("TCP")
	// ELBProtocolSSL defines the ELB API string representing the TLS protocol.
	ELBProtocolSSL = ELBProtocol("SSL")
	// ELBProtocolHTTP defines the ELB API string representing the HTTP protocol at L7.
	ELBProtocolHTTP = ELBProtocol("HTTP")
	// ELBProtocolHTTPS defines the ELB API string representing the HTTP protocol at L7.
	ELBProtocolHTTPS = ELBProtocol("HTTPS")
	// ELBProtocolTLS defines the NLB API string representing the TLS protocol.
	ELBProtocolTLS = ELBProtocol("TLS")
	// ELBProtocolUDP defines the NLB API string representing the UDP protocol.
	ELBProtocolUDP = ELBProtocol("UDP")
)

// TargetGroupHealthCheck defines health check settings for the target group.
type TargetGroupHealthCheck struct {
	Protocol                *string `json:"protocol,omitempty"`
	Path                    *string `json:"path,omitempty"`
	Port                    *string `json:"port,omitempty"`
	IntervalSeconds         *int64  `json:"intervalSeconds,omitempty"`
	TimeoutSeconds          *int64  `json:"timeoutSeconds,omitempty"`
	ThresholdCount          *int64  `json:"thresholdCount,omitempty"`
	UnhealthyThresholdCount *int64  `json:"unhealthyThresholdCount,omitempty"`
}

// TargetGroupHealthCheckAPISpec defines the optional health check settings for the API target group.
type TargetGroupHealthCheckAPISpec struct {
	// The approximate amount of time, in seconds, between health checks of an individual
	// target.
	// +kubebuilder:validation:Minimum=5
	// +kubebuilder:validation:Maximum=300
	// +optional
	IntervalSeconds *int64 `json:"intervalSeconds,omitempty"`

	// The amount of time, in seconds, during which no response from a target means
	// a failed health check.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=120
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// The number of consecutive health check successes required before considering
	// a target healthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	ThresholdCount *int64 `json:"thresholdCount,omitempty"`

	// The number of consecutive health check failures required before considering
	// a target unhealthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	UnhealthyThresholdCount *int64 `json:"unhealthyThresholdCount,omitempty"`
}

// TargetGroupHealthCheckAdditionalSpec defines the optional health check settings for the additional target groups.
type TargetGroupHealthCheckAdditionalSpec struct {
	// The protocol to use to health check connect with the target. When not specified the Protocol
	// will be the same of the listener.
	// +kubebuilder:validation:Enum=TCP;HTTP;HTTPS
	// +optional
	Protocol *string `json:"protocol,omitempty"`

	// The port the load balancer uses when performing health checks for additional target groups. When
	// not specified this value will be set for the same of listener port.
	// +optional
	Port *string `json:"port,omitempty"`

	// The destination for health checks on the targets when using the protocol HTTP or HTTPS,
	// otherwise the path will be ignored.
	// +optional
	Path *string `json:"path,omitempty"`
	// The approximate amount of time, in seconds, between health checks of an individual
	// target.
	// +kubebuilder:validation:Minimum=5
	// +kubebuilder:validation:Maximum=300
	// +optional
	IntervalSeconds *int64 `json:"intervalSeconds,omitempty"`

	// The amount of time, in seconds, during which no response from a target means
	// a failed health check.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=120
	// +optional
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// The number of consecutive health check successes required before considering
	// a target healthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	ThresholdCount *int64 `json:"thresholdCount,omitempty"`

	// The number of consecutive health check failures required before considering
	// a target unhealthy.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=10
	// +optional
	UnhealthyThresholdCount *int64 `json:"unhealthyThresholdCount,omitempty"`
}

// TargetGroupAttribute defines attribute key values for V2 Load Balancer Attributes.
type TargetGroupAttribute string

var (
	// TargetGroupAttributeEnablePreserveClientIP defines the attribute key for enabling preserve client IP.
	TargetGroupAttributeEnablePreserveClientIP = "preserve_client_ip.enabled"
)

// LoadBalancerAttribute defines a set of attributes for a V2 load balancer.
type LoadBalancerAttribute string

var (
	// LoadBalancerAttributeEnableLoadBalancingCrossZone defines the attribute key for enabling load balancing cross zone.
	LoadBalancerAttributeEnableLoadBalancingCrossZone = "load_balancing.cross_zone.enabled"
	// LoadBalancerAttributeIdleTimeTimeoutSeconds defines the attribute key for idle timeout.
	LoadBalancerAttributeIdleTimeTimeoutSeconds = "idle_timeout.timeout_seconds"
	// LoadBalancerAttributeIdleTimeDefaultTimeoutSecondsInSeconds defines the default idle timeout in seconds.
	LoadBalancerAttributeIdleTimeDefaultTimeoutSecondsInSeconds = "60"
)

// TargetGroupSpec specifies target group settings for a given listener.
// This is created first, and the ARN is then passed to the listener.
type TargetGroupSpec struct {
	// Name of the TargetGroup. Must be unique over the same group of listeners.
	// +kubebuilder:validation:MaxLength=32
	Name string `json:"name"`
	// Port is the exposed port
	Port int64 `json:"port"`
	// +kubebuilder:validation:Enum=tcp;tls;udp;TCP;TLS;UDP
	Protocol ELBProtocol `json:"protocol"`
	VpcID    string      `json:"vpcId"`
	// HealthCheck is the elb health check associated with the load balancer.
	HealthCheck *TargetGroupHealthCheck `json:"targetGroupHealthCheck,omitempty"`
}

// Listener defines an AWS network load balancer listener.
type Listener struct {
	Protocol    ELBProtocol     `json:"protocol"`
	Port        int64           `json:"port"`
	TargetGroup TargetGroupSpec `json:"targetGroup"`
}

// LoadBalancer defines an AWS load balancer.
type LoadBalancer struct {
	// ARN of the load balancer. Unlike the ClassicLB, ARN is used mostly
	// to define and get it.
	ARN string `json:"arn,omitempty"`
	// The name of the load balancer. It must be unique within the set of load balancers
	// defined in the region. It also serves as identifier.
	// +optional
	Name string `json:"name,omitempty"`

	// DNSName is the dns name of the load balancer.
	DNSName string `json:"dnsName,omitempty"`

	// Scheme is the load balancer scheme, either internet-facing or private.
	Scheme ELBScheme `json:"scheme,omitempty"`

	// AvailabilityZones is an array of availability zones in the VPC attached to the load balancer.
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// SubnetIDs is an array of subnets in the VPC attached to the load balancer.
	SubnetIDs []string `json:"subnetIds,omitempty"`

	// SecurityGroupIDs is an array of security groups assigned to the load balancer.
	SecurityGroupIDs []string `json:"securityGroupIds,omitempty"`

	// ClassicELBListeners is an array of classic elb listeners associated with the load balancer. There must be at least one.
	ClassicELBListeners []ClassicELBListener `json:"listeners,omitempty"`

	// HealthCheck is the classic elb health check associated with the load balancer.
	HealthCheck *ClassicELBHealthCheck `json:"healthChecks,omitempty"`

	// ClassicElbAttributes defines extra attributes associated with the load balancer.
	ClassicElbAttributes ClassicELBAttributes `json:"attributes,omitempty"`

	// Tags is a map of tags associated with the load balancer.
	Tags map[string]string `json:"tags,omitempty"`

	// ELBListeners is an array of listeners associated with the load balancer. There must be at least one.
	ELBListeners []Listener `json:"elbListeners,omitempty"`

	// ELBAttributes defines extra attributes associated with v2 load balancers.
	ELBAttributes map[string]*string `json:"elbAttributes,omitempty"`

	// LoadBalancerType sets the type for a load balancer. The default type is classic.
	// +kubebuilder:validation:Enum:=classic;elb;alb;nlb
	LoadBalancerType LoadBalancerType `json:"loadBalancerType,omitempty"`
}

// IsUnmanaged returns true if the Classic ELB is unmanaged.
func (b *LoadBalancer) IsUnmanaged(clusterName string) bool {
	return b.Name != "" && !Tags(b.Tags).HasOwned(clusterName)
}

// IsManaged returns true if Classic ELB is managed.
func (b *LoadBalancer) IsManaged(clusterName string) bool {
	return !b.IsUnmanaged(clusterName)
}

// ClassicELBAttributes defines extra attributes associated with a classic load balancer.
type ClassicELBAttributes struct {
	// IdleTimeout is time that the connection is allowed to be idle (no data
	// has been sent over the connection) before it is closed by the load balancer.
	IdleTimeout time.Duration `json:"idleTimeout,omitempty"`

	// CrossZoneLoadBalancing enables the classic load balancer load balancing.
	// +optional
	CrossZoneLoadBalancing bool `json:"crossZoneLoadBalancing,omitempty"`
}

// ClassicELBListener defines an AWS classic load balancer listener.
type ClassicELBListener struct {
	Protocol         ELBProtocol `json:"protocol"`
	Port             int64       `json:"port"`
	InstanceProtocol ELBProtocol `json:"instanceProtocol"`
	InstancePort     int64       `json:"instancePort"`
}

// ClassicELBHealthCheck defines an AWS classic load balancer health check.
type ClassicELBHealthCheck struct {
	Target             string        `json:"target"`
	Interval           time.Duration `json:"interval"`
	Timeout            time.Duration `json:"timeout"`
	HealthyThreshold   int64         `json:"healthyThreshold"`
	UnhealthyThreshold int64         `json:"unhealthyThreshold"`
}

// NetworkSpec encapsulates all things related to AWS network.
type NetworkSpec struct {
	// VPC configuration.
	// +optional
	VPC VPCSpec `json:"vpc,omitempty"`

	// Subnets configuration.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`

	// CNI configuration
	// +optional
	CNI *CNISpec `json:"cni,omitempty"`

	// SecurityGroupOverrides is an optional set of security groups to use for cluster instances
	// This is optional - if not provided new security groups will be created for the cluster
	// +optional
	SecurityGroupOverrides map[SecurityGroupRole]string `json:"securityGroupOverrides,omitempty"`

	// AdditionalControlPlaneIngressRules is an optional set of ingress rules to add to the control plane
	// +optional
	AdditionalControlPlaneIngressRules []IngressRule `json:"additionalControlPlaneIngressRules,omitempty"`
}

// IPv6 contains ipv6 specific settings for the network.
type IPv6 struct {
	// CidrBlock is the CIDR block provided by Amazon when VPC has enabled IPv6.
	// Mutually exclusive with IPAMPool.
	// +optional
	CidrBlock string `json:"cidrBlock,omitempty"`

	// PoolID is the IP pool which must be defined in case of BYO IP is defined.
	// Must be specified if CidrBlock is set.
	// Mutually exclusive with IPAMPool.
	// +optional
	PoolID string `json:"poolId,omitempty"`

	// EgressOnlyInternetGatewayID is the id of the egress only internet gateway associated with an IPv6 enabled VPC.
	// +optional
	EgressOnlyInternetGatewayID *string `json:"egressOnlyInternetGatewayId,omitempty"`

	// IPAMPool defines the IPAMv6 pool to be used for VPC.
	// Mutually exclusive with CidrBlock.
	// +optional
	IPAMPool *IPAMPool `json:"ipamPool,omitempty"`
}

// IPAMPool defines the IPAM pool to be used for VPC.
type IPAMPool struct {
	// ID is the ID of the IPAM pool this provider should use to create VPC.
	ID string `json:"id,omitempty"`
	// Name is the name of the IPAM pool this provider should use to create VPC.
	Name string `json:"name,omitempty"`
	// The netmask length of the IPv4 CIDR you want to allocate to VPC from
	// an Amazon VPC IP Address Manager (IPAM) pool.
	// Defaults to /16 for IPv4 if not specified.
	NetmaskLength int64 `json:"netmaskLength,omitempty"`
}

// VPCSpec configures an AWS VPC.
type VPCSpec struct {
	// ID is the vpc-id of the VPC this provider should use to create resources.
	ID string `json:"id,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	// Defaults to 10.0.0.0/16.
	// Mutually exclusive with IPAMPool.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// IPAMPool defines the IPAMv4 pool to be used for VPC.
	// Mutually exclusive with CidrBlock.
	IPAMPool *IPAMPool `json:"ipamPool,omitempty"`

	// IPv6 contains ipv6 specific settings for the network. Supported only in managed clusters.
	// This field cannot be set on AWSCluster object.
	// +optional
	IPv6 *IPv6 `json:"ipv6,omitempty"`

	// InternetGatewayID is the id of the internet gateway associated with the VPC.
	// +optional
	InternetGatewayID *string `json:"internetGatewayId,omitempty"`

	// CarrierGatewayID is the id of the internet gateway associated with the VPC,
	// for carrier network (Wavelength Zones).
	// +optional
	// +kubebuilder:validation:XValidation:rule="self.startsWith('cagw-')",message="Carrier Gateway ID must start with 'cagw-'"
	CarrierGatewayID *string `json:"carrierGatewayId,omitempty"`

	// Tags is a collection of tags describing the resource.
	Tags Tags `json:"tags,omitempty"`

	// AvailabilityZoneUsageLimit specifies the maximum number of availability zones (AZ) that
	// should be used in a region when automatically creating subnets. If a region has more
	// than this number of AZs then this number of AZs will be picked randomly when creating
	// default subnets. Defaults to 3
	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=1
	AvailabilityZoneUsageLimit *int `json:"availabilityZoneUsageLimit,omitempty"`

	// AvailabilityZoneSelection specifies how AZs should be selected if there are more AZs
	// in a region than specified by AvailabilityZoneUsageLimit. There are 2 selection schemes:
	// Ordered - selects based on alphabetical order
	// Random - selects AZs randomly in a region
	// Defaults to Ordered
	// +kubebuilder:default=Ordered
	// +kubebuilder:validation:Enum=Ordered;Random
	AvailabilityZoneSelection *AZSelectionScheme `json:"availabilityZoneSelection,omitempty"`

	// EmptyRoutesDefaultVPCSecurityGroup specifies whether the default VPC security group ingress
	// and egress rules should be removed.
	//
	// By default, when creating a VPC, AWS creates a security group called `default` with ingress and egress
	// rules that allow traffic from anywhere. The group could be used as a potential surface attack and
	// it's generally suggested that the group rules are removed or modified appropriately.
	//
	// NOTE: This only applies when the VPC is managed by the Cluster API AWS controller.
	//
	// +optional
	EmptyRoutesDefaultVPCSecurityGroup bool `json:"emptyRoutesDefaultVPCSecurityGroup,omitempty"`

	// PrivateDNSHostnameTypeOnLaunch is the type of hostname to assign to instances in the subnet at launch.
	// For IPv4-only and dual-stack (IPv4 and IPv6) subnets, an instance DNS name can be based on the instance IPv4 address (ip-name)
	// or the instance ID (resource-name). For IPv6 only subnets, an instance DNS name must be based on the instance ID (resource-name).
	// +optional
	// +kubebuilder:validation:Enum:=ip-name;resource-name
	PrivateDNSHostnameTypeOnLaunch *string `json:"privateDnsHostnameTypeOnLaunch,omitempty"`

	// ElasticIPPool contains specific configuration to allocate Public IPv4 address (Elastic IP) from user-defined pool
	// brought to AWS for core infrastructure resources, like NAT Gateways and Public Network Load Balancers for
	// the API Server.
	// +optional
	ElasticIPPool *ElasticIPPool `json:"elasticIpPool,omitempty"`

	// SubnetSchema specifies how CidrBlock should be divided on subnets in the VPC depending on the number of AZs.
	// PreferPrivate - one private subnet for each AZ plus one other subnet that will be further sub-divided for the public subnets.
	// PreferPublic - have the reverse logic of PreferPrivate, one public subnet for each AZ plus one other subnet
	// that will be further sub-divided for the private subnets.
	// Defaults to PreferPrivate
	// +optional
	// +kubebuilder:default=PreferPrivate
	// +kubebuilder:validation:Enum=PreferPrivate;PreferPublic
	SubnetSchema *SubnetSchemaType `json:"subnetSchema,omitempty"`
}

// String returns a string representation of the VPC.
func (v *VPCSpec) String() string {
	return fmt.Sprintf("id=%s", v.ID)
}

// IsUnmanaged returns true if the VPC is unmanaged.
func (v *VPCSpec) IsUnmanaged(clusterName string) bool {
	return v.ID != "" && !v.Tags.HasOwned(clusterName)
}

// IsManaged returns true if VPC is managed.
func (v *VPCSpec) IsManaged(clusterName string) bool {
	return !v.IsUnmanaged(clusterName)
}

// IsIPv6Enabled returns true if the IPv6 block is defined on the network spec.
func (v *VPCSpec) IsIPv6Enabled() bool {
	return v.IPv6 != nil
}

// GetElasticIPPool returns the custom Elastic IP Pool configuration when present.
func (v *VPCSpec) GetElasticIPPool() *ElasticIPPool {
	return v.ElasticIPPool
}

// GetPublicIpv4Pool returns the custom public IPv4 pool brought to AWS when present.
func (v *VPCSpec) GetPublicIpv4Pool() *string {
	if v.ElasticIPPool == nil {
		return nil
	}
	if v.ElasticIPPool.PublicIpv4Pool != nil {
		return v.ElasticIPPool.PublicIpv4Pool
	}
	return nil
}

// SubnetSpec configures an AWS Subnet.
type SubnetSpec struct {
	// ID defines a unique identifier to reference this resource.
	// If you're bringing your subnet, set the AWS subnet-id here, it must start with `subnet-`.
	//
	// When the VPC is managed by CAPA, and you'd like the provider to create a subnet for you,
	// the id can be set to any placeholder value that does not start with `subnet-`;
	// upon creation, the subnet AWS identifier will be populated in the `ResourceID` field and
	// the `id` field is going to be used as the subnet name. If you specify a tag
	// called `Name`, it takes precedence.
	ID string `json:"id"`

	// ResourceID is the subnet identifier from AWS, READ ONLY.
	// This field is populated when the provider manages the subnet.
	// +optional
	ResourceID string `json:"resourceID,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// IPv6CidrBlock is the IPv6 CIDR block to be used when the provider creates a managed VPC.
	// A subnet can have an IPv4 and an IPv6 address.
	// IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.
	// +optional
	IPv6CidrBlock string `json:"ipv6CidrBlock,omitempty"`

	// AvailabilityZone defines the availability zone to use for this subnet in the cluster's region.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// IsPublic defines the subnet as a public subnet. A subnet is public when it is associated with a route table that has a route to an internet gateway.
	// +optional
	IsPublic bool `json:"isPublic"`

	// IsIPv6 defines the subnet as an IPv6 subnet. A subnet is IPv6 when it is associated with a VPC that has IPv6 enabled.
	// IPv6 is only supported in managed clusters, this field cannot be set on AWSCluster object.
	// +optional
	IsIPv6 bool `json:"isIpv6,omitempty"`

	// RouteTableID is the routing table id associated with the subnet.
	// +optional
	RouteTableID *string `json:"routeTableId,omitempty"`

	// NatGatewayID is the NAT gateway id associated with the subnet.
	// Ignored unless the subnet is managed by the provider, in which case this is set on the public subnet where the NAT gateway resides. It is then used to determine routes for private subnets in the same AZ as the public subnet.
	// +optional
	NatGatewayID *string `json:"natGatewayId,omitempty"`

	// Tags is a collection of tags describing the resource.
	Tags Tags `json:"tags,omitempty"`

	// ZoneType defines the type of the zone where the subnet is created.
	//
	// The valid values are availability-zone, local-zone, and wavelength-zone.
	//
	// Subnet with zone type availability-zone (regular) is always selected to create cluster
	// resources, like Load Balancers, NAT Gateways, Contol Plane nodes, etc.
	//
	// Subnet with zone type local-zone or wavelength-zone is not eligible to automatically create
	// regular cluster resources.
	//
	// The public subnet in availability-zone or local-zone is associated with regular public
	// route table with default route entry to a Internet Gateway.
	//
	// The public subnet in wavelength-zone is associated with a carrier public
	// route table with default route entry to a Carrier Gateway.
	//
	// The private subnet in the availability-zone is associated with a private route table with
	// the default route entry to a NAT Gateway created in that zone.
	//
	// The private subnet in the local-zone or wavelength-zone is associated with a private route table with
	// the default route entry re-using the NAT Gateway in the Region (preferred from the
	// parent zone, the zone type availability-zone in the region, or first table available).
	//
	// +kubebuilder:validation:Enum=availability-zone;local-zone;wavelength-zone
	// +optional
	ZoneType *ZoneType `json:"zoneType,omitempty"`

	// ParentZoneName is the zone name where the current subnet's zone is tied when
	// the zone is a Local Zone.
	//
	// The subnets in Local Zone or Wavelength Zone locations consume the ParentZoneName
	// to select the correct private route table to egress traffic to the internet.
	//
	// +optional
	ParentZoneName *string `json:"parentZoneName,omitempty"`
}

// GetResourceID returns the identifier for this subnet,
// if the subnet was not created or reconciled, it returns the subnet ID.
func (s *SubnetSpec) GetResourceID() string {
	if s.ResourceID != "" {
		return s.ResourceID
	}
	return s.ID
}

// String returns a string representation of the subnet.
func (s *SubnetSpec) String() string {
	return fmt.Sprintf("id=%s/az=%s/public=%v", s.GetResourceID(), s.AvailabilityZone, s.IsPublic)
}

// IsEdge returns the true when the subnet is created in the edge zone,
// Local Zones.
func (s *SubnetSpec) IsEdge() bool {
	if s.ZoneType == nil {
		return false
	}
	if s.ZoneType.Equal(ZoneTypeLocalZone) {
		return true
	}
	if s.ZoneType.Equal(ZoneTypeWavelengthZone) {
		return true
	}
	return false
}

// IsEdgeWavelength returns true only when the subnet is created in Wavelength Zone.
func (s *SubnetSpec) IsEdgeWavelength() bool {
	if s.ZoneType == nil {
		return false
	}
	if *s.ZoneType == ZoneTypeWavelengthZone {
		return true
	}
	return false
}

// SetZoneInfo updates the subnets with zone information.
func (s *SubnetSpec) SetZoneInfo(zones []*ec2.AvailabilityZone) error {
	zoneInfo := func(zoneName string) *ec2.AvailabilityZone {
		for _, zone := range zones {
			if aws.StringValue(zone.ZoneName) == zoneName {
				return zone
			}
		}
		return nil
	}

	zone := zoneInfo(s.AvailabilityZone)
	if zone == nil {
		if len(s.AvailabilityZone) > 0 {
			return fmt.Errorf("unable to update zone information for subnet '%v' and zone '%v'", s.ID, s.AvailabilityZone)
		}
		return fmt.Errorf("unable to update zone information for subnet '%v'", s.ID)
	}
	if zone.ZoneType != nil {
		s.ZoneType = ptr.To(ZoneType(*zone.ZoneType))
	}
	if zone.ParentZoneName != nil {
		s.ParentZoneName = zone.ParentZoneName
	}
	return nil
}

// Subnets is a slice of Subnet.
// +listType=map
// +listMapKey=id
type Subnets []SubnetSpec

// ToMap returns a map from id to subnet.
func (s Subnets) ToMap() map[string]*SubnetSpec {
	res := make(map[string]*SubnetSpec)
	for i := range s {
		x := s[i]
		res[x.GetResourceID()] = &x
	}
	return res
}

// IDs returns a slice of the subnet ids.
func (s Subnets) IDs() []string {
	res := []string{}
	for _, subnet := range s {
		// Prevent returning edge zones (Local Zone) to regular Subnet IDs.
		// Edge zones should not deploy control plane nodes, and does not support Nat Gateway and
		// Network Load Balancers. Any resource for the core infrastructure should not consume edge
		// zones.
		if subnet.IsEdge() {
			continue
		}
		res = append(res, subnet.GetResourceID())
	}
	return res
}

// IDsWithEdge returns a slice of the subnet ids.
func (s Subnets) IDsWithEdge() []string {
	res := []string{}
	for _, subnet := range s {
		res = append(res, subnet.GetResourceID())
	}
	return res
}

// FindByID returns a single subnet matching the given id or nil.
//
// The returned pointer can be used to write back into the original slice.
func (s Subnets) FindByID(id string) *SubnetSpec {
	for i := range s {
		x := &(s[i]) // pointer to original structure
		if x.GetResourceID() == id {
			return x
		}
	}
	return nil
}

// FindEqual returns a subnet spec that is equal to the one passed in.
// Two subnets are defined equal to each other if their id is equal
// or if they are in the same vpc and the cidr block is the same.
//
// The returned pointer can be used to write back into the original slice.
func (s Subnets) FindEqual(spec *SubnetSpec) *SubnetSpec {
	for i := range s {
		x := &(s[i]) // pointer to original structure
		if (spec.GetResourceID() != "" && x.GetResourceID() == spec.GetResourceID()) ||
			(spec.CidrBlock == x.CidrBlock) ||
			(spec.IPv6CidrBlock != "" && spec.IPv6CidrBlock == x.IPv6CidrBlock) {
			return x
		}
	}
	return nil
}

// FilterPrivate returns a slice containing all subnets marked as private.
func (s Subnets) FilterPrivate() (res Subnets) {
	for _, x := range s {
		// Subnets in AWS Local Zones or Wavelength should not be used by core infrastructure.
		if x.IsEdge() {
			continue
		}
		if !x.IsPublic {
			res = append(res, x)
		}
	}
	return
}

// FilterNonCni returns the subnets that are NOT intended for usage with the CNI pod network
// (i.e. do NOT have the `sigs.k8s.io/cluster-api-provider-aws/association=secondary` tag).
func (s Subnets) FilterNonCni() (res Subnets) {
	for _, x := range s {
		if x.Tags[NameAWSSubnetAssociation] != SecondarySubnetTagValue {
			res = append(res, x)
		}
	}
	return
}

// FilterPublic returns a slice containing all subnets marked as public.
func (s Subnets) FilterPublic() (res Subnets) {
	for _, x := range s {
		// Subnets in AWS Local Zones or Wavelength should not be used by core infrastructure.
		if x.IsEdge() {
			continue
		}
		if x.IsPublic {
			res = append(res, x)
		}
	}
	return
}

// FilterByZone returns a slice containing all subnets that live in the availability zone specified.
func (s Subnets) FilterByZone(zone string) (res Subnets) {
	for _, x := range s {
		if x.AvailabilityZone == zone {
			res = append(res, x)
		}
	}
	return
}

// GetUniqueZones returns a slice containing the unique zones of the subnets.
func (s Subnets) GetUniqueZones() []string {
	keys := make(map[string]bool)
	zones := []string{}
	for _, x := range s {
		if _, value := keys[x.AvailabilityZone]; len(x.AvailabilityZone) > 0 && !value {
			keys[x.AvailabilityZone] = true
			zones = append(zones, x.AvailabilityZone)
		}
	}
	return zones
}

// SetZoneInfo updates the subnets with zone information.
func (s Subnets) SetZoneInfo(zones []*ec2.AvailabilityZone) error {
	for i := range s {
		if err := s[i].SetZoneInfo(zones); err != nil {
			return err
		}
	}
	return nil
}

// HasPublicSubnetWavelength returns true when there are subnets in Wavelength zone.
func (s Subnets) HasPublicSubnetWavelength() bool {
	for _, sub := range s {
		if sub.ZoneType == nil {
			return false
		}
		if sub.IsPublic && *sub.ZoneType == ZoneTypeWavelengthZone {
			return true
		}
	}
	return false
}

// CNISpec defines configuration for CNI.
type CNISpec struct {
	// CNIIngressRules specify rules to apply to control plane and worker node security groups.
	// The source for the rule will be set to control plane and worker security group IDs.
	CNIIngressRules CNIIngressRules `json:"cniIngressRules,omitempty"`
}

// CNIIngressRules is a slice of CNIIngressRule.
type CNIIngressRules []CNIIngressRule

// CNIIngressRule defines an AWS ingress rule for CNI requirements.
type CNIIngressRule struct {
	Description string                `json:"description"`
	Protocol    SecurityGroupProtocol `json:"protocol"`
	FromPort    int64                 `json:"fromPort"`
	ToPort      int64                 `json:"toPort"`
}

// RouteTable defines an AWS routing table.
type RouteTable struct {
	ID string `json:"id"`
}

// SecurityGroupRole defines the unique role of a security group.
// +kubebuilder:validation:Enum=bastion;node;controlplane;apiserver-lb;lb;node-eks-additional
type SecurityGroupRole string

var (
	// SecurityGroupBastion defines an SSH bastion role.
	SecurityGroupBastion = SecurityGroupRole("bastion")

	// SecurityGroupNode defines a Kubernetes workload node role.
	SecurityGroupNode = SecurityGroupRole("node")

	// SecurityGroupEKSNodeAdditional defines an extra node group from eks nodes.
	SecurityGroupEKSNodeAdditional = SecurityGroupRole("node-eks-additional")

	// SecurityGroupControlPlane defines a Kubernetes control plane node role.
	SecurityGroupControlPlane = SecurityGroupRole("controlplane")

	// SecurityGroupAPIServerLB defines a Kubernetes API Server Load Balancer role.
	SecurityGroupAPIServerLB = SecurityGroupRole("apiserver-lb")

	// SecurityGroupLB defines a container for the cloud provider to inject its load balancer ingress rules.
	SecurityGroupLB = SecurityGroupRole("lb")
)

// SecurityGroup defines an AWS security group.
type SecurityGroup struct {
	// ID is a unique identifier.
	ID string `json:"id"`

	// Name is the security group name.
	Name string `json:"name"`

	// IngressRules is the inbound rules associated with the security group.
	// +optional
	IngressRules IngressRules `json:"ingressRule,omitempty"`

	// Tags is a map of tags associated with the security group.
	Tags Tags `json:"tags,omitempty"`
}

// String returns a string representation of the security group.
func (s *SecurityGroup) String() string {
	return fmt.Sprintf("id=%s/name=%s", s.ID, s.Name)
}

// SecurityGroupProtocol defines the protocol type for a security group rule.
type SecurityGroupProtocol string

var (
	// SecurityGroupProtocolAll is a wildcard for all IP protocols.
	SecurityGroupProtocolAll = SecurityGroupProtocol("-1")

	// SecurityGroupProtocolIPinIP represents the IP in IP protocol in ingress rules.
	SecurityGroupProtocolIPinIP = SecurityGroupProtocol("4")

	// SecurityGroupProtocolTCP represents the TCP protocol in ingress rules.
	SecurityGroupProtocolTCP = SecurityGroupProtocol("tcp")

	// SecurityGroupProtocolUDP represents the UDP protocol in ingress rules.
	SecurityGroupProtocolUDP = SecurityGroupProtocol("udp")

	// SecurityGroupProtocolICMP represents the ICMP protocol in ingress rules.
	SecurityGroupProtocolICMP = SecurityGroupProtocol("icmp")

	// SecurityGroupProtocolICMPv6 represents the ICMPv6 protocol in ingress rules.
	SecurityGroupProtocolICMPv6 = SecurityGroupProtocol("58")

	// SecurityGroupProtocolESP represents the ESP protocol in ingress rules.
	SecurityGroupProtocolESP = SecurityGroupProtocol("50")
)

// IngressRule defines an AWS ingress rule for security groups.
type IngressRule struct {
	// Description provides extended information about the ingress rule.
	Description string `json:"description"`
	// Protocol is the protocol for the ingress rule. Accepted values are "-1" (all), "4" (IP in IP),"tcp", "udp", "icmp", and "58" (ICMPv6), "50" (ESP).
	// +kubebuilder:validation:Enum="-1";"4";tcp;udp;icmp;"58";"50"
	Protocol SecurityGroupProtocol `json:"protocol"`
	// FromPort is the start of port range.
	FromPort int64 `json:"fromPort"`
	// ToPort is the end of port range.
	ToPort int64 `json:"toPort"`

	// List of CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.
	// +optional
	CidrBlocks []string `json:"cidrBlocks,omitempty"`

	// List of IPv6 CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.
	// +optional
	IPv6CidrBlocks []string `json:"ipv6CidrBlocks,omitempty"`

	// The security group id to allow access from. Cannot be specified with CidrBlocks.
	// +optional
	SourceSecurityGroupIDs []string `json:"sourceSecurityGroupIds,omitempty"`

	// The security group role to allow access from. Cannot be specified with CidrBlocks.
	// The field will be combined with source security group IDs if specified.
	// +optional
	SourceSecurityGroupRoles []SecurityGroupRole `json:"sourceSecurityGroupRoles,omitempty"`

	// NatGatewaysIPsSource use the NAT gateways IPs as the source for the ingress rule.
	// +optional
	NatGatewaysIPsSource bool `json:"natGatewaysIPsSource,omitempty"`
}

// String returns a string representation of the ingress rule.
func (i IngressRule) String() string {
	return fmt.Sprintf("protocol=%s/range=[%d-%d]/description=%s", i.Protocol, i.FromPort, i.ToPort, i.Description)
}

// IngressRules is a slice of AWS ingress rules for security groups.
type IngressRules []IngressRule

// Difference returns the difference between this slice and the other slice.
func (i IngressRules) Difference(o IngressRules) (out IngressRules) {
	for index := range i {
		x := i[index]
		found := false
		for oIndex := range o {
			y := o[oIndex]
			if x.Equals(&y) {
				found = true
				break
			}
		}

		if !found {
			out = append(out, x)
		}
	}

	return
}

// Equals returns true if two IngressRule are equal.
func (i *IngressRule) Equals(o *IngressRule) bool {
	// ipv4
	if len(i.CidrBlocks) != len(o.CidrBlocks) {
		return false
	}

	sort.Strings(i.CidrBlocks)
	sort.Strings(o.CidrBlocks)

	for i, v := range i.CidrBlocks {
		if v != o.CidrBlocks[i] {
			return false
		}
	}
	// ipv6
	if len(i.IPv6CidrBlocks) != len(o.IPv6CidrBlocks) {
		return false
	}

	sort.Strings(i.IPv6CidrBlocks)
	sort.Strings(o.IPv6CidrBlocks)

	for i, v := range i.IPv6CidrBlocks {
		if v != o.IPv6CidrBlocks[i] {
			return false
		}
	}

	if len(i.SourceSecurityGroupIDs) != len(o.SourceSecurityGroupIDs) {
		return false
	}

	sort.Strings(i.SourceSecurityGroupIDs)
	sort.Strings(o.SourceSecurityGroupIDs)

	for i, v := range i.SourceSecurityGroupIDs {
		if v != o.SourceSecurityGroupIDs[i] {
			return false
		}
	}

	if i.Description != o.Description || i.Protocol != o.Protocol {
		return false
	}

	// AWS seems to ignore the From/To port when set on protocols where it doesn't apply, but
	// we avoid serializing it out for clarity's sake.
	// See: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_IpPermission.html
	switch i.Protocol {
	case SecurityGroupProtocolTCP,
		SecurityGroupProtocolUDP,
		SecurityGroupProtocolICMP,
		SecurityGroupProtocolICMPv6:
		return i.FromPort == o.FromPort && i.ToPort == o.ToPort
	case SecurityGroupProtocolAll, SecurityGroupProtocolIPinIP, SecurityGroupProtocolESP:
		// FromPort / ToPort are not applicable
	}

	return true
}

// ZoneType defines listener AWS Availability Zone type.
type ZoneType string

// String returns the string representation for the zone type.
func (z ZoneType) String() string {
	return string(z)
}

// Equal compares two zone types.
func (z ZoneType) Equal(other ZoneType) bool {
	return z == other
}

// ElasticIPPool allows configuring a Elastic IP pool for resources allocating
// public IPv4 addresses on public subnets.
type ElasticIPPool struct {
	// PublicIpv4Pool sets a custom Public IPv4 Pool used to create Elastic IP address for resources
	// created in public IPv4 subnets. Every IPv4 address, Elastic IP, will be allocated from the custom
	// Public IPv4 pool that you brought to AWS, instead of Amazon-provided pool. The public IPv4 pool
	// resource ID starts with 'ipv4pool-ec2'.
	//
	// +kubebuilder:validation:MaxLength=30
	// +optional
	PublicIpv4Pool *string `json:"publicIpv4Pool,omitempty"`

	// PublicIpv4PoolFallBackOrder defines the fallback action when the Public IPv4 Pool has been exhausted,
	// no more IPv4 address available in the pool.
	//
	// When set to 'amazon-pool', the controller check if the pool has available IPv4 address, when pool has reached the
	// IPv4 limit, the address will be claimed from Amazon-pool (default).
	//
	// When set to 'none', the controller will fail the Elastic IP allocation when the publicIpv4Pool is exhausted.
	//
	// +kubebuilder:validation:Enum:=amazon-pool;none
	// +optional
	PublicIpv4PoolFallBackOrder *PublicIpv4PoolFallbackOrder `json:"publicIpv4PoolFallbackOrder,omitempty"`

	// TODO(mtulio): add future support of user-defined Elastic IP to allow users to assign BYO Public IP from
	// 'static'/preallocated amazon-provided IPsstrucute currently holds only 'BYO Public IP from Public IPv4 Pool' (user brought to AWS),
	// although a dedicated structure would help to hold 'BYO Elastic IP' variants like:
	// - AllocationIdPoolApiLoadBalancer: an user-defined (static) IP address to the Public API Load Balancer.
	// - AllocationIdPoolNatGateways: an user-defined (static) IP address to allocate to NAT Gateways (egress traffic).
}

// PublicIpv4PoolFallbackOrder defines the list of available fallback action when the PublicIpv4Pool is exhausted.
// 'none' let the controllers return failures when the PublicIpv4Pool is exhausted - no more IPv4 available.
// 'amazon-pool' let the controllers to skip the PublicIpv4Pool and use the Amazon pool, the default.
// +kubebuilder:validation:XValidation:rule="self in ['none','amazon-pool']",message="allowed values are 'none' and 'amazon-pool'"
type PublicIpv4PoolFallbackOrder string

const (
	// PublicIpv4PoolFallbackOrderAmazonPool refers to use Amazon-pool Public IPv4 Pool as a fallback strategy.
	PublicIpv4PoolFallbackOrderAmazonPool = PublicIpv4PoolFallbackOrder("amazon-pool")

	// PublicIpv4PoolFallbackOrderNone refers to not use any fallback strategy.
	PublicIpv4PoolFallbackOrderNone = PublicIpv4PoolFallbackOrder("none")
)

func (r PublicIpv4PoolFallbackOrder) String() string {
	return string(r)
}

// Equal compares PublicIpv4PoolFallbackOrder types and return true if input param is equal.
func (r PublicIpv4PoolFallbackOrder) Equal(e PublicIpv4PoolFallbackOrder) bool {
	return r == e
}
