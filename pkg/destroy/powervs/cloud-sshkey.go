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

// listCloudSSHKeys lists images in the vpc.
func (o *ClusterUninstaller) listCloudSSHKeys() (cloudResources, error) {
	o.Logger.Debugf("Listing Cloud SSHKeys")

	// https://raw.githubusercontent.com/IBM/vpc-go-sdk/master/vpcv1/vpc_v1.go
	var (
		ctx              context.Context
		foundOne         bool  = false
		perPage          int64 = 20
		moreData         bool  = true
		listKeysOptions  *vpcv1.ListKeysOptions
		sshKeyCollection *vpcv1.KeyCollection
		detailedResponse *core.DetailedResponse
		err              error
		sshKey           vpcv1.Key
	)

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listCloudSSHKeys: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	listKeysOptions = o.vpcSvc.NewListKeysOptions()
	listKeysOptions.SetLimit(perPage)
	listKeysOptions.SetResourceGroupID(o.resourceGroupID)

	result := []cloudResource{}

	for moreData {
		sshKeyCollection, detailedResponse, err = o.vpcSvc.ListKeysWithContext(ctx, listKeysOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list Cloud ssh keys: %w and the response is: %s", err, detailedResponse)
		}

		for _, sshKey = range sshKeyCollection.Keys {
			if strings.Contains(*sshKey.Name, o.InfraID) {
				foundOne = true
				o.Logger.Debugf("listCloudSSHKeys: FOUND: %v", *sshKey.Name)
				result = append(result, cloudResource{
					key:      *sshKey.Name,
					name:     *sshKey.Name,
					status:   "",
					typeName: cloudSSHKeyTypeName,
					id:       *sshKey.ID,
				})
			}
		}

		if sshKeyCollection.First != nil {
			o.Logger.Debugf("listCloudSSHKeys: First = %v", *sshKeyCollection.First.Href)
		}
		if sshKeyCollection.Limit != nil {
			o.Logger.Debugf("listCloudSSHKeys: Limit = %v", *sshKeyCollection.Limit)
		}
		if sshKeyCollection.Next != nil {
			start, err := sshKeyCollection.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listCloudSSHKeys: err = %v", err)
				return nil, fmt.Errorf("listCloudSSHKeys: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listCloudSSHKeys: start = %v", *start)
				listKeysOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listCloudSSHKeys: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listCloudSSHKeys: NO matching sshKey against: %s", o.InfraID)

		listKeysOptions = o.vpcSvc.NewListKeysOptions()
		listKeysOptions.SetLimit(perPage)
		listKeysOptions.SetResourceGroupID(o.resourceGroupID)

		moreData = true

		for moreData {
			sshKeyCollection, detailedResponse, err = o.vpcSvc.ListKeysWithContext(ctx, listKeysOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list Cloud ssh keys: %w and the response is: %s", err, detailedResponse)
			}
			for _, sshKey = range sshKeyCollection.Keys {
				o.Logger.Debugf("listCloudSSHKeys: FOUND: %v", *sshKey.Name)
			}
			if sshKeyCollection.First != nil {
				o.Logger.Debugf("listCloudSSHKeys: First = %v", *sshKeyCollection.First.Href)
			}
			if sshKeyCollection.Limit != nil {
				o.Logger.Debugf("listCloudSSHKeys: Limit = %v", *sshKeyCollection.Limit)
			}
			if sshKeyCollection.Next != nil {
				start, err := sshKeyCollection.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listCloudSSHKeys: err = %v", err)
					return nil, fmt.Errorf("listCloudSSHKeys: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listCloudSSHKeys: start = %v", *start)
					listKeysOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listCloudSSHKeys: Next = nil")
				moreData = false
			}
		}
	}

	return cloudResources{}.insert(result...), nil
}

// deleteCloudSSHKey deletes a given ssh key.
func (o *ClusterUninstaller) deleteCloudSSHKey(item cloudResource) error {
	var (
		ctx              context.Context
		getKeyOptions    *vpcv1.GetKeyOptions
		deleteKeyOptions *vpcv1.DeleteKeyOptions
		err              error
	)

	ctx, cancel := contextWithTimeout()
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
