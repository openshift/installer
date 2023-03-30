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
package dataplex

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLLakeSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Dataplex/Lake",
			Description: "The Dataplex Lake resource",
			StructName:  "Lake",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Lake",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "lake",
						Required:    true,
						Description: "A full instance of a Lake",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Lake",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "lake",
						Required:    true,
						Description: "A full instance of a Lake",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Lake",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "lake",
						Required:    true,
						Description: "A full instance of a Lake",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Lake",
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
				Description: "The function used to list information about many Lake",
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
				"Lake": &dcl.Component{
					Title:           "Lake",
					ID:              "projects/{{project}}/locations/{{location}}/lakes/{{name}}",
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
							"assetStatus": &dcl.Property{
								Type:        "object",
								GoName:      "AssetStatus",
								GoType:      "LakeAssetStatus",
								ReadOnly:    true,
								Description: "Output only. Aggregated status of the underlying assets of the lake.",
								Properties: map[string]*dcl.Property{
									"activeAssets": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "ActiveAssets",
										Description: "Number of active assets.",
									},
									"securityPolicyApplyingAssets": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "SecurityPolicyApplyingAssets",
										Description: "Number of assets that are in process of updating the security policy on attached resources.",
									},
									"updateTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "UpdateTime",
										Description: "Last update time of the status.",
									},
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time when the lake was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. Description of the lake.",
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "Optional. User friendly display name.",
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Optional. User-defined labels for the lake.",
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
							},
							"metastore": &dcl.Property{
								Type:        "object",
								GoName:      "Metastore",
								GoType:      "LakeMetastore",
								Description: "Optional. Settings to manage lake and Dataproc Metastore service instance association.",
								Properties: map[string]*dcl.Property{
									"service": &dcl.Property{
										Type:        "string",
										GoName:      "Service",
										Description: "Optional. A relative reference to the Dataproc Metastore (https://cloud.google.com/dataproc-metastore/docs) service associated with the lake: `projects/{project_id}/locations/{location_id}/services/{service_id}`",
									},
								},
							},
							"metastoreStatus": &dcl.Property{
								Type:        "object",
								GoName:      "MetastoreStatus",
								GoType:      "LakeMetastoreStatus",
								ReadOnly:    true,
								Description: "Output only. Metastore status of the lake.",
								Properties: map[string]*dcl.Property{
									"endpoint": &dcl.Property{
										Type:        "string",
										GoName:      "Endpoint",
										Description: "The URI of the endpoint used to access the Metastore service.",
									},
									"message": &dcl.Property{
										Type:        "string",
										GoName:      "Message",
										Description: "Additional information about the current status.",
									},
									"state": &dcl.Property{
										Type:        "string",
										GoName:      "State",
										GoType:      "LakeMetastoreStatusStateEnum",
										Description: "Current state of association. Possible values: STATE_UNSPECIFIED, NONE, READY, UPDATING, ERROR",
										Enum: []string{
											"STATE_UNSPECIFIED",
											"NONE",
											"READY",
											"UPDATING",
											"ERROR",
										},
									},
									"updateTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "UpdateTime",
										Description: "Last update time of the metastore status of the lake.",
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of the lake.",
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Dataplex/Lake",
										Field:    "selfLink",
										Parent:   true,
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
							"serviceAccount": &dcl.Property{
								Type:        "string",
								GoName:      "ServiceAccount",
								ReadOnly:    true,
								Description: "Output only. Service account associated with this lake. This service account must be authorized to access or operate on resources managed by the lake.",
								Immutable:   true,
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "LakeStateEnum",
								ReadOnly:    true,
								Description: "Output only. Current state of the lake. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"ACTIVE",
									"CREATING",
									"DELETING",
									"ACTION_REQUIRED",
								},
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. System generated globally unique ID for the lake. This ID will be different if the lake is deleted and re-created with the same name.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time when the lake was last updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
