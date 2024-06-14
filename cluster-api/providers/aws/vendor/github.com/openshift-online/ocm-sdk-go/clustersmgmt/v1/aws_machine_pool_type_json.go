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

// MarshalAWSMachinePool writes a value of the 'AWS_machine_pool' type to the given writer.
func MarshalAWSMachinePool(object *AWSMachinePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAWSMachinePool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAWSMachinePool writes a value of the 'AWS_machine_pool' type to the given stream.
func writeAWSMachinePool(object *AWSMachinePool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(AWSMachinePoolLinkKind)
	} else {
		stream.WriteString(AWSMachinePoolKind)
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
		writeStringList(object.additionalSecurityGroupIds, stream)
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
	present_ = object.bitmap_&32 != 0 && object.spotMarketOptions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("spot_market_options")
		writeAWSSpotMarketOptions(object.spotMarketOptions, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.subnetOutposts != nil
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
	present_ = object.bitmap_&128 != 0 && object.tags != nil
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

// UnmarshalAWSMachinePool reads a value of the 'AWS_machine_pool' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAWSMachinePool(source interface{}) (object *AWSMachinePool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAWSMachinePool(iterator)
	err = iterator.Error
	return
}

// readAWSMachinePool reads a value of the 'AWS_machine_pool' type from the given iterator.
func readAWSMachinePool(iterator *jsoniter.Iterator) *AWSMachinePool {
	object := &AWSMachinePool{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AWSMachinePoolLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "additional_security_group_ids":
			value := readStringList(iterator)
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
		case "spot_market_options":
			value := readAWSSpotMarketOptions(iterator)
			object.spotMarketOptions = value
			object.bitmap_ |= 32
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
			object.bitmap_ |= 64
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
			object.bitmap_ |= 128
		default:
			iterator.ReadAny()
		}
	}
	return object
}
