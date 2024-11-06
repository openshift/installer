/*
Copyright 2024 The Kubernetes Authors.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

const (
	// OpenStackServerFinalizer allows ReconcileOpenStackServer to clean up resources associated with OpenStackServer before
	// removing it from the apiserver.
	OpenStackServerFinalizer = "openstackserver.infrastructure.cluster.x-k8s.io"
)

// OpenStackServerSpec defines the desired state of OpenStackServer.
// +kubebuilder:validation:XValidation:message="at least one of flavor or flavorID must be set",rule=(has(self.flavor) || has(self.flavorID))
type OpenStackServerSpec struct {
	// AdditionalBlockDevices is a list of specifications for additional block devices to attach to the server instance.
	// +listType=map
	// +listMapKey=name
	// +optional
	AdditionalBlockDevices []infrav1.AdditionalBlockDevice `json:"additionalBlockDevices,omitempty"`

	// AvailabilityZone is the availability zone in which to create the server instance.
	//+optional
	AvailabilityZone optional.String `json:"availabilityZone,omitempty"`

	// ConfigDrive is a flag to enable config drive for the server instance.
	// +optional
	ConfigDrive optional.Bool `json:"configDrive,omitempty"`

	// The flavor reference for the flavor for the server instance.
	// +optional
	// +kubebuilder:validation:MinLength=1
	Flavor *string `json:"flavor,omitempty"`

	// FlavorID allows flavors to be specified by ID.  This field takes precedence
	// over Flavor.
	// +optional
	// +kubebuilder:validation:MinLength=1
	FlavorID *string `json:"flavorID,omitempty"`

	// FloatingIPPoolRef is a reference to a FloatingIPPool to allocate a floating IP from.
	// +optional
	FloatingIPPoolRef *corev1.TypedLocalObjectReference `json:"floatingIPPoolRef,omitempty"`

	// IdentityRef is a reference to a secret holding OpenStack credentials.
	// +required
	IdentityRef infrav1.OpenStackIdentityReference `json:"identityRef"`

	// The image to use for the server instance.
	// +required
	Image infrav1.ImageParam `json:"image"`

	// Ports to be attached to the server instance.
	// +required
	Ports []infrav1.PortOpts `json:"ports"`

	// RootVolume is the specification for the root volume of the server instance.
	// +optional
	RootVolume *infrav1.RootVolume `json:"rootVolume,omitempty"`

	// SSHKeyName is the name of the SSH key to inject in the instance.
	// +required
	SSHKeyName string `json:"sshKeyName"`

	// SecurityGroups is a list of security groups names to assign to the instance.
	// +optional
	SecurityGroups []infrav1.SecurityGroupParam `json:"securityGroups,omitempty"`

	// ServerGroup is the server group to which the server instance belongs.
	// +optional
	ServerGroup *infrav1.ServerGroupParam `json:"serverGroup,omitempty"`

	// ServerMetadata is a map of key value pairs to add to the server instance.
	// +listType=map
	// +listMapKey=key
	// +optional
	ServerMetadata []infrav1.ServerMetadata `json:"serverMetadata,omitempty"`

	// Tags which will be added to the machine and all dependent resources
	// which support them. These are in addition to Tags defined on the
	// cluster.
	// Requires Nova api 2.52 minimum!
	// +listType=set
	Tags []string `json:"tags,omitempty"`

	// Trunk is a flag to indicate if the server instance is created on a trunk port or not.
	// +optional
	Trunk optional.Bool `json:"trunk,omitempty"`

	// UserDataRef is a reference to a secret containing the user data to
	// be injected into the server instance.
	// +optional
	UserDataRef *corev1.LocalObjectReference `json:"userDataRef,omitempty"`

	// SchedulerHintAdditionalProperties are arbitrary key/value pairs that provide additional hints
	// to the OpenStack scheduler. These hints can influence how instances are placed on the infrastructure,
	// such as specifying certain host aggregates or availability zones.
	// +optional
	// +listType=map
	// +listMapKey=name
	SchedulerHintAdditionalProperties []infrav1.SchedulerHintAdditionalProperty `json:"schedulerHintAdditionalProperties,omitempty"`
}

// OpenStackServerStatus defines the observed state of OpenStackServer.
type OpenStackServerStatus struct {
	// Ready is true when the OpenStack server is ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// InstanceID is the ID of the server instance.
	// +optional
	InstanceID optional.String `json:"instanceID,omitempty"`

	// InstanceState is the state of the server instance.
	// +optional
	InstanceState *infrav1.InstanceState `json:"instanceState,omitempty"`

	// Addresses is the list of addresses of the server instance.
	// +optional
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// Resolved contains parts of the machine spec with all external
	// references fully resolved.
	// +optional
	Resolved *ResolvedServerSpec `json:"resolved,omitempty"`

	// Resources contains references to OpenStack resources created for the machine.
	// +optional
	Resources *ServerResources `json:"resources,omitempty"`

	// Conditions defines current service state of the OpenStackServer.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=openstackservers,scope=Namespaced,categories=cluster-api,shortName=oss
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="InstanceState",type="string",JSONPath=".status.instanceState",description="OpenStack instance state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="OpenStack instance ready status"
// +kubebuilder:printcolumn:name="InstanceID",type="string",JSONPath=".status.instanceID",description="OpenStack instance ID"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of OpenStack instance"

// OpenStackServer is the Schema for the openstackservers API.
type OpenStackServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackServerSpec   `json:"spec,omitempty"`
	Status OpenStackServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackServerList contains a list of OpenStackServer.
type OpenStackServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackServer `json:"items"`
}

// GetConditions returns the observations of the operational state of the OpenStackServer resource.
func (r *OpenStackServer) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the OpenStackServer to the predescribed clusterv1.Conditions.
func (r *OpenStackServer) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

var _ infrav1.IdentityRefProvider = &OpenStackFloatingIPPool{}

// GetIdentifyRef returns the Server's namespace and IdentityRef.
func (r *OpenStackServer) GetIdentityRef() (*string, *infrav1.OpenStackIdentityReference) {
	return &r.Namespace, &r.Spec.IdentityRef
}

func init() {
	SchemeBuilder.Register(&OpenStackServer{}, &OpenStackServerList{})
}
