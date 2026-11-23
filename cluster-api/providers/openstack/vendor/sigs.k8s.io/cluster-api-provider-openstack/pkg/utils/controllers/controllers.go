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

package controllers

import (
	"context"
	"fmt"
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	conditions "sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// ValidateSubnets validates that all subnet CIDRs are parseable.
// Multiple subnets of the same IP version are allowed; use PrimarySubnet to
// specify which subnet should be used for load balancer VIP allocation and
// member registration when multiple subnets are present.
func ValidateSubnets(subnets []infrav1.Subnet) error {
	for _, subnet := range subnets {
		if _, _, err := net.ParseCIDR(subnet.CIDR); err != nil {
			return fmt.Errorf("invalid CIDR %q in subnet %q: %w", subnet.CIDR, subnet.ID, err)
		}
	}
	return nil
}

func GetInfraCluster(ctx context.Context, c client.Client, cluster *clusterv1.Cluster) (*infrav1.OpenStackCluster, error) {
	openStackCluster := &infrav1.OpenStackCluster{}
	openStackClusterName := client.ObjectKey{
		Namespace: cluster.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := c.Get(ctx, openStackClusterName, openStackCluster); err != nil {
		return nil, err
	}
	return openStackCluster, nil
}

const (
	// legacyConditionReasonTrue is the reason assigned to a legacy True condition that has no reason.
	legacyConditionReasonTrue = infrav1.ReadyConditionReason
	// legacyConditionReasonFalse is the reason assigned to a legacy False condition that has no reason.
	legacyConditionReasonFalse = "NotReady"
	// legacyConditionReasonUnknown is the reason assigned to a legacy Unknown condition that has no reason.
	legacyConditionReasonUnknown = "Unknown"
)

// EnsureConditionReasons sets a reason on every condition of obj that has an empty one.
//
// Conditions written by releases that used the Cluster API v1beta1 condition
// type (CAPO < 0.15) may have no reason, which is not valid for metav1.Condition
// and is rejected by the CRD schema. A status patch rewrites the full
// conditions list, so a single legacy condition would make every later patch fail.
//
// This must be called after the patch helper has been created, so that the
// fix is part of the computed patch. It returns true if any condition was changed.
func EnsureConditionReasons(obj conditions.Setter) bool {
	existing := obj.GetConditions()
	if len(existing) == 0 {
		return false
	}

	fixed := make([]metav1.Condition, len(existing))
	copy(fixed, existing)

	changed := false
	for i := range fixed {
		if fixed[i].Reason != "" {
			continue
		}
		switch fixed[i].Status {
		case metav1.ConditionTrue:
			fixed[i].Reason = legacyConditionReasonTrue
		case metav1.ConditionFalse:
			fixed[i].Reason = legacyConditionReasonFalse
		default:
			fixed[i].Reason = legacyConditionReasonUnknown
		}
		changed = true
	}

	if changed {
		obj.SetConditions(fixed)
	}
	return changed
}
