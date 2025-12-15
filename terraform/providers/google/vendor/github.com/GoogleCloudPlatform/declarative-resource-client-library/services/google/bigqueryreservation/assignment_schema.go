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
package bigqueryreservation

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLAssignmentSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "BigqueryReservation/Assignment",
			Description: "The BigqueryReservation Assignment resource",
			StructName:  "Assignment",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Assignment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "assignment",
						Required:    true,
						Description: "A full instance of a Assignment",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Assignment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "assignment",
						Required:    true,
						Description: "A full instance of a Assignment",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Assignment",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "assignment",
						Required:    true,
						Description: "A full instance of a Assignment",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Assignment",
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
						Name:     "reservation",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Assignment",
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
						Name:     "reservation",
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
				"Assignment": &dcl.Component{
					Title:           "Assignment",
					ID:              "projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"assignee",
							"jobType",
							"reservation",
						},
						Properties: map[string]*dcl.Property{
							"assignee": &dcl.Property{
								Type:        "string",
								GoName:      "Assignee",
								Description: "The resource which will use the reservation. E.g. projects/myproject, folders/123, organizations/456.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Folder",
										Field:    "name",
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Organization",
										Field:    "name",
									},
								},
							},
							"jobType": &dcl.Property{
								Type:        "string",
								GoName:      "JobType",
								GoType:      "AssignmentJobTypeEnum",
								Description: "Types of job, which could be specified when using the reservation. Possible values: JOB_TYPE_UNSPECIFIED, PIPELINE, QUERY",
								Immutable:   true,
								Enum: []string{
									"JOB_TYPE_UNSPECIFIED",
									"PIPELINE",
									"QUERY",
								},
							},
							"location": &dcl.Property{
								Type:           "string",
								GoName:         "Location",
								Description:    "The location for the resource",
								Immutable:      true,
								ExtractIfEmpty: true,
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "Output only. The resource name of the assignment.",
								Immutable:                true,
								ServerGeneratedParameter: true,
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
								ExtractIfEmpty: true,
							},
							"reservation": &dcl.Property{
								Type:        "string",
								GoName:      "Reservation",
								Description: "The reservation for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Bigqueryreservation/Reservation",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "AssignmentStateEnum",
								ReadOnly:    true,
								Description: "Assignment will remain in PENDING state if no active capacity commitment is present. It will become ACTIVE when some capacity commitment becomes active. Possible values: STATE_UNSPECIFIED, PENDING, ACTIVE",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"PENDING",
									"ACTIVE",
								},
							},
						},
					},
				},
			},
		},
	}
}
