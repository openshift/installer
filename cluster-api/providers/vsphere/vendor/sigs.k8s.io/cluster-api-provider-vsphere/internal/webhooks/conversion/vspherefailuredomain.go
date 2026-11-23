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

	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
)

// VSphereFailureDomain is a HubSpokeConverter for the VSphereFailureDomain API type.
var VSphereFailureDomain = conversion.NewHubSpokeConverter(&infrav1.VSphereFailureDomain{},
	conversion.NewSpokeConverter(&infrav1beta1.VSphereFailureDomain{}, ConvertVSphereFailureDomainHubToV1Beta1, ConvertVSphereFailureDomainV1Beta1ToHub),
)

// ConvertVSphereFailureDomainV1Beta1ToHub converts a v1beta1 VSphereFailureDomain to a hub VSphereFailureDomain.
func ConvertVSphereFailureDomainV1Beta1ToHub(_ context.Context, src *infrav1beta1.VSphereFailureDomain, dst *infrav1.VSphereFailureDomain) error {
	return infrav1beta1.Convert_v1beta1_VSphereFailureDomain_To_v1beta2_VSphereFailureDomain(src, dst, nil)
}

// ConvertVSphereFailureDomainHubToV1Beta1 converts a hub VSphereFailureDomain to a v1beta1 VSphereFailureDomain.
func ConvertVSphereFailureDomainHubToV1Beta1(_ context.Context, src *infrav1.VSphereFailureDomain, dst *infrav1beta1.VSphereFailureDomain) error {
	if err := infrav1beta1.Convert_v1beta2_VSphereFailureDomain_To_v1beta1_VSphereFailureDomain(src, dst, nil); err != nil {
		return err
	}

	if dst.Spec.Topology.ComputeCluster != nil && *dst.Spec.Topology.ComputeCluster == "" {
		dst.Spec.Topology.ComputeCluster = nil
	}
	return nil
}
