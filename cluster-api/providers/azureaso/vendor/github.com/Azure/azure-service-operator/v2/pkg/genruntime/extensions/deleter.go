/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package extensions

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// Deleter can be implemented to customize how the reconciler deletes resources from Azure.
// This extension is invoked when a resource has a deletion timestamp, before the standard ARM DELETE operation.
// Implement this extension when:
// - Pre-deletion operations are required (e.g., canceling subscriptions, disabling features)
// - Multiple API calls are needed for complete deletion
// - Custom error handling is needed during deletion
// - Resources should be preserved in Azure under certain conditions
type Deleter interface {
	// Delete deletes the resource from Azure, with the ability to perform custom logic before, during, or after deletion.
	// ctx is the current operation context.
	// log is a logger for the current operation.
	// resolver helps resolve resource references.
	// armClient allows making ARM API calls.
	// obj is the Kubernetes resource being deleted.
	// next is the default deletion implementation - call this to perform standard ARM DELETE.
	// Returns a reconciliation result (e.g., requeue timing) and an error if deletion fails.
	Delete(
		ctx context.Context,
		log logr.Logger,
		resolver *resolver.Resolver,
		armClient *genericarmclient.GenericClient,
		obj genruntime.ARMMetaObject,
		next DeleteFunc) (ctrl.Result, error)
}

// DeleteFunc is the signature of a function that can be used to create a default Deleter
type DeleteFunc = func(
	ctx context.Context,
	log logr.Logger,
	resolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	obj genruntime.ARMMetaObject) (ctrl.Result, error)

// CreateDeleter creates a DeleteFunc. If the resource in question has not implemented the Deleter interface
// the provided default DeleteFunc is run by default.
func CreateDeleter(
	host genruntime.ResourceExtension,
	next DeleteFunc,
) DeleteFunc {
	impl, ok := host.(Deleter)
	if !ok {
		return next
	}

	return func(ctx context.Context, log logr.Logger, resolver *resolver.Resolver, armClient *genericarmclient.GenericClient, obj genruntime.ARMMetaObject) (ctrl.Result, error) {
		log.V(Status).Info("Running customized deletion")
		return impl.Delete(ctx, log, resolver, armClient, obj, next)
	}
}
