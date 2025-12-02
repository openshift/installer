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
package privateca

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func DCLCertificateTemplateSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Privateca/CertificateTemplate",
			Description: "Certificate Authority Service provides reusable and parameterized templates that you can use for common certificate issuance scenarios. A certificate template represents a relatively static and well-defined certificate issuance schema within an organization.  A certificate template can essentially become a full-fledged vertical certificate issuance framework.",
			StructName:  "CertificateTemplate",
			Reference: &dcl.Link{
				Text: "REST API",
				URL:  "https://cloud.google.com/certificate-authority-service/docs/reference/rest/v1/projects.locations.certificateTemplates",
			},
			Guides: []*dcl.Link{
				&dcl.Link{
					Text: "Understanding Certificate Templates",
					URL:  "https://cloud.google.com/certificate-authority-service/docs/certificate-template",
				},
				&dcl.Link{
					Text: "Common configurations and Certificate Profiles",
					URL:  "https://cloud.google.com/certificate-authority-service/docs/certificate-profile",
				},
			},
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a CertificateTemplate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificateTemplate",
						Required:    true,
						Description: "A full instance of a CertificateTemplate",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a CertificateTemplate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificateTemplate",
						Required:    true,
						Description: "A full instance of a CertificateTemplate",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a CertificateTemplate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificateTemplate",
						Required:    true,
						Description: "A full instance of a CertificateTemplate",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all CertificateTemplate",
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
				Description: "The function used to list information about many CertificateTemplate",
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
				"CertificateTemplate": &dcl.Component{
					Title:           "CertificateTemplate",
					ID:              "projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this CertificateTemplate was created.",
								Immutable:   true,
							},
							"description": &dcl.Property{
								Type:        "string",
								GoName:      "Description",
								Description: "Optional. A human-readable description of scenarios this template is intended for.",
							},
							"identityConstraints": &dcl.Property{
								Type:        "object",
								GoName:      "IdentityConstraints",
								GoType:      "CertificateTemplateIdentityConstraints",
								Description: "Optional. Describes constraints on identities that may be appear in Certificates issued using this template. If this is omitted, then this template will not add restrictions on a certificate's identity.",
								Required: []string{
									"allowSubjectPassthrough",
									"allowSubjectAltNamesPassthrough",
								},
								Properties: map[string]*dcl.Property{
									"allowSubjectAltNamesPassthrough": &dcl.Property{
										Type:        "boolean",
										GoName:      "AllowSubjectAltNamesPassthrough",
										Description: "Required. If this is true, the SubjectAltNames extension may be copied from a certificate request into the signed certificate. Otherwise, the requested SubjectAltNames will be discarded.",
									},
									"allowSubjectPassthrough": &dcl.Property{
										Type:        "boolean",
										GoName:      "AllowSubjectPassthrough",
										Description: "Required. If this is true, the Subject field may be copied from a certificate request into the signed certificate. Otherwise, the requested Subject will be discarded.",
									},
									"celExpression": &dcl.Property{
										Type:        "object",
										GoName:      "CelExpression",
										GoType:      "CertificateTemplateIdentityConstraintsCelExpression",
										Description: "Optional. A CEL expression that may be used to validate the resolved X.509 Subject and/or Subject Alternative Name before a certificate is signed. To see the full allowed syntax and some examples, see https://cloud.google.com/certificate-authority-service/docs/using-cel",
										Properties: map[string]*dcl.Property{
											"description": &dcl.Property{
												Type:        "string",
												GoName:      "Description",
												Description: "Optional. Description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.",
											},
											"expression": &dcl.Property{
												Type:        "string",
												GoName:      "Expression",
												Description: "Textual representation of an expression in Common Expression Language syntax.",
											},
											"location": &dcl.Property{
												Type:        "string",
												GoName:      "Location",
												Description: "Optional. String indicating the location of the expression for error reporting, e.g. a file name and a position in the file.",
											},
											"title": &dcl.Property{
												Type:        "string",
												GoName:      "Title",
												Description: "Optional. Title for the expression, i.e. a short string describing its purpose. This can be used e.g. in UIs which allow to enter the expression.",
											},
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
								Description: "Optional. Labels with user-defined metadata.",
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
								Description: "The resource name for this CertificateTemplate in the format `projects/*/locations/*/certificateTemplates/*`.",
								Immutable:   true,
							},
							"passthroughExtensions": &dcl.Property{
								Type:        "object",
								GoName:      "PassthroughExtensions",
								GoType:      "CertificateTemplatePassthroughExtensions",
								Description: "Optional. Describes the set of X.509 extensions that may appear in a Certificate issued using this CertificateTemplate. If a certificate request sets extensions that don't appear in the passthrough_extensions, those extensions will be dropped. If the issuing CaPool's IssuancePolicy defines baseline_values that don't appear here, the certificate issuance request will fail. If this is omitted, then this template will not add restrictions on a certificate's X.509 extensions. These constraints do not apply to X.509 extensions set in this CertificateTemplate's predefined_values.",
								Properties: map[string]*dcl.Property{
									"additionalExtensions": &dcl.Property{
										Type:        "array",
										GoName:      "AdditionalExtensions",
										Description: "Optional. A set of ObjectIds identifying custom X.509 extensions. Will be combined with known_extensions to determine the full set of X.509 extensions.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "CertificateTemplatePassthroughExtensionsAdditionalExtensions",
											Required: []string{
												"objectIdPath",
											},
											Properties: map[string]*dcl.Property{
												"objectIdPath": &dcl.Property{
													Type:        "array",
													GoName:      "ObjectIdPath",
													Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "integer",
														Format: "int64",
														GoType: "int64",
													},
												},
											},
										},
									},
									"knownExtensions": &dcl.Property{
										Type:        "array",
										GoName:      "KnownExtensions",
										Description: "Optional. A set of named X.509 extensions. Will be combined with additional_extensions to determine the full set of X.509 extensions.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "CertificateTemplatePassthroughExtensionsKnownExtensionsEnum",
											Enum: []string{
												"KNOWN_CERTIFICATE_EXTENSION_UNSPECIFIED",
												"BASE_KEY_USAGE",
												"EXTENDED_KEY_USAGE",
												"CA_OPTIONS",
												"POLICY_IDS",
												"AIA_OCSP_SERVERS",
											},
										},
									},
								},
							},
							"predefinedValues": &dcl.Property{
								Type:        "object",
								GoName:      "PredefinedValues",
								GoType:      "CertificateTemplatePredefinedValues",
								Description: "Optional. A set of X.509 values that will be applied to all issued certificates that use this template. If the certificate request includes conflicting values for the same properties, they will be overwritten by the values defined here. If the issuing CaPool's IssuancePolicy defines conflicting baseline_values for the same properties, the certificate issuance request will fail.",
								Properties: map[string]*dcl.Property{
									"additionalExtensions": &dcl.Property{
										Type:        "array",
										GoName:      "AdditionalExtensions",
										Description: "Optional. Describes custom X.509 extensions.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "CertificateTemplatePredefinedValuesAdditionalExtensions",
											Required: []string{
												"objectId",
												"value",
											},
											Properties: map[string]*dcl.Property{
												"critical": &dcl.Property{
													Type:        "boolean",
													GoName:      "Critical",
													Description: "Optional. Indicates whether or not this extension is critical (i.e., if the client does not know how to handle this extension, the client should consider this to be an error).",
												},
												"objectId": &dcl.Property{
													Type:        "object",
													GoName:      "ObjectId",
													GoType:      "CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId",
													Description: "Required. The OID for this X.509 extension.",
													Required: []string{
														"objectIdPath",
													},
													Properties: map[string]*dcl.Property{
														"objectIdPath": &dcl.Property{
															Type:        "array",
															GoName:      "ObjectIdPath",
															Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "integer",
																Format: "int64",
																GoType: "int64",
															},
														},
													},
												},
												"value": &dcl.Property{
													Type:        "string",
													GoName:      "Value",
													Description: "Required. The value of this X.509 extension.",
												},
											},
										},
									},
									"aiaOcspServers": &dcl.Property{
										Type:        "array",
										GoName:      "AiaOcspServers",
										Description: "Optional. Describes Online Certificate Status Protocol (OCSP) endpoint addresses that appear in the \"Authority Information Access\" extension in the certificate.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
									"caOptions": &dcl.Property{
										Type:        "object",
										GoName:      "CaOptions",
										GoType:      "CertificateTemplatePredefinedValuesCaOptions",
										Description: "Optional. Describes options in this X509Parameters that are relevant in a CA certificate.",
										Properties: map[string]*dcl.Property{
											"isCa": &dcl.Property{
												Type:        "boolean",
												GoName:      "IsCa",
												Description: "Optional. Refers to the \"CA\" X.509 extension, which is a boolean value. When this value is missing, the extension will be omitted from the CA certificate.",
											},
											"maxIssuerPathLength": &dcl.Property{
												Type:        "integer",
												Format:      "int64",
												GoName:      "MaxIssuerPathLength",
												Description: "Optional. Refers to the path length restriction X.509 extension. For a CA certificate, this value describes the depth of subordinate CA certificates that are allowed. If this value is less than 0, the request will fail. If this value is missing, the max path length will be omitted from the CA certificate.",
											},
										},
									},
									"keyUsage": &dcl.Property{
										Type:        "object",
										GoName:      "KeyUsage",
										GoType:      "CertificateTemplatePredefinedValuesKeyUsage",
										Description: "Optional. Indicates the intended use for keys that correspond to a certificate.",
										Properties: map[string]*dcl.Property{
											"baseKeyUsage": &dcl.Property{
												Type:        "object",
												GoName:      "BaseKeyUsage",
												GoType:      "CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage",
												Description: "Describes high-level ways in which a key may be used.",
												Properties: map[string]*dcl.Property{
													"certSign": &dcl.Property{
														Type:        "boolean",
														GoName:      "CertSign",
														Description: "The key may be used to sign certificates.",
													},
													"contentCommitment": &dcl.Property{
														Type:        "boolean",
														GoName:      "ContentCommitment",
														Description: "The key may be used for cryptographic commitments. Note that this may also be referred to as \"non-repudiation\".",
													},
													"crlSign": &dcl.Property{
														Type:        "boolean",
														GoName:      "CrlSign",
														Description: "The key may be used sign certificate revocation lists.",
													},
													"dataEncipherment": &dcl.Property{
														Type:        "boolean",
														GoName:      "DataEncipherment",
														Description: "The key may be used to encipher data.",
													},
													"decipherOnly": &dcl.Property{
														Type:        "boolean",
														GoName:      "DecipherOnly",
														Description: "The key may be used to decipher only.",
													},
													"digitalSignature": &dcl.Property{
														Type:        "boolean",
														GoName:      "DigitalSignature",
														Description: "The key may be used for digital signatures.",
													},
													"encipherOnly": &dcl.Property{
														Type:        "boolean",
														GoName:      "EncipherOnly",
														Description: "The key may be used to encipher only.",
													},
													"keyAgreement": &dcl.Property{
														Type:        "boolean",
														GoName:      "KeyAgreement",
														Description: "The key may be used in a key agreement protocol.",
													},
													"keyEncipherment": &dcl.Property{
														Type:        "boolean",
														GoName:      "KeyEncipherment",
														Description: "The key may be used to encipher other keys.",
													},
												},
											},
											"extendedKeyUsage": &dcl.Property{
												Type:        "object",
												GoName:      "ExtendedKeyUsage",
												GoType:      "CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage",
												Description: "Detailed scenarios in which a key may be used.",
												Properties: map[string]*dcl.Property{
													"clientAuth": &dcl.Property{
														Type:        "boolean",
														GoName:      "ClientAuth",
														Description: "Corresponds to OID 1.3.6.1.5.5.7.3.2. Officially described as \"TLS WWW client authentication\", though regularly used for non-WWW TLS.",
													},
													"codeSigning": &dcl.Property{
														Type:        "boolean",
														GoName:      "CodeSigning",
														Description: "Corresponds to OID 1.3.6.1.5.5.7.3.3. Officially described as \"Signing of downloadable executable code client authentication\".",
													},
													"emailProtection": &dcl.Property{
														Type:        "boolean",
														GoName:      "EmailProtection",
														Description: "Corresponds to OID 1.3.6.1.5.5.7.3.4. Officially described as \"Email protection\".",
													},
													"ocspSigning": &dcl.Property{
														Type:        "boolean",
														GoName:      "OcspSigning",
														Description: "Corresponds to OID 1.3.6.1.5.5.7.3.9. Officially described as \"Signing OCSP responses\".",
													},
													"serverAuth": &dcl.Property{
														Type:        "boolean",
														GoName:      "ServerAuth",
														Description: "Corresponds to OID 1.3.6.1.5.5.7.3.1. Officially described as \"TLS WWW server authentication\", though regularly used for non-WWW TLS.",
													},
													"timeStamping": &dcl.Property{
														Type:        "boolean",
														GoName:      "TimeStamping",
														Description: "Corresponds to OID 1.3.6.1.5.5.7.3.8. Officially described as \"Binding the hash of an object to a time\".",
													},
												},
											},
											"unknownExtendedKeyUsages": &dcl.Property{
												Type:        "array",
												GoName:      "UnknownExtendedKeyUsages",
												Description: "Used to describe extended key usages that are not listed in the KeyUsage.ExtendedKeyUsageOptions message.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "object",
													GoType: "CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages",
													Required: []string{
														"objectIdPath",
													},
													Properties: map[string]*dcl.Property{
														"objectIdPath": &dcl.Property{
															Type:        "array",
															GoName:      "ObjectIdPath",
															Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "integer",
																Format: "int64",
																GoType: "int64",
															},
														},
													},
												},
											},
										},
									},
									"policyIds": &dcl.Property{
										Type:        "array",
										GoName:      "PolicyIds",
										Description: "Optional. Describes the X.509 certificate policy object identifiers, per https://tools.ietf.org/html/rfc5280#section-4.2.1.4.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "CertificateTemplatePredefinedValuesPolicyIds",
											Required: []string{
												"objectIdPath",
											},
											Properties: map[string]*dcl.Property{
												"objectIdPath": &dcl.Property{
													Type:        "array",
													GoName:      "ObjectIdPath",
													Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "integer",
														Format: "int64",
														GoType: "int64",
													},
												},
											},
										},
									},
								},
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
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this CertificateTemplate was updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
