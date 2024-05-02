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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalBreakGlassCredential writes a value of the 'break_glass_credential' type to the given writer.
func MarshalBreakGlassCredential(object *BreakGlassCredential, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeBreakGlassCredential(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeBreakGlassCredential writes a value of the 'break_glass_credential' type to the given stream.
func writeBreakGlassCredential(object *BreakGlassCredential, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(BreakGlassCredentialLinkKind)
	} else {
		stream.WriteString(BreakGlassCredentialKind)
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
		stream.WriteObjectField("expiration_timestamp")
		stream.WriteString((object.expirationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubeconfig")
		stream.WriteString(object.kubeconfig)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("revocation_timestamp")
		stream.WriteString((object.revocationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(string(object.status))
		count++
	}
	present_ = object.bitmap_&128 != 0
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
	object = readBreakGlassCredential(iterator)
	err = iterator.Error
	return
}

// readBreakGlassCredential reads a value of the 'break_glass_credential' type from the given iterator.
func readBreakGlassCredential(iterator *jsoniter.Iterator) *BreakGlassCredential {
	object := &BreakGlassCredential{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == BreakGlassCredentialLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "expiration_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.expirationTimestamp = value
			object.bitmap_ |= 8
		case "kubeconfig":
			value := iterator.ReadString()
			object.kubeconfig = value
			object.bitmap_ |= 16
		case "revocation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.revocationTimestamp = value
			object.bitmap_ |= 32
		case "status":
			text := iterator.ReadString()
			value := BreakGlassCredentialStatus(text)
			object.status = value
			object.bitmap_ |= 64
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.bitmap_ |= 128
		default:
			iterator.ReadAny()
		}
	}
	return object
}
