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

package clustermodule

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/clustermodules"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/session"
)

const validMachineTemplate = "VSphereMachineTemplate"

type service struct {
	ControllerManagerContext *capvcontext.ControllerManagerContext
	Client                   client.Client
}

// NewService returns a new Cluster Module service.
func NewService(controllerManagerCtx *capvcontext.ControllerManagerContext, client client.Client) Service {
	return &service{
		ControllerManagerContext: controllerManagerCtx,
		Client:                   client,
	}
}

func (s *service) Create(ctx context.Context, clusterCtx *capvcontext.ClusterContext, wrapper Wrapper) (string, error) {
	log := ctrl.LoggerFrom(ctx)

	templateRef, err := s.fetchTemplateRef(ctx, wrapper)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching template ref for object %s/%s", wrapper.GetNamespace(), wrapper.GetName())
	}
	if templateRef.Kind != validMachineTemplate {
		// since this is a heterogeneous cluster, we should skip cluster module creation for non VSphereMachine objects
		log.V(4).Info("Skipping module creation for non-VSphereMachine objects")
		return "", nil
	}

	template, err := s.fetchMachineTemplate(ctx, wrapper, templateRef.Name)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching machine template for object %s/%s", wrapper.GetNamespace(), wrapper.GetName())
	}
	if server := template.Spec.Template.Spec.Server; server != clusterCtx.VSphereCluster.Spec.Server {
		log.V(4).Info("Skipping module creation for object since template uses a different server", "server", server)
		return "", nil
	}

	vCenterSession, err := s.fetchSessionForObject(ctx, clusterCtx, template)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching session for object %s/%s", wrapper.GetNamespace(), wrapper.GetName())
	}

	// Fetch the compute cluster resource by tracing the owner of the resource pool in use.
	// TODO (srm09): How do we support Multi AZ scenarios here
	computeClusterRef, err := getComputeClusterResource(ctx, vCenterSession, template.Spec.Template.Spec.ResourcePool)
	if err != nil {
		return "", errors.Wrapf(err, "error fetching compute cluster resource")
	}

	provider := clustermodules.NewProvider(vCenterSession.TagManager.Client)
	moduleUUID, err := provider.CreateModule(ctx, computeClusterRef)
	if err != nil {
		return "", errors.Wrapf(err, "error creating cluster module")
	}
	return moduleUUID, nil
}

func (s *service) DoesExist(ctx context.Context, clusterCtx *capvcontext.ClusterContext, wrapper Wrapper, moduleUUID string) (bool, error) {
	templateRef, err := s.fetchTemplateRef(ctx, wrapper)
	if err != nil {
		return false, errors.Wrapf(err, "error fetching template ref for object %s/%s", wrapper.GetNamespace(), wrapper.GetName())
	}

	template, err := s.fetchMachineTemplate(ctx, wrapper, templateRef.Name)
	if err != nil {
		return false, errors.Wrapf(err, "error fetching machine template for object %s/%s", wrapper.GetNamespace(), wrapper.GetName())
	}

	vCenterSession, err := s.fetchSessionForObject(ctx, clusterCtx, template)
	if err != nil {
		return false, errors.Wrapf(err, "error fetching session for object %s/%s", wrapper.GetNamespace(), wrapper.GetName())
	}

	provider := clustermodules.NewProvider(vCenterSession.TagManager.Client)
	return provider.DoesModuleExist(ctx, moduleUUID)
}

func (s *service) Remove(ctx context.Context, clusterCtx *capvcontext.ClusterContext, moduleUUID string) error {
	params := s.newParams(*clusterCtx)
	vcenterSession, err := s.fetchSession(ctx, clusterCtx, params)
	if err != nil {
		return err
	}

	provider := clustermodules.NewProvider(vcenterSession.TagManager.Client)
	return provider.DeleteModule(ctx, moduleUUID)
}

func getComputeClusterResource(ctx context.Context, s *session.Session, resourcePool string) (types.ManagedObjectReference, error) {
	rp, err := s.Finder.ResourcePoolOrDefault(ctx, resourcePool)
	if err != nil {
		return types.ManagedObjectReference{}, err
	}

	cc, err := rp.Owner(ctx)
	if err != nil {
		return types.ManagedObjectReference{}, err
	}

	ownerPath, err := find.InventoryPath(ctx, s.Client.Client, cc.Reference())
	if err != nil {
		return types.ManagedObjectReference{}, err
	}
	if _, err = s.Finder.ClusterComputeResource(ctx, ownerPath); err != nil {
		return types.ManagedObjectReference{}, IncompatibleOwnerError{cc.Reference().Value}
	}

	return cc.Reference(), nil
}
