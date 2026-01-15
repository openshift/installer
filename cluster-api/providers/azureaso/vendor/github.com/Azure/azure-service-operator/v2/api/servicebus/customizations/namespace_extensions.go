/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	servicebus "github.com/Azure/azure-service-operator/v2/api/servicebus/v1api20240101/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	endpoint                  = "endpoint"
	primaryKey                = "primaryKey"
	primaryConnectionString   = "primaryConnectionString"
	secondaryKey              = "secondaryKey"
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
	namespace, ok := obj.(*servicebus.Namespace)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *servicebus.Namespace", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = namespace

	primarySecrets := namespaceSecretsSpecified(namespace)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)

	if len(requestedSecrets) == 0 {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(namespace)
	if err != nil {
		return nil, err
	}

	// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
	// connection each time through
	clientFactory, err := armservicebus.NewClientFactory(
		id.SubscriptionID,
		armClient.Creds(),
		armClient.ClientOptions())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create ARM servicebus client factory")
	}

	// This access rule always exists and provides management access to the namespace
	// See https://learn.microsoft.com/en-us/azure/service-bus-messaging/service-bus-sas
	const rootRuleName = "RootManageSharedAccessKey"

	client := clientFactory.NewNamespacesClient()
	options := armservicebus.NamespacesClientListKeysOptions{}
	response, err := client.ListKeys(
		ctx,
		id.ResourceGroupName,
		id.Name,
		rootRuleName,
		&options)
	if err != nil {
		return nil, eris.Wrapf(
			err,
			"failed to retrieve namespace management keys from authorization rule %q",
			rootRuleName)
	}

	secretSlice, err := namespaceSecretsToWrite(namespace, response)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := map[string]string{}
	if to.Value(namespace.Status.ServiceBusEndpoint) != "" {
		resolvedSecrets[endpoint] = to.Value(namespace.Status.ServiceBusEndpoint)
	}
	if to.Value(response.PrimaryKey) != "" {
		resolvedSecrets[primaryKey] = to.Value(response.PrimaryKey)
	}
	if to.Value(response.PrimaryConnectionString) != "" {
		resolvedSecrets[primaryConnectionString] = to.Value(response.PrimaryConnectionString)
	}
	if to.Value(response.SecondaryKey) != "" {
		resolvedSecrets[secondaryKey] = to.Value(response.SecondaryKey)
	}
	if to.Value(response.SecondaryConnectionString) != "" {
		resolvedSecrets[secondaryConnectionString] = to.Value(response.SecondaryConnectionString)
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func namespaceSecretsSpecified(obj *servicebus.Namespace) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	specSecrets := obj.Spec.OperatorSpec.Secrets

	result := make(set.Set[string])
	if specSecrets.Endpoint != nil {
		result.Add(endpoint)
	}
	if specSecrets.PrimaryKey != nil {
		result.Add(primaryKey)
	}
	if specSecrets.PrimaryConnectionString != nil {
		result.Add(primaryConnectionString)
	}
	if specSecrets.SecondaryKey != nil {
		result.Add(secondaryKey)
	}
	if specSecrets.SecondaryConnectionString != nil {
		result.Add(secondaryConnectionString)
	}
	return result
}

func namespaceSecretsToWrite(
	obj *servicebus.Namespace,
	response armservicebus.NamespacesClientListKeysResponse,
) ([]*v1.Secret, error) {
	specSecrets := obj.Spec.OperatorSpec.Secrets
	if specSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(specSecrets.Endpoint, to.Value(obj.Status.ServiceBusEndpoint))
	collector.AddValue(specSecrets.PrimaryKey, to.Value(response.PrimaryKey))
	collector.AddValue(specSecrets.PrimaryConnectionString, to.Value(response.PrimaryConnectionString))
	collector.AddValue(specSecrets.SecondaryKey, to.Value(response.SecondaryKey))
	collector.AddValue(specSecrets.SecondaryConnectionString, to.Value(response.SecondaryConnectionString))

	return collector.Values()
}
