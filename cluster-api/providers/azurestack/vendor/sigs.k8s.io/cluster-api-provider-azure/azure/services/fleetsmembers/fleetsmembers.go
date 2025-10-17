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
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/slice"
)

const serviceName = "fleetsmember"

// FleetsMemberScope defines the scope interface for a Fleet host service.
type FleetsMemberScope interface {
	aso.Scope
	AzureFleetsMemberSpec() []azure.ASOResourceSpecGetter[*asocontainerservicev1.FleetsMember]
}

// New creates a new service.
func New(scope FleetsMemberScope) *aso.Service[*asocontainerservicev1.FleetsMember, FleetsMemberScope] {
	svc := aso.NewService[*asocontainerservicev1.FleetsMember, FleetsMemberScope](serviceName, scope)
	svc.ListFunc = list
	svc.Specs = scope.AzureFleetsMemberSpec()
	svc.ConditionType = infrav1.FleetReadyCondition
	return svc
}

func list(ctx context.Context, client client.Client, opts ...client.ListOption) ([]*asocontainerservicev1.FleetsMember, error) {
	list := &asocontainerservicev1.FleetsMemberList{}
	err := client.List(ctx, list, opts...)
	return slice.ToPtrs(list.Items), err
}
