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

func DCLNotificationChannelSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Monitoring/NotificationChannel",
			Description: "The Monitoring NotificationChannel resource",
			StructName:  "NotificationChannel",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a NotificationChannel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "notificationChannel",
						Required:    true,
						Description: "A full instance of a NotificationChannel",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a NotificationChannel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "notificationChannel",
						Required:    true,
						Description: "A full instance of a NotificationChannel",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a NotificationChannel",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "notificationChannel",
						Required:    true,
						Description: "A full instance of a NotificationChannel",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all NotificationChannel",
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
				Description: "The function used to list information about many NotificationChannel",
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
				"NotificationChannel": &dcl.Component{
					Title:           "NotificationChannel",
					ID:              "projects/{{project}}/notificationChannels/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Properties: map[string]*dcl.Property{
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "An optional human-readable description of this notification channel. This description may provide additional details, beyond the display name, for the channel. This may not exceed 1024 Unicode characters.",
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "An optional human-readable name for this notification channel. It is recommended that you specify a non-empty and unique name in order to make it easier to identify the channels in your project, though this is not enforced. The display name is limited to 512 Unicode characters.",
							},
							"enabled": &dcl.Property{
								Type:        "boolean",
								GoName:      "Enabled",
								Description: "Whether notifications are forwarded to the described channel. This makes it possible to disable delivery of notifications to a particular channel without removing the channel from all alerting policies that reference the channel. This is a more convenient approach when the change is temporary and you want to receive notifications from the same set of alerting policies on the channel at some point in the future.",
								Default:     true,
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "Configuration fields that define the channel and its behavior. The permissible and required labels are specified in the [NotificationChannelDescriptor.labels][google.monitoring.v3.NotificationChannelDescriptor.labels] of the `NotificationChannelDescriptor` corresponding to the `type` field.",
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "The full REST resource name for this channel. The format is: projects/[PROJECT_ID_OR_NUMBER]/notificationChannels/[CHANNEL_ID] The `[CHANNEL_ID]` is automatically assigned by the server on creation.",
								Immutable:                true,
								ServerGeneratedParameter: true,
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project for this notification channel.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"type": &dcl.Property{
								Type:        "string",
								GoName:      "Type",
								Description: "The type of the notification channel. This field matches the value of the [NotificationChannelDescriptor.type][google.monitoring.v3.NotificationChannelDescriptor.type] field.",
							},
							"userLabels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "UserLabels",
								Description: "User-supplied key/value data that does not need to conform to the corresponding `NotificationChannelDescriptor`'s schema, unlike the `labels` field. This field is intended to be used for orv3nizing and identifying the `NotificationChannel` objects. The field can contain up to 64 entries. Each key and value is limited to 63 Unicode characters or 128 bytes, whichever is smaller. Labels and values can contain only lowercase letters, numerals, underscores, and dashes. Keys must begin with a letter.",
							},
							"verificationStatus": &dcl.Property{
								Type:        "string",
								GoName:      "VerificationStatus",
								GoType:      "NotificationChannelVerificationStatusEnum",
								ReadOnly:    true,
								Description: "Indicates whether this channel has been verified or not. On a [`ListNotificationChannels`][google.monitoring.v3.NotificationChannelService.ListNotificationChannels] or [`GetNotificationChannel`][google.monitoring.v3.NotificationChannelService.GetNotificationChannel] operation, this field is expected to be populated. If the value is `UNVERIFIED`, then it indicates that the channel is non-functioning (it both requires verification and lacks verification); otherwise, it is assumed that the channel works. If the channel is neither `VERIFIED` nor `UNVERIFIED`, it implies that the channel is of a type that does not require verification or that this specific channel has been exempted from verification because it was created prior to verification being required for channels of this type. This field cannot be modified using a standard [`UpdateNotificationChannel`][google.monitoring.v3.NotificationChannelService.UpdateNotificationChannel] operation. To change the value of this field, you must call [`VerifyNotificationChannel`][google.monitoring.v3.NotificationChannelService.VerifyNotificationChannel]. Possible values: VERIFICATION_STATUS_UNSPECIFIED, UNVERIFIED, VERIFIED",
								Immutable:   true,
								Enum: []string{
									"VERIFICATION_STATUS_UNSPECIFIED",
									"UNVERIFIED",
									"VERIFIED",
								},
							},
						},
					},
				},
			},
		},
	}
}
