package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	powerInstanceTypeName = "powerInstance"
)

// listPowerInstances lists instances in the power server.
func (o *ClusterUninstaller) listPowerInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing virtual Power service instances (%s)", o.InfraID)

	if o.instanceClient == nil {
		o.Logger.Infof("Skipping deleting Power service instances because no service instance was found")
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

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
			o.Logger.Debugf("listPowerInstances: only found virtual instance: %s", *instance.ServerName)
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
		o.Logger.Infof("Deleted Power Instance %q (%q)", item.name, item.status)
		return nil
	}

	o.Logger.Debugf("Deleting Power instance %q", item.name)

	err = o.instanceClient.Delete(item.id)
	if err != nil {
		o.Logger.Infof("Error: o.instanceClient.Delete: %q", err)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted Power Instance %q (%q)", item.name, item.status)

	return nil
}

// destroyPowerInstances searches for Power instances that have a name that starts with
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyPowerInstances() error {
	firstPassList, err := o.listPowerInstances()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(powerInstanceTypeName, firstPassList.list())

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyPowerInstances: case <-ctx.Done()")
			return o.Context.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.destroyPowerInstance(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyPowerInstances: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(powerInstanceTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyPowerInstances: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyPowerInstances: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listPowerInstances()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyPowerInstances: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyPowerInstances: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
