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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalClusterMetricsNodes writes a value of the 'cluster_metrics_nodes' type to the given writer.
func MarshalClusterMetricsNodes(object *ClusterMetricsNodes, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterMetricsNodes(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterMetricsNodes writes a value of the 'cluster_metrics_nodes' type to the given stream.
func writeClusterMetricsNodes(object *ClusterMetricsNodes, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute")
		stream.WriteFloat64(object.compute)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra")
		stream.WriteFloat64(object.infra)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master")
		stream.WriteFloat64(object.master)
		count++
	}
	present_ = object.bitmap_&8 != 0
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
	object = readClusterMetricsNodes(iterator)
	err = iterator.Error
	return
}

// readClusterMetricsNodes reads a value of the 'cluster_metrics_nodes' type from the given iterator.
func readClusterMetricsNodes(iterator *jsoniter.Iterator) *ClusterMetricsNodes {
	object := &ClusterMetricsNodes{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "compute":
			value := iterator.ReadFloat64()
			object.compute = value
			object.bitmap_ |= 1
		case "infra":
			value := iterator.ReadFloat64()
			object.infra = value
			object.bitmap_ |= 2
		case "master":
			value := iterator.ReadFloat64()
			object.master = value
			object.bitmap_ |= 4
		case "total":
			value := iterator.ReadFloat64()
			object.total = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
