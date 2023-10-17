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

package v1beta1

import (
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// InstanceProvisionFailedReason used for failures during instance provisioning.
	InstanceProvisionFailedReason = "InstanceProvisionFailed"
	// WaitingForClusterInfrastructureReason used when machine is waiting for cluster infrastructure to be ready before proceeding.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"
	// WaitingForBootstrapDataReason used when machine is waiting for bootstrap data to be ready before proceeding.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"
)

const (
	// InstanceStoppedReason instance is in a stopped state.
	InstanceStoppedReason = "InstanceStopped"

	// InstanceErroredReason instance is in a errored state.
	InstanceErroredReason = "InstanceErrored"

	// InstanceNotReadyReason used when the instance is in a pending state.
	InstanceNotReadyReason = "InstanceNotReady"
)

const (
	// InstanceReadyCondition reports on current status of the instance. Ready indicates the instance is in a Running state.
	InstanceReadyCondition capiv1beta1.ConditionType = "InstanceReady"
)

const (
	// WaitingForIBMPowerVSImageReason used when machine is waiting for powervs image to be ready before proceeding.
	WaitingForIBMPowerVSImageReason = "WaitingForIBMPowerVSImage"
)

const (
	// ImageNotReadyReason used when the image is in a queued state.
	ImageNotReadyReason = "ImageNotReady"

	// ImageImportFailedReason used when the image import is failed.
	ImageImportFailedReason = "ImageImportFailed"
)

const (
	// ImageReadyCondition reports on current status of the image. Ready indicates the image is in a active state.
	ImageReadyCondition capiv1beta1.ConditionType = "ImageReady"

	// ImageImportedCondition reports on current status of the image import job. Ready indicates the import job is finished.
	ImageImportedCondition capiv1beta1.ConditionType = "ImageImported"
)

const (
	// LoadBalancerNotReadyReason used when cluster is waiting for load balancer to be ready before proceeding.
	LoadBalancerNotReadyReason = "LoadBalancerNotReady"
)

const (
	// LoadBalancerReadyCondition reports on current status of the load balancer. Ready indicates the load balancer is in a active state.
	LoadBalancerReadyCondition capiv1beta1.ConditionType = "LoadBalancerReady"
)
