package core

// (C) Copyright IBM Corp. 2024.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"errors"
)

// SDKProblem provides a type suited to problems that
// occur in SDK projects. It extends the base
// "IBMProblem" type with a field to store the
// function being called when the problem occurs.
type SDKProblem struct {
	*IBMProblem

	// Function provides the name of the in-code
	// function or method in which the problem
	// occurred.
	Function string

	// A computed stack trace including the relevant
	// function names, files, and line numbers invoked
	// leading up to the origination of the problem.
	stack []sdkStackFrame
}

// GetConsoleMessage returns all public fields of
// the problem, formatted in YAML.
func (e *SDKProblem) GetConsoleMessage() string {
	return ComputeConsoleMessage(e)
}

// GetDebugMessage returns all information
// about the problem, formatted in YAML.
func (e *SDKProblem) GetDebugMessage() string {
	return ComputeDebugMessage(e)
}

// GetID returns the computed identifier, computed from the
// "Component", "discriminator", and "Function" fields, as well as the
// identifier of the "causedBy" problem, if it exists.
func (e *SDKProblem) GetID() string {
	return CreateIDHash("sdk", e.GetBaseSignature(), e.Function)
}

// Is allows an SDKProblem instance to be compared against another error for equality.
// An SDKProblem is considered equal to another error if 1) the error is also a Problem and
// 2) it has the same ID (i.e. it is the same problem scenario).
func (e *SDKProblem) Is(target error) bool {
	return is(target, e.GetID())
}

// GetConsoleOrderedMaps returns an ordered-map representation
// of an SDKProblem instance suited for a console message.
func (e *SDKProblem) GetConsoleOrderedMaps() *OrderedMaps {
	orderedMaps := NewOrderedMaps()

	orderedMaps.Add("id", e.GetID())
	orderedMaps.Add("summary", e.Summary)
	orderedMaps.Add("severity", e.Severity)
	orderedMaps.Add("function", e.Function)
	orderedMaps.Add("component", e.Component)

	return orderedMaps
}

// GetDebugOrderedMaps returns an ordered-map representation
// of an SDKProblem instance, with additional information
// suited for a debug message.
func (e *SDKProblem) GetDebugOrderedMaps() *OrderedMaps {
	orderedMaps := e.GetConsoleOrderedMaps()

	orderedMaps.Add("stack", e.stack)

	var orderableCausedBy OrderableProblem
	if errors.As(e.GetCausedBy(), &orderableCausedBy) {
		orderedMaps.Add("caused_by", orderableCausedBy.GetDebugOrderedMaps().GetMaps())
	}

	return orderedMaps
}

// SDKErrorf creates and returns a new instance of "SDKProblem" with "error" level severity.
func SDKErrorf(err error, summary, discriminator string, component *ProblemComponent) *SDKProblem {
	function := computeFunctionName(component.Name)
	stack := getStackInfo(component.Name)

	return &SDKProblem{
		IBMProblem: IBMErrorf(err, component, summary, discriminator),
		Function:   function,
		stack:      stack,
	}
}

// RepurposeSDKProblem provides a convenient way to take a problem from
// another function in the same component and contextualize it to the current
// function. Should only be used in public (exported) functions.
func RepurposeSDKProblem(err error, discriminator string) error {
	if err == nil {
		return nil
	}

	// It only makes sense to carry out this logic with SDK Errors.
	var sdkErr *SDKProblem
	if !errors.As(err, &sdkErr) {
		return err
	}

	// Special behavior to allow SDK problems coming from a method that wraps a
	// "*WithContext" method to maintain the discriminator of the originating
	// problem. Otherwise, we would lose all of that data in the wrap.
	if discriminator != "" {
		sdkErr.discriminator = discriminator
	}

	// Recompute the function to reflect this public boundary (but let the stack
	// remain as it is - it is the path to the original problem origination point).
	sdkErr.Function = computeFunctionName(sdkErr.Component.Name)

	return sdkErr
}
