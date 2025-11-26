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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalServerConfig writes a value of the 'server_config' type to the given writer.
func MarshalServerConfig(object *ServerConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteServerConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteServerConfig writes a value of the 'server_config' type to the given stream.
func WriteServerConfig(object *ServerConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ServerConfigLinkKind)
	} else {
		stream.WriteString(ServerConfigKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.awsShard != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_shard")
		WriteAWSShard(object.awsShard, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubeconfig")
		stream.WriteString(object.kubeconfig)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("server")
		stream.WriteString(object.server)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("topology")
		stream.WriteString(string(object.topology))
	}
	stream.WriteObjectEnd()
}

// UnmarshalServerConfig reads a value of the 'server_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalServerConfig(source interface{}) (object *ServerConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadServerConfig(iterator)
	err = iterator.Error
	return
}

// ReadServerConfig reads a value of the 'server_config' type from the given iterator.
func ReadServerConfig(iterator *jsoniter.Iterator) *ServerConfig {
	object := &ServerConfig{
		fieldSet_: make([]bool, 7),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ServerConfigLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "aws_shard":
			value := ReadAWSShard(iterator)
			object.awsShard = value
			object.fieldSet_[3] = true
		case "kubeconfig":
			value := iterator.ReadString()
			object.kubeconfig = value
			object.fieldSet_[4] = true
		case "server":
			value := iterator.ReadString()
			object.server = value
			object.fieldSet_[5] = true
		case "topology":
			text := iterator.ReadString()
			value := ProvisionShardTopology(text)
			object.topology = value
			object.fieldSet_[6] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
