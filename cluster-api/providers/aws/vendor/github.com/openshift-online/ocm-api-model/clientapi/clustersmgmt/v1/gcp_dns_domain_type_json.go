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

// MarshalGcpDnsDomain writes a value of the 'gcp_dns_domain' type to the given writer.
func MarshalGcpDnsDomain(object *GcpDnsDomain, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGcpDnsDomain(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGcpDnsDomain writes a value of the 'gcp_dns_domain' type to the given stream.
func WriteGcpDnsDomain(object *GcpDnsDomain, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("domain_prefix")
		stream.WriteString(object.domainPrefix)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("network_id")
		stream.WriteString(object.networkId)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("project_id")
		stream.WriteString(object.projectId)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGcpDnsDomain reads a value of the 'gcp_dns_domain' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGcpDnsDomain(source interface{}) (object *GcpDnsDomain, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGcpDnsDomain(iterator)
	err = iterator.Error
	return
}

// ReadGcpDnsDomain reads a value of the 'gcp_dns_domain' type from the given iterator.
func ReadGcpDnsDomain(iterator *jsoniter.Iterator) *GcpDnsDomain {
	object := &GcpDnsDomain{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "domain_prefix":
			value := iterator.ReadString()
			object.domainPrefix = value
			object.fieldSet_[0] = true
		case "network_id":
			value := iterator.ReadString()
			object.networkId = value
			object.fieldSet_[1] = true
		case "project_id":
			value := iterator.ReadString()
			object.projectId = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
