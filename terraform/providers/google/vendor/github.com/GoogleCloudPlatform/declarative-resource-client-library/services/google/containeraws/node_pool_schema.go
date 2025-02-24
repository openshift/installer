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
package containeraws

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLNodePoolSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "ContainerAws/NodePool",
			Description: "An Anthos node pool running on AWS.",
			StructName:  "NodePool",
			Reference: &dcl.Link{
				Text: "API reference",
				URL:  "https://cloud.google.com/kubernetes-engine/multi-cloud/docs/reference/rest/v1/projects.locations.awsClusters.awsNodePools",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Multicloud overview",
					URL:  "https://cloud.google.com/kubernetes-engine/multi-cloud/docs",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a NodePool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "nodePool",
						Required:    true,
						Description: "A full instance of a NodePool",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a NodePool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "nodePool",
						Required:    true,
						Description: "A full instance of a NodePool",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a NodePool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "nodePool",
						Required:    true,
						Description: "A full instance of a NodePool",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all NodePool",
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
					dcl.PathParameters{
						Name:     "cluster",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many NodePool",
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
					dcl.PathParameters{
						Name:     "cluster",
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
				"NodePool": &dcl.Component{
					Title:           "NodePool",
					ID:              "projects/{{project}}/locations/{{location}}/awsClusters/{{cluster}}/awsNodePools/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"version",
							"config",
							"autoscaling",
							"subnetId",
							"maxPodsConstraint",
							"project",
							"location",
							"cluster",
						},
						Properties: map[string]*dcl.Property{
							"annotations": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Annotations",
								Description: "Optional. Annotations on the node pool. This field has the same restrictions as Kubernetes annotations. The total size of all keys and values combined is limited to 256k. Key can have 2 segments: prefix (optional) and name (required), separated by a slash (/). Prefix must be a DNS subdomain. Name must be 63 characters or less, begin and end with alphanumerics, with dashes (-), underscores (_), dots (.), and alphanumerics between.",
							},
							"autoscaling": &dcl.Property{
								Type:        "object",
								GoName:      "Autoscaling",
								GoType:      "NodePoolAutoscaling",
								Description: "Autoscaler configuration for this node pool.",
								Required: []string{
									"minNodeCount",
									"maxNodeCount",
								},
								Properties: map[string]*dcl.Property{
									"maxNodeCount": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MaxNodeCount",
										Description: "Maximum number of nodes in the NodePool. Must be >= min_node_count.",
									},
									"minNodeCount": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MinNodeCount",
										Description: "Minimum number of nodes in the NodePool. Must be >= 1 and <= max_node_count.",
									},
								},
							},
							"cluster": &dcl.Property{
								Type:        "string",
								GoName:      "Cluster",
								Description: "The awsCluster for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Gkemulticloud/Cluster",
										Field:    "name",
										Parent:   true,
									},
								},
								Parameter: true,
							},
							"config": &dcl.Property{
								Type:        "object",
								GoName:      "Config",
								GoType:      "NodePoolConfig",
								Description: "The configuration of the node pool.",
								Required: []string{
									"iamInstanceProfile",
									"configEncryption",
								},
								Properties: map[string]*dcl.Property{
									"autoscalingMetricsCollection": &dcl.Property{
										Type:        "object",
										GoName:      "AutoscalingMetricsCollection",
										GoType:      "NodePoolConfigAutoscalingMetricsCollection",
										Description: "Optional. Configuration related to CloudWatch metrics collection on the Auto Scaling group of the node pool. When unspecified, metrics collection is disabled.",
										Required: []string{
											"granularity",
										},
										Properties: map[string]*dcl.Property{
											"granularity": &dcl.Property{
												Type:        "string",
												GoName:      "Granularity",
												Description: "The frequency at which EC2 Auto Scaling sends aggregated data to AWS CloudWatch. The only valid value is \"1Minute\".",
											},
											"metrics": &dcl.Property{
												Type:        "array",
												GoName:      "Metrics",
												Description: "The metrics to enable. For a list of valid metrics, see https://docs.aws.amazon.com/autoscaling/ec2/APIReference/API_EnableMetricsCollection.html. If you specify granularity and don't specify any metrics, all metrics are enabled.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
										},
									},
									"configEncryption": &dcl.Property{
										Type:        "object",
										GoName:      "ConfigEncryption",
										GoType:      "NodePoolConfigConfigEncryption",
										Description: "The ARN of the AWS KMS key used to encrypt node pool configuration.",
										Required: []string{
											"kmsKeyArn",
										},
										Properties: map[string]*dcl.Property{
											"kmsKeyArn": &dcl.Property{
												Type:        "string",
												GoName:      "KmsKeyArn",
												Description: "The ARN of the AWS KMS key used to encrypt node pool configuration.",
											},
										},
									},
									"iamInstanceProfile": &dcl.Property{
										Type:        "string",
										GoName:      "IamInstanceProfile",
										Description: "The name of the AWS IAM role assigned to nodes in the pool.",
									},
									"instanceType": &dcl.Property{
										Type:          "string",
										GoName:        "InstanceType",
										Description:   "Optional. The AWS instance type. When unspecified, it defaults to `m5.large`.",
										ServerDefault: true,
									},
									"labels": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:      "Labels",
										Description: "Optional. The initial labels assigned to nodes of this node pool. An object containing a list of \"key\": value pairs. Example: { \"name\": \"wrench\", \"mass\": \"1.3kg\", \"count\": \"3\" }.",
									},
									"proxyConfig": &dcl.Property{
										Type:        "object",
										GoName:      "ProxyConfig",
										GoType:      "NodePoolConfigProxyConfig",
										Description: "Proxy configuration for outbound HTTP(S) traffic.",
										Required: []string{
											"secretArn",
											"secretVersion",
										},
										Properties: map[string]*dcl.Property{
											"secretArn": &dcl.Property{
												Type:        "string",
												GoName:      "SecretArn",
												Description: "The ARN of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.",
											},
											"secretVersion": &dcl.Property{
												Type:        "string",
												GoName:      "SecretVersion",
												Description: "The version string of the AWS Secret Manager secret that contains the HTTP(S) proxy configuration.",
											},
										},
									},
									"rootVolume": &dcl.Property{
										Type:          "object",
										GoName:        "RootVolume",
										GoType:        "NodePoolConfigRootVolume",
										Description:   "Optional. Template for the root volume provisioned for node pool nodes. Volumes will be provisioned in the availability zone assigned to the node pool subnet. When unspecified, it defaults to 32 GiB with the GP2 volume type.",
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"iops": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "Iops",
												Description:   "Optional. The number of I/O operations per second (IOPS) to provision for GP3 volume.",
												ServerDefault: true,
											},
											"kmsKeyArn": &dcl.Property{
												Type:        "string",
												GoName:      "KmsKeyArn",
												Description: "Optional. The Amazon Resource Name (ARN) of the Customer Managed Key (CMK) used to encrypt AWS EBS volumes. If not specified, the default Amazon managed key associated to the AWS region where this cluster runs will be used.",
											},
											"sizeGib": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "SizeGib",
												Description:   "Optional. The size of the volume, in GiBs. When unspecified, a default value is provided. See the specific reference in the parent resource.",
												ServerDefault: true,
											},
											"throughput": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "Throughput",
												Description:   "Optional. The throughput to provision for the volume, in MiB/s. Only valid if the volume type is GP3. If volume type is gp3 and throughput is not specified, the throughput will defaults to 125.",
												ServerDefault: true,
											},
											"volumeType": &dcl.Property{
												Type:          "string",
												GoName:        "VolumeType",
												GoType:        "NodePoolConfigRootVolumeVolumeTypeEnum",
												Description:   "Optional. Type of the EBS volume. When unspecified, it defaults to GP2 volume. Possible values: VOLUME_TYPE_UNSPECIFIED, GP2, GP3",
												ServerDefault: true,
												Enum: []string{
													"VOLUME_TYPE_UNSPECIFIED",
													"GP2",
													"GP3",
												},
											},
										},
									},
									"securityGroupIds": &dcl.Property{
										Type:        "array",
										GoName:      "SecurityGroupIds",
										Description: "Optional. The IDs of additional security groups to add to nodes in this pool. The manager will automatically create security groups with minimum rules needed for a functioning cluster.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"sshConfig": &dcl.Property{
										Type:        "object",
										GoName:      "SshConfig",
										GoType:      "NodePoolConfigSshConfig",
										Description: "Optional. The SSH configuration.",
										Required: []string{
											"ec2KeyPair",
										},
										Properties: map[string]*dcl.Property{
											"ec2KeyPair": &dcl.Property{
												Type:        "string",
												GoName:      "Ec2KeyPair",
												Description: "The name of the EC2 key pair used to login into cluster machines.",
											},
										},
									},
									"tags": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:      "Tags",
										Description: "Optional. Key/value metadata to assign to each underlying AWS resource. Specify at most 50 pairs containing alphanumerics, spaces, and symbols (.+-=_:@/). Keys can be up to 127 Unicode characters. Values can be up to 255 Unicode characters.",
									},
									"taints": &dcl.Property{
										Type:        "array",
										GoName:      "Taints",
										Description: "Optional. The initial taints assigned to nodes of this node pool.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "NodePoolConfigTaints",
											Required: []string{
												"key",
												"value",
												"effect",
											},
											Properties: map[string]*dcl.Property{
												"effect": &dcl.Property{
													Type:        "string",
													GoName:      "Effect",
													GoType:      "NodePoolConfigTaintsEffectEnum",
													Description: "The taint effect. Possible values: EFFECT_UNSPECIFIED, NO_SCHEDULE, PREFER_NO_SCHEDULE, NO_EXECUTE",
													Immutable:   true,
													Enum: []string{
														"EFFECT_UNSPECIFIED",
														"NO_SCHEDULE",
														"PREFER_NO_SCHEDULE",
														"NO_EXECUTE",
													},
												},
												"key": &dcl.Property{
													Type:        "string",
													GoName:      "Key",
													Description: "Key for the taint.",
													Immutable:   true,
												},
												"value": &dcl.Property{
													Type:        "string",
													GoName:      "Value",
													Description: "Value for the taint.",
													Immutable:   true,
												},
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
								Description: "Output only. The time at which this node pool was created.",
								Immutable:   true,
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Allows clients to perform consistent read-modify-writes through optimistic concurrency control. May be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"management": &dcl.Property{
								Type:        "object",
								GoName:      "Management",
								GoType:      "NodePoolManagement",
								Description: "The Management configuration for this node pool.",
								Properties: map[string]*dcl.Property{
									"autoRepair": &dcl.Property{
										Type:        "boolean",
										GoName:      "AutoRepair",
										Description: "Optional. Whether or not the nodes will be automatically repaired.",
									},
								},
							},
							"maxPodsConstraint": &dcl.Property{
								Type:        "object",
								GoName:      "MaxPodsConstraint",
								GoType:      "NodePoolMaxPodsConstraint",
								Description: "The constraint on the maximum number of pods that can be run simultaneously on a node in the node pool.",
								Immutable:   true,
								Required: []string{
									"maxPodsPerNode",
								},
								Properties: map[string]*dcl.Property{
									"maxPodsPerNode": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MaxPodsPerNode",
										Description: "The maximum number of pods to schedule on a single node.",
										Immutable:   true,
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of this resource.",
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
							"reconciling": &dcl.Property{
								Type:        "boolean",
								GoName:      "Reconciling",
								ReadOnly:    true,
								Description: "Output only. If set, there are currently changes in flight to the node pool.",
								Immutable:   true,
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "NodePoolStateEnum",
								ReadOnly:    true,
								Description: "Output only. The lifecycle state of the node pool. Possible values: STATE_UNSPECIFIED, PROVISIONING, RUNNING, RECONCILING, STOPPING, ERROR, DEGRADED",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"PROVISIONING",
									"RUNNING",
									"RECONCILING",
									"STOPPING",
									"ERROR",
									"DEGRADED",
								},
							},
							"subnetId": &dcl.Property{
								Type:        "string",
								GoName:      "SubnetId",
								Description: "The subnet where the node pool node run.",
								Immutable:   true,
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. A globally unique identifier for the node pool.",
								Immutable:   true,
							},
							"updateSettings": &dcl.Property{
								Type:          "object",
								GoName:        "UpdateSettings",
								GoType:        "NodePoolUpdateSettings",
								Description:   "Optional. Update settings control the speed and disruption of the node pool update.",
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"surgeSettings": &dcl.Property{
										Type:          "object",
										GoName:        "SurgeSettings",
										GoType:        "NodePoolUpdateSettingsSurgeSettings",
										Description:   "Optional. Settings for surge update.",
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"maxSurge": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "MaxSurge",
												Description:   "Optional. The maximum number of nodes that can be created beyond the current size of the node pool during the update process.",
												ServerDefault: true,
											},
											"maxUnavailable": &dcl.Property{
												Type:          "integer",
												Format:        "int64",
												GoName:        "MaxUnavailable",
												Description:   "Optional. The maximum number of nodes that can be simultaneously unavailable during the update process. A node is considered unavailable if its status is not Ready.",
												ServerDefault: true,
											},
										},
									},
								},
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this node pool was last updated.",
								Immutable:   true,
							},
							"version": &dcl.Property{
								Type:        "string",
								GoName:      "Version",
								Description: "The Kubernetes version to run on this node pool (e.g. `1.19.10-gke.1000`). You can list all supported versions on a given Google Cloud region by calling GetAwsServerConfig.",
							},
						},
					},
				},
			},
		},
	}
}
