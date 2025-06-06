// Code generated by azure-service-operator-codegen. DO NOT EDIT.
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package storage

import (
	"fmt"
	storage "github.com/Azure/azure-service-operator/v2/api/apimanagement/v1api20220801/storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/configmaps"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/secrets"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
// Storage version of v1api20230501preview.ProductApi
// Generator information:
// - Generated from: /apimanagement/resource-manager/Microsoft.ApiManagement/preview/2023-05-01-preview/apimproducts.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/products/{productId}/apis/{apiId}
type ProductApi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProductApi_Spec   `json:"spec,omitempty"`
	Status            ProductApi_STATUS `json:"status,omitempty"`
}

var _ conditions.Conditioner = &ProductApi{}

// GetConditions returns the conditions of the resource
func (productApi *ProductApi) GetConditions() conditions.Conditions {
	return productApi.Status.Conditions
}

// SetConditions sets the conditions on the resource status
func (productApi *ProductApi) SetConditions(conditions conditions.Conditions) {
	productApi.Status.Conditions = conditions
}

var _ conversion.Convertible = &ProductApi{}

// ConvertFrom populates our ProductApi from the provided hub ProductApi
func (productApi *ProductApi) ConvertFrom(hub conversion.Hub) error {
	source, ok := hub.(*storage.ProductApi)
	if !ok {
		return fmt.Errorf("expected apimanagement/v1api20220801/storage/ProductApi but received %T instead", hub)
	}

	return productApi.AssignProperties_From_ProductApi(source)
}

// ConvertTo populates the provided hub ProductApi from our ProductApi
func (productApi *ProductApi) ConvertTo(hub conversion.Hub) error {
	destination, ok := hub.(*storage.ProductApi)
	if !ok {
		return fmt.Errorf("expected apimanagement/v1api20220801/storage/ProductApi but received %T instead", hub)
	}

	return productApi.AssignProperties_To_ProductApi(destination)
}

var _ configmaps.Exporter = &ProductApi{}

// ConfigMapDestinationExpressions returns the Spec.OperatorSpec.ConfigMapExpressions property
func (productApi *ProductApi) ConfigMapDestinationExpressions() []*core.DestinationExpression {
	if productApi.Spec.OperatorSpec == nil {
		return nil
	}
	return productApi.Spec.OperatorSpec.ConfigMapExpressions
}

var _ secrets.Exporter = &ProductApi{}

// SecretDestinationExpressions returns the Spec.OperatorSpec.SecretExpressions property
func (productApi *ProductApi) SecretDestinationExpressions() []*core.DestinationExpression {
	if productApi.Spec.OperatorSpec == nil {
		return nil
	}
	return productApi.Spec.OperatorSpec.SecretExpressions
}

var _ genruntime.KubernetesResource = &ProductApi{}

// AzureName returns the Azure name of the resource
func (productApi *ProductApi) AzureName() string {
	return productApi.Spec.AzureName
}

// GetAPIVersion returns the ARM API version of the resource. This is always "2023-05-01-preview"
func (productApi ProductApi) GetAPIVersion() string {
	return "2023-05-01-preview"
}

// GetResourceScope returns the scope of the resource
func (productApi *ProductApi) GetResourceScope() genruntime.ResourceScope {
	return genruntime.ResourceScopeResourceGroup
}

// GetSpec returns the specification of this resource
func (productApi *ProductApi) GetSpec() genruntime.ConvertibleSpec {
	return &productApi.Spec
}

// GetStatus returns the status of this resource
func (productApi *ProductApi) GetStatus() genruntime.ConvertibleStatus {
	return &productApi.Status
}

// GetSupportedOperations returns the operations supported by the resource
func (productApi *ProductApi) GetSupportedOperations() []genruntime.ResourceOperation {
	return []genruntime.ResourceOperation{
		genruntime.ResourceOperationDelete,
		genruntime.ResourceOperationHead,
		genruntime.ResourceOperationPut,
	}
}

// GetType returns the ARM Type of the resource. This is always "Microsoft.ApiManagement/service/products/apis"
func (productApi *ProductApi) GetType() string {
	return "Microsoft.ApiManagement/service/products/apis"
}

// NewEmptyStatus returns a new empty (blank) status
func (productApi *ProductApi) NewEmptyStatus() genruntime.ConvertibleStatus {
	return &ProductApi_STATUS{}
}

// Owner returns the ResourceReference of the owner
func (productApi *ProductApi) Owner() *genruntime.ResourceReference {
	group, kind := genruntime.LookupOwnerGroupKind(productApi.Spec)
	return productApi.Spec.Owner.AsResourceReference(group, kind)
}

// SetStatus sets the status of this resource
func (productApi *ProductApi) SetStatus(status genruntime.ConvertibleStatus) error {
	// If we have exactly the right type of status, assign it
	if st, ok := status.(*ProductApi_STATUS); ok {
		productApi.Status = *st
		return nil
	}

	// Convert status to required version
	var st ProductApi_STATUS
	err := status.ConvertStatusTo(&st)
	if err != nil {
		return errors.Wrap(err, "failed to convert status")
	}

	productApi.Status = st
	return nil
}

// AssignProperties_From_ProductApi populates our ProductApi from the provided source ProductApi
func (productApi *ProductApi) AssignProperties_From_ProductApi(source *storage.ProductApi) error {

	// ObjectMeta
	productApi.ObjectMeta = *source.ObjectMeta.DeepCopy()

	// Spec
	var spec ProductApi_Spec
	err := spec.AssignProperties_From_ProductApi_Spec(&source.Spec)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_From_ProductApi_Spec() to populate field Spec")
	}
	productApi.Spec = spec

	// Status
	var status ProductApi_STATUS
	err = status.AssignProperties_From_ProductApi_STATUS(&source.Status)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_From_ProductApi_STATUS() to populate field Status")
	}
	productApi.Status = status

	// Invoke the augmentConversionForProductApi interface (if implemented) to customize the conversion
	var productApiAsAny any = productApi
	if augmentedProductApi, ok := productApiAsAny.(augmentConversionForProductApi); ok {
		err := augmentedProductApi.AssignPropertiesFrom(source)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesFrom() for conversion")
		}
	}

	// No error
	return nil
}

// AssignProperties_To_ProductApi populates the provided destination ProductApi from our ProductApi
func (productApi *ProductApi) AssignProperties_To_ProductApi(destination *storage.ProductApi) error {

	// ObjectMeta
	destination.ObjectMeta = *productApi.ObjectMeta.DeepCopy()

	// Spec
	var spec storage.ProductApi_Spec
	err := productApi.Spec.AssignProperties_To_ProductApi_Spec(&spec)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_To_ProductApi_Spec() to populate field Spec")
	}
	destination.Spec = spec

	// Status
	var status storage.ProductApi_STATUS
	err = productApi.Status.AssignProperties_To_ProductApi_STATUS(&status)
	if err != nil {
		return errors.Wrap(err, "calling AssignProperties_To_ProductApi_STATUS() to populate field Status")
	}
	destination.Status = status

	// Invoke the augmentConversionForProductApi interface (if implemented) to customize the conversion
	var productApiAsAny any = productApi
	if augmentedProductApi, ok := productApiAsAny.(augmentConversionForProductApi); ok {
		err := augmentedProductApi.AssignPropertiesTo(destination)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesTo() for conversion")
		}
	}

	// No error
	return nil
}

// OriginalGVK returns a GroupValueKind for the original API version used to create the resource
func (productApi *ProductApi) OriginalGVK() *schema.GroupVersionKind {
	return &schema.GroupVersionKind{
		Group:   GroupVersion.Group,
		Version: productApi.Spec.OriginalVersion,
		Kind:    "ProductApi",
	}
}

// +kubebuilder:object:root=true
// Storage version of v1api20230501preview.ProductApi
// Generator information:
// - Generated from: /apimanagement/resource-manager/Microsoft.ApiManagement/preview/2023-05-01-preview/apimproducts.json
// - ARM URI: /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/products/{productId}/apis/{apiId}
type ProductApiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProductApi `json:"items"`
}

type augmentConversionForProductApi interface {
	AssignPropertiesFrom(src *storage.ProductApi) error
	AssignPropertiesTo(dst *storage.ProductApi) error
}

// Storage version of v1api20230501preview.ProductApi_Spec
type ProductApi_Spec struct {
	// AzureName: The name of the resource in Azure. This is often the same as the name of the resource in Kubernetes but it
	// doesn't have to be.
	AzureName       string                  `json:"azureName,omitempty"`
	OperatorSpec    *ProductApiOperatorSpec `json:"operatorSpec,omitempty"`
	OriginalVersion string                  `json:"originalVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Owner: The owner of the resource. The owner controls where the resource goes when it is deployed. The owner also
	// controls the resources lifecycle. When the owner is deleted the resource will also be deleted. Owner is expected to be a
	// reference to a apimanagement.azure.com/Product resource
	Owner       *genruntime.KnownResourceReference `group:"apimanagement.azure.com" json:"owner,omitempty" kind:"Product"`
	PropertyBag genruntime.PropertyBag             `json:"$propertyBag,omitempty"`
}

var _ genruntime.ConvertibleSpec = &ProductApi_Spec{}

// ConvertSpecFrom populates our ProductApi_Spec from the provided source
func (productApi *ProductApi_Spec) ConvertSpecFrom(source genruntime.ConvertibleSpec) error {
	src, ok := source.(*storage.ProductApi_Spec)
	if ok {
		// Populate our instance from source
		return productApi.AssignProperties_From_ProductApi_Spec(src)
	}

	// Convert to an intermediate form
	src = &storage.ProductApi_Spec{}
	err := src.ConvertSpecFrom(source)
	if err != nil {
		return errors.Wrap(err, "initial step of conversion in ConvertSpecFrom()")
	}

	// Update our instance from src
	err = productApi.AssignProperties_From_ProductApi_Spec(src)
	if err != nil {
		return errors.Wrap(err, "final step of conversion in ConvertSpecFrom()")
	}

	return nil
}

// ConvertSpecTo populates the provided destination from our ProductApi_Spec
func (productApi *ProductApi_Spec) ConvertSpecTo(destination genruntime.ConvertibleSpec) error {
	dst, ok := destination.(*storage.ProductApi_Spec)
	if ok {
		// Populate destination from our instance
		return productApi.AssignProperties_To_ProductApi_Spec(dst)
	}

	// Convert to an intermediate form
	dst = &storage.ProductApi_Spec{}
	err := productApi.AssignProperties_To_ProductApi_Spec(dst)
	if err != nil {
		return errors.Wrap(err, "initial step of conversion in ConvertSpecTo()")
	}

	// Update dst from our instance
	err = dst.ConvertSpecTo(destination)
	if err != nil {
		return errors.Wrap(err, "final step of conversion in ConvertSpecTo()")
	}

	return nil
}

// AssignProperties_From_ProductApi_Spec populates our ProductApi_Spec from the provided source ProductApi_Spec
func (productApi *ProductApi_Spec) AssignProperties_From_ProductApi_Spec(source *storage.ProductApi_Spec) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(source.PropertyBag)

	// AzureName
	productApi.AzureName = source.AzureName

	// OperatorSpec
	if source.OperatorSpec != nil {
		var operatorSpec ProductApiOperatorSpec
		err := operatorSpec.AssignProperties_From_ProductApiOperatorSpec(source.OperatorSpec)
		if err != nil {
			return errors.Wrap(err, "calling AssignProperties_From_ProductApiOperatorSpec() to populate field OperatorSpec")
		}
		productApi.OperatorSpec = &operatorSpec
	} else {
		productApi.OperatorSpec = nil
	}

	// OriginalVersion
	productApi.OriginalVersion = source.OriginalVersion

	// Owner
	if source.Owner != nil {
		owner := source.Owner.Copy()
		productApi.Owner = &owner
	} else {
		productApi.Owner = nil
	}

	// Update the property bag
	if len(propertyBag) > 0 {
		productApi.PropertyBag = propertyBag
	} else {
		productApi.PropertyBag = nil
	}

	// Invoke the augmentConversionForProductApi_Spec interface (if implemented) to customize the conversion
	var productApiAsAny any = productApi
	if augmentedProductApi, ok := productApiAsAny.(augmentConversionForProductApi_Spec); ok {
		err := augmentedProductApi.AssignPropertiesFrom(source)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesFrom() for conversion")
		}
	}

	// No error
	return nil
}

// AssignProperties_To_ProductApi_Spec populates the provided destination ProductApi_Spec from our ProductApi_Spec
func (productApi *ProductApi_Spec) AssignProperties_To_ProductApi_Spec(destination *storage.ProductApi_Spec) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(productApi.PropertyBag)

	// AzureName
	destination.AzureName = productApi.AzureName

	// OperatorSpec
	if productApi.OperatorSpec != nil {
		var operatorSpec storage.ProductApiOperatorSpec
		err := productApi.OperatorSpec.AssignProperties_To_ProductApiOperatorSpec(&operatorSpec)
		if err != nil {
			return errors.Wrap(err, "calling AssignProperties_To_ProductApiOperatorSpec() to populate field OperatorSpec")
		}
		destination.OperatorSpec = &operatorSpec
	} else {
		destination.OperatorSpec = nil
	}

	// OriginalVersion
	destination.OriginalVersion = productApi.OriginalVersion

	// Owner
	if productApi.Owner != nil {
		owner := productApi.Owner.Copy()
		destination.Owner = &owner
	} else {
		destination.Owner = nil
	}

	// Update the property bag
	if len(propertyBag) > 0 {
		destination.PropertyBag = propertyBag
	} else {
		destination.PropertyBag = nil
	}

	// Invoke the augmentConversionForProductApi_Spec interface (if implemented) to customize the conversion
	var productApiAsAny any = productApi
	if augmentedProductApi, ok := productApiAsAny.(augmentConversionForProductApi_Spec); ok {
		err := augmentedProductApi.AssignPropertiesTo(destination)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesTo() for conversion")
		}
	}

	// No error
	return nil
}

// Storage version of v1api20230501preview.ProductApi_STATUS
type ProductApi_STATUS struct {
	Conditions  []conditions.Condition `json:"conditions,omitempty"`
	PropertyBag genruntime.PropertyBag `json:"$propertyBag,omitempty"`
}

var _ genruntime.ConvertibleStatus = &ProductApi_STATUS{}

// ConvertStatusFrom populates our ProductApi_STATUS from the provided source
func (productApi *ProductApi_STATUS) ConvertStatusFrom(source genruntime.ConvertibleStatus) error {
	src, ok := source.(*storage.ProductApi_STATUS)
	if ok {
		// Populate our instance from source
		return productApi.AssignProperties_From_ProductApi_STATUS(src)
	}

	// Convert to an intermediate form
	src = &storage.ProductApi_STATUS{}
	err := src.ConvertStatusFrom(source)
	if err != nil {
		return errors.Wrap(err, "initial step of conversion in ConvertStatusFrom()")
	}

	// Update our instance from src
	err = productApi.AssignProperties_From_ProductApi_STATUS(src)
	if err != nil {
		return errors.Wrap(err, "final step of conversion in ConvertStatusFrom()")
	}

	return nil
}

// ConvertStatusTo populates the provided destination from our ProductApi_STATUS
func (productApi *ProductApi_STATUS) ConvertStatusTo(destination genruntime.ConvertibleStatus) error {
	dst, ok := destination.(*storage.ProductApi_STATUS)
	if ok {
		// Populate destination from our instance
		return productApi.AssignProperties_To_ProductApi_STATUS(dst)
	}

	// Convert to an intermediate form
	dst = &storage.ProductApi_STATUS{}
	err := productApi.AssignProperties_To_ProductApi_STATUS(dst)
	if err != nil {
		return errors.Wrap(err, "initial step of conversion in ConvertStatusTo()")
	}

	// Update dst from our instance
	err = dst.ConvertStatusTo(destination)
	if err != nil {
		return errors.Wrap(err, "final step of conversion in ConvertStatusTo()")
	}

	return nil
}

// AssignProperties_From_ProductApi_STATUS populates our ProductApi_STATUS from the provided source ProductApi_STATUS
func (productApi *ProductApi_STATUS) AssignProperties_From_ProductApi_STATUS(source *storage.ProductApi_STATUS) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(source.PropertyBag)

	// Conditions
	productApi.Conditions = genruntime.CloneSliceOfCondition(source.Conditions)

	// Update the property bag
	if len(propertyBag) > 0 {
		productApi.PropertyBag = propertyBag
	} else {
		productApi.PropertyBag = nil
	}

	// Invoke the augmentConversionForProductApi_STATUS interface (if implemented) to customize the conversion
	var productApiAsAny any = productApi
	if augmentedProductApi, ok := productApiAsAny.(augmentConversionForProductApi_STATUS); ok {
		err := augmentedProductApi.AssignPropertiesFrom(source)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesFrom() for conversion")
		}
	}

	// No error
	return nil
}

// AssignProperties_To_ProductApi_STATUS populates the provided destination ProductApi_STATUS from our ProductApi_STATUS
func (productApi *ProductApi_STATUS) AssignProperties_To_ProductApi_STATUS(destination *storage.ProductApi_STATUS) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(productApi.PropertyBag)

	// Conditions
	destination.Conditions = genruntime.CloneSliceOfCondition(productApi.Conditions)

	// Update the property bag
	if len(propertyBag) > 0 {
		destination.PropertyBag = propertyBag
	} else {
		destination.PropertyBag = nil
	}

	// Invoke the augmentConversionForProductApi_STATUS interface (if implemented) to customize the conversion
	var productApiAsAny any = productApi
	if augmentedProductApi, ok := productApiAsAny.(augmentConversionForProductApi_STATUS); ok {
		err := augmentedProductApi.AssignPropertiesTo(destination)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesTo() for conversion")
		}
	}

	// No error
	return nil
}

type augmentConversionForProductApi_Spec interface {
	AssignPropertiesFrom(src *storage.ProductApi_Spec) error
	AssignPropertiesTo(dst *storage.ProductApi_Spec) error
}

type augmentConversionForProductApi_STATUS interface {
	AssignPropertiesFrom(src *storage.ProductApi_STATUS) error
	AssignPropertiesTo(dst *storage.ProductApi_STATUS) error
}

// Storage version of v1api20230501preview.ProductApiOperatorSpec
// Details for configuring operator behavior. Fields in this struct are interpreted by the operator directly rather than being passed to Azure
type ProductApiOperatorSpec struct {
	ConfigMapExpressions []*core.DestinationExpression `json:"configMapExpressions,omitempty"`
	PropertyBag          genruntime.PropertyBag        `json:"$propertyBag,omitempty"`
	SecretExpressions    []*core.DestinationExpression `json:"secretExpressions,omitempty"`
}

// AssignProperties_From_ProductApiOperatorSpec populates our ProductApiOperatorSpec from the provided source ProductApiOperatorSpec
func (operator *ProductApiOperatorSpec) AssignProperties_From_ProductApiOperatorSpec(source *storage.ProductApiOperatorSpec) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(source.PropertyBag)

	// ConfigMapExpressions
	if source.ConfigMapExpressions != nil {
		configMapExpressionList := make([]*core.DestinationExpression, len(source.ConfigMapExpressions))
		for configMapExpressionIndex, configMapExpressionItem := range source.ConfigMapExpressions {
			// Shadow the loop variable to avoid aliasing
			configMapExpressionItem := configMapExpressionItem
			if configMapExpressionItem != nil {
				configMapExpression := *configMapExpressionItem.DeepCopy()
				configMapExpressionList[configMapExpressionIndex] = &configMapExpression
			} else {
				configMapExpressionList[configMapExpressionIndex] = nil
			}
		}
		operator.ConfigMapExpressions = configMapExpressionList
	} else {
		operator.ConfigMapExpressions = nil
	}

	// SecretExpressions
	if source.SecretExpressions != nil {
		secretExpressionList := make([]*core.DestinationExpression, len(source.SecretExpressions))
		for secretExpressionIndex, secretExpressionItem := range source.SecretExpressions {
			// Shadow the loop variable to avoid aliasing
			secretExpressionItem := secretExpressionItem
			if secretExpressionItem != nil {
				secretExpression := *secretExpressionItem.DeepCopy()
				secretExpressionList[secretExpressionIndex] = &secretExpression
			} else {
				secretExpressionList[secretExpressionIndex] = nil
			}
		}
		operator.SecretExpressions = secretExpressionList
	} else {
		operator.SecretExpressions = nil
	}

	// Update the property bag
	if len(propertyBag) > 0 {
		operator.PropertyBag = propertyBag
	} else {
		operator.PropertyBag = nil
	}

	// Invoke the augmentConversionForProductApiOperatorSpec interface (if implemented) to customize the conversion
	var operatorAsAny any = operator
	if augmentedOperator, ok := operatorAsAny.(augmentConversionForProductApiOperatorSpec); ok {
		err := augmentedOperator.AssignPropertiesFrom(source)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesFrom() for conversion")
		}
	}

	// No error
	return nil
}

// AssignProperties_To_ProductApiOperatorSpec populates the provided destination ProductApiOperatorSpec from our ProductApiOperatorSpec
func (operator *ProductApiOperatorSpec) AssignProperties_To_ProductApiOperatorSpec(destination *storage.ProductApiOperatorSpec) error {
	// Clone the existing property bag
	propertyBag := genruntime.NewPropertyBag(operator.PropertyBag)

	// ConfigMapExpressions
	if operator.ConfigMapExpressions != nil {
		configMapExpressionList := make([]*core.DestinationExpression, len(operator.ConfigMapExpressions))
		for configMapExpressionIndex, configMapExpressionItem := range operator.ConfigMapExpressions {
			// Shadow the loop variable to avoid aliasing
			configMapExpressionItem := configMapExpressionItem
			if configMapExpressionItem != nil {
				configMapExpression := *configMapExpressionItem.DeepCopy()
				configMapExpressionList[configMapExpressionIndex] = &configMapExpression
			} else {
				configMapExpressionList[configMapExpressionIndex] = nil
			}
		}
		destination.ConfigMapExpressions = configMapExpressionList
	} else {
		destination.ConfigMapExpressions = nil
	}

	// SecretExpressions
	if operator.SecretExpressions != nil {
		secretExpressionList := make([]*core.DestinationExpression, len(operator.SecretExpressions))
		for secretExpressionIndex, secretExpressionItem := range operator.SecretExpressions {
			// Shadow the loop variable to avoid aliasing
			secretExpressionItem := secretExpressionItem
			if secretExpressionItem != nil {
				secretExpression := *secretExpressionItem.DeepCopy()
				secretExpressionList[secretExpressionIndex] = &secretExpression
			} else {
				secretExpressionList[secretExpressionIndex] = nil
			}
		}
		destination.SecretExpressions = secretExpressionList
	} else {
		destination.SecretExpressions = nil
	}

	// Update the property bag
	if len(propertyBag) > 0 {
		destination.PropertyBag = propertyBag
	} else {
		destination.PropertyBag = nil
	}

	// Invoke the augmentConversionForProductApiOperatorSpec interface (if implemented) to customize the conversion
	var operatorAsAny any = operator
	if augmentedOperator, ok := operatorAsAny.(augmentConversionForProductApiOperatorSpec); ok {
		err := augmentedOperator.AssignPropertiesTo(destination)
		if err != nil {
			return errors.Wrap(err, "calling augmented AssignPropertiesTo() for conversion")
		}
	}

	// No error
	return nil
}

type augmentConversionForProductApiOperatorSpec interface {
	AssignPropertiesFrom(src *storage.ProductApiOperatorSpec) error
	AssignPropertiesTo(dst *storage.ProductApiOperatorSpec) error
}

func init() {
	SchemeBuilder.Register(&ProductApi{}, &ProductApiList{})
}
