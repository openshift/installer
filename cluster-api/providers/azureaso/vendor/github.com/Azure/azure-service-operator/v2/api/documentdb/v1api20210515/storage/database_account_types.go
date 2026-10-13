// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package storage

import (
	v20231115s "github.com/Azure/azure-service-operator/v2/api/documentdb/v1api20231115/storage"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
)

// Ensure databaseAccount implements augmentConversionForDatabaseAccount
var _ augmentConversionForDatabaseAccount_Spec = &DatabaseAccount_Spec{}

func (account *DatabaseAccount_Spec) AssignPropertiesFrom(src *v20231115s.DatabaseAccount_Spec) error {
	// Copy any references that point directly to ARM resources into NetworkAclBypassResourceIds
	if len(src.NetworkAclBypassResourceReferences) > 0 {
		ids := make([]string, 0, len(src.NetworkAclBypassResourceReferences))
		for _, ref := range src.NetworkAclBypassResourceReferences {
			ids = append(ids, ref.ARMID)
		}

		account.NetworkAclBypassResourceIds = ids
		account.PropertyBag.Remove("NetworkAclBypassResourceReferences")
	} else {
		account.NetworkAclBypassResourceIds = nil
	}

	return nil
}

func (account *DatabaseAccount_Spec) AssignPropertiesTo(dst *v20231115s.DatabaseAccount_Spec) error {
	// Copy all ARM IDs into NetworkAclBypassResourceReferences
	if len(account.NetworkAclBypassResourceIds) > 0 {
		references := make([]genruntime.ResourceReference, 0, len(account.NetworkAclBypassResourceIds))
		for _, ref := range account.NetworkAclBypassResourceIds {
			krr := genruntime.ResourceReference{
				ARMID: ref,
			}
			references = append(references, krr)
		}

		dst.NetworkAclBypassResourceReferences = references
	} else {
		dst.NetworkAclBypassResourceReferences = nil
	}

	return nil
}
