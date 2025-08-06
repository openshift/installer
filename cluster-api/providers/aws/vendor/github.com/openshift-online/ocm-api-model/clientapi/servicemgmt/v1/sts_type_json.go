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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalSTS writes a value of the 'STS' type to the given writer.
func MarshalSTS(object *STS, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSTS(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSTS writes a value of the 'STS' type to the given stream.
func WriteSTS(object *STS, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("oidc_endpoint_url")
		stream.WriteString(object.oidcEndpointURL)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.instanceIAMRoles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("instance_iam_roles")
		WriteInstanceIAMRoles(object.instanceIAMRoles, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.operatorIAMRoles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_iam_roles")
		WriteOperatorIAMRoleList(object.operatorIAMRoles, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_role_prefix")
		stream.WriteString(object.operatorRolePrefix)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_arn")
		stream.WriteString(object.roleARN)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("support_role_arn")
		stream.WriteString(object.supportRoleARN)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSTS reads a value of the 'STS' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSTS(source interface{}) (object *STS, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadSTS(iterator)
	err = iterator.Error
	return
}

// ReadSTS reads a value of the 'STS' type from the given iterator.
func ReadSTS(iterator *jsoniter.Iterator) *STS {
	object := &STS{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "oidc_endpoint_url":
			value := iterator.ReadString()
			object.oidcEndpointURL = value
			object.fieldSet_[0] = true
		case "instance_iam_roles":
			value := ReadInstanceIAMRoles(iterator)
			object.instanceIAMRoles = value
			object.fieldSet_[1] = true
		case "operator_iam_roles":
			value := ReadOperatorIAMRoleList(iterator)
			object.operatorIAMRoles = value
			object.fieldSet_[2] = true
		case "operator_role_prefix":
			value := iterator.ReadString()
			object.operatorRolePrefix = value
			object.fieldSet_[3] = true
		case "role_arn":
			value := iterator.ReadString()
			object.roleARN = value
			object.fieldSet_[4] = true
		case "support_role_arn":
			value := iterator.ReadString()
			object.supportRoleARN = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
