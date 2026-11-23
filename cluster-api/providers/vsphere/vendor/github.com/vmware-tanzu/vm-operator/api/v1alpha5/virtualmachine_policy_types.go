// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	vmopv1a5common "github.com/vmware-tanzu/vm-operator/api/v1alpha5/common"
)

type PolicySpec vmopv1a5common.LocalObjectRef

type PolicyStatus struct {
	PolicySpec `json:",inline"`

	// Generation describes the observed generation of the policy applied to
	// this VM.
	Generation int64 `json:"generation"`
}
