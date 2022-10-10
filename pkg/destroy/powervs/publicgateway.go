package powervs

import (
	"context"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "failed to list subnets")
	}

	result := []cloudResource{}
	for _, subnet := range resources.Subnets {
		if subnet.PublicGateway != nil && *subnet.PublicGateway.ID == publicGatewayID {
			result = append(result, cloudResource{
				key:      *subnet.ID,
				name:     *subnet.Name,
				status:   *subnet.Status,
				typeName: subnetTypeName,
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

	ctx, _ = o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listPublicGateways: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	listPublicGatewaysOptions = o.vpcSvc.NewListPublicGatewaysOptions()

	listPublicGatewaysOptions.SetLimit(perPage)

	result := []cloudResource{}

	for moreData {

		publicGatewayCollection, detailedResponse, err = o.vpcSvc.ListPublicGatewaysWithContext(ctx, listPublicGatewaysOptions)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to list publicGateways and the response is: %s", detailedResponse)
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
			o.Logger.Debugf("listPublicGateways: Next = %v", *publicGatewayCollection.Next.Href)
			listPublicGatewaysOptions.SetStart(*publicGatewayCollection.Next.Href)
		}

		moreData = publicGatewayCollection.Next != nil
		o.Logger.Debugf("listPublicGateways: moreData = %v", moreData)
	}
	if !foundOne {
		o.Logger.Debugf("listPublicGateways: NO matching publicGateway against: %s", o.InfraID)

		listPublicGatewaysOptions = o.vpcSvc.NewListPublicGatewaysOptions()
		listPublicGatewaysOptions.SetLimit(perPage)

		for moreData {
			publicGatewayCollection, detailedResponse, err = o.vpcSvc.ListPublicGatewaysWithContext(ctx, listPublicGatewaysOptions)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to list publicGateways and the response is: %s", detailedResponse)
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
				o.Logger.Debugf("listPublicGateways: Next = %v", *publicGatewayCollection.Next.Href)
				listPublicGatewaysOptions.SetStart(*publicGatewayCollection.Next.Href)
			}
			moreData = publicGatewayCollection.Next != nil
			o.Logger.Debugf("listPublicGateways: moreData = %v", moreData)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deletePublicGateway(item cloudResource) error {
	ctx, _ := o.contextWithTimeout()

	getPublicGatewayOptions := o.vpcSvc.NewGetPublicGatewayOptions(item.id)

	_, _, err := o.vpcSvc.GetPublicGatewayWithContext(ctx, getPublicGatewayOptions)
	if err != nil {
		o.Logger.Debugf("deletePublicGateway: publicGateway %q no longer exists", item.name)
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted publicGateway %q", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting publicGateway %q", item.name)

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deletePublicGateway: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	// Detach gateway from any subnets using it
	subnets, err := o.listAttachedSubnets(item.id)
	if err != nil {
		return errors.Wrapf(err, "failed to list subnets with gateway %s attached", item.name)
	}
	for _, subnet := range subnets {
		unsetSubnetPublicGatewayOptions := o.vpcSvc.NewUnsetSubnetPublicGatewayOptions(subnet.id)

		_, err = o.vpcSvc.UnsetSubnetPublicGatewayWithContext(ctx, unsetSubnetPublicGatewayOptions)
		if err != nil {
			return errors.Wrapf(err, "failed to detach publicGateway %s from subnet %s", item.name, subnet.id)
		}
	}

	deletePublicGatewayOptions := o.vpcSvc.NewDeletePublicGatewayOptions(item.id)

	_, err = o.vpcSvc.DeletePublicGatewayWithContext(ctx, deletePublicGatewayOptions)
	if err != nil {
		return errors.Wrapf(err, "failed to delete publicGateway %s", item.name)
	}

	return nil
}

// destroyPublicGateways removes all publicGateway resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyPublicGateways() error {
	found, err := o.listPublicGateways()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(publicGatewayTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyPublicGateways: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted publicGateway %q", item.name)
				continue
			}
			err := o.deletePublicGateway(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(publicGatewayTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(publicGatewayTypeName); len(items) > 0 {
		return errors.Errorf("destroyPublicGateways: %d undeleted items pending", len(items))
	}
	return nil
}
