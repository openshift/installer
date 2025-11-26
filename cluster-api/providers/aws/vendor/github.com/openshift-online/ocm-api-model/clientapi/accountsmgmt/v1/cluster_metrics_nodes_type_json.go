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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalClusterMetricsNodes writes a value of the 'cluster_metrics_nodes' type to the given writer.
func MarshalClusterMetricsNodes(object *ClusterMetricsNodes, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterMetricsNodes(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterMetricsNodes writes a value of the 'cluster_metrics_nodes' type to the given stream.
func WriteClusterMetricsNodes(object *ClusterMetricsNodes, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute")
		stream.WriteFloat64(object.compute)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra")
		stream.WriteFloat64(object.infra)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master")
		stream.WriteFloat64(object.master)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("total")
		stream.WriteFloat64(object.total)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterMetricsNodes reads a value of the 'cluster_metrics_nodes' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterMetricsNodes(source interface{}) (object *ClusterMetricsNodes, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterMetricsNodes(iterator)
	err = iterator.Error
	return
}

// ReadClusterMetricsNodes reads a value of the 'cluster_metrics_nodes' type from the given iterator.
func ReadClusterMetricsNodes(iterator *jsoniter.Iterator) *ClusterMetricsNodes {
	object := &ClusterMetricsNodes{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "compute":
			value := iterator.ReadFloat64()
			object.compute = value
			object.fieldSet_[0] = true
		case "infra":
			value := iterator.ReadFloat64()
			object.infra = value
			object.fieldSet_[1] = true
		case "master":
			value := iterator.ReadFloat64()
			object.master = value
			object.fieldSet_[2] = true
		case "total":
			value := iterator.ReadFloat64()
			object.total = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
