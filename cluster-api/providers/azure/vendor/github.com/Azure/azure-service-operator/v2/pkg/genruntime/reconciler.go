/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Reconciler performs create/delete actions against a particular kind of resource.
type Reconciler interface {
	// CreateOrUpdate performs create or update of the resource. This must be idempotent. In the event the CreateOrUpdate
	// takes a long time, CreateOrUpdate should return quickly but set an annotation or ready condition that can be used on subsequent
	// calls to monitor the ongoing CreateOrUpdate.
	CreateOrUpdate(
		ctx context.Context,
		log logr.Logger,
		eventRecorder record.EventRecorder,
		obj MetaObject) (ctrl.Result, error)

	// Delete performs deletion of the resource. This must be idempotent. Removal of the common finalizer is performed elsewhere.
	// Delete should concern itself with issuing and tracking the resource deletion.
	Delete(
		ctx context.Context,
		log logr.Logger,
		eventRecorder record.EventRecorder,
		obj MetaObject) (ctrl.Result, error)

	// Claim performs resource specific claim actions. This must be idempotent.
	// A standard finalizer is added to all resources, Claim
	// should deal with any resource specific claiming actions (such as setting a resource ID annotation, etc).
	// If Claim returns an error then reconciliation will be retried according to the returned Ready condition until
	// no error is returned. Once Claim succeeds CreateOrUpdate is called.
	Claim(
		ctx context.Context,
		log logr.Logger,
		eventRecorder record.EventRecorder,
		obj MetaObject) error

	// UpdateStatus fetches the resource's status but performs no other actions. This is primarily called if the
	// reconcile-policy annotation was set in such a way that it blocks CreateOrUpdate
	UpdateStatus(
		ctx context.Context,
		log logr.Logger,
		eventRecorder record.EventRecorder,
		obj MetaObject) error
}
