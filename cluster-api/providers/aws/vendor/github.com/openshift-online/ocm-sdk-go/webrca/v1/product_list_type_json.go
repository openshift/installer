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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalProductList writes a list of values of the 'product' type to
// the given writer.
func MarshalProductList(list []*Product, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeProductList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeProductList writes a list of value of the 'product' type to
// the given stream.
func writeProductList(list []*Product, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		writeProduct(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalProductList reads a list of values of the 'product' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalProductList(source interface{}) (items []*Product, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = readProductList(iterator)
	err = iterator.Error
	return
}

// readProductList reads list of values of the ”product' type from
// the given iterator.
func readProductList(iterator *jsoniter.Iterator) []*Product {
	list := []*Product{}
	for iterator.ReadArray() {
		item := readProduct(iterator)
		list = append(list, item)
	}
	return list
}
