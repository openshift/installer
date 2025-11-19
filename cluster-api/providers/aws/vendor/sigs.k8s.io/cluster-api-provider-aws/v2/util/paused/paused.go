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

// This is a copy of the paused helpers from cluster-api, without the
// v1beta2conditions requirement. Once the v1beta2conditions migration is
// complete, we will want to use the upstream package.

// Package paused implements paused helper functions.
package paused

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
)

// ConditionSetter combines the client.Object and Setter interface.
type ConditionSetter interface {
	conditions.Setter
	client.Object
}

// EnsurePausedCondition sets the paused condition on the object and returns if it should be considered as paused.
func EnsurePausedCondition(ctx context.Context, c client.Client, cluster *clusterv1.Cluster, obj ConditionSetter) (isPaused bool, conditionChanged bool, err error) {
	oldCondition := conditions.Get(obj, clusterv1.PausedV1Beta2Condition)
	newCondition := pausedCondition(c.Scheme(), cluster, obj, clusterv1.PausedV1Beta2Condition)

	isPaused = newCondition.Status == corev1.ConditionTrue

	log := ctrl.LoggerFrom(ctx)

	// Return early if the paused condition did not change.
	if oldCondition != nil && conditions.HasSameState(oldCondition, &newCondition) {
		if isPaused {
			log.V(6).Info("Reconciliation is paused for this object", "reason", newCondition.Message)
		}
		return isPaused, false, nil
	}

	patchHelper, err := patch.NewHelper(obj, c)
	if err != nil {
		return isPaused, false, err
	}

	if isPaused {
		log.V(4).Info("Pausing reconciliation for this object", "reason", newCondition.Message)
	} else {
		log.V(4).Info("Unpausing reconciliation for this object")
	}

	conditions.Set(obj, &newCondition)

	if err := patchHelper.Patch(ctx, obj, patch.WithOwnedV1Beta2Conditions{Conditions: []string{
		clusterv1.PausedV1Beta2Condition,
	}}); err != nil {
		return isPaused, false, err
	}

	return isPaused, true, nil
}

// pausedCondition sets the paused condition on the object and returns if it should be considered as paused.
func pausedCondition(scheme *runtime.Scheme, cluster *clusterv1.Cluster, obj ConditionSetter, targetConditionType string) clusterv1.Condition {
	if (cluster != nil && cluster.Spec.Paused) || annotations.HasPaused(obj) {
		var messages []string
		if cluster != nil && cluster.Spec.Paused {
			messages = append(messages, "Cluster spec.paused is set to true")
		}
		if annotations.HasPaused(obj) {
			kind := "Object"
			if gvk, err := apiutil.GVKForObject(obj, scheme); err == nil {
				kind = gvk.Kind
			}
			messages = append(messages, fmt.Sprintf("%s has the cluster.x-k8s.io/paused annotation", kind))
		}

		return clusterv1.Condition{
			Type:    clusterv1.ConditionType(targetConditionType),
			Status:  corev1.ConditionTrue,
			Reason:  clusterv1.PausedV1Beta2Reason,
			Message: strings.Join(messages, ", "),
		}
	}

	return clusterv1.Condition{
		Type:   clusterv1.ConditionType(targetConditionType),
		Status: corev1.ConditionFalse,
		Reason: clusterv1.NotPausedV1Beta2Reason,
	}
}
