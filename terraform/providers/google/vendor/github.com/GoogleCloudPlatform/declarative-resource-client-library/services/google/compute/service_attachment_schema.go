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

func DCLServiceAttachmentSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/ServiceAttachment",
			Description: "Represents a ServiceAttachment resource.",
			StructName:  "ServiceAttachment",
			Reference: &dcl.Link{
				Text: "API documentation",
				URL:  "https://cloud.google.com/compute/docs/reference/rest/beta/serviceAttachments",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Configuring Private Service Connect to access services",
					URL:  "https://cloud.google.com/vpc/docs/configure-private-service-connect-services",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a ServiceAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "serviceAttachment",
						Required:    true,
						Description: "A full instance of a ServiceAttachment",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a ServiceAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "serviceAttachment",
						Required:    true,
						Description: "A full instance of a ServiceAttachment",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a ServiceAttachment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "serviceAttachment",
						Required:    true,
						Description: "A full instance of a ServiceAttachment",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all ServiceAttachment",
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
				Description: "The function used to list information about many ServiceAttachment",
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
				"ServiceAttachment": &dcl.Component{
					Title: "ServiceAttachment",
					ID:    "projects/{{project}}/regions/{{location}}/serviceAttachments/{{name}}",
					Locations: []string{
						"region",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"targetService",
							"connectionPreference",
							"natSubnets",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"connectedEndpoints": &dcl.Property{
								Type:        "array",
								GoName:      "ConnectedEndpoints",
								ReadOnly:    true,
								Description: "An array of connections for all the consumers connected to this service attachment.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "ServiceAttachmentConnectedEndpoints",
									Properties: map[string]*dcl.Property{
										"endpoint": &dcl.Property{
											Type:        "string",
											GoName:      "Endpoint",
											Description: "The url of a connected endpoint.",
										},
										"pscConnectionId": &dcl.Property{
											Type:        "integer",
											Format:      "int64",
											GoName:      "PscConnectionId",
											Description: "The PSC connection id of the connected endpoint.",
										},
										"status": &dcl.Property{
											Type:        "string",
											GoName:      "Status",
											GoType:      "ServiceAttachmentConnectedEndpointsStatusEnum",
											Description: "The status of a connected endpoint to this service attachment. Possible values: PENDING, RUNNING, DONE",
											Enum: []string{
												"PENDING",
												"RUNNING",
												"DONE",
											},
										},
									},
								},
							},
							"connectionPreference": &dcl.Property{
								Type:        "string",
								GoName:      "ConnectionPreference",
								GoType:      "ServiceAttachmentConnectionPreferenceEnum",
								Description: "The connection preference of service attachment. The value can be set to `ACCEPT_AUTOMATIC`. An `ACCEPT_AUTOMATIC` service attachment is one that always accepts the connection from consumer forwarding rules. Possible values: CONNECTION_PREFERENCE_UNSPECIFIED, ACCEPT_AUTOMATIC, ACCEPT_MANUAL",
								Enum: []string{
									"CONNECTION_PREFERENCE_UNSPECIFIED",
									"ACCEPT_AUTOMATIC",
									"ACCEPT_MANUAL",
								},
							},
							"consumerAcceptLists": &dcl.Property{
								Type:        "array",
								GoName:      "ConsumerAcceptLists",
								Description: "Projects that are allowed to connect to this service attachment.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "ServiceAttachmentConsumerAcceptLists",
									Required: []string{
										"projectIdOrNum",
									},
									Properties: map[string]*dcl.Property{
										"connectionLimit": &dcl.Property{
											Type:        "integer",
											Format:      "int64",
											GoName:      "ConnectionLimit",
											Description: "The value of the limit to set.",
										},
										"projectIdOrNum": &dcl.Property{
											Type:        "string",
											GoName:      "ProjectIdOrNum",
											Description: "The project id or number for the project to set the limit for.",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Cloudresourcemanager/Project",
													Field:    "name",
												},
											},
										},
									},
								},
							},
							"consumerRejectLists": &dcl.Property{
								Type:        "array",
								GoName:      "ConsumerRejectLists",
								Description: "Projects that are not allowed to connect to this service attachment. The project can be specified using its id or number.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
									ResourceReferences: []*dcl.PropertyResourceReference{
										&dcl.PropertyResourceReference{
											Resource: "Cloudresourcemanager/Project",
											Field:    "name",
										},
									},
								},
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. Provide this property when you create the resource.",
							},
							"enableProxyProtocol": &dcl.Property{
								Type:        "boolean",
								GoName:      "EnableProxyProtocol",
								Description: "If true, enable the proxy protocol which is for supplying client TCP/IP address data in TCP connections that traverse proxies on their way to destination servers.",
								Immutable:   true,
							},
							"fingerprint": &dcl.Property{
								Type:        "string",
								GoName:      "Fingerprint",
								ReadOnly:    true,
								Description: "Fingerprint of this resource. This field is used internally during updates of this resource.",
								Immutable:   true,
							},
							"id": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "The unique identifier for the resource type. The server generates this identifier.",
								Immutable:   true,
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
								Description: "Name of the resource. Provided by the client when the resource is created. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt). Specifically, the name must be 1-63 characters long and match the regular expression `)?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash.",
								Immutable:   true,
							},
							"natSubnets": &dcl.Property{
								Type:        "array",
								GoName:      "NatSubnets",
								Description: "An array of URLs where each entry is the URL of a subnet provided by the service producer to use for NAT in this service attachment.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
									ResourceReferences: []*dcl.PropertyResourceReference{
										&dcl.PropertyResourceReference{
											Resource: "Compute/Subnetwork",
											Field:    "selfLink",
										},
									},
								},
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
							"pscServiceAttachmentId": &dcl.Property{
								Type:        "object",
								GoName:      "PscServiceAttachmentId",
								GoType:      "ServiceAttachmentPscServiceAttachmentId",
								ReadOnly:    true,
								Description: "An 128-bit global unique ID of the PSC service attachment.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"high": &dcl.Property{
										Type:      "integer",
										Format:    "int64",
										GoName:    "High",
										ReadOnly:  true,
										Immutable: true,
									},
									"low": &dcl.Property{
										Type:      "integer",
										Format:    "int64",
										GoName:    "Low",
										ReadOnly:  true,
										Immutable: true,
									},
								},
							},
							"region": &dcl.Property{
								Type:        "string",
								GoName:      "Region",
								ReadOnly:    true,
								Description: "URL of the region where the service attachment resides. This field applies only to the region resource. You must specify this field as part of the HTTP request URL. It is not settable as a field in the request body.",
								Immutable:   true,
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "Server-defined URL for the resource.",
								Immutable:   true,
							},
							"targetService": &dcl.Property{
								Type:        "string",
								GoName:      "TargetService",
								Description: "The URL of a service serving the endpoint identified by this service attachment.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/ForwardingRule",
										Field:    "selfLink",
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
