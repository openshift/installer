/*
Copyright 2021 The Kubernetes Authors.

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

package loadbalancers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// loadBalancingMode describes the load balancing mode that the backend performs.
type loadBalancingMode string

const (
	// Utilization determines how the traffic load is spread based on the
	// utilization of instances.
	loadBalancingModeUtilization = loadBalancingMode("UTILIZATION")

	// Connection determines how the traffic load is spread based on the
	// total number of connections that a backend can handle. This is
	// only mode available for passthrough Load Balancers.
	loadBalancingModeConnection = loadBalancingMode("CONNECTION")

	loadBalanceTrafficInternal = "INTERNAL"
)

// Reconcile reconcile cluster control-plane loadbalancer components.
func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling loadbalancer resources")

	// Creates instance groups used by load balancer(s)
	instancegroups, err := s.createOrGetInstanceGroups(ctx)
	if err != nil {
		return err
	}

	lbSpec := s.scope.LoadBalancer()
	lbType := ptr.Deref(lbSpec.LoadBalancerType, infrav1.External)
	// Create a Global External Proxy Load Balancer by default
	if lbType == infrav1.External || lbType == infrav1.InternalExternal {
		if err = s.createExternalLoadBalancer(ctx, lbType, instancegroups); err != nil {
			return err
		}
	}

	// Create a Regional Internal Passthrough Load Balancer if configured
	if lbType == infrav1.Internal || lbType == infrav1.InternalExternal {
		name := infrav1.InternalRoleTagValue
		if lbSpec.InternalLoadBalancer != nil {
			name = ptr.Deref(lbSpec.InternalLoadBalancer.Name, infrav1.InternalRoleTagValue)
		}
		if err = s.createInternalLoadBalancer(ctx, name, lbType, instancegroups); err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes cluster control-plane loadbalancer components.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	var allErrs []error
	lbSpec := s.scope.LoadBalancer()
	lbType := ptr.Deref(lbSpec.LoadBalancerType, infrav1.External)
	if lbType == infrav1.External || lbType == infrav1.InternalExternal {
		if err := s.deleteExternalLoadBalancer(ctx); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if lbType == infrav1.Internal || lbType == infrav1.InternalExternal {
		name := infrav1.InternalRoleTagValue
		if lbSpec.InternalLoadBalancer != nil {
			name = ptr.Deref(lbSpec.InternalLoadBalancer.Name, infrav1.InternalRoleTagValue)
		}
		if err := s.deleteInternalLoadBalancer(ctx, name); err != nil {
			allErrs = append(allErrs, err)
		}
	}
	if err := s.deleteInstanceGroups(ctx); err != nil {
		log.Error(err, "Error deleting instancegroup")
		allErrs = append(allErrs, err)
	}

	return errors.Join(allErrs...)
}

func (s *Service) deleteExternalLoadBalancer(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting external loadbalancer resources")
	name := infrav1.APIServerRoleTagValue
	if err := s.deleteForwardingRule(ctx, name); err != nil {
		return fmt.Errorf("deleting ForwardingRule: %w", err)
	}
	s.scope.Network().APIServerForwardingRule = nil

	if err := s.deleteAddress(ctx, name); err != nil {
		return fmt.Errorf("deleting Address: %w", err)
	}
	s.scope.Network().APIServerAddress = nil

	if err := s.deleteTargetTCPProxy(ctx); err != nil {
		return fmt.Errorf("deleting TargetTCPProxy: %w", err)
	}
	s.scope.Network().APIServerTargetProxy = nil

	if err := s.deleteBackendService(ctx, name); err != nil {
		return fmt.Errorf("deleting BackendService: %w", err)
	}
	s.scope.Network().APIServerBackendService = nil

	if err := s.deleteHealthCheck(ctx, name); err != nil {
		return fmt.Errorf("deleting HealthCheck: %w", err)
	}
	s.scope.Network().APIServerHealthCheck = nil

	return nil
}

func (s *Service) deleteInternalLoadBalancer(ctx context.Context, name string) error {
	log := log.FromContext(ctx)
	log.Info("Deleting internal loadbalancer resources")
	if err := s.deleteRegionalForwardingRule(ctx, name); err != nil {
		return fmt.Errorf("deleting ForwardingRule: %w", err)
	}
	s.scope.Network().APIInternalForwardingRule = nil

	if err := s.deleteInternalAddress(ctx, name); err != nil {
		return fmt.Errorf("deleting InternalAddress: %w", err)
	}
	s.scope.Network().APIInternalAddress = nil

	if err := s.deleteRegionalBackendService(ctx, name); err != nil {
		return fmt.Errorf("deleting RegionalBackendService: %w", err)
	}
	s.scope.Network().APIInternalBackendService = nil

	if err := s.deleteRegionalHealthCheck(ctx, name); err != nil {
		return fmt.Errorf("deleting RegionalHealthCheck: %w", err)
	}
	s.scope.Network().APIInternalHealthCheck = nil

	return nil
}

// createExternalLoadBalancer creates the components for a Global External Proxy LoadBalancer.
func (s *Service) createExternalLoadBalancer(ctx context.Context, lbType infrav1.LoadBalancerType, instancegroups []*compute.InstanceGroup) error {
	name := infrav1.APIServerRoleTagValue
	healthchecks, err := s.createOrGetHealthChecks(ctx, name)
	if err != nil {
		return err
	}
	for _, healthcheck := range healthchecks {
		if (strings.HasSuffix(healthcheck.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.DualStackType) ||
			(!strings.HasSuffix(healthcheck.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.IPv4OnlyStackType) {
			s.scope.Network().APIServerHealthCheck = ptr.To[string](healthcheck.SelfLink)
		}
	}

	// If an Internal LoadBalancer is being created, the BalancingMode must match the Internal LB.
	// which must be CONNECTION for Internal Proxy Load Balancers, see
	// https://cloud.google.com/load-balancing/docs/backend-service#balancing-mode-lb
	mode := loadBalancingModeUtilization
	if lbType == infrav1.InternalExternal {
		mode = loadBalancingModeConnection
	}
	backendsvc, err := s.createOrGetBackendService(ctx, name, mode, instancegroups, healthchecks)
	if err != nil {
		return err
	}
	s.scope.Network().APIServerBackendService = ptr.To[string](backendsvc.SelfLink)

	// Create TargetTCPProxy for Proxy Load Balancer
	target, err := s.createOrGetTargetTCPProxy(ctx, backendsvc)
	if err != nil {
		return err
	}
	s.scope.Network().APIServerTargetProxy = ptr.To[string](target.SelfLink)

	addrs, err := s.createOrGetAddress(ctx, name)
	if err != nil {
		return err
	}

	endpoint := s.scope.ControlPlaneEndpoint()
	for _, addr := range addrs {
		if (strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.DualStackType) ||
			(!strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.IPv4OnlyStackType) {
			s.scope.Network().APIServerAddress = ptr.To[string](addr.SelfLink)
			endpoint.Host = addr.Address
		}
	}
	s.scope.SetControlPlaneEndpoint(endpoint)

	forwardingrules, err := s.createOrGetForwardingRules(ctx, name, target, addrs)
	if err != nil {
		return err
	}
	for _, forwardingrule := range forwardingrules {
		if (strings.HasSuffix(forwardingrule.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.DualStackType) ||
			(!strings.HasSuffix(forwardingrule.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.IPv4OnlyStackType) {
			s.scope.Network().APIServerForwardingRule = ptr.To[string](forwardingrule.SelfLink)
		}
	}

	return nil
}

// createInternalLoadBalancer creates the components for a Regional Internal Passthrough LoadBalancer.
// Since this is a passthrough LoadBalancer the TargetTCPProxy resource is not created.
func (s *Service) createInternalLoadBalancer(ctx context.Context, name string, lbType infrav1.LoadBalancerType, instancegroups []*compute.InstanceGroup) error {
	healthchecks, err := s.createOrGetRegionalHealthChecks(ctx, name)
	if err != nil {
		return err
	}

	for _, healthcheck := range healthchecks {
		if (strings.HasSuffix(healthcheck.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.DualStackType) ||
			(!strings.HasSuffix(healthcheck.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.IPv4OnlyStackType) {
			s.scope.Network().APIInternalHealthCheck = ptr.To[string](healthcheck.SelfLink)
		}
	}

	backendsvc, err := s.createOrGetRegionalBackendService(ctx, name, instancegroups, healthchecks)
	if err != nil {
		return err
	}
	s.scope.Network().APIInternalBackendService = ptr.To[string](backendsvc.SelfLink)

	// Create an address on internal subnet.
	addrs, err := s.createOrGetInternalAddresses(ctx, name)
	if err != nil {
		return err
	}
	for _, addr := range addrs {
		if (strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.DualStackType) ||
			(!strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.IPv4OnlyStackType) {
			s.scope.Network().APIInternalAddress = ptr.To[string](addr.SelfLink)

			if lbType == infrav1.Internal {
				// If only creating an internal Load Balancer, set the control plane endpoint
				endpoint := s.scope.ControlPlaneEndpoint()
				endpoint.Host = addr.Address
				s.scope.SetControlPlaneEndpoint(endpoint)
			}
		}
	}

	// Create a regional forwarding rule to the backend service
	forwardingrules, err := s.createOrGetRegionalForwardingRules(ctx, name, backendsvc, addrs)
	if err != nil {
		return err
	}
	for _, forwardingrule := range forwardingrules {
		if (strings.HasSuffix(forwardingrule.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.DualStackType) ||
			(!strings.HasSuffix(forwardingrule.Name, infrav1.DualStackAdditionalResourceSuffix) && s.scope.StackType() == infrav1.IPv4OnlyStackType) {
			s.scope.Network().APIInternalForwardingRule = ptr.To[string](forwardingrule.SelfLink)
		}
	}

	return nil
}

func (s *Service) createOrGetInstanceGroups(ctx context.Context) ([]*compute.InstanceGroup, error) {
	log := log.FromContext(ctx)
	zones := s.scope.FailureDomains()

	groups := make([]*compute.InstanceGroup, 0, len(zones))
	groupsMap := s.scope.Network().APIServerInstanceGroups
	if groupsMap == nil {
		groupsMap = make(map[string]string)
	}

	for _, zone := range zones {
		instancegroupSpec := s.scope.InstanceGroupSpec(zone)
		log.V(2).Info("Looking for instancegroup in zone", "zone", zone, "name", instancegroupSpec.Name)
		instancegroup, err := s.instancegroups.Get(ctx, meta.ZonalKey(instancegroupSpec.Name, zone))
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error looking for instancegroup in zone", "zone", zone)
				return groups, err
			}

			log.V(2).Info("Creating instancegroup in zone", "zone", zone, "name", instancegroupSpec.Name)
			if err := s.instancegroups.Insert(ctx, meta.ZonalKey(instancegroupSpec.Name, zone), instancegroupSpec); err != nil {
				log.Error(err, "Error creating instancegroup", "name", instancegroupSpec.Name)
				return groups, err
			}

			instancegroup, err = s.instancegroups.Get(ctx, meta.ZonalKey(instancegroupSpec.Name, zone))
			if err != nil {
				return groups, err
			}
		}

		groups = append(groups, instancegroup)
		groupsMap[zone] = instancegroup.SelfLink
	}

	s.scope.Network().APIServerInstanceGroups = groupsMap
	return groups, nil
}

func (s *Service) createOrGetHealthChecks(ctx context.Context, lbname string) ([]*compute.HealthCheck, error) {
	log := log.FromContext(ctx)
	healthcheckSpecs := s.scope.HealthCheckSpecs(lbname)
	healthchecks := []*compute.HealthCheck{}
	for _, healthcheckSpec := range healthcheckSpecs {
		log.V(2).Info("Looking for healthcheck", "name", healthcheckSpec.Name)
		key := meta.GlobalKey(healthcheckSpec.Name)
		healthcheck, err := s.healthchecks.Get(ctx, key)
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error looking for healthcheck", "name", healthcheckSpec.Name)
				return nil, err
			}

			log.V(2).Info("Creating a healthcheck", "name", healthcheckSpec.Name)
			if err := s.healthchecks.Insert(ctx, key, healthcheckSpec); err != nil {
				log.Error(err, "Error creating a healthcheck", "name", healthcheckSpec.Name)
				return nil, err
			}

			healthcheck, err = s.healthchecks.Get(ctx, key)
			if err != nil {
				return nil, err
			}
		}
		healthchecks = append(healthchecks, healthcheck)
	}
	return healthchecks, nil
}

func (s *Service) createOrGetRegionalHealthChecks(ctx context.Context, lbname string) ([]*compute.HealthCheck, error) {
	log := log.FromContext(ctx)
	healthcheckSpecs := s.scope.HealthCheckSpecs(lbname)
	healthchecks := []*compute.HealthCheck{}
	for _, healthcheckSpec := range healthcheckSpecs {
		healthcheckSpec.Region = s.scope.Region()
		log.V(2).Info("Looking for regional healthcheck", "name", healthcheckSpec.Name)
		key := meta.RegionalKey(healthcheckSpec.Name, s.scope.Region())
		healthcheck, err := s.regionalhealthchecks.Get(ctx, key)
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error looking for regional healthcheck", "name", healthcheckSpec.Name)
				return nil, err
			}

			log.V(2).Info("Creating a regional healthcheck", "name", healthcheckSpec.Name)
			if err := s.regionalhealthchecks.Insert(ctx, key, healthcheckSpec); err != nil {
				log.Error(err, "Error creating a regional healthcheck", "name", healthcheckSpec.Name)
				return nil, err
			}

			healthcheck, err = s.regionalhealthchecks.Get(ctx, key)
			if err != nil {
				return nil, err
			}
		}
		healthchecks = append(healthchecks, healthcheck)
	}
	return healthchecks, nil
}

func (s *Service) createOrGetBackendService(ctx context.Context, lbname string, mode loadBalancingMode, instancegroups []*compute.InstanceGroup, healthchecks []*compute.HealthCheck) (*compute.BackendService, error) {
	log := log.FromContext(ctx)
	backends := make([]*compute.Backend, 0, len(instancegroups))
	for _, group := range instancegroups {
		be := &compute.Backend{
			BalancingMode: string(mode),
			Group:         group.SelfLink,
		}
		if mode == loadBalancingModeConnection {
			// Set max connections to a reasonable limit based
			// on database max connections https://cloud.google.com/sql/docs/postgres/flags#postgres-m
			be.MaxConnections = 1000
		}
		backends = append(backends, be)
	}

	backendsvcSpec := s.scope.BackendServiceSpec(lbname)
	backendsvcSpec.Backends = backends
	healthcheckLinks := []string{}
	for _, healthcheck := range healthchecks {
		healthcheckLinks = append(healthcheckLinks, healthcheck.SelfLink)
	}
	backendsvcSpec.HealthChecks = healthcheckLinks

	key := meta.GlobalKey(backendsvcSpec.Name)
	backendsvc, err := s.backendservices.Get(ctx, key)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}

		log.V(2).Info("Creating a backendservice", "name", backendsvcSpec.Name)
		if err := s.backendservices.Insert(ctx, key, backendsvcSpec); err != nil {
			log.Error(err, "Error creating a backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}

		backendsvc, err = s.backendservices.Get(ctx, key)
		if err != nil {
			return nil, err
		}
	}

	if len(backendsvc.Backends) != len(backendsvcSpec.Backends) {
		log.V(2).Info("Updating a backendservice", "name", backendsvcSpec.Name)
		backendsvc.Backends = backendsvcSpec.Backends
		if err := s.backendservices.Update(ctx, key, backendsvc); err != nil {
			log.Error(err, "Error updating a backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}
	}

	return backendsvc, nil
}

// createOrGetRegionalBackendService is used for internal passthrough load balancers.
func (s *Service) createOrGetRegionalBackendService(ctx context.Context, lbname string, instancegroups []*compute.InstanceGroup, healthchecks []*compute.HealthCheck) (*compute.BackendService, error) {
	log := log.FromContext(ctx)
	backends := make([]*compute.Backend, 0, len(instancegroups))
	for _, group := range instancegroups {
		be := &compute.Backend{
			// Always use connection mode for passthrough load balancer
			BalancingMode: string(loadBalancingModeConnection),
			Group:         group.SelfLink,
		}
		backends = append(backends, be)
	}

	backendsvcSpec := s.scope.BackendServiceSpec(lbname)
	backendsvcSpec.Backends = backends
	checks := []string{}
	for _, healthcheck := range healthchecks {
		checks = append(checks, healthcheck.SelfLink)
	}
	backendsvcSpec.HealthChecks = checks
	backendsvcSpec.Region = s.scope.Region()
	backendsvcSpec.LoadBalancingScheme = string(loadBalanceTrafficInternal)
	backendsvcSpec.PortName = ""
	network := s.scope.Network()
	if network.SelfLink != nil {
		backendsvcSpec.Network = *network.SelfLink
	}

	key := meta.RegionalKey(backendsvcSpec.Name, s.scope.Region())
	backendsvc, err := s.regionalbackendservices.Get(ctx, key)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for regional backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}

		log.V(2).Info("Creating a regional backendservice", "name", backendsvcSpec.Name)
		if err := s.regionalbackendservices.Insert(ctx, key, backendsvcSpec); err != nil {
			log.Error(err, "Error creating a regional backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}

		backendsvc, err = s.regionalbackendservices.Get(ctx, key)
		if err != nil {
			return nil, err
		}
	}

	if len(backendsvc.Backends) != len(backendsvcSpec.Backends) {
		log.V(2).Info("Updating a regional backendservice", "name", backendsvcSpec.Name)
		backendsvc.Backends = backendsvcSpec.Backends
		if err := s.regionalbackendservices.Update(ctx, key, backendsvc); err != nil {
			log.Error(err, "Error updating a regional backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}
	}

	return backendsvc, nil
}

func (s *Service) createOrGetTargetTCPProxy(ctx context.Context, service *compute.BackendService) (*compute.TargetTcpProxy, error) {
	log := log.FromContext(ctx)
	targetSpec := s.scope.TargetTCPProxySpec()
	targetSpec.Service = service.SelfLink
	key := meta.GlobalKey(targetSpec.Name)
	target, err := s.targettcpproxies.Get(ctx, key)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for targettcpproxy", "name", targetSpec.Name)
			return nil, err
		}

		log.V(2).Info("Creating a targettcpproxy", "name", targetSpec.Name)
		if err := s.targettcpproxies.Insert(ctx, key, targetSpec); err != nil {
			log.Error(err, "Error creating a targettcpproxy", "name", targetSpec.Name)
			return nil, err
		}

		target, err = s.targettcpproxies.Get(ctx, key)
		if err != nil {
			return nil, err
		}
	}

	return target, nil
}

// createOrGetAddress is used to obtain a Global address.
func (s *Service) createOrGetAddress(ctx context.Context, lbname string) ([]*compute.Address, error) {
	log := log.FromContext(ctx)
	addrSpecs := s.scope.AddressSpecs(lbname)

	addrs := []*compute.Address{}
	for _, addrSpec := range addrSpecs {
		log.V(2).Info("Looking for address", "name", addrSpec.Name)
		key := meta.GlobalKey(addrSpec.Name)
		addr, err := s.addresses.Get(ctx, key)
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error looking for address", "name", addrSpec.Name)
				return nil, err
			}

			log.V(2).Info("Creating an address", "name", addrSpec.Name)
			if err := s.addresses.Insert(ctx, key, addrSpec); err != nil {
				log.Error(err, "Error creating an address", "name", addrSpec.Name)
				return nil, err
			}

			addr, err = s.addresses.Get(ctx, key)
			if err != nil {
				return nil, err
			}
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// createOrGetInternalAddress is used to obtain an internal address.
func (s *Service) createOrGetInternalAddresses(ctx context.Context, lbname string) ([]*compute.Address, error) {
	log := log.FromContext(ctx)
	addrSpecs := s.scope.AddressSpecs(lbname)
	addrs := []*compute.Address{}
	for _, addrSpec := range addrSpecs {
		addrSpec.AddressType = string(loadBalanceTrafficInternal)
		addrSpec.Region = s.scope.Region()
		subnet, err := s.getSubnet(ctx)
		if err != nil {
			log.Error(err, "Error getting subnet for Internal Load Balancer")
			return nil, err
		}
		lbSpec := s.scope.LoadBalancer()
		if lbSpec.InternalLoadBalancer != nil && lbSpec.InternalLoadBalancer.IPAddress != nil {
			// If an IP address is configured, use it instead of creating a new one.
			addrSpec.Address = *lbSpec.InternalLoadBalancer.IPAddress
		}
		addrSpec.Subnetwork = subnet.SelfLink
		addrSpec.Purpose = "GCE_ENDPOINT"
		log.V(2).Info("Looking for internal address", "name", addrSpec.Name)
		key := meta.RegionalKey(addrSpec.Name, s.scope.Region())
		addr, err := s.internaladdresses.Get(ctx, key)
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error looking for internal address", "name", addrSpec.Name)
				return nil, err
			}

			log.V(2).Info("Creating an internal address", "name", addrSpec.Name)
			if err := s.internaladdresses.Insert(ctx, key, addrSpec); err != nil {
				log.Error(err, "Error creating an internal address", "name", addrSpec.Name)
				return nil, err
			}

			addr, err = s.internaladdresses.Get(ctx, key)
			if err != nil {
				return nil, err
			}
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

// createOrGetForwardingRule is used obtain a Global ForwardingRule.
func (s *Service) createOrGetForwardingRules(ctx context.Context, lbname string, target *compute.TargetTcpProxy, addrs []*compute.Address) ([]*compute.ForwardingRule, error) {
	log := log.FromContext(ctx)
	specs := s.scope.ForwardingRuleSpecs(lbname)

	forwardingRules := []*compute.ForwardingRule{}
	for _, spec := range specs {
		spec.Target = target.SelfLink
		for _, addr := range addrs {
			if (strings.HasSuffix(spec.Name, infrav1.DualStackAdditionalResourceSuffix) && strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix)) ||
				(!strings.HasSuffix(spec.Name, infrav1.DualStackAdditionalResourceSuffix) && !strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix)) {
				spec.IPAddress = addr.SelfLink
			}
		}

		key := meta.GlobalKey(spec.Name)
		log.V(2).Info("Looking for forwardingrule", "name", spec.Name)
		forwarding, err := s.forwardingrules.Get(ctx, key)
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error looking for forwardingrule", "name", spec.Name)
				return nil, err
			}

			log.V(2).Info("Creating a forwardingrule", "name", spec.Name)
			if err := s.forwardingrules.Insert(ctx, key, spec); err != nil {
				log.Error(err, "Error creating a forwardingrule", "name", spec.Name)
				return nil, err
			}

			forwarding, err = s.forwardingrules.Get(ctx, key)
			if err != nil {
				return nil, err
			}
		}

		// Labels on ForwardingRules must be added after resource is created
		labels := s.scope.AdditionalLabels()
		if !labels.Equals(forwarding.Labels) {
			setLabelsRequest := &compute.GlobalSetLabelsRequest{
				LabelFingerprint: forwarding.LabelFingerprint,
				Labels:           labels,
			}
			if err = s.forwardingrules.SetLabels(ctx, key, setLabelsRequest); err != nil {
				return nil, err
			}
		}
		forwardingRules = append(forwardingRules, forwarding)
	}

	return forwardingRules, nil
}

// createOrGetRegionalForwardingRule is used to obtain a Regional ForwardingRule.
func (s *Service) createOrGetRegionalForwardingRules(ctx context.Context, lbname string, backendSvc *compute.BackendService, addrs []*compute.Address) ([]*compute.ForwardingRule, error) {
	log := log.FromContext(ctx)
	specs := s.scope.ForwardingRuleSpecs(lbname)
	forwardingrules := []*compute.ForwardingRule{}
	for _, spec := range specs {
		spec.LoadBalancingScheme = string(loadBalanceTrafficInternal)
		spec.Region = s.scope.Region()
		spec.BackendService = backendSvc.SelfLink
		lbSpec := s.scope.LoadBalancer()
		if lbSpec.InternalLoadBalancer != nil && lbSpec.InternalLoadBalancer.InternalAccess == infrav1.InternalAccessGlobal {
			spec.AllowGlobalAccess = true
		}
		// Ports is used instead or PortRange for passthrough Load Balancer
		// Configure ports for k8s API to match the external API which is the first port of range
		var ports []string
		portList := strings.Split(spec.PortRange, "-")
		ports = append(ports, portList[0])
		// Also configure ignition port
		ports = append(ports, "22623")
		spec.Ports = ports
		spec.PortRange = ""
		subnet, err := s.getSubnet(ctx)
		if err != nil {
			log.Error(err, "Error getting subnet for regional forwardingrule")
			return nil, err
		}
		spec.Subnetwork = subnet.SelfLink
		for _, addr := range addrs {
			if (strings.HasSuffix(spec.Name, infrav1.DualStackAdditionalResourceSuffix) && strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix)) ||
				(!strings.HasSuffix(spec.Name, infrav1.DualStackAdditionalResourceSuffix) && !strings.HasSuffix(addr.Name, infrav1.DualStackAdditionalResourceSuffix)) {
				spec.IPAddress = addr.SelfLink
			}
		}

		key := meta.RegionalKey(spec.Name, s.scope.Region())
		log.V(2).Info("Looking for regional forwardingrule", "name", spec.Name)
		forwarding, err := s.regionalforwardingrules.Get(ctx, key)
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error looking for regional forwardingrule", "name", spec.Name)
				return nil, err
			}

			log.V(2).Info("Creating a regional forwardingrule", "name", spec.Name)
			if err := s.regionalforwardingrules.Insert(ctx, key, spec); err != nil {
				log.Error(err, "Error creating a regional forwardingrule", "name", spec.Name)
				return nil, err
			}

			forwarding, err = s.regionalforwardingrules.Get(ctx, key)
			if err != nil {
				return nil, err
			}
		}

		// Labels on ForwardingRules must be added after resource is created
		labels := s.scope.AdditionalLabels()
		if !labels.Equals(forwarding.Labels) {
			setLabelsRequest := &compute.RegionSetLabelsRequest{
				LabelFingerprint: forwarding.LabelFingerprint,
				Labels:           labels,
			}
			if err = s.regionalforwardingrules.SetLabels(ctx, key, setLabelsRequest); err != nil {
				return nil, err
			}
		}

		forwardingrules = append(forwardingrules, forwarding)
	}

	return forwardingrules, nil
}

func (s *Service) deleteForwardingRule(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	specs := s.scope.ForwardingRuleSpecs(lbname)
	for _, spec := range specs {
		key := meta.GlobalKey(spec.Name)
		log.V(2).Info("Deleting a forwardingrule", "name", spec.Name)
		if err := s.forwardingrules.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
			log.Error(err, "Error updating a forwardingrule", "name", spec.Name)
			return err
		}
	}
	return nil
}

func (s *Service) deleteRegionalForwardingRule(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	specs := s.scope.ForwardingRuleSpecs(lbname)
	for _, spec := range specs {
		key := meta.RegionalKey(spec.Name, s.scope.Region())
		log.V(2).Info("Deleting a regional forwardingrule", "name", spec.Name)
		if err := s.regionalforwardingrules.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
			log.Error(err, "Error updating a regional forwardingrule", "name", spec.Name)
			return err
		}
	}

	return nil
}

func (s *Service) deleteAddress(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	specs := s.scope.AddressSpecs(lbname)
	for _, spec := range specs {
		key := meta.GlobalKey(spec.Name)
		log.V(2).Info("Deleting a address", "name", spec.Name)
		if err := s.addresses.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (s *Service) deleteInternalAddress(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	specs := s.scope.AddressSpecs(lbname)
	for _, spec := range specs {
		key := meta.RegionalKey(spec.Name, s.scope.Region())
		log.V(2).Info("Deleting an internal address", "name", spec.Name)
		if err := s.internaladdresses.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
			return err
		}
	}
	return nil
}

func (s *Service) deleteTargetTCPProxy(ctx context.Context) error {
	log := log.FromContext(ctx)
	spec := s.scope.TargetTCPProxySpec()
	key := meta.GlobalKey(spec.Name)
	log.V(2).Info("Deleting a targettcpproxy", "name", spec.Name)
	if err := s.targettcpproxies.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
		log.Error(err, "Error deleting a targettcpproxy", "name", spec.Name)
		return err
	}

	return nil
}

func (s *Service) deleteBackendService(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	spec := s.scope.BackendServiceSpec(lbname)
	key := meta.GlobalKey(spec.Name)
	log.V(2).Info("Deleting a backendservice", "name", spec.Name)
	if err := s.backendservices.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
		log.Error(err, "Error deleting a backendservice", "name", spec.Name)
		return err
	}

	return nil
}

func (s *Service) deleteRegionalBackendService(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	spec := s.scope.BackendServiceSpec(lbname)
	key := meta.RegionalKey(spec.Name, s.scope.Region())
	log.V(2).Info("Deleting a regional backendservice", "name", spec.Name)
	if err := s.regionalbackendservices.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
		log.Error(err, "Error deleting a regional backendservice", "name", spec.Name)
		return err
	}

	return nil
}

func (s *Service) deleteHealthCheck(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	specs := s.scope.HealthCheckSpecs(lbname)
	for _, spec := range specs {
		key := meta.GlobalKey(spec.Name)
		log.V(2).Info("Deleting a healthcheck", "name", spec.Name)
		if err := s.healthchecks.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
			log.Error(err, "Error deleting a healthcheck", "name", spec.Name)
			return err
		}
	}
	return nil
}

func (s *Service) deleteRegionalHealthCheck(ctx context.Context, lbname string) error {
	log := log.FromContext(ctx)
	specs := s.scope.HealthCheckSpecs(lbname)
	for _, spec := range specs {
		key := meta.RegionalKey(spec.Name, s.scope.Region())
		log.V(2).Info("Deleting a regional healthcheck", "name", spec.Name)
		if err := s.regionalhealthchecks.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
			log.Error(err, "Error deleting a regional healthcheck", "name", spec.Name)
			return err
		}
	}
	return nil
}

func (s *Service) deleteInstanceGroups(ctx context.Context) error {
	log := log.FromContext(ctx)
	for zone := range s.scope.Network().APIServerInstanceGroups {
		spec := s.scope.InstanceGroupSpec(zone)
		key := meta.ZonalKey(spec.Name, zone)
		log.V(2).Info("Deleting a instancegroup", "name", spec.Name)
		if err := s.instancegroups.Delete(ctx, key); err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error deleting a instancegroup", "name", spec.Name)
				return err
			}

			delete(s.scope.Network().APIServerInstanceGroups, zone)
		}
	}

	return nil
}

// getSubnet gets the subnet to use for an internal Load Balancer.
func (s *Service) getSubnet(ctx context.Context) (*compute.Subnetwork, error) {
	log := log.FromContext(ctx)
	cfgSubnet := ""
	lbSpec := s.scope.LoadBalancer()
	if lbSpec.InternalLoadBalancer != nil {
		cfgSubnet = ptr.Deref(lbSpec.InternalLoadBalancer.Subnet, "")
	}
	for _, subnetSpec := range s.scope.SubnetSpecs() {
		log.V(2).Info("Looking for subnet for load balancer", "name", subnetSpec.Name)
		region := subnetSpec.Region
		if region == "" {
			region = s.scope.Region()
		}

		subnetKey := meta.RegionalKey(subnetSpec.Name, region)
		subnet, err := s.subnets.Get(ctx, subnetKey)
		if err != nil {
			return nil, err
		}
		// Return subnet that matches configuration, or first one if not configured
		if cfgSubnet == "" || strings.HasSuffix(subnet.Name, cfgSubnet) {
			return subnet, nil
		}
	}

	return nil, errors.New("could not find subnet")
}
