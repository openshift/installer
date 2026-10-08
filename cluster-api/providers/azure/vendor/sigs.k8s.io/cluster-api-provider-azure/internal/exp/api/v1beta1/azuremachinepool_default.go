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

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/uuid"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	infrav1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"
	apiinternal "sigs.k8s.io/cluster-api-provider-azure/internal/api/v1beta1"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	utilSSH "sigs.k8s.io/cluster-api-provider-azure/util/ssh"
)

// SetDefaults sets the default values for an AzureMachinePool.
func SetDefaults(amp *infrav1exp.AzureMachinePool, c client.Client) error {
	var errs []error
	if err := SetDefaultSSHPublicKey(amp); err != nil {
		errs = append(errs, errors.Wrap(err, "failed to set default SSH public key"))
	}

	if err := SetIdentityDefaults(amp, c); err != nil {
		errs = append(errs, errors.Wrap(err, "failed to set default managed identity defaults"))
	}
	SetDiagnosticsDefaults(amp)
	SetNetworkInterfacesDefaults(amp)
	SetOSDiskDefaults(amp)

	return kerrors.NewAggregate(errs)
}

// SetDefaultSSHPublicKey sets the default SSHPublicKey for an AzureMachinePool.
func SetDefaultSSHPublicKey(amp *infrav1exp.AzureMachinePool) error {
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
func SetIdentityDefaults(amp *infrav1exp.AzureMachinePool, c client.Client) error {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if amp.Spec.RoleAssignmentName != "" && amp.Spec.SystemAssignedIdentityRole != nil && amp.Spec.SystemAssignedIdentityRole.Name != "" { //nolint:staticcheck
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return nil
	}
	if amp.Spec.Identity == infrav1.VMIdentitySystemAssigned {
		machinePool, err := azureutil.FindParentMachinePoolWithRetryV1Beta1(amp.Name, c, 5)
		if err != nil {
			return errors.Wrap(err, "failed to find parent machine pool")
		}

		ownerAzureClusterName, ownerAzureClusterNamespace, err := apiinternal.GetOwnerAzureClusterNameAndNamespace(c, machinePool.Spec.ClusterName, machinePool.Namespace, 5)
		if err != nil {
			return errors.Wrap(err, "failed to get owner cluster")
		}

		subscriptionID, err := apiinternal.GetSubscriptionID(c, ownerAzureClusterName, ownerAzureClusterNamespace, 5)
		if err != nil {
			return errors.Wrap(err, "failed to get subscription ID")
		}

		if amp.Spec.SystemAssignedIdentityRole == nil {
			amp.Spec.SystemAssignedIdentityRole = &infrav1.SystemAssignedIdentityRole{}
		}
		if amp.Spec.RoleAssignmentName != "" { //nolint:staticcheck
			amp.Spec.SystemAssignedIdentityRole.Name = amp.Spec.RoleAssignmentName //nolint:staticcheck
			amp.Spec.RoleAssignmentName = ""                                       //nolint:staticcheck
		} else if amp.Spec.SystemAssignedIdentityRole.Name == "" {
			amp.Spec.SystemAssignedIdentityRole.Name = string(uuid.NewUUID())
		}
		if amp.Spec.SystemAssignedIdentityRole.Scope == "" {
			// Default scope to the subscription.
			amp.Spec.SystemAssignedIdentityRole.Scope = fmt.Sprintf("/subscriptions/%s/", subscriptionID)
		}
		if amp.Spec.SystemAssignedIdentityRole.DefinitionID == "" {
			// Default role definition ID to Contributor role.
			amp.Spec.SystemAssignedIdentityRole.DefinitionID = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Authorization/roleDefinitions/%s", subscriptionID, apiinternal.ContributorRoleID)
		}
	}
	return nil
}

// SetSpotEvictionPolicyDefaults sets the defaults for the spot VM eviction policy.
func SetSpotEvictionPolicyDefaults(amp *infrav1exp.AzureMachinePool) {
	if amp.Spec.Template.SpotVMOptions != nil && amp.Spec.Template.SpotVMOptions.EvictionPolicy == nil {
		defaultPolicy := infrav1.SpotEvictionPolicyDeallocate
		if amp.Spec.Template.OSDisk.DiffDiskSettings != nil && amp.Spec.Template.OSDisk.DiffDiskSettings.Option == "Local" {
			defaultPolicy = infrav1.SpotEvictionPolicyDelete
		}
		amp.Spec.Template.SpotVMOptions.EvictionPolicy = &defaultPolicy
	}
}

// SetDiagnosticsDefaults sets the defaults for Diagnostic settings for an AzureMachinePool.
func SetDiagnosticsDefaults(amp *infrav1exp.AzureMachinePool) {
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
func SetNetworkInterfacesDefaults(amp *infrav1exp.AzureMachinePool) {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if (amp.Spec.Template.SubnetName != "" || amp.Spec.Template.AcceleratedNetworking != nil) && len(amp.Spec.Template.NetworkInterfaces) > 0 { //nolint:staticcheck
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return
	}

	if len(amp.Spec.Template.NetworkInterfaces) == 0 {
		amp.Spec.Template.NetworkInterfaces = []infrav1.NetworkInterface{
			{
				SubnetName:            amp.Spec.Template.SubnetName,            //nolint:staticcheck
				AcceleratedNetworking: amp.Spec.Template.AcceleratedNetworking, //nolint:staticcheck
			},
		}
		amp.Spec.Template.SubnetName = ""             //nolint:staticcheck
		amp.Spec.Template.AcceleratedNetworking = nil //nolint:staticcheck
	}

	// Ensure that PrivateIPConfigs defaults to 1 if not specified.
	for i := 0; i < len(amp.Spec.Template.NetworkInterfaces); i++ {
		if amp.Spec.Template.NetworkInterfaces[i].PrivateIPConfigs == 0 {
			amp.Spec.Template.NetworkInterfaces[i].PrivateIPConfigs = 1
		}
	}
}

// SetOSDiskDefaults sets the defaults for the OSDisk.
func SetOSDiskDefaults(amp *infrav1exp.AzureMachinePool) {
	if amp.Spec.Template.OSDisk.OSType == "" {
		amp.Spec.Template.OSDisk.OSType = "Linux"
	}
	if amp.Spec.Template.OSDisk.CachingType == "" {
		amp.Spec.Template.OSDisk.CachingType = "None"
	}
}
