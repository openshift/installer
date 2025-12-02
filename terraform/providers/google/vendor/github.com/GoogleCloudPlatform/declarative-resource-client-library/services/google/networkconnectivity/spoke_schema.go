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
package networkconnectivity

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLSpokeSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "NetworkConnectivity/Spoke",
			Description: "The NetworkConnectivity Spoke resource",
			StructName:  "Spoke",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Spoke",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "spoke",
						Required:    true,
						Description: "A full instance of a Spoke",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Spoke",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "spoke",
						Required:    true,
						Description: "A full instance of a Spoke",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Spoke",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "spoke",
						Required:    true,
						Description: "A full instance of a Spoke",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Spoke",
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
				Description: "The function used to list information about many Spoke",
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
				"Spoke": &dcl.Component{
					Title:           "Spoke",
					ID:              "projects/{{project}}/locations/{{location}}/spokes/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"hub",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time the spoke was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of the spoke.",
							},
							"hub": &dcl.Property{
								Type:        "string",
								GoName:      "Hub",
								Description: "Immutable. The URI of the hub that this spoke is attached to.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Networkconnectivity/Hub",
										Field:    "name",
									},
								},
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Optional labels in key:value format. For more information about labels, see [Requirements for labels](https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements).",
							},
							"linkedInterconnectAttachments": &dcl.Property{
								Type:        "object",
								GoName:      "LinkedInterconnectAttachments",
								GoType:      "SpokeLinkedInterconnectAttachments",
								Description: "A collection of VLAN attachment resources. These resources should be redundant attachments that all advertise the same prefixes to Google Cloud. Alternatively, in active/passive configurations, all attachments should be capable of advertising the same prefixes.",
								Immutable:   true,
								Conflicts: []string{
									"linkedVpnTunnels",
									"linkedRouterApplianceInstances",
								},
								Required: []string{
									"uris",
									"siteToSiteDataTransfer",
								},
								Properties: map[string]*dcl.Property{
									"siteToSiteDataTransfer": &dcl.Property{
										Type:        "boolean",
										GoName:      "SiteToSiteDataTransfer",
										Description: "A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.",
										Immutable:   true,
									},
									"uris": &dcl.Property{
										Type:        "array",
										GoName:      "Uris",
										Description: "The URIs of linked interconnect attachment resources",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/InterconnectAttachment",
													Field:    "selfLink",
												},
											},
										},
									},
								},
							},
							"linkedRouterApplianceInstances": &dcl.Property{
								Type:        "object",
								GoName:      "LinkedRouterApplianceInstances",
								GoType:      "SpokeLinkedRouterApplianceInstances",
								Description: "The URIs of linked Router appliance resources",
								Immutable:   true,
								Conflicts: []string{
									"linkedVpnTunnels",
									"linkedInterconnectAttachments",
								},
								Required: []string{
									"instances",
									"siteToSiteDataTransfer",
								},
								Properties: map[string]*dcl.Property{
									"instances": &dcl.Property{
										Type:        "array",
										GoName:      "Instances",
										Description: "The list of router appliance instances",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "SpokeLinkedRouterApplianceInstancesInstances",
											Properties: map[string]*dcl.Property{
												"ipAddress": &dcl.Property{
													Type:        "string",
													GoName:      "IPAddress",
													Description: "The IP address on the VM to use for peering.",
													Immutable:   true,
												},
												"virtualMachine": &dcl.Property{
													Type:        "string",
													GoName:      "VirtualMachine",
													Description: "The URI of the virtual machine resource",
													Immutable:   true,
													ResourceReferences: []*dcl.PropertyResourceReference{
														&dcl.PropertyResourceReference{
															Resource: "Compute/Instance",
															Field:    "selfLink",
														},
													},
												},
											},
										},
									},
									"siteToSiteDataTransfer": &dcl.Property{
										Type:        "boolean",
										GoName:      "SiteToSiteDataTransfer",
										Description: "A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.",
										Immutable:   true,
									},
								},
							},
							"linkedVpnTunnels": &dcl.Property{
								Type:        "object",
								GoName:      "LinkedVpnTunnels",
								GoType:      "SpokeLinkedVpnTunnels",
								Description: "The URIs of linked VPN tunnel resources",
								Immutable:   true,
								Conflicts: []string{
									"linkedInterconnectAttachments",
									"linkedRouterApplianceInstances",
								},
								Required: []string{
									"uris",
									"siteToSiteDataTransfer",
								},
								Properties: map[string]*dcl.Property{
									"siteToSiteDataTransfer": &dcl.Property{
										Type:        "boolean",
										GoName:      "SiteToSiteDataTransfer",
										Description: "A value that controls whether site-to-site data transfer is enabled for these resources. Note that data transfer is available only in supported locations.",
										Immutable:   true,
									},
									"uris": &dcl.Property{
										Type:        "array",
										GoName:      "Uris",
										Description: "The URIs of linked VPN tunnel resources.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/VpnTunnel",
													Field:    "selfLink",
												},
											},
										},
									},
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Immutable. The name of the spoke. Spoke names must be unique.",
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
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "SpokeStateEnum",
								ReadOnly:    true,
								Description: "Output only. The current lifecycle state of this spoke. Possible values: STATE_UNSPECIFIED, CREATING, ACTIVE, DELETING",
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
								Description: "Output only. The Google-generated UUID for the spoke. This value is unique across all spoke resources. If a spoke is deleted and another with the same name is created, the new spoke is assigned a different unique_id.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time the spoke was last updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
