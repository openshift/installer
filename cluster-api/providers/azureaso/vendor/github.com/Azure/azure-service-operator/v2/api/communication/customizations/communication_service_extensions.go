/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/communication/armcommunication"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	communication "github.com/Azure/azure-service-operator/v2/api/communication/v20230401/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	primaryKey                = "primaryKey"
	secondaryKey              = "secondaryKey"
	primaryConnectionString   = "primaryConnectionString"
	secondaryConnectionString = "secondaryConnectionString"
)

var _ genruntime.KubernetesSecretExporter = &CommunicationServiceExtension{}

// ExportKubernetesSecrets implements genruntime.KubernetesSecretExporter.
func (ext *CommunicationServiceExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// Make sure we're working with the current hub version of the resource.
	// This will need to be updated if the hub version changes.
	typedObj, ok := obj.(*communication.CommunicationService)
	if !ok {
		return nil, eris.Errorf(
			"cannot run on unknown resource type %T, expected *communication.CommunicationService", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not.
	var _ conversion.Hub = typedObj

	primarySecrets := communicationServiceSecretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)

	if len(requestedSecrets) == 0 {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
	// connection each time through.
	client, err := armcommunication.NewServiceClient(id.SubscriptionID, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create ARM communication service client")
	}

	res, err := client.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
	if err != nil {
		return nil, eris.Wrapf(err, "failed to list keys")
	}

	secretSlice, err := communicationServiceSecretsToWrite(typedObj, res.ServiceKeys)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := makeCommunicationServiceResolvedSecretsMap(res.ServiceKeys)

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func communicationServiceSecretsSpecified(obj *communication.CommunicationService) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	specSecrets := obj.Spec.OperatorSpec.Secrets
	result := make(set.Set[string])
	if specSecrets.PrimaryKey != nil {
		result.Add(primaryKey)
	}
	if specSecrets.SecondaryKey != nil {
		result.Add(secondaryKey)
	}
	if specSecrets.PrimaryConnectionString != nil {
		result.Add(primaryConnectionString)
	}
	if specSecrets.SecondaryConnectionString != nil {
		result.Add(secondaryConnectionString)
	}

	return result
}

func communicationServiceSecretsToWrite(obj *communication.CommunicationService, keys armcommunication.ServiceKeys) ([]*v1.Secret, error) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil, nil
	}

	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(keys.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(keys.SecondaryKey))
	collector.AddValue(operatorSpecSecrets.PrimaryConnectionString, to.Value(keys.PrimaryConnectionString))
	collector.AddValue(operatorSpecSecrets.SecondaryConnectionString, to.Value(keys.SecondaryConnectionString))

	return collector.Values()
}

func makeCommunicationServiceResolvedSecretsMap(keys armcommunication.ServiceKeys) map[string]string {
	resolvedSecrets := map[string]string{}
	if to.Value(keys.PrimaryKey) != "" {
		resolvedSecrets[primaryKey] = to.Value(keys.PrimaryKey)
	}
	if to.Value(keys.SecondaryKey) != "" {
		resolvedSecrets[secondaryKey] = to.Value(keys.SecondaryKey)
	}
	if to.Value(keys.PrimaryConnectionString) != "" {
		resolvedSecrets[primaryConnectionString] = to.Value(keys.PrimaryConnectionString)
	}
	if to.Value(keys.SecondaryConnectionString) != "" {
		resolvedSecrets[secondaryConnectionString] = to.Value(keys.SecondaryConnectionString)
	}

	return resolvedSecrets
}
