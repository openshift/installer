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

// MarshalOidcConfigList writes a list of values of the 'oidc_config' type to
// the given writer.
func MarshalOidcConfigList(list []*OidcConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteOidcConfigList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteOidcConfigList writes a list of value of the 'oidc_config' type to
// the given stream.
func WriteOidcConfigList(list []*OidcConfig, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteOidcConfig(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalOidcConfigList reads a list of values of the 'oidc_config' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalOidcConfigList(source interface{}) (items []*OidcConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadOidcConfigList(iterator)
	err = iterator.Error
	return
}

// ReadOidcConfigList reads list of values of the ”oidc_config' type from
// the given iterator.
func ReadOidcConfigList(iterator *jsoniter.Iterator) []*OidcConfig {
	list := []*OidcConfig{}
	for iterator.ReadArray() {
		item := ReadOidcConfig(iterator)
		list = append(list, item)
	}
	return list
}
