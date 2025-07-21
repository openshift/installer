/*
Copyright 2020.

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
	"github.com/openshift/assisted-service/api/common"
	"github.com/openshift/assisted-service/models"
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SpecSyncedCondition conditionsv1.ConditionType = "SpecSynced"

	ConnectedCondition      conditionsv1.ConditionType = "Connected"
	AgentConnectedReason    string                     = "AgentIsConnected"
	AgentDisconnectedReason string                     = "AgentIsDisconnected"
	AgentConnectedMsg       string                     = "The agent's connection to the installation service is unimpaired"
	AgentDisonnectedMsg     string                     = "The agent has not contacted the installation service in some time, user action should be taken"

	InstalledCondition conditionsv1.ConditionType = "Installed"

	RequirementsMetCondition       conditionsv1.ConditionType = "RequirementsMet"
	AgentReadyReason               string                     = "AgentIsReady"
	AgentReadyMsg                  string                     = "The agent is ready to begin the installation"
	AgentNotReadyReason            string                     = "AgentNotReady"
	AgentNotReadyMsg               string                     = "The agent is not ready to begin the installation"
	AgentAlreadyInstallingReason   string                     = "AgentAlreadyInstalling"
	AgentAlreadyInstallingMsg      string                     = "Installation already started and is in progress"
	AgentIsNotApprovedReason       string                     = "AgentIsNotApproved"
	AgentIsNotApprovedMsg          string                     = "The agent is not approved"
	AgentInstallationStoppedReason string                     = "AgentInstallationStopped"
	AgentInstallationStoppedMsg    string                     = "The agent installation stopped"

	ValidatedCondition             conditionsv1.ConditionType = "Validated"
	AgentValidationsPassingMsg     string                     = "The agent's validations are passing"
	AgentValidationsUnknownMsg     string                     = "The agent's validations have not yet been calculated"
	AgentValidationsFailingMsg     string                     = "The agent's validations are failing:"
	AgentValidationsUserPendingMsg string                     = "The agent's validations are pending for user:"

	InstalledReason              string = "InstallationCompleted"
	InstalledMsg                 string = "The installation has completed:"
	InstallationFailedReason     string = "InstallationFailed"
	InstallationFailedMsg        string = "The installation has failed:"
	InstallationNotStartedReason string = "InstallationNotStarted"
	InstallationNotStartedMsg    string = "The installation has not yet started"
	InstallationInProgressReason string = "InstallationInProgress"
	InstallationInProgressMsg    string = "The installation is in progress:"
	UnknownStatusReason          string = "UnknownStatus"
	UnknownStatusMsg             string = "The installation status is currently not recognized:"

	ValidationsPassingReason     string = "ValidationsPassing"
	ValidationsUnknownReason     string = "ValidationsUnknown"
	ValidationsFailingReason     string = "ValidationsFailing"
	ValidationsUserPendingReason string = "ValidationsUserPending"

	NotAvailableReason string = "NotAvailable"
	NotAvailableMsg    string = "Information not available"

	SyncedOkReason     string = "SyncOK"
	SyncedOkMsg        string = "The Spec has been successfully applied"
	BackendErrorReason string = "BackendError"
	BackendErrorMsg    string = "The Spec could not be synced due to backend error:"
	InputErrorReason   string = "InputError"
	InputErrorMsg      string = "The Spec could not be synced due to an input error:"

	BoundCondition                   conditionsv1.ConditionType = "Bound"
	BoundReason                      string                     = "Bound"
	BoundMsg                         string                     = "The agent is bound to a cluster deployment"
	UnboundReason                    string                     = "Unbound"
	UnboundMsg                       string                     = "The agent is not bound to any cluster deployment"
	BindingReason                    string                     = "Binding"
	BindingMsg                       string                     = "The agent is currently binding to a cluster deployment"
	UnbindingReason                  string                     = "Unbinding"
	UnbindingMsg                     string                     = "The agent is currently unbinding from a cluster deployment"
	UnbindingPendingUserActionReason string                     = "UnbindingPendingUserAction"
	UnbindingPendingUserActionMsg    string                     = "The agent is currently unbinding; Pending host reboot from infraenv image"

	CleanupCondition    conditionsv1.ConditionType = "Cleanup"
	CleanupFailedReason string                     = "CleanupFailed"
)

type HostMemory struct {
	PhysicalBytes int64 `json:"physicalBytes,omitempty"`
	UsableBytes   int64 `json:"usableBytes,omitempty"`
}

type HostCPU struct {
	Count int64 `json:"count,omitempty"`
	// Name in REST API: frequency
	ClockMegahertz int64    `json:"clockMegahertz,omitempty"`
	Flags          []string `json:"flags,omitempty"`
	ModelName      string   `json:"modelName,omitempty"`
	Architecture   string   `json:"architecture,omitempty"`
}

type HostInterface struct {
	IPV6Addresses []string `json:"ipV6Addresses"`
	Vendor        string   `json:"vendor,omitempty"`
	Name          string   `json:"name,omitempty"`
	HasCarrier    bool     `json:"hasCarrier,omitempty"`
	Product       string   `json:"product,omitempty"`
	Mtu           int64    `json:"mtu,omitempty"`
	IPV4Addresses []string `json:"ipV4Addresses"`
	Biosdevname   string   `json:"biosDevName,omitempty"`
	ClientId      string   `json:"clientID,omitempty"`
	MacAddress    string   `json:"macAddress,omitempty"`
	Flags         []string `json:"flags"`
	SpeedMbps     int64    `json:"speedMbps,omitempty"`
}

type HostInstallationEligibility struct {
	Eligible           bool     `json:"eligible,omitempty"`
	NotEligibleReasons []string `json:"notEligibleReasons"`
}

type HostIOPerf struct {
	// 99th percentile of fsync duration in milliseconds
	SyncDurationMilliseconds int64 `json:"syncDurationMilliseconds,omitempty"`
}

type HostDisk struct {
	ID                      string                      `json:"id"`
	DriveType               string                      `json:"driveType,omitempty"`
	Vendor                  string                      `json:"vendor,omitempty"`
	Name                    string                      `json:"name,omitempty"`
	Path                    string                      `json:"path,omitempty"`
	Hctl                    string                      `json:"hctl,omitempty"`
	ByPath                  string                      `json:"byPath,omitempty"`
	ByID                    string                      `json:"byID,omitempty"`
	Model                   string                      `json:"model,omitempty"`
	Wwn                     string                      `json:"wwn,omitempty"`
	Serial                  string                      `json:"serial,omitempty"`
	SizeBytes               int64                       `json:"sizeBytes,omitempty"`
	Bootable                bool                        `json:"bootable,omitempty"`
	Smart                   string                      `json:"smart,omitempty"`
	InstallationEligibility HostInstallationEligibility `json:"installationEligibility,omitempty"`
	IoPerf                  HostIOPerf                  `json:"ioPerf,omitempty"`
}

type HostBoot struct {
	CurrentBootMode string `json:"currentBootMode,omitempty"`
	PxeInterface    string `json:"pxeInterface,omitempty"`
	DeviceType      string `json:"deviceType,omitempty"`
}

type HostSystemVendor struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`
	Virtual      bool   `json:"virtual,omitempty"`
}

type HostInventory struct {
	// Name in REST API: timestamp
	ReportTime   *metav1.Time     `json:"reportTime,omitempty"`
	Hostname     string           `json:"hostname,omitempty"`
	BmcAddress   string           `json:"bmcAddress,omitempty"`
	BmcV6address string           `json:"bmcV6Address,omitempty"`
	Memory       HostMemory       `json:"memory,omitempty"`
	Cpu          HostCPU          `json:"cpu,omitempty"`
	Interfaces   []HostInterface  `json:"interfaces,omitempty"`
	Disks        []HostDisk       `json:"disks,omitempty"`
	Boot         HostBoot         `json:"boot,omitempty"`
	SystemVendor HostSystemVendor `json:"systemVendor,omitempty"`
}

// AgentSpec defines the desired state of Agent
type AgentSpec struct {
	// +optional
	ClusterDeploymentName *ClusterReference `json:"clusterDeploymentName,omitempty"`
	Role                  models.HostRole   `json:"role" protobuf:"bytes,1,opt,name=role,casttype=HostRole"`
	Hostname              string            `json:"hostname,omitempty"`
	MachineConfigPool     string            `json:"machineConfigPool,omitempty"`
	Approved              bool              `json:"approved"`
	// InstallationDiskID defines the installation destination disk (must be equal to the inventory disk id).
	InstallationDiskID string `json:"installation_disk_id,omitempty"`
	// Json formatted string containing the user overrides for the host's coreos installer args
	InstallerArgs string `json:"installerArgs,omitempty"`
	// Json formatted string containing the user overrides for the host's ignition config
	IgnitionConfigOverrides string `json:"ignitionConfigOverrides,omitempty"`
	// IgnitionEndpointTokenReference references a secret containing an Authorization Bearer token to fetch the ignition from ignition_endpoint_url.
	IgnitionEndpointTokenReference *IgnitionEndpointTokenReference `json:"ignitionEndpointTokenReference,omitempty"`
	// IgnitionEndpointHTTPHeaders are the additional HTTP headers used when fetching the ignition.
	IgnitionEndpointHTTPHeaders map[string]string `json:"ignitionEndpointHTTPHeaders,omitempty"`
	// NodeLabels are the labels to be applied on the node associated with this agent
	NodeLabels map[string]string `json:"nodeLabels,omitempty"`
}

type IgnitionEndpointTokenReference struct {
	// Namespace is the namespace of the secret containing the ignition endpoint token.
	Namespace string `json:"namespace"`
	// Name is the name of the secret containing the ignition endpoint token.
	Name string `json:"name"`
}

type HostProgressInfo struct {
	// current installation stage
	CurrentStage models.HostStage `json:"currentStage,omitempty"`
	// All stages (ordered by their appearance) for this agent
	ProgressStages []models.HostStage `json:"progressStages,omitempty"`
	// Additional information for the current installation stage
	ProgressInfo string `json:"progressInfo,omitempty"`
	// Estimate progress (percentage)
	InstallationPercentage int64 `json:"installationPercentage,omitempty"`
	// host field: progress: stage_started_at
	StageStartTime *metav1.Time `json:"stageStartTime,omitempty"`
	// host field: progress: stage_updated_at
	StageUpdateTime *metav1.Time `json:"stageUpdateTime,omitempty"`
}

type HostNTPSources struct {
	SourceName  string             `json:"sourceName,omitempty"`
	SourceState models.SourceState `json:"sourceState,omitempty"`
}

type AgentDeprovisionInfo struct {
	ClusterName      string `json:"cluster_name,omitempty"`
	ClusterNamespace string `json:"cluster_namespace,omitempty"`
	NodeName         string `json:"node_name,omitempty"`
	Message          string `json:"message,omitempty"`
}

// AgentStatus defines the observed state of Agent
type AgentStatus struct {
	Bootstrap bool `json:"bootstrap,omitempty"`
	// +optional
	Role       models.HostRole          `json:"role" protobuf:"bytes,1,opt,name=role,casttype=HostRole,omitempty"`
	Inventory  HostInventory            `json:"inventory,omitempty"`
	Progress   HostProgressInfo         `json:"progress,omitempty"`
	NtpSources []HostNTPSources         `json:"ntpSources,omitempty"`
	Conditions []conditionsv1.Condition `json:"conditions,omitempty"`
	// DebugInfo includes information for debugging the installation process.
	// +optional
	DebugInfo DebugInfo `json:"debugInfo"`

	// ValidationsInfo is a JSON-formatted string containing the validation results for each validation id grouped by category (network, hosts-data, etc.)
	// +optional
	ValidationsInfo common.ValidationsStatus `json:"validationsInfo,omitempty"`

	// InstallationDiskID is the disk that will be used for the installation.
	// +optional
	InstallationDiskID string `json:"installation_disk_id,omitempty"`

	// DeprovisionInfo stores data related to the agent's previous cluster binding in order to clean up when the agent re-registers
	// +optional
	DeprovisionInfo *AgentDeprovisionInfo `json:"deprovision_info,omitempty"`

	// Kind corresponds to the same field in the model Host. It indicates the type of cluster the host is
	// being installed to; either an existing cluster (day-2) or a new cluster (day-1).
	// Value is one of: "AddToExistingClusterHost" (day-2) or "Host" (day-1)
	// +optional
	Kind string `json:"kind,omitempty"`
}

type DebugInfo struct {
	// EventsURL specifies an HTTP/S URL that contains events which occured during the cluster installation process
	// +optional
	EventsURL string `json:"eventsURL,omitempty"`
	// LogsURL specifies a url for download controller logs tar file.
	// +optional
	LogsURL string `json:"logsURL,omitempty"`
	// +optional
	// Current state of the Agent
	State string `json:"state,omitempty"`
	//Additional information pertaining to the status of the Agent
	// +optional
	StateInfo string `json:"stateInfo,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".spec.clusterDeploymentName.name",description="The name of the cluster the Agent registered to."
// +kubebuilder:printcolumn:name="Approved",type="boolean",JSONPath=".spec.approved",description="The `Approve` state of the Agent."
// +kubebuilder:printcolumn:name="Role",type="string",JSONPath=".status.role",description="The role (master/worker) of the Agent."
// +kubebuilder:printcolumn:name="Stage",type="string",JSONPath=".status.progress.currentStage",description="The HostStage of the Agent."
// +kubebuilder:printcolumn:name="Hostname",type="string",JSONPath=".status.inventory.hostname",description="The hostname of the Agent.",priority=1
// +kubebuilder:printcolumn:name="Requested Hostname",type="string",JSONPath=".spec.hostname",description="The requested hostname for the Agent.",priority=1

// Agent is the Schema for the hosts API
type Agent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AgentSpec   `json:"spec,omitempty"`
	Status AgentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AgentList contains a list of Agent
type AgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Agent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Agent{}, &AgentList{})
}
