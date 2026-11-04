// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package labels

import (
	"strings"

	. "github.com/Azure/azure-service-operator/v2/internal/logging"

	"github.com/go-logr/logr"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

const (
	OwnerNameLabel             = "serviceoperator.azure.com/owner-name"
	OwnerGroupKindLabel        = "serviceoperator.azure.com/owner-group-kind"
	OwnerUIDLabel              = "serviceoperator.azure.com/owner-uid"
	LastReconciledVersionLabel = "serviceoperator.azure.com/last-reconciled-version"
)

// SetOwnerNameLabel sets the owner name label on the given object, or truncates it to 63 characters if it exceeds the character limit.
func SetOwnerNameLabel(logger logr.Logger, obj genruntime.ARMMetaObject) {
	if obj.Owner() != nil && obj.Owner().Name != "" {
		ownerName := obj.Owner().Name
		if len(ownerName) > 63 {
			// Truncate the owner name to 63 characters if it exceeds the limit
			ownerName = ownerName[:63]
			logger.V(Status).Info("WARNING: Owner name label truncated to 63 characters", "ownerName", ownerName)
		}
		genruntime.AddLabel(obj, OwnerNameLabel, ownerName)
	}
}

// SetOwnGroupKindLabel sets the owner group kind label on the given object, or truncates it to 63 characters if it exceeds the character limit.
func SetOwnerGroupKindLabel(logger logr.Logger, obj genruntime.ARMMetaObject) {
	if obj.Owner() != nil && obj.Owner().IsKubernetesReference() {
		groupKind := obj.Owner().GroupKind().String()
		if len(groupKind) > 63 {
			// Truncate the groupKind to 63 characters if it exceeds the limit
			groupKind = groupKind[:63]
			logger.V(Status).Info("WARNING: GroupKind name truncated to 63 characters", "groupKind", groupKind)
		}

		genruntime.AddLabel(obj, OwnerGroupKindLabel, groupKind)
	}
}

// SetOwnerUIDLabel sets the owner UID label on the given object if the owner reference is found in the object's owner references.
func SetOwnerUIDLabel(obj genruntime.ARMMetaObject) {
	ownerRefs := obj.GetOwnerReferences()
	if len(ownerRefs) == 0 || obj.Owner() == nil || !obj.Owner().IsKubernetesReference() {
		return
	}

	groupKind := obj.Owner().GroupKind()

	for _, ref := range ownerRefs {
		ownerGroup := strings.Split(ref.APIVersion, "/")[0]
		if ref.Kind == groupKind.Kind && ownerGroup == groupKind.Group {
			// Set the label with the UID of the owner
			genruntime.AddLabel(obj, OwnerUIDLabel, string(ref.UID))
			return
		}
	}
}
