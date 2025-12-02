// Copyright 2023 Google LLC. All Rights Reserved.
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

func DCLNetworkFirewallPolicyRuleSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/NetworkFirewallPolicyRule",
			Description: "The Compute NetworkFirewallPolicyRule resource",
			StructName:  "NetworkFirewallPolicyRule",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a NetworkFirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicyRule",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicyRule",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a NetworkFirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicyRule",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicyRule",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a NetworkFirewallPolicyRule",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicyRule",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicyRule",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all NetworkFirewallPolicyRule",
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
				Description: "The function used to list information about many NetworkFirewallPolicyRule",
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
				"NetworkFirewallPolicyRule": &dcl.Component{
					Title: "NetworkFirewallPolicyRule",
					ID:    "projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/rules/{{priority}}",
					Locations: []string{
						"global",
						"region",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"priority",
							"match",
							"action",
							"direction",
							"firewallPolicy",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"action": &dcl.Property{
								Type:        "string",
								GoName:      "Action",
								Description: "The Action to perform when the client connection triggers the rule. Can currently be either \"allow\" or \"deny()\" where valid values for status are 403, 404, and 502.",
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description for this resource.",
							},
							"direction": &dcl.Property{
								Type:        "string",
								GoName:      "Direction",
								GoType:      "NetworkFirewallPolicyRuleDirectionEnum",
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
										Resource: "Compute/NetworkFirewallPolicy",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"kind": &dcl.Property{
								Type:        "string",
								GoName:      "Kind",
								ReadOnly:    true,
								Description: "Type of the resource. Always `compute#firewallPolicyRule` for firewall policy rules",
								Immutable:   true,
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location of this resource.",
								Immutable:   true,
							},
							"match": &dcl.Property{
								Type:        "object",
								GoName:      "Match",
								GoType:      "NetworkFirewallPolicyRuleMatch",
								Description: "A match condition that incoming traffic is evaluated against. If it evaluates to true, the corresponding 'action' is enforced.",
								Required: []string{
									"layer4Configs",
								},
								Properties: map[string]*dcl.Property{
									"destIPRanges": &dcl.Property{
										Type:        "array",
										GoName:      "DestIPRanges",
										Description: "CIDR IP address range. Maximum number of destination CIDR IP ranges allowed is 5000.",
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
											GoType: "NetworkFirewallPolicyRuleMatchLayer4Configs",
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
									"srcIPRanges": &dcl.Property{
										Type:        "array",
										GoName:      "SrcIPRanges",
										Description: "CIDR IP address range. Maximum number of source CIDR IP ranges allowed is 5000.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"srcSecureTags": &dcl.Property{
										Type:        "array",
										GoName:      "SrcSecureTags",
										Description: "List of secure tag values, which should be matched at the source of the traffic. For INGRESS rule, if all the <code>srcSecureTag</code> are INEFFECTIVE, and there is no <code>srcIpRange</code>, this rule will be ignored. Maximum number of source tag values allowed is 256.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "NetworkFirewallPolicyRuleMatchSrcSecureTags",
											Required: []string{
												"name",
											},
											Properties: map[string]*dcl.Property{
												"name": &dcl.Property{
													Type:        "string",
													GoName:      "Name",
													Description: "Name of the secure tag, created with TagManager's TagValue API. @pattern tagValues/[0-9]+",
													ResourceReferences: []*dcl.PropertyResourceReference{
														&dcl.PropertyResourceReference{
															Resource: "Cloudresourcemanager/TagValue",
															Field:    "namespacedName",
														},
													},
												},
												"state": &dcl.Property{
													Type:        "string",
													GoName:      "State",
													GoType:      "NetworkFirewallPolicyRuleMatchSrcSecureTagsStateEnum",
													ReadOnly:    true,
													Description: "[Output Only] State of the secure tag, either `EFFECTIVE` or `INEFFECTIVE`. A secure tag is `INEFFECTIVE` when it is deleted or its network is deleted.",
													Enum: []string{
														"EFFECTIVE",
														"INEFFECTIVE",
													},
												},
											},
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
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"ruleName": &dcl.Property{
								Type:        "string",
								GoName:      "RuleName",
								Description: "An optional name for the rule. This field is not a unique identifier and can be updated.",
							},
							"ruleTupleCount": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "RuleTupleCount",
								ReadOnly:    true,
								Description: "Calculation of the complexity of a single firewall policy rule.",
								Immutable:   true,
							},
							"targetSecureTags": &dcl.Property{
								Type:        "array",
								GoName:      "TargetSecureTags",
								Description: "A list of secure tags that controls which instances the firewall rule applies to. If <code>targetSecureTag</code> are specified, then the firewall rule applies only to instances in the VPC network that have one of those EFFECTIVE secure tags, if all the target_secure_tag are in INEFFECTIVE state, then this rule will be ignored. <code>targetSecureTag</code> may not be set at the same time as <code>targetServiceAccounts</code>. If neither <code>targetServiceAccounts</code> nor <code>targetSecureTag</code> are specified, the firewall rule applies to all instances on the specified network. Maximum number of target label tags allowed is 256.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "NetworkFirewallPolicyRuleTargetSecureTags",
									Required: []string{
										"name",
									},
									Properties: map[string]*dcl.Property{
										"name": &dcl.Property{
											Type:        "string",
											GoName:      "Name",
											Description: "Name of the secure tag, created with TagManager's TagValue API. @pattern tagValues/[0-9]+",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Cloudresourcemanager/TagValue",
													Field:    "namespacedName",
												},
											},
										},
										"state": &dcl.Property{
											Type:        "string",
											GoName:      "State",
											GoType:      "NetworkFirewallPolicyRuleTargetSecureTagsStateEnum",
											ReadOnly:    true,
											Description: "[Output Only] State of the secure tag, either `EFFECTIVE` or `INEFFECTIVE`. A secure tag is `INEFFECTIVE` when it is deleted or its network is deleted.",
											Enum: []string{
												"EFFECTIVE",
												"INEFFECTIVE",
											},
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
						},
					},
				},
			},
		},
	}
}
