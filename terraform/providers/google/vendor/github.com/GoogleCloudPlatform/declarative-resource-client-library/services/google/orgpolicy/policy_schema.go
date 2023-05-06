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
package orgpolicy

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLPolicySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "OrgPolicy/Policy",
			Description: "An organization policy gives you programmatic control over your organization's cloud resources.  Using Organization Policies, you will be able to configure constraints across your entire resource hierarchy.",
			StructName:  "Policy",
			Reference: &dcl.Link{
				Text: "REST API",
				URL:  "https://cloud.google.com/resource-manager/docs/reference/orgpolicy/rest/v2/organizations.policies",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Understanding Org Policy concepts",
					URL:  "https://cloud.google.com/resource-manager/docs/organization-policy/overview",
				},
				&dcl.Link{
					Text: "The resource hierarchy",
					URL:  "https://cloud.google.com/resource-manager/docs/cloud-platform-resource-hierarchy",
				},
				&dcl.Link{
					Text: "All valid constraints",
					URL:  "https://cloud.google.com/resource-manager/docs/organization-policy/org-policy-constraints",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Policy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "policy",
						Required:    true,
						Description: "A full instance of a Policy",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Policy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "policy",
						Required:    true,
						Description: "A full instance of a Policy",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Policy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "policy",
						Required:    true,
						Description: "A full instance of a Policy",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Policy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "parent",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Policy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "parent",
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
				"Policy": &dcl.Component{
					Title:     "Policy",
					ID:        "{{parent}}/policies/{{name}}",
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"parent",
						},
						Properties: map[string]*dcl.Property{
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Immutable. The resource name of the Policy. Must be one of the following forms, where constraint_name is the name of the constraint which this Policy configures: * `projects/{project_number}/policies/{constraint_name}` * `folders/{folder_id}/policies/{constraint_name}` * `organizations/{organization_id}/policies/{constraint_name}` For example, \"projects/123/policies/compute.disableSerialPortAccess\". Note: `projects/{project_id}/policies/{constraint_name}` is also an acceptable name for API requests, but responses will return the name using the equivalent project number.",
								Immutable:   true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "The parent of the resource.",
								Immutable:           true,
								ForwardSlashAllowed: true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Folder",
										Field:    "name",
										Parent:   true,
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Organization",
										Field:    "name",
										Parent:   true,
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"spec": &dcl.Property{
								Type:        "object",
								GoName:      "Spec",
								GoType:      "PolicySpec",
								Description: "Basic information about the Organization Policy.",
								Properties: map[string]*dcl.Property{
									"etag": &dcl.Property{
										Type:        "string",
										GoName:      "Etag",
										ReadOnly:    true,
										Description: "An opaque tag indicating the current version of the `Policy`, used for concurrency control. This field is ignored if used in a `CreatePolicy` request. When the `Policy` is returned from either a `GetPolicy` or a `ListPolicies` request, this `etag` indicates the version of the current `Policy` to use when executing a read-modify-write loop. When the `Policy` is returned from a `GetEffectivePolicy` request, the `etag` will be unset.",
									},
									"inheritFromParent": &dcl.Property{
										Type:        "boolean",
										GoName:      "InheritFromParent",
										Description: "Determines the inheritance behavior for this `Policy`. If `inherit_from_parent` is true, PolicyRules set higher up in the hierarchy (up to the closest root) are inherited and present in the effective policy. If it is false, then no rules are inherited, and this Policy becomes the new root for evaluation. This field can be set only for Policies which configure list constraints.",
									},
									"reset": &dcl.Property{
										Type:        "boolean",
										GoName:      "Reset",
										Description: "Ignores policies set above this resource and restores the `constraint_default` enforcement behavior of the specific `Constraint` at this resource. This field can be set in policies for either list or boolean constraints. If set, `rules` must be empty and `inherit_from_parent` must be set to false.",
									},
									"rules": &dcl.Property{
										Type:        "array",
										GoName:      "Rules",
										Description: "Up to 10 PolicyRules are allowed. In Policies for boolean constraints, the following requirements apply: - There must be one and only one PolicyRule where condition is unset. - BooleanPolicyRules with conditions must set `enforced` to the opposite of the PolicyRule without a condition. - During policy evaluation, PolicyRules with conditions that are true for a target resource take precedence.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "PolicySpecRules",
											Properties: map[string]*dcl.Property{
												"allowAll": &dcl.Property{
													Type:        "boolean",
													GoName:      "AllowAll",
													Description: "Setting this to true means that all values are allowed. This field can be set only in Policies for list constraints.",
													Conflicts: []string{
														"values",
														"denyAll",
														"enforce",
													},
												},
												"condition": &dcl.Property{
													Type:        "object",
													GoName:      "Condition",
													GoType:      "PolicySpecRulesCondition",
													Description: "A condition which determines whether this rule is used in the evaluation of the policy. When set, the `expression` field in the `Expr' must include from 1 to 10 subexpressions, joined by the \"||\" or \"&&\" operators. Each subexpression must be of the form \"resource.matchTag('/tag_key_short_name, 'tag_value_short_name')\". or \"resource.matchTagId('tagKeys/key_id', 'tagValues/value_id')\". where key_name and value_name are the resource names for Label Keys and Values. These names are available from the Tag Manager Service. An example expression is: \"resource.matchTag('123456789/environment, 'prod')\". or \"resource.matchTagId('tagKeys/123', 'tagValues/456')\".",
													Properties: map[string]*dcl.Property{
														"description": &dcl.Property{
															Type:        "string",
															GoName:      "Description",
															Description: "Optional. Description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.",
														},
														"expression": &dcl.Property{
															Type:        "string",
															GoName:      "Expression",
															Description: "Textual representation of an expression in Common Expression Language syntax.",
														},
														"location": &dcl.Property{
															Type:        "string",
															GoName:      "Location",
															Description: "Optional. String indicating the location of the expression for error reporting, e.g. a file name and a position in the file.",
														},
														"title": &dcl.Property{
															Type:        "string",
															GoName:      "Title",
															Description: "Optional. Title for the expression, i.e. a short string describing its purpose. This can be used e.g. in UIs which allow to enter the expression.",
														},
													},
												},
												"denyAll": &dcl.Property{
													Type:        "boolean",
													GoName:      "DenyAll",
													Description: "Setting this to true means that all values are denied. This field can be set only in Policies for list constraints.",
													Conflicts: []string{
														"values",
														"allowAll",
														"enforce",
													},
												},
												"enforce": &dcl.Property{
													Type:        "boolean",
													GoName:      "Enforce",
													Description: "If `true`, then the `Policy` is enforced. If `false`, then any configuration is acceptable. This field can be set only in Policies for boolean constraints.",
													Conflicts: []string{
														"values",
														"allowAll",
														"denyAll",
													},
												},
												"values": &dcl.Property{
													Type:        "object",
													GoName:      "Values",
													GoType:      "PolicySpecRulesValues",
													Description: "List of values to be used for this PolicyRule. This field can be set only in Policies for list constraints.",
													Conflicts: []string{
														"allowAll",
														"denyAll",
														"enforce",
													},
													Properties: map[string]*dcl.Property{
														"allowedValues": &dcl.Property{
															Type:        "array",
															GoName:      "AllowedValues",
															Description: "List of values allowed at this resource.",
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
														"deniedValues": &dcl.Property{
															Type:        "array",
															GoName:      "DeniedValues",
															Description: "List of values denied at this resource.",
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
										},
									},
									"updateTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "UpdateTime",
										ReadOnly:    true,
										Description: "Output only. The time stamp this was previously updated. This represents the last time a call to `CreatePolicy` or `UpdatePolicy` was made for that `Policy`.",
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
