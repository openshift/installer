package powervs

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	cloudSSHKeyTypeName = "cloudSshKey"
)

// listCloudSSHKeys lists ssh-keys matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listCloudSSHKeys() (cloudResources, error) {
	var (
		keyIDs   []string
		keyID    string
		ctx      context.Context
		cancel   context.CancelFunc
		result   = make([]cloudResource, 0, 1)
		options  *vpcv1.GetKeyOptions
		response *core.DetailedResponse
		sshKey   *vpcv1.Key
		err      error
	)

	if false { // @TODO o.searchByTag {
		// Should we list by tag matching?
		// @TODO keyIDs, err = o.listByTag(TagTypeCloudSshKey)
		err = fmt.Errorf("listByTag(TagTypeCloudSshKey) is not supported yet")
	} else {
		// Otherwise list will list by name matching.
		keyIDs, err = o.listCloudSSHKeysByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, keyID = range keyIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listCloudSSHKeys: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.vpcSvc.NewGetKeyOptions(keyID)

		sshKey, response, err = o.vpcSvc.GetKeyWithContext(ctx, options)
		if err != nil {
			return nil, fmt.Errorf("failed to get cloud ssh-key (%s): err = %w, response = %v", keyID, err, response)
		}

		result = append(result, cloudResource{
			key:      *sshKey.Name,
			name:     *sshKey.Name,
			status:   "",
			typeName: cloudSSHKeyTypeName,
			id:       *sshKey.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listCloudSSHKeysByName lists ssh-keys matching by name in the IBM Cloud.
func (o *ClusterUninstaller) listCloudSSHKeysByName() ([]string, error) {
	var (
		ctx              context.Context
		cancel           context.CancelFunc
		foundOne         bool  = false
		perPage          int64 = 20
		moreData         bool  = true
		listKeysOptions  *vpcv1.ListKeysOptions
		sshKeyCollection *vpcv1.KeyCollection
		detailedResponse *core.DetailedResponse
		err              error
		sshKey           vpcv1.Key
		result           = make([]string, 0, 20)
	)

	o.Logger.Debugf("Listing Cloud SSHKeys by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	listKeysOptions = o.vpcSvc.NewListKeysOptions()
	listKeysOptions.SetLimit(perPage)
	listKeysOptions.SetResourceGroupID(o.resourceGroupID)

	for moreData {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listCloudSSHKeysByName: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		sshKeyCollection, detailedResponse, err = o.vpcSvc.ListKeysWithContext(ctx, listKeysOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list Cloud ssh keys: %w and the response is: %s", err, detailedResponse)
		}

		for _, sshKey = range sshKeyCollection.Keys {
			if strings.Contains(*sshKey.Name, o.InfraID) {
				foundOne = true
				o.Logger.Debugf("listCloudSSHKeysByName: FOUND: %v", *sshKey.Name)
				result = append(result, *sshKey.ID)
			}
		}

		if sshKeyCollection.First != nil {
			o.Logger.Debugf("listCloudSSHKeysByName: First = %v", *sshKeyCollection.First.Href)
		}
		if sshKeyCollection.Limit != nil {
			o.Logger.Debugf("listCloudSSHKeysByName: Limit = %v", *sshKeyCollection.Limit)
		}
		if sshKeyCollection.Next != nil {
			start, err := sshKeyCollection.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listCloudSSHKeysByName: err = %v", err)
				return nil, fmt.Errorf("listCloudSSHKeysByName: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listCloudSSHKeysByName: start = %v", *start)
				listKeysOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listCloudSSHKeysByName: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listCloudSSHKeysByName: NO matching sshKey against: %s", o.InfraID)

		listKeysOptions = o.vpcSvc.NewListKeysOptions()
		listKeysOptions.SetLimit(perPage)
		listKeysOptions.SetResourceGroupID(o.resourceGroupID)

		moreData = true

		for moreData {
			select {
			case <-ctx.Done():
				o.Logger.Debugf("listCloudSSHKeysByName: case <-ctx.Done()")
				return nil, ctx.Err() // we're cancelled, abort
			default:
			}

			sshKeyCollection, detailedResponse, err = o.vpcSvc.ListKeysWithContext(ctx, listKeysOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list Cloud ssh keys: %w and the response is: %s", err, detailedResponse)
			}
			for _, sshKey = range sshKeyCollection.Keys {
				o.Logger.Debugf("listCloudSSHKeysByName: FOUND: %v", *sshKey.Name)
			}
			if sshKeyCollection.First != nil {
				o.Logger.Debugf("listCloudSSHKeysByName: First = %v", *sshKeyCollection.First.Href)
			}
			if sshKeyCollection.Limit != nil {
				o.Logger.Debugf("listCloudSSHKeysByName: Limit = %v", *sshKeyCollection.Limit)
			}
			if sshKeyCollection.Next != nil {
				start, err := sshKeyCollection.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listCloudSSHKeysByName: err = %v", err)
					return nil, fmt.Errorf("listCloudSSHKeysByName: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listCloudSSHKeysByName: start = %v", *start)
					listKeysOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listCloudSSHKeysByName: Next = nil")
				moreData = false
			}
		}
	}

	return result, nil
}

// deleteCloudSSHKey deletes a given ssh key.
func (o *ClusterUninstaller) deleteCloudSSHKey(item cloudResource) error {
	var (
		ctx              context.Context
		cancel           context.CancelFunc
		getKeyOptions    *vpcv1.GetKeyOptions
		deleteKeyOptions *vpcv1.DeleteKeyOptions
		err              error
	)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteCloudSSHKey: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	getKeyOptions = o.vpcSvc.NewGetKeyOptions(item.id)

	_, _, err = o.vpcSvc.GetKey(getKeyOptions)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Cloud SSHKey %q", item.name)
		return nil
	}

	deleteKeyOptions = o.vpcSvc.NewDeleteKeyOptions(item.id)

	_, err = o.vpcSvc.DeleteKeyWithContext(ctx, deleteKeyOptions)
	if err != nil {
		return fmt.Errorf("failed to delete sshKey %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Cloud SSHKey %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyCloudSSHKeys removes all key resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyCloudSSHKeys() error {
	firstPassList, err := o.listCloudSSHKeys()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(cloudSSHKeyTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyCloudSSHKeys: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deleteCloudSSHKey(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyCloudSSHKeys: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(cloudSSHKeyTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyCloudSSHKeys: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyCloudSSHKeys: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyCloudSSHKeys: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listCloudSSHKeys()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyCloudSSHKeys: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyCloudSSHKeys: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
