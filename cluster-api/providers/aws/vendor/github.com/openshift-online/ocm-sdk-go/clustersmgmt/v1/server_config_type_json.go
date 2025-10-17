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

// MarshalServerConfig writes a value of the 'server_config' type to the given writer.
func MarshalServerConfig(object *ServerConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeServerConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeServerConfig writes a value of the 'server_config' type to the given stream.
func writeServerConfig(object *ServerConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ServerConfigLinkKind)
	} else {
		stream.WriteString(ServerConfigKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0 && object.awsShard != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_shard")
		writeAWSShard(object.awsShard, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubeconfig")
		stream.WriteString(object.kubeconfig)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("server")
		stream.WriteString(object.server)
		count++
	}
	present_ = object.bitmap_&64 != 0
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
	object = readServerConfig(iterator)
	err = iterator.Error
	return
}

// readServerConfig reads a value of the 'server_config' type from the given iterator.
func readServerConfig(iterator *jsoniter.Iterator) *ServerConfig {
	object := &ServerConfig{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ServerConfigLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "aws_shard":
			value := readAWSShard(iterator)
			object.awsShard = value
			object.bitmap_ |= 8
		case "kubeconfig":
			value := iterator.ReadString()
			object.kubeconfig = value
			object.bitmap_ |= 16
		case "server":
			value := iterator.ReadString()
			object.server = value
			object.bitmap_ |= 32
		case "topology":
			text := iterator.ReadString()
			value := ProvisionShardTopology(text)
			object.topology = value
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
