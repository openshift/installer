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

package v1alpha3

import (
	"fmt"
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	clusterv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

const (
	// DefaultNameSuffix is the default suffix appended to all AWS IAM roles created by clusterawsadm.
	DefaultNameSuffix = ".cluster-api-provider-aws.sigs.k8s.io"
)

// AWSResourceReference is a reference to a specific AWS resource by ID, ARN, or filters.
// Only one of ID, ARN or Filters may be specified. Specifying more than one will result in
// a validation error.
type AWSResourceReference struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`

	// ARN of resource
	// +optional
	ARN *string `json:"arn,omitempty"`

	// Filters is a set of key/value pairs used to identify a resource
	// They are applied according to the rules defined by the AWS API:
	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Filtering.html
	// +optional
	Filters []Filter `json:"filters,omitempty"`
}

// AWSMachineTemplateResource describes the data needed to create am AWSMachine from a template
type AWSMachineTemplateResource struct {
	// Spec is the specification of the desired behavior of the machine.
	Spec AWSMachineSpec `json:"spec"`
}

// Filter is a filter used to identify an AWS resource
type Filter struct {
	// Name of the filter. Filter names are case-sensitive.
	Name string `json:"name"`

	// Values includes one or more filter values. Filter values are case-sensitive.
	Values []string `json:"values"`
}

// AWSMachineProviderConditionType is a valid value for AWSMachineProviderCondition.Type.
type AWSMachineProviderConditionType string

// Valid conditions for an AWS machine instance.
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated AWSMachineProviderConditionType = "MachineCreated"
)

// Network encapsulates AWS networking resources.
type Network struct {
	// SecurityGroups is a map from the role/kind of the security group to its unique name, if any.
	SecurityGroups map[SecurityGroupRole]SecurityGroup `json:"securityGroups,omitempty"`

	// APIServerELB is the Kubernetes api server classic load balancer.
	APIServerELB ClassicELB `json:"apiServerElb,omitempty"`
}

// ClassicELBScheme defines the scheme of a classic load balancer.
type ClassicELBScheme string

var (
	// ClassicELBSchemeInternetFacing defines an internet-facing, publicly
	// accessible AWS Classic ELB scheme.
	ClassicELBSchemeInternetFacing = ClassicELBScheme("internet-facing")

	// ClassicELBSchemeInternal defines an internal-only facing
	// load balancer internal to an ELB.
	ClassicELBSchemeInternal = ClassicELBScheme("internal")
)

func (e ClassicELBScheme) String() string {
	return string(e)
}

// ClassicELBProtocol defines listener protocols for a classic load balancer.
type ClassicELBProtocol string

var (
	// ClassicELBProtocolTCP defines the ELB API string representing the TCP protocol.
	ClassicELBProtocolTCP = ClassicELBProtocol("TCP")

	// ClassicELBProtocolSSL defines the ELB API string representing the TLS protocol.
	ClassicELBProtocolSSL = ClassicELBProtocol("SSL")

	// ClassicELBProtocolHTTP defines the ELB API string representing the HTTP protocol at L7.
	ClassicELBProtocolHTTP = ClassicELBProtocol("HTTP")

	// ClassicELBProtocolHTTPS defines the ELB API string representing the HTTP protocol at L7.
	ClassicELBProtocolHTTPS = ClassicELBProtocol("HTTPS")
)

// ClassicELB defines an AWS classic load balancer.
type ClassicELB struct {
	// The name of the load balancer. It must be unique within the set of load balancers
	// defined in the region. It also serves as identifier.
	Name string `json:"name,omitempty"`

	// DNSName is the dns name of the load balancer.
	DNSName string `json:"dnsName,omitempty"`

	// Scheme is the load balancer scheme, either internet-facing or private.
	Scheme ClassicELBScheme `json:"scheme,omitempty"`

	// AvailabilityZones is an array of availability zones in the VPC attached to the load balancer.
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// SubnetIDs is an array of subnets in the VPC attached to the load balancer.
	SubnetIDs []string `json:"subnetIds,omitempty"`

	// SecurityGroupIDs is an array of security groups assigned to the load balancer.
	SecurityGroupIDs []string `json:"securityGroupIds,omitempty"`

	// Listeners is an array of classic elb listeners associated with the load balancer. There must be at least one.
	Listeners []ClassicELBListener `json:"listeners,omitempty"`

	// HealthCheck is the classic elb health check associated with the load balancer.
	HealthCheck *ClassicELBHealthCheck `json:"healthChecks,omitempty"`

	// Attributes defines extra attributes associated with the load balancer.
	Attributes ClassicELBAttributes `json:"attributes,omitempty"`

	// Tags is a map of tags associated with the load balancer.
	Tags map[string]string `json:"tags,omitempty"`
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
	Protocol         ClassicELBProtocol `json:"protocol"`
	Port             int64              `json:"port"`
	InstanceProtocol ClassicELBProtocol `json:"instanceProtocol"`
	InstancePort     int64              `json:"instancePort"`
}

// ClassicELBHealthCheck defines an AWS classic load balancer health check.
type ClassicELBHealthCheck struct {
	Target             string        `json:"target"`
	Interval           time.Duration `json:"interval"`
	Timeout            time.Duration `json:"timeout"`
	HealthyThreshold   int64         `json:"healthyThreshold"`
	UnhealthyThreshold int64         `json:"unhealthyThreshold"`
}

// AZSelectionScheme defines the scheme of selecting AZs.
type AZSelectionScheme string

var (
	// AZSelectionSchemeOrdered will select AZs based on alphabetical order.
	AZSelectionSchemeOrdered = AZSelectionScheme("Ordered")

	// AZSelectionSchemeRandom will select AZs randomly.
	AZSelectionSchemeRandom = AZSelectionScheme("Random")
)

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
}

// VPCSpec configures an AWS VPC.
type VPCSpec struct {
	// ID is the vpc-id of the VPC this provider should use to create resources.
	ID string `json:"id,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	// Defaults to 10.0.0.0/16.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// InternetGatewayID is the id of the internet gateway associated with the VPC.
	// +optional
	InternetGatewayID *string `json:"internetGatewayId,omitempty"`

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

// SubnetSpec configures an AWS Subnet.
type SubnetSpec struct {
	// ID defines a unique identifier to reference this resource.
	ID string `json:"id,omitempty"`

	// CidrBlock is the CIDR block to be used when the provider creates a managed VPC.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// AvailabilityZone defines the availability zone to use for this subnet in the cluster's region.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// IsPublic defines the subnet as a public subnet. A subnet is public when it is associated with a route table that has a route to an internet gateway.
	// +optional
	IsPublic bool `json:"isPublic"`

	// RouteTableID is the routing table id associated with the subnet.
	// +optional
	RouteTableID *string `json:"routeTableId,omitempty"`

	// NatGatewayID is the NAT gateway id associated with the subnet.
	// Ignored unless the subnet is managed by the provider, in which case this is set on the public subnet where the NAT gateway resides. It is then used to determine routes for private subnets in the same AZ as the public subnet.
	// +optional
	NatGatewayID *string `json:"natGatewayId,omitempty"`

	// Tags is a collection of tags describing the resource.
	Tags Tags `json:"tags,omitempty"`
}

// String returns a string representation of the subnet.
func (s *SubnetSpec) String() string {
	return fmt.Sprintf("id=%s/az=%s/public=%v", s.ID, s.AvailabilityZone, s.IsPublic)
}

// Subnets is a slice of Subnet.
type Subnets []SubnetSpec

// ToMap returns a map from id to subnet.
func (s Subnets) ToMap() map[string]*SubnetSpec {
	res := make(map[string]*SubnetSpec)
	for i := range s {
		x := s[i]
		res[x.ID] = &x
	}
	return res
}

// IDs returns a slice of the subnet ids.
func (s Subnets) IDs() []string {
	res := []string{}
	for _, subnet := range s {
		res = append(res, subnet.ID)
	}
	return res
}

// FindByID returns a single subnet matching the given id or nil.
func (s Subnets) FindByID(id string) *SubnetSpec {
	for _, x := range s {
		if x.ID == id {
			return &x
		}
	}

	return nil
}

// FindEqual returns a subnet spec that is equal to the one passed in.
// Two subnets are defined equal to each other if their id is equal
// or if they are in the same vpc and the cidr block is the same.
func (s Subnets) FindEqual(spec *SubnetSpec) *SubnetSpec {
	for _, x := range s {
		if (spec.ID != "" && x.ID == spec.ID) || (spec.CidrBlock == x.CidrBlock) {
			return &x
		}
	}
	return nil
}

// FilterPrivate returns a slice containing all subnets marked as private.
func (s Subnets) FilterPrivate() (res Subnets) {
	for _, x := range s {
		if !x.IsPublic {
			res = append(res, x)
		}
	}
	return
}

// FilterPublic returns a slice containing all subnets marked as public.
func (s Subnets) FilterPublic() (res Subnets) {
	for _, x := range s {
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
		if _, value := keys[x.AvailabilityZone]; !value {
			keys[x.AvailabilityZone] = true
			zones = append(zones, x.AvailabilityZone)
		}
	}
	return zones
}

// CNISpec defines configuration for CNI.
type CNISpec struct {
	// CNIIngressRules specify rules to apply to control plane and worker node security groups.
	// The source for the rule will be set to control plane and worker security group IDs.
	CNIIngressRules CNIIngressRules `json:"cniIngressRules,omitempty"`
}

// CNIIngressRules is a slice of CNIIngressRule
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
)

// IngressRule defines an AWS ingress rule for security groups.
type IngressRule struct {
	Description string                `json:"description"`
	Protocol    SecurityGroupProtocol `json:"protocol"`
	FromPort    int64                 `json:"fromPort"`
	ToPort      int64                 `json:"toPort"`

	// List of CIDR blocks to allow access from. Cannot be specified with SourceSecurityGroupID.
	// +optional
	CidrBlocks []string `json:"cidrBlocks,omitempty"`

	// The security group id to allow access from. Cannot be specified with CidrBlocks.
	// +optional
	SourceSecurityGroupIDs []string `json:"sourceSecurityGroupIds,omitempty"`
}

// String returns a string representation of the ingress rule.
func (i *IngressRule) String() string {
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
	case SecurityGroupProtocolAll, SecurityGroupProtocolIPinIP:
		// FromPort / ToPort are not applicable
	}

	return true
}

// InstanceState describes the state of an AWS instance.
type InstanceState string

var (
	// InstanceStatePending is the string representing an instance in a pending state.
	InstanceStatePending = InstanceState("pending")

	// InstanceStateRunning is the string representing an instance in a running state.
	InstanceStateRunning = InstanceState("running")

	// InstanceStateShuttingDown is the string representing an instance shutting down.
	InstanceStateShuttingDown = InstanceState("shutting-down")

	// InstanceStateTerminated is the string representing an instance that has been terminated.
	InstanceStateTerminated = InstanceState("terminated")

	// InstanceStateStopping is the string representing an instance
	// that is in the process of being stopped and can be restarted.
	InstanceStateStopping = InstanceState("stopping")

	// InstanceStateStopped is the string representing an instance
	// that has been stopped and can be restarted.
	InstanceStateStopped = InstanceState("stopped")

	// InstanceRunningStates defines the set of states in which an EC2 instance is
	// running or going to be running soon.
	InstanceRunningStates = sets.NewString(
		string(InstanceStatePending),
		string(InstanceStateRunning),
	)

	// InstanceOperationalStates defines the set of states in which an EC2 instance is
	// or can return to running, and supports all EC2 operations.
	InstanceOperationalStates = InstanceRunningStates.Union(
		sets.NewString(
			string(InstanceStateStopping),
			string(InstanceStateStopped),
		),
	)

	// InstanceKnownStates represents all known EC2 instance states.
	InstanceKnownStates = InstanceOperationalStates.Union(
		sets.NewString(
			string(InstanceStateShuttingDown),
			string(InstanceStateTerminated),
		),
	)
)

// Instance describes an AWS instance.
type Instance struct {
	ID string `json:"id"`

	// The current state of the instance.
	State InstanceState `json:"instanceState,omitempty"`

	// The instance type.
	Type string `json:"type,omitempty"`

	// The ID of the subnet of the instance.
	SubnetID string `json:"subnetId,omitempty"`

	// The ID of the AMI used to launch the instance.
	ImageID string `json:"imageId,omitempty"`

	// The name of the SSH key pair.
	SSHKeyName *string `json:"sshKeyName,omitempty"`

	// SecurityGroupIDs are one or more security group IDs this instance belongs to.
	SecurityGroupIDs []string `json:"securityGroupIds,omitempty"`

	// UserData is the raw data script passed to the instance which is run upon bootstrap.
	// This field must not be base64 encoded and should only be used when running a new instance.
	UserData *string `json:"userData,omitempty"`

	// The name of the IAM instance profile associated with the instance, if applicable.
	IAMProfile string `json:"iamProfile,omitempty"`

	// Addresses contains the AWS instance associated addresses.
	Addresses []clusterv1alpha3.MachineAddress `json:"addresses,omitempty"`

	// The private IPv4 address assigned to the instance.
	PrivateIP *string `json:"privateIp,omitempty"`

	// The public IPv4 address assigned to the instance, if applicable.
	PublicIP *string `json:"publicIp,omitempty"`

	// Specifies whether enhanced networking with ENA is enabled.
	ENASupport *bool `json:"enaSupport,omitempty"`

	// Indicates whether the instance is optimized for Amazon EBS I/O.
	EBSOptimized *bool `json:"ebsOptimized,omitempty"`

	// Configuration options for the root storage volume.
	// +optional
	RootVolume *Volume `json:"rootVolume,omitempty"`

	// Configuration options for the non root storage volumes.
	// +optional
	NonRootVolumes []Volume `json:"nonRootVolumes,omitempty"`

	// Specifies ENIs attached to instance
	NetworkInterfaces []string `json:"networkInterfaces,omitempty"`

	// The tags associated with the instance.
	Tags map[string]string `json:"tags,omitempty"`

	// Availability zone of instance
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// SpotMarketOptions option for configuring instances to be run using AWS Spot instances.
	SpotMarketOptions *SpotMarketOptions `json:"spotMarketOptions,omitempty"`

	// Tenancy indicates if instance should run on shared or single-tenant hardware.
	// +optional
	Tenancy string `json:"tenancy,omitempty"`
}

// Volume encapsulates the configuration options for the storage device
type Volume struct {
	// Device name
	// +optional
	DeviceName string `json:"deviceName,omitempty"`

	// Size specifies size (in Gi) of the storage device.
	// Must be greater than the image snapshot size or 8 (whichever is greater).
	// +kubebuilder:validation:Minimum=8
	Size int64 `json:"size"`

	// Type is the type of the volume (e.g. gp2, io1, etc...).
	// +optional
	Type string `json:"type,omitempty"`

	// IOPS is the number of IOPS requested for the disk. Not applicable to all types.
	// +optional
	IOPS int64 `json:"iops,omitempty"`

	// Encrypted is whether the volume should be encrypted or not.
	// +optional
	Encrypted bool `json:"encrypted,omitempty"`

	// EncryptionKey is the KMS key to use to encrypt the volume. Can be either a KMS key ID or ARN.
	// If Encrypted is set and this is omitted, the default AWS key will be used.
	// The key must already exist and be accessible by the controller.
	// +optional
	EncryptionKey string `json:"encryptionKey,omitempty"`
}

// SpotMarketOptions defines the options available to a user when configuring
// Machines to run on Spot instances.
// Most users should provide an empty struct.
type SpotMarketOptions struct {
	// MaxPrice defines the maximum price the user is willing to pay for Spot VM instances
	// +optional
	// +kubebuilder:validation:pattern="^[0-9]+(\.[0-9]+)?$"
	MaxPrice *string `json:"maxPrice,omitempty"`
}
