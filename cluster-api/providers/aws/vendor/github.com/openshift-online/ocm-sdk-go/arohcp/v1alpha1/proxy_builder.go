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

// ProxyBuilder contains the data and logic needed to build 'proxy' objects.
//
// Proxy configuration of a cluster.
type ProxyBuilder struct {
	bitmap_    uint32
	httpProxy  string
	httpsProxy string
	noProxy    string
}

// NewProxy creates a new builder of 'proxy' objects.
func NewProxy() *ProxyBuilder {
	return &ProxyBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProxyBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// HTTPProxy sets the value of the 'HTTP_proxy' attribute to the given value.
func (b *ProxyBuilder) HTTPProxy(value string) *ProxyBuilder {
	b.httpProxy = value
	b.bitmap_ |= 1
	return b
}

// HTTPSProxy sets the value of the 'HTTPS_proxy' attribute to the given value.
func (b *ProxyBuilder) HTTPSProxy(value string) *ProxyBuilder {
	b.httpsProxy = value
	b.bitmap_ |= 2
	return b
}

// NoProxy sets the value of the 'no_proxy' attribute to the given value.
func (b *ProxyBuilder) NoProxy(value string) *ProxyBuilder {
	b.noProxy = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProxyBuilder) Copy(object *Proxy) *ProxyBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.httpProxy = object.httpProxy
	b.httpsProxy = object.httpsProxy
	b.noProxy = object.noProxy
	return b
}

// Build creates a 'proxy' object using the configuration stored in the builder.
func (b *ProxyBuilder) Build() (object *Proxy, err error) {
	object = new(Proxy)
	object.bitmap_ = b.bitmap_
	object.httpProxy = b.httpProxy
	object.httpsProxy = b.httpsProxy
	object.noProxy = b.noProxy
	return
}
