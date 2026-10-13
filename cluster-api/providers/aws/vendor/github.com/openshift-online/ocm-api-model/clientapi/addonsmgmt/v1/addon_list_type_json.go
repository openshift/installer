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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAddonList writes a list of values of the 'addon' type to
// the given writer.
func MarshalAddonList(list []*Addon, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonList writes a list of value of the 'addon' type to
// the given stream.
func WriteAddonList(list []*Addon, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteAddon(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalAddonList reads a list of values of the 'addon' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAddonList(source interface{}) (items []*Addon, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadAddonList(iterator)
	err = iterator.Error
	return
}

// ReadAddonList reads list of values of the ”addon' type from
// the given iterator.
func ReadAddonList(iterator *jsoniter.Iterator) []*Addon {
	list := []*Addon{}
	for iterator.ReadArray() {
		item := ReadAddon(iterator)
		list = append(list, item)
	}
	return list
}
