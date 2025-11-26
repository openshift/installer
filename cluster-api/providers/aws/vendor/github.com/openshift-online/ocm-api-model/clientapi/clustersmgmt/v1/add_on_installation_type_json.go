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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAddOnInstallation writes a value of the 'add_on_installation' type to the given writer.
func MarshalAddOnInstallation(object *AddOnInstallation, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddOnInstallation(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddOnInstallation writes a value of the 'add_on_installation' type to the given stream.
func WriteAddOnInstallation(object *AddOnInstallation, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AddOnInstallationLinkKind)
	} else {
		stream.WriteString(AddOnInstallationKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.addon != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addon")
		WriteAddOn(object.addon, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.addonVersion != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addon_version")
		WriteAddOnVersion(object.addonVersion, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.billing != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing")
		WriteAddOnInstallationBilling(object.billing, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_version")
		stream.WriteString(object.operatorVersion)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.parameters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("parameters")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteAddOnInstallationParameterList(object.parameters.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state_description")
		stream.WriteString(object.stateDescription)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_timestamp")
		stream.WriteString((object.updatedTimestamp).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnInstallation reads a value of the 'add_on_installation' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnInstallation(source interface{}) (object *AddOnInstallation, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddOnInstallation(iterator)
	err = iterator.Error
	return
}

// ReadAddOnInstallation reads a value of the 'add_on_installation' type from the given iterator.
func ReadAddOnInstallation(iterator *jsoniter.Iterator) *AddOnInstallation {
	object := &AddOnInstallation{
		fieldSet_: make([]bool, 12),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddOnInstallationLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "addon":
			value := ReadAddOn(iterator)
			object.addon = value
			object.fieldSet_[3] = true
		case "addon_version":
			value := ReadAddOnVersion(iterator)
			object.addonVersion = value
			object.fieldSet_[4] = true
		case "billing":
			value := ReadAddOnInstallationBilling(iterator)
			object.billing = value
			object.fieldSet_[5] = true
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.fieldSet_[6] = true
		case "operator_version":
			value := iterator.ReadString()
			object.operatorVersion = value
			object.fieldSet_[7] = true
		case "parameters":
			value := &AddOnInstallationParameterList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == AddOnInstallationParameterListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadAddOnInstallationParameterList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.parameters = value
			object.fieldSet_[8] = true
		case "state":
			text := iterator.ReadString()
			value := AddOnInstallationState(text)
			object.state = value
			object.fieldSet_[9] = true
		case "state_description":
			value := iterator.ReadString()
			object.stateDescription = value
			object.fieldSet_[10] = true
		case "updated_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedTimestamp = value
			object.fieldSet_[11] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
