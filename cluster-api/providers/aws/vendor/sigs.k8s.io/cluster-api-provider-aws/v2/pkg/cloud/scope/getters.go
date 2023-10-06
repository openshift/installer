/*
Copyright 2018 The Kubernetes Authors.

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

package scope

var (
	// DefaultClusterScopeGetter defines the default cluster scope getter.
	DefaultClusterScopeGetter ClusterScopeGetter = ClusterScopeGetterFunc(NewClusterScope)

	// DefaultMachineScopeGetter defines the default machine scope getter.
	DefaultMachineScopeGetter MachineScopeGetter = MachineScopeGetterFunc(NewMachineScope)
)

// ClusterScopeGetter defines the cluster scope getter interface.
type ClusterScopeGetter interface {
	ClusterScope(params ClusterScopeParams) (*ClusterScope, error)
}

// ClusterScopeGetterFunc defines handler types for cluster scope getters.
type ClusterScopeGetterFunc func(params ClusterScopeParams) (*ClusterScope, error)

// ClusterScope will return the cluster scope.
func (f ClusterScopeGetterFunc) ClusterScope(params ClusterScopeParams) (*ClusterScope, error) {
	return f(params)
}

// MachineScopeGetter defines the machine scope getter interface.
type MachineScopeGetter interface {
	MachineScope(params MachineScopeParams) (*MachineScope, error)
}

// MachineScopeGetterFunc defines handler types for machine scope getters.
type MachineScopeGetterFunc func(params MachineScopeParams) (*MachineScope, error)

// MachineScope will return the machine scope.
func (f MachineScopeGetterFunc) MachineScope(params MachineScopeParams) (*MachineScope, error) {
	return f(params)
}
