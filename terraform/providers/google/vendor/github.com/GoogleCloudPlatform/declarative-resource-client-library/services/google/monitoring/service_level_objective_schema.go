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
package monitoring

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLServiceLevelObjectiveSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Monitoring/ServiceLevelObjective",
			Description: "The Monitoring ServiceLevelObjective resource",
			StructName:  "ServiceLevelObjective",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a ServiceLevelObjective",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "serviceLevelObjective",
						Required:    true,
						Description: "A full instance of a ServiceLevelObjective",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a ServiceLevelObjective",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "serviceLevelObjective",
						Required:    true,
						Description: "A full instance of a ServiceLevelObjective",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a ServiceLevelObjective",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "serviceLevelObjective",
						Required:    true,
						Description: "A full instance of a ServiceLevelObjective",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all ServiceLevelObjective",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "service",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many ServiceLevelObjective",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "service",
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
				"ServiceLevelObjective": &dcl.Component{
					Title:           "ServiceLevelObjective",
					ID:              "projects/{{project}}/services/{{service}}/serviceLevelObjectives/{{name}}",
					ParentContainer: "project",
					LabelsField:     "userLabels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"goal",
							"project",
							"service",
						},
						Properties: map[string]*dcl.Property{
							"calendarPeriod": &dcl.Property{
								Type:        "string",
								GoName:      "CalendarPeriod",
								GoType:      "ServiceLevelObjectiveCalendarPeriodEnum",
								Description: "A calendar period, semantically \"since the start of the current ``\". At this time, only `DAY`, `WEEK`, `FORTNIGHT`, and `MONTH` are supported. Possible values: CALENDAR_PERIOD_UNSPECIFIED, DAY, WEEK, FORTNIGHT, MONTH, QUARTER, HALF, YEAR",
								Conflicts: []string{
									"rollingPeriod",
								},
								Enum: []string{
									"CALENDAR_PERIOD_UNSPECIFIED",
									"DAY",
									"WEEK",
									"FORTNIGHT",
									"MONTH",
									"QUARTER",
									"HALF",
									"YEAR",
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Time stamp of the `Create` or most recent `Update` command on this `Slo`.",
								Immutable:   true,
							},
							"deleteTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "DeleteTime",
								ReadOnly:    true,
								Description: "Time stamp of the `Update` or `Delete` command that made this no longer a current `Slo`. This field is not populated in `ServiceLevelObjective`s returned from calls to `GetServiceLevelObjective` and `ListServiceLevelObjectives`, because it is always empty in the current version. It is populated in `ServiceLevelObjective`s representing previous versions in the output of `ListServiceLevelObjectiveVersions`. Because all old configuration versions are stored, `Update` operations mark the obsoleted version as deleted.",
								Immutable:   true,
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "Name used for UI elements listing this SLO.",
							},
							"goal": &dcl.Property{
								Type:        "number",
								Format:      "double",
								GoName:      "Goal",
								Description: "The fraction of service that must be good in order for this objective to be met. `0 < goal <= 0.999`.",
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Resource name for this `ServiceLevelObjective`. The format is: projects/[PROJECT_ID_OR_NUMBER]/services/[SERVICE_ID]/serviceLevelObjectives/[SLO_NAME]",
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
							"rollingPeriod": &dcl.Property{
								Type:        "string",
								GoName:      "RollingPeriod",
								Description: "A rolling time period, semantically \"in the past ``\". Must be an integer multiple of 1 day no larger than 30 days.",
								Conflicts: []string{
									"calendarPeriod",
								},
							},
							"service": &dcl.Property{
								Type:        "string",
								GoName:      "Service",
								Description: "The service for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Monitoring/Service",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"serviceLevelIndicator": &dcl.Property{
								Type:        "object",
								GoName:      "ServiceLevelIndicator",
								GoType:      "ServiceLevelObjectiveServiceLevelIndicator",
								Description: "The definition of good service, used to measure and calculate the quality of the `Service`'s performance with respect to a single aspect of service quality.",
								Properties: map[string]*dcl.Property{
									"basicSli": &dcl.Property{
										Type:        "object",
										GoName:      "BasicSli",
										GoType:      "ServiceLevelObjectiveServiceLevelIndicatorBasicSli",
										Description: "Basic SLI on a well-known service type.",
										Conflicts: []string{
											"requestBased",
											"windowsBased",
										},
										Properties: map[string]*dcl.Property{
											"availability": &dcl.Property{
												Type:        "object",
												GoName:      "Availability",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorBasicSliAvailability",
												Description: "Good service is defined to be the count of requests made to this service that return successfully.",
												Conflicts: []string{
													"latency",
													"operationAvailability",
													"operationLatency",
												},
												Properties: map[string]*dcl.Property{},
											},
											"latency": &dcl.Property{
												Type:        "object",
												GoName:      "Latency",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatency",
												Description: "Good service is defined to be the count of requests made to this service that are fast enough with respect to `latency.threshold`.",
												Conflicts: []string{
													"availability",
													"operationAvailability",
													"operationLatency",
												},
												Properties: map[string]*dcl.Property{
													"experience": &dcl.Property{
														Type:        "string",
														GoName:      "Experience",
														GoType:      "ServiceLevelObjectiveServiceLevelIndicatorBasicSliLatencyExperienceEnum",
														Description: "A description of the experience associated with failing requests. Possible values: LATENCY_EXPERIENCE_UNSPECIFIED, DELIGHTING, SATISFYING, ANNOYING",
														Enum: []string{
															"LATENCY_EXPERIENCE_UNSPECIFIED",
															"DELIGHTING",
															"SATISFYING",
															"ANNOYING",
														},
													},
													"threshold": &dcl.Property{
														Type:        "string",
														GoName:      "Threshold",
														Description: "Good service is defined to be the count of requests made to this service that return in no more than `threshold`.",
													},
												},
											},
											"location": &dcl.Property{
												Type:        "array",
												GoName:      "Location",
												Description: "OPTIONAL: The set of locations to which this SLI is relevant. Telemetry from other locations will not be used to calculate performance for this SLI. If omitted, this SLI applies to all locations in which the Service has activity. For service types that don't support breaking down by location, setting this field will result in an error.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
											"method": &dcl.Property{
												Type:        "array",
												GoName:      "Method",
												Description: "OPTIONAL: The set of RPCs to which this SLI is relevant. Telemetry from other methods will not be used to calculate performance for this SLI. If omitted, this SLI applies to all the Service's methods. For service types that don't support breaking down by method, setting this field will result in an error.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
											"operationAvailability": &dcl.Property{
												Type:        "object",
												GoName:      "OperationAvailability",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationAvailability",
												Description: "Good service is defined to be the count of operations performed by this service that return successfully",
												Conflicts: []string{
													"availability",
													"latency",
													"operationLatency",
												},
												Properties: map[string]*dcl.Property{},
											},
											"operationLatency": &dcl.Property{
												Type:        "object",
												GoName:      "OperationLatency",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatency",
												Description: "Good service is defined to be the count of operations performed by this service that are fast enough with respect to `operation_latency.threshold`.",
												Conflicts: []string{
													"availability",
													"latency",
													"operationAvailability",
												},
												Properties: map[string]*dcl.Property{
													"experience": &dcl.Property{
														Type:        "string",
														GoName:      "Experience",
														GoType:      "ServiceLevelObjectiveServiceLevelIndicatorBasicSliOperationLatencyExperienceEnum",
														Description: "A description of the experience associated with failing requests. Possible values: LATENCY_EXPERIENCE_UNSPECIFIED, DELIGHTING, SATISFYING, ANNOYING",
														Enum: []string{
															"LATENCY_EXPERIENCE_UNSPECIFIED",
															"DELIGHTING",
															"SATISFYING",
															"ANNOYING",
														},
													},
													"threshold": &dcl.Property{
														Type:        "string",
														GoName:      "Threshold",
														Description: "Good service is defined to be the count of operations that are completed in no more than `threshold`.",
													},
												},
											},
											"version": &dcl.Property{
												Type:        "array",
												GoName:      "Version",
												Description: "OPTIONAL: The set of API versions to which this SLI is relevant. Telemetry from other API versions will not be used to calculate performance for this SLI. If omitted, this SLI applies to all API versions. For service types that don't support breaking down by version, setting this field will result in an error.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
										},
									},
									"requestBased": &dcl.Property{
										Type:        "object",
										GoName:      "RequestBased",
										GoType:      "ServiceLevelObjectiveServiceLevelIndicatorRequestBased",
										Description: "Request-based SLIs",
										Conflicts: []string{
											"basicSli",
											"windowsBased",
										},
										Properties: map[string]*dcl.Property{
											"distributionCut": &dcl.Property{
												Type:        "object",
												GoName:      "DistributionCut",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCut",
												Description: "`distribution_cut` is used when `good_service` is a count of values aggregated in a `Distribution` that fall into a good range. The `total_service` is the total count of all values aggregated in the `Distribution`.",
												Conflicts: []string{
													"goodTotalRatio",
												},
												Properties: map[string]*dcl.Property{
													"distributionFilter": &dcl.Property{
														Type:        "string",
														GoName:      "DistributionFilter",
														Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` aggregating values. Must have `ValueType = DISTRIBUTION` and `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
													},
													"range": &dcl.Property{
														Type:        "object",
														GoName:      "Range",
														GoType:      "ServiceLevelObjectiveServiceLevelIndicatorRequestBasedDistributionCutRange",
														Description: "Range of values considered \"good.\" For a one-sided range, set one bound to an infinite value.",
														Properties: map[string]*dcl.Property{
															"max": &dcl.Property{
																Type:        "number",
																Format:      "double",
																GoName:      "Max",
																Description: "Range maximum.",
															},
															"min": &dcl.Property{
																Type:        "number",
																Format:      "double",
																GoName:      "Min",
																Description: "Range minimum.",
															},
														},
													},
												},
											},
											"goodTotalRatio": &dcl.Property{
												Type:        "object",
												GoName:      "GoodTotalRatio",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorRequestBasedGoodTotalRatio",
												Description: "`good_total_ratio` is used when the ratio of `good_service` to `total_service` is computed from two `TimeSeries`.",
												Conflicts: []string{
													"distributionCut",
												},
												Properties: map[string]*dcl.Property{
													"badServiceFilter": &dcl.Property{
														Type:        "string",
														GoName:      "BadServiceFilter",
														Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` quantifying bad service, either demanded service that was not provided or demanded service that was of inadequate quality. Must have `ValueType = DOUBLE` or `ValueType = INT64` and must have `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
													},
													"goodServiceFilter": &dcl.Property{
														Type:        "string",
														GoName:      "GoodServiceFilter",
														Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` quantifying good service provided. Must have `ValueType = DOUBLE` or `ValueType = INT64` and must have `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
													},
													"totalServiceFilter": &dcl.Property{
														Type:        "string",
														GoName:      "TotalServiceFilter",
														Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` quantifying total demanded service. Must have `ValueType = DOUBLE` or `ValueType = INT64` and must have `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
													},
												},
											},
										},
									},
									"windowsBased": &dcl.Property{
										Type:        "object",
										GoName:      "WindowsBased",
										GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBased",
										Description: "Windows-based SLIs",
										Conflicts: []string{
											"basicSli",
											"requestBased",
										},
										Properties: map[string]*dcl.Property{
											"goodBadMetricFilter": &dcl.Property{
												Type:        "string",
												GoName:      "GoodBadMetricFilter",
												Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` with `ValueType = BOOL`. The window is good if any `true` values appear in the window.",
												Conflicts: []string{
													"goodTotalRatioThreshold",
													"metricMeanInRange",
													"metricSumInRange",
												},
											},
											"goodTotalRatioThreshold": &dcl.Property{
												Type:        "object",
												GoName:      "GoodTotalRatioThreshold",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThreshold",
												Description: "A window is good if its `performance` is high enough.",
												Conflicts: []string{
													"goodBadMetricFilter",
													"metricMeanInRange",
													"metricSumInRange",
												},
												Properties: map[string]*dcl.Property{
													"basicSliPerformance": &dcl.Property{
														Type:        "object",
														GoName:      "BasicSliPerformance",
														GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformance",
														Description: "`BasicSli` to evaluate to judge window quality.",
														Conflicts: []string{
															"performance",
														},
														Properties: map[string]*dcl.Property{
															"availability": &dcl.Property{
																Type:        "object",
																GoName:      "Availability",
																GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceAvailability",
																Description: "Good service is defined to be the count of requests made to this service that return successfully.",
																Conflicts: []string{
																	"latency",
																	"operationAvailability",
																	"operationLatency",
																},
																Properties: map[string]*dcl.Property{},
															},
															"latency": &dcl.Property{
																Type:        "object",
																GoName:      "Latency",
																GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatency",
																Description: "Good service is defined to be the count of requests made to this service that are fast enough with respect to `latency.threshold`.",
																Conflicts: []string{
																	"availability",
																	"operationAvailability",
																	"operationLatency",
																},
																Properties: map[string]*dcl.Property{
																	"experience": &dcl.Property{
																		Type:        "string",
																		GoName:      "Experience",
																		GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceLatencyExperienceEnum",
																		Description: "A description of the experience associated with failing requests. Possible values: LATENCY_EXPERIENCE_UNSPECIFIED, DELIGHTING, SATISFYING, ANNOYING",
																		Enum: []string{
																			"LATENCY_EXPERIENCE_UNSPECIFIED",
																			"DELIGHTING",
																			"SATISFYING",
																			"ANNOYING",
																		},
																	},
																	"threshold": &dcl.Property{
																		Type:        "string",
																		GoName:      "Threshold",
																		Description: "Good service is defined to be the count of requests made to this service that return in no more than `threshold`.",
																	},
																},
															},
															"location": &dcl.Property{
																Type:        "array",
																GoName:      "Location",
																Description: "OPTIONAL: The set of locations to which this SLI is relevant. Telemetry from other locations will not be used to calculate performance for this SLI. If omitted, this SLI applies to all locations in which the Service has activity. For service types that don't support breaking down by location, setting this field will result in an error.",
																SendEmpty:   true,
																ListType:    "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																},
															},
															"method": &dcl.Property{
																Type:        "array",
																GoName:      "Method",
																Description: "OPTIONAL: The set of RPCs to which this SLI is relevant. Telemetry from other methods will not be used to calculate performance for this SLI. If omitted, this SLI applies to all the Service's methods. For service types that don't support breaking down by method, setting this field will result in an error.",
																SendEmpty:   true,
																ListType:    "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																},
															},
															"operationAvailability": &dcl.Property{
																Type:        "object",
																GoName:      "OperationAvailability",
																GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationAvailability",
																Description: "Good service is defined to be the count of operations performed by this service that return successfully",
																Conflicts: []string{
																	"availability",
																	"latency",
																	"operationLatency",
																},
																Properties: map[string]*dcl.Property{},
															},
															"operationLatency": &dcl.Property{
																Type:        "object",
																GoName:      "OperationLatency",
																GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatency",
																Description: "Good service is defined to be the count of operations performed by this service that are fast enough with respect to `operation_latency.threshold`.",
																Conflicts: []string{
																	"availability",
																	"latency",
																	"operationAvailability",
																},
																Properties: map[string]*dcl.Property{
																	"experience": &dcl.Property{
																		Type:        "string",
																		GoName:      "Experience",
																		GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdBasicSliPerformanceOperationLatencyExperienceEnum",
																		Description: "A description of the experience associated with failing requests. Possible values: LATENCY_EXPERIENCE_UNSPECIFIED, DELIGHTING, SATISFYING, ANNOYING",
																		Enum: []string{
																			"LATENCY_EXPERIENCE_UNSPECIFIED",
																			"DELIGHTING",
																			"SATISFYING",
																			"ANNOYING",
																		},
																	},
																	"threshold": &dcl.Property{
																		Type:        "string",
																		GoName:      "Threshold",
																		Description: "Good service is defined to be the count of operations that are completed in no more than `threshold`.",
																	},
																},
															},
															"version": &dcl.Property{
																Type:        "array",
																GoName:      "Version",
																Description: "OPTIONAL: The set of API versions to which this SLI is relevant. Telemetry from other API versions will not be used to calculate performance for this SLI. If omitted, this SLI applies to all API versions. For service types that don't support breaking down by version, setting this field will result in an error.",
																SendEmpty:   true,
																ListType:    "list",
																Items: &dcl.Property{
																	Type:   "string",
																	GoType: "string",
																},
															},
														},
													},
													"performance": &dcl.Property{
														Type:        "object",
														GoName:      "Performance",
														GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformance",
														Description: "`RequestBasedSli` to evaluate to judge window quality.",
														Conflicts: []string{
															"basicSliPerformance",
														},
														Properties: map[string]*dcl.Property{
															"distributionCut": &dcl.Property{
																Type:        "object",
																GoName:      "DistributionCut",
																GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCut",
																Description: "`distribution_cut` is used when `good_service` is a count of values aggregated in a `Distribution` that fall into a good range. The `total_service` is the total count of all values aggregated in the `Distribution`.",
																Conflicts: []string{
																	"goodTotalRatio",
																},
																Properties: map[string]*dcl.Property{
																	"distributionFilter": &dcl.Property{
																		Type:        "string",
																		GoName:      "DistributionFilter",
																		Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` aggregating values. Must have `ValueType = DISTRIBUTION` and `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
																	},
																	"range": &dcl.Property{
																		Type:        "object",
																		GoName:      "Range",
																		GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceDistributionCutRange",
																		Description: "Range of values considered \"good.\" For a one-sided range, set one bound to an infinite value.",
																		Properties: map[string]*dcl.Property{
																			"max": &dcl.Property{
																				Type:        "number",
																				Format:      "double",
																				GoName:      "Max",
																				Description: "Range maximum.",
																			},
																			"min": &dcl.Property{
																				Type:        "number",
																				Format:      "double",
																				GoName:      "Min",
																				Description: "Range minimum.",
																			},
																		},
																	},
																},
															},
															"goodTotalRatio": &dcl.Property{
																Type:        "object",
																GoName:      "GoodTotalRatio",
																GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedGoodTotalRatioThresholdPerformanceGoodTotalRatio",
																Description: "`good_total_ratio` is used when the ratio of `good_service` to `total_service` is computed from two `TimeSeries`.",
																Conflicts: []string{
																	"distributionCut",
																},
																Properties: map[string]*dcl.Property{
																	"badServiceFilter": &dcl.Property{
																		Type:        "string",
																		GoName:      "BadServiceFilter",
																		Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` quantifying bad service, either demanded service that was not provided or demanded service that was of inadequate quality. Must have `ValueType = DOUBLE` or `ValueType = INT64` and must have `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
																	},
																	"goodServiceFilter": &dcl.Property{
																		Type:        "string",
																		GoName:      "GoodServiceFilter",
																		Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` quantifying good service provided. Must have `ValueType = DOUBLE` or `ValueType = INT64` and must have `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
																	},
																	"totalServiceFilter": &dcl.Property{
																		Type:        "string",
																		GoName:      "TotalServiceFilter",
																		Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying a `TimeSeries` quantifying total demanded service. Must have `ValueType = DOUBLE` or `ValueType = INT64` and must have `MetricKind = DELTA` or `MetricKind = CUMULATIVE`.",
																	},
																},
															},
														},
													},
													"threshold": &dcl.Property{
														Type:        "number",
														Format:      "double",
														GoName:      "Threshold",
														Description: "If window `performance >= threshold`, the window is counted as good.",
													},
												},
											},
											"metricMeanInRange": &dcl.Property{
												Type:        "object",
												GoName:      "MetricMeanInRange",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRange",
												Description: "A window is good if the metric's value is in a good range, averaged across returned streams.",
												Conflicts: []string{
													"goodBadMetricFilter",
													"goodTotalRatioThreshold",
													"metricSumInRange",
												},
												Properties: map[string]*dcl.Property{
													"range": &dcl.Property{
														Type:        "object",
														GoName:      "Range",
														GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricMeanInRangeRange",
														Description: "Range of values considered \"good.\" For a one-sided range, set one bound to an infinite value.",
														Properties: map[string]*dcl.Property{
															"max": &dcl.Property{
																Type:        "number",
																Format:      "double",
																GoName:      "Max",
																Description: "Range maximum.",
															},
															"min": &dcl.Property{
																Type:        "number",
																Format:      "double",
																GoName:      "Min",
																Description: "Range minimum.",
															},
														},
													},
													"timeSeries": &dcl.Property{
														Type:        "string",
														GoName:      "TimeSeries",
														Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying the `TimeSeries` to use for evaluating window quality.",
													},
												},
											},
											"metricSumInRange": &dcl.Property{
												Type:        "object",
												GoName:      "MetricSumInRange",
												GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRange",
												Description: "A window is good if the metric's value is in a good range, summed across returned streams.",
												Conflicts: []string{
													"goodBadMetricFilter",
													"goodTotalRatioThreshold",
													"metricMeanInRange",
												},
												Properties: map[string]*dcl.Property{
													"range": &dcl.Property{
														Type:        "object",
														GoName:      "Range",
														GoType:      "ServiceLevelObjectiveServiceLevelIndicatorWindowsBasedMetricSumInRangeRange",
														Description: "Range of values considered \"good.\" For a one-sided range, set one bound to an infinite value.",
														Properties: map[string]*dcl.Property{
															"max": &dcl.Property{
																Type:        "number",
																Format:      "double",
																GoName:      "Max",
																Description: "Range maximum.",
															},
															"min": &dcl.Property{
																Type:        "number",
																Format:      "double",
																GoName:      "Min",
																Description: "Range minimum.",
															},
														},
													},
													"timeSeries": &dcl.Property{
														Type:        "string",
														GoName:      "TimeSeries",
														Description: "A [monitoring filter](https://cloud.google.com/monitoring/api/v3/filters) specifying the `TimeSeries` to use for evaluating window quality.",
													},
												},
											},
											"windowPeriod": &dcl.Property{
												Type:        "string",
												GoName:      "WindowPeriod",
												Description: "Duration over which window quality is evaluated. Must be an integer fraction of a day and at least `60s`.",
											},
										},
									},
								},
							},
							"serviceManagementOwned": &dcl.Property{
								Type:        "boolean",
								GoName:      "ServiceManagementOwned",
								ReadOnly:    true,
								Description: "Output only. If set, this SLO is managed at the [Service Management](https://cloud.google.com/service-management/overview) level. Therefore the service yaml file is the source of truth for this SLO, and API `Update` and `Delete` operations are forbidden.",
								Immutable:   true,
							},
							"userLabels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "UserLabels",
								Description: "Labels which have been used to annotate the service-level objective. Label keys must start with a letter. Label keys and values may contain lowercase letters, numbers, underscores, and dashes. Label keys and values have a maximum length of 63 characters, and must be less than 128 bytes in size. Up to 64 label entries may be stored. For labels which do not have a semantic value, the empty string may be supplied for the label value.",
							},
						},
					},
				},
			},
		},
	}
}
