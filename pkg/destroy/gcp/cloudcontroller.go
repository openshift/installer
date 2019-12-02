package gcp

import (
	"fmt"

	compute "google.golang.org/api/compute/v1"
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

// listCloudControllerBackendServices returns backend services created by the cloud controller.
// It list all backend services matching the cloud controller name convention that contain
// only cluster instance groups.
func (o *ClusterUninstaller) listCloudControllerBackendServices(instanceGroups []cloudResource) ([]cloudResource, error) {
	urls := sets.NewString()
	for _, instanceGroup := range instanceGroups {
		urls.Insert(o.getInstanceGroupURL(instanceGroup))
	}
	filter := "name eq \"a[0-9a-f]{30,50}\""
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

// listCloudControllerTargetPools returns target pools created by the cloud controller.
// It list all target pools matching the cloud controller name convention that contain
// only cluster instances.
func (o *ClusterUninstaller) listCloudControllerTargetPools() ([]cloudResource, error) {
	filter := "name eq \"a[0-9a-f]{30,50}\""
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

// DiscoverCloudControllerLoadBalancerResources follows a similar procedure as:
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_internal.go#L222
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_external.go#L289
func (o *ClusterUninstaller) discoverCloudControllerLoadBalancerResources(loadBalancerName string) error {
	loadBalancerNameFilter := fmt.Sprintf("name eq \"%s\"", loadBalancerName)

	// Discover associated addresses: loadBalancerName
	found, err := o.listAddressesWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("address", found)

	// Discover associated firewall rules: loadBalancerName
	found, err = o.listFirewallsWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("firewall", found)

	// Discover associated firewall rules: loadBalancerName-hc
	filter := fmt.Sprintf("name eq \"%s-hc\"", loadBalancerName)
	found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("firewall", found)

	// Discover associated firewall rules: k8s-fw-loadBalancerName
	filter = fmt.Sprintf("name eq \"k8s-fw-%s\"", loadBalancerName)
	found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("firewall", found)

	// Discover associated firewall rules: k8s-loadBalancerName-http-hc
	filter = fmt.Sprintf("name eq \"k8s-%s-http-hc\"", loadBalancerName)
	found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("firewall", found)

	// Discover associated forwarding rules: loadBalancerName
	found, err = o.listForwardingRulesWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("forwardingrule", found)

	// Discover associated health checks: loadBalancerName
	found, err = o.listHealthChecksWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("healthcheck", found)

	// Discover associated http health checks: loadBalancerName
	found, err = o.listHTTPHealthChecksWithFilter("items(name),nextPageToken", loadBalancerNameFilter, nil)
	if err != nil {
		return err
	}
	o.insertPendingItems("httphealthcheck", found)

	return nil
}

// discoverCloudControllerResources finds resources associated with internal load balancers
// created by the kube cloud controller. It first finds instance groups associated with instances
// belonging to this cluster, then finds backend resources that point to these instance groups.
// For each of those backend services, resources like forwarding rules, firewalls, health checks and
// backend services are added to pendingItems
func (o *ClusterUninstaller) discoverCloudControllerResources() error {
	o.Logger.Debugf("Discovering cloud controller resources")
	errs := []error{}

	// Instance group related items
	instanceGroups, err := o.listCloudControllerInstanceGroups()
	if err != nil {
		return err
	}
	installerInstanceGroups, err := o.listInstanceGroups()
	if err != nil {
		return err
	}
	clusterInstanceGroups := append(instanceGroups, installerInstanceGroups...)
	if len(clusterInstanceGroups) != 0 {
		backends, err := o.listCloudControllerBackendServices(clusterInstanceGroups)
		if err != nil {
			return err
		}
		for _, backend := range backends {
			o.Logger.Debugf("Discovering cloud controller resources for %s", backend.name)
			err := o.discoverCloudControllerLoadBalancerResources(backend.name)
			if err != nil {
				errs = append(errs, err)
			}
		}
		o.insertPendingItems("backendservice", backends)
	}
	o.insertPendingItems("instancegroup", instanceGroups)

	// Target pool related items
	pools, err := o.listCloudControllerTargetPools()
	if err != nil {
		return err
	}
	for _, pool := range pools {
		o.Logger.Debugf("Discovering cloud controller resources for %s", pool.name)
		err := o.discoverCloudControllerLoadBalancerResources(pool.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	o.insertPendingItems("targetpool", pools)

	// cloudControllerUID related items
	if len(o.cloudControllerUID) > 0 {
		// Discover Cloud Controller health checks: k8s-cloudControllerUID-node
		filter := fmt.Sprintf("name eq \"k8s-%s-node\"", o.cloudControllerUID)
		found, err := o.listHealthChecksWithFilter("items(name),nextPageToken", filter, nil)
		if err != nil {
			return err
		}
		o.insertPendingItems("healthcheck", found)

		// Discover Cloud Controller http health checks: k8s-cloudControllerUID-node
		found, err = o.listHTTPHealthChecksWithFilter("items(name),nextPageToken", filter, nil)
		if err != nil {
			return err
		}
		o.insertPendingItems("httphealthcheck", found)

		// Discover Cloud Controller firewall rules: k8s-cloudControllerUID-node-hc, k8s-cloudControllerUID-node-http-hc
		filter = fmt.Sprintf("name eq \"k8s-%s-node-hc\"", o.cloudControllerUID)
		found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
		if err != nil {
			return err
		}
		o.insertPendingItems("firewall", found)

		filter = fmt.Sprintf("name eq \"k8s-%s-node-http-hc\"", o.cloudControllerUID)
		found, err = o.listFirewallsWithFilter("items(name),nextPageToken", filter, nil)
		if err != nil {
			return err
		}
		o.insertPendingItems("firewall", found)
	}

	return aggregateError(errs, 0)
}
