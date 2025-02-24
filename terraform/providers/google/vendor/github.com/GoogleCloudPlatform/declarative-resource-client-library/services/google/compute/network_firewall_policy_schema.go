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

func DCLNetworkFirewallPolicySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/NetworkFirewallPolicy",
			Description: "The Compute NetworkFirewallPolicy resource",
			StructName:  "NetworkFirewallPolicy",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a NetworkFirewallPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicy",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicy",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a NetworkFirewallPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicy",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicy",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a NetworkFirewallPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicy",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicy",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all NetworkFirewallPolicy",
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
				Description: "The function used to list information about many NetworkFirewallPolicy",
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
				"NetworkFirewallPolicy": &dcl.Component{
					Title: "NetworkFirewallPolicy",
					ID:    "projects/{{project}}/global/firewallPolicies/{{name}}",
					Locations: []string{
						"region",
						"global",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"creationTimestamp": &dcl.Property{
								Type:        "string",
								GoName:      "CreationTimestamp",
								ReadOnly:    true,
								Description: "Creation timestamp in RFC3339 text format.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. Provide this property when you create the resource.",
							},
							"fingerprint": &dcl.Property{
								Type:        "string",
								GoName:      "Fingerprint",
								ReadOnly:    true,
								Description: "Fingerprint of the resource. This field is used internally during updates of this resource.",
								Immutable:   true,
							},
							"id": &dcl.Property{
								Type:        "string",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "The unique identifier for the resource. This identifier is defined by the server.",
								Immutable:   true,
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location of this resource.",
								Immutable:   true,
								Parameter:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "User-provided name of the Network firewall policy. The name should be unique in the project in which the firewall policy is created. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
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
								Parameter: true,
							},
							"region": &dcl.Property{
								Type:        "string",
								GoName:      "Region",
								ReadOnly:    true,
								Description: "[Output Only] URL of the region where the regional firewall policy resides. This field is not applicable to global firewall policies. You must specify this field as part of the HTTP request URL. It is not settable as a field in the request body.",
								Immutable:   true,
							},
							"ruleTupleCount": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "RuleTupleCount",
								ReadOnly:    true,
								Description: "Total count of all firewall policy rule tuples. A firewall policy can not exceed a set number of tuples.",
								Immutable:   true,
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "Server-defined URL for the resource.",
								Immutable:   true,
							},
							"selfLinkWithId": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLinkWithId",
								ReadOnly:    true,
								Description: "Server-defined URL for this resource with the resource id.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
