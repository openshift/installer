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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// DNSBuilder contains the data and logic needed to build 'DNS' objects.
//
// DNS settings of the cluster.
type DNSBuilder struct {
	bitmap_    uint32
	baseDomain string
}

// NewDNS creates a new builder of 'DNS' objects.
func NewDNS() *DNSBuilder {
	return &DNSBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DNSBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// BaseDomain sets the value of the 'base_domain' attribute to the given value.
func (b *DNSBuilder) BaseDomain(value string) *DNSBuilder {
	b.baseDomain = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DNSBuilder) Copy(object *DNS) *DNSBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.baseDomain = object.baseDomain
	return b
}

// Build creates a 'DNS' object using the configuration stored in the builder.
func (b *DNSBuilder) Build() (object *DNS, err error) {
	object = new(DNS)
	object.bitmap_ = b.bitmap_
	object.baseDomain = b.baseDomain
	return
}
