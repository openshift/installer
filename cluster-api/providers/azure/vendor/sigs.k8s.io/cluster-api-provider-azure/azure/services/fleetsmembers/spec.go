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

package fleetsmembers

import (
	"context"

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20230315preview"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// AzureFleetsMemberSpec defines the specification for an Azure Fleets Member.
type AzureFleetsMemberSpec struct {
	Name                 string
	ClusterName          string
	ClusterResourceGroup string
	Group                string
	SubscriptionID       string
	ManagerName          string
	ManagerResourceGroup string
}

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *AzureFleetsMemberSpec) ResourceRef() *asocontainerservicev1.FleetsMember {
	return &asocontainerservicev1.FleetsMember{
		ObjectMeta: metav1.ObjectMeta{
			Name: s.Name,
		},
	}
}

// Parameters implements azure.ASOResourceSpecGetter.
func (s *AzureFleetsMemberSpec) Parameters(ctx context.Context, existingFleetsMember *asocontainerservicev1.FleetsMember) (parameters *asocontainerservicev1.FleetsMember, err error) {
	fleetsMember := &asocontainerservicev1.FleetsMember{}
	if existingFleetsMember != nil {
		fleetsMember = existingFleetsMember
	}

	fleetsMember.Spec.AzureName = s.Name
	fleetsMember.Spec.Owner = &genruntime.KnownResourceReference{
		ARMID: azure.FleetID(s.SubscriptionID, s.ManagerResourceGroup, s.ManagerName),
	}
	fleetsMember.Spec.Group = ptr.To(s.Group)
	fleetsMember.Spec.ClusterResourceReference = &genruntime.ResourceReference{
		ARMID: azure.ManagedClusterID(s.SubscriptionID, s.ClusterResourceGroup, s.ClusterName),
	}

	return fleetsMember, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *AzureFleetsMemberSpec) WasManaged(resource *asocontainerservicev1.FleetsMember) bool {
	// returns false because previous versions of CAPZ did not support Fleets.
	return false
}
