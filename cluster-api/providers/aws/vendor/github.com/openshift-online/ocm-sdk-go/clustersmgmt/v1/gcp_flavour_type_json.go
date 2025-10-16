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

// MarshalGCPFlavour writes a value of the 'GCP_flavour' type to the given writer.
func MarshalGCPFlavour(object *GCPFlavour, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGCPFlavour(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGCPFlavour writes a value of the 'GCP_flavour' type to the given stream.
func WriteGCPFlavour(object *GCPFlavour, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_instance_type")
		stream.WriteString(object.computeInstanceType)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra_instance_type")
		stream.WriteString(object.infraInstanceType)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.infraVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra_volume")
		WriteGCPVolume(object.infraVolume, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master_instance_type")
		stream.WriteString(object.masterInstanceType)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.masterVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master_volume")
		WriteGCPVolume(object.masterVolume, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.workerVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("worker_volume")
		WriteGCPVolume(object.workerVolume, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGCPFlavour reads a value of the 'GCP_flavour' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGCPFlavour(source interface{}) (object *GCPFlavour, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGCPFlavour(iterator)
	err = iterator.Error
	return
}

// ReadGCPFlavour reads a value of the 'GCP_flavour' type from the given iterator.
func ReadGCPFlavour(iterator *jsoniter.Iterator) *GCPFlavour {
	object := &GCPFlavour{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "compute_instance_type":
			value := iterator.ReadString()
			object.computeInstanceType = value
			object.bitmap_ |= 1
		case "infra_instance_type":
			value := iterator.ReadString()
			object.infraInstanceType = value
			object.bitmap_ |= 2
		case "infra_volume":
			value := ReadGCPVolume(iterator)
			object.infraVolume = value
			object.bitmap_ |= 4
		case "master_instance_type":
			value := iterator.ReadString()
			object.masterInstanceType = value
			object.bitmap_ |= 8
		case "master_volume":
			value := ReadGCPVolume(iterator)
			object.masterVolume = value
			object.bitmap_ |= 16
		case "worker_volume":
			value := ReadGCPVolume(iterator)
			object.workerVolume = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
