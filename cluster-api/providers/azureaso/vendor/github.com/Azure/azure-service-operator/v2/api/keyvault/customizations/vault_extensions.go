/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package customizations

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	keyvault "github.com/Azure/azure-service-operator/v2/api/keyvault/v1api20230701/storage"
	resources "github.com/Azure/azure-service-operator/v2/api/resources/v1api20200601/storage"

	"github.com/Azure/azure-service-operator/v2/internal/genericarmclient"
	"github.com/Azure/azure-service-operator/v2/internal/reflecthelpers"
	"github.com/Azure/azure-service-operator/v2/internal/resolver"
	"github.com/Azure/azure-service-operator/v2/internal/util/kubeclient"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/extensions"
)

var _ extensions.ARMResourceModifier = &VaultExtension{}

const (
	CreateMode_Default         = "default"
	CreateMode_Recover         = "recover"
	CreateMode_CreateOrRecover = "createOrRecover"
	CreateMode_PurgeThenCreate = "purgeThenCreate"
)

// ModifyARMResource implements extensions.ARMResourceModifier.
func (ex *VaultExtension) ModifyARMResource(
	ctx context.Context,
	armClient *genericarmclient.GenericClient,
	armObj genruntime.ARMResource,
	obj genruntime.ARMMetaObject,
	kubeClient kubeclient.Client,
	resolver *resolver.Resolver,
	log logr.Logger,
) (genruntime.ARMResource, error) {

	kv, ok := obj.(*keyvault.Vault)
	if !ok {
		return nil, errors.Errorf(
			"Cannot run VaultExtension.ModifyARMResource() with unexpected resource type %T",
			obj)
	}

	// Type assert that we are the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated to match
	var _ conversion.Hub = kv

	// If createMode is nil, nothing for us to do
	// (This shouldn't be possible, but better to hedge against it)
	if kv.Spec.Properties == nil || kv.Spec.Properties.CreateMode == nil {
		return armObj, nil
	}

	// Get the owner of the KeyVault, we need this resource group to determine the subscription
	owner, ownerErr := ex.getOwner(ctx, kv, resolver, log)
	if ownerErr != nil {
		return nil, errors.Wrapf(ownerErr, "unable to find owner of KeyVault %s", kv.Name)
	}

	// Parse the ID of the owner
	// (Can't use the KeyVault as we do this before the KV exists)
	id, err := genruntime.GetAndParseResourceID(owner)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get and parse resource ID from KeyVault owner")
	}

	vc, err := armkeyvault.NewVaultsClient(id.SubscriptionID, armClient.Creds(), armClient.ClientOptions())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new VaultsClient")
	}

	createMode := *kv.Spec.Properties.CreateMode
	if createMode == CreateMode_CreateOrRecover {
		createMode, err = ex.handleCreateOrRecover(ctx, kv, vc, resolver, log)
		if err != nil {
			return nil, errors.Wrapf(err, "error checking for existence of soft-deleted KeyVault")
		}
	}

	if createMode == CreateMode_PurgeThenCreate {
		err = ex.handlePurgeThenCreate(ctx, kv, vc, resolver, log)
		if err != nil {
			return nil, errors.Wrapf(err, "error purging soft-deleted KeyVault")
		}

		createMode = CreateMode_Default
	}

	// Modify the payload as necessary
	spec := armObj.Spec()
	err = reflecthelpers.SetProperty(spec, "Properties.CreateMode", &createMode)
	if err != nil {
		return nil, errors.Wrapf(err, "error setting CreateMode to %s", createMode)
	}

	return armObj, nil
}

func (ex *VaultExtension) handleCreateOrRecover(
	ctx context.Context,
	kv *keyvault.Vault,
	vc *armkeyvault.VaultsClient,
	resolver *resolver.Resolver,
	log logr.Logger,
) (string, error) {
	exists, err := ex.checkForExistenceOfDeletedKeyVault(ctx, kv, resolver, vc, log)
	if err != nil {
		return "", errors.Wrapf(err, "error checking for existence of soft-deleted KeyVault %s", kv.Name)
	}

	result := CreateMode_Default
	if exists {
		result = CreateMode_Recover
	}

	log.Info(
		"KeyVault reconciliation requested CreateOrRecover",
		"KeyVault", kv.Name,
		"softDeletedKeyvaultExists", exists,
		"createMode", result)

	return result, err
}

func (ex *VaultExtension) handlePurgeThenCreate(
	ctx context.Context,
	kv *keyvault.Vault,
	vc *armkeyvault.VaultsClient,
	resolver *resolver.Resolver,
	log logr.Logger,
) error {
	// Find out whether a soft-deleted KeyVault with the same name exists
	exists, err := ex.checkForExistenceOfDeletedKeyVault(ctx, kv, resolver, vc, log)
	if err != nil {
		// Could not determine whether a soft-deleted keyvault exists in the same subscription, assume it doesn't

		log.Error(err, "error checking for existence of soft-deleted KeyVault")
		return nil
	}

	log.Info(
		"KeyVault reconciliation requested PurgeThenCreate",
		"KeyVault", kv.Name,
		"softDeletedKeyVaultExists", exists)

	if exists {
		// Get the owner of the KeyVault, we need this resource group to determine the location
		owner, ownerErr := ex.getOwner(ctx, kv, resolver, log)
		if ownerErr != nil {
			return errors.Wrapf(ownerErr, "unable to find owner of KeyVault %s", kv.Name)
		}

		// if a soft-deleted KeyVault exists, we need to purge it before we can create a new one
		// Get the location of the KeyVault
		location, locationOk := ex.getLocation(kv, owner)
		if !locationOk {
			return errors.Errorf("unable to determine location of KeyVault %s", kv.Name)
		}

		poller, err := vc.BeginPurgeDeleted(ctx, kv.Name, location, &armkeyvault.VaultsClientBeginPurgeDeletedOptions{})
		if err != nil {
			return errors.Wrapf(err, "failed to begin purging deleted KeyVault %s", kv.Name)
		}

		_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: 10 * time.Second})
		if err != nil {
			return errors.Wrapf(err, "failed to purge deleted KeyVault %s", kv.Name)
		}
	}

	return nil
}

// checkForExistenceOfDeletedKeyVault checks to see whether there's a soft deleted KeyVault with the same name.
// This might be true if another party has deleted the KeyVault, even if we previously created it
func (ex *VaultExtension) checkForExistenceOfDeletedKeyVault(
	ctx context.Context,
	kv *keyvault.Vault,
	resolver *resolver.Resolver,
	vaultsClient *armkeyvault.VaultsClient,
	log logr.Logger,
) (bool, error) {
	// Get the owner of the KeyVault, we need this resource group to determine the subscription
	owner, ownerErr := ex.getOwner(ctx, kv, resolver, log)
	if ownerErr != nil {
		return false, errors.Wrapf(ownerErr, "unable to find owner of KeyVault %s", kv.Name)
	}

	// Get the location of the KeyVault
	location, locationOk := ex.getLocation(kv, owner)
	if !locationOk {
		return false, errors.Errorf("unable to determine location of KeyVault %s", kv.Name)
	}

	// Get the name of the KeyVault
	vaultName := kv.Spec.AzureName
	if vaultName == "" {
		vaultName = kv.Name
	}

	// Default to assuming a soft-deleted keyvault exists
	exists := true

	// Check to see if this is true
	_, err := vaultsClient.GetDeleted(ctx, vaultName, location, &armkeyvault.VaultsClientGetDeletedOptions{})
	if err != nil {
		var responseError *azcore.ResponseError
		if errors.As(err, &responseError) {
			if responseError.StatusCode != http.StatusNotFound {
				return false, errors.Wrapf(err, "failed to get deleted KeyVault %s, error %d", kv.Name, responseError.StatusCode)
			}

			// KeyVault doesn't exist,
			exists = false
		}
	}

	log.Info(
		"Checking for existence of soft-deleted KeyVault",
		"keyVault", kv.Name,
		"location", location,
		"softDeletedKeyvaultExists", exists,
	)

	return exists, nil
}

func (*VaultExtension) getOwner(
	ctx context.Context,
	kv *keyvault.Vault,
	resolver *resolver.Resolver,
	log logr.Logger,
) (*resources.ResourceGroup, error) {
	owner, err := resolver.ResolveOwner(ctx, kv)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to resolve owner of KeyVault %s", kv.Name)
	}

	// No need to wait for resources that don't have an owner
	if !owner.FoundKubernetesOwner() {
		log.Info(
			"KeyVault owner is not within the cluster, cannot determine subscription",
			"keyVault", kv.Name)
		return nil, errors.Errorf("owner of KeyVault %s is not within the cluster", kv.Name)
	}

	rg, ok := owner.Owner.(*resources.ResourceGroup)
	if !ok {
		return nil, errors.Errorf("expected owner of KeyVault %s to be a ResourceGroup", kv.Name)
	}

	// Type assert that the ResourceGroup is the hub type. This will fail to compile if
	// the hub type has been changed but this extension has not been updated to match
	var _ conversion.Hub = rg

	return rg, nil
}

// findKeyVaultLocation determines which location we're trying to create KeyVault within
func (*VaultExtension) getLocation(
	kv *keyvault.Vault,
	rg *resources.ResourceGroup,
) (string, bool) {
	// Prefer location on the KeyVault
	if kv.Spec.Location != nil && *kv.Spec.Location != "" {
		return *kv.Spec.Location, true
	}

	// Fallback to location on ResourceGroup
	if rg.Spec.Location != nil && *rg.Spec.Location != "" {
		return *rg.Spec.Location, true
	}

	return "", false
}
