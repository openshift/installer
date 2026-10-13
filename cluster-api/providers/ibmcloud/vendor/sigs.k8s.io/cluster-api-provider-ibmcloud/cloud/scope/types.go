/*
Copyright 2024 The Kubernetes Authors.

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

// ResourceNotFound is the string representing an error when a resource is not found in IBM Cloud.
type ResourceNotFound string

var (
	// ResourceNotFoundCode indicates the http status code when a resource does not exist.
	ResourceNotFoundCode = 404

	// DHCPServerNotFound is the error returned when a DHCP server is not found.
	DHCPServerNotFound = ResourceNotFound("dhcp server does not exist")
)
