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

// MarshalClusterStatus writes a value of the 'cluster_status' type to the given writer.
func MarshalClusterStatus(object *ClusterStatus, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterStatus(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterStatus writes a value of the 'cluster_status' type to the given stream.
func WriteClusterStatus(object *ClusterStatus, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ClusterStatusLinkKind)
	} else {
		stream.WriteString(ClusterStatusKind)
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
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dns_ready")
		stream.WriteBool(object.dnsReady)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("oidc_ready")
		stream.WriteBool(object.oidcReady)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("configuration_mode")
		stream.WriteString(string(object.configurationMode))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("current_compute")
		stream.WriteInt(object.currentCompute)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("limited_support_reason_count")
		stream.WriteInt(object.limitedSupportReasonCount)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_error_code")
		stream.WriteString(object.provisionErrorCode)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_error_message")
		stream.WriteString(object.provisionErrorMessage)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterStatus reads a value of the 'cluster_status' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterStatus(source interface{}) (object *ClusterStatus, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterStatus(iterator)
	err = iterator.Error
	return
}

// ReadClusterStatus reads a value of the 'cluster_status' type from the given iterator.
func ReadClusterStatus(iterator *jsoniter.Iterator) *ClusterStatus {
	object := &ClusterStatus{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ClusterStatusLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "dns_ready":
			value := iterator.ReadBool()
			object.dnsReady = value
			object.bitmap_ |= 8
		case "oidc_ready":
			value := iterator.ReadBool()
			object.oidcReady = value
			object.bitmap_ |= 16
		case "configuration_mode":
			text := iterator.ReadString()
			value := ClusterConfigurationMode(text)
			object.configurationMode = value
			object.bitmap_ |= 32
		case "current_compute":
			value := iterator.ReadInt()
			object.currentCompute = value
			object.bitmap_ |= 64
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 128
		case "limited_support_reason_count":
			value := iterator.ReadInt()
			object.limitedSupportReasonCount = value
			object.bitmap_ |= 256
		case "provision_error_code":
			value := iterator.ReadString()
			object.provisionErrorCode = value
			object.bitmap_ |= 512
		case "provision_error_message":
			value := iterator.ReadString()
			object.provisionErrorMessage = value
			object.bitmap_ |= 1024
		case "state":
			text := iterator.ReadString()
			value := ClusterState(text)
			object.state = value
			object.bitmap_ |= 2048
		default:
			iterator.ReadAny()
		}
	}
	return object
}
