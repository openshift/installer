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
package containerazure

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLAzureClientSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "ContainerAzure/Client",
			Description: "AzureClient resources hold client authentication information needed by the Anthos Multi-Cloud API to manage Azure resources on your Azure subscription.When an AzureCluster is created, an AzureClient resource needs to be provided and all operations on Azure resources associated to that cluster will authenticate to Azure services using the given client.AzureClient resources are immutable and cannot be modified upon creation.Each AzureClient resource is bound to a single Azure Active Directory Application and tenant.",
			StructName:  "AzureClient",
			Reference: &dcl.Link{
				Text: "API reference",
				URL:  "https://cloud.google.com/anthos/clusters/docs/multi-cloud/reference/rest/v1/projects.locations.azureClients",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Multicloud overview",
					URL:  "https://cloud.google.com/anthos/clusters/docs/multi-cloud",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Client",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "client",
						Required:    true,
						Description: "A full instance of a Client",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Client",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "client",
						Required:    true,
						Description: "A full instance of a Client",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Client",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "client",
						Required:    true,
						Description: "A full instance of a Client",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Client",
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
				Description: "The function used to list information about many Client",
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
				"Client": &dcl.Component{
					Title:           "AzureClient",
					ID:              "projects/{{project}}/locations/{{location}}/azureClients/{{name}}",
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"tenantId",
							"applicationId",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"applicationId": &dcl.Property{
								Type:        "string",
								GoName:      "ApplicationId",
								Description: "The Azure Active Directory Application ID.",
								Immutable:   true,
							},
							"certificate": &dcl.Property{
								Type:        "string",
								GoName:      "Certificate",
								ReadOnly:    true,
								Description: "Output only. The PEM encoded x509 certificate.",
								Immutable:   true,
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this resource was created.",
								Immutable:   true,
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
								Description: "The name of this resource.",
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
							"tenantId": &dcl.Property{
								Type:        "string",
								GoName:      "TenantId",
								Description: "The Azure Active Directory Tenant ID.",
								Immutable:   true,
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. A globally unique identifier for the client.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
