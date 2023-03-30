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

func DCLCertificateSchema() *dcl.Schema {
	return &dcl.Schema{
		Info: &dcl.Info{
			Title:       "Privateca/Certificate",
			Description: "The Privateca Certificate resource",
			StructName:  "Certificate",
		},
		Paths: &dcl.Paths{
			Get: &dcl.Path{
				Description: "The function used to get information about a Certificate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificate",
						Required:    true,
						Description: "A full instance of a Certificate",
					},
				},
			},
			Apply: &dcl.Path{
				Description: "The function used to apply information about a Certificate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificate",
						Required:    true,
						Description: "A full instance of a Certificate",
					},
				},
			},
			Delete: &dcl.Path{
				Description: "The function used to delete a Certificate",
				Parameters: []dcl.PathParameters{
					dcl.PathParameters{
						Name:        "certificate",
						Required:    true,
						Description: "A full instance of a Certificate",
					},
				},
			},
			DeleteAll: &dcl.Path{
				Description: "The function used to delete all Certificate",
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
				Description: "The function used to list information about many Certificate",
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
				"Certificate": &dcl.Component{
					Title:           "Certificate",
					ID:              "projects/{{project}}/locations/{{location}}/caPools/{{ca_pool}}/certificates/{{name}}",
					ParentContainer: "project",
					LabelsField:     "labels",
					HasCreate:       true,
					SchemaProperty: dcl.Property{
						Type: "object",
						Required: []string{
							"name",
							"lifetime",
							"project",
							"location",
							"caPool",
						},
						Properties: map[string]*dcl.Property{
							"caPool": &dcl.Property{
								Type:        "string",
								GoName:      "CaPool",
								Description: "The ca_pool for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Privateca/CaPool",
										Field:    "name",
										Parent:   true,
									},
								},
							},
							"certificateAuthority": &dcl.Property{
								Type:        "string",
								GoName:      "CertificateAuthority",
								Description: "The certificate authority for the resource",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Privateca/CertificateAuthority",
										Field:    "name",
									},
								},
							},
							"certificateDescription": &dcl.Property{
								Type:        "object",
								GoName:      "CertificateDescription",
								GoType:      "CertificateCertificateDescription",
								ReadOnly:    true,
								Description: "Output only. A structured description of the issued X.509 certificate.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"aiaIssuingCertificateUrls": &dcl.Property{
										Type:        "array",
										GoName:      "AiaIssuingCertificateUrls",
										Description: "Describes lists of issuer CA certificate URLs that appear in the \"Authority Information Access\" extension in the certificate.",
										Immutable:   true,
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
										GoType:      "CertificateCertificateDescriptionAuthorityKeyId",
										Description: "Identifies the subject_key_id of the parent certificate, per https://tools.ietf.org/html/rfc5280#section-4.2.1.1",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"keyId": &dcl.Property{
												Type:        "string",
												GoName:      "KeyId",
												Description: "Optional. The value of this KeyId encoded in lowercase hexadecimal. This is most likely the 160 bit SHA-1 hash of the public key.",
												Immutable:   true,
											},
										},
									},
									"certFingerprint": &dcl.Property{
										Type:        "object",
										GoName:      "CertFingerprint",
										GoType:      "CertificateCertificateDescriptionCertFingerprint",
										Description: "The hash of the x.509 certificate.",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"sha256Hash": &dcl.Property{
												Type:        "string",
												GoName:      "Sha256Hash",
												Description: "The SHA 256 hash, encoded in hexadecimal, of the DER x509 certificate.",
												Immutable:   true,
											},
										},
									},
									"crlDistributionPoints": &dcl.Property{
										Type:        "array",
										GoName:      "CrlDistributionPoints",
										Description: "Describes a list of locations to obtain CRL information, i.e. the DistributionPoint.fullName described by https://tools.ietf.org/html/rfc5280#section-4.2.1.13",
										Immutable:   true,
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
										GoType:      "CertificateCertificateDescriptionPublicKey",
										Description: "The public key that corresponds to an issued certificate.",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"format": &dcl.Property{
												Type:        "string",
												GoName:      "Format",
												GoType:      "CertificateCertificateDescriptionPublicKeyFormatEnum",
												Description: "Required. The format of the public key. Possible values: KEY_FORMAT_UNSPECIFIED, PEM",
												Immutable:   true,
												Enum: []string{
													"KEY_FORMAT_UNSPECIFIED",
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
									"subjectDescription": &dcl.Property{
										Type:        "object",
										GoName:      "SubjectDescription",
										GoType:      "CertificateCertificateDescriptionSubjectDescription",
										Description: "Describes some of the values in a certificate that are related to the subject and lifetime.",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"hexSerialNumber": &dcl.Property{
												Type:        "string",
												GoName:      "HexSerialNumber",
												Description: "The serial number encoded in lowercase hexadecimal.",
												Immutable:   true,
											},
											"lifetime": &dcl.Property{
												Type:        "string",
												GoName:      "Lifetime",
												Description: "For convenience, the actual lifetime of an issued certificate.",
												Immutable:   true,
											},
											"notAfterTime": &dcl.Property{
												Type:        "string",
												Format:      "date-time",
												GoName:      "NotAfterTime",
												Description: "The time after which the certificate is expired. Per RFC 5280, the validity period for a certificate is the period of time from not_before_time through not_after_time, inclusive. Corresponds to 'not_before_time' + 'lifetime' - 1 second.",
												Immutable:   true,
											},
											"notBeforeTime": &dcl.Property{
												Type:        "string",
												Format:      "date-time",
												GoName:      "NotBeforeTime",
												Description: "The time at which the certificate becomes valid.",
												Immutable:   true,
											},
											"subject": &dcl.Property{
												Type:        "object",
												GoName:      "Subject",
												GoType:      "CertificateCertificateDescriptionSubjectDescriptionSubject",
												Description: "Contains distinguished name fields such as the common name, location and / organization.",
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
												GoType:      "CertificateCertificateDescriptionSubjectDescriptionSubjectAltName",
												Description: "The subject alternative name fields.",
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
															GoType: "CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans",
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
																	GoType:      "CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId",
																	Description: "Required. The OID for this X.509 extension.",
																	Immutable:   true,
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
									"subjectKeyId": &dcl.Property{
										Type:        "object",
										GoName:      "SubjectKeyId",
										GoType:      "CertificateCertificateDescriptionSubjectKeyId",
										Description: "Provides a means of identifiying certificates that contain a particular public key, per https://tools.ietf.org/html/rfc5280#section-4.2.1.2.",
										Immutable:   true,
										Properties: map[string]*dcl.Property{
											"keyId": &dcl.Property{
												Type:        "string",
												GoName:      "KeyId",
												Description: "Optional. The value of this KeyId encoded in lowercase hexadecimal. This is most likely the 160 bit SHA-1 hash of the public key.",
												Immutable:   true,
											},
										},
									},
									"x509Description": &dcl.Property{
										Type:        "object",
										GoName:      "X509Description",
										GoType:      "CertificateCertificateDescriptionX509Description",
										Description: "Describes some of the technical X.509 fields in a certificate.",
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
													GoType: "CertificateCertificateDescriptionX509DescriptionAdditionalExtensions",
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
															GoType:      "CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId",
															Description: "Required. The OID for this X.509 extension.",
															Immutable:   true,
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
												Description: "Optional. Describes Online Certificate Status Protocol (OCSP) endpoint addresses that appear in the \"Authority Information Access\" extension in the certificate.",
												Immutable:   true,
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
												GoType:      "CertificateCertificateDescriptionX509DescriptionCaOptions",
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
												},
											},
											"keyUsage": &dcl.Property{
												Type:        "object",
												GoName:      "KeyUsage",
												GoType:      "CertificateCertificateDescriptionX509DescriptionKeyUsage",
												Description: "Optional. Indicates the intended use for keys that correspond to a certificate.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"baseKeyUsage": &dcl.Property{
														Type:        "object",
														GoName:      "BaseKeyUsage",
														GoType:      "CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage",
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
														GoType:      "CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage",
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
															GoType: "CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages",
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
													GoType: "CertificateCertificateDescriptionX509DescriptionPolicyIds",
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
							"certificateTemplate": &dcl.Property{
								Type:        "string",
								GoName:      "CertificateTemplate",
								Description: "Immutable. The resource name for a CertificateTemplate used to issue this certificate, in the format `projects/*/locations/*/certificateTemplates/*`. If this is specified, the caller must have the necessary permission to use this template. If this is omitted, no template will be used. This template must be in the same location as the Certificate.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Privateca/CertificateTemplate",
										Field:    "selfLink",
									},
								},
							},
							"config": &dcl.Property{
								Type:        "object",
								GoName:      "Config",
								GoType:      "CertificateConfig",
								Description: "Immutable. A description of the certificate and key that does not require X.509 or ASN.1.",
								Immutable:   true,
								Conflicts: []string{
									"pemCsr",
								},
								Required: []string{
									"subjectConfig",
									"x509Config",
								},
								Properties: map[string]*dcl.Property{
									"publicKey": &dcl.Property{
										Type:        "object",
										GoName:      "PublicKey",
										GoType:      "CertificateConfigPublicKey",
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
												GoType:      "CertificateConfigPublicKeyFormatEnum",
												Description: "Required. The format of the public key. Possible values: KEY_FORMAT_UNSPECIFIED, PEM",
												Immutable:   true,
												Enum: []string{
													"KEY_FORMAT_UNSPECIFIED",
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
										GoType:      "CertificateConfigSubjectConfig",
										Description: "Required. Specifies some of the values in a certificate that are related to the subject.",
										Immutable:   true,
										Required: []string{
											"subject",
										},
										Properties: map[string]*dcl.Property{
											"subject": &dcl.Property{
												Type:        "object",
												GoName:      "Subject",
												GoType:      "CertificateConfigSubjectConfigSubject",
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
												GoType:      "CertificateConfigSubjectConfigSubjectAltName",
												Description: "Optional. The subject alternative name fields.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
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
										GoType:      "CertificateConfigX509Config",
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
													GoType: "CertificateConfigX509ConfigAdditionalExtensions",
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
															GoType:      "CertificateConfigX509ConfigAdditionalExtensionsObjectId",
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
												Description: "Optional. Describes Online Certificate Status Protocol (OCSP) endpoint addresses that appear in the \"Authority Information Access\" extension in the certificate.",
												Immutable:   true,
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
												GoType:      "CertificateConfigX509ConfigCaOptions",
												Description: "Optional. Describes options in this X509Parameters that are relevant in a CA certificate.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"isCa": &dcl.Property{
														Type:        "boolean",
														GoName:      "IsCa",
														Description: "Optional. When true, the \"CA\" in Basic Constraints extension will be set to true.",
														Immutable:   true,
													},
													"maxIssuerPathLength": &dcl.Property{
														Type:        "integer",
														Format:      "int64",
														GoName:      "MaxIssuerPathLength",
														Description: "Optional. Refers to the \"path length constraint\" in Basic Constraints extension. For a CA certificate, this value describes the depth of subordinate CA certificates that are allowed. If this value is less than 0, the request will fail.",
														Immutable:   true,
													},
													"nonCa": &dcl.Property{
														Type:        "boolean",
														GoName:      "NonCa",
														Description: "Optional. When true, the \"CA\" in Basic Constraints extension will be set to false. If both `is_ca` and `non_ca` are unset, the extension will be omitted from the CA certificate.",
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
												GoType:      "CertificateConfigX509ConfigKeyUsage",
												Description: "Optional. Indicates the intended use for keys that correspond to a certificate.",
												Immutable:   true,
												Properties: map[string]*dcl.Property{
													"baseKeyUsage": &dcl.Property{
														Type:        "object",
														GoName:      "BaseKeyUsage",
														GoType:      "CertificateConfigX509ConfigKeyUsageBaseKeyUsage",
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
														GoType:      "CertificateConfigX509ConfigKeyUsageExtendedKeyUsage",
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
															GoType: "CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages",
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
													GoType: "CertificateConfigX509ConfigPolicyIds",
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
								Description: "Output only. The time at which this Certificate was created.",
								Immutable:   true,
							},
							"issuerCertificateAuthority": &dcl.Property{
								Type:        "string",
								GoName:      "IssuerCertificateAuthority",
								ReadOnly:    true,
								Description: "Output only. The resource name of the issuing CertificateAuthority in the format `projects/*/locations/*/caPools/*/certificateAuthorities/*`.",
								Immutable:   true,
								ResourceReferences: []*dcl.PropertyResourceReference{
									&dcl.PropertyResourceReference{
										Resource: "Privateca/CertificateAuthority",
										Field:    "selfLink",
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
								Description: "Required. Immutable. The desired lifetime of a certificate. Used to create the \"not_before_time\" and \"not_after_time\" fields inside an X.509 certificate. Note that the lifetime may be truncated if it would extend past the life of any certificate authority in the issuing chain.",
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
								Description: "The resource name for this Certificate in the format `projects/*/locations/*/caPools/*/certificates/*`.",
								Immutable:   true,
							},
							"pemCertificate": &dcl.Property{
								Type:        "string",
								GoName:      "PemCertificate",
								ReadOnly:    true,
								Description: "Output only. The pem-encoded, signed X.509 certificate.",
								Immutable:   true,
							},
							"pemCertificateChain": &dcl.Property{
								Type:        "array",
								GoName:      "PemCertificateChain",
								ReadOnly:    true,
								Description: "Output only. The chain that may be used to verify the X.509 certificate. Expected to be in issuer-to-root order according to RFC 5246.",
								Immutable:   true,
								ListType:    "list",
								Items: &dcl.Property{
									Type:   "string",
									GoType: "string",
								},
							},
							"pemCsr": &dcl.Property{
								Type:        "string",
								GoName:      "PemCsr",
								Description: "Immutable. A pem-encoded X.509 certificate signing request (CSR).",
								Immutable:   true,
								Conflicts: []string{
									"config",
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
							"revocationDetails": &dcl.Property{
								Type:        "object",
								GoName:      "RevocationDetails",
								GoType:      "CertificateRevocationDetails",
								ReadOnly:    true,
								Description: "Output only. Details regarding the revocation of this Certificate. This Certificate is considered revoked if and only if this field is present.",
								Immutable:   true,
								Properties: map[string]*dcl.Property{
									"revocationState": &dcl.Property{
										Type:        "string",
										GoName:      "RevocationState",
										GoType:      "CertificateRevocationDetailsRevocationStateEnum",
										Description: "Indicates why a Certificate was revoked. Possible values: REVOCATION_REASON_UNSPECIFIED, KEY_COMPROMISE, CERTIFICATE_AUTHORITY_COMPROMISE, AFFILIATION_CHANGED, SUPERSEDED, CESSATION_OF_OPERATION, CERTIFICATE_HOLD, PRIVILEGE_WITHDRAWN, ATTRIBUTE_AUTHORITY_COMPROMISE",
										Immutable:   true,
										Enum: []string{
											"REVOCATION_REASON_UNSPECIFIED",
											"KEY_COMPROMISE",
											"CERTIFICATE_AUTHORITY_COMPROMISE",
											"AFFILIATION_CHANGED",
											"SUPERSEDED",
											"CESSATION_OF_OPERATION",
											"CERTIFICATE_HOLD",
											"PRIVILEGE_WITHDRAWN",
											"ATTRIBUTE_AUTHORITY_COMPROMISE",
										},
									},
									"revocationTime": &dcl.Property{
										Type:        "string",
										Format:      "date-time",
										GoName:      "RevocationTime",
										Description: "The time at which this Certificate was revoked.",
										Immutable:   true,
									},
								},
							},
							"subjectMode": &dcl.Property{
								Type:        "string",
								GoName:      "SubjectMode",
								GoType:      "CertificateSubjectModeEnum",
								Description: "Immutable. Specifies how the Certificate's identity fields are to be decided. If this is omitted, the `DEFAULT` subject mode will be used. Possible values: SUBJECT_REQUEST_MODE_UNSPECIFIED, DEFAULT, REFLECTED_SPIFFE",
								Immutable:   true,
								Enum: []string{
									"SUBJECT_REQUEST_MODE_UNSPECIFIED",
									"DEFAULT",
									"REFLECTED_SPIFFE",
								},
							},
							"updateTime": &dcl.Property{
								Type:        "string",
								Format:      "date-time",
								GoName:      "UpdateTime",
								ReadOnly:    true,
								Description: "Output only. The time at which this Certificate was updated.",
								Immutable:   true,
							},
						},
					},
				},
			},
		},
	}
}
