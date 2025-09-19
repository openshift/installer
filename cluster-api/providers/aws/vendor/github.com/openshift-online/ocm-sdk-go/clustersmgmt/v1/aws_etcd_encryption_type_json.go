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

// MarshalAwsEtcdEncryption writes a value of the 'aws_etcd_encryption' type to the given writer.
func MarshalAwsEtcdEncryption(object *AwsEtcdEncryption, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAwsEtcdEncryption(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAwsEtcdEncryption writes a value of the 'aws_etcd_encryption' type to the given stream.
func WriteAwsEtcdEncryption(object *AwsEtcdEncryption, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_key_arn")
		stream.WriteString(object.kmsKeyARN)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAwsEtcdEncryption reads a value of the 'aws_etcd_encryption' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAwsEtcdEncryption(source interface{}) (object *AwsEtcdEncryption, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAwsEtcdEncryption(iterator)
	err = iterator.Error
	return
}

// ReadAwsEtcdEncryption reads a value of the 'aws_etcd_encryption' type from the given iterator.
func ReadAwsEtcdEncryption(iterator *jsoniter.Iterator) *AwsEtcdEncryption {
	object := &AwsEtcdEncryption{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kms_key_arn":
			value := iterator.ReadString()
			object.kmsKeyARN = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
