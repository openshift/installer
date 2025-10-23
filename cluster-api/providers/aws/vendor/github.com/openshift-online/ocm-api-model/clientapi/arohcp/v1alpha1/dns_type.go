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

// DNS represents the values of the 'DNS' type.
//
// DNS settings of the cluster.
type DNS struct {
	fieldSet_  []bool
	baseDomain string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *DNS) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// BaseDomain returns the value of the 'base_domain' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Base DNS domain of the cluster.
//
// During the installation of the cluster it is necessary to create multiple DNS records.
// They will be created as sub-domains of this domain. For example, if the domain_prefix of the
// cluster is `mycluster` and the base domain is `example.com` then the following DNS
// records will be created:
//
// ```
// mycluster-api.example.com
// mycluster-etcd-0.example.com
// mycluster-etcd-1.example.com
// mycluster-etcd-3.example.com
// ```
//
// The exact number, type and names of the created DNS record depends on the characteristics
// of the cluster, and may be different for different versions of _OpenShift_. Please don't
// rely on them. For example, to find what is the URL of the Kubernetes API server of the
// cluster don't assume that it will be `mycluster-api.example.com`. Instead of that use
// this API to retrieve the description of the cluster, and get it from the `api.url`
// attribute. For example, if the identifier of the cluster is `123` send a request like
// this:
//
// ```http
// GET /api/clusters_mgmt/v1/clusters/123 HTTP/1.1
// ```
//
// That will return a response like this, including the `api.url` attribute:
//
// ```json
//
//	{
//	    "kind": "Cluster",
//	    "id": "123",
//	    "href": "/api/clusters_mgmt/v1/clusters/123",
//	        "api": {
//	        "url": "https://mycluster-api.example.com:6443"
//	    },
//	    ...
//	}
//
// ```
//
// When the cluster is created in Amazon Web Services it is necessary to create this base
// DNS domain in advance, using AWS Route53 (https://console.aws.amazon.com/route53).
func (o *DNS) BaseDomain() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.baseDomain
	}
	return ""
}

// GetBaseDomain returns the value of the 'base_domain' attribute and
// a flag indicating if the attribute has a value.
//
// Base DNS domain of the cluster.
//
// During the installation of the cluster it is necessary to create multiple DNS records.
// They will be created as sub-domains of this domain. For example, if the domain_prefix of the
// cluster is `mycluster` and the base domain is `example.com` then the following DNS
// records will be created:
//
// ```
// mycluster-api.example.com
// mycluster-etcd-0.example.com
// mycluster-etcd-1.example.com
// mycluster-etcd-3.example.com
// ```
//
// The exact number, type and names of the created DNS record depends on the characteristics
// of the cluster, and may be different for different versions of _OpenShift_. Please don't
// rely on them. For example, to find what is the URL of the Kubernetes API server of the
// cluster don't assume that it will be `mycluster-api.example.com`. Instead of that use
// this API to retrieve the description of the cluster, and get it from the `api.url`
// attribute. For example, if the identifier of the cluster is `123` send a request like
// this:
//
// ```http
// GET /api/clusters_mgmt/v1/clusters/123 HTTP/1.1
// ```
//
// That will return a response like this, including the `api.url` attribute:
//
// ```json
//
//	{
//	    "kind": "Cluster",
//	    "id": "123",
//	    "href": "/api/clusters_mgmt/v1/clusters/123",
//	        "api": {
//	        "url": "https://mycluster-api.example.com:6443"
//	    },
//	    ...
//	}
//
// ```
//
// When the cluster is created in Amazon Web Services it is necessary to create this base
// DNS domain in advance, using AWS Route53 (https://console.aws.amazon.com/route53).
func (o *DNS) GetBaseDomain() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.baseDomain
	}
	return
}

// DNSListKind is the name of the type used to represent list of objects of
// type 'DNS'.
const DNSListKind = "DNSList"

// DNSListLinkKind is the name of the type used to represent links to list
// of objects of type 'DNS'.
const DNSListLinkKind = "DNSListLink"

// DNSNilKind is the name of the type used to nil lists of objects of
// type 'DNS'.
const DNSListNilKind = "DNSListNil"

// DNSList is a list of values of the 'DNS' type.
type DNSList struct {
	href  string
	link  bool
	items []*DNS
}

// Len returns the length of the list.
func (l *DNSList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *DNSList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *DNSList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *DNSList) SetItems(items []*DNS) {
	l.items = items
}

// Items returns the items of the list.
func (l *DNSList) Items() []*DNS {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *DNSList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *DNSList) Get(i int) *DNS {
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
func (l *DNSList) Slice() []*DNS {
	var slice []*DNS
	if l == nil {
		slice = make([]*DNS, 0)
	} else {
		slice = make([]*DNS, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *DNSList) Each(f func(item *DNS) bool) {
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
func (l *DNSList) Range(f func(index int, item *DNS) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
