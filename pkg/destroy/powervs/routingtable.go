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
	routingTableTypeName = "routingTable"
)

// listRoutingTables lists VPC routing tables whose name contains the cluster's infra ID.
// Routing tables are scoped to a VPC, so this function first resolves the cluster VPC ID(s)
// via listVPCsByName and then queries routing tables for each VPC found.
func (o *ClusterUninstaller) listRoutingTables() (cloudResources, error) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	// Resolve the VPC(s) that belong to this cluster.
	vpcIDs, err := o.listVPCsByName()
	if err != nil {
		return nil, fmt.Errorf("listRoutingTables: failed to list cluster VPCs: %w", err)
	}
	if len(vpcIDs) == 0 {
		o.Logger.Debugf("listRoutingTables: no cluster VPC found, skipping routing table listing")
		return cloudResources{}.insert(), nil
	}

	result := make([]cloudResource, 0, 2)
	perPage := int64(20)

	for _, vpcID := range vpcIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listRoutingTables: case <-ctx.Done()")
			return nil, ctx.Err()
		default:
		}

		var (
			listOptions      *vpcv1.ListVPCRoutingTablesOptions
			detailedResponse *core.DetailedResponse
			moreData         = true
		)

		listOptions = o.vpcSvc.NewListVPCRoutingTablesOptions(vpcID)
		listOptions.SetLimit(perPage)

		for moreData {
			collection, detResp, listErr := o.vpcSvc.ListVPCRoutingTablesWithContext(ctx, listOptions)
			detailedResponse = detResp
			if listErr != nil {
				return nil, fmt.Errorf("failed to list routing tables for VPC %s: %s: %w", vpcID, detailedResponse, listErr)
			}

			for _, rt := range collection.RoutingTables {
				if !strings.Contains(*rt.Name, o.InfraID) {
					continue
				}
				o.Logger.Debugf("listRoutingTables: FOUND: %s (%s) in VPC %s", *rt.Name, *rt.ID, vpcID)

				// Verify it still exists with a GET before adding to results.
				getOptions := o.vpcSvc.NewGetVPCRoutingTableOptions(vpcID, *rt.ID)
				fullRT, response, getErr := o.vpcSvc.GetVPCRoutingTableWithContext(ctx, getOptions)
				if getErr != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
					continue
				}
				if getErr != nil {
					return nil, fmt.Errorf("failed to get routing table %s: err = %w, response = %v", *rt.ID, getErr, response)
				}

				// Encode the VPC ID in the status field so deleteRoutingTable can use it.
				result = append(result, cloudResource{
					key:      *fullRT.Name,
					name:     *fullRT.Name,
					status:   vpcID,
					typeName: routingTableTypeName,
					id:       *fullRT.ID,
				})
			}

			if collection.Next != nil {
				start, nextErr := collection.GetNextStart()
				if nextErr != nil {
					return nil, fmt.Errorf("listRoutingTables: GetNextStart: %w", nextErr)
				}
				if start != nil {
					listOptions.SetStart(*start)
					continue
				}
			}
			moreData = false
		}
	}

	return cloudResources{}.insert(result...), nil
}

// deleteRoutingTable deletes the VPC routing table described by item.
// item.status holds the VPC ID (set by listRoutingTables).
func (o *ClusterUninstaller) deleteRoutingTable(item cloudResource) error {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteRoutingTable: case <-ctx.Done()")
		return ctx.Err()
	default:
	}

	vpcID := item.status

	// Confirm it still exists.
	getOptions := o.vpcSvc.NewGetVPCRoutingTableOptions(vpcID, item.id)
	_, _, err := o.vpcSvc.GetVPCRoutingTableWithContext(ctx, getOptions)
	if err != nil {
		o.Logger.Debugf("deleteRoutingTable: routing table %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted VPC Routing Table %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting VPC routing table %q (%s) in VPC %s", item.name, item.id, vpcID)

	deleteOptions := o.vpcSvc.NewDeleteVPCRoutingTableOptions(vpcID, item.id)
	response, err := o.vpcSvc.DeleteVPCRoutingTableWithContext(ctx, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete routing table %s: err = %w, response = %v", item.name, err, response)
	}

	o.deletePendingItems(item.typeName, []cloudResource{item})
	o.Logger.Infof("Deleted VPC Routing Table %q", item.name)

	return nil
}

// destroyRoutingTables removes all VPC Routing Table resources whose name contains the
// cluster's infra ID. This must run before VPC deletion.
func (o *ClusterUninstaller) destroyRoutingTables() error {
	firstPassList, err := o.listRoutingTables()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(routingTableTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyRoutingTables: case <-ctx.Done()")
			return ctx.Err()
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32,
		}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deleteRoutingTable(item)
			if err2 == nil {
				return true, nil
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyRoutingTables: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(routingTableTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyRoutingTables: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyRoutingTables: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyRoutingTables: case <-ctx.Done()")
		return ctx.Err()
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32,
	}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listRoutingTables()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyRoutingTables: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyRoutingTables: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
