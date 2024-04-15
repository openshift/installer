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

package scope

import (
	"context"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/go-autorest/autorest"
	"github.com/jongio/azidext/go/azidext"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AzureSecretKey is the value for they client secret key.
const AzureSecretKey = "clientSecret"

// CredentialsProvider defines the behavior for azure identity based credential providers.
type CredentialsProvider interface {
	GetAuthorizer(ctx context.Context, tokenCredential azcore.TokenCredential, tokenAudience string) (autorest.Authorizer, error)
	GetClientID() string
	GetClientSecret(ctx context.Context) (string, error)
	GetTenantID() string
	GetTokenCredential(ctx context.Context, resourceManagerEndpoint, activeDirectoryEndpoint, tokenAudience string) (azcore.TokenCredential, error)
}

// AzureCredentialsProvider represents a credential provider with azure cluster identity.
type AzureCredentialsProvider struct {
	Client   client.Client
	Identity *infrav1.AzureClusterIdentity
}

// AzureClusterCredentialsProvider wraps AzureCredentialsProvider with AzureCluster.
type AzureClusterCredentialsProvider struct {
	AzureCredentialsProvider
	AzureCluster *infrav1.AzureCluster
}

// ManagedControlPlaneCredentialsProvider wraps AzureCredentialsProvider with AzureManagedControlPlane.
type ManagedControlPlaneCredentialsProvider struct {
	AzureCredentialsProvider
	AzureManagedControlPlane *infrav1.AzureManagedControlPlane
}

var _ CredentialsProvider = (*AzureClusterCredentialsProvider)(nil)
var _ CredentialsProvider = (*ManagedControlPlaneCredentialsProvider)(nil)

// NewAzureClusterCredentialsProvider creates a new AzureClusterCredentialsProvider from the supplied inputs.
func NewAzureClusterCredentialsProvider(ctx context.Context, kubeClient client.Client, azureCluster *infrav1.AzureCluster) (*AzureClusterCredentialsProvider, error) {
	if azureCluster.Spec.IdentityRef == nil {
		return nil, errors.New("failed to generate new AzureClusterCredentialsProvider from empty identityName")
	}

	ref := azureCluster.Spec.IdentityRef
	// if the namespace isn't specified then assume it's in the same namespace as the AzureCluster
	namespace := ref.Namespace
	if namespace == "" {
		namespace = azureCluster.Namespace
	}
	identity := &infrav1.AzureClusterIdentity{}
	key := client.ObjectKey{Name: ref.Name, Namespace: namespace}
	if err := kubeClient.Get(ctx, key, identity); err != nil {
		return nil, errors.Errorf("failed to retrieve AzureClusterIdentity external object %q/%q: %v", key.Namespace, key.Name, err)
	}

	return &AzureClusterCredentialsProvider{
		AzureCredentialsProvider{
			Client:   kubeClient,
			Identity: identity,
		},
		azureCluster,
	}, nil
}

// GetAuthorizer returns an Azure authorizer based on the provided azure identity. It delegates to AzureCredentialsProvider with AzureCluster metadata.
func (p *AzureClusterCredentialsProvider) GetAuthorizer(ctx context.Context, tokenCredential azcore.TokenCredential, tokenAudience string) (autorest.Authorizer, error) {
	return p.AzureCredentialsProvider.GetAuthorizer(ctx, tokenCredential, tokenAudience)
}

// GetTokenCredential returns an Azure TokenCredential based on the provided azure identity.
func (p *AzureClusterCredentialsProvider) GetTokenCredential(ctx context.Context, resourceManagerEndpoint, activeDirectoryEndpoint, tokenAudience string) (azcore.TokenCredential, error) {
	return p.AzureCredentialsProvider.GetTokenCredential(ctx, resourceManagerEndpoint, activeDirectoryEndpoint, tokenAudience, p.AzureCluster.ObjectMeta)
}

// NewManagedControlPlaneCredentialsProvider creates a new ManagedControlPlaneCredentialsProvider from the supplied inputs.
func NewManagedControlPlaneCredentialsProvider(ctx context.Context, kubeClient client.Client, managedControlPlane *infrav1.AzureManagedControlPlane) (*ManagedControlPlaneCredentialsProvider, error) {
	if managedControlPlane.Spec.IdentityRef == nil {
		return nil, errors.New("failed to generate new ManagedControlPlaneCredentialsProvider from empty identityName")
	}

	ref := managedControlPlane.Spec.IdentityRef
	// if the namespace isn't specified then assume it's in the same namespace as the AzureManagedControlPlane
	namespace := ref.Namespace
	if namespace == "" {
		namespace = managedControlPlane.Namespace
	}
	identity := &infrav1.AzureClusterIdentity{}
	key := client.ObjectKey{Name: ref.Name, Namespace: namespace}
	if err := kubeClient.Get(ctx, key, identity); err != nil {
		return nil, errors.Errorf("failed to retrieve AzureClusterIdentity external object %q/%q: %v", key.Namespace, key.Name, err)
	}

	return &ManagedControlPlaneCredentialsProvider{
		AzureCredentialsProvider{
			Client:   kubeClient,
			Identity: identity,
		},
		managedControlPlane,
	}, nil
}

// GetAuthorizer returns an Azure authorizer based on the provided azure identity. It delegates to AzureCredentialsProvider with AzureManagedControlPlane metadata.
func (p *ManagedControlPlaneCredentialsProvider) GetAuthorizer(ctx context.Context, tokenCredential azcore.TokenCredential, tokenAudience string) (autorest.Authorizer, error) {
	return p.AzureCredentialsProvider.GetAuthorizer(ctx, tokenCredential, tokenAudience)
}

// GetTokenCredential returns an Azure TokenCredential based on the provided azure identity.
func (p *ManagedControlPlaneCredentialsProvider) GetTokenCredential(ctx context.Context, resourceManagerEndpoint, activeDirectoryEndpoint, tokenAudience string) (azcore.TokenCredential, error) {
	return p.AzureCredentialsProvider.GetTokenCredential(ctx, resourceManagerEndpoint, activeDirectoryEndpoint, tokenAudience, p.AzureManagedControlPlane.ObjectMeta)
}

// GetTokenCredential returns an Azure TokenCredential based on the provided azure identity.
func (p *AzureCredentialsProvider) GetTokenCredential(ctx context.Context, resourceManagerEndpoint, activeDirectoryEndpoint, tokenAudience string, clusterMeta metav1.ObjectMeta) (azcore.TokenCredential, error) {
	ctx, log, done := tele.StartSpanWithLogger(ctx, "azure.scope.AzureCredentialsProvider.GetTokenCredential")
	defer done()

	var authErr error
	var cred azcore.TokenCredential

	switch p.Identity.Spec.Type {
	case infrav1.WorkloadIdentity:
		azwiCredOptions, err := NewWorkloadIdentityCredentialOptions().
			WithTenantID(p.Identity.Spec.TenantID).
			WithClientID(p.Identity.Spec.ClientID).
			WithDefaults()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to setup azwi options for identity %s", p.Identity.Name)
		}
		cred, authErr = NewWorkloadIdentityCredential(azwiCredOptions)

	case infrav1.ManualServicePrincipal:
		log.Info("Identity type ManualServicePrincipal is deprecated and will be removed in a future release. See https://capz.sigs.k8s.io/topics/identities to find a supported identity type.")
		fallthrough
	case infrav1.ServicePrincipal:
		clientSecret, err := p.GetClientSecret(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get client secret")
		}
		options := azidentity.ClientSecretCredentialOptions{
			ClientOptions: azcore.ClientOptions{
				Cloud: cloud.Configuration{
					ActiveDirectoryAuthorityHost: activeDirectoryEndpoint,
					Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
						cloud.ResourceManager: {
							Audience: tokenAudience,
							Endpoint: resourceManagerEndpoint,
						},
					},
				},
			},
		}
		cred, authErr = azidentity.NewClientSecretCredential(p.GetTenantID(), p.Identity.Spec.ClientID, clientSecret, &options)

	case infrav1.ServicePrincipalCertificate:
		clientSecret, err := p.GetClientSecret(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get client secret")
		}
		certs, key, err := azidentity.ParseCertificates([]byte(clientSecret), nil)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse certificate data")
		}
		cred, authErr = azidentity.NewClientCertificateCredential(p.GetTenantID(), p.Identity.Spec.ClientID, certs, key, nil)

	case infrav1.UserAssignedMSI:
		options := azidentity.ManagedIdentityCredentialOptions{
			ID: azidentity.ClientID(p.Identity.Spec.ClientID),
		}
		cred, authErr = azidentity.NewManagedIdentityCredential(&options)

	default:
		return nil, errors.Errorf("identity type %s not supported", p.Identity.Spec.Type)
	}

	if authErr != nil {
		return nil, errors.Errorf("failed to create credential: %v", authErr)
	}

	return cred, nil
}

// GetAuthorizer returns an Azure authorizer based on the provided azure identity, cluster metadata, and tokenCredential.
func (p *AzureCredentialsProvider) GetAuthorizer(ctx context.Context, cred azcore.TokenCredential, tokenAudience string) (autorest.Authorizer, error) {
	// We must use TokenAudience for StackCloud, otherwise we get an
	// AADSTS500011 error from the API
	scope := tokenAudience
	if !strings.HasSuffix(scope, "/.default") {
		scope += "/.default"
	}
	authorizer := azidext.NewTokenCredentialAdapter(cred, []string{scope})
	return authorizer, nil
}

// GetClientID returns the Client ID associated with the AzureCredentialsProvider's Identity.
func (p *AzureCredentialsProvider) GetClientID() string {
	return p.Identity.Spec.ClientID
}

// GetClientSecret returns the Client Secret associated with the AzureCredentialsProvider's Identity.
// NOTE: this only works if the Identity references a Service Principal Client Secret.
// If using another type of credentials, such a Certificate, we return an empty string.
func (p *AzureCredentialsProvider) GetClientSecret(ctx context.Context) (string, error) {
	if p.hasClientSecret() {
		secretRef := p.Identity.Spec.ClientSecret
		key := types.NamespacedName{
			Namespace: secretRef.Namespace,
			Name:      secretRef.Name,
		}
		secret := &corev1.Secret{}

		if err := p.Client.Get(ctx, key, secret); err != nil {
			return "", errors.Wrap(err, "Unable to fetch ClientSecret")
		}
		return string(secret.Data[AzureSecretKey]), nil
	}
	return "", nil
}

// GetTenantID returns the Tenant ID associated with the AzureCredentialsProvider's Identity.
func (p *AzureCredentialsProvider) GetTenantID() string {
	return p.Identity.Spec.TenantID
}

// hasClientSecret returns true if the identity has a Service Principal Client Secret.
// This does not include managed identities.
func (p *AzureCredentialsProvider) hasClientSecret() bool {
	switch p.Identity.Spec.Type {
	case infrav1.ServicePrincipal, infrav1.ManualServicePrincipal, infrav1.ServicePrincipalCertificate:
		return true
	default:
		return false
	}
}

// IsClusterNamespaceAllowed indicates if the cluster namespace is allowed.
func IsClusterNamespaceAllowed(ctx context.Context, k8sClient client.Client, allowedNamespaces *infrav1.AllowedNamespaces, namespace string) bool {
	if allowedNamespaces == nil {
		return false
	}

	// empty value matches with all namespaces
	if reflect.DeepEqual(*allowedNamespaces, infrav1.AllowedNamespaces{}) {
		return true
	}

	for _, v := range allowedNamespaces.NamespaceList {
		if v == namespace {
			return true
		}
	}

	// Check if clusterNamespace is in the namespaces selected by the identity's allowedNamespaces selector.
	namespaces := &corev1.NamespaceList{}
	selector, err := metav1.LabelSelectorAsSelector(allowedNamespaces.Selector)
	if err != nil {
		return false
	}

	// If a Selector has a nil or empty selector, it should match nothing.
	if selector.Empty() {
		return false
	}

	if err := k8sClient.List(ctx, namespaces, client.MatchingLabelsSelector{Selector: selector}); err != nil {
		return false
	}

	for _, n := range namespaces.Items {
		if n.Name == namespace {
			return true
		}
	}

	return false
}
