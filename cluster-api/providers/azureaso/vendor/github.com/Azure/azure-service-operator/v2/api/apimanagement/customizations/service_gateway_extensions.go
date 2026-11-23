/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v2"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	apimanagement "github.com/Azure/azure-service-operator/v2/api/apimanagement/v20240501/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	gatewayPrimaryKey   = "primaryKey"
	gatewaySecondaryKey = "secondaryKey"
)

var _ genruntime.KubernetesSecretExporter = &ServiceGatewayExtension{}

func (ext *ServiceGatewayExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*apimanagement.ServiceGateway)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *apimanagement.ServiceGateway", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	primarySecrets := gatewaySecretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)
	if len(requestedSecrets) == 0 {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	if id.Parent == nil {
		return nil, eris.Errorf("APIM gateway had no parent ID: %s", id.String())
	}
	parentName := id.Parent.Name

	// Only bother calling ListKeys if there are secrets to retrieve
	var s armapimanagement.GatewayKeysContract
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var gatewayClient *armapimanagement.GatewayClient
		gatewayClient, err = armapimanagement.NewGatewayClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create new GatewayClient")
		}

		var resp armapimanagement.GatewayClientListKeysResponse
		resp, err = gatewayClient.ListKeys(ctx, id.ResourceGroupName, parentName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, eris.Wrapf(err, "failed listing keys")
		}

		s = resp.GatewayKeysContract
	}

	resolvedSecrets := map[string]string{}
	if to.Value(s.Primary) != "" {
		resolvedSecrets[gatewayPrimaryKey] = to.Value(s.Primary)
	}
	if to.Value(s.Secondary) != "" {
		resolvedSecrets[gatewaySecondaryKey] = to.Value(s.Secondary)
	}
	secretSlice, err := gatewaySecretsToWrite(typedObj, s)
	if err != nil {
		return nil, err
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func gatewaySecretsSpecified(obj *apimanagement.ServiceGateway) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	s := obj.Spec.OperatorSpec.Secrets
	result := make(set.Set[string])
	if s.PrimaryKey != nil {
		result.Add(gatewayPrimaryKey)
	}
	if s.SecondaryKey != nil {
		result.Add(gatewaySecondaryKey)
	}

	return result
}

func gatewaySecretsToWrite(obj *apimanagement.ServiceGateway, s armapimanagement.GatewayKeysContract) ([]*v1.Secret, error) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil, nil
	}

	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(s.Primary))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(s.Secondary))

	return collector.Values()
}
