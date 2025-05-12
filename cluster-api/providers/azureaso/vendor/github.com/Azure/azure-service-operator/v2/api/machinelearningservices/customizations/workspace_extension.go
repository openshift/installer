/*
* Copyright (c) Microsoft Corporation.
* Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/Azure/azure-service-operator/v2/api/machinelearningservices/v1api20240401/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	appInsightsInstrumentationKey = "appInsightsInstrumentationKey"
	primaryNotebookAccessKey      = "primaryNotebookAccessKey"
	secondaryNotebookAccessKey    = "secondaryNotebookAccessKey"
	userStorageKey                = "userStorageKey"
	containerRegistryPassword     = "containerRegistryPassword"
	containerRegistryPassword2    = "containerRegistryPassword2"
	containerRegistryUserName     = "containerRegistryUserName"
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
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *storage.Workspace", obj)
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

	var keys armmachinelearning.ListWorkspaceKeysResult
	// Only bother calling ListKeys if there are secrets to retrieve
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var workspacesClient *armmachinelearning.WorkspacesClient
		workspacesClient, err = armmachinelearning.NewWorkspacesClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new workspaceClient")
		}

		var resp armmachinelearning.WorkspacesClientListKeysResponse
		resp, err = workspacesClient.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing keys")
		}
		keys = resp.ListWorkspaceKeysResult
	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	resolvedSecrets := makeResolvedSecrets(keys)
	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *storage.Workspace) set.Set[string] {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil
	}

	operatorSecrets := obj.Spec.OperatorSpec.Secrets

	result := make(set.Set[string])
	if operatorSecrets.AppInsightsInstrumentationKey != nil {
		result.Add(appInsightsInstrumentationKey)
	}
	if operatorSecrets.PrimaryNotebookAccessKey != nil {
		result.Add(primaryNotebookAccessKey)
	}
	if operatorSecrets.SecondaryNotebookAccessKey != nil {
		result.Add(secondaryNotebookAccessKey)
	}
	if operatorSecrets.UserStorageKey != nil {
		result.Add(userStorageKey)
	}
	if operatorSecrets.ContainerRegistryPassword != nil {
		result.Add(containerRegistryPassword)
	}
	if operatorSecrets.ContainerRegistryPassword2 != nil {
		result.Add(containerRegistryPassword2)
	}
	if operatorSecrets.ContainerRegistryUserName != nil {
		result.Add(containerRegistryUserName)
	}

	return result
}

func secretsToWrite(obj *storage.Workspace, keysResp armmachinelearning.ListWorkspaceKeysResult) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)

	creds, crUsername := getContainerRegCreds(keysResp)

	collector.AddValue(operatorSpecSecrets.ContainerRegistryPassword, creds["password"])
	collector.AddValue(operatorSpecSecrets.ContainerRegistryPassword2, creds["password2"])
	collector.AddValue(operatorSpecSecrets.ContainerRegistryUserName, crUsername)
	collector.AddValue(operatorSpecSecrets.UserStorageKey, to.Value(keysResp.UserStorageKey))
	collector.AddValue(operatorSpecSecrets.AppInsightsInstrumentationKey, to.Value(keysResp.AppInsightsInstrumentationKey))

	if keysResp.NotebookAccessKeys != nil {
		collector.AddValue(operatorSpecSecrets.PrimaryNotebookAccessKey, to.Value(keysResp.NotebookAccessKeys.PrimaryAccessKey))
		collector.AddValue(operatorSpecSecrets.SecondaryNotebookAccessKey, to.Value(keysResp.NotebookAccessKeys.SecondaryAccessKey))
	}

	return collector.Values()
}

func makeResolvedSecrets(keys armmachinelearning.ListWorkspaceKeysResult) map[string]string {
	resolvedSecrets := map[string]string{}
	if to.Value(keys.AppInsightsInstrumentationKey) != "" {
		resolvedSecrets[appInsightsInstrumentationKey] = to.Value(keys.AppInsightsInstrumentationKey)
	}
	if keys.NotebookAccessKeys != nil {
		if to.Value(keys.NotebookAccessKeys.PrimaryAccessKey) != "" {
			resolvedSecrets[primaryNotebookAccessKey] = to.Value(keys.NotebookAccessKeys.PrimaryAccessKey)
		}
		if to.Value(keys.NotebookAccessKeys.SecondaryAccessKey) != "" {
			resolvedSecrets[secondaryNotebookAccessKey] = to.Value(keys.NotebookAccessKeys.SecondaryAccessKey)
		}
	}
	creds, crUsername := getContainerRegCreds(keys)
	if creds["password"] != "" {
		resolvedSecrets[containerRegistryPassword] = creds["password"]
	}
	if creds["password2"] != "" {
		resolvedSecrets[containerRegistryPassword2] = creds["password2"]
	}
	if crUsername != "" {
		resolvedSecrets[containerRegistryUserName] = crUsername
	}

	if to.Value(keys.UserStorageKey) != "" {
		resolvedSecrets[userStorageKey] = to.Value(keys.UserStorageKey)
	}

	return resolvedSecrets
}

func getContainerRegCreds(keysResp armmachinelearning.ListWorkspaceKeysResult) (map[string]string, string) {
	creds := make(map[string]string)

	if keysResp.ContainerRegistryCredentials == nil {
		return creds, ""
	}

	for _, password := range keysResp.ContainerRegistryCredentials.Passwords {
		if password.Name != nil && password.Value != nil {
			creds[to.Value(password.Name)] = to.Value(password.Value)
		}
	}
	return creds, to.Value(keysResp.ContainerRegistryCredentials.Username)
}

var _ extensions.ErrorClassifier = &WorkspaceExtension{}

func (ext *WorkspaceExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	// Override is to treat StorageAccountIsNotProvisioned as retryable for Workspaces
	// (as we may be waiting for the Storage Account to be provisioned)
	if isRetryableWorkspaceError(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

func isRetryableWorkspaceError(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	if err.Code() == "BadRequest" && strings.Contains(err.Message(), "StorageAccountIsNotProvisioned") {
		// This is a retryable error
		return true
	}

	return false
}
