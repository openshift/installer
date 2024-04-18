/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	documentdb "github.com/Azure/azure-service-operator/v2/api/documentdb/v1api20210515/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &DatabaseAccountExtension{}

func (ext *DatabaseAccountExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*documentdb.DatabaseAccount)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *documentdb.DatabaseAccount", obj)
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

	var keys armcosmos.DatabaseAccountListKeysResult
	// Only bother calling ListKeys if there are secrets to retrieve
	if hasSecrets {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var acctClient *armcosmos.DatabaseAccountsClient
		acctClient, err = armcosmos.NewDatabaseAccountsClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new DatabaseAccountClient")
		}

		// TODO: There is a ListReadOnlyKeys API that requires less permissions. We should consider determining
		// TODO: that we don't need to call the ListKeys API and install call the listReadOnlyKeys API.
		var resp armcosmos.DatabaseAccountsClientListKeysResponse
		resp, err = acctClient.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing keys")
		}

		keys = resp.DatabaseAccountListKeysResult
	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func secretsSpecified(obj *documentdb.DatabaseAccount) (bool, bool) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false, false
	}

	specSecrets := obj.Spec.OperatorSpec.Secrets
	hasSecrets := false
	hasEndpoints := false
	if specSecrets.PrimaryMasterKey != nil ||
		specSecrets.SecondaryMasterKey != nil ||
		specSecrets.PrimaryReadonlyMasterKey != nil ||
		specSecrets.SecondaryReadonlyMasterKey != nil {
		hasSecrets = true
	}

	if specSecrets.DocumentEndpoint != nil {
		hasEndpoints = true
	}

	return hasSecrets, hasEndpoints
}

func secretsToWrite(obj *documentdb.DatabaseAccount, accessKeys armcosmos.DatabaseAccountListKeysResult) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryMasterKey, to.Value(accessKeys.PrimaryMasterKey))
	collector.AddValue(operatorSpecSecrets.SecondaryMasterKey, to.Value(accessKeys.SecondaryMasterKey))
	collector.AddValue(operatorSpecSecrets.PrimaryReadonlyMasterKey, to.Value(accessKeys.PrimaryReadonlyMasterKey))
	collector.AddValue(operatorSpecSecrets.SecondaryReadonlyMasterKey, to.Value(accessKeys.SecondaryReadonlyMasterKey))
	collector.AddValue(operatorSpecSecrets.DocumentEndpoint, to.Value(obj.Status.DocumentEndpoint))

	return collector.Values()
}
