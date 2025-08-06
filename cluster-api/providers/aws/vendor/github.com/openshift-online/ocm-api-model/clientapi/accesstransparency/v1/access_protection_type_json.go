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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accesstransparency/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAccessProtection writes a value of the 'access_protection' type to the given writer.
func MarshalAccessProtection(object *AccessProtection, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccessProtection(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccessProtection writes a value of the 'access_protection' type to the given stream.
func WriteAccessProtection(object *AccessProtection, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccessProtection reads a value of the 'access_protection' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccessProtection(source interface{}) (object *AccessProtection, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAccessProtection(iterator)
	err = iterator.Error
	return
}

// ReadAccessProtection reads a value of the 'access_protection' type from the given iterator.
func ReadAccessProtection(iterator *jsoniter.Iterator) *AccessProtection {
	object := &AccessProtection{
		fieldSet_: make([]bool, 1),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
