/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

// +kubebuilder:validation:Optional
package genruntime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/Azure/azure-service-operator/v2/internal/set"
)

// KnownResourceReference is a resource reference to a known type.
// +kubebuilder:object:generate=true
type KnownResourceReference struct {
	// TODO: In practice this type is used only for Owner fields and so might more appropriately have been called OwnerReference
	// TODO: but changing to that would be a breaking change so avoiding it.

	// This is the name of the Kubernetes resource to reference.
	Name string `json:"name,omitempty"`

	// References across namespaces are not supported.
	// Note that ownership across namespaces in Kubernetes is not allowed, but technically resource
	// references are. There are RBAC considerations here though so probably easier to just start by
	// disallowing cross-namespace references for now

	// +kubebuilder:validation:Pattern="(?i)(^(/subscriptions/([^/]+)(/resourcegroups/([^/]+))?)?/providers/([^/]+)/([^/]+/[^/]+)(/([^/]+/[^/]+))*$|^/subscriptions/([^/]+)(/resourcegroups/([^/]+))?$)"
	ARMID string `json:"armId,omitempty"`
}

// AsResourceReference transforms this KnownResourceReference into a ResourceReference
func (ref KnownResourceReference) AsResourceReference(group string, kind string) *ResourceReference {
	if ref.Name == "" {
		group = ""
		kind = ""
	}

	return &ResourceReference{
		Group: group,
		Kind:  kind,
		Name:  ref.Name,
		ARMID: ref.ARMID,
	}
}

// KubernetesOwnerReference is a resource reference to a known type in Kuberentes. Most types support
// ARM references as well but some (such as SQL users) do not.
// +kubebuilder:object:generate=true
type KubernetesOwnerReference struct {
	// +kubebuilder:validation:Required
	// This is the name of the Kubernetes resource to reference.
	Name string `json:"name,omitempty"`
}

// AsResourceReference transforms this KnownResourceReference into a ResourceReference
func (ref KubernetesOwnerReference) AsResourceReference(group string, kind string) *ResourceReference {
	return &ResourceReference{
		Group: group,
		Kind:  kind,
		Name:  ref.Name,
	}
}

// TODO: This type and ResourceReference are almost exactly the same now...
// ArbitraryOwnerReference is an owner reference to an unknown type.
// +kubebuilder:object:generate=true
type ArbitraryOwnerReference struct {
	// This is the name of the Kubernetes resource to reference.
	Name string `json:"name,omitempty"`

	// Group is the Kubernetes group of the resource.
	Group string `json:"group,omitempty"`

	// Kind is the Kubernetes kind of the resource.
	Kind string `json:"kind,omitempty"`

	// Ownership across namespaces is not supported.
	// +kubebuilder:validation:Pattern="(?i)(^(/subscriptions/([^/]+)(/resourcegroups/([^/]+))?)?/providers/([^/]+)/([^/]+/[^/]+)(/([^/]+/[^/]+))*$|^/subscriptions/([^/]+)(/resourcegroups/([^/]+))?$)"
	ARMID string `json:"armId,omitempty"`
}

// AsResourceReference transforms this ArbitraryOwnerReference into a ResourceReference
func (ref ArbitraryOwnerReference) AsResourceReference() *ResourceReference {
	return &ResourceReference{
		Group: ref.Group,
		Kind:  ref.Kind,
		Name:  ref.Name,
		ARMID: ref.ARMID,
	}
}

var _ fmt.Stringer = ResourceReference{}

// ResourceReference represents a resource reference, either to a Kubernetes resource or directly to an Azure resource via ARMID
// +kubebuilder:object:generate=true
type ResourceReference struct {
	// Group is the Kubernetes group of the resource.
	Group string `json:"group,omitempty"`
	// Kind is the Kubernetes kind of the resource.
	Kind string `json:"kind,omitempty"`
	// Name is the Kubernetes name of the resource.
	Name string `json:"name,omitempty"`

	// Note: Version is not required here because references are all about linking one Kubernetes
	// resource to another, and Kubernetes resources are uniquely identified by group, kind, (optionally namespace) and
	// name - the versions are just giving a different view on the same resource
	// Here are some test patterns for it: https://regex101.com/r/K7l3sv/1

	// +kubebuilder:validation:Pattern="(?i)(^(/subscriptions/([^/]+)(/resourcegroups/([^/]+))?)?/providers/([^/]+)/([^/]+/[^/]+)(/([^/]+/[^/]+))*$|^/subscriptions/([^/]+)(/resourcegroups/([^/]+))?$)"
	// ARMID is a string of the form /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}.
	// The /resourcegroups/{resourceGroupName} bit is optional as some resources are scoped at the subscription level
	// ARMID is mutually exclusive with Group, Kind, Namespace and Name.
	ARMID string `json:"armId,omitempty"`
}

// CreateResourceReferenceFromARMID creates a new ResourceReference from a string representing an ARM ID
func CreateResourceReferenceFromARMID(armID string) ResourceReference {
	return ResourceReference{
		ARMID: armID,
	}
}

// IsDirectARMReference returns true if this ResourceReference is referring to an ARMID directly.
func (ref ResourceReference) IsDirectARMReference() bool {
	return ref.ARMID != "" && ref.Name == "" && ref.Group == "" && ref.Kind == ""
}

// IsKubernetesReference returns true if this ResourceReference is referring to a Kubernetes resource.
func (ref ResourceReference) IsKubernetesReference() bool {
	return ref.ARMID == "" && ref.Name != "" && ref.Group != "" && ref.Kind != ""
}

func (ref ResourceReference) String() string {
	if ref.IsDirectARMReference() {
		return ref.ARMID
	}

	if ref.IsKubernetesReference() {
		return fmt.Sprintf("%s, Group/Kind: %s/%s", ref.Name, ref.Group, ref.Kind)
	}

	// Printing all the fields here just in case something weird happens and we have an ARMID and also Kubernetes reference stuff
	return fmt.Sprintf("Group: %q, Kind: %q, Name: %q, ARMID: %q", ref.Group, ref.Kind, ref.Name, ref.ARMID)
}

// TODO: We wouldn't need this if controller-gen supported DUs or OneOf better, see: https://github.com/kubernetes-sigs/controller-tools/issues/461
// Validate validates the ResourceReference to ensure that it is structurally valid.
func (ref ResourceReference) Validate() (admission.Warnings, error) {
	if ref.ARMID == "" && ref.Name == "" && ref.Group == "" && ref.Kind == "" {
		return nil, errors.Errorf("at least one of ['ARMID'] or ['Group', 'Kind', 'Namespace', 'Name'] must be set for ResourceReference")
	}

	if ref.ARMID != "" && !ref.IsDirectARMReference() {
		return nil, errors.Errorf("the 'ARMID' field is mutually exclusive with 'Group', 'Kind', 'Namespace', and 'Name' for ResourceReference: %s", ref.String())
	}

	if ref.ARMID == "" && !ref.IsKubernetesReference() {
		return nil, errors.Errorf("when referencing a Kubernetes resource, 'Group', 'Kind', 'Namespace', and 'Name' must all be specified for ResourceReference: %s", ref.String())
	}

	return nil, nil
}

// AsNamespacedRef creates a NamespacedResourceReference from this reference.
func (ref ResourceReference) AsNamespacedRef(namespace string) NamespacedResourceReference {
	// If this is a direct ARM reference, don't append a namespace as it reads weird
	if ref.IsDirectARMReference() {
		return NamespacedResourceReference{
			ResourceReference: ref,
		}
	}

	return NamespacedResourceReference{
		ResourceReference: ref,
		Namespace:         namespace,
	}
}

// GroupKind returns the GroupKind of the resource reference
func (ref ResourceReference) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: ref.Group,
		Kind:  ref.Kind,
	}
}

// LookupOwnerGroupKind looks up an owners group and kind annotations using reflection.
// This is primarily used to convert from a KnownResourceReference to the more general
// ResourceReference
func LookupOwnerGroupKind(v interface{}) (string, string) {
	t := reflect.TypeOf(v)
	field, _ := t.FieldByName("Owner")

	group, ok := field.Tag.Lookup("group")
	if !ok {
		panic("Couldn't find owner group tag")
	}
	kind, ok := field.Tag.Lookup("kind")
	if !ok {
		panic("Couldn't find %s owner kind tag")
	}

	return group, kind
}

// Copy makes an independent copy of the KnownResourceReference
func (ref KnownResourceReference) Copy() KnownResourceReference {
	return ref
}

// Copy makes an independent copy of the ArbitraryOwnerReference
func (ref ArbitraryOwnerReference) Copy() ArbitraryOwnerReference {
	return ref
}

// Copy makes an independent copy of the ResourceReference
func (ref ResourceReference) Copy() ResourceReference {
	return ref
}

// Copy makes an independent copy of the KubernetesOwnerReference
func (ref KubernetesOwnerReference) Copy() KubernetesOwnerReference {
	return ref
}

// ValidateResourceReferences calls Validate on each ResourceReference
func ValidateResourceReferences(refs set.Set[ResourceReference]) (admission.Warnings, error) {
	errs := make([]error, 0, len(refs))
	var warnings admission.Warnings
	for ref := range refs {
		warning, err := ref.Validate()
		if warning != nil {
			warnings = append(warnings, warning...)
		}
		if err != nil {
			errs = append(errs, err)
		}
	}

	return nil, kerrors.NewAggregate(errs)
}

func VerifyResourceOwnerARMID(resource ARMMetaObject) error {
	owner := resource.Owner()
	if owner == nil {
		return nil
	}
	if !owner.IsDirectARMReference() {
		return nil
	}

	armID, err := arm.ParseResourceID(owner.ARMID)
	if err != nil {
		return err
	}

	provider, rootResourceTypes, err := GetResourceTypeAndProvider(resource)
	if err != nil {
		return err
	}
	expectedResourceTypesIncludedInARMID := rootResourceTypes[:len(rootResourceTypes)-1]

	// Ensure that the ARM ID actually has a suffix containing the resource types we expect
	if len(expectedResourceTypesIncludedInARMID) > 0 {
		if strings.ToLower(armID.ResourceType.Namespace) != strings.ToLower(provider) {
			return errors.Errorf(
				"expected owner ARM ID to be from provider %q, but was %q",
				provider,
				armID.ResourceType.Namespace)
		}
		expectedARMIDType := strings.Join(expectedResourceTypesIncludedInARMID, "/")
		if !strings.EqualFold(armID.ResourceType.Type, expectedARMIDType) {
			return errors.Errorf(
				"expected owner ARM ID to be of type %q, but was %q",
				fmt.Sprintf("%s/%s", provider, expectedARMIDType),
				armID.ResourceType.String())
		}
	} else if len(expectedResourceTypesIncludedInARMID) == 0 {
		scope := resource.GetResourceScope()
		if scope == ResourceScopeResourceGroup && armID.ResourceType.String() != "Microsoft.Resources/resourceGroups" {
			return errors.Errorf(
				"expected owner ARM ID to be for a resource group, but was %q",
				armID.ResourceType.String())
		}
	}

	return nil
}

// ValidateOwner calls Validate on the resource Owner
func ValidateOwner(obj ARMMetaObject) (admission.Warnings, error) {
	owner := obj.Owner()
	if owner == nil {
		return nil, nil
	}

	var warningsResult admission.Warnings

	warnings, err := owner.Validate()
	warningsResult = append(warningsResult, warnings...)
	if err != nil {
		return warningsResult, err
	}

	err = VerifyResourceOwnerARMID(obj)
	if err != nil {
		return warningsResult, err
	}

	return warningsResult, nil
}

// NamespacedResourceReference is a resource reference with namespace information included
type NamespacedResourceReference struct {
	ResourceReference
	Namespace string
}
