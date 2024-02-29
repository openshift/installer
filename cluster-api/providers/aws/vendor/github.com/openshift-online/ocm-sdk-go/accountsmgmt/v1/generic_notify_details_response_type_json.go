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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalGenericNotifyDetailsResponse writes a value of the 'generic_notify_details_response' type to the given writer.
func MarshalGenericNotifyDetailsResponse(object *GenericNotifyDetailsResponse, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeGenericNotifyDetailsResponse(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeGenericNotifyDetailsResponse writes a value of the 'generic_notify_details_response' type to the given stream.
func writeGenericNotifyDetailsResponse(object *GenericNotifyDetailsResponse, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(GenericNotifyDetailsResponseLinkKind)
	} else {
		stream.WriteString(GenericNotifyDetailsResponseKind)
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
	present_ = object.bitmap_&8 != 0 && object.associates != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("associates")
		writeStringList(object.associates, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.items != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("items")
		writeNotificationDetailsResponseList(object.items, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.recipients != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("recipients")
		writeStringList(object.recipients, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGenericNotifyDetailsResponse reads a value of the 'generic_notify_details_response' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGenericNotifyDetailsResponse(source interface{}) (object *GenericNotifyDetailsResponse, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readGenericNotifyDetailsResponse(iterator)
	err = iterator.Error
	return
}

// readGenericNotifyDetailsResponse reads a value of the 'generic_notify_details_response' type from the given iterator.
func readGenericNotifyDetailsResponse(iterator *jsoniter.Iterator) *GenericNotifyDetailsResponse {
	object := &GenericNotifyDetailsResponse{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == GenericNotifyDetailsResponseLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "associates":
			value := readStringList(iterator)
			object.associates = value
			object.bitmap_ |= 8
		case "items":
			value := readNotificationDetailsResponseList(iterator)
			object.items = value
			object.bitmap_ |= 16
		case "recipients":
			value := readStringList(iterator)
			object.recipients = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
