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

package instances

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/filter"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/compute/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile reconcile machine instance.
func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Reconciling instance resources")
	instance, err := s.createOrGetInstance(ctx)
	if err != nil {
		return err
	}

	addresses := make([]corev1.NodeAddress, 0, len(instance.NetworkInterfaces))
	for _, iface := range instance.NetworkInterfaces {
		addresses = append(addresses, corev1.NodeAddress{
			Type:    corev1.NodeInternalIP,
			Address: iface.NetworkIP,
		})

		for _, ac := range iface.AccessConfigs {
			addresses = append(addresses, corev1.NodeAddress{
				Type:    corev1.NodeExternalIP,
				Address: ac.NatIP,
			})
		}
	}

	machineName := s.scope.Name()
	zone := s.scope.Zone()
	project := s.scope.Project()

	// Since we don't know when the project was created, we must account for
	// both types of internal-dns:
	// https://cloud.google.com/compute/docs/internal-dns#instance-fully-qualified-domain-names
	// [INSTANCE_NAME].[ZONE].c.[PROJECT_ID].internal (newer)
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalDNS,
		Address: fmt.Sprintf("%s.%s.c.%s.internal", machineName, zone, project),
	})
	// [INSTANCE_NAME].c.[PROJECT_ID].internal
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalDNS,
		Address: fmt.Sprintf("%s.c.%s.internal", machineName, project),
	})
	// Add the machine's name as a known NodeInternalDNS because GCP platform
	// provides search paths to resolve those.
	// https://cloud.google.com/compute/docs/internal-dns#resolv.conf
	addresses = append(addresses, corev1.NodeAddress{
		Type:    corev1.NodeInternalDNS,
		Address: machineName,
	})

	s.scope.SetProviderID()
	s.scope.SetAddresses(addresses)
	s.scope.SetInstanceStatus(infrav1.InstanceStatus(instance.Status))

	if s.scope.IsControlPlane() {
		if err := s.registerControlPlaneInstance(ctx, instance); err != nil {
			return err
		}
	}

	return nil
}

// Delete delete machine instance.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	log.Info("Deleting instance resources")
	instanceSpec := s.scope.InstanceSpec(log)
	instanceName := instanceSpec.Name
	instanceKey := meta.ZonalKey(instanceName, s.scope.Zone())
	log.V(2).Info("Looking for instance before deleting", "name", instanceName, "zone", s.scope.Zone())
	instance, err := s.instances.Get(ctx, instanceKey)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for instance before deleting", "name", instanceName)
			return err
		}

		return nil
	}

	if s.scope.IsControlPlane() {
		if err := s.deregisterControlPlaneInstance(ctx, instance); err != nil {
			return err
		}
	}

	log.V(2).Info("Deleting instance", "name", instanceName, "zone", s.scope.Zone())
	return gcperrors.IgnoreNotFound(s.instances.Delete(ctx, instanceKey))
}

func (s *Service) createOrGetInstance(ctx context.Context) (*compute.Instance, error) {
	log := log.FromContext(ctx)
	log.V(2).Info("Getting bootstrap data for machine")
	bootstrapData, err := s.scope.GetBootstrapData()
	if err != nil {
		log.Error(err, "Error getting bootstrap data for machine")
		return nil, errors.Wrap(err, "failed to retrieve bootstrap data")
	}

	instanceSpec := s.scope.InstanceSpec(log)
	instanceName := instanceSpec.Name
	instanceKey := meta.ZonalKey(instanceName, s.scope.Zone())
	instanceSpec.Metadata.Items = append(instanceSpec.Metadata.Items, &compute.MetadataItems{
		Key:   "user-data",
		Value: ptr.To[string](bootstrapData),
	})

	log.V(2).Info("Looking for instance", "name", instanceName, "zone", s.scope.Zone())
	instance, err := s.instances.Get(ctx, instanceKey)
	if err != nil {
		if !gcperrors.IsNotFound(err) {
			log.Error(err, "Error looking for instance", "name", instanceName, "zone", s.scope.Zone())
			return nil, err
		}

		log.V(2).Info("Creating an instance", "name", instanceName, "zone", s.scope.Zone())
		if err := s.instances.Insert(ctx, instanceKey, instanceSpec); err != nil {
			log.Error(err, "Error creating an instance", "name", instanceName, "zone", s.scope.Zone())
			return nil, err
		}

		instance, err = s.instances.Get(ctx, instanceKey)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

func (s *Service) registerControlPlaneInstance(ctx context.Context, instance *compute.Instance) error {
	log := log.FromContext(ctx)
	instancegroupName := s.scope.ControlPlaneGroupName()
	log.V(2).Info("Ensuring instance already registered in the instancegroup", "name", instance.Name, "instancegroup", instancegroupName)
	instancegroupKey := meta.ZonalKey(instancegroupName, s.scope.Zone())
	instanceList, err := s.instancegroups.ListInstances(ctx, instancegroupKey, &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "RUNNING",
	}, filter.None)
	if err != nil {
		log.Error(err, "Error retrieving list of instances in the instancegroup", "instancegroup", instancegroupName)
		return err
	}

	instanceSets := sets.NewString()
	defer instanceSets.Delete()
	for _, i := range instanceList {
		instanceSets.Insert(i.Instance)
	}

	if !instanceSets.Has(instance.SelfLink) && instance.Status == string(infrav1.InstanceStatusRunning) {
		log.V(2).Info("Registering instance in the instancegroup", "name", instance.Name, "instancegroup", instancegroupName)
		if err := s.instancegroups.AddInstances(ctx, instancegroupKey, &compute.InstanceGroupsAddInstancesRequest{
			Instances: []*compute.InstanceReference{
				{
					Instance: instance.SelfLink,
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) deregisterControlPlaneInstance(ctx context.Context, instance *compute.Instance) error {
	log := log.FromContext(ctx)
	instancegroupName := s.scope.ControlPlaneGroupName()
	log.V(2).Info("Ensuring instance already registered in the instancegroup", "name", instance.Name, "instancegroup", instancegroupName)
	instancegroupKey := meta.ZonalKey(instancegroupName, s.scope.Zone())
	instanceList, err := s.instancegroups.ListInstances(ctx, instancegroupKey, &compute.InstanceGroupsListInstancesRequest{
		InstanceState: "RUNNING",
	}, filter.None)
	if err != nil {
		return gcperrors.IgnoreNotFound(err)
	}

	instanceSets := sets.NewString()
	defer instanceSets.Delete()
	for _, i := range instanceList {
		instanceSets.Insert(i.Instance)
	}

	if len(instanceSets.List()) > 0 && instanceSets.Has(instance.SelfLink) {
		log.V(2).Info("Deregistering instance in the instancegroup", "name", instance.Name, "instancegroup", instancegroupName)
		if err := s.instancegroups.RemoveInstances(ctx, instancegroupKey, &compute.InstanceGroupsRemoveInstancesRequest{
			Instances: []*compute.InstanceReference{
				{
					Instance: instance.SelfLink,
				},
			},
		}); err != nil {
			return gcperrors.IgnoreNotFound(err)
		}
	}

	return nil
}
