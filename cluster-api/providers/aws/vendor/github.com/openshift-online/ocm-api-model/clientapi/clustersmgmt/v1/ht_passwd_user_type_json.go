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

// MarshalHTPasswdUser writes a value of the 'HT_passwd_user' type to the given writer.
func MarshalHTPasswdUser(object *HTPasswdUser, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteHTPasswdUser(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteHTPasswdUser writes a value of the 'HT_passwd_user' type to the given stream.
func WriteHTPasswdUser(object *HTPasswdUser, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hashed_password")
		stream.WriteString(object.hashedPassword)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("password")
		stream.WriteString(object.password)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("username")
		stream.WriteString(object.username)
	}
	stream.WriteObjectEnd()
}

// UnmarshalHTPasswdUser reads a value of the 'HT_passwd_user' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalHTPasswdUser(source interface{}) (object *HTPasswdUser, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadHTPasswdUser(iterator)
	err = iterator.Error
	return
}

// ReadHTPasswdUser reads a value of the 'HT_passwd_user' type from the given iterator.
func ReadHTPasswdUser(iterator *jsoniter.Iterator) *HTPasswdUser {
	object := &HTPasswdUser{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.fieldSet_[0] = true
		case "hashed_password":
			value := iterator.ReadString()
			object.hashedPassword = value
			object.fieldSet_[1] = true
		case "password":
			value := iterator.ReadString()
			object.password = value
			object.fieldSet_[2] = true
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
