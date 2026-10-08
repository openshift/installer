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

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
)

// VSphereClusterIdentity is a HubSpokeConverter for the VSphereClusterIdentity API type.
var VSphereClusterIdentity = conversion.NewHubSpokeConverter(&infrav1.VSphereClusterIdentity{},
	conversion.NewSpokeConverter(&infrav1beta1.VSphereClusterIdentity{}, ConvertVSphereClusterIdentityHubToV1Beta1, ConvertVSphereClusterIdentityV1Beta1ToHub),
)

// ConvertVSphereClusterIdentityV1Beta1ToHub converts a v1beta1 VSphereClusterIdentity to a hub VSphereClusterIdentity.
func ConvertVSphereClusterIdentityV1Beta1ToHub(_ context.Context, src *infrav1beta1.VSphereClusterIdentity, dst *infrav1.VSphereClusterIdentity) error {
	if err := infrav1beta1.Convert_v1beta1_VSphereClusterIdentity_To_v1beta2_VSphereClusterIdentity(src, dst, nil); err != nil {
		return err
	}

	restored := &infrav1.VSphereClusterIdentity{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}

	clusterv1.Convert_bool_To_Pointer_bool(src.Status.Ready, ok, restored.Status.Ready, &dst.Status.Ready)
	return nil
}

// ConvertVSphereClusterIdentityHubToV1Beta1 converts a hub VSphereClusterIdentity to a v1beta1 VSphereClusterIdentity.
func ConvertVSphereClusterIdentityHubToV1Beta1(_ context.Context, src *infrav1.VSphereClusterIdentity, dst *infrav1beta1.VSphereClusterIdentity) error {
	if err := infrav1beta1.Convert_v1beta2_VSphereClusterIdentity_To_v1beta1_VSphereClusterIdentity(src, dst, nil); err != nil {
		return err
	}

	return utilconversion.MarshalDataUnsafeNoCopy(src, dst)
}
