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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	"io"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAddOnNamespace writes a value of the 'add_on_namespace' type to the given writer.
func MarshalAddOnNamespace(object *AddOnNamespace, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddOnNamespace(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddOnNamespace writes a value of the 'add_on_namespace' type to the given stream.
func WriteAddOnNamespace(object *AddOnNamespace, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AddOnNamespaceLinkKind)
	} else {
		stream.WriteString(AddOnNamespaceKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.annotations != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("annotations")
		if object.annotations != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.annotations))
			i := 0
			for key := range object.annotations {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.annotations[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		if object.labels != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.labels))
			i := 0
			for key := range object.labels {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.labels[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnNamespace reads a value of the 'add_on_namespace' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnNamespace(source interface{}) (object *AddOnNamespace, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddOnNamespace(iterator)
	err = iterator.Error
	return
}

// ReadAddOnNamespace reads a value of the 'add_on_namespace' type from the given iterator.
func ReadAddOnNamespace(iterator *jsoniter.Iterator) *AddOnNamespace {
	object := &AddOnNamespace{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddOnNamespaceLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "annotations":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.annotations = value
			object.fieldSet_[3] = true
		case "labels":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.labels = value
			object.fieldSet_[4] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
