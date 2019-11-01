package gcp

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	compute "google.golang.org/api/compute/v1"
)

// getInternalLBInstanceGroups finds instance groups created for kube cloud controller
// internal load balancers. They should be named "k8s-ig--{clusterid}":
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_naming.go#L33-L40
// where clusterid is an 8-char id generated and saved in a configmap named
// "ingress-uid" in kube-system, key "uid":
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_clusterid.go#L210-L238
// If no clusterID is set, we look into each instance group and determine if it contains instances prefixed with the cluster's infra ID
func (o *ClusterUninstaller) getInternalLBInstanceGroups() ([]cloudResource, error) {
	filter := "name eq \"k8s-ig--.*\""
	if len(o.cloudControllerUID) > 0 {
		filter = fmt.Sprintf("name eq \"k8s-ig--%s\"", o.cloudControllerUID)
	}
	candidates, err := o.listInstanceGroupsWithFilter(filter)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, nil
	}

	igName := ""
	if len(o.cloudControllerUID) > 0 {
		igName = fmt.Sprintf("k8s-ig--%s", o.cloudControllerUID)
	} else {
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
	}
	igs := []cloudResource{}
	if len(igName) > 0 {
		for _, ig := range candidates {
			if ig.name == igName {
				igs = append(igs, ig)
			}
		}
	}
	return igs, nil
}

func (o *ClusterUninstaller) listBackendServicesForInstanceGroups(igs []cloudResource) ([]cloudResource, error) {
	urls := sets.NewString()
	for _, ig := range igs {
		urls.Insert(o.getInstanceGroupURL(ig))
	}
	return o.listBackendServicesWithFilter("items(name,backends),nextPageToken", "name eq \"a[0-9a-f]{30,50}\"", func(item *compute.BackendService) bool {
		if len(item.Backends) == 0 {
			return false
		}
		for _, backend := range item.Backends {
			if !urls.Has(backend.Group) {
				return false
			}
		}
		o.Logger.Debugf("Found backend service %s", item.Name)
		return true
	})
}

// deleteInternalLoadBalancer follows a similar cleanup procedure as:
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_internal.go#L222
// TODO: add cleanup for shared mode resources (determine if it's supported in 4.2)
func (o *ClusterUninstaller) deleteInternalLoadBalancer(loadBalancerName string) error {
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
		key:      loadBalancerName,
		name:     loadBalancerName,
		typeName: "healthcheck",
	}
	if err := o.deleteHealthCheck(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName+"-hc",
		name:     loadBalancerName+"-hc",
		typeName: "healthcheck",
	}
	if err := o.deleteHealthCheck(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName+"-hc",
		name:     loadBalancerName+"-hc",
		typeName: "healthcheck",
	}
	if err := o.deleteHealthCheck(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName,
		name:     loadBalancerName,
		typeName: "firewall",
	}
	if err := o.deleteFirewall(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName+"-hc",
		name:     loadBalancerName+"-hc",
		typeName: "firewall",
	}
	if err := o.deleteFirewall(item, true); err != nil {
		return err
	}
	item = cloudResource{
		key:      loadBalancerName,
		name:     loadBalancerName,
		typeName: "backendservice",
	}
	if err := o.deleteBackendService(item); err != nil {
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
		o.Logger.Debugf("Did not find any internal load balancer instance groups")
		return nil
	}
	if o.cloudControllerUID == "" {
		o.cloudControllerUID = strings.TrimPrefix(groups[0].name, "k8s-ig--")
	}
	backends, err := o.listBackendServicesForInstanceGroups(groups)
	if err != nil {
		return err
	}

	errs := []error{}
	found := cloudResources{}
	// Each backend found represents an internal load balancer.
	// For each, remove related resources
	for _, backend := range backends {
		found.insert(backend)
		err := o.deleteInternalLoadBalancer(backend.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("internallb", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted internal load balancer %s", item.name)
	}
	if len(errs) > 0 {
		return aggregateError(errs)
	}

	// Finally, remove the instance groups {
	found = cloudResources{}
	for _, group := range groups {
		found.insert(group)
		err := o.deleteInstanceGroup(group)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted = o.setPendingItems("lbinstancegroup", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted instance group %s", item.name)
	}
	if len(o.cloudControllerUID) > 0 {
		item := cloudResource{
			key:      fmt.Sprintf("k8s-%s-node", o.cloudControllerUID),
			name:     fmt.Sprintf("k8s-%s-node", o.cloudControllerUID),
			typeName: "healthcheck",
		}
		if err := o.deleteHealthCheck(item, true); err != nil {
			return err
		}
		item = cloudResource{
			key:      fmt.Sprintf("k8s-%s-node-http-hc", o.cloudControllerUID),
			name:     fmt.Sprintf("k8s-%s-node-http-hc", o.cloudControllerUID),
			typeName: "firewall",
		}
		if err := o.deleteFirewall(item, true); err != nil {
			return err
		}
	}
	return aggregateError(errs, len(found))
}
