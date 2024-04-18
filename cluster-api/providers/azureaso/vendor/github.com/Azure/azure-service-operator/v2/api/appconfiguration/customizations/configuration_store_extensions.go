/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appconfiguration/armappconfiguration"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	storage "github.com/Azure/azure-service-operator/v2/api/appconfiguration/v1api20220501/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &ConfigurationStoreExtension{}

func (ext *ConfigurationStoreExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*storage.ConfigurationStore)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *appconfiguration.ConfigurationStore", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	hasSecrets := secretsSpecified(typedObj)
	if !hasSecrets {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	keys := make(map[string]armappconfiguration.APIKey)
	// Only bother calling ListKeys if there are secrets to retrieve
	if hasSecrets {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var confClient *armappconfiguration.ConfigurationStoresClient
		confClient, err = armappconfiguration.NewConfigurationStoresClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new ConfigurationStoresClient")
		}

		var pager *runtime.Pager[armappconfiguration.ConfigurationStoresClientListKeysResponse]
		var resp armappconfiguration.ConfigurationStoresClientListKeysResponse
		pager = confClient.NewListKeysPager(id.ResourceGroupName, typedObj.AzureName(), nil)
		for pager.More() {
			resp, err = pager.NextPage(ctx)
			addSecretsToMap(resp.Value, keys)

		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retreive response")
		}
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing keys")
		}

	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func secretsSpecified(obj *storage.ConfigurationStore) bool {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	secrets := obj.Spec.OperatorSpec.Secrets

	if secrets.PrimaryKeyID != nil ||
		secrets.SecondaryKeyID != nil ||
		secrets.PrimaryReadOnlyKeyID != nil ||
		secrets.SecondaryReadOnlyKeyID != nil ||
		secrets.PrimaryKey != nil ||
		secrets.SecondaryKey != nil ||
		secrets.PrimaryReadOnlyKey != nil ||
		secrets.SecondaryReadOnlyKey != nil ||
		secrets.PrimaryConnectionString != nil ||
		secrets.SecondaryConnectionString != nil ||
		secrets.PrimaryReadOnlyConnectionString != nil ||
		secrets.SecondaryReadOnlyConnectionString != nil {
		return true
	}

	return false
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
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	primary, ok := keys["Primary"]
	if ok {
		collector.AddValue(operatorSpecSecrets.PrimaryConnectionString, *primary.ConnectionString)
		collector.AddValue(operatorSpecSecrets.PrimaryKeyID, *primary.ID)
		collector.AddValue(operatorSpecSecrets.PrimaryKey, *primary.Value)
	}

	primaryReadOnly, ok := keys["Primary Read Only"]
	if ok {
		collector.AddValue(operatorSpecSecrets.PrimaryReadOnlyConnectionString, *primaryReadOnly.ConnectionString)
		collector.AddValue(operatorSpecSecrets.PrimaryReadOnlyKeyID, *primaryReadOnly.ID)
		collector.AddValue(operatorSpecSecrets.PrimaryReadOnlyKey, *primaryReadOnly.Value)
	}

	secondary, ok := keys["Secondary"]
	if ok {
		collector.AddValue(operatorSpecSecrets.SecondaryConnectionString, *secondary.ConnectionString)
		collector.AddValue(operatorSpecSecrets.SecondaryKeyID, *secondary.ID)
		collector.AddValue(operatorSpecSecrets.SecondaryKey, *secondary.Value)
	}

	secondaryReadOnly, ok := keys["Secondary Read Only"]
	if ok {
		collector.AddValue(operatorSpecSecrets.SecondaryReadOnlyConnectionString, *secondaryReadOnly.ConnectionString)
		collector.AddValue(operatorSpecSecrets.SecondaryReadOnlyKeyID, *secondaryReadOnly.ID)
		collector.AddValue(operatorSpecSecrets.SecondaryReadOnlyKey, *secondaryReadOnly.Value)
	}

	return collector.Values()
}
