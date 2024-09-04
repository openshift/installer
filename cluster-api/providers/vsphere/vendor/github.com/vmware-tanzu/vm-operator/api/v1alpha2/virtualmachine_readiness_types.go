// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	"k8s.io/apimachinery/pkg/util/intstr"
)

// VirtualMachineReadinessProbeSpec describes a probe used to determine if a VM
// is in a ready state. All probe actions are mutually exclusive.
type VirtualMachineReadinessProbeSpec struct {
	// TCPSocket specifies an action involving a TCP port.
	//
	// Deprecated: The TCPSocket action requires network connectivity that is not supported in all environments.
	// This field will be removed in a later API version.
	// +optional
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`

	// GuestHeartbeat specifies an action involving the guest heartbeat status.
	// +optional
	GuestHeartbeat *GuestHeartbeatAction `json:"guestHeartbeat,omitempty"`

	// GuestInfo specifies an action involving key/value pairs from GuestInfo.
	//
	// The elements are evaluated with the logical AND operator, meaning
	// all expressions must evaluate as true for the probe to succeed.
	//
	// For example, a VM resource's probe definition could be specified as the
	// following:
	//
	//         guestInfo:
	//         - key:   ready
	//           value: true
	//
	// With the above configuration in place, the VM would not be considered
	// ready until the GuestInfo key "ready" was set to the value "true".
	//
	// From within the guest operating system it is possible to set GuestInfo
	// key/value pairs using the program "vmware-rpctool," which is included
	// with VM Tools. For example, the following command will set the key
	// "guestinfo.ready" to the value "true":
	//
	//         vmware-rpctool "info-set guestinfo.ready true"
	//
	// Once executed, the VM's readiness probe will be signaled and the
	// VM resource will be marked as ready.
	//
	// +optional
	GuestInfo []GuestInfoAction `json:"guestInfo,omitempty"`

	// TimeoutSeconds specifies a number of seconds after which the probe times out.
	// Defaults to 10 seconds. Minimum value is 1.
	// +optional
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=60
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`

	// PeriodSeconds specifics how often (in seconds) to perform the probe.
	// Defaults to 10 seconds. Minimum value is 1.
	// +optional
	// +kubebuilder:validation:Minimum:=1
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
}

// TCPSocketAction describes an action based on opening a socket.
type TCPSocketAction struct {
	// Port specifies a number or name of the port to access on the VM.
	// If the format of port is a number, it must be in the range 1 to 65535.
	// If the format of name is a string, it must be an IANA_SVC_NAME.
	Port intstr.IntOrString `json:"port"`

	// Host is an optional host name to connect to. Host defaults to the VM IP.
	// +optional
	Host string `json:"host,omitempty"`
}

// GuestHeartbeatStatus is the guest heartbeat status.
type GuestHeartbeatStatus string

// See govmomi.vim25.types.ManagedEntityStatus
const (
	// GrayHeartbeatStatus means VMware Tools are not installed or not running.
	GrayHeartbeatStatus GuestHeartbeatStatus = "gray"
	// RedHeartbeatStatus means no heartbeat.
	// Guest operating system may have stopped responding.
	RedHeartbeatStatus GuestHeartbeatStatus = "red"
	// YellowHeartbeatStatus means an intermittent heartbeat.
	// This may be due to guest load.
	YellowHeartbeatStatus GuestHeartbeatStatus = "yellow"
	// GreenHeartbeatStatus means the guest operating system is responding normally.
	GreenHeartbeatStatus GuestHeartbeatStatus = "green"
)

// GuestHeartbeatAction describes an action based on the guest heartbeat.
type GuestHeartbeatAction struct {
	// ThresholdStatus is the value that the guest heartbeat status must be at or above to be
	// considered successful.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=green
	// +kubebuilder:validation:Enum=yellow;green
	ThresholdStatus GuestHeartbeatStatus `json:"thresholdStatus,omitempty"`
}

// GuestInfoAction describes a key from GuestInfo that must match the associated
// value expression.
type GuestInfoAction struct {
	// Key is the name of the GuestInfo key.
	//
	// The key is automatically prefixed with "guestinfo." before being
	// evaluated. Thus if the key "guestinfo.mykey" is provided, it will be
	// evaluated as "guestinfo.guestinfo.mykey".
	Key string `json:"key"`

	// Value is a regular expression that is matched against the value of the
	// specified key.
	//
	// An empty value is the equivalent of "match any" or ".*".
	//
	// All values must adhere to the RE2 regular expression syntax as documented
	// at https://golang.org/s/re2syntax. Invalid values may be rejected or
	// ignored depending on the implementation of this API. Either way, invalid
	// values will not be considered when evaluating the ready state of a VM.
	//
	// +optional
	Value string `json:"value,omitempty"`
}
