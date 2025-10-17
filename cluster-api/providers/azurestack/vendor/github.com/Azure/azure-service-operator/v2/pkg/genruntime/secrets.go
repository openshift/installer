/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"fmt"
)

// SecretReference is a reference to a Kubernetes secret and key in the same namespace as
// the resource it is on.
// +kubebuilder:object:generate=true
type SecretReference struct {
	// Name is the name of the Kubernetes secret being referenced.
	// The secret must be in the same namespace as the resource
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Key is the key in the Kubernetes secret being referenced
	// +kubebuilder:validation:Required
	Key string `json:"key"`

	// If we end up wanting to support secrets from KeyVault (or elsewhere) we should be able to add a
	// Type *SecretType
	// here and default it to Kubernetes if it's not set. See the secrets design for more details.
}

var _ Indexer = SecretReference{}

func (c SecretReference) Index() []string {
	return []string{c.Name}
}

// Copy makes an independent copy of the SecretReference
func (s SecretReference) Copy() SecretReference {
	return s
}

func (s SecretReference) String() string {
	return fmt.Sprintf("Name: %q, Key: %q", s.Name, s.Key)
}

// AsNamespacedRef creates a NamespacedSecretReference from this SecretReference in the given namespace
func (s SecretReference) AsNamespacedRef(namespace string) NamespacedSecretReference {
	return NamespacedSecretReference{
		SecretReference: s,
		Namespace:       namespace,
	}
}

// NamespacedSecretReference is a SecretReference with namespace information included
type NamespacedSecretReference struct {
	SecretReference
	Namespace string
}

func (s NamespacedSecretReference) String() string {
	return fmt.Sprintf("Namespace: %q, %s", s.Namespace, s.SecretReference)
}

// SecretMapReference is a reference to a Kubernetes secret in the same namespace as
// the resource it is on.
// +kubebuilder:object:generate=true
type SecretMapReference struct {
	// Name is the name of the Kubernetes secret being referenced.
	// The secret must be in the same namespace as the resource
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// If we end up wanting to support secrets from KeyVault (or elsewhere) we should be able to add a
	// Type *SecretType
	// here and default it to Kubernetes if it's not set. See the secrets design for more details.
	// TODO: If we wanted to do this, probably the secret would need to be JSON-formatted, as unlike
	// TODO: Kubernetes secrets there's no map[string]string structure inside of the secret. I think
	// TODO: that's still probably preferable to having a map[string]SecretReference from a UX perspective.
}

var _ Indexer = SecretMapReference{}

func (c SecretMapReference) Index() []string {
	return []string{c.Name}
}

// Copy makes an independent copy of the SecretMapReference
func (s SecretMapReference) Copy() SecretMapReference {
	return s
}

func (s SecretMapReference) String() string {
	return fmt.Sprintf("Name: %q", s.Name)
}

// AsNamespacedRef creates a NamespacedSecretReference from this SecretReference in the given namespace
func (s SecretMapReference) AsNamespacedRef(namespace string) NamespacedSecretMapReference {
	return NamespacedSecretMapReference{
		SecretMapReference: s,
		Namespace:          namespace,
	}
}

// NamespacedSecretMapReference is an SecretMapReference with namespace information included
type NamespacedSecretMapReference struct {
	SecretMapReference
	Namespace string
}

func (s NamespacedSecretMapReference) String() string {
	return fmt.Sprintf("Namespace: %q, %s", s.Namespace, s.SecretMapReference)
}

// SecretDestination describes the location to store a single secret value.
// Note: This is similar to: ConfigMapDestination in configmaps.go.
// Changes to one may need to be made to the others as well.
type SecretDestination struct {
	// Name is the name of the Kubernetes secret to write to.
	// The secret will be created in the same namespace as the resource.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`

	// Key is the key in the Kubernetes secret being referenced.
	// +kubebuilder:validation:Required
	Key string `json:"key,omitempty"`

	// This is a type separate from SecretReference as in the future we may want to support things like
	// customizable annotations or labels, instructions to not delete the secret when the resource is
	// deleted, etc. None of those things make sense for SecretReference so using the exact same type isn't
	// advisable.
}

// Copy makes an independent copy of the SecretDestination
func (s SecretDestination) Copy() SecretDestination {
	return s
}

func (s SecretDestination) String() string {
	return fmt.Sprintf("Name: %q, Key: %q", s.Name, s.Key)
}
