package ibmcloud

import (
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
		return errors.Wrapf(err, "Failed to delete resource group %s", item.name)
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
