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
package cloudbuildv2

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLConnectionSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Cloudbuildv2/Connection",
			Description: "The Cloudbuildv2 Connection resource",
			StructName:  "Connection",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Connection",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "connection",
						Required:    true,
						Description: "A full instance of a Connection",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Connection",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "connection",
						Required:    true,
						Description: "A full instance of a Connection",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Connection",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "connection",
						Required:    true,
						Description: "A full instance of a Connection",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Connection",
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
				Description: "The function used to list information about many Connection",
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
				"Connection": &dcl.Component{
					Title:           "Connection",
					ID:              "projects/{{project}}/locations/{{location}}/connections/{{name}}",
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
							"annotations": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Annotations",
								Description: "Allows clients to store small amounts of arbitrary data.",
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. Server assigned timestamp for when the connection was created.",
								Immutable:   true,
							},
							"disabled": &dcl.Property{
								Type:        "boolean",
								GoName:      "Disabled",
								Description: "If disabled is set to true, functionality is disabled for this connection. Repository based API methods and webhooks processing for repositories in this connection will be disabled.",
							},
							"etag": &dcl.Property{
								Type:        "string",
								GoName:      "Etag",
								ReadOnly:    true,
								Description: "This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.",
								Immutable:   true,
							},
							"githubConfig": &dcl.Property{
								Type:        "object",
								GoName:      "GithubConfig",
								GoType:      "ConnectionGithubConfig",
								Description: "Configuration for connections to github.com.",
								Conflicts: []string{
									"githubEnterpriseConfig",
									"gitlabConfig",
								},
								Properties: map[string]*dcl.Property{
									"appInstallationId": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "AppInstallationId",
										Description: "GitHub App installation id.",
									},
									"authorizerCredential": &dcl.Property{
										Type:        "object",
										GoName:      "AuthorizerCredential",
										GoType:      "ConnectionGithubConfigAuthorizerCredential",
										Description: "OAuth credential of the account that authorized the Cloud Build GitHub App. It is recommended to use a robot account instead of a human user account. The OAuth token must be tied to the Cloud Build GitHub App.",
										Properties: map[string]*dcl.Property{
											"oauthTokenSecretVersion": &dcl.Property{
												Type:        "string",
												GoName:      "OAuthTokenSecretVersion",
												Description: "A SecretManager resource containing the OAuth token that authorizes the Cloud Build connection. Format: `projects/*/secrets/*/versions/*`.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Secretmanager/SecretVersion",
														Field:    "selfLink",
													},
												},
											},
											"username": &dcl.Property{
												Type:        "string",
												GoName:      "Username",
												ReadOnly:    true,
												Description: "Output only. The username associated to this token.",
											},
										},
									},
								},
							},
							"githubEnterpriseConfig": &dcl.Property{
								Type:        "object",
								GoName:      "GithubEnterpriseConfig",
								GoType:      "ConnectionGithubEnterpriseConfig",
								Description: "Configuration for connections to an instance of GitHub Enterprise.",
								Conflicts: []string{
									"githubConfig",
									"gitlabConfig",
								},
								Required: []string{
									"hostUri",
								},
								Properties: map[string]*dcl.Property{
									"appId": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "AppId",
										Description: "Id of the GitHub App created from the manifest.",
									},
									"appInstallationId": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "AppInstallationId",
										Description: "ID of the installation of the GitHub App.",
									},
									"appSlug": &dcl.Property{
										Type:        "string",
										GoName:      "AppSlug",
										Description: "The URL-friendly name of the GitHub App.",
									},
									"hostUri": &dcl.Property{
										Type:        "string",
										GoName:      "HostUri",
										Description: "Required. The URI of the GitHub Enterprise host this connection is for.",
									},
									"privateKeySecretVersion": &dcl.Property{
										Type:        "string",
										GoName:      "PrivateKeySecretVersion",
										Description: "SecretManager resource containing the private key of the GitHub App, formatted as `projects/*/secrets/*/versions/*`.",
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Secretmanager/SecretVersion",
												Field:    "selfLink",
											},
										},
									},
									"serviceDirectoryConfig": &dcl.Property{
										Type:        "object",
										GoName:      "ServiceDirectoryConfig",
										GoType:      "ConnectionGithubEnterpriseConfigServiceDirectoryConfig",
										Description: "Configuration for using Service Directory to privately connect to a GitHub Enterprise server. This should only be set if the GitHub Enterprise server is hosted on-premises and not reachable by public internet. If this field is left empty, calls to the GitHub Enterprise server will be made over the public internet.",
										Required: []string{
											"service",
										},
										Properties: map[string]*dcl.Property{
											"service": &dcl.Property{
												Type:        "string",
												GoName:      "Service",
												Description: "Required. The Service Directory service name. Format: projects/{project}/locations/{location}/namespaces/{namespace}/services/{service}.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Servicedirectory/Service",
														Field:    "selfLink",
													},
												},
											},
										},
									},
									"sslCa": &dcl.Property{
										Type:        "string",
										GoName:      "SslCa",
										Description: "SSL certificate to use for requests to GitHub Enterprise.",
									},
									"webhookSecretSecretVersion": &dcl.Property{
										Type:        "string",
										GoName:      "WebhookSecretSecretVersion",
										Description: "SecretManager resource containing the webhook secret of the GitHub App, formatted as `projects/*/secrets/*/versions/*`.",
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Secretmanager/SecretVersion",
												Field:    "selfLink",
											},
										},
									},
								},
							},
							"gitlabConfig": &dcl.Property{
								Type:        "object",
								GoName:      "GitlabConfig",
								GoType:      "ConnectionGitlabConfig",
								Description: "Configuration for connections to gitlab.com or an instance of GitLab Enterprise.",
								Conflicts: []string{
									"githubConfig",
									"githubEnterpriseConfig",
								},
								Required: []string{
									"webhookSecretSecretVersion",
									"readAuthorizerCredential",
									"authorizerCredential",
								},
								Properties: map[string]*dcl.Property{
									"authorizerCredential": &dcl.Property{
										Type:        "object",
										GoName:      "AuthorizerCredential",
										GoType:      "ConnectionGitlabConfigAuthorizerCredential",
										Description: "Required. A GitLab personal access token with the `api` scope access.",
										Required: []string{
											"userTokenSecretVersion",
										},
										Properties: map[string]*dcl.Property{
											"userTokenSecretVersion": &dcl.Property{
												Type:        "string",
												GoName:      "UserTokenSecretVersion",
												Description: "Required. A SecretManager resource containing the user token that authorizes the Cloud Build connection. Format: `projects/*/secrets/*/versions/*`.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Secretmanager/SecretVersion",
														Field:    "selfLink",
													},
												},
											},
											"username": &dcl.Property{
												Type:        "string",
												GoName:      "Username",
												ReadOnly:    true,
												Description: "Output only. The username associated to this token.",
											},
										},
									},
									"hostUri": &dcl.Property{
										Type:          "string",
										GoName:        "HostUri",
										Description:   "The URI of the GitLab Enterprise host this connection is for. If not specified, the default value is https://gitlab.com.",
										ServerDefault: true,
									},
									"readAuthorizerCredential": &dcl.Property{
										Type:        "object",
										GoName:      "ReadAuthorizerCredential",
										GoType:      "ConnectionGitlabConfigReadAuthorizerCredential",
										Description: "Required. A GitLab personal access token with the minimum `read_api` scope access.",
										Required: []string{
											"userTokenSecretVersion",
										},
										Properties: map[string]*dcl.Property{
											"userTokenSecretVersion": &dcl.Property{
												Type:        "string",
												GoName:      "UserTokenSecretVersion",
												Description: "Required. A SecretManager resource containing the user token that authorizes the Cloud Build connection. Format: `projects/*/secrets/*/versions/*`.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Secretmanager/SecretVersion",
														Field:    "selfLink",
													},
												},
											},
											"username": &dcl.Property{
												Type:        "string",
												GoName:      "Username",
												ReadOnly:    true,
												Description: "Output only. The username associated to this token.",
											},
										},
									},
									"serverVersion": &dcl.Property{
										Type:        "string",
										GoName:      "ServerVersion",
										ReadOnly:    true,
										Description: "Output only. Version of the GitLab Enterprise server running on the `host_uri`.",
									},
									"serviceDirectoryConfig": &dcl.Property{
										Type:        "object",
										GoName:      "ServiceDirectoryConfig",
										GoType:      "ConnectionGitlabConfigServiceDirectoryConfig",
										Description: "Configuration for using Service Directory to privately connect to a GitLab Enterprise server. This should only be set if the GitLab Enterprise server is hosted on-premises and not reachable by public internet. If this field is left empty, calls to the GitLab Enterprise server will be made over the public internet.",
										Required: []string{
											"service",
										},
										Properties: map[string]*dcl.Property{
											"service": &dcl.Property{
												Type:        "string",
												GoName:      "Service",
												Description: "Required. The Service Directory service name. Format: projects/{project}/locations/{location}/namespaces/{namespace}/services/{service}.",
												ResourceReferences: []*dcl.PropertyResourceReference{
													&dcl.PropertyResourceReference{
														Resource: "Servicedirectory/Service",
														Field:    "selfLink",
													},
												},
											},
										},
									},
									"sslCa": &dcl.Property{
										Type:        "string",
										GoName:      "SslCa",
										Description: "SSL certificate to use for requests to GitLab Enterprise.",
									},
									"webhookSecretSecretVersion": &dcl.Property{
										Type:        "string",
										GoName:      "WebhookSecretSecretVersion",
										Description: "Required. Immutable. SecretManager resource containing the webhook secret of a GitLab Enterprise project, formatted as `projects/*/secrets/*/versions/*`.",
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Secretmanager/SecretVersion",
												Field:    "selfLink",
											},
										},
									},
								},
							},
							"installationState": &dcl.Property{
								Type:        "object",
								GoName:      "InstallationState",
								GoType:      "ConnectionInstallationState",
								ReadOnly:    true,
								Description: "Output only. Installation state of the Connection.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"actionUri": &dcl.Property{
										Type:        "string",
										GoName:      "ActionUri",
										ReadOnly:    true,
										Description: "Output only. Link to follow for next action. Empty string if the installation is already complete.",
										Immutable:   true,
									},
									"message": &dcl.Property{
										Type:        "string",
										GoName:      "Message",
										ReadOnly:    true,
										Description: "Output only. Message of what the user should do next to continue the installation. Empty string if the installation is already complete.",
										Immutable:   true,
									},
									"stage": &dcl.Property{
										Type:        "string",
										GoName:      "Stage",
										GoType:      "ConnectionInstallationStateStageEnum",
										ReadOnly:    true,
										Description: "Output only. Current step of the installation process. Possible values: STAGE_UNSPECIFIED, PENDING_CREATE_APP, PENDING_USER_OAUTH, PENDING_INSTALL_APP, COMPLETE",
										Immutable:   true,
										Enum: []string{
											"STAGE_UNSPECIFIED",
											"PENDING_CREATE_APP",
											"PENDING_USER_OAUTH",
											"PENDING_INSTALL_APP",
											"COMPLETE",
										},
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
								Description: "Immutable. The resource name of the connection, in the format `projects/{project}/locations/{location}/connections/{connection_id}`.",
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
							"reconciling": &dcl.Property{
								Type:        "boolean",
								GoName:      "Reconciling",
								ReadOnly:    true,
								Description: "Output only. Set to true when the connection is being set up or updated in the background.",
								Immutable:   true,
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. Server assigned timestamp for when the connection was updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
