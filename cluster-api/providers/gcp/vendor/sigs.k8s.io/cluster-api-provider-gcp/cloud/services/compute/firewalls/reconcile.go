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

package firewalls

import (
	"context"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile reconcile cluster firewall compoenents.
func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	if s.scope.IsSharedVpc() {
		log.V(2).Info("Shared VPC enabled. Ignore Reconciling firewall resources")
		return nil
	}
	log.Info("Reconciling firewall resources")
	for _, spec := range s.scope.FirewallRulesSpec() {
		log.V(2).Info("Looking firewall", "name", spec.Name)
		firewallKey := meta.GlobalKey(spec.Name)
		if _, err := s.firewalls.Get(ctx, firewallKey); err != nil {
			if !gcperrors.IsNotFound(err) {
				return err
			}

			log.V(2).Info("Creating firewall", "name", spec.Name)
			if err := s.firewalls.Insert(ctx, firewallKey, spec); err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete delete cluster firewall compoenents.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	if s.scope.IsSharedVpc() {
		log.V(2).Info("Shared VPC enabled. Ignore Deleting firewall resources")
		return nil
	}
	log.Info("Deleting firewall resources")
	for _, spec := range s.scope.FirewallRulesSpec() {
		log.V(2).Info("Deleting firewall", "name", spec.Name)
		firewallKey := meta.GlobalKey(spec.Name)
		if err := s.firewalls.Delete(ctx, firewallKey); err != nil {
			if !gcperrors.IsNotFound(err) {
				log.Error(err, "Error deleting firewall", "name", spec.Name)
				return err
			}
		}
	}

	return nil
}
