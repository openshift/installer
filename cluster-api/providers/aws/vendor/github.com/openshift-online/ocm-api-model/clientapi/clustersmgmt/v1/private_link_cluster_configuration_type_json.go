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

// MarshalPrivateLinkClusterConfiguration writes a value of the 'private_link_cluster_configuration' type to the given writer.
func MarshalPrivateLinkClusterConfiguration(object *PrivateLinkClusterConfiguration, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WritePrivateLinkClusterConfiguration(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WritePrivateLinkClusterConfiguration writes a value of the 'private_link_cluster_configuration' type to the given stream.
func WritePrivateLinkClusterConfiguration(object *PrivateLinkClusterConfiguration, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.principals != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("principals")
		WritePrivateLinkPrincipalList(object.principals, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalPrivateLinkClusterConfiguration reads a value of the 'private_link_cluster_configuration' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalPrivateLinkClusterConfiguration(source interface{}) (object *PrivateLinkClusterConfiguration, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadPrivateLinkClusterConfiguration(iterator)
	err = iterator.Error
	return
}

// ReadPrivateLinkClusterConfiguration reads a value of the 'private_link_cluster_configuration' type from the given iterator.
func ReadPrivateLinkClusterConfiguration(iterator *jsoniter.Iterator) *PrivateLinkClusterConfiguration {
	object := &PrivateLinkClusterConfiguration{
		fieldSet_: make([]bool, 1),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "principals":
			value := ReadPrivateLinkPrincipalList(iterator)
			object.principals = value
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
