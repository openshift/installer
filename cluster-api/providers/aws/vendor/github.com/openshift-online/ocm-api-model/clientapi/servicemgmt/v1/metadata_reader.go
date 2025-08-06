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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalMetadata writes a value of the metadata type to the given target, which
// can be a writer or a JSON encoder.
func MarshalMetadata(object *Metadata, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeMetadata(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}
func writeMetadata(object *Metadata, stream *jsoniter.Stream) {
	stream.WriteObjectStart()
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteObjectField("server_version")
		stream.WriteString(object.serverVersion)
	}
	stream.WriteObjectEnd()
}

// UnmarshalMetadata reads a value of the metadata type from the given source, which
// which can be a reader, a slice of byte or a string.
func UnmarshalMetadata(source interface{}) (object *Metadata, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readMetadata(iterator)
	err = iterator.Error
	return
}
func readMetadata(iterator *jsoniter.Iterator) *Metadata {
	object := &Metadata{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "server_version":
			object.serverVersion = iterator.ReadString()
			if len(object.fieldSet_) <= 0 {
				object.fieldSet_ = make([]bool, 1)
			}
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
