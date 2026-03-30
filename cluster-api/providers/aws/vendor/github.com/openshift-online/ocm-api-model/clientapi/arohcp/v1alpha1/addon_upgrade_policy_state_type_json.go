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

// MarshalAddonUpgradePolicyState writes a value of the 'addon_upgrade_policy_state' type to the given writer.
func MarshalAddonUpgradePolicyState(object *AddonUpgradePolicyState, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonUpgradePolicyState(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonUpgradePolicyState writes a value of the 'addon_upgrade_policy_state' type to the given stream.
func WriteAddonUpgradePolicyState(object *AddonUpgradePolicyState, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AddonUpgradePolicyStateLinkKind)
	} else {
		stream.WriteString(AddonUpgradePolicyStateKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteString(string(object.value))
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonUpgradePolicyState reads a value of the 'addon_upgrade_policy_state' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonUpgradePolicyState(source interface{}) (object *AddonUpgradePolicyState, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddonUpgradePolicyState(iterator)
	err = iterator.Error
	return
}

// ReadAddonUpgradePolicyState reads a value of the 'addon_upgrade_policy_state' type from the given iterator.
func ReadAddonUpgradePolicyState(iterator *jsoniter.Iterator) *AddonUpgradePolicyState {
	object := &AddonUpgradePolicyState{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddonUpgradePolicyStateLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[3] = true
		case "value":
			text := iterator.ReadString()
			value := UpgradePolicyStateValue(text)
			object.value = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
