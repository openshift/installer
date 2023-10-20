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
	"github.com/pkg/errors"
	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	vmwareutil "sigs.k8s.io/cluster-api-provider-vsphere/pkg/util/vmware"
)

// RPService represents the ability to reconcile a VirtualMachineSetResourcePolicy via vmoperator.
type RPService struct{}

// ReconcileResourcePolicy ensures that a VirtualMachineSetResourcePolicy exists for the cluster
// Returns the name of a policy if it exists, otherwise returns an error.
func (s RPService) ReconcileResourcePolicy(ctx *vmware.ClusterContext) (string, error) {
	resourcePolicy, err := s.getVirtualMachineSetResourcePolicy(ctx)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return "", errors.Errorf("unexpected error in getting the Resource policy: %+v", err)
		}
		resourcePolicy, err = s.createVirtualMachineSetResourcePolicy(ctx)
		if err != nil {
			return "", errors.Errorf("failed to create Resource Policy: %+v", err)
		}
	}

	return resourcePolicy.Name, nil
}

func (s RPService) newVirtualMachineSetResourcePolicy(ctx *vmware.ClusterContext) *vmoprv1.VirtualMachineSetResourcePolicy {
	return &vmoprv1.VirtualMachineSetResourcePolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ctx.Cluster.Namespace,
			Name:      ctx.Cluster.Name,
		},
	}
}

func (s RPService) getVirtualMachineSetResourcePolicy(ctx *vmware.ClusterContext) (*vmoprv1.VirtualMachineSetResourcePolicy, error) {
	vmResourcePolicy := &vmoprv1.VirtualMachineSetResourcePolicy{}
	vmResourcePolicyName := client.ObjectKey{
		Namespace: ctx.Cluster.Namespace,
		Name:      ctx.Cluster.Name,
	}
	err := ctx.Client.Get(ctx, vmResourcePolicyName, vmResourcePolicy)
	return vmResourcePolicy, err
}

func (s RPService) createVirtualMachineSetResourcePolicy(ctx *vmware.ClusterContext) (*vmoprv1.VirtualMachineSetResourcePolicy, error) {
	vmResourcePolicy := s.newVirtualMachineSetResourcePolicy(ctx)

	_, err := ctrlutil.CreateOrPatch(ctx, ctx.Client, vmResourcePolicy, func() error {
		vmResourcePolicy.Spec = vmoprv1.VirtualMachineSetResourcePolicySpec{
			ResourcePool: vmoprv1.ResourcePoolSpec{
				Name: ctx.Cluster.Name,
			},
			Folder: vmoprv1.FolderSpec{
				Name: ctx.Cluster.Name,
			},
			ClusterModules: []vmoprv1.ClusterModuleSpec{
				{
					GroupName: ControlPlaneVMClusterModuleGroupName,
				},
				{
					GroupName: vmwareutil.GetMachineDeploymentNameForCluster(ctx.Cluster),
				},
			},
		}
		// Ensure that the VirtualMachineSetResourcePolicy is owned by the VSphereCluster
		if err := ctrlutil.SetOwnerReference(
			ctx.VSphereCluster,
			vmResourcePolicy,
			ctx.Scheme,
		); err != nil {
			return errors.Wrapf(
				err,
				"error setting %s/%s as owner of %s/%s",
				ctx.VSphereCluster.Namespace,
				ctx.VSphereCluster.Name,
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
