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

// MarshalAddOnConfig writes a value of the 'add_on_config' type to the given writer.
func MarshalAddOnConfig(object *AddOnConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddOnConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddOnConfig writes a value of the 'add_on_config' type to the given stream.
func writeAddOnConfig(object *AddOnConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AddOnConfigLinkKind)
	} else {
		stream.WriteString(AddOnConfigKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0 && object.addOnEnvironmentVariables != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("add_on_environment_variables")
		writeAddOnEnvironmentVariableList(object.addOnEnvironmentVariables, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.secretPropagations != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret_propagations")
		writeAddOnSecretPropagationList(object.secretPropagations, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnConfig reads a value of the 'add_on_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnConfig(source interface{}) (object *AddOnConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddOnConfig(iterator)
	err = iterator.Error
	return
}

// readAddOnConfig reads a value of the 'add_on_config' type from the given iterator.
func readAddOnConfig(iterator *jsoniter.Iterator) *AddOnConfig {
	object := &AddOnConfig{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddOnConfigLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "add_on_environment_variables":
			value := readAddOnEnvironmentVariableList(iterator)
			object.addOnEnvironmentVariables = value
			object.bitmap_ |= 8
		case "secret_propagations":
			value := readAddOnSecretPropagationList(iterator)
			object.secretPropagations = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
