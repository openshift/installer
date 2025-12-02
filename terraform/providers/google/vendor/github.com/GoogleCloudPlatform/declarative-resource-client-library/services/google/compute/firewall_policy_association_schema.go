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

func DCLFirewallPolicyAssociationSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/FirewallPolicyAssociation",
			Description: "The Compute FirewallPolicyAssociation resource",
			StructName:  "FirewallPolicyAssociation",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a FirewallPolicyAssociation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicyAssociation",
						Required:    true,
						Description: "A full instance of a FirewallPolicyAssociation",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a FirewallPolicyAssociation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicyAssociation",
						Required:    true,
						Description: "A full instance of a FirewallPolicyAssociation",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a FirewallPolicyAssociation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "firewallPolicyAssociation",
						Required:    true,
						Description: "A full instance of a FirewallPolicyAssociation",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all FirewallPolicyAssociation",
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
				Description: "The function used to list information about many FirewallPolicyAssociation",
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
				"FirewallPolicyAssociation": &dcl.Component{
					Title: "FirewallPolicyAssociation",
					ID:    "locations/global/firewallPolicies/{{firewall_policy}}/associations/{{name}}",
					Locations: []string{
						"global",
					},
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"attachmentTarget",
							"firewallPolicy",
						},
						Properties: map[string]*dcl.Property{
							"attachmentTarget": &dcl.Property{
								Type:        "string",
								GoName:      "AttachmentTarget",
								Description: "The target that the firewall policy is attached to.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Folder",
										Field:    "name",
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Organization",
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
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name for an association.",
								Immutable:   true,
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
