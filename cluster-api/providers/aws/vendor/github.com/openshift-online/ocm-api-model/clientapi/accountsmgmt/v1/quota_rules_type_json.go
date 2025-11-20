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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
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
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_model")
		stream.WriteString(object.billingModel)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("byoc")
		stream.WriteString(object.byoc)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud")
		stream.WriteString(object.cloud)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cost")
		stream.WriteInt(object.cost)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product")
		stream.WriteString(object.product)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
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
	object := &QuotaRules{
		fieldSet_: make([]bool, 8),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.fieldSet_[0] = true
		case "billing_model":
			value := iterator.ReadString()
			object.billingModel = value
			object.fieldSet_[1] = true
		case "byoc":
			value := iterator.ReadString()
			object.byoc = value
			object.fieldSet_[2] = true
		case "cloud":
			value := iterator.ReadString()
			object.cloud = value
			object.fieldSet_[3] = true
		case "cost":
			value := iterator.ReadInt()
			object.cost = value
			object.fieldSet_[4] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[5] = true
		case "product":
			value := iterator.ReadString()
			object.product = value
			object.fieldSet_[6] = true
		case "quota_id":
			value := iterator.ReadString()
			object.quotaId = value
			object.fieldSet_[7] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
