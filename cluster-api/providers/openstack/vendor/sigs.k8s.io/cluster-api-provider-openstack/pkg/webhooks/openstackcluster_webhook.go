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

	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-openstackcluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackclusters,versions=v1beta2,name=validation.openstackcluster.v1beta2.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// SetupOpenStackClusterWebhook registers the validating webhook for OpenStackCluster with the manager.
func SetupOpenStackClusterWebhook(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr, &infrav1.OpenStackCluster{}).
		WithValidator(&openStackClusterWebhook{}).
		Complete()
}

type openStackClusterWebhook struct{}

var _ admission.Validator[*infrav1.OpenStackCluster] = &openStackClusterWebhook{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type.
func (*openStackClusterWebhook) ValidateCreate(_ context.Context, newObj *infrav1.OpenStackCluster) (admission.Warnings, error) {
	var allErrs field.ErrorList

	if newObj.Spec.ManagedSecurityGroups != nil {
		allErrs = append(allErrs, validateManagedSecurityGroupRules(newObj.Spec.ManagedSecurityGroups)...)
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

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type.
func (*openStackClusterWebhook) ValidateUpdate(_ context.Context, oldObj, newObj *infrav1.OpenStackCluster) (admission.Warnings, error) { //nolint:gocyclo,cyclop
	var allErrs field.ErrorList

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

	// Allow APIServerFixedIP to be set for the first time when floating IP is disabled.
	if oldObj.Spec.APIServer != nil && newObj.Spec.APIServer != nil {
		if !ptr.Deref(oldObj.Spec.APIServer.GetEnableFloatingIP(), true) && ptr.Deref(oldObj.Spec.APIServer.FixedIP, "") == "" {
			oldObj.Spec.APIServer.FixedIP = nil
			newObj.Spec.APIServer.FixedIP = nil
		}
		// If API Server floating IP is disabled, allow the change of the API Server port only for the first time.
		if !ptr.Deref(oldObj.Spec.APIServer.GetEnableFloatingIP(), true) && oldObj.Spec.APIServer.GetPort() == nil && newObj.Spec.APIServer.GetPort() != nil {
			newObj.Spec.APIServer.Port = nil
		}
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

	// Validate new security group rules before zeroing them out for immutability comparison.
	if newObj.Spec.ManagedSecurityGroups != nil {
		allErrs = append(allErrs, validateManagedSecurityGroupRules(newObj.Spec.ManagedSecurityGroups)...)
	}

	// Allow changes to the managed securityGroupRules.
	if newObj.Spec.ManagedSecurityGroups != nil {
		if oldObj.Spec.ManagedSecurityGroups == nil {
			oldObj.Spec.ManagedSecurityGroups = &infrav1.ManagedSecurityGroups{}
		}

		oldObj.Spec.ManagedSecurityGroups.ClusterNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}
		newObj.Spec.ManagedSecurityGroups.ClusterNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}

		oldObj.Spec.ManagedSecurityGroups.ControlPlaneNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}
		newObj.Spec.ManagedSecurityGroups.ControlPlaneNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}

		oldObj.Spec.ManagedSecurityGroups.WorkerNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}
		newObj.Spec.ManagedSecurityGroups.WorkerNodesSecurityGroupRules = []infrav1.SecurityGroupRuleSpec{}

		// Allow change to the allowAllInClusterTraffic.
		oldObj.Spec.ManagedSecurityGroups.AllowAllInClusterTraffic = false
		newObj.Spec.ManagedSecurityGroups.AllowAllInClusterTraffic = false
	}

	// Allow changes only to DNSNameservers in ManagedSubnets spec.
	// Subnets are matched by CIDR: the CIDR itself and AllocationPools are
	// immutable, but DNSNameservers can be updated. We zero out DNSNameservers
	// on both copies so the final DeepEqual ignores those changes.
	if newObj.Spec.ManagedSubnets != nil && oldObj.Spec.ManagedSubnets != nil {
		// Check if any fields other than DNSNameservers have changed.
		if len(oldObj.Spec.ManagedSubnets) != len(newObj.Spec.ManagedSubnets) {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "managedSubnets"), "cannot add or remove subnets"))
		} else {
			// Build maps of subnets by CIDR.
			oldSubnetMap := make(map[string]*infrav1.SubnetSpec)
			for i := range oldObj.Spec.ManagedSubnets {
				oldSubnet := &oldObj.Spec.ManagedSubnets[i]
				oldSubnetMap[oldSubnet.CIDR] = oldSubnet
			}

			// Compare each new subnet against the old map by CIDR.
			for i := range newObj.Spec.ManagedSubnets {
				newSubnet := &newObj.Spec.ManagedSubnets[i]

				// If the CIDR is not found in the old map, it's a new subnet.
				oldSubnet, exists := oldSubnetMap[newSubnet.CIDR]
				if !exists {
					allErrs = append(allErrs, field.Forbidden(
						field.NewPath("spec", "managedSubnets"),
						fmt.Sprintf("cannot change subnet CIDR from existing value to %s", newSubnet.CIDR),
					))
					continue
				}

				// Zero out DNSNameservers on both copies so they are ignored by DeepEqual.
				oldSubnet.DNSNameservers = nil
				newSubnet.DNSNameservers = nil
			}
		}
	}

	// Allow changes on AllowedCIDRs and APIServerLB monitors
	if oldLbSpec := oldObj.Spec.APIServer.GetManagedLoadBalancer(); oldLbSpec != nil {
		if newLbSpec := newObj.Spec.APIServer.GetManagedLoadBalancer(); newLbSpec != nil {
			oldLbSpec.AllowedCIDRs = []string{}
			newLbSpec.AllowedCIDRs = []string{}
			oldLbSpec.Monitor = &infrav1.APIServerLoadBalancerMonitor{}
			newLbSpec.Monitor = &infrav1.APIServerLoadBalancerMonitor{}
		}
	}

	// Allow changes to the availability zones.
	oldObj.Spec.ControlPlaneAvailabilityZones = []string{}
	newObj.Spec.ControlPlaneAvailabilityZones = []string{}

	// Allow the scheduling to be changed from CAPI managed to Nova and
	// vice versa.
	oldObj.Spec.ControlPlaneOmitAvailabilityZone = nil
	newObj.Spec.ControlPlaneOmitAvailabilityZone = nil

	// Allow change on the spec.APIServer.FloatingIP only if it matches the current api server loadbalancer IP.
	if oldObj.Status.APIServerManagedLoadBalancer != nil &&
		ptr.Deref(newObj.Spec.APIServer.GetFloatingIP(), "") == oldObj.Status.APIServerManagedLoadBalancer.IP {
		if newObj.Spec.APIServer != nil {
			newObj.Spec.APIServer.FloatingIP = nil
			// If APIServer is now empty after zeroing FloatingIP, nil it out
			// entirely so it compares equal to a nil APIServer on the old object.
			if reflect.DeepEqual(newObj.Spec.APIServer, &infrav1.APIServer{}) {
				newObj.Spec.APIServer = nil
			}
		}
		if oldObj.Spec.APIServer != nil {
			oldObj.Spec.APIServer.FloatingIP = nil
			if reflect.DeepEqual(oldObj.Spec.APIServer, &infrav1.APIServer{}) {
				oldObj.Spec.APIServer = nil
			}
		}
	}

	// Allow changes to primarySubnet to support directing new nodes to a
	// different subnet, e.g. when migrating away from an exhausted subnet.
	oldObj.Spec.PrimarySubnet = nil
	newObj.Spec.PrimarySubnet = nil

	// Allow transitioning spec.network from filter to id and spec.subnets from
	// filter to id. This lets users pin to resolved IDs after initial creation.
	if newObj.Spec.Network != nil && oldObj.Spec.Network != nil && oldObj.Status.Network != nil {
		if allowSubnetFilterToIDTransition(oldObj, newObj) {
			oldObj.Spec.Subnets = nil
			newObj.Spec.Subnets = nil
		}
		if ptr.Deref(newObj.Spec.Network.ID, "") == oldObj.Status.Network.ID {
			newObj.Spec.Network = nil
			oldObj.Spec.Network = nil
		}
	}

	// After zeroing out all mutable fields above, any remaining difference
	// in spec is an immutability violation.
	if !reflect.DeepEqual(oldObj.Spec, newObj.Spec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "cannot be modified"))
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type.
func (*openStackClusterWebhook) ValidateDelete(_ context.Context, _ *infrav1.OpenStackCluster) (admission.Warnings, error) {
	return nil, nil
}

// securityGroupRemoteFields returns whether each remote field is set on a SecurityGroupRuleSpec.
func securityGroupRemoteFields(r *infrav1.SecurityGroupRuleSpec) (bool, bool, bool) {
	return r.RemoteManagedGroups != nil, r.RemoteGroupID != nil, r.RemoteIPPrefix != nil
}

// validateManagedSecurityGroupRules validates that remote* fields are mutually exclusive
// across all security group rule lists in ManagedSecurityGroups.
func validateManagedSecurityGroupRules(msg *infrav1.ManagedSecurityGroups) field.ErrorList {
	var allErrs field.ErrorList //nolint:prealloc // We can't know the number of errors upfront
	allErrs = append(allErrs, validateSecurityGroupRulesRemoteMutualExclusion(
		msg.ClusterNodesSecurityGroupRules,
		field.NewPath("spec", "managedSecurityGroups", "clusterNodesSecurityGroupRules"),
		securityGroupRemoteFields,
	)...)
	allErrs = append(allErrs, validateSecurityGroupRulesRemoteMutualExclusion(
		msg.ControlPlaneNodesSecurityGroupRules,
		field.NewPath("spec", "managedSecurityGroups", "controlPlaneNodesSecurityGroupRules"),
		securityGroupRemoteFields,
	)...)
	allErrs = append(allErrs, validateSecurityGroupRulesRemoteMutualExclusion(
		msg.WorkerNodesSecurityGroupRules,
		field.NewPath("spec", "managedSecurityGroups", "workerNodesSecurityGroupRules"),
		securityGroupRemoteFields,
	)...)
	return allErrs
}
