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

// Proxy represents the values of the 'proxy' type.
//
// Proxy configuration of a cluster.
type Proxy struct {
	bitmap_    uint32
	httpProxy  string
	httpsProxy string
	noProxy    string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Proxy) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// HTTPProxy returns the value of the 'HTTP_proxy' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HTTPProxy is the URL of the proxy for HTTP requests.
func (o *Proxy) HTTPProxy() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.httpProxy
	}
	return ""
}

// GetHTTPProxy returns the value of the 'HTTP_proxy' attribute and
// a flag indicating if the attribute has a value.
//
// HTTPProxy is the URL of the proxy for HTTP requests.
func (o *Proxy) GetHTTPProxy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.httpProxy
	}
	return
}

// HTTPSProxy returns the value of the 'HTTPS_proxy' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HTTPSProxy is the URL of the proxy for HTTPS requests.
func (o *Proxy) HTTPSProxy() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.httpsProxy
	}
	return ""
}

// GetHTTPSProxy returns the value of the 'HTTPS_proxy' attribute and
// a flag indicating if the attribute has a value.
//
// HTTPSProxy is the URL of the proxy for HTTPS requests.
func (o *Proxy) GetHTTPSProxy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.httpsProxy
	}
	return
}

// NoProxy returns the value of the 'no_proxy' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// NoProxy is a comma-separated list of domains and CIDRs for which
// the proxy should not be used
func (o *Proxy) NoProxy() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.noProxy
	}
	return ""
}

// GetNoProxy returns the value of the 'no_proxy' attribute and
// a flag indicating if the attribute has a value.
//
// NoProxy is a comma-separated list of domains and CIDRs for which
// the proxy should not be used
func (o *Proxy) GetNoProxy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.noProxy
	}
	return
}

// ProxyListKind is the name of the type used to represent list of objects of
// type 'proxy'.
const ProxyListKind = "ProxyList"

// ProxyListLinkKind is the name of the type used to represent links to list
// of objects of type 'proxy'.
const ProxyListLinkKind = "ProxyListLink"

// ProxyNilKind is the name of the type used to nil lists of objects of
// type 'proxy'.
const ProxyListNilKind = "ProxyListNil"

// ProxyList is a list of values of the 'proxy' type.
type ProxyList struct {
	href  string
	link  bool
	items []*Proxy
}

// Len returns the length of the list.
func (l *ProxyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ProxyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ProxyList) Get(i int) *Proxy {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *ProxyList) Slice() []*Proxy {
	var slice []*Proxy
	if l == nil {
		slice = make([]*Proxy, 0)
	} else {
		slice = make([]*Proxy, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ProxyList) Each(f func(item *Proxy) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *ProxyList) Range(f func(index int, item *Proxy) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
