/*
Copyright 2021 The Kubernetes Authors.

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
	vmoprv1alpha5 "github.com/vmware-tanzu/vm-operator/api/v1alpha5"
	vmoprv1alpha6 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
	"k8s.io/component-base/featuregate"
)

const (
	// Every capv-specific feature gate should add method here following this template:
	//
	// // owner: @username
	// // alpha: v1.X
	// MyFeature featuregate.Feature = "MyFeature".

	// MultiNetworks is a feature gate for the MultiNetworks functionality for supervisor.
	//
	// alpha: v1.14
	MultiNetworks featuregate.Feature = "MultiNetworks"

	// NodeAntiAffinity is a feature gate for the NodeAntiAffinity functionality.
	//
	// alpha: v1.4
	NodeAntiAffinity featuregate.Feature = "NodeAntiAffinity"

	// NamespaceScopedZones is a feature gate for the NamespaceScopedZones functionality for supervisor.
	//
	// alpha: v1.11
	NamespaceScopedZones featuregate.Feature = "NamespaceScopedZones"

	// NodeAutoPlacement is a feature gate for the NodeAutoPlacement functionality for supervisor.
	//
	// alpha: v1.15
	NodeAutoPlacement featuregate.Feature = "NodeAutoPlacement"

	// InfrastructurePolicies is a feature gate for the Support for Supervisor infrastructure policies.
	// When enabled, VSphereMachine.spec.policies are validated on admission and
	// mapped to the underlying VirtualMachine spec.
	//
	// alpha: v1.17
	InfrastructurePolicies featuregate.Feature = "InfrastructurePolicies"

	// PriorityQueue is a feature gate that controls if the controller uses the controller-runtime PriorityQueue
	// instead of the default queue implementation.
	//
	// alpha: v1.13
	// beta: v1.16
	PriorityQueue featuregate.Feature = "PriorityQueue"

	// ReconcilerRateLimiting is a feature gate that controls if reconcilers are rate-limited.
	// Note: Currently the feature gate is rate-limiting to 1 request / 1 second.
	//
	// beta: v1.16
	ReconcilerRateLimiting featuregate.Feature = "ReconcilerRateLimiting"

	// IPv6DualStack enables IPv6 and dualstack support for clusters with NSX-VPC network provider.
	// Requires vm-operator v1alpha6 or later on the supervisor.
	//
	// alpha: v1.17
	IPv6DualStack featuregate.Feature = "IPv6DualStack"

	// VLANSubinterface is a feature gate for the VLAN sub-interface functionality for supervisor.
	//
	// alpha: v1.17
	VLANSubinterface featuregate.Feature = "VLANSubinterface"
)

var (
	commonGates = map[featuregate.Feature]featuregate.FeatureSpec{
		PriorityQueue:          {Default: true, PreRelease: featuregate.Beta},
		ReconcilerRateLimiting: {Default: true, PreRelease: featuregate.Beta},
	}

	govmomiGates = map[featuregate.Feature]featuregate.FeatureSpec{
		NodeAntiAffinity: {Default: false, PreRelease: featuregate.Alpha},
	}

	supervisorGates = map[featuregate.Feature]featuregate.FeatureSpec{
		NamespaceScopedZones: {Default: false, PreRelease: featuregate.Alpha},
		NodeAutoPlacement:    {Default: false, PreRelease: featuregate.Alpha},
		MultiNetworks:        {Default: false, PreRelease: featuregate.Alpha},
	}

	supervisorVersionedGates = map[featuregate.Feature]featuregate.VersionedSpecs{
		// NOTE: Supervisor features gates depending on a specific vm-operator version should be added here.
		// e.g.
		// FeatureDependingOnV1alpha5: {
		//  {Version: toFeatureVersion(vmoprv1alpha5.GroupVersion.Version), Default: false, PreRelease: featuregate.Alpha},
		// },
		// FeatureDependingOnV1alpha6: {
		// 	{Version: toFeatureVersion(vmoprv1alpha6.GroupVersion.Version), Default: false, PreRelease: featuregate.Alpha},
		// },
		InfrastructurePolicies: {
			{Version: toFeatureVersion(vmoprv1alpha5.GroupVersion.Version), Default: false, PreRelease: featuregate.Alpha},
		},
		IPv6DualStack: {
			{Version: toFeatureVersion(vmoprv1alpha6.GroupVersion.Version), Default: false, PreRelease: featuregate.Alpha},
		},
		VLANSubinterface: {
			{Version: toFeatureVersion(vmoprv1alpha6.GroupVersion.Version), Default: false, PreRelease: featuregate.Alpha},
		},
	}
)
