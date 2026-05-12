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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAzureShard writes a value of the 'azure_shard' type to the given writer.
func MarshalAzureShard(object *AzureShard, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureShard(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureShard writes a value of the 'azure_shard' type to the given stream.
func WriteAzureShard(object *AzureShard, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aks_management_cluster_resource_id")
		stream.WriteString(object.aksManagementClusterResourceId)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cx_managed_identities_key_vault_url")
		stream.WriteString(object.cxManagedIdentitiesKeyVaultUrl)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cx_secrets_key_vault_managed_identity_client_id")
		stream.WriteString(object.cxSecretsKeyVaultManagedIdentityClientId)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cx_secrets_key_vault_url")
		stream.WriteString(object.cxSecretsKeyVaultUrl)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("public_dns_zone_resource_id")
		stream.WriteString(object.publicDnsZoneResourceId)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureShard reads a value of the 'azure_shard' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureShard(source interface{}) (object *AzureShard, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureShard(iterator)
	err = iterator.Error
	return
}

// ReadAzureShard reads a value of the 'azure_shard' type from the given iterator.
func ReadAzureShard(iterator *jsoniter.Iterator) *AzureShard {
	object := &AzureShard{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "aks_management_cluster_resource_id":
			value := iterator.ReadString()
			object.aksManagementClusterResourceId = value
			object.fieldSet_[0] = true
		case "cx_managed_identities_key_vault_url":
			value := iterator.ReadString()
			object.cxManagedIdentitiesKeyVaultUrl = value
			object.fieldSet_[1] = true
		case "cx_secrets_key_vault_managed_identity_client_id":
			value := iterator.ReadString()
			object.cxSecretsKeyVaultManagedIdentityClientId = value
			object.fieldSet_[2] = true
		case "cx_secrets_key_vault_url":
			value := iterator.ReadString()
			object.cxSecretsKeyVaultUrl = value
			object.fieldSet_[3] = true
		case "public_dns_zone_resource_id":
			value := iterator.ReadString()
			object.publicDnsZoneResourceId = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
