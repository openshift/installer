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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddonConfigList writes a list of values of the 'addon_config' type to
// the given writer.
func MarshalAddonConfigList(list []*AddonConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddonConfigList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddonConfigList writes a list of value of the 'addon_config' type to
// the given stream.
func writeAddonConfigList(list []*AddonConfig, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		writeAddonConfig(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalAddonConfigList reads a list of values of the 'addon_config' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAddonConfigList(source interface{}) (items []*AddonConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = readAddonConfigList(iterator)
	err = iterator.Error
	return
}

// readAddonConfigList reads list of values of the ”addon_config' type from
// the given iterator.
func readAddonConfigList(iterator *jsoniter.Iterator) []*AddonConfig {
	list := []*AddonConfig{}
	for iterator.ReadArray() {
		item := readAddonConfig(iterator)
		list = append(list, item)
	}
	return list
}
