// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// VirtualMachineGroupPublishRequestConditionComplete is the Type for a
	// VirtualMachineGroupPublishRequest resource's status condition.
	//
	// The condition's status is set to true only when all other conditions
	// present on the resource have a truthy status.
	VirtualMachineGroupPublishRequestConditionComplete = "Complete"

	// VirtualMachineGroupPublishRequestConditionReasonPending indicates there are still pending
	// VirtualMachinePublishRequest to be completed.
	VirtualMachineGroupPublishRequestConditionReasonPending = "VirtualMachinePublishRequestsPending"
)

// VirtualMachineGroupPublishRequestSpec defines the desired state of a
// VirtualMachineGroupPublishRequest.
//
// All the fields in this spec are optional. This is especially useful when a
// DevOps persona wants to publish a VM Group without doing anything more than
// applying a VirtualMachineGroupPublishRequest resource that has the same name
// as said VMs in the same namespace as said VMs in VM Group.
type VirtualMachineGroupPublishRequestSpec struct {
	// +optional

	// Source is the name of the VirtualMachineGroup to be published.
	//
	// If this value is omitted then the publication controller checks to
	// see if there is a VirtualMachineGroup with the same name as this
	// VirtualMachineGroupPublishRequest resource. If such a resource exists,
	// then it is the source of the publication.
	Source string `json:"source,omitempty"`

	// +optional

	// VirtualMachines is a list of the VirtualMachine objects from the source
	// VirtualMachineGroup that are included in this publish request.
	//
	// If omitted, this field defaults to the names of all of the VMs currently
	// a member of the group, either directly or indirectly via a nested group.
	VirtualMachines []string `json:"virtualMachines,omitempty"`

	// +optional

	// Target is the name of the ContentLibrary resource to which the
	// VirtualMachines from the VirtualMachineGroup should be published.
	//
	// When this value is omitted, the controller attempts to identify the
	// target location by matching a ContentLibrary resource with the label
	// "imageregistry.vmware.com/default".
	//
	// Please note that while optional, if a VirtualMachineGroupPublishRequest
	// sans target information is applied to a namespace without a default
	// publication target, then the VirtualMachineGroupPublishRequest resource
	// will be marked in error.
	Target string `json:"target,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0

	// TTLSecondsAfterFinished is the time-to-live duration for how long this
	// resource will be allowed to exist once the publication operation
	// completes. After the TTL expires, the resource will be automatically
	// deleted without the user having to take any direct action.
	// This will be passed into each VirtualMachinePublishRequestSpec.
	//
	// If this field is unset then the request resource will not be
	// automatically deleted. If this field is set to zero then the request
	// resource is eligible for deletion immediately after it finishes.
	TTLSecondsAfterFinished *int64 `json:"ttlSecondsAfterFinished,omitempty"`
}

// VirtualMachineGroupPublishRequestStatus defines the observed state of a
// VirtualMachineGroupPublishRequest.
type VirtualMachineGroupPublishRequestStatus struct {
	// +optional

	// Source is the name of the published VirtualMachineGroup.
	Source string `json:"source,omitempty"`

	// +optional

	// Target is the name of the ContentLibrary to which the group is published.
	Target string `json:"target,omitempty"`

	// +optional

	// CompletionTime represents when the request was completed. It is not
	// guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	//
	// The value of this field should be equal to the value of the
	// LastTransitionTime for the status condition Type=Complete.
	CompletionTime metav1.Time `json:"completionTime,omitempty"`

	// +optional

	// StartTime represents when the request was acknowledged by the
	// controller. It is not guaranteed to be set in happens-before order
	// across separate operations. It is represented in RFC3339 form and is
	// in UTC.
	//
	// Please note that the group will not be published until the group's Ready
	// condition is true.
	StartTime metav1.Time `json:"startTime,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=source

	// Images describes the observed status of the individual VirtualMachine
	// publications.
	Images []VirtualMachineGroupPublishRequestImageStatus `json:"images,omitempty"`

	// +optional

	// Conditions is a list of the latest, available observations of the
	// request's current state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type VirtualMachineGroupPublishRequestImageStatus struct {
	// +required

	// Source is the name of the published VirtualMachine.
	Source string `json:"source"`

	// +required

	// PublishRequestName is the name of the VirtualMachinePublishRequest object
	// created to publish the VM.
	PublishRequestName string `json:"publishRequestName"`

	// +optional

	// ImageName is the name of the VirtualMachineImage resource that is
	// eventually realized after the publication operation completes.
	//
	// This field will not be set until the VirtualMachineImage resource
	// is realized.
	ImageName string `json:"imageName,omitempty"`

	// +optional

	// Conditions is a copy of the conditions from the
	// VirtualMachinePublishRequest object created to publish the VM.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmgpub
// +kubebuilder:subresource:status

// VirtualMachineGroupPublishRequest defines the information necessary to
// publish the VirtualMachines in a VirtualMachineGroup as VirtualMachineImages
// to an image registry.
type VirtualMachineGroupPublishRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineGroupPublishRequestSpec   `json:"spec,omitempty"`
	Status VirtualMachineGroupPublishRequestStatus `json:"status,omitempty"`
}

func (vmgrouppub *VirtualMachineGroupPublishRequest) GetConditions() []metav1.Condition {
	return vmgrouppub.Status.Conditions
}

func (vmgrouppub *VirtualMachineGroupPublishRequest) SetConditions(conditions []metav1.Condition) {
	vmgrouppub.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VirtualMachineGroupPublishRequestList contains a list of
// VirtualMachineGroupPublishRequest resources.
type VirtualMachineGroupPublishRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineGroupPublishRequest `json:"items"`
}

func init() {
	objectTypes = append(objectTypes,
		&VirtualMachineGroupPublishRequest{},
		&VirtualMachineGroupPublishRequestList{},
	)
}
