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

package v1beta1

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/uuid"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	utilSSH "sigs.k8s.io/cluster-api-provider-azure/util/ssh"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SetDefaults sets the default values for an AzureMachinePool.
func (amp *AzureMachinePool) SetDefaults(client client.Client) error {
	var errs []error
	if err := amp.SetDefaultSSHPublicKey(); err != nil {
		errs = append(errs, errors.Wrap(err, "failed to set default SSH public key"))
	}

	if err := amp.SetIdentityDefaults(client); err != nil {
		errs = append(errs, errors.Wrap(err, "failed to set default managed identity defaults"))
	}
	amp.SetDiagnosticsDefaults()
	amp.SetNetworkInterfacesDefaults()

	return kerrors.NewAggregate(errs)
}

// SetDefaultSSHPublicKey sets the default SSHPublicKey for an AzureMachinePool.
func (amp *AzureMachinePool) SetDefaultSSHPublicKey() error {
	if sshKeyData := amp.Spec.Template.SSHPublicKey; sshKeyData == "" {
		_, publicRsaKey, err := utilSSH.GenerateSSHKey()
		if err != nil {
			return err
		}

		amp.Spec.Template.SSHPublicKey = base64.StdEncoding.EncodeToString(ssh.MarshalAuthorizedKey(publicRsaKey))
	}
	return nil
}

// SetIdentityDefaults sets the defaults for VMSS Identity.
func (amp *AzureMachinePool) SetIdentityDefaults(client client.Client) error {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if amp.Spec.RoleAssignmentName != "" && amp.Spec.SystemAssignedIdentityRole != nil && amp.Spec.SystemAssignedIdentityRole.Name != "" {
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return nil
	}
	if amp.Spec.Identity == infrav1.VMIdentitySystemAssigned {
		machinePool, err := azureutil.FindParentMachinePoolWithRetry(amp.Name, client, 5)
		if err != nil {
			return errors.Wrap(err, "failed to find parent machine pool")
		}

		ownerAzureClusterName, ownerAzureClusterNamespace, err := infrav1.GetOwnerAzureClusterNameAndNamespace(client, machinePool.Spec.ClusterName, machinePool.Namespace, 5)
		if err != nil {
			return errors.Wrap(err, "failed to get owner cluster")
		}

		subscriptionID, err := infrav1.GetSubscriptionID(client, ownerAzureClusterName, ownerAzureClusterNamespace, 5)
		if err != nil {
			return errors.Wrap(err, "failed to get subscription ID")
		}

		if amp.Spec.SystemAssignedIdentityRole == nil {
			amp.Spec.SystemAssignedIdentityRole = &infrav1.SystemAssignedIdentityRole{}
		}
		if amp.Spec.RoleAssignmentName != "" {
			amp.Spec.SystemAssignedIdentityRole.Name = amp.Spec.RoleAssignmentName
			amp.Spec.RoleAssignmentName = ""
		} else if amp.Spec.SystemAssignedIdentityRole.Name == "" {
			amp.Spec.SystemAssignedIdentityRole.Name = string(uuid.NewUUID())
		}
		if amp.Spec.SystemAssignedIdentityRole.Scope == "" {
			// Default scope to the subscription.
			amp.Spec.SystemAssignedIdentityRole.Scope = fmt.Sprintf("/subscriptions/%s/", subscriptionID)
		}
		if amp.Spec.SystemAssignedIdentityRole.DefinitionID == "" {
			// Default role definition ID to Contributor role.
			amp.Spec.SystemAssignedIdentityRole.DefinitionID = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s", subscriptionID, infrav1.ContributorRoleID)
		}
	}
	return nil
}

// SetSpotEvictionPolicyDefaults sets the defaults for the spot VM eviction policy.
func (amp *AzureMachinePool) SetSpotEvictionPolicyDefaults() {
	if amp.Spec.Template.SpotVMOptions != nil && amp.Spec.Template.SpotVMOptions.EvictionPolicy == nil {
		defaultPolicy := infrav1.SpotEvictionPolicyDeallocate
		if amp.Spec.Template.OSDisk.DiffDiskSettings != nil && amp.Spec.Template.OSDisk.DiffDiskSettings.Option == "Local" {
			defaultPolicy = infrav1.SpotEvictionPolicyDelete
		}
		amp.Spec.Template.SpotVMOptions.EvictionPolicy = &defaultPolicy
	}
}

// SetDiagnosticsDefaults sets the defaults for Diagnostic settings for an AzureMachinePool.
func (amp *AzureMachinePool) SetDiagnosticsDefaults() {
	bootDefault := &infrav1.BootDiagnostics{
		StorageAccountType: infrav1.ManagedDiagnosticsStorage,
	}

	if amp.Spec.Template.Diagnostics == nil {
		amp.Spec.Template.Diagnostics = &infrav1.Diagnostics{
			Boot: bootDefault,
		}
	}

	if amp.Spec.Template.Diagnostics.Boot == nil {
		amp.Spec.Template.Diagnostics.Boot = bootDefault
	}
}

// SetNetworkInterfacesDefaults sets the defaults for the network interfaces.
func (amp *AzureMachinePool) SetNetworkInterfacesDefaults() {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if (amp.Spec.Template.SubnetName != "" || amp.Spec.Template.AcceleratedNetworking != nil) && len(amp.Spec.Template.NetworkInterfaces) > 0 {
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return
	}

	if len(amp.Spec.Template.NetworkInterfaces) == 0 {
		amp.Spec.Template.NetworkInterfaces = []infrav1.NetworkInterface{
			{
				SubnetName:            amp.Spec.Template.SubnetName,
				AcceleratedNetworking: amp.Spec.Template.AcceleratedNetworking,
			},
		}
		amp.Spec.Template.SubnetName = ""
		amp.Spec.Template.AcceleratedNetworking = nil
	}

	// Ensure that PrivateIPConfigs defaults to 1 if not specified.
	for i := 0; i < len(amp.Spec.Template.NetworkInterfaces); i++ {
		if amp.Spec.Template.NetworkInterfaces[i].PrivateIPConfigs == 0 {
			amp.Spec.Template.NetworkInterfaces[i].PrivateIPConfigs = 1
		}
	}
}
