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

// MarshalAddonVersion writes a value of the 'addon_version' type to the given writer.
func MarshalAddonVersion(object *AddonVersion, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddonVersion(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddonVersion writes a value of the 'addon_version' type to the given stream.
func writeAddonVersion(object *AddonVersion, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AddonVersionLinkKind)
	} else {
		stream.WriteString(AddonVersionKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0 && object.additionalCatalogSources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_catalog_sources")
		writeAdditionalCatalogSourceList(object.additionalCatalogSources, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.availableUpgrades != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("available_upgrades")
		writeStringList(object.availableUpgrades, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("channel")
		stream.WriteString(object.channel)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.config != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("config")
		writeAddonConfig(object.config, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.parameters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("parameters")
		writeAddonParameterList(object.parameters, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("pull_secret_name")
		stream.WriteString(object.pullSecretName)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.requirements != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requirements")
		writeAddonRequirementList(object.requirements, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("source_image")
		stream.WriteString(object.sourceImage)
		count++
	}
	present_ = object.bitmap_&4096 != 0 && object.subOperators != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sub_operators")
		writeAddonSubOperatorList(object.subOperators, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonVersion reads a value of the 'addon_version' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonVersion(source interface{}) (object *AddonVersion, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddonVersion(iterator)
	err = iterator.Error
	return
}

// readAddonVersion reads a value of the 'addon_version' type from the given iterator.
func readAddonVersion(iterator *jsoniter.Iterator) *AddonVersion {
	object := &AddonVersion{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddonVersionLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "additional_catalog_sources":
			value := readAdditionalCatalogSourceList(iterator)
			object.additionalCatalogSources = value
			object.bitmap_ |= 8
		case "available_upgrades":
			value := readStringList(iterator)
			object.availableUpgrades = value
			object.bitmap_ |= 16
		case "channel":
			value := iterator.ReadString()
			object.channel = value
			object.bitmap_ |= 32
		case "config":
			value := readAddonConfig(iterator)
			object.config = value
			object.bitmap_ |= 64
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 128
		case "parameters":
			value := readAddonParameterList(iterator)
			object.parameters = value
			object.bitmap_ |= 256
		case "pull_secret_name":
			value := iterator.ReadString()
			object.pullSecretName = value
			object.bitmap_ |= 512
		case "requirements":
			value := readAddonRequirementList(iterator)
			object.requirements = value
			object.bitmap_ |= 1024
		case "source_image":
			value := iterator.ReadString()
			object.sourceImage = value
			object.bitmap_ |= 2048
		case "sub_operators":
			value := readAddonSubOperatorList(iterator)
			object.subOperators = value
			object.bitmap_ |= 4096
		default:
			iterator.ReadAny()
		}
	}
	return object
}
