/*
Copyright 2020 The Kubernetes Authors.

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
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
)

// ErrUserAssignedIdentitiesNotFound is the error thrown when user assigned identities is not passed with the identity type being UserAssigned.
var ErrUserAssignedIdentitiesNotFound = errors.New("the user-assigned identity provider ids must not be null or empty for 'UserAssigned' identity type")

// VMIdentityToVMSDK converts CAPZ VM identity to Azure SDK identity.
func VMIdentityToVMSDK(identity infrav1.VMIdentity, uami []infrav1.UserAssignedIdentity) (*armcompute.VirtualMachineIdentity, error) {
	if identity == infrav1.VMIdentitySystemAssigned {
		return &armcompute.VirtualMachineIdentity{
			Type: ptr.To(armcompute.ResourceIdentityTypeSystemAssigned),
		}, nil
	}

	if identity == infrav1.VMIdentityUserAssigned {
		userIdentitiesMap, err := UserAssignedIdentitiesToVMSDK(uami)
		if err != nil {
			return nil, errors.Wrap(err, "failed to assign VM identity")
		}

		return &armcompute.VirtualMachineIdentity{
			Type:                   ptr.To(armcompute.ResourceIdentityTypeUserAssigned),
			UserAssignedIdentities: userIdentitiesMap,
		}, nil
	}

	return nil, nil
}

// UserAssignedIdentitiesToVMSDK converts CAPZ user assigned identities associated with the Virtual Machine to Azure SDK identities
// The user identity dictionary key references will be ARM resource ids in the form:
// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'.
func UserAssignedIdentitiesToVMSDK(identities []infrav1.UserAssignedIdentity) (map[string]*armcompute.UserAssignedIdentitiesValue, error) {
	if len(identities) == 0 {
		return nil, ErrUserAssignedIdentitiesNotFound
	}
	userIdentitiesMap := make(map[string]*armcompute.UserAssignedIdentitiesValue, len(identities))
	for _, id := range identities {
		key := sanitized(id.ProviderID)
		userIdentitiesMap[key] = &armcompute.UserAssignedIdentitiesValue{}
	}

	return userIdentitiesMap, nil
}

// UserAssignedIdentitiesToVMSSSDK converts CAPZ user-assigned identities associated with the Virtual Machine Scale Set to Azure SDK identities
// Similar to UserAssignedIdentitiesToVMSDK.
func UserAssignedIdentitiesToVMSSSDK(identities []infrav1.UserAssignedIdentity) (map[string]*armcompute.UserAssignedIdentitiesValue, error) {
	if len(identities) == 0 {
		return nil, ErrUserAssignedIdentitiesNotFound
	}
	userIdentitiesMap := make(map[string]*armcompute.UserAssignedIdentitiesValue, len(identities))
	for _, id := range identities {
		key := sanitized(id.ProviderID)
		userIdentitiesMap[key] = &armcompute.UserAssignedIdentitiesValue{}
	}

	return userIdentitiesMap, nil
}

// sanitized removes "azure://" prefix from the given id.
func sanitized(id string) string {
	return strings.TrimPrefix(id, azureutil.ProviderIDPrefix)
}
