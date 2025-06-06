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

// +kubebuilder:rbac:groups=dbforpostgresql.azure.com,resources=flexibleserversconfigurations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dbforpostgresql.azure.com,resources={flexibleserversconfigurations/status,flexibleserversconfigurations/finalizers},verbs=get;update;patch

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// Storage version of v1api20221201.FlexibleServersConfiguration
// Generator information:
// - Generated from: /postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2022-12-01/Configuration.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{serverName}/configurations/{configurationName}
type FlexibleServersConfiguration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FlexibleServersConfiguration_Spec   `json:"spec,omitempty"`
	Status            FlexibleServersConfiguration_STATUS `json:"status,omitempty"`
}

var _ conditions.Conditioner = &FlexibleServersConfiguration{}

// GetConditions returns the conditions of the resource
func (configuration *FlexibleServersConfiguration) GetConditions() conditions.Conditions {
	return configuration.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (configuration *FlexibleServersConfiguration) SetConditions(conditions conditions.Conditions) {
	configuration.Status.Conditions = conditions
}

var _ configmaps.Exporter = &FlexibleServersConfiguration{}

// ConfigMapDestinationExpressions returns the Spec.OperatorSpec.ConfigMapExpressions property
func (configuration *FlexibleServersConfiguration) ConfigMapDestinationExpressions() []*core.DestinationExpression {
	if configuration.Spec.OperatorSpec == nil {
		return nil
	}
	return configuration.Spec.OperatorSpec.ConfigMapExpressions
}

var _ secrets.Exporter = &FlexibleServersConfiguration{}

// SecretDestinationExpressions returns the Spec.OperatorSpec.SecretExpressions property
func (configuration *FlexibleServersConfiguration) SecretDestinationExpressions() []*core.DestinationExpression {
	if configuration.Spec.OperatorSpec == nil {
		return nil
	}
	return configuration.Spec.OperatorSpec.SecretExpressions
}

var _ genruntime.KubernetesResource = &FlexibleServersConfiguration{}

// AzureName returns the Azure name of the resource
func (configuration *FlexibleServersConfiguration) AzureName() string {
	return configuration.Spec.AzureName
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2022-12-01"
func (configuration FlexibleServersConfiguration) GetAPIVersion() string {
	return "2022-12-01"
}

// GetResourceScope returns the scope of the resource
func (configuration *FlexibleServersConfiguration) GetResourceScope() genruntime.ResourceScope {
	return genruntime.ResourceScopeResourceGroup
}

// GetSpec returns the specification of this resource
func (configuration *FlexibleServersConfiguration) GetSpec() genruntime.ConvertibleSpec {
	return &configuration.Spec
}

// GetStatus returns the status of this resource
func (configuration *FlexibleServersConfiguration) GetStatus() genruntime.ConvertibleStatus {
	return &configuration.Status
}

// GetSupportedOperations returns the operations supported by the resource
func (configuration *FlexibleServersConfiguration) GetSupportedOperations() []genruntime.ResourceOperation {
	return []genruntime.ResourceOperation{
		genruntime.ResourceOperationGet,
		genruntime.ResourceOperationPut,
	}
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.DBforPostgreSQL/flexibleServers/configurations"
func (configuration *FlexibleServersConfiguration) GetType() string {
	return "Microsoft.DBforPostgreSQL/flexibleServers/configurations"
}

// NewEmptyStatus returns a new empty (blank) status
func (configuration *FlexibleServersConfiguration) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &FlexibleServersConfiguration_STATUS{}
}

// Owner returns the ResourceReference of the owner
func (configuration *FlexibleServersConfiguration) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(configuration.Spec)
	return configuration.Spec.Owner.AsResourceReference(group, kind)
}

// SetStatus sets the status of this resource
func (configuration *FlexibleServersConfiguration) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*FlexibleServersConfiguration_STATUS); ok {
		configuration.Status = *st
		return nil
	}

	// Convert status to required version
	var st FlexibleServersConfiguration_STATUS
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	configuration.Status = st
	return nil
}

// Hub marks that this FlexibleServersConfiguration is the hub type for conversion
func (configuration *FlexibleServersConfiguration) Hub() {}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (configuration *FlexibleServersConfiguration) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: configuration.Spec.OriginalVersion,
		Kind:    "FlexibleServersConfiguration",
	}
}

// +kubebuilder:object:root=true
// Storage version of v1api20221201.FlexibleServersConfiguration
// Generator information:
// - Generated from: /postgresql/resource-manager/Microsoft.DBforPostgreSQL/stable/2022-12-01/Configuration.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBforPostgreSQL/flexibleServers/{serverName}/configurations/{configurationName}
type FlexibleServersConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FlexibleServersConfiguration `json:"items"`
}

// Storage version of v1api20221201.FlexibleServersConfiguration_Spec
type FlexibleServersConfiguration_Spec struct {
	// AzureName: The name of the resource in Azure. This is often the same as the name of the resource in Kubernetes but it
	// doesn't have to be.
	AzureName       string                                    `json:"azureName,omitempty"`
	OperatorSpec    *FlexibleServersConfigurationOperatorSpec `json:"operatorSpec,omitempty"`
	OriginalVersion string                                    `json:"originalVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	// controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	// reference to a dbforpostgresql.azure.com/FlexibleServer resource
	Owner       *genruntime.KnownResourceReference `group:"dbforpostgresql.azure.com" json:"owner,omitempty" kind:"FlexibleServer"`
	PropertyBag genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
	Source      *string                            `json:"source,omitempty"`
	Value       *string                            `json:"value,omitempty"`
}

var _ genruntime.ConvertibleSpec = &FlexibleServersConfiguration_Spec{}

// ConvertSpecFrom populates our FlexibleServersConfiguration_Spec from the provided source
func (configuration *FlexibleServersConfiguration_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	if source == configuration {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return source.ConvertSpecTo(configuration)
}

// ConvertSpecTo populates the provided destination from our FlexibleServersConfiguration_Spec
func (configuration *FlexibleServersConfiguration_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	if destination == configuration {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleSpec")
	}

	return destination.ConvertSpecFrom(configuration)
}

// Storage version of v1api20221201.FlexibleServersConfiguration_STATUS
type FlexibleServersConfiguration_STATUS struct {
	AllowedValues          *string                `json:"allowedValues,omitempty"`
	Conditions             []conditions.Condition `json:"conditions,omitempty"`
	DataType               *string                `json:"dataType,omitempty"`
	DefaultValue           *string                `json:"defaultValue,omitempty"`
	Description            *string                `json:"description,omitempty"`
	DocumentationLink      *string                `json:"documentationLink,omitempty"`
	Id                     *string                `json:"id,omitempty"`
	IsConfigPendingRestart *bool                  `json:"isConfigPendingRestart,omitempty"`
	IsDynamicConfig        *bool                  `json:"isDynamicConfig,omitempty"`
	IsReadOnly             *bool                  `json:"isReadOnly,omitempty"`
	Name                   *string                `json:"name,omitempty"`
	PropertyBag            genruntime.PropertyBag `json:"$propertyBag,omitempty"`
	Source                 *string                `json:"source,omitempty"`
	SystemData             *SystemData_STATUS     `json:"systemData,omitempty"`
	Type                   *string                `json:"type,omitempty"`
	Unit                   *string                `json:"unit,omitempty"`
	Value                  *string                `json:"value,omitempty"`
}

var _ genruntime.ConvertibleStatus = &FlexibleServersConfiguration_STATUS{}

// ConvertStatusFrom populates our FlexibleServersConfiguration_STATUS from the provided source
func (configuration *FlexibleServersConfiguration_STATUS) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	if source == configuration {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return source.ConvertStatusTo(configuration)
}

// ConvertStatusTo populates the provided destination from our FlexibleServersConfiguration_STATUS
func (configuration *FlexibleServersConfiguration_STATUS) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	if destination == configuration {
		return errors.New("attempted conversion between unrelated implementations of github.com/Azure/azure-service-operator/v2/pkg/genruntime/ConvertibleStatus")
	}

	return destination.ConvertStatusFrom(configuration)
}

// Storage version of v1api20221201.FlexibleServersConfigurationOperatorSpec
// Details for configuring operator behavior. Fields in this struct are interpreted by the operator directly rather than being passed to Azure
type FlexibleServersConfigurationOperatorSpec struct {
	ConfigMapExpressions []*core.DestinationExpression `json:"configMapExpressions,omitempty"`
	PropertyBag          genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	SecretExpressions    []*core.DestinationExpression `json:"secretExpressions,omitempty"`
}

func init() {
	SchemeBuilder.Register(&FlexibleServersConfiguration{}, &FlexibleServersConfigurationList{})
}
