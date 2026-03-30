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

// MarshalBreakGlassCredential writes a value of the 'break_glass_credential' type to the given writer.
func MarshalBreakGlassCredential(object *BreakGlassCredential, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteBreakGlassCredential(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteBreakGlassCredential writes a value of the 'break_glass_credential' type to the given stream.
func WriteBreakGlassCredential(object *BreakGlassCredential, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(BreakGlassCredentialLinkKind)
	} else {
		stream.WriteString(BreakGlassCredentialKind)
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
		stream.WriteObjectField("expiration_timestamp")
		stream.WriteString((object.expirationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubeconfig")
		stream.WriteString(object.kubeconfig)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("revocation_timestamp")
		stream.WriteString((object.revocationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(string(object.status))
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("username")
		stream.WriteString(object.username)
	}
	stream.WriteObjectEnd()
}

// UnmarshalBreakGlassCredential reads a value of the 'break_glass_credential' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalBreakGlassCredential(source interface{}) (object *BreakGlassCredential, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadBreakGlassCredential(iterator)
	err = iterator.Error
	return
}

// ReadBreakGlassCredential reads a value of the 'break_glass_credential' type from the given iterator.
func ReadBreakGlassCredential(iterator *jsoniter.Iterator) *BreakGlassCredential {
	object := &BreakGlassCredential{
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
			if value == BreakGlassCredentialLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "expiration_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.expirationTimestamp = value
			object.fieldSet_[3] = true
		case "kubeconfig":
			value := iterator.ReadString()
			object.kubeconfig = value
			object.fieldSet_[4] = true
		case "revocation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.revocationTimestamp = value
			object.fieldSet_[5] = true
		case "status":
			text := iterator.ReadString()
			value := BreakGlassCredentialStatus(text)
			object.status = value
			object.fieldSet_[6] = true
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.fieldSet_[7] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
