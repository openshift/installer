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

package vmware

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/vmware"
)

// VSphereMachine implements a validation and defaulting webhook for VSphereMachine.
type VSphereMachine struct{}

// SetupWebhookWithManager sets up VSphereMachine webhooks.
func (webhook *VSphereMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&vmware.VSphereMachine{}).SetupWebhookWithManager(mgr)
}

// VSphereMachineTemplate implements a validation webhook for VSphereMachineTemplate.
type VSphereMachineTemplate struct{}

// SetupWebhookWithManager sets up VSphereMachineTemplate webhooks.
func (webhook *VSphereMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return (&vmware.VSphereMachineTemplate{}).SetupWebhookWithManager(mgr)
}
