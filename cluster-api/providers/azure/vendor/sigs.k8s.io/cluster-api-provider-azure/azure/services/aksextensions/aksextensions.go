/*
Copyright 2024 The Kubernetes Authors.

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

package aksextensions

import (
	"context"

	asokubernetesconfigurationv1 "github.com/Azure/azure-service-operator/v2/api/kubernetesconfiguration/v1api20230501"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/slice"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const serviceName = "extension"

// AKSExtensionScope defines the scope interface for an AKS extensions service.
type AKSExtensionScope interface {
	azure.ClusterScoper
	aso.Scope
	AKSExtensionSpecs() []azure.ASOResourceSpecGetter[*asokubernetesconfigurationv1.Extension]
}

// Service provides operations on Azure resources.
type Service struct {
	Scope AKSExtensionScope
	*aso.Service[*asokubernetesconfigurationv1.Extension, AKSExtensionScope]
}

// New creates a new service.
func New(scope AKSExtensionScope) *Service {
	svc := aso.NewService[*asokubernetesconfigurationv1.Extension, AKSExtensionScope](serviceName, scope)
	svc.ListFunc = list
	svc.Specs = scope.AKSExtensionSpecs()
	svc.ConditionType = infrav1.AKSExtensionsReadyCondition
	return &Service{
		Scope:   scope,
		Service: svc,
	}
}

func list(ctx context.Context, client client.Client, opts ...client.ListOption) ([]*asokubernetesconfigurationv1.Extension, error) {
	list := &asokubernetesconfigurationv1.ExtensionList{}
	err := client.List(ctx, list, opts...)
	return slice.ToPtrs(list.Items), err
}
