// Copyright 2024 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package compute

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLForwardingRuleSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/ForwardingRule",
			Description: "The Compute ForwardingRule resource",
			StructName:  "ForwardingRule",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a ForwardingRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "forwardingRule",
						Required:    true,
						Description: "A full instance of a ForwardingRule",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a ForwardingRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "forwardingRule",
						Required:    true,
						Description: "A full instance of a ForwardingRule",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a ForwardingRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "forwardingRule",
						Required:    true,
						Description: "A full instance of a ForwardingRule",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all ForwardingRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many ForwardingRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"ForwardingRule": &dcl.Component{
					Title: "ForwardingRule",
					ID:    "projects/{{project}}/global/forwardingRules/{{name}}",
					Locations: []string{
						"region",
						"global",
					},
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"allPorts": &dcl.Property{
								Type:        "boolean",
								GoName:      "AllPorts",
								Description: "This field is used along with the `backend_service` field for internal load balancing or with the `target` field for internal TargetInstance. This field cannot be used with `port` or `portRange` fields. When the load balancing scheme is `INTERNAL` and protocol is TCP/UDP, specify this field to allow packets addressed to any ports will be forwarded to the backends configured with this forwarding rule.",
								Immutable:   true,
							},
							"allowGlobalAccess": &dcl.Property{
								Type:        "boolean",
								GoName:      "AllowGlobalAccess",
								Description: "This field is used along with the `backend_service` field for internal load balancing or with the `target` field for internal TargetInstance. If the field is set to `TRUE`, clients can access ILB from all regions. Otherwise only allows access from clients in the same region as the internal load balancer.",
							},
							"backendService": &dcl.Property{
								Type:        "string",
								GoName:      "BackendService",
								Description: "This field is only used for `INTERNAL` load balancing. For internal load balancing, this field identifies the BackendService resource to receive the matched traffic.",
								Immutable:   true,
							},
							"baseForwardingRule": &dcl.Property{
								Type:        "string",
								GoName:      "BaseForwardingRule",
								ReadOnly:    true,
								Description: "[Output Only] The URL for the corresponding base Forwarding Rule. By base Forwarding Rule, we mean the Forwarding Rule that has the same IP address, protocol, and port settings with the current Forwarding Rule, but without sourceIPRanges specified. Always empty if the current Forwarding Rule does not have sourceIPRanges specified.",
								Immutable:   true,
							},
							"creationTimestamp": &dcl.Property{
								Type:        "string",
								GoName:      "CreationTimestamp",
								ReadOnly:    true,
								Description: "[Output Only] Creation timestamp in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt) text format.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. Provide this property when you create the resource.",
								Immutable:   true,
							},
							"ipAddress": &dcl.Property{
								Type:          "string",
								GoName:        "IPAddress",
								Description:   "IP address that this forwarding rule serves. When a client sends traffic to this IP address, the forwarding rule directs the traffic to the target that you specify in the forwarding rule. If you don't specify a reserved IP address, an ephemeral IP address is assigned. Methods for specifying an IP address: * IPv4 dotted decimal, as in `100.1.2.3` * Full URL, as in `https://www.googleapis.com/compute/v1/projects/project_id/regions/region/addresses/address-name` * Partial URL or by name, as in: * `projects/project_id/regions/region/addresses/address-name` * `regions/region/addresses/address-name` * `global/addresses/address-name` * `address-name` The loadBalancingScheme and the forwarding rule's target determine the type of IP address that you can use. For detailed information, refer to [IP address specifications](/load-balancing/docs/forwarding-rule-concepts#ip_address_specifications).",
								Immutable:     true,
								ServerDefault: true,
							},
							"ipProtocol": &dcl.Property{
								Type:          "string",
								GoName:        "IPProtocol",
								GoType:        "ForwardingRuleIPProtocolEnum",
								Description:   "The IP protocol to which this rule applies. For protocol forwarding, valid options are `TCP`, `UDP`, `ESP`, `AH`, `SCTP` or `ICMP`. For Internal TCP/UDP Load Balancing, the load balancing scheme is `INTERNAL`, and one of `TCP` or `UDP` are valid. For Traffic Director, the load balancing scheme is `INTERNAL_SELF_MANAGED`, and only `TCP`is valid. For Internal HTTP(S) Load Balancing, the load balancing scheme is `INTERNAL_MANAGED`, and only `TCP` is valid. For HTTP(S), SSL Proxy, and TCP Proxy Load Balancing, the load balancing scheme is `EXTERNAL` and only `TCP` is valid. For Network TCP/UDP Load Balancing, the load balancing scheme is `EXTERNAL`, and one of `TCP` or `UDP` is valid.",
								Immutable:     true,
								ServerDefault: true,
								Enum: []string{
									"TCP",
									"UDP",
									"ESP",
									"AH",
									"SCTP",
									"ICMP",
									"L3_DEFAULT",
								},
							},
							"ipVersion": &dcl.Property{
								Type:        "string",
								GoName:      "IPVersion",
								GoType:      "ForwardingRuleIPVersionEnum",
								Description: "The IP Version that will be used by this forwarding rule. Valid options are `IPV4` or `IPV6`. This can only be specified for an external global forwarding rule. Possible values: UNSPECIFIED_VERSION, IPV4, IPV6",
								Immutable:   true,
								Enum: []string{
									"UNSPECIFIED_VERSION",
									"IPV4",
									"IPV6",
								},
							},
							"isMirroringCollector": &dcl.Property{
								Type:        "boolean",
								GoName:      "IsMirroringCollector",
								Description: "Indicates whether or not this load balancer can be used as a collector for packet mirroring. To prevent mirroring loops, instances behind this load balancer will not have their traffic mirrored even if a `PacketMirroring` rule applies to them. This can only be set to true for load balancers that have their `loadBalancingScheme` set to `INTERNAL`.",
								Immutable:   true,
							},
							"labelFingerprint": &dcl.Property{
								Type:        "string",
								GoName:      "LabelFingerprint",
								ReadOnly:    true,
								Description: "Used internally during label updates.",
								Immutable:   true,
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Labels to apply to this rule.",
							},
							"loadBalancingScheme": &dcl.Property{
								Type:        "string",
								GoName:      "LoadBalancingScheme",
								GoType:      "ForwardingRuleLoadBalancingSchemeEnum",
								Description: "Specifies the forwarding rule type.\n\n*   `EXTERNAL` is used for:\n    *   Classic Cloud VPN gateways\n    *   Protocol forwarding to VMs from an external IP address\n    *   The following load balancers: HTTP(S), SSL Proxy, TCP Proxy, and Network TCP/UDP\n*   `INTERNAL` is used for:\n    *   Protocol forwarding to VMs from an internal IP address\n    *   Internal TCP/UDP load balancers\n*   `INTERNAL_MANAGED` is used for:\n    *   Internal HTTP(S) load balancers\n*   `INTERNAL_SELF_MANAGED` is used for:\n    *   Traffic Director\n*   `EXTERNAL_MANAGED` is used for:\n    *   Global external HTTP(S) load balancers \n\nFor more information about forwarding rules, refer to [Forwarding rule concepts](/load-balancing/docs/forwarding-rule-concepts). Possible values: INVALID, INTERNAL, INTERNAL_MANAGED, INTERNAL_SELF_MANAGED, EXTERNAL, EXTERNAL_MANAGED",
								Immutable:   true,
								Enum: []string{
									"INVALID",
									"INTERNAL",
									"INTERNAL_MANAGED",
									"INTERNAL_SELF_MANAGED",
									"EXTERNAL",
									"EXTERNAL_MANAGED",
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location of this resource.",
								Immutable:   true,
								Parameter:   true,
							},
							"metadataFilter": &dcl.Property{
								Type:        "array",
								GoName:      "MetadataFilter",
								Description: "Opaque filter criteria used by Loadbalancer to restrict routing configuration to a limited set of [xDS](https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol) compliant clients. In their xDS requests to Loadbalancer, xDS clients present [node metadata](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto.html#config-core-v3-node). If a match takes place, the relevant configuration is made available to those proxies. Otherwise, all the resources (e.g. `TargetHttpProxy`, `UrlMap`) referenced by the `ForwardingRule` will not be visible to those proxies.\n\nFor each `metadataFilter` in this list, if its `filterMatchCriteria` is set to MATCH_ANY, at least one of the `filterLabel`s must match the corresponding label provided in the metadata. If its `filterMatchCriteria` is set to MATCH_ALL, then all of its `filterLabel`s must match with corresponding labels provided in the metadata.\n\n`metadataFilters` specified here will be applifed before those specified in the `UrlMap` that this `ForwardingRule` references.\n\n`metadataFilters` only applies to Loadbalancers that have their loadBalancingScheme set to `INTERNAL_SELF_MANAGED`.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "ForwardingRuleMetadataFilter",
									Required: []string{
										"filterMatchCriteria",
										"filterLabel",
									},
									Properties: map[string]*dcl.Property{
										"filterLabel": &dcl.Property{
											Type:        "array",
											GoName:      "FilterLabel",
											Description: "The list of label value pairs that must match labels in the provided metadata based on `filterMatchCriteria`\n\nThis list must not be empty and can have at the most 64 entries.",
											Immutable:   true,
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "object",
												GoType: "ForwardingRuleMetadataFilterFilterLabel",
												Required: []string{
													"name",
													"value",
												},
												Properties: map[string]*dcl.Property{
													"name": &dcl.Property{
														Type:        "string",
														GoName:      "Name",
														Description: "Name of metadata label.\n\nThe name can have a maximum length of 1024 characters and must be at least 1 character long.",
														Immutable:   true,
													},
													"value": &dcl.Property{
														Type:        "string",
														GoName:      "Value",
														Description: "The value of the label must match the specified value.\n\nvalue can have a maximum length of 1024 characters.",
														Immutable:   true,
													},
												},
											},
										},
										"filterMatchCriteria": &dcl.Property{
											Type:        "string",
											GoName:      "FilterMatchCriteria",
											GoType:      "ForwardingRuleMetadataFilterFilterMatchCriteriaEnum",
											Description: "Specifies how individual `filterLabel` matches within the list of `filterLabel`s contribute towards the overall `metadataFilter` match.\n\nSupported values are:\n\n*   MATCH_ANY: At least one of the `filterLabels` must have a matching label in the provided metadata.\n*   MATCH_ALL: All `filterLabels` must have matching labels in the provided metadata. Possible values: NOT_SET, MATCH_ALL, MATCH_ANY",
											Immutable:   true,
											Enum: []string{
												"NOT_SET",
												"MATCH_ALL",
												"MATCH_ANY",
											},
										},
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the resource; provided by the client when the resource is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the name must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
								Immutable:   true,
							},
							"network": &dcl.Property{
								Type:          "string",
								GoName:        "Network",
								Description:   "This field is not used for external load balancing. For `INTERNAL` and `INTERNAL_SELF_MANAGED` load balancing, this field identifies the network that the load balanced IP should belong to for this Forwarding Rule. If this field is not specified, the default network will be used.",
								Immutable:     true,
								ServerDefault: true,
								HasLongForm:   true,
							},
							"networkTier": &dcl.Property{
								Type:          "string",
								GoName:        "NetworkTier",
								GoType:        "ForwardingRuleNetworkTierEnum",
								Description:   "This signifies the networking tier used for configuring this load balancer and can only take the following values: `PREMIUM`, `STANDARD`. For regional ForwardingRule, the valid values are `PREMIUM` and `STANDARD`. For GlobalForwardingRule, the valid value is `PREMIUM`. If this field is not specified, it is assumed to be `PREMIUM`. If `IPAddress` is specified, this value must be equal to the networkTier of the Address.",
								Immutable:     true,
								ServerDefault: true,
								Enum: []string{
									"PREMIUM",
									"STANDARD",
								},
							},
							"portRange": &dcl.Property{
								Type:        "string",
								GoName:      "PortRange",
								Description: "When the load balancing scheme is `EXTERNAL`, `INTERNAL_SELF_MANAGED` and `INTERNAL_MANAGED`, you can specify a `port_range`. Use with a forwarding rule that points to a target proxy or a target pool. Do not use with a forwarding rule that points to a backend service. This field is used along with the `target` field for TargetHttpProxy, TargetHttpsProxy, TargetSslProxy, TargetTcpProxy, TargetVpnGateway, TargetPool, TargetInstance. Applicable only when `IPProtocol` is `TCP`, `UDP`, or `SCTP`, only packets addressed to ports in the specified range will be forwarded to `target`. Forwarding rules with the same `[IPAddress, IPProtocol]` pair must have disjoint port ranges. Some types of forwarding target have constraints on the acceptable ports:\n\n*   TargetHttpProxy: 80, 8080\n*   TargetHttpsProxy: 443\n*   TargetTcpProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1688, 1883, 5222\n*   TargetSslProxy: 25, 43, 110, 143, 195, 443, 465, 587, 700, 993, 995, 1688, 1883, 5222\n*   TargetVpnGateway: 500, 4500\n\n@pattern: d+(?:-d+)?",
								Immutable:   true,
							},
							"ports": &dcl.Property{
								Type:        "array",
								GoName:      "Ports",
								Description: "This field is used along with the `backend_service` field for internal load balancing. When the load balancing scheme is `INTERNAL`, a list of ports can be configured, for example, ['80'], ['8000','9000']. Only packets addressed to these ports are forwarded to the backends configured with the forwarding rule. If the forwarding rule's loadBalancingScheme is INTERNAL, you can specify ports in one of the following ways: * A list of up to five ports, which can be non-contiguous * Keyword `ALL`, which causes the forwarding rule to forward traffic on any port of the forwarding rule's protocol. @pattern: d+(?:-d+)? For more information, refer to [Port specifications](/load-balancing/docs/forwarding-rule-concepts#port_specifications).",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "set",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project this resource belongs in.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
								Parameter: true,
							},
							"pscConnectionId": &dcl.Property{
								Type:        "string",
								GoName:      "PscConnectionId",
								ReadOnly:    true,
								Description: "The PSC connection id of the PSC Forwarding Rule.",
								Immutable:   true,
							},
							"pscConnectionStatus": &dcl.Property{
								Type:        "string",
								GoName:      "PscConnectionStatus",
								GoType:      "ForwardingRulePscConnectionStatusEnum",
								ReadOnly:    true,
								Description: "The PSC connection status of the PSC Forwarding Rule. Possible values: STATUS_UNSPECIFIED, PENDING, ACCEPTED, REJECTED, CLOSED",
								Immutable:   true,
								Enum: []string{
									"STATUS_UNSPECIFIED",
									"PENDING",
									"ACCEPTED",
									"REJECTED",
									"CLOSED",
								},
							},
							"region": &dcl.Property{
								Type:        "string",
								GoName:      "Region",
								ReadOnly:    true,
								Description: "[Output Only] URL of the region where the regional forwarding rule resides. This field is not applicable to global forwarding rules. You must specify this field as part of the HTTP request URL. It is not settable as a field in the request body.",
								Immutable:   true,
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "[Output Only] Server-defined URL for the resource.",
								Immutable:   true,
							},
							"serviceDirectoryRegistrations": &dcl.Property{
								Type:          "array",
								GoName:        "ServiceDirectoryRegistrations",
								Description:   "Service Directory resources to register this forwarding rule with. Currently, only supports a single Service Directory resource.",
								Immutable:     true,
								ServerDefault: true,
								SendEmpty:     true,
								ListType:      "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "ForwardingRuleServiceDirectoryRegistrations",
									Properties: map[string]*dcl.Property{
										"namespace": &dcl.Property{
											Type:          "string",
											GoName:        "Namespace",
											Description:   "Service Directory namespace to register the forwarding rule under.",
											Immutable:     true,
											ServerDefault: true,
										},
										"service": &dcl.Property{
											Type:        "string",
											GoName:      "Service",
											Description: "Service Directory service to register the forwarding rule under.",
											Immutable:   true,
										},
									},
								},
							},
							"serviceLabel": &dcl.Property{
								Type:        "string",
								GoName:      "ServiceLabel",
								Description: "An optional prefix to the service name for this Forwarding Rule. If specified, the prefix is the first label of the fully qualified service name. The label must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the label must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash. This field is only used for internal load balancing.",
								Immutable:   true,
							},
							"serviceName": &dcl.Property{
								Type:        "string",
								GoName:      "ServiceName",
								ReadOnly:    true,
								Description: "[Output Only] The internal fully qualified service name for this Forwarding Rule. This field is only used for internal load balancing.",
								Immutable:   true,
							},
							"sourceIPRanges": &dcl.Property{
								Type:        "array",
								GoName:      "SourceIPRanges",
								Description: "If not empty, this Forwarding Rule will only forward the traffic when the source IP address matches one of the IP addresses or CIDR ranges set here. Note that a Forwarding Rule can only have up to 64 source IP ranges, and this field can only be used with a regional Forwarding Rule whose scheme is EXTERNAL. Each sourceIpRange entry should be either an IP address (for example, 1.2.3.4) or a CIDR range (for example, 1.2.3.0/24).",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"subnetwork": &dcl.Property{
								Type:          "string",
								GoName:        "Subnetwork",
								Description:   "This field is only used for `INTERNAL` load balancing. For internal load balancing, this field identifies the subnetwork that the load balanced IP should belong to for this Forwarding Rule. If the network specified is in auto subnet mode, this field is optional. However, if the network is in custom subnet mode, a subnetwork must be specified.",
								Immutable:     true,
								ServerDefault: true,
								HasLongForm:   true,
							},
							"target": &dcl.Property{
								Type:        "string",
								GoName:      "Target",
								Description: "The URL of the target resource to receive the matched traffic. For regional forwarding rules, this target must live in the same region as the forwarding rule. For global forwarding rules, this target must be a global load balancing resource. The forwarded traffic must be of a type appropriate to the target object. For `INTERNAL_SELF_MANAGED` load balancing, only `targetHttpProxy` is valid, not `targetHttpsProxy`.",
							},
						},
					},
				},
			},
		},
	}
}
