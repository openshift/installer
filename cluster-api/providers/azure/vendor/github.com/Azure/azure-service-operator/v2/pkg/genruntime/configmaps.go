/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"fmt"

	"github.com/pkg/errors"
)

// ConfigMapReference is a reference to a Kubernetes configmap and key in the same namespace as
// the resource it is on.
// +kubebuilder:object:generate=true
type ConfigMapReference struct {
	// Name is the name of the Kubernetes configmap being referenced.
	// The configmap must be in the same namespace as the resource
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`

	// Key is the key in the Kubernetes configmap being referenced
	// +kubebuilder:validation:Required
	Key string `json:"key,omitempty"`
}

var _ Indexer = ConfigMapReference{}

func (c ConfigMapReference) Index() []string {
	return []string{c.Name}
}

// Copy makes an independent copy of the ConfigMapReference
func (c ConfigMapReference) Copy() ConfigMapReference {
	return c
}

func (c ConfigMapReference) String() string {
	return fmt.Sprintf("Name: %q, Key: %q", c.Name, c.Key)
}

// AsNamespacedRef creates a NamespacedSecretReference from this SecretReference in the given namespace
func (c ConfigMapReference) AsNamespacedRef(namespace string) NamespacedConfigMapReference {
	return NamespacedConfigMapReference{
		ConfigMapReference: c,
		Namespace:          namespace,
	}
}

// NamespacedConfigMapReference is a ConfigMapReference with namespace information included
type NamespacedConfigMapReference struct {
	ConfigMapReference
	Namespace string
}

func (s NamespacedConfigMapReference) String() string {
	return fmt.Sprintf("Namespace: %q, %s", s.Namespace, s.ConfigMapReference)
}

// ConfigMapDestination describes the location to store a single configmap value
// Note: This is similar to SecretDestination in secrets.go. Changes to one should likely also be made to the other.
type ConfigMapDestination struct {
	// Note: We could embed ConfigMapReference here, but it makes our life harder because then our reflection based tools will "find" ConfigMapReferences's
	// inside of ConfigMapDestination and try to resolve them. It also gives a worse experience when using the Go Types (the YAML is the same either way).

	// Name is the name of the Kubernetes ConfigMap being referenced.
	// The ConfigMap must be in the same namespace as the resource
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Key is the key in the ConfigMap being referenced
	// +kubebuilder:validation:Required
	Key string `json:"key"`

	// This is a type separate from ConfigMapReference as in the future we may want to support things like
	// customizable annotations or labels, instructions to not delete the ConfigMap when the resource is
	// deleted, etc. None of those things make sense for ConfigMapReference so using the exact same type isn't
	// advisable.
}

// Copy makes an independent copy of the ConfigMapDestination
func (c ConfigMapDestination) Copy() ConfigMapDestination {
	return c
}

func (c ConfigMapDestination) String() string {
	return fmt.Sprintf("Name: %q, Key: %q", c.Name, c.Key)
}

// LookupOptionalConfigMapReferenceValue looks up a ConfigMapReference if it's not nil, or else returns the provided value
func LookupOptionalConfigMapReferenceValue(resolved Resolved[ConfigMapReference, string], ref *ConfigMapReference, value *string) (string, error) {
	if ref == nil && value == nil {
		return "", errors.Errorf("ref and value are both nil")
	}

	if ref != nil && value != nil {
		return "", errors.Errorf("ref and value cannot both be set")
	}

	if ref == nil {
		return *value, nil
	} else {
		return resolved.LookupFromPtr(ref)
	}
}
