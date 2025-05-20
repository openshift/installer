/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	signalr "github.com/Azure/azure-service-operator/v2/api/signalrservice/v1api20211001/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
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

var _ genruntime.KubernetesSecretExporter = &SignalRExtension{}

func (ext *SignalRExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// Make sure we're working with the current hub version of the resource
	// This will need to be updated if the hub version changes
	typedObj, ok := obj.(*signalr.SignalR)
	if !ok {
		return nil, errors.Errorf(
			"cannot run on unknown resource type %T, expected *signalr.SignalR", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	primarySecrets := secretsSpecified(typedObj)
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
	// connection each time through
	clientFactory, err := armsignalr.NewClientFactory(id.SubscriptionID, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create ARM signalR client factory")
	}

	res, err := clientFactory.NewClient().ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list keys")
	}

	secretSlice, err := secretsToWrite(typedObj, res.Keys)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := makeResolvedSecretsMap(res.Keys)

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *signalr.SignalR) set.Set[string] {
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

func secretsToWrite(obj *signalr.SignalR, accessKeys armsignalr.Keys) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(accessKeys.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(accessKeys.SecondaryKey))
	collector.AddValue(operatorSpecSecrets.PrimaryConnectionString, to.Value(accessKeys.PrimaryConnectionString))
	collector.AddValue(operatorSpecSecrets.SecondaryConnectionString, to.Value(accessKeys.SecondaryConnectionString))

	return collector.Values()
}

func makeResolvedSecretsMap(accessKeys armsignalr.Keys) map[string]string {
	resolvedSecrets := map[string]string{}
	if to.Value(accessKeys.PrimaryKey) != "" {
		resolvedSecrets[primaryKey] = to.Value(accessKeys.PrimaryKey)
	}
	if to.Value(accessKeys.SecondaryKey) != "" {
		resolvedSecrets[secondaryKey] = to.Value(accessKeys.SecondaryKey)
	}
	if to.Value(accessKeys.PrimaryConnectionString) != "" {
		resolvedSecrets[primaryConnectionString] = to.Value(accessKeys.PrimaryConnectionString)
	}
	if to.Value(accessKeys.SecondaryConnectionString) != "" {
		resolvedSecrets[secondaryConnectionString] = to.Value(accessKeys.SecondaryConnectionString)
	}

	return resolvedSecrets
}
