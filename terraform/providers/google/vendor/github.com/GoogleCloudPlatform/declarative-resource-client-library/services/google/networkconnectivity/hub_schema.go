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
package networkconnectivity

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLHubSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "NetworkConnectivity/Hub",
			Description: "The NetworkConnectivity Hub resource",
			StructName:  "Hub",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Hub",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "hub",
						Required:    true,
						Description: "A full instance of a Hub",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Hub",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "hub",
						Required:    true,
						Description: "A full instance of a Hub",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Hub",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "hub",
						Required:    true,
						Description: "A full instance of a Hub",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Hub",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Hub",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
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
				"Hub": &dcl.Component{
					Title: "Hub",
					ID:    "projects/{{project}}/locations/global/hubs/{{name}}",
					Locations: []string{
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
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time the hub was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of the hub.",
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).",
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Immutable. The name of the hub. Hub names must be unique. They use the following form: `projects/{project_number}/locations/global/hubs/{hub_id}`",
								Immutable:   true,
								HasLongForm: true,
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
							"routingVpcs": &dcl.Property{
								Type:        "array",
								GoName:      "RoutingVpcs",
								ReadOnly:    true,
								Description: "The VPC network associated with this hub's spokes. All of the VPN tunnels, VLAN attachments, and router appliance instances referenced by this hub's spokes must belong to this VPC network. This field is read-only. Network Connectivity Center automatically populates it based on the set of spokes attached to the hub.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "HubRoutingVpcs",
									Properties: map[string]*dcl.Property{
										"uri": &dcl.Property{
											Type:        "string",
											GoName:      "Uri",
											Description: "The URI of the VPC network.",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/Network",
													Field:    "selfLink",
												},
											},
										},
									},
								},
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "HubStateEnum",
								ReadOnly:    true,
								Description: "Output only. The current lifecycle state of this hub. Possible values: STATE_UNSPECIFIED, CREATING, ACTIVE, DELETING",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"CREATING",
									"ACTIVE",
									"DELETING",
								},
							},
							"uniqueId": &dcl.Property{
								Type:        "string",
								GoName:      "UniqueId",
								ReadOnly:    true,
								Description: "Output only. The Google-generated UUID for the hub. This value is unique across all hub resources. If a hub is deleted and another with the same name is created, the new hub is assigned a different unique_id.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time the hub was last updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
