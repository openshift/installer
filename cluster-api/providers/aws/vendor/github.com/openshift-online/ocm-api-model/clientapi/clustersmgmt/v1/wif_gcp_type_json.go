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

// MarshalWifGcp writes a value of the 'wif_gcp' type to the given writer.
func MarshalWifGcp(object *WifGcp, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifGcp(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifGcp writes a value of the 'wif_gcp' type to the given stream.
func WriteWifGcp(object *WifGcp, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("federated_project_id")
		stream.WriteString(object.federatedProjectId)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("federated_project_number")
		stream.WriteString(object.federatedProjectNumber)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("impersonator_email")
		stream.WriteString(object.impersonatorEmail)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("project_id")
		stream.WriteString(object.projectId)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("project_number")
		stream.WriteString(object.projectNumber)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_prefix")
		stream.WriteString(object.rolePrefix)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.serviceAccounts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_accounts")
		WriteWifServiceAccountList(object.serviceAccounts, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.support != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("support")
		WriteWifSupport(object.support, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.workloadIdentityPool != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("workload_identity_pool")
		WriteWifPool(object.workloadIdentityPool, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifGcp reads a value of the 'wif_gcp' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifGcp(source interface{}) (object *WifGcp, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifGcp(iterator)
	err = iterator.Error
	return
}

// ReadWifGcp reads a value of the 'wif_gcp' type from the given iterator.
func ReadWifGcp(iterator *jsoniter.Iterator) *WifGcp {
	object := &WifGcp{
		fieldSet_: make([]bool, 9),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "federated_project_id":
			value := iterator.ReadString()
			object.federatedProjectId = value
			object.fieldSet_[0] = true
		case "federated_project_number":
			value := iterator.ReadString()
			object.federatedProjectNumber = value
			object.fieldSet_[1] = true
		case "impersonator_email":
			value := iterator.ReadString()
			object.impersonatorEmail = value
			object.fieldSet_[2] = true
		case "project_id":
			value := iterator.ReadString()
			object.projectId = value
			object.fieldSet_[3] = true
		case "project_number":
			value := iterator.ReadString()
			object.projectNumber = value
			object.fieldSet_[4] = true
		case "role_prefix":
			value := iterator.ReadString()
			object.rolePrefix = value
			object.fieldSet_[5] = true
		case "service_accounts":
			value := ReadWifServiceAccountList(iterator)
			object.serviceAccounts = value
			object.fieldSet_[6] = true
		case "support":
			value := ReadWifSupport(iterator)
			object.support = value
			object.fieldSet_[7] = true
		case "workload_identity_pool":
			value := ReadWifPool(iterator)
			object.workloadIdentityPool = value
			object.fieldSet_[8] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
