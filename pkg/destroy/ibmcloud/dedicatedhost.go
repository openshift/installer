package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const (
	dedicatedHostTypeName      = "dedicated host"
	dedicatedHostGroupTypeName = "dedicated host group"
)

// listDedicatedHosts searches for dedicated host that have a name that
// starts with the cluster's infra ID.
func (o *ClusterUninstaller) listDedicatedHosts() (cloudResources, error) {
	o.Logger.Debugf("Listing dedicated hosts")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	resourceGroupID, err := o.ResourceGroupID()
	if err != nil {
		return nil, err
	}

	options := o.vpcSvc.NewListDedicatedHostsOptions()
	options.SetResourceGroupID(resourceGroupID)
	resources, _, err := o.vpcSvc.ListDedicatedHostsWithContext(ctx, options)
	if err != nil {
		return nil, err
	}

	result := []cloudResource{}
	for _, dhost := range resources.DedicatedHosts {
		if strings.HasPrefix(*dhost.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *dhost.ID,
				name:     *dhost.Name,
				status:   *dhost.State,
				typeName: dedicatedHostTypeName,
				id:       *dhost.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

// listDedicatedHostGroups searches for dedicated host groups that have a name
// that starts with the cluster's infra ID.
func (o *ClusterUninstaller) listDedicatedHostGroups() (cloudResources, error) {
	o.Logger.Debugf("Listing dedicated host groups")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	resourceGroupID, err := o.ResourceGroupID()
	if err != nil {
		return nil, err
	}

	options := o.vpcSvc.NewListDedicatedHostGroupsOptions()
	options.SetResourceGroupID(resourceGroupID)
	resources, _, err := o.vpcSvc.ListDedicatedHostGroupsWithContext(ctx, options)
	if err != nil {
		return nil, err
	}

	result := []cloudResource{}
	for _, dgroup := range resources.Groups {
		if strings.HasPrefix(*dgroup.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *dgroup.ID,
				name:     *dgroup.Name,
				status:   "",
				typeName: dedicatedHostGroupTypeName,
				id:       *dgroup.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

// deleteDedicatedHosts disables and deletes a dedicated host
func (o *ClusterUninstaller) deleteDedicatedHost(item cloudResource) error {
	if item.status == vpcv1.DedicatedHostLifecycleStateDeletingConst {
		o.Logger.Debugf("Waiting for dedicated host %q to delete", item.name)
		return nil
	}

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	o.Logger.Debugf("Getting dedicated host %q", item.name)

	getOpts := o.vpcSvc.NewGetDedicatedHostOptions(item.id)
	resource, _, err := o.vpcSvc.GetDedicatedHostWithContext(ctx, getOpts)
	if err != nil {
		return errors.Wrapf(err, "Failed to get dedicated host %q", item.name)
	}

	if *resource.InstancePlacementEnabled {
		if err := o.disableDedicatedHost(item); err != nil {
			return err
		}
	}

	o.Logger.Debugf("Deleting dedicated host %q", item.name)

	deleteOpts := o.vpcSvc.NewDeleteDedicatedHostOptions(item.id)
	details, err := o.vpcSvc.DeleteDedicatedHostWithContext(ctx, deleteOpts)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted dedicated host %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete dedicated host %q", item.name)
	}

	return nil
}

// deleteDedicatedHostGroups deletes a dedicated host group
func (o *ClusterUninstaller) deleteDedicatedHostGroup(item cloudResource) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	o.Logger.Debugf("Deleting dedicated host group %q", item.name)

	options := o.vpcSvc.NewDeleteDedicatedHostGroupOptions(item.id)
	details, err := o.vpcSvc.DeleteDedicatedHostGroupWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted dedicated host group %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete dedicated host group %q", item.name)
	}

	return nil
}

// disableDedicatedHosts disables a dedicated host
func (o *ClusterUninstaller) disableDedicatedHost(item cloudResource) error {
	o.Logger.Debugf("Disabling dedicated host %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewUpdateDedicatedHostOptions(item.id, map[string]interface{}{
		"instance_placement_enabled": core.BoolPtr(false),
	})
	_, _, err := o.vpcSvc.UpdateDedicatedHostWithContext(ctx, options)
	if err != nil {
		return errors.Wrapf(err, "Failed to disable dedicated host %q", item.name)
	}

	return nil
}

// destroyDedicatedHosts searches for dedicated hosts that have a name
// that starts with the cluster's infra ID.
func (o *ClusterUninstaller) destroyDedicatedHosts() error {
	found, err := o.listDedicatedHosts()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(dedicatedHostTypeName, found.list())
	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted dedicated host %q", item.name)
			continue
		}
		err := o.deleteDedicatedHost(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(dedicatedHostTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}

// destroyDedicatedHostGroups searches for dedicated host groups that have
// a name that starts with the cluster's infra ID.
func (o *ClusterUninstaller) destroyDedicatedHostGroups() error {
	found, err := o.listDedicatedHostGroups()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(dedicatedHostGroupTypeName, found.list())
	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted dedicated host group %q", item.name)
			continue
		}
		err := o.deleteDedicatedHostGroup(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(dedicatedHostGroupTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
