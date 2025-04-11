package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"k8s.io/apimachinery/pkg/util/wait"
)

const powerSSHKeyTypeName = "powerSshKey"

// listPowerSSHKeys lists ssh keys in the Power server.
func (o *ClusterUninstaller) listPowerSSHKeys() (cloudResources, error) {
	o.Logger.Debugf("Listing Power SSHKeys")

	if o.keyClient == nil {
		o.Logger.Infof("Skipping deleting Power sshkeys because no service instance was found")
		result := []cloudResource{}
		return cloudResources{}.insert(result...), nil
	}

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listPowerSSHKeys: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	var sshKeys *models.SSHKeys
	var err error

	sshKeys, err = o.keyClient.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list Power sshkeys: %w", err)
	}

	var sshKey *models.SSHKey
	var foundOne = false

	result := []cloudResource{}
	for _, sshKey = range sshKeys.SSHKeys {
		if strings.Contains(*sshKey.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listPowerSSHKeys: FOUND: %v", *sshKey.Name)
			result = append(result, cloudResource{
				key:      *sshKey.Name,
				name:     *sshKey.Name,
				status:   "",
				typeName: powerSSHKeyTypeName,
				id:       *sshKey.Name,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listPowerSSHKeys: NO matching sshKey against: %s", o.InfraID)
		for _, sshKey := range sshKeys.SSHKeys {
			o.Logger.Debugf("listPowerSSHKeys: sshKey: %s", *sshKey.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deletePowerSSHKey(item cloudResource) error {
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deletePowerSSHKey: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	_, err = o.keyClient.Get(item.id)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Power SSHKey %q", item.name)
		return nil
	}

	err = o.keyClient.Delete(item.id)
	if err != nil {
		return fmt.Errorf("failed to delete Power sshKey %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Power SSHKey %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyPowerSSHKeys removes all ssh keys that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyPowerSSHKeys() error {
	firstPassList, err := o.listPowerSSHKeys()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(powerSSHKeyTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyPowerSSHKeys: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deletePowerSSHKey(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyPowerSSHKeys: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(powerSSHKeyTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyPowerSSHKeys: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyPowerSSHKeys: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listPowerSSHKeys()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyPowerSSHKeys: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyPowerSSHKeys: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
