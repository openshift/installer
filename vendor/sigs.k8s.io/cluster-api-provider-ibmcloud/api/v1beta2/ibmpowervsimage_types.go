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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// IBMPowerVSImageFinalizer allows IBMPowerVSImageReconciler to clean up resources associated with IBMPowerVSImage before
	// removing it from the apiserver.
	IBMPowerVSImageFinalizer = "ibmpowervsimage.infrastructure.cluster.x-k8s.io"
)

// IBMPowerVSImageSpec defines the desired state of IBMPowerVSImage.
type IBMPowerVSImageSpec struct {

	// ClusterName is the name of the Cluster this object belongs to.
	// +kubebuilder:validation:MinLength=1
	ClusterName string `json:"clusterName"`

	// ServiceInstanceID is the id of the power cloud instance where the image will get imported.
	// Deprecated: use ServiceInstance instead
	ServiceInstanceID string `json:"serviceInstanceID"`

	// serviceInstance is the reference to the Power VS workspace on which the server instance(VM) will be created.
	// Power VS workspace is a container for all Power VS instances at a specific geographic region.
	// serviceInstance can be created via IBM Cloud catalog or CLI.
	// supported serviceInstance identifier in PowerVSResource are Name and ID and that can be obtained from IBM Cloud UI or IBM Cloud cli.
	// More detail about Power VS service instance.
	// https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server
	// when omitted system will dynamically create the service instance
	// +optional
	ServiceInstance *IBMPowerVSResourceReference `json:"serviceInstance,omitempty"`

	// Cloud Object Storage bucket name; bucket-name[/optional/folder]
	Bucket *string `json:"bucket"`

	// Cloud Object Storage image filename.
	Object *string `json:"object"`

	// Cloud Object Storage region.
	Region *string `json:"region"`

	// Type of storage, storage pool with the most available space will be selected.
	// +kubebuilder:default=tier1
	// +kubebuilder:validation:Enum=tier1;tier3
	// +optional
	StorageType string `json:"storageType,omitempty"`

	// DeletePolicy defines the policy used to identify images to be preserved beyond the lifecycle of associated cluster.
	// +kubebuilder:default=delete
	// +kubebuilder:validation:Enum=delete;retain
	// +optional
	DeletePolicy string `json:"deletePolicy,omitempty"`
}

// IBMPowerVSImageStatus defines the observed state of IBMPowerVSImage.
type IBMPowerVSImageStatus struct {

	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// ImageID is the id of the imported image.
	ImageID string `json:"imageID,omitempty"`

	// ImageState is the status of the imported image.
	// +optional
	ImageState PowerVSImageState `json:"imageState,omitempty"`

	// JobID is the job ID of an import operation.
	// +optional
	JobID string `json:"jobID,omitempty"`

	// Conditions defines current service state of the IBMPowerVSImage.
	// +optional
	Conditions capiv1beta1.Conditions `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.imageState",description="PowerVS image state"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Image is ready for IBM PowerVS instances"

// IBMPowerVSImage is the Schema for the ibmpowervsimages API.
type IBMPowerVSImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IBMPowerVSImageSpec   `json:"spec,omitempty"`
	Status IBMPowerVSImageStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the IBMPowerVSImage resource.
func (r *IBMPowerVSImage) GetConditions() capiv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the IBMPowerVSImage to the predescribed clusterv1.Conditions.
func (r *IBMPowerVSImage) SetConditions(conditions capiv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

//+kubebuilder:object:root=true

// IBMPowerVSImageList contains a list of IBMPowerVSImage.
type IBMPowerVSImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMPowerVSImage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IBMPowerVSImage{}, &IBMPowerVSImageList{})
}
