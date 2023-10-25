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

// AuditLogBuilder contains the data and logic needed to build 'audit_log' objects.
//
// Contains the necessary attributes to support audit log forwarding
type AuditLogBuilder struct {
	bitmap_ uint32
	roleArn string
}

// NewAuditLog creates a new builder of 'audit_log' objects.
func NewAuditLog() *AuditLogBuilder {
	return &AuditLogBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AuditLogBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// RoleArn sets the value of the 'role_arn' attribute to the given value.
func (b *AuditLogBuilder) RoleArn(value string) *AuditLogBuilder {
	b.roleArn = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AuditLogBuilder) Copy(object *AuditLog) *AuditLogBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.roleArn = object.roleArn
	return b
}

// Build creates a 'audit_log' object using the configuration stored in the builder.
func (b *AuditLogBuilder) Build() (object *AuditLog, err error) {
	object = new(AuditLog)
	object.bitmap_ = b.bitmap_
	object.roleArn = b.roleArn
	return
}
