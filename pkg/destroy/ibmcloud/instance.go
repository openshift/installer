package ibmcloud

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const (
	instanceTypeName       = "instance"
	instanceActionTypeName = "instance action"
)

func (o *ClusterUninstaller) listInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing virtual service instances")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListInstancesOptions()
	options.SetVPCName(fmt.Sprintf("%s-vpc", o.InfraID))
	resources, _, err := o.vpcSvc.ListInstancesWithContext(ctx, options)
	if err != nil {
		return nil, err
	}

	result := []cloudResource{}
	for _, instance := range resources.Instances {
		if strings.Contains(*instance.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *instance.ID,
				name:     *instance.Name,
				status:   *instance.Status,
				typeName: "instance",
				id:       *instance.ID,
			})
		}
	}
	return cloudResources{}.insert(result...), nil
}

// stopInstances searches for instances in the vpc that have a name that starts with
// the infra ID prefix and are not yet stopped. It then stops each instance found.
func (o *ClusterUninstaller) stopInstances() error {
	found, err := o.listInstances()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(instanceActionTypeName, found.list())
	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has been removed.
			o.deletePendingItems(instanceActionTypeName, []cloudResource{item})
			o.Logger.Warnf("Instance %q not found", item.name)
			continue
		}
		err := o.stopInstance(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(instanceActionTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}

func (o *ClusterUninstaller) stopInstance(item cloudResource) error {
	if item.status == vpcv1.InstanceStatusStoppingConst {
		o.Logger.Debugf("Waiting for instance %q to stop", item.name)
		return nil
	}

	if item.status == vpcv1.InstanceStatusStoppedConst {
		o.Logger.Infof("Stopped instance %q", item.name)
		o.deletePendingItems(instanceActionTypeName, []cloudResource{item})
		return nil
	}

	o.Logger.Debugf("Stopping instance %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := &vpcv1.CreateInstanceActionOptions{
		InstanceID: core.StringPtr(item.id),
		Type:       core.StringPtr(vpcv1.CreateInstanceActionOptionsTypeStopConst),
	}

	_, _, err := o.vpcSvc.CreateInstanceActionWithContext(ctx, options)
	if err != nil {
		return errors.Wrapf(err, "Failed to stop instance %q", item.name)
	}

	return nil
}

func (o *ClusterUninstaller) deleteInstance(item cloudResource) error {
	if item.status == vpcv1.InstanceStatusDeletingConst {
		o.Logger.Debugf("Waiting for instance %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting instance %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteInstanceOptions(item.id)
	details, err := o.vpcSvc.DeleteInstanceWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted instance %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete instance %q", item.name)
	}

	return nil
}

// destroyInstances searches for instances that have a name that starts with
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyInstances() error {
	found, err := o.listInstances()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(instanceTypeName, found.list())
	for _, item := range items {
		if _, ok := found[item.key]; !ok {
			// This item has finished deletion.
			o.deletePendingItems(item.typeName, []cloudResource{item})
			o.Logger.Infof("Deleted instance %q", item.name)
			continue
		}
		err := o.deleteInstance(item)
		if err != nil {
			o.errorTracker.suppressWarning(item.key, err, o.Logger)
		}
	}

	if items = o.getPendingItems(instanceTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
