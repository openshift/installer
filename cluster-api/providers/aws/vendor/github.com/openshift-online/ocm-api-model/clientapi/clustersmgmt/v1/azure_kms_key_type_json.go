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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAzureKmsKey writes a value of the 'azure_kms_key' type to the given writer.
func MarshalAzureKmsKey(object *AzureKmsKey, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureKmsKey(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureKmsKey writes a value of the 'azure_kms_key' type to the given stream.
func WriteAzureKmsKey(object *AzureKmsKey, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_name")
		stream.WriteString(object.keyName)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_vault_name")
		stream.WriteString(object.keyVaultName)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_version")
		stream.WriteString(object.keyVersion)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureKmsKey reads a value of the 'azure_kms_key' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureKmsKey(source interface{}) (object *AzureKmsKey, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureKmsKey(iterator)
	err = iterator.Error
	return
}

// ReadAzureKmsKey reads a value of the 'azure_kms_key' type from the given iterator.
func ReadAzureKmsKey(iterator *jsoniter.Iterator) *AzureKmsKey {
	object := &AzureKmsKey{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "key_name":
			value := iterator.ReadString()
			object.keyName = value
			object.fieldSet_[0] = true
		case "key_vault_name":
			value := iterator.ReadString()
			object.keyVaultName = value
			object.fieldSet_[1] = true
		case "key_version":
			value := iterator.ReadString()
			object.keyVersion = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
