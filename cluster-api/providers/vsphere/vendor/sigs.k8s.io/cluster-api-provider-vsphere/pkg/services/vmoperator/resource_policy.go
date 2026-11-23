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

package vmoperator

import (
	"context"
	"reflect"

	pkgerrors "github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
	conversionclient "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/client"
)

// RPService represents the ability to reconcile a VirtualMachineSetResourcePolicy via vmoperator.
type RPService struct {
	Client client.Client
}

// ReconcileResourcePolicy ensures that a VirtualMachineSetResourcePolicy exists for the cluster
// Returns the name of a policy if it exists, otherwise returns an error.
func (s *RPService) ReconcileResourcePolicy(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	if err := s.createOrPatchVirtualMachineSetResourcePolicy(ctx, clusterCtx); err != nil {
		return pkgerrors.Errorf("failed to create Resource Policy: %+v", err)
	}
	return nil
}

func (s *RPService) newVirtualMachineSetResourcePolicy(clusterCtx *vmware.ClusterContext) *vmoprvhub.VirtualMachineSetResourcePolicy {
	return &vmoprvhub.VirtualMachineSetResourcePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: clusterCtx.Cluster.Namespace,
			Name:      clusterCtx.Cluster.Name,
		},
	}
}

func (s *RPService) createOrPatchVirtualMachineSetResourcePolicy(ctx context.Context, clusterCtx *vmware.ClusterContext) error {
	vmResourcePolicy := s.newVirtualMachineSetResourcePolicy(clusterCtx)

	vmResourcePolicyExists := true
	if err := s.Client.Get(ctx, client.ObjectKeyFromObject(vmResourcePolicy), vmResourcePolicy); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
		vmResourcePolicyExists = false
	}
	originalResourcePolicy := vmResourcePolicy.DeepCopy()

	vmResourcePolicy.Spec = vmoprvhub.VirtualMachineSetResourcePolicySpec{
		ResourcePool: vmoprvhub.ResourcePoolSpec{
			Name: clusterCtx.Cluster.Name,
		},
		Folder: clusterCtx.Cluster.Name,
		ClusterModuleGroups: []string{
			ControlPlaneVMClusterModuleGroupName,
			getMachineDeploymentNameForCluster(clusterCtx.Cluster),
		},
	}
	// Ensure that the VirtualMachineSetResourcePolicy is owned by the VSphereCluster
	if err := ctrlutil.SetOwnerReference(
		clusterCtx.VSphereCluster,
		vmResourcePolicy,
		s.Client.Scheme(),
	); err != nil {
		return pkgerrors.Wrapf(
			err,
			"error setting %s/%s as owner of %s/%s",
			clusterCtx.VSphereCluster.Namespace,
			clusterCtx.VSphereCluster.Name,
			vmResourcePolicy.Namespace,
			vmResourcePolicy.Name,
		)
	}

	if !vmResourcePolicyExists {
		if err := s.Client.Create(ctx, vmResourcePolicy); err != nil {
			return err
		}
	} else if !reflect.DeepEqual(originalResourcePolicy, vmResourcePolicy) {
		patch, err := conversionclient.MergeFrom(ctx, s.Client, originalResourcePolicy)
		if err != nil {
			return pkgerrors.Wrapf(err, "failed to create patch for VirtualMachineSetResourcePolicy object")
		}
		if err := s.Client.Patch(ctx, vmResourcePolicy, patch); err != nil {
			return pkgerrors.Wrapf(err, "failed to patch VirtualMachineSetResourcePolicy object")
		}
	}

	return nil
}
