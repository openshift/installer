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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalSubnetwork writes a value of the 'subnetwork' type to the given writer.
func MarshalSubnetwork(object *Subnetwork, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSubnetwork(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSubnetwork writes a value of the 'subnetwork' type to the given stream.
func writeSubnetwork(object *Subnetwork, stream *jsoniter.Stream) {
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
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("public")
		stream.WriteBool(object.public)
		count++
	}
	present_ = object.bitmap_&8 != 0
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
	object = readSubnetwork(iterator)
	err = iterator.Error
	return
}

// readSubnetwork reads a value of the 'subnetwork' type from the given iterator.
func readSubnetwork(iterator *jsoniter.Iterator) *Subnetwork {
	object := &Subnetwork{}
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
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 2
		case "public":
			value := iterator.ReadBool()
			object.public = value
			object.bitmap_ |= 4
		case "subnet_id":
			value := iterator.ReadString()
			object.subnetID = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
