package powervs

import (
	"context"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	"strings"
)

const (
	cloudInstanceTypeName = "cloudInstance"
)

// listCloudInstances lists instances in the cloud server.
func (o *ClusterUninstaller) listCloudInstances() (cloudResources, error) {
	o.Logger.Debugf("Listing virtual Cloud service instances")

	ctx, _ := o.contextWithTimeout()

	options := o.vpcSvc.NewListInstancesOptions()

	// https://raw.githubusercontent.com/IBM/vpc-go-sdk/master/vpcv1/vpc_v1.go
	resources, _, err := o.vpcSvc.ListInstancesWithContext(ctx, options)
	if err != nil {
		o.Logger.Warnf("Error o.vpcSvc.ListInstancesWithContext: %v", err)
		return nil, err
	}

	var foundOne = false

	result := []cloudResource{}
	for _, instance := range resources.Instances {
		if strings.Contains(*instance.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listCloudInstances: FOUND: %s, %s, %s", *instance.ID, *instance.Name, *instance.Status)
			result = append(result, cloudResource{
				key:      *instance.ID,
				name:     *instance.Name,
				status:   *instance.Status,
				typeName: cloudInstanceTypeName,
				id:       *instance.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listCloudInstances: NO matching virtual instance against: %s", o.InfraID)
		for _, instance := range resources.Instances {
			o.Logger.Debugf("listInstances: only found virtual instance: %s", *instance.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

// destroyCloudInstance deletes a given instance.
func (o *ClusterUninstaller) destroyCloudInstance(item cloudResource) error {
	var (
		ctx                   context.Context
		err                   error
		getInstanceOptions    *vpcv1.GetInstanceOptions
		deleteInstanceOptions *vpcv1.DeleteInstanceOptions
		response              *core.DetailedResponse
	)

	ctx, _ = o.contextWithTimeout()

	getInstanceOptions = o.vpcSvc.NewGetInstanceOptions(item.id)

	_, _, err = o.vpcSvc.GetInstanceWithContext(ctx, getInstanceOptions)
	if err != nil {
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Cloud instance %q", item.name)
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
	o.Logger.Infof("Deleted Cloud instance %q", item.name)

	return nil
}

// destroyCloudInstances searches for Cloud instances that have a name that starts with
// the cluster's infra ID.
func (o *ClusterUninstaller) destroyCloudInstances() error {
	found, err := o.listCloudInstances()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(cloudInstanceTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyCloudInstances: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted Cloud instance %q", item.name)
				continue
			}
			err := o.destroyCloudInstance(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(cloudInstanceTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(cloudInstanceTypeName); len(items) > 0 {
		return errors.Errorf("destroyCloudInstances: %d undeleted items pending", len(items))
	}

	return nil
}
