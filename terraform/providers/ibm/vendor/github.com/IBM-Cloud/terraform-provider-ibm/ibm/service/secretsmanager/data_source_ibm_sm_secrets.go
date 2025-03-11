// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmSecrets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmSecretsRead,

		Schema: map[string]*schema.Schema{
			"sort": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort a collection of secrets by the specified field in ascending order. To sort in descending order use the `-` character. Available values: id | created_at | updated_at | expiration_date | secret_type | name",
			},
			"search": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Obtain a collection of secrets that contain the specified string in one or more of the fields: `id`, `name`, `description`,\n        `labels`, `secret_type`.",
			},
			"groups": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter secrets by groups. You can apply multiple filters by using a comma-separated list of secret group IDs. If you need to filter secrets that are in the default secret group, use the `default` keyword.",
			},
			"secret_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter secrets by secret types.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"match_all_labels": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter secrets by a label or a combination of labels.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources in a collection.",
			},
			"secrets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A collection of secret metadata.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier that is associated with the entity that created the secret.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when a resource was created. The date format follows RFC 3339.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A CRN that uniquely identifies an IBM Cloud resource.",
						},
						"custom_metadata": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "The secret metadata that a user can customize.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An extended description of your secret.To protect your privacy, do not use personal data, such as your name or location, as a description for your secret group.",
						},
						"downloaded": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the secret data that is associated with a secret version was retrieved in a call to the service API.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A v4 UUID identifier.",
						},
						"labels": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Labels that you can use to search for secrets in your instance.Up to 30 labels can be created.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"locks_total": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of locks of the secret.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The human-readable name of your secret.",
						},
						"secret_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A v4 UUID identifier, or `default` secret group.",
						},
						"secret_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
						},
						"state_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A text representation of the secret state.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
						},
						"versions_total": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of versions of the secret.",
						},
						"signing_algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identifier for the cryptographic algorithm that was used by the issuing certificate authority to sign a certificate.",
						},
						"alt_names": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "With the Subject Alternative Name field, you can specify additional host names to be protected by a single SSL certificate.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"common_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Common Name (AKA CN) represents the server name protected by the SSL certificate.",
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date a secret is expired. The date format follows RFC 3339.",
						},
						"intermediate_included": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the certificate was imported with an associated intermediate certificate.",
						},
						"issuer": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
						},
						"key_algorithm": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The identifier for the cryptographic algorithm used to generate the public key that is associated with the certificate.",
						},
						"private_key_included": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the certificate was imported with an associated private key.",
						},
						"serial_number": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique serial number that was assigned to a certificate by the issuing certificate authority.",
						},
						"validity": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The date and time that the certificate validity period begins and ends.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"not_before": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date-time format follows RFC 3339.",
									},
									"not_after": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date-time format follows RFC 3339.",
									},
								},
							},
						},
						"issuance_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Issuance information that is associated with your certificate.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_rotated": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the issued certificate is configured with an automatic rotation policy.",
									},
									"challenges": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The set of challenges. It is returned only when ordering public certificates by using manual DNS configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"domain": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The challenge domain.",
												},
												"expiration": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The challenge expiration date. The date format follows RFC 3339.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The challenge status.",
												},
												"txt_record_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The TXT record name.",
												},
												"txt_record_value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The TXT record value.",
												},
											},
										},
									},
									"dns_challenge_validation_time": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date that a user requests to validate DNS challenges for certificates that are ordered with a manual DNS provider. The date format follows RFC 3339.",
									},
									"error_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A code that identifies an issuance error.This field, along with `error_message`, is returned when Secrets Manager successfully processes your request, but the certificate authority is unable to issue a certificate.",
									},
									"error_message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A human-readable message that provides details about the issuance error.",
									},
									"ordered_on": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date when the certificate is ordered. The date format follows RFC 3339.",
									},
									"state": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The secret state that is based on NIST SP 800-57. States are integers and correspond to the `Pre-activation = 0`, `Active = 1`,  `Suspended = 2`, `Deactivated = 3`, and `Destroyed = 5` values.",
									},
									"state_description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A text representation of the secret state.",
									},
								},
							},
						},
						"rotation": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Determines whether Secrets Manager rotates your secrets automatically.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_rotate": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines whether Secrets Manager rotates your secret automatically.Default is `false`. If `auto_rotate` is set to `true` the service rotates your secret based on the defined interval.",
									},
									"interval": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The length of the secret rotation time interval.",
									},
									"unit": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The units for the secret rotation time interval.",
									},
									"rotate_keys": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Determines whether Secrets Manager rotates the private key for your public certificate automatically.Default is `false`. If it is set to `true`, the service generates and stores a new private key for your rotated certificate.",
									},
								},
							},
						},
						"bundle_certs": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether the issued certificate is bundled with intermediate certificates.",
						},
						"ca": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name that is assigned to the certificate authority configuration.",
						},
						"dns": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name that is assigned to the DNS provider configuration.",
						},
						"next_rotation_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that have an existing rotation policy.",
						},
						"ttl": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time-to-live (TTL) or lease duration to assign to generated credentials.For `iam_credentials` secrets, the TTL defines for how long each generated API key remains valid. The value can be either an integer that specifies the number of seconds, or the string representation of a duration, such as `120m` or `24h`.Minimum duration is 1 minute. Maximum is 90 days.",
						},
						"access_groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Access Groups that you can use for an `iam_credentials` secret.Up to 10 Access Groups can be used for each secret.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"api_key_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the API key that is generated for this secret.",
						},
						"service_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The service ID under which the API key (see the `api_key` field) is created.If you omit this parameter, Secrets Manager generates a new service ID for your secret at its creation and adds it to the access groups that you assign.Optionally, you can use this field to provide your own service ID if you prefer to manage its access directly or retain the service ID after your secret expires, is rotated, or deleted. If you provide a service ID, do not include the `access_groups` parameter.",
						},
						"service_id_is_static": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether an `iam_credentials` secret was created with a static service ID.If it is set to `true`, the service ID for the secret was provided by the user at secret creation. If it is set to `false`, the service ID was generated by Secrets Manager.",
						},
						"reuse_api_key": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether to use the same service ID and API key for future read operations on an`iam_credentials` secret. The value is always `true` for IAM credentials secrets managed by Terraform.",
						},
						"certificate_authority": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The intermediate certificate authority that signed this certificate.",
						},
						"certificate_template": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the certificate template.",
						},
						"revocation_time_seconds": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The timestamp of the certificate revocation.",
						},
						"revocation_time_rfc3339": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the certificate was revoked. The date format follows RFC 3339.",
						},
						"source_service": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The properties required for creating the service credentials for the specified source service instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The source service instance identifier.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "A CRN that uniquely identifies a service credentials target.",
												},
											},
										},
									},
									"role": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The service-specific custom role object, CRN role is accepted. Refer to the service’s documentation for supported roles.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN role identifier for creating a service-id.",
												},
											},
										},
									},
									"iam": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The source service IAM data is returned in case IAM credentials where created for this secret.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"apikey": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The IAM apikey metadata for the IAM credentials that were generated.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IAM API key name for the generated service credentials.",
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IAM API key description for the generated service credentials.",
															},
														},
													},
												},
												"role": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The IAM role for the generate service credentials.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"crn": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IAM role CRN assigned to the generated service credentials.",
															},
														},
													},
												},
												"serviceid": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The IAM serviceid for the generated service credentials.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"crn": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IAM Service ID CRN.",
															},
														},
													},
												},
											},
										},
									},
									"resource_key": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The source service resource key data of the generated service credentials.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource key CRN of the generated service credentials.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource key name of the generated service credentials.",
												},
											},
										},
									},
									"parameters": &schema.Schema{
										Type:        schema.TypeMap,
										Computed:    true,
										Description: "The collection of parameters for the service credentials target.",
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

func dataSourceIbmSmSecretsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", SecretsResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	listSecretsOptions := &secretsmanagerv2.ListSecretsOptions{}
	sort, ok := d.GetOk("sort")
	if ok {
		sortStr := sort.(string)
		listSecretsOptions.SetSort(sortStr)
	}

	search, ok := d.GetOk("search")
	if ok {
		searchStr := search.(string)
		listSecretsOptions.SetSearch(searchStr)
	}

	groups, ok := d.GetOk("groups")
	if ok {
		groupsStr := groups.(string)
		if groupsStr != "" {
			groupsList := strings.Split(groupsStr, ",")
			listSecretsOptions.SetGroups(groupsList)
		}
	}

	if _, ok := d.GetOk("secret_types"); ok {
		secretTypes := d.Get("secret_types").([]interface{})
		parsedTypes := make([]string, len(secretTypes))
		for i, v := range secretTypes {
			parsedTypes[i] = fmt.Sprint(v)
		}
		listSecretsOptions.SetSecretTypes(parsedTypes)
	}

	if _, ok := d.GetOk("match_all_labels"); ok {
		labels := d.Get("match_all_labels").([]interface{})
		parsedLabels := make([]string, len(labels))
		for i, v := range labels {
			parsedLabels[i] = fmt.Sprint(v)
		}
		listSecretsOptions.SetMatchAllLabels(parsedLabels)
	}

	var pager *secretsmanagerv2.SecretsPager
	pager, err = secretsManagerClient.NewSecretsPager(listSecretsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", SecretsResourceName), "read")
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] SecretsPager.GetAll() failed %s", err)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SecretsPager.GetAll() failed %s", err), fmt.Sprintf("(Data) %s", SecretsResourceName), "read")
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", region, instanceId))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIbmSmSecretsSecretMetadataToMap(modelItem)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", SecretsResourceName), "read")
			return tfErr.GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", SecretsResourceName), "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secrets", mapSlice); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secrets"), fmt.Sprintf("(Data) %s", SecretsResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("total_count", len(mapSlice)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting total_count"), fmt.Sprintf("(Data) %s", SecretsResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmSmSecretsSecretMetadataToMap(model secretsmanagerv2.SecretMetadataIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.ImportedCertificateMetadata); ok {
		return dataSourceIbmSmSecretsImportedCertificateMetadataToMap(model.(*secretsmanagerv2.ImportedCertificateMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PublicCertificateMetadata); ok {
		return dataSourceIbmSmSecretsPublicCertificateMetadataToMap(model.(*secretsmanagerv2.PublicCertificateMetadata))
	} else if _, ok := model.(*secretsmanagerv2.KVSecretMetadata); ok {
		return dataSourceIbmSmSecretsKVSecretMetadataToMap(model.(*secretsmanagerv2.KVSecretMetadata))
	} else if _, ok := model.(*secretsmanagerv2.UsernamePasswordSecretMetadata); ok {
		return dataSourceIbmSmSecretsUsernamePasswordSecretMetadataToMap(model.(*secretsmanagerv2.UsernamePasswordSecretMetadata))
	} else if _, ok := model.(*secretsmanagerv2.IAMCredentialsSecretMetadata); ok {
		return dataSourceIbmSmSecretsIAMCredentialsSecretMetadataToMap(model.(*secretsmanagerv2.IAMCredentialsSecretMetadata))
	} else if _, ok := model.(*secretsmanagerv2.ArbitrarySecretMetadata); ok {
		return dataSourceIbmSmSecretsArbitrarySecretMetadataToMap(model.(*secretsmanagerv2.ArbitrarySecretMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PrivateCertificateMetadata); ok {
		return dataSourceIbmSmSecretsPrivateCertificateMetadataToMap(model.(*secretsmanagerv2.PrivateCertificateMetadata))
	} else if _, ok := model.(*secretsmanagerv2.ServiceCredentialsSecretMetadata); ok {
		return dataSourceIbmSmSecretsServiceCredentialsSecretMetadataToMap(model.(*secretsmanagerv2.ServiceCredentialsSecretMetadata))
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.SecretMetadataIntf subtype encountered")
	}
}

func dataSourceIbmSmSecretsCertificateValidityToMap(model *secretsmanagerv2.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotBefore != nil {
		modelMap["not_before"] = model.NotBefore.String()
	}
	if model.NotAfter != nil {
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsCertificateIssuanceInfoToMap(model *secretsmanagerv2.CertificateIssuanceInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotated != nil {
		modelMap["auto_rotated"] = *model.AutoRotated
	}
	if model.Challenges != nil {
		challenges := []map[string]interface{}{}
		for _, challengesItem := range model.Challenges {
			challengesItemMap, err := dataSourceIbmSmSecretsChallengeResourceToMap(&challengesItem)
			if err != nil {
				return modelMap, err
			}
			challenges = append(challenges, challengesItemMap)
		}
		modelMap["challenges"] = challenges
	}
	if model.DnsChallengeValidationTime != nil {
		modelMap["dns_challenge_validation_time"] = model.DnsChallengeValidationTime.String()
	}
	if model.ErrorCode != nil {
		modelMap["error_code"] = *model.ErrorCode
	}
	if model.ErrorMessage != nil {
		modelMap["error_message"] = *model.ErrorMessage
	}
	if model.OrderedOn != nil {
		modelMap["ordered_on"] = model.OrderedOn.String()
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsChallengeResourceToMap(model *secretsmanagerv2.ChallengeResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Domain != nil {
		modelMap["domain"] = *model.Domain
	}
	if model.Expiration != nil {
		modelMap["expiration"] = model.Expiration.String()
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.TxtRecordName != nil {
		modelMap["txt_record_name"] = *model.TxtRecordName
	}
	if model.TxtRecordValue != nil {
		modelMap["txt_record_value"] = *model.TxtRecordValue
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsRotationPolicyToMap(model secretsmanagerv2.RotationPolicyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.CommonRotationPolicy); ok {
		return dataSourceIbmSmSecretsCommonRotationPolicyToMap(model.(*secretsmanagerv2.CommonRotationPolicy))
	} else if _, ok := model.(*secretsmanagerv2.PublicCertificateRotationPolicy); ok {
		return dataSourceIbmSmSecretsPublicCertificateRotationPolicyToMap(model.(*secretsmanagerv2.PublicCertificateRotationPolicy))
	} else if _, ok := model.(*secretsmanagerv2.RotationPolicy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv2.RotationPolicy)
		if model.AutoRotate != nil {
			modelMap["auto_rotate"] = *model.AutoRotate
		}
		if model.Interval != nil {
			modelMap["interval"] = *model.Interval
		}
		if model.Unit != nil {
			modelMap["unit"] = *model.Unit
		}
		if model.RotateKeys != nil {
			modelMap["rotate_keys"] = *model.RotateKeys
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.RotationPolicyIntf subtype encountered")
	}
}

func dataSourceIbmSmSecretsCommonRotationPolicyToMap(model *secretsmanagerv2.CommonRotationPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = *model.AutoRotate
	}
	if model.Interval != nil {
		modelMap["interval"] = *model.Interval
	}
	if model.Unit != nil {
		modelMap["unit"] = *model.Unit
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsPublicCertificateRotationPolicyToMap(model *secretsmanagerv2.PublicCertificateRotationPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = *model.AutoRotate
	}
	if model.RotateKeys != nil {
		modelMap["rotate_keys"] = *model.RotateKeys
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsArbitrarySecretMetadataToMap(model *secretsmanagerv2.ArbitrarySecretMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsIAMCredentialsSecretMetadataToMap(model *secretsmanagerv2.IAMCredentialsSecretMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.TTL != nil {
		modelMap["ttl"] = *model.TTL
	}
	if model.AccessGroups != nil {
		modelMap["access_groups"] = model.AccessGroups
	}
	if model.ApiKeyID != nil {
		modelMap["api_key_id"] = *model.ApiKeyID
	}
	if model.ServiceID != nil {
		modelMap["service_id"] = *model.ServiceID
	}
	if model.ServiceIdIsStatic != nil {
		modelMap["service_id_is_static"] = *model.ServiceIdIsStatic
	}
	if model.ReuseApiKey != nil {
		modelMap["reuse_api_key"] = *model.ReuseApiKey
	}
	if model.Rotation != nil {
		rotationMap, err := dataSourceIbmSmSecretsRotationPolicyToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.NextRotationDate != nil {
		modelMap["next_rotation_date"] = model.NextRotationDate.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsImportedCertificateMetadataToMap(model *secretsmanagerv2.ImportedCertificateMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.SigningAlgorithm != nil {
		modelMap["signing_algorithm"] = *model.SigningAlgorithm
	}
	if model.AltNames != nil {
		modelMap["alt_names"] = model.AltNames
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.IntermediateIncluded != nil {
		modelMap["intermediate_included"] = *model.IntermediateIncluded
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = *model.KeyAlgorithm
	}
	if model.PrivateKeyIncluded != nil {
		modelMap["private_key_included"] = *model.PrivateKeyIncluded
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.Validity != nil {
		validityMap, err := dataSourceIbmSmSecretsCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsPublicCertificateMetadataToMap(model *secretsmanagerv2.PublicCertificateMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.SigningAlgorithm != nil {
		modelMap["signing_algorithm"] = *model.SigningAlgorithm
	}
	if model.AltNames != nil {
		modelMap["alt_names"] = model.AltNames
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.IssuanceInfo != nil {
		issuanceInfoMap, err := dataSourceIbmSmSecretsCertificateIssuanceInfoToMap(model.IssuanceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["issuance_info"] = []map[string]interface{}{issuanceInfoMap}
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = *model.KeyAlgorithm
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.Validity != nil {
		validityMap, err := dataSourceIbmSmSecretsCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.Rotation != nil {
		rotationMap, err := dataSourceIbmSmSecretsRotationPolicyToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.BundleCerts != nil {
		modelMap["bundle_certs"] = *model.BundleCerts
	}
	if model.Ca != nil {
		modelMap["ca"] = *model.Ca
	}
	if model.Dns != nil {
		modelMap["dns"] = *model.Dns
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsKVSecretMetadataToMap(model *secretsmanagerv2.KVSecretMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsUsernamePasswordSecretMetadataToMap(model *secretsmanagerv2.UsernamePasswordSecretMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.Rotation != nil {
		rotationMap, err := dataSourceIbmSmSecretsRotationPolicyToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.NextRotationDate != nil {
		modelMap["next_rotation_date"] = model.NextRotationDate.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsPrivateCertificateMetadataToMap(model *secretsmanagerv2.PrivateCertificateMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.SigningAlgorithm != nil {
		modelMap["signing_algorithm"] = *model.SigningAlgorithm
	}
	if model.AltNames != nil {
		modelMap["alt_names"] = model.AltNames
	}
	if model.CertificateAuthority != nil {
		modelMap["certificate_authority"] = *model.CertificateAuthority
	}
	if model.CertificateTemplate != nil {
		modelMap["certificate_template"] = *model.CertificateTemplate
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.KeyAlgorithm != nil {
		modelMap["key_algorithm"] = *model.KeyAlgorithm
	}
	if model.NextRotationDate != nil {
		modelMap["next_rotation_date"] = model.NextRotationDate.String()
	}
	if model.Rotation != nil {
		rotationMap, err := dataSourceIbmSmSecretsRotationPolicyToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.SerialNumber != nil {
		modelMap["serial_number"] = *model.SerialNumber
	}
	if model.Validity != nil {
		validityMap, err := dataSourceIbmSmSecretsCertificateValidityToMap(model.Validity)
		if err != nil {
			return modelMap, err
		}
		modelMap["validity"] = []map[string]interface{}{validityMap}
	}
	if model.RevocationTimeSeconds != nil {
		modelMap["revocation_time_seconds"] = *model.RevocationTimeSeconds
	}
	if model.RevocationTimeRfc3339 != nil {
		modelMap["revocation_time_rfc3339"] = model.RevocationTimeRfc3339.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmSecretsServiceCredentialsSecretMetadataToMap(model *secretsmanagerv2.ServiceCredentialsSecretMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomMetadata != nil {
		customMetadataMap := make(map[string]interface{}, len(model.CustomMetadata))
		for k, v := range model.CustomMetadata {
			customMetadataMap[k] = v
		}
		modelMap["custom_metadata"] = flex.Flatten(customMetadataMap)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.Downloaded != nil {
		modelMap["downloaded"] = *model.Downloaded
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.LocksTotal != nil {
		modelMap["locks_total"] = *model.LocksTotal
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretGroupID != nil {
		modelMap["secret_group_id"] = *model.SecretGroupID
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.StateDescription != nil {
		modelMap["state_description"] = *model.StateDescription
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.VersionsTotal != nil {
		modelMap["versions_total"] = *model.VersionsTotal
	}
	if model.TTL != nil {
		modelMap["ttl"] = *model.TTL
	}
	if model.Rotation != nil {
		rotationMap, err := dataSourceIbmSmSecretsRotationPolicyToMap(model.Rotation)
		if err != nil {
			return modelMap, err
		}
		modelMap["rotation"] = []map[string]interface{}{rotationMap}
	}
	if model.NextRotationDate != nil {
		modelMap["next_rotation_date"] = model.NextRotationDate.String()
	}
	if model.SourceService != nil {
		sourceServiceMap, err := dataSourceIbmSmServiceCredentialsSecretMetadataSourceServiceToMap(model.SourceService)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_service"] = []map[string]interface{}{sourceServiceMap}
	}

	return modelMap, nil
}
