/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package secrets

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Azure/azure-service-operator/v2/internal/set"
)

func SliceToClientObjectSlice(s []*v1.Secret) []client.Object {
	if s == nil {
		return nil
	}

	result := make([]client.Object, 0, len(s))
	for _, s := range s {
		result = append(result, s)
	}

	return result
}

func SelectSecrets(requested set.Set[string], retrieved map[string]string) map[string]string {
	result := make(map[string]string, len(requested))
	for key, value := range retrieved {
		if requested.Contains(key) {
			result[key] = value
		}
	}

	return result
}
