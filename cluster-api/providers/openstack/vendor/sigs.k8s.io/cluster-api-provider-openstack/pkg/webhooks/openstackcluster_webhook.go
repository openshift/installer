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

// allowSubnetFilterToIDTransition checks if changes to OpenStackCluster.Spec.Subnets
// are transitioning from a Filter-based definition to an ID-based one, and whether
// those transitions are valid based on the current status.network.subnets.
//
// This function only allows Filter → ID transitions when the filter name in the old
// spec matches the subnet name in status, and the new ID matches the corresponding subnet ID.
//
// Returns true if all such transitions are valid; false otherwise.
func allowSubnetFilterToIDTransition(oldObj, newObj *infrav1.OpenStackCluster) bool {
	if newObj.Spec.Network == nil || oldObj.Spec.Network == nil || oldObj.Status.Network == nil {
		return false
	}

	if len(newObj.Spec.Subnets) != len(oldObj.Spec.Subnets) || len(oldObj.Status.Network.Subnets) == 0 {
		return false
	}

	for i := range newObj.Spec.Subnets {
		oldSubnet := oldObj.Spec.Subnets[i]
		newSubnet := newObj.Spec.Subnets[i]

		// Allow Filter → ID only if both values match a known subnet in status
		if oldSubnet.Filter != nil && newSubnet.ID != nil && newSubnet.Filter == nil {
			matchFound := false
			for _, statusSubnet := range oldObj.Status.Network.Subnets {
				if oldSubnet.Filter.Name == statusSubnet.Name && *newSubnet.ID == statusSubnet.ID {
					matchFound = true
					break
				}
			}
			if !matchFound {
				return false
			}
		}

		// Reject any change from ID → Filter
		if oldSubnet.ID != nil && newSubnet.Filter != nil {
			return false
		}

		// Reject changes to Filter or ID if they do not match the old values
		if oldSubnet.Filter != nil && newSubnet.Filter != nil &&
			oldSubnet.Filter.Name != newSubnet.Filter.Name {
			return false
		}
		if oldSubnet.ID != nil && newSubnet.ID != nil &&
			*oldSubnet.ID != *newSubnet.ID {
			return false
		}
	}

	return true
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

	// Allow changes to the managed securityGroupRules.
	if newObj.Spec.ManagedSecurityGroups != nil {
		if oldObj.Spec.ManagedSecurityGroups == nil {
			oldObj.Spec.ManagedSecurityGroups = &infrav1.ManagedSecurityGroups{}
		}

		oldObj.Spec.ManagedSecurityGroups.AllNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}
		newObj.Spec.ManagedSecurityGroups.AllNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}

		oldObj.Spec.ManagedSecurityGroups.ControlPlaneNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}
		newObj.Spec.ManagedSecurityGroups.ControlPlaneNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}

		oldObj.Spec.ManagedSecurityGroups.WorkerNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}
		newObj.Spec.ManagedSecurityGroups.WorkerNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}

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

	// Allow changes from filter to id for spec.network and spec.subnets
	if newObj.Spec.Network != nil && oldObj.Spec.Network != nil && oldObj.Status.Network != nil {
		// Allow change from spec.network.subnets from filter to id if it matches the current subnets.
		if allowSubnetFilterToIDTransition(oldObj, newObj) {
			oldObj.Spec.Subnets = nil
			newObj.Spec.Subnets = nil
		}
		// Allow change from spec.network.filter to spec.network.id only if it matches the current network.
		if ptr.Deref(newObj.Spec.Network.ID, "") == oldObj.Status.Network.ID {
			newObj.Spec.Network = nil
			oldObj.Spec.Network = nil
		}
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
