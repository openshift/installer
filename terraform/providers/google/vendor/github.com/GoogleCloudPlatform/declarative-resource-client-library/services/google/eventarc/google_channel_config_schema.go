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
package eventarc

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLGoogleChannelConfigSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Eventarc/GoogleChannelConfig",
			Description: "The Eventarc GoogleChannelConfig resource",
			StructName:  "GoogleChannelConfig",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a GoogleChannelConfig",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "googleChannelConfig",
						Required:    true,
						Description: "A full instance of a GoogleChannelConfig",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a GoogleChannelConfig",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "googleChannelConfig",
						Required:    true,
						Description: "A full instance of a GoogleChannelConfig",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a GoogleChannelConfig",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "googleChannelConfig",
						Required:    true,
						Description: "A full instance of a GoogleChannelConfig",
					},
				},
			},
		},
		Components: &dcl.Components{
			Schemas: map[string]*dcl.Component{
				"GoogleChannelConfig": &dcl.Component{
					Title:           "GoogleChannelConfig",
					ID:              "projects/{{project}}/locations/{{location}}/googleChannelConfig",
					ParentContainer: "project",
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"cryptoKeyName": &dcl.Property{
								Type:        "string",
								GoName:      "CryptoKeyName",
								Description: "Optional. Resource name of a KMS crypto key (managed by the user) used to encrypt/decrypt their event data. It must match the pattern `projects/*/locations/*/keyRings/*/cryptoKeys/*`.",
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudkms/CryptoKey",
										Field:    "name",
									},
								},
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
								Description: "Required. The resource name of the config. Must be in the format of, `projects/{project}/locations/{location}/googleChannelConfig`.",
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
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The last-modified time.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
