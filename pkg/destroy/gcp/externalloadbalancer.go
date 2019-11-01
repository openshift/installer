package gcp

import (
	"fmt"

	compute "google.golang.org/api/compute/v1"
)

// deleteExternalLoadBalancer follows a similar cleanup procedure as:
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_external.go#L289
// TODO: cleanup nodes health check (using clusterid)
func (o *ClusterUninstaller) deleteExternalLoadBalancer(loadBalancerName string) error {

	// Remove health checks from target pool so we can delete the health checks first
	// and leave the target pool for last. The target pool is the anchor for external load
	// balancers and without it we do not have a name to use when deleting them.
	if err := o.clearTargetPoolHealthChecks(loadBalancerName); err != nil {
		return err
	}
	item := cloudResource{
		key:      loadBalancerName,
		name:     loadBalancerName,
		typeName: "address",
	}
	if err := o.deleteAddress(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName,
		name:     loadBalancerName,
		typeName: "forwardingrule",
	}
	if err := o.deleteForwardingRule(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      fmt.Sprintf("k8s-fw-%s", loadBalancerName),
		name:     fmt.Sprintf("k8s-fw-%s", loadBalancerName),
		typeName: "firewall",
	}
	if err := o.deleteFirewall(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      fmt.Sprintf("k8s-%s-http-hc", loadBalancerName),
		name:     fmt.Sprintf("k8s-%s-http-hc", loadBalancerName),
		typeName: "firewall",
	}
	if err := o.deleteFirewall(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName,
		name:     loadBalancerName,
		typeName: "httphealthcheck",
	}
	if err := o.deleteHTTPHealthCheck(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName,
		name:     loadBalancerName,
		typeName: "targetpool",
	}
	if err := o.deleteTargetPool(item); err != nil {
		return err
	}
	return nil
}

// getExternalLBTargetPools returns all target pools that point to instances in the cluster
func (o *ClusterUninstaller) getExternalLBTargetPools() ([]cloudResource, error) {
	return o.listTargetPoolsWithFilter("items(name,instances),nextPageToken", "", func(pool *compute.TargetPool) bool {
		if len(pool.Instances) == 0 {
			return false
		}
		for _, instanceURL := range pool.Instances {
			name, _ := o.getInstanceNameAndZone(instanceURL)
			if !o.isClusterResource(name) {
				return false
			}
		}
		o.Logger.Debugf("Found external load balancer target pool: %s", pool.Name)
		return true
	})
}

// destroyCloudControllerExternalLBs removes resources associated with external load balancers
// created by the kube cloud controller. It first finds target pools associated with instances
// belonging to this cluster. For each of those target pools, it removes resources like
// addresses, forwarding rules, firewalls, health checks and target pools.
func (o *ClusterUninstaller) destroyCloudControllerExternalLBs() error {
	pools, err := o.getExternalLBTargetPools()
	if err != nil {
		return err
	}

	found := cloudResources{}
	errs := []error{}
	for _, pool := range pools {
		found.insert(pool)
		err := o.deleteExternalLoadBalancer(pool.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("externallb", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted external load balancer %s", item.name)
	}
	if len(o.cloudControllerUID) != 0 {
		item := cloudResource {
			key:      fmt.Sprintf("k8s-%s-node", o.cloudControllerUID),
			name:     fmt.Sprintf("k8s-%s-node", o.cloudControllerUID),
			typeName: "httphealthcheck",
		}
		if err := o.deleteHTTPHealthCheck(item, true); err != nil {
			return err
		}
		item = cloudResource {
			key:      fmt.Sprintf("k8s-%s-node-hc", o.cloudControllerUID),
			name:     fmt.Sprintf("k8s-%s-node-hc", o.cloudControllerUID),
			typeName: "firewall",
		}
		if err := o.deleteFirewall(item, true); err != nil {
			return err
		}
	}
	return aggregateError(errs, len(found))
}
