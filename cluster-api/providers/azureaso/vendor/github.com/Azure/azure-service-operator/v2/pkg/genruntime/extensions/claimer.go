/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package extensions

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// Claimer can be implemented to customize how the reconciler claims a resource.
// Claiming establishes the link between a Kubernetes resource and its Azure resource by setting the ARM ID.
// Most resources use the default claiming logic. Implement this extension when:
// - The resource's ARM ID doesn't follow standard construction patterns
// - Custom validation is required before claiming
// - Additional operations must be performed during the claim process
type Claimer interface {
	// Claim claims the resource by establishing its Azure Resource Manager ID.
	// ctx is the current operation context.
	// log is a logger for the current operation.
	// obj is the Kubernetes resource being claimed.
	// next is the default claim implementation - call this to use standard claiming logic.
	Claim(ctx context.Context, log logr.Logger, obj genruntime.ARMOwnedMetaObject, next ClaimFunc) error
}

// ClaimFunc is the signature of a function that can be used to create a default Claimer
type ClaimFunc = func(ctx context.Context, log logr.Logger, obj genruntime.ARMOwnedMetaObject) error

// CreateClaimer creates a ClaimFunc. If the resource in question has not implemented the Claimer interface
// the provided default ClaimFunc is run by default.
func CreateClaimer(
	host genruntime.ResourceExtension,
	next ClaimFunc,
) ClaimFunc {
	impl, ok := host.(Claimer)
	if !ok {
		return next
	}

	return func(ctx context.Context, log logr.Logger, obj genruntime.ARMOwnedMetaObject) error {
		log.V(Status).Info("Running customized claim")
		return impl.Claim(ctx, log, obj, next)
	}
}
