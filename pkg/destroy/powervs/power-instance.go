package powervs

import (
	"github.com/pkg/errors"
	"strings"
)

const (
	powerInstanceTypeName = "powerInstance"
)

// listPowerInstances lists instances in the Power server.
func (o *ClusterUninstaller) listPowerInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing virtual Power service instances")

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
			o.Logger.Debugf("listPowerInstances: FOUND: %s, %s, %s", *instance.PvmInstanceID, *instance.ServerName, *instance.Status)
			result = append(result, cloudResource{
				key:      *instance.PvmInstanceID,
				name:     *instance.ServerName,
				status:   *instance.Status,
				typeName: powerInstanceTypeName,
				id:       *instance.PvmInstanceID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listPowerInstances: NO matching virtual instance against: %s", o.InfraID)
		for _, instance := range instances.PvmInstances {
			o.Logger.Debugf("listInstances: only found virtual instance: %s", *instance.ServerName)
		}
	}

	return cloudResources{}.insert(result...), nil
}

// destroyPowerInstance deletes a given instance.
func (o *ClusterUninstaller) destroyPowerInstance(item cloudResource) error {
	var err error

	_, err = o.instanceClient.Get(item.id)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Power instance %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting Power instance %q", item.name)

	err = o.instanceClient.Delete(item.id)
	if err != nil {
		o.Logger.Infof("Error: o.instanceClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted Power instance %q", item.name)

	return nil
}

// destroyPowerInstances searches for Power instances that have a name that starts with
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyPowerInstances() error {
	found, err := o.listPowerInstances()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(powerInstanceTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyPowerInstances: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted Power instance %q", item.name)
				continue
			}
			err := o.destroyPowerInstance(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(powerInstanceTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(powerInstanceTypeName); len(items) > 0 {
		return errors.Errorf("destroyPowerInstances: %d undeleted items pending", len(items))
	}
	return nil
}
