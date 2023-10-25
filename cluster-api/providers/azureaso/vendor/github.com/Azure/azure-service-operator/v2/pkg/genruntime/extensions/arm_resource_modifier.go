/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package extensions

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// ARMResourceModifier provides a hook allowing resources to modify the payload that will be sent to ARM just before it is sent.
type ARMResourceModifier interface {
	// ModifyARMResource takes a genruntime.ARMResource and returns an updated genruntime.ARMResource. The updated resource
	// is then serialized and sent to ARM in the body of a PUT request.
	ModifyARMResource(
		ctx context.Context,
		armClient *genericarmclient.GenericClient,
		armObj genruntime.ARMResource,
		obj genruntime.ARMMetaObject,
		kubeClient kubeclient.Client,
		resolver *resolver.Resolver,
		log logr.Logger,
	) (genruntime.ARMResource, error)
}

type ARMResourceModifierFunc = func(ctx context.Context, obj genruntime.ARMMetaObject, armObj genruntime.ARMResource) (genruntime.ARMResource, error)

// CreateARMResourceModifier returns a function that performs per-resource modifications. If the ARMResourceModifier extension has
// not been implemented for the resource in question, the default behavior is to return the provided genruntime.ARMResource unmodified.
func CreateARMResourceModifier(
	host genruntime.ResourceExtension,
	armClient *genericarmclient.GenericClient,
	kubeClient kubeclient.Client,
	resolver *resolver.Resolver,
	log logr.Logger) ARMResourceModifierFunc {

	impl, ok := host.(ARMResourceModifier)
	if !ok {
		return func(ctx context.Context, obj genruntime.ARMMetaObject, armObj genruntime.ARMResource) (genruntime.ARMResource, error) {
			return armObj, nil
		}
	}

	return func(ctx context.Context, obj genruntime.ARMMetaObject, armObj genruntime.ARMResource) (genruntime.ARMResource, error) {
		log.V(Info).Info("Modifying ARM payload")

		armResource, err := impl.ModifyARMResource(ctx, armClient, armObj, obj, kubeClient, resolver, log)
		if err != nil {
			log.V(Status).Info("Failed to modify ARM resource payload", "error", err.Error())
			return nil, errors.Wrap(err, "failed to modify ARM payload")
		}

		log.V(Info).Info("Successfully modified ARM payload")
		return armResource, nil
	}
}
