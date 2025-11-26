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

// MarshalControlPlane writes a value of the 'control_plane' type to the given writer.
func MarshalControlPlane(object *ControlPlane, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteControlPlane(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteControlPlane writes a value of the 'control_plane' type to the given stream.
func WriteControlPlane(object *ControlPlane, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.backup != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("backup")
		WriteBackup(object.backup, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalControlPlane reads a value of the 'control_plane' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalControlPlane(source interface{}) (object *ControlPlane, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadControlPlane(iterator)
	err = iterator.Error
	return
}

// ReadControlPlane reads a value of the 'control_plane' type from the given iterator.
func ReadControlPlane(iterator *jsoniter.Iterator) *ControlPlane {
	object := &ControlPlane{
		fieldSet_: make([]bool, 1),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "backup":
			value := ReadBackup(iterator)
			object.backup = value
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
