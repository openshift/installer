/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redhatopenshift/armredhatopenshift"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	redhatopenshift "github.com/Azure/azure-service-operator/v2/api/redhatopenshift/v1api20231122/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	adminCredentialsKey = "adminCredentials"
	usernameKey         = "username"
	passwordKey         = "password"
)

var _ genruntime.KubernetesSecretExporter = &OpenShiftClusterExtension{}

func (ext *OpenShiftClusterExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*redhatopenshift.OpenShiftCluster)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *redhatopenshift.OpenShiftCluster", obj)
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
	var clusterClient *armredhatopenshift.OpenShiftClustersClient
	clusterClient, err = armredhatopenshift.NewOpenShiftClustersClient(subscription, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, eris.Wrapf(err, "failed to create new NewOpenShiftClustersClient")
	}

	var adminCredentials string
	if requestedSecrets.Contains(adminCredentialsKey) {
		var resp armredhatopenshift.OpenShiftClustersClientListAdminCredentialsResponse
		resp, err = clusterClient.ListAdminCredentials(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, eris.Wrapf(err, "failed listing admin credentials")
		}
		adminCredentials = to.Value(resp.Kubeconfig)
	}

	var username string
	var password string
	if requestedSecrets.Contains(usernameKey) || requestedSecrets.Contains(passwordKey) {
		var resp armredhatopenshift.OpenShiftClustersClientListCredentialsResponse
		resp, err = clusterClient.ListCredentials(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, eris.Wrapf(err, "failed listing credentials")
		}
		username = to.Value(resp.KubeadminUsername)
		password = to.Value(resp.KubeadminPassword)
	}

	secretSlice, err := secretsToWrite(typedObj, adminCredentials, username, password)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := map[string]string{}
	if adminCredentials != "" {
		resolvedSecrets[adminCredentialsKey] = adminCredentials
	}
	if username != "" {
		resolvedSecrets[usernameKey] = username
	}
	if password != "" {
		resolvedSecrets[passwordKey] = password
	}
	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *redhatopenshift.OpenShiftCluster) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	secrets := obj.Spec.OperatorSpec.Secrets
	result := set.Set[string]{}
	if secrets.AdminCredentials != nil {
		result.Add(adminCredentialsKey)
	}
	if secrets.Username != nil {
		result.Add(usernameKey)
	}
	if secrets.Password != nil {
		result.Add(passwordKey)
	}

	return result
}

func secretsToWrite(obj *redhatopenshift.OpenShiftCluster, adminCreds string, username string, password string) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.AdminCredentials, adminCreds)
	collector.AddValue(operatorSpecSecrets.Username, username)
	collector.AddValue(operatorSpecSecrets.Password, password)

	return collector.Values()
}
