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
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_instance_type")
		stream.WriteString(object.computeInstanceType)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra_instance_type")
		stream.WriteString(object.infraInstanceType)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.infraVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra_volume")
		WriteGCPVolume(object.infraVolume, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master_instance_type")
		stream.WriteString(object.masterInstanceType)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.masterVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master_volume")
		WriteGCPVolume(object.masterVolume, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.workerVolume != nil
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
	object := &GCPFlavour{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "compute_instance_type":
			value := iterator.ReadString()
			object.computeInstanceType = value
			object.fieldSet_[0] = true
		case "infra_instance_type":
			value := iterator.ReadString()
			object.infraInstanceType = value
			object.fieldSet_[1] = true
		case "infra_volume":
			value := ReadGCPVolume(iterator)
			object.infraVolume = value
			object.fieldSet_[2] = true
		case "master_instance_type":
			value := iterator.ReadString()
			object.masterInstanceType = value
			object.fieldSet_[3] = true
		case "master_volume":
			value := ReadGCPVolume(iterator)
			object.masterVolume = value
			object.fieldSet_[4] = true
		case "worker_volume":
			value := ReadGCPVolume(iterator)
			object.workerVolume = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
