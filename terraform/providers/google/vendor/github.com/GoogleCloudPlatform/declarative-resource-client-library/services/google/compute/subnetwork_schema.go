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

func DCLSubnetworkSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/Subnetwork",
			Description: "The Compute Subnetwork resource",
			StructName:  "Subnetwork",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Subnetwork",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "subnetwork",
						Required:    true,
						Description: "A full instance of a Subnetwork",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Subnetwork",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "subnetwork",
						Required:    true,
						Description: "A full instance of a Subnetwork",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Subnetwork",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "subnetwork",
						Required:    true,
						Description: "A full instance of a Subnetwork",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Subnetwork",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "region",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Subnetwork",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "region",
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
				"Subnetwork": &dcl.Component{
					Title: "Subnetwork",
					ID:    "projects/{{project}}/regions/{{region}}/subnetworks/{{name}}",
					Locations: []string{
						"region",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"ipCidrRange",
							"name",
							"network",
							"region",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"creationTimestamp": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreationTimestamp",
								ReadOnly:    true,
								Description: "Creation timestamp in RFC3339 text format.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource. Provide this property when you create the resource. This field can be set only at resource creation time. ",
								Immutable:   true,
							},
							"enableFlowLogs": &dcl.Property{
								Type:          "boolean",
								GoName:        "EnableFlowLogs",
								Description:   "Whether to enable flow logging for this subnetwork. If this field is not explicitly set, it will not appear in `get` listings. If not set the default behavior is to disable flow logging. This field isn't supported with the `purpose` field set to `INTERNAL_HTTPS_LOAD_BALANCER`.",
								Immutable:     true,
								ServerDefault: true,
							},
							"fingerprint": &dcl.Property{
								Type:        "string",
								GoName:      "Fingerprint",
								ReadOnly:    true,
								Description: "Fingerprint of this resource. This field is used internally during updates of this resource. ",
								Immutable:   true,
							},
							"gatewayAddress": &dcl.Property{
								Type:        "string",
								GoName:      "GatewayAddress",
								ReadOnly:    true,
								Description: "The gateway address for default routes to reach destination addresses outside this subnetwork. ",
								Immutable:   true,
							},
							"ipCidrRange": &dcl.Property{
								Type:        "string",
								GoName:      "IPCidrRange",
								Description: "The range of internal addresses that are owned by this subnetwork. Provide this property when you create the subnetwork. For example, 10.0.0.0/8 or 192.168.0.0/16. Ranges must be unique and non-overlapping within a network. Only IPv4 is supported. ",
							},
							"logConfig": &dcl.Property{
								Type:        "object",
								GoName:      "LogConfig",
								GoType:      "SubnetworkLogConfig",
								Description: "Denotes the logging options for the subnetwork flow logs. If logging is enabled logs will be exported to Cloud Logging. ",
								Properties: map[string]*dcl.Property{
									"aggregationInterval": &dcl.Property{
										Type:        "string",
										GoName:      "AggregationInterval",
										GoType:      "SubnetworkLogConfigAggregationIntervalEnum",
										Description: "Can only be specified if VPC flow logging for this subnetwork is enabled. Toggles the aggregation interval for collecting flow logs. Increasing the interval time will reduce the amount of generated flow logs for long lasting connections. Default is an interval of 5 seconds per connection. Possible values are INTERVAL_5_SEC, INTERVAL_30_SEC, INTERVAL_1_MIN, INTERVAL_5_MIN, INTERVAL_10_MIN, INTERVAL_15_MIN ",
										Default:     "INTERVAL_5_SEC",
										Enum: []string{
											"INTERVAL_5_SEC",
											"INTERVAL_30_SEC",
											"INTERVAL_1_MIN",
											"INTERVAL_5_MIN",
											"INTERVAL_10_MIN",
											"INTERVAL_15_MIN",
										},
									},
									"flowSampling": &dcl.Property{
										Type:        "number",
										Format:      "double",
										GoName:      "FlowSampling",
										Description: "Can only be specified if VPC flow logging for this subnetwork is enabled. The value of the field must be in [0, 1]. Set the sampling rate of VPC flow logs within the subnetwork where 1.0 means all collected logs are reported and 0.0 means no logs are reported. Default is 0.5 which means half of all collected logs are reported. ",
										Default:     0.5,
									},
									"metadata": &dcl.Property{
										Type:        "string",
										GoName:      "Metadata",
										GoType:      "SubnetworkLogConfigMetadataEnum",
										Description: "Can only be specified if VPC flow logging for this subnetwork is enabled. Configures whether metadata fields should be added to the reported VPC flow logs. Default is `INCLUDE_ALL_METADATA`.  Possible values: EXCLUDE_ALL_METADATA, INCLUDE_ALL_METADATA",
										Default:     "INCLUDE_ALL_METADATA",
										Enum: []string{
											"EXCLUDE_ALL_METADATA",
											"INCLUDE_ALL_METADATA",
										},
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of the resource, provided by the client when initially creating the resource. The name must be 1-63 characters long, and comply with RFC1035. Specifically, the name must be 1-63 characters long and match the regular expression `[a-z]([-a-z0-9]*[a-z0-9])?` which means the first character must be a lowercase letter, and all following characters must be a dash, lowercase letter, or digit, except the last character, which cannot be a dash. ",
								Immutable:   true,
							},
							"network": &dcl.Property{
								Type:        "string",
								GoName:      "Network",
								Description: "The network this subnet belongs to. Only networks that are in the distributed mode can have subnetworks. ",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/Network",
										Field:    "selfLink",
									},
								},
							},
							"privateIPGoogleAccess": &dcl.Property{
								Type:        "boolean",
								GoName:      "PrivateIPGoogleAccess",
								Description: "When enabled, VMs in this subnetwork without external IP addresses can access Google APIs and services by using Private Google Access. ",
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
								Parameter: true,
							},
							"purpose": &dcl.Property{
								Type:          "string",
								GoName:        "Purpose",
								GoType:        "SubnetworkPurposeEnum",
								Description:   "The purpose of the resource. This field can be either PRIVATE or INTERNAL_HTTPS_LOAD_BALANCER. A subnetwork with purpose set to INTERNAL_HTTPS_LOAD_BALANCER is a user-created subnetwork that is reserved for Internal HTTP(S) Load Balancing. If unspecified, the purpose defaults to PRIVATE.  If set to INTERNAL_HTTPS_LOAD_BALANCER you must also set the role. ",
								Immutable:     true,
								ServerDefault: true,
								Enum: []string{
									"INTERNAL_HTTPS_LOAD_BALANCER",
									"PRIVATE",
									"AGGREGATE",
									"PRIVATE_SERVICE_CONNECT",
									"CLOUD_EXTENSION",
									"PRIVATE_NAT",
								},
							},
							"region": &dcl.Property{
								Type:        "string",
								GoName:      "Region",
								Description: "The GCP region for this subnetwork. ",
								Immutable:   true,
								Parameter:   true,
							},
							"role": &dcl.Property{
								Type:        "string",
								GoName:      "Role",
								GoType:      "SubnetworkRoleEnum",
								Description: "The role of subnetwork. Currenly, this field is only used when purpose = INTERNAL_HTTPS_LOAD_BALANCER. The value can be set to ACTIVE or BACKUP. An ACTIVE subnetwork is one that is currently being used for Internal HTTP(S) Load Balancing. A BACKUP subnetwork is one that is ready to be promoted to ACTIVE or is currently draining. ",
								Enum: []string{
									"ACTIVE",
									"BACKUP",
								},
							},
							"secondaryIPRanges": &dcl.Property{
								Type:          "array",
								GoName:        "SecondaryIPRanges",
								Description:   "An array of configurations for secondary IP ranges for VM instances contained in this subnetwork. The primary IP of such VM must belong to the primary ipCidrRange of the subnetwork. The alias IPs may belong to either primary or secondary ranges. This field uses attr-as-block mode to avoid breaking users during the 0.12 upgrade. See [the Attr-as-Block page](https://www.terraform.io/docs/configuration/attr-as-blocks.html) for more details. ",
								ServerDefault: true,
								SendEmpty:     true,
								ListType:      "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "SubnetworkSecondaryIPRanges",
									Required: []string{
										"rangeName",
										"ipCidrRange",
									},
									Properties: map[string]*dcl.Property{
										"ipCidrRange": &dcl.Property{
											Type:        "string",
											GoName:      "IPCidrRange",
											Description: "The range of IP addresses belonging to this subnetwork secondary range. Provide this property when you create the subnetwork. Ranges must be unique and non-overlapping with all primary and secondary IP ranges within a network. Only IPv4 is supported. ",
										},
										"rangeName": &dcl.Property{
											Type:        "string",
											GoName:      "RangeName",
											Description: "The name associated with this subnetwork secondary range, used when adding an alias IP range to a VM instance. The name must be 1-63 characters long, and comply with RFC1035. The name must be unique within the subnetwork. ",
										},
									},
								},
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "[Output Only] Server-defined URL for the resource.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
