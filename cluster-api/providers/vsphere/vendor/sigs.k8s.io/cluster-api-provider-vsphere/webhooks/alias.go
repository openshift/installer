/*
Copyright 2025 The Kubernetes Authors.

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

package webhooks

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks"
)

// VSphereClusterTemplate implements a validation webhook for VSphereClusterTemplate.
type VSphereClusterTemplate struct{}

// SetupWebhookWithManager sets up VSphereClusterTemplate webhooks.
func (webhook *VSphereClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&webhooks.VSphereClusterTemplateWebhook{}).SetupWebhookWithManager(mgr)
}

// VSphereDeploymentZone implements a defaulting webhook for VSphereDeploymentZone.
type VSphereDeploymentZone struct{}

// SetupWebhookWithManager sets up VSphereDeploymentZone webhooks.
func (webhook *VSphereDeploymentZone) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&webhooks.VSphereDeploymentZoneWebhook{}).SetupWebhookWithManager(mgr)
}

// VSphereFailureDomain implements a validation and defaulting webhook for VSphereFailureDomain.
type VSphereFailureDomain struct{}

// SetupWebhookWithManager sets up VSphereFailureDomain webhooks.
func (webhook *VSphereFailureDomain) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&webhooks.VSphereFailureDomainWebhook{}).SetupWebhookWithManager(mgr)
}

// VSphereMachine implements a validation and defaulting webhook for VSphereMachine.
type VSphereMachine struct{}

// SetupWebhookWithManager sets up VSphereMachine webhooks.
func (webhook *VSphereMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&webhooks.VSphereMachineWebhook{}).SetupWebhookWithManager(mgr)
}

// VSphereMachineTemplate implements a validation webhook for VSphereMachineTemplate.
type VSphereMachineTemplate struct{}

// SetupWebhookWithManager sets up VSphereMachineTemplate webhooks.
func (webhook *VSphereMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&webhooks.VSphereMachineTemplateWebhook{}).SetupWebhookWithManager(mgr)
}

// VSphereVM implements a validation and defaulting webhook for VSphereVM.
type VSphereVM struct{}

// SetupWebhookWithManager sets up VSphereVM webhooks.
func (webhook *VSphereVM) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&webhooks.VSphereVMWebhook{}).SetupWebhookWithManager(mgr)
}
