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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalControlPlaneUpgradePolicy writes a value of the 'control_plane_upgrade_policy' type to the given writer.
func MarshalControlPlaneUpgradePolicy(object *ControlPlaneUpgradePolicy, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteControlPlaneUpgradePolicy(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteControlPlaneUpgradePolicy writes a value of the 'control_plane_upgrade_policy' type to the given stream.
func WriteControlPlaneUpgradePolicy(object *ControlPlaneUpgradePolicy, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ControlPlaneUpgradePolicyLinkKind)
	} else {
		stream.WriteString(ControlPlaneUpgradePolicyKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enable_minor_version_upgrades")
		stream.WriteBool(object.enableMinorVersionUpgrades)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_update_timestamp")
		stream.WriteString((object.lastUpdateTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("next_run")
		stream.WriteString((object.nextRun).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("schedule")
		stream.WriteString(object.schedule)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("schedule_type")
		stream.WriteString(string(object.scheduleType))
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.state != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		WriteUpgradePolicyState(object.state, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("upgrade_type")
		stream.WriteString(string(object.upgradeType))
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		stream.WriteString(object.version)
	}
	stream.WriteObjectEnd()
}

// UnmarshalControlPlaneUpgradePolicy reads a value of the 'control_plane_upgrade_policy' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalControlPlaneUpgradePolicy(source interface{}) (object *ControlPlaneUpgradePolicy, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadControlPlaneUpgradePolicy(iterator)
	err = iterator.Error
	return
}

// ReadControlPlaneUpgradePolicy reads a value of the 'control_plane_upgrade_policy' type from the given iterator.
func ReadControlPlaneUpgradePolicy(iterator *jsoniter.Iterator) *ControlPlaneUpgradePolicy {
	object := &ControlPlaneUpgradePolicy{
		fieldSet_: make([]bool, 13),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ControlPlaneUpgradePolicyLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.fieldSet_[3] = true
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.fieldSet_[4] = true
		case "enable_minor_version_upgrades":
			value := iterator.ReadBool()
			object.enableMinorVersionUpgrades = value
			object.fieldSet_[5] = true
		case "last_update_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdateTimestamp = value
			object.fieldSet_[6] = true
		case "next_run":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.nextRun = value
			object.fieldSet_[7] = true
		case "schedule":
			value := iterator.ReadString()
			object.schedule = value
			object.fieldSet_[8] = true
		case "schedule_type":
			text := iterator.ReadString()
			value := ScheduleType(text)
			object.scheduleType = value
			object.fieldSet_[9] = true
		case "state":
			value := ReadUpgradePolicyState(iterator)
			object.state = value
			object.fieldSet_[10] = true
		case "upgrade_type":
			text := iterator.ReadString()
			value := UpgradeType(text)
			object.upgradeType = value
			object.fieldSet_[11] = true
		case "version":
			value := iterator.ReadString()
			object.version = value
			object.fieldSet_[12] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
