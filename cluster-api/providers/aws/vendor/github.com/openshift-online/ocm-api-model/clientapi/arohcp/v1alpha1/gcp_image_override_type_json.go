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

// MarshalGCPImageOverride writes a value of the 'GCP_image_override' type to the given writer.
func MarshalGCPImageOverride(object *GCPImageOverride, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGCPImageOverride(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGCPImageOverride writes a value of the 'GCP_image_override' type to the given stream.
func WriteGCPImageOverride(object *GCPImageOverride, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(GCPImageOverrideLinkKind)
	} else {
		stream.WriteString(GCPImageOverrideKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.billingModel != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_model")
		v1.WriteBillingModelItem(object.billingModel, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("image_id")
		stream.WriteString(object.imageID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.product != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product")
		v1.WriteProduct(object.product, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("project_id")
		stream.WriteString(object.projectID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGCPImageOverride reads a value of the 'GCP_image_override' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGCPImageOverride(source interface{}) (object *GCPImageOverride, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGCPImageOverride(iterator)
	err = iterator.Error
	return
}

// ReadGCPImageOverride reads a value of the 'GCP_image_override' type from the given iterator.
func ReadGCPImageOverride(iterator *jsoniter.Iterator) *GCPImageOverride {
	object := &GCPImageOverride{
		fieldSet_: make([]bool, 7),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == GCPImageOverrideLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "billing_model":
			value := v1.ReadBillingModelItem(iterator)
			object.billingModel = value
			object.fieldSet_[3] = true
		case "image_id":
			value := iterator.ReadString()
			object.imageID = value
			object.fieldSet_[4] = true
		case "product":
			value := v1.ReadProduct(iterator)
			object.product = value
			object.fieldSet_[5] = true
		case "project_id":
			value := iterator.ReadString()
			object.projectID = value
			object.fieldSet_[6] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
