// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// VirtualMachineGroupMemberConditionGroupLinked indicates that the member
	// exists and has its "Spec.GroupName" field set to the group's name.
	VirtualMachineGroupMemberConditionGroupLinked = "GroupLinked"

	// VirtualMachineGroupMemberConditionPowerStateSynced indicates that the
	// member has been updated to match the group's power state.
	VirtualMachineGroupMemberConditionPowerStateSynced = "PowerStateSynced"

	// VirtualMachineGroupMemberConditionPlacementReady indicates that the
	// member has a placement decision ready.
	VirtualMachineGroupMemberConditionPlacementReady = "PlacementReady"
)

// GroupMember describes a member of a VirtualMachineGroup.
type GroupMember struct {
	// Name is the name of member of this group.
	Name string `json:"name"`

	// +optional
	// +kubebuilder:default=VirtualMachine
	// +kubebuilder:validation:Enum=VirtualMachine;VirtualMachineGroup

	// Kind is the kind of member of this group, which can be either
	// VirtualMachine or VirtualMachineGroup.
	//
	// If omitted, it defaults to VirtualMachine.
	Kind string `json:"kind,omitempty"`
}

// VirtualMachineGroupBootOrderGroup describes a boot order group within a
// VirtualMachineGroup.
type VirtualMachineGroupBootOrderGroup struct {
	// +optional
	// +listType=map
	// +listMapKey=kind
	// +listMapKey=name

	// Members describes the names of VirtualMachine or VirtualMachineGroup
	// objects that are members of this boot order group. The VM or VM Group
	// objects must be in the same namespace as this group.
	Members []GroupMember `json:"members,omitempty"`

	// +optional

	// PowerOnDelay is the amount of time to wait before powering on all the
	// members of this boot order group.
	//
	// If omitted, the members will be powered on immediately when the group's
	// power state changes to PoweredOn.
	PowerOnDelay *metav1.Duration `json:"powerOnDelay,omitempty"`
}

// VirtualMachineGroupSpec defines the desired state of VirtualMachineGroup.
type VirtualMachineGroupSpec struct {
	// +optional

	// GroupName describes the name of the group that this group belongs to.
	//
	// When this field is set to a valid group that contains this VM Group as a
	// member, an owner reference to that group is added to this VM Group.
	//
	// When this field is deleted or changed, any existing owner reference to
	// the previous group will be removed from this VM Group.
	GroupName string `json:"groupName,omitempty"`

	// +optional

	// BootOrder describes the boot sequence for this group members. Each boot
	// order contains a set of members that will be powered on simultaneously,
	// with an optional delay before powering on. The orders are processed
	// sequentially in the order they appear in this list, with delays being
	// cumulative across orders.
	//
	// When powering off, all members are stopped immediately without delays.
	BootOrder []VirtualMachineGroupBootOrderGroup `json:"bootOrder,omitempty"`

	// +optional

	// PowerState describes the desired power state of a VirtualMachineGroup.
	//
	// Please note this field may be omitted when creating a new VM group. This
	// ensures that the power states of any existing VMs that are added to the
	// group do not have their power states changed until the group's power
	// state is explicitly altered.
	//
	// However, once the field is set to a non-empty value, it may no longer be
	// set to an empty value. This means that if the group's power state is
	// PoweredOn, and a VM whose power state is PoweredOff is added to the
	// group, that VM will be powered on.
	PowerState VirtualMachinePowerState `json:"powerState,omitempty"`

	// +optional

	// NextForcePowerStateSyncTime may be used to force sync the power state of
	// the group to all of its members, by setting the value of this field to
	// "now" (case-insensitive).
	//
	// A mutating webhook changes this value to the current time (UTC), which
	// the VM Group controller then uses to trigger a sync of the group's power
	// state to its members.
	//
	// Please note it is not possible to schedule future syncs using this field.
	// The only value that users may set is the string "now" (case-insensitive).
	NextForcePowerStateSyncTime string `json:"nextForcePowerStateSyncTime,omitempty"`

	// +optional

	// PowerOffMode describes the desired behavior when powering off a VM Group.
	// Refer to the VirtualMachine.PowerOffMode field for more details.
	//
	// Please note this field is only propagated to the group's members when
	// the group's power state is changed or the nextForcePowerStateSyncTime
	// field is set to "now".
	PowerOffMode VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// +optional

	// SuspendMode describes the desired behavior when suspending a VM Group.
	// Refer to the VirtualMachine.SuspendMode field for more details.
	//
	// Please note this field is only propagated to the group's members when
	// the group's power state is changed or the nextForcePowerStateSyncTime
	// field is set to "now".
	SuspendMode VirtualMachinePowerOpMode `json:"suspendMode,omitempty"`
}

type VirtualMachineGroupPlacementDatastoreStatus struct {
	// Name describes the name of a datastore.
	Name string `json:"name"`

	// ID describes the datastore ID.
	ID string `json:"id,omitempty"`

	// URL describes the datastore URL.
	URL string `json:"url,omitempty"`

	// +optional

	// SupportedDiskFormat describes the list of disk formats supported by this
	// datastore.
	SupportedDiskFormats []string `json:"supportedDiskFormats,omitempty"`

	// +optional

	// DiskKey describes the device key to which this recommendation applies.
	// When omitted, this recommendation is for the VM's home directory.
	DiskKey *int32 `json:"diskKey,omitempty"`

	// +optional

	// TopLevelDirectoryCreateSupported indicates whether or not the datastore
	// supports creating a top-level directory or requires the use of the
	// namespace manager (i.e. vSAN).
	TopLevelDirectoryCreateSupported bool `json:"topLevelDirectoryCreateSupported,omitempty"`
}

type VirtualMachinePlacementStatus struct {
	// +optional

	// Zone describes the recommended zone for this VM.
	Zone string `json:"zoneID,omitempty"`

	// +optional

	// Node describes the recommended node for this VM.
	Node string `json:"node,omitempty"`

	// +optional

	// Pool describes the recommended resource pool for this VM.
	Pool string `json:"pool,omitempty"`

	// +optional

	// Datastores describe the recommended datastores for this VM.
	// This includes the recommendations for each of the VM's disks
	// and files.
	Datastores []VirtualMachineGroupPlacementDatastoreStatus `json:"datastores,omitempty"`
}

// VirtualMachineGroupMemberStatus describes the observed status of a group
// member.
type VirtualMachineGroupMemberStatus struct {
	// Name is the name of this member.
	Name string `json:"name"`

	// +kubebuilder:validation:Enum=VirtualMachine;VirtualMachineGroup

	// Kind is the kind of this member, which can be either VirtualMachine or
	// VirtualMachineGroup.
	Kind string `json:"kind"`

	// +optional

	// UID is the K8s metadata UID of this current member object.
	UID types.UID `json:"uid,omitempty"`

	// +optional

	// Placement describes the placement results for this member.
	//
	// Please note this field is only set for VirtualMachine members.
	Placement *VirtualMachinePlacementStatus `json:"placement,omitempty"`

	// +optional

	// PowerState describes the observed power state of this member.
	//
	// Please note this field is only set for VirtualMachine members.
	PowerState *VirtualMachinePowerState `json:"powerState,omitempty"`

	// +optional

	// Conditions describes any conditions associated with this member.
	//
	// - The GroupLinked condition is True when the member exists and has its
	//   "Spec.GroupName" field set to the group's name.
	// - The PowerStateSynced condition is True for the VirtualMachine member
	//   when the member's power state matches the group's power state.
	// - The PlacementReady condition is True for the VirtualMachine member
	//   when the member has a placement decision ready.
	// - The ReadyType condition is True for the VirtualMachineGroup member
	//   when all of its members' conditions are True.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// VirtualMachineGroupStatus defines the observed state of VirtualMachineGroup.
type VirtualMachineGroupStatus struct {
	// +optional
	// +listType=map
	// +listMapKey=name
	// +listMapKey=kind

	// Members describes the observed status of group members.
	Members []VirtualMachineGroupMemberStatus `json:"members,omitempty"`

	// +optional

	// LastUpdatedPowerStateTime describes the observed time when the power
	// state of the group was last updated.
	LastUpdatedPowerStateTime *metav1.Time `json:"lastUpdatedPowerStateTime,omitempty"`

	// +optional

	// Conditions describes any conditions associated with this VM Group.
	//
	// - The ReadyType condition is True when all of the group members have
	//   all of their expected conditions set to True.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmg
// +kubebuilder:subresource:status

// VirtualMachineGroup is the schema for the VirtualMachineGroup API and
// represents the desired state and observed status of a VirtualMachineGroup
// resource.
type VirtualMachineGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineGroupSpec   `json:"spec,omitempty"`
	Status VirtualMachineGroupStatus `json:"status,omitempty"`
}

func (vmg VirtualMachineGroup) GetMemberKind() string {
	return "VirtualMachineGroup"
}

func (vmg VirtualMachineGroup) GetGroupName() string {
	return vmg.Spec.GroupName
}

func (vmg *VirtualMachineGroup) SetGroupName(value string) {
	vmg.Spec.GroupName = value
}

func (vmg VirtualMachineGroup) GetPowerState() VirtualMachinePowerState {
	// VirtualMachineGroup does not have a power state recorded in its status.
	return ""
}

func (vmg *VirtualMachineGroup) SetPowerState(value VirtualMachinePowerState) {
	vmg.Spec.PowerState = value
}

func (vmg VirtualMachineGroup) GetConditions() []metav1.Condition {
	return vmg.Status.Conditions
}

func (vmg *VirtualMachineGroup) SetConditions(conditions []metav1.Condition) {
	vmg.Status.Conditions = conditions
}

func (m VirtualMachineGroupMemberStatus) GetConditions() []metav1.Condition {
	return m.Conditions
}

func (m *VirtualMachineGroupMemberStatus) SetConditions(conditions []metav1.Condition) {
	m.Conditions = conditions
}

// +kubebuilder:object:root=true

// VirtualMachineGroupList contains a list of VirtualMachineGroup.
type VirtualMachineGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineGroup `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VirtualMachineGroup{}, &VirtualMachineGroupList{})
}
