/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package customizations

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	containerservice "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001/storage"
	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/set"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
)

var _ genruntime.KubernetesExporter = &ManagedClusterExtension{}

func (ext *ManagedClusterExtension) ExportKubernetesResources(
	ctx context.Context,
	obj genruntime.MetaObject,
	armClient *genericarmclient.GenericClient,
	log logr.Logger) ([]client.Object, error) {

	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	typedObj, ok := obj.(*containerservice.ManagedCluster)
	if !ok {
		return nil, errors.Errorf("cannot run on unknown resource type %T, expected *containerservice.ManagedCluster", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = typedObj

	hasAdminCreds, hasUserCreds := secretsSpecified(typedObj)
	if !hasAdminCreds && !hasUserCreds {
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
	var mcClient *armcontainerservice.ManagedClustersClient
	mcClient, err = armcontainerservice.NewManagedClustersClient(subscription, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new ManagedClustersClient")
	}

	// TODO: In the future we may need variants of these secret properties that configure usage of the public FQDN rather than the private one, see:
	// TODO: https://docs.microsoft.com/en-us/answers/questions/670332/azure-aks-get-credentials-using-wrong-hostname-for.html
	var adminCredentials string
	if hasAdminCreds {
		var resp armcontainerservice.ManagedClustersClientListClusterAdminCredentialsResponse
		resp, err = mcClient.ListClusterAdminCredentials(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing admin credentials")
		}
		if len(resp.CredentialResults.Kubeconfigs) > 0 {
			// It's awkward that we're ignoring the other possible responses here, but that's what the AZ CLI does too:
			// https://github.com/Azure/azure-cli/blob/6786b5014ae71eb6d93f95e1ad123e9171368e8f/src/azure-cli/azure/cli/command_modules/acs/custom.py#L2166
			adminCredentials = string(resp.CredentialResults.Kubeconfigs[0].Value)
		}
	}

	var userCredentials string
	if hasUserCreds {
		var resp armcontainerservice.ManagedClustersClientListClusterUserCredentialsResponse
		resp, err = mcClient.ListClusterUserCredentials(ctx, id.ResourceGroupName, typedObj.AzureName(), nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed listing admin credentials")
		}
		if len(resp.CredentialResults.Kubeconfigs) > 0 {
			// It's awkward that we're ignoring the other possible responses here, but that's what the AZ CLI does too:
			// https://github.com/Azure/azure-cli/blob/6786b5014ae71eb6d93f95e1ad123e9171368e8f/src/azure-cli/azure/cli/command_modules/acs/custom.py#L2166
			userCredentials = string(resp.CredentialResults.Kubeconfigs[0].Value)
		}
	}

	secretSlice, err := secretsToWrite(typedObj, adminCredentials, userCredentials)
	if err != nil {
		return nil, err
	}

	return secrets.SliceToClientObjectSlice(secretSlice), nil
}

func secretsSpecified(obj *containerservice.ManagedCluster) (bool, bool) {
	if obj.Spec.OperatorSpec == nil || obj.Spec.OperatorSpec.Secrets == nil {
		return false, false
	}

	secrets := obj.Spec.OperatorSpec.Secrets
	return secrets.AdminCredentials != nil, secrets.UserCredentials != nil
}

func secretsToWrite(obj *containerservice.ManagedCluster, adminCreds string, userCreds string) ([]*v1.Secret, error) {
	operatorSpecSecrets := obj.Spec.OperatorSpec.Secrets
	if operatorSpecSecrets == nil {
		return nil, errors.Errorf("unexpected nil operatorspec")
	}

	collector := secrets.NewCollector(obj.Namespace)
	collector.AddValue(operatorSpecSecrets.AdminCredentials, adminCreds)
	collector.AddValue(operatorSpecSecrets.UserCredentials, userCreds)

	return collector.Values()
}

var _ extensions.PreReconciliationChecker = &ManagedClusterExtension{}

// If a managed cluster has a provisioningState not in this set, it will reject any attempt to PUT a new state out of
// hand; so there's no point in even trying. This is true even if the PUT we're doing will have no effect on the state
// of the cluster.
// These are all listed lowercase, so we can do a case-insensitive match.
var nonBlockingManagedClusterProvisioningStates = set.Make(
	"succeeded",
	"failed",
	"canceled",
)

func (ext *ManagedClusterExtension) PreReconcileCheck(
	_ context.Context,
	obj genruntime.MetaObject,
	_ genruntime.MetaObject,
	_ *resolver.Resolver,
	_ *genericarmclient.GenericClient,
	_ logr.Logger,
	_ extensions.PreReconcileCheckFunc,
) (extensions.PreReconcileCheckResult, error) {
	// This has to be the current hub storage version. It will need to be updated
	// if the hub storage version changes.
	managedCluster, ok := obj.(*containerservice.ManagedCluster)
	if !ok {
		return extensions.PreReconcileCheckResult{},
			errors.Errorf("cannot run on unknown resource type %T, expected *containerservice.ManagedCluster", obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not
	var _ conversion.Hub = managedCluster

	// If the cluster is in a state that will reject any PUT, then we should skip reconciliation
	// as there's no point in even trying.
	// This allows us to "play nice with others" and not use up request quota attempting to make changes when we
	// already know those attempts will fail.
	state := managedCluster.Status.ProvisioningState
	if state != nil && clusterProvisioningStateBlocksReconciliation(state) {
		return extensions.BlockReconcile(
				fmt.Sprintf("Managed cluster is in provisioning state %q", *state)),
			nil
	}

	return extensions.ProceedWithReconcile(), nil
}

func clusterProvisioningStateBlocksReconciliation(provisioningState *string) bool {
	if provisioningState == nil {
		return false
	}

	return !nonBlockingManagedClusterProvisioningStates.Contains(strings.ToLower(*provisioningState))
}

var _ extensions.ErrorClassifier = &ManagedClusterExtension{}

// ClassifyError evaluates the provided error, returning including whether it is fatal or can be retried.
// cloudError is the error returned from ARM.
// apiVersion is the ARM API version used for the request.
// log is a logger than can be used for telemetry.
// next is the next implementation to call.
func (ext *ManagedClusterExtension) ClassifyError(
	cloudError *genericarmclient.CloudError,
	apiVersion string,
	log logr.Logger,
	next extensions.ErrorClassifierFunc,
) (core.CloudErrorDetails, error) {
	details, err := next(cloudError)
	if err != nil {
		return core.CloudErrorDetails{}, err
	}

	if isRetryableClusterError(cloudError) {
		details.Classification = core.ErrorRetryable
	}

	return details, nil
}

func isRetryableClusterError(err *genericarmclient.CloudError) bool {
	if err == nil {
		return false
	}

	// A CustomKubeletIdentityMissingPermissionError can occur if the user-assigned identity required by the cluster
	// hasn't yet been provisioned; we want to retry so that we finish provisioning the cluster once it is available.
	if err.Code() == "CustomKubeletIdentityMissingPermissionError" {
		return true
	}

	return false
}
