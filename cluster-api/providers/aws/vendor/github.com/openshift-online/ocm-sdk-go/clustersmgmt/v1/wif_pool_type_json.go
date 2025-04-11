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
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalWifPool writes a value of the 'wif_pool' type to the given writer.
func MarshalWifPool(object *WifPool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeWifPool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeWifPool writes a value of the 'wif_pool' type to the given stream.
func writeWifPool(object *WifPool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.identityProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("identity_provider")
		writeWifIdentityProvider(object.identityProvider, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("pool_id")
		stream.WriteString(object.poolId)
		count++
	}
	present_ = object.bitmap_&4 != 0
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
	object = readWifPool(iterator)
	err = iterator.Error
	return
}

// readWifPool reads a value of the 'wif_pool' type from the given iterator.
func readWifPool(iterator *jsoniter.Iterator) *WifPool {
	object := &WifPool{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "identity_provider":
			value := readWifIdentityProvider(iterator)
			object.identityProvider = value
			object.bitmap_ |= 1
		case "pool_id":
			value := iterator.ReadString()
			object.poolId = value
			object.bitmap_ |= 2
		case "pool_name":
			value := iterator.ReadString()
			object.poolName = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
