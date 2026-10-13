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

// MarshalAWSBackupConfig writes a value of the 'AWS_backup_config' type to the given writer.
func MarshalAWSBackupConfig(object *AWSBackupConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWSBackupConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWSBackupConfig writes a value of the 'AWS_backup_config' type to the given stream.
func WriteAWSBackupConfig(object *AWSBackupConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("s3_bucket")
		stream.WriteString(object.s3Bucket)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountId)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("identity_provider_arn")
		stream.WriteString(object.identityProviderArn)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("management_cluster")
		stream.WriteString(object.managementCluster)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_arn")
		stream.WriteString(object.roleArn)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAWSBackupConfig reads a value of the 'AWS_backup_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSBackupConfig(source interface{}) (object *AWSBackupConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAWSBackupConfig(iterator)
	err = iterator.Error
	return
}

// ReadAWSBackupConfig reads a value of the 'AWS_backup_config' type from the given iterator.
func ReadAWSBackupConfig(iterator *jsoniter.Iterator) *AWSBackupConfig {
	object := &AWSBackupConfig{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "s3_bucket":
			value := iterator.ReadString()
			object.s3Bucket = value
			object.fieldSet_[0] = true
		case "account_id":
			value := iterator.ReadString()
			object.accountId = value
			object.fieldSet_[1] = true
		case "identity_provider_arn":
			value := iterator.ReadString()
			object.identityProviderArn = value
			object.fieldSet_[2] = true
		case "management_cluster":
			value := iterator.ReadString()
			object.managementCluster = value
			object.fieldSet_[3] = true
		case "role_arn":
			value := iterator.ReadString()
			object.roleArn = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
