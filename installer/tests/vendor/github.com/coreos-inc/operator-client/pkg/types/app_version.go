package types

import (
	"fmt"

	"k8s.io/client-go/pkg/api/unversioned"
	"k8s.io/client-go/pkg/api/v1"
)

const (
	// TectonicAPIGroup is the APIGroup of the Tectonic TPRs.
	TectonicAPIGroup = "coreos.com"
	// TectonicNamespace is the namespace of the Tectonic TPRs and operators.
	TectonicNamespace = "tectonic-system"
)

const (
	// TectonicVersionTPRVersion is the version of the TectonicVersion TPR.
	TectonicVersionTPRVersion = "v1"
	// TectonicVersionTPRKind is the Kind name of the TectonicVersion TPR.
	TectonicVersionTPRKind = "TectonicVersion"
)

const (
	// AppVersionTPRVersion is the version of the AppVersion TPR.
	AppVersionTPRVersion = "v1"
	// AppVersionTPRKind is the name of the AppVersion TPR.
	AppVersionTPRKind = "AppVersion"

	// AppVersionTPRNameTectonicCluster is the Object name of the AppVersion
	// TPR for the Tectonic cluster,
	// which contains the current/desired/target versions.
	AppVersionTPRNameTectonicCluster = "tectonic-cluster"

	// AppVersionTPRNameKubernetes is the Object name of the AppVersion
	// TPR for the Kubernetes operator,
	// which contains the current/desired/target versions.
	AppVersionTPRNameKubernetes = "kubernetes"
)

// Pre-defined failure types.
const (
	// FailureTypeUpdateFailed represents an error that occured during an
	// update from which we cannot recover.
	FailureTypeUpdateFailed FailureType = "Update failed" // This may be used to describe failures that can be recovered after restoring from backup.
	// FailureTypeUpdateCannotProceed represents a failure not caused by
	// an operator failure, but that otherwise causes the update to not
	// proceed.
	FailureTypeUpdateCannotProceed FailureType = "Update cannot proceed"
	// FailureTypeHumanDecisionNeeded represents an error which must be
	// resolved by a human before the update can proceed.
	FailureTypeHumanDecisionNeeded FailureType = "Human decision needed"
	// FailureTypeVoidedwarranty represents a failure to update due to
	// modifications to the cluster that have voided the warranty.
	FailureTypeVoidedwarranty FailureType = "Voided warranty"
	// FailureTypeUpdatesNotPossible represents a failure type which expresses
	// that an update is not possible, such as updating from an unsupported
	// version.
	FailureTypeUpdatesNotPossible FailureType = "Updates are not possible"
)

// AppVersion represents the AppVersion TPR object, it only
// contains the required fields. It doesn't include operator specific
// fields. So update and write back this object directly might cause information
// loss.
type AppVersion struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata"`

	Spec   AppVersionSpec   `json:"spec"`
	Status AppVersionStatus `json:"status"`
}

// AppVersionSpec is the "spec" part of the AppVersion TPR.
type AppVersionSpec struct {
	DesiredVersion string `json:"desiredVersion"`
	Paused         bool   `json:"paused"`
}

// AppVersionStatus is the "status" part of the AppVersion TPR.
type AppVersionStatus struct {
	CurrentVersion string `json:"currentVersion"`
	TargetVersion  string `json:"targetVersion"`
	// If non-empty, then the upgrade is considered as a failure.
	// Detailed information is embeded in this field.
	FailureStatus *FailureStatus `json:"failureStatus,omitempty"`
	Paused        bool           `json:"paused"`
	TaskStatuses  []TaskStatus   `json:"taskStatuses"`
}

// TaskState is the update state of an update task.
type TaskState string

const (
	// States of each update task.

	// TaskStateNotStarted means the update task has started yet.
	// All update tasks will be set to this state at the start of the update process.
	TaskStateNotStarted TaskState = "NotStarted"

	// TaskStateRunning means the update task is in progress.
	TaskStateRunning TaskState = "Running"

	// TaskStateCompleted means the update task is succesfully completed.
	TaskStateCompleted TaskState = "Completed"

	// TaskStateFailed means the update task failed.
	TaskStateFailed TaskState = "Failed"

	// TaskStateBackOff means the update task is in back-off because it
	// has failed last time and is going to try again.
	TaskStateBackOff TaskState = "BackOff"
)

// TaskStatus represents the status of an update task.
type TaskStatus struct {
	Name   string    `json:"name"`
	State  TaskState `json:"state"`
	Reason string    `json:"reason"`
}

// FailureType is a human readable string to for pre-defined failure types.
type FailureType string

// FailureStatus represents the failure information.
type FailureStatus struct {
	Type   FailureType `json:"type"`
	Reason string      `json:"reason"`
}

// String returns a stringified failure status including Type and Reason.
func (f FailureStatus) String() string {
	return fmt.Sprintf("%s: %s", f.Type, f.Reason)
}

// AppVersionList represents a list of AppVersion TPR objects that will
// be returned from a List() operation.
type AppVersionList struct {
	unversioned.TypeMeta `json:",inline"`
	v1.ObjectMeta        `json:"metadata"`

	Items []AppVersion `json:"items"`
}

// AppVersionModifier is a modifier function to be used when atomically
// updating an AppVersion TPR.
type AppVersionModifier func(*AppVersion) error
