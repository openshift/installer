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

// MarshalClusterOperatorsInfo writes a value of the 'cluster_operators_info' type to the given writer.
func MarshalClusterOperatorsInfo(object *ClusterOperatorsInfo, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterOperatorsInfo(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterOperatorsInfo writes a value of the 'cluster_operators_info' type to the given stream.
func writeClusterOperatorsInfo(object *ClusterOperatorsInfo, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.operators != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operators")
		writeClusterOperatorInfoList(object.operators, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterOperatorsInfo reads a value of the 'cluster_operators_info' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterOperatorsInfo(source interface{}) (object *ClusterOperatorsInfo, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readClusterOperatorsInfo(iterator)
	err = iterator.Error
	return
}

// readClusterOperatorsInfo reads a value of the 'cluster_operators_info' type from the given iterator.
func readClusterOperatorsInfo(iterator *jsoniter.Iterator) *ClusterOperatorsInfo {
	object := &ClusterOperatorsInfo{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "operators":
			value := readClusterOperatorInfoList(iterator)
			object.operators = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
