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

func DCLFirewallPolicySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/FirewallPolicy",
			Description: "The Compute FirewallPolicy resource",
			StructName:  "FirewallPolicy",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a FirewallPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicy",
						Required:    true,
						Description: "A full instance of a FirewallPolicy",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a FirewallPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicy",
						Required:    true,
						Description: "A full instance of a FirewallPolicy",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a FirewallPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicy",
						Required:    true,
						Description: "A full instance of a FirewallPolicy",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all FirewallPolicy",
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
				Description: "The function used to list information about many FirewallPolicy",
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
				"FirewallPolicy": &dcl.Component{
					Title: "FirewallPolicy",
					ID:    "locations/global/firewallPolicies/{{name}}",
					Locations: []string{
						"global",
					},
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"shortName",
							"parent",
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
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "Name of the resource. It is a numeric ID allocated by GCP which uniquely identifies the Firewall Policy.",
								Immutable:                true,
								ServerGeneratedParameter: true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "The parent of the firewall policy.",
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
								},
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
							"shortName": &dcl.Property{
								Type:        "string",
								GoName:      "ShortName",
								Description: "User-provided name of the Organization firewall policy. The name should be unique in the organization in which the firewall policy is created. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression [a-z]([-a-z0-9]*[a-z0-9])? which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
