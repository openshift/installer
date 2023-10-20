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

// MarshalSTSOperator writes a value of the 'STS_operator' type to the given writer.
func MarshalSTSOperator(object *STSOperator, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSTSOperator(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSTSOperator writes a value of the 'STS_operator' type to the given stream.
func writeSTSOperator(object *STSOperator, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_version")
		stream.WriteString(object.maxVersion)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("min_version")
		stream.WriteString(object.minVersion)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespace")
		stream.WriteString(object.namespace)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.serviceAccounts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_accounts")
		writeStringList(object.serviceAccounts, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSTSOperator reads a value of the 'STS_operator' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSTSOperator(source interface{}) (object *STSOperator, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSTSOperator(iterator)
	err = iterator.Error
	return
}

// readSTSOperator reads a value of the 'STS_operator' type from the given iterator.
func readSTSOperator(iterator *jsoniter.Iterator) *STSOperator {
	object := &STSOperator{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "max_version":
			value := iterator.ReadString()
			object.maxVersion = value
			object.bitmap_ |= 1
		case "min_version":
			value := iterator.ReadString()
			object.minVersion = value
			object.bitmap_ |= 2
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 4
		case "namespace":
			value := iterator.ReadString()
			object.namespace = value
			object.bitmap_ |= 8
		case "service_accounts":
			value := readStringList(iterator)
			object.serviceAccounts = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
