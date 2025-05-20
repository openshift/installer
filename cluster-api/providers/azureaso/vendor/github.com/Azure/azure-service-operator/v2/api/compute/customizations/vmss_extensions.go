/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

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

	compute "github.com/Azure/azure-service-operator/v2/api/compute/v1api20220301/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ErrorClassifier = &VirtualMachineScaleSetExtension{}

var (
	rawChildCollectionPath = []string{"properties", "virtualMachineProfile", "extensionProfile", "extensions"}
	childCollectionPathARM = []string{"Properties", "VirtualMachineProfile", "ExtensionProfile", "Extensions"}
)

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
func (e *VirtualMachineScaleSetExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// It's weird to retry on OperationNotAllowed as it certainly sounds like a fatal error, but
	// it primarily happens for quota errors on VM/VMSS, which we do want to retry on as the quota may free
	// up at some point in the future.
	if details.Code == "OperationNotAllowed" {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

// Attention: A lot of code in this file is very similar to the logic in network/network_security_group_extension.go, network/route_table_extensions.go, network/virtual_network_extensions.go and network/load_balancer_extension.go.
// The two should be kept in sync as much as possible.

func (e *VirtualMachineScaleSetExtension) ModifyARMResource(
	ctx context.Context,
	armClient *genericarmclient.GenericClient,
	armObj genruntime.ARMResource,
	obj genruntime.ARMMetaObject,
	kubeClient kubeclient.Client,
	resolver *resolver.Resolver,
	log logr.Logger,
) (genruntime.ARMResource, error) {
	typedObj, ok := obj.(*compute.VirtualMachineScaleSet)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *compute.VirtualMachineScaleSet", obj)
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

	azureExtensions, err := getRawChildCollection(raw, rawChildCollectionPath...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get VMSS Extensions")
	}

	// If the child collection is not defined, We return the arm object as is here.
	if azureExtensions == nil {
		log.V(Info).Info("Found no Extensions to include on VMSS")
		return armObj, nil
	}

	log.V(Info).Info("Found Extensions to include on VMSS", "count", len(azureExtensions), "names", genruntime.RawNames(azureExtensions))

	err = setChildCollection(armObj.Spec(), azureExtensions, childCollectionPathARM...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set VMSS Extensions")
	}

	return armObj, nil
}

func getExactParentRecursively(parentValue reflect.Value, childFieldName []string, i int) (ret reflect.Value, err error) {
	// Until we reach the exact parent
	if i == len(childFieldName)-1 {
		return parentValue, nil
	}

	fieldName := childFieldName[i]
	field := parentValue.FieldByName(fieldName)
	if !field.IsValid() {
		return reflect.Value{}, errors.Errorf("couldn't find %s field", fieldName)
	}

	propertiesValue := reflect.Indirect(field)
	if !propertiesValue.IsValid() {
		// If the properties field is nil, we must construct an entirely new properties and assign it here
		temp := reflect.New(field.Type().Elem())
		field.Set(temp)
		propertiesValue = reflect.Indirect(temp)
	}

	return getExactParentRecursively(propertiesValue, childFieldName, i+1)
}

func getChildCollectionField(parent any, fieldPath []string) (ret reflect.Value, err error) {
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

	exactParent, err := getExactParentRecursively(parentValue, fieldPath, 0)
	if err != nil {
		return exactParent, err
	}

	childFieldName := fieldPath[len(fieldPath)-1]
	childField := exactParent.FieldByName(childFieldName)
	if !childField.IsValid() {
		return reflect.Value{}, errors.Errorf("couldn't find %q field", fieldPath)
	}

	if childField.Type().Kind() != reflect.Slice {
		return reflect.Value{}, errors.Errorf("%q field was not of kind Slice", fieldPath)
	}

	return childField, nil
}

func getRawExactParentRecursively(parent map[string]any, fieldSlice []string, i int) (map[string]any, error) {
	// Until we reach the exact parent
	if i == len(fieldSlice)-1 {
		return parent, nil
	}

	prop := fieldSlice[i]

	props, ok := parent[prop]
	if !ok {
		// We don't want to return an error here since field might be empty or not defined initially.
		return nil, nil
	}

	propsMap, ok := props.(map[string]any)
	if !ok {
		return nil, errors.Errorf("%s field wasn't a map", prop)
	}

	return getRawExactParentRecursively(propsMap, fieldSlice, i+1)
}

func getRawChildCollection(parent map[string]any, fieldSlice ...string) ([]any, error) {
	exactParent, err := getRawExactParentRecursively(parent, fieldSlice, 0)
	if err != nil {
		return nil, err
	}

	if exactParent == nil {
		return nil, nil
	}

	childFieldName := fieldSlice[len(fieldSlice)-1]
	childField, ok := exactParent[childFieldName]
	if !ok {
		return nil, errors.Errorf("couldn't find %q field", fieldSlice)
	}

	childSlice, ok := childField.([]any)
	if !ok {
		return nil, errors.Errorf("%q field wasn't a slice", fieldSlice)
	}

	return childSlice, nil
}

func setChildCollection(parent genruntime.ARMResourceSpec, childCollectionFromAzure []any, childFieldPath ...string) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.Errorf("caught panic: %s", x)
		}
	}()

	childField, err := getChildCollectionField(parent, childFieldPath)
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

	childSlice, err = mergeExtensions(childField, childSlice)
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

func mergeExtensions(extensionField reflect.Value, azureExtensionsSlice reflect.Value) (reflect.Value, error) {
	if extensionField.Len() == 0 {
		return azureExtensionsSlice, nil
	}

	resultSlice := reflect.MakeSlice(extensionField.Type(), 0, 0)

	// The result slice should take every item from extensionField, and merge in anything from azureExtensionsSlice
	// with a different name
	for i := 0; i < extensionField.Len(); i++ {
		extension := extensionField.Index(i)
		resultSlice = reflect.Append(resultSlice, extension)
	}

	for i := 0; i < azureExtensionsSlice.Len(); i++ {
		azureExtension := azureExtensionsSlice.Index(i)
		newExtensionName, err := getNameField(azureExtension)
		if err != nil {
			return reflect.Value{}, errors.Wrapf(err, "failed to get name for new extension")
		}
		foundExistingExtension := false

		for j := 0; j < extensionField.Len(); j++ {
			existingExtension := extensionField.Index(j)
			var existingName reflect.Value
			existingName, err = getNameField(existingExtension)
			if err != nil {
				return reflect.Value{}, errors.Wrapf(err, "failed to get name for existing extension")
			}

			if existingName.String() == newExtensionName.String() {
				foundExistingExtension = true
				break
			}
		}

		if !foundExistingExtension {
			resultSlice = reflect.Append(resultSlice, azureExtension)
		}
	}

	return resultSlice, nil
}

func getNameField(natValue reflect.Value) (ret reflect.Value, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.Errorf("caught panic: %s", x)
		}
	}()

	nameField := natValue.FieldByName("Name")
	if !nameField.IsValid() {
		return nameField, errors.Errorf("couldn't find name field")
	}

	nameField = reflect.Indirect(nameField)

	return nameField, nil
}
