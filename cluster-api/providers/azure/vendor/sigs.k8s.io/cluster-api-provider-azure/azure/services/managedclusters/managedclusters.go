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
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/async"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/token"
	"sigs.k8s.io/cluster-api-provider-azure/util/reconciler"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	serviceName        = "managedcluster"
	kubeletIdentityKey = "kubeletidentity"

	// The aadResourceID is the application-id used by the server side. The access token accessing AKS clusters need to be issued for this app.
	// Refer: https://azure.github.io/kubelogin/concepts/aks.html?highlight=6dae42f8-4368-4678-94ff-3960e28e3630#azure-kubernetes-service-aad-server
	aadResourceID = "6dae42f8-4368-4678-94ff-3960e28e3630"
)

// ManagedClusterScope defines the scope interface for a managed cluster.
type ManagedClusterScope interface {
	azure.Authorizer
	azure.AsyncStatusUpdater
	ManagedClusterSpec() azure.ResourceSpecGetter
	SetControlPlaneEndpoint(clusterv1.APIEndpoint)
	SetKubeletIdentity(string)
	MakeEmptyKubeConfigSecret() corev1.Secret
	GetAdminKubeconfigData() []byte
	SetAdminKubeconfigData([]byte)
	GetUserKubeconfigData() []byte
	SetUserKubeconfigData([]byte)
	IsAADEnabled() bool
	AreLocalAccountsDisabled() bool
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
		adminKubeConfigData, userKubeConfigData, err := s.ReconcileKubeconfig(ctx, managedClusterSpec)
		if err != nil {
			return errors.Wrap(err, "error while reconciling adminKubeConfigData")
		}

		s.Scope.SetAdminKubeconfigData(adminKubeConfigData)
		s.Scope.SetUserKubeconfigData(userKubeConfigData)

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

// ReconcileKubeconfig will reconcile admin kubeconfig and user kubeconfig.
/*
  Returns the admin kubeconfig and user kubeconfig
  If aad is enabled a user kubeconfig will also get generated and stored in the secret <cluster-name>-kubeconfig-user
  If we disable local accounts for aad clusters we do not have access to admin kubeconfig, hence we need to create
  the admin kubeconfig by authenticating with the user credentials and retrieving the token for kubeconfig.
  The token is used to create the admin kubeconfig.
  The user needs to ensure to provide service principle with admin aad privileges.
*/
func (s *Service) ReconcileKubeconfig(ctx context.Context, managedClusterSpec azure.ResourceSpecGetter) (userKubeConfigData []byte, adminKubeConfigData []byte, err error) {
	if s.Scope.IsAADEnabled() {
		if userKubeConfigData, err = s.GetUserKubeconfigData(ctx, managedClusterSpec); err != nil {
			return nil, nil, errors.Wrap(err, "error while trying to get user kubeconfig")
		}
	}

	if s.Scope.AreLocalAccountsDisabled() {
		userKubeconfigWithToken, err := s.GetUserKubeConfigWithToken(userKubeConfigData, ctx, managedClusterSpec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error while trying to get user kubeconfig with token")
		}
		return userKubeconfigWithToken, userKubeConfigData, nil
	}

	adminKubeConfigData, err = s.GetCredentials(ctx, managedClusterSpec.ResourceGroupName(), managedClusterSpec.ResourceName())
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get credentials for managed cluster")
	}
	return adminKubeConfigData, userKubeConfigData, nil
}

// GetUserKubeconfigData gets user kubeconfig when aad is enabled for the aad clusters.
func (s *Service) GetUserKubeconfigData(ctx context.Context, managedClusterSpec azure.ResourceSpecGetter) ([]byte, error) {
	kubeConfigData, err := s.GetUserCredentials(ctx, managedClusterSpec.ResourceGroupName(), managedClusterSpec.ResourceName())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get credentials for managed cluster")
	}
	return kubeConfigData, nil
}

// GetUserKubeConfigWithToken returns the kubeconfig with user token, for capz to create the target cluster.
func (s *Service) GetUserKubeConfigWithToken(userKubeConfigData []byte, ctx context.Context, managedClusterSpec azure.ResourceSpecGetter) ([]byte, error) {
	tokenClient, err := token.NewClient(s.Scope)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting aad token client")
	}

	token, err := tokenClient.GetAzureActiveDirectoryToken(ctx, aadResourceID)
	if err != nil {
		return nil, errors.Wrap(err, "error while getting aad token for user kubeconfig")
	}

	return s.CreateUserKubeconfigWithToken(token, userKubeConfigData)
}

// CreateUserKubeconfigWithToken gets the kubeconfigdata for authenticating with target cluster.
func (s *Service) CreateUserKubeconfigWithToken(token string, userKubeConfigData []byte) ([]byte, error) {
	config, err := clientcmd.Load(userKubeConfigData)
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to unmarshal new user kubeconfig with token")
	}
	for _, auth := range config.AuthInfos {
		auth.Token = token
		auth.Exec = nil
	}
	kubeconfig, err := clientcmd.Write(*config)
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to marshal new user kubeconfig with token")
	}
	return kubeconfig, nil
}
