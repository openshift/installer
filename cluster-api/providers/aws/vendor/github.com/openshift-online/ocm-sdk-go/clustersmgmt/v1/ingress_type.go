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

// IngressKind is the name of the type used to represent objects
// of type 'ingress'.
const IngressKind = "Ingress"

// IngressLinkKind is the name of the type used to represent links
// to objects of type 'ingress'.
const IngressLinkKind = "IngressLink"

// IngressNilKind is the name of the type used to nil references
// to objects of type 'ingress'.
const IngressNilKind = "IngressNil"

// Ingress represents the values of the 'ingress' type.
//
// Representation of an ingress.
type Ingress struct {
	bitmap_                       uint32
	id                            string
	href                          string
	dnsName                       string
	clusterRoutesHostname         string
	clusterRoutesTlsSecretRef     string
	excludedNamespaces            []string
	listening                     ListeningMethod
	loadBalancerType              LoadBalancerFlavor
	routeNamespaceOwnershipPolicy NamespaceOwnershipPolicy
	routeSelectors                map[string]string
	routeWildcardPolicy           WildcardPolicy
	default_                      bool
}

// Kind returns the name of the type of the object.
func (o *Ingress) Kind() string {
	if o == nil {
		return IngressNilKind
	}
	if o.bitmap_&1 != 0 {
		return IngressLinkKind
	}
	return IngressKind
}

// Link returns true iif this is a link.
func (o *Ingress) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *Ingress) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Ingress) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Ingress) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Ingress) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Ingress) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// DNSName returns the value of the 'DNS_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DNS Name of the ingress.
func (o *Ingress) DNSName() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.dnsName
	}
	return ""
}

// GetDNSName returns the value of the 'DNS_name' attribute and
// a flag indicating if the attribute has a value.
//
// DNS Name of the ingress.
func (o *Ingress) GetDNSName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.dnsName
	}
	return
}

// ClusterRoutesHostname returns the value of the 'cluster_routes_hostname' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster routes hostname.
func (o *Ingress) ClusterRoutesHostname() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.clusterRoutesHostname
	}
	return ""
}

// GetClusterRoutesHostname returns the value of the 'cluster_routes_hostname' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster routes hostname.
func (o *Ingress) GetClusterRoutesHostname() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.clusterRoutesHostname
	}
	return
}

// ClusterRoutesTlsSecretRef returns the value of the 'cluster_routes_tls_secret_ref' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Cluster routes TLS Secret reference.
func (o *Ingress) ClusterRoutesTlsSecretRef() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.clusterRoutesTlsSecretRef
	}
	return ""
}

// GetClusterRoutesTlsSecretRef returns the value of the 'cluster_routes_tls_secret_ref' attribute and
// a flag indicating if the attribute has a value.
//
// Cluster routes TLS Secret reference.
func (o *Ingress) GetClusterRoutesTlsSecretRef() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.clusterRoutesTlsSecretRef
	}
	return
}

// Default returns the value of the 'default' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if this is the default ingress.
func (o *Ingress) Default() bool {
	if o != nil && o.bitmap_&64 != 0 {
		return o.default_
	}
	return false
}

// GetDefault returns the value of the 'default' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if this is the default ingress.
func (o *Ingress) GetDefault() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.default_
	}
	return
}

// ExcludedNamespaces returns the value of the 'excluded_namespaces' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A set of excluded namespaces for the ingress.
func (o *Ingress) ExcludedNamespaces() []string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.excludedNamespaces
	}
	return nil
}

// GetExcludedNamespaces returns the value of the 'excluded_namespaces' attribute and
// a flag indicating if the attribute has a value.
//
// A set of excluded namespaces for the ingress.
func (o *Ingress) GetExcludedNamespaces() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.excludedNamespaces
	}
	return
}

// Listening returns the value of the 'listening' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Listening method of the ingress
func (o *Ingress) Listening() ListeningMethod {
	if o != nil && o.bitmap_&256 != 0 {
		return o.listening
	}
	return ListeningMethod("")
}

// GetListening returns the value of the 'listening' attribute and
// a flag indicating if the attribute has a value.
//
// Listening method of the ingress
func (o *Ingress) GetListening() (value ListeningMethod, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.listening
	}
	return
}

// LoadBalancerType returns the value of the 'load_balancer_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Load Balancer type of the ingress
func (o *Ingress) LoadBalancerType() LoadBalancerFlavor {
	if o != nil && o.bitmap_&512 != 0 {
		return o.loadBalancerType
	}
	return LoadBalancerFlavor("")
}

// GetLoadBalancerType returns the value of the 'load_balancer_type' attribute and
// a flag indicating if the attribute has a value.
//
// Load Balancer type of the ingress
func (o *Ingress) GetLoadBalancerType() (value LoadBalancerFlavor, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.loadBalancerType
	}
	return
}

// RouteNamespaceOwnershipPolicy returns the value of the 'route_namespace_ownership_policy' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Namespace Ownership Policy for the ingress.
func (o *Ingress) RouteNamespaceOwnershipPolicy() NamespaceOwnershipPolicy {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.routeNamespaceOwnershipPolicy
	}
	return NamespaceOwnershipPolicy("")
}

// GetRouteNamespaceOwnershipPolicy returns the value of the 'route_namespace_ownership_policy' attribute and
// a flag indicating if the attribute has a value.
//
// Namespace Ownership Policy for the ingress.
func (o *Ingress) GetRouteNamespaceOwnershipPolicy() (value NamespaceOwnershipPolicy, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.routeNamespaceOwnershipPolicy
	}
	return
}

// RouteSelectors returns the value of the 'route_selectors' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A set of labels for the ingress.
func (o *Ingress) RouteSelectors() map[string]string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.routeSelectors
	}
	return nil
}

// GetRouteSelectors returns the value of the 'route_selectors' attribute and
// a flag indicating if the attribute has a value.
//
// A set of labels for the ingress.
func (o *Ingress) GetRouteSelectors() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.routeSelectors
	}
	return
}

// RouteWildcardPolicy returns the value of the 'route_wildcard_policy' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Wildcard policy for the ingress.
func (o *Ingress) RouteWildcardPolicy() WildcardPolicy {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.routeWildcardPolicy
	}
	return WildcardPolicy("")
}

// GetRouteWildcardPolicy returns the value of the 'route_wildcard_policy' attribute and
// a flag indicating if the attribute has a value.
//
// Wildcard policy for the ingress.
func (o *Ingress) GetRouteWildcardPolicy() (value WildcardPolicy, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.routeWildcardPolicy
	}
	return
}

// IngressListKind is the name of the type used to represent list of objects of
// type 'ingress'.
const IngressListKind = "IngressList"

// IngressListLinkKind is the name of the type used to represent links to list
// of objects of type 'ingress'.
const IngressListLinkKind = "IngressListLink"

// IngressNilKind is the name of the type used to nil lists of objects of
// type 'ingress'.
const IngressListNilKind = "IngressListNil"

// IngressList is a list of values of the 'ingress' type.
type IngressList struct {
	href  string
	link  bool
	items []*Ingress
}

// Kind returns the name of the type of the object.
func (l *IngressList) Kind() string {
	if l == nil {
		return IngressListNilKind
	}
	if l.link {
		return IngressListLinkKind
	}
	return IngressListKind
}

// Link returns true iif this is a link.
func (l *IngressList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *IngressList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *IngressList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *IngressList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *IngressList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *IngressList) Get(i int) *Ingress {
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
func (l *IngressList) Slice() []*Ingress {
	var slice []*Ingress
	if l == nil {
		slice = make([]*Ingress, 0)
	} else {
		slice = make([]*Ingress, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *IngressList) Each(f func(item *Ingress) bool) {
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
func (l *IngressList) Range(f func(index int, item *Ingress) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
