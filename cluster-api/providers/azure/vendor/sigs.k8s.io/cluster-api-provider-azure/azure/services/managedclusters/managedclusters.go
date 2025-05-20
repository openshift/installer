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
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	asocontainerservicev1hub "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20240901/storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/secret"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
)

const (
	serviceName        = "managedcluster"
	kubeletIdentityKey = "kubeletidentity"

	// The aadResourceID is the application-id used by the server side. The access token accessing AKS clusters need to be issued for this app.
	// Refer: https://azure.github.io/kubelogin/concepts/aks.html?highlight=6dae42f8-4368-4678-94ff-3960e28e3630#azure-kubernetes-service-aad-server
	aadResourceID = "6dae42f8-4368-4678-94ff-3960e28e3630"

	// oidcIssuerProfileUrl is a constant representing the key name for the oidc-issuer-profile-url config map.
	oidcIssuerProfileURL = "oidc-issuer-profile-url"
)

// ManagedClusterScope defines the scope interface for a managed cluster.
type ManagedClusterScope interface {
	aso.Scope
	azure.Authorizer
	ManagedClusterSpec() azure.ASOResourceSpecGetter[genruntime.MetaObject]
	SetControlPlaneEndpoint(clusterv1.APIEndpoint)
	MakeEmptyKubeConfigSecret() corev1.Secret
	GetAdminKubeconfigData() []byte
	SetAdminKubeconfigData([]byte)
	GetUserKubeconfigData() []byte
	SetUserKubeconfigData([]byte)
	IsAADEnabled() bool
	AreLocalAccountsDisabled() bool
	SetOIDCIssuerProfileStatus(*infrav1.OIDCIssuerProfileStatus)
	MakeClusterCA() *corev1.Secret
	StoreClusterInfo(context.Context, []byte) error
	SetAutoUpgradeVersionStatus(version string)
	SetVersionStatus(version string)
	IsManagedVersionUpgrade() bool
}

// New creates a new service.
func New(scope ManagedClusterScope) *aso.Service[genruntime.MetaObject, ManagedClusterScope] {
	// genruntime.MetaObject is used here instead of an *asocontainerservicev1.ManagedCluster to better
	// facilitate returning different API versions.
	svc := aso.NewService[genruntime.MetaObject](serviceName, scope)
	svc.Specs = []azure.ASOResourceSpecGetter[genruntime.MetaObject]{scope.ManagedClusterSpec()}
	svc.ConditionType = infrav1.ManagedClusterRunningCondition
	svc.PostCreateOrUpdateResourceHook = postCreateOrUpdateResourceHook
	return svc
}

func postCreateOrUpdateResourceHook(ctx context.Context, scope ManagedClusterScope, obj genruntime.MetaObject, err error) error {
	if err != nil {
		return err
	}

	managedCluster := &asocontainerservicev1hub.ManagedCluster{}
	if err := obj.(conversion.Convertible).ConvertTo(managedCluster); err != nil {
		return err
	}

	// Update control plane endpoint.
	endpoint := clusterv1.APIEndpoint{
		Host: ptr.Deref(managedCluster.Status.Fqdn, ""),
		Port: 443,
	}
	if managedCluster.Status.ApiServerAccessProfile != nil &&
		ptr.Deref(managedCluster.Status.ApiServerAccessProfile.EnablePrivateCluster, false) &&
		!ptr.Deref(managedCluster.Status.ApiServerAccessProfile.EnablePrivateClusterPublicFQDN, false) {
		endpoint = clusterv1.APIEndpoint{
			Host: ptr.Deref(managedCluster.Status.PrivateFQDN, ""),
			Port: 443,
		}
	}
	scope.SetControlPlaneEndpoint(endpoint)

	// Update kubeconfig data
	// Always fetch credentials in case of rotation
	adminKubeConfigData, userKubeConfigData, err := reconcileKubeconfig(ctx, scope, managedCluster.Namespace)
	if err != nil {
		return errors.Wrap(err, "error while reconciling kubeconfigs")
	}
	scope.SetAdminKubeconfigData(adminKubeConfigData)
	scope.SetUserKubeconfigData(userKubeConfigData)

	scope.SetOIDCIssuerProfileStatus(nil)
	if managedCluster.Status.OidcIssuerProfile != nil && managedCluster.Status.OidcIssuerProfile.IssuerURL != nil {
		scope.SetOIDCIssuerProfileStatus(&infrav1.OIDCIssuerProfileStatus{
			IssuerURL: managedCluster.Status.OidcIssuerProfile.IssuerURL,
		})
	}
	if managedCluster.Status.CurrentKubernetesVersion != nil {
		currentKubernetesVersion := fmt.Sprintf("v%s", *managedCluster.Status.CurrentKubernetesVersion)
		scope.SetVersionStatus(currentKubernetesVersion)
		if scope.IsManagedVersionUpgrade() {
			scope.SetAutoUpgradeVersionStatus(currentKubernetesVersion)
		}
	}

	return nil
}

// reconcileKubeconfig will reconcile admin kubeconfig and user kubeconfig.
/*
  Returns the admin kubeconfig and user kubeconfig
  If AAD is enabled a user kubeconfig will also get generated and stored in the secret <cluster-name>-kubeconfig-user
  If we disable local accounts for AAD clusters we do not have access to admin kubeconfig, hence we need to create
  the admin kubeconfig by authenticating with the user credentials and retrieving the token for kubeconfig.
  The token is used to create the admin kubeconfig.
  The user needs to ensure to provide service principal with admin AAD privileges.
*/
func reconcileKubeconfig(ctx context.Context, scope ManagedClusterScope, namespace string) (adminKubeConfigData []byte, userKubeConfigData []byte, err error) {
	if scope.IsAADEnabled() {
		if userKubeConfigData, err = getUserKubeconfigData(ctx, scope, namespace); err != nil {
			return nil, nil, errors.Wrap(err, "error while trying to get user kubeconfig")
		}
	}

	if scope.AreLocalAccountsDisabled() {
		userKubeconfigWithToken, err := getUserKubeConfigWithToken(ctx, userKubeConfigData, scope)
		if err != nil {
			return nil, nil, errors.Wrap(err, "error while trying to get user kubeconfig with token")
		}
		return userKubeconfigWithToken, userKubeConfigData, nil
	}

	asoSecret := &corev1.Secret{}
	err = scope.GetClient().Get(
		ctx,
		client.ObjectKey{
			Namespace: namespace,
			Name:      adminKubeconfigSecretName(scope.ClusterName()),
		},
		asoSecret,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get ASO admin kubeconfig secret")
	}
	adminKubeConfigData = asoSecret.Data[secret.KubeconfigDataName]
	return adminKubeConfigData, userKubeConfigData, nil
}

// getUserKubeconfigData gets user kubeconfig when aad is enabled for the aad clusters.
func getUserKubeconfigData(ctx context.Context, scope ManagedClusterScope, namespace string) ([]byte, error) {
	asoSecret := &corev1.Secret{}
	err := scope.GetClient().Get(
		ctx,
		client.ObjectKey{
			Namespace: namespace,
			Name:      userKubeconfigSecretName(scope.ClusterName()),
		},
		asoSecret,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ASO user kubeconfig secret")
	}
	kubeConfigData := asoSecret.Data[secret.KubeconfigDataName]
	return kubeConfigData, nil
}

// getUserKubeConfigWithToken returns the kubeconfig with user token, for capz to create the target cluster.
func getUserKubeConfigWithToken(ctx context.Context, userKubeConfigData []byte, auth azure.Authorizer) ([]byte, error) {
	token, err := auth.Token().GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{aadResourceID + "/.default"}})
	if err != nil {
		return nil, errors.Wrap(err, "error while getting aad token for user kubeconfig")
	}
	config, err := clientcmd.Load(userKubeConfigData)
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to unmarshal new user kubeconfig with token")
	}
	for _, auth := range config.AuthInfos {
		auth.Token = token.Token
		auth.Exec = nil
	}
	kubeconfig, err := clientcmd.Write(*config)
	if err != nil {
		return nil, errors.Wrap(err, "error while trying to marshal new user kubeconfig with token")
	}
	return kubeconfig, nil
}
