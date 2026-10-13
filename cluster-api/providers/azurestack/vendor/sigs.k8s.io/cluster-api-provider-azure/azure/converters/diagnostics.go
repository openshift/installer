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

package converters

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// GetDiagnosticsProfile converts a CAPZ Diagnostics option to a Azure SDK Diagnostics Profile.
func GetDiagnosticsProfile(diagnostics *infrav1.Diagnostics) *armcompute.DiagnosticsProfile {
	if diagnostics != nil && diagnostics.Boot != nil {
		switch diagnostics.Boot.StorageAccountType {
		case infrav1.DisabledDiagnosticsStorage:
			return &armcompute.DiagnosticsProfile{
				BootDiagnostics: &armcompute.BootDiagnostics{
					Enabled: ptr.To(false),
				},
			}
		case infrav1.ManagedDiagnosticsStorage:
			return &armcompute.DiagnosticsProfile{
				BootDiagnostics: &armcompute.BootDiagnostics{
					Enabled: ptr.To(true),
				},
			}
		case infrav1.UserManagedDiagnosticsStorage:
			if diagnostics.Boot.UserManaged != nil {
				return &armcompute.DiagnosticsProfile{
					BootDiagnostics: &armcompute.BootDiagnostics{
						Enabled:    ptr.To(true),
						StorageURI: &diagnostics.Boot.UserManaged.StorageAccountURI,
					},
				}
			}
		}
	}

	return nil
}
