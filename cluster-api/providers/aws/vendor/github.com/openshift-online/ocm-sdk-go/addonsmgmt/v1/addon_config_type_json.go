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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddonConfig writes a value of the 'addon_config' type to the given writer.
func MarshalAddonConfig(object *AddonConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddonConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddonConfig writes a value of the 'addon_config' type to the given stream.
func writeAddonConfig(object *AddonConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.addOnEnvironmentVariables != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("add_on_environment_variables")
		writeAddonEnvironmentVariableList(object.addOnEnvironmentVariables, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.addOnSecretPropagations != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("add_on_secret_propagations")
		writeAddonSecretPropagationList(object.addOnSecretPropagations, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonConfig reads a value of the 'addon_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonConfig(source interface{}) (object *AddonConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddonConfig(iterator)
	err = iterator.Error
	return
}

// readAddonConfig reads a value of the 'addon_config' type from the given iterator.
func readAddonConfig(iterator *jsoniter.Iterator) *AddonConfig {
	object := &AddonConfig{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "add_on_environment_variables":
			value := readAddonEnvironmentVariableList(iterator)
			object.addOnEnvironmentVariables = value
			object.bitmap_ |= 1
		case "add_on_secret_propagations":
			value := readAddonSecretPropagationList(iterator)
			object.addOnSecretPropagations = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
