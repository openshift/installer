package ibmcloud

import (
	"net/http"
	"strings"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/pkg/errors"
)

const loadBalancerTypeName = "load balancer"

// listLoadBalancers lists subnets in the vpc
func (o *ClusterUninstaller) listLoadBalancers() (cloudResources, error) {
	o.Logger.Debugf("Listing load balancers")
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewListLoadBalancersOptions()
	resources, _, err := o.vpcSvc.ListLoadBalancersWithContext(ctx, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list load balancers")
	}

	result := []cloudResource{}
	for _, loadbalancer := range resources.LoadBalancers {
		if strings.Contains(*loadbalancer.Name, o.InfraID) {
			result = append(result, cloudResource{
				key:      *loadbalancer.ID,
				name:     *loadbalancer.Name,
				status:   *loadbalancer.ProvisioningStatus,
				typeName: loadBalancerTypeName,
				id:       *loadbalancer.ID,
			})
		}
	}

	return cloudResources{}.insert(result...), nil
}

func (o *ClusterUninstaller) deleteLoadBalancer(item cloudResource) error {
	if item.status == vpcv1.LoadBalancerProvisioningStatusDeletePendingConst {
		o.Logger.Debugf("Waiting for load balancer %q to delete", item.name)
		return nil
	}

	o.Logger.Debugf("Deleting load balancer %q", item.name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	options := o.vpcSvc.NewDeleteLoadBalancerOptions(item.id)
	details, err := o.vpcSvc.DeleteLoadBalancerWithContext(ctx, options)

	if err != nil && details != nil && details.StatusCode == http.StatusNotFound {
		// The resource is gone.
		o.deletePendingItems(item.typeName, []cloudResource{item})
		o.Logger.Infof("Deleted load balancer %q", item.name)
		return nil
	}

	if err != nil && details != nil && details.StatusCode != http.StatusNotFound {
		return errors.Wrapf(err, "Failed to delete load balancer %s", item.name)
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
	for _, item := range items {
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

	if items = o.getPendingItems(loadBalancerTypeName); len(items) > 0 {
		return errors.Errorf("%d items pending", len(items))
	}
	return nil
}
