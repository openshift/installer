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

// MarshalAWSSpotMarketOptions writes a value of the 'AWS_spot_market_options' type to the given writer.
func MarshalAWSSpotMarketOptions(object *AWSSpotMarketOptions, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWSSpotMarketOptions(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWSSpotMarketOptions writes a value of the 'AWS_spot_market_options' type to the given stream.
func WriteAWSSpotMarketOptions(object *AWSSpotMarketOptions, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AWSSpotMarketOptionsLinkKind)
	} else {
		stream.WriteString(AWSSpotMarketOptionsKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_price")
		stream.WriteFloat64(object.maxPrice)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWSSpotMarketOptions reads a value of the 'AWS_spot_market_options' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSSpotMarketOptions(source interface{}) (object *AWSSpotMarketOptions, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAWSSpotMarketOptions(iterator)
	err = iterator.Error
	return
}

// ReadAWSSpotMarketOptions reads a value of the 'AWS_spot_market_options' type from the given iterator.
func ReadAWSSpotMarketOptions(iterator *jsoniter.Iterator) *AWSSpotMarketOptions {
	object := &AWSSpotMarketOptions{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AWSSpotMarketOptionsLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "max_price":
			value := iterator.ReadFloat64()
			object.maxPrice = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
