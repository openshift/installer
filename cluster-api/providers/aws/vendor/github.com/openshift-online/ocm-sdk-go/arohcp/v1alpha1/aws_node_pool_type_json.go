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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"io"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAWSNodePool writes a value of the 'AWS_node_pool' type to the given writer.
func MarshalAWSNodePool(object *AWSNodePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAWSNodePool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAWSNodePool writes a value of the 'AWS_node_pool' type to the given stream.
func WriteAWSNodePool(object *AWSNodePool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AWSNodePoolLinkKind)
	} else {
		stream.WriteString(AWSNodePoolKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0 && object.additionalSecurityGroupIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_security_group_ids")
		WriteStringList(object.additionalSecurityGroupIds, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.availabilityZoneTypes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone_types")
		if object.availabilityZoneTypes != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.availabilityZoneTypes))
			i := 0
			for key := range object.availabilityZoneTypes {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.availabilityZoneTypes[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ec2_metadata_http_tokens")
		stream.WriteString(string(object.ec2MetadataHttpTokens))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("instance_profile")
		stream.WriteString(object.instanceProfile)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("instance_type")
		stream.WriteString(object.instanceType)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.rootVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("root_volume")
		WriteAWSVolume(object.rootVolume, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.subnetOutposts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet_outposts")
		if object.subnetOutposts != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.subnetOutposts))
			i := 0
			for key := range object.subnetOutposts {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.subnetOutposts[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.tags != nil
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

// UnmarshalAWSNodePool reads a value of the 'AWS_node_pool' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSNodePool(source interface{}) (object *AWSNodePool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAWSNodePool(iterator)
	err = iterator.Error
	return
}

// ReadAWSNodePool reads a value of the 'AWS_node_pool' type from the given iterator.
func ReadAWSNodePool(iterator *jsoniter.Iterator) *AWSNodePool {
	object := &AWSNodePool{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AWSNodePoolLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "additional_security_group_ids":
			value := ReadStringList(iterator)
			object.additionalSecurityGroupIds = value
			object.bitmap_ |= 8
		case "availability_zone_types":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.availabilityZoneTypes = value
			object.bitmap_ |= 16
		case "ec2_metadata_http_tokens":
			text := iterator.ReadString()
			value := Ec2MetadataHttpTokens(text)
			object.ec2MetadataHttpTokens = value
			object.bitmap_ |= 32
		case "instance_profile":
			value := iterator.ReadString()
			object.instanceProfile = value
			object.bitmap_ |= 64
		case "instance_type":
			value := iterator.ReadString()
			object.instanceType = value
			object.bitmap_ |= 128
		case "root_volume":
			value := ReadAWSVolume(iterator)
			object.rootVolume = value
			object.bitmap_ |= 256
		case "subnet_outposts":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.subnetOutposts = value
			object.bitmap_ |= 512
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
			object.bitmap_ |= 1024
		default:
			iterator.ReadAny()
		}
	}
	return object
}
