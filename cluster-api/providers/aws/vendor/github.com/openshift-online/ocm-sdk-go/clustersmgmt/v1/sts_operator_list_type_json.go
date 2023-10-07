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

// MarshalSTSOperatorList writes a list of values of the 'STS_operator' type to
// the given writer.
func MarshalSTSOperatorList(list []*STSOperator, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSTSOperatorList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSTSOperatorList writes a list of value of the 'STS_operator' type to
// the given stream.
func writeSTSOperatorList(list []*STSOperator, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		writeSTSOperator(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalSTSOperatorList reads a list of values of the 'STS_operator' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalSTSOperatorList(source interface{}) (items []*STSOperator, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = readSTSOperatorList(iterator)
	err = iterator.Error
	return
}

// readSTSOperatorList reads list of values of the ”STS_operator' type from
// the given iterator.
func readSTSOperatorList(iterator *jsoniter.Iterator) []*STSOperator {
	list := []*STSOperator{}
	for iterator.ReadArray() {
		item := readSTSOperator(iterator)
		list = append(list, item)
	}
	return list
}
