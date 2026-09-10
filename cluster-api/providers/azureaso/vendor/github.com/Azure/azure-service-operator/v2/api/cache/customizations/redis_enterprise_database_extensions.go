/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	redisenterprise "github.com/Azure/azure-service-operator/v2/api/cache/v1api20250401/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesSecretExporter = &RedisEnterpriseDatabaseExtension{}

func (ext *RedisEnterpriseDatabaseExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*redisenterprise.RedisEnterpriseDatabase)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *redisenterprise.RedisEnterpriseDatabase", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	primarySecrets := redisEnterpriseSecretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)
	if len(requestedSecrets) == 0 {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	var accessKeys armredisenterprise.AccessKeys
	// Only bother calling ListKeys if there are secrets to retrieve
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var databaseClient *armredisenterprise.DatabasesClient
		databaseClient, err = armredisenterprise.NewDatabasesClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create new DatabasesClient")
		}

		var resp armredisenterprise.DatabasesClientListKeysResponse
		resp, err = databaseClient.ListKeys(ctx, id.ResourceGroupName, id.Parent.Name, typedObj.AzureName(), nil)
		if err != nil {
			return nil, eris.Wrapf(err, "failed listing keys")
		}
		accessKeys = resp.AccessKeys
	}
	secretSlice, err := redisEnterpriseSecretsToWrite(typedObj, accessKeys)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := map[string]string{}
	if to.Value(accessKeys.PrimaryKey) != "" {
		resolvedSecrets[primaryKey] = to.Value(accessKeys.PrimaryKey)
	}
	if to.Value(accessKeys.SecondaryKey) != "" {
		resolvedSecrets[secondaryKey] = to.Value(accessKeys.SecondaryKey)
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func redisEnterpriseSecretsSpecified(obj *redisenterprise.RedisEnterpriseDatabase) set.Set[string] {
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

	return result
}

func redisEnterpriseSecretsToWrite(obj *redisenterprise.RedisEnterpriseDatabase, accessKeys armredisenterprise.AccessKeys) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(accessKeys.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(accessKeys.SecondaryKey))

	return collector.Values()
}
