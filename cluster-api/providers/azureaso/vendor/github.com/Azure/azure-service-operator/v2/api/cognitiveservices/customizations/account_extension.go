/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	storage "github.com/Azure/azure-service-operator/v2/api/cognitiveservices/v1api20250601/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	accountKey1 = "key1"
	accountKey2 = "key2"
)

var _ genruntime.KubernetesSecretExporter = &AccountExtension{}

func (ext *AccountExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	typedObj, ok := obj.(*storage.Account)

	if !ok {
		return nil, eris.Errorf("Cannot run on unknown resource type %T, expected *storage.Account", obj)
	}

	var _ conversion.Hub = typedObj

	primarySecrets := secretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)

	if len(requestedSecrets) == 0 {
		log.V(1).Info("No secrets to retrieve as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	keys := make(map[string]string)

	if primarySecrets.Contains(accountKey1) || primarySecrets.Contains(accountKey2) {
		subscription := id.SubscriptionID
		acctClient, err := armcognitiveservices.NewAccountsClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "Failed to create new AccountsClient")
		}

		resp, err := acctClient.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, eris.Wrapf(err, "Failed listing keys")
		}

		keys[accountKey1] = to.Value(resp.Key1)
		keys[accountKey2] = to.Value(resp.Key2)
	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, keys),
	}, nil
}

func secretsSpecified(obj *storage.Account) set.Set[string] {
	result := make(set.Set[string])

	if obj.Spec.OperatorSpec != nil && obj.Spec.OperatorSpec.Secrets != nil {
		s := obj.Spec.OperatorSpec.Secrets

		if s.Key1 != nil {
			result.Add(accountKey1)
		}

		if s.Key2 != nil {
			result.Add(accountKey2)
		}
	}
	return result
}

func secretsToWrite(obj *storage.Account, keys map[string]string) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.Key1, keys[accountKey1])
	collector.AddValue(operatorSpecSecrets.Key2, keys[accountKey2])

	return collector.Values()
}
