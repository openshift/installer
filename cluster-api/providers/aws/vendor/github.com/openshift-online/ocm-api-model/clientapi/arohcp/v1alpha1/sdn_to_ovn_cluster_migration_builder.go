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

// Details for `SdnToOvn` cluster migrations.
type SdnToOvnClusterMigrationBuilder struct {
	fieldSet_      []bool
	joinIpv4       string
	masqueradeIpv4 string
	transitIpv4    string
}

// NewSdnToOvnClusterMigration creates a new builder of 'sdn_to_ovn_cluster_migration' objects.
func NewSdnToOvnClusterMigration() *SdnToOvnClusterMigrationBuilder {
	return &SdnToOvnClusterMigrationBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SdnToOvnClusterMigrationBuilder) Empty() bool {
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

// JoinIpv4 sets the value of the 'join_ipv_4' attribute to the given value.
func (b *SdnToOvnClusterMigrationBuilder) JoinIpv4(value string) *SdnToOvnClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.joinIpv4 = value
	b.fieldSet_[0] = true
	return b
}

// MasqueradeIpv4 sets the value of the 'masquerade_ipv_4' attribute to the given value.
func (b *SdnToOvnClusterMigrationBuilder) MasqueradeIpv4(value string) *SdnToOvnClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.masqueradeIpv4 = value
	b.fieldSet_[1] = true
	return b
}

// TransitIpv4 sets the value of the 'transit_ipv_4' attribute to the given value.
func (b *SdnToOvnClusterMigrationBuilder) TransitIpv4(value string) *SdnToOvnClusterMigrationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.transitIpv4 = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SdnToOvnClusterMigrationBuilder) Copy(object *SdnToOvnClusterMigration) *SdnToOvnClusterMigrationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.joinIpv4 = object.joinIpv4
	b.masqueradeIpv4 = object.masqueradeIpv4
	b.transitIpv4 = object.transitIpv4
	return b
}

// Build creates a 'sdn_to_ovn_cluster_migration' object using the configuration stored in the builder.
func (b *SdnToOvnClusterMigrationBuilder) Build() (object *SdnToOvnClusterMigration, err error) {
	object = new(SdnToOvnClusterMigration)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.joinIpv4 = b.joinIpv4
	object.masqueradeIpv4 = b.masqueradeIpv4
	object.transitIpv4 = b.transitIpv4
	return
}
