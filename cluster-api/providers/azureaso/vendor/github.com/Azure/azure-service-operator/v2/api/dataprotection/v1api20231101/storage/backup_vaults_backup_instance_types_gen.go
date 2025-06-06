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

// +kubebuilder:rbac:groups=dataprotection.azure.com,resources=backupvaultsbackupinstances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dataprotection.azure.com,resources={backupvaultsbackupinstances/status,backupvaultsbackupinstances/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// Storage version of v1api20231101.BackupVaultsBackupInstance
// Generator information:
// - Generated from: /dataprotection/resource-manager/Microsoft.DataProtection/stable/2023-11-01/dataprotection.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupInstances/{backupInstanceName}
type BackupVaultsBackupInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BackupVaultsBackupInstance_Spec   `json:"spec,omitempty"`
	Status            BackupVaultsBackupInstance_STATUS `json:"status,omitempty"`
}

var _ conditions.Conditioner = &BackupVaultsBackupInstance{}

// GetConditions returns the conditions of the resource
func (instance *BackupVaultsBackupInstance) GetConditions() conditions.Conditions {
	return instance.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (instance *BackupVaultsBackupInstance) SetConditions(conditions conditions.Conditions) {
	instance.Status.Conditions = conditions
}

var _ configmaps.Exporter = &BackupVaultsBackupInstance{}

// ConfigMapDestinationExpressions returns the Spec.OperatorSpec.ConfigMapExpressions property
func (instance *BackupVaultsBackupInstance) ConfigMapDestinationExpressions() []*core.DestinationExpression {
	if instance.Spec.OperatorSpec == nil {
		return nil
	}
	return instance.Spec.OperatorSpec.ConfigMapExpressions
}

var _ secrets.Exporter = &BackupVaultsBackupInstance{}

// SecretDestinationExpressions returns the Spec.OperatorSpec.SecretExpressions property
func (instance *BackupVaultsBackupInstance) SecretDestinationExpressions() []*core.DestinationExpression {
	if instance.Spec.OperatorSpec == nil {
		return nil
	}
	return instance.Spec.OperatorSpec.SecretExpressions
}

var _ genruntime.KubernetesResource = &BackupVaultsBackupInstance{}

// AzureName returns the Azure name of the resource
func (instance *BackupVaultsBackupInstance) AzureName() string {
	return instance.Spec.AzureName
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2023-11-01"
func (instance BackupVaultsBackupInstance) GetAPIVersion() string {
	return "2023-11-01"
}

// GetResourceScope returns the scope of the resource
func (instance *BackupVaultsBackupInstance) GetResourceScope() genruntime.ResourceScope {
	return genruntime.ResourceScopeResourceGroup
}

// GetSpec returns the specification of this resource
func (instance *BackupVaultsBackupInstance) GetSpec() genruntime.ConvertibleSpec {
	return &instance.Spec
}

// GetStatus returns the status of this resource
func (instance *BackupVaultsBackupInstance) GetStatus() genruntime.ConvertibleStatus {
	return &instance.Status
}

// GetSupportedOperations returns the operations supported by the resource
func (instance *BackupVaultsBackupInstance) GetSupportedOperations() []genruntime.ResourceOperation {
	return []genruntime.ResourceOperation{
		genruntime.ResourceOperationDelete,
		genruntime.ResourceOperationGet,
		genruntime.ResourceOperationPut,
	}
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.DataProtection/backupVaults/backupInstances"
func (instance *BackupVaultsBackupInstance) GetType() string {
	return "Microsoft.DataProtection/backupVaults/backupInstances"
}

// NewEmptyStatus returns a new empty (blank) status
func (instance *BackupVaultsBackupInstance) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &BackupVaultsBackupInstance_STATUS{}
}

// Owner returns the ResourceReference of the owner
func (instance *BackupVaultsBackupInstance) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(instance.Spec)
	return instance.Spec.Owner.AsResourceReference(group, kind)
}

// SetStatus sets the status of this resource
func (instance *BackupVaultsBackupInstance) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*BackupVaultsBackupInstance_STATUS); ok {
		instance.Status = *st
		return nil
	}

	// Convert status to required version
	var st BackupVaultsBackupInstance_STATUS
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	instance.Status = st
	return nil
}

// Hub marks that this BackupVaultsBackupInstance is the hub type for conversion
func (instance *BackupVaultsBackupInstance) Hub() {}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (instance *BackupVaultsBackupInstance) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: instance.Spec.OriginalVersion,
		Kind:    "BackupVaultsBackupInstance",
	}
}

// +kubebuilder:object:root=true
// Storage version of v1api20231101.BackupVaultsBackupInstance
// Generator information:
// - Generated from: /dataprotection/resource-manager/Microsoft.DataProtection/stable/2023-11-01/dataprotection.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataProtection/backupVaults/{vaultName}/backupInstances/{backupInstanceName}
type BackupVaultsBackupInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BackupVaultsBackupInstance `json:"items"`
}

// Storage version of v1api20231101.BackupVaultsBackupInstance_Spec
type BackupVaultsBackupInstance_Spec struct {
	// AzureName: The name of the resource in Azure. This is often the same as the name of the resource in Kubernetes but it
	// doesn't have to be.
	AzureName       string                                  `json:"azureName,omitempty"`
	OperatorSpec    *BackupVaultsBackupInstanceOperatorSpec `json:"operatorSpec,omitempty"`
	OriginalVersion string                                  `json:"originalVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	// controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	// reference to a dataprotection.azure.com/BackupVault resource
	Owner       *genruntime.KnownResourceReference `group:"dataprotection.azure.com" json:"owner,omitempty" kind:"BackupVault"`
	Properties  *BackupInstance                    `json:"properties,omitempty"`
	PropertyBag genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
	Tags        map[string]string                  `json:"tags,omitempty"`
}

var _ genruntime.ConvertibleSpec = &BackupVaultsBackupInstance_Spec{}

// ConvertSpecFrom populates our BackupVaultsBackupInstance_Spec from the provided source
func (instance *BackupVaultsBackupInstance_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	if source == instance {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return source.ConvertSpecTo(instance)
}

// ConvertSpecTo populates the provided destination from our BackupVaultsBackupInstance_Spec
func (instance *BackupVaultsBackupInstance_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	if destination == instance {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return destination.ConvertSpecFrom(instance)
}

// Storage version of v1api20231101.BackupVaultsBackupInstance_STATUS
type BackupVaultsBackupInstance_STATUS struct {
	Conditions  []conditions.Condition `json:"conditions,omitempty"`
	Id          *string                `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Properties  *BackupInstance_STATUS `json:"properties,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	SystemData  *SystemData_STATUS     `json:"systemData,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
	Type        *string                `json:"type,omitempty"`
}

var _ genruntime.ConvertibleStatus = &BackupVaultsBackupInstance_STATUS{}

// ConvertStatusFrom populates our BackupVaultsBackupInstance_STATUS from the provided source
func (instance *BackupVaultsBackupInstance_STATUS) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	if source == instance {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return source.ConvertStatusTo(instance)
}

// ConvertStatusTo populates the provided destination from our BackupVaultsBackupInstance_STATUS
func (instance *BackupVaultsBackupInstance_STATUS) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	if destination == instance {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return destination.ConvertStatusFrom(instance)
}

// Storage version of v1api20231101.BackupInstance
// Backup Instance
type BackupInstance struct {
	DataSourceInfo            *Datasource            `json:"dataSourceInfo,omitempty"`
	DataSourceSetInfo         *DatasourceSet         `json:"dataSourceSetInfo,omitempty"`
	DatasourceAuthCredentials *AuthCredentials       `json:"datasourceAuthCredentials,omitempty"`
	FriendlyName              *string                `json:"friendlyName,omitempty"`
	IdentityDetails           *IdentityDetails       `json:"identityDetails,omitempty"`
	ObjectType                *string                `json:"objectType,omitempty"`
	PolicyInfo                *PolicyInfo            `json:"policyInfo,omitempty"`
	PropertyBag               genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	ValidationType            *string                `json:"validationType,omitempty"`
}

// Storage version of v1api20231101.BackupInstance_STATUS
// Backup Instance
type BackupInstance_STATUS struct {
	CurrentProtectionState    *string                         `json:"currentProtectionState,omitempty"`
	DataSourceInfo            *Datasource_STATUS              `json:"dataSourceInfo,omitempty"`
	DataSourceSetInfo         *DatasourceSet_STATUS           `json:"dataSourceSetInfo,omitempty"`
	DatasourceAuthCredentials *AuthCredentials_STATUS         `json:"datasourceAuthCredentials,omitempty"`
	FriendlyName              *string                         `json:"friendlyName,omitempty"`
	IdentityDetails           *IdentityDetails_STATUS         `json:"identityDetails,omitempty"`
	ObjectType                *string                         `json:"objectType,omitempty"`
	PolicyInfo                *PolicyInfo_STATUS              `json:"policyInfo,omitempty"`
	PropertyBag               genruntime.PropertyBag          `json:"$propertyBag,omitempty"`
	ProtectionErrorDetails    *UserFacingError_STATUS         `json:"protectionErrorDetails,omitempty"`
	ProtectionStatus          *ProtectionStatusDetails_STATUS `json:"protectionStatus,omitempty"`
	ProvisioningState         *string                         `json:"provisioningState,omitempty"`
	ValidationType            *string                         `json:"validationType,omitempty"`
}

// Storage version of v1api20231101.BackupVaultsBackupInstanceOperatorSpec
// Details for configuring operator behavior. Fields in this struct are interpreted by the operator directly rather than being passed to Azure
type BackupVaultsBackupInstanceOperatorSpec struct {
	ConfigMapExpressions []*core.DestinationExpression `json:"configMapExpressions,omitempty"`
	PropertyBag          genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	SecretExpressions    []*core.DestinationExpression `json:"secretExpressions,omitempty"`
}

// Storage version of v1api20231101.AuthCredentials
type AuthCredentials struct {
	PropertyBag                     genruntime.PropertyBag           `json:"$propertyBag,omitempty"`
	SecretStoreBasedAuthCredentials *SecretStoreBasedAuthCredentials `json:"secretStoreBasedAuthCredentials,omitempty"`
}

// Storage version of v1api20231101.AuthCredentials_STATUS
type AuthCredentials_STATUS struct {
	PropertyBag                     genruntime.PropertyBag                  `json:"$propertyBag,omitempty"`
	SecretStoreBasedAuthCredentials *SecretStoreBasedAuthCredentials_STATUS `json:"secretStoreBasedAuthCredentials,omitempty"`
}

// Storage version of v1api20231101.Datasource
// Datasource to be backed up
type Datasource struct {
	DatasourceType     *string                 `json:"datasourceType,omitempty"`
	ObjectType         *string                 `json:"objectType,omitempty"`
	PropertyBag        genruntime.PropertyBag  `json:"$propertyBag,omitempty"`
	ResourceLocation   *string                 `json:"resourceLocation,omitempty"`
	ResourceName       *string                 `json:"resourceName,omitempty"`
	ResourceProperties *BaseResourceProperties `json:"resourceProperties,omitempty"`

	// +kubebuilder:validation:Required
	// ResourceReference: Full ARM ID of the resource. For azure resources, this is ARM ID. For non azure resources, this will
	// be the ID created by backup service via Fabric/Vault.
	ResourceReference *genruntime.ResourceReference `armReference:"ResourceID" json:"resourceReference,omitempty"`
	ResourceType      *string                       `json:"resourceType,omitempty"`
	ResourceUri       *string                       `json:"resourceUri,omitempty"`
}

// Storage version of v1api20231101.Datasource_STATUS
// Datasource to be backed up
type Datasource_STATUS struct {
	DatasourceType     *string                        `json:"datasourceType,omitempty"`
	ObjectType         *string                        `json:"objectType,omitempty"`
	PropertyBag        genruntime.PropertyBag         `json:"$propertyBag,omitempty"`
	ResourceID         *string                        `json:"resourceID,omitempty"`
	ResourceLocation   *string                        `json:"resourceLocation,omitempty"`
	ResourceName       *string                        `json:"resourceName,omitempty"`
	ResourceProperties *BaseResourceProperties_STATUS `json:"resourceProperties,omitempty"`
	ResourceType       *string                        `json:"resourceType,omitempty"`
	ResourceUri        *string                        `json:"resourceUri,omitempty"`
}

// Storage version of v1api20231101.DatasourceSet
// DatasourceSet details of datasource to be backed up
type DatasourceSet struct {
	DatasourceType     *string                 `json:"datasourceType,omitempty"`
	ObjectType         *string                 `json:"objectType,omitempty"`
	PropertyBag        genruntime.PropertyBag  `json:"$propertyBag,omitempty"`
	ResourceLocation   *string                 `json:"resourceLocation,omitempty"`
	ResourceName       *string                 `json:"resourceName,omitempty"`
	ResourceProperties *BaseResourceProperties `json:"resourceProperties,omitempty"`

	// +kubebuilder:validation:Required
	// ResourceReference: Full ARM ID of the resource. For azure resources, this is ARM ID. For non azure resources, this will
	// be the ID created by backup service via Fabric/Vault.
	ResourceReference *genruntime.ResourceReference `armReference:"ResourceID" json:"resourceReference,omitempty"`
	ResourceType      *string                       `json:"resourceType,omitempty"`
	ResourceUri       *string                       `json:"resourceUri,omitempty"`
}

// Storage version of v1api20231101.DatasourceSet_STATUS
// DatasourceSet details of datasource to be backed up
type DatasourceSet_STATUS struct {
	DatasourceType     *string                        `json:"datasourceType,omitempty"`
	ObjectType         *string                        `json:"objectType,omitempty"`
	PropertyBag        genruntime.PropertyBag         `json:"$propertyBag,omitempty"`
	ResourceID         *string                        `json:"resourceID,omitempty"`
	ResourceLocation   *string                        `json:"resourceLocation,omitempty"`
	ResourceName       *string                        `json:"resourceName,omitempty"`
	ResourceProperties *BaseResourceProperties_STATUS `json:"resourceProperties,omitempty"`
	ResourceType       *string                        `json:"resourceType,omitempty"`
	ResourceUri        *string                        `json:"resourceUri,omitempty"`
}

// Storage version of v1api20231101.IdentityDetails
type IdentityDetails struct {
	PropertyBag                genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	UseSystemAssignedIdentity  *bool                  `json:"useSystemAssignedIdentity,omitempty"`
	UserAssignedIdentityArmUrl *string                `json:"userAssignedIdentityArmUrl,omitempty"`
}

// Storage version of v1api20231101.IdentityDetails_STATUS
type IdentityDetails_STATUS struct {
	PropertyBag                genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	UseSystemAssignedIdentity  *bool                  `json:"useSystemAssignedIdentity,omitempty"`
	UserAssignedIdentityArmUrl *string                `json:"userAssignedIdentityArmUrl,omitempty"`
}

// Storage version of v1api20231101.PolicyInfo
// Policy Info in backupInstance
type PolicyInfo struct {
	PolicyParameters *PolicyParameters `json:"policyParameters,omitempty"`

	// +kubebuilder:validation:Required
	PolicyReference *genruntime.ResourceReference `armReference:"PolicyId" json:"policyReference,omitempty"`
	PropertyBag     genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.PolicyInfo_STATUS
// Policy Info in backupInstance
type PolicyInfo_STATUS struct {
	PolicyId         *string                  `json:"policyId,omitempty"`
	PolicyParameters *PolicyParameters_STATUS `json:"policyParameters,omitempty"`
	PolicyVersion    *string                  `json:"policyVersion,omitempty"`
	PropertyBag      genruntime.PropertyBag   `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.ProtectionStatusDetails_STATUS
// Protection status details
type ProtectionStatusDetails_STATUS struct {
	ErrorDetails *UserFacingError_STATUS `json:"errorDetails,omitempty"`
	PropertyBag  genruntime.PropertyBag  `json:"$propertyBag,omitempty"`
	Status       *string                 `json:"status,omitempty"`
}

// Storage version of v1api20231101.UserFacingError_STATUS
// Error object used by layers that have access to localized content, and propagate that to user
type UserFacingError_STATUS struct {
	Code              *string                           `json:"code,omitempty"`
	Details           []UserFacingError_STATUS_Unrolled `json:"details,omitempty"`
	InnerError        *InnerError_STATUS                `json:"innerError,omitempty"`
	IsRetryable       *bool                             `json:"isRetryable,omitempty"`
	IsUserError       *bool                             `json:"isUserError,omitempty"`
	Message           *string                           `json:"message,omitempty"`
	Properties        map[string]string                 `json:"properties,omitempty"`
	PropertyBag       genruntime.PropertyBag            `json:"$propertyBag,omitempty"`
	RecommendedAction []string                          `json:"recommendedAction,omitempty"`
	Target            *string                           `json:"target,omitempty"`
}

// Storage version of v1api20231101.BaseResourceProperties
type BaseResourceProperties struct {
	DefaultResourceProperties *DefaultResourceProperties `json:"defaultResourceProperties,omitempty"`
	PropertyBag               genruntime.PropertyBag     `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.BaseResourceProperties_STATUS
type BaseResourceProperties_STATUS struct {
	DefaultResourceProperties *DefaultResourceProperties_STATUS `json:"defaultResourceProperties,omitempty"`
	PropertyBag               genruntime.PropertyBag            `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.InnerError_STATUS
// Inner Error
type InnerError_STATUS struct {
	AdditionalInfo     map[string]string           `json:"additionalInfo,omitempty"`
	Code               *string                     `json:"code,omitempty"`
	EmbeddedInnerError *InnerError_STATUS_Unrolled `json:"embeddedInnerError,omitempty"`
	PropertyBag        genruntime.PropertyBag      `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.PolicyParameters
// Parameters in Policy
type PolicyParameters struct {
	BackupDatasourceParametersList []BackupDatasourceParameters `json:"backupDatasourceParametersList,omitempty"`
	DataStoreParametersList        []DataStoreParameters        `json:"dataStoreParametersList,omitempty"`
	PropertyBag                    genruntime.PropertyBag       `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.PolicyParameters_STATUS
// Parameters in Policy
type PolicyParameters_STATUS struct {
	BackupDatasourceParametersList []BackupDatasourceParameters_STATUS `json:"backupDatasourceParametersList,omitempty"`
	DataStoreParametersList        []DataStoreParameters_STATUS        `json:"dataStoreParametersList,omitempty"`
	PropertyBag                    genruntime.PropertyBag              `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.SecretStoreBasedAuthCredentials
type SecretStoreBasedAuthCredentials struct {
	ObjectType          *string                `json:"objectType,omitempty"`
	PropertyBag         genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	SecretStoreResource *SecretStoreResource   `json:"secretStoreResource,omitempty"`
}

// Storage version of v1api20231101.SecretStoreBasedAuthCredentials_STATUS
type SecretStoreBasedAuthCredentials_STATUS struct {
	ObjectType          *string                     `json:"objectType,omitempty"`
	PropertyBag         genruntime.PropertyBag      `json:"$propertyBag,omitempty"`
	SecretStoreResource *SecretStoreResource_STATUS `json:"secretStoreResource,omitempty"`
}

// Storage version of v1api20231101.UserFacingError_STATUS_Unrolled
type UserFacingError_STATUS_Unrolled struct {
	Code              *string                `json:"code,omitempty"`
	InnerError        *InnerError_STATUS     `json:"innerError,omitempty"`
	IsRetryable       *bool                  `json:"isRetryable,omitempty"`
	IsUserError       *bool                  `json:"isUserError,omitempty"`
	Message           *string                `json:"message,omitempty"`
	Properties        map[string]string      `json:"properties,omitempty"`
	PropertyBag       genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	RecommendedAction []string               `json:"recommendedAction,omitempty"`
	Target            *string                `json:"target,omitempty"`
}

// Storage version of v1api20231101.BackupDatasourceParameters
type BackupDatasourceParameters struct {
	Blob              *BlobBackupDatasourceParameters              `json:"blobBackupDatasourceParameters,omitempty"`
	KubernetesCluster *KubernetesClusterBackupDatasourceParameters `json:"kubernetesClusterBackupDatasourceParameters,omitempty"`
	PropertyBag       genruntime.PropertyBag                       `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.BackupDatasourceParameters_STATUS
type BackupDatasourceParameters_STATUS struct {
	Blob              *BlobBackupDatasourceParameters_STATUS              `json:"blobBackupDatasourceParameters,omitempty"`
	KubernetesCluster *KubernetesClusterBackupDatasourceParameters_STATUS `json:"kubernetesClusterBackupDatasourceParameters,omitempty"`
	PropertyBag       genruntime.PropertyBag                              `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.DataStoreParameters
type DataStoreParameters struct {
	AzureOperationalStoreParameters *AzureOperationalStoreParameters `json:"azureOperationalStoreParameters,omitempty"`
	PropertyBag                     genruntime.PropertyBag           `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.DataStoreParameters_STATUS
type DataStoreParameters_STATUS struct {
	AzureOperationalStoreParameters *AzureOperationalStoreParameters_STATUS `json:"azureOperationalStoreParameters,omitempty"`
	PropertyBag                     genruntime.PropertyBag                  `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.DefaultResourceProperties
type DefaultResourceProperties struct {
	ObjectType  *string                `json:"objectType,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.DefaultResourceProperties_STATUS
type DefaultResourceProperties_STATUS struct {
	ObjectType  *string                `json:"objectType,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.InnerError_STATUS_Unrolled
type InnerError_STATUS_Unrolled struct {
	AdditionalInfo map[string]string      `json:"additionalInfo,omitempty"`
	Code           *string                `json:"code,omitempty"`
	PropertyBag    genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.SecretStoreResource
// Class representing a secret store resource.
type SecretStoreResource struct {
	PropertyBag     genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	SecretStoreType *string                `json:"secretStoreType,omitempty"`
	Uri             *string                `json:"uri,omitempty"`
	Value           *string                `json:"value,omitempty"`
}

// Storage version of v1api20231101.SecretStoreResource_STATUS
// Class representing a secret store resource.
type SecretStoreResource_STATUS struct {
	PropertyBag     genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	SecretStoreType *string                `json:"secretStoreType,omitempty"`
	Uri             *string                `json:"uri,omitempty"`
	Value           *string                `json:"value,omitempty"`
}

// Storage version of v1api20231101.AzureOperationalStoreParameters
type AzureOperationalStoreParameters struct {
	DataStoreType *string                `json:"dataStoreType,omitempty"`
	ObjectType    *string                `json:"objectType,omitempty"`
	PropertyBag   genruntime.PropertyBag `json:"$propertyBag,omitempty"`

	// ResourceGroupReference: Gets or sets the Snapshot Resource Group Uri.
	ResourceGroupReference *genruntime.ResourceReference `armReference:"ResourceGroupId" json:"resourceGroupReference,omitempty"`
}

// Storage version of v1api20231101.AzureOperationalStoreParameters_STATUS
type AzureOperationalStoreParameters_STATUS struct {
	DataStoreType   *string                `json:"dataStoreType,omitempty"`
	ObjectType      *string                `json:"objectType,omitempty"`
	PropertyBag     genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	ResourceGroupId *string                `json:"resourceGroupId,omitempty"`
}

// Storage version of v1api20231101.BlobBackupDatasourceParameters
type BlobBackupDatasourceParameters struct {
	ContainersList []string               `json:"containersList,omitempty"`
	ObjectType     *string                `json:"objectType,omitempty"`
	PropertyBag    genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.BlobBackupDatasourceParameters_STATUS
type BlobBackupDatasourceParameters_STATUS struct {
	ContainersList []string               `json:"containersList,omitempty"`
	ObjectType     *string                `json:"objectType,omitempty"`
	PropertyBag    genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.KubernetesClusterBackupDatasourceParameters
type KubernetesClusterBackupDatasourceParameters struct {
	BackupHookReferences         []NamespacedNameResource `json:"backupHookReferences,omitempty"`
	ExcludedNamespaces           []string                 `json:"excludedNamespaces,omitempty"`
	ExcludedResourceTypes        []string                 `json:"excludedResourceTypes,omitempty"`
	IncludeClusterScopeResources *bool                    `json:"includeClusterScopeResources,omitempty"`
	IncludedNamespaces           []string                 `json:"includedNamespaces,omitempty"`
	IncludedResourceTypes        []string                 `json:"includedResourceTypes,omitempty"`
	LabelSelectors               []string                 `json:"labelSelectors,omitempty"`
	ObjectType                   *string                  `json:"objectType,omitempty"`
	PropertyBag                  genruntime.PropertyBag   `json:"$propertyBag,omitempty"`
	SnapshotVolumes              *bool                    `json:"snapshotVolumes,omitempty"`
}

// Storage version of v1api20231101.KubernetesClusterBackupDatasourceParameters_STATUS
type KubernetesClusterBackupDatasourceParameters_STATUS struct {
	BackupHookReferences         []NamespacedNameResource_STATUS `json:"backupHookReferences,omitempty"`
	ExcludedNamespaces           []string                        `json:"excludedNamespaces,omitempty"`
	ExcludedResourceTypes        []string                        `json:"excludedResourceTypes,omitempty"`
	IncludeClusterScopeResources *bool                           `json:"includeClusterScopeResources,omitempty"`
	IncludedNamespaces           []string                        `json:"includedNamespaces,omitempty"`
	IncludedResourceTypes        []string                        `json:"includedResourceTypes,omitempty"`
	LabelSelectors               []string                        `json:"labelSelectors,omitempty"`
	ObjectType                   *string                         `json:"objectType,omitempty"`
	PropertyBag                  genruntime.PropertyBag          `json:"$propertyBag,omitempty"`
	SnapshotVolumes              *bool                           `json:"snapshotVolumes,omitempty"`
}

// Storage version of v1api20231101.NamespacedNameResource
// Class to refer resources which contains namespace and name
type NamespacedNameResource struct {
	Name        *string                `json:"name,omitempty"`
	Namespace   *string                `json:"namespace,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20231101.NamespacedNameResource_STATUS
// Class to refer resources which contains namespace and name
type NamespacedNameResource_STATUS struct {
	Name        *string                `json:"name,omitempty"`
	Namespace   *string                `json:"namespace,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

func init() {
	SchemeBuilder.Register(&BackupVaultsBackupInstance{}, &BackupVaultsBackupInstanceList{})
}
