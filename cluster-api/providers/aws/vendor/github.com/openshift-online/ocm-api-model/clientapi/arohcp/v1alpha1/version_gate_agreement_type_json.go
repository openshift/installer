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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalVersionGateAgreement writes a value of the 'version_gate_agreement' type to the given writer.
func MarshalVersionGateAgreement(object *VersionGateAgreement, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteVersionGateAgreement(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteVersionGateAgreement writes a value of the 'version_gate_agreement' type to the given stream.
func WriteVersionGateAgreement(object *VersionGateAgreement, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(VersionGateAgreementLinkKind)
	} else {
		stream.WriteString(VersionGateAgreementKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("agreed_timestamp")
		stream.WriteString((object.agreedTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.versionGate != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version_gate")
		WriteVersionGate(object.versionGate, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalVersionGateAgreement reads a value of the 'version_gate_agreement' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalVersionGateAgreement(source interface{}) (object *VersionGateAgreement, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadVersionGateAgreement(iterator)
	err = iterator.Error
	return
}

// ReadVersionGateAgreement reads a value of the 'version_gate_agreement' type from the given iterator.
func ReadVersionGateAgreement(iterator *jsoniter.Iterator) *VersionGateAgreement {
	object := &VersionGateAgreement{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == VersionGateAgreementLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "agreed_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.agreedTimestamp = value
			object.fieldSet_[3] = true
		case "version_gate":
			value := ReadVersionGate(iterator)
			object.versionGate = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
