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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAWS writes a value of the 'AWS' type to the given writer.
func MarshalAWS(object *AWS, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWS(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWS writes a value of the 'AWS' type to the given stream.
func WriteAWS(object *AWS, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_key_arn")
		stream.WriteString(object.kmsKeyArn)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.sts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sts")
		WriteSTS(object.sts, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("access_key_id")
		stream.WriteString(object.accessKeyID)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.additionalAllowedPrincipals != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_allowed_principals")
		WriteStringList(object.additionalAllowedPrincipals, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.additionalComputeSecurityGroupIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_compute_security_group_ids")
		WriteStringList(object.additionalComputeSecurityGroupIds, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.additionalControlPlaneSecurityGroupIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_control_plane_security_group_ids")
		WriteStringList(object.additionalControlPlaneSecurityGroupIds, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.additionalInfraSecurityGroupIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_infra_security_group_ids")
		WriteStringList(object.additionalInfraSecurityGroupIds, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.auditLog != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("audit_log")
		WriteAuditLog(object.auditLog, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.autoNode != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auto_node")
		WriteAwsAutoNode(object.autoNode, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_account_id")
		stream.WriteString(object.billingAccountID)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ec2_metadata_http_tokens")
		stream.WriteString(string(object.ec2MetadataHttpTokens))
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.etcdEncryption != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("etcd_encryption")
		WriteAwsEtcdEncryption(object.etcdEncryption, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hcp_internal_communication_hosted_zone_id")
		stream.WriteString(object.hcpInternalCommunicationHostedZoneId)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_hosted_zone_id")
		stream.WriteString(object.privateHostedZoneID)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_hosted_zone_role_arn")
		stream.WriteString(object.privateHostedZoneRoleARN)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_link")
		stream.WriteBool(object.privateLink)
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17] && object.privateLinkConfiguration != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_link_configuration")
		WritePrivateLinkClusterConfiguration(object.privateLinkConfiguration, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 18 && object.fieldSet_[18]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret_access_key")
		stream.WriteString(object.secretAccessKey)
		count++
	}
	present_ = len(object.fieldSet_) > 19 && object.fieldSet_[19] && object.subnetIDs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet_ids")
		WriteStringList(object.subnetIDs, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 20 && object.fieldSet_[20] && object.tags != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("tags")
		if object.tags != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.tags))
			i := 0
			for key := range object.tags {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.tags[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 21 && object.fieldSet_[21]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("vpc_endpoint_role_arn")
		stream.WriteString(object.vpcEndpointRoleArn)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWS reads a value of the 'AWS' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWS(source interface{}) (object *AWS, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAWS(iterator)
	err = iterator.Error
	return
}

// ReadAWS reads a value of the 'AWS' type from the given iterator.
func ReadAWS(iterator *jsoniter.Iterator) *AWS {
	object := &AWS{
		fieldSet_: make([]bool, 22),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kms_key_arn":
			value := iterator.ReadString()
			object.kmsKeyArn = value
			object.fieldSet_[0] = true
		case "sts":
			value := ReadSTS(iterator)
			object.sts = value
			object.fieldSet_[1] = true
		case "access_key_id":
			value := iterator.ReadString()
			object.accessKeyID = value
			object.fieldSet_[2] = true
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.fieldSet_[3] = true
		case "additional_allowed_principals":
			value := ReadStringList(iterator)
			object.additionalAllowedPrincipals = value
			object.fieldSet_[4] = true
		case "additional_compute_security_group_ids":
			value := ReadStringList(iterator)
			object.additionalComputeSecurityGroupIds = value
			object.fieldSet_[5] = true
		case "additional_control_plane_security_group_ids":
			value := ReadStringList(iterator)
			object.additionalControlPlaneSecurityGroupIds = value
			object.fieldSet_[6] = true
		case "additional_infra_security_group_ids":
			value := ReadStringList(iterator)
			object.additionalInfraSecurityGroupIds = value
			object.fieldSet_[7] = true
		case "audit_log":
			value := ReadAuditLog(iterator)
			object.auditLog = value
			object.fieldSet_[8] = true
		case "auto_node":
			value := ReadAwsAutoNode(iterator)
			object.autoNode = value
			object.fieldSet_[9] = true
		case "billing_account_id":
			value := iterator.ReadString()
			object.billingAccountID = value
			object.fieldSet_[10] = true
		case "ec2_metadata_http_tokens":
			text := iterator.ReadString()
			value := Ec2MetadataHttpTokens(text)
			object.ec2MetadataHttpTokens = value
			object.fieldSet_[11] = true
		case "etcd_encryption":
			value := ReadAwsEtcdEncryption(iterator)
			object.etcdEncryption = value
			object.fieldSet_[12] = true
		case "hcp_internal_communication_hosted_zone_id":
			value := iterator.ReadString()
			object.hcpInternalCommunicationHostedZoneId = value
			object.fieldSet_[13] = true
		case "private_hosted_zone_id":
			value := iterator.ReadString()
			object.privateHostedZoneID = value
			object.fieldSet_[14] = true
		case "private_hosted_zone_role_arn":
			value := iterator.ReadString()
			object.privateHostedZoneRoleARN = value
			object.fieldSet_[15] = true
		case "private_link":
			value := iterator.ReadBool()
			object.privateLink = value
			object.fieldSet_[16] = true
		case "private_link_configuration":
			value := ReadPrivateLinkClusterConfiguration(iterator)
			object.privateLinkConfiguration = value
			object.fieldSet_[17] = true
		case "secret_access_key":
			value := iterator.ReadString()
			object.secretAccessKey = value
			object.fieldSet_[18] = true
		case "subnet_ids":
			value := ReadStringList(iterator)
			object.subnetIDs = value
			object.fieldSet_[19] = true
		case "tags":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.tags = value
			object.fieldSet_[20] = true
		case "vpc_endpoint_role_arn":
			value := iterator.ReadString()
			object.vpcEndpointRoleArn = value
			object.fieldSet_[21] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
