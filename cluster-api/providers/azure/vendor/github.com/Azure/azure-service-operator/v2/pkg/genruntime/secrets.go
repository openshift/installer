/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"fmt"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"github.com/Azure/azure-service-operator/v2/internal/set"
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

// SecretDestination describes the location to store a single secret value.
// Note: This is similar to ConfigMapDestination in configmaps.go. Changes to one should likely also be made to the other.
type SecretDestination struct {
	// Note: We could embed SecretReference here, but it makes our life harder because then our reflection based tools will "find" SecretReference's
	// inside of SecretDestination and try to resolve them. It also gives a worse experience when using the Go Types (the YAML is the same either way).

	// Name is the name of the Kubernetes secret being referenced.
	// The secret must be in the same namespace as the resource
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Key is the key in the Kubernetes secret being referenced
	// +kubebuilder:validation:Required
	Key string `json:"key"`

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

type keyPair struct {
	name string
	key  string
}

func makeKeyPairFromSecret(dest *SecretDestination) keyPair {
	return keyPair{
		name: dest.Name,
		key:  dest.Key,
	}
}

// ValidateSecretDestinations checks that no two destinations are writing to the same secret/key, as that could cause
// those secrets to overwrite one another.
func ValidateSecretDestinations(destinations []*SecretDestination) (admission.Warnings, error) {
	// Map of secret -> keys
	locations := set.Make[keyPair]()

	for _, dest := range destinations {
		if dest == nil {
			continue
		}

		pair := makeKeyPairFromSecret(dest)
		if locations.Contains(pair) {
			return nil, errors.Errorf("cannot write more than one secret to destination %s", dest.String())
		}

		locations.Add(pair)
	}

	return nil, nil
}
