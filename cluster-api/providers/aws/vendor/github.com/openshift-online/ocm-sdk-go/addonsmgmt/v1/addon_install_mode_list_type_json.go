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

// MarshalAddonInstallModeList writes a list of values of the 'addon_install_mode' type to
// the given writer.
func MarshalAddonInstallModeList(list []AddonInstallMode, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonInstallModeList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonInstallModeList writes a list of value of the 'addon_install_mode' type to
// the given stream.
func WriteAddonInstallModeList(list []AddonInstallMode, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		stream.WriteString(string(value))
	}
	stream.WriteArrayEnd()
}

// UnmarshalAddonInstallModeList reads a list of values of the 'addon_install_mode' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalAddonInstallModeList(source interface{}) (items []AddonInstallMode, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadAddonInstallModeList(iterator)
	err = iterator.Error
	return
}

// ReadAddonInstallModeList reads list of values of the ”addon_install_mode' type from
// the given iterator.
func ReadAddonInstallModeList(iterator *jsoniter.Iterator) []AddonInstallMode {
	list := []AddonInstallMode{}
	for iterator.ReadArray() {
		text := iterator.ReadString()
		item := AddonInstallMode(text)
		list = append(list, item)
	}
	return list
}
