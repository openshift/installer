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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

import (
	"io"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalMetricsFederation writes a value of the 'metrics_federation' type to the given writer.
func MarshalMetricsFederation(object *MetricsFederation, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteMetricsFederation(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteMetricsFederation writes a value of the 'metrics_federation' type to the given stream.
func WriteMetricsFederation(object *MetricsFederation, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.matchLabels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("match_labels")
		if object.matchLabels != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.matchLabels))
			i := 0
			for key := range object.matchLabels {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.matchLabels[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.matchNames != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("match_names")
		WriteStringList(object.matchNames, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespace")
		stream.WriteString(object.namespace)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("port_name")
		stream.WriteString(object.portName)
	}
	stream.WriteObjectEnd()
}

// UnmarshalMetricsFederation reads a value of the 'metrics_federation' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalMetricsFederation(source interface{}) (object *MetricsFederation, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadMetricsFederation(iterator)
	err = iterator.Error
	return
}

// ReadMetricsFederation reads a value of the 'metrics_federation' type from the given iterator.
func ReadMetricsFederation(iterator *jsoniter.Iterator) *MetricsFederation {
	object := &MetricsFederation{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "match_labels":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.matchLabels = value
			object.fieldSet_[0] = true
		case "match_names":
			value := ReadStringList(iterator)
			object.matchNames = value
			object.fieldSet_[1] = true
		case "namespace":
			value := iterator.ReadString()
			object.namespace = value
			object.fieldSet_[2] = true
		case "port_name":
			value := iterator.ReadString()
			object.portName = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
