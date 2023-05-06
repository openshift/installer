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

func DCLLogViewSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Logging/LogView",
			Description: "The Logging LogView resource",
			StructName:  "LogView",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a LogView",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logView",
						Required:    true,
						Description: "A full instance of a LogView",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a LogView",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logView",
						Required:    true,
						Description: "A full instance of a LogView",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a LogView",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logView",
						Required:    true,
						Description: "A full instance of a LogView",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all LogView",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "bucket",
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
				Description: "The function used to list information about many LogView",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:     "location",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
					dcl.PathParameters{
						Name:     "bucket",
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
				"LogView": &dcl.Component{
					Title:     "LogView",
					ID:        "{{parent}}/locations/{{location}}/buckets/{{bucket}}/views/{{name}}",
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"bucket",
						},
						Properties: map[string]*dcl.Property{
							"bucket": &dcl.Property{
								Type:        "string",
								GoName:      "Bucket",
								Description: "The bucket of the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Logging/LogBucket",
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
								Description: "Output only. The creation timestamp of the view.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Describes this view.",
							},
							"filter": &dcl.Property{
								Type:        "string",
								GoName:      "Filter",
								Description: "Filter that restricts which log entries in a bucket are visible in this view. Filters are restricted to be a logical AND of ==/!= of any of the following: - originating project/folder/organization/billing account. - resource type - log id For example: SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\" AND LOG_ID(\"stdout\")",
							},
							"location": &dcl.Property{
								Type:           "string",
								GoName:         "Location",
								Description:    "The location of the resource. The supported locations are: global, us-central1, us-east1, us-west1, asia-east1, europe-west1.",
								Immutable:      true,
								ExtractIfEmpty: true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The resource name of the view. For example: `projects/my-project/locations/global/buckets/my-bucket/views/my-view`",
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
								ExtractIfEmpty: true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The last update timestamp of the view.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
