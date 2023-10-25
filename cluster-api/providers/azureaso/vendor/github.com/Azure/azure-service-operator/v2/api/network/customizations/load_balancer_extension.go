// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package customizations

import (
	"context"
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

// Attention: A lot of code in this file is very similar to the logic in network_security_group_extension.go, route_table_extensions.go and virtual_network_extensions.go.
// The two should be kept in sync as much as possible.
// NOTE: This wouldn't work without adding indexes in 'getGeneratedStorageTypes' method in controller_resources.go

var _ extensions.ARMResourceModifier = &LoadBalancerExtension{}

func (extension *LoadBalancerExtension) ModifyARMResource(
	ctx context.Context,
	armClient *genericarmclient.GenericClient,
	armObj genruntime.ARMResource,
	obj genruntime.ARMMetaObject,
	kubeClient kubeclient.Client,
	resolver *resolver.Resolver,
	log logr.Logger,
) (genruntime.ARMResource, error) {

	typedObj, ok := obj.(*network.LoadBalancer)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *network.LoadBalancer", obj)
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

	azureInboundNatRules, err := getRawChildCollection(raw, "inboundNatRules")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get inboundNatRules")
	}

	log.V(Info).Info("Found InboundNatRules to include on LoadBalancer", "count", len(azureInboundNatRules), "names", genruntime.RawNames(azureInboundNatRules))

	err = setInboundNatRules(armObj.Spec(), azureInboundNatRules)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set inboundNatRules")
	}

	return armObj, nil
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

func setInboundNatRules(lb genruntime.ARMResourceSpec, azureInboundNatRules []any) (err error) {
	defer func() {
		if x := recover(); x != nil {
			err = errors.Errorf("caught panic: %s", x)
		}
	}()

	inboundNatRulesField, err := getChildCollectionField(lb, "InboundNatRules")
	if err != nil {
		return err
	}

	elemType := inboundNatRulesField.Type().Elem()
	inboundNatRulesSlice := reflect.MakeSlice(inboundNatRulesField.Type(), 0, 0)

	for _, inboundNatRule := range azureInboundNatRules {
		embeddedInboundNatRules := reflect.New(elemType)
		err = fuzzySetResource(inboundNatRule, embeddedInboundNatRules)
		if err != nil {
			return err
		}

		inboundNatRulesSlice = reflect.Append(inboundNatRulesSlice, reflect.Indirect(embeddedInboundNatRules))
	}

	inboundNatRulesSlice, err = mergeNatRules(inboundNatRulesField, inboundNatRulesSlice)
	if err != nil {
		return errors.Wrapf(err, "failed to merge NAT rules")
	}

	inboundNatRulesField.Set(inboundNatRulesSlice)

	return nil
}

func mergeNatRules(inboundNatRulesField reflect.Value, azureInboundNatRules reflect.Value) (reflect.Value, error) {
	if inboundNatRulesField.Len() == 0 {
		return azureInboundNatRules, nil
	}

	resultSlice := reflect.MakeSlice(inboundNatRulesField.Type(), 0, 0)

	// The result slice should take every item from inboundNatRulesField, and merge in anything from azureInboundNatRules
	// with a different name
	for i := 0; i < inboundNatRulesField.Len(); i++ {
		inboundNatRule := inboundNatRulesField.Index(i)
		resultSlice = reflect.Append(resultSlice, inboundNatRule)
	}

	for i := 0; i < azureInboundNatRules.Len(); i++ {
		inboundNatRule := azureInboundNatRules.Index(i)
		newRuleName, err := getNameField(inboundNatRule)
		if err != nil {
			return reflect.Value{}, errors.Wrapf(err, "failed to get name for new rule")
		}
		foundExistingRule := false

		for j := 0; j < inboundNatRulesField.Len(); j++ {
			existingInboundNatRule := inboundNatRulesField.Index(j)
			var existingName reflect.Value
			existingName, err = getNameField(existingInboundNatRule)
			if err != nil {
				return reflect.Value{}, errors.Wrapf(err, "failed to get name for existing rule")
			}

			if existingName.String() == newRuleName.String() {
				foundExistingRule = true
				break
			}
		}

		if !foundExistingRule {
			resultSlice = reflect.Append(resultSlice, inboundNatRule)
		}
	}

	return resultSlice, nil
}
