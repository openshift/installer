/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package reconcilers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/ownerutil"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// LogObj logs the obj
func LogObj(log logr.Logger, level int, note string, obj genruntime.MetaObject) {
	if log.V(level).Enabled() {
		ourAnnotations := make(map[string]string)
		for key, value := range obj.GetAnnotations() {
			if strings.HasPrefix(key, "serviceoperator.azure.com") {
				ourAnnotations[key] = value
			}
		}

		keysAndValues := []interface{}{
			"kind", obj.GetObjectKind(),
			"resourceVersion", obj.GetResourceVersion(),
			"generation", obj.GetGeneration(),
			"uid", obj.GetUID(),
			"ownerReferences", obj.GetOwnerReferences(),
			"creationTimestamp", obj.GetCreationTimestamp(),
			"deletionTimestamp", obj.GetDeletionTimestamp(),
			"finalizers", obj.GetFinalizers(),
			"annotations", ourAnnotations,
			// Use fmt here to ensure the output uses the String() method, which log.Info doesn't seem to do by default
			"conditions", fmt.Sprintf("%s", obj.GetConditions()),
		}

		if armObj, ok := obj.(genruntime.ARMMetaObject); ok {
			keysAndValues = append(keysAndValues, "owner", armObj.Owner())
		}

		// Log just what we're interested in. We avoid logging the whole obj
		// due to possible risk of disclosing secrets or other data that is "private" and users may
		// not want in logs.
		// nolint: logrlint // We can see the keys and values, above, but the linter can't
		log.V(level).Info(note, keysAndValues...)
	}
}

type ARMOwnedResourceReconcilerCommon struct {
	ReconcilerCommon
	ResourceResolver *resolver.Resolver
}

// NeedsToWaitForOwner returns false if the owner doesn't need to be waited for, and true if it does.
func (r *ARMOwnedResourceReconcilerCommon) NeedsToWaitForOwner(ctx context.Context, log logr.Logger, obj genruntime.ARMOwnedMetaObject) (bool, error) {
	ownerDetails, err := r.ResourceResolver.ResolveOwner(ctx, obj)
	if err != nil {
		var typedErr *core.ReferenceNotFound
		if errors.As(err, &typedErr) {
			log.V(Info).Info("Owner does not yet exist", "NamespacedName", typedErr.NamespacedName)
			return true, nil
		}

		return true, errors.Wrap(err, "failed to get owner")
	}

	// No need to wait for resources that don't have an owner
	if !ownerDetails.FoundKubernetesOwner() {
		return false, nil
	}

	// If the owner isn't ready, wait
	ready := genruntime.GetReadyCondition(ownerDetails.Owner)
	isOwnerReady := ready != nil && ready.Status == metav1.ConditionTrue
	if !isOwnerReady {
		var readyStr string
		if ready == nil {
			readyStr = "<nil>"
		} else {
			readyStr = ready.String()
		}
		log.V(Info).Info("Owner exists but is not ready. Current condition", "ready", readyStr)
		return true, nil
	}

	return false, nil
}

func (r *ARMOwnedResourceReconcilerCommon) ApplyOwnership(ctx context.Context, log logr.Logger, obj genruntime.ARMOwnedMetaObject) error {
	ownerDetails, err := r.ResourceResolver.ResolveOwner(ctx, obj)
	if err != nil {
		return errors.Wrap(err, "failed to get owner")
	}

	if !ownerDetails.FoundKubernetesOwner() {
		// If no owner is expected or the owner is only in ARM, no need to assign ownership in Kubernetes
		return nil
	}

	ownerRef := ownerutil.MakeOwnerReference(ownerDetails.Owner)

	obj.SetOwnerReferences(ownerutil.EnsureOwnerRef(obj.GetOwnerReferences(), ownerRef))
	log.V(Info).Info(
		"Set owner reference",
		"ownerGvk", ownerDetails.Owner.GetObjectKind().GroupVersionKind(),
		"ownerName", ownerDetails.Owner.GetName())

	return nil
}

// ClaimResource ensures that the owner reference is set
func (r *ARMOwnedResourceReconcilerCommon) ClaimResource(ctx context.Context, log logr.Logger, obj genruntime.ARMOwnedMetaObject) error {
	log.V(Info).Info("applying ownership")
	waitForOwner, err := r.NeedsToWaitForOwner(ctx, log, obj)
	if err != nil {
		return err
	}

	if waitForOwner {
		err = errors.Errorf("Owner %q cannot be found. Progress is blocked until the owner is created.", obj.Owner().String())
		err = conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonWaitingForOwner)
		return err
	}

	// Short circuit here if there's no owner management to do
	if obj.Owner() == nil {
		return nil
	}

	err = r.ApplyOwnership(ctx, log, obj)
	if err != nil {
		return err
	}

	return nil
}

type ReconcilerCommon struct {
	KubeClient         kubeclient.Client
	PositiveConditions *conditions.PositiveConditionBuilder
}

func ClassifyResolverError(err error) error {
	// If it's specifically secret not found, say so
	var secretErr *core.SecretNotFound
	if errors.As(err, &secretErr) {
		return conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonSecretNotFound)
	}

	// If it's specifically configmap not found, say so
	var configMapErr *core.ConfigMapNotFound
	if errors.As(err, &configMapErr) {
		return conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonConfigMapNotFound)
	}

	// If it's subscription mismatch, classify that
	var subscriptionMismatchErr *core.SubscriptionMismatch
	if errors.As(err, &subscriptionMismatchErr) {
		return conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityError, conditions.ReasonFailed)
	}

	// Everything else is ReferenceNotFound. This is maybe a bit of a lie but secrets are also references and we want to make sure
	// everything is classified as something, so for now it's good enough.
	return conditions.NewReadyConditionImpactingError(err, conditions.ConditionSeverityWarning, conditions.ReasonReferenceNotFound)
}
