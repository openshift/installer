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

func DCLFolderSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "CloudResourceManager/Folder",
			Description: "The CloudResourceManager Folder resource",
			StructName:  "Folder",
			HasIAM:      true,
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Folder",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "folder",
						Required:    true,
						Description: "A full instance of a Folder",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Folder",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "folder",
						Required:    true,
						Description: "A full instance of a Folder",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Folder",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "folder",
						Required:    true,
						Description: "A full instance of a Folder",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Folder",
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
				Description: "The function used to list information about many Folder",
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
				"Folder": &dcl.Component{
					Title:     "Folder",
					ID:        "folders/{{name}}",
					HasCreate: true,
					HasIAM:    true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"parent",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. Timestamp when the Folder was created. Assigned by the server.",
								Immutable:   true,
							},
							"deleteTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "DeleteTime",
								ReadOnly:    true,
								Description: "Output only. Timestamp when the Folder was requested to be deleted.",
								Immutable:   true,
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "The folder's display name. A folder's display name must be unique amongst its siblings, e.g. no two folders with the same parent can share the same display name. The display name must start and end with a letter or digit, may contain letters, digits, spaces, hyphens and underscores and can be no longer than 30 characters. This is captured by the regular expression: `[p{L}p{N}]([p{L}p{N}_- ]{0,28}[p{L}p{N}])?`.",
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Output only. A checksum computed by the server based on the current value of the Folder resource. This may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "Output only. The resource name of the Folder.",
								Immutable:                true,
								ServerGeneratedParameter: true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "Required. The Folder's parent's resource name. Updates to the folder's parent must be performed via MoveFolder.",
								ForwardSlashAllowed: true,
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "FolderStateEnum",
								ReadOnly:    true,
								Description: "Output only. The lifecycle state of the folder. Possible values: LIFECYCLE_STATE_UNSPECIFIED, ACTIVE, DELETE_REQUESTED",
								Immutable:   true,
								Enum: []string{
									"LIFECYCLE_STATE_UNSPECIFIED",
									"ACTIVE",
									"DELETE_REQUESTED",
								},
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Timestamp when the Folder was last modified.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
