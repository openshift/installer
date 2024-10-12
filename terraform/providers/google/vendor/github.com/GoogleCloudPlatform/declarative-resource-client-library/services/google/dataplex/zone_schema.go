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

func DCLZoneSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Dataplex/Zone",
			Description: "The Dataplex Zone resource",
			StructName:  "Zone",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Zone",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "zone",
						Required:    true,
						Description: "A full instance of a Zone",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Zone",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "zone",
						Required:    true,
						Description: "A full instance of a Zone",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Zone",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "zone",
						Required:    true,
						Description: "A full instance of a Zone",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Zone",
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
						Name:     "lake",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Zone",
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
				"Zone": &dcl.Component{
					Title:           "Zone",
					ID:              "projects/{{project}}/locations/{{location}}/lakes/{{lake}}/zones/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"type",
							"discoverySpec",
							"resourceSpec",
							"project",
							"location",
							"lake",
						},
						Properties: map[string]*dcl.Property{
							"assetStatus": &dcl.Property{
								Type:        "object",
								GoName:      "AssetStatus",
								GoType:      "ZoneAssetStatus",
								ReadOnly:    true,
								Description: "Output only. Aggregated status of the underlying assets of the zone.",
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
								Description: "Output only. The time when the zone was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. Description of the zone.",
							},
							"discoverySpec": &dcl.Property{
								Type:        "object",
								GoName:      "DiscoverySpec",
								GoType:      "ZoneDiscoverySpec",
								Description: "Required. Specification of the discovery feature applied to data in this zone.",
								Required: []string{
									"enabled",
								},
								Properties: map[string]*dcl.Property{
									"csvOptions": &dcl.Property{
										Type:          "object",
										GoName:        "CsvOptions",
										GoType:        "ZoneDiscoverySpecCsvOptions",
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
										GoType:        "ZoneDiscoverySpecJsonOptions",
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
										Type:          "string",
										GoName:        "Schedule",
										Description:   "Optional. Cron schedule (https://en.wikipedia.org/wiki/Cron) for running discovery periodically. Successive discovery runs must be scheduled at least 60 minutes apart. The default value is to run discovery every 60 minutes. To explicitly set a timezone to the cron tab, apply a prefix in the cron tab: \"CRON_TZ=${IANA_TIME_ZONE}\" or TZ=${IANA_TIME_ZONE}\". The ${IANA_TIME_ZONE} may only be a valid string from IANA time zone database. For example, \"CRON_TZ=America/New_York 1 * * * *\", or \"TZ=America/New_York 1 * * * *\".",
										ServerDefault: true,
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
								Description: "Optional. User defined labels for the zone.",
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
								Description: "The name of the zone.",
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Dataplex/Zone",
										Field:    "selfLink",
										Parent:   true,
									},
								},
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
								GoType:      "ZoneResourceSpec",
								Description: "Required. Immutable. Specification of the resources that are referenced by the assets within this zone.",
								Immutable:   true,
								Required: []string{
									"locationType",
								},
								Properties: map[string]*dcl.Property{
									"locationType": &dcl.Property{
										Type:        "string",
										GoName:      "LocationType",
										GoType:      "ZoneResourceSpecLocationTypeEnum",
										Description: "Required. Immutable. The location type of the resources that are allowed to be attached to the assets within this zone. Possible values: LOCATION_TYPE_UNSPECIFIED, SINGLE_REGION, MULTI_REGION",
										Immutable:   true,
										Enum: []string{
											"LOCATION_TYPE_UNSPECIFIED",
											"SINGLE_REGION",
											"MULTI_REGION",
										},
									},
								},
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "ZoneStateEnum",
								ReadOnly:    true,
								Description: "Output only. Current state of the zone. Possible values: STATE_UNSPECIFIED, ACTIVE, CREATING, DELETING, ACTION_REQUIRED",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"ACTIVE",
									"CREATING",
									"DELETING",
									"ACTION_REQUIRED",
								},
							},
							"type": &dcl.Property{
								Type:        "string",
								GoName:      "Type",
								GoType:      "ZoneTypeEnum",
								Description: "Required. Immutable. The type of the zone. Possible values: TYPE_UNSPECIFIED, RAW, CURATED",
								Immutable:   true,
								Enum: []string{
									"TYPE_UNSPECIFIED",
									"RAW",
									"CURATED",
								},
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. System generated globally unique ID for the zone. This ID will be different if the zone is deleted and re-created with the same name.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time when the zone was last updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
