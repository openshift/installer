/*
Copyright 2024 The Kubernetes Authors.

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
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

func RegisterAllWithManager(mgr manager.Manager) []error {
	var errs []error

	// Register webhooks for all types with custom validators.
	for _, webhook := range []struct {
		name  string
		setup func(ctrl.Manager) error
	}{
		{"OpenStackCluster", SetupOpenStackClusterWebhook},
		{"OpenStackClusterTemplate", SetupOpenStackClusterTemplateWebhook},
		{"OpenStackMachine", SetupOpenStackMachineWebhook},
		{"OpenStackMachineTemplate", SetupOpenStackMachineTemplateWebhook},
		{"OpenStackServer", SetupOpenStackServerWebhook},
	} {
		if err := webhook.setup(mgr); err != nil {
			errs = append(errs, fmt.Errorf("creating webhook for %s: %v", webhook.name, err))
		}
	}

	// Additionally register webhooks for other types so they get conversion webhooks.
	for _, conversionOnlyType := range []conversion.Hub{
		&infrav1.OpenStackClusterList{},
		&infrav1.OpenStackClusterTemplateList{},
		&infrav1.OpenStackMachineList{},
		&infrav1.OpenStackMachineTemplateList{},
	} {
		if err := builder.WebhookManagedBy(mgr).
			For(conversionOnlyType).
			Complete(); err != nil {
			errs = append(errs, fmt.Errorf("creating webhook for %T: %v", conversionOnlyType, err))
		}
	}

	return errs
}
