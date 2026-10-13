package gcp

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
	"k8s.io/apimachinery/pkg/util/sets"
)

type resourceFilterFunc func(string) bool

func (o *ClusterUninstaller) createLoadBalancerFilterFunc(loadBalancerName string) resourceFilterFunc {
	return func(itemName string) bool {
		return strings.HasPrefix(itemName, loadBalancerName)
	}
}

// listCloudControllerInstanceGroups returns instance groups created by the cloud controller.
// It list all instance groups matching the cloud controller name convention.
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_naming.go#L33-L40
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_clusterid.go#L210-L238
func (o *ClusterUninstaller) listCloudControllerInstanceGroups(ctx context.Context) ([]cloudResource, error) {
	return o.listInstanceGroupsWithFilter(ctx, "items/*/instanceGroups(name,selfLink,zone),nextPageToken", func(itemName string) bool {
		// TODO: Why does this have an extra - ??
		return itemName == fmt.Sprintf("k8s-ig--%s", o.cloudControllerUID)
	})
}

// listCloudControllerBackendServices returns backend services created by the cloud controller.
// It list all backend services matching the cloud controller name convention that contain
// only cluster instance groups.
func (o *ClusterUninstaller) listCloudControllerBackendServices(ctx context.Context, instanceGroups []cloudResource) ([]cloudResource, error) {
	return o.listBackendServicesWithFilter(ctx, regionBackendServiceResource, "items(name,backends),nextPageToken",
		func(item *compute.BackendService) bool {
			filter := regexp.MustCompile(`a[0-9a-f]{30,50}`)
			if !filter.MatchString(item.Name) {
				return false
			}

			// During cluster creation, the backends may be empty in the Backend Service. When this is the case,
			// the cluster uninstaller does not know if the resource should actually be considered for deletion.
			// Search for a resource that should be related to the Backend Service, a Firewall Rule. The firewall
			// rule should have tags that include the cluster ID in the name.
			// TODO: when/if the shared tag is used this could pose future problems (until the backend service can be tagged).
			// TODO: When backend services can be tagged, use the resource manager to get tags for the backend service.

			fwList, err := o.computeSvc.Firewalls.List(o.ProjectID).Fields(googleapi.Field("items(name,targetTags),nextPageToken")).Context(ctx).Do()
			if err != nil {
				o.Logger.Debugf("failed to list firewall rules associated with backend service %s: %v", item.Name, err)
				return false
			}
			for _, fw := range fwList.Items {
				if strings.Contains(fw.Name, item.Name) {
					for _, tag := range fw.TargetTags {
						// These tags are in the form {o.ClusterID}-worker
						if strings.Contains(tag, o.ClusterID) {
							return true
						}
					}
				}
			}

			urls := sets.Set[string]{}
			for _, instanceGroup := range instanceGroups {
				urls.Insert(instanceGroup.url)
			}
			if len(item.Backends) == 0 {
				return false
			}

			// If the backends for the Backend Service are not empty, compare them to the instance
			// group urls to determine if they are associated with each other. A match indicates that
			// this Backend Service is part of the cluster being destroyed.
			for _, backend := range item.Backends {
				if !urls.Has(backend.Group) {
					return false
				}
			}
			return true
		})
}

// listCloudControllerTargetPools returns target pools created by the cloud controller or owned by the cloud controller.
// It lists all target pools matching the cloud controller name convention that contain
// only cluster instances or cluster instances that were owned by the cluster.
func (o *ClusterUninstaller) listCloudControllerTargetPools(ctx context.Context, instances []cloudResource) ([]cloudResource, error) {
	return o.listTargetPoolsWithFilter(ctx, "items(name,instances),nextPageToken", func(item *compute.TargetPool) bool {
		filter := regexp.MustCompile(`a[0-9a-f]{30,50}`)
		if !filter.MatchString(item.Name) || len(item.Instances) == 0 {
			return false
		}
		for _, instanceURL := range item.Instances {
			name, _ := o.getInstanceNameAndZone(instanceURL)
			if !o.isClusterResource(name) {
				foundClusterResource := false
				for _, instance := range instances {
					if instance.name == name {
						foundClusterResource = true
						break
					}
				}

				if !foundClusterResource {
					o.Logger.Debugf("Skipping target pool instance %s because it is not a cluster resource", item.Name)
					return false
				}
			}
		}
		return true
	})
}

// DiscoverCloudControllerLoadBalancerResources follows a similar procedure as:
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_internal.go#L222
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_external.go#L289
func (o *ClusterUninstaller) discoverCloudControllerLoadBalancerResources(ctx context.Context, loadBalancerName string) error {
	loadBalancerFilterFunc := o.createLoadBalancerFilterFunc(loadBalancerName)

	// Discover associated addresses: loadBalancerName
	found, err := o.listAddressesWithFilter(ctx, regionalAddressResource, "items(name),nextPageToken", loadBalancerFilterFunc)
	if err != nil {
		return err
	}
	o.insertPendingItems(regionalAddressResource, found)

	// Discover associated firewall rules:
	// 1. loadBalancerName
	// 2. loadBalancerName-hc
	// 3. k8s-fw-loadBalancerName
	// 4. k8s-loadBalancerName-http-hc
	// 5. k8s-%s-node-hc
	// 6. k8s-%s-node-http-hc
	found, err = o.listFirewallsWithFilter(ctx, "items(name,targetTags),nextPageToken", o.firewallFilterFunc)
	if err != nil {
		return err
	}
	o.insertPendingItems(firewallResourceName, found)

	// Discover associated forwarding rules: loadBalancerName
	found, err = o.listForwardingRulesWithFilter(ctx, regionForwardingRuleResource, "items(name),nextPageToken", loadBalancerFilterFunc)
	if err != nil {
		return err
	}
	o.insertPendingItems(regionForwardingRuleResource, found)

	// Discover associated target tcp proxies: loadBalancerName
	found, err = o.listTargetTCPProxiesWithFilter(ctx, globalTargetTCPProxyResource, "items(name),nextPageToken", loadBalancerFilterFunc)
	if err != nil {
		return err
	}
	o.insertPendingItems(globalTargetTCPProxyResource, found)

	// Discover associated health checks: loadBalancerName
	found, err = o.listHealthChecksWithFilter(ctx, regionHealthCheckResource, "items(name),nextPageToken", loadBalancerFilterFunc)
	if err != nil {
		return err
	}
	o.insertPendingItems(regionHealthCheckResource, found)

	// Discover associated health checks: loadBalancerName - GLOBAL
	found, err = o.listHealthChecksWithFilter(ctx, globalHealthCheckResource, "items(name),nextPageToken", loadBalancerFilterFunc)
	if err != nil {
		return err
	}
	o.insertPendingItems(globalHealthCheckResource, found)

	// Discover associated http health checks: loadBalancerName
	found, err = o.listHTTPHealthChecksWithFilter(ctx, "items(name),nextPageToken", loadBalancerFilterFunc)
	if err != nil {
		return err
	}
	o.insertPendingItems(httpHealthCheckResourceName, found)

	return nil
}

// discoverCloudControllerResources finds resources associated with internal load balancers
// created by the kube cloud controller. It first finds instance groups associated with instances
// belonging to this cluster, then finds backend resources that point to these instance groups.
// For each of those backend services, resources like forwarding rules, firewalls, health checks and
// backend services are added to pendingItems
func (o *ClusterUninstaller) discoverCloudControllerResources(ctx context.Context) error {
	o.Logger.Debugf("Discovering cloud controller resources")
	errs := []error{}

	// Instance group related items
	instanceGroups, err := o.listCloudControllerInstanceGroups(ctx)
	if err != nil {
		return err
	}
	installerInstanceGroups, err := o.listInstanceGroups(ctx)
	if err != nil {
		return err
	}
	clusterInstanceGroups := append(instanceGroups, installerInstanceGroups...)
	if len(clusterInstanceGroups) != 0 {
		backends, err := o.listCloudControllerBackendServices(ctx, clusterInstanceGroups)
		if err != nil {
			return err
		}
		for _, backend := range backends {
			o.Logger.Debugf("Discovering cloud controller resources for %s", backend.name)
			err := o.discoverCloudControllerLoadBalancerResources(ctx, backend.name)
			if err != nil {
				errs = append(errs, err)
			}
		}
		o.insertPendingItems(regionBackendServiceResource, backends)
	}
	o.insertPendingItems(instanceGroupResourceName, instanceGroups)

	// Get a list of known cluster instances
	instances, err := o.listInstances(ctx)
	if err != nil {
		return err
	}

	// Target pool related items
	pools, err := o.listCloudControllerTargetPools(ctx, instances)
	if err != nil {
		return err
	}
	for _, pool := range pools {
		o.Logger.Debugf("Discovering cloud controller resources for %s", pool.name)
		err := o.discoverCloudControllerLoadBalancerResources(ctx, pool.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	o.insertPendingItems(targetPoolResourceName, pools)

	// cloudControllerUID related items
	if len(o.cloudControllerUID) > 0 {
		// Discover Cloud Controller health checks: k8s-cloudControllerUID-node
		found, err := o.listHealthChecksWithFilter(ctx, regionHealthCheckResource, "items(name),nextPageToken",
			o.createLoadBalancerFilterFunc(fmt.Sprintf("k8s-%s-node", o.cloudControllerUID)),
		)
		if err != nil {
			return err
		}
		o.insertPendingItems(regionHealthCheckResource, found)

		// Discover Cloud Controller http health checks: k8s-cloudControllerUID-node
		found, err = o.listHTTPHealthChecksWithFilter(ctx, "items(name),nextPageToken",
			o.createLoadBalancerFilterFunc(fmt.Sprintf("k8s-%s-node", o.cloudControllerUID)),
		)
		if err != nil {
			return err
		}
		o.insertPendingItems(httpHealthCheckResourceName, found)
	}

	return aggregateError(errs, 0)
}
