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
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ClusterStatusLinkKind)
	} else {
		stream.WriteString(ClusterStatusKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dns_ready")
		stream.WriteBool(object.dnsReady)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("oidc_ready")
		stream.WriteBool(object.oidcReady)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("configuration_mode")
		stream.WriteString(string(object.configurationMode))
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("current_compute")
		stream.WriteInt(object.currentCompute)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("limited_support_reason_count")
		stream.WriteInt(object.limitedSupportReasonCount)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_error_code")
		stream.WriteString(object.provisionErrorCode)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_error_message")
		stream.WriteString(object.provisionErrorMessage)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
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
	object := &ClusterStatus{
		fieldSet_: make([]bool, 12),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ClusterStatusLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "dns_ready":
			value := iterator.ReadBool()
			object.dnsReady = value
			object.fieldSet_[3] = true
		case "oidc_ready":
			value := iterator.ReadBool()
			object.oidcReady = value
			object.fieldSet_[4] = true
		case "configuration_mode":
			text := iterator.ReadString()
			value := ClusterConfigurationMode(text)
			object.configurationMode = value
			object.fieldSet_[5] = true
		case "current_compute":
			value := iterator.ReadInt()
			object.currentCompute = value
			object.fieldSet_[6] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[7] = true
		case "limited_support_reason_count":
			value := iterator.ReadInt()
			object.limitedSupportReasonCount = value
			object.fieldSet_[8] = true
		case "provision_error_code":
			value := iterator.ReadString()
			object.provisionErrorCode = value
			object.fieldSet_[9] = true
		case "provision_error_message":
			value := iterator.ReadString()
			object.provisionErrorMessage = value
			object.fieldSet_[10] = true
		case "state":
			text := iterator.ReadString()
			value := ClusterState(text)
			object.state = value
			object.fieldSet_[11] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
