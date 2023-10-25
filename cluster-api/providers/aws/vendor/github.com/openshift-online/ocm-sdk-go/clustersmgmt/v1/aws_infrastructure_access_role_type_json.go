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

// MarshalAWSInfrastructureAccessRole writes a value of the 'AWS_infrastructure_access_role' type to the given writer.
func MarshalAWSInfrastructureAccessRole(object *AWSInfrastructureAccessRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAWSInfrastructureAccessRole(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAWSInfrastructureAccessRole writes a value of the 'AWS_infrastructure_access_role' type to the given stream.
func writeAWSInfrastructureAccessRole(object *AWSInfrastructureAccessRole, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AWSInfrastructureAccessRoleLinkKind)
	} else {
		stream.WriteString(AWSInfrastructureAccessRoleKind)
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
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWSInfrastructureAccessRole reads a value of the 'AWS_infrastructure_access_role' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSInfrastructureAccessRole(source interface{}) (object *AWSInfrastructureAccessRole, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAWSInfrastructureAccessRole(iterator)
	err = iterator.Error
	return
}

// readAWSInfrastructureAccessRole reads a value of the 'AWS_infrastructure_access_role' type from the given iterator.
func readAWSInfrastructureAccessRole(iterator *jsoniter.Iterator) *AWSInfrastructureAccessRole {
	object := &AWSInfrastructureAccessRole{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AWSInfrastructureAccessRoleLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.bitmap_ |= 8
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.bitmap_ |= 16
		case "state":
			text := iterator.ReadString()
			value := AWSInfrastructureAccessRoleState(text)
			object.state = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
