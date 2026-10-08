// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// VirtualMachinePublishRequestConditionSourceValid is the Type for a
	// VirtualMachinePublishRequest resource's status condition.
	//
	// The condition's status is set to true only when the information
	// that describes the source side of the publication has been validated.
	VirtualMachinePublishRequestConditionSourceValid = "SourceValid"

	// VirtualMachinePublishRequestConditionTargetValid is the Type for a
	// VirtualMachinePublishRequest resource's status condition.
	//
	// The condition's status is set to true only when the information
	// that describes the target side of the publication has been
	// validated.
	VirtualMachinePublishRequestConditionTargetValid = "TargetValid"

	// VirtualMachinePublishRequestConditionUploaded is the Type for a
	// VirtualMachinePublishRequest resource's status condition.
	//
	// The condition's status is set to true only when the VM being
	// published has been successfully uploaded.
	VirtualMachinePublishRequestConditionUploaded = "Uploaded"

	// VirtualMachinePublishRequestConditionImageAvailable is the Type for a
	// VirtualMachinePublishRequest resource's status condition.
	//
	// The condition's status is set to true only when a new
	// VirtualMachineImage resource has been realized from the published
	// VM.
	VirtualMachinePublishRequestConditionImageAvailable = "ImageAvailable"

	// VirtualMachinePublishRequestConditionComplete is the Type for a
	// VirtualMachinePublishRequest resource's status condition.
	//
	// The condition's status is set to true only when all other conditions
	// present on the resource have a truthy status.
	VirtualMachinePublishRequestConditionComplete = "Complete"
)

// Condition.Reason for Conditions related to VirtualMachinePublishRequest.
const (
	// SourceVirtualMachineNotExistReason documents that the source VM of
	// the VirtualMachinePublishRequest doesn't exist.
	SourceVirtualMachineNotExistReason = "SourceVirtualMachineNotExist"

	// SourceVirtualMachineNotCreatedReason documents that the source VM of
	// the VirtualMachinePublishRequest hasn't been created.
	SourceVirtualMachineNotCreatedReason = "SourceVirtualMachineNotCreated"

	// SourceVirtualMachineUnsupportedEncryptionReason documents that the
	// source VM of the VirtualMachinePublishRequest has a mix of encrypted
	// and unencrypted disks.
	SourceVirtualMachineUnsupportedEncryptionReason = "SourceVirtualMachineUnsupportedEncryption"

	// TargetContentLibraryNotExistReason documents that the target content
	// library of the VirtualMachinePublishRequest doesn't exist.
	TargetContentLibraryNotExistReason = "TargetContentLibraryNotExist"

	// TargetContentLibraryNotWritableReason documents that the target content
	// library of the VirtualMachinePublishRequest isn't writable.
	TargetContentLibraryNotWritableReason = "TargetContentLibraryNotWritable"

	// TargetContentLibraryNotReadyReason documents that the target content
	// library of the VirtualMachinePublishRequest isn't ready.
	TargetContentLibraryNotReadyReason = "TargetContentLibraryNotReady"

	// TargetContentLibraryFailedQuotaCheckReason documents that the target
	// content library of the VirtualMachinePublishRequest failed quota
	// validation.
	TargetContentLibraryFailedQuotaCheckReason = "TargetContentLibraryFailedQuotaCheck"

	// TargetContentLibraryIncompatibleStorageClassReason documents that the
	// target content library of the VirtualMachinePublishRequest has a storage
	// class which is incompatible with the encryption status of the source VM.
	TargetContentLibraryIncompatibleStorageClassReason = "TargetContentLibraryIncompatibleStorageClass"

	// TargetItemAlreadyExistsReason documents that an item with the same name
	// as the VirtualMachinePublishRequest's target item name exists in
	// the target content library.
	TargetItemAlreadyExistsReason = "TargetItemAlreadyExists"

	// TargetVirtualMachineImageNotFoundReason documents that the expected
	// VirtualMachineImage resource corresponding to the VirtualMachinePublishRequest's
	// target item is not found in the namespace.
	TargetVirtualMachineImageNotFoundReason = "VirtualMachineImageNotFound"

	// UploadTaskNotStartedReason documents that the VM publish task hasn't started.
	UploadTaskNotStartedReason = "NotStarted"

	// UploadTaskQueuedReason documents that the VM publish task is in queued status.
	UploadTaskQueuedReason = "Queued"

	// UploadingReason documents that the VM publish task is in running status
	// and the published item is being uploaded to the target location.
	UploadingReason = "Uploading"

	// UploadItemIDInvalidReason documents that the VM publish task result
	// returns an invalid Item id.
	UploadItemIDInvalidReason = "ItemIDInvalid"

	// UploadFailureReason documents that uploading published item to the
	// target location failed.
	UploadFailureReason = "UploadFailure"

	// HasNotBeenUploadedReason documents that the VirtualMachinePublishRequest
	// hasn't completed because the published item hasn't been uploaded
	// to the target location.
	HasNotBeenUploadedReason = "HasNotBeenUploaded"

	// ImageUnavailableReason documents that the VirtualMachinePublishRequest
	// hasn't been completed because the expected VirtualMachineImage resource
	// isn't available yet.
	ImageUnavailableReason = "ImageUnavailable"

	// FatalReason documents that the outcome of the VirtualMachinePublishRequest
	// will not be retried.
	FatalReason = "Fatal"
)

const (
	// VirtualMachinePublishRequestManagedByLabelKey is the label key to identify VirtualMachinePublishRequest that
	// is created as a part of the VirtualMachineGroupPublishRequest.
	VirtualMachinePublishRequestManagedByLabelKey = GroupName + "/virtualmachinepublishrequest-managed-by"

	// VirtualMachinePublishRequestUUIDExtraConfigKey is the ExtraConfig key to store UUID of the
	// VirtualMachinePublishRequest that clones a source virtual machine, as a virtual machine template.
	VirtualMachinePublishRequestUUIDExtraConfigKey = "vmservice.vmpub.uuid"
)

// VirtualMachinePublishRequestSource is the source of a publication request,
// typically a VirtualMachine resource.
type VirtualMachinePublishRequestSource struct {
	// +optional

	// Name is the name of the referenced object.
	//
	// If omitted this value defaults to the name of the
	// VirtualMachinePublishRequest resource.
	Name string `json:"name,omitempty"`

	// +optional
	// +kubebuilder:default=vmoperator.vmware.com/v1alpha1

	// APIVersion is the API version of the referenced object.
	APIVersion string `json:"apiVersion,omitempty"`

	// +optional
	// +kubebuilder:default=VirtualMachine

	// Kind is the kind of referenced object.
	Kind string `json:"kind,omitempty"`
}

// VirtualMachinePublishRequestTargetItem is the item part of a
// publication request's target.
type VirtualMachinePublishRequestTargetItem struct {
	// +optional

	// Name is the name of the published object.
	//
	// If the spec.target.location.apiVersion equals
	// imageregistry.vmware.com/v1alpha1 and the spec.target.location.kind
	// equals ContentLibrary, then this should be the name that will
	// show up in vCenter Content Library, not the custom resource name
	// in the namespace.
	//
	// If omitted then the controller will use spec.source.name + "-image".
	Name string `json:"name,omitempty"`

	// +optional

	// Description is the description to assign to the published object.
	Description string `json:"description,omitempty"`
}

// VirtualMachinePublishRequestTargetLocation is the location part of a
// publication request's target.
type VirtualMachinePublishRequestTargetLocation struct {
	// +optional

	// Name is the name of the referenced object.
	//
	// Please note an error will be returned if this field is not
	// set in a namespace that lacks a default publication target.
	//
	// A default publication target is a resource with an API version
	// equal to spec.target.location.apiVersion, a kind equal to
	// spec.target.location.kind, and has the label
	// "imageregistry.vmware.com/default".
	Name string `json:"name,omitempty"`

	// +optional
	// +kubebuilder:default=imageregistry.vmware.com/v1alpha1

	// APIVersion is the API version of the referenced object.
	APIVersion string `json:"apiVersion,omitempty"`

	// +optional
	// +kubebuilder:default=ContentLibrary

	// Kind is the kind of referenced object.
	Kind string `json:"kind,omitempty"`
}

// VirtualMachinePublishRequestTarget is the target of a publication request,
// typically a ContentLibrary resource.
type VirtualMachinePublishRequestTarget struct {
	// +optional

	// Item contains information about the name of the object to which
	// the VM is published.
	//
	// Please note this value is optional and if omitted, the controller
	// will use spec.source.name + "-image" as the name of the published
	// item.
	Item VirtualMachinePublishRequestTargetItem `json:"item,omitempty"`

	// +optional

	// Location contains information about the location to which to publish
	// the VM.
	Location VirtualMachinePublishRequestTargetLocation `json:"location,omitempty"`
}

// VirtualMachinePublishRequestSpec defines the desired state of a
// VirtualMachinePublishRequest.
//
// All the fields in this spec are optional. This is especially useful when a
// DevOps persona wants to publish a VM without doing anything more than
// applying a VirtualMachinePublishRequest resource that has the same name
// as said VM in the same namespace as said VM.
type VirtualMachinePublishRequestSpec struct {
	// +optional

	// Source is the source of the publication request, ex. a VirtualMachine
	// resource.
	//
	// If this value is omitted then the publication controller checks to
	// see if there is a resource with the same name as this
	// VirtualMachinePublishRequest resource, an API version equal to
	// spec.source.apiVersion, and a kind equal to spec.source.kind. If such
	// a resource exists, then it is the source of the publication.
	Source VirtualMachinePublishRequestSource `json:"source,omitempty"`

	// +optional

	// Target is the target of the publication request, ex. item
	// information and a ContentLibrary resource.
	//
	// If this value is omitted, the controller uses spec.source.name + "-image"
	// as the name of the published item. Additionally, when omitted the
	// controller attempts to identify the target location by matching a
	// resource with an API version equal to spec.target.location.apiVersion, a
	// kind equal to spec.target.location.kind, w/ the label
	// "imageregistry.vmware.com/default".
	//
	// Please note that while optional, if a VirtualMachinePublishRequest sans
	// target information is applied to a namespace without a default
	// publication target, then the VirtualMachinePublishRequest resource
	// will be marked in error.
	Target VirtualMachinePublishRequestTarget `json:"target,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=0

	// TTLSecondsAfterFinished is the time-to-live duration for how long this
	// resource will be allowed to exist once the publication operation
	// completes. After the TTL expires, the resource will be automatically
	// deleted without the user having to take any direct action.
	//
	// If this field is unset then the request resource will not be
	// automatically deleted. If this field is set to zero then the request
	// resource is eligible for deletion immediately after it finishes.
	TTLSecondsAfterFinished *int64 `json:"ttlSecondsAfterFinished,omitempty"`

	// +kubebuilder:default=3
	// +kubebuilder:validation:Minimum=0

	// BackoffLimit is the number of status.attempts that should be allowed
	// before failing the VirtualMachinePublishRequest and halting any
	// future processing.
	BackoffLimit int64 `json:"backoffLimit,omitempty"`
}

// VirtualMachinePublishRequestStatus defines the observed state of a
// VirtualMachinePublishRequest.
type VirtualMachinePublishRequestStatus struct {
	// +optional

	// SourceRef is the reference to the source of the publication request,
	// ex. a VirtualMachine resource.
	SourceRef *VirtualMachinePublishRequestSource `json:"sourceRef,omitempty"`

	// +optional

	// TargetRef is the reference to the target of the publication request,
	// ex. item information and a ContentLibrary resource.
	TargetRef *VirtualMachinePublishRequestTarget `json:"targetRef,omitempty"`

	// +optional

	// CompletionTime represents time when the request was completed. It is not
	// guaranteed to be set in happens-before order across separate operations.
	// It is represented in RFC3339 form and is in UTC.
	//
	// The value of this field should be equal to the value of the
	// LastTransitionTime for the status condition Type=Complete.
	CompletionTime metav1.Time `json:"completionTime,omitempty"`

	// +optional

	// StartTime represents time when the request was acknowledged by the
	// controller. It is not guaranteed to be set in happens-before order
	// across separate operations. It is represented in RFC3339 form and is
	// in UTC.
	StartTime metav1.Time `json:"startTime,omitempty"`

	// +optional

	// Attempts represents the number of times the request to publish the VM
	// has been attempted.
	Attempts int64 `json:"attempts,omitempty"`

	// +optional

	// LastAttemptTime represents the time when the latest request was sent.
	LastAttemptTime metav1.Time `json:"lastAttemptTime,omitempty"`

	// +optional

	// ImageName is the name of the VirtualMachineImage resource that is
	// eventually realized in the same namespace as the VM and publication
	// request after the publication operation completes.
	//
	// This field will not be set until the VirtualMachineImage resource
	// is realized.
	ImageName string `json:"imageName,omitempty"`

	// +optional

	// Ready is set to true only when the VM has been published successfully
	// and the new VirtualMachineImage resource is ready.
	//
	// Readiness is determined by waiting until there is status condition
	// Type=Complete and ensuring it and all other status conditions present
	// have a Status=True. The conditions present will be:
	//
	//   * SourceValid
	//   * TargetValid
	//   * Uploaded
	//   * ImageAvailable
	//   * Complete
	Ready bool `json:"ready,omitempty"`

	// +optional

	// Conditions is a list of the latest, available observations of the
	// request's current state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmpub
// +kubebuilder:subresource:status

// VirtualMachinePublishRequest defines the information necessary to publish a
// VirtualMachine as a VirtualMachineImage to an image registry.
type VirtualMachinePublishRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachinePublishRequestSpec   `json:"spec,omitempty"`
	Status VirtualMachinePublishRequestStatus `json:"status,omitempty"`
}

func (vmpub *VirtualMachinePublishRequest) GetConditions() []metav1.Condition {
	return vmpub.Status.Conditions
}

func (vmpub *VirtualMachinePublishRequest) SetConditions(conditions []metav1.Condition) {
	vmpub.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VirtualMachinePublishRequestList contains a list of
// VirtualMachinePublishRequest resources.
type VirtualMachinePublishRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachinePublishRequest `json:"items"`
}

func init() {
	objectTypes = append(objectTypes,
		&VirtualMachinePublishRequest{},
		&VirtualMachinePublishRequestList{},
	)
}
