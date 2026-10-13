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

// Package feature provides a feature-gate implementation for capa.
package feature

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/component-base/featuregate"
)

const (
	// Every capa-specific feature gate should add method here following this template:
	//
	// // owner: @username
	// // alpha: v1.X
	// MyFeature featuregate.Feature = "MyFeature".

	// EKS is used to enable EKS support
	// owner: @richardcase
	// alpha: v0.4
	EKS featuregate.Feature = "EKS"

	// EKSEnableIAM will enable the IAM resource creation/modification
	// owner: @richardcase
	// alpha: v0.4
	EKSEnableIAM featuregate.Feature = "EKSEnableIAM"

	// EKSAllowAddRoles is used to enable the usage of additional IAM roles
	// owner: @richardcase
	// alpha: v0.4
	EKSAllowAddRoles featuregate.Feature = "EKSAllowAddRoles"

	// EKSFargate is used to enable the usage of EKS fargate profiles
	// owner: @richardcase
	// alpha: v0.4
	EKSFargate featuregate.Feature = "EKSFargate"

	// MachinePool is used to enable ASG support
	// owner: @mytunguyen
	// alpha: v0.1
	MachinePool featuregate.Feature = "MachinePool"

	// MachinePoolMachines is a feature gate that enables creation of AWSMachine objects for AWSMachinePool and AWSManagedMachinePool.
	//
	// owner: @AndiDog
	// alpha: v2.8
	MachinePoolMachines featuregate.Feature = "MachinePoolMachines"

	// EventBridgeInstanceState will use Event Bridge and notifications to keep instance state up-to-date
	// owner: @gab-satchi
	// alpha: v0.7?
	EventBridgeInstanceState featuregate.Feature = "EventBridgeInstanceState"

	// AutoControllerIdentityCreator will create AWSClusterControllerIdentity instance that allows all namespaces to use it.
	// owner: @sedefsavas
	// alpha: v0.6
	AutoControllerIdentityCreator featuregate.Feature = "AutoControllerIdentityCreator"

	// BootstrapFormatIgnition will allow an user to enable alternate machine bootstrap format, viz. Ignition.
	BootstrapFormatIgnition featuregate.Feature = "BootstrapFormatIgnition"

	// ExternalResourceGC is used to enable the garbage collection of external resources like NLB/ALB on deletion
	// owner: @richardcase
	// alpha: v1.5
	ExternalResourceGC featuregate.Feature = "ExternalResourceGC"

	// AlternativeGCStrategy is used to enable garbage collection of external resources to be performed without resource group tagging API. It is usually needed in airgap env when tagging API is not available.
	// owner: @wyike
	// alpha: v2.0
	AlternativeGCStrategy featuregate.Feature = "AlternativeGCStrategy"

	// TagUnmanagedNetworkResources is used to disable tagging unmanaged networking resources.
	// owner: @skarlso
	// alpha: v2.0
	TagUnmanagedNetworkResources featuregate.Feature = "TagUnmanagedNetworkResources"

	// ROSA is used to enable ROSA support
	// owner: @enxebre
	// alpha: v2.2
	ROSA featuregate.Feature = "ROSA"
)

func init() {
	runtime.Must(MutableGates.Add(defaultCAPAFeatureGates))
}

// defaultCAPAFeatureGates consists of all known capa-specific feature keys.
// To add a new feature, define a key for it above and add it here.
var defaultCAPAFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	// Every feature should be initiated here:
	EKS:                           {Default: true, PreRelease: featuregate.Beta},
	EKSEnableIAM:                  {Default: false, PreRelease: featuregate.Beta},
	EKSAllowAddRoles:              {Default: false, PreRelease: featuregate.Beta},
	EKSFargate:                    {Default: false, PreRelease: featuregate.Alpha},
	EventBridgeInstanceState:      {Default: false, PreRelease: featuregate.Alpha},
	MachinePool:                   {Default: true, PreRelease: featuregate.Beta},
	MachinePoolMachines:           {Default: false, PreRelease: featuregate.Alpha},
	AutoControllerIdentityCreator: {Default: true, PreRelease: featuregate.Alpha},
	BootstrapFormatIgnition:       {Default: false, PreRelease: featuregate.Alpha},
	ExternalResourceGC:            {Default: true, PreRelease: featuregate.Beta},
	AlternativeGCStrategy:         {Default: false, PreRelease: featuregate.Beta},
	TagUnmanagedNetworkResources:  {Default: true, PreRelease: featuregate.Alpha},
	ROSA:                          {Default: false, PreRelease: featuregate.Alpha},
}
