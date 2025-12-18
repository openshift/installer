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

func DCLNetworkSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/Network",
			Description: "The Compute Network resource",
			StructName:  "Network",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Network",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "network",
						Required:    true,
						Description: "A full instance of a Network",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Network",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "network",
						Required:    true,
						Description: "A full instance of a Network",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Network",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "network",
						Required:    true,
						Description: "A full instance of a Network",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Network",
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
				Description: "The function used to list information about many Network",
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
				"Network": &dcl.Component{
					Title: "Network",
					ID:    "projects/{{project}}/global/networks/{{name}}",
					Locations: []string{
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
							"autoCreateSubnetworks": &dcl.Property{
								Type:          "boolean",
								GoName:        "AutoCreateSubnetworks",
								Description:   "When set to `true`, the network is created in \"auto subnet mode\" and it will create a subnet for each region automatically across the `10.128.0.0/9` address range.  When set to `false`, the network is created in \"custom subnet mode\" so the user can explicitly connect subnetwork resources. ",
								Immutable:     true,
								Default:       true,
								ServerDefault: true,
								SendEmpty:     true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. The resource must be recreated to modify this field. ",
								Immutable:   true,
							},
							"gatewayIPv4": &dcl.Property{
								Type:        "string",
								GoName:      "GatewayIPv4",
								ReadOnly:    true,
								Description: "The gateway address for default routing out of the network. This value is selected by GCP. ",
								Immutable:   true,
							},
							"mtu": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Mtu",
								Description: "Maximum Transmission Unit in bytes. The minimum value for this field is 1460 and the maximum value is 1500 bytes.",
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash. ",
								Immutable:   true,
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project id of the resource.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"routingConfig": &dcl.Property{
								Type:          "object",
								GoName:        "RoutingConfig",
								GoType:        "NetworkRoutingConfig",
								Description:   "The network-level routing configuration for this network. Used by Cloud Router to determine what type of network-wide routing behavior to enforce. ",
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"routingMode": &dcl.Property{
										Type:          "string",
										GoName:        "RoutingMode",
										GoType:        "NetworkRoutingConfigRoutingModeEnum",
										Description:   "The network-wide routing mode to use. If set to `REGIONAL`, this network's cloud routers will only advertise routes with subnetworks of this network in the same region as the router. If set to `GLOBAL`, this network's cloud routers will advertise routes with all subnetworks of this network, across regions. ",
										ServerDefault: true,
										Enum: []string{
											"REGIONAL",
											"GLOBAL",
										},
									},
								},
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
								Description: "Server-defined URL for the resource containing the network ID.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
