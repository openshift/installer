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

// MarshalUpgradePolicyState writes a value of the 'upgrade_policy_state' type to the given writer.
func MarshalUpgradePolicyState(object *UpgradePolicyState, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteUpgradePolicyState(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteUpgradePolicyState writes a value of the 'upgrade_policy_state' type to the given stream.
func WriteUpgradePolicyState(object *UpgradePolicyState, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(UpgradePolicyStateLinkKind)
	} else {
		stream.WriteString(UpgradePolicyStateKind)
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

// UnmarshalUpgradePolicyState reads a value of the 'upgrade_policy_state' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalUpgradePolicyState(source interface{}) (object *UpgradePolicyState, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadUpgradePolicyState(iterator)
	err = iterator.Error
	return
}

// ReadUpgradePolicyState reads a value of the 'upgrade_policy_state' type from the given iterator.
func ReadUpgradePolicyState(iterator *jsoniter.Iterator) *UpgradePolicyState {
	object := &UpgradePolicyState{
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
			if value == UpgradePolicyStateLinkKind {
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
