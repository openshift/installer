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
package cloudresourcemanager

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLProjectSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "CloudResourceManager/Project",
			Description: "The CloudResourceManager Project resource",
			StructName:  "Project",
			HasIAM:      true,
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Project",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "project",
						Required:    true,
						Description: "A full instance of a Project",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Project",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "project",
						Required:    true,
						Description: "A full instance of a Project",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Project",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "project",
						Required:    true,
						Description: "A full instance of a Project",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Project",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "parent",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many Project",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "parent",
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
				"Project": &dcl.Component{
					Title:       "Project",
					ID:          "projects/{{name}}",
					LabelsField: "labels",
					HasCreate:   true,
					HasIAM:      true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Properties: map[string]*dcl.Property{
							"displayname": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "The optional user-assigned display name of the Project. When present it must be between 4 to 30 characters. Allowed characters are: lowercase and uppercase letters, numbers, hyphen, single-quote, double-quote, space, and exclamation point. Example: `My Project` Read-write.",
								Immutable:   true,
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "User-specified labels.",
							},
							"lifecycleState": &dcl.Property{
								Type:        "string",
								GoName:      "LifecycleState",
								GoType:      "ProjectLifecycleStateEnum",
								ReadOnly:    true,
								Description: "The Project lifecycle state. Read-only. Possible values: LIFECYCLE_STATE_UNSPECIFIED, ACTIVE, DELETE_REQUESTED, DELETE_IN_PROGRESS",
								Immutable:   true,
								Enum: []string{
									"LIFECYCLE_STATE_UNSPECIFIED",
									"ACTIVE",
									"DELETE_REQUESTED",
									"DELETE_IN_PROGRESS",
								},
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The unique, user-assigned ID of the Project. It must be 6 to 30 lowercase letters, digits, or hyphens. It must start with a letter. Trailing hyphens are prohibited. Example: `tokyo-rain-123` Read-only after creation.",
								Immutable:   true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "An optional reference to a parent Resource. Supported values include organizations/<org_id> and folders/<folder_id>. Once set, the parent cannot be cleared. The `parent` can be set on creation or using the `UpdateProject` method; the end user must have the `resourcemanager.projects.create` permission on the parent. Read-write. ",
								Immutable:           true,
								ForwardSlashAllowed: true,
							},
							"projectNumber": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "ProjectNumber",
								ReadOnly:    true,
								Description: "The number uniquely identifying the project. Example: `415104041262` Read-only. ",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
