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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAzureOperatorsAuthenticationManagedIdentities writes a value of the 'azure_operators_authentication_managed_identities' type to the given writer.
func MarshalAzureOperatorsAuthenticationManagedIdentities(object *AzureOperatorsAuthenticationManagedIdentities, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzureOperatorsAuthenticationManagedIdentities(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzureOperatorsAuthenticationManagedIdentities writes a value of the 'azure_operators_authentication_managed_identities' type to the given stream.
func WriteAzureOperatorsAuthenticationManagedIdentities(object *AzureOperatorsAuthenticationManagedIdentities, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.controlPlaneOperatorsManagedIdentities != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("control_plane_operators_managed_identities")
		if object.controlPlaneOperatorsManagedIdentities != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.controlPlaneOperatorsManagedIdentities))
			i := 0
			for key := range object.controlPlaneOperatorsManagedIdentities {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.controlPlaneOperatorsManagedIdentities[key]
				stream.WriteObjectField(key)
				WriteAzureControlPlaneManagedIdentity(item, stream)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.dataPlaneOperatorsManagedIdentities != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("data_plane_operators_managed_identities")
		if object.dataPlaneOperatorsManagedIdentities != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.dataPlaneOperatorsManagedIdentities))
			i := 0
			for key := range object.dataPlaneOperatorsManagedIdentities {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.dataPlaneOperatorsManagedIdentities[key]
				stream.WriteObjectField(key)
				WriteAzureDataPlaneManagedIdentity(item, stream)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_identities_data_plane_identity_url")
		stream.WriteString(object.managedIdentitiesDataPlaneIdentityUrl)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.serviceManagedIdentity != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_managed_identity")
		WriteAzureServiceManagedIdentity(object.serviceManagedIdentity, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzureOperatorsAuthenticationManagedIdentities reads a value of the 'azure_operators_authentication_managed_identities' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzureOperatorsAuthenticationManagedIdentities(source interface{}) (object *AzureOperatorsAuthenticationManagedIdentities, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzureOperatorsAuthenticationManagedIdentities(iterator)
	err = iterator.Error
	return
}

// ReadAzureOperatorsAuthenticationManagedIdentities reads a value of the 'azure_operators_authentication_managed_identities' type from the given iterator.
func ReadAzureOperatorsAuthenticationManagedIdentities(iterator *jsoniter.Iterator) *AzureOperatorsAuthenticationManagedIdentities {
	object := &AzureOperatorsAuthenticationManagedIdentities{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "control_plane_operators_managed_identities":
			value := map[string]*AzureControlPlaneManagedIdentity{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := ReadAzureControlPlaneManagedIdentity(iterator)
				value[key] = item
			}
			object.controlPlaneOperatorsManagedIdentities = value
			object.fieldSet_[0] = true
		case "data_plane_operators_managed_identities":
			value := map[string]*AzureDataPlaneManagedIdentity{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := ReadAzureDataPlaneManagedIdentity(iterator)
				value[key] = item
			}
			object.dataPlaneOperatorsManagedIdentities = value
			object.fieldSet_[1] = true
		case "managed_identities_data_plane_identity_url":
			value := iterator.ReadString()
			object.managedIdentitiesDataPlaneIdentityUrl = value
			object.fieldSet_[2] = true
		case "service_managed_identity":
			value := ReadAzureServiceManagedIdentity(iterator)
			object.serviceManagedIdentity = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
