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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddonNamespace writes a value of the 'addon_namespace' type to the given writer.
func MarshalAddonNamespace(object *AddonNamespace, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddonNamespace(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddonNamespace writes a value of the 'addon_namespace' type to the given stream.
func writeAddonNamespace(object *AddonNamespace, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.annotations != nil
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
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.labels != nil
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
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonNamespace reads a value of the 'addon_namespace' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonNamespace(source interface{}) (object *AddonNamespace, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddonNamespace(iterator)
	err = iterator.Error
	return
}

// readAddonNamespace reads a value of the 'addon_namespace' type from the given iterator.
func readAddonNamespace(iterator *jsoniter.Iterator) *AddonNamespace {
	object := &AddonNamespace{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
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
			object.bitmap_ |= 1
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 2
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
			object.bitmap_ |= 4
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
