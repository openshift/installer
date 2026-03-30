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

// MarshalWifPool writes a value of the 'wif_pool' type to the given writer.
func MarshalWifPool(object *WifPool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifPool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifPool writes a value of the 'wif_pool' type to the given stream.
func WriteWifPool(object *WifPool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.identityProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("identity_provider")
		WriteWifIdentityProvider(object.identityProvider, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("pool_id")
		stream.WriteString(object.poolId)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("pool_name")
		stream.WriteString(object.poolName)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifPool reads a value of the 'wif_pool' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifPool(source interface{}) (object *WifPool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifPool(iterator)
	err = iterator.Error
	return
}

// ReadWifPool reads a value of the 'wif_pool' type from the given iterator.
func ReadWifPool(iterator *jsoniter.Iterator) *WifPool {
	object := &WifPool{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "identity_provider":
			value := ReadWifIdentityProvider(iterator)
			object.identityProvider = value
			object.fieldSet_[0] = true
		case "pool_id":
			value := iterator.ReadString()
			object.poolId = value
			object.fieldSet_[1] = true
		case "pool_name":
			value := iterator.ReadString()
			object.poolName = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
