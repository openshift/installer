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

const loadBalancerTypeName = "load balancer"

// listLoadBalancers lists load balancers matching either name or tag in the IBM Cloud.
func (o *ClusterUninstaller) listLoadBalancers() (cloudResources, error) {
	var (
		lbIDs    []string
		lbID     string
		ctx      context.Context
		cancel   context.CancelFunc
		result   = make([]cloudResource, 0, 3)
		options  *vpcv1.GetLoadBalancerOptions
		lb       *vpcv1.LoadBalancer
		response *core.DetailedResponse
		err      error
	)

	if o.searchByTag {
		// Should we list by tag matching?
		lbIDs, err = o.listByTag(TagTypeLoadBalancer)
	} else {
		// Otherwise list will list by name matching.
		lbIDs, err = o.listLoadBalancersByName()
	}
	if err != nil {
		return nil, err
	}

	ctx, cancel = contextWithTimeout()
	defer cancel()

	for _, lbID = range lbIDs {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("listLoadBalancers: case <-ctx.Done()")
			return nil, ctx.Err() // we're cancelled, abort
		default:
		}

		options = o.vpcSvc.NewGetLoadBalancerOptions(lbID)

		lb, response, err = o.vpcSvc.GetLoadBalancerWithContext(ctx, options)
		if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
			// The load balancer could have been deleted just after a list was created.
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get load balancer (%s): err = %w, response = %v", lbID, err, response)
		}

		result = append(result, cloudResource{
			key:      *lb.ID,
			name:     *lb.Name,
			status:   *lb.ProvisioningStatus,
			typeName: loadBalancerTypeName,
			id:       *lb.ID,
		})
	}

	return cloudResources{}.insert(result...), nil
}

// listLoadBalancersByName list the load balancers matching by name in the IBM Cloud.
func (o *ClusterUninstaller) listLoadBalancersByName() ([]string, error) {
	var (
		ctx          context.Context
		cancel       context.CancelFunc
		options      *vpcv1.ListLoadBalancersOptions
		lbCollection *vpcv1.LoadBalancerCollection
		response     *core.DetailedResponse
		foundOne     = false
		result       = make([]string, 0, 3)
		lb           vpcv1.LoadBalancer
		err          error
	)

	o.Logger.Debugf("Listing load balancers by NAME")

	ctx, cancel = contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listLoadBalancersByName: case <-ctx.Done()")
		return nil, ctx.Err() // we're cancelled, abort
	default:
	}

	options = o.vpcSvc.NewListLoadBalancersOptions()
	// @WHY options.SetResourceGroupID(o.resourceGroupID)

	lbCollection, response, err = o.vpcSvc.ListLoadBalancersWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list load balancers: err = %w, response = %v", err, response)
	}

	for _, lb = range lbCollection.LoadBalancers {
		if strings.Contains(*lb.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listLoadBalancersByName: FOUND: %s, %s, %s", *lb.ID, *lb.Name, *lb.ProvisioningStatus)
			result = append(result, *lb.ID)
		}
	}
	if !foundOne {
		o.Logger.Debugf("listLoadBalancersByName: NO matching loadbalancers against: %s", o.InfraID)
		for _, loadbalancer := range lbCollection.LoadBalancers {
			o.Logger.Debugf("listLoadBalancersByName: loadbalancer: %s", *loadbalancer.Name)
		}
	}

	return result, nil
}

// deleteLoadBalancer deletes the load balancer specified.
func (o *ClusterUninstaller) deleteLoadBalancer(item cloudResource) error {
	var getOptions *vpcv1.GetLoadBalancerOptions
	var lb *vpcv1.LoadBalancer
	var response *core.DetailedResponse
	var err error

	ctx, cancel := contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteLoadBalancer: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	getOptions = o.vpcSvc.NewGetLoadBalancerOptions(item.id)
	lb, response, err = o.vpcSvc.GetLoadBalancer(getOptions)

	if err == nil && response.StatusCode == gohttp.StatusNoContent {
		return nil
	}
	if err != nil && response != nil && response.StatusCode == gohttp.StatusNotFound {
		// The resource is gone.
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Load Balancer %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == gohttp.StatusInternalServerError {
		o.Logger.Infof("deleteLoadBalancer: internal server error")
		return nil
	}
	if lb == nil {
		o.Logger.Debugf("deleteLoadBalancer: lb = %v", lb)
		o.Logger.Debugf("deleteLoadBalancer: response = %v", response)
		o.Logger.Debugf("deleteLoadBalancer: err = %v", err)
		o.Logger.Debugf("Rate and unhandled code, please investigate further")
		return nil
	}

	if *lb.ProvisioningStatus == vpcv1.LoadBalancerProvisioningStatusDeletePendingConst {
		o.Logger.Debugf("Waiting for load balancer %q to delete", item.name)
		return nil
	}

	deleteOptions := o.vpcSvc.NewDeleteLoadBalancerOptions(item.id)

	_, err = o.vpcSvc.DeleteLoadBalancerWithContext(ctx, deleteOptions)
	if err != nil {
		return fmt.Errorf("failed to delete load balancer %s: %w", item.name, err)
	}

	o.Logger.Infof("Deleted Load Balancer %q", item.name)
	o.deletePendingItems(item.typeName, []cloudResource{item})

	return nil
}

// destroyLoadBalancers removes all load balancer resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyLoadBalancers() error {
	firstPassList, err := o.listLoadBalancers()
	if err != nil {
		return err
	}

	if len(firstPassList.list()) == 0 {
		return nil
	}

	items := o.insertPendingItems(loadBalancerTypeName, firstPassList.list())

	ctx, cancel := contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyLoadBalancers: case <-ctx.Done()")
			return ctx.Err() // we're cancelled, abort
		default:
		}

		backoff := wait.Backoff{
			Duration: 15 * time.Second,
			Factor:   1.1,
			Cap:      leftInContext(ctx),
			Steps:    math.MaxInt32}
		err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
			err2 := o.deleteLoadBalancer(item)
			if err2 == nil {
				return true, err2
			}
			o.errorTracker.suppressWarning(item.key, err2, o.Logger)
			return false, err2
		})
		if err != nil {
			o.Logger.Fatal("destroyLoadBalancers: ExponentialBackoffWithContext (destroy) returns ", err)
		}
	}

	if items = o.getPendingItems(loadBalancerTypeName); len(items) > 0 {
		for _, item := range items {
			o.Logger.Debugf("destroyLoadBalancers: found %s in pending items", item.name)
		}
		return fmt.Errorf("destroyLoadBalancers: %d undeleted items pending", len(items))
	}

	select {
	case <-ctx.Done():
		o.Logger.Debugf("destroyLoadBalancers: case <-ctx.Done()")
		return ctx.Err() // we're cancelled, abort
	default:
	}

	backoff := wait.Backoff{
		Duration: 15 * time.Second,
		Factor:   1.1,
		Cap:      leftInContext(ctx),
		Steps:    math.MaxInt32}
	err = wait.ExponentialBackoffWithContext(ctx, backoff, func(context.Context) (bool, error) {
		secondPassList, err2 := o.listLoadBalancers()
		if err2 != nil {
			return false, err2
		}
		if len(secondPassList) == 0 {
			// We finally don't see any remaining instances!
			return true, nil
		}
		for _, item := range secondPassList {
			o.Logger.Debugf("destroyLoadBalancers: found %s in second pass", item.name)
		}
		return false, nil
	})
	if err != nil {
		o.Logger.Fatal("destroyLoadBalancers: ExponentialBackoffWithContext (list) returns ", err)
	}

	return nil
}
