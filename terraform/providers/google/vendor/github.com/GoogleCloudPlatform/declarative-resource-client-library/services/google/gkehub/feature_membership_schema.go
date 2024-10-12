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
package gkehub

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLFeatureMembershipSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "GkeHub/FeatureMembership",
			Description: "The GkeHub FeatureMembership resource",
			StructName:  "FeatureMembership",
			Mutex:       "{{project}}/{{location}}/{{feature}}",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a FeatureMembership",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "featureMembership",
						Required:    true,
						Description: "A full instance of a FeatureMembership",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a FeatureMembership",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "featureMembership",
						Required:    true,
						Description: "A full instance of a FeatureMembership",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a FeatureMembership",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "featureMembership",
						Required:    true,
						Description: "A full instance of a FeatureMembership",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all FeatureMembership",
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
						Name:     "feature",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many FeatureMembership",
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
						Name:     "feature",
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
				"FeatureMembership": &dcl.Component{
					Title:           "FeatureMembership",
					ID:              "projects/{{project}}/locations/{{location}}/features/{{feature}}/memberships/{{membership}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"project",
							"location",
							"feature",
							"membership",
						},
						Properties: map[string]*dcl.Property{
							"configmanagement": &dcl.Property{
								Type:        "object",
								GoName:      "Configmanagement",
								GoType:      "FeatureMembershipConfigmanagement",
								Description: "Config Management-specific spec.",
								Properties: map[string]*dcl.Property{
									"binauthz": &dcl.Property{
										Type:          "object",
										GoName:        "Binauthz",
										GoType:        "FeatureMembershipConfigmanagementBinauthz",
										Description:   "**DEPRECATED** Binauthz configuration for the cluster. This field will be ignored and should not be set.",
										ServerDefault: true,
										Properties: map[string]*dcl.Property{
											"enabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "Enabled",
												Description: "Whether binauthz is enabled in this cluster.",
											},
										},
										Parameter: true,
									},
									"configSync": &dcl.Property{
										Type:        "object",
										GoName:      "ConfigSync",
										GoType:      "FeatureMembershipConfigmanagementConfigSync",
										Description: "Config Sync configuration for the cluster.",
										SendEmpty:   true,
										Properties: map[string]*dcl.Property{
											"git": &dcl.Property{
												Type:   "object",
												GoName: "Git",
												GoType: "FeatureMembershipConfigmanagementConfigSyncGit",
												Properties: map[string]*dcl.Property{
													"gcpServiceAccountEmail": &dcl.Property{
														Type:        "string",
														GoName:      "GcpServiceAccountEmail",
														Description: "The GCP Service Account Email used for auth when secretType is gcpServiceAccount.",
														ResourceReferences: []*dcl.PropertyResourceReference{
															&dcl.PropertyResourceReference{
																Resource: "Iam/ServiceAccount",
																Field:    "email",
															},
														},
													},
													"httpsProxy": &dcl.Property{
														Type:        "string",
														GoName:      "HttpsProxy",
														Description: "URL for the HTTPS proxy to be used when communicating with the Git repo.",
													},
													"policyDir": &dcl.Property{
														Type:        "string",
														GoName:      "PolicyDir",
														Description: "The path within the Git repository that represents the top level of the repo to sync. Default: the root directory of the repository.",
													},
													"secretType": &dcl.Property{
														Type:        "string",
														GoName:      "SecretType",
														Description: "Type of secret configured for access to the Git repo. Must be one of ssh, cookiefile, gcenode, token, gcpserviceaccount or none. The validation of this is case-sensitive.",
													},
													"syncBranch": &dcl.Property{
														Type:        "string",
														GoName:      "SyncBranch",
														Description: "The branch of the repository to sync from. Default: master.",
													},
													"syncRepo": &dcl.Property{
														Type:        "string",
														GoName:      "SyncRepo",
														Description: "The URL of the Git repository to use as the source of truth.",
													},
													"syncRev": &dcl.Property{
														Type:        "string",
														GoName:      "SyncRev",
														Description: "Git revision (tag or hash) to check out. Default HEAD.",
													},
													"syncWaitSecs": &dcl.Property{
														Type:        "string",
														GoName:      "SyncWaitSecs",
														Description: "Period in seconds between consecutive syncs. Default: 15.",
													},
												},
											},
											"metricsGcpServiceAccountEmail": &dcl.Property{
												Type:        "string",
												GoName:      "MetricsGcpServiceAccountEmail",
												Description: "The Email of the Google Cloud Service Account (GSA) used for exporting Config Sync metrics to Cloud Monitoring. The GSA should have the Monitoring Metric Writer(roles/monitoring.metricWriter) IAM role. The Kubernetes ServiceAccount `default` in the namespace `config-management-monitoring` should be bound to the GSA.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Iam/ServiceAccount",
														Field:    "email",
													},
												},
											},
											"oci": &dcl.Property{
												Type:   "object",
												GoName: "Oci",
												GoType: "FeatureMembershipConfigmanagementConfigSyncOci",
												Properties: map[string]*dcl.Property{
													"gcpServiceAccountEmail": &dcl.Property{
														Type:        "string",
														GoName:      "GcpServiceAccountEmail",
														Description: "The GCP Service Account Email used for auth when secret_type is gcpserviceaccount. ",
														ResourceReferences: []*dcl.PropertyResourceReference{
															&dcl.PropertyResourceReference{
																Resource: "Iam/ServiceAccount",
																Field:    "email",
															},
														},
													},
													"policyDir": &dcl.Property{
														Type:        "string",
														GoName:      "PolicyDir",
														Description: "The absolute path of the directory that contains the local resources. Default: the root directory of the image.",
													},
													"secretType": &dcl.Property{
														Type:        "string",
														GoName:      "SecretType",
														Description: "Type of secret configured for access to the OCI Image. Must be one of gcenode, gcpserviceaccount or none. The validation of this is case-sensitive.",
													},
													"syncRepo": &dcl.Property{
														Type:        "string",
														GoName:      "SyncRepo",
														Description: "The OCI image repository URL for the package to sync from. e.g. LOCATION-docker.pkg.dev/PROJECT_ID/REPOSITORY_NAME/PACKAGE_NAME.",
													},
													"syncWaitSecs": &dcl.Property{
														Type:        "string",
														GoName:      "SyncWaitSecs",
														Description: "Period in seconds(int64 format) between consecutive syncs. Default: 15.",
													},
												},
											},
											"preventDrift": &dcl.Property{
												Type:          "boolean",
												GoName:        "PreventDrift",
												Description:   "Set to true to enable the Config Sync admission webhook to prevent drifts. If set to `false`, disables the Config Sync admission webhook and does not prevent drifts.",
												ServerDefault: true,
											},
											"sourceFormat": &dcl.Property{
												Type:        "string",
												GoName:      "SourceFormat",
												Description: "Specifies whether the Config Sync Repo is in \"hierarchical\" or \"unstructured\" mode.",
											},
										},
									},
									"hierarchyController": &dcl.Property{
										Type:        "object",
										GoName:      "HierarchyController",
										GoType:      "FeatureMembershipConfigmanagementHierarchyController",
										Description: "Hierarchy Controller configuration for the cluster.",
										SendEmpty:   true,
										Properties: map[string]*dcl.Property{
											"enableHierarchicalResourceQuota": &dcl.Property{
												Type:        "boolean",
												GoName:      "EnableHierarchicalResourceQuota",
												Description: "Whether hierarchical resource quota is enabled in this cluster.",
												SendEmpty:   true,
											},
											"enablePodTreeLabels": &dcl.Property{
												Type:        "boolean",
												GoName:      "EnablePodTreeLabels",
												Description: "Whether pod tree labels are enabled in this cluster.",
												SendEmpty:   true,
											},
											"enabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "Enabled",
												Description: "Whether Hierarchy Controller is enabled in this cluster.",
												SendEmpty:   true,
											},
										},
									},
									"policyController": &dcl.Property{
										Type:        "object",
										GoName:      "PolicyController",
										GoType:      "FeatureMembershipConfigmanagementPolicyController",
										Description: "Policy Controller configuration for the cluster.",
										Properties: map[string]*dcl.Property{
											"auditIntervalSeconds": &dcl.Property{
												Type:        "string",
												GoName:      "AuditIntervalSeconds",
												Description: "Sets the interval for Policy Controller Audit Scans (in seconds). When set to 0, this disables audit functionality altogether.",
											},
											"enabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "Enabled",
												Description: "Enables the installation of Policy Controller. If false, the rest of PolicyController fields take no effect.",
											},
											"exemptableNamespaces": &dcl.Property{
												Type:        "array",
												GoName:      "ExemptableNamespaces",
												Description: "The set of namespaces that are excluded from Policy Controller checks. Namespaces do not need to currently exist on the cluster.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
											"logDeniesEnabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "LogDeniesEnabled",
												Description: "Logs all denies and dry run failures.",
											},
											"monitoring": &dcl.Property{
												Type:          "object",
												GoName:        "Monitoring",
												GoType:        "FeatureMembershipConfigmanagementPolicyControllerMonitoring",
												Description:   "Specifies the backends Policy Controller should export metrics to. For example, to specify metrics should be exported to Cloud Monitoring and Prometheus, specify backends: [\"cloudmonitoring\", \"prometheus\"]. Default: [\"cloudmonitoring\", \"prometheus\"]",
												ServerDefault: true,
												Properties: map[string]*dcl.Property{
													"backends": &dcl.Property{
														Type:          "array",
														GoName:        "Backends",
														Description:   " Specifies the list of backends Policy Controller will export to. Specifying an empty value `[]` disables metrics export.",
														ServerDefault: true,
														SendEmpty:     true,
														ListType:      "list",
														Items: &dcl.Property{
															Type:   "string",
															GoType: "FeatureMembershipConfigmanagementPolicyControllerMonitoringBackendsEnum",
															Enum: []string{
																"MONITORING_BACKEND_UNSPECIFIED",
																"PROMETHEUS",
																"CLOUD_MONITORING",
															},
														},
													},
												},
											},
											"mutationEnabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "MutationEnabled",
												Description: "Enable or disable mutation in policy controller. If true, mutation CRDs, webhook and controller deployment will be deployed to the cluster.",
											},
											"referentialRulesEnabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "ReferentialRulesEnabled",
												Description: "Enables the ability to use Constraint Templates that reference to objects other than the object currently being evaluated.",
											},
											"templateLibraryInstalled": &dcl.Property{
												Type:        "boolean",
												GoName:      "TemplateLibraryInstalled",
												Description: "Installs the default template library along with Policy Controller.",
											},
										},
									},
									"version": &dcl.Property{
										Type:          "string",
										GoName:        "Version",
										Description:   "Optional. Version of ACM to install. Defaults to the latest version.",
										ServerDefault: true,
									},
								},
							},
							"feature": &dcl.Property{
								Type:        "string",
								GoName:      "Feature",
								Description: "The name of the feature",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Gkehub/Feature",
										Field:    "name",
										Parent:   true,
									},
								},
								Parameter: true,
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location of the feature",
								Immutable:   true,
								Parameter:   true,
							},
							"membership": &dcl.Property{
								Type:        "string",
								GoName:      "Membership",
								Description: "The name of the membership",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Gkehub/Membership",
										Field:    "name",
									},
								},
								Parameter: true,
							},
							"membershipLocation": &dcl.Property{
								Type:        "string",
								GoName:      "MembershipLocation",
								Description: "The location of the membership",
								Immutable:   true,
								Parameter:   true,
							},
							"mesh": &dcl.Property{
								Type:        "object",
								GoName:      "Mesh",
								GoType:      "FeatureMembershipMesh",
								Description: "Manage Mesh Features",
								Properties: map[string]*dcl.Property{
									"controlPlane": &dcl.Property{
										Type:        "string",
										GoName:      "ControlPlane",
										GoType:      "FeatureMembershipMeshControlPlaneEnum",
										Description: "**DEPRECATED** Whether to automatically manage Service Mesh control planes. Possible values: CONTROL_PLANE_MANAGEMENT_UNSPECIFIED, AUTOMATIC, MANUAL",
										Enum: []string{
											"CONTROL_PLANE_MANAGEMENT_UNSPECIFIED",
											"AUTOMATIC",
											"MANUAL",
										},
									},
									"management": &dcl.Property{
										Type:        "string",
										GoName:      "Management",
										GoType:      "FeatureMembershipMeshManagementEnum",
										Description: "Whether to automatically manage Service Mesh. Possible values: MANAGEMENT_UNSPECIFIED, MANAGEMENT_AUTOMATIC, MANAGEMENT_MANUAL",
										Enum: []string{
											"MANAGEMENT_UNSPECIFIED",
											"MANAGEMENT_AUTOMATIC",
											"MANAGEMENT_MANUAL",
										},
									},
								},
							},
							"policycontroller": &dcl.Property{
								Type:        "object",
								GoName:      "Policycontroller",
								GoType:      "FeatureMembershipPolicycontroller",
								Description: "Policy Controller-specific spec.",
								Required: []string{
									"policyControllerHubConfig",
								},
								Properties: map[string]*dcl.Property{
									"policyControllerHubConfig": &dcl.Property{
										Type:        "object",
										GoName:      "PolicyControllerHubConfig",
										GoType:      "FeatureMembershipPolicycontrollerPolicyControllerHubConfig",
										Description: "Policy Controller configuration for the cluster.",
										Properties: map[string]*dcl.Property{
											"auditIntervalSeconds": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "AuditIntervalSeconds",
												Description: "Sets the interval for Policy Controller Audit Scans (in seconds). When set to 0, this disables audit functionality altogether.",
											},
											"constraintViolationLimit": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "ConstraintViolationLimit",
												Description: "The maximum number of audit violations to be stored in a constraint. If not set, the internal default of 20 will be used.",
											},
											"deploymentConfigs": &dcl.Property{
												Type: "object",
												AdditionalProperties: &dcl.Property{
													Type:   "object",
													GoType: "FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigs",
													Properties: map[string]*dcl.Property{
														"containerResources": &dcl.Property{
															Type:        "object",
															GoName:      "ContainerResources",
															GoType:      "FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResources",
															Description: "Container resource requirements.",
															Conflicts: []string{
																"replicaCount",
																"podAffinity",
																"podTolerations",
															},
															Properties: map[string]*dcl.Property{
																"limits": &dcl.Property{
																	Type:        "object",
																	GoName:      "Limits",
																	GoType:      "FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesLimits",
																	Description: "Limits describes the maximum amount of compute resources allowed for use by the running container.",
																	Properties: map[string]*dcl.Property{
																		"cpu": &dcl.Property{
																			Type:        "string",
																			GoName:      "Cpu",
																			Description: "CPU requirement expressed in Kubernetes resource units.",
																		},
																		"memory": &dcl.Property{
																			Type:        "string",
																			GoName:      "Memory",
																			Description: "Memory requirement expressed in Kubernetes resource units.",
																		},
																	},
																},
																"requests": &dcl.Property{
																	Type:        "object",
																	GoName:      "Requests",
																	GoType:      "FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsContainerResourcesRequests",
																	Description: "Requests describes the amount of compute resources reserved for the container by the kube-scheduler.",
																	Properties: map[string]*dcl.Property{
																		"cpu": &dcl.Property{
																			Type:        "string",
																			GoName:      "Cpu",
																			Description: "CPU requirement expressed in Kubernetes resource units.",
																		},
																		"memory": &dcl.Property{
																			Type:        "string",
																			GoName:      "Memory",
																			Description: "Memory requirement expressed in Kubernetes resource units.",
																		},
																	},
																},
															},
														},
														"podAffinity": &dcl.Property{
															Type:        "string",
															GoName:      "PodAffinity",
															GoType:      "FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodAffinityEnum",
															Description: "Pod affinity configuration. Possible values: AFFINITY_UNSPECIFIED, NO_AFFINITY, ANTI_AFFINITY",
															Conflicts: []string{
																"replicaCount",
																"containerResources",
																"podTolerations",
															},
															Enum: []string{
																"AFFINITY_UNSPECIFIED",
																"NO_AFFINITY",
																"ANTI_AFFINITY",
															},
														},
														"podTolerations": &dcl.Property{
															Type:        "array",
															GoName:      "PodTolerations",
															Description: "Pod tolerations of node taints.",
															Conflicts: []string{
																"replicaCount",
																"containerResources",
																"podAffinity",
															},
															SendEmpty: true,
															ListType:  "list",
															Items: &dcl.Property{
																Type:   "object",
																GoType: "FeatureMembershipPolicycontrollerPolicyControllerHubConfigDeploymentConfigsPodTolerations",
																Properties: map[string]*dcl.Property{
																	"effect": &dcl.Property{
																		Type:        "string",
																		GoName:      "Effect",
																		Description: "Matches a taint effect.",
																	},
																	"key": &dcl.Property{
																		Type:        "string",
																		GoName:      "Key",
																		Description: "Matches a taint key (not necessarily unique).",
																	},
																	"operator": &dcl.Property{
																		Type:        "string",
																		GoName:      "Operator",
																		Description: "Matches a taint operator.",
																	},
																	"value": &dcl.Property{
																		Type:        "string",
																		GoName:      "Value",
																		Description: "Matches a taint value.",
																	},
																},
															},
														},
														"replicaCount": &dcl.Property{
															Type:        "integer",
															Format:      "int64",
															GoName:      "ReplicaCount",
															Description: "Pod replica count.",
															Conflicts: []string{
																"containerResources",
																"podAffinity",
																"podTolerations",
															},
														},
													},
												},
												GoName:        "DeploymentConfigs",
												Description:   "Map of deployment configs to deployments (\"admission\", \"audit\", \"mutation\").",
												ServerDefault: true,
											},
											"exemptableNamespaces": &dcl.Property{
												Type:        "array",
												GoName:      "ExemptableNamespaces",
												Description: "The set of namespaces that are excluded from Policy Controller checks. Namespaces do not need to currently exist on the cluster.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
											"installSpec": &dcl.Property{
												Type:        "string",
												GoName:      "InstallSpec",
												GoType:      "FeatureMembershipPolicycontrollerPolicyControllerHubConfigInstallSpecEnum",
												Description: "Configures the mode of the Policy Controller installation. Possible values: INSTALL_SPEC_UNSPECIFIED, INSTALL_SPEC_NOT_INSTALLED, INSTALL_SPEC_ENABLED, INSTALL_SPEC_SUSPENDED, INSTALL_SPEC_DETACHED",
												Enum: []string{
													"INSTALL_SPEC_UNSPECIFIED",
													"INSTALL_SPEC_NOT_INSTALLED",
													"INSTALL_SPEC_ENABLED",
													"INSTALL_SPEC_SUSPENDED",
													"INSTALL_SPEC_DETACHED",
												},
											},
											"logDeniesEnabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "LogDeniesEnabled",
												Description: "Logs all denies and dry run failures.",
											},
											"monitoring": &dcl.Property{
												Type:          "object",
												GoName:        "Monitoring",
												GoType:        "FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoring",
												Description:   "Specifies the backends Policy Controller should export metrics to. For example, to specify metrics should be exported to Cloud Monitoring and Prometheus, specify backends: [\"cloudmonitoring\", \"prometheus\"]. Default: [\"cloudmonitoring\", \"prometheus\"]",
												ServerDefault: true,
												Properties: map[string]*dcl.Property{
													"backends": &dcl.Property{
														Type:          "array",
														GoName:        "Backends",
														Description:   " Specifies the list of backends Policy Controller will export to. Specifying an empty value `[]` disables metrics export.",
														ServerDefault: true,
														SendEmpty:     true,
														ListType:      "list",
														Items: &dcl.Property{
															Type:   "string",
															GoType: "FeatureMembershipPolicycontrollerPolicyControllerHubConfigMonitoringBackendsEnum",
															Enum: []string{
																"MONITORING_BACKEND_UNSPECIFIED",
																"PROMETHEUS",
																"CLOUD_MONITORING",
															},
														},
													},
												},
											},
											"mutationEnabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "MutationEnabled",
												Description: "Enables the ability to mutate resources using Policy Controller.",
											},
											"policyContent": &dcl.Property{
												Type:          "object",
												GoName:        "PolicyContent",
												GoType:        "FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContent",
												Description:   "Specifies the desired policy content on the cluster.",
												ServerDefault: true,
												Properties: map[string]*dcl.Property{
													"bundles": &dcl.Property{
														Type: "object",
														AdditionalProperties: &dcl.Property{
															Type:   "object",
															GoType: "FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentBundles",
															Properties: map[string]*dcl.Property{
																"exemptedNamespaces": &dcl.Property{
																	Type:        "array",
																	GoName:      "ExemptedNamespaces",
																	Description: "The set of namespaces to be exempted from the bundle.",
																	SendEmpty:   true,
																	ListType:    "list",
																	Items: &dcl.Property{
																		Type:   "string",
																		GoType: "string",
																	},
																},
															},
														},
														GoName:      "Bundles",
														Description: "map of bundle name to BundleInstallSpec. The bundle name maps to the `bundleName` key in the `policycontroller.gke.io/constraintData` annotation on a constraint.",
													},
													"templateLibrary": &dcl.Property{
														Type:          "object",
														GoName:        "TemplateLibrary",
														GoType:        "FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibrary",
														Description:   "Configures the installation of the Template Library.",
														ServerDefault: true,
														Properties: map[string]*dcl.Property{
															"installation": &dcl.Property{
																Type:        "string",
																GoName:      "Installation",
																GoType:      "FeatureMembershipPolicycontrollerPolicyControllerHubConfigPolicyContentTemplateLibraryInstallationEnum",
																Description: "Configures the manner in which the template library is installed on the cluster. Possible values: INSTALLATION_UNSPECIFIED, NOT_INSTALLED, ALL",
																Enum: []string{
																	"INSTALLATION_UNSPECIFIED",
																	"NOT_INSTALLED",
																	"ALL",
																},
															},
														},
													},
												},
											},
											"referentialRulesEnabled": &dcl.Property{
												Type:        "boolean",
												GoName:      "ReferentialRulesEnabled",
												Description: "Enables the ability to use Constraint Templates that reference to objects other than the object currently being evaluated.",
											},
										},
									},
									"version": &dcl.Property{
										Type:          "string",
										GoName:        "Version",
										Description:   "Optional. Version of Policy Controller to install. Defaults to the latest version.",
										ServerDefault: true,
									},
								},
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project of the feature",
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
						},
					},
				},
			},
		},
	}
}
