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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalDNS writes a value of the 'DNS' type to the given writer.
func MarshalDNS(object *DNS, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeDNS(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeDNS writes a value of the 'DNS' type to the given stream.
func writeDNS(object *DNS, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("base_domain")
		stream.WriteString(object.baseDomain)
	}
	stream.WriteObjectEnd()
}

// UnmarshalDNS reads a value of the 'DNS' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalDNS(source interface{}) (object *DNS, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readDNS(iterator)
	err = iterator.Error
	return
}

// readDNS reads a value of the 'DNS' type from the given iterator.
func readDNS(iterator *jsoniter.Iterator) *DNS {
	object := &DNS{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "base_domain":
			value := iterator.ReadString()
			object.baseDomain = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
