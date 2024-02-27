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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAWS writes a value of the 'AWS' type to the given writer.
func MarshalAWS(object *AWS, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAWS(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAWS writes a value of the 'AWS' type to the given stream.
func writeAWS(object *AWS, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_key_arn")
		stream.WriteString(object.kmsKeyArn)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.sts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sts")
		writeSTS(object.sts, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("access_key_id")
		stream.WriteString(object.accessKeyID)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.additionalComputeSecurityGroupIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_compute_security_group_ids")
		writeStringList(object.additionalComputeSecurityGroupIds, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.additionalControlPlaneSecurityGroupIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_control_plane_security_group_ids")
		writeStringList(object.additionalControlPlaneSecurityGroupIds, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.additionalInfraSecurityGroupIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_infra_security_group_ids")
		writeStringList(object.additionalInfraSecurityGroupIds, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.auditLog != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("audit_log")
		writeAuditLog(object.auditLog, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_account_id")
		stream.WriteString(object.billingAccountID)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ec2_metadata_http_tokens")
		stream.WriteString(string(object.ec2MetadataHttpTokens))
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.etcdEncryption != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("etcd_encryption")
		writeAwsEtcdEncryption(object.etcdEncryption, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_hosted_zone_id")
		stream.WriteString(object.privateHostedZoneID)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_hosted_zone_role_arn")
		stream.WriteString(object.privateHostedZoneRoleARN)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_link")
		stream.WriteBool(object.privateLink)
		count++
	}
	present_ = object.bitmap_&16384 != 0 && object.privateLinkConfiguration != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_link_configuration")
		writePrivateLinkClusterConfiguration(object.privateLinkConfiguration, stream)
		count++
	}
	present_ = object.bitmap_&32768 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret_access_key")
		stream.WriteString(object.secretAccessKey)
		count++
	}
	present_ = object.bitmap_&65536 != 0 && object.subnetIDs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet_ids")
		writeStringList(object.subnetIDs, stream)
		count++
	}
	present_ = object.bitmap_&131072 != 0 && object.tags != nil
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
	object = readAWS(iterator)
	err = iterator.Error
	return
}

// readAWS reads a value of the 'AWS' type from the given iterator.
func readAWS(iterator *jsoniter.Iterator) *AWS {
	object := &AWS{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kms_key_arn":
			value := iterator.ReadString()
			object.kmsKeyArn = value
			object.bitmap_ |= 1
		case "sts":
			value := readSTS(iterator)
			object.sts = value
			object.bitmap_ |= 2
		case "access_key_id":
			value := iterator.ReadString()
			object.accessKeyID = value
			object.bitmap_ |= 4
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.bitmap_ |= 8
		case "additional_compute_security_group_ids":
			value := readStringList(iterator)
			object.additionalComputeSecurityGroupIds = value
			object.bitmap_ |= 16
		case "additional_control_plane_security_group_ids":
			value := readStringList(iterator)
			object.additionalControlPlaneSecurityGroupIds = value
			object.bitmap_ |= 32
		case "additional_infra_security_group_ids":
			value := readStringList(iterator)
			object.additionalInfraSecurityGroupIds = value
			object.bitmap_ |= 64
		case "audit_log":
			value := readAuditLog(iterator)
			object.auditLog = value
			object.bitmap_ |= 128
		case "billing_account_id":
			value := iterator.ReadString()
			object.billingAccountID = value
			object.bitmap_ |= 256
		case "ec2_metadata_http_tokens":
			text := iterator.ReadString()
			value := Ec2MetadataHttpTokens(text)
			object.ec2MetadataHttpTokens = value
			object.bitmap_ |= 512
		case "etcd_encryption":
			value := readAwsEtcdEncryption(iterator)
			object.etcdEncryption = value
			object.bitmap_ |= 1024
		case "private_hosted_zone_id":
			value := iterator.ReadString()
			object.privateHostedZoneID = value
			object.bitmap_ |= 2048
		case "private_hosted_zone_role_arn":
			value := iterator.ReadString()
			object.privateHostedZoneRoleARN = value
			object.bitmap_ |= 4096
		case "private_link":
			value := iterator.ReadBool()
			object.privateLink = value
			object.bitmap_ |= 8192
		case "private_link_configuration":
			value := readPrivateLinkClusterConfiguration(iterator)
			object.privateLinkConfiguration = value
			object.bitmap_ |= 16384
		case "secret_access_key":
			value := iterator.ReadString()
			object.secretAccessKey = value
			object.bitmap_ |= 32768
		case "subnet_ids":
			value := readStringList(iterator)
			object.subnetIDs = value
			object.bitmap_ |= 65536
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
			object.bitmap_ |= 131072
		default:
			iterator.ReadAny()
		}
	}
	return object
}
