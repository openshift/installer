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

func DCLMonitoredProjectSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Monitoring/MonitoredProject",
			Description: "Monitored Project allows you to set a project as monitored by a _metrics scope_, which is a term for a project used to group the metrics of multiple projects, potentially across multiple organizations.  This enables you to view these groups in the Monitoring page of the cloud console.",
			StructName:  "MonitoredProject",
			Reference: &dcl.Link{
				Text: "REST API",
				URL:  "https://cloud.google.com/monitoring/api/ref_v3/rest/v1/locations.global.metricsScopes",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Understanding metrics scopes",
					URL:  "https://cloud.google.com/monitoring/settings#concept-scope",
				},
				&dcl.Link{
					Text: "API notes",
					URL:  "https://cloud.google.com/monitoring/settings/manage-api",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a MonitoredProject",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "monitoredProject",
						Required:    true,
						Description: "A full instance of a MonitoredProject",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a MonitoredProject",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "monitoredProject",
						Required:    true,
						Description: "A full instance of a MonitoredProject",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a MonitoredProject",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "monitoredProject",
						Required:    true,
						Description: "A full instance of a MonitoredProject",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all MonitoredProject",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "metricsScope",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many MonitoredProject",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "metricsScope",
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
				"MonitoredProject": &dcl.Component{
					Title: "MonitoredProject",
					ID:    "locations/global/metricsScopes/{{metrics_scope}}/projects/{{name}}",
					Locations: []string{
						"global",
					},
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"metricsScope",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time when this `MonitoredProject` was created.",
								Immutable:   true,
							},
							"metricsScope": &dcl.Property{
								Type:        "string",
								GoName:      "MetricsScope",
								Description: "Required. The resource name of the existing Metrics Scope that will monitor this project. Example: locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}",
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
			},
		},
	}
}
