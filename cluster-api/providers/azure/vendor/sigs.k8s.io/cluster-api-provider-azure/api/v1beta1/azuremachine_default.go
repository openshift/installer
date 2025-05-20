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
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/uuid"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	utilSSH "sigs.k8s.io/cluster-api-provider-azure/util/ssh"
)

// ContributorRoleID is the ID of the built-in "Contributor" role.
const ContributorRoleID = "b24988ac-6180-42a0-ab88-20f7382dd24c"

// SetDefaultSSHPublicKey sets the default SSHPublicKey for an AzureMachine.
func (s *AzureMachineSpec) SetDefaultSSHPublicKey() error {
	if sshKeyData := s.SSHPublicKey; sshKeyData == "" {
		_, publicRsaKey, err := utilSSH.GenerateSSHKey()
		if err != nil {
			return err
		}

		s.SSHPublicKey = base64.StdEncoding.EncodeToString(ssh.MarshalAuthorizedKey(publicRsaKey))
	}
	return nil
}

// SetDataDisksDefaults sets the data disk defaults for an AzureMachine.
func (s *AzureMachineSpec) SetDataDisksDefaults() {
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

// SetIdentityDefaults sets the defaults for VM Identity.
func (s *AzureMachineSpec) SetIdentityDefaults(subscriptionID string) {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if s.RoleAssignmentName != "" && s.SystemAssignedIdentityRole != nil && s.SystemAssignedIdentityRole.Name != "" {
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return
	}
	if s.Identity == VMIdentitySystemAssigned {
		if s.SystemAssignedIdentityRole == nil {
			s.SystemAssignedIdentityRole = &SystemAssignedIdentityRole{}
		}
		if s.RoleAssignmentName != "" {
			// Move the existing value from the deprecated RoleAssignmentName field.
			s.SystemAssignedIdentityRole.Name = s.RoleAssignmentName
			s.RoleAssignmentName = ""
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

// SetSpotEvictionPolicyDefaults sets the defaults for the spot VM eviction policy.
func (s *AzureMachineSpec) SetSpotEvictionPolicyDefaults() {
	if s.SpotVMOptions != nil && s.SpotVMOptions.EvictionPolicy == nil {
		defaultPolicy := SpotEvictionPolicyDeallocate
		if s.OSDisk.DiffDiskSettings != nil && s.OSDisk.DiffDiskSettings.Option == "Local" {
			defaultPolicy = SpotEvictionPolicyDelete
		}
		s.SpotVMOptions.EvictionPolicy = &defaultPolicy
	}
}

// SetDiagnosticsDefaults sets the defaults for Diagnostic settings for an AzureMachinePool.
func (s *AzureMachineSpec) SetDiagnosticsDefaults() {
	bootDiagnosticsDefault := &BootDiagnostics{
		StorageAccountType: ManagedDiagnosticsStorage,
	}

	diagnosticsDefault := &Diagnostics{Boot: bootDiagnosticsDefault}

	if s.Diagnostics == nil {
		s.Diagnostics = diagnosticsDefault
	}

	if s.Diagnostics.Boot == nil {
		s.Diagnostics.Boot = bootDiagnosticsDefault
	}
}

// SetNetworkInterfacesDefaults sets the defaults for the network interfaces.
func (s *AzureMachineSpec) SetNetworkInterfacesDefaults() {
	// Ensure the deprecated fields and new fields are not populated simultaneously
	if (s.SubnetName != "" || s.AcceleratedNetworking != nil) && len(s.NetworkInterfaces) > 0 {
		// Both the deprecated and the new fields are both set, return without changes
		// and reject the request in the validating webhook which runs later.
		return
	}

	if len(s.NetworkInterfaces) == 0 {
		s.NetworkInterfaces = []NetworkInterface{
			{
				SubnetName:            s.SubnetName,
				AcceleratedNetworking: s.AcceleratedNetworking,
			},
		}
		s.SubnetName = ""
		s.AcceleratedNetworking = nil
	}

	// Ensure that PrivateIPConfigs defaults to 1 if not specified.
	for i := 0; i < len(s.NetworkInterfaces); i++ {
		if s.NetworkInterfaces[i].PrivateIPConfigs == 0 {
			s.NetworkInterfaces[i].PrivateIPConfigs = 1
		}
	}
}

// GetOwnerAzureClusterNameAndNamespace returns the owner azure cluster's name and namespace for the given cluster name and namespace.
func GetOwnerAzureClusterNameAndNamespace(cli client.Client, clusterName string, namespace string, maxAttempts int) (azureClusterName string, azureClusterNamespace string, err error) {
	ctx := context.Background()

	ownerCluster := &clusterv1.Cluster{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      clusterName,
	}

	for i := 1; ; i++ {
		if err := cli.Get(ctx, key, ownerCluster); err != nil {
			if i > maxAttempts {
				return "", "", errors.Wrapf(err, "failed to find owner cluster for %s/%s", namespace, clusterName)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return ownerCluster.Spec.InfrastructureRef.Name, ownerCluster.Spec.InfrastructureRef.Namespace, nil
}

// GetSubscriptionID returns the subscription ID for the AzureCluster given the cluster name and namespace.
func GetSubscriptionID(cli client.Client, ownerAzureClusterName string, ownerAzureClusterNamespace string, maxAttempts int) (string, error) {
	ctx := context.Background()

	ownerAzureCluster := &AzureCluster{}
	key := client.ObjectKey{
		Namespace: ownerAzureClusterNamespace,
		Name:      ownerAzureClusterName,
	}
	for i := 1; ; i++ {
		if err := cli.Get(ctx, key, ownerAzureCluster); err != nil {
			if i >= maxAttempts {
				return "", errors.Wrapf(err, "failed to find AzureCluster for owner cluster %s/%s", ownerAzureClusterNamespace, ownerAzureClusterName)
			}
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return ownerAzureCluster.Spec.SubscriptionID, nil
}

// SetDefaults sets to the defaults for the AzureMachineSpec.
func (m *AzureMachine) SetDefaults(client client.Client) error {
	var errs []error
	if err := m.Spec.SetDefaultSSHPublicKey(); err != nil {
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

	m.Spec.SetDataDisksDefaults()
	m.Spec.SetIdentityDefaults(subscriptionID)
	m.Spec.SetSpotEvictionPolicyDefaults()
	m.Spec.SetDiagnosticsDefaults()
	m.Spec.SetNetworkInterfacesDefaults()

	return kerrors.NewAggregate(errs)
}
