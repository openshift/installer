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

// SdnToOvnClusterMigrationBuilder contains the data and logic needed to build 'sdn_to_ovn_cluster_migration' objects.
//
// Details for `SdnToOvn` cluster migrations.
type SdnToOvnClusterMigrationBuilder struct {
	bitmap_        uint32
	joinIpv4       string
	masqueradeIpv4 string
	transitIpv4    string
}

// NewSdnToOvnClusterMigration creates a new builder of 'sdn_to_ovn_cluster_migration' objects.
func NewSdnToOvnClusterMigration() *SdnToOvnClusterMigrationBuilder {
	return &SdnToOvnClusterMigrationBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SdnToOvnClusterMigrationBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// JoinIpv4 sets the value of the 'join_ipv_4' attribute to the given value.
func (b *SdnToOvnClusterMigrationBuilder) JoinIpv4(value string) *SdnToOvnClusterMigrationBuilder {
	b.joinIpv4 = value
	b.bitmap_ |= 1
	return b
}

// MasqueradeIpv4 sets the value of the 'masquerade_ipv_4' attribute to the given value.
func (b *SdnToOvnClusterMigrationBuilder) MasqueradeIpv4(value string) *SdnToOvnClusterMigrationBuilder {
	b.masqueradeIpv4 = value
	b.bitmap_ |= 2
	return b
}

// TransitIpv4 sets the value of the 'transit_ipv_4' attribute to the given value.
func (b *SdnToOvnClusterMigrationBuilder) TransitIpv4(value string) *SdnToOvnClusterMigrationBuilder {
	b.transitIpv4 = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SdnToOvnClusterMigrationBuilder) Copy(object *SdnToOvnClusterMigration) *SdnToOvnClusterMigrationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.joinIpv4 = object.joinIpv4
	b.masqueradeIpv4 = object.masqueradeIpv4
	b.transitIpv4 = object.transitIpv4
	return b
}

// Build creates a 'sdn_to_ovn_cluster_migration' object using the configuration stored in the builder.
func (b *SdnToOvnClusterMigrationBuilder) Build() (object *SdnToOvnClusterMigration, err error) {
	object = new(SdnToOvnClusterMigration)
	object.bitmap_ = b.bitmap_
	object.joinIpv4 = b.joinIpv4
	object.masqueradeIpv4 = b.masqueradeIpv4
	object.transitIpv4 = b.transitIpv4
	return
}
