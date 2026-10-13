/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// AzureNodePoolEncryptionAtHost represents the values of the 'azure_node_pool_encryption_at_host' type.
//
// AzureNodePoolEncryptionAtHost defines the encryption setting for Encryption At Host.
// If not specified, Encryption at Host is not enabled.
type AzureNodePoolEncryptionAtHost = api_v1.AzureNodePoolEncryptionAtHost

// AzureNodePoolEncryptionAtHostListKind is the name of the type used to represent list of objects of
// type 'azure_node_pool_encryption_at_host'.
const AzureNodePoolEncryptionAtHostListKind = api_v1.AzureNodePoolEncryptionAtHostListKind

// AzureNodePoolEncryptionAtHostListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_node_pool_encryption_at_host'.
const AzureNodePoolEncryptionAtHostListLinkKind = api_v1.AzureNodePoolEncryptionAtHostListLinkKind

// AzureNodePoolEncryptionAtHostNilKind is the name of the type used to nil lists of objects of
// type 'azure_node_pool_encryption_at_host'.
const AzureNodePoolEncryptionAtHostListNilKind = api_v1.AzureNodePoolEncryptionAtHostListNilKind

type AzureNodePoolEncryptionAtHostList = api_v1.AzureNodePoolEncryptionAtHostList
