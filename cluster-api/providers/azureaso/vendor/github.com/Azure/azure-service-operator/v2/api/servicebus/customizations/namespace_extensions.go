/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"

	servicebus "github.com/Azure/azure-service-operator/v2/api/servicebus/v1api20211101/storage"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &NamespaceExtension{}

func (ext *NamespaceExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	namespace, ok := obj.(*servicebus.Namespace)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *servicebus.Namespace", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = namespace

	hasSecrets := namespaceSecretsSpecified(namespace)
	if !hasSecrets {
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
		return nil, errors.Wrapf(err, "failed to create ARM servicebus client factory")
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
		return nil, errors.Wrapf(
			err,
			"failed to retrieve namespace management keys from authorization rule %q",
			rootRuleName)
	}

	secretSlice, err := namespaceSecretsToWrite(namespace, response)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func namespaceSecretsSpecified(obj *servicebus.Namespace) bool {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	specSecrets := obj.Spec.OperatorSpec.Secrets

	return specSecrets.Endpoint != nil ||
		specSecrets.PrimaryKey != nil ||
		specSecrets.PrimaryConnectionString != nil ||
		specSecrets.SecondaryKey != nil ||
		specSecrets.SecondaryConnectionString != nil

}

func namespaceSecretsToWrite(
	obj *servicebus.Namespace,
	response armservicebus.NamespacesClientListKeysResponse,
) ([]*v1.Secret, error) {
	specSecrets := obj.Spec.OperatorSpec.Secrets
	if specSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(specSecrets.Endpoint, to.Value(obj.Status.ServiceBusEndpoint))
	collector.AddValue(specSecrets.PrimaryKey, *response.PrimaryKey)
	collector.AddValue(specSecrets.PrimaryConnectionString, *response.PrimaryConnectionString)
	collector.AddValue(specSecrets.SecondaryKey, *response.SecondaryKey)
	collector.AddValue(specSecrets.SecondaryConnectionString, *response.SecondaryConnectionString)

	return collector.Values()
}
