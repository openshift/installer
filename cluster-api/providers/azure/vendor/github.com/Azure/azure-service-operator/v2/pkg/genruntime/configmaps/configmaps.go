/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package configmaps

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SliceToClientObjectSlice(s []*v1.ConfigMap) []client.Object {
	result := make([]client.Object, 0, len(s))
	for _, s := range s {
		result = append(result, s)
	}

	return result
}
