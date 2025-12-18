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

func DCLGroupSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Monitoring/Group",
			Description: "The Monitoring Group resource",
			StructName:  "Group",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Group",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "group",
						Required:    true,
						Description: "A full instance of a Group",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Group",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "group",
						Required:    true,
						Description: "A full instance of a Group",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Group",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "group",
						Required:    true,
						Description: "A full instance of a Group",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Group",
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
				Description: "The function used to list information about many Group",
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
				"Group": &dcl.Component{
					Title:           "Group",
					ID:              "projects/{{project}}/groups/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"displayName",
							"filter",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "A user-assigned name for this group, used only for display purposes.",
							},
							"filter": &dcl.Property{
								Type:        "string",
								GoName:      "Filter",
								Description: "The filter used to determine which monitored resources belong to this group.",
							},
							"isCluster": &dcl.Property{
								Type:        "boolean",
								GoName:      "IsCluster",
								Description: "If true, the members of this group are considered to be a cluster. The system can perform additional analysis on groups that are clusters.",
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "Output only. The name of this group. The format is: `projects/{{project}}/groups/{{name}}`, which is generated automatically.",
								Immutable:                true,
								ServerGeneratedParameter: true,
							},
							"parentName": &dcl.Property{
								Type:        "string",
								GoName:      "ParentName",
								Description: "The name of the group's parent, if it has one. The format is: projects/ For groups with no parent, `parent_name` is the empty string, ``.",
								SendEmpty:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Monitoring/Group",
										Field:    "name",
									},
								},
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project of the group",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
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
