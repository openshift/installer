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

// Representation of a Control Plane
type ControlPlaneBuilder struct {
	fieldSet_     []bool
	backup        *BackupBuilder
	logForwarders *LogForwarderListBuilder
}

// NewControlPlane creates a new builder of 'control_plane' objects.
func NewControlPlane() *ControlPlaneBuilder {
	return &ControlPlaneBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ControlPlaneBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Backup sets the value of the 'backup' attribute to the given value.
//
// Representation of a Backup.
func (b *ControlPlaneBuilder) Backup(value *BackupBuilder) *ControlPlaneBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.backup = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// LogForwarders sets the value of the 'log_forwarders' attribute to the given values.
func (b *ControlPlaneBuilder) LogForwarders(value *LogForwarderListBuilder) *ControlPlaneBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.logForwarders = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ControlPlaneBuilder) Copy(object *ControlPlane) *ControlPlaneBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.backup != nil {
		b.backup = NewBackup().Copy(object.backup)
	} else {
		b.backup = nil
	}
	if object.logForwarders != nil {
		b.logForwarders = NewLogForwarderList().Copy(object.logForwarders)
	} else {
		b.logForwarders = nil
	}
	return b
}

// Build creates a 'control_plane' object using the configuration stored in the builder.
func (b *ControlPlaneBuilder) Build() (object *ControlPlane, err error) {
	object = new(ControlPlane)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.backup != nil {
		object.backup, err = b.backup.Build()
		if err != nil {
			return
		}
	}
	if b.logForwarders != nil {
		object.logForwarders, err = b.logForwarders.Build()
		if err != nil {
			return
		}
	}
	return
}
