/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redis/armredis"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	redis "github.com/Azure/azure-service-operator/v2/api/cache/v1api20230801/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	primaryKey   = "primaryKey"
	secondaryKey = "secondaryKey"
)

var _ genruntime.KubernetesSecretExporter = &RedisExtension{}

func (ext *RedisExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*redis.Redis)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *redis.Redis", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	primarySecrets, hasEndpoints := secretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)
	if len(requestedSecrets) == 0 && !hasEndpoints {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	var accessKeys armredis.AccessKeys
	// Only bother calling ListKeys if there are secrets to retrieve
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var redisClient *armredis.Client
		redisClient, err = armredis.NewClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new new RedisClient")
		}

		var resp armredis.ClientListKeysResponse
		resp, err = redisClient.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing keys")
		}
		accessKeys = resp.AccessKeys
	}
	secretSlice, err := secretsToWrite(typedObj, accessKeys)
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

func secretsSpecified(obj *redis.Redis) (set.Set[string], bool) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil, false
	}

	secrets := obj.Spec.OperatorSpec.Secrets
	result := make(set.Set[string])
	hasEndpoints := false

	if secrets.PrimaryKey != nil {
		result.Add(primaryKey)
	}
	if secrets.SecondaryKey != nil {
		result.Add(secondaryKey)
	}

	if secrets.HostName != nil ||
		secrets.Port != nil ||
		secrets.SSLPort != nil {
		hasEndpoints = true
	}

	return result, hasEndpoints
}

func secretsToWrite(obj *redis.Redis, accessKeys armredis.AccessKeys) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(accessKeys.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(accessKeys.SecondaryKey))
	collector.AddValue(operatorSpecSecrets.HostName, to.Value(obj.Status.HostName))
	collector.AddValue(operatorSpecSecrets.Port, intPtrToString(obj.Status.Port))
	collector.AddValue(operatorSpecSecrets.SSLPort, intPtrToString(obj.Status.SslPort))

	return collector.Values()
}

func intPtrToString(i *int) string {
	if i == nil {
		return ""
	}

	return strconv.Itoa(*i)
}
