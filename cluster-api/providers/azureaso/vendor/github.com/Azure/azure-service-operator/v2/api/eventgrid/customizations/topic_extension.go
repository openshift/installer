/* Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-service-operator/v2/api/eventgrid/v1api20200601/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	key1 = "key1"
	key2 = "key2"
)

var _ genruntime.KubernetesSecretExporter = &TopicExtension{}

func (ext *TopicExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*storage.Topic)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *eventgrid.Topic", obj)
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

	subscription := id.SubscriptionID
	// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
	// connection each time through
	var confClient *armeventgrid.TopicsClient
	confClient, err = armeventgrid.NewTopicsClient(subscription, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create new TopicsClient")
	}

	var resp armeventgrid.TopicsClientListSharedAccessKeysResponse
	resp, err = confClient.ListSharedAccessKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
	if err != nil {
		return nil, eris.Wrapf(err, "failed listing keys")
	}

	secretSlice, err := secretsToWrite(typedObj, resp)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := map[string]string{}
	if to.Value(resp.Key1) != "" {
		resolvedSecrets[key1] = to.Value(resp.Key1)
	}
	if to.Value(resp.Key2) != "" {
		resolvedSecrets[key2] = to.Value(resp.Key2)
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *storage.Topic) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	secrets := obj.Spec.OperatorSpec.Secrets
	result := make(set.Set[string])
	if secrets.Key1 != nil {
		result.Add(key1)
	}
	if secrets.Key2 != nil {
		result.Add(key2)
	}

	return result
}

func secretsToWrite(obj *storage.Topic, keys armeventgrid.TopicsClientListSharedAccessKeysResponse) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.Key1, to.Value(keys.Key1))
	collector.AddValue(operatorSpecSecrets.Key2, to.Value(keys.Key2))

	return collector.Values()
}
