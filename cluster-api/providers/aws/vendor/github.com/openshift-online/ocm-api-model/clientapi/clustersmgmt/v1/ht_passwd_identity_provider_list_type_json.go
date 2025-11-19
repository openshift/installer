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

// MarshalHTPasswdIdentityProviderList writes a list of values of the 'HT_passwd_identity_provider' type to
// the given writer.
func MarshalHTPasswdIdentityProviderList(list []*HTPasswdIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteHTPasswdIdentityProviderList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteHTPasswdIdentityProviderList writes a list of value of the 'HT_passwd_identity_provider' type to
// the given stream.
func WriteHTPasswdIdentityProviderList(list []*HTPasswdIdentityProvider, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteHTPasswdIdentityProvider(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalHTPasswdIdentityProviderList reads a list of values of the 'HT_passwd_identity_provider' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalHTPasswdIdentityProviderList(source interface{}) (items []*HTPasswdIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadHTPasswdIdentityProviderList(iterator)
	err = iterator.Error
	return
}

// ReadHTPasswdIdentityProviderList reads list of values of the ”HT_passwd_identity_provider' type from
// the given iterator.
func ReadHTPasswdIdentityProviderList(iterator *jsoniter.Iterator) []*HTPasswdIdentityProvider {
	list := []*HTPasswdIdentityProvider{}
	for iterator.ReadArray() {
		item := ReadHTPasswdIdentityProvider(iterator)
		list = append(list, item)
	}
	return list
}
