/*
Copyright 2025 The Kubernetes Authors.

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

package instancegroupmanagers

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/cluster-api-provider-gcp/pkg/gcp"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile reconciles GCP instanceGroupManager resources.
func (s *Service) Reconcile(ctx context.Context, instanceTemplateKey *meta.Key) (*compute.InstanceGroupManager, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling instanceGroupManager resources")
	igm, err := s.createOrGet(ctx, instanceTemplateKey)
	if err != nil {
		return nil, err
	}
	log.V(2).Info("Reconciled instanceGroupManager", "selfLink", igm.SelfLink)

	return igm, nil
}

// Delete deletes the GCP instanceGroupManager resource.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)

	igmKey, err := s.scope.InstanceGroupManagerResourceName()
	if err != nil {
		return err
	}

	selfLink := gcp.FormatKey("instanceGroupManagers", igmKey)

	log = log.WithValues("instanceGroupManager", selfLink)
	log.Info("Deleting instanceGroupManager resources")

	log.V(2).Info("Looking for instanceGroupManager before deleting")
	existing, err := s.instanceGroupManagers.Get(ctx, igmKey)
	if err != nil {
		if gcperrors.IsNotFound(err) {
			existing = nil
		} else {
			log.Error(err, "Error looking for instanceGroupManager before deleting")
			return fmt.Errorf("getting instanceGroupManager: %w", err)
		}
	}

	if existing != nil {
		// There are no labels on the instanceGroupManager itself, so we cannot filter by
		// ownership. Instead, we rely on the fact that the instanceGroupManager name is
		// derived from the cluster and machine pool name, so if it exists, we assume
		// we own it.

		log.V(2).Info("Deleting instanceGroupManager", "selfLink", existing.SelfLink)
		if err := s.instanceGroupManagers.Delete(ctx, igmKey); err != nil {
			if gcperrors.IsNotFound(err) {
				log.V(2).Info("instanceGroupManager not found, assuming already deleted", "selfLink", existing.SelfLink)
			} else {
				log.Error(err, "Error deleting instanceGroupManager", "selfLink", existing.SelfLink)
				return fmt.Errorf("deleting instanceGroupManager: %w", err)
			}
		}
	}

	return nil
}

func (s *Service) createOrGet(ctx context.Context, instanceTemplateKey *meta.Key) (*compute.InstanceGroupManager, error) {
	log := log.FromContext(ctx)

	igmKey, err := s.scope.InstanceGroupManagerResourceName()
	if err != nil {
		return nil, err
	}

	selfLink := gcp.FormatKey("instanceGroupManagers", igmKey)

	log = log.WithValues("instanceGroupManager", selfLink)
	log.Info("Getting instanceGroupManager resources")

	desired, err := s.scope.InstanceGroupManagerResource(instanceTemplateKey)
	if err != nil {
		return nil, err
	}

	log.V(2).Info("Looking for instanceGroupManager")
	actual, err := s.instanceGroupManagers.Get(ctx, igmKey)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for instanceGroupManager")
			return nil, fmt.Errorf("getting instanceGroupManager %v: %w", selfLink, err)
		}

		log.V(2).Info("Creating instanceGroupManager")
		if err := s.instanceGroupManagers.Insert(ctx, igmKey, desired); err != nil {
			log.Error(err, "creating instanceGroupManager")
			return nil, fmt.Errorf("creating instanceGroupManager %v: %w", selfLink, err)
		}

		// We fetch the created instanceGroupManager to return the complete object.
		// Also, we continue to go through the comparisons below so we have one code path,
		// (there are no API calls so this is cheap),
		// although in practice we don't expect to need further updates.
		actual, err = s.instanceGroupManagers.Get(ctx, igmKey)
		if err != nil {
			return nil, fmt.Errorf("getting instanceGroupManager %v: %w", selfLink, err)
		}
	}

	if desired.TargetSize != actual.TargetSize {
		log.V(2).Info("resizing instanceGroupManager", "targetSize", desired.TargetSize)
		if err := s.instanceGroupManagers.Resize(ctx, igmKey, desired.TargetSize); err != nil {
			log.Error(err, "resizing instanceGroupManager")
			return nil, fmt.Errorf("resizing instanceGroupManager %v: %w", selfLink, err)
		}

		actual.TargetSize = desired.TargetSize
	}

	if desired.InstanceTemplate != actual.InstanceTemplate {
		log.V(2).Info("updating instanceTemplate for instanceGroupManager", "desired.instanceTemplate", desired.InstanceTemplate, "actual.instanceTemplate", actual.InstanceTemplate)
		if err := s.instanceGroupManagers.SetInstanceTemplate(ctx, igmKey, &compute.InstanceGroupManagersSetInstanceTemplateRequest{
			InstanceTemplate: desired.InstanceTemplate,
		}); err != nil {
			log.Error(err, "updating instanceTemplate for instanceGroupManager")
			return nil, fmt.Errorf("updating instanceTemplate for instanceGroupManager %v: %w", selfLink, err)
		}

		actual.InstanceTemplate = desired.InstanceTemplate
	}

	return actual, nil
}

// ListInstances lists instances in the the instanceGroup linked to the passed instanceGroupManager.
func (s *Service) ListInstances(ctx context.Context, instanceGroupManager *compute.InstanceGroupManager) ([]*compute.InstanceWithNamedPorts, error) {
	log := log.FromContext(ctx)

	var igKey *meta.Key
	{
		instanceGroup := instanceGroupManager.InstanceGroup
		instanceGroup = strings.TrimPrefix(instanceGroup, "https://www.googleapis.com/")
		instanceGroup = strings.TrimPrefix(instanceGroup, "compute/v1/")
		tokens := strings.Split(instanceGroup, "/")
		if len(tokens) == 6 && tokens[0] == "projects" && tokens[2] == "zones" && tokens[4] == "instanceGroups" {
			igKey = meta.ZonalKey(tokens[5], tokens[3])
		} else {
			return nil, fmt.Errorf("unexpected format for instanceGroup: %q", instanceGroup)
		}
	}

	log.Info("Listing instances in instanceGroup", "instanceGroup", instanceGroupManager.InstanceGroup)
	listInstancesRequest := &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "ALL",
	}
	instances, err := s.instanceGroups.ListInstances(ctx, igKey, listInstancesRequest, filter.None)
	if err != nil {
		log.Error(err, "Error listing instances in instanceGroup", "instanceGroup", instanceGroupManager.InstanceGroup)
		return nil, fmt.Errorf("listing instances in instanceGroup %q: %w", instanceGroupManager.InstanceGroup, err)
	}

	return instances, nil
}
