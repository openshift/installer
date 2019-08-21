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
	req := o.computeSvc.Firewalls.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
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

func (o *ClusterUninstaller) deleteFirewall(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting firewall rule %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Firewalls.Delete(o.ProjectID, name).RequestId(o.requestID("firewall", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("firewall", name)
		return errors.Wrapf(err, "failed to delete firewall %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("firewall", name)
		return errors.Errorf("failed to delete firewall %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of firewall %s is pending", name)
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
	found := make([]string, 0, len(firewalls))
	errs := []error{}
	for _, firewall := range firewalls {
		found = append(found, firewall)
		err := o.deleteFirewall(firewall, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("firewall", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted firewall %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listAddresses() ([]string, error) {
	o.Logger.Debugf("Listing addresses")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Addresses.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
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

func (o *ClusterUninstaller) deleteAddress(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting address %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.Addresses.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("address", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("address", name)
		return errors.Wrapf(err, "failed to delete address %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("address", name)
		return errors.Errorf("failed to delete address %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of address %s is pending", name)
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
	found := make([]string, 0, len(addresses))
	errs := []error{}
	for _, address := range addresses {
		found = append(found, address)
		err := o.deleteAddress(address, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("address", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted address %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listForwardingRules() ([]string, error) {
	o.Logger.Debugf("Listing forwarding rules")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.ForwardingRules.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
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

func (o *ClusterUninstaller) deleteForwardingRule(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting forwarding rule %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.ForwardingRules.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("forwardingrule", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("forwardingrule", name)
		return errors.Wrapf(err, "failed to delete forwarding rule %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("forwardingrule", name)
		return errors.Errorf("failed to delete forwarding rule %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of forwarding rule %s is pending", name)
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
	found := make([]string, 0, len(forwardingRules))
	errs := []error{}
	for _, forwardingRule := range forwardingRules {
		found = append(found, forwardingRule)
		err := o.deleteForwardingRule(forwardingRule, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("forwardingrule", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted forwarding rule %s", item)
	}
	return aggregateError(errs, len(found))
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
	req := o.computeSvc.RegionBackendServices.List(o.ProjectID, o.Region).Fields(googleapi.Field(fields))
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
	op, err := o.computeSvc.RegionBackendServices.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("backendservice", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("backendservice", name)
		return errors.Wrapf(err, "failed to delete backend service %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("backendservice", name)
		return errors.Errorf("failed to delete backend service %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(backendServices))
	errs := []error{}
	for _, backendService := range backendServices {
		found = append(found, backendService)
		err := o.deleteBackendService(backendService)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("backendservice", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted backend service %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listHealthChecks() ([]string, error) {
	o.Logger.Debugf("Listing health checks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.HealthChecks.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
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

func (o *ClusterUninstaller) deleteHealthCheck(name string, errorOnPending bool) error {
	o.Logger.Debugf("Deleting health check %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.HealthChecks.Delete(o.ProjectID, name).RequestId(o.requestID("healthcheck", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("healthcheck", name)
		return errors.Wrapf(err, "failed to delete health check %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("healthcheck", name)
		return errors.Errorf("failed to delete health check %s with error: %s", name, operationErrorMessage(op))
	}
	if errorOnPending && op != nil && op.Status != "DONE" {
		return errors.Errorf("deletion of forwarding rule %s is pending", name)
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
	found := make([]string, 0, len(healthChecks))
	errs := []error{}
	for _, healthCheck := range healthChecks {
		found = append(found, healthCheck)
		err := o.deleteHealthCheck(healthCheck, false)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("healthcheck", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted health check %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listHTTPHealthChecks() ([]string, error) {
	o.Logger.Debugf("Listing HTTP health checks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.HttpHealthChecks.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
	err := req.Pages(ctx, func(list *compute.HttpHealthCheckList) error {
		for _, healthCheck := range list.Items {
			o.Logger.Debugf("Found HTTP health check: %s", healthCheck.Name)
			result = append(result, healthCheck.Name)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list HTTP health checks")
	}
	return result, nil
}

func (o *ClusterUninstaller) deleteHTTPHealthCheck(name string) error {
	o.Logger.Debugf("Deleting HTTP health check %s", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	op, err := o.computeSvc.HttpHealthChecks.Delete(o.ProjectID, name).RequestId(o.requestID("httphealthcheck", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("httphealthcheck", name)
		return errors.Wrapf(err, "failed to delete HTTP health check %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("httphealthcheck", name)
		return errors.Errorf("failed to delete HTTP health check %s with error: %s", name, operationErrorMessage(op))
	}
	return nil
}

// destroyHTTPHealthChecks removes all HTTP health check resources that have a name prefixed
// with the cluster's infra ID
func (o *ClusterUninstaller) destroyHTTPHealthChecks() error {
	healthChecks, err := o.listHTTPHealthChecks()
	if err != nil {
		return err
	}
	found := make([]string, 0, len(healthChecks))
	errs := []error{}
	for _, healthCheck := range healthChecks {
		found = append(found, healthCheck)
		err := o.deleteHTTPHealthCheck(healthCheck)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("httphealthcheck", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted HTTP health check %s", item)
	}
	return aggregateError(errs, len(found))
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
	op, err := o.computeSvc.TargetPools.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("targetpool", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("targetpool", name)
		return errors.Wrapf(err, "failed to delete target pool %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("targetpool", name)
		return errors.Errorf("failed to delete route %s with error: %s", name, operationErrorMessage(op))
	}
	o.Logger.Infof("Deleted target pool %s", name)
	return nil
}

func (o *ClusterUninstaller) clearTargetPoolHealthChecks(name string) error {
	o.Logger.Debugf("Clearing target pool %s health checks", name)
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	targetPool, err := o.computeSvc.TargetPools.Get(o.ProjectID, o.Region, name).Context(ctx).Do()
	if isNotFound(err) {
		return nil
	}
	if err != nil {
		return errors.Wrapf(err, "cannot retrieve target pool %s", name)
	}
	if len(targetPool.HealthChecks) == 0 {
		o.Logger.Debugf("Target pool %s has no health checks to clear", name)
		return nil
	}
	hcRemoveRequest := &compute.TargetPoolsRemoveHealthCheckRequest{}
	for _, hc := range targetPool.HealthChecks {
		hcRemoveRequest.HealthChecks = append(hcRemoveRequest.HealthChecks, &compute.HealthCheckReference{
			HealthCheck: hc,
		})
	}
	op, err := o.computeSvc.TargetPools.RemoveHealthCheck(o.ProjectID, o.Region, name, hcRemoveRequest).Context(ctx).RequestId(o.requestID("cleartargetpool", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("cleartargetpool", name)
		return errors.Wrapf(err, "failed to clear target pool %s health checks", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("cleartargetpool", name)
		return errors.Errorf("failed to clear target pool %s health checks with error: %s", name, operationErrorMessage(op))
	}
	if op != nil && op.Status != "DONE" {
		return errors.Errorf("target pool pending to be cleared of health checks")
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
	found := make([]string, 0, len(targetPools))
	errs := []error{}
	for _, targetPool := range targetPools {
		found = append(found, targetPool)
		err := o.deleteTargetPool(targetPool)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("targetpool", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted target pool %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listSubNetworks() ([]string, error) {
	o.Logger.Debugf("Listing subnetworks")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Subnetworks.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
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
	op, err := o.computeSvc.Subnetworks.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("subnetwork", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("subnetwork", name)
		return errors.Wrapf(err, "failed to delete subnetwork %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("subnetwork", name)
		return errors.Errorf("failed to delete subnetwork %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(subNetworks))
	errs := []error{}
	for _, subNetwork := range subNetworks {
		found = append(found, subNetwork)
		err := o.deleteSubNetwork(subNetwork)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("subnetwork", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted subnetwork %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listNetworks() ([]nameAndURL, error) {
	o.Logger.Debugf("Listing networks")
	result := []nameAndURL{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Networks.List(o.ProjectID).Fields("items(name,selfLink),nextPageToken").Filter(o.clusterIDFilter())
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
	ctx, cancel := o.contextWithTimeout()
	defer cancel()

	o.Logger.Debugf("Deleting network %s", name)
	op, err := o.computeSvc.Networks.Delete(o.ProjectID, name).RequestId(o.requestID("network", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("network", name)
		return errors.Wrapf(err, "failed to delete network %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("network", name)
		return errors.Errorf("failed to delete network %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(networks))
	errs := []error{}
	for _, network := range networks {
		found = append(found, network.name)
		// destroy any network routes that are not named with the infra ID
		routes, err := o.listNetworkRoutes(network.url)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, route := range routes {
			err := o.deleteRoute(route)
			if err != nil {
				o.Logger.Debugf("Failed to delete route %s: %v", route, err)
			}
		}

		err = o.deleteNetwork(network.name)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("network", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted network %s", item)
	}
	return aggregateError(errs, len(found))
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
	req := o.computeSvc.Routes.List(o.ProjectID).Fields("items(name),nextPageToken").Filter(filter)
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
	op, err := o.computeSvc.Routes.Delete(o.ProjectID, name).RequestId(o.requestID("route", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("route", name)
		return errors.Wrapf(err, "failed to delete route %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("route", name)
		return errors.Errorf("failed to delete route %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(routes))
	errs := []error{}
	for _, route := range routes {
		found = append(found, route)
		err := o.deleteRoute(route)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("route", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted route %s", item)
	}
	return aggregateError(errs, len(found))
}

func (o *ClusterUninstaller) listRouters() ([]string, error) {
	o.Logger.Debug("Listing routers")
	result := []string{}
	ctx, cancel := o.contextWithTimeout()
	defer cancel()
	req := o.computeSvc.Routers.List(o.ProjectID, o.Region).Fields("items(name),nextPageToken").Filter(o.clusterIDFilter())
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
	op, err := o.computeSvc.Routers.Delete(o.ProjectID, o.Region, name).RequestId(o.requestID("router", name)).Context(ctx).Do()
	if err != nil && !isNoOp(err) {
		o.resetRequestID("router", name)
		return errors.Wrapf(err, "failed to delete router %s", name)
	}
	if op != nil && op.Status == "DONE" && isErrorStatus(op.HttpErrorStatusCode) {
		o.resetRequestID("router", name)
		return errors.Errorf("failed to delete router %s with error: %s", name, operationErrorMessage(op))
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
	found := make([]string, 0, len(routers))
	errs := []error{}
	for _, router := range routers {
		found = append(found, router)
		err := o.deleteRouter(router)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("router", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted router %s", item)
	}
	return aggregateError(errs, len(found))
}

// getInternalLBInstanceGroups finds instance groups created for kube cloud controller
// internal load balancers. They should be named "k8s-ig--{clusterid}":
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_loadbalancer_naming.go#L33-L40
// where clusterid is an 8-char id generated and saved in a configmap named
// "ingress-uid" in kube-system, key "uid":
// https://github.com/openshift/kubernetes/blob/1e5983903742f64bca36a464582178c940353e9a/pkg/cloudprovider/providers/gce/gce_clusterid.go#L210-L238
// If no clusterID is set, we look into each instance group and determine if it contains instances prefixed with the cluster's infra ID
func (o *ClusterUninstaller) getInternalLBInstanceGroups() ([]nameAndZone, error) {
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
	return o.listBackendServicesWithFilter("items(name,backends)", "name eq \"a[0-9a-f]{30,50}\"", func(item *compute.BackendService) bool {
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
	if err := o.deleteAddress(loadBalancerName, true); err != nil {
		return err
	}
	if err := o.deleteForwardingRule(loadBalancerName, true); err != nil {
		return err
	}
	if err := o.deleteHealthCheck(loadBalancerName, true); err != nil {
		return err
	}
	if err := o.deleteHealthCheck(loadBalancerName+"-hc", true); err != nil {
		return err
	}
	if err := o.deleteFirewall(loadBalancerName, true); err != nil {
		return err
	}
	if err := o.deleteFirewall(loadBalancerName+"-hc", true); err != nil {
		return err
	}
	if err := o.deleteBackendService(loadBalancerName); err != nil {
		return err
	}
	return nil
}

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
	if err := o.deleteAddress(loadBalancerName, true); err != nil {
		return err
	}
	if err := o.deleteForwardingRule(loadBalancerName, true); err != nil {
		return err
	}
	if err := o.deleteFirewall(fmt.Sprintf("k8s-fw-%s", loadBalancerName), true); err != nil {
		return err
	}
	if err := o.deleteFirewall(fmt.Sprintf("k8s-%s-http-hc", loadBalancerName), true); err != nil {
		return err
	}
	if err := o.deleteHealthCheck(loadBalancerName, true); err != nil {
		return err
	}
	if err := o.deleteTargetPool(loadBalancerName); err != nil {
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
	found := make([]string, 0, len(backends))
	// Each backend found represents an internal load balancer.
	// For each, remove related resources
	for _, backend := range backends {
		found = append(found, backend)
		err := o.deleteInternalLoadBalancer(backend)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("internallb", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted internal load balancer %s", item)
	}
	if len(errs) > 0 {
		return aggregateError(errs)
	}

	// Finally, remove the instance groups {
	found = make([]string, len(groups))
	for _, group := range groups {
		found = append(found, fmt.Sprintf("%s/%s", group.zone, group.name))
		err := o.deleteInstanceGroup(group)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted = o.setPendingItems("lbinstancegroup", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted instance group %s", item)
	}
	if len(o.cloudControllerUID) > 0 {
		if err := o.deleteHealthCheck(fmt.Sprintf("k8s-%s-node", o.cloudControllerUID), true); err != nil {
			return err
		}
		if err := o.deleteFirewall(fmt.Sprintf("k8s-%s-node-http-hc", o.cloudControllerUID), true); err != nil {
			return err
		}
	}
	return aggregateError(errs, len(found))
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

	found := make([]string, 0, len(pools))
	errs := []error{}
	for _, loadBalancerName := range pools {
		found = append(found, loadBalancerName)
		err := o.deleteExternalLoadBalancer(loadBalancerName)
		if err != nil {
			errs = append(errs, err)
		}
	}
	deleted := o.setPendingItems("externallb", found)
	for _, item := range deleted {
		o.Logger.Infof("Deleted external load balancer %s", item)
	}
	if len(o.cloudControllerUID) != 0 {
		if err := o.deleteHealthCheck(fmt.Sprintf("k8s-%s-node-hc", o.cloudControllerUID), true); err != nil {
			return err
		}
		if err := o.deleteFirewall(fmt.Sprintf("k8s-%s-node-hc", o.cloudControllerUID), true); err != nil {
			return err
		}
	}
	return aggregateError(errs, len(found))
}
