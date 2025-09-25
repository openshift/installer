/*
Copyright (c) 2020 Red Hat, Inc.

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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// MachineTypeSize represents the values of the 'machine_type_size' enumerated type.
type MachineTypeSize string

const (
	// Large machine type (e.g. c5.4xlarge, custom-16-65536)
	MachineTypeSizeLarge MachineTypeSize = "large"
	// Medium machine type (e.g. r5.2xlarge, custom-8-32768)
	MachineTypeSizeMedium MachineTypeSize = "medium"
	// Small machine type (e.g. m5.xlarge, custom-4-16384)
	MachineTypeSizeSmall MachineTypeSize = "small"
)
