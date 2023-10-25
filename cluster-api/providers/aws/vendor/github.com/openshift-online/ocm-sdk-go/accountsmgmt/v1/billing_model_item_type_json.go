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

// MarshalBillingModelItem writes a value of the 'billing_model_item' type to the given writer.
func MarshalBillingModelItem(object *BillingModelItem, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeBillingModelItem(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeBillingModelItem writes a value of the 'billing_model_item' type to the given stream.
func writeBillingModelItem(object *BillingModelItem, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(BillingModelItemLinkKind)
	} else {
		stream.WriteString(BillingModelItemKind)
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
		stream.WriteObjectField("billing_model_type")
		stream.WriteString(object.billingModelType)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("marketplace")
		stream.WriteString(object.marketplace)
	}
	stream.WriteObjectEnd()
}

// UnmarshalBillingModelItem reads a value of the 'billing_model_item' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalBillingModelItem(source interface{}) (object *BillingModelItem, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readBillingModelItem(iterator)
	err = iterator.Error
	return
}

// readBillingModelItem reads a value of the 'billing_model_item' type from the given iterator.
func readBillingModelItem(iterator *jsoniter.Iterator) *BillingModelItem {
	object := &BillingModelItem{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == BillingModelItemLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "billing_model_type":
			value := iterator.ReadString()
			object.billingModelType = value
			object.bitmap_ |= 8
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 16
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.bitmap_ |= 32
		case "marketplace":
			value := iterator.ReadString()
			object.marketplace = value
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
