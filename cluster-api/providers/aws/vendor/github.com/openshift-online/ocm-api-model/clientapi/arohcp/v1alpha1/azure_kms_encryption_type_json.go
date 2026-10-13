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

// MarshalAzureKmsEncryption writes a value of the 'azure_kms_encryption' type to the given writer.
func MarshalAzureKmsEncryption(object *AzureKmsEncryption, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureKmsEncryption(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureKmsEncryption writes a value of the 'azure_kms_encryption' type to the given stream.
func WriteAzureKmsEncryption(object *AzureKmsEncryption, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.activeKey != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("active_key")
		WriteAzureKmsKey(object.activeKey, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureKmsEncryption reads a value of the 'azure_kms_encryption' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureKmsEncryption(source interface{}) (object *AzureKmsEncryption, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureKmsEncryption(iterator)
	err = iterator.Error
	return
}

// ReadAzureKmsEncryption reads a value of the 'azure_kms_encryption' type from the given iterator.
func ReadAzureKmsEncryption(iterator *jsoniter.Iterator) *AzureKmsEncryption {
	object := &AzureKmsEncryption{
		fieldSet_: make([]bool, 1),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "active_key":
			value := ReadAzureKmsKey(iterator)
			object.activeKey = value
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
