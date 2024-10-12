// Copyright 2024 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

type Resource interface {
	Describe() ServiceTypeVersion
}

// ServiceTypeVersion is a tuple that can uniquely identify a
// DCL resource type.
type ServiceTypeVersion struct {
	// Service indicates the service to which this resource
	// belongs, e.g., "compute". It is roughly analogous to the
	// K8S "Group" identifier.
	Service string

	// Type identifies the particular type of this resource,
	// e.g., "ComputeInstance". It maps to the K8S "Kind".
	Type string

	// Version is the DCL version of the resource, e.g.,
	// "beta" or "ga".
	Version string
}
