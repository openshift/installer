/*
Copyright 2019 The Kubernetes Authors.

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

package virtualmachines

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	azprovider "sigs.k8s.io/cloud-provider-azure/pkg/provider"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/identities"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/networkinterfaces"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/publicips"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const serviceName = "virtualmachine"
const vmMissingUAI = "VM is missing expected user assigned identity with client ID: "

// VMScope defines the scope interface for a virtual machines service.
type VMScope interface {
	azure.Authorizer
	azure.AsyncStatusUpdater
	VMSpec() azure.ResourceSpecGetter
	SetAnnotation(string, string)
	SetProviderID(string)
	SetAddresses([]corev1.NodeAddress)
	SetVMState(infrav1.ProvisioningState)
	SetConditionFalse(clusterv1.ConditionType, string, clusterv1.ConditionSeverity, string)
}

// Service provides operations on Azure resources.
type Service struct {
	Scope VMScope
	async.Reconciler
	interfacesGetter async.Getter
	publicIPsGetter  async.Getter
	identitiesGetter identities.Client
}

// New creates a new service.
func New(scope VMScope) (*Service, error) {
	Client, err := NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	identitiesSvc, err := identities.NewClient(scope)
	if err != nil {
		return nil, err
	}
	interfacesSvc, err := networkinterfaces.NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	publicIPsSvc, err := publicips.NewClient(scope, scope.DefaultedAzureCallTimeout())
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope:            scope,
		interfacesGetter: interfacesSvc,
		publicIPsGetter:  publicIPsSvc,
		identitiesGetter: identitiesSvc,
		Reconciler: async.New[armcompute.VirtualMachinesClientCreateOrUpdateResponse,
			armcompute.VirtualMachinesClientDeleteResponse](scope, Client, Client),
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a virtual machine.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualmachines.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	vmSpec := s.Scope.VMSpec()
	if vmSpec == nil {
		return nil
	}

	result, err := s.CreateOrUpdateResource(ctx, vmSpec, serviceName)
	s.Scope.UpdatePutStatus(infrav1.VMRunningCondition, serviceName, err)
	// Set the DiskReady condition here since the disk gets created with the VM.
	s.Scope.UpdatePutStatus(infrav1.DisksReadyCondition, serviceName, err)
	if err == nil && result != nil {
		vm, ok := result.(armcompute.VirtualMachine)
		if !ok {
			return errors.Errorf("%T is not an armcompute.VirtualMachine", result)
		}
		infraVM := converters.SDKToVM(vm)
		// Transform the VM resource representation to conform to the cloud-provider-azure representation
		providerID, err := azprovider.ConvertResourceGroupNameToLower(azureutil.ProviderIDPrefix + infraVM.ID)
		if err != nil {
			return errors.Wrapf(err, "failed to parse VM ID %s", infraVM.ID)
		}
		s.Scope.SetProviderID(providerID)
		s.Scope.SetAnnotation("cluster-api-provider-azure", "true")

		// Discover addresses for NICs associated with the VM
		addresses, err := s.getAddresses(ctx, vm, vmSpec.ResourceGroupName())
		if err != nil {
			return errors.Wrap(err, "failed to fetch VM addresses")
		}
		s.Scope.SetAddresses(addresses)
		s.Scope.SetVMState(infraVM.State)

		spec, ok := vmSpec.(*VMSpec)
		if !ok {
			return errors.Errorf("%T is not a valid VM spec", vmSpec)
		}

		err = s.checkUserAssignedIdentities(ctx, spec.UserAssignedIdentities, infraVM.UserAssignedIdentities)
		if err != nil {
			return errors.Wrap(err, "failed to check user assigned identities")
		}
	}
	return err
}

// Delete deletes the virtual machine with the provided name.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualmachines.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, s.Scope.DefaultedAzureServiceReconcileTimeout())
	defer cancel()

	vmSpec := s.Scope.VMSpec()
	if vmSpec == nil {
		return nil
	}

	err := s.DeleteResource(ctx, vmSpec, serviceName)
	if err != nil {
		s.Scope.SetVMState(infrav1.Deleting)
	} else {
		s.Scope.SetVMState(infrav1.Deleted)
	}
	s.Scope.UpdateDeleteStatus(infrav1.VMRunningCondition, serviceName, err)
	return err
}

func (s *Service) checkUserAssignedIdentities(ctx context.Context, specIdentities []infrav1.UserAssignedIdentity, vmIdentities []infrav1.UserAssignedIdentity) error {
	expectedMap := make(map[string]struct{})
	actualMap := make(map[string]struct{})

	// Create a map of the expected identities. The ProviderID is converted to match the format of the VM identity.
	for _, expectedIdentity := range specIdentities {
		identitiesClient := s.identitiesGetter
		parsed, err := azureutil.ParseResourceID(expectedIdentity.ProviderID)
		if err != nil {
			return err
		}
		if parsed.SubscriptionID != s.Scope.SubscriptionID() {
			identitiesClient, err = identities.NewClientBySub(s.Scope, parsed.SubscriptionID)
			if err != nil {
				return errors.Wrapf(err, "failed to create identities client from subscription ID %s", parsed.SubscriptionID)
			}
		}
		expectedClientID, err := identitiesClient.GetClientID(ctx, expectedIdentity.ProviderID)
		if err != nil {
			return errors.Wrap(err, "failed to get client ID")
		}
		expectedMap[expectedClientID] = struct{}{}
	}

	// Create a map of the actual identities from the vm.
	for _, actualIdentity := range vmIdentities {
		actualMap[actualIdentity.ProviderID] = struct{}{}
	}

	// Check if the expected identities are present in the vm.
	for expectedKey := range expectedMap {
		_, exists := actualMap[expectedKey]
		if !exists {
			s.Scope.SetConditionFalse(infrav1.VMIdentitiesReadyCondition, infrav1.UserAssignedIdentityMissingReason, clusterv1.ConditionSeverityWarning, vmMissingUAI+expectedKey)
			return nil
		}
	}

	return nil
}

func (s *Service) getAddresses(ctx context.Context, vm armcompute.VirtualMachine, rgName string) ([]corev1.NodeAddress, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualmachines.Service.getAddresses")
	defer done()

	addresses := []corev1.NodeAddress{
		{
			Type:    corev1.NodeInternalDNS,
			Address: ptr.Deref(vm.Name, ""),
		},
	}
	if vm.Properties.NetworkProfile.NetworkInterfaces == nil {
		return addresses, nil
	}
	for _, nicRef := range vm.Properties.NetworkProfile.NetworkInterfaces {
		// The full ID includes the name at the very end. Split the string and pull the last element
		// Ex: /subscriptions/$SUB/resourceGroups/$RG/providers/Microsoft.Network/networkInterfaces/$NICNAME
		// We'll check to see if ID is nil and bail early if we don't have it
		if nicRef.ID == nil {
			continue
		}
		nicName := getResourceNameByID(ptr.Deref(nicRef.ID, ""))

		// Fetch nic and append its addresses
		existingNic, err := s.interfacesGetter.Get(ctx, &networkinterfaces.NICSpec{
			Name:          nicName,
			ResourceGroup: rgName,
		})
		if err != nil {
			return addresses, err
		}

		nic, ok := existingNic.(armnetwork.Interface)
		if !ok {
			return nil, errors.Errorf("%T is not an armnetwork.Interface", existingNic)
		}

		if nic.Properties.IPConfigurations == nil {
			continue
		}
		for _, ipConfig := range nic.Properties.IPConfigurations {
			if ipConfig != nil && ipConfig.Properties != nil && ipConfig.Properties.PrivateIPAddress != nil {
				addresses = append(addresses,
					corev1.NodeAddress{
						Type:    corev1.NodeInternalIP,
						Address: ptr.Deref(ipConfig.Properties.PrivateIPAddress, ""),
					},
				)
			}

			if ipConfig.Properties.PublicIPAddress == nil {
				continue
			}
			// ID is the only field populated in PublicIPAddress sub-resource.
			// Thus, we have to go fetch the publicIP with the name.
			publicIPName := getResourceNameByID(ptr.Deref(ipConfig.Properties.PublicIPAddress.ID, ""))
			publicNodeAddress, err := s.getPublicIPAddress(ctx, publicIPName, rgName)
			if err != nil {
				return addresses, err
			}
			addresses = append(addresses, publicNodeAddress)
		}
	}

	return addresses, nil
}

// getPublicIPAddress will fetch a public ip address resource by name and return a nodeaddresss representation.
func (s *Service) getPublicIPAddress(ctx context.Context, publicIPAddressName string, rgName string) (corev1.NodeAddress, error) {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "virtualmachines.Service.getPublicIPAddress")
	defer done()

	retAddress := corev1.NodeAddress{}
	result, err := s.publicIPsGetter.Get(ctx, &publicips.PublicIPSpec{
		Name:          publicIPAddressName,
		ResourceGroup: rgName,
	})
	if err != nil {
		return retAddress, err
	}

	publicIP, ok := result.(armnetwork.PublicIPAddress)
	if !ok {
		return retAddress, errors.Errorf("%T is not an armnetwork.PublicIPAddress", result)
	}

	retAddress.Type = corev1.NodeExternalIP
	retAddress.Address = ptr.Deref(publicIP.Properties.IPAddress, "")

	return retAddress, nil
}

// getResourceNameById takes a resource ID like
// `/subscriptions/$SUB/resourceGroups/$RG/providers/Microsoft.Network/networkInterfaces/$NICNAME`
// and parses out the string after the last slash.
func getResourceNameByID(resourceID string) string {
	explodedResourceID := strings.Split(resourceID, "/")
	resourceName := explodedResourceID[len(explodedResourceID)-1]
	return resourceName
}

// IsManaged returns always returns true as CAPZ does not support BYO VM.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
