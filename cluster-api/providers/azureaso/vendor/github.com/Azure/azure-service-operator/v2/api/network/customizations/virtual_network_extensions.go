// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package customizations

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	network "github.com/Azure/azure-service-operator/v2/api/network/v1api20201101storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

// Attention: A lot of code in this file is very similar to the logic in network_security_group_extension.go, load_balancer_extensions.go and route_table_extensions.go.
// The two should be kept in sync as much as possible.
// NOTE: This wouldn't work without adding indexes in 'getGeneratedStorageTypes' method in controller_resources.goould be kept in sync as much as possible.

var _ extensions.ARMResourceModifier = &VirtualNetworkExtension{}

func (extension *VirtualNetworkExtension) ModifyARMResource(
	ctx context.Context,
	armClient *genericarmclient.GenericClient,
	armObj genruntime.ARMResource,
	obj genruntime.ARMMetaObject,
	kubeClient kubeclient.Client,
	resolver *resolver.Resolver,
	log logr.Logger,
) (genruntime.ARMResource, error) {
	typedObj, ok := obj.(*network.VirtualNetwork)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *network.VirtualNetwork", obj)
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

	azureSubnets, err := getRawChildCollection(raw, "subnets")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get subnets")
	}

	log.V(Info).Info("Found subnets to include on VNET", "count", len(azureSubnets), "names", genruntime.RawNames(azureSubnets))

	err = setChildCollection(armObj.Spec(), azureSubnets, "Subnets")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set subnets")
	}

	return armObj, nil
}

func getChildCollectionField(parent any, childFieldName string) (ret reflect.Value, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.Errorf("caught panic: %s", x)
		}
	}()

	// Here be dragons
	parentValue := reflect.ValueOf(parent)
	parentValue = reflect.Indirect(parentValue)
	if !parentValue.IsValid() {
		return reflect.Value{}, errors.Errorf("cannot assign to nil parent")
	}

	propertiesField := parentValue.FieldByName("Properties")
	if !propertiesField.IsValid() {
		return reflect.Value{}, errors.Errorf("couldn't find properties field")
	}

	propertiesValue := reflect.Indirect(propertiesField)
	if !propertiesValue.IsValid() {
		// If the properties field is nil, we must construct an entirely new properties and assign it here
		temp := reflect.New(propertiesField.Type().Elem())
		propertiesField.Set(temp)
		propertiesValue = reflect.Indirect(temp)
	}

	childField := propertiesValue.FieldByName(childFieldName)
	if !childField.IsValid() {
		return reflect.Value{}, errors.Errorf("couldn't find %q field", childFieldName)
	}

	if childField.Type().Kind() != reflect.Slice {
		return reflect.Value{}, errors.Errorf("%q field was not of kind Slice", childFieldName)
	}

	return childField, nil
}

func getRawChildCollection(parent map[string]any, childJSONName string) ([]any, error) {
	props, ok := parent["properties"]
	if !ok {
		return nil, errors.Errorf("couldn't find properties field")
	}

	propsMap, ok := props.(map[string]any)
	if !ok {
		return nil, errors.Errorf("properties field wasn't a map")
	}

	childField, ok := propsMap[childJSONName]
	if !ok {
		return nil, errors.Errorf("couldn't find %q field", childJSONName)
	}

	childSlice, ok := childField.([]any)
	if !ok {
		return nil, errors.Errorf("%q field wasn't a slice", childJSONName)
	}

	return childSlice, nil
}

func setChildCollection(parent genruntime.ARMResourceSpec, childCollectionFromAzure []any, childFieldName string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.Errorf("caught panic: %s", x)
		}
	}()

	childField, err := getChildCollectionField(parent, childFieldName)
	if err != nil {
		return err
	}

	elemType := childField.Type().Elem()
	childSlice := reflect.MakeSlice(childField.Type(), 0, 0)

	for _, child := range childCollectionFromAzure {
		embeddedResource := reflect.New(elemType)
		err = fuzzySetResource(child, embeddedResource)
		if err != nil {
			return err
		}

		childSlice = reflect.Append(childSlice, reflect.Indirect(embeddedResource))
	}

	childField.Set(childSlice)

	return nil
}

func fuzzySetResource(resource any, embeddedResource reflect.Value) error {
	resourceJSON, err := json.Marshal(resource)
	if err != nil {
		return errors.Wrap(err, "failed to marshal resource JSON")
	}

	err = json.Unmarshal(resourceJSON, embeddedResource.Interface())
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal resource JSON")
	}

	// TODO: Can't do a trivial fuzzyEqualityComparison here because we don't know which fields are readonly
	// TODO: and which are not. This results in mismatches like dropping etag and other fields.

	return nil
}
