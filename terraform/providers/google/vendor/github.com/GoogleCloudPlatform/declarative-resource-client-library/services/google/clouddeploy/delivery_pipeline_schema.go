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
package clouddeploy

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLDeliveryPipelineSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Clouddeploy/DeliveryPipeline",
			Description: "The Cloud Deploy `DeliveryPipeline` resource",
			StructName:  "DeliveryPipeline",
			Reference: &dcl.Link{
				Text: "REST API",
				URL:  "https://cloud.google.com/deploy/docs/api/reference/rest/v1/projects.locations.deliveryPipelines",
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a DeliveryPipeline",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "deliveryPipeline",
						Required:    true,
						Description: "A full instance of a DeliveryPipeline",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a DeliveryPipeline",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "deliveryPipeline",
						Required:    true,
						Description: "A full instance of a DeliveryPipeline",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a DeliveryPipeline",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "deliveryPipeline",
						Required:    true,
						Description: "A full instance of a DeliveryPipeline",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all DeliveryPipeline",
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
				Description: "The function used to list information about many DeliveryPipeline",
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
				"DeliveryPipeline": &dcl.Component{
					Title:           "DeliveryPipeline",
					ID:              "projects/{{project}}/locations/{{location}}/deliveryPipelines/{{name}}",
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
								Description: "User annotations. These attributes can only be set and used by the user, and not by Google Cloud Deploy. See https://google.aip.dev/128#annotations for more details such as format and size limitations.",
							},
							"condition": &dcl.Property{
								Type:        "object",
								GoName:      "Condition",
								GoType:      "DeliveryPipelineCondition",
								ReadOnly:    true,
								Description: "Output only. Information around the state of the Delivery Pipeline.",
								Properties: map[string]*dcl.Property{
									"pipelineReadyCondition": &dcl.Property{
										Type:        "object",
										GoName:      "PipelineReadyCondition",
										GoType:      "DeliveryPipelineConditionPipelineReadyCondition",
										Description: "Details around the Pipeline's overall status.",
										Properties: map[string]*dcl.Property{
											"status": &dcl.Property{
												Type:        "boolean",
												GoName:      "Status",
												Description: "True if the Pipeline is in a valid state. Otherwise at least one condition in `PipelineCondition` is in an invalid state. Iterate over those conditions and see which condition(s) has status = false to find out what is wrong with the Pipeline.",
											},
											"updateTime": &dcl.Property{
												Type:        "string",
												Format:      "date-time",
												GoName:      "UpdateTime",
												Description: "Last time the condition was updated.",
											},
										},
									},
									"targetsPresentCondition": &dcl.Property{
										Type:        "object",
										GoName:      "TargetsPresentCondition",
										GoType:      "DeliveryPipelineConditionTargetsPresentCondition",
										Description: "Details around targets enumerated in the pipeline.",
										Properties: map[string]*dcl.Property{
											"missingTargets": &dcl.Property{
												Type:        "array",
												GoName:      "MissingTargets",
												Description: "The list of Target names that are missing. For example, projects/{project_id}/locations/{location_name}/targets/{target_name}.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
													ResourceReferences: []*dcl.PropertyResourceReference{
														&dcl.PropertyResourceReference{
															Resource: "Clouddeploy/Target",
															Field:    "selfLink",
														},
													},
												},
											},
											"status": &dcl.Property{
												Type:        "boolean",
												GoName:      "Status",
												Description: "True if there aren't any missing Targets.",
											},
											"updateTime": &dcl.Property{
												Type:        "string",
												Format:      "date-time",
												GoName:      "UpdateTime",
												Description: "Last time the condition was updated.",
											},
										},
									},
									"targetsTypeCondition": &dcl.Property{
										Type:        "object",
										GoName:      "TargetsTypeCondition",
										GoType:      "DeliveryPipelineConditionTargetsTypeCondition",
										Description: "Details on the whether the targets enumerated in the pipeline are of the same type.",
										Properties: map[string]*dcl.Property{
											"errorDetails": &dcl.Property{
												Type:        "string",
												GoName:      "ErrorDetails",
												Description: "Human readable error message.",
											},
											"status": &dcl.Property{
												Type:        "boolean",
												GoName:      "Status",
												Description: "True if the targets are all a comparable type. For example this is true if all targets are GKE clusters. This is false if some targets are Cloud Run targets and others are GKE clusters.",
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
								Description: "Output only. Time at which the pipeline was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Description of the `DeliveryPipeline`. Max length is 255 characters.",
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Labels are attributes that can be set and used by both the user and by Google Cloud Deploy. Labels must meet the following constraints: * Keys and values can contain only lowercase letters, numeric characters, underscores, and dashes. * All characters must use UTF-8 encoding, and international characters are allowed. * Keys must start with a lowercase letter or international character. * Each resource is limited to a maximum of 64 labels. Both keys and values are additionally constrained to be <= 128 bytes.",
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the `DeliveryPipeline`. Format is `[a-z]([a-z0-9-]{0,61}[a-z0-9])?`.",
								Immutable:   true,
								Parameter:   true,
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
							"serialPipeline": &dcl.Property{
								Type:        "object",
								GoName:      "SerialPipeline",
								GoType:      "DeliveryPipelineSerialPipeline",
								Description: "SerialPipeline defines a sequential set of stages for a `DeliveryPipeline`.",
								Properties: map[string]*dcl.Property{
									"stages": &dcl.Property{
										Type:        "array",
										GoName:      "Stages",
										Description: "Each stage specifies configuration for a `Target`. The ordering of this list defines the promotion flow.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "DeliveryPipelineSerialPipelineStages",
											Properties: map[string]*dcl.Property{
												"deployParameters": &dcl.Property{
													Type:        "array",
													GoName:      "DeployParameters",
													Description: "Optional. The deploy parameters to use for the target in this stage.",
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "object",
														GoType: "DeliveryPipelineSerialPipelineStagesDeployParameters",
														Required: []string{
															"values",
														},
														Properties: map[string]*dcl.Property{
															"matchTargetLabels": &dcl.Property{
																Type: "object",
																AdditionalProperties: &dcl.Property{
																	Type: "string",
																},
																GoName:      "MatchTargetLabels",
																Description: "Optional. Deploy parameters are applied to targets with match labels. If unspecified, deploy parameters are applied to all targets (including child targets of a multi-target).",
															},
															"values": &dcl.Property{
																Type: "object",
																AdditionalProperties: &dcl.Property{
																	Type: "string",
																},
																GoName:      "Values",
																Description: "Required. Values are deploy parameters in key-value pairs.",
															},
														},
													},
												},
												"profiles": &dcl.Property{
													Type:        "array",
													GoName:      "Profiles",
													Description: "Skaffold profiles to use when rendering the manifest for this stage's `Target`.",
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"strategy": &dcl.Property{
													Type:        "object",
													GoName:      "Strategy",
													GoType:      "DeliveryPipelineSerialPipelineStagesStrategy",
													Description: "Optional. The strategy to use for a `Rollout` to this stage.",
													Properties: map[string]*dcl.Property{
														"canary": &dcl.Property{
															Type:        "object",
															GoName:      "Canary",
															GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanary",
															Description: "Canary deployment strategy provides progressive percentage based deployments to a Target.",
															Properties: map[string]*dcl.Property{
																"canaryDeployment": &dcl.Property{
																	Type:        "object",
																	GoName:      "CanaryDeployment",
																	GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeployment",
																	Description: "Configures the progressive based deployment for a Target.",
																	Conflicts: []string{
																		"customCanaryDeployment",
																	},
																	Required: []string{
																		"percentages",
																	},
																	Properties: map[string]*dcl.Property{
																		"percentages": &dcl.Property{
																			Type:        "array",
																			GoName:      "Percentages",
																			Description: "Required. The percentage based deployments that will occur as a part of a `Rollout`. List is expected in ascending order and each integer n is 0 <= n < 100.",
																			SendEmpty:   true,
																			ListType:    "list",
																			Items: &dcl.Property{
																				Type:   "integer",
																				Format: "int64",
																				GoType: "int64",
																			},
																		},
																		"postdeploy": &dcl.Property{
																			Type:        "object",
																			GoName:      "Postdeploy",
																			GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPostdeploy",
																			Description: "Optional. Configuration for the postdeploy job of the last phase. If this is not configured, postdeploy job will not be present.",
																			Properties: map[string]*dcl.Property{
																				"actions": &dcl.Property{
																					Type:        "array",
																					GoName:      "Actions",
																					Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the postdeploy job.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																			},
																		},
																		"predeploy": &dcl.Property{
																			Type:        "object",
																			GoName:      "Predeploy",
																			GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryCanaryDeploymentPredeploy",
																			Description: "Optional. Configuration for the predeploy job of the first phase. If this is not configured, predeploy job will not be present.",
																			Properties: map[string]*dcl.Property{
																				"actions": &dcl.Property{
																					Type:        "array",
																					GoName:      "Actions",
																					Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the predeploy job.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																			},
																		},
																		"verify": &dcl.Property{
																			Type:        "boolean",
																			GoName:      "Verify",
																			Description: "Whether to run verify tests after each percentage deployment.",
																		},
																	},
																},
																"customCanaryDeployment": &dcl.Property{
																	Type:        "object",
																	GoName:      "CustomCanaryDeployment",
																	GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeployment",
																	Description: "Configures the progressive based deployment for a Target, but allows customizing at the phase level where a phase represents each of the percentage deployments.",
																	Conflicts: []string{
																		"canaryDeployment",
																	},
																	Required: []string{
																		"phaseConfigs",
																	},
																	Properties: map[string]*dcl.Property{
																		"phaseConfigs": &dcl.Property{
																			Type:        "array",
																			GoName:      "PhaseConfigs",
																			Description: "Required. Configuration for each phase in the canary deployment in the order executed.",
																			SendEmpty:   true,
																			ListType:    "list",
																			Items: &dcl.Property{
																				Type:   "object",
																				GoType: "DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigs",
																				Required: []string{
																					"phaseId",
																					"percentage",
																				},
																				Properties: map[string]*dcl.Property{
																					"percentage": &dcl.Property{
																						Type:        "integer",
																						Format:      "int64",
																						GoName:      "Percentage",
																						Description: "Required. Percentage deployment for the phase.",
																					},
																					"phaseId": &dcl.Property{
																						Type:        "string",
																						GoName:      "PhaseId",
																						Description: "Required. The ID to assign to the `Rollout` phase. This value must consist of lower-case letters, numbers, and hyphens, start with a letter and end with a letter or a number, and have a max length of 63 characters. In other words, it must match the following regex: `^[a-z]([a-z0-9-]{0,61}[a-z0-9])?$`.",
																					},
																					"postdeploy": &dcl.Property{
																						Type:        "object",
																						GoName:      "Postdeploy",
																						GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPostdeploy",
																						Description: "Optional. Configuration for the postdeploy job of this phase. If this is not configured, postdeploy job will not be present for this phase.",
																						Properties: map[string]*dcl.Property{
																							"actions": &dcl.Property{
																								Type:        "array",
																								GoName:      "Actions",
																								Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the postdeploy job.",
																								SendEmpty:   true,
																								ListType:    "list",
																								Items: &dcl.Property{
																									Type:   "string",
																									GoType: "string",
																								},
																							},
																						},
																					},
																					"predeploy": &dcl.Property{
																						Type:        "object",
																						GoName:      "Predeploy",
																						GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryCustomCanaryDeploymentPhaseConfigsPredeploy",
																						Description: "Optional. Configuration for the predeploy job of this phase. If this is not configured, predeploy job will not be present for this phase.",
																						Properties: map[string]*dcl.Property{
																							"actions": &dcl.Property{
																								Type:        "array",
																								GoName:      "Actions",
																								Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the predeploy job.",
																								SendEmpty:   true,
																								ListType:    "list",
																								Items: &dcl.Property{
																									Type:   "string",
																									GoType: "string",
																								},
																							},
																						},
																					},
																					"profiles": &dcl.Property{
																						Type:        "array",
																						GoName:      "Profiles",
																						Description: "Skaffold profiles to use when rendering the manifest for this phase. These are in addition to the profiles list specified in the `DeliveryPipeline` stage.",
																						SendEmpty:   true,
																						ListType:    "list",
																						Items: &dcl.Property{
																							Type:   "string",
																							GoType: "string",
																						},
																					},
																					"verify": &dcl.Property{
																						Type:        "boolean",
																						GoName:      "Verify",
																						Description: "Whether to run verify tests after the deployment.",
																					},
																				},
																			},
																		},
																	},
																},
																"runtimeConfig": &dcl.Property{
																	Type:        "object",
																	GoName:      "RuntimeConfig",
																	GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfig",
																	Description: "Optional. Runtime specific configurations for the deployment strategy. The runtime configuration is used to determine how Cloud Deploy will split traffic to enable a progressive deployment.",
																	Properties: map[string]*dcl.Property{
																		"cloudRun": &dcl.Property{
																			Type:        "object",
																			GoName:      "CloudRun",
																			GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigCloudRun",
																			Description: "Cloud Run runtime configuration.",
																			Conflicts: []string{
																				"kubernetes",
																			},
																			Properties: map[string]*dcl.Property{
																				"automaticTrafficControl": &dcl.Property{
																					Type:        "boolean",
																					GoName:      "AutomaticTrafficControl",
																					Description: "Whether Cloud Deploy should update the traffic stanza in a Cloud Run Service on the user's behalf to facilitate traffic splitting. This is required to be true for CanaryDeployments, but optional for CustomCanaryDeployments.",
																				},
																				"canaryRevisionTags": &dcl.Property{
																					Type:        "array",
																					GoName:      "CanaryRevisionTags",
																					Description: "Optional. A list of tags that are added to the canary revision while the canary phase is in progress.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"priorRevisionTags": &dcl.Property{
																					Type:        "array",
																					GoName:      "PriorRevisionTags",
																					Description: "Optional. A list of tags that are added to the prior revision while the canary phase is in progress.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"stableRevisionTags": &dcl.Property{
																					Type:        "array",
																					GoName:      "StableRevisionTags",
																					Description: "Optional. A list of tags that are added to the final stable revision when the stable phase is applied.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																			},
																		},
																		"kubernetes": &dcl.Property{
																			Type:        "object",
																			GoName:      "Kubernetes",
																			GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetes",
																			Description: "Kubernetes runtime configuration.",
																			Conflicts: []string{
																				"cloudRun",
																			},
																			Properties: map[string]*dcl.Property{
																				"gatewayServiceMesh": &dcl.Property{
																					Type:        "object",
																					GoName:      "GatewayServiceMesh",
																					GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesGatewayServiceMesh",
																					Description: "Kubernetes Gateway API service mesh configuration.",
																					Conflicts: []string{
																						"serviceNetworking",
																					},
																					Required: []string{
																						"httpRoute",
																						"service",
																						"deployment",
																					},
																					Properties: map[string]*dcl.Property{
																						"deployment": &dcl.Property{
																							Type:        "string",
																							GoName:      "Deployment",
																							Description: "Required. Name of the Kubernetes Deployment whose traffic is managed by the specified HTTPRoute and Service.",
																						},
																						"httpRoute": &dcl.Property{
																							Type:        "string",
																							GoName:      "HttpRoute",
																							Description: "Required. Name of the Gateway API HTTPRoute.",
																						},
																						"podSelectorLabel": &dcl.Property{
																							Type:        "string",
																							GoName:      "PodSelectorLabel",
																							Description: "Optional. The label to use when selecting Pods for the Deployment and Service resources. This label must already be present in both resources.",
																						},
																						"routeUpdateWaitTime": &dcl.Property{
																							Type:        "string",
																							GoName:      "RouteUpdateWaitTime",
																							Description: "Optional. The time to wait for route updates to propagate. The maximum configurable time is 3 hours, in seconds format. If unspecified, there is no wait time.",
																						},
																						"service": &dcl.Property{
																							Type:        "string",
																							GoName:      "Service",
																							Description: "Required. Name of the Kubernetes Service.",
																						},
																						"stableCutbackDuration": &dcl.Property{
																							Type:        "string",
																							GoName:      "StableCutbackDuration",
																							Description: "Optional. The amount of time to migrate traffic back from the canary Service to the original Service during the stable phase deployment. If specified, must be between 15s and 3600s. If unspecified, there is no cutback time.",
																						},
																					},
																				},
																				"serviceNetworking": &dcl.Property{
																					Type:        "object",
																					GoName:      "ServiceNetworking",
																					GoType:      "DeliveryPipelineSerialPipelineStagesStrategyCanaryRuntimeConfigKubernetesServiceNetworking",
																					Description: "Kubernetes Service networking configuration.",
																					Conflicts: []string{
																						"gatewayServiceMesh",
																					},
																					Required: []string{
																						"service",
																						"deployment",
																					},
																					Properties: map[string]*dcl.Property{
																						"deployment": &dcl.Property{
																							Type:        "string",
																							GoName:      "Deployment",
																							Description: "Required. Name of the Kubernetes Deployment whose traffic is managed by the specified Service.",
																						},
																						"disablePodOverprovisioning": &dcl.Property{
																							Type:        "boolean",
																							GoName:      "DisablePodOverprovisioning",
																							Description: "Optional. Whether to disable Pod overprovisioning. If Pod overprovisioning is disabled then Cloud Deploy will limit the number of total Pods used for the deployment strategy to the number of Pods the Deployment has on the cluster.",
																						},
																						"podSelectorLabel": &dcl.Property{
																							Type:        "string",
																							GoName:      "PodSelectorLabel",
																							Description: "Optional. The label to use when selecting Pods for the Deployment resource. This label must already be present in the Deployment.",
																						},
																						"service": &dcl.Property{
																							Type:        "string",
																							GoName:      "Service",
																							Description: "Required. Name of the Kubernetes Service.",
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
														"standard": &dcl.Property{
															Type:        "object",
															GoName:      "Standard",
															GoType:      "DeliveryPipelineSerialPipelineStagesStrategyStandard",
															Description: "Standard deployment strategy executes a single deploy and allows verifying the deployment.",
															Properties: map[string]*dcl.Property{
																"postdeploy": &dcl.Property{
																	Type:        "object",
																	GoName:      "Postdeploy",
																	GoType:      "DeliveryPipelineSerialPipelineStagesStrategyStandardPostdeploy",
																	Description: "Optional. Configuration for the postdeploy job. If this is not configured, postdeploy job will not be present.",
																	Properties: map[string]*dcl.Property{
																		"actions": &dcl.Property{
																			Type:        "array",
																			GoName:      "Actions",
																			Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the postdeploy job.",
																			SendEmpty:   true,
																			ListType:    "list",
																			Items: &dcl.Property{
																				Type:   "string",
																				GoType: "string",
																			},
																		},
																	},
																},
																"predeploy": &dcl.Property{
																	Type:        "object",
																	GoName:      "Predeploy",
																	GoType:      "DeliveryPipelineSerialPipelineStagesStrategyStandardPredeploy",
																	Description: "Optional. Configuration for the predeploy job. If this is not configured, predeploy job will not be present.",
																	Properties: map[string]*dcl.Property{
																		"actions": &dcl.Property{
																			Type:        "array",
																			GoName:      "Actions",
																			Description: "Optional. A sequence of skaffold custom actions to invoke during execution of the predeploy job.",
																			SendEmpty:   true,
																			ListType:    "list",
																			Items: &dcl.Property{
																				Type:   "string",
																				GoType: "string",
																			},
																		},
																	},
																},
																"verify": &dcl.Property{
																	Type:        "boolean",
																	GoName:      "Verify",
																	Description: "Whether to verify a deployment.",
																},
															},
														},
													},
												},
												"targetId": &dcl.Property{
													Type:        "string",
													GoName:      "TargetId",
													Description: "The target_id to which this stage points. This field refers exclusively to the last segment of a target name. For example, this field would just be `my-target` (rather than `projects/project/locations/location/targets/my-target`). The location of the `Target` is inferred to be the same as the location of the `DeliveryPipeline` that contains this `Stage`.",
												},
											},
										},
									},
								},
							},
							"suspended": &dcl.Property{
								Type:        "boolean",
								GoName:      "Suspended",
								Description: "When suspended, no new releases or rollouts can be created, but in-progress ones will complete.",
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. Unique identifier of the `DeliveryPipeline`.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Most recent time at which the pipeline was updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
