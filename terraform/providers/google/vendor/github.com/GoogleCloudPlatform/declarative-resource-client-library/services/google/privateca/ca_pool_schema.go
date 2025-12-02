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

func DCLCaPoolSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Privateca/CaPool",
			Description: "The Privateca CaPool resource",
			StructName:  "CaPool",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a CaPool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "caPool",
						Required:    true,
						Description: "A full instance of a CaPool",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a CaPool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "caPool",
						Required:    true,
						Description: "A full instance of a CaPool",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a CaPool",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "caPool",
						Required:    true,
						Description: "A full instance of a CaPool",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all CaPool",
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
				Description: "The function used to list information about many CaPool",
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
				"CaPool": &dcl.Component{
					Title:           "CaPool",
					ID:              "projects/{{project}}/locations/{{location}}/caPools/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"tier",
							"project",
							"location",
						},
						Properties: map[string]*dcl.Property{
							"issuancePolicy": &dcl.Property{
								Type:        "object",
								GoName:      "IssuancePolicy",
								GoType:      "CaPoolIssuancePolicy",
								Description: "Optional. The IssuancePolicy to control how Certificates will be issued from this CaPool.",
								Properties: map[string]*dcl.Property{
									"allowedIssuanceModes": &dcl.Property{
										Type:        "object",
										GoName:      "AllowedIssuanceModes",
										GoType:      "CaPoolIssuancePolicyAllowedIssuanceModes",
										Description: "Optional. If specified, then only methods allowed in the IssuanceModes may be used to issue Certificates.",
										Properties: map[string]*dcl.Property{
											"allowConfigBasedIssuance": &dcl.Property{
												Type:        "boolean",
												GoName:      "AllowConfigBasedIssuance",
												Description: "Optional. When true, allows callers to create Certificates by specifying a CertificateConfig.",
											},
											"allowCsrBasedIssuance": &dcl.Property{
												Type:        "boolean",
												GoName:      "AllowCsrBasedIssuance",
												Description: "Optional. When true, allows callers to create Certificates by specifying a CSR.",
											},
										},
									},
									"allowedKeyTypes": &dcl.Property{
										Type:        "array",
										GoName:      "AllowedKeyTypes",
										Description: "Optional. If any AllowedKeyType is specified, then the certificate request's public key must match one of the key types listed here. Otherwise, any key may be used.",
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "object",
											GoType: "CaPoolIssuancePolicyAllowedKeyTypes",
											Properties: map[string]*dcl.Property{
												"ellipticCurve": &dcl.Property{
													Type:        "object",
													GoName:      "EllipticCurve",
													GoType:      "CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve",
													Description: "Represents an allowed Elliptic Curve key type.",
													Conflicts: []string{
														"rsa",
													},
													Properties: map[string]*dcl.Property{
														"signatureAlgorithm": &dcl.Property{
															Type:        "string",
															GoName:      "SignatureAlgorithm",
															GoType:      "CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum",
															Description: "Optional. A signature algorithm that must be used. If this is omitted, any EC-based signature algorithm will be allowed. Possible values: EC_SIGNATURE_ALGORITHM_UNSPECIFIED, ECDSA_P256, ECDSA_P384, EDDSA_25519",
															Enum: []string{
																"EC_SIGNATURE_ALGORITHM_UNSPECIFIED",
																"ECDSA_P256",
																"ECDSA_P384",
																"EDDSA_25519",
															},
														},
													},
												},
												"rsa": &dcl.Property{
													Type:        "object",
													GoName:      "Rsa",
													GoType:      "CaPoolIssuancePolicyAllowedKeyTypesRsa",
													Description: "Represents an allowed RSA key type.",
													Conflicts: []string{
														"ellipticCurve",
													},
													Properties: map[string]*dcl.Property{
														"maxModulusSize": &dcl.Property{
															Type:        "integer",
															Format:      "int64",
															GoName:      "MaxModulusSize",
															Description: "Optional. The maximum allowed RSA modulus size, in bits. If this is not set, or if set to zero, the service will not enforce an explicit upper bound on RSA modulus sizes.",
														},
														"minModulusSize": &dcl.Property{
															Type:        "integer",
															Format:      "int64",
															GoName:      "MinModulusSize",
															Description: "Optional. The minimum allowed RSA modulus size, in bits. If this is not set, or if set to zero, the service-level min RSA modulus size will continue to apply.",
														},
													},
												},
											},
										},
									},
									"baselineValues": &dcl.Property{
										Type:        "object",
										GoName:      "BaselineValues",
										GoType:      "CaPoolIssuancePolicyBaselineValues",
										Description: "Optional. A set of X.509 values that will be applied to all certificates issued through this CaPool. If a certificate request includes conflicting values for the same properties, they will be overwritten by the values defined here. If a certificate request uses a CertificateTemplate that defines conflicting predefined_values for the same properties, the certificate issuance request will fail.",
										Properties: map[string]*dcl.Property{
											"additionalExtensions": &dcl.Property{
												Type:        "array",
												GoName:      "AdditionalExtensions",
												Description: "Optional. Describes custom X.509 extensions.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "object",
													GoType: "CaPoolIssuancePolicyBaselineValuesAdditionalExtensions",
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
															GoType:      "CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId",
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
												GoType:      "CaPoolIssuancePolicyBaselineValuesCaOptions",
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
												GoType:      "CaPoolIssuancePolicyBaselineValuesKeyUsage",
												Description: "Optional. Indicates the intended use for keys that correspond to a certificate.",
												Properties: map[string]*dcl.Property{
													"baseKeyUsage": &dcl.Property{
														Type:        "object",
														GoName:      "BaseKeyUsage",
														GoType:      "CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage",
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
														GoType:      "CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage",
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
															GoType: "CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages",
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
													GoType: "CaPoolIssuancePolicyBaselineValuesPolicyIds",
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
									"identityConstraints": &dcl.Property{
										Type:        "object",
										GoName:      "IdentityConstraints",
										GoType:      "CaPoolIssuancePolicyIdentityConstraints",
										Description: "Optional. Describes constraints on identities that may appear in Certificates issued through this CaPool. If this is omitted, then this CaPool will not add restrictions on a certificate's identity.",
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
												GoType:      "CaPoolIssuancePolicyIdentityConstraintsCelExpression",
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
									"maximumLifetime": &dcl.Property{
										Type:        "string",
										GoName:      "MaximumLifetime",
										Description: "Optional. The maximum lifetime allowed for issued Certificates. Note that if the issuing CertificateAuthority expires before a Certificate's requested maximum_lifetime, the effective lifetime will be explicitly truncated to match it.",
									},
									"passthroughExtensions": &dcl.Property{
										Type:        "object",
										GoName:      "PassthroughExtensions",
										GoType:      "CaPoolIssuancePolicyPassthroughExtensions",
										Description: "Optional. Describes the set of X.509 extensions that may appear in a Certificate issued through this CaPool. If a certificate request sets extensions that don't appear in the passthrough_extensions, those extensions will be dropped. If a certificate request uses a CertificateTemplate with predefined_values that don't appear here, the certificate issuance request will fail. If this is omitted, then this CaPool will not add restrictions on a certificate's X.509 extensions. These constraints do not apply to X.509 extensions set in this CaPool's baseline_values.",
										Properties: map[string]*dcl.Property{
											"additionalExtensions": &dcl.Property{
												Type:        "array",
												GoName:      "AdditionalExtensions",
												Description: "Optional. A set of ObjectIds identifying custom X.509 extensions. Will be combined with known_extensions to determine the full set of X.509 extensions.",
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "object",
													GoType: "CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions",
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
													GoType: "CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum",
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
								Description: "The resource name for this CaPool in the format `projects/*/locations/*/caPools/*`.",
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
							"publishingOptions": &dcl.Property{
								Type:        "object",
								GoName:      "PublishingOptions",
								GoType:      "CaPoolPublishingOptions",
								Description: "Optional. The PublishingOptions to follow when issuing Certificates from any CertificateAuthority in this CaPool.",
								Properties: map[string]*dcl.Property{
									"publishCaCert": &dcl.Property{
										Type:        "boolean",
										GoName:      "PublishCaCert",
										Description: "Optional. When true, publishes each CertificateAuthority's CA certificate and includes its URL in the \"Authority Information Access\" X.509 extension in all issued Certificates. If this is false, the CA certificate will not be published and the corresponding X.509 extension will not be written in issued certificates.",
									},
									"publishCrl": &dcl.Property{
										Type:        "boolean",
										GoName:      "PublishCrl",
										Description: "Optional. When true, publishes each CertificateAuthority's CRL and includes its URL in the \"CRL Distribution Points\" X.509 extension in all issued Certificates. If this is false, CRLs will not be published and the corresponding X.509 extension will not be written in issued certificates. CRLs will expire 7 days from their creation. However, we will rebuild daily. CRLs are also rebuilt shortly after a certificate is revoked.",
									},
								},
							},
							"tier": &dcl.Property{
								Type:        "string",
								GoName:      "Tier",
								GoType:      "CaPoolTierEnum",
								Description: "Required. Immutable. The Tier of this CaPool. Possible values: TIER_UNSPECIFIED, ENTERPRISE, DEVOPS",
								Immutable:   true,
								Enum: []string{
									"TIER_UNSPECIFIED",
									"ENTERPRISE",
									"DEVOPS",
								},
							},
						},
					},
				},
			},
		},
	}
}
