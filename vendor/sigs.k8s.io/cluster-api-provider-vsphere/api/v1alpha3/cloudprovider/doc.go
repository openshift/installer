/*
Copyright 2019 The Kubernetes Authors.

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

// Package cloudprovider contains API types for the vSphere cloud provider.
//
// The configuration may be unmarshalled from an INI-style configuration using
// the "gopkg.in/gcfg.v1" package.
//
// The configuration may be marshalled to an INI-style configuration using a Go
// template.
//
// The "gopkg.in/go-ini/ini.v1" package was investigated, but it does not
// support reflecting a struct with a field of type "map[string]TYPE" to INI.
//
// +kubebuilder:object:generate=true
package cloudprovider
