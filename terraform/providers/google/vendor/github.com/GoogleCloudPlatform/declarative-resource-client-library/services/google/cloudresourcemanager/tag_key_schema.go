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
package cloudresourcemanager

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLTagKeySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "CloudResourceManager/TagKey",
			Description: "The CloudResourceManager TagKey resource",
			StructName:  "TagKey",
			HasIAM:      true,
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a TagKey",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "tagKey",
						Required:    true,
						Description: "A full instance of a TagKey",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a TagKey",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "tagKey",
						Required:    true,
						Description: "A full instance of a TagKey",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a TagKey",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "tagKey",
						Required:    true,
						Description: "A full instance of a TagKey",
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"TagKey": &dcl.Component{
					Title:     "TagKey",
					ID:        "tagKeys/{{name}}",
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
								Description: "Optional. User-assigned description of the TagKey. Must not exceed 256 characters. Read-write.",
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
								Description:              "Immutable. The generated numeric id for the TagKey.",
								Immutable:                true,
								ServerGeneratedParameter: true,
							},
							"namespacedName": &dcl.Property{
								Type:        "string",
								GoName:      "NamespacedName",
								ReadOnly:    true,
								Description: "Output only. Immutable. Namespaced name of the TagKey.",
								Immutable:   true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "Immutable. The resource name of the new TagKey's parent. Must be of the form `organizations/{org_id}`.",
								Immutable:           true,
								ForwardSlashAllowed: true,
							},
							"purpose": &dcl.Property{
								Type:        "string",
								GoName:      "Purpose",
								GoType:      "TagKeyPurposeEnum",
								Description: "Optional. A purpose denotes that this Tag is intended for use in policies of a specific policy engine, and will involve that policy engine in management operations involving this Tag. A purpose does not grant a policy engine exclusive rights to the Tag, and it may be referenced by other policy engines. A purpose cannot be changed once set. Possible values: PURPOSE_UNSPECIFIED, GCE_FIREWALL",
								Immutable:   true,
								Enum: []string{
									"PURPOSE_UNSPECIFIED",
									"GCE_FIREWALL",
								},
							},
							"purposeData": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "PurposeData",
								Description: "Optional. Purpose data corresponds to the policy system that the tag is intended for. See documentation for `Purpose` for formatting of this field. Purpose data cannot be changed once set.",
								Immutable:   true,
							},
							"shortName": &dcl.Property{
								Type:        "string",
								GoName:      "ShortName",
								Description: "Required. Immutable. The user friendly name for a TagKey. The short name should be unique for TagKeys within the same tag namespace. The short name must be 1-63 characters, beginning and ending with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.",
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
