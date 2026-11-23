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

// Package ignition provides Go type definitions for the Coreos Ignition v2
// machine bootstrap configuration format.
//
// Why this is a local copy:
// The upstream library github.com/coreos/ignition/v2 only ships v3.x config
// types (v3_0 through v3_6); v2.x types were removed from the library entirely.
// There is no upstream Go package for Ignition v2 types, so this local
// definition is the only option for generating Ignition v2 redirect documents
// when creating PowerVS machine user-data.
package ignition
