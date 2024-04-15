/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/signalr/armsignalr"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	signalr "github.com/Azure/azure-service-operator/v2/api/signalrservice/v1api20211001/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

// Ensure that SignalRAuthorizationExtension implements the KubernetesExporter interface
var _ genruntime.KubernetesExporter = &SignalRExtension{}

// ExportKubernetesResources implements genruntime.KubernetesExporter
func (*SignalRExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {
	// Make sure we're working with the current hub version of the resource
	// This will need to be updated if the hub version changes
	typedObj, ok := obj.(*signalr.SignalR)
	if !ok {
		return nil, errors.Errorf(
			"cannot run on unknown resource type %T, expected *signalr.SignalR", obj)
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

	// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
	// connection each time through
	clientFactory, err := armsignalr.NewClientFactory(id.SubscriptionID, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create ARM signalR client factory")
	}

	res, err := clientFactory.NewClient().ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list keys")
	}

	secretSlice, err := secretsToWrite(typedObj, res.Keys)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil

}

func secretsSpecified(obj *signalr.SignalR) bool {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false
	}

	specSecrets := obj.Spec.OperatorSpec.Secrets
	hasSecrets := false
	if specSecrets.PrimaryKey != nil ||
		specSecrets.SecondaryKey != nil ||
		specSecrets.PrimaryConnectionString != nil ||
		specSecrets.SecondaryConnectionString != nil {
		hasSecrets = true
	}

	return hasSecrets
}

func secretsToWrite(obj *signalr.SignalR, accessKeys armsignalr.Keys) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryKey, to.Value(accessKeys.PrimaryKey))
	collector.AddValue(operatorSpecSecrets.SecondaryKey, to.Value(accessKeys.SecondaryKey))
	collector.AddValue(operatorSpecSecrets.PrimaryConnectionString, to.Value(accessKeys.PrimaryConnectionString))
	collector.AddValue(operatorSpecSecrets.SecondaryConnectionString, to.Value(accessKeys.SecondaryConnectionString))

	return collector.Values()
}
