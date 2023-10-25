// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package randextensions

import (
	"fmt"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Namespace is the ASOv2 UUIDv5 namespace UUID
var Namespace = uuid.Must(uuid.Parse("9a329043-7ad7-4b1c-812f-9c7a93d6392a"))

// MakeUniqueResourceString generates a string that uniquely identifies a cluster resource.
func MakeUniqueResourceString(group string, kind string, namespace string, name string) string {
	return fmt.Sprintf("%s/%s:%s/%s", group, kind, namespace, name)
}

// TODO: Fix name
// MakeUniqueResourceString generates a string that uniquely identifies a cluster resource.
func MakeUniqueResourceString2(ownerGK schema.GroupKind, ownerName string, gk schema.GroupKind, namespace string, name string) string {
	return fmt.Sprintf("%s/%s:%s/%s:%s/%s:%s/%s", ownerGK.Group, ownerGK.Kind, namespace, ownerName, gk.Kind, gk.Group, namespace, name)
}

// MakeUUIDName creates a stable UUID (v5) based on the group, kind, namespace, and name of a resource and its owner.
// If the name of the resource is already a compliant UUID, we just use that. If the name is not a UUID one is
// generated and returned.
func MakeUUIDName(ownerGK schema.GroupKind, ownerName string, gk schema.GroupKind, namespace string, name string) string {
	// If name is already a UUID we can just use that
	_, err := uuid.Parse(name)
	if err == nil {
		return name
	}

	uniqueStr := MakeUniqueResourceString2(ownerGK, ownerName, gk, namespace, name)
	return uuid.NewSHA1(Namespace, []byte(uniqueStr)).String()
}
