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
							"organization",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"billingAccount": &dcl.Property{
								Type:        "string",
								GoName:      "BillingAccount",
								Description: "Optional. Input only. The billing account used for the resources which are direct children of workload. This billing account is initially associated with the resources created as part of Workload creation. After the initial creation of these resources, the customer can change the assigned billing account. The resource name has the form `billingAccounts/{billing_account_id}`. For example, `billingAccounts/012345-567890-ABCDEF`.",
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
								Description: "Required. Immutable. Compliance Regime associated with this workload. Possible values: COMPLIANCE_REGIME_UNSPECIFIED, IL4, CJIS, FEDRAMP_HIGH, FEDRAMP_MODERATE, US_REGIONAL_ACCESS, HIPAA, HITRUST, EU_REGIONS_AND_SUPPORT, CA_REGIONS_AND_SUPPORT, ITAR, AU_REGIONS_AND_US_SUPPORT, ASSURED_WORKLOADS_FOR_PARTNERS, ISR_REGIONS, ISR_REGIONS_AND_SUPPORT, CA_PROTECTED_B, IL5, IL2, JP_REGIONS_AND_SUPPORT, KSA_REGIONS_AND_SUPPORT_WITH_SOVEREIGNTY_CONTROLS, REGIONAL_CONTROLS, HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS, HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS_WITH_US_SUPPORT",
								Immutable:   true,
								Enum: []string{
									"COMPLIANCE_REGIME_UNSPECIFIED",
									"IL4",
									"CJIS",
									"FEDRAMP_HIGH",
									"FEDRAMP_MODERATE",
									"US_REGIONAL_ACCESS",
									"HIPAA",
									"HITRUST",
									"EU_REGIONS_AND_SUPPORT",
									"CA_REGIONS_AND_SUPPORT",
									"ITAR",
									"AU_REGIONS_AND_US_SUPPORT",
									"ASSURED_WORKLOADS_FOR_PARTNERS",
									"ISR_REGIONS",
									"ISR_REGIONS_AND_SUPPORT",
									"CA_PROTECTED_B",
									"IL5",
									"IL2",
									"JP_REGIONS_AND_SUPPORT",
									"KSA_REGIONS_AND_SUPPORT_WITH_SOVEREIGNTY_CONTROLS",
									"REGIONAL_CONTROLS",
									"HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS",
									"HEALTHCARE_AND_LIFE_SCIENCES_CONTROLS_WITH_US_SUPPORT",
								},
							},
							"complianceStatus": &dcl.Property{
								Type:        "object",
								GoName:      "ComplianceStatus",
								GoType:      "WorkloadComplianceStatus",
								ReadOnly:    true,
								Description: "Output only. Count of active Violations in the Workload.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"acknowledgedViolationCount": &dcl.Property{
										Type:        "array",
										GoName:      "AcknowledgedViolationCount",
										Description: "Number of current orgPolicy violations which are acknowledged.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "integer",
											Format: "int64",
											GoType: "int64",
										},
									},
									"activeViolationCount": &dcl.Property{
										Type:        "array",
										GoName:      "ActiveViolationCount",
										Description: "Number of current orgPolicy violations which are not acknowledged.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "integer",
											Format: "int64",
											GoType: "int64",
										},
									},
								},
							},
							"compliantButDisallowedServices": &dcl.Property{
								Type:        "array",
								GoName:      "CompliantButDisallowedServices",
								ReadOnly:    true,
								Description: "Output only. Urls for services which are compliant for this Assured Workload, but which are currently disallowed by the ResourceUsageRestriction org policy. Invoke workloads.restrictAllowedResources endpoint to allow your project developers to use these services in their environment.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
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
							"ekmProvisioningResponse": &dcl.Property{
								Type:        "object",
								GoName:      "EkmProvisioningResponse",
								GoType:      "WorkloadEkmProvisioningResponse",
								ReadOnly:    true,
								Description: "Optional. Represents the Ekm Provisioning State of the given workload.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"ekmProvisioningErrorDomain": &dcl.Property{
										Type:        "string",
										GoName:      "EkmProvisioningErrorDomain",
										GoType:      "WorkloadEkmProvisioningResponseEkmProvisioningErrorDomainEnum",
										Description: "Indicates Ekm provisioning error if any. Possible values: EKM_PROVISIONING_ERROR_DOMAIN_UNSPECIFIED, UNSPECIFIED_ERROR, GOOGLE_SERVER_ERROR, EXTERNAL_USER_ERROR, EXTERNAL_PARTNER_ERROR, TIMEOUT_ERROR",
										Immutable:   true,
										Enum: []string{
											"EKM_PROVISIONING_ERROR_DOMAIN_UNSPECIFIED",
											"UNSPECIFIED_ERROR",
											"GOOGLE_SERVER_ERROR",
											"EXTERNAL_USER_ERROR",
											"EXTERNAL_PARTNER_ERROR",
											"TIMEOUT_ERROR",
										},
									},
									"ekmProvisioningErrorMapping": &dcl.Property{
										Type:        "string",
										GoName:      "EkmProvisioningErrorMapping",
										GoType:      "WorkloadEkmProvisioningResponseEkmProvisioningErrorMappingEnum",
										Description: "Detailed error message if Ekm provisioning fails Possible values: EKM_PROVISIONING_ERROR_MAPPING_UNSPECIFIED, INVALID_SERVICE_ACCOUNT, MISSING_METRICS_SCOPE_ADMIN_PERMISSION, MISSING_EKM_CONNECTION_ADMIN_PERMISSION",
										Immutable:   true,
										Enum: []string{
											"EKM_PROVISIONING_ERROR_MAPPING_UNSPECIFIED",
											"INVALID_SERVICE_ACCOUNT",
											"MISSING_METRICS_SCOPE_ADMIN_PERMISSION",
											"MISSING_EKM_CONNECTION_ADMIN_PERMISSION",
										},
									},
									"ekmProvisioningState": &dcl.Property{
										Type:        "string",
										GoName:      "EkmProvisioningState",
										GoType:      "WorkloadEkmProvisioningResponseEkmProvisioningStateEnum",
										Description: "Indicates Ekm enrollment Provisioning of a given workload. Possible values: EKM_PROVISIONING_STATE_UNSPECIFIED, EKM_PROVISIONING_STATE_PENDING, EKM_PROVISIONING_STATE_FAILED, EKM_PROVISIONING_STATE_COMPLETED",
										Immutable:   true,
										Enum: []string{
											"EKM_PROVISIONING_STATE_UNSPECIFIED",
											"EKM_PROVISIONING_STATE_PENDING",
											"EKM_PROVISIONING_STATE_FAILED",
											"EKM_PROVISIONING_STATE_COMPLETED",
										},
									},
								},
							},
							"enableSovereignControls": &dcl.Property{
								Type:        "boolean",
								GoName:      "EnableSovereignControls",
								Description: "Optional. Indicates the sovereignty status of the given workload. Currently meant to be used by Europe/Canada customers.",
								Immutable:   true,
							},
							"kajEnrollmentState": &dcl.Property{
								Type:        "string",
								GoName:      "KajEnrollmentState",
								GoType:      "WorkloadKajEnrollmentStateEnum",
								ReadOnly:    true,
								Description: "Output only. Represents the KAJ enrollment state of the given workload. Possible values: KAJ_ENROLLMENT_STATE_UNSPECIFIED, KAJ_ENROLLMENT_STATE_PENDING, KAJ_ENROLLMENT_STATE_COMPLETE",
								Immutable:   true,
								Enum: []string{
									"KAJ_ENROLLMENT_STATE_UNSPECIFIED",
									"KAJ_ENROLLMENT_STATE_PENDING",
									"KAJ_ENROLLMENT_STATE_COMPLETE",
								},
							},
							"kmsSettings": &dcl.Property{
								Type:        "object",
								GoName:      "KmsSettings",
								GoType:      "WorkloadKmsSettings",
								Description: "**DEPRECATED** Input only. Settings used to create a CMEK crypto key. When set, a project with a KMS CMEK key is provisioned. This field is deprecated as of Feb 28, 2022. In order to create a Keyring, callers should specify, ENCRYPTION_KEYS_PROJECT or KEYRING in ResourceSettings.resource_type field.",
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
								Parameter:   true,
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "Output only. The resource name of the workload.",
								Immutable:                true,
								ServerGeneratedParameter: true,
								HasLongForm:              true,
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
								Parameter: true,
							},
							"partner": &dcl.Property{
								Type:        "string",
								GoName:      "Partner",
								GoType:      "WorkloadPartnerEnum",
								Description: "Optional. Partner regime associated with this workload. Possible values: PARTNER_UNSPECIFIED, LOCAL_CONTROLS_BY_S3NS, SOVEREIGN_CONTROLS_BY_T_SYSTEMS, SOVEREIGN_CONTROLS_BY_SIA_MINSAIT, SOVEREIGN_CONTROLS_BY_PSN, SOVEREIGN_CONTROLS_BY_CNTXT, SOVEREIGN_CONTROLS_BY_CNTXT_NO_EKM",
								Immutable:   true,
								Enum: []string{
									"PARTNER_UNSPECIFIED",
									"LOCAL_CONTROLS_BY_S3NS",
									"SOVEREIGN_CONTROLS_BY_T_SYSTEMS",
									"SOVEREIGN_CONTROLS_BY_SIA_MINSAIT",
									"SOVEREIGN_CONTROLS_BY_PSN",
									"SOVEREIGN_CONTROLS_BY_CNTXT",
									"SOVEREIGN_CONTROLS_BY_CNTXT_NO_EKM",
								},
							},
							"partnerPermissions": &dcl.Property{
								Type:        "object",
								GoName:      "PartnerPermissions",
								GoType:      "WorkloadPartnerPermissions",
								Description: "Optional. Permissions granted to the AW Partner SA account for the customer workload",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"assuredWorkloadsMonitoring": &dcl.Property{
										Type:        "boolean",
										GoName:      "AssuredWorkloadsMonitoring",
										Description: "Optional. Allow partner to view violation alerts.",
										Immutable:   true,
									},
									"dataLogsViewer": &dcl.Property{
										Type:        "boolean",
										GoName:      "DataLogsViewer",
										Description: "Allow the partner to view inspectability logs and monitoring violations.",
										Immutable:   true,
									},
									"serviceAccessApprover": &dcl.Property{
										Type:        "boolean",
										GoName:      "ServiceAccessApprover",
										Description: "Optional. Allow partner to view access approval logs.",
										Immutable:   true,
									},
								},
							},
							"partnerServicesBillingAccount": &dcl.Property{
								Type:        "string",
								GoName:      "PartnerServicesBillingAccount",
								Description: "Optional. Input only. Billing account necessary for purchasing services from Sovereign Partners. This field is required for creating SIA/PSN/CNTXT partner workloads. The caller should have 'billing.resourceAssociations.create' IAM permission on this billing-account. The format of this string is billingAccounts/AAAAAA-BBBBBB-CCCCCC.",
								Immutable:   true,
								Unreadable:  true,
							},
							"provisionedResourcesParent": &dcl.Property{
								Type:        "string",
								GoName:      "ProvisionedResourcesParent",
								Description: "Input only. The parent resource for the resources managed by this Assured Workload. May be either empty or a folder resource which is a child of the Workload parent. If not specified all resources are created under the parent organization. Format: folders/{folder_id}",
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
										"displayName": &dcl.Property{
											Type:        "string",
											GoName:      "DisplayName",
											Description: "User-assigned resource display name. If not empty it will be used to create a resource with the specified name.",
											Immutable:   true,
										},
										"resourceId": &dcl.Property{
											Type:        "string",
											GoName:      "ResourceId",
											Description: "Resource identifier. For a project this represents projectId. If the project is already taken, the workload creation will fail. For KeyRing, this represents the keyring_id. For a folder, don't set this value as folder_id is assigned by Google.",
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
							"saaEnrollmentResponse": &dcl.Property{
								Type:        "object",
								GoName:      "SaaEnrollmentResponse",
								GoType:      "WorkloadSaaEnrollmentResponse",
								ReadOnly:    true,
								Description: "Output only. Represents the SAA enrollment response of the given workload. SAA enrollment response is queried during workloads.get call. In failure cases, user friendly error message is shown in SAA details page.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"setupErrors": &dcl.Property{
										Type:        "array",
										GoName:      "SetupErrors",
										Description: "Indicates SAA enrollment setup error if any.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "WorkloadSaaEnrollmentResponseSetupErrorsEnum",
											Enum: []string{
												"SETUP_ERROR_UNSPECIFIED",
												"ERROR_INVALID_BASE_SETUP",
												"ERROR_MISSING_EXTERNAL_SIGNING_KEY",
												"ERROR_NOT_ALL_SERVICES_ENROLLED",
												"ERROR_SETUP_CHECK_FAILED",
											},
										},
									},
									"setupStatus": &dcl.Property{
										Type:        "string",
										GoName:      "SetupStatus",
										GoType:      "WorkloadSaaEnrollmentResponseSetupStatusEnum",
										Description: "Indicates SAA enrollment status of a given workload. Possible values: SETUP_STATE_UNSPECIFIED, STATUS_PENDING, STATUS_COMPLETE",
										Immutable:   true,
										Enum: []string{
											"SETUP_STATE_UNSPECIFIED",
											"STATUS_PENDING",
											"STATUS_COMPLETE",
										},
									},
								},
							},
							"violationNotificationsEnabled": &dcl.Property{
								Type:        "boolean",
								GoName:      "ViolationNotificationsEnabled",
								Description: "Optional. Indicates whether the e-mail notification for a violation is enabled for a workload. This value will be by default True, and if not present will be considered as true. This should only be updated via updateWorkload call. Any Changes to this field during the createWorkload call will not be honored. This will always be true while creating the workload.",
								Immutable:   true,
							},
							"workloadOptions": &dcl.Property{
								Type:        "object",
								GoName:      "WorkloadOptions",
								GoType:      "WorkloadWorkloadOptions",
								Description: "Optional. Used to specify certain options for a workload during workload creation - currently only supporting KAT Optionality for Regional Controls workloads.",
								Immutable:   true,
								Unreadable:  true,
								Properties: map[string]*dcl.Property{
									"kajEnrollmentType": &dcl.Property{
										Type:        "string",
										GoName:      "KajEnrollmentType",
										GoType:      "WorkloadWorkloadOptionsKajEnrollmentTypeEnum",
										Description: "Indicates type of KAJ enrollment for the workload. Currently, only specifiying KEY_ACCESS_TRANSPARENCY_OFF is implemented to not enroll in KAT-level KAJ enrollment for Regional Controls workloads. Possible values: KAJ_ENROLLMENT_TYPE_UNSPECIFIED, FULL_KAJ, EKM_ONLY, KEY_ACCESS_TRANSPARENCY_OFF",
										Immutable:   true,
										Enum: []string{
											"KAJ_ENROLLMENT_TYPE_UNSPECIFIED",
											"FULL_KAJ",
											"EKM_ONLY",
											"KEY_ACCESS_TRANSPARENCY_OFF",
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
