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

// MarshalRoleBinding writes a value of the 'role_binding' type to the given writer.
func MarshalRoleBinding(object *RoleBinding, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteRoleBinding(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteRoleBinding writes a value of the 'role_binding' type to the given stream.
func WriteRoleBinding(object *RoleBinding, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(RoleBindingLinkKind)
	} else {
		stream.WriteString(RoleBindingKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.account != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account")
		WriteAccount(object.account, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("config_managed")
		stream.WriteBool(object.configManaged)
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
		stream.WriteString(object.managedBy)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.organization != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization")
		WriteOrganization(object.organization, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationID)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.role != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role")
		WriteRole(object.role, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_id")
		stream.WriteString(object.roleID)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.subscription != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription")
		WriteSubscription(object.subscription, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionID)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalRoleBinding reads a value of the 'role_binding' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRoleBinding(source interface{}) (object *RoleBinding, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadRoleBinding(iterator)
	err = iterator.Error
	return
}

// ReadRoleBinding reads a value of the 'role_binding' type from the given iterator.
func ReadRoleBinding(iterator *jsoniter.Iterator) *RoleBinding {
	object := &RoleBinding{
		fieldSet_: make([]bool, 16),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == RoleBindingLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "account":
			value := ReadAccount(iterator)
			object.account = value
			object.fieldSet_[3] = true
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.fieldSet_[4] = true
		case "config_managed":
			value := iterator.ReadBool()
			object.configManaged = value
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
			value := iterator.ReadString()
			object.managedBy = value
			object.fieldSet_[7] = true
		case "organization":
			value := ReadOrganization(iterator)
			object.organization = value
			object.fieldSet_[8] = true
		case "organization_id":
			value := iterator.ReadString()
			object.organizationID = value
			object.fieldSet_[9] = true
		case "role":
			value := ReadRole(iterator)
			object.role = value
			object.fieldSet_[10] = true
		case "role_id":
			value := iterator.ReadString()
			object.roleID = value
			object.fieldSet_[11] = true
		case "subscription":
			value := ReadSubscription(iterator)
			object.subscription = value
			object.fieldSet_[12] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.fieldSet_[13] = true
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.fieldSet_[14] = true
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.fieldSet_[15] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
