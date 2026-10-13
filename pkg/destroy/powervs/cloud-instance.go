package powervs

import (
	"context"
	"fmt"
	"math"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	cloudInstanceTypeName = "cloudInstance"
)

// listCloudInstances lists VM instances matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listCloudInstances() (cloudResources, error) {
	var (
		instanceIDs []string
		instanceID  string
		ctx         context.Context
		cancel      context.CancelFunc
		result      = make([]cloudResource, 0, 1)
		options     *vpcv1.GetInstanceOptions
		instance    *vpcv1.Instance
		response    *core.DetailedResponse
		err         error
	)

	if o.searchByTag {
		// Should we list by tag matching?
		instanceIDs, err = o.listByTag(TagTypeCloudInstance)
	} else {
		// Otherwise list will list by name matching.
		instanceIDs, err = o.listCloudInstancesByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, instanceID = range instanceIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listCloudInstances: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.vpcSvc.NewGetInstanceOptions(instanceID)

		instance, response, err = o.vpcSvc.GetInstanceWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The vpc could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get instance (%s): err = %w, response = %v", instanceID, err, response)
		}

		result = append(result, cloudResource{
			key:      *instance.ID,
			name:     *instance.Name,
			status:   *instance.Status,
			typeName: cloudInstanceTypeName,
			id:       *instance.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listCloudInstances lists VM instances matching by name in the IBM Cloud.
func (o *ClusterUninstaller) listCloudInstancesByName() ([]string, error) {
	var (
		ctx                context.Context
		cancel             context.CancelFunc
		options            *vpcv1.ListInstancesOptions
		instanceCollection *vpcv1.InstanceCollection
		response           *core.DetailedResponse
		foundOne           = false
		result             = make([]string, 0, 3)
		instance           vpcv1.Instance
		err                error
	)

	o.Logger.Debugf("Listing virtual Cloud service instances by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listCloudInstancesByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	options = o.vpcSvc.NewListInstancesOptions()
	options.SetResourceGroupID(o.resourceGroupID)

	// https://raw.githubusercontent.com/IBM/vpc-go-sdk/master/vpcv1/vpc_v1.go
	instanceCollection, response, err = o.vpcSvc.ListInstancesWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list cloud instances: err = %w, response = %v", err, response)
	}

	for _, instance = range instanceCollection.Instances {
		if strings.Contains(*instance.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listCloudInstances: FOUND: %s, %s, %s", *instance.ID, *instance.Name, *instance.Status)
			result = append(result, *instance.ID)
		}
	}
	if !foundOne {
		o.Logger.Debugf("listCloudInstances: NO matching virtual instance against: %s", o.InfraID)
		for _, instance = range instanceCollection.Instances {
			o.Logger.Debugf("listCloudInstances: only found virtual instance: %s", *instance.Name)
		}
	}

	return result, nil
}

// destroyCloudInstance deletes a given instance.
func (o *ClusterUninstaller) destroyCloudInstance(item cloudResource) error {
	var (
		ctx                   context.Context
		cancel                context.CancelFunc
		err                   error
		getInstanceOptions    *vpcv1.GetInstanceOptions
		deleteInstanceOptions *vpcv1.DeleteInstanceOptions
		response              *core.DetailedResponse
	)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyCloudInstance: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	getInstanceOptions = o.vpcSvc.NewGetInstanceOptions(item.id)

	_, _, err = o.vpcSvc.GetInstanceWithContext(ctx, getInstanceOptions)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Cloud Instance %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting Cloud instance %q", item.name)

	deleteInstanceOptions = o.vpcSvc.NewDeleteInstanceOptions(item.id)

	response, err = o.vpcSvc.DeleteInstanceWithContext(ctx, deleteInstanceOptions)
	if err != nil {
		o.Logger.Infof("Error: o.vpcSvc.DeleteInstanceWithContext: %q %q", err, response)
		return err
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted Cloud Instance %q", item.name)

	return nil
}

// destroyCloudInstances searches for Cloud instances that have a name that starts with
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyCloudInstances() error {
	firstPassList, err := o.listCloudInstances()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(cloudInstanceTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyCloudInstances: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.destroyCloudInstance(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyCloudInstances: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(cloudInstanceTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyCloudInstances: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyCloudInstances: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyCloudInstances: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listCloudInstances()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyCloudInstances: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyCloudInstances: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
