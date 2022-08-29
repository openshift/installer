package powervs

import (
	"strings"

	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/pkg/errors"
)

const powerSSHKeyTypeName = "powerSshKey"

// listPowerSSHKeys lists ssh keys in the Power server.
func (o *ClusterUninstaller) listPowerSSHKeys() (cloudResources, error) {
	o.Logger.Debugf("Listing Power SSHKeys")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listPowerSSHKeys: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	var sshKeys *models.SSHKeys
	var err error

	sshKeys, err = o.keyClient.GetAll()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list Power sshkeys: %v", err)
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

// deletePowerSSHKey deleted a given ssh key.
func (o *ClusterUninstaller) deletePowerSSHKey(item cloudResource) error {
	var err error

	_, err = o.keyClient.Get(item.id)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Power sshKey %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting Power sshKey %q", item.name)

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deletePowerSSHKey: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	err = o.keyClient.Delete(item.id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete Power sshKey %s", item.name)
	}

	return nil
}

// destroyPowerSSHKeys removes all ssh keys that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyPowerSSHKeys() error {
	found, err := o.listPowerSSHKeys()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(powerSSHKeyTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyPowerSSHKeys: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted sshKey %q", item.name)
				continue
			}
			err := o.deletePowerSSHKey(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(powerSSHKeyTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(powerSSHKeyTypeName); len(items) > 0 {
		return errors.Errorf("destroyPower1SSHKeys: %d undeleted items pending", len(items))
	}
	return nil
}
