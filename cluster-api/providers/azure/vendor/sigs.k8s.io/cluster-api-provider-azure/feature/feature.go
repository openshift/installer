/*
Copyright 2020 The Kubernetes Authors.

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

package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	//nolint:godot // Ignore "Comment should end in a period" check.
	// Every capz-specific feature gate should add a method here following this template:
	//
	// // MyFeature is the feature gate for my feature.
	// // owner: @username
	// // alpha: v1.X
	// MyFeature featuregate.Feature = "MyFeature"

	//nolint:godot // Ignore "Comment should end in a period" check.
	// AKS is the feature gate for AKS managed clusters.
	// owner: @jackfrancis @nojnhuh
	// alpha: v0.4
	// GA: v1.8
	AKS featuregate.Feature = "AKS"

	// AKSResourceHealth is the feature gate for reporting Azure Resource Health
	// on AKS managed clusters.
	// owner: @nojnhuh
	// alpha: v1.7
	AKSResourceHealth featuregate.Feature = "AKSResourceHealth"

	// EdgeZone is the feature gate for creating clusters on public MEC.
	// owner: @upxinxin
	// alpha: v1.8
	EdgeZone featuregate.Feature = "EdgeZone"

	// ASOAPI is the feature gate for enabling the AzureASO... APIs.
	// owner: @nojnhuh
	// alpha: v1.15
	ASOAPI featuregate.Feature = "ASOAPI"
)

func init() {
	runtime.Must(MutableGates.Add(defaultCAPZFeatureGates))
}

// defaultCAPZFeatureGates consists of all known capz-specific feature keys.
// To add a new feature, define a key for it above and add it here.
var defaultCAPZFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	// Every feature should be initiated here:
	AKS:               {Default: true, PreRelease: featuregate.GA, LockToDefault: true}, // Remove in 1.12
	AKSResourceHealth: {Default: false, PreRelease: featuregate.Alpha},
	EdgeZone:          {Default: false, PreRelease: featuregate.Alpha},
	ASOAPI:            {Default: true, PreRelease: featuregate.Alpha},
}
