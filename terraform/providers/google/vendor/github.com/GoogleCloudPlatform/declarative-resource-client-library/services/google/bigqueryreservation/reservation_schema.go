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

func DCLReservationSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "BigqueryReservation/Reservation",
			Description: "The BigqueryReservation Reservation resource",
			StructName:  "Reservation",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Reservation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "reservation",
						Required:    true,
						Description: "A full instance of a Reservation",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Reservation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "reservation",
						Required:    true,
						Description: "A full instance of a Reservation",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Reservation",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "reservation",
						Required:    true,
						Description: "A full instance of a Reservation",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Reservation",
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
				Description: "The function used to list information about many Reservation",
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
				"Reservation": &dcl.Component{
					Title:           "Reservation",
					ID:              "projects/{{project}}/locations/{{location}}/reservations/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"creationTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreationTime",
								ReadOnly:    true,
								Description: "Output only. Creation time of the reservation.",
								Immutable:   true,
							},
							"ignoreIdleSlots": &dcl.Property{
								Type:        "boolean",
								GoName:      "IgnoreIdleSlots",
								Description: "If false, any query using this reservation will use idle slots from other reservations within the same admin project. If true, a query using this reservation will execute with the slot capacity specified above at most.",
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The resource name of the reservation.",
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
							"slotCapacity": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "SlotCapacity",
								Description: "Minimum slots available to this reservation. A slot is a unit of computational power in BigQuery, and serves as the unit of parallelism. Queries using this reservation might use more slots during runtime if ignore_idle_slots is set to false. If the new reservation's slot capacity exceed the parent's slot capacity or if total slot capacity of the new reservation and its siblings exceeds the parent's slot capacity, the request will fail with `google.rpc.Code.RESOURCE_EXHAUSTED`.",
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Last update time of the reservation.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
