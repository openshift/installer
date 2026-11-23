// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha6

import (
	vmopv1common "github.com/vmware-tanzu/vm-operator/api/v1alpha6/common"
)

type PolicySpec vmopv1common.LocalObjectRef

type PolicyStatus struct {
	PolicySpec `json:",inline"`

	// Generation describes the observed generation of the policy applied to
	// this VM.
	Generation int64 `json:"generation"`
}
