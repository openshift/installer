/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package conditions

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/retry"
)

type Reason struct {
	Name                string
	RetryClassification retry.Classification
}

// Auth reasons
var ReasonSubscriptionMismatch = Reason{Name: "SubscriptionMismatch", RetryClassification: retry.Fast}

// Precondition reasons
var (
	ReasonSecretNotFound    = Reason{Name: "SecretNotFound", RetryClassification: retry.Fast}
	ReasonConfigMapNotFound = Reason{Name: "ConfigMapNotFound", RetryClassification: retry.Fast}
	ReasonReferenceNotFound = Reason{Name: "ReferenceNotFound", RetryClassification: retry.Fast}
	ReasonWaitingForOwner   = Reason{Name: "WaitingForOwner", RetryClassification: retry.Fast}
)

// Post-ARM PUT reasons
var (
	ReasonAzureResourceNotFound               = Reason{Name: "AzureResourceNotFound", RetryClassification: retry.Slow}
	ReasonAdditionalKubernetesObjWriteFailure = Reason{Name: "FailedWritingAdditionalKubernetesObjects", RetryClassification: retry.Slow}
)

// Other reasons
var (
	ReasonReconciling                     = Reason{Name: "Reconciling", RetryClassification: retry.Fast}
	ReasonDeleting                        = Reason{Name: "Deleting", RetryClassification: retry.Fast}
	ReasonReconciliationFailedPermanently = Reason{Name: "ReconciliationFailedPermanently", RetryClassification: retry.None}
	ReasonReconcileBlocked                = Reason{Name: "ReconciliationBlocked", RetryClassification: retry.Slow}
	ReasonReconcilePostponed              = Reason{Name: "ReconciliationPostponed", RetryClassification: retry.Slow}
	ReasonPostReconcileFailure            = Reason{Name: "PostReconciliationFailure", RetryClassification: retry.Slow}
)

// ReasonFailed is a catch-all error code for when we don't have a more specific error classification
var ReasonFailed = Reason{Name: "Failed", RetryClassification: retry.Slow}

func MakeReason(reason string, retryClassification retry.Classification) Reason {
	if retryClassification == "" || retryClassification == retry.None { // Unset and none default to slow
		retryClassification = retry.Slow
	}

	return Reason{Name: reason, RetryClassification: retryClassification}
}

func NewReadyConditionBuilder(builder PositiveConditionBuilderInterface) *ReadyConditionBuilder {
	return &ReadyConditionBuilder{
		builder: builder,
	}
}

type ReadyConditionBuilder struct {
	builder PositiveConditionBuilderInterface
}

func (b *ReadyConditionBuilder) ReadyCondition(severity ConditionSeverity, observedGeneration int64, reason string, message string) Condition {
	return b.builder.MakeFalseCondition(
		ConditionTypeReady,
		severity,
		observedGeneration,
		reason,
		message)
}

func (b *ReadyConditionBuilder) Reconciling(observedGeneration int64) Condition {
	return b.builder.MakeFalseCondition(
		ConditionTypeReady,
		ConditionSeverityInfo,
		observedGeneration,
		ReasonReconciling.Name,
		"The resource is in the process of being reconciled by the operator")
}

func (b *ReadyConditionBuilder) Deleting(observedGeneration int64) Condition {
	return b.builder.MakeFalseCondition(
		ConditionTypeReady,
		ConditionSeverityInfo,
		observedGeneration,
		ReasonDeleting.Name,
		"The resource is being deleted")
}

func (b *ReadyConditionBuilder) Succeeded(observedGeneration int64) Condition {
	return b.builder.MakeTrueCondition(ConditionTypeReady, observedGeneration)
}
