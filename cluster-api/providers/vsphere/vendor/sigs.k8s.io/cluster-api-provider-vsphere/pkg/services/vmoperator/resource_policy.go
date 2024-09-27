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

	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
)

// RPService represents the ability to reconcile a VirtualMachineSetResourcePolicy via vmoperator.
type RPService struct {
	Client client.Client
}

// ReconcileResourcePolicy ensures that a VirtualMachineSetResourcePolicy exists for the cluster
// Returns the name of a policy if it exists, otherwise returns an error.
func (s *RPService) ReconcileResourcePolicy(ctx context.Context, clusterCtx *vmware.ClusterContext) (string, error) {
	resourcePolicy, err := s.getVirtualMachineSetResourcePolicy(ctx, clusterCtx)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return "", errors.Errorf("unexpected error in getting the Resource policy: %+v", err)
		}
		resourcePolicy, err = s.createVirtualMachineSetResourcePolicy(ctx, clusterCtx)
		if err != nil {
			return "", errors.Errorf("failed to create Resource Policy: %+v", err)
		}
	}

	return resourcePolicy.Name, nil
}

func (s *RPService) newVirtualMachineSetResourcePolicy(clusterCtx *vmware.ClusterContext) *vmoprv1.VirtualMachineSetResourcePolicy {
	return &vmoprv1.VirtualMachineSetResourcePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: clusterCtx.Cluster.Namespace,
			Name:      clusterCtx.Cluster.Name,
		},
	}
}

func (s *RPService) getVirtualMachineSetResourcePolicy(ctx context.Context, clusterCtx *vmware.ClusterContext) (*vmoprv1.VirtualMachineSetResourcePolicy, error) {
	vmResourcePolicy := &vmoprv1.VirtualMachineSetResourcePolicy{}
	vmResourcePolicyName := client.ObjectKey{
		Namespace: clusterCtx.Cluster.Namespace,
		Name:      clusterCtx.Cluster.Name,
	}
	err := s.Client.Get(ctx, vmResourcePolicyName, vmResourcePolicy)
	return vmResourcePolicy, err
}

func (s *RPService) createVirtualMachineSetResourcePolicy(ctx context.Context, clusterCtx *vmware.ClusterContext) (*vmoprv1.VirtualMachineSetResourcePolicy, error) {
	vmResourcePolicy := s.newVirtualMachineSetResourcePolicy(clusterCtx)

	_, err := ctrlutil.CreateOrPatch(ctx, s.Client, vmResourcePolicy, func() error {
		vmResourcePolicy.Spec = vmoprv1.VirtualMachineSetResourcePolicySpec{
			ResourcePool: vmoprv1.ResourcePoolSpec{
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
			return errors.Wrapf(
				err,
				"error setting %s/%s as owner of %s/%s",
				clusterCtx.VSphereCluster.Namespace,
				clusterCtx.VSphereCluster.Name,
				vmResourcePolicy.Namespace,
				vmResourcePolicy.Name,
			)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vmResourcePolicy, nil
}
