/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v2"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	apimanagement "github.com/Azure/azure-service-operator/v2/api/apimanagement/v1api20220801/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	primaryKey   = "primaryKey"
	secondaryKey = "secondaryKey"
)

var _ genruntime.KubernetesSecretExporter = &SubscriptionExtension{}

func (ext *SubscriptionExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*apimanagement.Subscription)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *apimanagement.Subscription", obj)
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

	if id.Parent == nil {
		return nil, eris.Errorf("APIM subscription had no parent ID: %s", id.String())
	}
	parentName := id.Parent.Name

	// Only bother calling ListSecrets if there are secrets to retrieve
	var s armapimanagement.SubscriptionKeysContract
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var subClient *armapimanagement.SubscriptionClient
		subClient, err = armapimanagement.NewSubscriptionClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create new SubscriptionClient")
		}

		var resp armapimanagement.SubscriptionClientListSecretsResponse
		resp, err = subClient.ListSecrets(ctx, id.ResourceGroupName, parentName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, eris.Wrapf(err, "failed listing secrets")
		}

		s = resp.SubscriptionKeysContract
	}

	resolvedSecrets := map[string]string{}
	if to.Value(s.PrimaryKey) != "" {
		resolvedSecrets[primaryKey] = to.Value(s.PrimaryKey)
	}
	if to.Value(s.SecondaryKey) != "" {
		resolvedSecrets[secondaryKey] = to.Value(s.SecondaryKey)
	}
	secretSlice, err := secretsToWrite(typedObj, s)
	if err != nil {
		return nil, err
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *apimanagement.Subscription) set.Set[string] {
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

func secretsToWrite(obj *apimanagement.Subscription, s armapimanagement.SubscriptionKeysContract) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(s.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(s.SecondaryKey))

	return collector.Values()
}
