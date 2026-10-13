package ibmcloud

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const resourceGroupTypeName = "resource group"

// listResourceGroups lists resource groups in the account
func (o *ClusterUninstaller) listResourceGroups() (cloudResources, error) {
	o.Logger.Debugf("Listing resource groups")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.managementSvc.NewListResourceGroupsOptions()
	options.SetAccountID(o.AccountID)
	resources, _, err := o.managementSvc.ListResourceGroupsWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list resource groups")
	}

	result := []cloudResource{}
	for _, resourceGroup := range resources.Resources {
		if *resourceGroup.Name == o.ResourceGroupName {
			result = append(result, cloudResource{
				key:      *resourceGroup.ID,
				name:     *resourceGroup.Name,
				status:   "",
				typeName: resourceGroupTypeName,
				id:       *resourceGroup.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) getBlockingServiceInstances(item cloudResource) ([]string, error) {
	o.Logger.Debugf("Collecting remaining service-instances in %s", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	// Find the Resource Group Id first
	listRGOptions := o.managementSvc.NewListResourceGroupsOptions()
	listRGOptions = listRGOptions.SetName(item.name)
	listRGOptions = listRGOptions.SetAccountID(o.AccountID)

	result, response, err := o.managementSvc.ListResourceGroupsWithContext(ctx, listRGOptions)
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			// The resource group is gone
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list existing resource groups: %w", err)
	}

	if len(result.Resources) != 1 {
		return nil, fmt.Errorf("unexpected number of Resource Groups found for %s: %d", item.name, len(result.Resources))
	}
	resourceGroupID := *result.Resources[0].ID

	// Get existing service-instances for Resource Group
	listRIOptions := o.controllerSvc.NewListResourceInstancesOptions()
	listRIOptions = listRIOptions.SetResourceGroupID(resourceGroupID)

	instancesResult, instancesResponse, err := o.controllerSvc.ListResourceInstancesWithContext(ctx, listRIOptions)
	if err != nil {
		if instancesResponse != nil && instancesResponse.StatusCode == http.StatusNotFound {
			// No service-instances or resource group could not be found; potentially ready to be or being deleted
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list service-instances in %s: %w", item.name, err)
	}

	if len(instancesResult.Resources) < 1 {
		return nil, fmt.Errorf("no service instances found for %s: unknown failure deleting Resource Group", item.name)
	}

	instanceCRNS := []string{}

	for _, instance := range instancesResult.Resources {
		instanceCRNS = append(instanceCRNS, *instance.CRN)
	}
	return instanceCRNS, nil
}

func (o *ClusterUninstaller) deleteResourceGroup(item cloudResource) error {
	o.Logger.Debugf("Deleting resource group %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.managementSvc.NewDeleteResourceGroupOptions(item.id)
	details, err := o.managementSvc.DeleteResourceGroupWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted resource group %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		serviceInstances, siErr := o.getBlockingServiceInstances(item)
		if err != nil { // nolint: gocritic
			return fmt.Errorf("failed attempting to identify instances blocking resource group deletion: %s: %w", item.name, siErr)
		} else if serviceInstances != nil {
			o.Logger.Warnf("resource group contains existing service instances: %s", serviceInstances)
			return fmt.Errorf("failed to delete resource group %s due to existing service instances", item.name)
		}
		return fmt.Errorf("failed to delete resource group, another attempt may be successful: %s", item.name)
	}

	return nil
}

// destroyResourceGroups removes the installer-generated resource group. If the
// resource group is user-provided, it will not be removed.
func (o *ClusterUninstaller) destroyResourceGroups() error {
	if o.ResourceGroupName != o.InfraID {
		o.Logger.Infof("Skipping deletion of user-provided resource group %v", o.ResourceGroupName)
		return nil
	}

	found, err := o.listResourceGroups()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(resourceGroupTypeName, found.list())

	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted resource group %q", item.name)
			continue
		}
		err = o.deleteResourceGroup(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(resourceGroupTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
