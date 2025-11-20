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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalK8sServiceAccountOperatorIdentityRequirement writes a value of the 'K8s_service_account_operator_identity_requirement' type to the given writer.
func MarshalK8sServiceAccountOperatorIdentityRequirement(object *K8sServiceAccountOperatorIdentityRequirement, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteK8sServiceAccountOperatorIdentityRequirement(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteK8sServiceAccountOperatorIdentityRequirement writes a value of the 'K8s_service_account_operator_identity_requirement' type to the given stream.
func WriteK8sServiceAccountOperatorIdentityRequirement(object *K8sServiceAccountOperatorIdentityRequirement, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespace")
		stream.WriteString(object.namespace)
	}
	stream.WriteObjectEnd()
}

// UnmarshalK8sServiceAccountOperatorIdentityRequirement reads a value of the 'K8s_service_account_operator_identity_requirement' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalK8sServiceAccountOperatorIdentityRequirement(source interface{}) (object *K8sServiceAccountOperatorIdentityRequirement, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadK8sServiceAccountOperatorIdentityRequirement(iterator)
	err = iterator.Error
	return
}

// ReadK8sServiceAccountOperatorIdentityRequirement reads a value of the 'K8s_service_account_operator_identity_requirement' type from the given iterator.
func ReadK8sServiceAccountOperatorIdentityRequirement(iterator *jsoniter.Iterator) *K8sServiceAccountOperatorIdentityRequirement {
	object := &K8sServiceAccountOperatorIdentityRequirement{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[0] = true
		case "namespace":
			value := iterator.ReadString()
			object.namespace = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
