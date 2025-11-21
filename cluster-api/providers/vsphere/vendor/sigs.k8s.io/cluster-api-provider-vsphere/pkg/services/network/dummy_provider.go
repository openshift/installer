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

package network

import (
	"context"

	vmoprv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context/vmware"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services"
)

// dummyNetworkProvider doesn't provision network resource.
type dummyNetworkProvider struct{}

// DummyNetworkProvider returns an instance of dummy network provider.
func DummyNetworkProvider() services.NetworkProvider {
	return &dummyNetworkProvider{}
}

func (np *dummyNetworkProvider) HasLoadBalancer() bool {
	return false
}

func (np *dummyNetworkProvider) SupportsVMReadinessProbe() bool {
	return true
}

func (np *dummyNetworkProvider) ProvisionClusterNetwork(_ context.Context, clusterCtx *vmware.ClusterContext) error {
	v1beta2conditions.Set(clusterCtx.VSphereCluster, metav1.Condition{
		Type:   vmwarev1.VSphereClusterNetworkReadyV1Beta2Condition,
		Status: metav1.ConditionTrue,
		Reason: vmwarev1.VSphereClusterNetworkReadyV1Beta2Reason,
	})
	return nil
}

func (np *dummyNetworkProvider) GetClusterNetworkName(_ context.Context, _ *vmware.ClusterContext) (string, error) {
	return "", nil
}

func (np *dummyNetworkProvider) ConfigureVirtualMachine(_ context.Context, _ *vmware.ClusterContext, _ *vmwarev1.VSphereMachine, _ *vmoprv1.VirtualMachine) error {
	return nil
}

func (np *dummyNetworkProvider) GetVMServiceAnnotations(_ context.Context, _ *vmware.ClusterContext) (map[string]string, error) {
	return map[string]string{}, nil
}

func (np *dummyNetworkProvider) VerifyNetworkStatus(_ context.Context, _ *vmware.ClusterContext, _ runtime.Object) error {
	return nil
}

type dummyLBNetworkProvider struct {
	dummyNetworkProvider
}

// DummyLBNetworkProvider returns an instance of dummy network provider that has a LB.
func DummyLBNetworkProvider() services.NetworkProvider {
	return &dummyLBNetworkProvider{}
}

func (np *dummyLBNetworkProvider) HasLoadBalancer() bool {
	return true
}
