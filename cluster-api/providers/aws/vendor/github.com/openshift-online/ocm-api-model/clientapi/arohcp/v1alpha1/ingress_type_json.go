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

import (
	"io"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalIngress writes a value of the 'ingress' type to the given writer.
func MarshalIngress(object *Ingress, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteIngress(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteIngress writes a value of the 'ingress' type to the given stream.
func WriteIngress(object *Ingress, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(IngressLinkKind)
	} else {
		stream.WriteString(IngressKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dns_name")
		stream.WriteString(object.dnsName)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_routes_hostname")
		stream.WriteString(object.clusterRoutesHostname)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_routes_tls_secret_ref")
		stream.WriteString(object.clusterRoutesTlsSecretRef)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.componentRoutes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("component_routes")
		if object.componentRoutes != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.componentRoutes))
			i := 0
			for key := range object.componentRoutes {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.componentRoutes[key]
				stream.WriteObjectField(key)
				WriteComponentRoute(item, stream)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("default")
		stream.WriteBool(object.default_)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.excludedNamespaces != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("excluded_namespaces")
		WriteStringList(object.excludedNamespaces, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("listening")
		stream.WriteString(string(object.listening))
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("load_balancer_type")
		stream.WriteString(string(object.loadBalancerType))
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("route_namespace_ownership_policy")
		stream.WriteString(string(object.routeNamespaceOwnershipPolicy))
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.routeSelectors != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("route_selectors")
		if object.routeSelectors != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.routeSelectors))
			i := 0
			for key := range object.routeSelectors {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.routeSelectors[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("route_wildcard_policy")
		stream.WriteString(string(object.routeWildcardPolicy))
	}
	stream.WriteObjectEnd()
}

// UnmarshalIngress reads a value of the 'ingress' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalIngress(source interface{}) (object *Ingress, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadIngress(iterator)
	err = iterator.Error
	return
}

// ReadIngress reads a value of the 'ingress' type from the given iterator.
func ReadIngress(iterator *jsoniter.Iterator) *Ingress {
	object := &Ingress{
		fieldSet_: make([]bool, 14),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == IngressLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "dns_name":
			value := iterator.ReadString()
			object.dnsName = value
			object.fieldSet_[3] = true
		case "cluster_routes_hostname":
			value := iterator.ReadString()
			object.clusterRoutesHostname = value
			object.fieldSet_[4] = true
		case "cluster_routes_tls_secret_ref":
			value := iterator.ReadString()
			object.clusterRoutesTlsSecretRef = value
			object.fieldSet_[5] = true
		case "component_routes":
			value := map[string]*ComponentRoute{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := ReadComponentRoute(iterator)
				value[key] = item
			}
			object.componentRoutes = value
			object.fieldSet_[6] = true
		case "default":
			value := iterator.ReadBool()
			object.default_ = value
			object.fieldSet_[7] = true
		case "excluded_namespaces":
			value := ReadStringList(iterator)
			object.excludedNamespaces = value
			object.fieldSet_[8] = true
		case "listening":
			text := iterator.ReadString()
			value := ListeningMethod(text)
			object.listening = value
			object.fieldSet_[9] = true
		case "load_balancer_type":
			text := iterator.ReadString()
			value := LoadBalancerFlavor(text)
			object.loadBalancerType = value
			object.fieldSet_[10] = true
		case "route_namespace_ownership_policy":
			text := iterator.ReadString()
			value := NamespaceOwnershipPolicy(text)
			object.routeNamespaceOwnershipPolicy = value
			object.fieldSet_[11] = true
		case "route_selectors":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.routeSelectors = value
			object.fieldSet_[12] = true
		case "route_wildcard_policy":
			text := iterator.ReadString()
			value := WildcardPolicy(text)
			object.routeWildcardPolicy = value
			object.fieldSet_[13] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
