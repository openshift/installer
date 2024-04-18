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
	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
)

// ServiceName is the name of this service.
const ServiceName = "privateendpoints"

// PrivateEndpointScope defines the scope interface for a private endpoint.
type PrivateEndpointScope interface {
	aso.Scope
	PrivateEndpointSpecs() []azure.ASOResourceSpecGetter[*asonetworkv1.PrivateEndpoint]
}

// New creates a new service.
func New(scope PrivateEndpointScope) *aso.Service[*asonetworkv1.PrivateEndpoint, PrivateEndpointScope] {
	svc := aso.NewService[*asonetworkv1.PrivateEndpoint, PrivateEndpointScope](ServiceName, scope)
	svc.ConditionType = infrav1.PrivateEndpointsReadyCondition
	svc.Specs = scope.PrivateEndpointSpecs()
	return svc
}
