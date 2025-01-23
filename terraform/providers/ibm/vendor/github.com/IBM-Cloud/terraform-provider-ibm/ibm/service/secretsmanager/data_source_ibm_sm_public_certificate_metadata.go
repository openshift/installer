// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmPublicCertificateMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmPublicCertificateMetadataRead,

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the secret.",
			},
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
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
			},
			"key_algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier for the cryptographic algorithm to be used to generate the public key that is associated with the certificate.The algorithm that you select determines the encryption algorithm (`RSA` or `ECDSA`) and key size to be used to generate keys and sign certificates. For longer living certificates, it is recommended to use longer keys to provide more encryption protection. Allowed values:  RSA2048, RSA4096, EC256, EC384.",
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
				Description: "The name of the certificate authority configuration.",
			},
			"dns": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the DNS provider configuration.",
			},
		},
	}
}

func dataSourceIbmSmPublicCertificateMetadataRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getSecretMetadataOptions := &secretsmanagerv2.GetSecretMetadataOptions{}

	secretId := d.Get("secret_id").(string)
	getSecretMetadataOptions.SetID(secretId)

	publicCertificateMetadataIntf, response, err := secretsManagerClient.GetSecretMetadataWithContext(context, getSecretMetadataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetSecretMetadataWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSecretMetadataWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	publicCertificateMetadata := publicCertificateMetadataIntf.(*secretsmanagerv2.PublicCertificateMetadata)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, secretId))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", publicCertificateMetadata.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(publicCertificateMetadata.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", publicCertificateMetadata.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if publicCertificateMetadata.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(publicCertificateMetadata.CustomMetadata))
		for k, v := range publicCertificateMetadata.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", publicCertificateMetadata.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", publicCertificateMetadata.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("locks_total", flex.IntValue(publicCertificateMetadata.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", publicCertificateMetadata.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", publicCertificateMetadata.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", publicCertificateMetadata.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(publicCertificateMetadata.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", publicCertificateMetadata.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(publicCertificateMetadata.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(publicCertificateMetadata.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("signing_algorithm", publicCertificateMetadata.SigningAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting signing_algorithm"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("common_name", publicCertificateMetadata.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(publicCertificateMetadata.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	issuanceInfo := []map[string]interface{}{}
	if publicCertificateMetadata.IssuanceInfo != nil {
		modelMap, err := dataSourceIbmSmPublicCertificateMetadataCertificateIssuanceInfoToMap(publicCertificateMetadata.IssuanceInfo)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		issuanceInfo = append(issuanceInfo, modelMap)
	}
	if err = d.Set("issuance_info", issuanceInfo); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuance_info"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("issuer", publicCertificateMetadata.Issuer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuer"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("key_algorithm", publicCertificateMetadata.KeyAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_algorithm"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("serial_number", publicCertificateMetadata.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	validity := []map[string]interface{}{}
	if publicCertificateMetadata.Validity != nil {
		modelMap, err := dataSourceIbmSmPublicCertificateMetadataCertificateValidityToMap(publicCertificateMetadata.Validity)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		validity = append(validity, modelMap)
	}
	if err = d.Set("validity", validity); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validity"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	rotation := []map[string]interface{}{}
	if publicCertificateMetadata.Rotation != nil {
		modelMap, err := dataSourceIbmSmPublicCertificateMetadataRotationPolicyToMap(publicCertificateMetadata.Rotation)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		rotation = append(rotation, modelMap)
	}
	if err = d.Set("rotation", rotation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rotation"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("bundle_certs", publicCertificateMetadata.BundleCerts); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting bundle_certs"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("ca", publicCertificateMetadata.Ca); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ca"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("dns", publicCertificateMetadata.Dns); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting dns"), fmt.Sprintf("(Data) %s_metadata", PublicCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmSmPublicCertificateMetadataCertificateIssuanceInfoToMap(model *secretsmanagerv2.CertificateIssuanceInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotated != nil {
		modelMap["auto_rotated"] = *model.AutoRotated
	}
	if model.Challenges != nil {
		challenges := []map[string]interface{}{}
		for _, challengesItem := range model.Challenges {
			challengesItemMap, err := dataSourceIbmSmPublicCertificateMetadataChallengeResourceToMap(&challengesItem)
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

func dataSourceIbmSmPublicCertificateMetadataChallengeResourceToMap(model *secretsmanagerv2.ChallengeResource) (map[string]interface{}, error) {
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

func dataSourceIbmSmPublicCertificateMetadataCertificateValidityToMap(model *secretsmanagerv2.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotBefore != nil {
		modelMap["not_before"] = model.NotBefore.String()
	}
	if model.NotAfter != nil {
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmPublicCertificateMetadataRotationPolicyToMap(model secretsmanagerv2.RotationPolicyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.CommonRotationPolicy); ok {
		return dataSourceIbmSmPublicCertificateMetadataCommonRotationPolicyToMap(model.(*secretsmanagerv2.CommonRotationPolicy))
	} else if _, ok := model.(*secretsmanagerv2.PublicCertificateRotationPolicy); ok {
		return dataSourceIbmSmPublicCertificateMetadataPublicCertificateRotationPolicyToMap(model.(*secretsmanagerv2.PublicCertificateRotationPolicy))
	} else if _, ok := model.(*secretsmanagerv2.RotationPolicy); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv2.RotationPolicy)
		if model.AutoRotate != nil {
			modelMap["auto_rotate"] = *model.AutoRotate
		}
		if model.RotateKeys != nil {
			modelMap["rotate_keys"] = *model.RotateKeys
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.RotationPolicyIntf subtype encountered")
	}
}

func dataSourceIbmSmPublicCertificateMetadataCommonRotationPolicyToMap(model *secretsmanagerv2.CommonRotationPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = *model.AutoRotate
	}
	return modelMap, nil
}

func dataSourceIbmSmPublicCertificateMetadataPublicCertificateRotationPolicyToMap(model *secretsmanagerv2.PublicCertificateRotationPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoRotate != nil {
		modelMap["auto_rotate"] = *model.AutoRotate
	}
	if model.RotateKeys != nil {
		modelMap["rotate_keys"] = *model.RotateKeys
	}
	return modelMap, nil
}
