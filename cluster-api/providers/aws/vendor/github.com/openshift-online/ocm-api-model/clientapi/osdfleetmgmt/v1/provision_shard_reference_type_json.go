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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/osdfleetmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalProvisionShardReference writes a value of the 'provision_shard_reference' type to the given writer.
func MarshalProvisionShardReference(object *ProvisionShardReference, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteProvisionShardReference(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteProvisionShardReference writes a value of the 'provision_shard_reference' type to the given stream.
func WriteProvisionShardReference(object *ProvisionShardReference, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
	}
	stream.WriteObjectEnd()
}

// UnmarshalProvisionShardReference reads a value of the 'provision_shard_reference' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalProvisionShardReference(source interface{}) (object *ProvisionShardReference, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadProvisionShardReference(iterator)
	err = iterator.Error
	return
}

// ReadProvisionShardReference reads a value of the 'provision_shard_reference' type from the given iterator.
func ReadProvisionShardReference(iterator *jsoniter.Iterator) *ProvisionShardReference {
	object := &ProvisionShardReference{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "href":
			value := iterator.ReadString()
			object.href = value
			object.fieldSet_[0] = true
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
