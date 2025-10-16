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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalOidcConfig writes a value of the 'oidc_config' type to the given writer.
func MarshalOidcConfig(object *OidcConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteOidcConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteOidcConfig writes a value of the 'oidc_config' type to the given stream.
func WriteOidcConfig(object *OidcConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("installer_role_arn")
		stream.WriteString(object.installerRoleArn)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("issuer_url")
		stream.WriteString(object.issuerUrl)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_update_timestamp")
		stream.WriteString((object.lastUpdateTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_used_timestamp")
		stream.WriteString((object.lastUsedTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed")
		stream.WriteBool(object.managed)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("organization_id")
		stream.WriteString(object.organizationId)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reusable")
		stream.WriteBool(object.reusable)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret_arn")
		stream.WriteString(object.secretArn)
	}
	stream.WriteObjectEnd()
}

// UnmarshalOidcConfig reads a value of the 'oidc_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalOidcConfig(source interface{}) (object *OidcConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadOidcConfig(iterator)
	err = iterator.Error
	return
}

// ReadOidcConfig reads a value of the 'oidc_config' type from the given iterator.
func ReadOidcConfig(iterator *jsoniter.Iterator) *OidcConfig {
	object := &OidcConfig{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "href":
			value := iterator.ReadString()
			object.href = value
			object.bitmap_ |= 1
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 2
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 4
		case "installer_role_arn":
			value := iterator.ReadString()
			object.installerRoleArn = value
			object.bitmap_ |= 8
		case "issuer_url":
			value := iterator.ReadString()
			object.issuerUrl = value
			object.bitmap_ |= 16
		case "last_update_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdateTimestamp = value
			object.bitmap_ |= 32
		case "last_used_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUsedTimestamp = value
			object.bitmap_ |= 64
		case "managed":
			value := iterator.ReadBool()
			object.managed = value
			object.bitmap_ |= 128
		case "organization_id":
			value := iterator.ReadString()
			object.organizationId = value
			object.bitmap_ |= 256
		case "reusable":
			value := iterator.ReadBool()
			object.reusable = value
			object.bitmap_ |= 512
		case "secret_arn":
			value := iterator.ReadString()
			object.secretArn = value
			object.bitmap_ |= 1024
		default:
			iterator.ReadAny()
		}
	}
	return object
}
