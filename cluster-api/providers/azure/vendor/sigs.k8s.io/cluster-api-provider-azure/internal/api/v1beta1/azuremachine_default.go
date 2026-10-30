/*
Copyright The Kubernetes Authors.

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

package v1beta1

import (
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/uuid"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	utilSSH "sigs.k8s.io/cluster-api-provider-azure/util/ssh"
)

// ContributorRoleID is the ID of the built-in "Contributor" role.
const ContributorRoleID = "b24988ac-6180-42a0-ab88-20f7382dd24c"

// SetDefaultAzureMachineSpecSSHPublicKey sets the default SSHPublicKey for an AzureMachine.
func SetDefaultAzureMachineSpecSSHPublicKey(s *infrav1.AzureMachineSpec) error {
	if sshKeyData := s.SSHPublicKey; sshKeyData == "" {
		_, publicRsaKey, err := utilSSH.GenerateSSHKey()
		if err != nil {
			return err
		}

		s.SSHPublicKey = base64.StdEncoding.EncodeToString(ssh.MarshalAuthorizedKey(publicRsaKey))
	}
	return nil
}

// SetDefaultAzureMachineSpecDataDisks sets the data disk defaults for an AzureMachine.
func SetDefaultAzureMachineSpecDataDisks(s *infrav1.AzureMachineSpec) {
	set := make(map[int32]struct{})
	// populate all the existing values in the set
	for _, disk := range s.DataDisks {
		if disk.Lun != nil {
			set[*disk.Lun] = struct{}{}
		}
	}
	// Look for unique values for unassigned LUNs
	for i, disk := range s.DataDisks {
		if disk.Lun == nil {
			for l := range s.DataDisks {
				lun := int32(l)
				if _, ok := set[lun]; !ok {
					s.DataDisks[i].Lun = &lun
					set[lun] = struct{}{}
					break
				}
			}
		}
		if disk.CachingType == "" {
			if s.DataDisks[i].ManagedDisk != nil &&
				s.DataDisks[i].ManagedDisk.StorageAccountType == string(armcompute.StorageAccountTypesUltraSSDLRS) {
				s.DataDisks[i].CachingType = string(armcompute.CachingTypesNone)
			} else {
				s.DataDisks[i].CachingType = string(armcompute.CachingTypesReadWrite)
			}
		}
	}
}

// SetDefaultAzureMachineSpecIdentity sets the defaults for VM Identity.
func SetDefaultAzureMachineSpecIdentity(s *infrav1.AzureMachineSpec, subscriptionID string) {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if s.RoleAssignmentName != "" && s.SystemAssignedIdentityRole != nil && s.SystemAssignedIdentityRole.Name != "" { //nolint:staticcheck
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return
	}
	if s.Identity == infrav1.VMIdentitySystemAssigned {
		if s.SystemAssignedIdentityRole == nil {
			s.SystemAssignedIdentityRole = &infrav1.SystemAssignedIdentityRole{}
		}
		if s.RoleAssignmentName != "" { //nolint:staticcheck
			// Move the existing value from the deprecated RoleAssignmentName field.
			s.SystemAssignedIdentityRole.Name = s.RoleAssignmentName //nolint:staticcheck
			s.RoleAssignmentName = ""                                //nolint:staticcheck
		} else if s.SystemAssignedIdentityRole.Name == "" {
			// Default role name to a generated UUID.
			s.SystemAssignedIdentityRole.Name = string(uuid.NewUUID())
		}
		if s.SystemAssignedIdentityRole.Scope == "" && subscriptionID != "" {
			// Default scope to the subscription.
			s.SystemAssignedIdentityRole.Scope = fmt.Sprintf("/subscriptions/%s/", subscriptionID)
		}
		if s.SystemAssignedIdentityRole.DefinitionID == "" && subscriptionID != "" {
			// Default role definition ID to Contributor role.
			s.SystemAssignedIdentityRole.DefinitionID = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s", subscriptionID, ContributorRoleID)
		}
	}
}

// setDefaultAzureMachineSpecSpotEvictionPolicy sets the defaults for the spot VM eviction policy.
func setDefaultAzureMachineSpecSpotEvictionPolicy(s *infrav1.AzureMachineSpec) {
	if s.SpotVMOptions != nil && s.SpotVMOptions.EvictionPolicy == nil {
		defaultPolicy := infrav1.SpotEvictionPolicyDeallocate
		if s.OSDisk.DiffDiskSettings != nil && s.OSDisk.DiffDiskSettings.Option == "Local" {
			defaultPolicy = infrav1.SpotEvictionPolicyDelete
		}
		s.SpotVMOptions.EvictionPolicy = &defaultPolicy
	}
}

// setDefaultAzureMachineSpecDiagnostics sets the defaults for Diagnostic settings for an AzureMachine.
func setDefaultAzureMachineSpecDiagnostics(s *infrav1.AzureMachineSpec) {
	bootDiagnosticsDefault := &infrav1.BootDiagnostics{
		StorageAccountType: infrav1.ManagedDiagnosticsStorage,
	}

	diagnosticsDefault := &infrav1.Diagnostics{Boot: bootDiagnosticsDefault}

	if s.Diagnostics == nil {
		s.Diagnostics = diagnosticsDefault
	}

	if s.Diagnostics.Boot == nil {
		s.Diagnostics.Boot = bootDiagnosticsDefault
	}
}

// SetDefaultAzureMachineSpecNetworkInterfaces sets the defaults for the network interfaces.
func SetDefaultAzureMachineSpecNetworkInterfaces(s *infrav1.AzureMachineSpec) {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if (s.SubnetName != "" || s.AcceleratedNetworking != nil) && len(s.NetworkInterfaces) > 0 { //nolint:staticcheck
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return
	}

	if len(s.NetworkInterfaces) == 0 {
		s.NetworkInterfaces = []infrav1.NetworkInterface{
			{
				SubnetName:            s.SubnetName,            //nolint:staticcheck
				AcceleratedNetworking: s.AcceleratedNetworking, //nolint:staticcheck
			},
		}
		s.SubnetName = ""             //nolint:staticcheck
		s.AcceleratedNetworking = nil //nolint:staticcheck
	}

	// Ensure that PrivateIPConfigs defaults to 1 if not specified.
	for i := 0; i < len(s.NetworkInterfaces); i++ {
		if s.NetworkInterfaces[i].PrivateIPConfigs == 0 {
			s.NetworkInterfaces[i].PrivateIPConfigs = 1
		}
	}
}

// SetDefaultsAzureMachine sets the defaults for the AzureMachine.
func SetDefaultsAzureMachine(m *infrav1.AzureMachine, client client.Client) error {
	var errs []error
	if err := SetDefaultAzureMachineSpecSSHPublicKey(&m.Spec); err != nil {
		errs = append(errs, errors.Wrap(err, "failed to set default SSH public key"))
	}

	// Fetch the Cluster.
	clusterName, ok := m.Labels[clusterv1.ClusterNameLabel]
	if !ok {
		errs = append(errs, errors.Errorf("failed to fetch ClusterName for AzureMachine %s/%s", m.Namespace, m.Name))
	}

	ownerAzureClusterName, ownerAzureClusterNamespace, err := GetOwnerAzureClusterNameAndNamespace(client, clusterName, m.Namespace, 5)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "failed to fetch owner cluster for AzureMachine %s/%s", m.Namespace, m.Name))
	}

	subscriptionID, err := GetSubscriptionID(client, ownerAzureClusterName, ownerAzureClusterNamespace, 5)
	if err != nil {
		errs = append(errs, errors.Wrapf(err, "failed to fetch subscription ID for AzureMachine %s/%s", m.Namespace, m.Name))
	}

	SetDefaultAzureMachineSpecDataDisks(&m.Spec)
	SetDefaultAzureMachineSpecIdentity(&m.Spec, subscriptionID)
	setDefaultAzureMachineSpecSpotEvictionPolicy(&m.Spec)
	setDefaultAzureMachineSpecDiagnostics(&m.Spec)
	SetDefaultAzureMachineSpecNetworkInterfaces(&m.Spec)

	return kerrors.NewAggregate(errs)
}
