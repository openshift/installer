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

package v1alpha1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

const (
	// OpenStackFloatingIPPoolFinalizer allows ReconcileOpenStackFloatingIPPool to clean up resources associated with OpenStackFloatingIPPool before
	// removing it from the apiserver.
	OpenStackFloatingIPPoolFinalizer = "openstackfloatingippool.infrastructure.cluster.x-k8s.io"

	OpenStackFloatingIPPoolNameIndex = "spec.poolRef.name"

	// OpenStackFloatingIPPoolIP.
	DeleteFloatingIPFinalizer = "openstackfloatingippool.infrastructure.cluster.x-k8s.io/delete-floating-ip"
)

// ReclaimPolicy is a string type alias to represent reclaim policies for floating ips.
type ReclaimPolicy string

const (
	// ReclaimDelete is the reclaim policy for floating ips.
	ReclaimDelete ReclaimPolicy = "Delete"
	// ReclaimRetain is the reclaim policy for floating ips.
	ReclaimRetain ReclaimPolicy = "Retain"
)

// OpenStackFloatingIPPoolSpec defines the desired state of OpenStackFloatingIPPool.
type OpenStackFloatingIPPoolSpec struct {
	// PreAllocatedFloatingIPs is a list of floating IPs precreated in OpenStack that should be used by this pool.
	// These are used before allocating new ones and are not deleted from OpenStack when the pool is deleted.
	PreAllocatedFloatingIPs []string `json:"preAllocatedFloatingIPs,omitempty"`

	// MaxIPs is the maximum number of floating ips that can be allocated from this pool, if nil there is no limit.
	// If set, the pool will stop allocating floating ips when it reaches this number of ClaimedIPs.
	// +optional
	MaxIPs *int `json:"maxIPs,omitempty"`

	// IdentityRef is a reference to a identity to be used when reconciling this pool.
	// +kubebuilder:validation:Required
	IdentityRef infrav1.OpenStackIdentityReference `json:"identityRef"`

	// FloatingIPNetwork is the external network to use for floating ips, if there's only one external network it will be used by default
	// +optional
	FloatingIPNetwork *infrav1.NetworkParam `json:"floatingIPNetwork"`

	// The stratergy to use for reclaiming floating ips when they are released from a machine
	// +kubebuilder:validation:Enum=Retain;Delete
	ReclaimPolicy ReclaimPolicy `json:"reclaimPolicy"`
}

// OpenStackFloatingIPPoolStatus defines the observed state of OpenStackFloatingIPPool.
type OpenStackFloatingIPPoolStatus struct {
	// +kubebuilder:default={}
	// +optional
	ClaimedIPs []string `json:"claimedIPs"`

	// +kubebuilder:default={}
	// +optional
	AvailableIPs []string `json:"availableIPs"`

	// FailedIPs contains a list of floating ips that failed to be allocated
	// +optional
	FailedIPs []string `json:"failedIPs,omitempty"`

	// floatingIPNetwork contains information about the network used for floating ips
	// +optional
	FloatingIPNetwork *infrav1.NetworkStatus `json:"floatingIPNetwork,omitempty"`

	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
// +kubebuilder:storageversion
//+kubebuilder:subresource:status

// OpenStackFloatingIPPool is the Schema for the openstackfloatingippools API.
type OpenStackFloatingIPPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackFloatingIPPoolSpec   `json:"spec,omitempty"`
	Status OpenStackFloatingIPPoolStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenStackFloatingIPPoolList contains a list of OpenStackFloatingIPPool.
type OpenStackFloatingIPPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackFloatingIPPool `json:"items"`
}

// GetConditions returns the observations of the operational state of the OpenStackFloatingIPPool resource.
func (r *OpenStackFloatingIPPool) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the OpenStackFloatingIPPool to the predescribed clusterv1.Conditions.
func (r *OpenStackFloatingIPPool) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func (r *OpenStackFloatingIPPool) GetFloatingIPTag() string {
	return fmt.Sprintf("cluster-api-provider-openstack-fip-pool-%s", r.Name)
}

var _ infrav1.IdentityRefProvider = &OpenStackFloatingIPPool{}

// GetIdentifyRef returns the FloatingIPPool's namespace and IdentityRef.
func (r *OpenStackFloatingIPPool) GetIdentityRef() (*string, *infrav1.OpenStackIdentityReference) {
	return &r.Namespace, &r.Spec.IdentityRef
}

func init() {
	SchemeBuilder.Register(&OpenStackFloatingIPPool{}, &OpenStackFloatingIPPoolList{})
}
