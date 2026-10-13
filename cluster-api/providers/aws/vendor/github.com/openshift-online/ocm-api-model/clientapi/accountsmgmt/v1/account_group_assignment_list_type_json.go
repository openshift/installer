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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAccountGroupAssignmentList writes a list of values of the 'account_group_assignment' type to
// the given writer.
func MarshalAccountGroupAssignmentList(list []*AccountGroupAssignment, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAccountGroupAssignmentList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAccountGroupAssignmentList writes a list of value of the 'account_group_assignment' type to
// the given stream.
func WriteAccountGroupAssignmentList(list []*AccountGroupAssignment, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteAccountGroupAssignment(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalAccountGroupAssignmentList reads a list of values of the 'account_group_assignment' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAccountGroupAssignmentList(source interface{}) (items []*AccountGroupAssignment, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadAccountGroupAssignmentList(iterator)
	err = iterator.Error
	return
}

// ReadAccountGroupAssignmentList reads list of values of the ”account_group_assignment' type from
// the given iterator.
func ReadAccountGroupAssignmentList(iterator *jsoniter.Iterator) []*AccountGroupAssignment {
	list := []*AccountGroupAssignment{}
	for iterator.ReadArray() {
		item := ReadAccountGroupAssignment(iterator)
		list = append(list, item)
	}
	return list
}
