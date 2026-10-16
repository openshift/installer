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

// Package hub defines the hub version of vm-operator core types.
//
// The hub version of vm-operator core types act as a version used internally by CAPV; more specifically:
// - When CAPV reads and write vm-operator resources, it will use vm-operator's preferred version for the specific environment.
// - vm-operator's preferred version depends primarily on the version of vm-operator that exists in the environment.
// - Apart from the read and write operations, CAPV will use its own hub types for processing.
//
// Notably:
// - CAPV hub types are a subset of vm-operator's types, with only the fields that are relevant for CAPV.
//
// +kubebuilder:object:generate=true
package hub
