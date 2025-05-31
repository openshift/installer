/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventhub/armeventhub"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-service-operator/v2/api/eventhub/v1api20211101/storage"
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

var _ genruntime.KubernetesSecretExporter = &NamespacesEventhubsAuthorizationRuleExtension{}

func (ext *NamespacesEventhubsAuthorizationRuleExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*storage.NamespacesEventhubsAuthorizationRule)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *eventhub.NamespacesEventhubsAuthorizationRule", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	primarySecrets := namespacesEventHubAuthorizationRuleSecretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)
	if len(requestedSecrets) == 0 {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	// Only bother calling ListKeys if there are secrets to retrieve
	var res armeventhub.EventHubsClientListKeysResponse
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var confClient *armeventhub.EventHubsClient
		confClient, err = armeventhub.NewEventHubsClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new EventHubsClient")
		}

		res, err = confClient.ListKeys(ctx, id.ResourceGroupName, id.Parent.Parent.Name, id.Parent.Name, typedObj.Name, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retreive response")
		}
	}

	secretSlice, err := namespacesEventHubAuthorizationRuleSecretsToWrite(typedObj, res.AccessKeys)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := makeNamespacesEventHubAuthorizationRuleResolvedSecretsMap(res.AccessKeys)

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func namespacesEventHubAuthorizationRuleSecretsSpecified(obj *storage.NamespacesEventhubsAuthorizationRule) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	secrets := obj.Spec.OperatorSpec.Secrets

	result := make(set.Set[string])
	if secrets.PrimaryKey != nil {
		result.Add(primaryKey)
	}
	if secrets.SecondaryKey != nil {
		result.Add(secondaryKey)
	}
	if secrets.PrimaryConnectionString != nil {
		result.Add(primaryConnectionString)
	}
	if secrets.SecondaryConnectionString != nil {
		result.Add(secondaryConnectionString)
	}

	return result
}

func namespacesEventHubAuthorizationRuleSecretsToWrite(obj *storage.NamespacesEventhubsAuthorizationRule, keys armeventhub.AccessKeys) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)

	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(keys.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(keys.SecondaryKey))
	collector.AddValue(operatorSpecSecrets.PrimaryConnectionString, to.Value(keys.PrimaryConnectionString))
	collector.AddValue(operatorSpecSecrets.SecondaryConnectionString, to.Value(keys.SecondaryConnectionString))

	return collector.Values()
}

func makeNamespacesEventHubAuthorizationRuleResolvedSecretsMap(accessKeys armeventhub.AccessKeys) map[string]string {
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
