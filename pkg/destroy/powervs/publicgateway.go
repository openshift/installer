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
	publicGatewayTypeName = "publicGateway"
)

// listAttachedSubnets lists subnets attached to the specified publicGateway.
func (o *ClusterUninstaller) listAttachedSubnets(publicGatewayID string) (cloudResources, error) {
	o.Logger.Debugf("Finding subnets attached to public gateway %s", publicGatewayID)

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListSubnetsOptions()
	resources, _, err := o.vpcSvc.ListSubnetsWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: %w", err)
	}

	result := []cloudResource{}
	for _, subnet := range resources.Subnets {
		if subnet.PublicGateway != nil && *subnet.PublicGateway.ID == publicGatewayID {
			result = append(result, cloudResource{
				key:      *subnet.ID,
				name:     *subnet.Name,
				status:   *subnet.Status,
				typeName: cloudSubnetTypeName,
				id:       *subnet.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

// listPublicGateways lists publicGateways in the vpc.
func (o *ClusterUninstaller) listPublicGateways() (cloudResources, error) {
	var (
		ctx context.Context
		// https://raw.githubusercontent.com/IBM/vpc-go-sdk/master/vpcv1/vpc_v1.go
		listPublicGatewaysOptions *vpcv1.ListPublicGatewaysOptions
		publicGatewayCollection   *vpcv1.PublicGatewayCollection
		detailedResponse          *core.DetailedResponse
		err                       error
		moreData                  bool  = true
		foundOne                  bool  = false
		perPage                   int64 = 20
	)

	o.Logger.Debugf("Listing publicGateways")

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listPublicGateways: case <-ctx.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	listPublicGatewaysOptions = o.vpcSvc.NewListPublicGatewaysOptions()

	listPublicGatewaysOptions.SetLimit(perPage)

	result := []cloudResource{}

	for moreData {

		publicGatewayCollection, detailedResponse, err = o.vpcSvc.ListPublicGatewaysWithContext(ctx, listPublicGatewaysOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to list publicGateways and the response is: %s: %w", detailedResponse, err)
		}

		for _, publicGateway := range publicGatewayCollection.PublicGateways {
			if strings.Contains(*publicGateway.Name, o.InfraID) {
				foundOne = true
				o.Logger.Debugf("listPublicGateways: FOUND: %s", *publicGateway.Name)
				result = append(result, cloudResource{
					key:      *publicGateway.Name,
					name:     *publicGateway.Name,
					status:   "",
					typeName: publicGatewayTypeName,
					id:       *publicGateway.ID,
				})
			}
		}

		if publicGatewayCollection.First != nil {
			o.Logger.Debugf("listPublicGateways: First = %v", *publicGatewayCollection.First.Href)
		}
		if publicGatewayCollection.Limit != nil {
			o.Logger.Debugf("listPublicGateways: Limit = %v", *publicGatewayCollection.Limit)
		}
		if publicGatewayCollection.Next != nil {
			start, err := publicGatewayCollection.GetNextStart()
			if err != nil {
				o.Logger.Debugf("listPublicGateways: err = %v", err)
				return nil, fmt.Errorf("listPublicGateways: failed to GetNextStart: %w", err)
			}
			if start != nil {
				o.Logger.Debugf("listPublicGateways: start = %v", *start)
				listPublicGatewaysOptions.SetStart(*start)
			}
		} else {
			o.Logger.Debugf("listPublicGateways: Next = nil")
			moreData = false
		}
	}
	if !foundOne {
		o.Logger.Debugf("listPublicGateways: NO matching publicGateway against: %s", o.InfraID)

		listPublicGatewaysOptions = o.vpcSvc.NewListPublicGatewaysOptions()
		listPublicGatewaysOptions.SetLimit(perPage)

		for moreData {
			publicGatewayCollection, detailedResponse, err = o.vpcSvc.ListPublicGatewaysWithContext(ctx, listPublicGatewaysOptions)
			if err != nil {
				return nil, fmt.Errorf("failed to list publicGateways and the response is: %s: %w", detailedResponse, err)
			}

			for _, publicGateway := range publicGatewayCollection.PublicGateways {
				o.Logger.Debugf("listPublicGateways: FOUND: %s", *publicGateway.Name)
			}
			if publicGatewayCollection.First != nil {
				o.Logger.Debugf("listPublicGateways: First = %v", *publicGatewayCollection.First.Href)
			}
			if publicGatewayCollection.Limit != nil {
				o.Logger.Debugf("listPublicGateways: Limit = %v", *publicGatewayCollection.Limit)
			}
			if publicGatewayCollection.Next != nil {
				start, err := publicGatewayCollection.GetNextStart()
				if err != nil {
					o.Logger.Debugf("listPublicGateways: err = %v", err)
					return nil, fmt.Errorf("listPublicGateways: failed to GetNextStart: %w", err)
				}
				if start != nil {
					o.Logger.Debugf("listPublicGateways: start = %v", *start)
					listPublicGatewaysOptions.SetStart(*start)
				}
			} else {
				o.Logger.Debugf("listPublicGateways: Next = nil")
				moreData = false
			}
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deletePublicGateway(item cloudResource) error {
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deletePublicGateway: case <-ctx.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	getPublicGatewayOptions := o.vpcSvc.NewGetPublicGatewayOptions(item.id)

	_, _, err := o.vpcSvc.GetPublicGatewayWithContext(ctx, getPublicGatewayOptions)
	if err != nil {
		o.Logger.Debugf("deletePublicGateway: publicGateway %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Public Gateway %q", item.name)
		return nil
	}

	// Detach gateway from any subnets using it
	subnets, err := o.listAttachedSubnets(item.id)
	if err != nil {
		return fmt.Errorf("failed to list subnets with gateway %s attached: %w", item.name, err)
	}
	for _, subnet := range subnets {
		unsetSubnetPublicGatewayOptions := o.vpcSvc.NewUnsetSubnetPublicGatewayOptions(subnet.id)

		_, err = o.vpcSvc.UnsetSubnetPublicGatewayWithContext(ctx, unsetSubnetPublicGatewayOptions)
		if err != nil {
			return fmt.Errorf("failed to detach publicGateway %s from subnet %s: %w", item.name, subnet.id, err)
		}
	}

	deletePublicGatewayOptions := o.vpcSvc.NewDeletePublicGatewayOptions(item.id)

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

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyPublicGateways: case <-ctx.Done()")
			return o.Context.Err() // we're cancelled, abort
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
