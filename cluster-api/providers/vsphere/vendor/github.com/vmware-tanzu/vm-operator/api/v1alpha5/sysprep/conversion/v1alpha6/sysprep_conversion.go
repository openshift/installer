// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	"unsafe"

	apiconversion "k8s.io/apimachinery/pkg/conversion"

	vmopv1a5common "github.com/vmware-tanzu/vm-operator/api/v1alpha5/common"
	vmopv1a5sysprep "github.com/vmware-tanzu/vm-operator/api/v1alpha5/sysprep"
	vmopv1sysprep "github.com/vmware-tanzu/vm-operator/api/v1alpha6/sysprep"
)

// Convert_sysprep_Sysprep_To_sysprep_Sysprep converts the Sysprep from v1alpha6
// to v1alpha5.
// Please see https://github.com/kubernetes/code-generator/issues/172 for why
// this function exists in this directory structure.
func Convert_sysprep_Sysprep_To_sysprep_Sysprep(
	in *vmopv1sysprep.Sysprep, out *vmopv1a5sysprep.Sysprep, s apiconversion.Scope) error {

	out.GUIRunOnce = (*vmopv1a5sysprep.GUIRunOnce)(unsafe.Pointer(in.GUIRunOnce))
	out.GUIUnattended = (*vmopv1a5sysprep.GUIUnattended)(unsafe.Pointer(in.GUIUnattended))
	out.LicenseFilePrintData = (*vmopv1a5sysprep.LicenseFilePrintData)(unsafe.Pointer(in.LicenseFilePrintData))
	out.UserData = *(*vmopv1a5sysprep.UserData)(unsafe.Pointer(&in.UserData))
	out.ExpirePasswordAfterNextLogin = in.ExpirePasswordAfterNextLogin

	if sc := in.ScriptText; sc != nil {
		out.ScriptText = &vmopv1a5common.ValueOrSecretKeySelector{
			Value: sc.Value,
		}
		if sc.From != nil {
			out.ScriptText.From = &vmopv1a5common.SecretKeySelector{
				Name: sc.From.Name,
				Key:  sc.From.Key,
			}
		}
	}

	if id := in.Identification; id != nil {
		out.Identification = &vmopv1a5sysprep.Identification{
			DomainAdmin:         id.DomainAdmin,
			DomainAdminPassword: (*vmopv1a5sysprep.DomainPasswordSecretKeySelector)(unsafe.Pointer(id.DomainAdminPassword)),
			DomainOU:            id.DomainOU,
			JoinWorkgroup:       id.JoinWorkgroup,
		}
	}

	return nil
}
