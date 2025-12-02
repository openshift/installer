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

func DCLUptimeCheckConfigSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Monitoring/UptimeCheckConfig",
			Description: "The Monitoring UptimeCheckConfig resource",
			StructName:  "UptimeCheckConfig",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a UptimeCheckConfig",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "uptimeCheckConfig",
						Required:    true,
						Description: "A full instance of a UptimeCheckConfig",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a UptimeCheckConfig",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "uptimeCheckConfig",
						Required:    true,
						Description: "A full instance of a UptimeCheckConfig",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a UptimeCheckConfig",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "uptimeCheckConfig",
						Required:    true,
						Description: "A full instance of a UptimeCheckConfig",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all UptimeCheckConfig",
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
				Description: "The function used to list information about many UptimeCheckConfig",
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
				"UptimeCheckConfig": &dcl.Component{
					Title:           "UptimeCheckConfig",
					ID:              "projects/{{project}}/uptimeCheckConfigs/{{name}}",
					UsesStateHint:   true,
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"displayName",
							"timeout",
						},
						Properties: map[string]*dcl.Property{
							"contentMatchers": &dcl.Property{
								Type:        "array",
								GoName:      "ContentMatchers",
								Description: "The content that is expected to appear in the data returned by the target server against which the check is run.  Currently, only the first entry in the `content_matchers` list is supported, and additional entries will be ignored. This field is optional and should only be specified if a content match is required as part of the/ Uptime check.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "UptimeCheckConfigContentMatchers",
									Required: []string{
										"content",
									},
									Properties: map[string]*dcl.Property{
										"content": &dcl.Property{
											Type:   "string",
											GoName: "Content",
										},
										"matcher": &dcl.Property{
											Type:        "string",
											GoName:      "Matcher",
											GoType:      "UptimeCheckConfigContentMatchersMatcherEnum",
											Description: " Possible values: CONTENT_MATCHER_OPTION_UNSPECIFIED, CONTAINS_STRING, NOT_CONTAINS_STRING, MATCHES_REGEX, NOT_MATCHES_REGEX",
											Default:     "CONTAINS_STRING",
											Enum: []string{
												"CONTENT_MATCHER_OPTION_UNSPECIFIED",
												"CONTAINS_STRING",
												"NOT_CONTAINS_STRING",
												"MATCHES_REGEX",
												"NOT_MATCHES_REGEX",
											},
										},
									},
								},
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "A human-friendly name for the Uptime check configuration. The display name should be unique within a Stackdriver Workspace in order to make it easier to identify; however, uniqueness is not enforced. Required.",
							},
							"httpCheck": &dcl.Property{
								Type:        "object",
								GoName:      "HttpCheck",
								GoType:      "UptimeCheckConfigHttpCheck",
								Description: "Contains information needed to make an HTTP or HTTPS check.",
								Conflicts: []string{
									"tcpCheck",
								},
								Properties: map[string]*dcl.Property{
									"authInfo": &dcl.Property{
										Type:        "object",
										GoName:      "AuthInfo",
										GoType:      "UptimeCheckConfigHttpCheckAuthInfo",
										Description: "The authentication information. Optional when creating an HTTP check; defaults to empty.",
										Required: []string{
											"username",
											"password",
										},
										Properties: map[string]*dcl.Property{
											"password": &dcl.Property{
												Type:       "string",
												GoName:     "Password",
												Sensitive:  true,
												Unreadable: true,
											},
											"username": &dcl.Property{
												Type:   "string",
												GoName: "Username",
											},
										},
									},
									"body": &dcl.Property{
										Type:        "string",
										GoName:      "Body",
										Description: "The request body associated with the HTTP POST request. If `content_type` is `URL_ENCODED`, the body passed in must be URL-encoded. Users can provide a `Content-Length` header via the `headers` field or the API will do so. If the `request_method` is `GET` and `body` is not empty, the API will return an error. The maximum byte size is 1 megabyte. Note: As with all `bytes` fields JSON representations are base64 encoded. e.g.: \"foo=bar\" in URL-encoded form is \"foo%3Dbar\" and in base64 encoding is \"Zm9vJTI1M0RiYXI=\".",
									},
									"contentType": &dcl.Property{
										Type:        "string",
										GoName:      "ContentType",
										GoType:      "UptimeCheckConfigHttpCheckContentTypeEnum",
										Description: "The content type to use for the check.  Possible values: TYPE_UNSPECIFIED, URL_ENCODED",
										Immutable:   true,
										Enum: []string{
											"TYPE_UNSPECIFIED",
											"URL_ENCODED",
										},
									},
									"headers": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:        "Headers",
										Description:   "The list of headers to send as part of the Uptime check request. If two headers have the same key and different values, they should be entered as a single header, with the value being a comma-separated list of all the desired values as described at https://www.w3.org/Protocols/rfc2616/rfc2616.txt (page 31). Entering two separate headers with the same key in a Create call will cause the first to be overwritten by the second. The maximum number of headers allowed is 100.",
										ServerDefault: true,
										Unreadable:    true,
									},
									"maskHeaders": &dcl.Property{
										Type:        "boolean",
										GoName:      "MaskHeaders",
										Description: "Boolean specifying whether to encrypt the header information. Encryption should be specified for any headers related to authentication that you do not wish to be seen when retrieving the configuration. The server will be responsible for encrypting the headers. On Get/List calls, if `mask_headers` is set to `true` then the headers will be obscured with `******.`",
										Immutable:   true,
									},
									"path": &dcl.Property{
										Type:        "string",
										GoName:      "Path",
										Description: "Optional (defaults to \"/\"). The path to the page against which to run the check. Will be combined with the `host` (specified within the `monitored_resource`) and `port` to construct the full URL. If the provided path does not begin with \"/\", a \"/\" will be prepended automatically.",
										Default:     "/",
									},
									"port": &dcl.Property{
										Type:          "integer",
										Format:        "int64",
										GoName:        "Port",
										Description:   "Optional (defaults to 80 when `use_ssl` is `false`, and 443 when `use_ssl` is `true`). The TCP port on the HTTP server against which to run the check. Will be combined with host (specified within the `monitored_resource`) and `path` to construct the full URL.",
										ServerDefault: true,
									},
									"requestMethod": &dcl.Property{
										Type:        "string",
										GoName:      "RequestMethod",
										GoType:      "UptimeCheckConfigHttpCheckRequestMethodEnum",
										Description: "The HTTP request method to use for the check. If set to `METHOD_UNSPECIFIED` then `request_method` defaults to `GET`.",
										Immutable:   true,
										Default:     "GET",
										Enum: []string{
											"METHOD_UNSPECIFIED",
											"GET",
											"POST",
										},
									},
									"useSsl": &dcl.Property{
										Type:        "boolean",
										GoName:      "UseSsl",
										Description: "If `true`, use HTTPS instead of HTTP to run the check.",
									},
									"validateSsl": &dcl.Property{
										Type:        "boolean",
										GoName:      "ValidateSsl",
										Description: "Boolean specifying whether to include SSL certificate validation as a part of the Uptime check. Only applies to checks where `monitored_resource` is set to `uptime_url`. If `use_ssl` is `false`, setting `validate_ssl` to `true` has no effect.",
									},
								},
							},
							"monitoredResource": &dcl.Property{
								Type:        "object",
								GoName:      "MonitoredResource",
								GoType:      "UptimeCheckConfigMonitoredResource",
								Description: "The [monitored resource](https://cloud.google.com/monitoring/api/resources) associated with the configuration. The following monitored resource types are supported for Uptime checks:   `uptime_url`,   `gce_instance`,   `gae_app`,   `aws_ec2_instance`,   `aws_elb_load_balancer`",
								Immutable:   true,
								Conflicts: []string{
									"resourceGroup",
								},
								Required: []string{
									"type",
									"filterLabels",
								},
								Properties: map[string]*dcl.Property{
									"filterLabels": &dcl.Property{
										Type: "object",
										AdditionalProperties: &dcl.Property{
											Type: "string",
										},
										GoName:    "FilterLabels",
										Immutable: true,
									},
									"type": &dcl.Property{
										Type:      "string",
										GoName:    "Type",
										Immutable: true,
									},
								},
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "A unique resource name for this Uptime check configuration. The format is: projects/[PROJECT_ID_OR_NUMBER]/uptimeCheckConfigs/[UPTIME_CHECK_ID] This field should be omitted when creating the Uptime check configuration; on create, the resource name is assigned by the server and included in the response.",
								Immutable:                true,
								ServerGeneratedParameter: true,
							},
							"period": &dcl.Property{
								Type:        "string",
								GoName:      "Period",
								Description: "How often, in seconds, the Uptime check is performed. Currently, the only supported values are `60s` (1 minute), `300s` (5 minutes), `600s` (10 minutes), and `900s` (15 minutes). Optional, defaults to `60s`.",
								Default:     "60s",
							},
							"project": &dcl.Property{
								Type:        "string",
								GoName:      "Project",
								Description: "The project for this uptime check config.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Cloudresourcemanager/Project",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"resourceGroup": &dcl.Property{
								Type:        "object",
								GoName:      "ResourceGroup",
								GoType:      "UptimeCheckConfigResourceGroup",
								Description: "The group resource associated with the configuration.",
								Immutable:   true,
								Conflicts: []string{
									"monitoredResource",
								},
								Properties: map[string]*dcl.Property{
									"groupId": &dcl.Property{
										Type:        "string",
										GoName:      "GroupId",
										Description: "The group of resources being monitored. Should be only the `[GROUP_ID]`, and not the full-path `projects/[PROJECT_ID_OR_NUMBER]/groups/[GROUP_ID]`.",
										Immutable:   true,
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Monitoring/Group",
												Field:    "name",
											},
										},
									},
									"resourceType": &dcl.Property{
										Type:        "string",
										GoName:      "ResourceType",
										GoType:      "UptimeCheckConfigResourceGroupResourceTypeEnum",
										Description: "The resource type of the group members. Possible values: RESOURCE_TYPE_UNSPECIFIED, INSTANCE, AWS_ELB_LOAD_BALANCER",
										Immutable:   true,
										Enum: []string{
											"RESOURCE_TYPE_UNSPECIFIED",
											"INSTANCE",
											"AWS_ELB_LOAD_BALANCER",
										},
									},
								},
							},
							"selectedRegions": &dcl.Property{
								Type:        "array",
								GoName:      "SelectedRegions",
								Description: "The list of regions from which the check will be run. Some regions contain one location, and others contain more than one. If this field is specified, enough regions must be provided to include a minimum of 3 locations.  Not specifying this field will result in Uptime checks running from all available regions.",
								SendEmpty:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"tcpCheck": &dcl.Property{
								Type:        "object",
								GoName:      "TcpCheck",
								GoType:      "UptimeCheckConfigTcpCheck",
								Description: "Contains information needed to make a TCP check.",
								Conflicts: []string{
									"httpCheck",
								},
								Required: []string{
									"port",
								},
								Properties: map[string]*dcl.Property{
									"port": &dcl.Property{
										Type:        "integer",
										Format:      "int64",
										GoName:      "Port",
										Description: "The TCP port on the server against which to run the check. Will be combined with host (specified within the `monitored_resource`) to construct the full URL. Required.",
									},
								},
							},
							"timeout": &dcl.Property{
								Type:        "string",
								GoName:      "Timeout",
								Description: "The maximum amount of time to wait for the request to complete (must be between 1 and 60 seconds). Required.",
							},
						},
					},
				},
			},
		},
	}
}
