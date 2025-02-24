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
package dataproc

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLAutoscalingPolicySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Dataproc/AutoscalingPolicy",
			Description: "The Dataproc AutoscalingPolicy resource",
			StructName:  "AutoscalingPolicy",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a AutoscalingPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "autoscalingPolicy",
						Required:    true,
						Description: "A full instance of a AutoscalingPolicy",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a AutoscalingPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "autoscalingPolicy",
						Required:    true,
						Description: "A full instance of a AutoscalingPolicy",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a AutoscalingPolicy",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "autoscalingPolicy",
						Required:    true,
						Description: "A full instance of a AutoscalingPolicy",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all AutoscalingPolicy",
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
				Description: "The function used to list information about many AutoscalingPolicy",
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
				"AutoscalingPolicy": &dcl.Component{
					Title:           "AutoscalingPolicy",
					ID:              "projects/{{project}}/locations/{{location}}/autoscalingPolicies/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"basicAlgorithm",
							"workerConfig",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"basicAlgorithm": &dcl.Property{
								Type:   "object",
								GoName: "BasicAlgorithm",
								GoType: "AutoscalingPolicyBasicAlgorithm",
								Required: []string{
									"yarnConfig",
								},
								Properties: map[string]*dcl.Property{
									"cooldownPeriod": &dcl.Property{
										Type:          "string",
										GoName:        "CooldownPeriod",
										Description:   "Optional. Duration between scaling events. A scaling period starts after the update operation from the previous event has completed. Bounds: . Default: 2m.",
										ServerDefault: true,
									},
									"yarnConfig": &dcl.Property{
										Type:        "object",
										GoName:      "YarnConfig",
										GoType:      "AutoscalingPolicyBasicAlgorithmYarnConfig",
										Description: "Required. YARN autoscaling configuration.",
										Required: []string{
											"gracefulDecommissionTimeout",
											"scaleUpFactor",
											"scaleDownFactor",
										},
										Properties: map[string]*dcl.Property{
											"gracefulDecommissionTimeout": &dcl.Property{
												Type:        "string",
												GoName:      "GracefulDecommissionTimeout",
												Description: "Required. Timeout for YARN graceful decommissioning of Node Managers. Specifies the duration to wait for jobs to complete before forcefully removing workers (and potentially interrupting jobs). Only applicable to downscaling operations.",
											},
											"scaleDownFactor": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "ScaleDownFactor",
												Description: "Required. Fraction of average YARN pending memory in the last cooldown period for which to remove workers. A scale-down factor of 1 will result in scaling down so that there is no available memory remaining after the update (more aggressive scaling). A scale-down factor of 0 disables removing workers, which can be beneficial for autoscaling a single job. See .",
											},
											"scaleDownMinWorkerFraction": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "ScaleDownMinWorkerFraction",
												Description: "Optional. Minimum scale-down threshold as a fraction of total cluster size before scaling occurs. For example, in a 20-worker cluster, a threshold of 0.1 means the autoscaler must recommend at least a 2 worker scale-down for the cluster to scale. A threshold of 0 means the autoscaler will scale down on any recommended change. Bounds: . Default: 0.0.",
											},
											"scaleUpFactor": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "ScaleUpFactor",
												Description: "Required. Fraction of average YARN pending memory in the last cooldown period for which to add workers. A scale-up factor of 1.0 will result in scaling up so that there is no pending memory remaining after the update (more aggressive scaling). A scale-up factor closer to 0 will result in a smaller magnitude of scaling up (less aggressive scaling). See .",
											},
											"scaleUpMinWorkerFraction": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "ScaleUpMinWorkerFraction",
												Description: "Optional. Minimum scale-up threshold as a fraction of total cluster size before scaling occurs. For example, in a 20-worker cluster, a threshold of 0.1 means the autoscaler must recommend at least a 2-worker scale-up for the cluster to scale. A threshold of 0 means the autoscaler will scale up on any recommended change. Bounds: . Default: 0.0.",
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
								Parameter:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The \"resource name\" of the autoscaling policy, as described in https://cloud.google.com/apis/design/resource_names. * For `projects.regions.autoscalingPolicies`, the resource name of the policy has the following format: `projects/{project_id}/regions/{region}/autoscalingPolicies/{policy_id}` * For `projects.locations.autoscalingPolicies`, the resource name of the policy has the following format: `projects/{project_id}/locations/{location}/autoscalingPolicies/{policy_id}`",
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
								Parameter: true,
							},
							"secondaryWorkerConfig": &dcl.Property{
								Type:        "object",
								GoName:      "SecondaryWorkerConfig",
								GoType:      "AutoscalingPolicySecondaryWorkerConfig",
								Description: "Optional. Describes how the autoscaler will operate for secondary workers.",
								Properties: map[string]*dcl.Property{
									"maxInstances": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MaxInstances",
										Description: "Optional. Maximum number of instances for this group. Note that by default, clusters will not use secondary workers. Required for secondary workers if the minimum secondary instances is set. Primary workers - Bounds: [min_instances, ). Secondary workers - Bounds: [min_instances, ). Default: 0.",
									},
									"minInstances": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MinInstances",
										Description: "Optional. Minimum number of instances for this group. Primary workers - Bounds: . Default: 0.",
									},
									"weight": &dcl.Property{
										Type:          "integer",
										Format:        "int64",
										GoName:        "Weight",
										Description:   "Optional. Weight for the instance group, which is used to determine the fraction of total workers in the cluster from this instance group. For example, if primary workers have weight 2, and secondary workers have weight 1, the cluster will have approximately 2 primary workers for each secondary worker. The cluster may not reach the specified balance if constrained by min/max bounds or other autoscaling settings. For example, if `max_instances` for secondary workers is 0, then only primary workers will be added. The cluster can also be out of balance when created. If weight is not set on any instance group, the cluster will default to equal weight for all groups: the cluster will attempt to maintain an equal number of workers in each group within the configured size bounds for each group. If weight is set for one group only, the cluster will default to zero weight on the unset group. For example if weight is set only on primary workers, the cluster will use primary workers only and no secondary workers.",
										ServerDefault: true,
									},
								},
							},
							"workerConfig": &dcl.Property{
								Type:        "object",
								GoName:      "WorkerConfig",
								GoType:      "AutoscalingPolicyWorkerConfig",
								Description: "Required. Describes how the autoscaler will operate for primary workers.",
								Required: []string{
									"maxInstances",
								},
								Properties: map[string]*dcl.Property{
									"maxInstances": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "MaxInstances",
										Description: "Required. Maximum number of instances for this group. Required for primary workers. Note that by default, clusters will not use secondary workers. Required for secondary workers if the minimum secondary instances is set. Primary workers - Bounds: [min_instances, ). Secondary workers - Bounds: [min_instances, ). Default: 0.",
									},
									"minInstances": &dcl.Property{
										Type:          "integer",
										Format:        "int64",
										GoName:        "MinInstances",
										Description:   "Optional. Minimum number of instances for this group. Primary workers - Bounds: . Default: 0.",
										ServerDefault: true,
									},
									"weight": &dcl.Property{
										Type:          "integer",
										Format:        "int64",
										GoName:        "Weight",
										Description:   "Optional. Weight for the instance group, which is used to determine the fraction of total workers in the cluster from this instance group. For example, if primary workers have weight 2, and secondary workers have weight 1, the cluster will have approximately 2 primary workers for each secondary worker. The cluster may not reach the specified balance if constrained by min/max bounds or other autoscaling settings. For example, if `max_instances` for secondary workers is 0, then only primary workers will be added. The cluster can also be out of balance when created. If weight is not set on any instance group, the cluster will default to equal weight for all groups: the cluster will attempt to maintain an equal number of workers in each group within the configured size bounds for each group. If weight is set for one group only, the cluster will default to zero weight on the unset group. For example if weight is set only on primary workers, the cluster will use primary workers only and no secondary workers.",
										ServerDefault: true,
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
