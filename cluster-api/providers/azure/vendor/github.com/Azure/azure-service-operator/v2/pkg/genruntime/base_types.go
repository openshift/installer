/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

type ResourceScope string

const (
	// ResourceScopeLocation is a resource that is deployed into a location
	ResourceScopeLocation = ResourceScope("location")
	// ResourceScopeResourceGroup is a resource that is deployed into a resource group
	ResourceScopeResourceGroup = ResourceScope("resourcegroup")
	// ResourceScopeExtension is an extension resource. Extension resources can have any resource as their parent.
	ResourceScopeExtension = ResourceScope("extension")
	// ResourceScopeTenant is an Azure resource rooted to the tenant (examples include subscription, managementGroup, etc)
	ResourceScopeTenant = ResourceScope("tenant")
)

type ResourceOperation string

const (
	ResourceOperationGet    = ResourceOperation("GET")
	ResourceOperationHead   = ResourceOperation("HEAD")
	ResourceOperationPut    = ResourceOperation("PUT")
	ResourceOperationDelete = ResourceOperation("DELETE")
)

func (o ResourceOperation) IsSupportedBy(obj SupportedResourceOperations) bool {
	for _, op := range obj.GetSupportedOperations() {
		if op == o {
			return true
		}
	}

	return false
}

// TODO: It's weird that this is isn't with the other annotations
// TODO: Should we move them all here (so they're exported?) Or shold we move them
// TODO: to serviceoperator-internal.azure.com to signify they are internal?
const (
	ResourceIDAnnotation = "serviceoperator.azure.com/resource-id"

	// ChildResourceIDOverrideAnnotation is an annotation that can be used to force child resources
	// to be owned by a different resource ID than it would normally. This is primarily used for
	// resources like SubscriptionAlias + Subscription, where the create API doesn't use the same
	// ResourceID as needed by child resources of the subscription.
	// When present, this takes precedent over the resources AzureName() and Type.
	// TODO: Currently this annotation can only be used on the root resource in a resource hierarchy.
	// TODO: For example if A owns B owns C, this annotation can be used on A but not on B or C.
	ChildResourceIDOverrideAnnotation = "serviceoperator.azure.com/child-resource-id-override"
)

// MetaObject represents an arbitrary ASO custom resource
type MetaObject interface {
	runtime.Object
	metav1.Object
	conditions.Conditioner
}

// ARMMetaObject represents an arbitrary ASO resource that is an ARM resource
type ARMMetaObject interface {
	MetaObject
	KubernetesResource
}

// ARMOwnedMetaObject represents an arbitrary ASO resource that is owned by an ARM resource
type ARMOwnedMetaObject interface {
	MetaObject
	ARMOwned
}

// AddAnnotation adds the specified annotation to the object.
// Empty string annotations are not allowed. Attempting to add an annotation with a value
// of empty string will result in the removal of that annotation.
func AddAnnotation(obj MetaObject, k string, v string) {
	annotations := obj.GetAnnotations()
	annotations = AddToMap(annotations, k, v)
	obj.SetAnnotations(annotations)
}

// RemoveAnnotation removes the specified annotation from the object
func RemoveAnnotation(obj MetaObject, k string) {
	AddAnnotation(obj, k, "")
}

// AddLabel adds the specified label to the object.
// Empty string labels are not allowed. Attempting to add a label with a value
// of empty string will result in the removal of that label.
func AddLabel(obj MetaObject, k string, v string) {
	labels := obj.GetLabels()
	labels = AddToMap(labels, k, v)
	obj.SetLabels(labels)
}

func AddToMap(m map[string]string, k string, v string) map[string]string {
	if m == nil {
		m = map[string]string{}
	}
	// I think this is the behavior we want...
	if v == "" {
		delete(m, k)
	} else {
		m[k] = v
	}
	return m
}

// RemoveLabel removes the specified label from the object
func RemoveLabel(obj MetaObject, k string) {
	AddLabel(obj, k, "")
}

// ARMResourceSpec is an ARM resource specification. This interface contains
// methods to access properties common to all ARM Resource Specs. An Azure
// Deployment is made of these.
type ARMResourceSpec interface {
	GetAPIVersion() string

	GetType() string

	GetName() string
}

// ARMResourceStatus is an ARM resource status
type ARMResourceStatus interface { // TODO: Unsure what the actual content of this interface needs to be.
	// TODO: We need to define it and generate the code for it
	// GetId() string
}

type ARMResource interface {
	Spec() ARMResourceSpec
	Status() ARMResourceStatus

	GetID() string // TODO: Should this be on Status instead?
}

func NewARMResource(spec ARMResourceSpec, status ARMResourceStatus, id string) ARMResource {
	return &armResourceImpl{
		spec:   spec,
		status: status,
		Id:     id,
	}
}

type armResourceImpl struct {
	spec   ARMResourceSpec
	status ARMResourceStatus
	Id     string
}

var _ ARMResource = &armResourceImpl{}

func (resource *armResourceImpl) Spec() ARMResourceSpec {
	return resource.spec
}

func (resource *armResourceImpl) Status() ARMResourceStatus {
	return resource.status
}

func (resource *armResourceImpl) GetID() string {
	return resource.Id
}

// GetReadyCondition gets the ready condition from the object
func GetReadyCondition(obj conditions.Conditioner) *conditions.Condition {
	for _, c := range obj.GetConditions() {
		if c.Type == conditions.ConditionTypeReady {
			return &c
		}
	}

	return nil
}
