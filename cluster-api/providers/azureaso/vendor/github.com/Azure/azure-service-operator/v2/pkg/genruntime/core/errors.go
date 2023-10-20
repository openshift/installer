/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package core

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
)

func AsTypedError[T error](err error) (T, bool) {
	var typedErr T
	if errors.As(err, &typedErr) {
		return typedErr, true
	}

	// Also deal with the possibility that this is a kerrors.Aggregate
	var aggregate kerrors.Aggregate
	if errors.As(err, &aggregate) {
		for _, e := range aggregate.Errors() {
			// This is a bit hacky but allows us to pick out the first error and raise on that
			if result, ok := AsTypedError[T](e); ok {
				return result, true
			}
		}
	}

	return typedErr, false
}

func AsNotOwnedError(err error) (*NotOwnedError, bool) {
	return AsTypedError[*NotOwnedError](err)
}

// NotOwnedError indicates the target resource is not owned by the resource attempting to write it
type NotOwnedError struct {
	Namespace  string
	TargetName string
	TargetType string
	SourceName string
	SourceType string
}

func NewNotOwnedError(namespace string, name string, gvk schema.GroupVersionKind, sourceName string, sourceGvk schema.GroupVersionKind) error {
	kindStr := gvk.GroupKind().String()
	sourceKindStr := sourceGvk.GroupKind().String()
	return &NotOwnedError{
		Namespace:  namespace,
		TargetName: name,
		TargetType: kindStr,
		SourceName: sourceName,
		SourceType: sourceKindStr,
	}
}

var _ error = &NotOwnedError{}

func (e *NotOwnedError) Error() string {
	return fmt.Sprintf("cannot overwrite %s %s/%s which is not owned by %s %s/%s",
		e.TargetType,
		e.Namespace,
		e.TargetName,
		e.SourceType,
		e.Namespace,
		e.SourceName)
}

type causer interface {
	error
	Cause() error

	// Note that we use Cause() and not Unwrap here because we don't want these errors mistakenly classified as generic
	// NotFound errors (which are ignored and retried).
}

type ReferenceNotFound struct {
	NamespacedName types.NamespacedName
	cause          error
}

func NewReferenceNotFoundError(name types.NamespacedName, cause error) *ReferenceNotFound {
	return &ReferenceNotFound{
		NamespacedName: name,
		cause:          cause,
	}
}

var _ error = &ReferenceNotFound{}
var _ causer = &ReferenceNotFound{}

func (e *ReferenceNotFound) Error() string {
	return fmt.Sprintf("%s does not exist (%s)", e.NamespacedName, e.cause)
}

func (e *ReferenceNotFound) Is(err error) bool {
	var typedErr *ReferenceNotFound
	if errors.As(err, &typedErr) {
		return e.NamespacedName == typedErr.NamespacedName
	}
	return false
}

func (e *ReferenceNotFound) Cause() error {
	return e.cause
}

func (e *ReferenceNotFound) Format(s fmt.State, verb rune) {
	format(e, s, verb)
}

// SecretNotFound error is used when secret or its expected keys are not found
type SecretNotFound struct {
	NamespacedName types.NamespacedName
	cause          error
}

func NewSecretNotFoundError(name types.NamespacedName, cause error) *SecretNotFound {
	return &SecretNotFound{
		NamespacedName: name,
		cause:          cause,
	}
}

var _ error = &SecretNotFound{}
var _ causer = &SecretNotFound{}

func (e *SecretNotFound) Error() string {
	return fmt.Sprintf("%s does not exist (%s)", e.NamespacedName, e.cause)
}

func (e *SecretNotFound) Is(err error) bool {
	var typedErr *SecretNotFound
	if errors.As(err, &typedErr) {
		return e.NamespacedName == typedErr.NamespacedName
	}
	return false
}

func (e *SecretNotFound) Cause() error {
	return e.cause
}

func (e *SecretNotFound) Format(s fmt.State, verb rune) {
	format(e, s, verb)
}

// ConfigMapNotFound error is used when configmap or its expected keys are not found
type ConfigMapNotFound struct {
	NamespacedName types.NamespacedName
	cause          error
}

func NewConfigMapNotFoundError(name types.NamespacedName, cause error) *ConfigMapNotFound {
	return &ConfigMapNotFound{
		NamespacedName: name,
		cause:          cause,
	}
}

var _ error = &ConfigMapNotFound{}
var _ causer = &ConfigMapNotFound{}

func (e *ConfigMapNotFound) Error() string {
	return fmt.Sprintf("%s does not exist (%s)", e.NamespacedName, e.cause)
}

func (e *ConfigMapNotFound) Is(err error) bool {
	var typedErr *ConfigMapNotFound
	if errors.As(err, &typedErr) {
		return e.NamespacedName == typedErr.NamespacedName
	}
	return false
}

func (e *ConfigMapNotFound) Cause() error {
	return e.cause
}

func (e *ConfigMapNotFound) Format(s fmt.State, verb rune) {
	format(e, s, verb)
}

// SubscriptionMismatch error is used when a child resource and parent resource subscription don't match
type SubscriptionMismatch struct {
	ExpectedSubscription string
	ActualSubscription   string
	inner                error
}

func NewSubscriptionMismatchError(expectedSub string, actualSub string) *SubscriptionMismatch {
	err := errors.Errorf(
		"resource subscription %q does not match parent subscription %q",
		actualSub,
		expectedSub)

	return &SubscriptionMismatch{
		ExpectedSubscription: expectedSub,
		ActualSubscription:   actualSub,
		inner:                err,
	}
}

var _ error = &SubscriptionMismatch{}
var _ causer = &SubscriptionMismatch{}

func (e *SubscriptionMismatch) Error() string {
	return e.inner.Error()
}

func (e *SubscriptionMismatch) Is(err error) bool {
	var typedErr *SubscriptionMismatch
	if errors.As(err, &typedErr) {
		return true
	}
	return false
}

func (e *SubscriptionMismatch) Cause() error {
	return e.inner
}

func (e *SubscriptionMismatch) Format(s fmt.State, verb rune) {
	format(e, s, verb)
}

// This was adapted from the function in errors
func format(e causer, s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			n, _ := fmt.Fprintf(s, "%s", e.Cause())
			if n > 0 {
				_, _ = fmt.Fprintf(s, "\n")
			}
			_, _ = io.WriteString(s, e.Error())
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}
