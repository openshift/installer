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
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/net"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/record"
	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/names"
	openstackutil "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/openstack"
	capostrings "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/strings"
)

const (
	networkPrefix   string = "k8s-clusterapi"
	kubeapiLBSuffix string = "kubeapi"
)

const loadBalancerProvisioningStatusActive = "ACTIVE"

func (s *Service) ReconcileLoadBalancer(openStackCluster *infrav1.OpenStackCluster, clusterName string, apiServerPort int) (bool, error) {
	loadBalancerName := getLoadBalancerName(clusterName)
	s.scope.Logger().Info("Reconciling load balancer", "name", loadBalancerName)

	var fixedIPAddress string
	switch {
	case openStackCluster.Spec.APIServerFixedIP != "":
		fixedIPAddress = openStackCluster.Spec.APIServerFixedIP
	case openStackCluster.Spec.DisableAPIServerFloatingIP && openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
		fixedIPAddress = openStackCluster.Spec.ControlPlaneEndpoint.Host
	}

	providers, err := s.loadbalancerClient.ListLoadBalancerProviders()
	if err != nil {
		return false, err
	}

	// Choose the selected provider if it is set in cluster spec, if not, omit the field and Octavia will use the default provider.
	lbProvider := ""
	if openStackCluster.Spec.APIServerLoadBalancer.Provider != "" {
		for _, v := range providers {
			if v.Name == openStackCluster.Spec.APIServerLoadBalancer.Provider {
				lbProvider = v.Name
				break
			}
		}
		if lbProvider == "" {
			record.Warnf(openStackCluster, "OctaviaProviderNotFound", "Provider specified for Octavia not found.")
			record.Eventf(openStackCluster, "OctaviaProviderNotFound", "Provider %s specified for Octavia not found, using the default provider.", openStackCluster.Spec.APIServerLoadBalancer.Provider)
		}
	}

	lb, err := s.getOrCreateLoadBalancer(openStackCluster, loadBalancerName, openStackCluster.Status.Network.Subnets[0].ID, clusterName, fixedIPAddress, lbProvider)
	if err != nil {
		return false, err
	}
	if err := s.waitForLoadBalancerActive(lb.ID); err != nil {
		return false, fmt.Errorf("load balancer %q with id %s is not active after timeout: %v", loadBalancerName, lb.ID, err)
	}

	var lbFloatingIP string
	if !openStackCluster.Spec.DisableAPIServerFloatingIP {
		var floatingIPAddress string
		switch {
		case openStackCluster.Spec.APIServerFloatingIP != "":
			floatingIPAddress = openStackCluster.Spec.APIServerFloatingIP
		case openStackCluster.Spec.ControlPlaneEndpoint.IsValid():
			floatingIPAddress = openStackCluster.Spec.ControlPlaneEndpoint.Host
		}
		fp, err := s.networkingService.GetOrCreateFloatingIP(openStackCluster, openStackCluster, clusterName, floatingIPAddress)
		if err != nil {
			return false, err
		}
		if err = s.networkingService.AssociateFloatingIP(openStackCluster, fp, lb.VipPortID); err != nil {
			return false, err
		}
		lbFloatingIP = fp.FloatingIP
	}

	allowedCIDRs := []string{}
	// To reduce API calls towards OpenStack API, let's handle the CIDR support verification for all Ports only once.
	allowedCIDRsSupported := false
	octaviaVersions, err := s.loadbalancerClient.ListOctaviaVersions()
	if err != nil {
		return false, err
	}
	// The current version is always the last one in the list.
	octaviaVersion := octaviaVersions[len(octaviaVersions)-1].ID
	if openstackutil.IsOctaviaFeatureSupported(octaviaVersion, openstackutil.OctaviaFeatureVIPACL, lbProvider) {
		allowedCIDRsSupported = true
	}

	portList := []int{apiServerPort}
	portList = append(portList, openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts...)
	for _, port := range portList {
		lbPortObjectsName := fmt.Sprintf("%s-%d", loadBalancerName, port)

		listener, err := s.getOrCreateListener(openStackCluster, lbPortObjectsName, lb.ID, port)
		if err != nil {
			return false, err
		}

		pool, err := s.getOrCreatePool(openStackCluster, lbPortObjectsName, listener.ID, lb.ID, lb.Provider)
		if err != nil {
			return false, err
		}

		if err := s.getOrCreateMonitor(openStackCluster, lbPortObjectsName, pool.ID, lb.ID); err != nil {
			return false, err
		}

		if allowedCIDRsSupported {
			// Skip reconciliation if network status is nil (e.g. during clusterctl move)
			if openStackCluster.Status.Network != nil {
				if err := s.getOrUpdateAllowedCIDRS(openStackCluster, listener); err != nil {
					return false, err
				}
				allowedCIDRs = listener.AllowedCIDRs
			}
		}
	}

	openStackCluster.Status.APIServerLoadBalancer = &infrav1.LoadBalancer{
		Name:         lb.Name,
		ID:           lb.ID,
		InternalIP:   lb.VipAddress,
		IP:           lbFloatingIP,
		AllowedCIDRs: allowedCIDRs,
		Tags:         lb.Tags,
	}
	return false, nil
}

func (s *Service) getOrCreateLoadBalancer(openStackCluster *infrav1.OpenStackCluster, loadBalancerName, subnetID, clusterName, vipAddress, provider string) (*loadbalancers.LoadBalancer, error) {
	lb, err := s.checkIfLbExists(loadBalancerName)
	if err != nil {
		return nil, err
	}

	if lb != nil {
		return lb, nil
	}

	s.scope.Logger().Info("Creating load balancer in subnet", "subnetID", subnetID, "name", loadBalancerName)

	lbCreateOpts := loadbalancers.CreateOpts{
		Name:        loadBalancerName,
		VipSubnetID: subnetID,
		VipAddress:  vipAddress,
		Description: names.GetDescription(clusterName),
		Provider:    provider,
		Tags:        openStackCluster.Spec.Tags,
	}
	lb, err = s.loadbalancerClient.CreateLoadBalancer(lbCreateOpts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateLoadBalancer", "Failed to create load balancer %s: %v", loadBalancerName, err)
		return nil, err
	}

	record.Eventf(openStackCluster, "SuccessfulCreateLoadBalancer", "Created load balancer %s with id %s", loadBalancerName, lb.ID)
	return lb, nil
}

func (s *Service) getOrCreateListener(openStackCluster *infrav1.OpenStackCluster, listenerName, lbID string, port int) (*listeners.Listener, error) {
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
	}
	listener, err = s.loadbalancerClient.CreateListener(listenerCreateOpts)
	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateListener", "Failed to create listener %s: %v", listenerName, err)
		return nil, err
	}

	if err := s.waitForLoadBalancerActive(lbID); err != nil {
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

func (s *Service) getOrUpdateAllowedCIDRS(openStackCluster *infrav1.OpenStackCluster, listener *listeners.Listener) error {
	allowedCIDRs := []string{}

	if len(openStackCluster.Spec.APIServerLoadBalancer.AllowedCIDRs) > 0 {
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

			if len(openStackCluster.Status.Router.IPs) > 0 {
				allowedCIDRs = append(allowedCIDRs, openStackCluster.Status.Router.IPs...)
			}
		}
	}

	// Validate CIDRs and convert any given IP into a CIDR.
	allowedCIDRs = validateIPs(openStackCluster, allowedCIDRs)

	// Remove duplicates.
	allowedCIDRs = capostrings.Unique(allowedCIDRs)
	listener.AllowedCIDRs = capostrings.Unique(listener.AllowedCIDRs)

	if !reflect.DeepEqual(allowedCIDRs, listener.AllowedCIDRs) {
		s.scope.Logger().Info("CIDRs do not match, updating listener", "expectedCIDRs", allowedCIDRs, "currentCIDRs", listener.AllowedCIDRs)
		listenerUpdateOpts := listeners.UpdateOpts{
			AllowedCIDRs: &allowedCIDRs,
		}

		listener, err := s.loadbalancerClient.UpdateListener(listener.ID, listenerUpdateOpts)
		if err != nil {
			record.Warnf(openStackCluster, "FailedUpdateListener", "Failed to update listener %s: %v", listener.Name, err)
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

// validateIPs validates given IPs/CIDRs and removes non valid network objects.
func validateIPs(openStackCluster *infrav1.OpenStackCluster, definedCIDRs []string) []string {
	marshaledCIDRs := []string{}

	for _, v := range definedCIDRs {
		switch {
		case net.IsIPv4String(v):
			marshaledCIDRs = append(marshaledCIDRs, v+"/32")
		case net.IsIPv4CIDRString(v):
			marshaledCIDRs = append(marshaledCIDRs, v)
		default:
			record.Warnf(openStackCluster, "FailedIPAddressValidation", "%s is not a valid IPv4 nor CIDR address and will not get applied to allowed_cidrs", v)
		}
	}

	return marshaledCIDRs
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

	if err := s.waitForLoadBalancerActive(lbID); err != nil {
		record.Warnf(openStackCluster, "FailedCreatePool", "Failed to create pool %s with id %s: wait for load balancer active %s: %v", poolName, pool.ID, lbID, err)
		return nil, err
	}

	record.Eventf(openStackCluster, "SuccessfulCreatePool", "Created pool %s with id %s", poolName, pool.ID)
	return pool, nil
}

func (s *Service) getOrCreateMonitor(openStackCluster *infrav1.OpenStackCluster, monitorName, poolID, lbID string) error {
	monitor, err := s.checkIfMonitorExists(monitorName)
	if err != nil {
		return err
	}

	if monitor != nil {
		return nil
	}

	s.scope.Logger().Info("Creating load balancer monitor for pool", "loadBalancerID", lbID, "name", monitorName, "poolID", poolID)

	monitorCreateOpts := monitors.CreateOpts{
		Name:           monitorName,
		PoolID:         poolID,
		Type:           "TCP",
		Delay:          10,
		MaxRetries:     5,
		MaxRetriesDown: 3,
		Timeout:        5,
	}
	monitor, err = s.loadbalancerClient.CreateMonitor(monitorCreateOpts)
	// Skip creating monitor if it is not supported by Octavia provider
	if capoerrors.IsNotImplementedError(err) {
		record.Warnf(openStackCluster, "SkippedCreateMonitor", "Health Monitor is not created as it's not implemented with the current Octavia provider.")
		return nil
	}

	if err != nil {
		record.Warnf(openStackCluster, "FailedCreateMonitor", "Failed to create monitor %s: %v", monitorName, err)
		return err
	}

	if err = s.waitForLoadBalancerActive(lbID); err != nil {
		record.Warnf(openStackCluster, "FailedCreateMonitor", "Failed to create monitor %s with id %s: wait for load balancer active %s: %v", monitorName, monitor.ID, lbID, err)
		return err
	}

	record.Eventf(openStackCluster, "SuccessfulCreateMonitor", "Created monitor %s with id %s", monitorName, monitor.ID)
	return nil
}

func (s *Service) ReconcileLoadBalancerMember(openStackCluster *infrav1.OpenStackCluster, openStackMachine *infrav1.OpenStackMachine, clusterName, ip string) error {
	if openStackCluster.Status.Network == nil {
		return errors.New("network is not yet available in openStackCluster.Status")
	}
	if len(openStackCluster.Status.Network.Subnets) == 0 {
		return errors.New("network.Subnets are not yet available in openStackCluster.Status")
	}
	if openStackCluster.Status.APIServerLoadBalancer == nil {
		return errors.New("network.APIServerLoadBalancer is not yet available in openStackCluster.Status")
	}

	loadBalancerName := getLoadBalancerName(clusterName)
	s.scope.Logger().Info("Reconciling load balancer member", "loadBalancerName", loadBalancerName)

	lbID := openStackCluster.Status.APIServerLoadBalancer.ID
	portList := []int{int(openStackCluster.Spec.ControlPlaneEndpoint.Port)}
	portList = append(portList, openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts...)
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
			err = s.waitForLoadBalancerActive(lbID)
			if err != nil {
				return err
			}
			if err := s.loadbalancerClient.DeletePoolMember(pool.ID, lbMember.ID); err != nil {
				return err
			}
			err = s.waitForLoadBalancerActive(lbID)
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

		if err := s.waitForLoadBalancerActive(lbID); err != nil {
			return err
		}

		if _, err := s.loadbalancerClient.CreatePoolMember(pool.ID, lbMemberOpts); err != nil {
			return err
		}

		if err := s.waitForLoadBalancerActive(lbID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) DeleteLoadBalancer(openStackCluster *infrav1.OpenStackCluster, clusterName string) error {
	loadBalancerName := getLoadBalancerName(clusterName)
	lb, err := s.checkIfLbExists(loadBalancerName)
	if err != nil {
		return err
	}

	if lb == nil {
		return nil
	}

	if lb.VipPortID != "" {
		fip, err := s.networkingService.GetFloatingIPByPortID(lb.VipPortID)
		if err != nil {
			return err
		}

		if fip != nil && fip.FloatingIP != "" {
			if err = s.networkingService.DisassociateFloatingIP(openStackCluster, fip.FloatingIP); err != nil {
				return err
			}
			if err = s.networkingService.DeleteFloatingIP(openStackCluster, fip.FloatingIP); err != nil {
				return err
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
		return err
	}

	record.Eventf(openStackCluster, "SuccessfulDeleteLoadBalancer", "Deleted load balancer %s with id %s", lb.Name, lb.ID)
	return nil
}

func (s *Service) DeleteLoadBalancerMember(openStackCluster *infrav1.OpenStackCluster, machine *clusterv1.Machine, openStackMachine *infrav1.OpenStackMachine, clusterName string) error {
	if openStackMachine == nil || !util.IsControlPlaneMachine(machine) {
		return nil
	}

	loadBalancerName := getLoadBalancerName(clusterName)
	lb, err := s.checkIfLbExists(loadBalancerName)
	if err != nil {
		return err
	}
	if lb == nil {
		// nothing to do
		return nil
	}

	lbID := lb.ID

	portList := []int{int(openStackCluster.Spec.ControlPlaneEndpoint.Port)}
	portList = append(portList, openStackCluster.Spec.APIServerLoadBalancer.AdditionalPorts...)
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
			err = s.waitForLoadBalancerActive(lbID)
			if err != nil {
				return err
			}
			if err := s.loadbalancerClient.DeletePoolMember(pool.ID, lbMember.ID); err != nil {
				return err
			}
			err = s.waitForLoadBalancerActive(lbID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getLoadBalancerName(clusterName string) string {
	return fmt.Sprintf("%s-cluster-%s-%s", networkPrefix, clusterName, kubeapiLBSuffix)
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
func (s *Service) waitForLoadBalancerActive(id string) error {
	s.scope.Logger().Info("Waiting for load balancer", "id", id, "targetStatus", "ACTIVE")
	return wait.ExponentialBackoff(backoff, func() (bool, error) {
		lb, err := s.loadbalancerClient.GetLoadBalancer(id)
		if err != nil {
			return false, err
		}
		return lb.ProvisioningStatus == loadBalancerProvisioningStatusActive, nil
	})
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
