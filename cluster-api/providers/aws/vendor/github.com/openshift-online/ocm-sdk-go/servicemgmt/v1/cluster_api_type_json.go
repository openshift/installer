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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalClusterAPI writes a value of the 'cluster_API' type to the given writer.
func MarshalClusterAPI(object *ClusterAPI, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterAPI(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterAPI writes a value of the 'cluster_API' type to the given stream.
func writeClusterAPI(object *ClusterAPI, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("listening")
		stream.WriteString(string(object.listening))
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterAPI reads a value of the 'cluster_API' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterAPI(source interface{}) (object *ClusterAPI, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readClusterAPI(iterator)
	err = iterator.Error
	return
}

// readClusterAPI reads a value of the 'cluster_API' type from the given iterator.
func readClusterAPI(iterator *jsoniter.Iterator) *ClusterAPI {
	object := &ClusterAPI{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "listening":
			text := iterator.ReadString()
			value := ListeningMethod(text)
			object.listening = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
