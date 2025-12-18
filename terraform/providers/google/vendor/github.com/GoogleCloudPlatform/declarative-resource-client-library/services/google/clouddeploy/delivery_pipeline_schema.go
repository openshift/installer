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
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the `DeliveryPipeline`. Format is [a-z][a-z0-9\\-]{0,62}.",
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
