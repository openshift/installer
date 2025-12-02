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
package apikeys

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLKeySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Apikeys/Key",
			Description: "The Apikeys Key resource",
			StructName:  "Key",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Key",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "key",
						Required:    true,
						Description: "A full instance of a Key",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Key",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "key",
						Required:    true,
						Description: "A full instance of a Key",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Key",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "key",
						Required:    true,
						Description: "A full instance of a Key",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Key",
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
				Description: "The function used to list information about many Key",
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
				"Key": &dcl.Component{
					Title: "Key",
					ID:    "projects/{{project}}/locations/global/keys/{{name}}",
					Locations: []string{
						"global",
					},
					ParentContainer: "project",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "Human-readable display name of this API key. Modifiable by user.",
							},
							"keyString": &dcl.Property{
								Type:        "string",
								GoName:      "KeyString",
								ReadOnly:    true,
								Description: "Output only. An encrypted and signed value held by this key. This field can be accessed only through the `GetKeyString` method.",
								Immutable:   true,
								Sensitive:   true,
							},
							"name": &dcl.Property{
								Type:        "string",
								GoName:      "Name",
								Description: "The resource name of the key. The name must be unique within the project, must conform with RFC-1034, is restricted to lower-cased letters, and has a maximum length of 63 characters. In another word, the name must match the regular expression: `[a-z]([a-z0-9-]{0,61}[a-z0-9])?`.",
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
							"restrictions": &dcl.Property{
								Type:        "object",
								GoName:      "Restrictions",
								GoType:      "KeyRestrictions",
								Description: "Key restrictions.",
								Properties: map[string]*dcl.Property{
									"androidKeyRestrictions": &dcl.Property{
										Type:        "object",
										GoName:      "AndroidKeyRestrictions",
										GoType:      "KeyRestrictionsAndroidKeyRestrictions",
										Description: "The Android apps that are allowed to use the key.",
										Conflicts: []string{
											"browserKeyRestrictions",
											"serverKeyRestrictions",
											"iosKeyRestrictions",
										},
										Required: []string{
											"allowedApplications",
										},
										Properties: map[string]*dcl.Property{
											"allowedApplications": &dcl.Property{
												Type:        "array",
												GoName:      "AllowedApplications",
												Description: "A list of Android applications that are allowed to make API calls with this key.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "object",
													GoType: "KeyRestrictionsAndroidKeyRestrictionsAllowedApplications",
													Required: []string{
														"sha1Fingerprint",
														"packageName",
													},
													Properties: map[string]*dcl.Property{
														"packageName": &dcl.Property{
															Type:        "string",
															GoName:      "PackageName",
															Description: "The package name of the application.",
														},
														"sha1Fingerprint": &dcl.Property{
															Type:        "string",
															GoName:      "Sha1Fingerprint",
															Description: "The SHA1 fingerprint of the application. For example, both sha1 formats are acceptable : DA:39:A3:EE:5E:6B:4B:0D:32:55:BF:EF:95:60:18:90:AF:D8:07:09 or DA39A3EE5E6B4B0D3255BFEF95601890AFD80709. Output format is the latter.",
														},
													},
												},
											},
										},
									},
									"apiTargets": &dcl.Property{
										Type:        "array",
										GoName:      "ApiTargets",
										Description: "A restriction for a specific service and optionally one or more specific methods. Requests are allowed if they match any of these restrictions. If no restrictions are specified, all targets are allowed.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "KeyRestrictionsApiTargets",
											Required: []string{
												"service",
											},
											Properties: map[string]*dcl.Property{
												"methods": &dcl.Property{
													Type:        "array",
													GoName:      "Methods",
													Description: "Optional. List of one or more methods that can be called. If empty, all methods for the service are allowed. A wildcard (*) can be used as the last symbol. Valid examples: `google.cloud.translate.v2.TranslateService.GetSupportedLanguage` `TranslateText` `Get*` `translate.googleapis.com.Get*`",
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"service": &dcl.Property{
													Type:        "string",
													GoName:      "Service",
													Description: "The service for this restriction. It should be the canonical service name, for example: `translate.googleapis.com`. You can use `gcloud services list` to get a list of services that are enabled in the project.",
												},
											},
										},
									},
									"browserKeyRestrictions": &dcl.Property{
										Type:        "object",
										GoName:      "BrowserKeyRestrictions",
										GoType:      "KeyRestrictionsBrowserKeyRestrictions",
										Description: "The HTTP referrers (websites) that are allowed to use the key.",
										Conflicts: []string{
											"serverKeyRestrictions",
											"androidKeyRestrictions",
											"iosKeyRestrictions",
										},
										Required: []string{
											"allowedReferrers",
										},
										Properties: map[string]*dcl.Property{
											"allowedReferrers": &dcl.Property{
												Type:        "array",
												GoName:      "AllowedReferrers",
												Description: "A list of regular expressions for the referrer URLs that are allowed to make API calls with this key.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
										},
									},
									"iosKeyRestrictions": &dcl.Property{
										Type:        "object",
										GoName:      "IosKeyRestrictions",
										GoType:      "KeyRestrictionsIosKeyRestrictions",
										Description: "The iOS apps that are allowed to use the key.",
										Conflicts: []string{
											"browserKeyRestrictions",
											"serverKeyRestrictions",
											"androidKeyRestrictions",
										},
										Required: []string{
											"allowedBundleIds",
										},
										Properties: map[string]*dcl.Property{
											"allowedBundleIds": &dcl.Property{
												Type:        "array",
												GoName:      "AllowedBundleIds",
												Description: "A list of bundle IDs that are allowed when making API calls with this key.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
										},
									},
									"serverKeyRestrictions": &dcl.Property{
										Type:        "object",
										GoName:      "ServerKeyRestrictions",
										GoType:      "KeyRestrictionsServerKeyRestrictions",
										Description: "The IP addresses of callers that are allowed to use the key.",
										Conflicts: []string{
											"browserKeyRestrictions",
											"androidKeyRestrictions",
											"iosKeyRestrictions",
										},
										Required: []string{
											"allowedIps",
										},
										Properties: map[string]*dcl.Property{
											"allowedIps": &dcl.Property{
												Type:        "array",
												GoName:      "AllowedIps",
												Description: "A list of the caller IP addresses that are allowed to make API calls with this key.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
										},
									},
								},
							},
							"uid": &dcl.Property{
								Type:        "string",
								GoName:      "Uid",
								ReadOnly:    true,
								Description: "Output only. Unique id in UUID4 format.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
