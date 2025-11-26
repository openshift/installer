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

package aksextensions

import (
	"context"

	asokubernetesconfigurationv1 "github.com/Azure/azure-service-operator/v2/api/kubernetesconfiguration/v1api20230501"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
)

// AKSExtensionSpec defines the specification for an AKS Extension.
type AKSExtensionSpec struct {
	Name                    string
	Namespace               string
	AKSAssignedIdentityType infrav1.AKSAssignedIdentity
	ExtensionIdentity       infrav1.ExtensionIdentity
	AutoUpgradeMinorVersion *bool
	ConfigurationSettings   map[string]string
	ExtensionType           *string
	ReleaseTrain            *string
	Version                 *string
	Owner                   string
	OwnerRef                metav1.OwnerReference
	Plan                    *infrav1.ExtensionPlan
	Scope                   infrav1.ExtensionScope
}

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *AKSExtensionSpec) ResourceRef() *asokubernetesconfigurationv1.Extension {
	return &asokubernetesconfigurationv1.Extension{
		ObjectMeta: metav1.ObjectMeta{
			Name:      azure.GetNormalizedKubernetesName(s.Name),
			Namespace: s.Namespace,
		},
	}
}

// Parameters implements azure.ASOResourceSpecGetter.
func (s *AKSExtensionSpec) Parameters(_ context.Context, existingAKSExtension *asokubernetesconfigurationv1.Extension) (parameters *asokubernetesconfigurationv1.Extension, err error) {
	aksExtension := &asokubernetesconfigurationv1.Extension{}
	if existingAKSExtension != nil {
		aksExtension = existingAKSExtension
	}

	aksExtension.Spec = asokubernetesconfigurationv1.Extension_Spec{}
	aksExtension.Spec.AzureName = s.Name
	aksExtension.Spec.AutoUpgradeMinorVersion = s.AutoUpgradeMinorVersion
	aksExtension.Spec.ConfigurationSettings = s.ConfigurationSettings
	aksExtension.Spec.ExtensionType = s.ExtensionType
	aksExtension.Spec.ReleaseTrain = s.ReleaseTrain
	aksExtension.Spec.Version = s.Version
	aksExtension.Spec.Owner = &genruntime.ArbitraryOwnerReference{
		ARMID: s.Owner,
	}

	if s.Plan != nil {
		aksExtension.Spec.Plan = &asokubernetesconfigurationv1.Plan{
			Name:      ptr.To(s.Plan.Name),
			Product:   ptr.To(s.Plan.Product),
			Publisher: ptr.To(s.Plan.Publisher),
			Version:   ptr.To(s.Plan.Version),
		}
	}
	if s.ExtensionIdentity != "" {
		aksExtension.Spec.Identity = &asokubernetesconfigurationv1.Identity{
			Type: (*asokubernetesconfigurationv1.Identity_Type)(ptr.To(s.ExtensionIdentity)),
		}
	}
	if s.AKSAssignedIdentityType != "" {
		aksExtension.Spec.AksAssignedIdentity = &asokubernetesconfigurationv1.Extension_Properties_AksAssignedIdentity_Spec{
			Type: (*asokubernetesconfigurationv1.Extension_Properties_AksAssignedIdentity_Type_Spec)(ptr.To(s.AKSAssignedIdentityType)),
		}
	}
	switch s.Scope.ScopeType {
	case infrav1.ExtensionScopeCluster:
		aksExtension.Spec.Scope = &asokubernetesconfigurationv1.Scope{
			Cluster: &asokubernetesconfigurationv1.ScopeCluster{
				ReleaseNamespace: ptr.To(s.Scope.ReleaseNamespace),
			},
		}
	case infrav1.ExtensionScopeNamespace:
		aksExtension.Spec.Scope = &asokubernetesconfigurationv1.Scope{
			Namespace: &asokubernetesconfigurationv1.ScopeNamespace{
				TargetNamespace: ptr.To(s.Scope.TargetNamespace),
			},
		}
	}

	return aksExtension, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *AKSExtensionSpec) WasManaged(_ *asokubernetesconfigurationv1.Extension) bool {
	// returns always returns true as CAPZ does not support BYO extension.
	return true
}
