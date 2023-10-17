/*
Copyright 2022 The Kubernetes Authors.

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

package subnets

import (
	"context"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile reconcile cluster network components.
func (s *Service) Reconcile(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling subnetwork resources")

	// reconcile subnets
	if _, err := s.createOrGetSubnets(ctx); err != nil {
		return err
	}

	return nil
}

// Delete deletes cluster subnetwork components.
func (s *Service) Delete(ctx context.Context) error {
	logger := log.FromContext(ctx)
	for _, subnetSpec := range s.scope.SubnetSpecs() {
		logger.V(2).Info("Deleting a subnet", "name", subnetSpec.Name)
		subnetKey := meta.RegionalKey(subnetSpec.Name, s.scope.Region())
		err := s.subnets.Delete(ctx, subnetKey)
		if err != nil && !gcperrors.IsNotFound(err) {
			logger.Error(err, "Error deleting subnet", "name", subnetSpec.Name)
			return err
		}
	}

	return nil
}

// createOrGetSubnets creates the subnetworks if they don't exist otherwise return the existing ones.
func (s *Service) createOrGetSubnets(ctx context.Context) ([]*compute.Subnetwork, error) {
	logger := log.FromContext(ctx)
	subnets := []*compute.Subnetwork{}
	for _, subnetSpec := range s.scope.SubnetSpecs() {
		logger.V(2).Info("Looking for subnet", "name", subnetSpec.Name)
		subnetKey := meta.RegionalKey(subnetSpec.Name, s.scope.Region())
		subnet, err := s.subnets.Get(ctx, subnetKey)
		if err != nil {
			if !gcperrors.IsNotFound(err) {
				logger.Error(err, "Error looking for subnet", "name", subnetSpec.Name)
				return subnets, err
			}

			// Subnet was not found, let's create it
			logger.V(2).Info("Creating a subnet", "name", subnetSpec.Name)
			if err := s.subnets.Insert(ctx, subnetKey, subnetSpec); err != nil {
				logger.Error(err, "Error creating a subnet", "name", subnetSpec.Name)
				return subnets, err
			}

			subnet, err = s.subnets.Get(ctx, subnetKey)
			if err != nil {
				logger.Error(err, "Error getting existing subnet", "name", subnetSpec.Name)
				return subnets, err
			}
		}
		subnets = append(subnets, subnet)
	}

	return subnets, nil
}
