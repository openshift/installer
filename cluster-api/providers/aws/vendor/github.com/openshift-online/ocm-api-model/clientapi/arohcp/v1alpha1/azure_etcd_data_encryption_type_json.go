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

// MarshalAzureEtcdDataEncryption writes a value of the 'azure_etcd_data_encryption' type to the given writer.
func MarshalAzureEtcdDataEncryption(object *AzureEtcdDataEncryption, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureEtcdDataEncryption(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureEtcdDataEncryption writes a value of the 'azure_etcd_data_encryption' type to the given stream.
func WriteAzureEtcdDataEncryption(object *AzureEtcdDataEncryption, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.customerManaged != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("customer_managed")
		WriteAzureEtcdDataEncryptionCustomerManaged(object.customerManaged, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_management_mode")
		stream.WriteString(object.keyManagementMode)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureEtcdDataEncryption reads a value of the 'azure_etcd_data_encryption' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureEtcdDataEncryption(source interface{}) (object *AzureEtcdDataEncryption, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureEtcdDataEncryption(iterator)
	err = iterator.Error
	return
}

// ReadAzureEtcdDataEncryption reads a value of the 'azure_etcd_data_encryption' type from the given iterator.
func ReadAzureEtcdDataEncryption(iterator *jsoniter.Iterator) *AzureEtcdDataEncryption {
	object := &AzureEtcdDataEncryption{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "customer_managed":
			value := ReadAzureEtcdDataEncryptionCustomerManaged(iterator)
			object.customerManaged = value
			object.fieldSet_[0] = true
		case "key_management_mode":
			value := iterator.ReadString()
			object.keyManagementMode = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
