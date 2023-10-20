/*
Copyright 2023 The Kubernetes Authors.

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

package privateendpoints

import (
	"context"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
)

// PrivateLinkServiceConnection defines the specification for a private link service connection associated with a private endpoint.
type PrivateLinkServiceConnection struct {
	Name                 string
	PrivateLinkServiceID string
	GroupIDs             []string
	RequestMessage       string
}

// PrivateEndpointSpec defines the specification for a private endpoint.
type PrivateEndpointSpec struct {
	Name                          string
	ResourceGroup                 string
	Location                      string
	CustomNetworkInterfaceName    string
	PrivateIPAddresses            []string
	SubnetID                      string
	ApplicationSecurityGroups     []string
	ManualApproval                bool
	PrivateLinkServiceConnections []PrivateLinkServiceConnection
	AdditionalTags                infrav1.Tags
	ClusterName                   string
}

// ResourceName returns the name of the private endpoint.
func (s *PrivateEndpointSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *PrivateEndpointSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for private endpoints.
func (s *PrivateEndpointSpec) OwnerResourceName() string {
	return ""
}

// Parameters returns the parameters for the PrivateEndpointSpec.
func (s *PrivateEndpointSpec) Parameters(ctx context.Context, existing interface{}) (interface{}, error) {
	_, log, done := tele.StartSpanWithLogger(ctx, "privateendpoints.Service.Parameters")
	defer done()

	privateEndpointProperties := armnetwork.PrivateEndpointProperties{
		Subnet: &armnetwork.Subnet{
			ID: &s.SubnetID,
			Properties: &armnetwork.SubnetPropertiesFormat{
				PrivateEndpointNetworkPolicies:    ptr.To(armnetwork.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
				PrivateLinkServiceNetworkPolicies: ptr.To(armnetwork.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled),
			},
		},
	}

	privateEndpointProperties.CustomNetworkInterfaceName = ptr.To(s.CustomNetworkInterfaceName)

	privateIPAddresses := make([]*armnetwork.PrivateEndpointIPConfiguration, 0, len(s.PrivateIPAddresses))
	for _, address := range s.PrivateIPAddresses {
		ipConfig := &armnetwork.PrivateEndpointIPConfigurationProperties{PrivateIPAddress: ptr.To(address)}

		privateIPAddresses = append(privateIPAddresses, &armnetwork.PrivateEndpointIPConfiguration{
			Properties: ipConfig,
		})
	}
	privateEndpointProperties.IPConfigurations = privateIPAddresses

	privateLinkServiceConnections := make([]*armnetwork.PrivateLinkServiceConnection, 0, len(s.PrivateLinkServiceConnections))
	for _, privateLinkServiceConnection := range s.PrivateLinkServiceConnections {
		linkServiceConnection := &armnetwork.PrivateLinkServiceConnection{
			Name: ptr.To(privateLinkServiceConnection.Name),
			Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
				PrivateLinkServiceID: ptr.To(privateLinkServiceConnection.PrivateLinkServiceID),
			},
		}

		if len(privateLinkServiceConnection.GroupIDs) > 0 {
			linkServiceConnection.Properties.GroupIDs = azure.PtrSlice(&privateLinkServiceConnection.GroupIDs)
		}

		if privateLinkServiceConnection.RequestMessage != "" {
			linkServiceConnection.Properties.RequestMessage = ptr.To(privateLinkServiceConnection.RequestMessage)
		}
		privateLinkServiceConnections = append(privateLinkServiceConnections, linkServiceConnection)
	}

	if s.ManualApproval {
		privateEndpointProperties.ManualPrivateLinkServiceConnections = privateLinkServiceConnections
		privateEndpointProperties.PrivateLinkServiceConnections = []*armnetwork.PrivateLinkServiceConnection{}
	} else {
		privateEndpointProperties.PrivateLinkServiceConnections = privateLinkServiceConnections
		privateEndpointProperties.ManualPrivateLinkServiceConnections = []*armnetwork.PrivateLinkServiceConnection{}
	}

	applicationSecurityGroups := make([]armnetwork.ApplicationSecurityGroup, 0, len(s.ApplicationSecurityGroups))

	for _, applicationSecurityGroup := range s.ApplicationSecurityGroups {
		applicationSecurityGroups = append(applicationSecurityGroups, armnetwork.ApplicationSecurityGroup{
			ID: ptr.To(applicationSecurityGroup),
		})
	}

	privateEndpointProperties.ApplicationSecurityGroups = azure.PtrSlice(&applicationSecurityGroups)

	newPrivateEndpoint := armnetwork.PrivateEndpoint{
		Name:       ptr.To(s.Name),
		Properties: &privateEndpointProperties,
		Tags: converters.TagsToMap(infrav1.Build(infrav1.BuildParams{
			ClusterName: s.ClusterName,
			Lifecycle:   infrav1.ResourceLifecycleOwned,
			Name:        ptr.To(s.Name),
			Additional:  s.AdditionalTags,
		})),
	}

	if s.Location != "" {
		newPrivateEndpoint.Location = ptr.To(s.Location)
	}

	if existing != nil {
		existingPE, ok := existing.(armnetwork.PrivateEndpoint)
		if !ok {
			return nil, errors.Errorf("%T is not a network.PrivateEndpoint", existing)
		}

		ps := ptr.Deref(existingPE.Properties.ProvisioningState, "")
		if string(ps) != string(infrav1.Canceled) && string(ps) != string(infrav1.Failed) && string(ps) != string(infrav1.Succeeded) {
			return nil, azure.WithTransientError(errors.Errorf("Unable to update existing private endpoint in non-terminal state. Service Endpoint must be in one of the following provisioning states: Canceled, Failed, or Succeeded. Actual state: %s", ps), 20*time.Second)
		}

		normalizedExistingPE := normalizePrivateEndpoint(existingPE, newPrivateEndpoint)
		normalizedExistingPE = sortSlicesPrivateEndpoint(normalizedExistingPE)

		newPrivateEndpoint = sortSlicesPrivateEndpoint(newPrivateEndpoint)

		diff := cmp.Diff(&normalizedExistingPE, &newPrivateEndpoint)
		if diff == "" {
			// PrivateEndpoint is up-to-date, nothing to do
			log.V(4).Info("no changes found between user-updated spec and existing spec")
			return nil, nil
		}
		log.V(4).Info("found a diff between the desired spec and the existing privateendpoint", "difference", diff)
	}

	return newPrivateEndpoint, nil
}

func normalizePrivateEndpoint(existingPE, newPrivateEndpoint armnetwork.PrivateEndpoint) armnetwork.PrivateEndpoint {
	normalizedExistingPE := armnetwork.PrivateEndpoint{
		Name:     existingPE.Name,
		Location: existingPE.Location,
		Properties: &armnetwork.PrivateEndpointProperties{
			Subnet: &armnetwork.Subnet{
				ID: existingPE.Properties.Subnet.ID,
				Properties: &armnetwork.SubnetPropertiesFormat{
					PrivateEndpointNetworkPolicies:    newPrivateEndpoint.Properties.Subnet.Properties.PrivateEndpointNetworkPolicies,
					PrivateLinkServiceNetworkPolicies: newPrivateEndpoint.Properties.Subnet.Properties.PrivateLinkServiceNetworkPolicies,
				},
			},
			ApplicationSecurityGroups:  existingPE.Properties.ApplicationSecurityGroups,
			IPConfigurations:           existingPE.Properties.IPConfigurations,
			CustomNetworkInterfaceName: existingPE.Properties.CustomNetworkInterfaceName,
		},
		Tags: existingPE.Tags,
	}
	if existingPE.Properties != nil && existingPE.Properties.Subnet != nil && existingPE.Properties.Subnet.Properties != nil {
		normalizedExistingPE.Properties.Subnet.Properties.PrivateEndpointNetworkPolicies = existingPE.Properties.Subnet.Properties.PrivateEndpointNetworkPolicies
		normalizedExistingPE.Properties.Subnet.Properties.PrivateLinkServiceNetworkPolicies = existingPE.Properties.Subnet.Properties.PrivateLinkServiceNetworkPolicies
	}

	existingPrivateLinkServiceConnections := make([]*armnetwork.PrivateLinkServiceConnection, 0, len(existingPE.Properties.PrivateLinkServiceConnections))
	for _, privateLinkServiceConnection := range existingPE.Properties.PrivateLinkServiceConnections {
		existingPrivateLinkServiceConnections = append(existingPrivateLinkServiceConnections, &armnetwork.PrivateLinkServiceConnection{
			Name: privateLinkServiceConnection.Name,
			Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
				PrivateLinkServiceID: privateLinkServiceConnection.Properties.PrivateLinkServiceID,
				RequestMessage:       privateLinkServiceConnection.Properties.RequestMessage,
				GroupIDs:             privateLinkServiceConnection.Properties.GroupIDs,
			},
		})
	}
	normalizedExistingPE.Properties.PrivateLinkServiceConnections = existingPrivateLinkServiceConnections

	existingManualPrivateLinkServiceConnections := make([]*armnetwork.PrivateLinkServiceConnection, 0, len(existingPE.Properties.ManualPrivateLinkServiceConnections))
	for _, manualPrivateLinkServiceConnection := range existingPE.Properties.ManualPrivateLinkServiceConnections {
		existingManualPrivateLinkServiceConnections = append(existingManualPrivateLinkServiceConnections, &armnetwork.PrivateLinkServiceConnection{
			Name: manualPrivateLinkServiceConnection.Name,
			Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
				PrivateLinkServiceID: manualPrivateLinkServiceConnection.Properties.PrivateLinkServiceID,
				RequestMessage:       manualPrivateLinkServiceConnection.Properties.RequestMessage,
				GroupIDs:             manualPrivateLinkServiceConnection.Properties.GroupIDs,
			},
		})
	}
	normalizedExistingPE.Properties.ManualPrivateLinkServiceConnections = existingManualPrivateLinkServiceConnections

	return normalizedExistingPE
}

// Sort all slices in order to get the same order of elements for both new and existing private endpoints.
func sortSlicesPrivateEndpoint(privateEndpoint armnetwork.PrivateEndpoint) armnetwork.PrivateEndpoint {
	// Sort ManualPrivateLinkServiceConnections
	if privateEndpoint.Properties.ManualPrivateLinkServiceConnections != nil {
		sort.SliceStable(privateEndpoint.Properties.ManualPrivateLinkServiceConnections, func(i, j int) bool {
			return *privateEndpoint.Properties.ManualPrivateLinkServiceConnections[i].Name < *privateEndpoint.Properties.ManualPrivateLinkServiceConnections[j].Name
		})
	}

	// Sort PrivateLinkServiceConnections
	if privateEndpoint.Properties.PrivateLinkServiceConnections != nil {
		sort.SliceStable(privateEndpoint.Properties.PrivateLinkServiceConnections, func(i, j int) bool {
			return *privateEndpoint.Properties.PrivateLinkServiceConnections[i].Name < *privateEndpoint.Properties.PrivateLinkServiceConnections[j].Name
		})
	}

	// Sort IPConfigurations
	if privateEndpoint.Properties.IPConfigurations != nil {
		sort.SliceStable(privateEndpoint.Properties.IPConfigurations, func(i, j int) bool {
			return *privateEndpoint.Properties.IPConfigurations[i].Properties.PrivateIPAddress < *privateEndpoint.Properties.IPConfigurations[j].Properties.PrivateIPAddress
		})
	}

	// Sort ApplicationSecurityGroups
	if privateEndpoint.Properties.ApplicationSecurityGroups != nil {
		sort.SliceStable(privateEndpoint.Properties.ApplicationSecurityGroups, func(i, j int) bool {
			return *privateEndpoint.Properties.ApplicationSecurityGroups[i].Name < *privateEndpoint.Properties.ApplicationSecurityGroups[j].Name
		})
	}

	return privateEndpoint
}
