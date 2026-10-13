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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalDataPlaneOperatorIdentityRequirementList writes a list of values of the 'data_plane_operator_identity_requirement' type to
// the given writer.
func MarshalDataPlaneOperatorIdentityRequirementList(list []*DataPlaneOperatorIdentityRequirement, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteDataPlaneOperatorIdentityRequirementList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteDataPlaneOperatorIdentityRequirementList writes a list of value of the 'data_plane_operator_identity_requirement' type to
// the given stream.
func WriteDataPlaneOperatorIdentityRequirementList(list []*DataPlaneOperatorIdentityRequirement, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteDataPlaneOperatorIdentityRequirement(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalDataPlaneOperatorIdentityRequirementList reads a list of values of the 'data_plane_operator_identity_requirement' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalDataPlaneOperatorIdentityRequirementList(source interface{}) (items []*DataPlaneOperatorIdentityRequirement, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadDataPlaneOperatorIdentityRequirementList(iterator)
	err = iterator.Error
	return
}

// ReadDataPlaneOperatorIdentityRequirementList reads list of values of the ”data_plane_operator_identity_requirement' type from
// the given iterator.
func ReadDataPlaneOperatorIdentityRequirementList(iterator *jsoniter.Iterator) []*DataPlaneOperatorIdentityRequirement {
	list := []*DataPlaneOperatorIdentityRequirement{}
	for iterator.ReadArray() {
		item := ReadDataPlaneOperatorIdentityRequirement(iterator)
		list = append(list, item)
	}
	return list
}
