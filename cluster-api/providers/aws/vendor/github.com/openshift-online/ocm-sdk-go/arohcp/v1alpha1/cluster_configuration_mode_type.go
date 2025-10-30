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

// ClusterConfigurationMode represents the values of the 'cluster_configuration_mode' enumerated type.
type ClusterConfigurationMode string

const (
	// Full configuration (default).
	ClusterConfigurationModeFull ClusterConfigurationMode = "full"
	// Only read configuration operations are supported.
	// The cluster can't be deleted, reshaped, configure IDPs, add/remove users, etc.
	ClusterConfigurationModeReadOnly ClusterConfigurationMode = "read_only"
)
