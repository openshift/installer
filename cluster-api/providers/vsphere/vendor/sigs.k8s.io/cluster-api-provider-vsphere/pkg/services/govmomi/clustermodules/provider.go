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

package clustermodules

import (
	"context"

	"github.com/vmware/govmomi/vapi/cluster"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/util"
)

var log = logf.Log.V(5).WithName("govmomi").WithName("clustermodule")

// Provider exposes methods to interact with the cluster module vCenter API
// TODO (srm09): Rethink and merge with ClusterModuleService.
type Provider interface {
	CreateModule(ctx context.Context, clusterRef types.ManagedObjectReference) (string, error)
	DeleteModule(ctx context.Context, moduleID string) error
	DoesModuleExist(ctx context.Context, moduleID string, cluster types.ManagedObjectReference) (bool, error)

	IsMoRefModuleMember(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) (bool, error)
	AddMoRefToModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error
	RemoveMoRefFromModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error
}

type provider struct {
	manager *cluster.Manager
}

func NewProvider(restClient *rest.Client) Provider {
	return &provider{
		manager: cluster.NewManager(restClient),
	}
}

func (cm *provider) CreateModule(ctx context.Context, clusterRef types.ManagedObjectReference) (string, error) {
	log.Info("Creating cluster module", "cluster", clusterRef)

	moduleID, err := cm.manager.CreateModule(ctx, clusterRef)
	if err != nil {
		return "", err
	}

	log.Info("Created cluster module", "moduleID", moduleID)
	return moduleID, nil
}

func (cm *provider) DeleteModule(ctx context.Context, moduleID string) error {
	log.Info("Deleting cluster module", "moduleID", moduleID)

	err := cm.manager.DeleteModule(ctx, moduleID)
	if err != nil && !util.IsNotFoundError(err) {
		return err
	}

	log.Info("Deleted cluster module", "moduleID", moduleID)
	return nil
}

func (cm *provider) DoesModuleExist(ctx context.Context, moduleID string, clusterRef types.ManagedObjectReference) (bool, error) {
	log.V(4).Info("Checking if cluster module exists", "moduleID", moduleID, "clusterRef", clusterRef)

	if moduleID == "" {
		return false, nil
	}

	modules, err := cm.manager.ListModules(ctx)
	if err != nil {
		return false, err
	}

	for _, mod := range modules {
		if mod.Cluster == clusterRef.Value && mod.Module == moduleID {
			return true, nil
		}
	}

	log.V(4).Info("Cluster module doesn't exist", "moduleID", moduleID, "clusterRef", clusterRef)
	return false, nil
}

func (cm *provider) IsMoRefModuleMember(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) (bool, error) {
	moduleMembers, err := cm.manager.ListModuleMembers(ctx, moduleID)
	if err != nil {
		return false, err
	}

	for _, member := range moduleMembers {
		if member.Reference() == moRef.Reference() {
			return true, nil
		}
	}

	return false, nil
}

func (cm *provider) AddMoRefToModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error {
	isMember, err := cm.IsMoRefModuleMember(ctx, moduleID, moRef)
	if err != nil {
		return err
	}

	if !isMember {
		log.Info("Adding moRef to cluster module", "moduleID", moduleID, "moRef", moRef)
		// TODO: Should we just skip the IsMoRefModuleMember() and always call this since we're already
		// ignoring the first return value?
		_, err := cm.manager.AddModuleMembers(ctx, moduleID, moRef.Reference())
		if err != nil {
			return err
		}
	}

	return nil
}

func (cm *provider) RemoveMoRefFromModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error {
	log.Info("Removing moRef from cluster module", "moduleID", moduleID, "moRef", moRef)

	_, err := cm.manager.RemoveModuleMembers(ctx, moduleID, moRef)
	if err != nil {
		return err
	}

	log.Info("Removed moRef from cluster module", "moduleID", moduleID, "moRef", moRef)
	return nil
}
