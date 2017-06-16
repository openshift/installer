package types

import (
	"fmt"

	"k8s.io/client-go/pkg/api/unversioned"
	v1api "k8s.io/client-go/pkg/api/v1"
	v1beta1extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	// APIGroup of the Tectonic TPRs.
	TectonicAPIGroup = "coreos.com"
	// Namespace of the Tectonic TPRs and operators.
	TectonicNamespace = "tectonic-system"
)

const (
	// Version of the TCO Config TPR.
	TCOConfigTPRVersion = "v1"
	// Kind name of the TCO Confg TPR.
	TCOConfigTPRKind = "ChannelOperatorConfig"
)

const (
	// Version of the TectonicVersion TPR.
	TectonicVersionTPRVersion = "v1"
	// Kind name of the TectonicVersion TPR.
	TectonicVersionTPRKind = "TectonicVersion"
)

const (
	// Version of the AppVersion TPR.
	AppVersionTPRVersion = "v1"
	// Kind name of the AppVersion TPR.
	AppVersionTPRKind = "AppVersion"

	// Object name of the AppVersion TPR for the Tectonic cluster,
	// which contains the current/desired/target versions.
	AppVersionTPRNameTectonicCluster = "tectonic-cluster"
)

const (
	// States of each update task.

	// TaskStateNotStarted means the update task has started yet.
	// All update tasks will be set to this state at the start of the update process.
	TaskStateNotStarted TaskState = "NotStarted"

	// TaskStateUpdating means the update task is in progress.
	TaskStateRunning TaskState = "Running"

	// TaskStateCompleted means the update task is succesfully completed.
	TaskStateCompleted TaskState = "Completed"

	// TaskStateFailed means the update task failed.
	TaskStateFailed TaskState = "Failed"

	// TaskStateBackOff means the update task is in back-off because it
	// has failed last time and is going to try again.
	TaskStateBackOff TaskState = "BackOff"
)

const (
	// Some agreed field names in the AppVersion TPR.
	// A typical AppVersion TPR should look like:
	//
	// ```json
	// {
	//   "apiVersion": "coreos.com/v1",
	//   "kind": "AppVersion",
	//   "metadata": {
	//     "name": "tectonic-cluster"
	//   },
	//   "status": {
	//     "currentVersion": "1.4.3",
	//     "targetVersion: "1.4.4",
	//     "paused": false,
	//   },
	//   "spec": {
	//     "desiredVersion": "1.4.4",
	//     "paused": false
	//   }
	// }
	// ```

	JSONNameVersionStatus               = "status"
	JSONNameVersionStatusCurrentVersion = "currentVersion"
	JSONNameVersionStatusTargetVersion  = "targetVersion"
	JSONNameVersionStatusPaused         = "paused"

	JSONNameVersionSpec               = "spec"
	JSONNameVersionSpecDesiredVersion = "desiredVersion"
	JSONNameVersionSpecPaused         = "paused"
)

const (
	// The label on the resource that indicates the resource is managed
	// by the channel operator. Used for garbage collection.
	LabelKeyManagedByChannelOperator   = "managed-by-channel-operator"
	LabelValueManagedByChannelOperator = "true"
)

const (
	// Pre-defined failure types.
	FailureTypeUpdateFailed        FailureType = "Update failed" // This may be used to describe failures that can be recovered after restoring from backup.
	FailureTypeUpdateCannotProceed FailureType = "Update cannot proceed"
	FailureTypeHumanDecisionNeeded FailureType = "Human decision needed"
	FailureTypeVoidedwarranty      FailureType = "Voided warranty"
	FailureTypeUpdatesNotPossible  FailureType = "Updates are not possible"
)

const (
	TectonicEventRecorderName = "tectonic-channel-operator"
)

// ChannelOperatorConfig defines the config of
// the Tectonic Channel Operator.
// The config should be stored as a third party resource, and
// the operator will read it after it's launched.
type ChannelOperatorConfig struct {
	unversioned.TypeMeta `json:",inline"`
	v1api.ObjectMeta     `json:"metadata"`

	// Omaha server URL.
	Server string `json:"server"`

	// Channel to subscribe of the Omaha server.
	Channel string `json:"channel"`

	// AppID to use when subscribing the server.
	AppID string `json:"appID"`

	// RawURL points to the URL that has the update payload.
	// If RawURL is specified, the channel operator will ignore
	// the Server/Channel/AppID.
	RawURL string `json:"rawURL"`

	// Not being used now.
	PublicKey string `json:"publicKey"`

	// If true, then update will be triggered automatically.
	AutoUpdate bool `json:"automaticUpdate"`

	// If 'AutoUpdate' is false, then update will only be
	// triggered when this field is set to true.
	TriggerUpdate bool `json:"triggerUpdate"`

	// The interval for checking updates from core update server in seconds.
	UpdateCheckInterval int `json:"updateCheckInterval"`

	// If true, then check the update now instead of waiting for UpdateCheckInterval.
	TriggerUpdateCheck bool `json:"triggerUpdateCheck"`
}

// DesiredVersion defines the name of the Version AppVersion TPR for an
// operator, and its desired version.
type DesiredVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// TectonicVersion defines the update specs of different components.
// Tectonic Channel Operator expects to receive the spec from omaha server.
type TectonicVersion struct {
	unversioned.TypeMeta `json:",inline"`
	v1api.ObjectMeta     `json:"metadata"`

	// The new version of this update.
	Version string `json:"version"`

	// The new deployment specs for the operators.
	Deployments []v1beta1extensions.Deployment `json:"deployments"`

	// The new desired versions that will be consumed by the operators.
	DesiredVersions []DesiredVersion `json:"desiredVersions"`
}

// AppVersionSpec is the "spec" part of the AppVersion TPR.
type AppVersionSpec struct {
	DesiredVersion string `json:"desiredVersion"`
	Paused         bool   `json:"paused"`
}

// TaskState is the update state of an update task.
type TaskState string

// TaskStatus represents the status of an update task.
type TaskStatus struct {
	Name   string    `json:"name"`
	State  TaskState `json:"state"`
	Reason string    `json:"reason"`
}

// TaskStatusList represents a list of update task statuses.
type TaskStatusList []TaskStatus

// FailureType is a human readable string to for pre-defined failure types.
type FailureType string

// FailureStatus represents the failure information.
type FailureStatus struct {
	Type   FailureType `json:"type"`
	Reason string      `json:"reason"`
}

func (f FailureStatus) String() string {
	return fmt.Sprintf("%s: %s", f.Type, f.Reason)
}

// AppVersionStatus is the "status" part of the AppVersion TPR.
type AppVersionStatus struct {
	CurrentVersion string `json:"currentVersion"`
	TargetVersion  string `json:"targetVersion"`
	// If non-empty, then the upgrade is considered as a failure.
	// Detailed information is embeded in this field.
	FailureStatus *FailureStatus `json:"failureStatus,omitempty"`
	Paused        bool           `json:"paused"`
	TaskStatuses  TaskStatusList `json:"taskStatuses"`
}

// AppVersion represents the AppVersion TPR object, it only
// contains the required fields. It doesn't include operator specific
// fields. So update and write back this object directly might cause information
// loss.
type AppVersion struct {
	unversioned.TypeMeta `json:",inline"`
	v1api.ObjectMeta     `json:"metadata"`

	Spec   AppVersionSpec   `json:"spec"`
	Status AppVersionStatus `json:"status"`
}

// AppVersionList represents a list of AppVersion TPR objects that will
// be returned from a List() operation.
type AppVersionList struct {
	unversioned.TypeMeta `json:",inline"`
	v1api.ObjectMeta     `json:"metadata"`

	Items []AppVersion `json:"items"`
}

// TectonicVersionList represents a list of TectonicVersion TPR objects that
// will be returned from a List() operation.
type TectonicVersionList struct {
	unversioned.TypeMeta `json:",inline"`
	v1api.ObjectMeta     `json:"metadata"`

	Items []TectonicVersion `json:"items"`
}
