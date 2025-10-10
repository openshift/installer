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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
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
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("oidc_endpoint_url")
		stream.WriteString(object.oidcEndpointURL)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auto_mode")
		stream.WriteBool(object.autoMode)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_id")
		stream.WriteString(object.externalID)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.instanceIAMRoles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("instance_iam_roles")
		WriteInstanceIAMRoles(object.instanceIAMRoles, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_policies")
		stream.WriteBool(object.managedPolicies)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.oidcConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("oidc_config")
		WriteOidcConfig(object.oidcConfig, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.operatorIAMRoles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_iam_roles")
		WriteOperatorIAMRoleList(object.operatorIAMRoles, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_role_prefix")
		stream.WriteString(object.operatorRolePrefix)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("permission_boundary")
		stream.WriteString(object.permissionBoundary)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_arn")
		stream.WriteString(object.roleARN)
		count++
	}
	present_ = object.bitmap_&2048 != 0
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
	object := &STS{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "oidc_endpoint_url":
			value := iterator.ReadString()
			object.oidcEndpointURL = value
			object.bitmap_ |= 1
		case "auto_mode":
			value := iterator.ReadBool()
			object.autoMode = value
			object.bitmap_ |= 2
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 4
		case "external_id":
			value := iterator.ReadString()
			object.externalID = value
			object.bitmap_ |= 8
		case "instance_iam_roles":
			value := ReadInstanceIAMRoles(iterator)
			object.instanceIAMRoles = value
			object.bitmap_ |= 16
		case "managed_policies":
			value := iterator.ReadBool()
			object.managedPolicies = value
			object.bitmap_ |= 32
		case "oidc_config":
			value := ReadOidcConfig(iterator)
			object.oidcConfig = value
			object.bitmap_ |= 64
		case "operator_iam_roles":
			value := ReadOperatorIAMRoleList(iterator)
			object.operatorIAMRoles = value
			object.bitmap_ |= 128
		case "operator_role_prefix":
			value := iterator.ReadString()
			object.operatorRolePrefix = value
			object.bitmap_ |= 256
		case "permission_boundary":
			value := iterator.ReadString()
			object.permissionBoundary = value
			object.bitmap_ |= 512
		case "role_arn":
			value := iterator.ReadString()
			object.roleARN = value
			object.bitmap_ |= 1024
		case "support_role_arn":
			value := iterator.ReadString()
			object.supportRoleARN = value
			object.bitmap_ |= 2048
		default:
			iterator.ReadAny()
		}
	}
	return object
}
