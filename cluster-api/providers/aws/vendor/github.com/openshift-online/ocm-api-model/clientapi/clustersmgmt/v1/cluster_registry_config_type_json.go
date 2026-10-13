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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalClusterRegistryConfig writes a value of the 'cluster_registry_config' type to the given writer.
func MarshalClusterRegistryConfig(object *ClusterRegistryConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterRegistryConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterRegistryConfig writes a value of the 'cluster_registry_config' type to the given stream.
func WriteClusterRegistryConfig(object *ClusterRegistryConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.additionalTrustedCa != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_trusted_ca")
		if object.additionalTrustedCa != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.additionalTrustedCa))
			i := 0
			for key := range object.additionalTrustedCa {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.additionalTrustedCa[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.allowedRegistriesForImport != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed_registries_for_import")
		WriteRegistryLocationList(object.allowedRegistriesForImport, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.platformAllowlist != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("platform_allowlist")
		WriteRegistryAllowlist(object.platformAllowlist, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.registrySources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("registry_sources")
		WriteRegistrySources(object.registrySources, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterRegistryConfig reads a value of the 'cluster_registry_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterRegistryConfig(source interface{}) (object *ClusterRegistryConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterRegistryConfig(iterator)
	err = iterator.Error
	return
}

// ReadClusterRegistryConfig reads a value of the 'cluster_registry_config' type from the given iterator.
func ReadClusterRegistryConfig(iterator *jsoniter.Iterator) *ClusterRegistryConfig {
	object := &ClusterRegistryConfig{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "additional_trusted_ca":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.additionalTrustedCa = value
			object.fieldSet_[0] = true
		case "allowed_registries_for_import":
			value := ReadRegistryLocationList(iterator)
			object.allowedRegistriesForImport = value
			object.fieldSet_[1] = true
		case "platform_allowlist":
			value := ReadRegistryAllowlist(iterator)
			object.platformAllowlist = value
			object.fieldSet_[2] = true
		case "registry_sources":
			value := ReadRegistrySources(iterator)
			object.registrySources = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
