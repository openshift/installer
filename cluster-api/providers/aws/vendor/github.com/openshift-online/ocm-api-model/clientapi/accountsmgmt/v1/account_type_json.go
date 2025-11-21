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
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AccountLinkKind)
	} else {
		stream.WriteString(AccountKind)
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
		stream.WriteObjectField("ban_code")
		stream.WriteString(object.banCode)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ban_description")
		stream.WriteString(object.banDescription)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("banned")
		stream.WriteBool(object.banned)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.capabilities != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("capabilities")
		WriteCapabilityList(object.capabilities, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("email")
		stream.WriteString(object.email)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("first_name")
		stream.WriteString(object.firstName)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		WriteLabelList(object.labels, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_name")
		stream.WriteString(object.lastName)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.organization != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization")
		WriteOrganization(object.organization, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rhit_account_id")
		stream.WriteString(object.rhitAccountID)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rhit_web_user_id")
		stream.WriteString(object.rhitWebUserId)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_account")
		stream.WriteBool(object.serviceAccount)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17]
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
	object := &Account{
		fieldSet_: make([]bool, 18),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AccountLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "ban_code":
			value := iterator.ReadString()
			object.banCode = value
			object.fieldSet_[3] = true
		case "ban_description":
			value := iterator.ReadString()
			object.banDescription = value
			object.fieldSet_[4] = true
		case "banned":
			value := iterator.ReadBool()
			object.banned = value
			object.fieldSet_[5] = true
		case "capabilities":
			value := ReadCapabilityList(iterator)
			object.capabilities = value
			object.fieldSet_[6] = true
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.fieldSet_[7] = true
		case "email":
			value := iterator.ReadString()
			object.email = value
			object.fieldSet_[8] = true
		case "first_name":
			value := iterator.ReadString()
			object.firstName = value
			object.fieldSet_[9] = true
		case "labels":
			value := ReadLabelList(iterator)
			object.labels = value
			object.fieldSet_[10] = true
		case "last_name":
			value := iterator.ReadString()
			object.lastName = value
			object.fieldSet_[11] = true
		case "organization":
			value := ReadOrganization(iterator)
			object.organization = value
			object.fieldSet_[12] = true
		case "rhit_account_id":
			value := iterator.ReadString()
			object.rhitAccountID = value
			object.fieldSet_[13] = true
		case "rhit_web_user_id":
			value := iterator.ReadString()
			object.rhitWebUserId = value
			object.fieldSet_[14] = true
		case "service_account":
			value := iterator.ReadBool()
			object.serviceAccount = value
			object.fieldSet_[15] = true
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.fieldSet_[16] = true
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.fieldSet_[17] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
