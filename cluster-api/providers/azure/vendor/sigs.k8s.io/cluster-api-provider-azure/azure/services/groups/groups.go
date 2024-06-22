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

package groups

import (
	"context"

	asoresourcesv1 "github.com/Azure/azure-service-operator/v2/api/resources/v1api20200601"
	asoannotations "github.com/Azure/azure-service-operator/v2/pkg/common/annotations"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/slice"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceName is the name of this service.
const ServiceName = "group"

// Service provides operations on Azure resources.
type Service struct {
	Scope GroupScope
	*aso.Service[*asoresourcesv1.ResourceGroup, GroupScope]
}

// GroupScope defines the scope interface for a group service.
type GroupScope interface {
	aso.Scope
	GroupSpecs() []azure.ASOResourceSpecGetter[*asoresourcesv1.ResourceGroup]
}

// New creates a new service.
func New(scope GroupScope) *Service {
	svc := aso.NewService[*asoresourcesv1.ResourceGroup](ServiceName, scope)
	svc.ListFunc = list
	svc.Specs = scope.GroupSpecs()
	svc.ConditionType = infrav1.ResourceGroupReadyCondition
	return &Service{
		Scope:   scope,
		Service: svc,
	}
}

// IsManaged returns true if all resource groups are
// managed and reconciled by ASO, meaning that we can rely on a single resource
// group delete operation as opposed to deleting every individual resource.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	// Unless all resource groups are managed by CAPZ and reconciled by ASO, resources need to be deleted individually.
	for _, spec := range s.Specs {
		managed, err := aso.IsManaged(ctx, s.Scope.GetClient(), spec.ResourceRef(), s.Scope.ASOOwner())
		if err != nil || !managed {
			return managed, err
		}

		// For ASO, "managed" only tells us that we're allowed to delete the ASO
		// resource. We also need to check that deleting the ASO resource will really
		// delete the underlying resource group by checking the ASO reconcile-policy.
		group := spec.ResourceRef()
		groupName := azure.GetNormalizedKubernetesName(group.Name)
		groupNamespace := s.Scope.ASOOwner().GetNamespace()
		err = s.Scope.GetClient().Get(ctx, client.ObjectKey{Namespace: groupNamespace, Name: groupName}, group)
		if err != nil || group.GetAnnotations()[asoannotations.ReconcilePolicy] != string(asoannotations.ReconcilePolicyManage) {
			return false, err
		}
	}
	return true, nil
}

func list(ctx context.Context, client client.Client, opts ...client.ListOption) ([]*asoresourcesv1.ResourceGroup, error) {
	list := &asoresourcesv1.ResourceGroupList{}
	err := client.List(ctx, list, opts...)
	return slice.ToPtrs(list.Items), err
}
