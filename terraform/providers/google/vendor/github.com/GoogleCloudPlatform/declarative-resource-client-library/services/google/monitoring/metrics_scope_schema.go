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

func DCLMetricsScopeSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Monitoring/MetricsScope",
			Description: "The Monitoring MetricsScope resource",
			StructName:  "MetricsScope",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a MetricsScope",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "metricsScope",
						Required:    true,
						Description: "A full instance of a MetricsScope",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a MetricsScope",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "metricsScope",
						Required:    true,
						Description: "A full instance of a MetricsScope",
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"MetricsScope": &dcl.Component{
					Title: "MetricsScope",
					ID:    "locations/global/metricsScopes/{{name}}",
					Locations: []string{
						"global",
					},
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time when this `Metrics Scope` was created.",
								Immutable:   true,
							},
							"monitoredProjects": &dcl.Property{
								Type:        "array",
								GoName:      "MonitoredProjects",
								ReadOnly:    true,
								Description: "Output only. The list of projects monitored by this `Metrics Scope`.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "MetricsScopeMonitoredProjects",
									Properties: map[string]*dcl.Property{
										"createTime": &dcl.Property{
											Type:        "string",
											Format:      "date-time",
											GoName:      "CreateTime",
											ReadOnly:    true,
											Description: "Output only. The time when this `MonitoredProject` was created.",
											Immutable:   true,
										},
										"name": &dcl.Property{
											Type:        "string",
											GoName:      "Name",
											Description: "Immutable. The resource name of the `MonitoredProject`. On input, the resource name includes the scoping project ID and monitored project ID. On output, it contains the equivalent project numbers. Example: `locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}/projects/{MONITORED_PROJECT_ID_OR_NUMBER}`",
											Immutable:   true,
										},
									},
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Immutable. The resource name of the Monitoring Metrics Scope. On input, the resource name can be specified with the scoping project ID or number. On output, the resource name is specified with the scoping project number. Example: `locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}`",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time when this `Metrics Scope` record was last updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
