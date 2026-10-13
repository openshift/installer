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

// MarshalAddOnInstallationParameterList writes a list of values of the 'add_on_installation_parameter' type to
// the given writer.
func MarshalAddOnInstallationParameterList(list []*AddOnInstallationParameter, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddOnInstallationParameterList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddOnInstallationParameterList writes a list of value of the 'add_on_installation_parameter' type to
// the given stream.
func WriteAddOnInstallationParameterList(list []*AddOnInstallationParameter, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteAddOnInstallationParameter(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalAddOnInstallationParameterList reads a list of values of the 'add_on_installation_parameter' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAddOnInstallationParameterList(source interface{}) (items []*AddOnInstallationParameter, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadAddOnInstallationParameterList(iterator)
	err = iterator.Error
	return
}

// ReadAddOnInstallationParameterList reads list of values of the ”add_on_installation_parameter' type from
// the given iterator.
func ReadAddOnInstallationParameterList(iterator *jsoniter.Iterator) []*AddOnInstallationParameter {
	list := []*AddOnInstallationParameter{}
	for iterator.ReadArray() {
		item := ReadAddOnInstallationParameter(iterator)
		list = append(list, item)
	}
	return list
}
