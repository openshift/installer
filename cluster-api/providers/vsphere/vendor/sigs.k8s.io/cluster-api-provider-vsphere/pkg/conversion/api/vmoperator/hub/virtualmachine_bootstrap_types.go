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

// VirtualMachineBootstrapSpec defines the desired state of a VM's bootstrap
// configuration.
type VirtualMachineBootstrapSpec struct {

	// +optional

	// CloudInit may be used to bootstrap Linux guests with Cloud-Init or
	// Windows guests that support Cloudbase-Init.
	//
	// The guest's networking stack is configured by Cloud-Init on Linux guests
	// and Cloudbase-Init on Windows guests.
	//
	// Please note this bootstrap provider may not be used in conjunction with
	// the other bootstrap providers.
	CloudInit *VirtualMachineBootstrapCloudInitSpec `json:"cloudInit,omitempty"`
}

// VirtualMachineBootstrapCloudInitSpec describes the CloudInit configuration
// used to bootstrap the VM.
type VirtualMachineBootstrapCloudInitSpec struct {

	// +optional

	// RawCloudConfig describes a key in a Secret resource that contains the
	// CloudConfig data used to bootstrap the VM.
	//
	// The CloudConfig data specified by the key may be plain-text,
	// base64-encoded, or gzipped and base64-encoded.
	//
	// Please note this field and CloudConfig are mutually exclusive.
	RawCloudConfig *SecretKeySelector `json:"rawCloudConfig,omitempty"`
}
