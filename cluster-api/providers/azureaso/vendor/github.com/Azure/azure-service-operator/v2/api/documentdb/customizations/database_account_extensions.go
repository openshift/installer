/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"strings"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/go-logr/logr"
	"github.com/rotisserie/eris"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	documentdb "github.com/Azure/azure-service-operator/v2/api/documentdb/v1api20240815/storage"
	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/internal/util/to"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

const (
	primaryMasterKeyKey           = "primaryMasterKey"
	secondaryMasterKeyKey         = "secondaryMasterKey"
	primaryReadonlyMasterKeyKey   = "primaryReadonlyMasterKey"
	secondaryReadonlyMasterKeyKey = "secondaryReadonlyMasterKey"
)

var _ genruntime.KubernetesSecretExporter = &DatabaseAccountExtension{}

func (ext *DatabaseAccountExtension) ExportKubernetesSecrets(
	ctx context.Context,
	obj genruntime.MetaObject,
	additionalSecrets set.Set[string],
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
) (*genruntime.KubernetesSecretExportResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*documentdb.DatabaseAccount)
	if !ok {
		return nil, eris.Errorf("cannot run on unknown resource type %T, expected *documentdb.DatabaseAccount", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	primarySecrets, hasEndpoints := secretsSpecified(typedObj)
	requestedSecrets := set.Union(primarySecrets, additionalSecrets)
	if len(requestedSecrets) == 0 && !hasEndpoints {
		log.V(Debug).Info("No secrets retrieval to perform as operatorSpec is empty")
		return nil, nil
	}

	id, err := genruntime.GetAndParseResourceID(typedObj)
	if err != nil {
		return nil, err
	}

	var keys armcosmos.DatabaseAccountListKeysResult
	// Only bother calling ListKeys if there are secrets to retrieve
	if len(requestedSecrets) > 0 {
		subscription := id.SubscriptionID
		// Using armClient.ClientOptions() here ensures we share the same HTTP connection, so this is not opening a new
		// connection each time through
		var acctClient *armcosmos.DatabaseAccountsClient
		acctClient, err = armcosmos.NewDatabaseAccountsClient(subscription, armClient.Creds(), armClient.ClientOptions())
		if err != nil {
			return nil, eris.Wrapf(err, "failed to create new DatabaseAccountClient")
		}

		// TODO: There is a ListReadOnlyKeys API that requires less permissions. We should consider determining
		// TODO: that we don't need to call the ListKeys API and install call the listReadOnlyKeys API.
		var resp armcosmos.DatabaseAccountsClientListKeysResponse
		resp, err = acctClient.ListKeys(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, eris.Wrapf(err, "failed listing keys")
		}

		keys = resp.DatabaseAccountListKeysResult
	}

	resolvedSecrets := map[string]string{}
	if to.Value(keys.PrimaryMasterKey) != "" {
		resolvedSecrets[primaryMasterKeyKey] = to.Value(keys.PrimaryMasterKey)
	}
	if to.Value(keys.SecondaryMasterKey) != "" {
		resolvedSecrets[secondaryMasterKeyKey] = to.Value(keys.SecondaryMasterKey)
	}
	if to.Value(keys.PrimaryReadonlyMasterKey) != "" {
		resolvedSecrets[primaryReadonlyMasterKeyKey] = to.Value(keys.PrimaryReadonlyMasterKey)
	}
	if to.Value(keys.SecondaryReadonlyMasterKey) != "" {
		resolvedSecrets[secondaryReadonlyMasterKeyKey] = to.Value(keys.SecondaryReadonlyMasterKey)
	}

	secretSlice, err := secretsToWrite(typedObj, keys)
	if err != nil {
		return nil, err
	}

	return &genruntime.KubernetesSecretExportResult{
		Objs:       secrets.SliceToClientObjectSlice(secretSlice),
		RawSecrets: secrets.SelectSecrets(additionalSecrets, resolvedSecrets),
	}, nil
}

func secretsSpecified(obj *documentdb.DatabaseAccount) (set.Set[string], bool) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return nil, false
	}

	specSecrets := obj.Spec.OperatorSpec.Secrets
	hasEndpoints := false
	result := make(set.Set[string])
	if specSecrets.PrimaryMasterKey != nil {
		result.Add(primaryMasterKeyKey)
	}
	if specSecrets.SecondaryMasterKey != nil {
		result.Add(secondaryMasterKeyKey)
	}
	if specSecrets.PrimaryReadonlyMasterKey != nil {
		result.Add(primaryReadonlyMasterKeyKey)
	}
	if specSecrets.SecondaryReadonlyMasterKey != nil {
		result.Add(secondaryReadonlyMasterKeyKey)
	}

	if specSecrets.DocumentEndpoint != nil {
		hasEndpoints = true
	}

	return result, hasEndpoints
}

func secretsToWrite(obj *documentdb.DatabaseAccount, accessKeys armcosmos.DatabaseAccountListKeysResult) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, nil
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.PrimaryMasterKey, to.Value(accessKeys.PrimaryMasterKey))
	collector.AddValue(operatorSpecSecrets.SecondaryMasterKey, to.Value(accessKeys.SecondaryMasterKey))
	collector.AddValue(operatorSpecSecrets.PrimaryReadonlyMasterKey, to.Value(accessKeys.PrimaryReadonlyMasterKey))
	collector.AddValue(operatorSpecSecrets.SecondaryReadonlyMasterKey, to.Value(accessKeys.SecondaryReadonlyMasterKey))
	collector.AddValue(operatorSpecSecrets.DocumentEndpoint, to.Value(obj.Status.DocumentEndpoint))

	return collector.Values()
}

var _ extensions.ErrorClassifier = &DatabaseAccountExtension{}

// ClassifyError evaluates the provided error, returning whether it is fatal or can be retried.
// A "ServiceUnavailable" error would usually be retryable, but in the case of DatabaseAccount, when coupled with
// a "high demand" message, it means that the region is capacity constrained and cannot have DatabseAccounts allocated.
// If we retry on this error CosmosDB will start returning a new BadRequest error stating
// DatabaseAccount is in a failed provisioning state because the previous attempt to create it was not successful.
// Please delete the previous instance before attempting to recreate this account."
// Since we can't retry anyway, we mark the original ServiceUnavailable error as fatal so the user has a clearer error message.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (ext *DatabaseAccountExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	if isCapacityError(cloudError) {
		details.Classification = core.ErrorFatal
	}

	return details, nil
}

func isCapacityError(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	return err.Code() == "ServiceUnavailable" && strings.Contains(err.Message(), "currently experiencing high demand")
}

var _ extensions.PreReconciliationChecker = &DatabaseAccountExtension{}

// PreReconcileCheck does a pre-reconcile check to see if the resource is in a state that can be reconciled.
// ARM resources should implement this to avoid reconciliation attempts that cannot possibly succeed.
// Returns ProceedWithReconcile if the reconciliation should go ahead.
// Returns BlockReconcile and a human-readable reason if the reconciliation should be skipped.
// ctx is the current operation context.
// obj is the resource about to be reconciled. The resource's State will be freshly updated.
// kubeClient allows access to the cluster for any required queries.
// armClient allows access to ARM for any required queries.
// log is the logger for the current operation.
// next is the next (nested) implementation to call.
func (ext *DatabaseAccountExtension) PreReconcileCheck(
	ctx context.Context,
	obj genruntime.MetaObject,
	owner genruntime.MetaObject,
	resourceResolver *resolver.Resolver,
	armClient *genericarmclient.GenericClient,
	log logr.Logger,
	next extensions.PreReconcileCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// This has to be the current hub storage version of the account.
	// It will need to be updated if the hub storage version changes.
	account, ok := obj.(*documentdb.DatabaseAccount)
	if !ok {
		return extensions.PreReconcileCheckResult{}, eris.Errorf("cannot run on unknown resource type %T, expected *documentdb.DatabaseAccount", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = account

	// If the account is already deleting, we have to wait for that to finish
	// before trying anything else
	if account.Status.ProvisioningState != nil && strings.EqualFold(*account.Status.ProvisioningState, "Deleting") {
		return extensions.BlockReconcile("reconcile blocked while account is at status deleting"), nil
	}

	return next(ctx, obj, owner, resourceResolver, armClient, log)
}
