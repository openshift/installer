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

// MarshalNamespaceSelectorList writes a list of values of the 'namespace_selector' type to
// the given writer.
func MarshalNamespaceSelectorList(list []*NamespaceSelector, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNamespaceSelectorList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNamespaceSelectorList writes a list of value of the 'namespace_selector' type to
// the given stream.
func WriteNamespaceSelectorList(list []*NamespaceSelector, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteNamespaceSelector(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalNamespaceSelectorList reads a list of values of the 'namespace_selector' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalNamespaceSelectorList(source interface{}) (items []*NamespaceSelector, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadNamespaceSelectorList(iterator)
	err = iterator.Error
	return
}

// ReadNamespaceSelectorList reads list of values of the ”namespace_selector' type from
// the given iterator.
func ReadNamespaceSelectorList(iterator *jsoniter.Iterator) []*NamespaceSelector {
	list := []*NamespaceSelector{}
	for iterator.ReadArray() {
		item := ReadNamespaceSelector(iterator)
		list = append(list, item)
	}
	return list
}
