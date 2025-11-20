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

// MarshalAzure writes a value of the 'azure' type to the given writer.
func MarshalAzure(object *Azure, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAzure(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAzure writes a value of the 'azure' type to the given stream.
func WriteAzure(object *Azure, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.etcdEncryption != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("etcd_encryption")
		WriteAzureEtcdEncryption(object.etcdEncryption, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_resource_group_name")
		stream.WriteString(object.managedResourceGroupName)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("network_security_group_resource_id")
		stream.WriteString(object.networkSecurityGroupResourceID)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.nodesOutboundConnectivity != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("nodes_outbound_connectivity")
		WriteAzureNodesOutboundConnectivity(object.nodesOutboundConnectivity, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.operatorsAuthentication != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operators_authentication")
		WriteAzureOperatorsAuthentication(object.operatorsAuthentication, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_group_name")
		stream.WriteString(object.resourceGroupName)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_name")
		stream.WriteString(object.resourceName)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet_resource_id")
		stream.WriteString(object.subnetResourceID)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription_id")
		stream.WriteString(object.subscriptionID)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("tenant_id")
		stream.WriteString(object.tenantID)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAzure reads a value of the 'azure' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAzure(source interface{}) (object *Azure, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAzure(iterator)
	err = iterator.Error
	return
}

// ReadAzure reads a value of the 'azure' type from the given iterator.
func ReadAzure(iterator *jsoniter.Iterator) *Azure {
	object := &Azure{
		fieldSet_: make([]bool, 10),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "etcd_encryption":
			value := ReadAzureEtcdEncryption(iterator)
			object.etcdEncryption = value
			object.fieldSet_[0] = true
		case "managed_resource_group_name":
			value := iterator.ReadString()
			object.managedResourceGroupName = value
			object.fieldSet_[1] = true
		case "network_security_group_resource_id":
			value := iterator.ReadString()
			object.networkSecurityGroupResourceID = value
			object.fieldSet_[2] = true
		case "nodes_outbound_connectivity":
			value := ReadAzureNodesOutboundConnectivity(iterator)
			object.nodesOutboundConnectivity = value
			object.fieldSet_[3] = true
		case "operators_authentication":
			value := ReadAzureOperatorsAuthentication(iterator)
			object.operatorsAuthentication = value
			object.fieldSet_[4] = true
		case "resource_group_name":
			value := iterator.ReadString()
			object.resourceGroupName = value
			object.fieldSet_[5] = true
		case "resource_name":
			value := iterator.ReadString()
			object.resourceName = value
			object.fieldSet_[6] = true
		case "subnet_resource_id":
			value := iterator.ReadString()
			object.subnetResourceID = value
			object.fieldSet_[7] = true
		case "subscription_id":
			value := iterator.ReadString()
			object.subscriptionID = value
			object.fieldSet_[8] = true
		case "tenant_id":
			value := iterator.ReadString()
			object.tenantID = value
			object.fieldSet_[9] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
