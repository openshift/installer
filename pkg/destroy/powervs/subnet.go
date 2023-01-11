package powervs

import (
	"math"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

const subnetTypeName = "subnet"

// listSubnets lists subnets in the cloud.
func (o *ClusterUninstaller) listSubnets() (cloudResources, error) {
	o.Logger.Debugf("Listing Subnets")

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listSubnets: case <-ctx.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListSubnetsOptions()
	subnets, detailedResponse, err := o.vpcSvc.ListSubnets(options)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to list subnets and the response is: %s", detailedResponse)
	}

	var foundOne = false

	result := []cloudResource{}
	for _, subnet := range subnets.Subnets {
		if strings.Contains(*subnet.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listSubnets: FOUND: %s, %s", *subnet.ID, *subnet.Name)
			result = append(result, cloudResource{
				key:      *subnet.ID,
				name:     *subnet.Name,
				status:   "",
				typeName: subnetTypeName,
				id:       *subnet.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listSubnets: NO matching subnet against: %s", o.InfraID)
		for _, subnet := range subnets.Subnets {
			o.Logger.Debugf("listSubnets: subnet: %s", *subnet.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteSubnet(item cloudResource) error {
	var getOptions *vpcv1.GetSubnetOptions
	var response *core.DetailedResponse
	var err error

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteSubnet: case <-ctx.Done()")
		return o.Context.Err() // we're cancelled, abort
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
		o.Logger.Infof("deleteSubnet: internal server error")
		return nil
	}

	deleteOptions := o.vpcSvc.NewDeleteSubnetOptions(item.id)
	_, err = o.vpcSvc.DeleteSubnetWithContext(ctx, deleteOptions)
	if err != nil {
		return errors.Wrapf(err, "failed to delete subnet %s", item.name)
	}

	o.Logger.Infof("Deleted Subnet %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroySubnets removes all subnet resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroySubnets() error {
	firstPassList, err := o.listSubnets()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(subnetTypeName, firstPassList.list())

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroySubnets: case <-ctx.Done()")
			return o.Context.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func() (bool, error) {
			err2 := o.deleteSubnet(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroySubnets: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(subnetTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroySubnets: found %s in pending items", item.name)
		}
		return errors.Errorf("destroySubnets: %d undeleted items pending", len(items))
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func() (bool, error) {
		secondPassList, err2 := o.listSubnets()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroySubnets: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroySubnets: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
