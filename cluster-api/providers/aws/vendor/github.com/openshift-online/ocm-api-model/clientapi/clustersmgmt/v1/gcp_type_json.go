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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalGCP writes a value of the 'GCP' type to the given writer.
func MarshalGCP(object *GCP, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGCP(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGCP writes a value of the 'GCP' type to the given stream.
func WriteGCP(object *GCP, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auth_uri")
		stream.WriteString(object.authURI)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auth_provider_x509_cert_url")
		stream.WriteString(object.authProviderX509CertURL)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.authentication != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("authentication")
		WriteGcpAuthentication(object.authentication, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_id")
		stream.WriteString(object.clientID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_x509_cert_url")
		stream.WriteString(object.clientX509CertURL)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("client_email")
		stream.WriteString(object.clientEmail)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_key")
		stream.WriteString(object.privateKey)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_key_id")
		stream.WriteString(object.privateKeyID)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.privateServiceConnect != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_service_connect")
		WriteGcpPrivateServiceConnect(object.privateServiceConnect, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("project_id")
		stream.WriteString(object.projectID)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.security != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("security")
		WriteGcpSecurity(object.security, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("token_uri")
		stream.WriteString(object.tokenURI)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGCP reads a value of the 'GCP' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGCP(source interface{}) (object *GCP, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGCP(iterator)
	err = iterator.Error
	return
}

// ReadGCP reads a value of the 'GCP' type from the given iterator.
func ReadGCP(iterator *jsoniter.Iterator) *GCP {
	object := &GCP{
		fieldSet_: make([]bool, 13),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "auth_uri":
			value := iterator.ReadString()
			object.authURI = value
			object.fieldSet_[0] = true
		case "auth_provider_x509_cert_url":
			value := iterator.ReadString()
			object.authProviderX509CertURL = value
			object.fieldSet_[1] = true
		case "authentication":
			value := ReadGcpAuthentication(iterator)
			object.authentication = value
			object.fieldSet_[2] = true
		case "client_id":
			value := iterator.ReadString()
			object.clientID = value
			object.fieldSet_[3] = true
		case "client_x509_cert_url":
			value := iterator.ReadString()
			object.clientX509CertURL = value
			object.fieldSet_[4] = true
		case "client_email":
			value := iterator.ReadString()
			object.clientEmail = value
			object.fieldSet_[5] = true
		case "private_key":
			value := iterator.ReadString()
			object.privateKey = value
			object.fieldSet_[6] = true
		case "private_key_id":
			value := iterator.ReadString()
			object.privateKeyID = value
			object.fieldSet_[7] = true
		case "private_service_connect":
			value := ReadGcpPrivateServiceConnect(iterator)
			object.privateServiceConnect = value
			object.fieldSet_[8] = true
		case "project_id":
			value := iterator.ReadString()
			object.projectID = value
			object.fieldSet_[9] = true
		case "security":
			value := ReadGcpSecurity(iterator)
			object.security = value
			object.fieldSet_[10] = true
		case "token_uri":
			value := iterator.ReadString()
			object.tokenURI = value
			object.fieldSet_[11] = true
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.fieldSet_[12] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
