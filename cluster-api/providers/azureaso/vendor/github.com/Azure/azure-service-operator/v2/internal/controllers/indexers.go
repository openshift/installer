/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package controllers

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	mysql "github.com/Azure/azure-service-operator/v2/api/dbformysql/v1"
	postgresql "github.com/Azure/azure-service-operator/v2/api/dbforpostgresql/v1"
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

// indexPostgreSqlUserPassword an index function for postgresql user passwords
func indexPostgreSqlUserPassword(rawObj client.Object) []string {
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
