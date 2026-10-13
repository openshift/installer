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

// MarshalGitlabIdentityProviderList writes a list of values of the 'gitlab_identity_provider' type to
// the given writer.
func MarshalGitlabIdentityProviderList(list []*GitlabIdentityProvider, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGitlabIdentityProviderList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGitlabIdentityProviderList writes a list of value of the 'gitlab_identity_provider' type to
// the given stream.
func WriteGitlabIdentityProviderList(list []*GitlabIdentityProvider, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteGitlabIdentityProvider(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalGitlabIdentityProviderList reads a list of values of the 'gitlab_identity_provider' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalGitlabIdentityProviderList(source interface{}) (items []*GitlabIdentityProvider, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadGitlabIdentityProviderList(iterator)
	err = iterator.Error
	return
}

// ReadGitlabIdentityProviderList reads list of values of the ”gitlab_identity_provider' type from
// the given iterator.
func ReadGitlabIdentityProviderList(iterator *jsoniter.Iterator) []*GitlabIdentityProvider {
	list := []*GitlabIdentityProvider{}
	for iterator.ReadArray() {
		item := ReadGitlabIdentityProvider(iterator)
		list = append(list, item)
	}
	return list
}
