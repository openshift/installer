package powervs

import (
	"strings"

	"github.com/pkg/errors"
)

const (
	instanceTypeName       = "instance"
	instanceActionTypeName = "instance action"
)

// listInstances lists instances in the vpc.
func (o *ClusterUninstaller) listInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing virtual service instances")

	instances, err := o.instanceClient.GetAll()
	if err != nil {
		o.Logger.Warnf("Error instanceClient.GetAll: %v", err)
		return nil, err
	}

	var foundOne = false

	result := []cloudResource{}
	for _, instance := range instances.PvmInstances {
		// https://github.com/IBM-Cloud/power-go-client/blob/master/power/models/p_vm_instance.go
		if strings.Contains(*instance.ServerName, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listInstances: FOUND: %s, %s, %s", *instance.PvmInstanceID, *instance.ServerName, *instance.Status)
			result = append(result, cloudResource{
				key:      *instance.PvmInstanceID,
				name:     *instance.ServerName,
				status:   *instance.Status,
				typeName: "instance",
				id:       *instance.PvmInstanceID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listInstances: NO matching virtual instance against: %s", o.InfraID)
		for _, instance := range instances.PvmInstances {
			o.Logger.Debugf("listInstances: only found virtual instance: %s", *instance.ServerName)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) destroyInstance(item cloudResource) error {
	var err error

	_, err = o.instanceClient.Get(item.id)
	if err != nil {
		o.deletePendingItems(instanceActionTypeName, []cloudResource{item})
		o.Logger.Infof("Deleted instance %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting instance %q", item.name)

	err = o.instanceClient.Delete(item.id)
	if err != nil {
		o.Logger.Infof("Error: o.instanceClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted instance %q", item.name)

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

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyInstances: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted instance %q", item.name)
				continue
			}
			err := o.destroyInstance(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(instanceTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(instanceTypeName); len(items) > 0 {
		return errors.Errorf("destroyInstances: %d undeleted items pending", len(items))
	}
	return nil
}
