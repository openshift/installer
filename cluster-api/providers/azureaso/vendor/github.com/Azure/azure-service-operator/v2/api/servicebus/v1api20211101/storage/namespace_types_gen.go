// Code generated by azure-service-operator-codegen. DO NOT EDIT.
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package storage

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/configmaps"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// +kubebuilder:rbac:groups=servicebus.azure.com,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=servicebus.azure.com,resources={namespaces/status,namespaces/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// Storage version of v1api20211101.Namespace
// Generator information:
// - Generated from: /servicebus/resource-manager/Microsoft.ServiceBus/stable/2021-11-01/namespace-preview.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{namespaceName}
type Namespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Namespace_Spec   `json:"spec,omitempty"`
	Status            Namespace_STATUS `json:"status,omitempty"`
}

var _ conditions.Conditioner = &Namespace{}

// GetConditions returns the conditions of the resource
func (namespace *Namespace) GetConditions() conditions.Conditions {
	return namespace.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (namespace *Namespace) SetConditions(conditions conditions.Conditions) {
	namespace.Status.Conditions = conditions
}

var _ configmaps.Exporter = &Namespace{}

// ConfigMapDestinationExpressions returns the Spec.OperatorSpec.ConfigMapExpressions property
func (namespace *Namespace) ConfigMapDestinationExpressions() []*core.DestinationExpression {
	if namespace.Spec.OperatorSpec == nil {
		return nil
	}
	return namespace.Spec.OperatorSpec.ConfigMapExpressions
}

var _ secrets.Exporter = &Namespace{}

// SecretDestinationExpressions returns the Spec.OperatorSpec.SecretExpressions property
func (namespace *Namespace) SecretDestinationExpressions() []*core.DestinationExpression {
	if namespace.Spec.OperatorSpec == nil {
		return nil
	}
	return namespace.Spec.OperatorSpec.SecretExpressions
}

var _ genruntime.KubernetesResource = &Namespace{}

// AzureName returns the Azure name of the resource
func (namespace *Namespace) AzureName() string {
	return namespace.Spec.AzureName
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2021-11-01"
func (namespace Namespace) GetAPIVersion() string {
	return "2021-11-01"
}

// GetResourceScope returns the scope of the resource
func (namespace *Namespace) GetResourceScope() genruntime.ResourceScope {
	return genruntime.ResourceScopeResourceGroup
}

// GetSpec returns the specification of this resource
func (namespace *Namespace) GetSpec() genruntime.ConvertibleSpec {
	return &namespace.Spec
}

// GetStatus returns the status of this resource
func (namespace *Namespace) GetStatus() genruntime.ConvertibleStatus {
	return &namespace.Status
}

// GetSupportedOperations returns the operations supported by the resource
func (namespace *Namespace) GetSupportedOperations() []genruntime.ResourceOperation {
	return []genruntime.ResourceOperation{
		genruntime.ResourceOperationDelete,
		genruntime.ResourceOperationGet,
		genruntime.ResourceOperationPut,
	}
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.ServiceBus/namespaces"
func (namespace *Namespace) GetType() string {
	return "Microsoft.ServiceBus/namespaces"
}

// NewEmptyStatus returns a new empty (blank) status
func (namespace *Namespace) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &Namespace_STATUS{}
}

// Owner returns the ResourceReference of the owner
func (namespace *Namespace) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(namespace.Spec)
	return namespace.Spec.Owner.AsResourceReference(group, kind)
}

// SetStatus sets the status of this resource
func (namespace *Namespace) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*Namespace_STATUS); ok {
		namespace.Status = *st
		return nil
	}

	// Convert status to required version
	var st Namespace_STATUS
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	namespace.Status = st
	return nil
}

// Hub marks that this Namespace is the hub type for conversion
func (namespace *Namespace) Hub() {}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (namespace *Namespace) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: namespace.Spec.OriginalVersion,
		Kind:    "Namespace",
	}
}

// +kubebuilder:object:root=true
// Storage version of v1api20211101.Namespace
// Generator information:
// - Generated from: /servicebus/resource-manager/Microsoft.ServiceBus/stable/2021-11-01/namespace-preview.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{namespaceName}
type NamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Namespace `json:"items"`
}

// Storage version of v1api20211101.APIVersion
// +kubebuilder:validation:Enum={"2021-11-01"}
type APIVersion string

const APIVersion_Value = APIVersion("2021-11-01")

// Storage version of v1api20211101.Namespace_Spec
type Namespace_Spec struct {
	AlternateName *string `json:"alternateName,omitempty"`

	// AzureName: The name of the resource in Azure. This is often the same as the name of the resource in Kubernetes but it
	// doesn't have to be.
	AzureName        string                 `json:"azureName,omitempty"`
	DisableLocalAuth *bool                  `json:"disableLocalAuth,omitempty"`
	Encryption       *Encryption            `json:"encryption,omitempty"`
	Identity         *Identity              `json:"identity,omitempty"`
	Location         *string                `json:"location,omitempty"`
	OperatorSpec     *NamespaceOperatorSpec `json:"operatorSpec,omitempty"`
	OriginalVersion  string                 `json:"originalVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	// controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	// reference to a resources.azure.com/ResourceGroup resource
	Owner         *genruntime.KnownResourceReference `group:"resources.azure.com" json:"owner,omitempty" kind:"ResourceGroup"`
	PropertyBag   genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
	Sku           *SBSku                             `json:"sku,omitempty"`
	Tags          map[string]string                  `json:"tags,omitempty"`
	ZoneRedundant *bool                              `json:"zoneRedundant,omitempty"`
}

var _ genruntime.ConvertibleSpec = &Namespace_Spec{}

// ConvertSpecFrom populates our Namespace_Spec from the provided source
func (namespace *Namespace_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	if source == namespace {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return source.ConvertSpecTo(namespace)
}

// ConvertSpecTo populates the provided destination from our Namespace_Spec
func (namespace *Namespace_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	if destination == namespace {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return destination.ConvertSpecFrom(namespace)
}

// Storage version of v1api20211101.Namespace_STATUS
type Namespace_STATUS struct {
	AlternateName              *string                            `json:"alternateName,omitempty"`
	Conditions                 []conditions.Condition             `json:"conditions,omitempty"`
	CreatedAt                  *string                            `json:"createdAt,omitempty"`
	DisableLocalAuth           *bool                              `json:"disableLocalAuth,omitempty"`
	Encryption                 *Encryption_STATUS                 `json:"encryption,omitempty"`
	Id                         *string                            `json:"id,omitempty"`
	Identity                   *Identity_STATUS                   `json:"identity,omitempty"`
	Location                   *string                            `json:"location,omitempty"`
	MetricId                   *string                            `json:"metricId,omitempty"`
	Name                       *string                            `json:"name,omitempty"`
	PrivateEndpointConnections []PrivateEndpointConnection_STATUS `json:"privateEndpointConnections,omitempty"`
	PropertyBag                genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
	ProvisioningState          *string                            `json:"provisioningState,omitempty"`
	ServiceBusEndpoint         *string                            `json:"serviceBusEndpoint,omitempty"`
	Sku                        *SBSku_STATUS                      `json:"sku,omitempty"`
	Status                     *string                            `json:"status,omitempty"`
	SystemData                 *SystemData_STATUS                 `json:"systemData,omitempty"`
	Tags                       map[string]string                  `json:"tags,omitempty"`
	Type                       *string                            `json:"type,omitempty"`
	UpdatedAt                  *string                            `json:"updatedAt,omitempty"`
	ZoneRedundant              *bool                              `json:"zoneRedundant,omitempty"`
}

var _ genruntime.ConvertibleStatus = &Namespace_STATUS{}

// ConvertStatusFrom populates our Namespace_STATUS from the provided source
func (namespace *Namespace_STATUS) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	if source == namespace {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return source.ConvertStatusTo(namespace)
}

// ConvertStatusTo populates the provided destination from our Namespace_STATUS
func (namespace *Namespace_STATUS) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	if destination == namespace {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return destination.ConvertStatusFrom(namespace)
}

// Storage version of v1api20211101.Encryption
// Properties to configure Encryption
type Encryption struct {
	KeySource                       *string                `json:"keySource,omitempty"`
	KeyVaultProperties              []KeyVaultProperties   `json:"keyVaultProperties,omitempty"`
	PropertyBag                     genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	RequireInfrastructureEncryption *bool                  `json:"requireInfrastructureEncryption,omitempty"`
}

// Storage version of v1api20211101.Encryption_STATUS
// Properties to configure Encryption
type Encryption_STATUS struct {
	KeySource                       *string                     `json:"keySource,omitempty"`
	KeyVaultProperties              []KeyVaultProperties_STATUS `json:"keyVaultProperties,omitempty"`
	PropertyBag                     genruntime.PropertyBag      `json:"$propertyBag,omitempty"`
	RequireInfrastructureEncryption *bool                       `json:"requireInfrastructureEncryption,omitempty"`
}

// Storage version of v1api20211101.Identity
// Properties to configure User Assigned Identities for Bring your Own Keys
type Identity struct {
	PropertyBag            genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	Type                   *string                       `json:"type,omitempty"`
	UserAssignedIdentities []UserAssignedIdentityDetails `json:"userAssignedIdentities,omitempty"`
}

// Storage version of v1api20211101.Identity_STATUS
// Properties to configure User Assigned Identities for Bring your Own Keys
type Identity_STATUS struct {
	PrincipalId            *string                                `json:"principalId,omitempty"`
	PropertyBag            genruntime.PropertyBag                 `json:"$propertyBag,omitempty"`
	TenantId               *string                                `json:"tenantId,omitempty"`
	Type                   *string                                `json:"type,omitempty"`
	UserAssignedIdentities map[string]UserAssignedIdentity_STATUS `json:"userAssignedIdentities,omitempty"`
}

// Storage version of v1api20211101.NamespaceOperatorSpec
// Details for configuring operator behavior. Fields in this struct are interpreted by the operator directly rather than being passed to Azure
type NamespaceOperatorSpec struct {
	ConfigMapExpressions []*core.DestinationExpression `json:"configMapExpressions,omitempty"`
	PropertyBag          genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	SecretExpressions    []*core.DestinationExpression `json:"secretExpressions,omitempty"`
	Secrets              *NamespaceOperatorSecrets     `json:"secrets,omitempty"`
}

// Storage version of v1api20211101.PrivateEndpointConnection_STATUS
// Properties of the PrivateEndpointConnection.
type PrivateEndpointConnection_STATUS struct {
	Id          *string                `json:"id,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20211101.SBSku
// SKU of the namespace.
type SBSku struct {
	Capacity    *int                   `json:"capacity,omitempty"`
	Name        *string                `json:"name,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Tier        *string                `json:"tier,omitempty"`
}

// Storage version of v1api20211101.SBSku_STATUS
// SKU of the namespace.
type SBSku_STATUS struct {
	Capacity    *int                   `json:"capacity,omitempty"`
	Name        *string                `json:"name,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Tier        *string                `json:"tier,omitempty"`
}

// Storage version of v1api20211101.SystemData_STATUS
// Metadata pertaining to creation and last modification of the resource.
type SystemData_STATUS struct {
	CreatedAt          *string                `json:"createdAt,omitempty"`
	CreatedBy          *string                `json:"createdBy,omitempty"`
	CreatedByType      *string                `json:"createdByType,omitempty"`
	LastModifiedAt     *string                `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string                `json:"lastModifiedBy,omitempty"`
	LastModifiedByType *string                `json:"lastModifiedByType,omitempty"`
	PropertyBag        genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20211101.KeyVaultProperties
// Properties to configure keyVault Properties
type KeyVaultProperties struct {
	Identity    *UserAssignedIdentityProperties `json:"identity,omitempty"`
	KeyName     *string                         `json:"keyName,omitempty"`
	KeyVaultUri *string                         `json:"keyVaultUri,omitempty"`
	KeyVersion  *string                         `json:"keyVersion,omitempty"`
	PropertyBag genruntime.PropertyBag          `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20211101.KeyVaultProperties_STATUS
// Properties to configure keyVault Properties
type KeyVaultProperties_STATUS struct {
	Identity    *UserAssignedIdentityProperties_STATUS `json:"identity,omitempty"`
	KeyName     *string                                `json:"keyName,omitempty"`
	KeyVaultUri *string                                `json:"keyVaultUri,omitempty"`
	KeyVersion  *string                                `json:"keyVersion,omitempty"`
	PropertyBag genruntime.PropertyBag                 `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20211101.NamespaceOperatorSecrets
type NamespaceOperatorSecrets struct {
	Endpoint                  *genruntime.SecretDestination `json:"endpoint,omitempty"`
	PrimaryConnectionString   *genruntime.SecretDestination `json:"primaryConnectionString,omitempty"`
	PrimaryKey                *genruntime.SecretDestination `json:"primaryKey,omitempty"`
	PropertyBag               genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	SecondaryConnectionString *genruntime.SecretDestination `json:"secondaryConnectionString,omitempty"`
	SecondaryKey              *genruntime.SecretDestination `json:"secondaryKey,omitempty"`
}

// Storage version of v1api20211101.UserAssignedIdentity_STATUS
// Recognized Dictionary value.
type UserAssignedIdentity_STATUS struct {
	ClientId    *string                `json:"clientId,omitempty"`
	PrincipalId *string                `json:"principalId,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20211101.UserAssignedIdentityDetails
// Information about the user assigned identity for the resource
type UserAssignedIdentityDetails struct {
	PropertyBag genruntime.PropertyBag       `json:"$propertyBag,omitempty"`
	Reference   genruntime.ResourceReference `armReference:"Reference" json:"reference,omitempty"`
}

// Storage version of v1api20211101.UserAssignedIdentityProperties
type UserAssignedIdentityProperties struct {
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`

	// UserAssignedIdentityReference: ARM ID of user Identity selected for encryption
	UserAssignedIdentityReference *genruntime.ResourceReference `armReference:"UserAssignedIdentity" json:"userAssignedIdentityReference,omitempty"`
}

// Storage version of v1api20211101.UserAssignedIdentityProperties_STATUS
type UserAssignedIdentityProperties_STATUS struct {
	PropertyBag          genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	UserAssignedIdentity *string                `json:"userAssignedIdentity,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Namespace{}, &NamespaceList{})
}
