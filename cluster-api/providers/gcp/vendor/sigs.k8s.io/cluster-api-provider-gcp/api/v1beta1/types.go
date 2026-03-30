/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// GCPMachineTemplateResource describes the data needed to create am GCPMachine from a template.
type GCPMachineTemplateResource struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	ObjectMeta clusterv1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the desired behavior of the machine.
	Spec GCPMachineSpec `json:"spec"`
}

// Filter is a filter used to identify an GCP resource.
type Filter struct {
	// Name of the filter. Filter names are case-sensitive.
	Name string `json:"name"`

	// Values includes one or more filter values. Filter values are case-sensitive.
	Values []string `json:"values"`
}

// Network encapsulates GCP networking resources.
type Network struct {
	// SelfLink is the link to the Network used for this cluster.
	SelfLink *string `json:"selfLink,omitempty"`

	// FirewallRules is a map from the name of the rule to its full reference.
	// +optional
	FirewallRules map[string]string `json:"firewallRules,omitempty"`

	// Router is the full reference to the router created within the network
	// it'll contain the cloud nat gateway
	// +optional
	Router *string `json:"router,omitempty"`

	// APIServerAddress is the IPV4 global address assigned to the load balancer
	// created for the API Server.
	// +optional
	APIServerAddress *string `json:"apiServerIpAddress,omitempty"`

	// APIServerHealthCheck is the full reference to the health check
	// created for the API Server.
	// +optional
	APIServerHealthCheck *string `json:"apiServerHealthCheck,omitempty"`

	// APIServerInstanceGroups is a map from zone to the full reference
	// to the instance groups created for the control plane nodes created in the same zone.
	// +optional
	APIServerInstanceGroups map[string]string `json:"apiServerInstanceGroups,omitempty"`

	// APIServerBackendService is the full reference to the backend service
	// created for the API Server.
	// +optional
	APIServerBackendService *string `json:"apiServerBackendService,omitempty"`

	// APIServerTargetProxy is the full reference to the target proxy
	// created for the API Server.
	// +optional
	APIServerTargetProxy *string `json:"apiServerTargetProxy,omitempty"`

	// APIServerForwardingRule is the full reference to the forwarding rule
	// created for the API Server.
	// +optional
	APIServerForwardingRule *string `json:"apiServerForwardingRule,omitempty"`

	// APIInternalAddress is the IPV4 regional address assigned to the
	// internal Load Balancer.
	// +optional
	APIInternalAddress *string `json:"apiInternalIpAddress,omitempty"`

	// APIInternalHealthCheck is the full reference to the health check
	// created for the internal Load Balancer.
	// +optional
	APIInternalHealthCheck *string `json:"apiInternalHealthCheck,omitempty"`

	// APIInternalBackendService is the full reference to the backend service
	// created for the internal Load Balancer.
	// +optional
	APIInternalBackendService *string `json:"apiInternalBackendService,omitempty"`

	// APIInternalForwardingRule is the full reference to the forwarding rule
	// created for the internal Load Balancer.
	// +optional
	APIInternalForwardingRule *string `json:"apiInternalForwardingRule,omitempty"`
}

// FirewallSpec contains configuration for the firewall.
type FirewallSpec struct {
	// DefaultRulesManagement determines the management policy for the default firewall rules
	// created by the controller. DefaultRulesManagement has no effect on user specified firewall
	// rules. DefaultRulesManagement has no effect when a HostProject is specified.
	// "Managed": The controller will create and manage firewall rules.
	// "Unmanaged": The controller will not create or modify any firewall rules. If
	// the RulesManagement is changed from Managed to Unmanaged after the firewall rules
	// have been created, then the firewall rules will not be deleted.
	//
	// Defaults to "Managed".
	// +optional
	// +kubebuilder:default:="Managed"
	DefaultRulesManagement RulesManagementPolicy `json:"defaultRulesManagement,omitempty"`
}

// RulesManagementPolicy is a string enum type for managing firewall rules.
// +kubebuilder:validation:Enum=Managed;Unmanaged
type RulesManagementPolicy string

const (
	// RulesManagementManaged indicates that the controller should create and manage
	// firewall rules. This is the default behavior.
	RulesManagementManaged RulesManagementPolicy = "Managed"

	// RulesManagementUnmanaged indicates that the controller should not create or manage
	// any firewall rules. If rules already exist, they will be left as-is.
	RulesManagementUnmanaged RulesManagementPolicy = "Unmanaged"
)

// NetworkSpec encapsulates all things related to a GCP network.
type NetworkSpec struct {
	// Name is the name of the network to be used.
	// +optional
	Name *string `json:"name,omitempty"`

	// AutoCreateSubnetworks: When set to true, the VPC network is created
	// in "auto" mode. When set to false, the VPC network is created in
	// "custom" mode.
	//
	// An auto mode VPC network starts with one subnet per region. Each
	// subnet has a predetermined range as described in Auto mode VPC
	// network IP ranges.
	//
	// Defaults to true.
	// +optional
	AutoCreateSubnetworks *bool `json:"autoCreateSubnetworks,omitempty"`

	// Subnets configuration.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`

	// Allow for configuration of load balancer backend (useful for changing apiserver port)
	// +optional
	LoadBalancerBackendPort *int32 `json:"loadBalancerBackendPort,omitempty"`

	// HostProject is the name of the project hosting the shared VPC network resources.
	// +optional
	HostProject *string `json:"hostProject,omitempty"`

	// Firewall configuration.
	// +optional
	Firewall FirewallSpec `json:"firewall,omitempty,omitzero"`

	// Mtu: Maximum Transmission Unit in bytes. The minimum value for this field is
	// 1300 and the maximum value is 8896. The suggested value is 1500, which is
	// the default MTU used on the Internet, or 8896 if you want to use Jumbo
	// frames. If unspecified, the value defaults to 1460.
	// More info: https://pkg.go.dev/google.golang.org/api/compute/v1#Network
	// +kubebuilder:validation:Minimum:=1300
	// +kubebuilder:validation:Maximum:=8896
	// +kubebuilder:default:=1460
	// +optional
	Mtu int64 `json:"mtu,omitempty"`

	// MinPortsPerVM: Minimum number of ports allocated to a VM from this NAT
	// config. If not set, a default number of ports is allocated to a VM. This is
	// rounded up to the nearest power of 2. For example, if the value of this
	// field is 50, at least 64 ports are allocated to a VM.
	// +kubebuilder:validation:Minimum:=2
	// +kubebuilder:validation:Maximum:=65536
	// +kubebuilder:default:=64
	// +optional
	MinPortsPerVM int64 `json:"minPortsPerVm,omitempty"`
}

// LoadBalancerType defines the Load Balancer that should be created.
type LoadBalancerType string

var (
	// External creates a Global External Proxy Load Balancer
	// to manage traffic to backends in multiple regions. This is the default Load
	// Balancer and will be created if no LoadBalancerType is defined.
	External = LoadBalancerType("External")

	// Internal creates a Regional Internal Passthrough Load
	// Balancer to manage traffic to backends in the configured region.
	Internal = LoadBalancerType("Internal")

	// InternalExternal creates both External and Internal Load Balancers to provide
	// separate endpoints for managing both external and internal traffic.
	InternalExternal = LoadBalancerType("InternalExternal")
)

// LoadBalancerSpec contains configuration for one or more LoadBalancers.
type LoadBalancerSpec struct {
	// APIServerInstanceGroupTagOverride overrides the default setting for the
	// tag used when creating the API Server Instance Group.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=16
	// +kubebuilder:validation:Pattern=`(^[1-9][0-9]{0,31}$)|(^[a-z][a-z0-9-]{4,28}[a-z0-9]$)`
	// +optional
	APIServerInstanceGroupTagOverride *string `json:"apiServerInstanceGroupTagOverride,omitempty"`

	// LoadBalancerType defines the type of Load Balancer that should be created.
	// If not set, a Global External Proxy Load Balancer will be created by default.
	// +optional
	LoadBalancerType *LoadBalancerType `json:"loadBalancerType,omitempty"`

	// InternalLoadBalancer is the configuration for an Internal Passthrough Network Load Balancer.
	// +optional
	InternalLoadBalancer *LoadBalancer `json:"internalLoadBalancer,omitempty"`
}

// SubnetSpec configures an GCP Subnet.
type SubnetSpec struct {
	// Name defines a unique identifier to reference this resource.
	Name string `json:"name,omitempty"`

	// CidrBlock is the range of internal addresses that are owned by this
	// subnetwork. Provide this property when you create the subnetwork. For
	// example, 10.0.0.0/8 or 192.168.0.0/16. Ranges must be unique and
	// non-overlapping within a network. Only IPv4 is supported. This field
	// can be set only at resource creation time.
	CidrBlock string `json:"cidrBlock,omitempty"`

	// Description is an optional description associated with the resource.
	// +optional
	Description *string `json:"description,omitempty"`

	// SecondaryCidrBlocks defines secondary CIDR ranges,
	// from which secondary IP ranges of a VM may be allocated
	// +optional
	SecondaryCidrBlocks map[string]string `json:"secondaryCidrBlocks,omitempty"`

	// Region is the name of the region where the Subnetwork resides.
	Region string `json:"region,omitempty"`

	// PrivateGoogleAccess defines whether VMs in this subnet can access
	// Google services without assigning external IP addresses
	// +optional
	PrivateGoogleAccess *bool `json:"privateGoogleAccess,omitempty"`

	// EnableFlowLogs: Whether to enable flow logging for this subnetwork.
	// If this field is not explicitly set, it will not appear in get
	// listings. If not set the default behavior is to disable flow logging.
	// +optional
	EnableFlowLogs *bool `json:"enableFlowLogs,omitempty"`

	// Purpose: The purpose of the resource.
	// If unspecified, the purpose defaults to PRIVATE_RFC_1918.
	// The enableFlowLogs field isn't supported with the purpose field set to INTERNAL_HTTPS_LOAD_BALANCER.
	//
	// Possible values:
	//   "INTERNAL_HTTPS_LOAD_BALANCER" - Subnet reserved for Internal
	// HTTP(S) Load Balancing.
	//   "PRIVATE" - Regular user created or automatically created subnet.
	//   "PRIVATE_RFC_1918" - Regular user created or automatically created
	// subnet.
	//   "PRIVATE_SERVICE_CONNECT" - Subnetworks created for Private Service
	// Connect in the producer network.
	//   "REGIONAL_MANAGED_PROXY" - Subnetwork used for Regional
	// Internal/External HTTP(S) Load Balancing.
	// +kubebuilder:validation:Enum=INTERNAL_HTTPS_LOAD_BALANCER;PRIVATE_RFC_1918;PRIVATE;PRIVATE_SERVICE_CONNECT;REGIONAL_MANAGED_PROXY
	// +kubebuilder:default=PRIVATE_RFC_1918
	// +optional
	Purpose *string `json:"purpose,omitempty"`

	// StackType: The stack type for the subnet. If set to IPV4_ONLY, new VMs in
	// the subnet are assigned IPv4 addresses only. If set to IPV4_IPV6, new VMs in
	// the subnet can be assigned both IPv4 and IPv6 addresses. If not specified,
	// IPV4_ONLY is used. This field can be both set at resource creation time and
	// updated using patch.
	//
	// Possible values:
	//   "IPV4_IPV6" - New VMs in this subnet can have both IPv4 and IPv6
	// addresses.
	//   "IPV4_ONLY" - New VMs in this subnet will only be assigned IPv4 addresses.
	//   "IPV6_ONLY" - New VMs in this subnet will only be assigned IPv6 addresses.
	// +kubebuilder:validation:Enum=IPV4_ONLY;IPV4_IPV6;IPV6_ONLY
	// +kubebuilder:default=IPV4_ONLY
	// +optional
	StackType string `json:"stackType,omitempty"`
}

// String returns a string representation of the subnet.
func (s *SubnetSpec) String() string {
	return fmt.Sprintf("name=%s/region=%s", s.Name, s.Region)
}

// Subnets is a slice of Subnet.
type Subnets []SubnetSpec

// ToMap returns a map from name to subnet.
func (s Subnets) ToMap() map[string]*SubnetSpec {
	res := make(map[string]*SubnetSpec)
	for i := range s {
		x := s[i]
		res[x.Name] = &x
	}

	return res
}

// FindByName returns a single subnet matching the given name or nil.
func (s Subnets) FindByName(name string) *SubnetSpec {
	for _, x := range s {
		if x.Name == name {
			return &x
		}
	}

	return nil
}

// FilterByRegion returns a slice containing all subnets that live in the specified region.
func (s Subnets) FilterByRegion(region string) (res Subnets) {
	for _, x := range s {
		if x.Region == region {
			res = append(res, x)
		}
	}

	return
}

// InstanceStatus describes the state of an GCP instance.
type InstanceStatus string

var (
	// InstanceStatusProvisioning is the string representing an instance in a provisioning state.
	InstanceStatusProvisioning = InstanceStatus("PROVISIONING")

	// InstanceStatusRepairing is the string representing an instance in a repairing state.
	InstanceStatusRepairing = InstanceStatus("REPAIRING")

	// InstanceStatusRunning is the string representing an instance in a pending state.
	InstanceStatusRunning = InstanceStatus("RUNNING")

	// InstanceStatusStaging is the string representing an instance in a staging state.
	InstanceStatusStaging = InstanceStatus("STAGING")

	// InstanceStatusStopped is the string representing an instance
	// that has been stopped and can be restarted.
	InstanceStatusStopped = InstanceStatus("STOPPED")

	// InstanceStatusStopping is the string representing an instance
	// that is in the process of being stopped and can be restarted.
	InstanceStatusStopping = InstanceStatus("STOPPING")

	// InstanceStatusSuspended is the string representing an instance
	// that is suspended.
	InstanceStatusSuspended = InstanceStatus("SUSPENDED")

	// InstanceStatusSuspending is the string representing an instance
	// that is in the process of being suspended.
	InstanceStatusSuspending = InstanceStatus("SUSPENDING")

	// InstanceStatusTerminated is the string representing an instance that has been terminated.
	InstanceStatusTerminated = InstanceStatus("TERMINATED")
)

// ServiceAccount describes compute.serviceAccount.
type ServiceAccount struct {
	// Email: Email address of the service account.
	Email string `json:"email,omitempty"`

	// Scopes: The list of scopes to be made available for this service
	// account.
	Scopes []string `json:"scopes,omitempty"`
}

// ObjectReference is a reference to another Kubernetes object instance.
type ObjectReference struct {
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// InternalAccess defines the access for the Internal Passthrough Load Balancer.
type InternalAccess string

const (
	// InternalAccessRegional restricts traffic to clients within the same region as the internal load balancer.
	InternalAccessRegional = InternalAccess("Regional")

	// InternalAccessGlobal allows traffic from any region to access the internal load balancer.
	InternalAccessGlobal = InternalAccess("Global")
)

// LoadBalancer specifies the configuration of a LoadBalancer.
type LoadBalancer struct {
	// Name is the name of the Load Balancer. If not set a default name
	// will be used. For an Internal Load Balancer service the default
	// name is "api-internal".
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`(^[1-9][0-9]{0,31}$)|(^[a-z][a-z0-9-]{4,28}[a-z0-9]$)`
	// +optional
	Name *string `json:"name,omitempty"`

	// Subnet is the name of the subnet to use for a regional Load Balancer. A subnet is
	// required for the Load Balancer, if not defined the first configured subnet will be
	// used.
	Subnet *string `json:"subnet,omitempty"`

	// InternalAccess defines the access for the Internal Passthrough Load Balancer.
	// It determines whether the load balancer allows global access,
	// or restricts traffic to clients within the same region as the load balancer.
	// If unspecified, the value defaults to "Regional".
	//
	// Possible values:
	//   "Regional" - Only clients in the same region as the load balancer can access it.
	//   "Global" - Clients from any region can access the load balancer.
	// +kubebuilder:validation:Enum=Regional;Global
	// +kubebuilder:default=Regional
	// +optional
	InternalAccess InternalAccess `json:"internalAccess,omitempty"`

	// IPAddress is the static IP address to use for the Load Balancer.
	// If not set, a new static IP address will be allocated.
	// If set, it must be a valid free IP address from the LoadBalancer Subnet.
	// +optional
	IPAddress *string `json:"ipAddress,omitempty"`
}
