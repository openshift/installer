// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	"maps"

	apiconversion "k8s.io/apimachinery/pkg/conversion"

	vmopv1a5common "github.com/vmware-tanzu/vm-operator/api/v1alpha5/common"
	vmopv1common "github.com/vmware-tanzu/vm-operator/api/v1alpha6/common"
)

// Convert_common_ObjectMeta_To_common_ObjectMeta converts the ObjectMeta from
// v1alpha6 to v1alpha3.
// Please see https://github.com/kubernetes/code-generator/issues/172 for why
// this function exists in this directory structure.
func Convert_common_ObjectMeta_To_common_ObjectMeta(
	in *vmopv1common.ObjectMeta, out *vmopv1a5common.ObjectMeta, s apiconversion.Scope) error {

	if in.Annotations == nil {
		out.Annotations = nil
	} else {
		out.Annotations = maps.Clone(in.Annotations)
	}

	if in.Labels == nil {
		out.Labels = nil
	} else {
		out.Labels = maps.Clone(in.Labels)
	}

	return nil
}
