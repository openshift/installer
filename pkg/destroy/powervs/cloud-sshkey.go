package powervs

import (
	"context"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const (
	cloudSSHKeyTypeName = "cloudSshKey"
)

// listCloudSSHKeys lists images in the vpc.
func (o *ClusterUninstaller) listCloudSSHKeys() (cloudResources, error) {
	o.Logger.Debugf("Listing Cloud SSHKeys")

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listCloudSSHKeys: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

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

	ctx, _ = o.contextWithTimeout()

	listKeysOptions = o.vpcSvc.NewListKeysOptions()
	listKeysOptions.SetLimit(perPage)

	result := []cloudResource{}

	for moreData {
		sshKeyCollection, detailedResponse, err = o.vpcSvc.ListKeysWithContext(ctx, listKeysOptions)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list Cloud ssh keys: %v and the response is: %s", err, detailedResponse)
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
			o.Logger.Debugf("listCloudSSHKeys: Next = %v", *sshKeyCollection.Next.Href)
			listKeysOptions.SetStart(*sshKeyCollection.Next.Href)
		}

		moreData = sshKeyCollection.Next != nil
		o.Logger.Debugf("listCloudSSHKeys: moreData = %v", moreData)
	}
	if !foundOne {
		o.Logger.Debugf("listCloudSSHKeys: NO matching sshKey against: %s", o.InfraID)

		listKeysOptions = o.vpcSvc.NewListKeysOptions()
		listKeysOptions.SetLimit(perPage)
		moreData = true

		for moreData {
			sshKeyCollection, detailedResponse, err = o.vpcSvc.ListKeysWithContext(ctx, listKeysOptions)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to list Cloud ssh keys: %v and the response is: %s", err, detailedResponse)
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
				o.Logger.Debugf("listCloudSSHKeys: Next = %v", *sshKeyCollection.Next.Href)
				listKeysOptions.SetStart(*sshKeyCollection.Next.Href)
			}
			moreData = sshKeyCollection.Next != nil
			o.Logger.Debugf("listCloudSSHKeys: moreData = %v", moreData)
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

	ctx, _ = o.contextWithTimeout()

	getKeyOptions = o.vpcSvc.NewGetKeyOptions(item.id)

	_, _, err = o.vpcSvc.GetKey(getKeyOptions)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Cloud sshKey %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting Cloud sshKey %q", item.name)

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteCloudSSHKey: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	deleteKeyOptions = o.vpcSvc.NewDeleteKeyOptions(item.id)

	_, err = o.vpcSvc.DeleteKeyWithContext(ctx, deleteKeyOptions)
	if err != nil {
		return errors.Wrapf(err, "failed to delete sshKey %s", item.name)
	}

	return nil
}

// destroyCloudSSHKeys removes all image resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyCloudSSHKeys() error {
	found, err := o.listCloudSSHKeys()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(cloudSSHKeyTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyCloudSSHKeys: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted sshKey %q", item.name)
				continue
			}
			err := o.deleteCloudSSHKey(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(cloudSSHKeyTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(cloudSSHKeyTypeName); len(items) > 0 {
		return errors.Errorf("destroyCloudSSHKeys: %d undeleted items pending", len(items))
	}
	return nil
}
