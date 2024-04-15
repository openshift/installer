// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	storage "github.com/Azure/azure-service-operator/v2/api/subscription/v1api20211001/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.Deleter = &AliasExtension{}

func (extension *AliasExtension) Delete(
	ctx context.Context,
	log logr.Logger,
	resolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	obj genruntime.ARMMetaObject,
	next extensions.DeleteFunc) (ctrl.Result, error) {

	// First cancel the subscription, then delete the alias
	typedObj, ok := obj.(*storage.Alias)
	if !ok {
		return ctrl.Result{}, errors.Errorf("cannot run on unknown resource type %T, expected *subscription.Alias", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated to match
	var _ conversion.Hub = typedObj

	// Get the subscription ID
	subscriptionID, ok := getSubscriptionID(typedObj)
	if !ok {
		// SubscriptionID isn't populated, allow deletion to proceed
		return next(ctx, log, resolver, armClient, obj)
	}

	// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
	// connection each time through
	subscriptionClient, err := armsubscription.NewClient(armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to create new workspaceClient")
	}

	// Don't need to do anything with the response here so just ignore it.
	// Note that this operation is not asynchronous. It completes synchronously.
	_, err = subscriptionClient.Cancel(ctx, subscriptionID, nil)
	if err != nil {
		// TODO: May need to set condition error here
		return ctrl.Result{}, errors.Wrapf(err, "failed to cancel subscription %q", subscriptionID)
	}

	return next(ctx, log, resolver, armClient, obj)
}

var _ extensions.SuccessfulCreationHandler = &AliasExtension{}

func (extension *AliasExtension) Success(obj genruntime.ARMMetaObject) error {
	typedObj, ok := obj.(*storage.Alias)
	if !ok {
		return errors.Errorf("cannot run on unknown resource type %T, expected *subscription.Alias", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated to match
	var _ conversion.Hub = typedObj

	// Get the subscription ID
	subscriptionID, ok := getSubscriptionID(typedObj)
	if !ok {
		// SubscriptionID isn't populated. That's a problem
		return errors.Errorf("SubscriptionID field not populated")
	}

	genruntime.SetChildResourceIDOverride(typedObj, genericarmclient.MakeSubscriptionID(subscriptionID))

	return nil
}

func getSubscriptionID(typedObj *storage.Alias) (string, bool) {
	// Get the subscription ID
	if typedObj.Status.Properties == nil || typedObj.Status.Properties.SubscriptionId == nil || *typedObj.Status.Properties.SubscriptionId == "" {
		return "", false
	}

	return *typedObj.Status.Properties.SubscriptionId, true
}
