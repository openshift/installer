/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package conditions

import (
	"time"

	"github.com/benbjohnson/clock"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const ReasonSucceeded = "Succeeded"

// TODO: name?
type PositiveConditionBuilderInterface interface {
	MakeTrueCondition(conditionType ConditionType, observedGeneration int64) Condition
	MakeFalseCondition(conditionType ConditionType, severity ConditionSeverity, observedGeneration int64, reason string, message string) Condition
	MakeUnknownCondition(conditionType ConditionType, observedGeneration int64, reason string, message string) Condition
}

var _ PositiveConditionBuilderInterface = &PositiveConditionBuilder{}

// PositiveConditionBuilder makes positive polarity conditions. A positive polarity condition
// is a condition where Status==True is the "happy path" state, and Status==False is a failure
// state.
type PositiveConditionBuilder struct {
	clock clock.Clock

	Ready *ReadyConditionBuilder
}

// NewPositiveConditionBuilder creates a new PositiveConditionBuilder for creating positive polarity conditions.
func NewPositiveConditionBuilder(clock clock.Clock) *PositiveConditionBuilder {
	result := &PositiveConditionBuilder{
		clock: clock,
	}

	result.Ready = NewReadyConditionBuilder(result)
	return result
}

func (b *PositiveConditionBuilder) now() metav1.Time {
	return metav1.NewTime(b.clock.Now().UTC().Truncate(time.Second))
}

// MakeTrueCondition makes a condition whose Status is True
func (b *PositiveConditionBuilder) MakeTrueCondition(conditionType ConditionType, observedGeneration int64) Condition {
	return Condition{
		Type:               conditionType,
		Status:             metav1.ConditionTrue,
		Severity:           ConditionSeverityNone,
		LastTransitionTime: b.now(),
		ObservedGeneration: observedGeneration,
		Reason:             ReasonSucceeded,
	}
}

// MakeFalseCondition makes a condition whose Status is False. A severity, reason, and message must be provided.
func (b *PositiveConditionBuilder) MakeFalseCondition(conditionType ConditionType, severity ConditionSeverity, observedGeneration int64, reason string, message string) Condition {
	return Condition{
		Type:               conditionType,
		Status:             metav1.ConditionFalse,
		Severity:           severity,
		LastTransitionTime: b.now(),
		ObservedGeneration: observedGeneration,
		Reason:             reason,
		Message:            message,
	}
}

// MakeUnknownCondition makes a condition whose Status is Unknown. A reason and message must be provided. No severity
// is required as conditions in Status Unknown do not have a known severity either.
func (b *PositiveConditionBuilder) MakeUnknownCondition(conditionType ConditionType, observedGeneration int64, reason string, message string) Condition {
	return Condition{
		Type:               conditionType,
		Status:             metav1.ConditionUnknown,
		Severity:           ConditionSeverityNone,
		LastTransitionTime: b.now(),
		ObservedGeneration: observedGeneration,
		Reason:             reason,
		Message:            message,
	}
}
