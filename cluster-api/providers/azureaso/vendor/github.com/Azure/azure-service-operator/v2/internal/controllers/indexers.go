/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysql "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1"
	postgresql "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1"
	azuresql "github.com/Azure/azure-service-operator/v2/api/sql/v1"
)

// indexMySQLUserPassword an index function for mysql user passwords
func indexMySQLUserPassword(rawObj client.Object) []string {
	obj, ok := rawObj.(*mysql.User)
	if !ok {
		return nil
	}
	if obj.Spec.LocalUser == nil {
		return nil
	}
	if obj.Spec.LocalUser.Password == nil {
		return nil
	}
	return []string{obj.Spec.LocalUser.Password.Name}
}

// indexPostgreSQLUserPassword an index function for postgresql user passwords
func indexPostgreSQLUserPassword(rawObj client.Object) []string {
	obj, ok := rawObj.(*postgresql.User)
	if !ok {
		return nil
	}
	if obj.Spec.LocalUser == nil {
		return nil
	}
	if obj.Spec.LocalUser.Password == nil {
		return nil
	}
	return []string{obj.Spec.LocalUser.Password.Name}
}

// indexAzureSQLUserPassword an index function for azure sql user passwords
func indexAzureSQLUserPassword(rawObj client.Object) []string {
	obj, ok := rawObj.(*azuresql.User)
	if !ok {
		return nil
	}
	if obj.Spec.LocalUser == nil {
		return nil
	}
	if obj.Spec.LocalUser.Password == nil {
		return nil
	}
	return []string{obj.Spec.LocalUser.Password.Name}
}
