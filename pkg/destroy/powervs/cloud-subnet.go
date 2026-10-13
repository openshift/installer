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

const cloudSubnetTypeName = "cloudSubnet"

// listCloudSubnets lists subnets matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listCloudSubnets() (cloudResources, error) {
	var (
		subnetIDs []string
		subnetID  string
		ctx       context.Context
		cancel    context.CancelFunc
		result    = make([]cloudResource, 0, 1)
		options   *vpcv1.GetSubnetOptions
		subnet    *vpcv1.Subnet
		response  *core.DetailedResponse
		err       error
	)

	if false { // @TODO o.searchByTag {
		// Should we list by tag matching?
		// @TODO subnetIDs, err = o.listByTag(TagTypeCloudSubnet)
		err = fmt.Errorf("listByTag(TagTypeCloudSubnet) is not supported yet")
	} else {
		// Otherwise list will list by name matching.
		subnetIDs, err = o.listCloudSubnetsByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, subnetID = range subnetIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listCloudSubnets: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.vpcSvc.NewGetSubnetOptions(subnetID)

		subnet, response, err = o.vpcSvc.GetSubnetWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The subnet could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get cloud subnet (%s): err = %w, response = %v", subnetID, err, response)
		}

		result = append(result, cloudResource{
			key:      *subnet.ID,
			name:     *subnet.Name,
			status:   "",
			typeName: cloudSubnetTypeName,
			id:       *subnet.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listCloudSubnetsByName lists subnets matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listCloudSubnetsByName() ([]string, error) {
	var (
		ctx              context.Context
		cancel           context.CancelFunc
		options          *vpcv1.ListSubnetsOptions
		subnetCollection *vpcv1.SubnetCollection
		response         *core.DetailedResponse
		foundOne         = false
		result           = make([]string, 0, 1)
		subnet           vpcv1.Subnet
		err              error
	)

	o.Logger.Debugf("Listing virtual Cloud Subnets by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listCloudSubnetsByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	options = o.vpcSvc.NewListSubnetsOptions()
	options.SetResourceGroupID(o.resourceGroupID)

	subnetCollection, response, err = o.vpcSvc.ListSubnetsWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list cloud subnets and the response is: %s: %w", response, err)
	}

	for _, subnet = range subnetCollection.Subnets {
		if strings.Contains(*subnet.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listCloudSubnetsByName: FOUND: %s, %s", *subnet.ID, *subnet.Name)
			result = append(result, *subnet.ID)
		}
	}
	if !foundOne {
		o.Logger.Debugf("listCloudSubnetsByName: NO matching subnet against: %s", o.InfraID)
		for _, subnet = range subnetCollection.Subnets {
			o.Logger.Debugf("listCloudSubnetsByName: subnet: %s", *subnet.Name)
		}
	}

	return result, nil
}

// deleteCloudSubnet deletes the cloud subnet specified.
func (o *ClusterUninstaller) deleteCloudSubnet(item cloudResource) error {
	var getOptions *vpcv1.GetSubnetOptions
	var response *core.DetailedResponse
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteCloudSubnet: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	getOptions = o.vpcSvc.NewGetSubnetOptions(item.id)
	_, response, err = o.vpcSvc.GetSubnet(getOptions)

	if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Subnet %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == gohttp.StatusInternalServerError {
		o.Logger.Infof("deleteCloudSubnet: internal server error")
		return nil
	}

	deleteOptions := o.vpcSvc.NewDeleteSubnetOptions(item.id)

	_, err = o.vpcSvc.DeleteSubnetWithContext(ctx, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete subnet %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Subnet %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyCloudSubnets removes all subnet resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyCloudSubnets() error {
	firstPassList, err := o.listCloudSubnets()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(cloudSubnetTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyCloudSubnets: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deleteCloudSubnet(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyCloudSubnets: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(cloudSubnetTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyCloudSubnets: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyCloudSubnets: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyCloudSubnets: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listCloudSubnets()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyCloudSubnets: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyCloudSubnets: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
