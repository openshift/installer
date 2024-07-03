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

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
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
}

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
}

// LoadBalancerSpec contains configuration for one or more LoadBalancers.
type LoadBalancerSpec struct {
	// APIServerInstanceGroupTagOverride overrides the default setting for the
	// tag used when creating the API Server Instance Group.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxLength=16
	// +kubebuilder:validation:Pattern=`(^[1-9][0-9]{0,31}$)|(^[a-z][a-z0-9-]{4,28}[a-z0-9]$)`
	// +optional
	APIServerInstanceGroupTagOverride *string `json:"apiServerInstanceGroupTagOverride,omitempty"`
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
