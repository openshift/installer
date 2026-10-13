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

// MarshalSubnetwork writes a value of the 'subnetwork' type to the given writer.
func MarshalSubnetwork(object *Subnetwork, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSubnetwork(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSubnetwork writes a value of the 'subnetwork' type to the given stream.
func WriteSubnetwork(object *Subnetwork, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cidr_block")
		stream.WriteString(object.cidrBlock)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("public")
		stream.WriteBool(object.public)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("red_hat_managed")
		stream.WriteBool(object.redHatManaged)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet_id")
		stream.WriteString(object.subnetID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSubnetwork reads a value of the 'subnetwork' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSubnetwork(source interface{}) (object *Subnetwork, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadSubnetwork(iterator)
	err = iterator.Error
	return
}

// ReadSubnetwork reads a value of the 'subnetwork' type from the given iterator.
func ReadSubnetwork(iterator *jsoniter.Iterator) *Subnetwork {
	object := &Subnetwork{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cidr_block":
			value := iterator.ReadString()
			object.cidrBlock = value
			object.fieldSet_[0] = true
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.fieldSet_[1] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[2] = true
		case "public":
			value := iterator.ReadBool()
			object.public = value
			object.fieldSet_[3] = true
		case "red_hat_managed":
			value := iterator.ReadBool()
			object.redHatManaged = value
			object.fieldSet_[4] = true
		case "subnet_id":
			value := iterator.ReadString()
			object.subnetID = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
