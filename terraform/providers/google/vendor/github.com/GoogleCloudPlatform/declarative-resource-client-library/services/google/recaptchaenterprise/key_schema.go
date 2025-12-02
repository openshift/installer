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
package recaptchaenterprise

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLKeySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "RecaptchaEnterprise/Key",
			Description: "The RecaptchaEnterprise Key resource",
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
					Title:           "Key",
					ID:              "projects/{{project}}/keys/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"displayName",
							"project",
						},
						Properties: map[string]*dcl.Property{
							"androidSettings": &dcl.Property{
								Type:        "object",
								GoName:      "AndroidSettings",
								GoType:      "KeyAndroidSettings",
								Description: "Settings for keys that can be used by Android apps.",
								Conflicts: []string{
									"webSettings",
									"iosSettings",
								},
								Properties: map[string]*dcl.Property{
									"allowAllPackageNames": &dcl.Property{
										Type:        "boolean",
										GoName:      "AllowAllPackageNames",
										Description: "If set to true, it means allowed_package_names will not be enforced.",
									},
									"allowedPackageNames": &dcl.Property{
										Type:        "array",
										GoName:      "AllowedPackageNames",
										Description: "Android package names of apps allowed to use the key. Example: 'com.companyname.appname'",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "The timestamp corresponding to the creation of this Key.",
								Immutable:   true,
							},
							"displayName": &dcl.Property{
								Type:        "string",
								GoName:      "DisplayName",
								Description: "Human-readable display name of this key. Modifiable by user.",
							},
							"iosSettings": &dcl.Property{
								Type:        "object",
								GoName:      "IosSettings",
								GoType:      "KeyIosSettings",
								Description: "Settings for keys that can be used by iOS apps.",
								Conflicts: []string{
									"webSettings",
									"androidSettings",
								},
								Properties: map[string]*dcl.Property{
									"allowAllBundleIds": &dcl.Property{
										Type:        "boolean",
										GoName:      "AllowAllBundleIds",
										Description: "If set to true, it means allowed_bundle_ids will not be enforced.",
									},
									"allowedBundleIds": &dcl.Property{
										Type:        "array",
										GoName:      "AllowedBundleIds",
										Description: "iOS bundle ids of apps allowed to use the key. Example: 'com.companyname.productname.appname'",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
								},
							},
							"labels": &dcl.Property{
								Type: "object",
								AdditionalProperties: &dcl.Property{
									Type: "string",
								},
								GoName:      "Labels",
								Description: "See [Creating and managing labels](https://cloud.google.com/recaptcha-enterprise/docs/labels).",
							},
							"name": &dcl.Property{
								Type:                     "string",
								GoName:                   "Name",
								Description:              "The resource name for the Key in the format \"projects/{project}/keys/{key}\".",
								Immutable:                true,
								ServerGeneratedParameter: true,
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
							"testingOptions": &dcl.Property{
								Type:        "object",
								GoName:      "TestingOptions",
								GoType:      "KeyTestingOptions",
								Description: "Options for user acceptance testing.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"testingChallenge": &dcl.Property{
										Type:          "string",
										GoName:        "TestingChallenge",
										GoType:        "KeyTestingOptionsTestingChallengeEnum",
										Description:   "For challenge-based keys only (CHECKBOX, INVISIBLE), all challenge requests for this site will return nocaptcha if NOCAPTCHA, or an unsolvable challenge if UNSOLVABLE_CHALLENGE. Possible values: TESTING_CHALLENGE_UNSPECIFIED, NOCAPTCHA, UNSOLVABLE_CHALLENGE",
										Immutable:     true,
										ServerDefault: true,
										Enum: []string{
											"TESTING_CHALLENGE_UNSPECIFIED",
											"NOCAPTCHA",
											"UNSOLVABLE_CHALLENGE",
										},
									},
									"testingScore": &dcl.Property{
										Type:        "number",
										Format:      "double",
										GoName:      "TestingScore",
										Description: "All assessments for this Key will return this score. Must be between 0 (likely not legitimate) and 1 (likely legitimate) inclusive.",
										Immutable:   true,
									},
								},
							},
							"webSettings": &dcl.Property{
								Type:        "object",
								GoName:      "WebSettings",
								GoType:      "KeyWebSettings",
								Description: "Settings for keys that can be used by websites.",
								Conflicts: []string{
									"androidSettings",
									"iosSettings",
								},
								Required: []string{
									"integrationType",
								},
								Properties: map[string]*dcl.Property{
									"allowAllDomains": &dcl.Property{
										Type:        "boolean",
										GoName:      "AllowAllDomains",
										Description: "If set to true, it means allowed_domains will not be enforced.",
									},
									"allowAmpTraffic": &dcl.Property{
										Type:        "boolean",
										GoName:      "AllowAmpTraffic",
										Description: "If set to true, the key can be used on AMP (Accelerated Mobile Pages) websites. This is supported only for the SCORE integration type.",
									},
									"allowedDomains": &dcl.Property{
										Type:        "array",
										GoName:      "AllowedDomains",
										Description: "Domains or subdomains of websites allowed to use the key. All subdomains of an allowed domain are automatically allowed. A valid domain requires a host and must not include any path, port, query or fragment. Examples: 'example.com' or 'subdomain.example.com'",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"challengeSecurityPreference": &dcl.Property{
										Type:          "string",
										GoName:        "ChallengeSecurityPreference",
										GoType:        "KeyWebSettingsChallengeSecurityPreferenceEnum",
										Description:   "Settings for the frequency and difficulty at which this key triggers captcha challenges. This should only be specified for IntegrationTypes CHECKBOX and INVISIBLE. Possible values: CHALLENGE_SECURITY_PREFERENCE_UNSPECIFIED, USABILITY, BALANCE, SECURITY",
										ServerDefault: true,
										Enum: []string{
											"CHALLENGE_SECURITY_PREFERENCE_UNSPECIFIED",
											"USABILITY",
											"BALANCE",
											"SECURITY",
										},
									},
									"integrationType": &dcl.Property{
										Type:        "string",
										GoName:      "IntegrationType",
										GoType:      "KeyWebSettingsIntegrationTypeEnum",
										Description: "Required. Describes how this key is integrated with the website. Possible values: SCORE, CHECKBOX, INVISIBLE",
										Immutable:   true,
										Enum: []string{
											"SCORE",
											"CHECKBOX",
											"INVISIBLE",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
