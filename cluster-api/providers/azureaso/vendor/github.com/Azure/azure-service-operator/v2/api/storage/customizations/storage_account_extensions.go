/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-service-operator/v2/api/storage/v1api20230101/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &StorageAccountExtension{}

func (ext *StorageAccountExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*storage.StorageAccount)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *storage.StorageAccount", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	hasSecrets, hasEndpoints := secretsSpecified(typedObj)
	if !hasSecrets && !hasEndpoints {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	keys := make(map[string]string)
	// Only bother calling ListKeys if there are secrets to retrieve
	if hasSecrets {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var acctClient *armstorage.AccountsClient
		acctClient, err = armstorage.NewAccountsClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new AccountsClient")
		}

		var resp armstorage.AccountsClientListKeysResponse
		resp, err = acctClient.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing keys")
		}

		keys = secretsByName(resp.Keys)
	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func secretsSpecified(obj *storage.StorageAccount) (bool, bool) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false, false
	}

	hasSecrets := false
	hasEndpoints := false
	secrets := obj.Spec.OperatorSpec.Secrets
	if secrets.Key1 != nil || secrets.Key2 != nil {
		hasSecrets = true
	}

	if secrets.BlobEndpoint != nil ||
		secrets.QueueEndpoint != nil ||
		secrets.TableEndpoint != nil ||
		secrets.FileEndpoint != nil ||
		secrets.WebEndpoint != nil ||
		secrets.DfsEndpoint != nil {
		hasEndpoints = true
	}

	return hasSecrets, hasEndpoints
}

func secretsByName(keys []*armstorage.AccountKey) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if key.KeyName == nil || key.Value == nil {
			continue
		}
		result[*key.KeyName] = *key.Value
	}

	return result
}

func secretsToWrite(obj *storage.StorageAccount, keys map[string]string) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.Key1, keys["key1"])
	collector.AddValue(operatorSpecSecrets.Key2, keys["key2"])
	// There are tons of different endpoints we could write, including secondary endpoints.
	// For now we're just exposing the main ones. See:
	// https://docs.microsoft.com/en-us/rest/api/storagerp/storage-accounts/get-properties for more details
	if obj.Status.PrimaryEndpoints != nil {
		eps := obj.Status.PrimaryEndpoints
		collector.AddValue(operatorSpecSecrets.BlobEndpoint, to.Value(eps.Blob))
		collector.AddValue(operatorSpecSecrets.QueueEndpoint, to.Value(eps.Queue))
		collector.AddValue(operatorSpecSecrets.TableEndpoint, to.Value(eps.Table))
		collector.AddValue(operatorSpecSecrets.FileEndpoint, to.Value(eps.File))
		collector.AddValue(operatorSpecSecrets.WebEndpoint, to.Value(eps.Web))
		collector.AddValue(operatorSpecSecrets.DfsEndpoint, to.Value(eps.Dfs))
	}

	return collector.Values()
}
