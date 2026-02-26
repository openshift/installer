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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalDNSDomain writes a value of the 'DNS_domain' type to the given writer.
func MarshalDNSDomain(object *DNSDomain, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteDNSDomain(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteDNSDomain writes a value of the 'DNS_domain' type to the given stream.
func WriteDNSDomain(object *DNSDomain, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(DNSDomainLinkKind)
	} else {
		stream.WriteString(DNSDomainKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.cluster != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster")
		WriteClusterLink(object.cluster, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_arch")
		stream.WriteString(string(object.clusterArch))
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.organization != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization")
		WriteOrganizationLink(object.organization, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reserved_at_timestamp")
		stream.WriteString((object.reservedAtTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("user_defined")
		stream.WriteBool(object.userDefined)
	}
	stream.WriteObjectEnd()
}

// UnmarshalDNSDomain reads a value of the 'DNS_domain' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalDNSDomain(source interface{}) (object *DNSDomain, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadDNSDomain(iterator)
	err = iterator.Error
	return
}

// ReadDNSDomain reads a value of the 'DNS_domain' type from the given iterator.
func ReadDNSDomain(iterator *jsoniter.Iterator) *DNSDomain {
	object := &DNSDomain{
		fieldSet_: make([]bool, 8),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == DNSDomainLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "cluster":
			value := ReadClusterLink(iterator)
			object.cluster = value
			object.fieldSet_[3] = true
		case "cluster_arch":
			text := iterator.ReadString()
			value := ClusterArchitecture(text)
			object.clusterArch = value
			object.fieldSet_[4] = true
		case "organization":
			value := ReadOrganizationLink(iterator)
			object.organization = value
			object.fieldSet_[5] = true
		case "reserved_at_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.reservedAtTimestamp = value
			object.fieldSet_[6] = true
		case "user_defined":
			value := iterator.ReadBool()
			object.userDefined = value
			object.fieldSet_[7] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
