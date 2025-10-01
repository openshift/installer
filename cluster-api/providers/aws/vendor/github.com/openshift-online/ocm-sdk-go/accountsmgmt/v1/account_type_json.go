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

// MarshalAccount writes a value of the 'account' type to the given writer.
func MarshalAccount(object *Account, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccount(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccount writes a value of the 'account' type to the given stream.
func WriteAccount(object *Account, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AccountLinkKind)
	} else {
		stream.WriteString(AccountKind)
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
		stream.WriteObjectField("ban_code")
		stream.WriteString(object.banCode)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ban_description")
		stream.WriteString(object.banDescription)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("banned")
		stream.WriteBool(object.banned)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.capabilities != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("capabilities")
		WriteCapabilityList(object.capabilities, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("email")
		stream.WriteString(object.email)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("first_name")
		stream.WriteString(object.firstName)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		WriteLabelList(object.labels, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_name")
		stream.WriteString(object.lastName)
		count++
	}
	present_ = object.bitmap_&4096 != 0 && object.organization != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization")
		WriteOrganization(object.organization, stream)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rhit_account_id")
		stream.WriteString(object.rhitAccountID)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rhit_web_user_id")
		stream.WriteString(object.rhitWebUserId)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_account")
		stream.WriteBool(object.serviceAccount)
		count++
	}
	present_ = object.bitmap_&65536 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&131072 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("username")
		stream.WriteString(object.username)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccount reads a value of the 'account' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccount(source interface{}) (object *Account, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAccount(iterator)
	err = iterator.Error
	return
}

// ReadAccount reads a value of the 'account' type from the given iterator.
func ReadAccount(iterator *jsoniter.Iterator) *Account {
	object := &Account{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AccountLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "ban_code":
			value := iterator.ReadString()
			object.banCode = value
			object.bitmap_ |= 8
		case "ban_description":
			value := iterator.ReadString()
			object.banDescription = value
			object.bitmap_ |= 16
		case "banned":
			value := iterator.ReadBool()
			object.banned = value
			object.bitmap_ |= 32
		case "capabilities":
			value := ReadCapabilityList(iterator)
			object.capabilities = value
			object.bitmap_ |= 64
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 128
		case "email":
			value := iterator.ReadString()
			object.email = value
			object.bitmap_ |= 256
		case "first_name":
			value := iterator.ReadString()
			object.firstName = value
			object.bitmap_ |= 512
		case "labels":
			value := ReadLabelList(iterator)
			object.labels = value
			object.bitmap_ |= 1024
		case "last_name":
			value := iterator.ReadString()
			object.lastName = value
			object.bitmap_ |= 2048
		case "organization":
			value := ReadOrganization(iterator)
			object.organization = value
			object.bitmap_ |= 4096
		case "rhit_account_id":
			value := iterator.ReadString()
			object.rhitAccountID = value
			object.bitmap_ |= 8192
		case "rhit_web_user_id":
			value := iterator.ReadString()
			object.rhitWebUserId = value
			object.bitmap_ |= 16384
		case "service_account":
			value := iterator.ReadBool()
			object.serviceAccount = value
			object.bitmap_ |= 32768
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 65536
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.bitmap_ |= 131072
		default:
			iterator.ReadAny()
		}
	}
	return object
}
