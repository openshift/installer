// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package labels

import (
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

const (
	OwnerNameLabel      = "serviceoperator.azure.com/owner-name"
	OwnerGroupKindLabel = "serviceoperator.azure.com/owner-group-kind"
)

func SetOwnerNameLabel(obj genruntime.ARMMetaObject) {
	if obj.Owner() != nil && obj.Owner().Name != "" {
		genruntime.AddLabel(obj, OwnerNameLabel, obj.Owner().Name)
	}
}

func SetOwnerGroupKindLabel(obj genruntime.ARMMetaObject) {
	if obj.Owner() != nil && obj.Owner().IsKubernetesReference() {
		genruntime.AddLabel(obj, OwnerGroupKindLabel, obj.Owner().GroupKind().String())
	}
}
