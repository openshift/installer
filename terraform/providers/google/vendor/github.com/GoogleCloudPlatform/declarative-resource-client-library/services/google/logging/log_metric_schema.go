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
package logging

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLLogMetricSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Logging/LogMetric",
			Description: "The Logging LogMetric resource",
			StructName:  "LogMetric",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a LogMetric",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logMetric",
						Required:    true,
						Description: "A full instance of a LogMetric",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a LogMetric",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logMetric",
						Required:    true,
						Description: "A full instance of a LogMetric",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a LogMetric",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logMetric",
						Required:    true,
						Description: "A full instance of a LogMetric",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all LogMetric",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many LogMetric",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "project",
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
				"LogMetric": &dcl.Component{
					Title:           "LogMetric",
					ID:              "projects/{{project}}/metrics/{{name}}",
					UsesStateHint:   true,
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"filter",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"bucketOptions": &dcl.Property{
								Type:        "object",
								GoName:      "BucketOptions",
								GoType:      "LogMetricBucketOptions",
								Description: "Optional. The `bucket_options` are required when the logs-based metric is using a DISTRIBUTION value type and it describes the bucket boundaries used to create a histogram of the extracted values.",
								Properties: map[string]*dcl.Property{
									"explicitBuckets": &dcl.Property{
										Type:        "object",
										GoName:      "ExplicitBuckets",
										GoType:      "LogMetricBucketOptionsExplicitBuckets",
										Description: "The explicit buckets.",
										Conflicts: []string{
											"linearBuckets",
											"exponentialBuckets",
										},
										Properties: map[string]*dcl.Property{
											"bounds": &dcl.Property{
												Type:        "array",
												GoName:      "Bounds",
												Description: "The values must be monotonically increasing.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "number",
													Format: "double",
													GoType: "float64",
												},
											},
										},
									},
									"exponentialBuckets": &dcl.Property{
										Type:        "object",
										GoName:      "ExponentialBuckets",
										GoType:      "LogMetricBucketOptionsExponentialBuckets",
										Description: "The exponential buckets.",
										Conflicts: []string{
											"linearBuckets",
											"explicitBuckets",
										},
										Properties: map[string]*dcl.Property{
											"growthFactor": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "GrowthFactor",
												Description: "Must be greater than 1.",
											},
											"numFiniteBuckets": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "NumFiniteBuckets",
												Description: "Must be greater than 0.",
											},
											"scale": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "Scale",
												Description: "Must be greater than 0.",
											},
										},
									},
									"linearBuckets": &dcl.Property{
										Type:        "object",
										GoName:      "LinearBuckets",
										GoType:      "LogMetricBucketOptionsLinearBuckets",
										Description: "The linear bucket.",
										Conflicts: []string{
											"exponentialBuckets",
											"explicitBuckets",
										},
										Properties: map[string]*dcl.Property{
											"numFiniteBuckets": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "NumFiniteBuckets",
												Description: "Must be greater than 0.",
											},
											"offset": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "Offset",
												Description: "Lower bound of the first bucket.",
											},
											"width": &dcl.Property{
												Type:        "number",
												Format:      "double",
												GoName:      "Width",
												Description: "Must be greater than 0.",
											},
										},
									},
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The creation timestamp of the metric. This field may not be present for older metrics.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. A description of this metric, which is used in documentation. The maximum length of the description is 8000 characters.",
							},
							"disabled": &dcl.Property{
								Type:        "boolean",
								GoName:      "Disabled",
								Description: "Optional. If set to True, then this metric is disabled and it does not generate any points.",
							},
							"filter": &dcl.Property{
								Type:        "string",
								GoName:      "Filter",
								Description: "Required. An [advanced logs filter](https://cloud.google.com/logging/docs/view/advanced_filters) which is used to match log entries. Example: \"resource.type=gae_app AND severity>=ERROR\" The maximum length of the filter is 20000 characters.",
							},
							"labelExtractors": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "LabelExtractors",
								Description: "Optional. A map from a label key string to an extractor expression which is used to extract data from a log entry field and assign as the label value. Each label key specified in the LabelDescriptor must have an associated extractor expression in this map. The syntax of the extractor expression is the same as for the `value_extractor` field. The extracted value is converted to the type defined in the label descriptor. If the either the extraction or the type conversion fails, the label will have a default value. The default value for a string label is an empty string, for an integer label its 0, and for a boolean label its `false`. Note that there are upper bounds on the maximum number of labels and the number of active time series that are allowed in a project.",
							},
							"metricDescriptor": &dcl.Property{
								Type:        "object",
								GoName:      "MetricDescriptor",
								GoType:      "LogMetricMetricDescriptor",
								Description: "Optional. The metric descriptor associated with the logs-based metric. If unspecified, it uses a default metric descriptor with a DELTA metric kind, INT64 value type, with no labels and a unit of \"1\". Such a metric counts the number of log entries matching the `filter` expression. The `name`, `type`, and `description` fields in the `metric_descriptor` are output only, and is constructed using the `name` and `description` field in the LogMetric. To create a logs-based metric that records a distribution of log values, a DELTA metric kind with a DISTRIBUTION value type must be used along with a `value_extractor` expression in the LogMetric. Each label in the metric descriptor must have a matching label name as the key and an extractor expression as the value in the `label_extractors` map. The `metric_kind` and `value_type` fields in the `metric_descriptor` cannot be updated once initially configured. New labels can be added in the `metric_descriptor`, but existing labels cannot be modified except for their description.",
								Properties: map[string]*dcl.Property{
									"description": &dcl.Property{
										Type:        "string",
										GoName:      "Description",
										ReadOnly:    true,
										Description: "A detailed description of the metric, which can be used in documentation.",
									},
									"displayName": &dcl.Property{
										Type:        "string",
										GoName:      "DisplayName",
										Description: "A concise name for the metric, which can be displayed in user interfaces. Use sentence case without an ending period, for example \"Request count\". This field is optional but it is recommended to be set for any metrics associated with user-visible concepts, such as Quota.",
									},
									"labels": &dcl.Property{
										Type:        "array",
										GoName:      "Labels",
										Description: "The set of labels that can be used to describe a specific instance of this metric type. For example, the `appengine.googleapis.com/http/server/response_latencies` metric type has a label for the HTTP response code, `response_code`, so you can look at latencies for successful responses or just for responses that failed.",
										SendEmpty:   true,
										ListType:    "set",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "LogMetricMetricDescriptorLabels",
											Properties: map[string]*dcl.Property{
												"description": &dcl.Property{
													Type:        "string",
													GoName:      "Description",
													Description: "A human-readable description for the label.",
													Immutable:   true,
												},
												"key": &dcl.Property{
													Type:        "string",
													GoName:      "Key",
													Description: "The label key.",
													Immutable:   true,
												},
												"valueType": &dcl.Property{
													Type:        "string",
													GoName:      "ValueType",
													GoType:      "LogMetricMetricDescriptorLabelsValueTypeEnum",
													Description: "The type of data that can be assigned to the label. Possible values: STRING, BOOL, INT64, DOUBLE, DISTRIBUTION, MONEY",
													Immutable:   true,
													Enum: []string{
														"STRING",
														"BOOL",
														"INT64",
														"DOUBLE",
														"DISTRIBUTION",
														"MONEY",
													},
												},
											},
										},
									},
									"launchStage": &dcl.Property{
										Type:        "string",
										GoName:      "LaunchStage",
										GoType:      "LogMetricMetricDescriptorLaunchStageEnum",
										Description: "Optional. The launch stage of the metric definition. Possible values: UNIMPLEMENTED, PRELAUNCH, EARLY_ACCESS, ALPHA, BETA, GA, DEPRECATED",
										Enum: []string{
											"UNIMPLEMENTED",
											"PRELAUNCH",
											"EARLY_ACCESS",
											"ALPHA",
											"BETA",
											"GA",
											"DEPRECATED",
										},
										Unreadable: true,
									},
									"metadata": &dcl.Property{
										Type:        "object",
										GoName:      "Metadata",
										GoType:      "LogMetricMetricDescriptorMetadata",
										Description: "Optional. Metadata which can be used to guide usage of the metric.",
										Unreadable:  true,
										Properties: map[string]*dcl.Property{
											"ingestDelay": &dcl.Property{
												Type:        "string",
												GoName:      "IngestDelay",
												Description: "The delay of data points caused by ingestion. Data points older than this age are guaranteed to be ingested and available to be read, excluding data loss due to errors.",
											},
											"samplePeriod": &dcl.Property{
												Type:        "string",
												GoName:      "SamplePeriod",
												Description: "The sampling period of metric data points. For metrics which are written periodically, consecutive data points are stored at this time interval, excluding data loss due to errors. Metrics with a higher granularity have a smaller sampling period.",
											},
										},
									},
									"metricKind": &dcl.Property{
										Type:        "string",
										GoName:      "MetricKind",
										GoType:      "LogMetricMetricDescriptorMetricKindEnum",
										Description: "Whether the metric records instantaneous values, changes to a value, etc. Some combinations of `metric_kind` and `value_type` might not be supported. Possible values: GAUGE, DELTA, CUMULATIVE",
										Immutable:   true,
										Enum: []string{
											"GAUGE",
											"DELTA",
											"CUMULATIVE",
										},
									},
									"monitoredResourceTypes": &dcl.Property{
										Type:        "array",
										GoName:      "MonitoredResourceTypes",
										ReadOnly:    true,
										Description: "Read-only. If present, then a time series, which is identified partially by a metric type and a MonitoredResourceDescriptor, that is associated with this metric type can only be associated with one of the monitored resource types listed here.",
										Immutable:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"name": &dcl.Property{
										Type:        "string",
										GoName:      "Name",
										ReadOnly:    true,
										Description: "The resource name of the metric descriptor.",
										Immutable:   true,
									},
									"type": &dcl.Property{
										Type:        "string",
										GoName:      "Type",
										ReadOnly:    true,
										Description: "The metric type, including its DNS name prefix. The type is not URL-encoded. All user-defined metric types have the DNS name `custom.googleapis.com` or `external.googleapis.com`. Metric types should use a natural hierarchical grouping. For example: \"custom.googleapis.com/invoice/paid/amount\" \"external.googleapis.com/prometheus/up\" \"appengine.googleapis.com/http/server/response_latencies\"",
										Immutable:   true,
									},
									"unit": &dcl.Property{
										Type:          "string",
										GoName:        "Unit",
										Description:   "The units in which the metric value is reported. It is only applicable if the `value_type` is `INT64`, `DOUBLE`, or `DISTRIBUTION`. The `unit` defines the representation of the stored metric values. Different systems might scale the values to be more easily displayed (so a value of `0.02kBy` _might_ be displayed as `20By`, and a value of `3523kBy` _might_ be displayed as `3.5MBy`). However, if the `unit` is `kBy`, then the value of the metric is always in thousands of bytes, no matter how it might be displayed. If you want a custom metric to record the exact number of CPU-seconds used by a job, you can create an `INT64 CUMULATIVE` metric whose `unit` is `s{CPU}` (or equivalently `1s{CPU}` or just `s`). If the job uses 12,005 CPU-seconds, then the value is written as `12005`. Alternatively, if you want a custom metric to record data in a more granular way, you can create a `DOUBLE CUMULATIVE` metric whose `unit` is `ks{CPU}`, and then write the value `12.005` (which is `12005/1000`), or use `Kis{CPU}` and write `11.723` (which is `12005/1024`). The supported units are a subset of [The Unified Code for Units of Measure](https://unitsofmeasure.org/ucum.html) standard: **Basic units (UNIT)** * `bit` bit * `By` byte * `s` second * `min` minute * `h` hour * `d` day * `1` dimensionless **Prefixes (PREFIX)** * `k` kilo (10^3) * `M` mega (10^6) * `G` giga (10^9) * `T` tera (10^12) * `P` peta (10^15) * `E` exa (10^18) * `Z` zetta (10^21) * `Y` yotta (10^24) * `m` milli (10^-3) * `u` micro (10^-6) * `n` nano (10^-9) * `p` pico (10^-12) * `f` femto (10^-15) * `a` atto (10^-18) * `z` zepto (10^-21) * `y` yocto (10^-24) * `Ki` kibi (2^10) * `Mi` mebi (2^20) * `Gi` gibi (2^30) * `Ti` tebi (2^40) * `Pi` pebi (2^50) **Grammar** The grammar also includes these connectors: * `/` division or ratio (as an infix operator). For examples, `kBy/{email}` or `MiBy/10ms` (although you should almost never have `/s` in a metric `unit`; rates should always be computed at query time from the underlying cumulative or delta value). * `.` multiplication or composition (as an infix operator). For examples, `GBy.d` or `k{watt}.h`. The grammar for a unit is as follows: Expression = Component: { \".\" Component } { \"/\" Component } ; Component = ( [ PREFIX ] UNIT | \"%\" ) [ Annotation ] | Annotation | \"1\" ; Annotation = \"{\" NAME \"}\" ; Notes: * `Annotation` is just a comment if it follows a `UNIT`. If the annotation is used alone, then the unit is equivalent to `1`. For examples, `{request}/s == 1/s`, `By{transmitted}/s == By/s`. * `NAME` is a sequence of non-blank printable ASCII characters not containing `{` or `}`. * `1` represents a unitary [dimensionless unit](https://en.wikipedia.org/wiki/Dimensionless_quantity) of 1, such as in `1/s`. It is typically used when none of the basic units are appropriate. For example, \"new users per day\" can be represented as `1/d` or `{new-users}/d` (and a metric value `5` would mean \"5 new users). Alternatively, \"thousands of page views per day\" would be represented as `1000/d` or `k1/d` or `k{page_views}/d` (and a metric value of `5.3` would mean \"5300 page views per day\"). * `%` represents dimensionless value of 1/100, and annotates values giving a percentage (so the metric values are typically in the range of 0..100, and a metric value `3` means \"3 percent\"). * `10^2.%` indicates a metric contains a ratio, typically in the range 0..1, that will be multiplied by 100 and displayed as a percentage (so a metric value `0.03` means \"3 percent\").",
										ServerDefault: true,
									},
									"valueType": &dcl.Property{
										Type:        "string",
										GoName:      "ValueType",
										GoType:      "LogMetricMetricDescriptorValueTypeEnum",
										Description: "Whether the measurement is an integer, a floating-point number, etc. Some combinations of `metric_kind` and `value_type` might not be supported. Possible values: STRING, BOOL, INT64, DOUBLE, DISTRIBUTION, MONEY",
										Immutable:   true,
										Enum: []string{
											"STRING",
											"BOOL",
											"INT64",
											"DOUBLE",
											"DISTRIBUTION",
											"MONEY",
										},
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Required. The client-assigned metric identifier. Examples: `\"error_count\"`, `\"nginx/requests\"`. Metric identifiers are limited to 100 characters and can include only the following characters: `A-Z`, `a-z`, `0-9`, and the special characters `_-.,+!*',()%/`. The forward-slash character (`/`) denotes a hierarchy of name pieces, and it cannot be the first character of the name. The metric identifier in this field must not be [URL-encoded](https://en.wikipedia.org/wiki/Percent-encoding). However, when the metric identifier appears as the `[METRIC_ID]` part of a `metric_name` API parameter, then the metric identifier must be URL-encoded. Example: `\"projects/my-project/metrics/nginx%2Frequests\"`.",
								Immutable:   true,
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The resource name of the project in which to create the metric.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The last update timestamp of the metric. This field may not be present for older metrics.",
								Immutable:   true,
							},
							"valueExtractor": &dcl.Property{
								Type:        "string",
								GoName:      "ValueExtractor",
								Description: "Optional. A `value_extractor` is required when using a distribution logs-based metric to extract the values to record from a log entry. Two functions are supported for value extraction: `EXTRACT(field)` or `REGEXP_EXTRACT(field, regex)`. The argument are: 1. field: The name of the log entry field from which the value is to be extracted. 2. regex: A regular expression using the Google RE2 syntax (https://github.com/google/re2/wiki/Syntax) with a single capture group to extract data from the specified log entry field. The value of the field is converted to a string before applying the regex. It is an error to specify a regex that does not include exactly one capture group. The result of the extraction must be convertible to a double type, as the distribution always records double values. If either the extraction or the conversion to double fails, then those values are not recorded in the distribution. Example: `REGEXP_EXTRACT(jsonPayload.request, \".*quantity=(d+).*\")`",
							},
						},
					},
				},
			},
		},
	}
}
