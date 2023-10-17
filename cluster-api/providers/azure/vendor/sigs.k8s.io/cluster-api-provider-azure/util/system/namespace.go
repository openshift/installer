/*
Copyright 2021 The Kubernetes Authors.

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

package system

import "os"

const (
	// NamespaceEnvVarName is the env var coming from DownwardAPI in the manager manifest.
	NamespaceEnvVarName = "POD_NAMESPACE"
	// DefaultNamespace is the default value from manifest.
	DefaultNamespace = "capz-system"
)

// GetManagerNamespace returns the namespace where the controller is running.
func GetManagerNamespace() string {
	managerNamespace := os.Getenv(NamespaceEnvVarName)
	if managerNamespace == "" {
		managerNamespace = DefaultNamespace
	}
	return managerNamespace
}
