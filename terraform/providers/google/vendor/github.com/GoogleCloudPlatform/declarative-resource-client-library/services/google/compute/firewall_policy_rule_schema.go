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

func DCLFirewallPolicyRuleSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/FirewallPolicyRule",
			Description: "The Compute FirewallPolicyRule resource",
			StructName:  "FirewallPolicyRule",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a FirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicyRule",
						Required:    true,
						Description: "A full instance of a FirewallPolicyRule",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a FirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicyRule",
						Required:    true,
						Description: "A full instance of a FirewallPolicyRule",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a FirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicyRule",
						Required:    true,
						Description: "A full instance of a FirewallPolicyRule",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all FirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "firewallPolicy",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many FirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "firewallPolicy",
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
				"FirewallPolicyRule": &dcl.Component{
					Title: "FirewallPolicyRule",
					ID:    "locations/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}",
					Locations: []string{
						"global",
					},
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"priority",
							"match",
							"action",
							"direction",
							"firewallPolicy",
						},
						Properties: map[string]*dcl.Property{
							"action": &dcl.Property{
								Type:        "string",
								GoName:      "Action",
								Description: "The Action to perform when the client connection triggers the rule. Valid actions are \"allow\", \"deny\", \"goto_next\" and \"apply_security_profile_group\".",
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description for this resource.",
							},
							"direction": &dcl.Property{
								Type:        "string",
								GoName:      "Direction",
								GoType:      "FirewallPolicyRuleDirectionEnum",
								Description: "The direction in which this rule applies. Possible values: INGRESS, EGRESS",
								Enum: []string{
									"INGRESS",
									"EGRESS",
								},
							},
							"disabled": &dcl.Property{
								Type:        "boolean",
								GoName:      "Disabled",
								Description: "Denotes whether the firewall policy rule is disabled. When set to true, the firewall policy rule is not enforced and traffic behaves as if it did not exist. If this is unspecified, the firewall policy rule will be enabled.",
							},
							"enableLogging": &dcl.Property{
								Type:        "boolean",
								GoName:      "EnableLogging",
								Description: "Denotes whether to enable logging for a particular rule. If logging is enabled, logs will be exported to the configured export destination in Stackdriver. Logs may be exported to BigQuery or Pub/Sub. Note: you cannot enable logging on \"goto_next\" rules.",
							},
							"firewallPolicy": &dcl.Property{
								Type:        "string",
								GoName:      "FirewallPolicy",
								Description: "The firewall policy of the resource.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/FirewallPolicy",
										Field:    "name",
										Parent:   true,
									},
								},
								HasLongForm: true,
							},
							"kind": &dcl.Property{
								Type:        "string",
								GoName:      "Kind",
								ReadOnly:    true,
								Description: "Type of the resource. Always `compute#firewallPolicyRule` for firewall policy rules",
								Immutable:   true,
							},
							"match": &dcl.Property{
								Type:        "object",
								GoName:      "Match",
								GoType:      "FirewallPolicyRuleMatch",
								Description: "A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding 'action' is enforced.",
								Required: []string{
									"layer4Configs",
								},
								Properties: map[string]*dcl.Property{
									"destAddressGroups": &dcl.Property{
										Type:        "array",
										GoName:      "DestAddressGroups",
										Description: "Address groups which should be matched against the traffic destination. Maximum number of destination address groups is 10. Destination address groups is only supported in Egress rules.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"destFqdns": &dcl.Property{
										Type:        "array",
										GoName:      "DestFqdns",
										Description: "Domain names that will be used to match against the resolved domain name of destination of traffic. Can only be specified if DIRECTION is egress.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"destIPRanges": &dcl.Property{
										Type:        "array",
										GoName:      "DestIPRanges",
										Description: "CIDR IP address range. Maximum number of destination CIDR IP ranges allowed is 256.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"destRegionCodes": &dcl.Property{
										Type:        "array",
										GoName:      "DestRegionCodes",
										Description: "The Unicode country codes whose IP addresses will be used to match against the source of traffic. Can only be specified if DIRECTION is egress.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"destThreatIntelligences": &dcl.Property{
										Type:        "array",
										GoName:      "DestThreatIntelligences",
										Description: "Name of the Google Cloud Threat Intelligence list.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"layer4Configs": &dcl.Property{
										Type:        "array",
										GoName:      "Layer4Configs",
										Description: "Pairs of IP protocols and ports that the rule should match.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "FirewallPolicyRuleMatchLayer4Configs",
											Required: []string{
												"ipProtocol",
											},
											Properties: map[string]*dcl.Property{
												"ipProtocol": &dcl.Property{
													Type:        "string",
													GoName:      "IPProtocol",
													Description: "The IP protocol to which this rule applies. The protocol type is required when creating a firewall rule. This value can either be one of the following well known protocol strings (`tcp`, `udp`, `icmp`, `esp`, `ah`, `ipip`, `sctp`), or the IP protocol number.",
												},
												"ports": &dcl.Property{
													Type:        "array",
													GoName:      "Ports",
													Description: "An optional list of ports to which this rule applies. This field is only applicable for UDP or TCP protocol. Each entry must be either an integer or a range. If not specified, this rule applies to connections through any port. Example inputs include: ``.",
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
											},
										},
									},
									"srcAddressGroups": &dcl.Property{
										Type:        "array",
										GoName:      "SrcAddressGroups",
										Description: "Address groups which should be matched against the traffic source. Maximum number of source address groups is 10. Source address groups is only supported in Ingress rules.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"srcFqdns": &dcl.Property{
										Type:        "array",
										GoName:      "SrcFqdns",
										Description: "Domain names that will be used to match against the resolved domain name of source of traffic. Can only be specified if DIRECTION is ingress.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"srcIPRanges": &dcl.Property{
										Type:        "array",
										GoName:      "SrcIPRanges",
										Description: "CIDR IP address range. Maximum number of source CIDR IP ranges allowed is 256.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"srcRegionCodes": &dcl.Property{
										Type:        "array",
										GoName:      "SrcRegionCodes",
										Description: "The Unicode country codes whose IP addresses will be used to match against the source of traffic. Can only be specified if DIRECTION is ingress.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"srcThreatIntelligences": &dcl.Property{
										Type:        "array",
										GoName:      "SrcThreatIntelligences",
										Description: "Name of the Google Cloud Threat Intelligence list.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
								},
							},
							"priority": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Priority",
								Description: "An integer indicating the priority of a rule in the list. The priority must be a positive value between 0 and 2147483647. Rules are evaluated from highest to lowest priority where 0 is the highest priority and 2147483647 is the lowest prority.",
								Immutable:   true,
							},
							"ruleTupleCount": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "RuleTupleCount",
								ReadOnly:    true,
								Description: "Calculation of the complexity of a single firewall policy rule.",
							},
							"securityProfileGroup": &dcl.Property{
								Type:        "string",
								GoName:      "SecurityProfileGroup",
								Description: "A fully-qualified URL of a SecurityProfileGroup resource. Example: https://networksecurity.googleapis.com/v1/organizations/{organizationId}/locations/global/securityProfileGroups/my-security-profile-group. It must be specified if action = 'apply_security_profile_group' and cannot be specified for other actions.",
							},
							"targetResources": &dcl.Property{
								Type:        "array",
								GoName:      "TargetResources",
								Description: "A list of network resource URLs to which this rule applies. This field allows you to control which network's VMs get this rule. If this field is left blank, all VMs within the organization will receive the rule.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
									ResourceReferences: []*dcl.PropertyResourceReference{
										&dcl.PropertyResourceReference{
											Resource: "Compute/Network",
											Field:    "selfLink",
										},
									},
								},
							},
							"targetServiceAccounts": &dcl.Property{
								Type:        "array",
								GoName:      "TargetServiceAccounts",
								Description: "A list of service accounts indicating the sets of instances that are applied with this rule.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
									ResourceReferences: []*dcl.PropertyResourceReference{
										&dcl.PropertyResourceReference{
											Resource: "Iam/ServiceAccount",
											Field:    "name",
										},
									},
								},
							},
							"tlsInspect": &dcl.Property{
								Type:        "boolean",
								GoName:      "TlsInspect",
								Description: "Boolean flag indicating if the traffic should be TLS decrypted. It can be set only if action = 'apply_security_profile_group' and cannot be set for other actions.",
							},
						},
					},
				},
			},
		},
	}
}
