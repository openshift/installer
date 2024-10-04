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

func DCLTagValueSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "CloudResourceManager/TagValue",
			Description: "The CloudResourceManager TagValue resource",
			StructName:  "TagValue",
			HasIAM:      true,
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a TagValue",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "tagValue",
						Required:    true,
						Description: "A full instance of a TagValue",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a TagValue",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "tagValue",
						Required:    true,
						Description: "A full instance of a TagValue",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a TagValue",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "tagValue",
						Required:    true,
						Description: "A full instance of a TagValue",
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"TagValue": &dcl.Component{
					Title:     "TagValue",
					ID:        "tagValues/{{name}}",
					HasCreate: true,
					HasIAM:    true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"parent",
							"shortName",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. Creation time.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. User-assigned description of the TagValue. Must not exceed 256 characters. Read-write.",
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "Optional. Entity tag which users can pass to prevent race conditions.",
								Immutable:   true,
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "Immutable. The generated numeric id for the TagValue.",
								Immutable:                true,
								ServerGeneratedParameter: true,
								HasLongForm:              true,
							},
							"namespacedName": &dcl.Property{
								Type:        "string",
								GoName:      "NamespacedName",
								ReadOnly:    true,
								Description: "Output only. Immutable. Namespaced name of the TagValue.",
								Immutable:   true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "Immutable. The resource name of the new TagValue's parent. Must be of the form `tagKeys/{tag_key_id}`.",
								Immutable:           true,
								ForwardSlashAllowed: true,
							},
							"shortName": &dcl.Property{
								Type:        "string",
								GoName:      "ShortName",
								Description: "Required. Immutable. The user friendly name for a TagValue. The short name should be unique for TagValuess within the same parent TagKey. The short name must be 1-63 characters, beginning and ending with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Update time.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
