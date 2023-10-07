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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

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
	present_ = object.bitmap_&1 != 0 && object.sts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sts")
		writeSTS(object.sts, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("access_key_id")
		stream.WriteString(object.accessKeyID)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("account_id")
		stream.WriteString(object.accountID)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("private_link")
		stream.WriteBool(object.privateLink)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secret_access_key")
		stream.WriteString(object.secretAccessKey)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.subnetIDs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet_ids")
		writeStringList(object.subnetIDs, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.tags != nil
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
		case "sts":
			value := readSTS(iterator)
			object.sts = value
			object.bitmap_ |= 1
		case "access_key_id":
			value := iterator.ReadString()
			object.accessKeyID = value
			object.bitmap_ |= 2
		case "account_id":
			value := iterator.ReadString()
			object.accountID = value
			object.bitmap_ |= 4
		case "private_link":
			value := iterator.ReadBool()
			object.privateLink = value
			object.bitmap_ |= 8
		case "secret_access_key":
			value := iterator.ReadString()
			object.secretAccessKey = value
			object.bitmap_ |= 16
		case "subnet_ids":
			value := readStringList(iterator)
			object.subnetIDs = value
			object.bitmap_ |= 32
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
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
