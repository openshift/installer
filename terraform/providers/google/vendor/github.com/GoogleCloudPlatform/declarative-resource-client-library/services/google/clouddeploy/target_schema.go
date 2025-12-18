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
package clouddeploy

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLTargetSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Clouddeploy/Target",
			Description: "The Cloud Deploy `Target` resource",
			StructName:  "Target",
			Reference: &dcl.Link{
				Text: "REST API",
				URL:  "https://cloud.google.com/deploy/docs/api/reference/rest/v1/projects.locations.targets",
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Target",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "target",
						Required:    true,
						Description: "A full instance of a Target",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Target",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "target",
						Required:    true,
						Description: "A full instance of a Target",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Target",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "target",
						Required:    true,
						Description: "A full instance of a Target",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Target",
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
				Description: "The function used to list information about many Target",
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
				"Target": &dcl.Component{
					Title:           "Target",
					ID:              "projects/{{project}}/locations/{{location}}/targets/{{name}}",
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
								Description: "Optional. User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.",
							},
							"anthosCluster": &dcl.Property{
								Type:        "object",
								GoName:      "AnthosCluster",
								GoType:      "TargetAnthosCluster",
								Description: "Information specifying an Anthos Cluster.",
								Conflicts: []string{
									"gke",
								},
								Properties: map[string]*dcl.Property{
									"membership": &dcl.Property{
										Type:        "string",
										GoName:      "Membership",
										Description: "Membership of the GKE Hub-registered cluster to which to apply the Skaffold configuration. Format is `projects/{project}/locations/{location}/memberships/{membership_name}`.",
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Gkehub/Membership",
												Field:    "selfLink",
											},
										},
									},
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. Time at which the `Target` was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. Description of the `Target`. Max length is 255 characters.",
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Optional. This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"executionConfigs": &dcl.Property{
								Type:          "array",
								GoName:        "ExecutionConfigs",
								Description:   "Configurations for all execution that relates to this `Target`. Each `ExecutionEnvironmentUsage` value may only be used in a single configuration; using the same value multiple times is an error. When one or more configurations are specified, they must include the `RENDER` and `DEPLOY` `ExecutionEnvironmentUsage` values. When no configurations are specified, execution will use the default specified in `DefaultPool`.",
								ServerDefault: true,
								SendEmpty:     true,
								ListType:      "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "TargetExecutionConfigs",
									Required: []string{
										"usages",
									},
									Properties: map[string]*dcl.Property{
										"artifactStorage": &dcl.Property{
											Type:          "string",
											GoName:        "ArtifactStorage",
											Description:   "Optional. Cloud Storage location in which to store execution outputs. This can either be a bucket (\"gs://my-bucket\") or a path within a bucket (\"gs://my-bucket/my-dir\"). If unspecified, a default bucket located in the same region will be used.",
											ServerDefault: true,
										},
										"executionTimeout": &dcl.Property{
											Type:          "string",
											GoName:        "ExecutionTimeout",
											Description:   "Optional. Execution timeout for a Cloud Build Execution. This must be between 10m and 24h in seconds format. If unspecified, a default timeout of 1h is used.",
											ServerDefault: true,
										},
										"serviceAccount": &dcl.Property{
											Type:          "string",
											GoName:        "ServiceAccount",
											Description:   "Optional. Google service account to use for execution. If unspecified, the project execution service account (-compute@developer.gserviceaccount.com) is used.",
											ServerDefault: true,
										},
										"usages": &dcl.Property{
											Type:        "array",
											GoName:      "Usages",
											Description: "Required. Usages when this configuration should be applied.",
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "string",
												GoType: "TargetExecutionConfigsUsagesEnum",
												Enum: []string{
													"EXECUTION_ENVIRONMENT_USAGE_UNSPECIFIED",
													"RENDER",
													"DEPLOY",
												},
											},
										},
										"workerPool": &dcl.Property{
											Type:        "string",
											GoName:      "WorkerPool",
											Description: "Optional. The resource name of the `WorkerPool`, with the format `projects/{project}/locations/{location}/workerPools/{worker_pool}`. If this optional field is unspecified, the default Cloud Build pool will be used.",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Cloudbuild/WorkerPool",
													Field:    "selfLink",
												},
											},
										},
									},
								},
							},
							"gke": &dcl.Property{
								Type:        "object",
								GoName:      "Gke",
								GoType:      "TargetGke",
								Description: "Information specifying a GKE Cluster.",
								Conflicts: []string{
									"anthosCluster",
								},
								Properties: map[string]*dcl.Property{
									"cluster": &dcl.Property{
										Type:        "string",
										GoName:      "Cluster",
										Description: "Information specifying a GKE Cluster. Format is `projects/{project_id}/locations/{location_id}/clusters/{cluster_id}.",
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Container/Cluster",
												Field:    "selfLink",
											},
										},
									},
									"internalIP": &dcl.Property{
										Type:        "boolean",
										GoName:      "InternalIP",
										Description: "Optional. If true, `cluster` is accessed using the private IP address of the control plane endpoint. Otherwise, the default IP address of the control plane endpoint is used. The default IP address is the private IP address for clusters with private control-plane endpoints and the public IP address otherwise. Only specify this option when `cluster` is a [private GKE cluster](https://cloud.google.com/kubernetes-engine/docs/concepts/private-cluster-concept).",
									},
								},
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Optional. Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.",
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
								Description: "Name of the `Target`. Format is [a-z][a-z0-9\\-]{0,62}.",
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
							"requireApproval": &dcl.Property{
								Type:        "boolean",
								GoName:      "RequireApproval",
								Description: "Optional. Whether or not the `Target` requires approval.",
							},
							"targetId": &dcl.Property{
								Type:        "string",
								GoName:      "TargetId",
								ReadOnly:    true,
								Description: "Output only. Resource id of the `Target`.",
								Immutable:   true,
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. Unique identifier of the `Target`.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Most recent time at which the `Target` was updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
