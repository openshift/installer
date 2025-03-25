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

func DCLNetworkFirewallPolicyAssociationSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/NetworkFirewallPolicyAssociation",
			Description: "The Compute NetworkFirewallPolicyAssociation resource",
			StructName:  "NetworkFirewallPolicyAssociation",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a NetworkFirewallPolicyAssociation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicyAssociation",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicyAssociation",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a NetworkFirewallPolicyAssociation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicyAssociation",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicyAssociation",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a NetworkFirewallPolicyAssociation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "networkFirewallPolicyAssociation",
						Required:    true,
						Description: "A full instance of a NetworkFirewallPolicyAssociation",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all NetworkFirewallPolicyAssociation",
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
				Description: "The function used to list information about many NetworkFirewallPolicyAssociation",
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
				"NetworkFirewallPolicyAssociation": &dcl.Component{
					Title: "NetworkFirewallPolicyAssociation",
					ID:    "projects/{{project}}/global/firewallPolicies/{{firewall_policy}}/getAssociation?name={{name}}",
					Locations: []string{
						"global",
						"region",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"attachmentTarget",
							"firewallPolicy",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"attachmentTarget": &dcl.Property{
								Type:        "string",
								GoName:      "AttachmentTarget",
								Description: "The target that the firewall policy is attached to.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/Network",
										Field:    "name",
									},
								},
							},
							"firewallPolicy": &dcl.Property{
								Type:        "string",
								GoName:      "FirewallPolicy",
								Description: "The firewall policy ID of the association.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/FirewallPolicy",
										Field:    "name",
										Parent:   true,
									},
								},
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
								Description: "The name for an association.",
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
							"shortName": &dcl.Property{
								Type:        "string",
								GoName:      "ShortName",
								ReadOnly:    true,
								Description: "The short name of the firewall policy of the association.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
