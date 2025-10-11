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

// MarshalSdnToOvnClusterMigration writes a value of the 'sdn_to_ovn_cluster_migration' type to the given writer.
func MarshalSdnToOvnClusterMigration(object *SdnToOvnClusterMigration, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSdnToOvnClusterMigration(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSdnToOvnClusterMigration writes a value of the 'sdn_to_ovn_cluster_migration' type to the given stream.
func WriteSdnToOvnClusterMigration(object *SdnToOvnClusterMigration, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("join_ipv4")
		stream.WriteString(object.joinIpv4)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("masquerade_ipv4")
		stream.WriteString(object.masqueradeIpv4)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("transit_ipv4")
		stream.WriteString(object.transitIpv4)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSdnToOvnClusterMigration reads a value of the 'sdn_to_ovn_cluster_migration' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSdnToOvnClusterMigration(source interface{}) (object *SdnToOvnClusterMigration, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadSdnToOvnClusterMigration(iterator)
	err = iterator.Error
	return
}

// ReadSdnToOvnClusterMigration reads a value of the 'sdn_to_ovn_cluster_migration' type from the given iterator.
func ReadSdnToOvnClusterMigration(iterator *jsoniter.Iterator) *SdnToOvnClusterMigration {
	object := &SdnToOvnClusterMigration{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "join_ipv4":
			value := iterator.ReadString()
			object.joinIpv4 = value
			object.bitmap_ |= 1
		case "masquerade_ipv4":
			value := iterator.ReadString()
			object.masqueradeIpv4 = value
			object.bitmap_ |= 2
		case "transit_ipv4":
			value := iterator.ReadString()
			object.transitIpv4 = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
