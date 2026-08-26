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

// Package feature handles feature gates.
package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// Every capo-specific feature gate should add method here following this template:
	//
	// // owner: @username
	// // alpha: v1.X
	// MyFeature featuregate.Feature = "MyFeature".

	// PriorityQueue is a feature gate that controls if the controller uses the controller-runtime PriorityQueue
	// instead of the default queue implementation.
	//
	// alpha: v0.14
	PriorityQueue featuregate.Feature = "PriorityQueue"

	// AutoScaleFromZero is a feature gate that enables the OpenStackMachineTemplate controller that adds
	// information in OpenStackMachineTemplate.status required by the cluster-autoscaler to scale from zero
	// without the addition of labels
	//
	// alpha: v0.14
	AutoScaleFromZero featuregate.Feature = "AutoScaleFromZero"
)

func init() {
	runtime.Must(MutableGates.Add(defaultCAPOFeatureGates))
}

// defaultCAPOFeatureGates consists of all known capo-specific feature keys.
// To add a new feature, define a key for it above and add it here.
var defaultCAPOFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	// Every feature should be initiated here:
	PriorityQueue:     {Default: false, PreRelease: featuregate.Alpha},
	AutoScaleFromZero: {Default: false, PreRelease: featuregate.Alpha},
}
