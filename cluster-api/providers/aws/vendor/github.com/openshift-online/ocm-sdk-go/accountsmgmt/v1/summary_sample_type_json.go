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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalSummarySample writes a value of the 'summary_sample' type to the given writer.
func MarshalSummarySample(object *SummarySample, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSummarySample(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSummarySample writes a value of the 'summary_sample' type to the given stream.
func WriteSummarySample(object *SummarySample, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("time")
		stream.WriteString(object.time)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteFloat64(object.value)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSummarySample reads a value of the 'summary_sample' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSummarySample(source interface{}) (object *SummarySample, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadSummarySample(iterator)
	err = iterator.Error
	return
}

// ReadSummarySample reads a value of the 'summary_sample' type from the given iterator.
func ReadSummarySample(iterator *jsoniter.Iterator) *SummarySample {
	object := &SummarySample{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "time":
			value := iterator.ReadString()
			object.time = value
			object.bitmap_ |= 1
		case "value":
			value := iterator.ReadFloat64()
			object.value = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
