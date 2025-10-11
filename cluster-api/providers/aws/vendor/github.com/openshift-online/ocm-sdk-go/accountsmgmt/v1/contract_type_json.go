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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalContract writes a value of the 'contract' type to the given writer.
func MarshalContract(object *Contract, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteContract(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteContract writes a value of the 'contract' type to the given stream.
func WriteContract(object *Contract, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.dimensions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dimensions")
		WriteContractDimensionList(object.dimensions, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("end_date")
		stream.WriteString((object.endDate).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("start_date")
		stream.WriteString((object.startDate).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalContract reads a value of the 'contract' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalContract(source interface{}) (object *Contract, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadContract(iterator)
	err = iterator.Error
	return
}

// ReadContract reads a value of the 'contract' type from the given iterator.
func ReadContract(iterator *jsoniter.Iterator) *Contract {
	object := &Contract{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "dimensions":
			value := ReadContractDimensionList(iterator)
			object.dimensions = value
			object.bitmap_ |= 1
		case "end_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.endDate = value
			object.bitmap_ |= 2
		case "start_date":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.startDate = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
