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

func DCLPacketMirroringSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/PacketMirroring",
			Description: "Packet Mirroring mirrors traffic to and from particular VM instances. You can use the collected traffic to help you detect security threats and monitor application performance.",
			StructName:  "PacketMirroring",
			Reference: &dcl.Link{
				Text: "API documentation",
				URL:  "https://cloud.google.com/compute/docs/reference/rest/beta/packetMirrorings",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Using Packet Mirroring",
					URL:  "https://cloud.google.com/vpc/docs/using-packet-mirroring",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a PacketMirroring",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "packetMirroring",
						Required:    true,
						Description: "A full instance of a PacketMirroring",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a PacketMirroring",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "packetMirroring",
						Required:    true,
						Description: "A full instance of a PacketMirroring",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a PacketMirroring",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "packetMirroring",
						Required:    true,
						Description: "A full instance of a PacketMirroring",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all PacketMirroring",
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
				Description: "The function used to list information about many PacketMirroring",
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
				"PacketMirroring": &dcl.Component{
					Title: "PacketMirroring",
					ID:    "projects/{{project}}/regions/{{location}}/packetMirrorings/{{name}}",
					Locations: []string{
						"region",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"network",
							"collectorIlb",
							"mirroredResources",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"collectorIlb": &dcl.Property{
								Type:        "object",
								GoName:      "CollectorIlb",
								GoType:      "PacketMirroringCollectorIlb",
								Description: "The Forwarding Rule resource of type `loadBalancingScheme=INTERNAL` that will be used as collector for mirrored traffic. The specified forwarding rule must have `isMirroringCollector` set to true.",
								Required: []string{
									"url",
								},
								Properties: map[string]*dcl.Property{
									"canonicalUrl": &dcl.Property{
										Type:        "string",
										GoName:      "CanonicalUrl",
										ReadOnly:    true,
										Description: "Output only. Unique identifier for the forwarding rule; defined by the server.",
									},
									"url": &dcl.Property{
										Type:        "string",
										GoName:      "Url",
										Description: "Resource URL to the forwarding rule representing the ILB configured as destination of the mirrored traffic.",
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Compute/ForwardingRule",
												Field:    "selfLink",
											},
										},
									},
								},
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. Provide this property when you create the resource.",
							},
							"enable": &dcl.Property{
								Type:          "string",
								GoName:        "Enable",
								GoType:        "PacketMirroringEnableEnum",
								Description:   "Indicates whether or not this packet mirroring takes effect. If set to FALSE, this packet mirroring policy will not be enforced on the network. The default is TRUE.",
								ServerDefault: true,
								Enum: []string{
									"TRUE",
									"FALSE",
								},
							},
							"filter": &dcl.Property{
								Type:          "object",
								GoName:        "Filter",
								GoType:        "PacketMirroringFilter",
								Description:   "Filter for mirrored traffic. If unspecified, all traffic is mirrored.",
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"cidrRanges": &dcl.Property{
										Type:        "array",
										GoName:      "CidrRanges",
										Description: "IP CIDR ranges that apply as filter on the source (ingress) or destination (egress) IP in the IP header. Only IPv4 is supported. If no ranges are specified, all traffic that matches the specified IPProtocols is mirrored. If neither cidrRanges nor IPProtocols is specified, all traffic is mirrored.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"direction": &dcl.Property{
										Type:        "string",
										GoName:      "Direction",
										GoType:      "PacketMirroringFilterDirectionEnum",
										Description: "Direction of traffic to mirror, either INGRESS, EGRESS, or BOTH. The default is BOTH.",
										Enum: []string{
											"INGRESS",
											"EGRESS",
										},
									},
									"ipProtocols": &dcl.Property{
										Type:        "array",
										GoName:      "IPProtocols",
										Description: "Protocols that apply as filter on mirrored traffic. If no protocols are specified, all traffic that matches the specified CIDR ranges is mirrored. If neither cidrRanges nor IPProtocols is specified, all traffic is mirrored.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
								},
							},
							"id": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "Output only. The unique identifier for the resource. This identifier is defined by the server.",
								Immutable:   true,
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"mirroredResources": &dcl.Property{
								Type:        "object",
								GoName:      "MirroredResources",
								GoType:      "PacketMirroringMirroredResources",
								Description: "PacketMirroring mirroredResourceInfos. MirroredResourceInfo specifies a set of mirrored VM instances, subnetworks and/or tags for which traffic from/to all VM instances will be mirrored.",
								Properties: map[string]*dcl.Property{
									"instances": &dcl.Property{
										Type:        "array",
										GoName:      "Instances",
										Description: "A set of virtual machine instances that are being mirrored. They must live in zones contained in the same region as this packetMirroring. Note that this config will apply only to those network interfaces of the Instances that belong to the network specified in this packetMirroring. You may specify a maximum of 50 Instances.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "PacketMirroringMirroredResourcesInstances",
											Properties: map[string]*dcl.Property{
												"canonicalUrl": &dcl.Property{
													Type:        "string",
													GoName:      "CanonicalUrl",
													ReadOnly:    true,
													Description: "Output only. Unique identifier for the instance; defined by the server.",
													Immutable:   true,
												},
												"url": &dcl.Property{
													Type:        "string",
													GoName:      "Url",
													Description: "Resource URL to the virtual machine instance which is being mirrored.",
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
									"subnetworks": &dcl.Property{
										Type:        "array",
										GoName:      "Subnetworks",
										Description: "A set of subnetworks for which traffic from/to all VM instances will be mirrored. They must live in the same region as this packetMirroring. You may specify a maximum of 5 subnetworks.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "PacketMirroringMirroredResourcesSubnetworks",
											Properties: map[string]*dcl.Property{
												"canonicalUrl": &dcl.Property{
													Type:        "string",
													GoName:      "CanonicalUrl",
													ReadOnly:    true,
													Description: "Output only. Unique identifier for the subnetwork; defined by the server.",
													Immutable:   true,
												},
												"url": &dcl.Property{
													Type:        "string",
													GoName:      "Url",
													Description: "Resource URL to the subnetwork for which traffic from/to all VM instances will be mirrored.",
													Immutable:   true,
													ResourceReferences: []*dcl.PropertyResourceReference{
														&dcl.PropertyResourceReference{
															Resource: "Compute/Subnetwork",
															Field:    "selfLink",
														},
													},
												},
											},
										},
									},
									"tags": &dcl.Property{
										Type:        "array",
										GoName:      "Tags",
										Description: "A set of mirrored tags. Traffic from/to all VM instances that have one or more of these tags will be mirrored.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the resource; provided by the client when the resource is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the name must be 1-63 characters long and match the regular expression `)?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
								Immutable:   true,
							},
							"network": &dcl.Property{
								Type:        "object",
								GoName:      "Network",
								GoType:      "PacketMirroringNetwork",
								Description: "Specifies the mirrored VPC network. Only packets in this network will be mirrored. All mirrored VMs should have a NIC in the given network. All mirrored subnetworks should belong to the given network.",
								Immutable:   true,
								Required: []string{
									"url",
								},
								Properties: map[string]*dcl.Property{
									"canonicalUrl": &dcl.Property{
										Type:        "string",
										GoName:      "CanonicalUrl",
										ReadOnly:    true,
										Description: "Output only. Unique identifier for the network; defined by the server.",
										Immutable:   true,
									},
									"url": &dcl.Property{
										Type:        "string",
										GoName:      "Url",
										Description: "URL of the network resource.",
										Immutable:   true,
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Compute/Network",
												Field:    "selfLink",
											},
										},
									},
								},
							},
							"priority": &dcl.Property{
								Type:          "integer",
								Format:        "int64",
								GoName:        "Priority",
								Description:   "The priority of applying this configuration. Priority is used to break ties in cases where there is more than one matching rule. In the case of two rules that apply for a given Instance, the one with the lowest-numbered priority value wins. Default value is 1000. Valid range is 0 through 65535.",
								ServerDefault: true,
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
								Description: "URI of the region where the packetMirroring resides.",
								Immutable:   true,
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "Server-defined URL for the resource.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
