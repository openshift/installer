package gcp

import (
	"fmt"

	compute "google.golang.org/api/compute/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

// listCloudControllerTargetPools returns target pools created by the cloud controller.
// It list all target pools matching the cloud controller name convention that contain
// only cluster instances.
func (o *ClusterUninstaller) listCloudControllerTargetPools() ([]cloudResource, error) {
	filter := fmt.Sprintf("name eq \"a[0-9a-f]{30,50}\"", o.cloudControllerUID)
	return o.listTargetPoolsWithFilter("items(name,instances),nextPageToken", filter, func(pool *compute.TargetPool) bool {
		if len(pool.Instances) == 0 {
			return false
		}
		for _, instanceURL := range pool.Instances {
			name, _ := o.getInstanceNameAndZone(instanceURL)
			if !o.isClusterResource(name) {
				return false
			}
		}
		return true
	})
}

// deleteExternalLoadBalancer follows a similar cleanup procedure as:
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_external.go#L289
func (o *ClusterUninstaller) deleteExternalLoadBalancer(loadBalancerName string) error {
	loadBalancerNameFilter := fmt.Sprintf("name eq \"%s\"", loadBalancerName)
	errs := []error{}

	// Delete associated address
	found, err := o.listAddressesWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	items := o.insertPendingItems("address", found)
	for _, item := range items {
		err := o.deleteAddress(item)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Delete associated firwall rules
	filter := fmt.Sprintf("name eq \"k8s-fw-%s\"", loadBalancerName)
	found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
	if err != nil {
		return err
	}
	items = o.insertPendingItems("firewall", found)
	filter = fmt.Sprintf("name eq \"k8s-%s-http-hc\"", loadBalancerName)
	found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
	if err != nil {
		return err
	}
	items = o.insertPendingItems("firewall", found)
	for _, item := range items {
		err = o.deleteFirewall(item)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Delete associated forwarding rules
	found, err = o.listForwardingRulesWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	items = o.insertPendingItems("forwardingrule", found)
	for _, item := range items {
		err = o.deleteForwardingRule(item)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Delete associated http health checks
	found, err = o.listHttpHealthChecksWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	items = o.insertPendingItems("httphealthcheck", found)
	for _, item := range items {
		err = o.deleteHttpHealthCheck(item)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// Delete associated target pools
	found, err = o.listTargetPoolsWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	items = o.insertPendingItems("targetpool", found)
	for _, item := range items {
		err := o.deleteTargetPool(item)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}

// destroyCloudControllerExternalLBs removes resources associated with external load balancers
// created by the kube cloud controller. It first finds target pools associated with instances
// belonging to this cluster. For each of those target pools, it removes resources like
// addresses, forwarding rules, firewalls, health checks and target pools.
func (o *ClusterUninstaller) destroyCloudControllerExternalLBs() error {
	errs := []error{}
	pools, err := o.listCloudControllerTargetPools()
	if err != nil {
		return err
	}
	for _, pool := range pools {
		err := o.deleteExternalLoadBalancer(pool.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(o.cloudControllerUID) != 0 {
		// Delete Cloud Controller http health checks
		filter := fmt.Sprintf("name eq \"k8s-%s-node\"", o.cloudControllerUID)
		found, err := o.listHttpHealthChecksWithFilter("items(name),nextPageToken", filter, nil)
		if err != nil {
			return err
		}
		items := o.insertPendingItems("httphealthcheck", found)
		for _, item := range items {
			err = o.deleteHttpHealthCheck(item)
			if err != nil {
				errs = append(errs, err)
			}
		}
		// Delete Cloud Controller firewall rules
		filter = fmt.Sprintf("name eq \"k8s-%s-node-hc\"", o.cloudControllerUID)
		found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
		if err != nil {
			return err
		}
		items = o.insertPendingItems("firewall", found)
		for _, item := range items {
			err = o.deleteFirewall(item)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return aggregateError(errs, 0)
}
