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
	publicGatewayTypeName = "publicGateway"
)

// listAttachedSubnetsByName lists subnets attached to the specified publicGateway.
func (o *ClusterUninstaller) listAttachedSubnetsByName(publicGatewayID string) ([]string, error) {
	var (
		ctx              context.Context
		cancel           context.CancelFunc
		listOptions      *vpcv1.ListSubnetsOptions
		subnetCollection *vpcv1.SubnetCollection
		subnet           vpcv1.Subnet
		result           = make([]string, 0, 1)
		response         *core.DetailedResponse
		err              error
	)

	o.Logger.Debugf("Finding subnets attached to public gateway %s", publicGatewayID)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listAttachedSubnetsByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	listOptions = o.vpcSvc.NewListSubnetsOptions()
	listOptions.SetResourceGroupID(o.resourceGroupID)

	subnetCollection, response, err = o.vpcSvc.ListSubnetsWithContext(ctx, listOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: response = %v, err = %w", response, err)
	}

	for _, subnet = range subnetCollection.Subnets {
		if subnet.PublicGateway != nil && *subnet.PublicGateway.ID == publicGatewayID {
			result = append(result, *subnet.ID)
		}
	}

	return result, nil
}

// listPublicGateways lists publicGateways matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listPublicGateways() (cloudResources, error) {
	var (
		gatewayIDs    []string
		gatewayID     string
		ctx           context.Context
		cancel        context.CancelFunc
		result        = make([]cloudResource, 0, 1)
		options       *vpcv1.GetPublicGatewayOptions
		publicGateway *vpcv1.PublicGateway
		response      *core.DetailedResponse
		err           error
	)

	if o.searchByTag {
		// Should we list by tag matching?
		gatewayIDs, err = o.listByTag(TagTypePublicGateway)
	} else {
		// Otherwise list will list by name matching.
		gatewayIDs, err = o.listPublicGatewaysByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, gatewayID = range gatewayIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listPublicGateways: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.vpcSvc.NewGetPublicGatewayOptions(gatewayID)

		publicGateway, response, err = o.vpcSvc.GetPublicGatewayWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The public gateway could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get public gateway (%s): err = %w, response = %v", gatewayID, err, response)
		}

		result = append(result, cloudResource{
			key:      *publicGateway.Name,
			name:     *publicGateway.Name,
			status:   "",
			typeName: publicGatewayTypeName,
			id:       *publicGateway.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listPublicGatewaysByName lists publicGateways matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listPublicGatewaysByName() ([]string, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		// https://raw.githubusercontent.com/IBM/vpc-go-sdk/master/vpcv1/vpc_v1.go
		listPublicGatewaysOptions *vpcv1.ListPublicGatewaysOptions
		publicGatewayCollection   *vpcv1.PublicGatewayCollection
		publicGateway             vpcv1.PublicGateway
		detailedResponse          *core.DetailedResponse
		err                       error
		moreData                  bool  = true
		foundOne                  bool  = false
		perPage                   int64 = 20
		result                          = make([]string, 0, 1)
	)

	o.Logger.Debugf("Listing publicGateways by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listPublicGatewaysByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	listPublicGatewaysOptions = o.vpcSvc.NewListPublicGatewaysOptions()
	listPublicGatewaysOptions.SetLimit(perPage)
	listPublicGatewaysOptions.SetResourceGroupID(o.resourceGroupID)

	for moreData {

		publicGatewayCollection, detailedResponse, err = o.vpcSvc.ListPublicGatewaysWithContext(ctx, listPublicGatewaysOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list publicGateways and the response is: %s: %w", detailedResponse, err)
		}

		for _, publicGateway = range publicGatewayCollection.PublicGateways {
			if strings.Contains(*publicGateway.Name, o.InfraID) {
				foundOne = true
				o.Logger.Debugf("listPublicGatewaysByName: FOUND: %s", *publicGateway.Name)
				result = append(result, *publicGateway.ID)
			}
		}

		if publicGatewayCollection.First != nil {
			o.Logger.Debugf("listPublicGatewaysByName: First = %v", *publicGatewayCollection.First.Href)
		}
		if publicGatewayCollection.Limit != nil {
			o.Logger.Debugf("listPublicGatewaysByName: Limit = %v", *publicGatewayCollection.Limit)
		}
		if publicGatewayCollection.Next != nil {
			start, err := publicGatewayCollection.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listPublicGatewaysByName: err = %v", err)
				return nil, fmt.Errorf("listPublicGatewaysByName: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listPublicGatewaysByName: start = %v", *start)
				listPublicGatewaysOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listPublicGatewaysByName: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listPublicGatewaysByName: NO matching publicGateway against: %s", o.InfraID)

		listPublicGatewaysOptions = o.vpcSvc.NewListPublicGatewaysOptions()
		listPublicGatewaysOptions.SetLimit(perPage)
		listPublicGatewaysOptions.SetResourceGroupID(o.resourceGroupID)

		for moreData {
			publicGatewayCollection, detailedResponse, err = o.vpcSvc.ListPublicGatewaysWithContext(ctx, listPublicGatewaysOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list publicGateways and the response is: %s: %w", detailedResponse, err)
			}

			for _, publicGateway = range publicGatewayCollection.PublicGateways {
				o.Logger.Debugf("listPublicGatewaysByName: FOUND: %s", *publicGateway.Name)
			}
			if publicGatewayCollection.First != nil {
				o.Logger.Debugf("listPublicGatewaysByName: First = %v", *publicGatewayCollection.First.Href)
			}
			if publicGatewayCollection.Limit != nil {
				o.Logger.Debugf("listPublicGatewaysByName: Limit = %v", *publicGatewayCollection.Limit)
			}
			if publicGatewayCollection.Next != nil {
				start, err := publicGatewayCollection.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listPublicGatewaysByName: err = %v", err)
					return nil, fmt.Errorf("listPublicGatewaysByName: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listPublicGatewaysByName: start = %v", *start)
					listPublicGatewaysOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listPublicGatewaysByName: Next = nil")
				moreData = false
			}
		}
	}

	return result, nil
}

// deletePublicGateway deletes the public gateway specified.
func (o *ClusterUninstaller) deletePublicGateway(item cloudResource) error {
	var (
		ctx                             context.Context
		cancel                          context.CancelFunc
		getPublicGatewayOptions         *vpcv1.GetPublicGatewayOptions
		subnets                         []string
		subnetID                        string
		unsetSubnetPublicGatewayOptions *vpcv1.UnsetSubnetPublicGatewayOptions
		deletePublicGatewayOptions      *vpcv1.DeletePublicGatewayOptions
		err                             error
	)

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deletePublicGateway: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	getPublicGatewayOptions = o.vpcSvc.NewGetPublicGatewayOptions(item.id)

	_, _, err = o.vpcSvc.GetPublicGatewayWithContext(ctx, getPublicGatewayOptions)
	if err != nil {
		o.Logger.Debugf("deletePublicGateway: publicGateway %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Public Gateway %q", item.name)
		return nil
	}

	// Detach gateway from any subnets using it
	subnets, err = o.listAttachedSubnetsByName(item.id)
	if err != nil {
		return fmt.Errorf("failed to list subnets with gateway %s attached: %w", item.name, err)
	}

	for _, subnetID = range subnets {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("deletePublicGateway: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		unsetSubnetPublicGatewayOptions = o.vpcSvc.NewUnsetSubnetPublicGatewayOptions(subnetID)

		_, err = o.vpcSvc.UnsetSubnetPublicGatewayWithContext(ctx, unsetSubnetPublicGatewayOptions)
		if err != nil {
			return fmt.Errorf("failed to detach publicGateway %s from subnet %s: %w", item.name, subnetID, err)
		}
	}

	deletePublicGatewayOptions = o.vpcSvc.NewDeletePublicGatewayOptions(item.id)

	_, err = o.vpcSvc.DeletePublicGatewayWithContext(ctx, deletePublicGatewayOptions)
	if err != nil {
		return fmt.Errorf("failed to delete publicGateway %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Public Gateway %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyPublicGateways removes all publicGateway resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyPublicGateways() error {
	firstPassList, err := o.listPublicGateways()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(publicGatewayTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyPublicGateways: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deletePublicGateway(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyPublicGateways: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(publicGatewayTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyPublicGateways: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyPublicGateways: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyPublicGateways: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listPublicGateways()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyPublicGateways: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyPublicGateways: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
