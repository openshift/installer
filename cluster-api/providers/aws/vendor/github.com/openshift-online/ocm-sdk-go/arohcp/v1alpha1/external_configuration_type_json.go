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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/helpers"
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
	present_ = object.bitmap_&1 != 0 && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteLabelList(object.labels.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.manifests != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("manifests")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteManifestList(object.manifests.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.syncsets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("syncsets")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteSyncsetList(object.syncsets.Items(), stream)
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
	object := &ExternalConfiguration{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "labels":
			value := &v1.LabelList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.LabelListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadLabelList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.labels = value
			object.bitmap_ |= 1
		case "manifests":
			value := &v1.ManifestList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.ManifestListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadManifestList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.manifests = value
			object.bitmap_ |= 2
		case "syncsets":
			value := &v1.SyncsetList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.SyncsetListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadSyncsetList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.syncsets = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
