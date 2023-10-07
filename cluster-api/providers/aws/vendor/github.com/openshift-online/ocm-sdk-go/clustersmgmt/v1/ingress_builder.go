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

// IngressBuilder contains the data and logic needed to build 'ingress' objects.
//
// Representation of an ingress.
type IngressBuilder struct {
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

// NewIngress creates a new builder of 'ingress' objects.
func NewIngress() *IngressBuilder {
	return &IngressBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *IngressBuilder) Link(value bool) *IngressBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *IngressBuilder) ID(value string) *IngressBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *IngressBuilder) HREF(value string) *IngressBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *IngressBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// DNSName sets the value of the 'DNS_name' attribute to the given value.
func (b *IngressBuilder) DNSName(value string) *IngressBuilder {
	b.dnsName = value
	b.bitmap_ |= 8
	return b
}

// ClusterRoutesHostname sets the value of the 'cluster_routes_hostname' attribute to the given value.
func (b *IngressBuilder) ClusterRoutesHostname(value string) *IngressBuilder {
	b.clusterRoutesHostname = value
	b.bitmap_ |= 16
	return b
}

// ClusterRoutesTlsSecretRef sets the value of the 'cluster_routes_tls_secret_ref' attribute to the given value.
func (b *IngressBuilder) ClusterRoutesTlsSecretRef(value string) *IngressBuilder {
	b.clusterRoutesTlsSecretRef = value
	b.bitmap_ |= 32
	return b
}

// Default sets the value of the 'default' attribute to the given value.
func (b *IngressBuilder) Default(value bool) *IngressBuilder {
	b.default_ = value
	b.bitmap_ |= 64
	return b
}

// ExcludedNamespaces sets the value of the 'excluded_namespaces' attribute to the given values.
func (b *IngressBuilder) ExcludedNamespaces(values ...string) *IngressBuilder {
	b.excludedNamespaces = make([]string, len(values))
	copy(b.excludedNamespaces, values)
	b.bitmap_ |= 128
	return b
}

// Listening sets the value of the 'listening' attribute to the given value.
//
// Cluster components listening method.
func (b *IngressBuilder) Listening(value ListeningMethod) *IngressBuilder {
	b.listening = value
	b.bitmap_ |= 256
	return b
}

// LoadBalancerType sets the value of the 'load_balancer_type' attribute to the given value.
//
// Type of load balancer for AWS cloud provider parameters.
func (b *IngressBuilder) LoadBalancerType(value LoadBalancerFlavor) *IngressBuilder {
	b.loadBalancerType = value
	b.bitmap_ |= 512
	return b
}

// RouteNamespaceOwnershipPolicy sets the value of the 'route_namespace_ownership_policy' attribute to the given value.
//
// Type of Namespace Ownership Policy.
func (b *IngressBuilder) RouteNamespaceOwnershipPolicy(value NamespaceOwnershipPolicy) *IngressBuilder {
	b.routeNamespaceOwnershipPolicy = value
	b.bitmap_ |= 1024
	return b
}

// RouteSelectors sets the value of the 'route_selectors' attribute to the given value.
func (b *IngressBuilder) RouteSelectors(value map[string]string) *IngressBuilder {
	b.routeSelectors = value
	if value != nil {
		b.bitmap_ |= 2048
	} else {
		b.bitmap_ &^= 2048
	}
	return b
}

// RouteWildcardPolicy sets the value of the 'route_wildcard_policy' attribute to the given value.
//
// Type of wildcard policy.
func (b *IngressBuilder) RouteWildcardPolicy(value WildcardPolicy) *IngressBuilder {
	b.routeWildcardPolicy = value
	b.bitmap_ |= 4096
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *IngressBuilder) Copy(object *Ingress) *IngressBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.dnsName = object.dnsName
	b.clusterRoutesHostname = object.clusterRoutesHostname
	b.clusterRoutesTlsSecretRef = object.clusterRoutesTlsSecretRef
	b.default_ = object.default_
	if object.excludedNamespaces != nil {
		b.excludedNamespaces = make([]string, len(object.excludedNamespaces))
		copy(b.excludedNamespaces, object.excludedNamespaces)
	} else {
		b.excludedNamespaces = nil
	}
	b.listening = object.listening
	b.loadBalancerType = object.loadBalancerType
	b.routeNamespaceOwnershipPolicy = object.routeNamespaceOwnershipPolicy
	if len(object.routeSelectors) > 0 {
		b.routeSelectors = map[string]string{}
		for k, v := range object.routeSelectors {
			b.routeSelectors[k] = v
		}
	} else {
		b.routeSelectors = nil
	}
	b.routeWildcardPolicy = object.routeWildcardPolicy
	return b
}

// Build creates a 'ingress' object using the configuration stored in the builder.
func (b *IngressBuilder) Build() (object *Ingress, err error) {
	object = new(Ingress)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.dnsName = b.dnsName
	object.clusterRoutesHostname = b.clusterRoutesHostname
	object.clusterRoutesTlsSecretRef = b.clusterRoutesTlsSecretRef
	object.default_ = b.default_
	if b.excludedNamespaces != nil {
		object.excludedNamespaces = make([]string, len(b.excludedNamespaces))
		copy(object.excludedNamespaces, b.excludedNamespaces)
	}
	object.listening = b.listening
	object.loadBalancerType = b.loadBalancerType
	object.routeNamespaceOwnershipPolicy = b.routeNamespaceOwnershipPolicy
	if b.routeSelectors != nil {
		object.routeSelectors = make(map[string]string)
		for k, v := range b.routeSelectors {
			object.routeSelectors[k] = v
		}
	}
	object.routeWildcardPolicy = b.routeWildcardPolicy
	return
}
