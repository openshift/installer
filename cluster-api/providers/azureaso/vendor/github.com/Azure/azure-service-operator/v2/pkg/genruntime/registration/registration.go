/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package registration

import (
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// Index describes an index registration.
// See controller-runtime mgr.GetFieldIndexer().IndexField() for more details.
type Index struct {
	Key  string
	Func func(rawObj client.Object) []string
}

type EventHandlerFactory func(client client.Client, log logr.Logger) handler.EventHandler

// Watch describes a watch registration.
// See controller-runtime builder.Watches() for more details.
type Watch struct {
	Type             client.Object
	MakeEventHandler EventHandlerFactory
}

// StorageType describes a storage type which will be reconciled.
type StorageType struct {
	// Obj is the object whose Kind should be registered as the storage type.
	Obj client.Object
	// Indexes are additional indexes which must be set up in client-go for this type. These indexes will
	// be used by the reconciliation loop for this resource.
	Indexes []Index
	// Watches are additional event sources that trigger the reconciliation process. This is commonly
	// used when multiple resources in Kubernetes combine to work together. For example, when a
	// database takes a Kubernetes secret as an input, that secret must also be watched in addition
	// to the database itself, so that changes to the secret are correctly propagated.
	Watches []Watch
	// Reconciler is the reconciler instance for resources of this type.
	Reconciler genruntime.Reconciler
	// Predicate determines which events trigger reconciliation for this type
	Predicate predicate.Predicate
	// Name is the friendly name of this storage type
	Name string
}

// NewStorageType makes a new storage type for the specified object
func NewStorageType(obj client.Object) *StorageType {
	return &StorageType{
		Obj: obj,
	}
}
