package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

type nameAndURL struct {
	name string
	url  string
}

func (n nameAndURL) String() string {
	return fmt.Sprintf("Name: %s, URL: %s\n", n.name, n.url)
}

func (o *ClusterUninstaller) listFirewalls() ([]string, error) {
	o.Logger.Debugf("Listing firewall rules")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Firewalls.List(o.ProjectID).Fields("items(name)").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.FirewallList) error {
		for _, firewall := range list.Items {
			o.Logger.Debugf("Found firewall rule: %s", firewall.Name)
			result = append(result, firewall.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list firewall rules")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteFirewall(name string) error {
	o.Logger.Debugf("Deleting firewall rule %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.Firewalls.Delete(o.ProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete firewall %s", name)
	}
	return nil
}

// destroyFirewalls removes all firewall resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyFirewalls() error {
	firewalls, err := o.listFirewalls()
	if err != nil {
		return err
	}
	for _, firewall := range firewalls {
		err = o.deleteFirewall(firewall)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listAddresses() ([]string, error) {
	o.Logger.Debugf("Listing addresses")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Addresses.List(o.ProjectID, o.Region).Fields("items(name)").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.AddressList) error {
		for _, address := range list.Items {
			o.Logger.Debugf("Found address: %s", address.Name)
			result = append(result, address.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list addresses")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteAddress(name string) error {
	o.Logger.Debugf("Deleting address %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete address %s", name)
	}
	return nil
}

// destroyAddresses removes all address resources that have a name prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyAddresses() error {
	addresses, err := o.listAddresses()
	if err != nil {
		return err
	}
	for _, address := range addresses {
		err = o.deleteAddress(address)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listForwardingRules() ([]string, error) {
	o.Logger.Debugf("Listing forwarding rules")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.ForwardingRules.List(o.ProjectID, o.Region).Fields("items(name)").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.ForwardingRuleList) error {
		for _, forwardingRule := range list.Items {
			o.Logger.Debugf("Found forwarding rule: %s", forwardingRule.Name)
			result = append(result, forwardingRule.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list forwarding rules")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteForwardingRule(name string) error {
	o.Logger.Debugf("Deleting forwarding rule %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.ForwardingRules.Delete(o.ProjectID, o.Region, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete forwarding rule %s", name)
	}
	return nil
}

// destroyForwardingRules removes all forwarding rules with a name prefixed with the
// cluster's infra ID.
func (o *ClusterUninstaller) destroyForwardingRules() error {
	forwardingRules, err := o.listForwardingRules()
	if err != nil {
		return err
	}
	for _, forwardingRule := range forwardingRules {
		err = o.deleteForwardingRule(forwardingRule)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listBackendServices() ([]string, error) {
	return o.listBackendServicesWithFilter("items(name)", o.clusterIDFilter(), nil)
}

// listBackendServicesWithFilter lists backend services in the project that satisfy the filter criteria.
// The fields parameter specifies which fields should be returned in the result, the filter string contains
// a filter string passed to the API to filter results. The filterFunc is a client-side filtering function
// that determines whether a particular result should be returned or not.
func (o *ClusterUninstaller) listBackendServicesWithFilter(fields string, filter string, filterFunc func(*compute.BackendService) bool) ([]string, error) {
	o.Logger.Debugf("Listing backend services")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.BackendServices.List(o.ProjectID).Fields(googleapi.Field(fields))
	if len(filter) > 0 {
		req = req.Filter(filter)
	}
	err := req.Pages(ctx, func(list *compute.BackendServiceList) error {
		for _, backendService := range list.Items {
			if filterFunc == nil || filterFunc != nil && filterFunc(backendService) {
				o.Logger.Debugf("Found backend service: %s", backendService.Name)
				result = append(result, backendService.Name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list backend services")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteBackendService(name string) error {
	o.Logger.Debugf("Deleting backend service %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.BackendServices.Delete(o.ProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete backend service %s", name)
	}
	return nil
}

// destroyBackendServices removes backend services with a name prefixed by the
// cluster's infra ID.
func (o *ClusterUninstaller) destroyBackendServices() error {
	backendServices, err := o.listBackendServices()
	if err != nil {
		return err
	}
	for _, backendService := range backendServices {
		err = o.deleteBackendService(backendService)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listHealthChecks() ([]string, error) {
	o.Logger.Debugf("Listing health checks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.HealthChecks.List(o.ProjectID).Fields("items(name)").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.HealthCheckList) error {
		for _, healthCheck := range list.Items {
			o.Logger.Debugf("Found health check: %s", healthCheck.Name)
			result = append(result, healthCheck.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list health checks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHealthCheck(name string) error {
	o.Logger.Debugf("Deleting health check %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.HealthChecks.Delete(o.ProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete health check %s", name)
	}
	return nil
}

// destroyHealthChecks removes all health check resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyHealthChecks() error {
	healthChecks, err := o.listHealthChecks()
	if err != nil {
		return err
	}
	for _, healthCheck := range healthChecks {
		err = o.deleteHealthCheck(healthCheck)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listTargetPools() ([]string, error) {
	return o.listTargetPoolsWithFilter("items(name)", o.clusterIDFilter(), nil)
}

// listTargetPoolsWithFilter lists target pools in the project. The field parameter allows
// specifying which fields to return. The filter parameter specifies a server-side filter for the
// GCP API (preferred). The filterFunc specifies a client-side filtering function for each TargetPool.
func (o *ClusterUninstaller) listTargetPoolsWithFilter(field string, filter string, filterFunc func(*compute.TargetPool) bool) ([]string, error) {
	o.Logger.Debugf("Listing target pools")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.TargetPools.List(o.ProjectID, o.Region).Fields(googleapi.Field(field)).Filter(filter)
	err := req.Pages(ctx, func(list *compute.TargetPoolList) error {
		for _, targetPool := range list.Items {
			if filterFunc == nil || (filterFunc != nil && filterFunc(targetPool)) {
				o.Logger.Debugf("Found target pool: %s", targetPool.Name)
				result = append(result, targetPool.Name)
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list target pools")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteTargetPool(name string) error {
	o.Logger.Debugf("Deleting target pool %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.TargetPools.Delete(o.ProjectID, o.Region, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete target pool %s", name)
	}
	return nil
}

// destroyTargetPools removes target pools created for external load balancers that have
// a name that starts with the cluster infra ID. These are load balancers created by the
// installer or cluster operators.
func (o *ClusterUninstaller) destroyTargetPools() error {
	targetPools, err := o.listTargetPools()
	if err != nil {
		return err
	}
	for _, targetPool := range targetPools {
		err = o.deleteTargetPool(targetPool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listSubNetworks() ([]string, error) {
	o.Logger.Debugf("Listing subnetworks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Subnetworks.List(o.ProjectID, o.Region).Fields("items(name)").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.SubnetworkList) error {
		for _, subNetwork := range list.Items {
			o.Logger.Debugf("Found subnetwork: %s", subNetwork.Name)
			result = append(result, subNetwork.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list subnetworks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteSubNetwork(name string) error {
	o.Logger.Debugf("Deleting subnetwork %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.Subnetworks.Delete(o.ProjectID, o.Region, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete subnetwork %s", name)
	}
	return nil
}

// destroySubNetworks removes all subnetwork resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroySubNetworks() error {
	subNetworks, err := o.listSubNetworks()
	if err != nil {
		return err
	}
	for _, subNetwork := range subNetworks {
		err = o.deleteSubNetwork(subNetwork)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listNetworks() ([]nameAndURL, error) {
	o.Logger.Debugf("Listing networks")
	result := []nameAndURL{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Networks.List(o.ProjectID).Fields("items(name,selfLink)").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.NetworkList) error {
		for _, network := range list.Items {
			o.Logger.Debugf("Found network: %s", network.Name)
			result = append(result, nameAndURL{
				name: network.Name,
				url:  network.SelfLink,
			})
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list networks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteNetwork(name string) error {
	o.Logger.Debugf("Deleting network %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.Networks.Delete(o.ProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete network %s", name)
	}
	return nil
}

// destroyNetworks removes all vpc network resources prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyNetworks() error {
	networks, err := o.listNetworks()
	if err != nil {
		return err
	}
	for _, network := range networks {
		// destroy any network routes that are not named with the infra ID
		routes, err := o.listNetworkRoutes(network.url)
		if err != nil {
			return err
		}
		for _, route := range routes {
			if err = o.deleteRoute(route); err != nil {
				o.Logger.Debug("error deleting route %s: %v", route, err)
			}
		}

		err = o.deleteNetwork(network.name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listNetworkRoutes(networkURL string) ([]string, error) {
	return o.listRoutesWithFilter(fmt.Sprintf("network eq %q", networkURL))
}

func (o *ClusterUninstaller) listRoutes() ([]string, error) {
	return o.listRoutesWithFilter(o.clusterIDFilter())
}

func (o *ClusterUninstaller) listRoutesWithFilter(filter string) ([]string, error) {
	o.Logger.Debugf("Listing routes")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Routes.List(o.ProjectID).Fields("items(name)").Filter(filter)
	err := req.Pages(ctx, func(list *compute.RouteList) error {
		for _, route := range list.Items {
			o.Logger.Debugf("Found route: %s", route.Name)
			result = append(result, route.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list routes")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteRoute(name string) error {
	o.Logger.Debugf("Deleting route %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.Routes.Delete(o.ProjectID, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete route %s", name)
	}
	return nil
}

// destroyRutes removes all route resources that have a name prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyRoutes() error {
	routes, err := o.listRoutes()
	if err != nil {
		return err
	}
	for _, route := range routes {
		err = o.deleteRoute(route)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *ClusterUninstaller) listRouters() ([]string, error) {
	o.Logger.Debug("Listing routers")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Routers.List(o.ProjectID, o.Region).Fields("items(name)").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.RouterList) error {
		for _, router := range list.Items {
			o.Logger.Debugf("Found router: %s", router.Name)
			result = append(result, router.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list routers")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteRouter(name string) error {
	o.Logger.Debugf("Deleting router %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	_, err := o.computeSvc.Routers.Delete(o.ProjectID, o.Region, name).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		return errors.Wrapf(err, "failed to delete router %s", name)
	}
	return nil
}

// destroyRouters removes all router resources that have a name prefixed with the
// cluster's infra ID
func (o *ClusterUninstaller) destroyRouters() error {
	routers, err := o.listRouters()
	if err != nil {
		return err
	}
	for _, router := range routers {
		err = o.deleteRouter(router)
		if err != nil {
			return err
		}
	}
	return nil
}

// getInternalLBInstanceGroups finds instance groups created for kube cloud controller
// internal load balancers. They should be named "k8s-ig--{clusterid}":
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_naming.go#L33-L40
// where clusterid is an 8-char id generated and saved in a configmap named
// "ingress-uid" in kube-system, key "uid":
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_clusterid.go#L210-L238
// TODO: Have the installer generate this UID, save it to metadata and populate it in the cluster. With that, we can just directly get
// an instanceGroup named "k8s-ig--{clusterid}"
// For now, we look into each instance group and determine if it contains instances prefixed with the cluster's infra ID
func (o *ClusterUninstaller) getInternalLBInstanceGroups() ([]nameAndZone, error) {
	candidates, err := o.listInstanceGroupsWithFilter("name eq \"k8s-ig--.*\"")
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, nil
	}

	igName := ""
	for _, ig := range candidates {
		instances, err := o.listInstanceGroupInstances(ig)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get internal LB instance group instances for %s", ig.name)
		}
		if o.areAllClusterInstances(instances) {
			igName = ig.name
			break
		}
	}
	igs := []nameAndZone{}
	if len(igName) > 0 {
		for _, ig := range candidates {
			if ig.name == igName {
				igs = append(igs, ig)
			}
		}
	}

	return igs, nil
}

func (o *ClusterUninstaller) areAllClusterInstances(instances []nameAndZone) bool {
	for _, instance := range instances {
		if !o.isClusterResource(instance.name) {
			return false
		}
	}
	return true
}

func (o *ClusterUninstaller) getInstanceGroupURL(ig nameAndZone) string {
	return fmt.Sprintf("%s%s/zones/%s/instanceGroups/%s", o.computeSvc.BasePath, o.ProjectID, ig.zone, ig.name)
}

func (o *ClusterUninstaller) listBackendServicesForInstanceGroups(igs []nameAndZone) ([]string, error) {
	urls := sets.NewString()
	for _, ig := range igs {
		urls.Insert(o.getInstanceGroupURL(ig))
	}
	return o.listBackendServicesWithFilter("items(name,backends)", "", func(item *compute.BackendService) bool {
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
// TODO: add cleanup for shared mode resources (determine if it's supported in 4.2)
func (o *ClusterUninstaller) deleteInternalLoadBalancer(clusterID string, loadBalancerName string) error {
	if err := o.deleteAddress(loadBalancerName); err != nil {
		return err
	}
	if err := o.deleteForwardingRule(loadBalancerName); err != nil {
		return err
	}
	// TODO: Figure out a way to preserve the backend service name after this is
	// gone. Otherwise we can't find this internal load balancer again. However,
	// the backend service must be deleted first before health checks.
	if err := o.deleteBackendService(loadBalancerName); err != nil {
		return err
	}
	if err := o.deleteHealthCheck(loadBalancerName); err != nil {
		return err
	}
	if err := o.deleteHealthCheck(loadBalancerName + "-hc"); err != nil {
		return err
	}
	if err := o.deleteFirewall(loadBalancerName + "-hc"); err != nil {
		return err
	}
	return nil
}

// deleteExternalLoadBalancer follows a similar cleanup procedure as:
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_external.go#L289
// TODO: cleanup nodes health check (using clusterid)
func (o *ClusterUninstaller) deleteExternalLoadBalancer(loadBalancerName string) error {
	if err := o.deleteAddress(loadBalancerName); err != nil {
		return err
	}
	if err := o.deleteForwardingRule(loadBalancerName); err != nil {
		return err
	}
	if err := o.deleteFirewall(fmt.Sprintf("k8s-fw-%s", loadBalancerName)); err != nil {
		return err
	}
	if err := o.deleteTargetPool(loadBalancerName); err != nil {
		return err
	}
	if err := o.deleteHealthCheck(loadBalancerName); err != nil {
		return err
	}
	return nil
}

// destroyCloudControllerInternalLBs removes resources associated with internal load balancers
// created by the kube cloud controller. It first finds instance groups associated with instances
// belonging to this cluster, then finds backend resources that point to these instance groups.
// For each of those backend services, resources like forwarding rules, firewalls, health checks and
// backend services are removed.
func (o *ClusterUninstaller) destroyCloudControllerInternalLBs() error {
	groups, err := o.getInternalLBInstanceGroups()
	if err != nil {
		return err
	}
	if len(groups) == 0 {
		return nil
	}
	clusterID := strings.TrimPrefix(groups[0].name, "k8s-ig--")
	backends, err := o.listBackendServicesForInstanceGroups(groups)
	if err != nil {
		return err
	}

	// Each backend found represents an internal load balancer.
	// For each, remove related resources
	for _, backend := range backends {
		err := o.deleteInternalLoadBalancer(clusterID, backend)
		if err != nil {
			return err
		}
	}

	// Finally, remove the instance groups {
	for _, group := range groups {
		err := o.deleteInstanceGroup(group)
		if err != nil {
			return err
		}
	}
	return nil
}

// getExternalLBTargetPools returns all target pools that point to instances in the cluster
func (o *ClusterUninstaller) getExternalLBTargetPools() ([]string, error) {
	return o.listTargetPoolsWithFilter("items(name,instances)", "", func(pool *compute.TargetPool) bool {
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

// destroyCloudControllerExternalLBs removes resources associated with external load balancers
// created by the kube cloud controller. It first finds target pools associated with instances
// belonging to this cluster. For each of those target pools, it removes resources like
// addresses, forwarding rules, firewalls, health checks and target pools.
func (o *ClusterUninstaller) destroyCloudControllerExternalLBs() error {
	pools, err := o.getExternalLBTargetPools()
	if err != nil {
		return err
	}

	for _, loadBalancerName := range pools {
		err = o.deleteExternalLoadBalancer(loadBalancerName)
		if err != nil {
			return err
		}
	}

	return nil
}
