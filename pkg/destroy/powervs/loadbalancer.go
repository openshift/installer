package powervs

import (
	"net/http"
	"strings"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const loadBalancerTypeName = "load balancer"

// listLoadBalancers lists subnets in the vpc.
func (o *ClusterUninstaller) listLoadBalancers() (cloudResources, error) {
	o.Logger.Debugf("Listing load balancers")

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("listLoadBalancers: case <-o.Context.Done()")
		return nil, o.Context.Err() // we're cancelled, abort
	default:
	}

	options := o.vpcSvc.NewListLoadBalancersOptions()

	resources, _, err := o.vpcSvc.ListLoadBalancersWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list load balancers")
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

	getOptions = o.vpcSvc.NewGetLoadBalancerOptions(item.id)
	lb, response, err = o.vpcSvc.GetLoadBalancer(getOptions)

	if err != nil && response != nil && response.StatusCode == http.StatusNotFound {
		// The resource is gone.
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted load balancer %q", item.name)
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

	o.Logger.Debugf("Deleting load balancer %q", item.name)

	ctx, _ := o.contextWithTimeout()

	select {
	case <-o.Context.Done():
		o.Logger.Debugf("deleteLoadBalancer: case <-o.Context.Done()")
		return o.Context.Err() // we're cancelled, abort
	default:
	}

	deleteOptions := o.vpcSvc.NewDeleteLoadBalancerOptions(item.id)
	_, err = o.vpcSvc.DeleteLoadBalancerWithContext(ctx, deleteOptions)

	if err != nil {
		return errors.Wrapf(err, "failed to delete load balancer %s", item.name)
	}

	return nil
}

// destroyLoadBalancers removes all load balancer resources that have a name prefixed
// with the cluster's infra ID.
func (o *ClusterUninstaller) destroyLoadBalancers() error {
	found, err := o.listLoadBalancers()
	if err != nil {
		return err
	}

	items := o.insertPendingItems(loadBalancerTypeName, found.list())

	ctx, _ := o.contextWithTimeout()

	for !o.timeout(ctx) {
		for _, item := range items {
			select {
			case <-o.Context.Done():
				o.Logger.Debugf("destroyLoadBalancers: case <-o.Context.Done()")
				return o.Context.Err() // we're cancelled, abort
			default:
			}

			if _, ok := found[item.key]; !ok {
				// This item has finished deletion.
				o.deletePendingItems(item.typeName, []cloudResource{item})
				o.Logger.Infof("Deleted load balancer %q", item.name)
				continue
			}
			err := o.deleteLoadBalancer(item)
			if err != nil {
				o.errorTracker.suppressWarning(item.key, err, o.Logger)
			}
		}

		items = o.getPendingItems(loadBalancerTypeName)
		if len(items) == 0 {
			break
		}
	}

	if items = o.getPendingItems(loadBalancerTypeName); len(items) > 0 {
		return errors.Errorf("destroyLoadBalancers: %d undeleted items pending", len(items))
	}
	return nil
}
