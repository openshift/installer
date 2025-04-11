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

// MarshalRegistryLocation writes a value of the 'registry_location' type to the given writer.
func MarshalRegistryLocation(object *RegistryLocation, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeRegistryLocation(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeRegistryLocation writes a value of the 'registry_location' type to the given stream.
func writeRegistryLocation(object *RegistryLocation, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("domain_name")
		stream.WriteString(object.domainName)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("insecure")
		stream.WriteBool(object.insecure)
	}
	stream.WriteObjectEnd()
}

// UnmarshalRegistryLocation reads a value of the 'registry_location' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRegistryLocation(source interface{}) (object *RegistryLocation, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readRegistryLocation(iterator)
	err = iterator.Error
	return
}

// readRegistryLocation reads a value of the 'registry_location' type from the given iterator.
func readRegistryLocation(iterator *jsoniter.Iterator) *RegistryLocation {
	object := &RegistryLocation{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "domain_name":
			value := iterator.ReadString()
			object.domainName = value
			object.bitmap_ |= 1
		case "insecure":
			value := iterator.ReadBool()
			object.insecure = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
