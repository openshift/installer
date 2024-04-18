// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	storage "github.com/Azure/azure-service-operator/v2/api/apimanagement/v1api20220801/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.Deleter = &ProductExtension{}

func (extension *ProductExtension) Delete(
	ctx context.Context,
	log logr.Logger,
	resolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	obj genruntime.ARMMetaObject,
	next extensions.DeleteFunc) (ctrl.Result, error) {

	typedObj, ok := obj.(*storage.Product)
	if !ok {
		return ctrl.Result{}, errors.Errorf("cannot run on unknown resource type %T, expected *apiManagement.Product", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated to match
	var _ conversion.Hub = typedObj

	productName := typedObj.GetName()
	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to get the ARM ResourceId for %s", productName)
	}

	if id.Parent == nil {
		return ctrl.Result{}, errors.Wrapf(err, ". APIM Product had no parent ID: %s", id.String())
	}
	parentName := id.Parent.Name

	// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
	// connection each time through
	clientFactory, err := armapimanagement.NewClientFactory(id.SubscriptionID, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to create new apimClient")
	}

	// This is a synchronous operation
	_, err = clientFactory.NewProductClient().Delete(
		ctx,
		id.ResourceGroupName,
		parentName,
		productName,
		"*",
		&armapimanagement.ProductClientDeleteOptions{
			DeleteSubscriptions: to.Ptr(true),
		})
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to delete product %q", productName)
	}

	return next(ctx, log, resolver, armClient, obj)
}
