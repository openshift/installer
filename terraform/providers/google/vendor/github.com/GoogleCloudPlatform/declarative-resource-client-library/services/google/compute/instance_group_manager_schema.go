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
package compute

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLInstanceGroupManagerSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Compute/InstanceGroupManager",
			Description: "The Compute InstanceGroupManager resource",
			StructName:  "InstanceGroupManager",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a InstanceGroupManager",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "instanceGroupManager",
						Required:    true,
						Description: "A full instance of a InstanceGroupManager",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a InstanceGroupManager",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "instanceGroupManager",
						Required:    true,
						Description: "A full instance of a InstanceGroupManager",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a InstanceGroupManager",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "instanceGroupManager",
						Required:    true,
						Description: "A full instance of a InstanceGroupManager",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all InstanceGroupManager",
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
				Description: "The function used to list information about many InstanceGroupManager",
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
				"InstanceGroupManager": &dcl.Component{
					Title: "InstanceGroupManager",
					Locations: []string{
						"zone",
						"region",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"targetSize",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"autoHealingPolicies": &dcl.Property{
								Type:        "array",
								GoName:      "AutoHealingPolicies",
								Description: "The autohealing policy for this managed instance group. You can specify only one value.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "InstanceGroupManagerAutoHealingPolicies",
									Properties: map[string]*dcl.Property{
										"healthCheck": &dcl.Property{
											Type:        "string",
											GoName:      "HealthCheck",
											Description: "The URL for the health check that signals autohealing.",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/HealthCheck",
													Field:    "selfLink",
												},
											},
										},
										"initialDelaySec": &dcl.Property{
											Type:        "integer",
											Format:      "int64",
											GoName:      "InitialDelaySec",
											Description: "The number of seconds that the managed instance group waits before it applies autohealing policies to new instances or recently recreated instances. This initial delay allows instances to initialize and run their startup scripts before the instance group determines that they are UNHEALTHY. This prevents the managed instance group from recreating its instances prematurely. This value must be from range [0, 3600].",
										},
									},
								},
							},
							"baseInstanceName": &dcl.Property{
								Type:          "string",
								GoName:        "BaseInstanceName",
								Description:   "The base instance name to use for instances in this group. The value must be 1-58 characters long. Instances are named by appending a hyphen and a random four-character string to the base instance name. The base instance name must comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt).",
								ServerDefault: true,
							},
							"creationTimestamp": &dcl.Property{
								Type:        "string",
								GoName:      "CreationTimestamp",
								ReadOnly:    true,
								Description: "The creation timestamp for this managed instance group in \\[RFC3339\\](https://www.ietf.org/rfc/rfc3339.txt) text format.",
								Immutable:   true,
							},
							"currentActions": &dcl.Property{
								Type:        "object",
								GoName:      "CurrentActions",
								GoType:      "InstanceGroupManagerCurrentActions",
								ReadOnly:    true,
								Description: "[Output Only] The list of instance actions and the number of instances in this managed instance group that are scheduled for each of those actions.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"abandoning": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Abandoning",
										ReadOnly:    true,
										Description: "[Output Only] The total number of instances in the managed instance group that are scheduled to be abandoned. Abandoning an instance removes it from the managed instance group without deleting it.",
										Immutable:   true,
									},
									"creating": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Creating",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances in the managed instance group that are scheduled to be created or are currently being created. If the group fails to create any of these instances, it tries again until it creates the instance successfully. If you have disabled creation retries, this field will not be populated; instead, the `creatingWithoutRetries` field will be populated.",
										Immutable:   true,
									},
									"creatingWithoutRetries": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "CreatingWithoutRetries",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances that the managed instance group will attempt to create. The group attempts to create each instance only once. If the group fails to create any of these instances, it decreases the group's `targetSize` value accordingly.",
										Immutable:   true,
									},
									"deleting": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Deleting",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances in the managed instance group that are scheduled to be deleted or are currently being deleted.",
										Immutable:   true,
									},
									"none": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "None",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances in the managed instance group that are running and have no scheduled actions.",
										Immutable:   true,
									},
									"recreating": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Recreating",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances in the managed instance group that are scheduled to be recreated or are currently being being recreated. Recreating an instance deletes the existing root persistent disk and creates a new disk from the image that is defined in the instance template.",
										Immutable:   true,
									},
									"refreshing": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Refreshing",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances in the managed instance group that are being reconfigured with properties that do not require a restart or a recreate action. For example, setting or removing target pools for the instance.",
										Immutable:   true,
									},
									"restarting": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Restarting",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances in the managed instance group that are scheduled to be restarted or are currently being restarted.",
										Immutable:   true,
									},
									"verifying": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Verifying",
										ReadOnly:    true,
										Description: "[Output Only] The number of instances in the managed instance group that are being verified. See the `managedInstances[].currentAction` property in the `listManagedInstances` method documentation.",
										Immutable:   true,
									},
								},
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional description of this resource.",
								Immutable:   true,
							},
							"distributionPolicy": &dcl.Property{
								Type:          "object",
								GoName:        "DistributionPolicy",
								GoType:        "InstanceGroupManagerDistributionPolicy",
								Description:   "Policy specifying the intended distribution of managed instances across zones in a regional managed instance group.",
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"targetShape": &dcl.Property{
										Type:        "string",
										GoName:      "TargetShape",
										GoType:      "InstanceGroupManagerDistributionPolicyTargetShapeEnum",
										Description: "The distribution shape to which the group converges either proactively or on resize events (depending on the value set in `updatePolicy.instanceRedistributionType`). Possible values: TARGET_SHAPE_UNSPECIFIED, ANY, BALANCED, ANY_SINGLE_ZONE",
										Enum: []string{
											"TARGET_SHAPE_UNSPECIFIED",
											"ANY",
											"BALANCED",
											"ANY_SINGLE_ZONE",
										},
									},
									"zones": &dcl.Property{
										Type:        "array",
										GoName:      "Zones",
										Description: "Zones where the regional managed instance group will create and manage its instances.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "InstanceGroupManagerDistributionPolicyZones",
											Properties: map[string]*dcl.Property{
												"zone": &dcl.Property{
													Type:        "string",
													GoName:      "Zone",
													Description: "The URL of the [zone](/compute/docs/regions-zones/#available). The zone must exist in the region where the managed instance group is located.",
													Immutable:   true,
												},
											},
										},
									},
								},
							},
							"fingerprint": &dcl.Property{
								Type:        "string",
								GoName:      "Fingerprint",
								ReadOnly:    true,
								Description: "Fingerprint of this resource. This field may be used in optimistic locking. It will be ignored when inserting an InstanceGroupManager. An up-to-date fingerprint must be provided in order to update the InstanceGroupManager, otherwise the request will fail with error `412 conditionNotMet`. To see the latest fingerprint, make a `get()` request to retrieve an InstanceGroupManager.",
								Immutable:   true,
							},
							"id": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "Id",
								ReadOnly:    true,
								Description: "[Output Only] A unique identifier for this resource type. The server generates this identifier.",
								Immutable:   true,
							},
							"instanceGroup": &dcl.Property{
								Type:        "string",
								GoName:      "InstanceGroup",
								ReadOnly:    true,
								Description: "[Output Only] The URL of the Instance Group resource.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/InstanceGroup",
										Field:    "selfLink",
									},
								},
							},
							"instanceTemplate": &dcl.Property{
								Type:        "string",
								GoName:      "InstanceTemplate",
								Description: "The URL of the instance template that is specified for this managed instance group. The group uses this template to create all new instances in the managed instance group. The templates for existing instances in the group do not change unless you run `recreateInstances`, run `applyUpdatesToInstances`, or set the group's `updatePolicy.type` to `PROACTIVE`.",
								Conflicts: []string{
									"versions",
								},
								ServerDefault: true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Compute/InstanceTemplate",
										Field:    "selfLink",
									},
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location of this resource.",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The name of the managed instance group. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt).",
								Immutable:   true,
							},
							"namedPorts": &dcl.Property{
								Type:        "array",
								GoName:      "NamedPorts",
								Description: "Named ports configured for the Instance Groups complementary to this Instance Group Manager.",
								Immutable:   true,
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "InstanceGroupManagerNamedPorts",
									Properties: map[string]*dcl.Property{
										"name": &dcl.Property{
											Type:        "string",
											GoName:      "Name",
											Description: "The name for this named port. The name must be 1-63 characters long, and comply with [RFC1035](https://www.ietf.org/rfc/rfc1035.txt).",
											Immutable:   true,
										},
										"port": &dcl.Property{
											Type:        "integer",
											Format:      "int64",
											GoName:      "Port",
											Description: "The port number, which can be a value between 1 and 65535.",
											Immutable:   true,
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
									},
								},
							},
							"region": &dcl.Property{
								Type:        "string",
								GoName:      "Region",
								ReadOnly:    true,
								Description: "[Output Only] The URL of the [region](/compute/docs/regions-zones/#available) where the managed instance group resides (for regional resources).",
								Immutable:   true,
							},
							"selfLink": &dcl.Property{
								Type:        "string",
								GoName:      "SelfLink",
								ReadOnly:    true,
								Description: "[Output Only] The URL for this managed instance group. The server defines this URL.",
								Immutable:   true,
							},
							"statefulPolicy": &dcl.Property{
								Type:        "object",
								GoName:      "StatefulPolicy",
								GoType:      "InstanceGroupManagerStatefulPolicy",
								Description: "Stateful configuration for this Instanced Group Manager",
								Properties: map[string]*dcl.Property{
									"preservedState": &dcl.Property{
										Type:   "object",
										GoName: "PreservedState",
										GoType: "InstanceGroupManagerStatefulPolicyPreservedState",
										Properties: map[string]*dcl.Property{
											"disks": &dcl.Property{
												Type: "object",
												AdditionalProperties: &dcl.Property{
													Type:   "object",
													GoType: "InstanceGroupManagerStatefulPolicyPreservedStateDisks",
													Properties: map[string]*dcl.Property{
														"autoDelete": &dcl.Property{
															Type:        "string",
															GoName:      "AutoDelete",
															GoType:      "InstanceGroupManagerStatefulPolicyPreservedStateDisksAutoDeleteEnum",
															Description: "These stateful disks will never be deleted during autohealing, update or VM instance recreate operations. This flag is used to configure if the disk should be deleted after it is no longer used by the group, e.g. when the given instance or the whole group is deleted. Note: disks attached in READ_ONLY mode cannot be auto-deleted. Possible values: NEVER, ON_PERMANENT_INSTANCE_DELETION",
															Enum: []string{
																"NEVER",
																"ON_PERMANENT_INSTANCE_DELETION",
															},
														},
													},
												},
												GoName:      "Disks",
												Description: "Disks created on the instances that will be preserved on instance delete, update, etc. This map is keyed with the device names of the disks.",
											},
										},
									},
								},
							},
							"status": &dcl.Property{
								Type:        "object",
								GoName:      "Status",
								GoType:      "InstanceGroupManagerStatus",
								ReadOnly:    true,
								Description: "[Output Only] The status of this managed instance group.",
								Properties: map[string]*dcl.Property{
									"autoscaler": &dcl.Property{
										Type:        "string",
										GoName:      "Autoscaler",
										ReadOnly:    true,
										Description: "[Output Only] The URL of the [Autoscaler](/compute/docs/autoscaler/) that targets this instance group manager.",
										Immutable:   true,
									},
									"isStable": &dcl.Property{
										Type:        "boolean",
										GoName:      "IsStable",
										ReadOnly:    true,
										Description: "[Output Only] A bit indicating whether the managed instance group is in a stable state. A stable state means that: none of the instances in the managed instance group is currently undergoing any type of change (for example, creation, restart, or deletion); no future changes are scheduled for instances in the managed instance group; and the managed instance group itself is not being modified.",
										Immutable:   true,
									},
									"stateful": &dcl.Property{
										Type:        "object",
										GoName:      "Stateful",
										GoType:      "InstanceGroupManagerStatusStateful",
										ReadOnly:    true,
										Description: "[Output Only] Stateful status of the given Instance Group Manager.",
										Properties: map[string]*dcl.Property{
											"hasStatefulConfig": &dcl.Property{
												Type:        "boolean",
												GoName:      "HasStatefulConfig",
												ReadOnly:    true,
												Description: "[Output Only] A bit indicating whether the managed instance group has stateful configuration, that is, if you have configured any items in a stateful policy or in per-instance configs. The group might report that it has no stateful config even when there is still some preserved state on a managed instance, for example, if you have deleted all PICs but not yet applied those deletions.",
												Immutable:   true,
											},
											"perInstanceConfigs": &dcl.Property{
												Type:        "object",
												GoName:      "PerInstanceConfigs",
												GoType:      "InstanceGroupManagerStatusStatefulPerInstanceConfigs",
												ReadOnly:    true,
												Description: "[Output Only] Status of per-instance configs on the instance.",
												Properties: map[string]*dcl.Property{
													"allEffective": &dcl.Property{
														Type:        "boolean",
														GoName:      "AllEffective",
														Description: "A bit indicating if all of the group's per-instance configs (listed in the output of a listPerInstanceConfigs API call) have status `EFFECTIVE` or there are no per-instance-configs.",
													},
												},
											},
										},
									},
									"versionTarget": &dcl.Property{
										Type:        "object",
										GoName:      "VersionTarget",
										GoType:      "InstanceGroupManagerStatusVersionTarget",
										ReadOnly:    true,
										Description: "[Output Only] A status of consistency of Instances' versions with their target version specified by `version` field on Instance Group Manager.",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"isReached": &dcl.Property{
												Type:        "boolean",
												GoName:      "IsReached",
												ReadOnly:    true,
												Description: "[Output Only] A bit indicating whether version target has been reached in this managed instance group, i.e. all instances are in their target version. Instances' target version are specified by `version` field on Instance Group Manager.",
												Immutable:   true,
											},
										},
									},
								},
							},
							"targetPools": &dcl.Property{
								Type:        "array",
								GoName:      "TargetPools",
								Description: "The URLs for all TargetPool resources to which instances in the `instanceGroup` field are added. The target pools automatically apply to all of the instances in the managed instance group.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
									ResourceReferences: []*dcl.PropertyResourceReference{
										&dcl.PropertyResourceReference{
											Resource: "Compute/TargetPool",
											Field:    "selfLink",
										},
									},
								},
							},
							"targetSize": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "TargetSize",
								Description: "The target number of running instances for this managed instance group. You can reduce this number by using the instanceGroupManager deleteInstances or abandonInstances methods. Resizing the group also changes this number.",
							},
							"updatePolicy": &dcl.Property{
								Type:          "object",
								GoName:        "UpdatePolicy",
								GoType:        "InstanceGroupManagerUpdatePolicy",
								Description:   "The update policy for this managed instance group.",
								ServerDefault: true,
								Properties: map[string]*dcl.Property{
									"instanceRedistributionType": &dcl.Property{
										Type:        "string",
										GoName:      "InstanceRedistributionType",
										GoType:      "InstanceGroupManagerUpdatePolicyInstanceRedistributionTypeEnum",
										Description: "The [instance redistribution policy](/compute/docs/instance-groups/regional-migs#proactive_instance_redistribution) for regional managed instance groups. Valid values are: - `PROACTIVE` (default): The group attempts to maintain an even distribution of VM instances across zones in the region. - `NONE`: For non-autoscaled groups, proactive redistribution is disabled.",
										Enum: []string{
											"NONE",
											"PROACTIVE",
										},
									},
									"maxSurge": &dcl.Property{
										Type:        "object",
										GoName:      "MaxSurge",
										GoType:      "InstanceGroupManagerUpdatePolicyMaxSurge",
										Description: "The maximum number of instances that can be created above the specified `targetSize` during the update process. This value can be either a fixed number or, if the group has 10 or more instances, a percentage. If you set a percentage, the number of instances is rounded if necessary. The default value for `maxSurge` is a fixed value equal to the number of zones in which the managed instance group operates. At least one of either `maxSurge` or `maxUnavailable` must be greater than 0. Learn more about [`maxSurge`](/compute/docs/instance-groups/rolling-out-updates-to-managed-instance-groups#max_surge).",
										SendEmpty:   true,
										Properties: map[string]*dcl.Property{
											"calculated": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Calculated",
												ReadOnly:    true,
												Description: "[Output Only] Absolute value of VM instances calculated based on the specific mode. - If the value is `fixed`, then the `calculated` value is equal to the `fixed` value. - If the value is a `percent`, then the `calculated` value is `percent`/100 * `targetSize`. For example, the `calculated` value of a 80% of a managed instance group with 150 instances would be (80/100 * 150) = 120 VM instances. If there is a remainder, the number is rounded.",
											},
											"fixed": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Fixed",
												Description: "Specifies a fixed number of VM instances. This must be a positive integer.",
												SendEmpty:   true,
											},
											"percent": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Percent",
												Description: "Specifies a percentage of instances between 0 to 100%, inclusive. For example, specify `80` for 80%.",
												SendEmpty:   true,
											},
										},
									},
									"maxUnavailable": &dcl.Property{
										Type:        "object",
										GoName:      "MaxUnavailable",
										GoType:      "InstanceGroupManagerUpdatePolicyMaxUnavailable",
										Description: "The maximum number of instances that can be unavailable during the update process. An instance is considered available if all of the following conditions are satisfied: - The instance's [status](/compute/docs/instances/checking-instance-status) is `RUNNING`. - If there is a [health check](/compute/docs/instance-groups/autohealing-instances-in-migs) on the instance group, the instance's health check status must be `HEALTHY` at least once. If there is no health check on the group, then the instance only needs to have a status of `RUNNING` to be considered available. This value can be either a fixed number or, if the group has 10 or more instances, a percentage. If you set a percentage, the number of instances is rounded if necessary. The default value for `maxUnavailable` is a fixed value equal to the number of zones in which the managed instance group operates. At least one of either `maxSurge` or `maxUnavailable` must be greater than 0. Learn more about [`maxUnavailable`](/compute/docs/instance-groups/rolling-out-updates-to-managed-instance-groups#max_unavailable).",
										Properties: map[string]*dcl.Property{
											"calculated": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Calculated",
												ReadOnly:    true,
												Description: "[Output Only] Absolute value of VM instances calculated based on the specific mode. - If the value is `fixed`, then the `calculated` value is equal to the `fixed` value. - If the value is a `percent`, then the `calculated` value is `percent`/100 * `targetSize`. For example, the `calculated` value of a 80% of a managed instance group with 150 instances would be (80/100 * 150) = 120 VM instances. If there is a remainder, the number is rounded.",
											},
											"fixed": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Fixed",
												Description: "Specifies a fixed number of VM instances. This must be a positive integer.",
												SendEmpty:   true,
											},
											"percent": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "Percent",
												Description: "Specifies a percentage of instances between 0 to 100%, inclusive. For example, specify `80` for 80%.",
												SendEmpty:   true,
											},
										},
									},
									"minimalAction": &dcl.Property{
										Type:        "string",
										GoName:      "MinimalAction",
										GoType:      "InstanceGroupManagerUpdatePolicyMinimalActionEnum",
										Description: "Minimal action to be taken on an instance. You can specify either `RESTART` to restart existing instances or `REPLACE` to delete and create new instances from the target template. If you specify a `RESTART`, the Updater will attempt to perform that action only. However, if the Updater determines that the minimal action you specify is not enough to perform the update, it might perform a more disruptive action.",
										Enum: []string{
											"REPLACE",
											"RESTART",
											"REFRESH",
											"NONE",
										},
									},
									"replacementMethod": &dcl.Property{
										Type:        "string",
										GoName:      "ReplacementMethod",
										GoType:      "InstanceGroupManagerUpdatePolicyReplacementMethodEnum",
										Description: "What action should be used to replace instances. See minimal_action.REPLACE Possible values: SUBSTITUTE, RECREATE",
										Enum: []string{
											"SUBSTITUTE",
											"RECREATE",
										},
									},
									"type": &dcl.Property{
										Type:        "string",
										GoName:      "Type",
										GoType:      "InstanceGroupManagerUpdatePolicyTypeEnum",
										Description: "The type of update process. You can specify either `PROACTIVE` so that the instance group manager proactively executes actions in order to bring instances to their target versions or `OPPORTUNISTIC` so that no action is proactively executed but the update will be performed as part of other actions (for example, resizes or `recreateInstances` calls).",
										Enum: []string{
											"OPPORTUNISTIC",
											"PROACTIVE",
										},
									},
								},
							},
							"versions": &dcl.Property{
								Type:        "array",
								GoName:      "Versions",
								Description: "Specifies the instance templates used by this managed instance group to create instances. Each version is defined by an `instanceTemplate` and a `name`. Every version can appear at most once per instance group. This field overrides the top-level `instanceTemplate` field. Read more about the [relationships between these fields](/compute/docs/instance-groups/rolling-out-updates-to-managed-instance-groups#relationship_between_versions_and_instancetemplate_properties_for_a_managed_instance_group). Exactly one `version` must leave the `targetSize` field unset. That version will be applied to all remaining instances. For more information, read about [canary updates](/compute/docs/instance-groups/rolling-out-updates-to-managed-instance-groups#starting_a_canary_update).",
								Conflicts: []string{
									"instanceTemplate",
								},
								ServerDefault: true,
								SendEmpty:     true,
								ListType:      "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "InstanceGroupManagerVersions",
									Properties: map[string]*dcl.Property{
										"instanceTemplate": &dcl.Property{
											Type:        "string",
											GoName:      "InstanceTemplate",
											Description: "The URL of the instance template that is specified for this managed instance group. The group uses this template to create new instances in the managed instance group until the `targetSize` for this version is reached. The templates for existing instances in the group do not change unless you run `recreateInstances`, run `applyUpdatesToInstances`, or set the group's `updatePolicy.type` to `PROACTIVE`; in those cases, existing instances are updated until the `targetSize` for this version is reached.",
											ResourceReferences: []*dcl.PropertyResourceReference{
												&dcl.PropertyResourceReference{
													Resource: "Compute/InstanceTemplate",
													Field:    "selfLink",
												},
											},
										},
										"name": &dcl.Property{
											Type:        "string",
											GoName:      "Name",
											Description: "Name of the version. Unique among all versions in the scope of this managed instance group.",
										},
										"targetSize": &dcl.Property{
											Type:        "object",
											GoName:      "TargetSize",
											GoType:      "InstanceGroupManagerVersionsTargetSize",
											Description: "Specifies the intended number of instances to be created from the `instanceTemplate`. The final number of instances created from the template will be equal to: - If expressed as a fixed number, the minimum of either `targetSize.fixed` or `instanceGroupManager.targetSize` is used. - if expressed as a `percent`, the `targetSize` would be `(targetSize.percent/100 * InstanceGroupManager.targetSize)` If there is a remainder, the number is rounded. If unset, this version will update any remaining instances not updated by another `version`. Read [Starting a canary update](/compute/docs/instance-groups/rolling-out-updates-to-managed-instance-groups#starting_a_canary_update) for more information.",
											Properties: map[string]*dcl.Property{
												"calculated": &dcl.Property{
													Type:        "integer",
													Format:      "int64",
													GoName:      "Calculated",
													ReadOnly:    true,
													Description: "[Output Only] Absolute value of VM instances calculated based on the specific mode. - If the value is `fixed`, then the `calculated` value is equal to the `fixed` value. - If the value is a `percent`, then the `calculated` value is `percent`/100 * `targetSize`. For example, the `calculated` value of a 80% of a managed instance group with 150 instances would be (80/100 * 150) = 120 VM instances. If there is a remainder, the number is rounded.",
												},
												"fixed": &dcl.Property{
													Type:        "integer",
													Format:      "int64",
													GoName:      "Fixed",
													Description: "Specifies a fixed number of VM instances. This must be a positive integer.",
													SendEmpty:   true,
												},
												"percent": &dcl.Property{
													Type:        "integer",
													Format:      "int64",
													GoName:      "Percent",
													Description: "Specifies a percentage of instances between 0 to 100%, inclusive. For example, specify `80` for 80%.",
													SendEmpty:   true,
												},
											},
										},
									},
								},
							},
							"zone": &dcl.Property{
								Type:        "string",
								GoName:      "Zone",
								ReadOnly:    true,
								Description: "[Output Only] The URL of a [zone](/compute/docs/regions-zones/#available) where the managed instance group is located (for zonal resources).",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
