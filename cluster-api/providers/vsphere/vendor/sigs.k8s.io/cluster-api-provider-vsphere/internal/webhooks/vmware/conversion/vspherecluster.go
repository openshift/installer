/*
Copyright 2026 The Kubernetes Authors.

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

package conversion

import (
	"context"
	"reflect"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	vmwarev1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
)

// VSphereCluster is a HubSpokeConverter for the VSphereCluster API type.
var VSphereCluster = conversion.NewHubSpokeConverter(&vmwarev1.VSphereCluster{},
	conversion.NewSpokeConverter(&vmwarev1beta1.VSphereCluster{}, ConvertVSphereClusterHubToV1Beta1, ConvertVSphereClusterV1Beta1ToHub),
)

// ConvertVSphereClusterV1Beta1ToHub converts a v1beta1 VSphereCluster to a hub VSphereCluster.
func ConvertVSphereClusterV1Beta1ToHub(_ context.Context, src *vmwarev1beta1.VSphereCluster, dst *vmwarev1.VSphereCluster) error {
	if err := vmwarev1beta1.Convert_v1beta1_VSphereCluster_To_v1beta2_VSphereCluster(src, dst, nil); err != nil {
		return err
	}

	restored := &vmwarev1.VSphereCluster{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}

	initialization := vmwarev1.VSphereClusterInitializationStatus{}
	clusterv1.Convert_bool_To_Pointer_bool(src.Status.Ready, ok, restored.Status.Initialization.Provisioned, &initialization.Provisioned)
	if !reflect.DeepEqual(initialization, vmwarev1.VSphereClusterInitializationStatus{}) {
		dst.Status.Initialization = initialization
	}
	return nil
}

// ConvertVSphereClusterHubToV1Beta1 converts a hub VSphereCluster to a v1beta1 VSphereCluster.
func ConvertVSphereClusterHubToV1Beta1(_ context.Context, src *vmwarev1.VSphereCluster, dst *vmwarev1beta1.VSphereCluster) error {
	if err := vmwarev1beta1.Convert_v1beta2_VSphereCluster_To_v1beta1_VSphereCluster(src, dst, nil); err != nil {
		return err
	}

	return utilconversion.MarshalDataUnsafeNoCopy(src, dst)
}
