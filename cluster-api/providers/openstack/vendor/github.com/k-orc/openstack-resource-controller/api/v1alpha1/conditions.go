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

package v1alpha1

import (
	"errors"
	"slices"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Condition Reasons are machine-readable. Treat them like HTTP status
	// codes: they categorise Reasons by 'what should we do about it', not
	// 'what went wrong'. The latter should be in the human-readable
	// Message.

	// Normal progress: continue waiting.
	OpenStackConditionReasonProgressing = "Progressing"

	// The user must fix the configuration before trying again.
	OpenStackConditionReasonInvalidConfiguration = "InvalidConfiguration"

	// An error occurred which we can't recover from. It must be addressed
	// before we can continue.
	OpenStackConditionReasonUnrecoverableError = "UnrecoverableError"

	// An error occurred which may go away eventually if we keep trying. The
	// user likely wants to know about this if it persists.
	OpenStackConditionReasonTransientError = "TransientError"

	// The resource is ready for use.
	OpenStackConditionReasonSuccess = "Success"
)

const (
	OpenStackConditionAvailable   = "Available"
	OpenStackConditionProgressing = "Progressing"
)

// IsConditionReasonTerminal returns true if the given reason represents an error which should prevent further reconciliation.
func IsConditionReasonTerminal(reason string) bool {
	return slices.Contains(
		[]string{
			OpenStackConditionReasonInvalidConfiguration,
			OpenStackConditionReasonUnrecoverableError,
		}, reason)
}

// ObjectWithConditions is a metav1.Object which also stores conditions in its status.
// +kubebuilder:object:generate:=false
type ObjectWithConditions interface {
	metav1.Object
	GetConditions() []metav1.Condition
}

// getUpToDateProgressing returns the progressing condition if and only if it
// exists and is up to date.
func getUpToDateProgressing(obj ObjectWithConditions) *metav1.Condition {
	conditions := obj.GetConditions()
	progressing := meta.FindStatusCondition(conditions, OpenStackConditionProgressing)

	// Not complete if Progressing condition does not exist
	if progressing == nil {
		return nil
	}

	// Not complete if status is out of date
	if progressing.ObservedGeneration != obj.GetGeneration() {
		return nil
	}

	return progressing
}

// IsReconciliationComplete returns true if the given set of conditions indicate that reconciliation is complete, either success or failure.
func IsReconciliationComplete(obj ObjectWithConditions) bool {
	progressing := getUpToDateProgressing(obj)
	if progressing == nil {
		return false
	}

	// Complete if we've either succeeded or failed terminally
	return progressing.Reason == OpenStackConditionReasonSuccess || IsConditionReasonTerminal(progressing.Reason)
}

// GetTerminalError returns an error containing a descriptive message if reconciliation has failed terminally, or nil otherwise.
func GetTerminalError(obj ObjectWithConditions) error {
	progressing := getUpToDateProgressing(obj)
	if progressing == nil {
		return nil
	}

	if IsConditionReasonTerminal(progressing.Reason) {
		return errors.New(progressing.Message)
	}

	return nil
}

func IsAvailable(obj ObjectWithConditions) bool {
	conditions := obj.GetConditions()
	available := meta.FindStatusCondition(conditions, OpenStackConditionAvailable)

	return available != nil && available.Status == metav1.ConditionTrue
}
