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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	redis "github.com/Azure/azure-service-operator/v2/api/cache/v1api20230401/storage"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &RedisExtension{}

func (ext *RedisExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*redis.Redis)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *redis.Redis", obj)
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

	var accessKeys armredis.AccessKeys
	// Only bother calling ListKeys if there are secrets to retrieve
	if hasSecrets {
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

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func secretsSpecified(obj *redis.Redis) (bool, bool) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false, false
	}

	secrets := obj.Spec.OperatorSpec.Secrets
	hasSecrets := false
	hasEndpoints := false

	if secrets.PrimaryKey != nil ||
		secrets.SecondaryKey != nil {
		hasSecrets = true
	}

	if secrets.HostName != nil ||
		secrets.Port != nil ||
		secrets.SSLPort != nil {
		hasEndpoints = true
	}

	return hasSecrets, hasEndpoints
}

func secretsToWrite(obj *redis.Redis, accessKeys armredis.AccessKeys) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
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
