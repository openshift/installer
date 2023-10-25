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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalRegistryCredential writes a value of the 'registry_credential' type to the given writer.
func MarshalRegistryCredential(object *RegistryCredential, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeRegistryCredential(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeRegistryCredential writes a value of the 'registry_credential' type to the given stream.
func writeRegistryCredential(object *RegistryCredential, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(RegistryCredentialLinkKind)
	} else {
		stream.WriteString(RegistryCredentialKind)
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
	present_ = object.bitmap_&8 != 0 && object.account != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account")
		writeAccount(object.account, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("created_at")
		stream.WriteString((object.createdAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_resource_id")
		stream.WriteString(object.externalResourceID)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.registry != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("registry")
		writeRegistry(object.registry, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("token")
		stream.WriteString(object.token)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_at")
		stream.WriteString((object.updatedAt).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("username")
		stream.WriteString(object.username)
	}
	stream.WriteObjectEnd()
}

// UnmarshalRegistryCredential reads a value of the 'registry_credential' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRegistryCredential(source interface{}) (object *RegistryCredential, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readRegistryCredential(iterator)
	err = iterator.Error
	return
}

// readRegistryCredential reads a value of the 'registry_credential' type from the given iterator.
func readRegistryCredential(iterator *jsoniter.Iterator) *RegistryCredential {
	object := &RegistryCredential{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == RegistryCredentialLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "account":
			value := readAccount(iterator)
			object.account = value
			object.bitmap_ |= 8
		case "created_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.createdAt = value
			object.bitmap_ |= 16
		case "external_resource_id":
			value := iterator.ReadString()
			object.externalResourceID = value
			object.bitmap_ |= 32
		case "registry":
			value := readRegistry(iterator)
			object.registry = value
			object.bitmap_ |= 64
		case "token":
			value := iterator.ReadString()
			object.token = value
			object.bitmap_ |= 128
		case "updated_at":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedAt = value
			object.bitmap_ |= 256
		case "username":
			value := iterator.ReadString()
			object.username = value
			object.bitmap_ |= 512
		default:
			iterator.ReadAny()
		}
	}
	return object
}
