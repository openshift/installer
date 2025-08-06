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

// MarshalExternalConfiguration writes a value of the 'external_configuration' type to the given writer.
func MarshalExternalConfiguration(object *ExternalConfiguration, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteExternalConfiguration(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteExternalConfiguration writes a value of the 'external_configuration' type to the given stream.
func WriteExternalConfiguration(object *ExternalConfiguration, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteLabelList(object.labels.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.manifests != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("manifests")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteManifestList(object.manifests.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.syncsets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("syncsets")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteSyncsetList(object.syncsets.Items(), stream)
		stream.WriteObjectEnd()
	}
	stream.WriteObjectEnd()
}

// UnmarshalExternalConfiguration reads a value of the 'external_configuration' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalExternalConfiguration(source interface{}) (object *ExternalConfiguration, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadExternalConfiguration(iterator)
	err = iterator.Error
	return
}

// ReadExternalConfiguration reads a value of the 'external_configuration' type from the given iterator.
func ReadExternalConfiguration(iterator *jsoniter.Iterator) *ExternalConfiguration {
	object := &ExternalConfiguration{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "labels":
			value := &LabelList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == LabelListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadLabelList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.labels = value
			object.fieldSet_[0] = true
		case "manifests":
			value := &ManifestList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == ManifestListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadManifestList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.manifests = value
			object.fieldSet_[1] = true
		case "syncsets":
			value := &SyncsetList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == SyncsetListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadSyncsetList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.syncsets = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
