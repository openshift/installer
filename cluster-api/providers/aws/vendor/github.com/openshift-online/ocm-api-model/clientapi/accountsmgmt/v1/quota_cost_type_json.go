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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalQuotaCost writes a value of the 'quota_cost' type to the given writer.
func MarshalQuotaCost(object *QuotaCost, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteQuotaCost(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteQuotaCost writes a value of the 'quota_cost' type to the given stream.
func WriteQuotaCost(object *QuotaCost, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed")
		stream.WriteInt(object.allowed)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.cloudAccounts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_accounts")
		WriteCloudAccountList(object.cloudAccounts, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("consumed")
		stream.WriteInt(object.consumed)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("quota_id")
		stream.WriteString(object.quotaID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.relatedResources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("related_resources")
		WriteRelatedResourceList(object.relatedResources, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		stream.WriteString(object.version)
	}
	stream.WriteObjectEnd()
}

// UnmarshalQuotaCost reads a value of the 'quota_cost' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalQuotaCost(source interface{}) (object *QuotaCost, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadQuotaCost(iterator)
	err = iterator.Error
	return
}

// ReadQuotaCost reads a value of the 'quota_cost' type from the given iterator.
func ReadQuotaCost(iterator *jsoniter.Iterator) *QuotaCost {
	object := &QuotaCost{
		fieldSet_: make([]bool, 7),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "allowed":
			value := iterator.ReadInt()
			object.allowed = value
			object.fieldSet_[0] = true
		case "cloud_accounts":
			value := ReadCloudAccountList(iterator)
			object.cloudAccounts = value
			object.fieldSet_[1] = true
		case "consumed":
			value := iterator.ReadInt()
			object.consumed = value
			object.fieldSet_[2] = true
		case "organization_id":
			value := iterator.ReadString()
			object.organizationID = value
			object.fieldSet_[3] = true
		case "quota_id":
			value := iterator.ReadString()
			object.quotaID = value
			object.fieldSet_[4] = true
		case "related_resources":
			value := ReadRelatedResourceList(iterator)
			object.relatedResources = value
			object.fieldSet_[5] = true
		case "version":
			value := iterator.ReadString()
			object.version = value
			object.fieldSet_[6] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
