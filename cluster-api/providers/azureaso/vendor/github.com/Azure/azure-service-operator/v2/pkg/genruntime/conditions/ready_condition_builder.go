/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package conditions

type RetryClassification string

const (
	// RetryNone means that this classification is not expected to ever retry (it's only ever set on for fatal errors)
	RetryNone = RetryClassification("None") // TODO: ??
	RetryFast = RetryClassification("RetryFast")
	RetrySlow = RetryClassification("RetrySlow")
)

type Reason struct {
	Name                string
	RetryClassification RetryClassification
}

// Auth reasons
var ReasonSubscriptionMismatch = Reason{Name: "SubscriptionMismatch", RetryClassification: RetryFast}

// Precondition reasons
var ReasonSecretNotFound = Reason{Name: "SecretNotFound", RetryClassification: RetryFast}
var ReasonConfigMapNotFound = Reason{Name: "ConfigMapNotFound", RetryClassification: RetryFast}
var ReasonReferenceNotFound = Reason{Name: "ReferenceNotFound", RetryClassification: RetryFast}
var ReasonWaitingForOwner = Reason{Name: "WaitingForOwner", RetryClassification: RetryFast}

// Post-ARM PUT reasons
var ReasonAzureResourceNotFound = Reason{Name: "AzureResourceNotFound", RetryClassification: RetrySlow}
var ReasonAdditionalKubernetesObjWriteFailure = Reason{Name: "FailedWritingAdditionalKubernetesObjects", RetryClassification: RetrySlow}

// Other reasons
var ReasonReconciling = Reason{Name: "Reconciling", RetryClassification: RetryFast}
var ReasonDeleting = Reason{Name: "Deleting", RetryClassification: RetryFast}
var ReasonReconciliationFailedPermanently = Reason{Name: "ReconciliationFailedPermanently", RetryClassification: RetryNone}
var ReasonReconcileBlocked = Reason{Name: "ReconciliationBlocked", RetryClassification: RetrySlow}
var ReasonReconcilePostponed = Reason{Name: "ReconciliationPostponed", RetryClassification: RetrySlow}
var ReasonPostReconcileFailure = Reason{Name: "PostReconciliationFailure", RetryClassification: RetrySlow}

// ReasonFailed is a catch-all error code for when we don't have a more specific error classification
var ReasonFailed = Reason{Name: "Failed", RetryClassification: RetrySlow}

func MakeReason(reason string) Reason {
	return Reason{Name: reason, RetryClassification: RetrySlow} // Always classify custom reasons as Slow retry
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
