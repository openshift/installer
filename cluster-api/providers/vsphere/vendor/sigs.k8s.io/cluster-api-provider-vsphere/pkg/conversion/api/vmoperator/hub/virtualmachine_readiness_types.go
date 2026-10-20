/*
Copyright 2026 The Kubernetes Authors.

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

package hub

import "k8s.io/apimachinery/pkg/util/intstr"

// VirtualMachineReadinessProbeSpec describes a probe used to determine if a VM
// is in a ready state. All probe actions are mutually exclusive.
type VirtualMachineReadinessProbeSpec struct {
	// +optional

	// TCPSocket specifies an action involving a TCP port.
	//
	// Deprecated: The TCPSocket action requires network connectivity that is not supported in all environments.
	// This field will be removed in a later API version.
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty"`
}

// TCPSocketAction describes an action based on opening a socket.
type TCPSocketAction struct {
	// Port specifies a number or name of the port to access on the VM.
	// If the format of port is a number, it must be in the range 1 to 65535.
	// If the format of name is a string, it must be an IANA_SVC_NAME.
	Port intstr.IntOrString `json:"port"`

	// +optional

	// Host is an optional host name to connect to. Host defaults to the VM IP.
	Host string `json:"host,omitempty"`
}
