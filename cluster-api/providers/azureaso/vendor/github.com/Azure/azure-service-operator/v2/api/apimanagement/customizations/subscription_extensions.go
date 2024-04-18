/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	apimanagement "github.com/Azure/azure-service-operator/v2/api/apimanagement/v1api20220801/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &SubscriptionExtension{}

func (ext *SubscriptionExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*apimanagement.Subscription)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *apimanagement.Subscription", obj)
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

	if id.Parent == nil {
		return nil, errors.Errorf("APIM subscription had no parent ID: %s", id.String())
	}
	parentName := id.Parent.Name

	// Only bother calling ListSecrets if there are secrets to retrieve
	var s armapimanagement.SubscriptionKeysContract
	if hasSecrets {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var subClient *armapimanagement.SubscriptionClient
		subClient, err = armapimanagement.NewSubscriptionClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new SubscriptionClient")
		}

		var resp armapimanagement.SubscriptionClientListSecretsResponse
		resp, err = subClient.ListSecrets(ctx, id.ResourceGroupName, parentName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing secrets")
		}

		s = resp.SubscriptionKeysContract
	}

	secretSlice, err := secretsToWrite(typedObj, s)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func secretsSpecified(obj *apimanagement.Subscription) bool {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	hasSecrets := false
	secrets := obj.Spec.OperatorSpec.Secrets
	if secrets.PrimaryKey != nil || secrets.SecondaryKey != nil {
		hasSecrets = true
	}

	return hasSecrets
}

func secretsToWrite(obj *apimanagement.Subscription, s armapimanagement.SubscriptionKeysContract) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(s.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(s.SecondaryKey))

	return collector.Values()
}
