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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalRoleBinding writes a value of the 'role_binding' type to the given writer.
func MarshalRoleBinding(object *RoleBinding, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeRoleBinding(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeRoleBinding writes a value of the 'role_binding' type to the given stream.
func writeRoleBinding(object *RoleBinding, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(RoleBindingLinkKind)
	} else {
		stream.WriteString(RoleBindingKind)
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
	present_ = object.bitmap_&8 != 0 && object.account != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account")
		writeAccount(object.account, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("config_managed")
		stream.WriteBool(object.configManaged)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_by")
		stream.WriteString(object.managedBy)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.organization != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization")
		writeOrganization(object.organization, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationID)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.role != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role")
		writeRole(object.role, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_id")
		stream.WriteString(object.roleID)
		count++
	}
	present_ = object.bitmap_&4096 != 0 && object.subscription != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription")
		writeSubscription(object.subscription, stream)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionID)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
		count++
	}
	present_ = object.bitmap_&32768 != 0
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
	object = readRoleBinding(iterator)
	err = iterator.Error
	return
}

// readRoleBinding reads a value of the 'role_binding' type from the given iterator.
func readRoleBinding(iterator *jsoniter.Iterator) *RoleBinding {
	object := &RoleBinding{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == RoleBindingLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "account":
			value := readAccount(iterator)
			object.account = value
			object.bitmap_ |= 8
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.bitmap_ |= 16
		case "config_managed":
			value := iterator.ReadBool()
			object.configManaged = value
			object.bitmap_ |= 32
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 64
		case "managed_by":
			value := iterator.ReadString()
			object.managedBy = value
			object.bitmap_ |= 128
		case "organization":
			value := readOrganization(iterator)
			object.organization = value
			object.bitmap_ |= 256
		case "organization_id":
			value := iterator.ReadString()
			object.organizationID = value
			object.bitmap_ |= 512
		case "role":
			value := readRole(iterator)
			object.role = value
			object.bitmap_ |= 1024
		case "role_id":
			value := iterator.ReadString()
			object.roleID = value
			object.bitmap_ |= 2048
		case "subscription":
			value := readSubscription(iterator)
			object.subscription = value
			object.bitmap_ |= 4096
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.bitmap_ |= 8192
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.bitmap_ |= 16384
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 32768
		default:
			iterator.ReadAny()
		}
	}
	return object
}
