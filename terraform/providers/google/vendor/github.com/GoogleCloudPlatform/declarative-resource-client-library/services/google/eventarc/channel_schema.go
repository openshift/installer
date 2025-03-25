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
package eventarc

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLChannelSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Eventarc/Channel",
			Description: "The Eventarc Channel resource",
			StructName:  "Channel",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Channel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "channel",
						Required:    true,
						Description: "A full instance of a Channel",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Channel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "channel",
						Required:    true,
						Description: "A full instance of a Channel",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Channel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "channel",
						Required:    true,
						Description: "A full instance of a Channel",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Channel",
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
				Description: "The function used to list information about many Channel",
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
				"Channel": &dcl.Component{
					Title:           "Channel",
					ID:              "projects/{{project}}/locations/{{location}}/channels/{{name}}",
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
							"activationToken": &dcl.Property{
								Type:        "string",
								GoName:      "ActivationToken",
								ReadOnly:    true,
								Description: "Output only. The activation token for the channel. The token must be used by the provider to register the channel for publishing.",
								Immutable:   true,
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The creation time.",
								Immutable:   true,
							},
							"cryptoKeyName": &dcl.Property{
								Type:        "string",
								GoName:      "CryptoKeyName",
								Description: "Optional. Resource name of a KMS crypto key (managed by the user) used to encrypt/decrypt their event data. It must match the pattern `projects/*/locations/*/keyRings/*/cryptoKeys/*`.",
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudkms/CryptoKey",
										Field:    "selfLink",
									},
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location for the resource",
								Immutable:   true,
								Parameter:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Required. The resource name of the channel. Must be unique within the location on the project.",
								Immutable:   true,
								HasLongForm: true,
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
								Parameter: true,
							},
							"pubsubTopic": &dcl.Property{
								Type:        "string",
								GoName:      "PubsubTopic",
								ReadOnly:    true,
								Description: "Output only. The name of the Pub/Sub topic created and managed by Eventarc system as a transport for the event delivery. Format: `projects/{project}/topics/{topic_id}`.",
								Immutable:   true,
							},
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "ChannelStateEnum",
								ReadOnly:    true,
								Description: "Output only. The state of a Channel. Possible values: STATE_UNSPECIFIED, PENDING, ACTIVE, INACTIVE",
								Immutable:   true,
								Enum: []string{
									"STATE_UNSPECIFIED",
									"PENDING",
									"ACTIVE",
									"INACTIVE",
								},
							},
							"thirdPartyProvider": &dcl.Property{
								Type:        "string",
								GoName:      "ThirdPartyProvider",
								Description: "The name of the event provider (e.g. Eventarc SaaS partner) associated with the channel. This provider will be granted permissions to publish events to the channel. Format: `projects/{project}/locations/{location}/providers/{provider_id}`.",
								Immutable:   true,
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. Server assigned unique identifier for the channel. The value is a UUID4 string and guaranteed to remain unchanged until the resource is deleted.",
								Immutable:   true,
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
