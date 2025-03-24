/*
Copyright 2018 The Kubernetes Authors.

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

package azure

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// Reconciler is a generic interface for a controller reconciler which has Reconcile and Delete methods.
type Reconciler interface {
	Reconcile(ctx context.Context) error
	Delete(ctx context.Context) error
}

// Pauser may be implemented for a ServiceReconciler that requires additional work to stop reconciliation.
type Pauser interface {
	Pause(context.Context) error
}

// ServiceReconciler is an Azure service reconciler which can reconcile an Azure service.
type ServiceReconciler interface {
	Name() string
	Reconciler
}

// Authorizer is an interface which can get details such as subscription ID, base URI, and token
// for authorizing to an Azure service.
type Authorizer interface {
	SubscriptionID() string
	ClientID() string
	ClientSecret() string
	CloudEnvironment() string
	TenantID() string
	BaseURI() string
	HashKey() string
	Token() azcore.TokenCredential
}

// NetworkDescriber is an interface which can get common Azure Cluster Networking information.
type NetworkDescriber interface {
	Vnet() *infrav1.VnetSpec
	IsVnetManaged() bool
	ControlPlaneSubnet() infrav1.SubnetSpec
	Subnets() infrav1.Subnets
	Subnet(string) infrav1.SubnetSpec
	NodeSubnets() []infrav1.SubnetSpec
	SetSubnet(infrav1.SubnetSpec)
	IsIPv6Enabled() bool
	ControlPlaneRouteTable() infrav1.RouteTable
	APIServerLB() *infrav1.LoadBalancerSpec
	APIServerLBName() string
	APIServerLBPoolName() string
	IsAPIServerPrivate() bool
	GetPrivateDNSZoneName() string
	OutboundLBName(string) string
	OutboundPoolName(string) string
}

// ClusterDescriber is an interface which can get common Azure Cluster information.
type ClusterDescriber interface {
	Authorizer
	ResourceGroup() string
	NodeResourceGroup() string
	ClusterName() string
	Location() string
	ExtendedLocation() *infrav1.ExtendedLocationSpec
	ExtendedLocationName() string
	ExtendedLocationType() string
	AdditionalTags() infrav1.Tags
	AvailabilitySetEnabled() bool
	CloudProviderConfigOverrides() *infrav1.CloudProviderConfigOverrides
	FailureDomains() []*string
}

// AsyncStatusUpdater is an interface used to keep track of long running operations in Status that has Conditions and Futures.
type AsyncStatusUpdater interface {
	SetLongRunningOperationState(*infrav1.Future)
	GetLongRunningOperationState(string, string, string) *infrav1.Future
	DeleteLongRunningOperationState(string, string, string)
	UpdatePutStatus(clusterv1.ConditionType, string, error)
	UpdateDeleteStatus(clusterv1.ConditionType, string, error)
	UpdatePatchStatus(clusterv1.ConditionType, string, error)
	AsyncReconciler
}

// AsyncReconciler is an interface used to get the default timeouts and requeue time for a reconciler that reconciles services asynchronously.
type AsyncReconciler interface {
	DefaultedAzureCallTimeout() time.Duration
	DefaultedAzureServiceReconcileTimeout() time.Duration
	DefaultedReconcilerRequeue() time.Duration
}

// ClusterScoper combines the ClusterDescriber and NetworkDescriber interfaces.
type ClusterScoper interface {
	ClusterDescriber
	NetworkDescriber
	AsyncStatusUpdater
	GetClient() client.Client
	GetDeletionTimestamp() *metav1.Time
}

// ManagedClusterScoper defines the interface for ManagedClusterScope.
type ManagedClusterScoper interface {
	ClusterDescriber
	NodeResourceGroup() string
	AsyncReconciler
}

// ResourceSpecGetter is an interface for getting all the required information to create/update/delete an Azure resource.
type ResourceSpecGetter interface {
	// ResourceName returns the name of the resource.
	ResourceName() string
	// OwnerResourceName returns the name of the resource that owns the resource
	// in the case that the resource is an Azure subresource.
	OwnerResourceName() string
	// ResourceGroupName returns the name of the resource group the resource is in.
	ResourceGroupName() string
	// Parameters takes the existing resource and returns the desired parameters of the resource.
	// If the resource does not exist, or we do not care about existing parameters to update the resource, existing should be nil.
	// If no update is needed on the resource, Parameters should return nil.
	Parameters(ctx context.Context, existing interface{}) (params interface{}, err error)
}

// ResourceSpecGetterWithHeaders is a ResourceSpecGetter that can return custom headers to be added to API calls.
type ResourceSpecGetterWithHeaders interface {
	ResourceSpecGetter
	// CustomHeaders returns the headers that should be added to Azure API calls.
	CustomHeaders() map[string]string
}

// ASOResourceSpecGetter is an interface for getting all the required information to create/update/delete an Azure resource.
type ASOResourceSpecGetter[T genruntime.MetaObject] interface {
	// ResourceRef returns a concrete, named ASO resource type to facilitate a
	// strongly-typed GET. Namespace is not read if set here and is instead
	// derived from OwnerReferences.
	ResourceRef() T
	// Parameters returns a modified object if it points to a non-nil resource.
	// Otherwise it returns an unmodified object if no updates are needed.
	Parameters(ctx context.Context, existing T) (T, error)
	// WasManaged returns whether or not the given resource was managed by a
	// non-ASO-backed CAPZ and should be considered eligible for adoption.
	WasManaged(T) bool
}

// CredentialCache caches azcore.TokenCredentials.
type CredentialCache interface {
	GetOrStoreClientSecret(tenantID, clientID, clientSecret string, opts *azidentity.ClientSecretCredentialOptions) (azcore.TokenCredential, error)
	GetOrStoreClientCert(tenantID, clientID string, cert, certPassword []byte, opts *azidentity.ClientCertificateCredentialOptions) (azcore.TokenCredential, error)
	GetOrStoreManagedIdentity(opts *azidentity.ManagedIdentityCredentialOptions) (azcore.TokenCredential, error)
	GetOrStoreWorkloadIdentity(opts *azidentity.WorkloadIdentityCredentialOptions) (azcore.TokenCredential, error)
}
