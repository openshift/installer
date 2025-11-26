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

// Proxy configuration of a cluster.
type ProxyBuilder struct {
	fieldSet_  []bool
	httpProxy  string
	httpsProxy string
	noProxy    string
}

// NewProxy creates a new builder of 'proxy' objects.
func NewProxy() *ProxyBuilder {
	return &ProxyBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProxyBuilder) Empty() bool {
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

// HTTPProxy sets the value of the 'HTTP_proxy' attribute to the given value.
func (b *ProxyBuilder) HTTPProxy(value string) *ProxyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.httpProxy = value
	b.fieldSet_[0] = true
	return b
}

// HTTPSProxy sets the value of the 'HTTPS_proxy' attribute to the given value.
func (b *ProxyBuilder) HTTPSProxy(value string) *ProxyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.httpsProxy = value
	b.fieldSet_[1] = true
	return b
}

// NoProxy sets the value of the 'no_proxy' attribute to the given value.
func (b *ProxyBuilder) NoProxy(value string) *ProxyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.noProxy = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProxyBuilder) Copy(object *Proxy) *ProxyBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.httpProxy = object.httpProxy
	b.httpsProxy = object.httpsProxy
	b.noProxy = object.noProxy
	return b
}

// Build creates a 'proxy' object using the configuration stored in the builder.
func (b *ProxyBuilder) Build() (object *Proxy, err error) {
	object = new(Proxy)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.httpProxy = b.httpProxy
	object.httpsProxy = b.httpsProxy
	object.noProxy = b.noProxy
	return
}
