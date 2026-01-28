/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appconfiguration/armappconfiguration"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	storage "github.com/Azure/azure-service-operator/v2/api/appconfiguration/v1api20220501/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	primaryKeyID                      = "primaryKeyID"
	secondaryKeyID                    = "secondaryKeyID"
	primaryReadOnlyKeyID              = "primaryReadOnlyKeyID"
	secondaryReadOnlyKeyID            = "secondaryReadOnlyKeyID"
	primaryKey                        = "primaryKey"
	secondaryKey                      = "secondaryKey"
	primaryReadOnlyKey                = "primaryReadOnlyKey"
	secondaryReadOnlyKey              = "secondaryReadOnlyKey"
	primaryConnectionString           = "primaryConnectionString"
	secondaryConnectionString         = "secondaryConnectionString"
	primaryReadOnlyConnectionString   = "primaryReadOnlyConnectionString"
	secondaryReadOnlyConnectionString = "secondaryReadOnlyConnectionString"
)

var _ genruntime.KubernetesSecretExporter = &ConfigurationStoreExtension{}

func (ext *ConfigurationStoreExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*storage.ConfigurationStore)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *appconfiguration.ConfigurationStore", obj)
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

	keys := make(map[string]armappconfiguration.APIKey)
	// Only bother calling ListKeys if there are secrets to retrieve
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var confClient *armappconfiguration.ConfigurationStoresClient
		confClient, err = armappconfiguration.NewConfigurationStoresClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create new ConfigurationStoresClient")
		}

		var pager *runtime.Pager[armappconfiguration.ConfigurationStoresClientListKeysResponse]
		var resp armappconfiguration.ConfigurationStoresClientListKeysResponse
		pager = confClient.NewListKeysPager(id.ResourceGroupName, typedObj.AzureName(), nil)
		for pager.More() {
			resp, err = pager.NextPage(ctx)
			if err != nil {
				return nil, eris.Wrapf(err, "failed to retreive response")
			}
			addSecretsToMap(resp.Value, keys)
		}
	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := makeResolvedSecretsMap(keys)

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *storage.ConfigurationStore) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	secrets := obj.Spec.OperatorSpec.Secrets

	result := make(set.Set[string])
	if secrets.PrimaryKeyID != nil {
		result.Add(primaryKeyID)
	}
	if secrets.SecondaryKeyID != nil {
		result.Add(secondaryKeyID)
	}
	if secrets.PrimaryReadOnlyKeyID != nil {
		result.Add(primaryReadOnlyKeyID)
	}
	if secrets.SecondaryReadOnlyKeyID != nil {
		result.Add(secondaryReadOnlyKeyID)
	}
	if secrets.PrimaryKey != nil {
		result.Add(primaryKey)
	}
	if secrets.SecondaryKey != nil {
		result.Add(secondaryKey)
	}
	if secrets.PrimaryReadOnlyKey != nil {
		result.Add(primaryReadOnlyKey)
	}
	if secrets.SecondaryReadOnlyKey != nil {
		result.Add(secondaryReadOnlyKey)
	}
	if secrets.PrimaryConnectionString != nil {
		result.Add(primaryConnectionString)
	}
	if secrets.SecondaryConnectionString != nil {
		result.Add(secondaryConnectionString)
	}
	if secrets.PrimaryReadOnlyConnectionString != nil {
		result.Add(primaryReadOnlyConnectionString)
	}
	if secrets.SecondaryReadOnlyConnectionString != nil {
		result.Add(secondaryReadOnlyConnectionString)
	}
	return result
}

func addSecretsToMap(keys []*armappconfiguration.APIKey, result map[string]armappconfiguration.APIKey) {
	for _, key := range keys {
		if key == nil || key.Name == nil {
			continue
		}
		result[*key.Name] = *key
	}
}

func secretsToWrite(obj *storage.ConfigurationStore, keys map[string]armappconfiguration.APIKey) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	primary, ok := keys["Primary"]
	if ok {
		collector.AddValue(operatorSpecSecrets.PrimaryConnectionString, to.Value(primary.ConnectionString))
		collector.AddValue(operatorSpecSecrets.PrimaryKeyID, to.Value(primary.ID))
		collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(primary.Value))
	}

	primaryReadOnly, ok := keys["Primary Read Only"]
	if ok {
		collector.AddValue(operatorSpecSecrets.PrimaryReadOnlyConnectionString, to.Value(primaryReadOnly.ConnectionString))
		collector.AddValue(operatorSpecSecrets.PrimaryReadOnlyKeyID, to.Value(primaryReadOnly.ID))
		collector.AddValue(operatorSpecSecrets.PrimaryReadOnlyKey, to.Value(primaryReadOnly.Value))
	}

	secondary, ok := keys["Secondary"]
	if ok {
		collector.AddValue(operatorSpecSecrets.SecondaryConnectionString, to.Value(secondary.ConnectionString))
		collector.AddValue(operatorSpecSecrets.SecondaryKeyID, to.Value(secondary.ID))
		collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(secondary.Value))
	}

	secondaryReadOnly, ok := keys["Secondary Read Only"]
	if ok {
		collector.AddValue(operatorSpecSecrets.SecondaryReadOnlyConnectionString, to.Value(secondaryReadOnly.ConnectionString))
		collector.AddValue(operatorSpecSecrets.SecondaryReadOnlyKeyID, to.Value(secondaryReadOnly.ID))
		collector.AddValue(operatorSpecSecrets.SecondaryReadOnlyKey, to.Value(secondaryReadOnly.Value))
	}

	return collector.Values()
}

func makeResolvedSecretsMap(keys map[string]armappconfiguration.APIKey) map[string]string {
	result := make(map[string]string)
	primary, ok := keys["Primary"]
	if ok {
		result[primaryConnectionString] = to.Value(primary.ConnectionString)
		result[primaryKeyID] = to.Value(primary.ID)
		result[primaryKey] = to.Value(primary.Value)
	}

	primaryReadOnly, ok := keys["Primary Read Only"]
	if ok {
		result[primaryReadOnlyConnectionString] = to.Value(primaryReadOnly.ConnectionString)
		result[primaryReadOnlyKeyID] = to.Value(primaryReadOnly.ID)
		result[primaryReadOnlyKey] = to.Value(primaryReadOnly.Value)
	}

	secondary, ok := keys["Secondary"]
	if ok {
		result[secondaryConnectionString] = to.Value(secondary.ConnectionString)
		result[secondaryKeyID] = to.Value(secondary.ID)
		result[secondaryKey] = to.Value(secondary.Value)
	}

	secondaryReadOnly, ok := keys["Secondary Read Only"]
	if ok {
		result[secondaryReadOnlyConnectionString] = to.Value(secondaryReadOnly.ConnectionString)
		result[secondaryReadOnlyKeyID] = to.Value(secondaryReadOnly.ID)
		result[secondaryReadOnlyKey] = to.Value(secondaryReadOnly.Value)
	}

	return result
}
