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
package cloudbuildv2

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLRepositorySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Cloudbuildv2/Repository",
			Description: "The Cloudbuildv2 Repository resource",
			StructName:  "Repository",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Repository",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "repository",
						Required:    true,
						Description: "A full instance of a Repository",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Repository",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "repository",
						Required:    true,
						Description: "A full instance of a Repository",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Repository",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "repository",
						Required:    true,
						Description: "A full instance of a Repository",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Repository",
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
						Name:     "connection",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Repository",
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
						Name:     "connection",
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
				"Repository": &dcl.Component{
					Title:           "Repository",
					ID:              "projects/{{project}}/locations/{{location}}/connections/{{connection}}/repositories/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"remoteUri",
							"project",
							"location",
							"connection",
						},
						Properties: map[string]*dcl.Property{
							"annotations": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Annotations",
								Description: "Allows clients to store small amounts of arbitrary data.",
								Immutable:   true,
							},
							"connection": &dcl.Property{
								Type:        "string",
								GoName:      "Connection",
								Description: "The connection for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudbuildv2/Connection",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. Server assigned timestamp for when the connection was created.",
								Immutable:   true,
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"location": &dcl.Property{
								Type:           "string",
								GoName:         "Location",
								Description:    "The location for the resource",
								Immutable:      true,
								ExtractIfEmpty: true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Name of the repository.",
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
								ExtractIfEmpty: true,
							},
							"remoteUri": &dcl.Property{
								Type:        "string",
								GoName:      "RemoteUri",
								Description: "Required. Git Clone HTTPS URI.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Server assigned timestamp for when the connection was updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
