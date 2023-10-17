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

package managedclusters

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const serviceName = "managedcluster"

const kubeletIdentityKey = "kubeletidentity"

// ManagedClusterScope defines the scope interface for a managed cluster.
type ManagedClusterScope interface {
	azure.Authorizer
	azure.AsyncStatusUpdater
	ManagedClusterSpec() azure.ResourceSpecGetter
	SetControlPlaneEndpoint(clusterv1.APIEndpoint)
	SetKubeletIdentity(string)
	MakeEmptyKubeConfigSecret() corev1.Secret
	GetKubeConfigData() []byte
	SetKubeConfigData([]byte)
	SetOIDCIssuerProfileStatus(*infrav1.OIDCIssuerProfileStatus)
}

// Service provides operations on azure resources.
type Service struct {
	Scope ManagedClusterScope
	async.Reconciler
	CredentialGetter
}

// New creates a new service.
func New(scope ManagedClusterScope) (*Service, error) {
	client, err := newClient(scope)
	if err != nil {
		return nil, err
	}
	return &Service{
		Scope: scope,
		Reconciler: async.New[armcontainerservice.ManagedClustersClientCreateOrUpdateResponse,
			armcontainerservice.ManagedClustersClientDeleteResponse](scope, client, client),
		CredentialGetter: client,
	}, nil
}

// Name returns the service name.
func (s *Service) Name() string {
	return serviceName
}

// Reconcile idempotently creates or updates a managed cluster.
func (s *Service) Reconcile(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "managedclusters.Service.Reconcile")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAKSServiceReconcileTimeout)
	defer cancel()

	managedClusterSpec := s.Scope.ManagedClusterSpec()
	if managedClusterSpec == nil {
		return nil
	}

	result, resultErr := s.CreateOrUpdateResource(ctx, managedClusterSpec, serviceName)
	if resultErr == nil {
		managedCluster, ok := result.(armcontainerservice.ManagedCluster)
		if !ok {
			return errors.Errorf("%T is not an armcontainerservice.ManagedCluster\n%v\n%v", result, result, managedCluster)
		}
		// Update control plane endpoint.
		endpoint := clusterv1.APIEndpoint{
			Host: ptr.Deref(managedCluster.Properties.Fqdn, ""),
			Port: 443,
		}
		s.Scope.SetControlPlaneEndpoint(endpoint)

		// Update kubeconfig data
		// Always fetch credentials in case of rotation
		kubeConfigData, err := s.GetCredentials(ctx, managedClusterSpec.ResourceGroupName(), managedClusterSpec.ResourceName())
		if err != nil {
			return errors.Wrap(err, "failed to get credentials for managed cluster")
		}
		s.Scope.SetKubeConfigData(kubeConfigData)

		// This field gets populated by AKS when not set by the user. Persist AKS's value so for future diffs,
		// the "before" reflects the correct value.
		if id := managedCluster.Properties.IdentityProfile[kubeletIdentityKey]; id != nil && id.ResourceID != nil {
			s.Scope.SetKubeletIdentity(*id.ResourceID)
		}

		s.Scope.SetOIDCIssuerProfileStatus(nil)
		if managedCluster.Properties.OidcIssuerProfile != nil && managedCluster.Properties.OidcIssuerProfile.IssuerURL != nil {
			s.Scope.SetOIDCIssuerProfileStatus(&infrav1.OIDCIssuerProfileStatus{
				IssuerURL: managedCluster.Properties.OidcIssuerProfile.IssuerURL,
			})
		}
	}
	s.Scope.UpdatePutStatus(infrav1.ManagedClusterRunningCondition, serviceName, resultErr)
	return resultErr
}

// Delete deletes the managed cluster.
func (s *Service) Delete(ctx context.Context) error {
	ctx, _, done := tele.StartSpanWithLogger(ctx, "managedclusters.Service.Delete")
	defer done()

	ctx, cancel := context.WithTimeout(ctx, reconciler.DefaultAzureServiceReconcileTimeout)
	defer cancel()

	managedClusterSpec := s.Scope.ManagedClusterSpec()
	if managedClusterSpec == nil {
		return nil
	}

	err := s.DeleteResource(ctx, managedClusterSpec, serviceName)
	s.Scope.UpdateDeleteStatus(infrav1.ManagedClusterRunningCondition, serviceName, err)
	return err
}

// IsManaged returns always returns true as CAPZ does not support BYO managed cluster.
func (s *Service) IsManaged(ctx context.Context) (bool, error) {
	return true, nil
}
