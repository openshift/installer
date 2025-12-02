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

func DCLServiceSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Monitoring/Service",
			Description: "The Monitoring Service resource",
			StructName:  "Service",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Service",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "service",
						Required:    true,
						Description: "A full instance of a Service",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Service",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "service",
						Required:    true,
						Description: "A full instance of a Service",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Service",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "service",
						Required:    true,
						Description: "A full instance of a Service",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Service",
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
				Description: "The function used to list information about many Service",
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
				"Service": &dcl.Component{
					Title:           "Service",
					ID:              "projects/{{project}}/services/{{name}}",
					ParentContainer: "project",
					LabelsField:     "userLabels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "Name used for UI elements listing this Service.",
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Resource name for this Service. The format is: projects/[PROJECT_ID_OR_NUMBER]/services/[SERVICE_ID]",
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
							"telemetry": &dcl.Property{
								Type:        "object",
								GoName:      "Telemetry",
								GoType:      "ServiceTelemetry",
								Description: "Configuration for how to query telemetry on a Service.",
								Properties: map[string]*dcl.Property{
									"resourceName": &dcl.Property{
										Type:        "string",
										GoName:      "ResourceName",
										Description: "The full name of the resource that defines this service. Formatted as described in https://cloud.google.com/apis/design/resource_names.",
									},
								},
							},
							"userLabels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "UserLabels",
								Description: "Labels which have been used to annotate the service. Label keys must start with a letter. Label keys and values may contain lowercase letters, numbers, underscores, and dashes. Label keys and values have a maximum length of 63 characters, and must be less than 128 bytes in size. Up to 64 label entries may be stored. For labels which do not have a semantic value, the empty string may be supplied for the label value.",
							},
						},
					},
				},
			},
		},
	}
}
