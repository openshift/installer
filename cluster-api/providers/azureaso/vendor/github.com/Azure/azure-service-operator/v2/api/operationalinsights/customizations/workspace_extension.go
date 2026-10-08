/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-service-operator/v2/api/operationalinsights/v20250701/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	primarySharedKey   = "primarySharedKey"
	secondarySharedKey = "secondarySharedKey"
)

var _ genruntime.KubernetesSecretExporter = &WorkspaceExtension{}

func (ext *WorkspaceExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*storage.Workspace)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *storage.Workspace", obj)
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
	var sharedKeysClient *armoperationalinsights.SharedKeysClient
	sharedKeysClient, err = armoperationalinsights.NewSharedKeysClient(subscription, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create new SharedKeysClient")
	}

	var resp armoperationalinsights.SharedKeysClientGetSharedKeysResponse
	resp, err = sharedKeysClient.GetSharedKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
	if err != nil {
		return nil, eris.Wrapf(err, "failed getting shared keys")
	}

	secretSlice, err := secretsToWrite(typedObj, resp)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := map[string]string{}
	if to.Value(resp.PrimarySharedKey) != "" {
		resolvedSecrets[primarySharedKey] = to.Value(resp.PrimarySharedKey)
	}
	if to.Value(resp.SecondarySharedKey) != "" {
		resolvedSecrets[secondarySharedKey] = to.Value(resp.SecondarySharedKey)
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *storage.Workspace) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	secrets := obj.Spec.OperatorSpec.Secrets
	result := make(set.Set[string])
	if secrets.PrimarySharedKey != nil {
		result.Add(primarySharedKey)
	}
	if secrets.SecondarySharedKey != nil {
		result.Add(secondarySharedKey)
	}

	return result
}

func secretsToWrite(obj *storage.Workspace, keys armoperationalinsights.SharedKeysClientGetSharedKeysResponse) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimarySharedKey, to.Value(keys.PrimarySharedKey))
	collector.AddValue(operatorSpecSecrets.SecondarySharedKey, to.Value(keys.SecondarySharedKey))

	return collector.Values()
}
