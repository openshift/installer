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
package dataplex

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLAssetSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Dataplex/Asset",
			Description: "The Dataplex Asset resource",
			StructName:  "Asset",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Asset",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "asset",
						Required:    true,
						Description: "A full instance of a Asset",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Asset",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "asset",
						Required:    true,
						Description: "A full instance of a Asset",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Asset",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "asset",
						Required:    true,
						Description: "A full instance of a Asset",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Asset",
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
						Name:     "dataplexZone",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "lake",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Asset",
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
						Name:     "dataplexZone",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "lake",
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
				"Asset": &dcl.Component{
					Title:           "Asset",
					ID:              "projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{dataplex_zone}}/assets/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"resourceSpec",
							"discoverySpec",
							"project",
							"location",
							"lake",
							"dataplexZone",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time when the asset was created.",
								Immutable:   true,
							},
							"dataplexZone": &dcl.Property{
								Type:        "string",
								GoName:      "DataplexZone",
								Description: "The zone for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. Description of the asset.",
							},
							"discoverySpec": &dcl.Property{
								Type:        "object",
								GoName:      "DiscoverySpec",
								GoType:      "AssetDiscoverySpec",
								Description: "Required. Specification of the discovery feature applied to data referenced by this asset. When this spec is left unset, the asset will use the spec set on the parent zone.",
								Required: []string{
									"enabled",
								},
								Properties: map[string]*dcl.Property{
									"csvOptions": &dcl.Property{
										Type:          "object",
										GoName:        "CsvOptions",
										GoType:        "AssetDiscoverySpecCsvOptions",
										Description:   "Optional. Configuration for CSV data.",
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"delimiter": &dcl.Property{
												Type:        "string",
												GoName:      "Delimiter",
												Description: "Optional. The delimiter being used to separate values. This defaults to ','.",
											},
											"disableTypeInference": &dcl.Property{
												Type:        "boolean",
												GoName:      "DisableTypeInference",
												Description: "Optional. Whether to disable the inference of data type for CSV data. If true, all columns will be registered as strings.",
											},
											"encoding": &dcl.Property{
												Type:        "string",
												GoName:      "Encoding",
												Description: "Optional. The character encoding of the data. The default is UTF-8.",
											},
											"headerRows": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "HeaderRows",
												Description: "Optional. The number of rows to interpret as header rows that should be skipped when reading data rows.",
											},
										},
									},
									"enabled": &dcl.Property{
										Type:        "boolean",
										GoName:      "Enabled",
										Description: "Required. Whether discovery is enabled.",
									},
									"excludePatterns": &dcl.Property{
										Type:        "array",
										GoName:      "ExcludePatterns",
										Description: "Optional. The list of patterns to apply for selecting data to exclude during discovery. For Cloud Storage bucket assets, these are interpreted as glob patterns used to match object names. For BigQuery dataset assets, these are interpreted as patterns to match table names.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"includePatterns": &dcl.Property{
										Type:        "array",
										GoName:      "IncludePatterns",
										Description: "Optional. The list of patterns to apply for selecting data to include during discovery if only a subset of the data should considered. For Cloud Storage bucket assets, these are interpreted as glob patterns used to match object names. For BigQuery dataset assets, these are interpreted as patterns to match table names.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"jsonOptions": &dcl.Property{
										Type:          "object",
										GoName:        "JsonOptions",
										GoType:        "AssetDiscoverySpecJsonOptions",
										Description:   "Optional. Configuration for Json data.",
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"disableTypeInference": &dcl.Property{
												Type:        "boolean",
												GoName:      "DisableTypeInference",
												Description: "Optional. Whether to disable the inference of data type for Json data. If true, all columns will be registered as their primitive types (strings, number or boolean).",
											},
											"encoding": &dcl.Property{
												Type:        "string",
												GoName:      "Encoding",
												Description: "Optional. The character encoding of the data. The default is UTF-8.",
											},
										},
									},
									"schedule": &dcl.Property{
										Type:        "string",
										GoName:      "Schedule",
										Description: "Optional. Cron schedule (https://en.wikipedia.org/wiki/Cron) for running discovery periodically. Successive discovery runs must be scheduled at least 60 minutes apart. The default value is to run discovery every 60 minutes. To explicitly set a timezone to the cron tab, apply a prefix in the cron tab: \"CRON_TZ=${IANA_TIME_ZONE}\" or TZ=${IANA_TIME_ZONE}\". The ${IANA_TIME_ZONE} may only be a valid string from IANA time zone database. For example, \"CRON_TZ=America/New_York 1 * * * *\", or \"TZ=America/New_York 1 * * * *\".",
									},
								},
							},
							"discoveryStatus": &dcl.Property{
								Type:        "object",
								GoName:      "DiscoveryStatus",
								GoType:      "AssetDiscoveryStatus",
								ReadOnly:    true,
								Description: "Output only. Status of the discovery feature applied to data referenced by this asset.",
								Properties: map[string]*dcl.Property{
									"lastRunDuration": &dcl.Property{
										Type:        "string",
										GoName:      "LastRunDuration",
										Description: "The duration of the last discovery run.",
									},
									"lastRunTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "LastRunTime",
										Description: "The start time of the last discovery run.",
									},
									"message": &dcl.Property{
										Type:        "string",
										GoName:      "Message",
										Description: "Additional information about the current state.",
									},
									"state": &dcl.Property{
										Type:        "string",
										GoName:      "State",
										GoType:      "AssetDiscoveryStatusStateEnum",
										Description: "The current status of the discovery feature. Possible values: STATE_UNSPECIFIED, SCHEDULED, IN_PROGRESS, PAUSED, DISABLED",
										Enum: []string{
											"STATE_UNSPECIFIED",
											"SCHEDULED",
											"IN_PROGRESS",
											"PAUSED",
											"DISABLED",
										},
									},
									"stats": &dcl.Property{
										Type:        "object",
										GoName:      "Stats",
										GoType:      "AssetDiscoveryStatusStats",
										Description: "Data Stats of the asset reported by discovery.",
										Properties: map[string]*dcl.Property{
											"dataItems": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "DataItems",
												Description: "The count of data items within the referenced resource.",
											},
											"dataSize": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "DataSize",
												Description: "The number of stored data bytes within the referenced resource.",
											},
											"filesets": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Filesets",
												Description: "The count of fileset entities within the referenced resource.",
											},
											"tables": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Tables",
												Description: "The count of table entities within the referenced resource.",
											},
										},
									},
									"updateTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "UpdateTime",
										Description: "Last update time of the status.",
									},
								},
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
								Description: "Optional. User defined labels for the asset.",
							},
							"lake": &dcl.Property{
								Type:        "string",
								GoName:      "Lake",
								Description: "The lake for the resource",
								Immutable:   true,
								Parameter:   true,
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
								Description: "The name of the asset.",
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
							"resourceSpec": &dcl.Property{
								Type:        "object",
								GoName:      "ResourceSpec",
								GoType:      "AssetResourceSpec",
								Description: "Required. Immutable. Specification of the resource that is referenced by this asset.",
								Required: []string{
									"type",
								},
								Properties: map[string]*dcl.Property{
									"name": &dcl.Property{
										Type:        "string",
										GoName:      "Name",
										Description: "Immutable. Relative name of the cloud resource that contains the data that is being managed within a lake. For example: `projects/{project_number}/buckets/{bucket_id}` `projects/{project_number}/datasets/{dataset_id}`",
										Immutable:   true,
									},
									"readAccessMode": &dcl.Property{
										Type:          "string",
										GoName:        "ReadAccessMode",
										GoType:        "AssetResourceSpecReadAccessModeEnum",
										Description:   "Optional. Determines how read permissions are handled for each asset and their associated tables. Only available to storage buckets assets. Possible values: DIRECT, MANAGED",
										ServerDefault: true,
										Enum: []string{
											"DIRECT",
											"MANAGED",
										},
									},
									"type": &dcl.Property{
										Type:        "string",
										GoName:      "Type",
										GoType:      "AssetResourceSpecTypeEnum",
										Description: "Required. Immutable. Type of resource. Possible values: STORAGE_BUCKET, BIGQUERY_DATASET",
										Immutable:   true,
										Enum: []string{
											"STORAGE_BUCKET",
											"BIGQUERY_DATASET",
										},
									},
								},
							},
							"resourceStatus": &dcl.Property{
								Type:        "object",
								GoName:      "ResourceStatus",
								GoType:      "AssetResourceStatus",
								ReadOnly:    true,
								Description: "Output only. Status of the resource referenced by this asset.",
								Properties: map[string]*dcl.Property{
									"message": &dcl.Property{
										Type:        "string",
										GoName:      "Message",
										Description: "Additional information about the current state.",
									},
									"state": &dcl.Property{
										Type:        "string",
										GoName:      "State",
										GoType:      "AssetResourceStatusStateEnum",
										Description: "The current state of the managed resource. Possible values: STATE_UNSPECIFIED, READY, ERROR",
										Enum: []string{
											"STATE_UNSPECIFIED",
											"READY",
											"ERROR",
										},
									},
									"updateTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "UpdateTime",
										Description: "Last update time of the status.",
									},
								},
							},
							"securityStatus": &dcl.Property{
								Type:        "object",
								GoName:      "SecurityStatus",
								GoType:      "AssetSecurityStatus",
								ReadOnly:    true,
								Description: "Output only. Status of the security policy applied to resource referenced by this asset.",
								Properties: map[string]*dcl.Property{
									"message": &dcl.Property{
										Type:        "string",
										GoName:      "Message",
										Description: "Additional information about the current state.",
									},
									"state": &dcl.Property{
										Type:        "string",
										GoName:      "State",
										GoType:      "AssetSecurityStatusStateEnum",
										Description: "The current state of the security policy applied to the attached resource. Possible values: STATE_UNSPECIFIED, READY, APPLYING, ERROR",
										Enum: []string{
											"STATE_UNSPECIFIED",
											"READY",
											"APPLYING",
											"ERROR",
										},
									},
									"updateTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "UpdateTime",
										Description: "Last update time of the status.",
									},
								},
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "AssetStateEnum",
								ReadOnly:    true,
								Description: "Output only. Current state of the asset. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED",
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
								Description: "Output only. System generated globally unique ID for the asset. This ID will be different if the asset is deleted and re-created with the same name.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time when the asset was last updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
