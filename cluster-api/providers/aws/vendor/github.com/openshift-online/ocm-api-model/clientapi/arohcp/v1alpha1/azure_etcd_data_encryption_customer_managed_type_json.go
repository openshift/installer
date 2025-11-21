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

// MarshalAzureEtcdDataEncryptionCustomerManaged writes a value of the 'azure_etcd_data_encryption_customer_managed' type to the given writer.
func MarshalAzureEtcdDataEncryptionCustomerManaged(object *AzureEtcdDataEncryptionCustomerManaged, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureEtcdDataEncryptionCustomerManaged(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureEtcdDataEncryptionCustomerManaged writes a value of the 'azure_etcd_data_encryption_customer_managed' type to the given stream.
func WriteAzureEtcdDataEncryptionCustomerManaged(object *AzureEtcdDataEncryptionCustomerManaged, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("encryption_type")
		stream.WriteString(object.encryptionType)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.kms != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms")
		WriteAzureKmsEncryption(object.kms, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureEtcdDataEncryptionCustomerManaged reads a value of the 'azure_etcd_data_encryption_customer_managed' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureEtcdDataEncryptionCustomerManaged(source interface{}) (object *AzureEtcdDataEncryptionCustomerManaged, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureEtcdDataEncryptionCustomerManaged(iterator)
	err = iterator.Error
	return
}

// ReadAzureEtcdDataEncryptionCustomerManaged reads a value of the 'azure_etcd_data_encryption_customer_managed' type from the given iterator.
func ReadAzureEtcdDataEncryptionCustomerManaged(iterator *jsoniter.Iterator) *AzureEtcdDataEncryptionCustomerManaged {
	object := &AzureEtcdDataEncryptionCustomerManaged{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "encryption_type":
			value := iterator.ReadString()
			object.encryptionType = value
			object.fieldSet_[0] = true
		case "kms":
			value := ReadAzureKmsEncryption(iterator)
			object.kms = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
