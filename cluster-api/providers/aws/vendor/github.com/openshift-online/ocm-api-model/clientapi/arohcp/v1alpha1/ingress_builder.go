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

// Representation of an ingress.
type IngressBuilder struct {
	fieldSet_                     []bool
	id                            string
	href                          string
	dnsName                       string
	clusterRoutesHostname         string
	clusterRoutesTlsSecretRef     string
	componentRoutes               map[string]*ComponentRouteBuilder
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
	return &IngressBuilder{
		fieldSet_: make([]bool, 14),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *IngressBuilder) Link(value bool) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *IngressBuilder) ID(value string) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *IngressBuilder) HREF(value string) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *IngressBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// DNSName sets the value of the 'DNS_name' attribute to the given value.
func (b *IngressBuilder) DNSName(value string) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.dnsName = value
	b.fieldSet_[3] = true
	return b
}

// ClusterRoutesHostname sets the value of the 'cluster_routes_hostname' attribute to the given value.
func (b *IngressBuilder) ClusterRoutesHostname(value string) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.clusterRoutesHostname = value
	b.fieldSet_[4] = true
	return b
}

// ClusterRoutesTlsSecretRef sets the value of the 'cluster_routes_tls_secret_ref' attribute to the given value.
func (b *IngressBuilder) ClusterRoutesTlsSecretRef(value string) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.clusterRoutesTlsSecretRef = value
	b.fieldSet_[5] = true
	return b
}

// ComponentRoutes sets the value of the 'component_routes' attribute to the given value.
func (b *IngressBuilder) ComponentRoutes(value map[string]*ComponentRouteBuilder) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.componentRoutes = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Default sets the value of the 'default' attribute to the given value.
func (b *IngressBuilder) Default(value bool) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.default_ = value
	b.fieldSet_[7] = true
	return b
}

// ExcludedNamespaces sets the value of the 'excluded_namespaces' attribute to the given values.
func (b *IngressBuilder) ExcludedNamespaces(values ...string) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.excludedNamespaces = make([]string, len(values))
	copy(b.excludedNamespaces, values)
	b.fieldSet_[8] = true
	return b
}

// Listening sets the value of the 'listening' attribute to the given value.
//
// Cluster components listening method.
func (b *IngressBuilder) Listening(value ListeningMethod) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.listening = value
	b.fieldSet_[9] = true
	return b
}

// LoadBalancerType sets the value of the 'load_balancer_type' attribute to the given value.
//
// Type of load balancer for AWS cloud provider parameters.
func (b *IngressBuilder) LoadBalancerType(value LoadBalancerFlavor) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.loadBalancerType = value
	b.fieldSet_[10] = true
	return b
}

// RouteNamespaceOwnershipPolicy sets the value of the 'route_namespace_ownership_policy' attribute to the given value.
//
// Type of Namespace Ownership Policy.
func (b *IngressBuilder) RouteNamespaceOwnershipPolicy(value NamespaceOwnershipPolicy) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.routeNamespaceOwnershipPolicy = value
	b.fieldSet_[11] = true
	return b
}

// RouteSelectors sets the value of the 'route_selectors' attribute to the given value.
func (b *IngressBuilder) RouteSelectors(value map[string]string) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.routeSelectors = value
	if value != nil {
		b.fieldSet_[12] = true
	} else {
		b.fieldSet_[12] = false
	}
	return b
}

// RouteWildcardPolicy sets the value of the 'route_wildcard_policy' attribute to the given value.
//
// Type of wildcard policy.
func (b *IngressBuilder) RouteWildcardPolicy(value WildcardPolicy) *IngressBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 14)
	}
	b.routeWildcardPolicy = value
	b.fieldSet_[13] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *IngressBuilder) Copy(object *Ingress) *IngressBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.dnsName = object.dnsName
	b.clusterRoutesHostname = object.clusterRoutesHostname
	b.clusterRoutesTlsSecretRef = object.clusterRoutesTlsSecretRef
	if len(object.componentRoutes) > 0 {
		b.componentRoutes = map[string]*ComponentRouteBuilder{}
		for k, v := range object.componentRoutes {
			b.componentRoutes[k] = NewComponentRoute().Copy(v)
		}
	} else {
		b.componentRoutes = nil
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.dnsName = b.dnsName
	object.clusterRoutesHostname = b.clusterRoutesHostname
	object.clusterRoutesTlsSecretRef = b.clusterRoutesTlsSecretRef
	if b.componentRoutes != nil {
		object.componentRoutes = make(map[string]*ComponentRoute)
		for k, v := range b.componentRoutes {
			object.componentRoutes[k], err = v.Build()
			if err != nil {
				return
			}
		}
	}
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
