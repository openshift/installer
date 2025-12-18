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

func DCLRouteSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/Route",
			Description: "The Compute Route resource",
			StructName:  "Route",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Route",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "route",
						Required:    true,
						Description: "A full instance of a Route",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Route",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "route",
						Required:    true,
						Description: "A full instance of a Route",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Route",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "route",
						Required:    true,
						Description: "A full instance of a Route",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Route",
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
				Description: "The function used to list information about many Route",
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
				"Route": &dcl.Component{
					Title: "Route",
					ID:    "projects/{{project}}/global/routes/{{name}}",
					Locations: []string{
						"global",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"network",
							"destRange",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. Provide this field when you\ncreate the resource.",
								Immutable:   true,
							},
							"destRange": &dcl.Property{
								Type:        "string",
								GoName:      "DestRange",
								Description: "The destination range of the route.",
								Immutable:   true,
							},
							"id": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "[Output Only] The unique identifier for the resource. This identifier is\ndefined by the server.",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the resource. Provided by the client when the resource is created.\nThe name must be 1-63 characters long, and comply with\n<a href=\"https://www.ietf.org/rfc/rfc1035.txt\" target=\"_blank\">RFC1035</a>.\nSpecifically, the name must be 1-63 characters long and match the regular\nexpression `[a-z]([-a-z0-9]*[a-z0-9])?`. The first character must be a\nlowercase letter, and all following characters (except for the last\ncharacter) must be a dash, lowercase letter, or digit. The last character\nmust be a lowercase letter or digit.",
								Immutable:   true,
							},
							"network": &dcl.Property{
								Type:        "string",
								GoName:      "Network",
								Description: "Fully-qualified URL of the network that this route applies to.",
								Immutable:   true,
							},
							"nextHopGateway": &dcl.Property{
								Type:        "string",
								GoName:      "NextHopGateway",
								Description: "The URL to a gateway that should handle matching packets.\nYou can only specify the internet gateway using a full or\npartial valid URL: </br>\n<code>projects/<var\nclass=\"apiparam\">project</var>/global/gateways/default-internet-gateway</code>",
								Immutable:   true,
								Conflicts: []string{
									"nextHopVpnTunnel",
									"nextHopIP",
									"nextHopInstance",
									"nextHopIlb",
								},
							},
							"nextHopIP": &dcl.Property{
								Type:        "string",
								GoName:      "NextHopIP",
								Description: "The network IP address of an instance that should handle matching packets.\nOnly IPv4 is supported.",
								Immutable:   true,
								Conflicts: []string{
									"nextHopVpnTunnel",
									"nextHopInstance",
									"nextHopGateway",
									"nextHopIlb",
								},
								ServerDefault: true,
							},
							"nextHopIlb": &dcl.Property{
								Type:        "string",
								GoName:      "NextHopIlb",
								Description: "The URL to a forwarding rule of type\n<code>loadBalancingScheme=INTERNAL</code> that should handle matching\npackets. You can only specify the forwarding rule as a partial or full\nURL. For example, the following are all valid URLs:\n<ul>\n   <li><code>https://www.googleapis.com/compute/v1/projects/<var\n   class=\"apiparam\">project</var>/regions/<var\n   class=\"apiparam\">region</var>/forwardingRules/<var\n   class=\"apiparam\">forwardingRule</var></code></li> <li><code>regions/<var\n   class=\"apiparam\">region</var>/forwardingRules/<var\n   class=\"apiparam\">forwardingRule</var></code></li>\n</ul>",
								Immutable:   true,
								Conflicts: []string{
									"nextHopVpnTunnel",
									"nextHopIP",
									"nextHopInstance",
									"nextHopGateway",
								},
							},
							"nextHopInstance": &dcl.Property{
								Type:        "string",
								GoName:      "NextHopInstance",
								Description: "The URL to an instance that should handle matching packets. You can specify\nthis as a full or partial URL.\nFor example: <br />\n<code>https://www.googleapis.com/compute/v1/projects/<var\nclass=\"apiparam\">project</var>/zones/<var\nclass=\"apiparam\">zone</var>/instances/<instance-name></code>",
								Immutable:   true,
								Conflicts: []string{
									"nextHopVpnTunnel",
									"nextHopIP",
									"nextHopGateway",
									"nextHopIlb",
								},
							},
							"nextHopNetwork": &dcl.Property{
								Type:        "string",
								GoName:      "NextHopNetwork",
								ReadOnly:    true,
								Description: "The URL of the local network if it should handle matching packets.",
								Immutable:   true,
							},
							"nextHopPeering": &dcl.Property{
								Type:        "string",
								GoName:      "NextHopPeering",
								ReadOnly:    true,
								Description: "[Output Only] The network peering name that should handle matching packets,\nwhich should conform to RFC1035.",
								Immutable:   true,
							},
							"nextHopVpnTunnel": &dcl.Property{
								Type:        "string",
								GoName:      "NextHopVpnTunnel",
								Description: "The URL to a VpnTunnel that should handle matching packets.",
								Immutable:   true,
								Conflicts: []string{
									"nextHopIP",
									"nextHopInstance",
									"nextHopGateway",
									"nextHopIlb",
								},
							},
							"priority": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Priority",
								Description: "The priority of the peering route.",
								Immutable:   true,
								Default:     1000,
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
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "[Output Only] Server-defined fully-qualified URL for this resource.",
								Immutable:   true,
							},
							"tag": &dcl.Property{
								Type:        "array",
								GoName:      "Tag",
								Description: "A list of instance tags to which this route applies.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"warning": &dcl.Property{
								Type:        "array",
								GoName:      "Warning",
								ReadOnly:    true,
								Description: "[Output Only] If potential misconfigurations are detected for this\nroute, this field will be populated with warning messages.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "RouteWarning",
									Properties: map[string]*dcl.Property{
										"code": &dcl.Property{
											Type:        "string",
											GoName:      "Code",
											GoType:      "RouteWarningCodeEnum",
											ReadOnly:    true,
											Description: "[Output Only] A warning code, if applicable. For example, Compute\nEngine returns <code>NO_RESULTS_ON_PAGE</code> if there\nare no results in the response. Possible values: BAD_REQUEST, FORBIDDEN, NOT_FOUND, CONFLICT, GONE, PRECONDITION_FAILED, INTERNAL_ERROR, SERVICE_UNAVAILABLE",
											Immutable:   true,
											Enum: []string{
												"BAD_REQUEST",
												"FORBIDDEN",
												"NOT_FOUND",
												"CONFLICT",
												"GONE",
												"PRECONDITION_FAILED",
												"INTERNAL_ERROR",
												"SERVICE_UNAVAILABLE",
											},
										},
										"data": &dcl.Property{
											Type: "object",
											AdditionalProperties: &dcl.Property{
												Type: "string",
											},
											GoName:      "Data",
											ReadOnly:    true,
											Description: "[Output Only] Metadata about this warning in <code class=\"lang-html\">key:\nvalue</code> format. For example:\n<pre class=\"lang-html prettyprint\">\"data\": [\n : {\n   \"key\": \"scope\",\n   \"value\": \"zones/us-east1-d\"\n  }</pre>",
											Immutable:   true,
										},
										"message": &dcl.Property{
											Type:        "string",
											GoName:      "Message",
											ReadOnly:    true,
											Description: "[Output Only] A human-readable description of the warning code.",
											Immutable:   true,
										},
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
