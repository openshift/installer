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

package aso

import (
	"context"

	"github.com/pkg/errors"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	reconcilerutil "sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// Service provides operations on Azure resources.
type Service[T deepCopier[T], S Scope] struct {
	Reconciler[T]

	Scope S
	Specs []azure.ASOResourceSpecGetter[T]

	ConditionType                  clusterv1.ConditionType
	PostCreateOrUpdateResourceHook func(ctx context.Context, scope S, result T, err error) error
	PostReconcileHook              func(ctx context.Context, scope S, err error) error
	PostDeleteHook                 func(ctx context.Context, scope S, err error) error

	name string
}

// NewService creates a new Service.
func NewService[T deepCopier[T], S Scope](name string, scope S) *Service[T, S] {
	return &Service[T, S]{
		Reconciler: New[T](scope.GetClient(), scope.ClusterName()),
		Scope:      scope,
		name:       name,
	}
}

// Name returns the service name.
func (s *Service[T, S]) Name() string {
	return s.name
}

// Reconcile idempotently creates or updates the resources.
func (s *Service[T, S]) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "aso.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconcilerutil.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	// We go through the list of Specs to reconcile each one, independently of the result of the previous one.
	// If multiple errors occur, we return the most pressing one.
	// Order of precedence (highest -> lowest) is:
	//   - error that is not an operationNotDoneError (i.e. error creating)
	//   - operationNotDoneError (i.e. creating in progress)
	//   - no error (i.e. created)
	var resultErr error
	for _, spec := range s.Specs {
		result, err := s.CreateOrUpdateResource(ctx, spec, s.Name())
		if s.PostCreateOrUpdateResourceHook != nil {
			err = s.PostCreateOrUpdateResourceHook(ctx, s.Scope, result, err)
		}
		if err != nil && (!azure.IsOperationNotDoneError(err) || resultErr == nil) {
			resultErr = err
		}
	}

	if s.PostReconcileHook != nil {
		resultErr = s.PostReconcileHook(ctx, s.Scope, resultErr)
	}
	s.Scope.UpdatePutStatus(s.ConditionType, s.Name(), resultErr)
	return resultErr
}

// Delete deletes the resources.
func (s *Service[T, S]) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "aso.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconcilerutil.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	// We go through the list of Specs to delete each one, independently of the resultErr of the previous one.
	// If multiple errors occur, we return the most pressing one.
	// Order of precedence (highest -> lowest) is:
	//   - error that is not an operationNotDoneError (i.e. error deleting)
	//   - operationNotDoneError (i.e. deleting in progress)
	//   - no error (i.e. deleted)
	var resultErr error
	for _, spec := range s.Specs {
		err := s.DeleteResource(ctx, spec, s.Name())
		if err != nil && (!azure.IsOperationNotDoneError(err) || resultErr == nil) {
			resultErr = err
		}
	}

	if s.PostDeleteHook != nil {
		resultErr = s.PostDeleteHook(ctx, s.Scope, resultErr)
	}
	s.Scope.UpdateDeleteStatus(s.ConditionType, s.Name(), resultErr)
	return resultErr
}

// Pause implements azure.Pauser.
func (s *Service[T, S]) Pause(ctx context.Context) error {
	var _ azure.Pauser = (*Service[T, S])(nil)

	ctx, _, done := tele.StartSpanWithLogger(ctx, "aso.Service.Pause")
	defer done()

	for _, spec := range s.Specs {
		if err := s.PauseResource(ctx, spec, s.Name()); err != nil {
			ref := spec.ResourceRef()
			return errors.Wrapf(err, "failed to pause ASO resource %s %s/%s", ref.GetObjectKind().GroupVersionKind(), ref.GetNamespace(), ref.GetName())
		}
	}

	return nil
}
