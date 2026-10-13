/*
Copyright 2023 The Kubernetes Authors.

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

package v1alpha1

type ServerStatusError string

const (
	// OpenstackFloatingIPPoolReadyCondition reports on the current status of the floating ip pool. Ready indicates that the pool is ready to be used.
	OpenstackFloatingIPPoolReadyCondition = "OpenstackFloatingIPPoolReadyCondition"

	// MaxIPsReachedReason is set when the maximum number of floating IPs has been reached.
	MaxIPsReachedReason = "MaxIPsReached"

	// UnableToFindFloatingIPNetworkReason is used when the floating ip network is not found.
	UnableToFindNetwork = "UnableToFindNetwork"

	CreateServerError ServerStatusError = "CreateError"
)
