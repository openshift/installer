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

	jsoniter "github.com/json-iterator/go"
	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAMIOverride writes a value of the 'AMI_override' type to the given writer.
func MarshalAMIOverride(object *AMIOverride, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAMIOverride(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAMIOverride writes a value of the 'AMI_override' type to the given stream.
func WriteAMIOverride(object *AMIOverride, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AMIOverrideLinkKind)
	} else {
		stream.WriteString(AMIOverrideKind)
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
		stream.WriteObjectField("ami")
		stream.WriteString(object.ami)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.product != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product")
		v1.WriteProduct(object.product, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.region != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		v1.WriteCloudRegion(object.region, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAMIOverride reads a value of the 'AMI_override' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAMIOverride(source interface{}) (object *AMIOverride, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAMIOverride(iterator)
	err = iterator.Error
	return
}

// ReadAMIOverride reads a value of the 'AMI_override' type from the given iterator.
func ReadAMIOverride(iterator *jsoniter.Iterator) *AMIOverride {
	object := &AMIOverride{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AMIOverrideLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "ami":
			value := iterator.ReadString()
			object.ami = value
			object.fieldSet_[3] = true
		case "product":
			value := v1.ReadProduct(iterator)
			object.product = value
			object.fieldSet_[4] = true
		case "region":
			value := v1.ReadCloudRegion(iterator)
			object.region = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
