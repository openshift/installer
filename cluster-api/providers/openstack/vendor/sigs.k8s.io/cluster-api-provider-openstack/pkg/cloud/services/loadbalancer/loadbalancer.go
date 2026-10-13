/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package loadbalancer

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"net"
	"slices"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"k8s.io/apimachinery/pkg/util/wait"
	utilsnet "k8s.io/utils/net"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
	openstackutil "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/openstack"
	capostrings "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/strings"
)

const (
	networkPrefix           string = "k8s-clusterapi"
	kubeapiLBSuffix         string = "kubeapi"
	resolvedMsg             string = "ControlPlaneEndpoint.Host is not an IP address, using the first resolved IP address"
	waitForOctaviaLBCleanup        = 15 * time.Second
)

const (
	loadBalancerProvisioningStatusActive        = "ACTIVE"
	loadBalancerProvisioningStatusPendingDelete = "PENDING_DELETE"
)

// Default values for Monitor, sync with `kubebuilder:default` annotations on APIServerLoadBalancerMonitor object.
const (
	defaultMonitorDelay          = 10
	defaultMonitorTimeout        = 5
	defaultMonitorMaxRetries     = 5
	defaultMonitorMaxRetriesDown = 3
)

// We wrap the LookupHost function in a variable to allow overriding it in unit tests.
//
//nolint:gocritic
var lookupHost = func(host string) (*string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	ips, err := net.DefaultResolver.LookupHost(ctx, host)
	if err != nil {
		return nil, err
	}
	if ip := net.ParseIP(ips[0]); ip == nil {
		return nil, fmt.Errorf("failed to resolve IP address for host %s", host)
	}
	return &ips[0], nil
}

// ReconcileLoadBalancer reconciles the load balancer for the given cluster.
func (s *Service) ReconcileLoadBalancer(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string, apiServerPort int) (bool, error) {
	lbSpec := openStackCluster.Spec.APIServerLoadBalancer
	if !lbSpec.IsEnabled() {
		return false, nil
	}

	loadBalancerName := getLoadBalancerName(clusterResourceName)
	s.scope.Logger().Info("Reconciling load balancer", "name", loadBalancerName)

	lbStatus := openStackCluster.Status.APIServerLoadBalancer
	if lbStatus == nil {
		lbStatus = &infrav1.LoadBalancer{}
		openStackCluster.Status.APIServerLoadBalancer = lbStatus
	}

	lb, err := s.getOrCreateAPILoadBalancer(openStackCluster, clusterResourceName)
	if err != nil {
		if errors.Is(err, capoerrors.ErrFilterMatch) {
			return true, err
		}
		return false, err
	}

	lbStatus.Name = lb.Name
	lbStatus.ID = lb.ID
	lbStatus.InternalIP = lb.VipAddress
	lbStatus.Tags = lb.Tags

	if lb.ProvisioningStatus != loadBalancerProvisioningStatusActive {
		var err error
		lbID := lb.ID
		lb, err = s.waitForLoadBalancerActive(lbID)
		if err != nil {
			return false, fmt.Errorf("load balancer %q with id %s is not active after timeout: %v", loadBalancerName, lbID, err)
		}
	}

	if !ptr.Deref(openStackCluster.Spec.DisableAPIServerFloatingIP, false) {
		floatingIPAddress, err := getAPIServerFloatingIP(openStackCluster)
		if err != nil {
			return false, err
		}

		fp, err := s.networkingService.GetOrCreateFloatingIP(openStackCluster, openStackCluster, clusterResourceName, floatingIPAddress)
		if err != nil {
			if errors.Is(err, capoerrors.ErrFilterMatch) {
				return true, err
			}
			return false, err
		}

		// Write the floating IP to the status immediately so we won't
		// create a new floating IP on the next reconcile if something
		// fails below.
		lbStatus.IP = fp.FloatingIP

		if err = s.networkingService.AssociateFloatingIP(openStackCluster, fp, lb.VipPortID); err != nil {
			return false, err
		}
	}

	allowedCIDRsSupported, err := s.isAllowsCIDRSSupported(lb)
	if err != nil {
		return false, err
	}

	// AllowedCIDRs will be nil if allowed CIDRs is not supported by the Octavia provider
	if allowedCIDRsSupported {
		lbStatus.AllowedCIDRs = getCanonicalAllowedCIDRs(openStackCluster)
	} else {
		lbStatus.AllowedCIDRs = nil
	}

	portList := []int{apiServerPort}
	portList = append(portList, lbSpec.AdditionalPorts...)
	for _, port := range portList {
		if err := s.reconcileAPILoadBalancerListener(lb, openStackCluster, clusterResourceName, port); err != nil {
			return false, err
		}
	}

	return false, nil
}

// getAPIServerVIPAddress gets the VIP address for the API server from wherever it is specified.
// Returns an empty string if the VIP address is not specified and it should be allocated automatically.
func getAPIServerVIPAddress(openStackCluster *infrav1.OpenStackCluster) (*string, error) {
	switch {
	// We only use call this function when creating the loadbalancer, so this case should never be used
	case openStackCluster.Status.APIServerLoadBalancer != nil && openStackCluster.Status.APIServerLoadBalancer.InternalIP != "":
		return &openStackCluster.Status.APIServerLoadBalancer.InternalIP, nil

	// Explicit fixed IP in the cluster spec
	case openStackCluster.Spec.APIServerFixedIP != nil:
		return openStackCluster.Spec.APIServerFixedIP, nil

	// If we are using the VIP as the control plane endpoint, use any value explicitly set on the control plane endpoint
	case ptr.Deref(openStackCluster.Spec.DisableAPIServerFloatingIP, false) && openStackCluster.Spec.ControlPlaneEndpoint != nil && openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
		fixedIPAddress, err := lookupHost(openStackCluster.Spec.ControlPlaneEndpoint.Host)
		if err != nil {
			return nil, fmt.Errorf("lookup host: %w", err)
		}
		return fixedIPAddress, nil
	}

	return nil, nil
}

// getAPIServerFloatingIP gets the floating IP from wherever it is specified.
// Returns an empty string if the floating IP is not specified and it should be allocated automatically.
func getAPIServerFloatingIP(openStackCluster *infrav1.OpenStackCluster) (*string, error) {
	switch {
	// The floating IP was created previously
	case openStackCluster.Status.APIServerLoadBalancer != nil && openStackCluster.Status.APIServerLoadBalancer.IP != "":
		return &openStackCluster.Status.APIServerLoadBalancer.IP, nil

	// Explicit floating IP in the cluster spec
	case openStackCluster.Spec.APIServerFloatingIP != nil:
		return openStackCluster.Spec.APIServerFloatingIP, nil

	// An IP address is specified explicitly in the control plane endpoint
	case openStackCluster.Spec.ControlPlaneEndpoint != nil && openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
		floatingIPAddress, err := lookupHost(openStackCluster.Spec.ControlPlaneEndpoint.Host)
		if err != nil {
			return nil, fmt.Errorf("lookup host: %w", err)
		}
		return floatingIPAddress, nil
	}

	return nil, nil
}

// getCanonicalAllowedCIDRs gets a filtered list of CIDRs which should be allowed to access the API server loadbalancer.
// Invalid CIDRs are filtered from the list and emil a warning event.
// It returns a canonical representation that can be directly compared with other canonicalized lists.
func getCanonicalAllowedCIDRs(openStackCluster *infrav1.OpenStackCluster) []string {
	allowedCIDRs := []string{}

	if openStackCluster.Spec.APIServerLoadBalancer != nil && len(openStackCluster.Spec.APIServerLoadBalancer.AllowedCIDRs) > 0 {
		allowedCIDRs = append(allowedCIDRs, openStackCluster.Spec.APIServerLoadBalancer.AllowedCIDRs...)

		// In the first reconciliation loop, only the Ready field is set in openStackCluster.Status
		// All other fields are empty/nil
		if openStackCluster.Status.Bastion != nil {
			if openStackCluster.Status.Bastion.FloatingIP != "" {
				allowedCIDRs = append(allowedCIDRs, openStackCluster.Status.Bastion.FloatingIP)
			}

			if openStackCluster.Status.Bastion.IP != "" {
				allowedCIDRs = append(allowedCIDRs, openStackCluster.Status.Bastion.IP)
			}
		}

		if openStackCluster.Status.Network != nil {
			for _, subnet := range openStackCluster.Status.Network.Subnets {
				if subnet.CIDR != "" {
					allowedCIDRs = append(allowedCIDRs, subnet.CIDR)
				}
			}

			if openStackCluster.Status.Router != nil && len(openStackCluster.Status.Router.IPs) > 0 {
				allowedCIDRs = append(allowedCIDRs, openStackCluster.Status.Router.IPs...)
			}
		}
	}

	// Filter invalid CIDRs and convert any IPs into CIDRs.
	validCIDRs := []string{}
	for _, v := range allowedCIDRs {
		switch {
		case utilsnet.IsIPv4String(v):
			validCIDRs = append(validCIDRs, v+"/32")
		case utilsnet.IsIPv4CIDRString(v):
			validCIDRs = append(validCIDRs, v)
		default:
			record.Warnf(openStackCluster, "FailedIPAddressValidation", "%s is not a valid IPv4 nor CIDR address and will not get applied to allowed_cidrs", v)
		}
	}

	// Sort and remove duplicates
	return capostrings.Canonicalize(validCIDRs)
}

// isAllowsCIDRSSupported returns true if Octavia supports allowed CIDRs for the loadbalancer provider in use.
func (s *Service) isAllowsCIDRSSupported(lb *loadbalancers.LoadBalancer) (bool, error) {
	octaviaVersions, err := s.loadbalancerClient.ListOctaviaVersions()
	if err != nil {
		return false, err
	}
	// The current version is always the last one in the list.
	octaviaVersion := octaviaVersions[len(octaviaVersions)-1].ID

	return openstackutil.IsOctaviaFeatureSupported(octaviaVersion, openstackutil.OctaviaFeatureVIPACL, lb.Provider), nil
}

// getOrCreateAPILoadBalancer returns an existing API loadbalancer if it already exists, or creates a new one if it does not.
func (s *Service) getOrCreateAPILoadBalancer(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) (*loadbalancers.LoadBalancer, error) {
	loadBalancerName := getLoadBalancerName(clusterResourceName)
	lb, err := s.checkIfLbExists(loadBalancerName)
	if err != nil {
		return nil, err
	}
	if lb != nil {
		return lb, nil
	}

	if openStackCluster.Status.Network == nil {
		return nil, fmt.Errorf("network is not yet available in OpenStackCluster.Status")
	}

	if openStackCluster.Status.APIServerLoadBalancer == nil {
		return nil, fmt.Errorf("apiserver loadbalancer network is not yet available in OpenStackCluster.Status")
	}

	lbNetwork := openStackCluster.Status.APIServerLoadBalancer.LoadBalancerNetwork
	if lbNetwork == nil {
		lbNetwork = &infrav1.NetworkStatusWithSubnets{}
		openStackCluster.Status.APIServerLoadBalancer.LoadBalancerNetwork = lbNetwork
	}

	var vipNetworkID, vipSubnetID string
	if lbNetwork.ID != "" {
		vipNetworkID = lbNetwork.ID
	}

	if len(lbNetwork.Subnets) > 0 {
		// Currently only the first subnet is taken into account.
		// This can be fixed as soon as we switch over to gophercloud release that
		// contains AdditionalVips field.
		vipSubnetID = lbNetwork.Subnets[0].ID
	}

	if vipNetworkID == "" && vipSubnetID == "" {
		// keep the default and create the VIP on the first cluster subnet
		vipSubnetID = openStackCluster.Status.Network.Subnets[0].ID
		s.scope.Logger().Info("No load balancer network specified, creating load balancer in the default subnet", "subnetID", vipSubnetID, "name", loadBalancerName)
	} else {
		s.scope.Logger().Info("Creating load balancer in subnet", "subnetID", vipSubnetID, "name", loadBalancerName)
	}

	providers, err := s.loadbalancerClient.ListLoadBalancerProviders()
	if err != nil {
		return nil, err
	}

	// Choose the selected provider and flavor if set in cluster spec, if not, omit these fields and Octavia will use the default values.
	lbProvider := ""
	lbFlavorID := ""
	var availabilityZone *string
	if openStackCluster.Spec.APIServerLoadBalancer != nil {
		if openStackCluster.Spec.APIServerLoadBalancer.Provider != nil {
			for _, v := range providers {
				if v.Name == *openStackCluster.Spec.APIServerLoadBalancer.Provider {
					lbProvider = v.Name
					break
				}
			}
			if lbProvider == "" {
				record.Warnf(openStackCluster, "OctaviaProviderNotFound", "Provider specified for Octavia not found.")
				record.Eventf(openStackCluster, "OctaviaProviderNotFound", "Provider %s specified for Octavia not found, using the default provider.", openStackCluster.Spec.APIServerLoadBalancer.Provider)
			}
		}
		if openStackCluster.Spec.APIServerLoadBalancer.Flavor != nil {
			// Gophercloud does not support filtering loadbalancer flavors by name and status (enabled) so we have to get all available flavors
			// and filter them localy. There is a feature request in Gophercloud to implement this functionality:
			// https://github.com/gophercloud/gophercloud/v2/issues/3049
			flavors, err := s.loadbalancerClient.ListLoadBalancerFlavors()
			if err != nil {
				return nil, err
			}

			for _, v := range flavors {
				if v.Enabled && v.Name == *openStackCluster.Spec.APIServerLoadBalancer.Flavor {
					lbFlavorID = v.ID
					break
				}
			}
			if lbFlavorID == "" {
				record.Warnf(openStackCluster, "OctaviaFlavorNotFound", "Flavor %s specified for Octavia not found, using the default flavor.", *openStackCluster.Spec.APIServerLoadBalancer.Flavor)
			}
		}

		availabilityZone = openStackCluster.Spec.APIServerLoadBalancer.AvailabilityZone
	}

	vipAddress, err := getAPIServerVIPAddress(openStackCluster)
	if err != nil {
		return nil, err
	}

	lbCreateOpts := loadbalancers.CreateOpts{
		Name:         loadBalancerName,
		VipSubnetID:  vipSubnetID,
		VipNetworkID: vipNetworkID,
		Description:  names.GetDescription(clusterResourceName),
		Provider:     lbProvider,
		FlavorID:     lbFlavorID,
		Tags:         openStackCluster.Spec.Tags,
	}
	if availabilityZone != nil {
		lbCreateOpts.AvailabilityZone = *availabilityZone
	}
	if vipAddress != nil {
		lbCreateOpts.VipAddress = *vipAddress
	}

	lb, err = s.loadbalancerClient.CreateLoadBalancer(lbCreateOpts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateLoadBalancer", "Failed to create load balancer %s: %v", loadBalancerName, err)
		return nil, err
	}

	record.Eventf(openStackCluster, "SuccessfulCreateLoadBalancer", "Created load balancer %s with id %s", loadBalancerName, lb.ID)
	return lb, nil
}

// reconcileAPILoadBalancerListener ensures that the listener on the given port exists and is configured correctly.
func (s *Service) reconcileAPILoadBalancerListener(lb *loadbalancers.LoadBalancer, openStackCluster *infrav1.OpenStackCluster, clusterResourceName string, port int) error {
	loadBalancerName := getLoadBalancerName(clusterResourceName)
	lbPortObjectsName := fmt.Sprintf("%s-%d", loadBalancerName, port)

	if openStackCluster.Status.APIServerLoadBalancer == nil {
		return fmt.Errorf("APIServerLoadBalancer is not yet available in OpenStackCluster.Status")
	}

	allowedCIDRs := openStackCluster.Status.APIServerLoadBalancer.AllowedCIDRs

	listener, err := s.getOrCreateListener(openStackCluster, lbPortObjectsName, lb.ID, allowedCIDRs, port)
	if err != nil {
		return err
	}

	pool, err := s.getOrCreatePool(openStackCluster, lbPortObjectsName, listener.ID, lb.ID, lb.Provider)
	if err != nil {
		return err
	}
	if err := s.ensureMonitor(openStackCluster, lbPortObjectsName, pool.ID, lb.ID); err != nil {
		return err
	}

	// allowedCIDRs is nil if allowedCIDRs is not supported by the Octavia provider
	// A non-nil empty slice is an explicitly empty list
	if allowedCIDRs != nil {
		if err := s.getOrUpdateAllowedCIDRs(openStackCluster, listener, allowedCIDRs); err != nil {
			return err
		}
	}

	return nil
}

// getOrCreateListener returns an existing listener for the given loadbalancer
// and port if it already exists, or creates a new one if it does not.
func (s *Service) getOrCreateListener(openStackCluster *infrav1.OpenStackCluster, listenerName, lbID string, allowedCIDRs []string, port int) (*listeners.Listener, error) {
	listener, err := s.checkIfListenerExists(listenerName)
	if err != nil {
		return nil, err
	}

	if listener != nil {
		return listener, nil
	}

	s.scope.Logger().Info("Creating load balancer listener", "name", listenerName, "loadBalancerID", lbID)

	listenerCreateOpts := listeners.CreateOpts{
		Name:           listenerName,
		Protocol:       "TCP",
		ProtocolPort:   port,
		LoadbalancerID: lbID,
		Tags:           openStackCluster.Spec.Tags,
		AllowedCIDRs:   allowedCIDRs,
	}
	listener, err = s.loadbalancerClient.CreateListener(listenerCreateOpts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateListener", "Failed to create listener %s: %v", listenerName, err)
		return nil, err
	}

	if _, err := s.waitForLoadBalancerActive(lbID); err != nil {
		record.Warnf(openStackCluster, "FailedCreateListener", "Failed to create listener %s with id %s: wait for load balancer active %s: %v", listenerName, listener.ID, lbID, err)
		return nil, err
	}

	if err := s.waitForListener(listener.ID, "ACTIVE"); err != nil {
		record.Warnf(openStackCluster, "FailedCreateListener", "Failed to create listener %s with id %s: wait for listener active: %v", listenerName, listener.ID, err)
		return nil, err
	}

	record.Eventf(openStackCluster, "SuccessfulCreateListener", "Created listener %s with id %s", listenerName, listener.ID)
	return listener, nil
}

// getOrUpdateAllowedCIDRs ensures that the allowed CIDRs configured on a listener correspond to the expected list.
func (s *Service) getOrUpdateAllowedCIDRs(openStackCluster *infrav1.OpenStackCluster, listener *listeners.Listener, allowedCIDRs []string) error {
	// Sort and remove duplicates
	listener.AllowedCIDRs = capostrings.Canonicalize(listener.AllowedCIDRs)

	if !slices.Equal(allowedCIDRs, listener.AllowedCIDRs) {
		s.scope.Logger().Info("CIDRs do not match, updating listener", "expectedCIDRs", allowedCIDRs, "currentCIDRs", listener.AllowedCIDRs)
		listenerUpdateOpts := listeners.UpdateOpts{
			AllowedCIDRs: &allowedCIDRs,
		}

		listenerID := listener.ID
		listener, err := s.loadbalancerClient.UpdateListener(listener.ID, listenerUpdateOpts)
		if err != nil {
			record.Warnf(openStackCluster, "FailedUpdateListener", "Failed to update listener %s: %v", listenerID, err)
			return err
		}

		if err := s.waitForListener(listener.ID, "ACTIVE"); err != nil {
			record.Warnf(openStackCluster, "FailedUpdateListener", "Failed to update listener %s with id %s: wait for listener active: %v", listener.Name, listener.ID, err)
			return err
		}

		record.Eventf(openStackCluster, "SuccessfulUpdateListener", "Updated allowed_cidrs %s for listener %s with id %s", listener.AllowedCIDRs, listener.Name, listener.ID)
	}
	return nil
}

func (s *Service) getOrCreatePool(openStackCluster *infrav1.OpenStackCluster, poolName, listenerID, lbID string, lbProvider string) (*pools.Pool, error) {
	pool, err := s.checkIfPoolExists(poolName)
	if err != nil {
		return nil, err
	}

	if pool != nil {
		return pool, nil
	}

	s.scope.Logger().Info("Creating load balancer pool for listener", "loadBalancerID", lbID, "listenerID", listenerID, "name", poolName)

	method := pools.LBMethodRoundRobin

	if lbProvider == "ovn" {
		method = pools.LBMethodSourceIpPort
	}

	poolCreateOpts := pools.CreateOpts{
		Name:       poolName,
		Protocol:   "TCP",
		LBMethod:   method,
		ListenerID: listenerID,
		Tags:       openStackCluster.Spec.Tags,
	}
	pool, err = s.loadbalancerClient.CreatePool(poolCreateOpts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreatePool", "Failed to create pool %s: %v", poolName, err)
		return nil, err
	}

	if _, err := s.waitForLoadBalancerActive(lbID); err != nil {
		record.Warnf(openStackCluster, "FailedCreatePool", "Failed to create pool %s with id %s: wait for load balancer active %s: %v", poolName, pool.ID, lbID, err)
		return nil, err
	}

	record.Eventf(openStackCluster, "SuccessfulCreatePool", "Created pool %s with id %s", poolName, pool.ID)
	return pool, nil
}

func (s *Service) ensureMonitor(openStackCluster *infrav1.OpenStackCluster, monitorName, poolID, lbID string) error {
	var cfg infrav1.APIServerLoadBalancerMonitor

	if openStackCluster.Spec.APIServerLoadBalancer.Monitor != nil {
		cfg = *openStackCluster.Spec.APIServerLoadBalancer.Monitor
	}

	cfg.Delay = cmp.Or(cfg.Delay, defaultMonitorDelay)
	cfg.Timeout = cmp.Or(cfg.Timeout, defaultMonitorTimeout)
	cfg.MaxRetries = cmp.Or(cfg.MaxRetries, defaultMonitorMaxRetries)
	cfg.MaxRetriesDown = cmp.Or(cfg.MaxRetriesDown, defaultMonitorMaxRetriesDown)

	monitor, err := s.checkIfMonitorExists(monitorName)
	if err != nil {
		return err
	}

	if monitor != nil {
		needsUpdate := false
		monitorUpdateOpts := monitors.UpdateOpts{}

		if monitor.Delay != cfg.Delay {
			s.scope.Logger().Info("Monitor delay needs update", "current", monitor.Delay, "desired", cfg.Delay)
			monitorUpdateOpts.Delay = cfg.Delay
			needsUpdate = true
		}

		if monitor.Timeout != cfg.Timeout {
			s.scope.Logger().Info("Monitor timeout needs update", "current", monitor.Timeout, "desired", cfg.Timeout)
			monitorUpdateOpts.Timeout = cfg.Timeout
			needsUpdate = true
		}

		if monitor.MaxRetries != cfg.MaxRetries {
			s.scope.Logger().Info("Monitor maxRetries needs update", "current", monitor.MaxRetries, "desired", cfg.MaxRetries)
			monitorUpdateOpts.MaxRetries = cfg.MaxRetries
			needsUpdate = true
		}

		if monitor.MaxRetriesDown != cfg.MaxRetriesDown {
			s.scope.Logger().Info("Monitor maxRetriesDown needs update", "current", monitor.MaxRetriesDown, "desired", cfg.MaxRetriesDown)
			monitorUpdateOpts.MaxRetriesDown = cfg.MaxRetriesDown
			needsUpdate = true
		}

		if needsUpdate {
			s.scope.Logger().Info("Updating load balancer monitor", "loadBalancerID", lbID, "name", monitorName, "monitorID", monitor.ID)

			updatedMonitor, err := s.loadbalancerClient.UpdateMonitor(monitor.ID, monitorUpdateOpts)
			if err != nil {
				record.Warnf(openStackCluster, "FailedUpdateMonitor", "Failed to update monitor %s with id %s: %v", monitorName, monitor.ID, err)
				return err
			}

			if _, err = s.waitForLoadBalancerActive(lbID); err != nil {
				record.Warnf(openStackCluster, "FailedUpdateMonitor", "Failed to update monitor %s with id %s: wait for load balancer active %s: %v", monitorName, monitor.ID, lbID, err)
				return err
			}

			record.Eventf(openStackCluster, "SuccessfulUpdateMonitor", "Updated monitor %s with id %s", monitorName, updatedMonitor.ID)
		}

		return nil
	}

	s.scope.Logger().Info("Creating load balancer monitor for pool", "loadBalancerID", lbID, "name", monitorName, "poolID", poolID)

	monitor, err = s.loadbalancerClient.CreateMonitor(monitors.CreateOpts{
		Name:           monitorName,
		PoolID:         poolID,
		Type:           "TCP",
		Delay:          cfg.Delay,
		Timeout:        cfg.Timeout,
		MaxRetries:     cfg.MaxRetries,
		MaxRetriesDown: cfg.MaxRetriesDown,
	})
	if err != nil {
		if capoerrors.IsNotImplementedError(err) {
			record.Warnf(openStackCluster, "SkippedCreateMonitor", "Health Monitor is not created as it's not implemented with the current Octavia provider.")
			return nil
		}

		record.Warnf(openStackCluster, "FailedCreateMonitor", "Failed to create monitor %s: %v", monitorName, err)
		return err
	}

	if _, err = s.waitForLoadBalancerActive(lbID); err != nil {
		record.Warnf(openStackCluster, "FailedCreateMonitor", "Failed to create monitor %s with id %s: wait for load balancer active %s: %v", monitorName, monitor.ID, lbID, err)
		return err
	}

	record.Eventf(openStackCluster, "SuccessfulCreateMonitor", "Created monitor %s with id %s", monitorName, monitor.ID)
	return nil
}

func (s *Service) ReconcileLoadBalancerMember(openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, clusterResourceName, ip string) error {
	if openStackCluster.Status.Network == nil {
		return errors.New("network is not yet available in openStackCluster.Status")
	}
	if len(openStackCluster.Status.Network.Subnets) == 0 {
		return errors.New("network.Subnets are not yet available in openStackCluster.Status")
	}
	if openStackCluster.Status.APIServerLoadBalancer == nil {
		return errors.New("network.APIServerLoadBalancer is not yet available in openStackCluster.Status")
	}
	if openStackCluster.Spec.ControlPlaneEndpoint == nil || !openStackCluster.Spec.ControlPlaneEndpoint.IsValid() {
		return errors.New("ControlPlaneEndpoint is not yet set in openStackCluster.Spec")
	}

	loadBalancerName := getLoadBalancerName(clusterResourceName)
	s.scope.Logger().Info("Reconciling load balancer member", "loadBalancerName", loadBalancerName)

	lbID := openStackCluster.Status.APIServerLoadBalancer.ID
	var portList []int
	if openStackCluster.Spec.ControlPlaneEndpoint != nil {
		portList = append(portList, int(openStackCluster.Spec.ControlPlaneEndpoint.Port))
	}
	if openStackCluster.Spec.APIServerLoadBalancer != nil {
		portList = append(portList, openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts...)
	}
	for _, port := range portList {
		lbPortObjectsName := fmt.Sprintf("%s-%d", loadBalancerName, port)
		name := lbPortObjectsName + "-" + openStackMachine.Name

		pool, err := s.checkIfPoolExists(lbPortObjectsName)
		if err != nil {
			return err
		}
		if pool == nil {
			return errors.New("load balancer pool does not exist yet")
		}

		lbMember, err := s.checkIfLbMemberExists(pool.ID, name)
		if err != nil {
			return err
		}

		if lbMember != nil {
			// check if we have to recreate the LB Member
			if lbMember.Address == ip {
				// nothing to do continue to next port
				continue
			}

			s.scope.Logger().Info("Deleting load balancer member because the IP of the machine changed", "name", name)

			// lb member changed so let's delete it so we can create it again with the correct IP
			_, err = s.waitForLoadBalancerActive(lbID)
			if err != nil {
				return err
			}
			if err := s.loadbalancerClient.DeletePoolMember(pool.ID, lbMember.ID); err != nil {
				return err
			}
			_, err = s.waitForLoadBalancerActive(lbID)
			if err != nil {
				return err
			}
		}

		s.scope.Logger().Info("Creating load balancer member", "name", name)

		// if we got to this point we should either create or re-create the lb member
		lbMemberOpts := pools.CreateMemberOpts{
			Name:         name,
			ProtocolPort: port,
			Address:      ip,
			Tags:         openStackCluster.Spec.Tags,
		}

		if _, err := s.waitForLoadBalancerActive(lbID); err != nil {
			return err
		}

		if _, err := s.loadbalancerClient.CreatePoolMember(pool.ID, lbMemberOpts); err != nil {
			return err
		}

		if _, err := s.waitForLoadBalancerActive(lbID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) DeleteLoadBalancer(openStackCluster *infrav1.OpenStackCluster, clusterResourceName string) (result *ctrl.Result, reterr error) {
	loadBalancerName := getLoadBalancerName(clusterResourceName)
	lb, err := s.checkIfLbExists(loadBalancerName)
	if err != nil {
		return nil, err
	}

	if lb == nil {
		return nil, nil
	}

	// If the load balancer is already in PENDING_DELETE state, we don't need to do anything.
	// However we should still wait for the load balancer to be deleted which is why we
	// request a new reconcile after a certain amount of time.
	if lb.ProvisioningStatus == loadBalancerProvisioningStatusPendingDelete {
		s.scope.Logger().Info("Load balancer is already in PENDING_DELETE state", "name", loadBalancerName)
		return &ctrl.Result{RequeueAfter: waitForOctaviaLBCleanup}, nil
	}

	if lb.VipPortID != "" {
		fip, err := s.networkingService.GetFloatingIPByPortID(lb.VipPortID)
		if err != nil {
			return nil, err
		}

		if fip != nil && fip.FloatingIP != "" {
			if err = s.networkingService.DisassociateFloatingIP(openStackCluster, fip.FloatingIP); err != nil {
				return nil, err
			}

			// If the floating is user-provider (BYO floating IP), don't delete it.
			if openStackCluster.Spec.APIServerFloatingIP == nil || *openStackCluster.Spec.APIServerFloatingIP != fip.FloatingIP {
				if err = s.networkingService.DeleteFloatingIP(openStackCluster, fip.FloatingIP); err != nil {
					return nil, err
				}
			} else {
				s.scope.Logger().Info("Skipping load balancer floating IP deletion as it's a user-provided resource", "name", loadBalancerName, "fip", fip.FloatingIP)
			}
		}
	}

	deleteOpts := loadbalancers.DeleteOpts{
		Cascade: true,
	}
	s.scope.Logger().Info("Deleting load balancer", "name", loadBalancerName, "cascade", deleteOpts.Cascade)
	err = s.loadbalancerClient.DeleteLoadBalancer(lb.ID, deleteOpts)
	if err != nil && !capoerrors.IsNotFound(err) {
		record.Warnf(openStackCluster, "FailedDeleteLoadBalancer", "Failed to delete load balancer %s with id %s: %v", lb.Name, lb.ID, err)
		return nil, err
	}

	record.Eventf(openStackCluster, "SuccessfulDeleteLoadBalancer", "Deleted load balancer %s with id %s", lb.Name, lb.ID)

	// If we have reached this point, that means that the load balancer wasn't initially deleted but the request to delete it didn't return an error.
	// So we want to requeue until the load balancer and its associated ports are actually deleted.
	return &ctrl.Result{RequeueAfter: waitForOctaviaLBCleanup}, nil
}

func (s *Service) DeleteLoadBalancerMember(openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, clusterResourceName string) error {
	if openStackMachine == nil {
		return errors.New("openStackMachine is nil")
	}

	loadBalancerName := getLoadBalancerName(clusterResourceName)
	lb, err := s.checkIfLbExists(loadBalancerName)
	if err != nil {
		return err
	}
	if lb == nil {
		// nothing to do
		return nil
	}

	lbID := lb.ID

	var portList []int
	if openStackCluster.Spec.ControlPlaneEndpoint != nil {
		portList = append(portList, int(openStackCluster.Spec.ControlPlaneEndpoint.Port))
	}
	if openStackCluster.Spec.APIServerLoadBalancer != nil {
		portList = append(portList, openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts...)
	}
	for _, port := range portList {
		lbPortObjectsName := fmt.Sprintf("%s-%d", loadBalancerName, port)
		name := lbPortObjectsName + "-" + openStackMachine.Name

		pool, err := s.checkIfPoolExists(lbPortObjectsName)
		if err != nil {
			return err
		}
		if pool == nil {
			s.scope.Logger().Info("Load balancer pool does not exist", "name", lbPortObjectsName)
			continue
		}

		lbMember, err := s.checkIfLbMemberExists(pool.ID, name)
		if err != nil {
			return err
		}

		if lbMember != nil {
			// lb member changed so let's delete it so we can create it again with the correct IP
			_, err = s.waitForLoadBalancerActive(lbID)
			if err != nil {
				return err
			}
			if err := s.loadbalancerClient.DeletePoolMember(pool.ID, lbMember.ID); err != nil {
				return err
			}
			_, err = s.waitForLoadBalancerActive(lbID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getLoadBalancerName(clusterResourceName string) string {
	return fmt.Sprintf("%s-cluster-%s-%s", networkPrefix, clusterResourceName, kubeapiLBSuffix)
}

func (s *Service) checkIfLbExists(name string) (*loadbalancers.LoadBalancer, error) {
	lbList, err := s.loadbalancerClient.ListLoadBalancers(loadbalancers.ListOpts{Name: name})
	if err != nil {
		return nil, err
	}
	if len(lbList) == 0 {
		return nil, nil
	}
	return &lbList[0], nil
}

func (s *Service) checkIfListenerExists(name string) (*listeners.Listener, error) {
	listenerList, err := s.loadbalancerClient.ListListeners(listeners.ListOpts{Name: name})
	if err != nil {
		return nil, err
	}
	if len(listenerList) == 0 {
		return nil, nil
	}
	return &listenerList[0], nil
}

func (s *Service) checkIfPoolExists(name string) (*pools.Pool, error) {
	poolList, err := s.loadbalancerClient.ListPools(pools.ListOpts{Name: name})
	if err != nil {
		return nil, err
	}
	if len(poolList) == 0 {
		return nil, nil
	}
	return &poolList[0], nil
}

func (s *Service) checkIfMonitorExists(name string) (*monitors.Monitor, error) {
	monitorList, err := s.loadbalancerClient.ListMonitors(monitors.ListOpts{Name: name})
	if err != nil {
		return nil, err
	}
	if len(monitorList) == 0 {
		return nil, nil
	}
	return &monitorList[0], nil
}

func (s *Service) checkIfLbMemberExists(poolID, name string) (*pools.Member, error) {
	lbMemberList, err := s.loadbalancerClient.ListPoolMember(poolID, pools.ListMembersOpts{Name: name})
	if err != nil {
		return nil, err
	}
	if len(lbMemberList) == 0 {
		return nil, nil
	}
	return &lbMemberList[0], nil
}

var backoff = wait.Backoff{
	Steps:    20,
	Duration: time.Second,
	Factor:   1.25,
	Jitter:   0.1,
}

// Possible LoadBalancer states are documented here: https://docs.openstack.org/api-ref/load-balancer/v2/index.html#prov-status
func (s *Service) waitForLoadBalancerActive(id string) (*loadbalancers.LoadBalancer, error) {
	var lb *loadbalancers.LoadBalancer

	s.scope.Logger().Info("Waiting for load balancer", "id", id, "targetStatus", "ACTIVE")
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		var err error
		lb, err = s.loadbalancerClient.GetLoadBalancer(id)
		if err != nil {
			return false, err
		}
		return lb.ProvisioningStatus == loadBalancerProvisioningStatusActive, nil
	})
	if err != nil {
		return nil, err
	}
	return lb, nil
}

func (s *Service) waitForListener(id, target string) error {
	s.scope.Logger().Info("Waiting for load balancer listener", "id", id, "targetStatus", target)
	return wait.ExponentialBackoff(backoff, func() (bool, error) {
		_, err := s.loadbalancerClient.GetListener(id)
		if err != nil {
			return false, err
		}
		// The listener resource has no Status attribute, so a successful Get is the best we can do
		return true, nil
	})
}
