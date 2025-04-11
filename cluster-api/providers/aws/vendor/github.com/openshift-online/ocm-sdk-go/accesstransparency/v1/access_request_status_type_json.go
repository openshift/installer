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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAccessRequestStatus writes a value of the 'access_request_status' type to the given writer.
func MarshalAccessRequestStatus(object *AccessRequestStatus, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAccessRequestStatus(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAccessRequestStatus writes a value of the 'access_request_status' type to the given stream.
func writeAccessRequestStatus(object *AccessRequestStatus, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("expires_at")
		stream.WriteString((object.expiresAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
	}
	stream.WriteObjectEnd()
}

// UnmarshalAccessRequestStatus reads a value of the 'access_request_status' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAccessRequestStatus(source interface{}) (object *AccessRequestStatus, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAccessRequestStatus(iterator)
	err = iterator.Error
	return
}

// readAccessRequestStatus reads a value of the 'access_request_status' type from the given iterator.
func readAccessRequestStatus(iterator *jsoniter.Iterator) *AccessRequestStatus {
	object := &AccessRequestStatus{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "expires_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.expiresAt = value
			object.bitmap_ |= 1
		case "state":
			text := iterator.ReadString()
			value := AccessRequestState(text)
			object.state = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
