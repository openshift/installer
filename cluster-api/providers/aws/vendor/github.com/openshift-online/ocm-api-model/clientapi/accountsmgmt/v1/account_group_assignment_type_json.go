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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAccountGroupAssignment writes a value of the 'account_group_assignment' type to the given writer.
func MarshalAccountGroupAssignment(object *AccountGroupAssignment, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccountGroupAssignment(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccountGroupAssignment writes a value of the 'account_group_assignment' type to the given stream.
func WriteAccountGroupAssignment(object *AccountGroupAssignment, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AccountGroupAssignmentLinkKind)
	} else {
		stream.WriteString(AccountGroupAssignmentKind)
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
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.accountGroup != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_group")
		WriteAccountGroup(object.accountGroup, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_group_id")
		stream.WriteString(object.accountGroupID)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_by")
		stream.WriteString(string(object.managedBy))
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccountGroupAssignment reads a value of the 'account_group_assignment' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccountGroupAssignment(source interface{}) (object *AccountGroupAssignment, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAccountGroupAssignment(iterator)
	err = iterator.Error
	return
}

// ReadAccountGroupAssignment reads a value of the 'account_group_assignment' type from the given iterator.
func ReadAccountGroupAssignment(iterator *jsoniter.Iterator) *AccountGroupAssignment {
	object := &AccountGroupAssignment{
		fieldSet_: make([]bool, 9),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AccountGroupAssignmentLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.fieldSet_[3] = true
		case "account_group":
			value := ReadAccountGroup(iterator)
			object.accountGroup = value
			object.fieldSet_[4] = true
		case "account_group_id":
			value := iterator.ReadString()
			object.accountGroupID = value
			object.fieldSet_[5] = true
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.fieldSet_[6] = true
		case "managed_by":
			text := iterator.ReadString()
			value := AccountGroupAssignmentManagedBy(text)
			object.managedBy = value
			object.fieldSet_[7] = true
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.fieldSet_[8] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
