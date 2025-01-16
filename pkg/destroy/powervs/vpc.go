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

// listVPCs lists VPCs in the cloud.
func (o *ClusterUninstaller) listVPCs() (cloudResources, error) {
	o.Logger.Debugf("Listing VPCs")

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listVPCs: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListVpcsOptions()
	options.SetResourceGroupID(o.resourceGroupID)

	vpcs, _, err := o.vpcSvc.ListVpcs(options)
	if err != nil {
		return nil, fmt.Errorf("failed to list vps: %w", err)
	}

	var foundOne = false

	result := []cloudResource{}
	for _, vpc := range vpcs.Vpcs {
		if strings.Contains(*vpc.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listVPCs: FOUND: %s, %s", *vpc.ID, *vpc.Name)
			result = append(result, cloudResource{
				key:      *vpc.ID,
				name:     *vpc.Name,
				status:   "",
				typeName: vpcTypeName,
				id:       *vpc.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listVPCs: NO matching vpc against: %s", o.InfraID)
		for _, vpc := range vpcs.Vpcs {
			o.Logger.Debugf("listVPCs: vpc: %s", *vpc.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

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

	o.Logger.Debugf("deleteVPC: getResponse = %v", getResponse)
	o.Logger.Debugf("deleteVPC: err = %v", err)

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

	deleteOptions := o.vpcSvc.NewDeleteVPCOptions(item.id)
	deleteResponse, err = o.vpcSvc.DeleteVPCWithContext(ctx, deleteOptions)
	o.Logger.Debugf("deleteVPC: DeleteVPCWithContext returns %+v", deleteResponse)

	if err != nil {
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
