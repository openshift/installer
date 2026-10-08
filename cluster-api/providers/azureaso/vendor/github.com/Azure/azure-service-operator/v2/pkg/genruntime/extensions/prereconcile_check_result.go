/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"github.com/rotisserie/eris"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

type PreReconcileCheckResult struct {
	action   preReconcileCheckResultType
	severity conditions.ConditionSeverity
	reason   conditions.Reason
	message  string
}

// ProceedWithReconcile indicates that a resource is ready for reconciliation by returning a PreReconcileCheckResult
// with action `Proceed`.
func ProceedWithReconcile() PreReconcileCheckResult {
	return PreReconcileCheckResult{
		action: preReconcileCheckResultTypeProceed,
	}
}

// BlockReconcile indicates reconciliation of a resource is currently blocked by returning a PreReconcileCheckResult
// with action `Block`. The reconciliation will automatically be retried after a short delay.
// reason is an explanatory reason to show to the user via a warning condition on the resource.
func BlockReconcile(reason string) PreReconcileCheckResult {
	return PreReconcileCheckResult{
		action:   preReconcileCheckResultTypeBlock,
		severity: conditions.ConditionSeverityWarning,
		reason:   conditions.ReasonReconcileBlocked,
		message:  reason,
	}
}

// PostponeReconcile indicates reconciliation of a resource is not currently required by returning a
// PreReconcileCheckResult with action `Postpone`. Reconciliation will not be retried until the usual scheduled check.
func PostponeReconcile() PreReconcileCheckResult {
	return PreReconcileCheckResult{
		action:   preReconcileCheckResultTypePostpone,
		severity: conditions.ConditionSeverityInfo,
		reason:   conditions.ReasonReconcilePostponed,
	}
}

func (r PreReconcileCheckResult) PostponeReconciliation() bool {
	return r.action == preReconcileCheckResultTypePostpone
}

func (r PreReconcileCheckResult) BlockReconciliation() bool {
	return r.action == preReconcileCheckResultTypeBlock
}

func (r PreReconcileCheckResult) Message() string {
	return r.message
}

// CreateConditionError returns an error that can be used to set a condition on the resource.
func (r PreReconcileCheckResult) CreateConditionError() error {
	return conditions.NewReadyConditionImpactingError(
		eris.New(r.message),
		r.severity,
		r.reason)
}

// PreReconcileCheckResultType is the type of result returned by PreReconcileCheck.
type preReconcileCheckResultType string

const (
	preReconcileCheckResultTypeBlock    preReconcileCheckResultType = "Block"
	preReconcileCheckResultTypeProceed  preReconcileCheckResultType = "Proceed"
	preReconcileCheckResultTypePostpone preReconcileCheckResultType = "Postpone"
)
