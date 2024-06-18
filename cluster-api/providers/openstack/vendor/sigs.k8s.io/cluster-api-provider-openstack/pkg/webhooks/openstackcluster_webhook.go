/*
Copyright 2023 The Kubernetes Authors.

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
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-openstackcluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackclusters,versions=v1beta1,name=validation.openstackcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

func SetupOpenStackClusterWebhook(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(&infrav1.OpenStackCluster{}).
		WithValidator(&openStackClusterWebhook{}).
		Complete()
}

type openStackClusterWebhook struct{}

// Compile-time assertion that openStackClusterWebhook implements webhook.CustomValidator.
var _ webhook.CustomValidator = &openStackClusterWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackClusterWebhook) ValidateCreate(_ context.Context, objRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList

	newObj, err := castToOpenStackCluster(objRaw)
	if err != nil {
		return nil, err
	}

	if newObj.Spec.ManagedSecurityGroups != nil {
		for _, rule := range newObj.Spec.ManagedSecurityGroups.AllNodesSecurityGroupRules {
			if rule.RemoteManagedGroups != nil && (rule.RemoteGroupID != nil || rule.RemoteIPPrefix != nil) {
				allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "managedSecurityGroups", "allNodesSecurityGroupRules"), "remoteManagedGroups cannot be used with remoteGroupID or remoteIPPrefix"))
			}
			if rule.RemoteGroupID != nil && (rule.RemoteManagedGroups != nil || rule.RemoteIPPrefix != nil) {
				allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "managedSecurityGroups", "allNodesSecurityGroupRules"), "remoteGroupID cannot be used with remoteManagedGroups or remoteIPPrefix"))
			}
			if rule.RemoteIPPrefix != nil && (rule.RemoteManagedGroups != nil || rule.RemoteGroupID != nil) {
				allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "managedSecurityGroups", "allNodesSecurityGroupRules"), "remoteIPPrefix cannot be used with remoteManagedGroups or remoteGroupID"))
			}
		}
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackClusterWebhook) ValidateUpdate(_ context.Context, oldObjRaw, newObjRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	oldObj, err := castToOpenStackCluster(oldObjRaw)
	if err != nil {
		return nil, err
	}
	newObj, err := castToOpenStackCluster(newObjRaw)
	if err != nil {
		return nil, err
	}

	// Allow changes to Spec.IdentityRef
	oldObj.Spec.IdentityRef = infrav1.OpenStackIdentityReference{}
	newObj.Spec.IdentityRef = infrav1.OpenStackIdentityReference{}

	// Allow changes to ControlPlaneEndpoint fields only if it was not
	// previously set, or did not previously have a host.
	if oldObj.Spec.ControlPlaneEndpoint == nil {
		newObj.Spec.ControlPlaneEndpoint = nil
	} else if oldObj.Spec.ControlPlaneEndpoint.Host == "" {
		oldObj.Spec.ControlPlaneEndpoint = nil
		newObj.Spec.ControlPlaneEndpoint = nil
	}

	// Allow change only for the first time.
	if ptr.Deref(oldObj.Spec.DisableAPIServerFloatingIP, false) && ptr.Deref(oldObj.Spec.APIServerFixedIP, "") == "" {
		oldObj.Spec.APIServerFixedIP = nil
		newObj.Spec.APIServerFixedIP = nil
	}

	// If API Server floating IP is disabled, allow the change of the API Server port only for the first time.
	if ptr.Deref(oldObj.Spec.DisableAPIServerFloatingIP, false) && oldObj.Spec.APIServerPort == nil && newObj.Spec.APIServerPort != nil {
		newObj.Spec.APIServerPort = nil
	}

	// Allow to remove the bastion spec only if it was disabled before.
	if newObj.Spec.Bastion == nil {
		if oldObj.Spec.Bastion.IsEnabled() {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "bastion"), "cannot be removed before disabling it"))
		}
	}

	// Allow changes to the bastion spec.
	oldObj.Spec.Bastion = &infrav1.Bastion{}
	newObj.Spec.Bastion = &infrav1.Bastion{}

	// Allow changes to the managed allNodesSecurityGroupRules.
	if newObj.Spec.ManagedSecurityGroups != nil {
		if oldObj.Spec.ManagedSecurityGroups == nil {
			oldObj.Spec.ManagedSecurityGroups = &infrav1.ManagedSecurityGroups{}
		}

		oldObj.Spec.ManagedSecurityGroups.AllNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}
		newObj.Spec.ManagedSecurityGroups.AllNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}

		// Allow change to the allowAllInClusterTraffic.
		oldObj.Spec.ManagedSecurityGroups.AllowAllInClusterTraffic = false
		newObj.Spec.ManagedSecurityGroups.AllowAllInClusterTraffic = false
	}

	// Allow changes on AllowedCIDRs
	if newObj.Spec.APIServerLoadBalancer != nil && oldObj.Spec.APIServerLoadBalancer != nil {
		oldObj.Spec.APIServerLoadBalancer.AllowedCIDRs = []string{}
		newObj.Spec.APIServerLoadBalancer.AllowedCIDRs = []string{}
	}

	// Allow changes to the availability zones.
	oldObj.Spec.ControlPlaneAvailabilityZones = []string{}
	newObj.Spec.ControlPlaneAvailabilityZones = []string{}

	// Allow the scheduling to be changed from CAPI managed to Nova and
	// vice versa.
	oldObj.Spec.ControlPlaneOmitAvailabilityZone = nil
	newObj.Spec.ControlPlaneOmitAvailabilityZone = nil

	// Allow change on the spec.APIServerFloatingIP only if it matches the current api server loadbalancer IP.
	if oldObj.Status.APIServerLoadBalancer != nil && ptr.Deref(newObj.Spec.APIServerFloatingIP, "") == oldObj.Status.APIServerLoadBalancer.IP {
		newObj.Spec.APIServerFloatingIP = nil
		oldObj.Spec.APIServerFloatingIP = nil
	}

	if !reflect.DeepEqual(oldObj.Spec, newObj.Spec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "cannot be modified"))
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackClusterWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func castToOpenStackCluster(obj runtime.Object) (*infrav1.OpenStackCluster, error) {
	cast, ok := obj.(*infrav1.OpenStackCluster)
	if !ok {
		return nil, fmt.Errorf("expected an OpenStackCluster but got a %T", obj)
	}
	return cast, nil
}
