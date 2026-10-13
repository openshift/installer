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

// MarshalOidcThumbprint writes a value of the 'oidc_thumbprint' type to the given writer.
func MarshalOidcThumbprint(object *OidcThumbprint, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteOidcThumbprint(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteOidcThumbprint writes a value of the 'oidc_thumbprint' type to the given stream.
func WriteOidcThumbprint(object *OidcThumbprint, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterId)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kind")
		stream.WriteString(object.kind)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("oidc_config_id")
		stream.WriteString(object.oidcConfigId)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("thumbprint")
		stream.WriteString(object.thumbprint)
	}
	stream.WriteObjectEnd()
}

// UnmarshalOidcThumbprint reads a value of the 'oidc_thumbprint' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalOidcThumbprint(source interface{}) (object *OidcThumbprint, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadOidcThumbprint(iterator)
	err = iterator.Error
	return
}

// ReadOidcThumbprint reads a value of the 'oidc_thumbprint' type from the given iterator.
func ReadOidcThumbprint(iterator *jsoniter.Iterator) *OidcThumbprint {
	object := &OidcThumbprint{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "href":
			value := iterator.ReadString()
			object.href = value
			object.fieldSet_[0] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterId = value
			object.fieldSet_[1] = true
		case "kind":
			value := iterator.ReadString()
			object.kind = value
			object.fieldSet_[2] = true
		case "oidc_config_id":
			value := iterator.ReadString()
			object.oidcConfigId = value
			object.fieldSet_[3] = true
		case "thumbprint":
			value := iterator.ReadString()
			object.thumbprint = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
