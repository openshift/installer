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

package v1beta1

const (
	// Represents that the combination of configuration in the MachineSpec
	// is not supported by this cluster. This is not a transient error, but
	// indicates a state that must be fixed before progress can be made.
	//
	// Example: the ProviderSpec specifies an instance type that doesn't exist,
	InvalidConfigurationMachineError MachineStatusError = "InvalidConfiguration"

	// There was an error while trying to create a Node to match this
	// Machine. This may indicate a transient problem that will be fixed
	// automatically with time, such as a service outage, or a terminal
	// error during creation that doesn't match a more specific
	// MachineStatusError value.
	//
	// Example: timeout trying to connect to GCE.
	CreateMachineError MachineStatusError = "CreateError"

	// There was an error while trying to update a Node that this
	// Machine represents. This may indicate a transient problem that will be
	// fixed automatically with time, such as a service outage,
	//
	// Example: error updating load balancers
	UpdateMachineError MachineStatusError = "UpdateError"

	// An error was encountered while trying to delete the Node that this
	// Machine represents. This could be a transient or terminal error, but
	// will only be observable if the provider's Machine controller has
	// added a finalizer to the object to more gracefully handle deletions.
	//
	// Example: cannot resolve EC2 IP address.
	DeleteMachineError MachineStatusError = "DeleteError"

	// TemplateClonedFromGroupKindAnnotation is the infrastructure machine annotation that stores the group-kind of the infrastructure template resource
	// that was cloned for the machine. This annotation is set only during cloning a template. Older/adopted machines will not have this annotation.
	TemplateClonedFromGroupKindAnnotation = "machine.openshift.io/cloned-from-groupkind"

	// TemplateClonedFromNameAnnotation is the infrastructure machine annotation that stores the name of the infrastructure template resource
	// that was cloned for the machine. This annotation is set only during cloning a template. Older/adopted machines will not have this annotation.
	TemplateClonedFromNameAnnotation = "machine.openshift.io/cloned-from-name"
)

// Constants aren't automatically generated for unversioned packages.
// Instead share the same constant for all versioned packages
type MachineStatusError string

type MachineSetStatusError string
