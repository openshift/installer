// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package storage

import (
	"encoding/json"

	"github.com/rotisserie/eris"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	storage "github.com/Azure/azure-service-operator/v2/api/dbformysql/v20241230/storage"
)

var _ augmentConversionForMySQLServerIdentity_STATUS = &MySQLServerIdentity_STATUS{}

// AssignPropertiesFrom implements [augmentConversionForMySQLServerIdentity_STATUS].
func (identity *MySQLServerIdentity_STATUS) AssignPropertiesFrom(src *storage.MySQLServerIdentity_STATUS) error {
	// We're lucky in that the JSON representation of storage.UserAssignedIdentity_STATUS
	// is exactly what we want to stash as JSON here

	resultMap := make(map[string]v1.JSON, len(src.UserAssignedIdentities))
	for key, id := range src.UserAssignedIdentities {
		jsonBytes, err := json.Marshal(id)
		if err != nil {
			return eris.Wrapf(
				err,
				"unable to marshal JSON for user assigned identity with key %q in MySQLServerIdentity_STATUS.UserAssignedIdentities",
				key,
			)
		}

		resultMap[key] = v1.JSON{Raw: jsonBytes}
	}

	identity.UserAssignedIdentities = resultMap

	return nil
}

// AssignPropertiesTo implements [augmentConversionForMySQLServerIdentity_STATUS].
func (identity *MySQLServerIdentity_STATUS) AssignPropertiesTo(dst *storage.MySQLServerIdentity_STATUS) error {
	// We're lucky in that the JSON representation of storage.UserAssignedIdentity_STATUS
	// is exactly what's stashed as JSON here
	resultMap := make(map[string]storage.UserAssignedIdentity_STATUS, len(identity.UserAssignedIdentities))
	for key, jsonBytes := range identity.UserAssignedIdentities {
		var id storage.UserAssignedIdentity_STATUS
		err := json.Unmarshal(jsonBytes.Raw, &id)
		if err != nil {
			return eris.Wrapf(
				err,
				"unable to unmarshal JSON for user assigned identity with key %q in MySQLServerIdentity_STATUS.UserAssignedIdentities",
				key,
			)
		}

		resultMap[key] = id
	}

	dst.UserAssignedIdentities = resultMap

	return nil
}
