// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package customizations

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	network "github.com/Azure/azure-service-operator/v2/api/network/v1api20240301/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

// Attention: A lot of code in this file is very similar to the logic in load_balancer_extension.go, route_table_extensions.go, virtual_network_extensions.go and compute/vmss_extensions.go.
// The two should be kept in sync as much as possible.

var _ extensions.ARMResourceModifier = &NetworkSecurityGroupExtension{}

func (extension *NetworkSecurityGroupExtension) ModifyARMResource(
	ctx context.Context,
	armClient *genericarmclient.GenericClient,
	armObj genruntime.ARMResource,
	obj genruntime.ARMMetaObject,
	kubeClient kubeclient.Client,
	resolver *resolver.Resolver,
	log logr.Logger,
) (genruntime.ARMResource, error) {
	typedObj, ok := obj.(*network.NetworkSecurityGroup)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *network.NetworkSecurityGroup", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated to match
	var _ conversion.Hub = typedObj

	resourceID, hasResourceID := genruntime.GetResourceID(obj)
	if !hasResourceID {
		// If we don't have an ARM ID yet, we've not been claimed so just return armObj as is
		return armObj, nil
	}

	apiVersion, err := genruntime.GetAPIVersion(typedObj, kubeClient.Scheme())
	if err != nil {
		return nil, errors.Wrapf(err, "error getting api version for resource %s while getting status", obj.GetName())
	}

	// Get the raw resource
	raw := make(map[string]any)
	_, err = armClient.GetByID(ctx, resourceID, apiVersion, &raw)
	if err != nil {
		// If the error is NotFound, the resource we're trying to Create doesn't exist and so no modification is needed
		var responseError *azcore.ResponseError
		if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound {
			return armObj, nil
		}
		return nil, errors.Wrapf(err, "getting resource with ID: %q", resourceID)
	}

	azureSecurityRules, err := getRawChildCollection(raw, "securityRules")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get SecurityRules")
	}

	log.V(Info).Info("Found security rules to include on NSG", "count", len(azureSecurityRules), "names", genruntime.RawNames(azureSecurityRules))

	err = setChildCollection(armObj.Spec(), azureSecurityRules, "SecurityRules")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set SecurityRules")
	}

	return armObj, nil
}
