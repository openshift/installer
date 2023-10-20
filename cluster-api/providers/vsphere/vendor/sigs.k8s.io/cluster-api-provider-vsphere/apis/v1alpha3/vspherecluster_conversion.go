/*
Copyright 2019 The Kubernetes Authors.

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

//nolint:forcetypeassert,golint,revive,stylecheck
package v1alpha3

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereCluster to the Hub version (v1beta1).
func (src *VSphereCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereCluster)
	if err := Convert_v1alpha3_VSphereCluster_To_v1beta1_VSphereCluster(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1beta1.VSphereCluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	if restored.Spec.IdentityRef != nil {
		dst.Spec.IdentityRef = restored.Spec.IdentityRef
	}
	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereCluster.
func (dst *VSphereCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta1.VSphereCluster)
	if err := Convert_v1beta1_VSphereCluster_To_v1alpha3_VSphereCluster(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}
	return nil
}

// ConvertTo converts this VSphereClusterList to the Hub version (v1beta1).
func (src *VSphereClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereClusterList)
	return Convert_v1alpha3_VSphereClusterList_To_v1beta1_VSphereClusterList(src, dst, nil)
}

// ConvertFrom converts this VSphereVM to the Hub version (v1beta1).
func (dst *VSphereClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta1.VSphereClusterList)
	return Convert_v1beta1_VSphereClusterList_To_v1alpha3_VSphereClusterList(src, dst, nil)
}

func Convert_v1alpha3_VSphereClusterSpec_To_v1beta1_VSphereClusterSpec(in *VSphereClusterSpec, out *infrav1beta1.VSphereClusterSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_VSphereClusterSpec_To_v1beta1_VSphereClusterSpec(in, out, s)
}
