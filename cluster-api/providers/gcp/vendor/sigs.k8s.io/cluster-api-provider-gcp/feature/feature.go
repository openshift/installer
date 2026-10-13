/*
Copyright 2022 The Kubernetes Authors.
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

// Package feature implements feature functionality.
package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// Every capg-specific feature gate should add method here following this template:
	//
	// // owner: @username
	// // alpha: v1.X
	// MyFeature featuregate.Feature = "MyFeature".

	// GKE is used to enable GKE support
	// owner: @richardchen331 & @richardcase
	// alpha: v0.1
	GKE featuregate.Feature = "GKE"
)

func init() {
	runtime.Must(MutableGates.Add(defaultCAPGFeatureGates))
}

// defaultCAPGFeatureGates consists of all known capg-specific feature keys.
// To add a new feature, define a key for it above and add it here.
var defaultCAPGFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	GKE: {Default: false, PreRelease: featuregate.Alpha},
}
