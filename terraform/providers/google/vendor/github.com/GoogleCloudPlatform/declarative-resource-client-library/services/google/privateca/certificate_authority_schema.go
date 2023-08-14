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

func DCLCertificateAuthoritySchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Privateca/CertificateAuthority",
			Description: "The Privateca CertificateAuthority resource",
			StructName:  "CertificateAuthority",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a CertificateAuthority",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificateAuthority",
						Required:    true,
						Description: "A full instance of a CertificateAuthority",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a CertificateAuthority",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificateAuthority",
						Required:    true,
						Description: "A full instance of a CertificateAuthority",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a CertificateAuthority",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificateAuthority",
						Required:    true,
						Description: "A full instance of a CertificateAuthority",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all CertificateAuthority",
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
					dcl.PathParameters{
						Name:     "caPool",
						Required: true,
						Schema: &dcl.PathParametersSchema{
							Type: "string",
						},
					},
				},
			},
			List: &dcl.Path{
				Description: "The function used to list information about many CertificateAuthority",
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
					dcl.PathParameters{
						Name:     "caPool",
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
				"CertificateAuthority": &dcl.Component{
					Title:           "CertificateAuthority",
					ID:              "projects/{{project}}/locations/{{location}}/caPools/{{ca_pool}}/certificateAuthorities/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"type",
							"config",
							"lifetime",
							"keySpec",
							"project",
							"location",
							"caPool",
						},
						Properties: map[string]*dcl.Property{
							"accessUrls": &dcl.Property{
								Type:        "object",
								GoName:      "AccessUrls",
								GoType:      "CertificateAuthorityAccessUrls",
								ReadOnly:    true,
								Description: "Output only. URLs for accessing content published by this CA, such as the CA certificate and CRLs.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"caCertificateAccessUrl": &dcl.Property{
										Type:        "string",
										GoName:      "CaCertificateAccessUrl",
										Description: "The URL where this CertificateAuthority's CA certificate is published. This will only be set for CAs that have been activated.",
										Immutable:   true,
									},
									"crlAccessUrls": &dcl.Property{
										Type:        "array",
										GoName:      "CrlAccessUrls",
										Description: "The URLs where this CertificateAuthority's CRLs are published. This will only be set for CAs that have been activated.",
										Immutable:   true,
										SendEmpty:   true,
										ListType:    "list",
										Items: &dcl.Property{
											Type:   "string",
											GoType: "string",
										},
									},
								},
							},
							"caCertificateDescriptions": &dcl.Property{
								Type:        "array",
								GoName:      "CaCertificateDescriptions",
								ReadOnly:    true,
								Description: "Output only. A structured description of this CertificateAuthority's CA certificate and its issuers. Ordered as self-to-root.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "object",
									GoType: "CertificateAuthorityCaCertificateDescriptions",
									Properties: map[string]*dcl.Property{
										"aiaIssuingCertificateUrls": &dcl.Property{
											Type:        "array",
											GoName:      "AiaIssuingCertificateUrls",
											Description: "Describes lists of issuer CA certificate URLs that appear in the \"Authority Information Access\" extension in the certificate.",
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "string",
												GoType: "string",
											},
										},
										"authorityKeyId": &dcl.Property{
											Type:        "object",
											GoName:      "AuthorityKeyId",
											GoType:      "CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId",
											Description: "Identifies the subject_key_id of the parent certificate, per https://tools.ietf.org/html/rfc5280#section-4.2.1.1",
											Properties: map[string]*dcl.Property{
												"keyId": &dcl.Property{
													Type:        "string",
													GoName:      "KeyId",
													Description: "Optional. The value of this KeyId encoded in lowercase hexadecimal. This is most likely the 160 bit SHA-1 hash of the public key.",
												},
											},
										},
										"certFingerprint": &dcl.Property{
											Type:        "object",
											GoName:      "CertFingerprint",
											GoType:      "CertificateAuthorityCaCertificateDescriptionsCertFingerprint",
											Description: "The hash of the x.509 certificate.",
											Properties: map[string]*dcl.Property{
												"sha256Hash": &dcl.Property{
													Type:        "string",
													GoName:      "Sha256Hash",
													Description: "The SHA 256 hash, encoded in hexadecimal, of the DER x509 certificate.",
												},
											},
										},
										"crlDistributionPoints": &dcl.Property{
											Type:        "array",
											GoName:      "CrlDistributionPoints",
											Description: "Describes a list of locations to obtain CRL information, i.e. the DistributionPoint.fullName described by https://tools.ietf.org/html/rfc5280#section-4.2.1.13",
											SendEmpty:   true,
											ListType:    "list",
											Items: &dcl.Property{
												Type:   "string",
												GoType: "string",
											},
										},
										"publicKey": &dcl.Property{
											Type:        "object",
											GoName:      "PublicKey",
											GoType:      "CertificateAuthorityCaCertificateDescriptionsPublicKey",
											Description: "The public key that corresponds to an issued certificate.",
											Required: []string{
												"key",
												"format",
											},
											Properties: map[string]*dcl.Property{
												"format": &dcl.Property{
													Type:        "string",
													GoName:      "Format",
													GoType:      "CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum",
													Description: "Required. The format of the public key. Possible values: PEM",
													Enum: []string{
														"PEM",
													},
												},
												"key": &dcl.Property{
													Type:        "string",
													GoName:      "Key",
													Description: "Required. A public key. The padding and encoding must match with the `KeyFormat` value specified for the `format` field.",
												},
											},
										},
										"subjectDescription": &dcl.Property{
											Type:        "object",
											GoName:      "SubjectDescription",
											GoType:      "CertificateAuthorityCaCertificateDescriptionsSubjectDescription",
											Description: "Describes some of the values in a certificate that are related to the subject and lifetime.",
											Properties: map[string]*dcl.Property{
												"hexSerialNumber": &dcl.Property{
													Type:        "string",
													GoName:      "HexSerialNumber",
													Description: "The serial number encoded in lowercase hexadecimal.",
												},
												"lifetime": &dcl.Property{
													Type:        "string",
													GoName:      "Lifetime",
													Description: "For convenience, the actual lifetime of an issued certificate.",
												},
												"notAfterTime": &dcl.Property{
													Type:        "string",
													Format:      "date-time",
													GoName:      "NotAfterTime",
													Description: "The time after which the certificate is expired. Per RFC 5280, the validity period for a certificate is the period of time from not_before_time through not_after_time, inclusive. Corresponds to 'not_before_time' + 'lifetime' - 1 second.",
												},
												"notBeforeTime": &dcl.Property{
													Type:        "string",
													Format:      "date-time",
													GoName:      "NotBeforeTime",
													Description: "The time at which the certificate becomes valid.",
												},
												"subject": &dcl.Property{
													Type:        "object",
													GoName:      "Subject",
													GoType:      "CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject",
													Description: "Contains distinguished name fields such as the common name, location and organization.",
													Properties: map[string]*dcl.Property{
														"commonName": &dcl.Property{
															Type:        "string",
															GoName:      "CommonName",
															Description: "The \"common name\" of the subject.",
														},
														"countryCode": &dcl.Property{
															Type:        "string",
															GoName:      "CountryCode",
															Description: "The country code of the subject.",
														},
														"locality": &dcl.Property{
															Type:        "string",
															GoName:      "Locality",
															Description: "The locality or city of the subject.",
														},
														"organization": &dcl.Property{
															Type:        "string",
															GoName:      "Organization",
															Description: "The organization of the subject.",
														},
														"organizationalUnit": &dcl.Property{
															Type:        "string",
															GoName:      "OrganizationalUnit",
															Description: "The organizational_unit of the subject.",
														},
														"postalCode": &dcl.Property{
															Type:        "string",
															GoName:      "PostalCode",
															Description: "The postal code of the subject.",
														},
														"province": &dcl.Property{
															Type:        "string",
															GoName:      "Province",
															Description: "The province, territory, or regional state of the subject.",
														},
														"streetAddress": &dcl.Property{
															Type:        "string",
															GoName:      "StreetAddress",
															Description: "The street address of the subject.",
														},
													},
												},
												"subjectAltName": &dcl.Property{
													Type:        "object",
													GoName:      "SubjectAltName",
													GoType:      "CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName",
													Description: "The subject alternative name fields.",
													Properties: map[string]*dcl.Property{
														"customSans": &dcl.Property{
															Type:        "array",
															GoName:      "CustomSans",
															Description: "Contains additional subject alternative name values.",
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "object",
																GoType: "CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans",
																Required: []string{
																	"objectId",
																	"critical",
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
																		GoType:      "CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId",
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
														"dnsNames": &dcl.Property{
															Type:        "array",
															GoName:      "DnsNames",
															Description: "Contains only valid, fully-qualified host names.",
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
														"emailAddresses": &dcl.Property{
															Type:        "array",
															GoName:      "EmailAddresses",
															Description: "Contains only valid RFC 2822 E-mail addresses.",
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
														"ipAddresses": &dcl.Property{
															Type:        "array",
															GoName:      "IPAddresses",
															Description: "Contains only valid 32-bit IPv4 addresses or RFC 4291 IPv6 addresses.",
															SendEmpty:   true,
															ListType:    "list",
															Items: &dcl.Property{
																Type:   "string",
																GoType: "string",
															},
														},
														"uris": &dcl.Property{
															Type:        "array",
															GoName:      "Uris",
															Description: "Contains only valid RFC 3986 URIs.",
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
										"subjectKeyId": &dcl.Property{
											Type:        "object",
											GoName:      "SubjectKeyId",
											GoType:      "CertificateAuthorityCaCertificateDescriptionsSubjectKeyId",
											Description: "Provides a means of identifiying certificates that contain a particular public key, per https://tools.ietf.org/html/rfc5280#section-4.2.1.2.",
											Properties: map[string]*dcl.Property{
												"keyId": &dcl.Property{
													Type:        "string",
													GoName:      "KeyId",
													Description: "Optional. The value of this KeyId encoded in lowercase hexadecimal. This is most likely the 160 bit SHA-1 hash of the public key.",
												},
											},
										},
										"x509Description": &dcl.Property{
											Type:        "object",
											GoName:      "X509Description",
											GoType:      "CertificateAuthorityCaCertificateDescriptionsX509Description",
											Description: "Describes some of the technical X.509 fields in a certificate.",
											Properties: map[string]*dcl.Property{
												"additionalExtensions": &dcl.Property{
													Type:        "array",
													GoName:      "AdditionalExtensions",
													Description: "Optional. Describes custom X.509 extensions.",
													SendEmpty:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "object",
														GoType: "CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions",
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
																GoType:      "CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId",
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
													ReadOnly:    true,
													Description: "Optional. Describes Online Certificate Status Protocol (OCSP) endpoint addresses that appear in the \"Authority Information Access\" extension in the certificate.",
													Immutable:   true,
													ListType:    "list",
													Items: &dcl.Property{
														Type:   "string",
														GoType: "string",
													},
												},
												"caOptions": &dcl.Property{
													Type:        "object",
													GoName:      "CaOptions",
													GoType:      "CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions",
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
													GoType:      "CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage",
													Description: "Optional. Indicates the intended use for keys that correspond to a certificate.",
													Properties: map[string]*dcl.Property{
														"baseKeyUsage": &dcl.Property{
															Type:        "object",
															GoName:      "BaseKeyUsage",
															GoType:      "CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage",
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
															GoType:      "CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage",
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
																GoType: "CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages",
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
														GoType: "CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds",
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
									},
								},
							},
							"caPool": &dcl.Property{
								Type:        "string",
								GoName:      "CaPool",
								Description: "The caPool for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Privateca/CaPool",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"config": &dcl.Property{
								Type:        "object",
								GoName:      "Config",
								GoType:      "CertificateAuthorityConfig",
								Description: "Required. Immutable. The config used to create a self-signed X.509 certificate or CSR.",
								Immutable:   true,
								Required: []string{
									"subjectConfig",
									"x509Config",
								},
								Properties: map[string]*dcl.Property{
									"publicKey": &dcl.Property{
										Type:        "object",
										GoName:      "PublicKey",
										GoType:      "CertificateAuthorityConfigPublicKey",
										ReadOnly:    true,
										Description: "Optional. The public key that corresponds to this config. This is, for example, used when issuing Certificates, but not when creating a self-signed CertificateAuthority or CertificateAuthority CSR.",
										Immutable:   true,
										Required: []string{
											"key",
											"format",
										},
										Properties: map[string]*dcl.Property{
											"format": &dcl.Property{
												Type:        "string",
												GoName:      "Format",
												GoType:      "CertificateAuthorityConfigPublicKeyFormatEnum",
												Description: "Required. The format of the public key. Possible values: PEM",
												Immutable:   true,
												Enum: []string{
													"PEM",
												},
											},
											"key": &dcl.Property{
												Type:        "string",
												GoName:      "Key",
												Description: "Required. A public key. The padding and encoding must match with the `KeyFormat` value specified for the `format` field.",
												Immutable:   true,
											},
										},
									},
									"subjectConfig": &dcl.Property{
										Type:        "object",
										GoName:      "SubjectConfig",
										GoType:      "CertificateAuthorityConfigSubjectConfig",
										Description: "Required. Specifies some of the values in a certificate that are related to the subject.",
										Immutable:   true,
										Required: []string{
											"subject",
										},
										Properties: map[string]*dcl.Property{
											"subject": &dcl.Property{
												Type:        "object",
												GoName:      "Subject",
												GoType:      "CertificateAuthorityConfigSubjectConfigSubject",
												Description: "Required. Contains distinguished name fields such as the common name, location and organization.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"commonName": &dcl.Property{
														Type:        "string",
														GoName:      "CommonName",
														Description: "The \"common name\" of the subject.",
														Immutable:   true,
													},
													"countryCode": &dcl.Property{
														Type:        "string",
														GoName:      "CountryCode",
														Description: "The country code of the subject.",
														Immutable:   true,
													},
													"locality": &dcl.Property{
														Type:        "string",
														GoName:      "Locality",
														Description: "The locality or city of the subject.",
														Immutable:   true,
													},
													"organization": &dcl.Property{
														Type:        "string",
														GoName:      "Organization",
														Description: "The organization of the subject.",
														Immutable:   true,
													},
													"organizationalUnit": &dcl.Property{
														Type:        "string",
														GoName:      "OrganizationalUnit",
														Description: "The organizational_unit of the subject.",
														Immutable:   true,
													},
													"postalCode": &dcl.Property{
														Type:        "string",
														GoName:      "PostalCode",
														Description: "The postal code of the subject.",
														Immutable:   true,
													},
													"province": &dcl.Property{
														Type:        "string",
														GoName:      "Province",
														Description: "The province, territory, or regional state of the subject.",
														Immutable:   true,
													},
													"streetAddress": &dcl.Property{
														Type:        "string",
														GoName:      "StreetAddress",
														Description: "The street address of the subject.",
														Immutable:   true,
													},
												},
											},
											"subjectAltName": &dcl.Property{
												Type:        "object",
												GoName:      "SubjectAltName",
												GoType:      "CertificateAuthorityConfigSubjectConfigSubjectAltName",
												Description: "Optional. The subject alternative name fields.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"customSans": &dcl.Property{
														Type:        "array",
														GoName:      "CustomSans",
														Description: "Contains additional subject alternative name values.",
														Immutable:   true,
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "object",
															GoType: "CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans",
															Required: []string{
																"objectId",
																"value",
															},
															Properties: map[string]*dcl.Property{
																"critical": &dcl.Property{
																	Type:        "boolean",
																	GoName:      "Critical",
																	Description: "Optional. Indicates whether or not this extension is critical (i.e., if the client does not know how to handle this extension, the client should consider this to be an error).",
																	Immutable:   true,
																},
																"objectId": &dcl.Property{
																	Type:        "object",
																	GoName:      "ObjectId",
																	GoType:      "CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId",
																	Description: "Required. The OID for this X.509 extension.",
																	Immutable:   true,
																	Required: []string{
																		"objectIdPath",
																	},
																	Properties: map[string]*dcl.Property{
																		"objectIdPath": &dcl.Property{
																			Type:        "array",
																			GoName:      "ObjectIdPath",
																			Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
																			Immutable:   true,
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
																	Immutable:   true,
																},
															},
														},
													},
													"dnsNames": &dcl.Property{
														Type:        "array",
														GoName:      "DnsNames",
														Description: "Contains only valid, fully-qualified host names.",
														Immutable:   true,
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "string",
															GoType: "string",
														},
													},
													"emailAddresses": &dcl.Property{
														Type:        "array",
														GoName:      "EmailAddresses",
														Description: "Contains only valid RFC 2822 E-mail addresses.",
														Immutable:   true,
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "string",
															GoType: "string",
														},
													},
													"ipAddresses": &dcl.Property{
														Type:        "array",
														GoName:      "IPAddresses",
														Description: "Contains only valid 32-bit IPv4 addresses or RFC 4291 IPv6 addresses.",
														Immutable:   true,
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "string",
															GoType: "string",
														},
													},
													"uris": &dcl.Property{
														Type:        "array",
														GoName:      "Uris",
														Description: "Contains only valid RFC 3986 URIs.",
														Immutable:   true,
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
									"x509Config": &dcl.Property{
										Type:        "object",
										GoName:      "X509Config",
										GoType:      "CertificateAuthorityConfigX509Config",
										Description: "Required. Describes how some of the technical X.509 fields in a certificate should be populated.",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"additionalExtensions": &dcl.Property{
												Type:        "array",
												GoName:      "AdditionalExtensions",
												Description: "Optional. Describes custom X.509 extensions.",
												Immutable:   true,
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "object",
													GoType: "CertificateAuthorityConfigX509ConfigAdditionalExtensions",
													Required: []string{
														"objectId",
														"value",
													},
													Properties: map[string]*dcl.Property{
														"critical": &dcl.Property{
															Type:        "boolean",
															GoName:      "Critical",
															Description: "Optional. Indicates whether or not this extension is critical (i.e., if the client does not know how to handle this extension, the client should consider this to be an error).",
															Immutable:   true,
														},
														"objectId": &dcl.Property{
															Type:        "object",
															GoName:      "ObjectId",
															GoType:      "CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId",
															Description: "Required. The OID for this X.509 extension.",
															Immutable:   true,
															Required: []string{
																"objectIdPath",
															},
															Properties: map[string]*dcl.Property{
																"objectIdPath": &dcl.Property{
																	Type:        "array",
																	GoName:      "ObjectIdPath",
																	Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
																	Immutable:   true,
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
															Immutable:   true,
														},
													},
												},
											},
											"aiaOcspServers": &dcl.Property{
												Type:        "array",
												GoName:      "AiaOcspServers",
												ReadOnly:    true,
												Description: "Optional. Describes Online Certificate Status Protocol (OCSP) endpoint addresses that appear in the \"Authority Information Access\" extension in the certificate.",
												Immutable:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "string",
													GoType: "string",
												},
											},
											"caOptions": &dcl.Property{
												Type:        "object",
												GoName:      "CaOptions",
												GoType:      "CertificateAuthorityConfigX509ConfigCaOptions",
												Description: "Optional. Describes options in this X509Parameters that are relevant in a CA certificate.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"isCa": &dcl.Property{
														Type:        "boolean",
														GoName:      "IsCa",
														Description: "Optional. Refers to the \"CA\" X.509 extension, which is a boolean value. When this value is missing, the extension will be omitted from the CA certificate.",
														Immutable:   true,
													},
													"maxIssuerPathLength": &dcl.Property{
														Type:        "integer",
														Format:      "int64",
														GoName:      "MaxIssuerPathLength",
														Description: "Optional. Refers to the path length restriction X.509 extension. For a CA certificate, this value describes the depth of subordinate CA certificates that are allowed. If this value is less than 0, the request will fail. If this value is missing, the max path length will be omitted from the CA certificate.",
														Immutable:   true,
													},
													"zeroMaxIssuerPathLength": &dcl.Property{
														Type:        "boolean",
														GoName:      "ZeroMaxIssuerPathLength",
														Description: "Optional. When true, the \"path length constraint\" in Basic Constraints extension will be set to 0. if both max_issuer_path_length and zero_max_issuer_path_length are unset, the max path length will be omitted from the CA certificate.",
														Immutable:   true,
													},
												},
											},
											"keyUsage": &dcl.Property{
												Type:        "object",
												GoName:      "KeyUsage",
												GoType:      "CertificateAuthorityConfigX509ConfigKeyUsage",
												Description: "Optional. Indicates the intended use for keys that correspond to a certificate.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"baseKeyUsage": &dcl.Property{
														Type:        "object",
														GoName:      "BaseKeyUsage",
														GoType:      "CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage",
														Description: "Describes high-level ways in which a key may be used.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"certSign": &dcl.Property{
																Type:        "boolean",
																GoName:      "CertSign",
																Description: "The key may be used to sign certificates.",
																Immutable:   true,
															},
															"contentCommitment": &dcl.Property{
																Type:        "boolean",
																GoName:      "ContentCommitment",
																Description: "The key may be used for cryptographic commitments. Note that this may also be referred to as \"non-repudiation\".",
																Immutable:   true,
															},
															"crlSign": &dcl.Property{
																Type:        "boolean",
																GoName:      "CrlSign",
																Description: "The key may be used sign certificate revocation lists.",
																Immutable:   true,
															},
															"dataEncipherment": &dcl.Property{
																Type:        "boolean",
																GoName:      "DataEncipherment",
																Description: "The key may be used to encipher data.",
																Immutable:   true,
															},
															"decipherOnly": &dcl.Property{
																Type:        "boolean",
																GoName:      "DecipherOnly",
																Description: "The key may be used to decipher only.",
																Immutable:   true,
															},
															"digitalSignature": &dcl.Property{
																Type:        "boolean",
																GoName:      "DigitalSignature",
																Description: "The key may be used for digital signatures.",
																Immutable:   true,
															},
															"encipherOnly": &dcl.Property{
																Type:        "boolean",
																GoName:      "EncipherOnly",
																Description: "The key may be used to encipher only.",
																Immutable:   true,
															},
															"keyAgreement": &dcl.Property{
																Type:        "boolean",
																GoName:      "KeyAgreement",
																Description: "The key may be used in a key agreement protocol.",
																Immutable:   true,
															},
															"keyEncipherment": &dcl.Property{
																Type:        "boolean",
																GoName:      "KeyEncipherment",
																Description: "The key may be used to encipher other keys.",
																Immutable:   true,
															},
														},
													},
													"extendedKeyUsage": &dcl.Property{
														Type:        "object",
														GoName:      "ExtendedKeyUsage",
														GoType:      "CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage",
														Description: "Detailed scenarios in which a key may be used.",
														Immutable:   true,
														Properties: map[string]*dcl.Property{
															"clientAuth": &dcl.Property{
																Type:        "boolean",
																GoName:      "ClientAuth",
																Description: "Corresponds to OID 1.3.6.1.5.5.7.3.2. Officially described as \"TLS WWW client authentication\", though regularly used for non-WWW TLS.",
																Immutable:   true,
															},
															"codeSigning": &dcl.Property{
																Type:        "boolean",
																GoName:      "CodeSigning",
																Description: "Corresponds to OID 1.3.6.1.5.5.7.3.3. Officially described as \"Signing of downloadable executable code client authentication\".",
																Immutable:   true,
															},
															"emailProtection": &dcl.Property{
																Type:        "boolean",
																GoName:      "EmailProtection",
																Description: "Corresponds to OID 1.3.6.1.5.5.7.3.4. Officially described as \"Email protection\".",
																Immutable:   true,
															},
															"ocspSigning": &dcl.Property{
																Type:        "boolean",
																GoName:      "OcspSigning",
																Description: "Corresponds to OID 1.3.6.1.5.5.7.3.9. Officially described as \"Signing OCSP responses\".",
																Immutable:   true,
															},
															"serverAuth": &dcl.Property{
																Type:        "boolean",
																GoName:      "ServerAuth",
																Description: "Corresponds to OID 1.3.6.1.5.5.7.3.1. Officially described as \"TLS WWW server authentication\", though regularly used for non-WWW TLS.",
																Immutable:   true,
															},
															"timeStamping": &dcl.Property{
																Type:        "boolean",
																GoName:      "TimeStamping",
																Description: "Corresponds to OID 1.3.6.1.5.5.7.3.8. Officially described as \"Binding the hash of an object to a time\".",
																Immutable:   true,
															},
														},
													},
													"unknownExtendedKeyUsages": &dcl.Property{
														Type:        "array",
														GoName:      "UnknownExtendedKeyUsages",
														Description: "Used to describe extended key usages that are not listed in the KeyUsage.ExtendedKeyUsageOptions message.",
														Immutable:   true,
														SendEmpty:   true,
														ListType:    "list",
														Items: &dcl.Property{
															Type:   "object",
															GoType: "CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages",
															Required: []string{
																"objectIdPath",
															},
															Properties: map[string]*dcl.Property{
																"objectIdPath": &dcl.Property{
																	Type:        "array",
																	GoName:      "ObjectIdPath",
																	Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
																	Immutable:   true,
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
												Immutable:   true,
												SendEmpty:   true,
												ListType:    "list",
												Items: &dcl.Property{
													Type:   "object",
													GoType: "CertificateAuthorityConfigX509ConfigPolicyIds",
													Required: []string{
														"objectIdPath",
													},
													Properties: map[string]*dcl.Property{
														"objectIdPath": &dcl.Property{
															Type:        "array",
															GoName:      "ObjectIdPath",
															Description: "Required. The parts of an OID path. The most significant parts of the path come first.",
															Immutable:   true,
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
								},
							},
							"createTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "CreateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this CertificateAuthority was created.",
								Immutable:   true,
							},
							"deleteTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "DeleteTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this CertificateAuthority was soft deleted, if it is in the DELETED state.",
								Immutable:   true,
							},
							"expireTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "ExpireTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this CertificateAuthority will be permanently purged, if it is in the DELETED state.",
								Immutable:   true,
							},
							"gcsBucket": &dcl.Property{
								Type:        "string",
								GoName:      "GcsBucket",
								Description: "Immutable. The name of a Cloud Storage bucket where this CertificateAuthority will publish content, such as the CA certificate and CRLs. This must be a bucket name, without any prefixes (such as `gs://`) or suffixes (such as `.googleapis.com`). For example, to use a bucket named `my-bucket`, you would simply specify `my-bucket`. If not specified, a managed bucket will be created.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Storage/Bucket",
										Field:    "name",
									},
								},
							},
							"keySpec": &dcl.Property{
								Type:        "object",
								GoName:      "KeySpec",
								GoType:      "CertificateAuthorityKeySpec",
								Description: "Required. Immutable. Used when issuing certificates for this CertificateAuthority. If this CertificateAuthority is a self-signed CertificateAuthority, this key is also used to sign the self-signed CA certificate. Otherwise, it is used to sign a CSR.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"algorithm": &dcl.Property{
										Type:        "string",
										GoName:      "Algorithm",
										GoType:      "CertificateAuthorityKeySpecAlgorithmEnum",
										Description: "The algorithm to use for creating a managed Cloud KMS key for a for a simplified experience. All managed keys will be have their ProtectionLevel as `HSM`. Possible values: RSA_PSS_2048_SHA256, RSA_PSS_3072_SHA256, RSA_PSS_4096_SHA256, RSA_PKCS1_2048_SHA256, RSA_PKCS1_3072_SHA256, RSA_PKCS1_4096_SHA256, EC_P256_SHA256, EC_P384_SHA384",
										Immutable:   true,
										Conflicts: []string{
											"cloudKmsKeyVersion",
										},
										Enum: []string{
											"RSA_PSS_2048_SHA256",
											"RSA_PSS_3072_SHA256",
											"RSA_PSS_4096_SHA256",
											"RSA_PKCS1_2048_SHA256",
											"RSA_PKCS1_3072_SHA256",
											"RSA_PKCS1_4096_SHA256",
											"EC_P256_SHA256",
											"EC_P384_SHA384",
										},
									},
									"cloudKmsKeyVersion": &dcl.Property{
										Type:        "string",
										GoName:      "CloudKmsKeyVersion",
										Description: "The resource name for an existing Cloud KMS CryptoKeyVersion in the format `projects/*/locations/*/keyRings/*/cryptoKeys/*/cryptoKeyVersions/*`. This option enables full flexibility in the key's capabilities and properties.",
										Immutable:   true,
										Conflicts: []string{
											"algorithm",
										},
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Cloudkms/CryptoKeyVersion",
												Field:    "name",
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
							"lifetime": &dcl.Property{
								Type:        "string",
								GoName:      "Lifetime",
								Description: "Required. The desired lifetime of the CA certificate. Used to create the \"not_before_time\" and \"not_after_time\" fields inside an X.509 certificate.",
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
								Description: "The resource name for this CertificateAuthority in the format `projects/*/locations/*/caPools/*/certificateAuthorities/*`.",
								Immutable:   true,
							},
							"pemCaCertificates": &dcl.Property{
								Type:        "array",
								GoName:      "PemCaCertificates",
								ReadOnly:    true,
								Description: "Output only. This CertificateAuthority's certificate chain, including the current CertificateAuthority's certificate. Ordered such that the root issuer is the final element (consistent with RFC 5246). For a self-signed CA, this will only list the current CertificateAuthority's certificate.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
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
							"state": &dcl.Property{
								Type:        "string",
								GoName:      "State",
								GoType:      "CertificateAuthorityStateEnum",
								ReadOnly:    true,
								Description: "Output only. The State for this CertificateAuthority. Possible values: ENABLED, DISABLED, STAGED, AWAITING_USER_ACTIVATION, DELETED",
								Immutable:   true,
								Enum: []string{
									"ENABLED",
									"DISABLED",
									"STAGED",
									"AWAITING_USER_ACTIVATION",
									"DELETED",
								},
							},
							"subordinateConfig": &dcl.Property{
								Type:        "object",
								GoName:      "SubordinateConfig",
								GoType:      "CertificateAuthoritySubordinateConfig",
								ReadOnly:    true,
								Description: "Optional. If this is a subordinate CertificateAuthority, this field will be set with the subordinate configuration, which describes its issuers. This may be updated, but this CertificateAuthority must continue to validate.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"certificateAuthority": &dcl.Property{
										Type:        "string",
										GoName:      "CertificateAuthority",
										Description: "Required. This can refer to a CertificateAuthority in the same project that was used to create a subordinate CertificateAuthority. This field is used for information and usability purposes only. The resource name is in the format `projects/*/locations/*/caPools/*/certificateAuthorities/*`.",
										Immutable:   true,
										Conflicts: []string{
											"pemIssuerChain",
										},
										ResourceReferences: []*dcl.PropertyResourceReference{
											&dcl.PropertyResourceReference{
												Resource: "Privateca/CertificateAuthority",
												Field:    "selfLink",
											},
										},
									},
									"pemIssuerChain": &dcl.Property{
										Type:        "object",
										GoName:      "PemIssuerChain",
										GoType:      "CertificateAuthoritySubordinateConfigPemIssuerChain",
										Description: "Required. Contains the PEM certificate chain for the issuers of this CertificateAuthority, but not pem certificate for this CA itself.",
										Immutable:   true,
										Conflicts: []string{
											"certificateAuthority",
										},
										Required: []string{
											"pemCertificates",
										},
										Properties: map[string]*dcl.Property{
											"pemCertificates": &dcl.Property{
												Type:        "array",
												GoName:      "PemCertificates",
												Description: "Required. Expected to be in leaf-to-root order according to RFC 5246.",
												Immutable:   true,
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
							"tier": &dcl.Property{
								Type:        "string",
								GoName:      "Tier",
								GoType:      "CertificateAuthorityTierEnum",
								ReadOnly:    true,
								Description: "Output only. The CaPool.Tier of the CaPool that includes this CertificateAuthority. Possible values: ENTERPRISE, DEVOPS",
								Immutable:   true,
								Enum: []string{
									"ENTERPRISE",
									"DEVOPS",
								},
							},
							"type": &dcl.Property{
								Type:        "string",
								GoName:      "Type",
								GoType:      "CertificateAuthorityTypeEnum",
								Description: "Required. Immutable. The Type of this CertificateAuthority. Possible values: SELF_SIGNED, SUBORDINATE",
								Immutable:   true,
								Enum: []string{
									"SELF_SIGNED",
									"SUBORDINATE",
								},
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this CertificateAuthority was last updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
