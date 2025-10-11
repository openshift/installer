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

// MarshalQuotaRules writes a value of the 'quota_rules' type to the given writer.
func MarshalQuotaRules(object *QuotaRules, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteQuotaRules(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteQuotaRules writes a value of the 'quota_rules' type to the given stream.
func WriteQuotaRules(object *QuotaRules, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_model")
		stream.WriteString(object.billingModel)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("byoc")
		stream.WriteString(object.byoc)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud")
		stream.WriteString(object.cloud)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cost")
		stream.WriteInt(object.cost)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product")
		stream.WriteString(object.product)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("quota_id")
		stream.WriteString(object.quotaId)
	}
	stream.WriteObjectEnd()
}

// UnmarshalQuotaRules reads a value of the 'quota_rules' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalQuotaRules(source interface{}) (object *QuotaRules, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadQuotaRules(iterator)
	err = iterator.Error
	return
}

// ReadQuotaRules reads a value of the 'quota_rules' type from the given iterator.
func ReadQuotaRules(iterator *jsoniter.Iterator) *QuotaRules {
	object := &QuotaRules{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.bitmap_ |= 1
		case "billing_model":
			value := iterator.ReadString()
			object.billingModel = value
			object.bitmap_ |= 2
		case "byoc":
			value := iterator.ReadString()
			object.byoc = value
			object.bitmap_ |= 4
		case "cloud":
			value := iterator.ReadString()
			object.cloud = value
			object.bitmap_ |= 8
		case "cost":
			value := iterator.ReadInt()
			object.cost = value
			object.bitmap_ |= 16
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 32
		case "product":
			value := iterator.ReadString()
			object.product = value
			object.bitmap_ |= 64
		case "quota_id":
			value := iterator.ReadString()
			object.quotaId = value
			object.bitmap_ |= 128
		default:
			iterator.ReadAny()
		}
	}
	return object
}
