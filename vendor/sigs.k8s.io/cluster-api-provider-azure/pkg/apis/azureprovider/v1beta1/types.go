/*
Copyright 2019 The Kubernetes Authors.

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
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: Write type tests

// AzureResourceReference is a reference to a specific Azure resource by ID
type AzureResourceReference struct {
	// ID of resource
	// +optional
	ID *string `json:"id,omitempty"`
	// TODO: Investigate if we should reference resources in other ways
}

// TODO: Investigate resource filters

// AzureMachineProviderConditionType is a valid value for AzureMachineProviderCondition.Type
type AzureMachineProviderConditionType string

// Valid conditions for an Azure machine instance
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated AzureMachineProviderConditionType = "MachineCreated"
)

// AzureMachineProviderCondition is a condition in a AzureMachineProviderStatus
type AzureMachineProviderCondition struct {
	// Type is the type of the condition.
	Type AzureMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message"`
}

const (
	// ControlPlane machine label
	ControlPlane string = "master"
	// Node machine label
	Node string = "worker"
	// MachineRoleLabel machine label to determine the role
	MachineRoleLabel = "machine.openshift.io/cluster-api-machine-role"
)

// Network encapsulates Azure networking resources.
type Network struct {
	// SecurityGroups is a map from the role/kind of the security group to its unique name, if any.
	SecurityGroups map[SecurityGroupRole]*SecurityGroup `json:"securityGroups,omitempty"`

	// APIServerLB is the Kubernetes API server load balancer.
	APIServerLB LoadBalancer `json:"apiServerLb,omitempty"`

	// APIServerIP is the Kubernetes API server public IP address.
	APIServerIP PublicIP `json:"apiServerIp,omitempty"`
}

// TODO: Implement tagging
/*
// Tags defines resource tags.
type Tags map[string]*string
*/

// Subnets is a slice of Subnet.
type Subnets []*SubnetSpec

// TODO
// ToMap returns a map from id to subnet.
func (s Subnets) ToMap() map[string]*SubnetSpec {
	res := make(map[string]*SubnetSpec)
	for _, x := range s {
		res[x.ID] = x
	}
	return res
}

// SecurityGroupRole defines the unique role of a security group.
type SecurityGroupRole string

var (
	// SecurityGroupBastion defines an SSH bastion role
	SecurityGroupBastion = SecurityGroupRole("bastion")

	// SecurityGroupNode defines a Kubernetes workload node role
	SecurityGroupNode = SecurityGroupRole(Node)

	// SecurityGroupControlPlane defines a Kubernetes control plane node role
	SecurityGroupControlPlane = SecurityGroupRole(ControlPlane)
)

// SecurityGroup defines an Azure security group.
type SecurityGroup struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	IngressRules IngressRules `json:"ingressRule"`
	// TODO: Uncomment once tagging is implemented.
	//Tags         *Tags        `json:"tags"`
}

/*
// TODO
// String returns a string representation of the security group.
func (s *SecurityGroup) String() string {
	return fmt.Sprintf("id=%s/name=%s", s.ID, s.Name)
}
*/

// SecurityGroupProtocol defines the protocol type for a security group rule.
type SecurityGroupProtocol string

var (
	// SecurityGroupProtocolAll is a wildcard for all IP protocols
	SecurityGroupProtocolAll = SecurityGroupProtocol("*")

	// SecurityGroupProtocolTCP represents the TCP protocol in ingress rules
	SecurityGroupProtocolTCP = SecurityGroupProtocol("Tcp")

	// SecurityGroupProtocolUDP represents the UDP protocol in ingress rules
	SecurityGroupProtocolUDP = SecurityGroupProtocol("Udp")
)

// TODO
// IngressRule defines an Azure ingress rule for security groups.
type IngressRule struct {
	Description string                `json:"description"`
	Protocol    SecurityGroupProtocol `json:"protocol"`

	// SourcePorts - The source port or range. Integer or range between 0 and 65535. Asterix '*' can also be used to match all ports.
	SourcePorts *string `json:"sourcePorts,omitempty"`

	// DestinationPorts - The destination port or range. Integer or range between 0 and 65535. Asterix '*' can also be used to match all ports.
	DestinationPorts *string `json:"destinationPorts,omitempty"`

	// Source - The CIDR or source IP range. Asterix '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used. If this is an ingress rule, specifies where network traffic originates from.
	Source *string `json:"source,omitempty"`

	// Destination - The destination address prefix. CIDR or destination IP range. Asterix '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used.
	Destination *string `json:"destination,omitempty"`
}

// TODO
// String returns a string representation of the ingress rule.
/*
func (i *IngressRule) String() string {
	return fmt.Sprintf("protocol=%s/range=[%d-%d]/description=%s", i.Protocol, i.FromPort, i.ToPort, i.Description)
}
*/

// TODO
// IngressRules is a slice of Azure ingress rules for security groups.
type IngressRules []*IngressRule

// TODO
// Difference returns the difference between this slice and the other slice.
/*
func (i IngressRules) Difference(o IngressRules) (out IngressRules) {
	for _, x := range i {
		found := false
		for _, y := range o {
			sort.Strings(x.CidrBlocks)
			sort.Strings(y.CidrBlocks)
			sort.Strings(x.SourceSecurityGroupIDs)
			sort.Strings(y.SourceSecurityGroupIDs)
			if reflect.DeepEqual(x, y) {
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
*/

// PublicIP defines an Azure public IP address.
// TODO: Remove once load balancer is implemented.
type PublicIP struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	DNSName   string `json:"dnsName,omitempty"`
}

// TODO
// LoadBalancer defines an Azure load balancer.
type LoadBalancer struct {
	ID               string           `json:"id,omitempty"`
	Name             string           `json:"name,omitempty"`
	SKU              SKU              `json:"sku,omitempty"`
	FrontendIPConfig FrontendIPConfig `json:"frontendIpConfig,omitempty"`
	BackendPool      BackendPool      `json:"backendPool,omitempty"`
	// TODO: Uncomment once tagging is implemented.
	//Tags             Tags             `json:"tags,omitempty"`
	/*
		// FrontendIPConfigurations - Object representing the frontend IPs to be used for the load balancer
		FrontendIPConfigurations *[]FrontendIPConfiguration `json:"frontendIPConfigurations,omitempty"`
		// BackendAddressPools - Collection of backend address pools used by a load balancer
		BackendAddressPools *[]BackendAddressPool `json:"backendAddressPools,omitempty"`
		// LoadBalancingRules - Object collection representing the load balancing rules Gets the provisioning
		LoadBalancingRules *[]LoadBalancingRule `json:"loadBalancingRules,omitempty"`
		// Probes - Collection of probe objects used in the load balancer
		Probes *[]Probe `json:"probes,omitempty"`
		// InboundNatRules - Collection of inbound NAT Rules used by a load balancer. Defining inbound NAT rules on your load balancer is mutually exclusive with defining an inbound NAT pool. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an Inbound NAT pool. They have to reference individual inbound NAT rules.
		InboundNatRules *[]InboundNatRule `json:"inboundNatRules,omitempty"`
		// InboundNatPools - Defines an external port range for inbound NAT to a single backend port on NICs associated with a load balancer. Inbound NAT rules are created automatically for each NIC associated with the Load Balancer using an external port from this range. Defining an Inbound NAT pool on your Load Balancer is mutually exclusive with defining inbound Nat rules. Inbound NAT pools are referenced from virtual machine scale sets. NICs that are associated with individual virtual machines cannot reference an inbound NAT pool. They have to reference individual inbound NAT rules.
		InboundNatPools *[]InboundNatPool `json:"inboundNatPools,omitempty"`
		// OutboundRules - The outbound rules.
		OutboundRules *[]OutboundRule `json:"outboundRules,omitempty"`
		// ResourceGUID - The resource GUID property of the load balancer resource.
		ResourceGUID *string `json:"resourceGuid,omitempty"`
		// ProvisioningState - Gets the provisioning state of the PublicIP resource. Possible values are: 'Updating', 'Deleting', and 'Failed'.
		ProvisioningState *string `json:"provisioningState,omitempty"`
	*/
}

// LoadBalancerSKU enumerates the values for load balancer sku name.
type SKU string

var (
	SKUBasic    = SKU("Basic")
	SKUStandard = SKU("Standard")
)

type FrontendIPConfig struct {
	/*
		// FrontendIPConfigurationPropertiesFormat - Properties of the load balancer probe.
		*FrontendIPConfigurationPropertiesFormat `json:"properties,omitempty"`
		// Name - The name of the resource that is unique within a resource group. This name can be used to access the resource.
		Name *string `json:"name,omitempty"`
		// Etag - A unique read-only string that changes whenever the resource is updated.
		Etag *string `json:"etag,omitempty"`
		// Zones - A list of availability zones denoting the IP allocated for the resource needs to come from.
		Zones *[]string `json:"zones,omitempty"`
		// ID - Resource ID.
		ID *string `json:"id,omitempty"`
	*/
}

type BackendPool struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// TODO
// LoadBalancerProtocol defines listener protocols for a load balancer.
type LoadBalancerProtocol string

// TODO
var (
	// LoadBalancerProtocolTCP defines the LB API string representing the TCP protocol
	LoadBalancerProtocolTCP = LoadBalancerProtocol("TCP")

	// LoadBalancerProtocolSSL defines the LB API string representing the TLS protocol
	LoadBalancerProtocolSSL = LoadBalancerProtocol("SSL")

	// LoadBalancerProtocolHTTP defines the LB API string representing the HTTP protocol at L7
	LoadBalancerProtocolHTTP = LoadBalancerProtocol("HTTP")

	// LoadBalancerProtocolHTTPS defines the LB API string representing the HTTP protocol at L7
	LoadBalancerProtocolHTTPS = LoadBalancerProtocol("HTTPS")
)

// TODO
// LoadBalancerListener defines an Azure load balancer listener.
type LoadBalancerListener struct {
	Protocol         LoadBalancerProtocol `json:"protocol"`
	Port             int64                `json:"port"`
	InstanceProtocol LoadBalancerProtocol `json:"instanceProtocol"`
	InstancePort     int64                `json:"instancePort"`
}

// TODO
// LoadBalancerHealthCheck defines an Azure load balancer health check.
type LoadBalancerHealthCheck struct {
	Target             string        `json:"target"`
	Interval           time.Duration `json:"interval"`
	Timeout            time.Duration `json:"timeout"`
	HealthyThreshold   int64         `json:"healthyThreshold"`
	UnhealthyThreshold int64         `json:"unhealthyThreshold"`
}

// VMState describes the state of an Azure virtual machine.
type VMState string

var (
	// VMStateCreating ...
	VMStateCreating = VMState("Creating")
	// VMStateDeleting ...
	VMStateDeleting = VMState("Deleting")
	// VMStateFailed ...
	VMStateFailed = VMState("Failed")
	// VMStateMigrating ...
	VMStateMigrating = VMState("Migrating")
	// VMStateSucceeded ...
	VMStateSucceeded = VMState("Succeeded")
	// VMStateUpdating ...
	VMStateUpdating = VMState("Updating")

	// VMStateStarting ...
	VMStateStarting = VMState("Starting")
	// VMStateRunning ...
	VMStateRunning = VMState("Running")
	// VMStateStopping ...
	VMStateStopping = VMState("Stopping")
	// VMStateStopped ...
	VMStateStopped = VMState("Stopped")
	// VMStateDeallocating ...
	VMStateDeallocating = VMState("Deallocating")
	// VMStateDeallocated ...
	VMStateDeallocated = VMState("Deallocated")

	// VMStateUnknown ...
	VMStateUnknown = VMState("Unknown")
)

// VM describes an Azure virtual machine.
type VM struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`

	// Hardware profile
	VMSize string `json:"vmSize,omitempty"`

	// Storage profile
	Image  Image  `json:"image,omitempty"`
	OSDisk OSDisk `json:"osDisk,omitempty"`

	StartupScript string `json:"startupScript,omitempty"`

	// State - The provisioning state, which only appears in the response.
	State    VMState    `json:"vmState,omitempty"`
	Identity VMIdentity `json:"identity,omitempty"`

	// TODO: Uncomment once tagging is implemented.
	//Tags Tags `json:"tags,omitempty"`

	// HardwareProfile - Specifies the hardware settings for the virtual machine.
	//HardwareProfile *HardwareProfile `json:"hardwareProfile,omitempty"`

	// StorageProfile - Specifies the storage settings for the virtual machine disks.
	//StorageProfile *StorageProfile `json:"storageProfile,omitempty"`

	// AdditionalCapabilities - Specifies additional capabilities enabled or disabled on the virtual machine.
	//AdditionalCapabilities *AdditionalCapabilities `json:"additionalCapabilities,omitempty"`

	// OsProfile - Specifies the operating system settings for the virtual machine.
	//OsProfile *OSProfile `json:"osProfile,omitempty"`
	// NetworkProfile - Specifies the network interfaces of the virtual machine.
	//NetworkProfile *NetworkProfile `json:"networkProfile,omitempty"`

	//AvailabilitySet *SubResource `json:"availabilitySet,omitempty"`
}

// Image is a mirror of azure sdk compute.ImageReference
type Image struct {
	// Fields below refer to os images in marketplace
	Publisher string `json:"publisher"`
	Offer     string `json:"offer"`
	SKU       string `json:"sku"`
	Version   string `json:"version"`
	// ResourceID represents the location of OS Image in azure subscription
	ResourceID string `json:"resourceID"`
}

// VMIdentity defines the identity of the virtual machine, if configured.
type VMIdentity string

type OSDisk struct {
	OSType      string                `json:"osType"`
	ManagedDisk ManagedDiskParameters `json:"managedDisk"`
	DiskSizeGB  int32                 `json:"diskSizeGB"`
}

type ManagedDiskParameters struct {
	StorageAccountType string                       `json:"storageAccountType"`
	DiskEncryptionSet  *DiskEncryptionSetParameters `json:"diskEncryptionSet,omitempty"`
}

type DiskEncryptionSetParameters struct {
	ID string `json:"id,omitempty"`
}
