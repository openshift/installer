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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalNetwork writes a value of the 'network' type to the given writer.
func MarshalNetwork(object *Network, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNetwork(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNetwork writes a value of the 'network' type to the given stream.
func WriteNetwork(object *Network, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("host_prefix")
		stream.WriteInt(object.hostPrefix)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("machine_cidr")
		stream.WriteString(object.machineCIDR)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("pod_cidr")
		stream.WriteString(object.podCIDR)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_cidr")
		stream.WriteString(object.serviceCIDR)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
	}
	stream.WriteObjectEnd()
}

// UnmarshalNetwork reads a value of the 'network' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalNetwork(source interface{}) (object *Network, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadNetwork(iterator)
	err = iterator.Error
	return
}

// ReadNetwork reads a value of the 'network' type from the given iterator.
func ReadNetwork(iterator *jsoniter.Iterator) *Network {
	object := &Network{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "host_prefix":
			value := iterator.ReadInt()
			object.hostPrefix = value
			object.fieldSet_[0] = true
		case "machine_cidr":
			value := iterator.ReadString()
			object.machineCIDR = value
			object.fieldSet_[1] = true
		case "pod_cidr":
			value := iterator.ReadString()
			object.podCIDR = value
			object.fieldSet_[2] = true
		case "service_cidr":
			value := iterator.ReadString()
			object.serviceCIDR = value
			object.fieldSet_[3] = true
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
