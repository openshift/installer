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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAddOn writes a value of the 'add_on' type to the given writer.
func MarshalAddOn(object *AddOn, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddOn(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddOn writes a value of the 'add_on' type to the given stream.
func WriteAddOn(object *AddOn, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AddOnLinkKind)
	} else {
		stream.WriteString(AddOnKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.commonAnnotations != nil
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
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.commonLabels != nil
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
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.config != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("config")
		WriteAddOnConfig(object.config, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.credentialsRequests != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("credentials_requests")
		WriteCredentialRequestList(object.credentialsRequests, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("docs_link")
		stream.WriteString(object.docsLink)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("has_external_resources")
		stream.WriteBool(object.hasExternalResources)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hidden")
		stream.WriteBool(object.hidden)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("icon")
		stream.WriteString(object.icon)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("install_mode")
		stream.WriteString(string(object.installMode))
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("label")
		stream.WriteString(object.label)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_service")
		stream.WriteBool(object.managedService)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17] && object.namespaces != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespaces")
		WriteAddOnNamespaceList(object.namespaces, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 18 && object.fieldSet_[18]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_name")
		stream.WriteString(object.operatorName)
		count++
	}
	present_ = len(object.fieldSet_) > 19 && object.fieldSet_[19] && object.parameters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("parameters")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteAddOnParameterList(object.parameters.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 20 && object.fieldSet_[20] && object.requirements != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requirements")
		WriteAddOnRequirementList(object.requirements, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 21 && object.fieldSet_[21]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_cost")
		stream.WriteFloat64(object.resourceCost)
		count++
	}
	present_ = len(object.fieldSet_) > 22 && object.fieldSet_[22]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_name")
		stream.WriteString(object.resourceName)
		count++
	}
	present_ = len(object.fieldSet_) > 23 && object.fieldSet_[23] && object.subOperators != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sub_operators")
		WriteAddOnSubOperatorList(object.subOperators, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 24 && object.fieldSet_[24]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("target_namespace")
		stream.WriteString(object.targetNamespace)
		count++
	}
	present_ = len(object.fieldSet_) > 25 && object.fieldSet_[25] && object.version != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		WriteAddOnVersion(object.version, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOn reads a value of the 'add_on' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOn(source interface{}) (object *AddOn, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddOn(iterator)
	err = iterator.Error
	return
}

// ReadAddOn reads a value of the 'add_on' type from the given iterator.
func ReadAddOn(iterator *jsoniter.Iterator) *AddOn {
	object := &AddOn{
		fieldSet_: make([]bool, 26),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddOnLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
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
			object.fieldSet_[3] = true
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
			object.fieldSet_[4] = true
		case "config":
			value := ReadAddOnConfig(iterator)
			object.config = value
			object.fieldSet_[5] = true
		case "credentials_requests":
			value := ReadCredentialRequestList(iterator)
			object.credentialsRequests = value
			object.fieldSet_[6] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[7] = true
		case "docs_link":
			value := iterator.ReadString()
			object.docsLink = value
			object.fieldSet_[8] = true
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[9] = true
		case "has_external_resources":
			value := iterator.ReadBool()
			object.hasExternalResources = value
			object.fieldSet_[10] = true
		case "hidden":
			value := iterator.ReadBool()
			object.hidden = value
			object.fieldSet_[11] = true
		case "icon":
			value := iterator.ReadString()
			object.icon = value
			object.fieldSet_[12] = true
		case "install_mode":
			text := iterator.ReadString()
			value := AddOnInstallMode(text)
			object.installMode = value
			object.fieldSet_[13] = true
		case "label":
			value := iterator.ReadString()
			object.label = value
			object.fieldSet_[14] = true
		case "managed_service":
			value := iterator.ReadBool()
			object.managedService = value
			object.fieldSet_[15] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[16] = true
		case "namespaces":
			value := ReadAddOnNamespaceList(iterator)
			object.namespaces = value
			object.fieldSet_[17] = true
		case "operator_name":
			value := iterator.ReadString()
			object.operatorName = value
			object.fieldSet_[18] = true
		case "parameters":
			value := &AddOnParameterList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == AddOnParameterListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadAddOnParameterList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.parameters = value
			object.fieldSet_[19] = true
		case "requirements":
			value := ReadAddOnRequirementList(iterator)
			object.requirements = value
			object.fieldSet_[20] = true
		case "resource_cost":
			value := iterator.ReadFloat64()
			object.resourceCost = value
			object.fieldSet_[21] = true
		case "resource_name":
			value := iterator.ReadString()
			object.resourceName = value
			object.fieldSet_[22] = true
		case "sub_operators":
			value := ReadAddOnSubOperatorList(iterator)
			object.subOperators = value
			object.fieldSet_[23] = true
		case "target_namespace":
			value := iterator.ReadString()
			object.targetNamespace = value
			object.fieldSet_[24] = true
		case "version":
			value := ReadAddOnVersion(iterator)
			object.version = value
			object.fieldSet_[25] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
