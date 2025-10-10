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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalUpgradePolicy writes a value of the 'upgrade_policy' type to the given writer.
func MarshalUpgradePolicy(object *UpgradePolicy, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteUpgradePolicy(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteUpgradePolicy writes a value of the 'upgrade_policy' type to the given stream.
func WriteUpgradePolicy(object *UpgradePolicy, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(UpgradePolicyLinkKind)
	} else {
		stream.WriteString(UpgradePolicyKind)
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
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enable_minor_version_upgrades")
		stream.WriteBool(object.enableMinorVersionUpgrades)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("next_run")
		stream.WriteString((object.nextRun).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("schedule")
		stream.WriteString(object.schedule)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("schedule_type")
		stream.WriteString(string(object.scheduleType))
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("upgrade_type")
		stream.WriteString(string(object.upgradeType))
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		stream.WriteString(object.version)
	}
	stream.WriteObjectEnd()
}

// UnmarshalUpgradePolicy reads a value of the 'upgrade_policy' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalUpgradePolicy(source interface{}) (object *UpgradePolicy, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadUpgradePolicy(iterator)
	err = iterator.Error
	return
}

// ReadUpgradePolicy reads a value of the 'upgrade_policy' type from the given iterator.
func ReadUpgradePolicy(iterator *jsoniter.Iterator) *UpgradePolicy {
	object := &UpgradePolicy{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == UpgradePolicyLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.bitmap_ |= 8
		case "enable_minor_version_upgrades":
			value := iterator.ReadBool()
			object.enableMinorVersionUpgrades = value
			object.bitmap_ |= 16
		case "next_run":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.nextRun = value
			object.bitmap_ |= 32
		case "schedule":
			value := iterator.ReadString()
			object.schedule = value
			object.bitmap_ |= 64
		case "schedule_type":
			text := iterator.ReadString()
			value := ScheduleType(text)
			object.scheduleType = value
			object.bitmap_ |= 128
		case "upgrade_type":
			text := iterator.ReadString()
			value := UpgradeType(text)
			object.upgradeType = value
			object.bitmap_ |= 256
		case "version":
			value := iterator.ReadString()
			object.version = value
			object.bitmap_ |= 512
		default:
			iterator.ReadAny()
		}
	}
	return object
}
