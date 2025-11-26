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

// MarshalAuditLog writes a value of the 'audit_log' type to the given writer.
func MarshalAuditLog(object *AuditLog, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAuditLog(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAuditLog writes a value of the 'audit_log' type to the given stream.
func WriteAuditLog(object *AuditLog, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_arn")
		stream.WriteString(object.roleArn)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAuditLog reads a value of the 'audit_log' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAuditLog(source interface{}) (object *AuditLog, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAuditLog(iterator)
	err = iterator.Error
	return
}

// ReadAuditLog reads a value of the 'audit_log' type from the given iterator.
func ReadAuditLog(iterator *jsoniter.Iterator) *AuditLog {
	object := &AuditLog{
		fieldSet_: make([]bool, 1),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "role_arn":
			value := iterator.ReadString()
			object.roleArn = value
			object.fieldSet_[0] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
