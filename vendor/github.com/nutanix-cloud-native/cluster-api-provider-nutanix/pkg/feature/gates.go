/*
Copyright 2026 Nutanix

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

//nolint:gochecknoinits // Standard Kubernetes pattern for feature gate registration.
func init() {
	runtime.Must(MutableGates.Add(defaultFeatureGates()))
}

var (
	// MutableGates is a mutable version of Gates.
	// Only top-level commands/options setup should make use of this.
	// Tests that need to modify feature gates for the duration of their test should use:
	//   featuregatetesting "k8s.io/component-base/featuregate/testing"
	//   featuregatetesting.SetFeatureGateDuringTest(t, feature.Gates, feature.<FeatureName>, <value>)
	MutableGates featuregate.MutableFeatureGate = featuregate.NewFeatureGate()

	// Gates is a shared global FeatureGate.
	// Read feature gate state using: feature.Gates.Enabled(feature.<FeatureName>)
	Gates featuregate.FeatureGate = MutableGates
)
