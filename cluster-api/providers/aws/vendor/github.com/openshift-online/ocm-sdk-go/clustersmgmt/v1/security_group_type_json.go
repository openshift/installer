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

// MarshalSecurityGroup writes a value of the 'security_group' type to the given writer.
func MarshalSecurityGroup(object *SecurityGroup, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSecurityGroup(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSecurityGroup writes a value of the 'security_group' type to the given stream.
func writeSecurityGroup(object *SecurityGroup, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
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
		stream.WriteObjectField("red_hat_managed")
		stream.WriteBool(object.redHatManaged)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSecurityGroup reads a value of the 'security_group' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSecurityGroup(source interface{}) (object *SecurityGroup, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSecurityGroup(iterator)
	err = iterator.Error
	return
}

// readSecurityGroup reads a value of the 'security_group' type from the given iterator.
func readSecurityGroup(iterator *jsoniter.Iterator) *SecurityGroup {
	object := &SecurityGroup{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 1
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 2
		case "red_hat_managed":
			value := iterator.ReadBool()
			object.redHatManaged = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
