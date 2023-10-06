package powervs

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const loadBalancerTypeName = "load balancer"

// listLoadBalancers lists load balancers in the vpc.
func (o *ClusterUninstaller) listLoadBalancers() (cloudResources, error) {
	o.Logger.Debugf("Listing load balancers")

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("listLoadBalancers: case <-ctx.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListLoadBalancersOptions()

	resources, _, err := o.vpcSvc.ListLoadBalancersWithContext(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to list load balancers: %w", err)
	}

	var foundOne = false

	result := []cloudResource{}
	for _, loadbalancer := range resources.LoadBalancers {
		if strings.Contains(*loadbalancer.Name, o.InfraID) {
			foundOne = true
			o.Logger.Debugf("listLoadBalancers: FOUND: %s, %s, %s", *loadbalancer.ID, *loadbalancer.Name, *loadbalancer.ProvisioningStatus)
			result = append(result, cloudResource{
				key:      *loadbalancer.ID,
				name:     *loadbalancer.Name,
				status:   *loadbalancer.ProvisioningStatus,
				typeName: loadBalancerTypeName,
				id:       *loadbalancer.ID,
			})
		}
	}
	if !foundOne {
		o.Logger.Debugf("listLoadBalancers: NO matching loadbalancers against: %s", o.InfraID)
		for _, loadbalancer := range resources.LoadBalancers {
			o.Logger.Debugf("listLoadBalancers: loadbalancer: %s", *loadbalancer.Name)
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteLoadBalancer(item cloudResource) error {
	var getOptions *vpcv1.GetLoadBalancerOptions
	var lb *vpcv1.LoadBalancer
	var response *core.DetailedResponse
	var err error

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	select {
	case <-ctx.Done():
		o.Logger.Debugf("deleteLoadBalancer: case <-ctx.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	getOptions = o.vpcSvc.NewGetLoadBalancerOptions(item.id)
	lb, response, err = o.vpcSvc.GetLoadBalancer(getOptions)

	if err == nil && response.StatusCode == http.StatusNoContent {
		return nil
	}
	if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
		// The resource is gone.
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted Load Balancer %q", item.name)
		return nil
	}
	if err != nil && response != nil && response.StatusCode == http.StatusInternalServerError {
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

	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	for _, item := range items {
		select {
		case <-ctx.Done():
			o.Logger.Debugf("destroyLoadBalancers: case <-ctx.Done()")
			return o.Context.Err() // we're cancelled, abort
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
