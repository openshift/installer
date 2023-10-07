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

// MarshalHypershiftConfig writes a value of the 'hypershift_config' type to the given writer.
func MarshalHypershiftConfig(object *HypershiftConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeHypershiftConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeHypershiftConfig writes a value of the 'hypershift_config' type to the given stream.
func writeHypershiftConfig(object *HypershiftConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("management_cluster")
		stream.WriteString(object.managementCluster)
	}
	stream.WriteObjectEnd()
}

// UnmarshalHypershiftConfig reads a value of the 'hypershift_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalHypershiftConfig(source interface{}) (object *HypershiftConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readHypershiftConfig(iterator)
	err = iterator.Error
	return
}

// readHypershiftConfig reads a value of the 'hypershift_config' type from the given iterator.
func readHypershiftConfig(iterator *jsoniter.Iterator) *HypershiftConfig {
	object := &HypershiftConfig{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 1
		case "management_cluster":
			value := iterator.ReadString()
			object.managementCluster = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
