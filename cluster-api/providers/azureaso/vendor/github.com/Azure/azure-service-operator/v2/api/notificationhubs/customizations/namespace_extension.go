/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/notificationhubs/armnotificationhubs"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-service-operator/v2/api/notificationhubs/v1api20230901/storage"
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

var _ genruntime.KubernetesSecretExporter = &NamespaceExtension{}

func (ext *NamespaceExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*storage.Namespace)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *notificationhubs.Namespace", obj)
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

	var resp armnotificationhubs.NamespacesClientListKeysResponse
	// Only bother calling ListKeys if there are secrets to retrieve
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var nsClient *armnotificationhubs.NamespacesClient
		nsClient, err = armnotificationhubs.NewNamespacesClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create new AccountsClient")
		}

		// "RootManageSharedAccessKey" is the default AuthorizationRule
		resp, err = nsClient.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), "RootManageSharedAccessKey", nil)
		if err != nil {
			return nil, eris.Wrapf(err, "failed listing keys")
		}
	}

	secretSlice, err := secretsToWrite(typedObj, resp)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := map[string]string{}
	addSecretsToMap(resp, resolvedSecrets)

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func addSecretsToMap(resp armnotificationhubs.NamespacesClientListKeysResponse, result map[string]string) {
	if to.Value(resp.PrimaryKey) != "" {
		result[primaryKey] = to.Value(resp.PrimaryConnectionString)
	}
	if to.Value(resp.SecondaryKey) != "" {
		result[secondaryKey] = to.Value(resp.SecondaryKey)
	}
	if to.Value(resp.PrimaryConnectionString) != "" {
		result[primaryConnectionString] = to.Value(resp.PrimaryConnectionString)
	}
	if to.Value(resp.SecondaryConnectionString) != "" {
		result[secondaryConnectionString] = to.Value(resp.SecondaryConnectionString)
	}
}

func secretsSpecified(obj *storage.Namespace) set.Set[string] {
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

func secretsToWrite(obj *storage.Namespace, resp armnotificationhubs.NamespacesClientListKeysResponse) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryConnectionString, to.Value(resp.PrimaryConnectionString))
	collector.AddValue(operatorSpecSecrets.SecondaryConnectionString, to.Value(resp.SecondaryConnectionString))
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(resp.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(resp.SecondaryKey))

	return collector.Values()
}
