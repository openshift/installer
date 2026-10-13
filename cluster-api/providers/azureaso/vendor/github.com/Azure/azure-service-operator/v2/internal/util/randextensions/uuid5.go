// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package randextensions

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// Namespace is the ASOv2 UUIDv5 namespace UUID
var Namespace = uuid.Must(uuid.Parse("9a329043-7ad7-4b1c-812f-9c7a93d6392a"))

// MakeUniqueOwnerScopedStringLegacy preserves the old buggy behavior with ARM ID owners, for now...
// Deprecated: use MakeUniqueOwnerScopedString instead
func MakeUniqueOwnerScopedStringLegacy(owner *genruntime.ResourceReference, gk schema.GroupKind, namespace string, name string) string {
	var prefix string
	if owner != nil {
		prefix = fmt.Sprintf("%s/%s:%s/%s", owner.Group, owner.Kind, namespace, owner.Name)
	}

	var parts []string
	if prefix != "" {
		parts = append(parts, prefix)
	}

	parts = append(parts, fmt.Sprintf("%s/%s", gk.Kind, gk.Group))
	parts = append(parts, fmt.Sprintf("%s/%s", namespace, name))

	return strings.Join(parts, ":")
}

// MakeUniqueOwnerScopedString generates a string that uniquely identifies a cluster resource.  It includes the
// following distinguishing parts:
// * Owner (either group, kind, namespace, name or raw ARM ID)
// * Group
// * Kind
// * Namespace
// * Name
func MakeUniqueOwnerScopedString(owner *genruntime.ResourceReference, gk schema.GroupKind, namespace string, name string) string {
	prefix := makeUniqueOwnerSubString(owner, namespace)

	var parts []string
	if prefix != "" {
		parts = append(parts, prefix)
	}

	parts = append(parts, fmt.Sprintf("%s/%s", gk.Kind, gk.Group))
	parts = append(parts, fmt.Sprintf("%s/%s", namespace, name))

	return strings.Join(parts, ":")
}

func makeUniqueOwnerSubString(owner *genruntime.ResourceReference, namespace string) string {
	if owner == nil {
		return ""
	}

	if owner.IsDirectARMReference() {
		return owner.ARMID
	}

	if owner.IsKubernetesReference() {
		return fmt.Sprintf("%s/%s:%s/%s", owner.Group, owner.Kind, namespace, owner.Name)
	}

	return "" // This is not expected
}

// MakeUUIDName creates a stable UUID (v5) if the provided name is not already a UUID based on the specified
// uniqueString.
func MakeUUIDName(name string, uniqueString string) string {
	// If name is already a UUID we can just use that
	_, err := uuid.Parse(name)
	if err == nil {
		return name
	}

	return uuid.NewSHA1(Namespace, []byte(uniqueString)).String()
}

// MakeRandomUUID creates a random UUID string. The uuid generated will always be distinct.
func MakeRandomUUID() string {
	// Error is always nil
	newUUID, _ := uuid.NewUUID()
	return newUUID.String()
}
