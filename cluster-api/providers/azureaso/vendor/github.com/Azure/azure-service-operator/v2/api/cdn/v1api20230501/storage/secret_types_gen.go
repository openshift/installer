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

// +kubebuilder:rbac:groups=cdn.azure.com,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cdn.azure.com,resources={secrets/status,secrets/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// Storage version of v1api20230501.Secret
// Generator information:
// - Generated from: /cdn/resource-manager/Microsoft.Cdn/stable/2023-05-01/afdx.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/secrets/{secretName}
type Secret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              Secret_Spec   `json:"spec,omitempty"`
	Status            Secret_STATUS `json:"status,omitempty"`
}

var _ conditions.Conditioner = &Secret{}

// GetConditions returns the conditions of the resource
func (secret *Secret) GetConditions() conditions.Conditions {
	return secret.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (secret *Secret) SetConditions(conditions conditions.Conditions) {
	secret.Status.Conditions = conditions
}

var _ configmaps.Exporter = &Secret{}

// ConfigMapDestinationExpressions returns the Spec.OperatorSpec.ConfigMapExpressions property
func (secret *Secret) ConfigMapDestinationExpressions() []*core.DestinationExpression {
	if secret.Spec.OperatorSpec == nil {
		return nil
	}
	return secret.Spec.OperatorSpec.ConfigMapExpressions
}

var _ secrets.Exporter = &Secret{}

// SecretDestinationExpressions returns the Spec.OperatorSpec.SecretExpressions property
func (secret *Secret) SecretDestinationExpressions() []*core.DestinationExpression {
	if secret.Spec.OperatorSpec == nil {
		return nil
	}
	return secret.Spec.OperatorSpec.SecretExpressions
}

var _ genruntime.KubernetesResource = &Secret{}

// AzureName returns the Azure name of the resource
func (secret *Secret) AzureName() string {
	return secret.Spec.AzureName
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2023-05-01"
func (secret Secret) GetAPIVersion() string {
	return "2023-05-01"
}

// GetResourceScope returns the scope of the resource
func (secret *Secret) GetResourceScope() genruntime.ResourceScope {
	return genruntime.ResourceScopeResourceGroup
}

// GetSpec returns the specification of this resource
func (secret *Secret) GetSpec() genruntime.ConvertibleSpec {
	return &secret.Spec
}

// GetStatus returns the status of this resource
func (secret *Secret) GetStatus() genruntime.ConvertibleStatus {
	return &secret.Status
}

// GetSupportedOperations returns the operations supported by the resource
func (secret *Secret) GetSupportedOperations() []genruntime.ResourceOperation {
	return []genruntime.ResourceOperation{
		genruntime.ResourceOperationDelete,
		genruntime.ResourceOperationGet,
		genruntime.ResourceOperationPut,
	}
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.Cdn/profiles/secrets"
func (secret *Secret) GetType() string {
	return "Microsoft.Cdn/profiles/secrets"
}

// NewEmptyStatus returns a new empty (blank) status
func (secret *Secret) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &Secret_STATUS{}
}

// Owner returns the ResourceReference of the owner
func (secret *Secret) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(secret.Spec)
	return secret.Spec.Owner.AsResourceReference(group, kind)
}

// SetStatus sets the status of this resource
func (secret *Secret) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*Secret_STATUS); ok {
		secret.Status = *st
		return nil
	}

	// Convert status to required version
	var st Secret_STATUS
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	secret.Status = st
	return nil
}

// Hub marks that this Secret is the hub type for conversion
func (secret *Secret) Hub() {}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (secret *Secret) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: secret.Spec.OriginalVersion,
		Kind:    "Secret",
	}
}

// +kubebuilder:object:root=true
// Storage version of v1api20230501.Secret
// Generator information:
// - Generated from: /cdn/resource-manager/Microsoft.Cdn/stable/2023-05-01/afdx.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}/secrets/{secretName}
type SecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Secret `json:"items"`
}

// Storage version of v1api20230501.Secret_Spec
type Secret_Spec struct {
	// AzureName: The name of the resource in Azure. This is often the same as the name of the resource in Kubernetes but it
	// doesn't have to be.
	AzureName       string              `json:"azureName,omitempty"`
	OperatorSpec    *SecretOperatorSpec `json:"operatorSpec,omitempty"`
	OriginalVersion string              `json:"originalVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	// controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	// reference to a cdn.azure.com/Profile resource
	Owner       *genruntime.KnownResourceReference `group:"cdn.azure.com" json:"owner,omitempty" kind:"Profile"`
	Parameters  *SecretParameters                  `json:"parameters,omitempty"`
	PropertyBag genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
}

var _ genruntime.ConvertibleSpec = &Secret_Spec{}

// ConvertSpecFrom populates our Secret_Spec from the provided source
func (secret *Secret_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	if source == secret {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return source.ConvertSpecTo(secret)
}

// ConvertSpecTo populates the provided destination from our Secret_Spec
func (secret *Secret_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	if destination == secret {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return destination.ConvertSpecFrom(secret)
}

// Storage version of v1api20230501.Secret_STATUS
type Secret_STATUS struct {
	Conditions        []conditions.Condition   `json:"conditions,omitempty"`
	DeploymentStatus  *string                  `json:"deploymentStatus,omitempty"`
	Id                *string                  `json:"id,omitempty"`
	Name              *string                  `json:"name,omitempty"`
	Parameters        *SecretParameters_STATUS `json:"parameters,omitempty"`
	ProfileName       *string                  `json:"profileName,omitempty"`
	PropertyBag       genruntime.PropertyBag   `json:"$propertyBag,omitempty"`
	ProvisioningState *string                  `json:"provisioningState,omitempty"`
	SystemData        *SystemData_STATUS       `json:"systemData,omitempty"`
	Type              *string                  `json:"type,omitempty"`
}

var _ genruntime.ConvertibleStatus = &Secret_STATUS{}

// ConvertStatusFrom populates our Secret_STATUS from the provided source
func (secret *Secret_STATUS) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	if source == secret {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return source.ConvertStatusTo(secret)
}

// ConvertStatusTo populates the provided destination from our Secret_STATUS
func (secret *Secret_STATUS) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	if destination == secret {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return destination.ConvertStatusFrom(secret)
}

// Storage version of v1api20230501.SecretOperatorSpec
// Details for configuring operator behavior. Fields in this struct are interpreted by the operator directly rather than being passed to Azure
type SecretOperatorSpec struct {
	ConfigMapExpressions []*core.DestinationExpression `json:"configMapExpressions,omitempty"`
	PropertyBag          genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	SecretExpressions    []*core.DestinationExpression `json:"secretExpressions,omitempty"`
}

// Storage version of v1api20230501.SecretParameters
type SecretParameters struct {
	AzureFirstPartyManagedCertificate *AzureFirstPartyManagedCertificateParameters `json:"azureFirstPartyManagedCertificate,omitempty"`
	CustomerCertificate               *CustomerCertificateParameters               `json:"customerCertificate,omitempty"`
	ManagedCertificate                *ManagedCertificateParameters                `json:"managedCertificate,omitempty"`
	PropertyBag                       genruntime.PropertyBag                       `json:"$propertyBag,omitempty"`
	UrlSigningKey                     *UrlSigningKeyParameters                     `json:"urlSigningKey,omitempty"`
}

// Storage version of v1api20230501.SecretParameters_STATUS
type SecretParameters_STATUS struct {
	AzureFirstPartyManagedCertificate *AzureFirstPartyManagedCertificateParameters_STATUS `json:"azureFirstPartyManagedCertificate,omitempty"`
	CustomerCertificate               *CustomerCertificateParameters_STATUS               `json:"customerCertificate,omitempty"`
	ManagedCertificate                *ManagedCertificateParameters_STATUS                `json:"managedCertificate,omitempty"`
	PropertyBag                       genruntime.PropertyBag                              `json:"$propertyBag,omitempty"`
	UrlSigningKey                     *UrlSigningKeyParameters_STATUS                     `json:"urlSigningKey,omitempty"`
}

// Storage version of v1api20230501.AzureFirstPartyManagedCertificateParameters
type AzureFirstPartyManagedCertificateParameters struct {
	PropertyBag             genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	SubjectAlternativeNames []string               `json:"subjectAlternativeNames,omitempty"`
	Type                    *string                `json:"type,omitempty"`
}

// Storage version of v1api20230501.AzureFirstPartyManagedCertificateParameters_STATUS
type AzureFirstPartyManagedCertificateParameters_STATUS struct {
	CertificateAuthority    *string                   `json:"certificateAuthority,omitempty"`
	ExpirationDate          *string                   `json:"expirationDate,omitempty"`
	PropertyBag             genruntime.PropertyBag    `json:"$propertyBag,omitempty"`
	SecretSource            *ResourceReference_STATUS `json:"secretSource,omitempty"`
	Subject                 *string                   `json:"subject,omitempty"`
	SubjectAlternativeNames []string                  `json:"subjectAlternativeNames,omitempty"`
	Thumbprint              *string                   `json:"thumbprint,omitempty"`
	Type                    *string                   `json:"type,omitempty"`
}

// Storage version of v1api20230501.CustomerCertificateParameters
type CustomerCertificateParameters struct {
	PropertyBag             genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	SecretSource            *ResourceReference     `json:"secretSource,omitempty"`
	SecretVersion           *string                `json:"secretVersion,omitempty"`
	SubjectAlternativeNames []string               `json:"subjectAlternativeNames,omitempty"`
	Type                    *string                `json:"type,omitempty"`
	UseLatestVersion        *bool                  `json:"useLatestVersion,omitempty"`
}

// Storage version of v1api20230501.CustomerCertificateParameters_STATUS
type CustomerCertificateParameters_STATUS struct {
	CertificateAuthority    *string                   `json:"certificateAuthority,omitempty"`
	ExpirationDate          *string                   `json:"expirationDate,omitempty"`
	PropertyBag             genruntime.PropertyBag    `json:"$propertyBag,omitempty"`
	SecretSource            *ResourceReference_STATUS `json:"secretSource,omitempty"`
	SecretVersion           *string                   `json:"secretVersion,omitempty"`
	Subject                 *string                   `json:"subject,omitempty"`
	SubjectAlternativeNames []string                  `json:"subjectAlternativeNames,omitempty"`
	Thumbprint              *string                   `json:"thumbprint,omitempty"`
	Type                    *string                   `json:"type,omitempty"`
	UseLatestVersion        *bool                     `json:"useLatestVersion,omitempty"`
}

// Storage version of v1api20230501.ManagedCertificateParameters
type ManagedCertificateParameters struct {
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Type        *string                `json:"type,omitempty"`
}

// Storage version of v1api20230501.ManagedCertificateParameters_STATUS
type ManagedCertificateParameters_STATUS struct {
	ExpirationDate *string                `json:"expirationDate,omitempty"`
	PropertyBag    genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Subject        *string                `json:"subject,omitempty"`
	Type           *string                `json:"type,omitempty"`
}

// Storage version of v1api20230501.UrlSigningKeyParameters
type UrlSigningKeyParameters struct {
	KeyId         *string                `json:"keyId,omitempty"`
	PropertyBag   genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	SecretSource  *ResourceReference     `json:"secretSource,omitempty"`
	SecretVersion *string                `json:"secretVersion,omitempty"`
	Type          *string                `json:"type,omitempty"`
}

// Storage version of v1api20230501.UrlSigningKeyParameters_STATUS
type UrlSigningKeyParameters_STATUS struct {
	KeyId         *string                   `json:"keyId,omitempty"`
	PropertyBag   genruntime.PropertyBag    `json:"$propertyBag,omitempty"`
	SecretSource  *ResourceReference_STATUS `json:"secretSource,omitempty"`
	SecretVersion *string                   `json:"secretVersion,omitempty"`
	Type          *string                   `json:"type,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Secret{}, &SecretList{})
}
