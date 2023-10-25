/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package conditions

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Conditions []Condition

// FindIndexByType returns the index of the condition with the given ConditionType if it exists.
// If the Condition with the specified ConditionType is doesn't exist, the second boolean parameter is false.
func (c Conditions) FindIndexByType(conditionType ConditionType) (int, bool) {
	for i := range c {
		condition := c[i]
		if condition.Type == conditionType {
			return i, true
		}
	}

	return -1, false
}

// TODO: Hah, name...
type Conditioner interface {
	GetConditions() Conditions
	SetConditions(conditions Conditions)
}

// ConditionSeverity expresses the severity of a Condition.
type ConditionSeverity string

const (
	// ConditionSeverityError specifies that a failure of a condition type
	// should be viewed as an error. Errors are fatal to reconciliation and
	// mean that the user must take some action to resolve
	// the problem before reconciliation will be attempted again.
	ConditionSeverityError ConditionSeverity = "Error"

	// ConditionSeverityWarning specifies that a failure of a condition type
	// should be viewed as a warning. Warnings are informational. The operator
	// may be able to retry and resolve the warning without any action from the user, but
	// in some cases user action to resolve the warning will be required.
	ConditionSeverityWarning ConditionSeverity = "Warning"

	// ConditionSeverityInfo specifies that a failure of a condition type
	// should be viewed as purely informational. Things are working.
	// This is the happy path.
	ConditionSeverityInfo ConditionSeverity = "Info"

	// ConditionSeverityNone specifies that there is no condition severity.
	// For conditions which have positive polarity (Status == True is their normal/healthy state), this will set when Status == True
	// For conditions which have negative polarity (Status == False is their normal/healthy state), this will be set when Status == False.
	// Conditions in Status == Unknown always have a severity of None as well.
	// This is the default state for conditions.
	ConditionSeverityNone ConditionSeverity = ""
)

type ConditionType string

const (
	// ConditionTypeReady is a condition indicating if the resource is ready or not.
	// A ready resource is one that has been successfully provisioned to Azure according to the
	// resource spec. It has reached the goal state. This usually means that the resource is ready
	// to use, but the exact meaning of Ready may vary slightly from resource to resource. Resources with
	// caveats to Ready's meaning will call that out in the resource specific documentation.
	ConditionTypeReady = "Ready"
)

var _ fmt.Stringer = Condition{}

// Condition defines an extension to status (an observation) of a resource
// +kubebuilder:object:generate=true
type Condition struct {
	// Type of condition.
	// +kubebuilder:validation:Required
	Type ConditionType `json:"type"`

	// Status of the condition, one of True, False, or Unknown.
	// +kubebuilder:validation:Required
	Status metav1.ConditionStatus `json:"status"`

	// Severity with which to treat failures of this type of condition.
	// For conditions which have positive polarity (Status == True is their normal/healthy state), this will be omitted when Status == True
	// For conditions which have negative polarity (Status == False is their normal/healthy state), this will be omitted when Status == False.
	// This is omitted in all cases when Status == Unknown
	// +kubebuilder:validation:Optional
	Severity ConditionSeverity `json:"severity,omitempty"`

	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +kubebuilder:validation:Required
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// Note: see the https://github.com/kubernetes/enhancements/tree/master/keps/sig-api-machinery/1623-standardize-conditions
	// KEP for details about ObservedGeneration

	// ObservedGeneration is the .metadata.generation that the condition was set based upon. For instance, if
	// .metadata.generation is currently 12, but the .status.condition[x].observedGeneration is 9, the condition is out of date
	// with respect to the current state of the instance.
	// +kubebuilder:validation:Optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Reason for the condition's last transition.
	// Reasons are upper CamelCase (PascalCase) with no spaces. A reason is always provided, this field will not be empty.
	// +kubebuilder:validation:Required
	Reason string `json:"reason"`

	// Message is a human readable message indicating details about the transition. This field may be empty.
	// +kubebuilder:validation:Optional
	Message string `json:"message,omitempty"`
}

// IsEquivalent returns true if this condition is equivalent to the passed in condition.
// Two conditions are equivalent if all of their fields EXCEPT LastTransitionTime are the same.
func (c Condition) IsEquivalent(other Condition) bool {
	return c.Type == other.Type &&
		c.Status == other.Status &&
		c.Severity == other.Severity &&
		c.Reason == other.Reason &&
		c.ObservedGeneration == other.ObservedGeneration &&
		c.Message == other.Message
}

// ShouldOverwrite determines if this condition should overwrite the other condition.
func (c Condition) ShouldOverwrite(other Condition) bool {
	// Safety check that the two conditions are of the same type. If not they certainly shouldn't overwrite
	if c.Type != other.Type {
		return false
	}
	// If this condition corresponds to a newer generation than the previous condition always overwrite
	// as other is out of date
	if c.ObservedGeneration > other.ObservedGeneration {
		return true
	}
	// If the Conditions are equivalent, don't overwrite. We want to keep the first occurrence of the condition
	// so that the LastTransitionTime is correct
	if c.IsEquivalent(other) {
		return false
	}
	// At this point the Conditions are the same type, same generation, so the winning Condition must be chosen
	// based on priority

	if c.priority() >= other.priority() {
		return true
	} else {
		return false
	}
}

// priority of the condition for overwrite purposes if Condition ObservedGeneration's are the same. Higher is more important.
// The result of this is the following:
//  1. Status == True conditions, and Status == False conditions with Severity == Warning or Error are all the highest priority.
//     This means that the most recent update with any of those states will overwrite anything.
//  2. Status == False conditions with Severity == Info will only overwrite other Status == False conditions with Severity == Info.
//  3. Status == Unknown conditions will not overwrite anything.
//
// Keep in mind that this priority is specifically for comparing Conditions with the same ObservedGeneration. If the ObservedGeneration
// is different, the newer one always wins.
func (c Condition) priority() int {
	switch c.Status {
	case metav1.ConditionTrue:
		return 5
	case metav1.ConditionFalse:
		switch c.Severity {
		case ConditionSeverityError:
			return 5
		case ConditionSeverityWarning:
			return 5
		case ConditionSeverityInfo:
			return 4
		case ConditionSeverityNone:
			// This shouldn't happen as a Condition with Status False should always specify a severity.
			// In the interest of safety though, we set this to 5 so if this DOES somehow happen it ties
			// or wins against most other things and users will see it
			return 5
		}
	case metav1.ConditionUnknown:
		return 3
	}

	// This shouldn't happen
	return 0
}

// Copy returns an independent copy of the Condition
func (c Condition) Copy() Condition {
	// NB: If you change this to a non-simple copy
	// you will need to update genruntime.CloneSliceOfCondition
	return c
}

// String returns a string representation of this condition
func (c Condition) String() string {
	return fmt.Sprintf(
		"Condition [%s], Status = %q, ObservedGeneration = %d, Severity = %q, Reason = %q, Message = %q, LastTransitionTime = %q",
		c.Type,
		c.Status,
		c.ObservedGeneration,
		c.Severity,
		c.Reason,
		c.Message,
		c.LastTransitionTime)
}

// SetCondition sets the provided Condition on the Conditioner. The condition is only
// set if the new condition is different from the existing condition of the same type.
// See Condition.IsEquivalent and Condition.ShouldOverwrite for more details.
func SetCondition(o Conditioner, new Condition) {
	setCondition(o, new, func(new Condition, old Condition) bool { return new.ShouldOverwrite(old) })
}

// Reasons other than those explicitly called out here have the default priority of 0.
// These are given a negative priority so that they "lose" to the default of 0 and are overwritten.
// Take care to not modify this structure (Golang doesn't support a readonly map). We could use
// v2/tools/generator/internal/readonly/readonly_map.go but given
// Golang generics bugs for now we go with the simpler approach
var reasonPriority = map[string]int{
	ReasonReferenceNotFound.Name: -2,
	ReasonSecretNotFound.Name:    -2,
	ReasonConfigMapNotFound.Name: -2,
	ReasonWaitingForOwner.Name:   -2,
	ReasonReconciling.Name:       -1,
}

// SetConditionReasonAware sets the provided Condition on the Conditioner. This is similar to SetCondition
// with one difference: SetConditionReasonAware understands common Reasons used by ASO and allows some of them to
// modify the standard Condition priority rules. This is primarily used to allow the Reconciling condition to overwrite
// Warning conditions raised by the operator that have been fixed. This is useful because sometimes getting a success or
// error from Azure can take a long time, and workflows like: submit -> warning -> fix warning -> call Azure -> wait -> success
// otherwise would continue reporting the Warning Condition until the final success step (possibly many minutes after the
// warning was resolved).
func SetConditionReasonAware(o Conditioner, new Condition) {
	shouldOverwrite := func(new Condition, old Condition) bool {
		if new.ShouldOverwrite(old) {
			return true
		}

		// If we normally wouldn't overwrite, check the reason of the old and new condition and compare their priorities
		oldPriority := reasonPriority[old.Reason] // Default is 0 if not mapped
		newPriority := reasonPriority[new.Reason] // Default is 0 if not mapped
		return newPriority > oldPriority          // Just > rather than >= here to prevent things overwriting themselves
	}

	setCondition(o, new, shouldOverwrite)
}

func setCondition(o Conditioner, new Condition, shouldOverwrite func(new Condition, old Condition) bool) {
	if o == nil {
		return
	}

	conditions := o.GetConditions()
	i, exists := conditions.FindIndexByType(new.Type)
	if exists {
		if !shouldOverwrite(new, conditions[i]) {
			// Nothing to do, the new condition is not supposed to overwrite
			return
		}
		conditions[i] = new
	} else {
		conditions = append(conditions, new)
	}

	// TODO: do we sort conditions here? CAPI does.

	o.SetConditions(conditions)
}

// GetCondition gets the Condition with the specified type from the provided Conditioner.
// Returns the Condition and true if a Condition with the specified type is found, or an empty Condition
// and false if not.
func GetCondition(o Conditioner, conditionType ConditionType) (Condition, bool) {
	if o == nil {
		return Condition{}, false
	}

	conditions := o.GetConditions()
	i, exists := conditions.FindIndexByType(conditionType)
	if exists {
		return conditions[i], true
	}

	return Condition{}, false
}
