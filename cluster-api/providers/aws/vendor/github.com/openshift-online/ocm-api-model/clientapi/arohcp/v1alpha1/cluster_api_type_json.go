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

// MarshalClusterAPI writes a value of the 'cluster_API' type to the given writer.
func MarshalClusterAPI(object *ClusterAPI, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterAPI(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterAPI writes a value of the 'cluster_API' type to the given stream.
func WriteClusterAPI(object *ClusterAPI, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.cidrBlockAccess != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cidr_block_access")
		WriteCIDRBlockAccess(object.cidrBlockAccess, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("url")
		stream.WriteString(object.url)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
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
	object = ReadClusterAPI(iterator)
	err = iterator.Error
	return
}

// ReadClusterAPI reads a value of the 'cluster_API' type from the given iterator.
func ReadClusterAPI(iterator *jsoniter.Iterator) *ClusterAPI {
	object := &ClusterAPI{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cidr_block_access":
			value := ReadCIDRBlockAccess(iterator)
			object.cidrBlockAccess = value
			object.fieldSet_[0] = true
		case "url":
			value := iterator.ReadString()
			object.url = value
			object.fieldSet_[1] = true
		case "listening":
			text := iterator.ReadString()
			value := ListeningMethod(text)
			object.listening = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
