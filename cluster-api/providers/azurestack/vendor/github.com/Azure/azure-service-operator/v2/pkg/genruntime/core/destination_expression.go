/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package core

import (
	"fmt"
)

// DestinationExpression is a CEL expression and a destination to store the result in. The destination may
// be a secret or a configmap. The value of the expression is stored at the specified location in
// the destination.
// +kubebuilder:object:generate=true
type DestinationExpression struct {
	// Name is the name of the Kubernetes configmap or secret to write to.
	// The configmap or secret will be created in the same namespace as the resource.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`

	// Key is the key in the ConfigMap or Secret being written to. If the CEL expression in Value returns a string
	// this is required to identify what key to write to. If the CEL expression in Value returns a map[string]string
	// Key must not be set, instead the keys written will be determined dynamically based on the keys of the resulting
	// map[string]string.
	Key string `json:"key,omitempty"`

	// Value is a CEL expression. The CEL expression may return a string or a map[string]string. For more information
	// on CEL in ASO see https://azure.github.io/azure-service-operator/guide/expressions/
	// +kubebuilder:validation:Required
	Value string `json:"value,omitempty"`
}

func (s DestinationExpression) String() string {
	return fmt.Sprintf("Name: %q, Key: %q, Value: %q", s.Name, s.Key, s.Value)
}
