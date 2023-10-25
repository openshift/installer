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

// MarshalAddOnSecretPropagation writes a value of the 'add_on_secret_propagation' type to the given writer.
func MarshalAddOnSecretPropagation(object *AddOnSecretPropagation, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddOnSecretPropagation(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddOnSecretPropagation writes a value of the 'add_on_secret_propagation' type to the given stream.
func writeAddOnSecretPropagation(object *AddOnSecretPropagation, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("destination_secret")
		stream.WriteString(object.destinationSecret)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("source_secret")
		stream.WriteString(object.sourceSecret)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnSecretPropagation reads a value of the 'add_on_secret_propagation' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnSecretPropagation(source interface{}) (object *AddOnSecretPropagation, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddOnSecretPropagation(iterator)
	err = iterator.Error
	return
}

// readAddOnSecretPropagation reads a value of the 'add_on_secret_propagation' type from the given iterator.
func readAddOnSecretPropagation(iterator *jsoniter.Iterator) *AddOnSecretPropagation {
	object := &AddOnSecretPropagation{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 1
		case "destination_secret":
			value := iterator.ReadString()
			object.destinationSecret = value
			object.bitmap_ |= 2
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 4
		case "source_secret":
			value := iterator.ReadString()
			object.sourceSecret = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
