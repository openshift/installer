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

// MarshalProvisionShardMaestroConfig writes a value of the 'provision_shard_maestro_config' type to the given writer.
func MarshalProvisionShardMaestroConfig(object *ProvisionShardMaestroConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteProvisionShardMaestroConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteProvisionShardMaestroConfig writes a value of the 'provision_shard_maestro_config' type to the given stream.
func WriteProvisionShardMaestroConfig(object *ProvisionShardMaestroConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("consumer_name")
		stream.WriteString(object.consumerName)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.grpcApiConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("grpc_api_config")
		WriteProvisionShardMaestroGrpcApiConfig(object.grpcApiConfig, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.restApiConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rest_api_config")
		WriteProvisionShardMaestroRestApiConfig(object.restApiConfig, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalProvisionShardMaestroConfig reads a value of the 'provision_shard_maestro_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalProvisionShardMaestroConfig(source interface{}) (object *ProvisionShardMaestroConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadProvisionShardMaestroConfig(iterator)
	err = iterator.Error
	return
}

// ReadProvisionShardMaestroConfig reads a value of the 'provision_shard_maestro_config' type from the given iterator.
func ReadProvisionShardMaestroConfig(iterator *jsoniter.Iterator) *ProvisionShardMaestroConfig {
	object := &ProvisionShardMaestroConfig{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "consumer_name":
			value := iterator.ReadString()
			object.consumerName = value
			object.fieldSet_[0] = true
		case "grpc_api_config":
			value := ReadProvisionShardMaestroGrpcApiConfig(iterator)
			object.grpcApiConfig = value
			object.fieldSet_[1] = true
		case "rest_api_config":
			value := ReadProvisionShardMaestroRestApiConfig(iterator)
			object.restApiConfig = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
