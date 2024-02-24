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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddon writes a value of the 'addon' type to the given writer.
func MarshalAddon(object *Addon, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddon(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddon writes a value of the 'addon' type to the given stream.
func writeAddon(object *Addon, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AddonLinkKind)
	} else {
		stream.WriteString(AddonKind)
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
	present_ = object.bitmap_&8 != 0 && object.commonAnnotations != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("common_annotations")
		if object.commonAnnotations != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.commonAnnotations))
			i := 0
			for key := range object.commonAnnotations {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.commonAnnotations[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.commonLabels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("common_labels")
		if object.commonLabels != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.commonLabels))
			i := 0
			for key := range object.commonLabels {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.commonLabels[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.config != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("config")
		writeAddonConfig(object.config, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.credentialsRequests != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("credentials_requests")
		writeCredentialRequestList(object.credentialsRequests, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("docs_link")
		stream.WriteString(object.docsLink)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("has_external_resources")
		stream.WriteBool(object.hasExternalResources)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hidden")
		stream.WriteBool(object.hidden)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("icon")
		stream.WriteString(object.icon)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("install_mode")
		stream.WriteString(string(object.installMode))
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("label")
		stream.WriteString(object.label)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_service")
		stream.WriteBool(object.managedService)
		count++
	}
	present_ = object.bitmap_&65536 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&131072 != 0 && object.namespaces != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespaces")
		writeAddonNamespaceList(object.namespaces, stream)
		count++
	}
	present_ = object.bitmap_&262144 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_name")
		stream.WriteString(object.operatorName)
		count++
	}
	present_ = object.bitmap_&524288 != 0 && object.parameters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("parameters")
		writeAddonParameters(object.parameters, stream)
		count++
	}
	present_ = object.bitmap_&1048576 != 0 && object.requirements != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requirements")
		writeAddonRequirementList(object.requirements, stream)
		count++
	}
	present_ = object.bitmap_&2097152 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_cost")
		stream.WriteFloat64(object.resourceCost)
		count++
	}
	present_ = object.bitmap_&4194304 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_name")
		stream.WriteString(object.resourceName)
		count++
	}
	present_ = object.bitmap_&8388608 != 0 && object.subOperators != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sub_operators")
		writeAddonSubOperatorList(object.subOperators, stream)
		count++
	}
	present_ = object.bitmap_&16777216 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("target_namespace")
		stream.WriteString(object.targetNamespace)
		count++
	}
	present_ = object.bitmap_&33554432 != 0 && object.version != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		writeAddonVersion(object.version, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddon reads a value of the 'addon' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddon(source interface{}) (object *Addon, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddon(iterator)
	err = iterator.Error
	return
}

// readAddon reads a value of the 'addon' type from the given iterator.
func readAddon(iterator *jsoniter.Iterator) *Addon {
	object := &Addon{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddonLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "common_annotations":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.commonAnnotations = value
			object.bitmap_ |= 8
		case "common_labels":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.commonLabels = value
			object.bitmap_ |= 16
		case "config":
			value := readAddonConfig(iterator)
			object.config = value
			object.bitmap_ |= 32
		case "credentials_requests":
			value := readCredentialRequestList(iterator)
			object.credentialsRequests = value
			object.bitmap_ |= 64
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 128
		case "docs_link":
			value := iterator.ReadString()
			object.docsLink = value
			object.bitmap_ |= 256
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 512
		case "has_external_resources":
			value := iterator.ReadBool()
			object.hasExternalResources = value
			object.bitmap_ |= 1024
		case "hidden":
			value := iterator.ReadBool()
			object.hidden = value
			object.bitmap_ |= 2048
		case "icon":
			value := iterator.ReadString()
			object.icon = value
			object.bitmap_ |= 4096
		case "install_mode":
			text := iterator.ReadString()
			value := AddonInstallMode(text)
			object.installMode = value
			object.bitmap_ |= 8192
		case "label":
			value := iterator.ReadString()
			object.label = value
			object.bitmap_ |= 16384
		case "managed_service":
			value := iterator.ReadBool()
			object.managedService = value
			object.bitmap_ |= 32768
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 65536
		case "namespaces":
			value := readAddonNamespaceList(iterator)
			object.namespaces = value
			object.bitmap_ |= 131072
		case "operator_name":
			value := iterator.ReadString()
			object.operatorName = value
			object.bitmap_ |= 262144
		case "parameters":
			value := readAddonParameters(iterator)
			object.parameters = value
			object.bitmap_ |= 524288
		case "requirements":
			value := readAddonRequirementList(iterator)
			object.requirements = value
			object.bitmap_ |= 1048576
		case "resource_cost":
			value := iterator.ReadFloat64()
			object.resourceCost = value
			object.bitmap_ |= 2097152
		case "resource_name":
			value := iterator.ReadString()
			object.resourceName = value
			object.bitmap_ |= 4194304
		case "sub_operators":
			value := readAddonSubOperatorList(iterator)
			object.subOperators = value
			object.bitmap_ |= 8388608
		case "target_namespace":
			value := iterator.ReadString()
			object.targetNamespace = value
			object.bitmap_ |= 16777216
		case "version":
			value := readAddonVersion(iterator)
			object.version = value
			object.bitmap_ |= 33554432
		default:
			iterator.ReadAny()
		}
	}
	return object
}
