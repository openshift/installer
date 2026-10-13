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

// Package clustermodules contains tools for handling Cluster Modules.
package clustermodules

import (
	"context"
	"net/http"

	"github.com/vmware/govmomi/vapi/cluster"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Provider exposes methods to interact with the cluster module vCenter API
// TODO (srm09): Rethink and merge with ClusterModuleService.
type Provider interface {
	CreateModule(ctx context.Context, clusterRef types.ManagedObjectReference) (string, error)
	DeleteModule(ctx context.Context, moduleID string) error
	DoesModuleExist(ctx context.Context, moduleID string) (bool, error)

	IsMoRefModuleMember(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) (bool, error)
	AddMoRefToModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error
	RemoveMoRefFromModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error
}

type provider struct {
	manager *cluster.Manager
}

// NewProvider returns a new Cluster Module provider.
func NewProvider(restClient *rest.Client) Provider {
	return &provider{
		manager: cluster.NewManager(restClient),
	}
}

// CreateModule creates a new Cluster Module and returns its ID.
func (cm *provider) CreateModule(ctx context.Context, clusterRef types.ManagedObjectReference) (string, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Creating cluster module", "computeClusterRef", clusterRef)

	moduleUUID, err := cm.manager.CreateModule(ctx, clusterRef)
	if err != nil {
		return "", err
	}

	log.Info("Created cluster module", "computeClusterRef", clusterRef, "moduleUUID", moduleUUID)
	return moduleUUID, nil
}

// DeleteModule deletes a  Cluster Module by ID.
func (cm *provider) DeleteModule(ctx context.Context, moduleUUID string) error {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Deleting cluster module")

	err := cm.manager.DeleteModule(ctx, moduleUUID)
	if err != nil && !rest.IsStatusError(err, http.StatusNotFound) {
		return err
	}

	log.Info("Deleted cluster module")
	return nil
}

// DoesModuleExist checks whether a module with a given moduleUUID exists.
func (cm *provider) DoesModuleExist(ctx context.Context, moduleUUID string) (bool, error) {
	log := ctrl.LoggerFrom(ctx)
	log.V(4).Info("Checking if cluster module exists")

	if moduleUUID == "" {
		return false, nil
	}

	_, err := cm.manager.ListModuleMembers(ctx, moduleUUID)
	if err == nil {
		log.V(4).Info("Cluster module exists")
		return true, nil
	}

	if rest.IsStatusError(err, http.StatusNotFound) {
		log.V(4).Info("Cluster module doesn't exist")
		return false, nil
	}

	return false, err
}

// IsMoRefModuleMember checks whether the passed managed object reference is in the ClusterModule.
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

// AddMoRefToModule adds the object to the ClusterModule if it is not already a member.
func (cm *provider) AddMoRefToModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error {
	log := ctrl.LoggerFrom(ctx)
	isMember, err := cm.IsMoRefModuleMember(ctx, moduleID, moRef)
	if err != nil {
		return err
	}

	if !isMember {
		log.Info("Adding moRef to the cluster module", "moRef", moRef)
		// TODO: Should we just skip the IsMoRefModuleMember() and always call this since we're already
		// ignoring the first return value?
		_, err := cm.manager.AddModuleMembers(ctx, moduleID, moRef.Reference())
		if err != nil {
			return err
		}

		log.Info("Added moRef to the cluster module", "moRef", moRef)
	}

	return nil
}

// RemoveMoRefFromModule removes the object from the ClusterModule.
func (cm *provider) RemoveMoRefFromModule(ctx context.Context, moduleID string, moRef types.ManagedObjectReference) error {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Removing moRef from the cluster module", "moRef", moRef)

	_, err := cm.manager.RemoveModuleMembers(ctx, moduleID, moRef)
	if err != nil {
		return err
	}

	log.Info("Removed moRef from the cluster module", "moRef", moRef)
	return nil
}
