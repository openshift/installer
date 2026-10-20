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

	vmwarev1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
)

// ProviderServiceAccount is a HubSpokeConverter for the ProviderServiceAccount API type.
var ProviderServiceAccount = conversion.NewHubSpokeConverter(&vmwarev1.ProviderServiceAccount{},
	conversion.NewSpokeConverter(&vmwarev1beta1.ProviderServiceAccount{}, ConvertProviderServiceAccountHubToV1Beta1, ConvertProviderServiceAccountV1Beta1ToHub),
)

// ConvertProviderServiceAccountV1Beta1ToHub converts a v1beta1 ProviderServiceAccount to a hub ProviderServiceAccount.
func ConvertProviderServiceAccountV1Beta1ToHub(_ context.Context, src *vmwarev1beta1.ProviderServiceAccount, dst *vmwarev1.ProviderServiceAccount) error {
	return vmwarev1beta1.Convert_v1beta1_ProviderServiceAccount_To_v1beta2_ProviderServiceAccount(src, dst, nil)
}

// ConvertProviderServiceAccountHubToV1Beta1 converts a hub ProviderServiceAccount to a v1beta1 ProviderServiceAccount.
func ConvertProviderServiceAccountHubToV1Beta1(_ context.Context, src *vmwarev1.ProviderServiceAccount, dst *vmwarev1beta1.ProviderServiceAccount) error {
	return vmwarev1beta1.Convert_v1beta2_ProviderServiceAccount_To_v1beta1_ProviderServiceAccount(src, dst, nil)
}
