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
package logging

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLLogBucketSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Logging/LogBucket",
			Description: "The Logging LogBucket resource",
			StructName:  "LogBucket",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a LogBucket",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logBucket",
						Required:    true,
						Description: "A full instance of a LogBucket",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a LogBucket",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logBucket",
						Required:    true,
						Description: "A full instance of a LogBucket",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a LogBucket",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logBucket",
						Required:    true,
						Description: "A full instance of a LogBucket",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all LogBucket",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
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
				Description: "The function used to list information about many LogBucket",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
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
				"LogBucket": &dcl.Component{
					Title:     "LogBucket",
					ID:        "{{parent}}/locations/{{location}}/buckets/{{name}}",
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"parent",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The creation timestamp of the bucket. This is not set for any of the default buckets.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Describes this bucket.",
							},
							"lifecycleState": &dcl.Property{
								Type:        "string",
								GoName:      "LifecycleState",
								GoType:      "LogBucketLifecycleStateEnum",
								ReadOnly:    true,
								Description: "Output only. The bucket lifecycle state. Possible values: LIFECYCLE_STATE_UNSPECIFIED, ACTIVE, DELETE_REQUESTED",
								Immutable:   true,
								Enum: []string{
									"LIFECYCLE_STATE_UNSPECIFIED",
									"ACTIVE",
									"DELETE_REQUESTED",
								},
							},
							"location": &dcl.Property{
								Type:        "string",
								GoName:      "Location",
								Description: "The location of the resource. The supported locations are: global, us-central1, us-east1, us-west1, asia-east1, europe-west1.",
								Immutable:   true,
							},
							"locked": &dcl.Property{
								Type:        "boolean",
								GoName:      "Locked",
								Description: "Whether the bucket has been locked. The retention period on a locked bucket may not be changed. Locked buckets may only be deleted if they are empty.",
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The resource name of the bucket. For example: \"projects/my-project-id/locations/my-location/buckets/my-bucket-id\" The supported locations are: `global`, `us-central1`, `us-east1`, `us-west1`, `asia-east1`, `europe-west1`. For the location of `global` it is unspecified where logs are actually stored. Once a bucket has been created, the location can not be changed.",
								Immutable:   true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "The parent of the resource.",
								Immutable:           true,
								ForwardSlashAllowed: true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/BillingAccount",
										Field:    "name",
										Parent:   true,
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Folder",
										Field:    "name",
										Parent:   true,
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Organization",
										Field:    "name",
										Parent:   true,
									},
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"retentionDays": &dcl.Property{
								Type:        "integer",
								Format:      "int64",
								GoName:      "RetentionDays",
								Description: "Logs will be retained by default for this amount of time, after which they will automatically be deleted. The minimum retention period is 1 day. If this value is set to zero at bucket creation time, the default time of 30 days will be used.",
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The last update timestamp of the bucket.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
