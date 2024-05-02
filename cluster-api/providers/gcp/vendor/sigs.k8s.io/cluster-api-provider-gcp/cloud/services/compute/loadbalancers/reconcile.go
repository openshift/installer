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

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"

	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile reconcile cluster control-plane loadbalancer compoenents.
func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling loadbalancer resources")
	instancegroups, err := s.createOrGetInstanceGroups(ctx)
	if err != nil {
		return err
	}

	healthcheck, err := s.createOrGetHealthCheck(ctx)
	if err != nil {
		return err
	}

	backendsvc, err := s.createOrGetBackendService(ctx, instancegroups, healthcheck)
	if err != nil {
		return err
	}

	target, err := s.createOrGetTargetTCPProxy(ctx, backendsvc)
	if err != nil {
		return err
	}

	addr, err := s.createOrGetAddress(ctx)
	if err != nil {
		return err
	}

	return s.createForwardingRule(ctx, target, addr)
}

// Delete delete cluster control-plane loadbalancer compoenents.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting loadbalancer resources")
	if err := s.deleteForwardingRule(ctx); err != nil {
		return err
	}

	if err := s.deleteAddress(ctx); err != nil {
		return err
	}

	if err := s.deleteTargetTCPProxy(ctx); err != nil {
		return err
	}

	if err := s.deleteBackendService(ctx); err != nil {
		return err
	}

	if err := s.deleteHealthCheck(ctx); err != nil {
		return err
	}

	return s.deleteInstanceGroups(ctx)
}

func (s *Service) createOrGetInstanceGroups(ctx context.Context) ([]*compute.InstanceGroup, error) {
	log := log.FromContext(ctx)
	fd := s.scope.FailureDomains()
	zones := make([]string, 0, len(fd))
	for zone := range fd {
		zones = append(zones, zone)
	}

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

func (s *Service) createOrGetHealthCheck(ctx context.Context) (*compute.HealthCheck, error) {
	log := log.FromContext(ctx)
	healthcheckSpec := s.scope.HealthCheckSpec()
	log.V(2).Info("Looking for healthcheck", "name", healthcheckSpec.Name)
	healthcheck, err := s.healthchecks.Get(ctx, meta.GlobalKey(healthcheckSpec.Name))
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for healthcheck", "name", healthcheckSpec.Name)
			return nil, err
		}

		log.V(2).Info("Creating a healthcheck", "name", healthcheckSpec.Name)
		if err := s.healthchecks.Insert(ctx, meta.GlobalKey(healthcheckSpec.Name), healthcheckSpec); err != nil {
			log.Error(err, "Error creating a healthcheck", "name", healthcheckSpec.Name)
			return nil, err
		}

		healthcheck, err = s.healthchecks.Get(ctx, meta.GlobalKey(healthcheckSpec.Name))
		if err != nil {
			return nil, err
		}
	}

	s.scope.Network().APIServerHealthCheck = ptr.To[string](healthcheck.SelfLink)
	return healthcheck, nil
}

func (s *Service) createOrGetBackendService(ctx context.Context, instancegroups []*compute.InstanceGroup, healthcheck *compute.HealthCheck) (*compute.BackendService, error) {
	log := log.FromContext(ctx)
	backends := make([]*compute.Backend, 0, len(instancegroups))
	for _, group := range instancegroups {
		backends = append(backends, &compute.Backend{
			BalancingMode: "UTILIZATION",
			Group:         group.SelfLink,
		})
	}

	backendsvcSpec := s.scope.BackendServiceSpec()
	backendsvcSpec.Backends = backends
	backendsvcSpec.HealthChecks = []string{healthcheck.SelfLink}
	backendsvc, err := s.backendservices.Get(ctx, meta.GlobalKey(backendsvcSpec.Name))
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}

		log.V(2).Info("Creating a backendservice", "name", backendsvcSpec.Name)
		if err := s.backendservices.Insert(ctx, meta.GlobalKey(backendsvcSpec.Name), backendsvcSpec); err != nil {
			log.Error(err, "Error creating a backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}

		backendsvc, err = s.backendservices.Get(ctx, meta.GlobalKey(backendsvcSpec.Name))
		if err != nil {
			return nil, err
		}
	}

	if len(backendsvc.Backends) != len(backendsvcSpec.Backends) {
		log.V(2).Info("Updating a backendservice", "name", backendsvcSpec.Name)
		backendsvc.Backends = backendsvcSpec.Backends
		if err := s.backendservices.Update(ctx, meta.GlobalKey(backendsvcSpec.Name), backendsvc); err != nil {
			log.Error(err, "Error updating a backendservice", "name", backendsvcSpec.Name)
			return nil, err
		}
	}

	s.scope.Network().APIServerBackendService = ptr.To[string](backendsvc.SelfLink)
	return backendsvc, nil
}

func (s *Service) createOrGetTargetTCPProxy(ctx context.Context, service *compute.BackendService) (*compute.TargetTcpProxy, error) {
	log := log.FromContext(ctx)
	targetSpec := s.scope.TargetTCPProxySpec()
	targetSpec.Service = service.SelfLink
	target, err := s.targettcpproxies.Get(ctx, meta.GlobalKey(targetSpec.Name))
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for targettcpproxy", "name", targetSpec.Name)
			return nil, err
		}

		log.V(2).Info("Creating a targettcpproxy", "name", targetSpec.Name)
		if err := s.targettcpproxies.Insert(ctx, meta.GlobalKey(targetSpec.Name), targetSpec); err != nil {
			log.Error(err, "Error creating a targettcpproxy", "name", targetSpec.Name)
			return nil, err
		}

		target, err = s.targettcpproxies.Get(ctx, meta.GlobalKey(targetSpec.Name))
		if err != nil {
			return nil, err
		}
	}

	s.scope.Network().APIServerTargetProxy = ptr.To[string](target.SelfLink)
	return target, nil
}

func (s *Service) createOrGetAddress(ctx context.Context) (*compute.Address, error) {
	log := log.FromContext(ctx)
	addrSpec := s.scope.AddressSpec()
	log.V(2).Info("Looking for address", "name", addrSpec.Name)
	addr, err := s.addresses.Get(ctx, meta.GlobalKey(addrSpec.Name))
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for address", "name", addrSpec.Name)
			return nil, err
		}

		log.V(2).Info("Creating an address", "name", addrSpec.Name)
		if err := s.addresses.Insert(ctx, meta.GlobalKey(addrSpec.Name), addrSpec); err != nil {
			log.Error(err, "Error creating an address", "name", addrSpec.Name)
			return nil, err
		}

		addr, err = s.addresses.Get(ctx, meta.GlobalKey(addrSpec.Name))
		if err != nil {
			return nil, err
		}
	}

	s.scope.Network().APIServerAddress = ptr.To[string](addr.SelfLink)
	endpoint := s.scope.ControlPlaneEndpoint()
	endpoint.Host = addr.Address
	s.scope.SetControlPlaneEndpoint(endpoint)
	return addr, nil
}

func (s *Service) createForwardingRule(ctx context.Context, target *compute.TargetTcpProxy, addr *compute.Address) error {
	log := log.FromContext(ctx)
	spec := s.scope.ForwardingRuleSpec()
	key := meta.GlobalKey(spec.Name)
	spec.IPAddress = addr.SelfLink
	spec.Target = target.SelfLink
	log.V(2).Info("Looking for forwardingrule", "name", spec.Name)
	forwarding, err := s.forwardingrules.Get(ctx, key)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for forwardingrule", "name", spec.Name)
			return err
		}

		log.V(2).Info("Creating a forwardingrule", "name", spec.Name)
		if err := s.forwardingrules.Insert(ctx, key, spec); err != nil {
			log.Error(err, "Error creating a forwardingrule", "name", spec.Name)
			return err
		}

		forwarding, err = s.forwardingrules.Get(ctx, key)
		if err != nil {
			return err
		}
	}

	s.scope.Network().APIServerForwardingRule = ptr.To[string](forwarding.SelfLink)
	return nil
}

func (s *Service) deleteForwardingRule(ctx context.Context) error {
	log := log.FromContext(ctx)
	spec := s.scope.ForwardingRuleSpec()
	key := meta.GlobalKey(spec.Name)
	log.V(2).Info("Deleting a forwardingrule", "name", spec.Name)
	if err := s.forwardingrules.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
		log.Error(err, "Error updating a forwardingrule", "name", spec.Name)
		return err
	}

	s.scope.Network().APIServerForwardingRule = nil
	return nil
}

func (s *Service) deleteAddress(ctx context.Context) error {
	log := log.FromContext(ctx)
	spec := s.scope.AddressSpec()
	key := meta.GlobalKey(spec.Name)
	log.V(2).Info("Deleting a address", "name", spec.Name)
	if err := s.addresses.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
		return err
	}

	s.scope.Network().APIServerAddress = nil
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

	s.scope.Network().APIServerTargetProxy = nil
	return nil
}

func (s *Service) deleteBackendService(ctx context.Context) error {
	log := log.FromContext(ctx)
	spec := s.scope.BackendServiceSpec()
	key := meta.GlobalKey(spec.Name)
	log.V(2).Info("Deleting a backendservice", "name", spec.Name)
	if err := s.backendservices.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
		log.Error(err, "Error deleting a backendservice", "name", spec.Name)
		return err
	}

	s.scope.Network().APIServerBackendService = nil
	return nil
}

func (s *Service) deleteHealthCheck(ctx context.Context) error {
	log := log.FromContext(ctx)
	spec := s.scope.HealthCheckSpec()
	key := meta.GlobalKey(spec.Name)
	log.V(2).Info("Deleting a healthcheck", "name", spec.Name)
	if err := s.healthchecks.Delete(ctx, key); err != nil && !gcperrors.IsNotFound(err) {
		log.Error(err, "Error deleting a healthcheck", "name", spec.Name)
		return err
	}

	s.scope.Network().APIServerHealthCheck = nil
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
