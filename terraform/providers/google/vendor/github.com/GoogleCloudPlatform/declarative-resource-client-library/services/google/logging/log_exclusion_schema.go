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

func DCLLogExclusionSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Logging/LogExclusion",
			Description: "The Logging LogExclusion resource",
			StructName:  "LogExclusion",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a LogExclusion",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logExclusion",
						Required:    true,
						Description: "A full instance of a LogExclusion",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a LogExclusion",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logExclusion",
						Required:    true,
						Description: "A full instance of a LogExclusion",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a LogExclusion",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "logExclusion",
						Required:    true,
						Description: "A full instance of a LogExclusion",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all LogExclusion",
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
				Description: "The function used to list information about many LogExclusion",
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
				"LogExclusion": &dcl.Component{
					Title:     "LogExclusion",
					ID:        "{{parent}}/exclusions/{{name}}",
					HasCreate: true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"filter",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The creation timestamp of the exclusion. This field may not be present for older exclusions.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. A description of this exclusion.",
							},
							"disabled": &dcl.Property{
								Type:        "boolean",
								GoName:      "Disabled",
								Description: "Optional. If set to True, then this exclusion is disabled and it does not exclude any log entries. You can update an exclusion to change the value of this field.",
							},
							"filter": &dcl.Property{
								Type:        "string",
								GoName:      "Filter",
								Description: "Required. An (https://cloud.google.com/logging/docs/view/advanced-queries#sample), you can exclude less than 100% of the matching log entries. For example, the following query matches 99% of low-severity log entries from Google Cloud Storage buckets: `\"resource.type=gcs_bucket severity",
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "Required. A client-assigned identifier, such as `\"load-balancer-exclusion\"`. Identifiers are limited to 100 characters and can include only letters, digits, underscores, hyphens, and periods. First character has to be alphanumeric.",
								Immutable:   true,
							},
							"parent": &dcl.Property{
								Type:                "string",
								GoName:              "Parent",
								Description:         "The parent resource in which to create the exclusion: \"projects/[PROJECT_ID]\" \"organizations/[ORGANIZATION_ID]\" \"billingAccounts/[BILLING_ACCOUNT_ID]\" \"folders/[FOLDER_ID]\" Examples: \"projects/my-logging-project\", \"organizations/123456789\". Authorization requires the following IAM permission on the specified resource parent: logging.exclusions.create",
								Immutable:           true,
								ForwardSlashAllowed: true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
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
										Resource: "Cloudresourcemanager/BillingAccount",
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
								Description: "Output only. The last update timestamp of the exclusion. This field may not be present for older exclusions.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
