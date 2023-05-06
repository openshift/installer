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
package assuredworkloads

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLWorkloadSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "AssuredWorkloads/Workload",
			Description: "The AssuredWorkloads Workload resource",
			StructName:  "Workload",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Workload",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workload",
						Required:    true,
						Description: "A full instance of a Workload",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Workload",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workload",
						Required:    true,
						Description: "A full instance of a Workload",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Workload",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "workload",
						Required:    true,
						Description: "A full instance of a Workload",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Workload",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "organization",
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
				Description: "The function used to list information about many Workload",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "organization",
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
				"Workload": &dcl.Component{
					Title:           "Workload",
					ID:              "organizations/{{organization}}/locations/{{location}}/workloads/{{name}}",
					UsesStateHint:   true,
					ParentContainer: "organization",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"displayName",
							"complianceRegime",
							"billingAccount",
							"organization",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"billingAccount": &dcl.Property{
								Type:        "string",
								GoName:      "BillingAccount",
								Description: "Required. Input only. The billing account used for the resources which are direct children of workload. This billing account is initially associated with the resources created as part of Workload creation. After the initial creation of these resources, the customer can change the assigned billing account. The resource name has the form `billingAccounts/{billing_account_id}`. For example, 'billingAccounts/012345-567890-ABCDEF`.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/BillingAccount",
										Field:    "name",
									},
								},
								Unreadable: true,
							},
							"complianceRegime": &dcl.Property{
								Type:        "string",
								GoName:      "ComplianceRegime",
								GoType:      "WorkloadComplianceRegimeEnum",
								Description: "Required. Immutable. Compliance Regime associated with this workload. Possible values: COMPLIANCE_REGIME_UNSPECIFIED, IL4, CJIS, FEDRAMP_HIGH, FEDRAMP_MODERATE, US_REGIONAL_ACCESS",
								Immutable:   true,
								Enum: []string{
									"COMPLIANCE_REGIME_UNSPECIFIED",
									"IL4",
									"CJIS",
									"FEDRAMP_HIGH",
									"FEDRAMP_MODERATE",
									"US_REGIONAL_ACCESS",
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. Immutable. The Workload creation timestamp.",
								Immutable:   true,
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "Required. The user-assigned display name of the Workload. When present it must be between 4 to 30 characters. Allowed characters are: lowercase and uppercase letters, numbers, hyphen, and spaces. Example: My Workload",
							},
							"kmsSettings": &dcl.Property{
								Type:        "object",
								GoName:      "KmsSettings",
								GoType:      "WorkloadKmsSettings",
								Description: "Input only. Settings used to create a CMEK crypto key. When set a project with a KMS CMEK key is provisioned. This field is mandatory for a subset of Compliance Regimes.",
								Immutable:   true,
								Unreadable:  true,
								Required: []string{
									"nextRotationTime",
									"rotationPeriod",
								},
								Properties: map[string]*dcl.Property{
									"nextRotationTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "NextRotationTime",
										Description: "Required. Input only. Immutable. The time at which the Key Management Service will automatically create a new version of the crypto key and mark it as the primary.",
										Immutable:   true,
									},
									"rotationPeriod": &dcl.Property{
										Type:        "string",
										GoName:      "RotationPeriod",
										Description: "Required. Input only. Immutable. will be advanced by this period when the Key Management Service automatically rotates a key. Must be at least 24 hours and at most 876,000 hours.",
										Immutable:   true,
									},
								},
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Optional. Labels applied to the workload.",
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "Output only. The resource name of the workload.",
								Immutable:                true,
								ServerGeneratedParameter: true,
							},
							"organization": &dcl.Property{
								Type:        "string",
								GoName:      "Organization",
								Description: "The organization for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Organization",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"provisionedResourcesParent": &dcl.Property{
								Type:        "string",
								GoName:      "ProvisionedResourcesParent",
								Description: "Input only. The parent resource for the resources managed by this Assured Workload. May be either an organization or a folder. Must be the same or a child of the Workload parent. If not specified all resources are created under the Workload parent. Formats: folders/{folder_id}, organizations/{organization_id}",
								Immutable:   true,
								Unreadable:  true,
							},
							"resourceSettings": &dcl.Property{
								Type:        "array",
								GoName:      "ResourceSettings",
								Description: "Input only. Resource properties that are used to customize workload resources. These properties (such as custom project id) will be used to create workload resources if possible. This field is optional.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "WorkloadResourceSettings",
									Properties: map[string]*dcl.Property{
										"resourceId": &dcl.Property{
											Type:        "string",
											GoName:      "ResourceId",
											Description: "Resource identifier. For a project this represents project_number. If the project is already taken, the workload creation will fail.",
											Immutable:   true,
										},
										"resourceType": &dcl.Property{
											Type:        "string",
											GoName:      "ResourceType",
											GoType:      "WorkloadResourceSettingsResourceTypeEnum",
											Description: "Indicates the type of resource. This field should be specified to correspond the id to the right project type (CONSUMER_PROJECT or ENCRYPTION_KEYS_PROJECT) Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER",
											Immutable:   true,
											Enum: []string{
												"RESOURCE_TYPE_UNSPECIFIED",
												"CONSUMER_PROJECT",
												"ENCRYPTION_KEYS_PROJECT",
												"KEYRING",
												"CONSUMER_FOLDER",
											},
										},
									},
								},
								Unreadable: true,
							},
							"resources": &dcl.Property{
								Type:        "array",
								GoName:      "Resources",
								ReadOnly:    true,
								Description: "Output only. The resources associated with this workload. These resources will be created when creating the workload. If any of the projects already exist, the workload creation will fail. Always read only.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "WorkloadResources",
									Properties: map[string]*dcl.Property{
										"resourceId": &dcl.Property{
											Type:        "integer",
											Format:      "int64",
											GoName:      "ResourceId",
											Description: "Resource identifier. For a project this represents project_number.",
											Immutable:   true,
										},
										"resourceType": &dcl.Property{
											Type:        "string",
											GoName:      "ResourceType",
											GoType:      "WorkloadResourcesResourceTypeEnum",
											Description: "Indicates the type of resource. Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER",
											Immutable:   true,
											Enum: []string{
												"RESOURCE_TYPE_UNSPECIFIED",
												"CONSUMER_PROJECT",
												"ENCRYPTION_KEYS_PROJECT",
												"KEYRING",
												"CONSUMER_FOLDER",
											},
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
