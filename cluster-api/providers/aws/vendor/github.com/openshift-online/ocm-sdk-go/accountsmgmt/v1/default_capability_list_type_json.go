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

// MarshalDefaultCapabilityList writes a list of values of the 'default_capability' type to
// the given writer.
func MarshalDefaultCapabilityList(list []*DefaultCapability, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeDefaultCapabilityList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeDefaultCapabilityList writes a list of value of the 'default_capability' type to
// the given stream.
func writeDefaultCapabilityList(list []*DefaultCapability, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		writeDefaultCapability(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalDefaultCapabilityList reads a list of values of the 'default_capability' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalDefaultCapabilityList(source interface{}) (items []*DefaultCapability, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = readDefaultCapabilityList(iterator)
	err = iterator.Error
	return
}

// readDefaultCapabilityList reads list of values of the ”default_capability' type from
// the given iterator.
func readDefaultCapabilityList(iterator *jsoniter.Iterator) []*DefaultCapability {
	list := []*DefaultCapability{}
	for iterator.ReadArray() {
		item := readDefaultCapability(iterator)
		list = append(list, item)
	}
	return list
}
