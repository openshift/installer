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

const vpcTypeName = "vpc"

// listVPCs lists VPCs matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listVPCs() (cloudResources, error) {
	var (
		vpcIDs   []string
		vpcID    string
		ctx      context.Context
		cancel   context.CancelFunc
		result   = make([]cloudResource, 0, 1)
		options  *vpcv1.GetVPCOptions
		vpc      *vpcv1.VPC
		response *core.DetailedResponse
		err      error
	)

	if o.searchByTag {
		// Should we list by tag matching?
		vpcIDs, err = o.listByTag(TagTypeVPC)
	} else {
		// Otherwise list will list by name matching.
		vpcIDs, err = o.listVPCsByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, vpcID = range vpcIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listVPCs: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.vpcSvc.NewGetVPCOptions(vpcID)

		vpc, response, err = o.vpcSvc.GetVPCWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The VPC could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get vpc (%s): err = %w, response = %v", vpcID, err, response)
		}

		result = append(result, cloudResource{
			key:      *vpc.ID,
			name:     *vpc.Name,
			status:   "",
			typeName: vpcTypeName,
			id:       *vpc.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listVPCsByName lists VPCs matching by name in the IBM Cloud.
func (o *ClusterUninstaller) listVPCsByName() ([]string, error) {
	var (
		ctx           context.Context
		cancel        context.CancelFunc
		options       *vpcv1.ListVpcsOptions
		vpcCollection *vpcv1.VPCCollection
		response      *core.DetailedResponse
		foundOne      = false
		result        = make([]string, 0)
		vpc           vpcv1.VPC
		err           error
	)

	o.Logger.Debugf("Listing VPCs by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listVPCsByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	options = o.vpcSvc.NewListVpcsOptions()
	options.SetResourceGroupID(o.resourceGroupID)

	vpcCollection, response, err = o.vpcSvc.ListVpcsWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list vps: err = %w, response = %v", err, response)
	}

	for _, vpc = range vpcCollection.Vpcs {
		if strings.Contains(*vpc.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listVPCsByName: FOUND: %s, %s", *vpc.ID, *vpc.Name)
			result = append(result, *vpc.ID)
		}
	}
	if !foundOne {
		o.Logger.Debugf("listVPCsByName: NO matching vpc against: %s", o.InfraID)
		for _, vpc = range vpcCollection.Vpcs {
			o.Logger.Debugf("listVPCsByName: vpc: %s", *vpc.Name)
		}
	}

	return result, nil
}

// deleteVPC deletes the VPC specified.
func (o *ClusterUninstaller) deleteVPC(item cloudResource) error {
	var getOptions *vpcv1.GetVPCOptions
	var getResponse *core.DetailedResponse
	var deleteResponse *core.DetailedResponse
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteVPC: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	getOptions = o.vpcSvc.NewGetVPCOptions(item.id)
	_, getResponse, err = o.vpcSvc.GetVPC(getOptions)

	// Sadly, there is no way to get the status of this VPC to check on the results of the
	// delete call.

	if err == nil && getResponse.StatusCode == gohttp.StatusNoContent {
		return nil
	}
	if err != nil && getResponse != nil && getResponse.StatusCode == gohttp.StatusNotFound {
		// The resource is gone
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted VPC %q", item.name)
		return nil
	}
	if err != nil && getResponse != nil && getResponse.StatusCode == gohttp.StatusInternalServerError {
		o.Logger.Infof("deleteVPC: internal server error")
		return nil
	}
	if err != nil {
		o.Logger.Debugf("deleteVPC: getResponse = %v", getResponse)
		o.Logger.Debugf("deleteVPC: err = %v", err)
		return err
	}

	deleteOptions := o.vpcSvc.NewDeleteVPCOptions(item.id)

	deleteResponse, err = o.vpcSvc.DeleteVPCWithContext(ctx, deleteOptions)
	if err != nil {
		o.Logger.Debugf("deleteVPC: DeleteVPCWithContext returns %+v", deleteResponse)
		return fmt.Errorf("failed to delete vpc %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted VPC %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyVPCs removes all vpc resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyVPCs() error {
	firstPassList, err := o.listVPCs()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(vpcTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyVPCs: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deleteVPC(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyVPCs: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(vpcTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyVPCs: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyVPCs: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyVPCs: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listVPCs()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyVPCs: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyVPCs: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
