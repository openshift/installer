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

package util

import (
	"context"

	"github.com/pkg/errors"
	apitypes "k8s.io/apimachinery/pkg/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
)

// GetVSphereClusterFromVMwareMachine gets the vmware.infrastructure.cluster.x-k8s.io.VSphereCluster resource for the given VSphereMachine.
// TODO (srm09): Rename this to a more appropriate name.
func GetVSphereClusterFromVMwareMachine(ctx context.Context, c client.Client, machine *vmwarev1.VSphereMachine) (*vmwarev1.VSphereCluster, error) {
	clusterName := machine.Labels[clusterv1.ClusterNameLabel]
	if clusterName == "" {
		return nil, errors.Errorf("error getting VSphereCluster name from VSphereMachine %s/%s",
			machine.Namespace, machine.Name)
	}
	namespacedName := apitypes.NamespacedName{
		Namespace: machine.Namespace,
		Name:      clusterName,
	}
	cluster := &clusterv1.Cluster{}
	if err := c.Get(ctx, namespacedName, cluster); err != nil {
		return nil, err
	}

	if !cluster.Spec.InfrastructureRef.IsDefined() {
		return nil, errors.Errorf("error getting VSphereCluster name from VSphereMachine %s/%s: Cluster.spec.infrastructureRef not yet set",
			machine.Namespace, machine.Name)
	}

	vsphereClusterKey := apitypes.NamespacedName{
		Namespace: machine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	vsphereCluster := &vmwarev1.VSphereCluster{}
	err := c.Get(ctx, vsphereClusterKey, vsphereCluster)
	return vsphereCluster, err
}

// GetVSphereClusterFromVSphereMachine gets the infrastructure.cluster.x-k8s.io.VSphereCluster resource for the given VSphereMachine.
func GetVSphereClusterFromVSphereMachine(ctx context.Context, c client.Client, machine *infrav1.VSphereMachine) (*infrav1.VSphereCluster, error) {
	clusterName := machine.Labels[clusterv1.ClusterNameLabel]
	if clusterName == "" {
		return nil, errors.Errorf("error getting VSphereCluster name from VSphereMachine %s/%s",
			machine.Namespace, machine.Name)
	}
	namespacedName := apitypes.NamespacedName{
		Namespace: machine.Namespace,
		Name:      clusterName,
	}
	cluster := &clusterv1.Cluster{}
	if err := c.Get(ctx, namespacedName, cluster); err != nil {
		return nil, err
	}

	if !cluster.Spec.InfrastructureRef.IsDefined() {
		return nil, errors.Errorf("error getting VSphereCluster name from VSphereMachine %s/%s: Cluster.spec.infrastructureRef not yet set",
			machine.Namespace, machine.Name)
	}
	vsphereClusterKey := apitypes.NamespacedName{
		Namespace: machine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	vsphereCluster := &infrav1.VSphereCluster{}
	err := c.Get(ctx, vsphereClusterKey, vsphereCluster)
	return vsphereCluster, err
}
