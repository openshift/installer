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
package osconfig

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLOSPolicyAssignmentSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "OSConfig/OSPolicyAssignment",
			Description: "Represents an OSPolicyAssignment resource.",
			StructName:  "OSPolicyAssignment",
			Reference: &dcl.Link{
				Text: "API documentation",
				URL:  "https://cloud.google.com/compute/docs/osconfig/rest/v1/projects.locations.osPolicyAssignments",
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a OSPolicyAssignment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "oSPolicyAssignment",
						Required:    true,
						Description: "A full instance of a OSPolicyAssignment",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a OSPolicyAssignment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "oSPolicyAssignment",
						Required:    true,
						Description: "A full instance of a OSPolicyAssignment",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a OSPolicyAssignment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "oSPolicyAssignment",
						Required:    true,
						Description: "A full instance of a OSPolicyAssignment",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all OSPolicyAssignment",
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
				Description: "The function used to list information about many OSPolicyAssignment",
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
				"OSPolicyAssignment": &dcl.Component{
					Title:           "OSPolicyAssignment",
					ID:              "projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}",
					UsesStateHint:   true,
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"osPolicies",
							"instanceFilter",
							"rollout",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"baseline": &dcl.Property{
								Type:        "boolean",
								GoName:      "Baseline",
								ReadOnly:    true,
								Description: "Output only. Indicates that this revision has been successfully rolled out in this zone and new VMs will be assigned OS policies from this revision. For a given OS policy assignment, there is only one revision with a value of `true` for this field.",
								Immutable:   true,
							},
							"deleted": &dcl.Property{
								Type:        "boolean",
								GoName:      "Deleted",
								ReadOnly:    true,
								Description: "Output only. Indicates that this revision deletes the OS policy assignment.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "OS policy assignment description. Length of the description is limited to 1024 characters.",
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "The etag for this OS policy assignment. If this is provided on update, it must match the server's etag.",
								Immutable:   true,
							},
							"instanceFilter": &dcl.Property{
								Type:        "object",
								GoName:      "InstanceFilter",
								GoType:      "OSPolicyAssignmentInstanceFilter",
								Description: "Required. Filter to select VMs.",
								Properties: map[string]*dcl.Property{
									"all": &dcl.Property{
										Type:        "boolean",
										GoName:      "All",
										Description: "Target all VMs in the project. If true, no other criteria is permitted.",
										SendEmpty:   true,
									},
									"exclusionLabels": &dcl.Property{
										Type:        "array",
										GoName:      "ExclusionLabels",
										Description: "List of label sets used for VM exclusion. If the list has more than one label set, the VM is excluded if any of the label sets are applicable for the VM.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "OSPolicyAssignmentInstanceFilterExclusionLabels",
											Properties: map[string]*dcl.Property{
												"labels": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Labels",
													Description: "Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.",
												},
											},
										},
									},
									"inclusionLabels": &dcl.Property{
										Type:        "array",
										GoName:      "InclusionLabels",
										Description: "List of label sets used for VM inclusion. If the list has more than one `LabelSet`, the VM is included if any of the label sets are applicable for the VM.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "OSPolicyAssignmentInstanceFilterInclusionLabels",
											Properties: map[string]*dcl.Property{
												"labels": &dcl.Property{
													Type: "object",
													AdditionalProperties: &dcl.Property{
														Type: "string",
													},
													GoName:      "Labels",
													Description: "Labels are identified by key/value pairs in this map. A VM should contain all the key/value pairs specified in this map to be selected.",
												},
											},
										},
									},
									"inventories": &dcl.Property{
										Type:        "array",
										GoName:      "Inventories",
										Description: "List of inventories to select VMs. A VM is selected if its inventory data matches at least one of the following inventories.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "OSPolicyAssignmentInstanceFilterInventories",
											Required: []string{
												"osShortName",
											},
											Properties: map[string]*dcl.Property{
												"osShortName": &dcl.Property{
													Type:        "string",
													GoName:      "OSShortName",
													Description: "Required. The OS short name",
												},
												"osVersion": &dcl.Property{
													Type:        "string",
													GoName:      "OSVersion",
													Description: "The OS version Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of `7`, specify the following value for this field `7.*` An empty string matches all OS versions.",
												},
											},
										},
									},
								},
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
								Description: "Resource name.",
								Immutable:   true,
							},
							"osPolicies": &dcl.Property{
								Type:        "array",
								GoName:      "OSPolicies",
								Description: "Required. List of OS policies to be applied to the VMs.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "OSPolicyAssignmentOSPolicies",
									Required: []string{
										"id",
										"mode",
										"resourceGroups",
									},
									Properties: map[string]*dcl.Property{
										"allowNoResourceGroupMatch": &dcl.Property{
											Type:        "boolean",
											GoName:      "AllowNoResourceGroupMatch",
											Description: "This flag determines the OS policy compliance status when none of the resource groups within the policy are applicable for a VM. Set this value to `true` if the policy needs to be reported as compliant even if the policy has nothing to validate or enforce.",
										},
										"description": &dcl.Property{
											Type:        "string",
											GoName:      "Description",
											Description: "Policy description. Length of the description is limited to 1024 characters.",
										},
										"id": &dcl.Property{
											Type:        "string",
											GoName:      "Id",
											Description: "Required. The id of the OS policy with the following restrictions: * Must contain only lowercase letters, numbers, and hyphens. * Must start with a letter. * Must be between 1-63 characters. * Must end with a number or a letter. * Must be unique within the assignment.",
										},
										"mode": &dcl.Property{
											Type:        "string",
											GoName:      "Mode",
											GoType:      "OSPolicyAssignmentOSPoliciesModeEnum",
											Description: "Required. Policy mode Possible values: MODE_UNSPECIFIED, VALIDATION, ENFORCEMENT",
											Enum: []string{
												"MODE_UNSPECIFIED",
												"VALIDATION",
												"ENFORCEMENT",
											},
										},
										"resourceGroups": &dcl.Property{
											Type:        "array",
											GoName:      "ResourceGroups",
											Description: "Required. List of resource groups for the policy. For a particular VM, resource groups are evaluated in the order specified and the first resource group that is applicable is selected and the rest are ignored. If none of the resource groups are applicable for a VM, the VM is considered to be non-compliant w.r.t this policy. This behavior can be toggled by the flag `allow_no_resource_group_match`",
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "object",
												GoType: "OSPolicyAssignmentOSPoliciesResourceGroups",
												Required: []string{
													"resources",
												},
												Properties: map[string]*dcl.Property{
													"inventoryFilters": &dcl.Property{
														Type:        "array",
														GoName:      "InventoryFilters",
														Description: "List of inventory filters for the resource group. The resources in this resource group are applied to the target VM if it satisfies at least one of the following inventory filters. For example, to apply this resource group to VMs running either `RHEL` or `CentOS` operating systems, specify 2 items for the list with following values: inventory_filters[0].os_short_name='rhel' and inventory_filters[1].os_short_name='centos' If the list is empty, this resource group will be applied to the target VM unconditionally.",
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "object",
															GoType: "OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters",
															Required: []string{
																"osShortName",
															},
															Properties: map[string]*dcl.Property{
																"osShortName": &dcl.Property{
																	Type:        "string",
																	GoName:      "OSShortName",
																	Description: "Required. The OS short name",
																},
																"osVersion": &dcl.Property{
																	Type:        "string",
																	GoName:      "OSVersion",
																	Description: "The OS version Prefix matches are supported if asterisk(*) is provided as the last character. For example, to match all versions with a major version of `7`, specify the following value for this field `7.*` An empty string matches all OS versions.",
																},
															},
														},
													},
													"resources": &dcl.Property{
														Type:        "array",
														GoName:      "Resources",
														Description: "Required. List of resources configured for this resource group. The resources are executed in the exact order specified here.",
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "object",
															GoType: "OSPolicyAssignmentOSPoliciesResourceGroupsResources",
															Required: []string{
																"id",
															},
															Properties: map[string]*dcl.Property{
																"exec": &dcl.Property{
																	Type:        "object",
																	GoName:      "Exec",
																	GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec",
																	Description: "Exec resource",
																	Conflicts: []string{
																		"pkg",
																		"repository",
																		"file",
																	},
																	Required: []string{
																		"validate",
																	},
																	Properties: map[string]*dcl.Property{
																		"enforce": &dcl.Property{
																			Type:        "object",
																			GoName:      "Enforce",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce",
																			Description: "What to run to bring this resource into the desired state. An exit code of 100 indicates \"success\", any other exit code indicates a failure running enforce.",
																			Required: []string{
																				"interpreter",
																			},
																			Properties: map[string]*dcl.Property{
																				"args": &dcl.Property{
																					Type:        "array",
																					GoName:      "Args",
																					Description: "Optional arguments to pass to the source during execution.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"file": &dcl.Property{
																					Type:        "object",
																					GoName:      "File",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile",
																					Description: "A remote or local file.",
																					Conflicts: []string{
																						"script",
																					},
																					Properties: map[string]*dcl.Property{
																						"allowInsecure": &dcl.Property{
																							Type:        "boolean",
																							GoName:      "AllowInsecure",
																							Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
																						},
																						"gcs": &dcl.Property{
																							Type:        "object",
																							GoName:      "Gcs",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs",
																							Description: "A Cloud Storage object.",
																							Conflicts: []string{
																								"remote",
																								"localPath",
																							},
																							Required: []string{
																								"bucket",
																								"object",
																							},
																							Properties: map[string]*dcl.Property{
																								"bucket": &dcl.Property{
																									Type:        "string",
																									GoName:      "Bucket",
																									Description: "Required. Bucket of the Cloud Storage object.",
																								},
																								"generation": &dcl.Property{
																									Type:        "integer",
																									Format:      "int64",
																									GoName:      "Generation",
																									Description: "Generation number of the Cloud Storage object.",
																								},
																								"object": &dcl.Property{
																									Type:        "string",
																									GoName:      "Object",
																									Description: "Required. Name of the Cloud Storage object.",
																								},
																							},
																						},
																						"localPath": &dcl.Property{
																							Type:        "string",
																							GoName:      "LocalPath",
																							Description: "A local path within the VM to use.",
																							Conflicts: []string{
																								"remote",
																								"gcs",
																							},
																						},
																						"remote": &dcl.Property{
																							Type:        "object",
																							GoName:      "Remote",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote",
																							Description: "A generic remote file.",
																							Conflicts: []string{
																								"gcs",
																								"localPath",
																							},
																							Required: []string{
																								"uri",
																							},
																							Properties: map[string]*dcl.Property{
																								"sha256Checksum": &dcl.Property{
																									Type:        "string",
																									GoName:      "Sha256Checksum",
																									Description: "SHA256 checksum of the remote file.",
																								},
																								"uri": &dcl.Property{
																									Type:        "string",
																									GoName:      "Uri",
																									Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
																								},
																							},
																						},
																					},
																				},
																				"interpreter": &dcl.Property{
																					Type:        "string",
																					GoName:      "Interpreter",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum",
																					Description: "Required. The script interpreter to use. Possible values: INTERPRETER_UNSPECIFIED, NONE, SHELL, POWERSHELL",
																					Enum: []string{
																						"INTERPRETER_UNSPECIFIED",
																						"NONE",
																						"SHELL",
																						"POWERSHELL",
																					},
																				},
																				"outputFilePath": &dcl.Property{
																					Type:        "string",
																					GoName:      "OutputFilePath",
																					Description: "Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.",
																				},
																				"script": &dcl.Property{
																					Type:        "string",
																					GoName:      "Script",
																					Description: "An inline script. The size of the script is limited to 1024 characters.",
																					Conflicts: []string{
																						"file",
																					},
																				},
																			},
																		},
																		"validate": &dcl.Property{
																			Type:        "object",
																			GoName:      "Validate",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate",
																			Description: "Required. What to run to validate this resource is in the desired state. An exit code of 100 indicates \"in desired state\", and exit code of 101 indicates \"not in desired state\". Any other exit code indicates a failure running validate.",
																			Required: []string{
																				"interpreter",
																			},
																			Properties: map[string]*dcl.Property{
																				"args": &dcl.Property{
																					Type:        "array",
																					GoName:      "Args",
																					Description: "Optional arguments to pass to the source during execution.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"file": &dcl.Property{
																					Type:        "object",
																					GoName:      "File",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile",
																					Description: "A remote or local file.",
																					Conflicts: []string{
																						"script",
																					},
																					Properties: map[string]*dcl.Property{
																						"allowInsecure": &dcl.Property{
																							Type:        "boolean",
																							GoName:      "AllowInsecure",
																							Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
																						},
																						"gcs": &dcl.Property{
																							Type:        "object",
																							GoName:      "Gcs",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs",
																							Description: "A Cloud Storage object.",
																							Conflicts: []string{
																								"remote",
																								"localPath",
																							},
																							Required: []string{
																								"bucket",
																								"object",
																							},
																							Properties: map[string]*dcl.Property{
																								"bucket": &dcl.Property{
																									Type:        "string",
																									GoName:      "Bucket",
																									Description: "Required. Bucket of the Cloud Storage object.",
																								},
																								"generation": &dcl.Property{
																									Type:        "integer",
																									Format:      "int64",
																									GoName:      "Generation",
																									Description: "Generation number of the Cloud Storage object.",
																								},
																								"object": &dcl.Property{
																									Type:        "string",
																									GoName:      "Object",
																									Description: "Required. Name of the Cloud Storage object.",
																								},
																							},
																						},
																						"localPath": &dcl.Property{
																							Type:        "string",
																							GoName:      "LocalPath",
																							Description: "A local path within the VM to use.",
																							Conflicts: []string{
																								"remote",
																								"gcs",
																							},
																						},
																						"remote": &dcl.Property{
																							Type:        "object",
																							GoName:      "Remote",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote",
																							Description: "A generic remote file.",
																							Conflicts: []string{
																								"gcs",
																								"localPath",
																							},
																							Required: []string{
																								"uri",
																							},
																							Properties: map[string]*dcl.Property{
																								"sha256Checksum": &dcl.Property{
																									Type:        "string",
																									GoName:      "Sha256Checksum",
																									Description: "SHA256 checksum of the remote file.",
																								},
																								"uri": &dcl.Property{
																									Type:        "string",
																									GoName:      "Uri",
																									Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
																								},
																							},
																						},
																					},
																				},
																				"interpreter": &dcl.Property{
																					Type:        "string",
																					GoName:      "Interpreter",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum",
																					Description: "Required. The script interpreter to use. Possible values: INTERPRETER_UNSPECIFIED, NONE, SHELL, POWERSHELL",
																					Enum: []string{
																						"INTERPRETER_UNSPECIFIED",
																						"NONE",
																						"SHELL",
																						"POWERSHELL",
																					},
																				},
																				"outputFilePath": &dcl.Property{
																					Type:        "string",
																					GoName:      "OutputFilePath",
																					Description: "Only recorded for enforce Exec. Path to an output file (that is created by this Exec) whose content will be recorded in OSPolicyResourceCompliance after a successful run. Absence or failure to read this file will result in this ExecResource being non-compliant. Output file size is limited to 100K bytes.",
																				},
																				"script": &dcl.Property{
																					Type:        "string",
																					GoName:      "Script",
																					Description: "An inline script. The size of the script is limited to 1024 characters.",
																					Conflicts: []string{
																						"file",
																					},
																				},
																			},
																		},
																	},
																},
																"file": &dcl.Property{
																	Type:        "object",
																	GoName:      "File",
																	GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile",
																	Description: "File resource",
																	Conflicts: []string{
																		"pkg",
																		"repository",
																		"exec",
																	},
																	Required: []string{
																		"path",
																		"state",
																	},
																	Properties: map[string]*dcl.Property{
																		"content": &dcl.Property{
																			Type:        "string",
																			GoName:      "Content",
																			Description: "A a file with this content. The size of the content is limited to 1024 characters.",
																			Conflicts: []string{
																				"file",
																			},
																		},
																		"file": &dcl.Property{
																			Type:        "object",
																			GoName:      "File",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile",
																			Description: "A remote or local source.",
																			Conflicts: []string{
																				"content",
																			},
																			Properties: map[string]*dcl.Property{
																				"allowInsecure": &dcl.Property{
																					Type:        "boolean",
																					GoName:      "AllowInsecure",
																					Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
																				},
																				"gcs": &dcl.Property{
																					Type:        "object",
																					GoName:      "Gcs",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs",
																					Description: "A Cloud Storage object.",
																					Conflicts: []string{
																						"remote",
																						"localPath",
																					},
																					Required: []string{
																						"bucket",
																						"object",
																					},
																					Properties: map[string]*dcl.Property{
																						"bucket": &dcl.Property{
																							Type:        "string",
																							GoName:      "Bucket",
																							Description: "Required. Bucket of the Cloud Storage object.",
																						},
																						"generation": &dcl.Property{
																							Type:        "integer",
																							Format:      "int64",
																							GoName:      "Generation",
																							Description: "Generation number of the Cloud Storage object.",
																						},
																						"object": &dcl.Property{
																							Type:        "string",
																							GoName:      "Object",
																							Description: "Required. Name of the Cloud Storage object.",
																						},
																					},
																				},
																				"localPath": &dcl.Property{
																					Type:        "string",
																					GoName:      "LocalPath",
																					Description: "A local path within the VM to use.",
																					Conflicts: []string{
																						"remote",
																						"gcs",
																					},
																				},
																				"remote": &dcl.Property{
																					Type:        "object",
																					GoName:      "Remote",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote",
																					Description: "A generic remote file.",
																					Conflicts: []string{
																						"gcs",
																						"localPath",
																					},
																					Required: []string{
																						"uri",
																					},
																					Properties: map[string]*dcl.Property{
																						"sha256Checksum": &dcl.Property{
																							Type:        "string",
																							GoName:      "Sha256Checksum",
																							Description: "SHA256 checksum of the remote file.",
																						},
																						"uri": &dcl.Property{
																							Type:        "string",
																							GoName:      "Uri",
																							Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
																						},
																					},
																				},
																			},
																		},
																		"path": &dcl.Property{
																			Type:        "string",
																			GoName:      "Path",
																			Description: "Required. The absolute path of the file within the VM.",
																		},
																		"permissions": &dcl.Property{
																			Type:        "string",
																			GoName:      "Permissions",
																			ReadOnly:    true,
																			Description: "Consists of three octal digits which represent, in order, the permissions of the owner, group, and other users for the file (similarly to the numeric mode used in the linux chmod utility). Each digit represents a three bit number with the 4 bit corresponding to the read permissions, the 2 bit corresponds to the write bit, and the one bit corresponds to the execute permission. Default behavior is 755. Below are some examples of permissions and their associated values: read, write, and execute: 7 read and execute: 5 read and write: 6 read only: 4",
																		},
																		"state": &dcl.Property{
																			Type:        "string",
																			GoName:      "State",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum",
																			Description: "Required. Desired state of the file. Possible values: OS_POLICY_COMPLIANCE_STATE_UNSPECIFIED, COMPLIANT, NON_COMPLIANT, UNKNOWN, NO_OS_POLICIES_APPLICABLE",
																			Enum: []string{
																				"OS_POLICY_COMPLIANCE_STATE_UNSPECIFIED",
																				"COMPLIANT",
																				"NON_COMPLIANT",
																				"UNKNOWN",
																				"NO_OS_POLICIES_APPLICABLE",
																			},
																		},
																	},
																},
																"id": &dcl.Property{
																	Type:        "string",
																	GoName:      "Id",
																	Description: "Required. The id of the resource with the following restrictions: * Must contain only lowercase letters, numbers, and hyphens. * Must start with a letter. * Must be between 1-63 characters. * Must end with a number or a letter. * Must be unique within the OS policy.",
																},
																"pkg": &dcl.Property{
																	Type:        "object",
																	GoName:      "Pkg",
																	GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg",
																	Description: "Package resource",
																	Conflicts: []string{
																		"repository",
																		"exec",
																		"file",
																	},
																	Required: []string{
																		"desiredState",
																	},
																	Properties: map[string]*dcl.Property{
																		"apt": &dcl.Property{
																			Type:        "object",
																			GoName:      "Apt",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt",
																			Description: "A package managed by Apt.",
																			Conflicts: []string{
																				"deb",
																				"yum",
																				"zypper",
																				"rpm",
																				"googet",
																				"msi",
																			},
																			Required: []string{
																				"name",
																			},
																			Properties: map[string]*dcl.Property{
																				"name": &dcl.Property{
																					Type:        "string",
																					GoName:      "Name",
																					Description: "Required. Package name.",
																				},
																			},
																		},
																		"deb": &dcl.Property{
																			Type:        "object",
																			GoName:      "Deb",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb",
																			Description: "A deb package file.",
																			Conflicts: []string{
																				"apt",
																				"yum",
																				"zypper",
																				"rpm",
																				"googet",
																				"msi",
																			},
																			Required: []string{
																				"source",
																			},
																			Properties: map[string]*dcl.Property{
																				"pullDeps": &dcl.Property{
																					Type:        "boolean",
																					GoName:      "PullDeps",
																					Description: "Whether dependencies should also be installed. - install when false: `dpkg -i package` - install when true: `apt-get update && apt-get -y install package.deb`",
																				},
																				"source": &dcl.Property{
																					Type:        "object",
																					GoName:      "Source",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource",
																					Description: "Required. A deb package.",
																					Properties: map[string]*dcl.Property{
																						"allowInsecure": &dcl.Property{
																							Type:        "boolean",
																							GoName:      "AllowInsecure",
																							Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
																						},
																						"gcs": &dcl.Property{
																							Type:        "object",
																							GoName:      "Gcs",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs",
																							Description: "A Cloud Storage object.",
																							Conflicts: []string{
																								"remote",
																								"localPath",
																							},
																							Required: []string{
																								"bucket",
																								"object",
																							},
																							Properties: map[string]*dcl.Property{
																								"bucket": &dcl.Property{
																									Type:        "string",
																									GoName:      "Bucket",
																									Description: "Required. Bucket of the Cloud Storage object.",
																								},
																								"generation": &dcl.Property{
																									Type:        "integer",
																									Format:      "int64",
																									GoName:      "Generation",
																									Description: "Generation number of the Cloud Storage object.",
																								},
																								"object": &dcl.Property{
																									Type:        "string",
																									GoName:      "Object",
																									Description: "Required. Name of the Cloud Storage object.",
																								},
																							},
																						},
																						"localPath": &dcl.Property{
																							Type:        "string",
																							GoName:      "LocalPath",
																							Description: "A local path within the VM to use.",
																							Conflicts: []string{
																								"remote",
																								"gcs",
																							},
																						},
																						"remote": &dcl.Property{
																							Type:        "object",
																							GoName:      "Remote",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote",
																							Description: "A generic remote file.",
																							Conflicts: []string{
																								"gcs",
																								"localPath",
																							},
																							Required: []string{
																								"uri",
																							},
																							Properties: map[string]*dcl.Property{
																								"sha256Checksum": &dcl.Property{
																									Type:        "string",
																									GoName:      "Sha256Checksum",
																									Description: "SHA256 checksum of the remote file.",
																								},
																								"uri": &dcl.Property{
																									Type:        "string",
																									GoName:      "Uri",
																									Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"desiredState": &dcl.Property{
																			Type:        "string",
																			GoName:      "DesiredState",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum",
																			Description: "Required. The desired state the agent should maintain for this package. Possible values: DESIRED_STATE_UNSPECIFIED, INSTALLED, REMOVED",
																			Enum: []string{
																				"DESIRED_STATE_UNSPECIFIED",
																				"INSTALLED",
																				"REMOVED",
																			},
																		},
																		"googet": &dcl.Property{
																			Type:        "object",
																			GoName:      "Googet",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget",
																			Description: "A package managed by GooGet.",
																			Conflicts: []string{
																				"apt",
																				"deb",
																				"yum",
																				"zypper",
																				"rpm",
																				"msi",
																			},
																			Required: []string{
																				"name",
																			},
																			Properties: map[string]*dcl.Property{
																				"name": &dcl.Property{
																					Type:        "string",
																					GoName:      "Name",
																					Description: "Required. Package name.",
																				},
																			},
																		},
																		"msi": &dcl.Property{
																			Type:        "object",
																			GoName:      "Msi",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi",
																			Description: "An MSI package.",
																			Conflicts: []string{
																				"apt",
																				"deb",
																				"yum",
																				"zypper",
																				"rpm",
																				"googet",
																			},
																			Required: []string{
																				"source",
																			},
																			Properties: map[string]*dcl.Property{
																				"properties": &dcl.Property{
																					Type:        "array",
																					GoName:      "Properties",
																					Description: "Additional properties to use during installation. This should be in the format of Property=Setting. Appended to the defaults of `ACTION=INSTALL REBOOT=ReallySuppress`.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"source": &dcl.Property{
																					Type:        "object",
																					GoName:      "Source",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource",
																					Description: "Required. The MSI package.",
																					Properties: map[string]*dcl.Property{
																						"allowInsecure": &dcl.Property{
																							Type:        "boolean",
																							GoName:      "AllowInsecure",
																							Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
																						},
																						"gcs": &dcl.Property{
																							Type:        "object",
																							GoName:      "Gcs",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs",
																							Description: "A Cloud Storage object.",
																							Conflicts: []string{
																								"remote",
																								"localPath",
																							},
																							Required: []string{
																								"bucket",
																								"object",
																							},
																							Properties: map[string]*dcl.Property{
																								"bucket": &dcl.Property{
																									Type:        "string",
																									GoName:      "Bucket",
																									Description: "Required. Bucket of the Cloud Storage object.",
																								},
																								"generation": &dcl.Property{
																									Type:        "integer",
																									Format:      "int64",
																									GoName:      "Generation",
																									Description: "Generation number of the Cloud Storage object.",
																								},
																								"object": &dcl.Property{
																									Type:        "string",
																									GoName:      "Object",
																									Description: "Required. Name of the Cloud Storage object.",
																								},
																							},
																						},
																						"localPath": &dcl.Property{
																							Type:        "string",
																							GoName:      "LocalPath",
																							Description: "A local path within the VM to use.",
																							Conflicts: []string{
																								"remote",
																								"gcs",
																							},
																						},
																						"remote": &dcl.Property{
																							Type:        "object",
																							GoName:      "Remote",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote",
																							Description: "A generic remote file.",
																							Conflicts: []string{
																								"gcs",
																								"localPath",
																							},
																							Required: []string{
																								"uri",
																							},
																							Properties: map[string]*dcl.Property{
																								"sha256Checksum": &dcl.Property{
																									Type:        "string",
																									GoName:      "Sha256Checksum",
																									Description: "SHA256 checksum of the remote file.",
																								},
																								"uri": &dcl.Property{
																									Type:        "string",
																									GoName:      "Uri",
																									Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"rpm": &dcl.Property{
																			Type:        "object",
																			GoName:      "Rpm",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm",
																			Description: "An rpm package file.",
																			Conflicts: []string{
																				"apt",
																				"deb",
																				"yum",
																				"zypper",
																				"googet",
																				"msi",
																			},
																			Required: []string{
																				"source",
																			},
																			Properties: map[string]*dcl.Property{
																				"pullDeps": &dcl.Property{
																					Type:        "boolean",
																					GoName:      "PullDeps",
																					Description: "Whether dependencies should also be installed. - install when false: `rpm --upgrade --replacepkgs package.rpm` - install when true: `yum -y install package.rpm` or `zypper -y install package.rpm`",
																				},
																				"source": &dcl.Property{
																					Type:        "object",
																					GoName:      "Source",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource",
																					Description: "Required. An rpm package.",
																					Properties: map[string]*dcl.Property{
																						"allowInsecure": &dcl.Property{
																							Type:        "boolean",
																							GoName:      "AllowInsecure",
																							Description: "Defaults to false. When false, files are subject to validations based on the file type: Remote: A checksum must be specified. Cloud Storage: An object generation number must be specified.",
																						},
																						"gcs": &dcl.Property{
																							Type:        "object",
																							GoName:      "Gcs",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs",
																							Description: "A Cloud Storage object.",
																							Conflicts: []string{
																								"remote",
																								"localPath",
																							},
																							Required: []string{
																								"bucket",
																								"object",
																							},
																							Properties: map[string]*dcl.Property{
																								"bucket": &dcl.Property{
																									Type:        "string",
																									GoName:      "Bucket",
																									Description: "Required. Bucket of the Cloud Storage object.",
																								},
																								"generation": &dcl.Property{
																									Type:        "integer",
																									Format:      "int64",
																									GoName:      "Generation",
																									Description: "Generation number of the Cloud Storage object.",
																								},
																								"object": &dcl.Property{
																									Type:        "string",
																									GoName:      "Object",
																									Description: "Required. Name of the Cloud Storage object.",
																								},
																							},
																						},
																						"localPath": &dcl.Property{
																							Type:        "string",
																							GoName:      "LocalPath",
																							Description: "A local path within the VM to use.",
																							Conflicts: []string{
																								"remote",
																								"gcs",
																							},
																						},
																						"remote": &dcl.Property{
																							Type:        "object",
																							GoName:      "Remote",
																							GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote",
																							Description: "A generic remote file.",
																							Conflicts: []string{
																								"gcs",
																								"localPath",
																							},
																							Required: []string{
																								"uri",
																							},
																							Properties: map[string]*dcl.Property{
																								"sha256Checksum": &dcl.Property{
																									Type:        "string",
																									GoName:      "Sha256Checksum",
																									Description: "SHA256 checksum of the remote file.",
																								},
																								"uri": &dcl.Property{
																									Type:        "string",
																									GoName:      "Uri",
																									Description: "Required. URI from which to fetch the object. It should contain both the protocol and path following the format `{protocol}://{location}`.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"yum": &dcl.Property{
																			Type:        "object",
																			GoName:      "Yum",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum",
																			Description: "A package managed by YUM.",
																			Conflicts: []string{
																				"apt",
																				"deb",
																				"zypper",
																				"rpm",
																				"googet",
																				"msi",
																			},
																			Required: []string{
																				"name",
																			},
																			Properties: map[string]*dcl.Property{
																				"name": &dcl.Property{
																					Type:        "string",
																					GoName:      "Name",
																					Description: "Required. Package name.",
																				},
																			},
																		},
																		"zypper": &dcl.Property{
																			Type:        "object",
																			GoName:      "Zypper",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper",
																			Description: "A package managed by Zypper.",
																			Conflicts: []string{
																				"apt",
																				"deb",
																				"yum",
																				"rpm",
																				"googet",
																				"msi",
																			},
																			Required: []string{
																				"name",
																			},
																			Properties: map[string]*dcl.Property{
																				"name": &dcl.Property{
																					Type:        "string",
																					GoName:      "Name",
																					Description: "Required. Package name.",
																				},
																			},
																		},
																	},
																},
																"repository": &dcl.Property{
																	Type:        "object",
																	GoName:      "Repository",
																	GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository",
																	Description: "Package repository resource",
																	Conflicts: []string{
																		"pkg",
																		"exec",
																		"file",
																	},
																	Properties: map[string]*dcl.Property{
																		"apt": &dcl.Property{
																			Type:        "object",
																			GoName:      "Apt",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt",
																			Description: "An Apt Repository.",
																			Conflicts: []string{
																				"yum",
																				"zypper",
																				"goo",
																			},
																			Required: []string{
																				"archiveType",
																				"uri",
																				"distribution",
																				"components",
																			},
																			Properties: map[string]*dcl.Property{
																				"archiveType": &dcl.Property{
																					Type:        "string",
																					GoName:      "ArchiveType",
																					GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum",
																					Description: "Required. Type of archive files in this repository. Possible values: ARCHIVE_TYPE_UNSPECIFIED, DEB, DEB_SRC",
																					Enum: []string{
																						"ARCHIVE_TYPE_UNSPECIFIED",
																						"DEB",
																						"DEB_SRC",
																					},
																				},
																				"components": &dcl.Property{
																					Type:        "array",
																					GoName:      "Components",
																					Description: "Required. List of components for this repository. Must contain at least one item.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"distribution": &dcl.Property{
																					Type:        "string",
																					GoName:      "Distribution",
																					Description: "Required. Distribution of this repository.",
																				},
																				"gpgKey": &dcl.Property{
																					Type:        "string",
																					GoName:      "GpgKey",
																					Description: "URI of the key file for this repository. The agent maintains a keyring at `/etc/apt/trusted.gpg.d/osconfig_agent_managed.gpg`.",
																				},
																				"uri": &dcl.Property{
																					Type:        "string",
																					GoName:      "Uri",
																					Description: "Required. URI for this repository.",
																				},
																			},
																		},
																		"goo": &dcl.Property{
																			Type:        "object",
																			GoName:      "Goo",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo",
																			Description: "A Goo Repository.",
																			Conflicts: []string{
																				"apt",
																				"yum",
																				"zypper",
																			},
																			Required: []string{
																				"name",
																				"url",
																			},
																			Properties: map[string]*dcl.Property{
																				"name": &dcl.Property{
																					Type:        "string",
																					GoName:      "Name",
																					Description: "Required. The name of the repository.",
																				},
																				"url": &dcl.Property{
																					Type:        "string",
																					GoName:      "Url",
																					Description: "Required. The url of the repository.",
																				},
																			},
																		},
																		"yum": &dcl.Property{
																			Type:        "object",
																			GoName:      "Yum",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum",
																			Description: "A Yum Repository.",
																			Conflicts: []string{
																				"apt",
																				"zypper",
																				"goo",
																			},
																			Required: []string{
																				"id",
																				"baseUrl",
																			},
																			Properties: map[string]*dcl.Property{
																				"baseUrl": &dcl.Property{
																					Type:        "string",
																					GoName:      "BaseUrl",
																					Description: "Required. The location of the repository directory.",
																				},
																				"displayName": &dcl.Property{
																					Type:        "string",
																					GoName:      "DisplayName",
																					Description: "The display name of the repository.",
																				},
																				"gpgKeys": &dcl.Property{
																					Type:        "array",
																					GoName:      "GpgKeys",
																					Description: "URIs of GPG keys.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"id": &dcl.Property{
																					Type:        "string",
																					GoName:      "Id",
																					Description: "Required. A one word, unique name for this repository. This is the `repo id` in the yum config file and also the `display_name` if `display_name` is omitted. This id is also used as the unique identifier when checking for resource conflicts.",
																				},
																			},
																		},
																		"zypper": &dcl.Property{
																			Type:        "object",
																			GoName:      "Zypper",
																			GoType:      "OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper",
																			Description: "A Zypper Repository.",
																			Conflicts: []string{
																				"apt",
																				"yum",
																				"goo",
																			},
																			Required: []string{
																				"id",
																				"baseUrl",
																			},
																			Properties: map[string]*dcl.Property{
																				"baseUrl": &dcl.Property{
																					Type:        "string",
																					GoName:      "BaseUrl",
																					Description: "Required. The location of the repository directory.",
																				},
																				"displayName": &dcl.Property{
																					Type:        "string",
																					GoName:      "DisplayName",
																					Description: "The display name of the repository.",
																				},
																				"gpgKeys": &dcl.Property{
																					Type:        "array",
																					GoName:      "GpgKeys",
																					Description: "URIs of GPG keys.",
																					SendEmpty:   true,
																					ListType:    "list",
																					Items: &dcl.Property{
																						Type:   "string",
																						GoType: "string",
																					},
																				},
																				"id": &dcl.Property{
																					Type:        "string",
																					GoName:      "Id",
																					Description: "Required. A one word, unique name for this repository. This is the `repo id` in the zypper config file and also the `display_name` if `display_name` is omitted. This id is also used as the unique identifier when checking for GuestPolicy conflicts.",
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
										},
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
							"reconciling": &dcl.Property{
								Type:        "boolean",
								GoName:      "Reconciling",
								ReadOnly:    true,
								Description: "Output only. Indicates that reconciliation is in progress for the revision. This value is `true` when the `rollout_state` is one of: * IN_PROGRESS * CANCELLING",
								Immutable:   true,
							},
							"revisionCreateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "RevisionCreateTime",
								ReadOnly:    true,
								Description: "Output only. The timestamp that the revision was created.",
								Immutable:   true,
							},
							"revisionId": &dcl.Property{
								Type:        "string",
								GoName:      "RevisionId",
								ReadOnly:    true,
								Description: "Output only. The assignment revision ID A new revision is committed whenever a rollout is triggered for a OS policy assignment",
								Immutable:   true,
							},
							"rollout": &dcl.Property{
								Type:        "object",
								GoName:      "Rollout",
								GoType:      "OSPolicyAssignmentRollout",
								Description: "Required. Rollout to deploy the OS policy assignment. A rollout is triggered in the following situations: 1) OSPolicyAssignment is created. 2) OSPolicyAssignment is updated and the update contains changes to one of the following fields: - instance_filter - os_policies 3) OSPolicyAssignment is deleted.",
								Required: []string{
									"disruptionBudget",
									"minWaitDuration",
								},
								Properties: map[string]*dcl.Property{
									"disruptionBudget": &dcl.Property{
										Type:        "object",
										GoName:      "DisruptionBudget",
										GoType:      "OSPolicyAssignmentRolloutDisruptionBudget",
										Description: "Required. The maximum number (or percentage) of VMs per zone to disrupt at any given moment.",
										Properties: map[string]*dcl.Property{
											"fixed": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Fixed",
												Description: "Specifies a fixed value.",
												Conflicts: []string{
													"percent",
												},
											},
											"percent": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Percent",
												Description: "Specifies the relative value defined as a percentage, which will be multiplied by a reference value.",
												Conflicts: []string{
													"fixed",
												},
											},
										},
									},
									"minWaitDuration": &dcl.Property{
										Type:        "string",
										GoName:      "MinWaitDuration",
										Description: "Required. This determines the minimum duration of time to wait after the configuration changes are applied through the current rollout. A VM continues to count towards the `disruption_budget` at least until this duration of time has passed after configuration changes are applied.",
									},
								},
							},
							"rolloutState": &dcl.Property{
								Type:        "string",
								GoName:      "RolloutState",
								GoType:      "OSPolicyAssignmentRolloutStateEnum",
								ReadOnly:    true,
								Description: "Output only. OS policy assignment rollout state Possible values: ROLLOUT_STATE_UNSPECIFIED, IN_PROGRESS, CANCELLING, CANCELLED, SUCCEEDED",
								Immutable:   true,
								Enum: []string{
									"ROLLOUT_STATE_UNSPECIFIED",
									"IN_PROGRESS",
									"CANCELLING",
									"CANCELLED",
									"SUCCEEDED",
								},
							},
							"skipAwaitRollout": &dcl.Property{
								Type:        "boolean",
								GoName:      "SkipAwaitRollout",
								Description: "Set to true to skip awaiting rollout during resource creation and update.",
								Unreadable:  true,
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. Server generated unique id for the OS policy assignment resource.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
