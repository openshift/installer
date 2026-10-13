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

// MarshalAzureEtcdEncryption writes a value of the 'azure_etcd_encryption' type to the given writer.
func MarshalAzureEtcdEncryption(object *AzureEtcdEncryption, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureEtcdEncryption(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureEtcdEncryption writes a value of the 'azure_etcd_encryption' type to the given stream.
func WriteAzureEtcdEncryption(object *AzureEtcdEncryption, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.dataEncryption != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("data_encryption")
		WriteAzureEtcdDataEncryption(object.dataEncryption, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureEtcdEncryption reads a value of the 'azure_etcd_encryption' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureEtcdEncryption(source interface{}) (object *AzureEtcdEncryption, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureEtcdEncryption(iterator)
	err = iterator.Error
	return
}

// ReadAzureEtcdEncryption reads a value of the 'azure_etcd_encryption' type from the given iterator.
func ReadAzureEtcdEncryption(iterator *jsoniter.Iterator) *AzureEtcdEncryption {
	object := &AzureEtcdEncryption{
		fieldSet_: make([]bool, 1),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "data_encryption":
			value := ReadAzureEtcdDataEncryption(iterator)
			object.dataEncryption = value
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
