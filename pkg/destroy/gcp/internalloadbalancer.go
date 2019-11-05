package gcp

import (
	"fmt"

	compute "google.golang.org/api/compute/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// listCloudControllerInstanceGroups returns instance groups created by the cloud controller.
// It list all instance groups matching the cloud controller name convention.
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_naming.go#L33-L40
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_clusterid.go#L210-L238
func (o *ClusterUninstaller) listCloudControllerInstanceGroups() ([]cloudResource, error) {
	filter := fmt.Sprintf("name eq \"k8s-ig--%s\"", o.cloudControllerUID)
	return o.listInstanceGroupsWithFilter("items/*/instanceGroups(name,zone),nextPageToken", filter, nil)
}

// listBackendServicesForInstanceGroups returns backend services created by the cloud controller.
// It list all backend services matching the cloud controller name convention that contain
// only cluster instance groups.
func (o *ClusterUninstaller) listBackendServicesForInstanceGroups(instanceGroups []cloudResource) ([]cloudResource, error) {
	urls := sets.NewString()
	for _, instanceGroup := range instanceGroups {
		urls.Insert(o.getInstanceGroupURL(instanceGroup))
	}
	filter := fmt.Sprintf("name eq \"a[0-9a-f]{30,50}\"", o.cloudControllerUID)
	return o.listBackendServicesWithFilter("items(name,backends),nextPageToken", filter, func(item *compute.BackendService) bool {
		if len(item.Backends) == 0 {
			return false
		}
		for _, backend := range item.Backends {
			if !urls.Has(backend.Group) {
				return false
			}
		}
		return true
	})
}

// deleteInternalLoadBalancer follows a similar cleanup procedure as:
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_internal.go#L222
func (o *ClusterUninstaller) deleteInternalLoadBalancer(loadBalancerName string) error {
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
	found, err = o.listFirewallsWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	items = o.insertPendingItems("firewall", found)
	filter := fmt.Sprintf("name eq \"%s-hc\"", loadBalancerName)
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

	// Delete associated backend services
	found, err = o.listBackendServicesWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	items = o.insertPendingItems("backendservice", found)
	for _, item := range items {
		err := o.deleteBackendService(item)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return utilerrors.NewAggregate(errs)
}

// destroyCloudControllerInternalLBs removes resources associated with internal load balancers
// created by the kube cloud controller. It first finds instance groups associated with instances
// belonging to this cluster, then finds backend resources that point to these instance groups.
// For each of those backend services, resources like forwarding rules, firewalls, health checks and
// backend services are removed.
func (o *ClusterUninstaller) destroyCloudControllerInternalLBs() error {
	errs := []error{}
	groups, err := o.listCloudControllerInstanceGroups()
	if err != nil {
		return err
	}
	if len(groups) == 0 {
		return nil
	}
	backends, err := o.listBackendServicesForInstanceGroups(groups)
	if err != nil {
		return err
	}
	for _, backend := range backends {
		err := o.deleteInternalLoadBalancer(backend.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, group := range groups {
		err := o.deleteInstanceGroup(group)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(o.cloudControllerUID) > 0 {
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
