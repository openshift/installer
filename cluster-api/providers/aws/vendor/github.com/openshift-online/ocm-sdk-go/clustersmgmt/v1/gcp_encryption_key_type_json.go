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

// MarshalGCPEncryptionKey writes a value of the 'GCP_encryption_key' type to the given writer.
func MarshalGCPEncryptionKey(object *GCPEncryptionKey, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGCPEncryptionKey(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGCPEncryptionKey writes a value of the 'GCP_encryption_key' type to the given stream.
func WriteGCPEncryptionKey(object *GCPEncryptionKey, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_key_service_account")
		stream.WriteString(object.kmsKeyServiceAccount)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_location")
		stream.WriteString(object.keyLocation)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_name")
		stream.WriteString(object.keyName)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_ring")
		stream.WriteString(object.keyRing)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGCPEncryptionKey reads a value of the 'GCP_encryption_key' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGCPEncryptionKey(source interface{}) (object *GCPEncryptionKey, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGCPEncryptionKey(iterator)
	err = iterator.Error
	return
}

// ReadGCPEncryptionKey reads a value of the 'GCP_encryption_key' type from the given iterator.
func ReadGCPEncryptionKey(iterator *jsoniter.Iterator) *GCPEncryptionKey {
	object := &GCPEncryptionKey{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kms_key_service_account":
			value := iterator.ReadString()
			object.kmsKeyServiceAccount = value
			object.bitmap_ |= 1
		case "key_location":
			value := iterator.ReadString()
			object.keyLocation = value
			object.bitmap_ |= 2
		case "key_name":
			value := iterator.ReadString()
			object.keyName = value
			object.bitmap_ |= 4
		case "key_ring":
			value := iterator.ReadString()
			object.keyRing = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
