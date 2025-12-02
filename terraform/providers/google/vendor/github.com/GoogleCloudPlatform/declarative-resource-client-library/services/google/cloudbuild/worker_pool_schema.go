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
package cloudbuild

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLWorkerPoolSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "CloudBuild/WorkerPool",
			Description: "The CloudBuild WorkerPool resource",
			StructName:  "WorkerPool",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a WorkerPool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workerPool",
						Required:    true,
						Description: "A full instance of a WorkerPool",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a WorkerPool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workerPool",
						Required:    true,
						Description: "A full instance of a WorkerPool",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a WorkerPool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workerPool",
						Required:    true,
						Description: "A full instance of a WorkerPool",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all WorkerPool",
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
				Description: "The function used to list information about many WorkerPool",
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
				"WorkerPool": &dcl.Component{
					Title:           "WorkerPool",
					ID:              "projects/{{project}}/locations/{{location}}/workerPools/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"annotations": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Annotations",
								Description: "User specified annotations. See https://google.aip.dev/128#annotations for more details such as format and size limitations.",
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. Time at which the request to create the `WorkerPool` was received.",
								Immutable:   true,
							},
							"deleteTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "DeleteTime",
								ReadOnly:    true,
								Description: "Output only. Time at which the request to delete the `WorkerPool` was received.",
								Immutable:   true,
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "A user-specified, human-readable name for the `WorkerPool`. If provided, this value must be 1-63 characters.",
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Output only. Checksum computed by the server. May be sent on update and delete requests to ensure that the client has an up-to-date value before proceeding.",
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
								Description: "User-defined name of the `WorkerPool`.",
								Immutable:   true,
							},
							"networkConfig": &dcl.Property{
								Type:        "object",
								GoName:      "NetworkConfig",
								GoType:      "WorkerPoolNetworkConfig",
								Description: "Network configuration for the `WorkerPool`.",
								Immutable:   true,
								Required: []string{
									"peeredNetwork",
								},
								Properties: map[string]*dcl.Property{
									"peeredNetwork": &dcl.Property{
										Type:        "string",
										GoName:      "PeeredNetwork",
										Description: "Required. Immutable. The network definition that the workers are peered to. If this section is left empty, the workers will be peered to `WorkerPool.project_id` on the service producer network. Must be in the format `projects/{project}/global/networks/{network}`, where `{project}` is a project number, such as `12345`, and `{network}` is the name of a VPC network in the project. See [Understanding network configuration options](https://cloud.google.com/cloud-build/docs/custom-workers/set-up-custom-worker-pool-environment#understanding_the_network_configuration_options)",
										Immutable:   true,
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Compute/Network",
												Field:    "selfLink",
											},
										},
									},
									"peeredNetworkIPRange": &dcl.Property{
										Type:        "string",
										GoName:      "PeeredNetworkIPRange",
										Description: "Optional. Immutable. Subnet IP range within the peered network. This is specified in CIDR notation with a slash and the subnet prefix size. You can optionally specify an IP address before the subnet prefix value. e.g. `192.168.0.0/29` would specify an IP range starting at 192.168.0.0 with a prefix size of 29 bits. `/16` would specify a prefix size of 16 bits, with an automatically determined IP within the peered VPC. If unspecified, a value of `/24` will be used.",
										Immutable:   true,
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
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "WorkerPoolStateEnum",
								ReadOnly:    true,
								Description: "Output only. `WorkerPool` state. Possible values: STATE_UNSPECIFIED, PENDING, APPROVED, REJECTED, CANCELLED",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"PENDING",
									"APPROVED",
									"REJECTED",
									"CANCELLED",
								},
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. A unique identifier for the `WorkerPool`.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Time at which the request to update the `WorkerPool` was received.",
								Immutable:   true,
							},
							"workerConfig": &dcl.Property{
								Type:          "object",
								GoName:        "WorkerConfig",
								GoType:        "WorkerPoolWorkerConfig",
								Description:   "Configuration to be used for a creating workers in the `WorkerPool`.",
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"diskSizeGb": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "DiskSizeGb",
										Description: "Size of the disk attached to the worker, in GB. See [Worker pool config file](https://cloud.google.com/cloud-build/docs/custom-workers/worker-pool-config-file). Specify a value of up to 1000. If `0` is specified, Cloud Build will use a standard disk size.",
									},
									"machineType": &dcl.Property{
										Type:        "string",
										GoName:      "MachineType",
										Description: "Machine type of a worker, such as `n1-standard-1`. See [Worker pool config file](https://cloud.google.com/cloud-build/docs/custom-workers/worker-pool-config-file). If left blank, Cloud Build will use `n1-standard-1`.",
									},
									"noExternalIP": &dcl.Property{
										Type:          "boolean",
										GoName:        "NoExternalIP",
										Description:   "If true, workers are created without any public address, which prevents network egress to public IPs.",
										ServerDefault: true,
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
