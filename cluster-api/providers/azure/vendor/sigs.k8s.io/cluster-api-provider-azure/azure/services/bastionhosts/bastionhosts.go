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

package bastionhosts

import (
	"context"

	asonetworkv1 "github.com/Azure/azure-service-operator/v2/api/network/v1api20220701"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/slice"
)

const serviceName = "bastionhosts"

// BastionScope defines the scope interface for a bastion host service.
type BastionScope interface {
	aso.Scope
	AzureBastionSpec() azure.ASOResourceSpecGetter[*asonetworkv1.BastionHost]
}

// New creates a new service.
func New(scope BastionScope) *aso.Service[*asonetworkv1.BastionHost, BastionScope] {
	svc := aso.NewService[*asonetworkv1.BastionHost, BastionScope](serviceName, scope)
	spec := scope.AzureBastionSpec()
	if spec != nil {
		svc.Specs = []azure.ASOResourceSpecGetter[*asonetworkv1.BastionHost]{spec}
	}
	svc.ListFunc = list
	svc.ConditionType = infrav1.BastionHostReadyCondition
	return svc
}

func list(ctx context.Context, client client.Client, opts ...client.ListOption) ([]*asonetworkv1.BastionHost, error) {
	list := &asonetworkv1.BastionHostList{}
	err := client.List(ctx, list, opts...)
	return slice.ToPtrs(list.Items), err
}
