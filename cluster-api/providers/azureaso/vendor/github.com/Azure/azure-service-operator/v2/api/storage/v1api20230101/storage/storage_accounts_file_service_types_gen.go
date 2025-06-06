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

// +kubebuilder:rbac:groups=storage.azure.com,resources=storageaccountsfileservices,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=storage.azure.com,resources={storageaccountsfileservices/status,storageaccountsfileservices/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// Storage version of v1api20230101.StorageAccountsFileService
// Generator information:
// - Generated from: /storage/resource-manager/Microsoft.Storage/stable/2023-01-01/file.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/fileServices/default
type StorageAccountsFileService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StorageAccountsFileService_Spec   `json:"spec,omitempty"`
	Status            StorageAccountsFileService_STATUS `json:"status,omitempty"`
}

var _ conditions.Conditioner = &StorageAccountsFileService{}

// GetConditions returns the conditions of the resource
func (service *StorageAccountsFileService) GetConditions() conditions.Conditions {
	return service.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (service *StorageAccountsFileService) SetConditions(conditions conditions.Conditions) {
	service.Status.Conditions = conditions
}

var _ configmaps.Exporter = &StorageAccountsFileService{}

// ConfigMapDestinationExpressions returns the Spec.OperatorSpec.ConfigMapExpressions property
func (service *StorageAccountsFileService) ConfigMapDestinationExpressions() []*core.DestinationExpression {
	if service.Spec.OperatorSpec == nil {
		return nil
	}
	return service.Spec.OperatorSpec.ConfigMapExpressions
}

var _ secrets.Exporter = &StorageAccountsFileService{}

// SecretDestinationExpressions returns the Spec.OperatorSpec.SecretExpressions property
func (service *StorageAccountsFileService) SecretDestinationExpressions() []*core.DestinationExpression {
	if service.Spec.OperatorSpec == nil {
		return nil
	}
	return service.Spec.OperatorSpec.SecretExpressions
}

var _ genruntime.KubernetesResource = &StorageAccountsFileService{}

// AzureName returns the Azure name of the resource (always "default")
func (service *StorageAccountsFileService) AzureName() string {
	return "default"
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2023-01-01"
func (service StorageAccountsFileService) GetAPIVersion() string {
	return "2023-01-01"
}

// GetResourceScope returns the scope of the resource
func (service *StorageAccountsFileService) GetResourceScope() genruntime.ResourceScope {
	return genruntime.ResourceScopeResourceGroup
}

// GetSpec returns the specification of this resource
func (service *StorageAccountsFileService) GetSpec() genruntime.ConvertibleSpec {
	return &service.Spec
}

// GetStatus returns the status of this resource
func (service *StorageAccountsFileService) GetStatus() genruntime.ConvertibleStatus {
	return &service.Status
}

// GetSupportedOperations returns the operations supported by the resource
func (service *StorageAccountsFileService) GetSupportedOperations() []genruntime.ResourceOperation {
	return []genruntime.ResourceOperation{
		genruntime.ResourceOperationGet,
		genruntime.ResourceOperationPut,
	}
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.Storage/storageAccounts/fileServices"
func (service *StorageAccountsFileService) GetType() string {
	return "Microsoft.Storage/storageAccounts/fileServices"
}

// NewEmptyStatus returns a new empty (blank) status
func (service *StorageAccountsFileService) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &StorageAccountsFileService_STATUS{}
}

// Owner returns the ResourceReference of the owner
func (service *StorageAccountsFileService) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(service.Spec)
	return service.Spec.Owner.AsResourceReference(group, kind)
}

// SetStatus sets the status of this resource
func (service *StorageAccountsFileService) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*StorageAccountsFileService_STATUS); ok {
		service.Status = *st
		return nil
	}

	// Convert status to required version
	var st StorageAccountsFileService_STATUS
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	service.Status = st
	return nil
}

// Hub marks that this StorageAccountsFileService is the hub type for conversion
func (service *StorageAccountsFileService) Hub() {}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (service *StorageAccountsFileService) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: service.Spec.OriginalVersion,
		Kind:    "StorageAccountsFileService",
	}
}

// +kubebuilder:object:root=true
// Storage version of v1api20230101.StorageAccountsFileService
// Generator information:
// - Generated from: /storage/resource-manager/Microsoft.Storage/stable/2023-01-01/file.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Storage/storageAccounts/{accountName}/fileServices/default
type StorageAccountsFileServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StorageAccountsFileService `json:"items"`
}

// Storage version of v1api20230101.StorageAccountsFileService_Spec
type StorageAccountsFileService_Spec struct {
	Cors            *CorsRules                              `json:"cors,omitempty"`
	OperatorSpec    *StorageAccountsFileServiceOperatorSpec `json:"operatorSpec,omitempty"`
	OriginalVersion string                                  `json:"originalVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	// controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	// reference to a storage.azure.com/StorageAccount resource
	Owner                      *genruntime.KnownResourceReference `group:"storage.azure.com" json:"owner,omitempty" kind:"StorageAccount"`
	PropertyBag                genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
	ProtocolSettings           *ProtocolSettings                  `json:"protocolSettings,omitempty"`
	ShareDeleteRetentionPolicy *DeleteRetentionPolicy             `json:"shareDeleteRetentionPolicy,omitempty"`
}

var _ genruntime.ConvertibleSpec = &StorageAccountsFileService_Spec{}

// ConvertSpecFrom populates our StorageAccountsFileService_Spec from the provided source
func (service *StorageAccountsFileService_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	if source == service {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return source.ConvertSpecTo(service)
}

// ConvertSpecTo populates the provided destination from our StorageAccountsFileService_Spec
func (service *StorageAccountsFileService_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	if destination == service {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return destination.ConvertSpecFrom(service)
}

// Storage version of v1api20230101.StorageAccountsFileService_STATUS
type StorageAccountsFileService_STATUS struct {
	Conditions                 []conditions.Condition        `json:"conditions,omitempty"`
	Cors                       *CorsRules_STATUS             `json:"cors,omitempty"`
	Id                         *string                       `json:"id,omitempty"`
	Name                       *string                       `json:"name,omitempty"`
	PropertyBag                genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	ProtocolSettings           *ProtocolSettings_STATUS      `json:"protocolSettings,omitempty"`
	ShareDeleteRetentionPolicy *DeleteRetentionPolicy_STATUS `json:"shareDeleteRetentionPolicy,omitempty"`
	Sku                        *Sku_STATUS                   `json:"sku,omitempty"`
	Type                       *string                       `json:"type,omitempty"`
}

var _ genruntime.ConvertibleStatus = &StorageAccountsFileService_STATUS{}

// ConvertStatusFrom populates our StorageAccountsFileService_STATUS from the provided source
func (service *StorageAccountsFileService_STATUS) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	if source == service {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return source.ConvertStatusTo(service)
}

// ConvertStatusTo populates the provided destination from our StorageAccountsFileService_STATUS
func (service *StorageAccountsFileService_STATUS) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	if destination == service {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return destination.ConvertStatusFrom(service)
}

// Storage version of v1api20230101.ProtocolSettings
// Protocol settings for file service
type ProtocolSettings struct {
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Smb         *SmbSetting            `json:"smb,omitempty"`
}

// Storage version of v1api20230101.ProtocolSettings_STATUS
// Protocol settings for file service
type ProtocolSettings_STATUS struct {
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Smb         *SmbSetting_STATUS     `json:"smb,omitempty"`
}

// Storage version of v1api20230101.StorageAccountsFileServiceOperatorSpec
// Details for configuring operator behavior. Fields in this struct are interpreted by the operator directly rather than being passed to Azure
type StorageAccountsFileServiceOperatorSpec struct {
	ConfigMapExpressions []*core.DestinationExpression `json:"configMapExpressions,omitempty"`
	PropertyBag          genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	SecretExpressions    []*core.DestinationExpression `json:"secretExpressions,omitempty"`
}

// Storage version of v1api20230101.SmbSetting
// Setting for SMB protocol
type SmbSetting struct {
	AuthenticationMethods    *string                `json:"authenticationMethods,omitempty"`
	ChannelEncryption        *string                `json:"channelEncryption,omitempty"`
	KerberosTicketEncryption *string                `json:"kerberosTicketEncryption,omitempty"`
	Multichannel             *Multichannel          `json:"multichannel,omitempty"`
	PropertyBag              genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Versions                 *string                `json:"versions,omitempty"`
}

// Storage version of v1api20230101.SmbSetting_STATUS
// Setting for SMB protocol
type SmbSetting_STATUS struct {
	AuthenticationMethods    *string                `json:"authenticationMethods,omitempty"`
	ChannelEncryption        *string                `json:"channelEncryption,omitempty"`
	KerberosTicketEncryption *string                `json:"kerberosTicketEncryption,omitempty"`
	Multichannel             *Multichannel_STATUS   `json:"multichannel,omitempty"`
	PropertyBag              genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Versions                 *string                `json:"versions,omitempty"`
}

// Storage version of v1api20230101.Multichannel
// Multichannel setting. Applies to Premium FileStorage only.
type Multichannel struct {
	Enabled     *bool                  `json:"enabled,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

// Storage version of v1api20230101.Multichannel_STATUS
// Multichannel setting. Applies to Premium FileStorage only.
type Multichannel_STATUS struct {
	Enabled     *bool                  `json:"enabled,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

func init() {
	SchemeBuilder.Register(&StorageAccountsFileService{}, &StorageAccountsFileServiceList{})
}
