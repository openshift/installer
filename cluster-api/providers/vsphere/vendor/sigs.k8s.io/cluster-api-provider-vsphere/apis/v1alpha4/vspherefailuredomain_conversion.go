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

package v1alpha4

import (
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereFailureDomain to the Hub version (v1beta1).
func (src *VSphereFailureDomain) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereFailureDomain)

	if err := Convert_v1alpha4_VSphereFailureDomain_To_v1beta1_VSphereFailureDomain(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1.VSphereFailureDomain{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Topology.NetworkConfigurations = restored.Spec.Topology.NetworkConfigurations

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereFailureDomain.
func (dst *VSphereFailureDomain) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereFailureDomain)

	if err := Convert_v1beta1_VSphereFailureDomain_To_v1alpha4_VSphereFailureDomain(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this VSphereFailureDomainList to the Hub version (v1beta1).
func (src *VSphereFailureDomainList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereFailureDomainList)
	return Convert_v1alpha4_VSphereFailureDomainList_To_v1beta1_VSphereFailureDomainList(src, dst, nil)
}

// ConvertFrom converts this VSphereFailureDomainList to the Hub version (v1beta1).
func (dst *VSphereFailureDomainList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereFailureDomainList)
	return Convert_v1beta1_VSphereFailureDomainList_To_v1alpha4_VSphereFailureDomainList(src, dst, nil)
}
