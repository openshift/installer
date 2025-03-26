// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmPrivateCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmPrivateCertificateRead,

		Schema: map[string]*schema.Schema{
			"secret_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"secret_id", "name"},
				Description:  "The ID of the secret.",
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"secret_id", "name"},
				RequiredWith: []string{"secret_group_name"},
				Description:  "The human-readable name of your secret.",
			},
			"secret_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A v4 UUID identifier, or `default` secret group.",
			},
			"secret_group_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"name"},
				Description:  "The human-readable name of your secret group.",
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
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Common Name (AKA CN) represents the server name that is protected by the SSL certificate.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
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
			"next_rotation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date that the secret is scheduled for automatic rotation.The service automatically creates a new version of the secret on its next rotation date. This field exists only for secrets that have an existing rotation policy.",
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
					},
				},
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
			"certificate": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The PEM-encoded contents of your certificate.",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "(Optional) The PEM-encoded private key to associate with the certificate.",
			},
			"issuing_ca": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The PEM-encoded certificate of the certificate authority that signed and issued this certificate.",
			},
			"ca_chain": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "The chain of certificate authorities that are associated with the certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceIbmSmPrivateCertificateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	privateCertificateIntf, region, instanceId, diagError := getSecretByIdOrByName(context, d, meta, PrivateCertSecretType, PrivateCertSecretResourceName)
	if diagError != nil {
		return diagError
	}

	privateCertificate := privateCertificateIntf.(*secretsmanagerv2.PrivateCertificate)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *privateCertificate.ID))

	var err error
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", privateCertificate.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(privateCertificate.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", privateCertificate.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crn"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if privateCertificate.CustomMetadata != nil {
		convertedMap := make(map[string]interface{}, len(privateCertificate.CustomMetadata))
		for k, v := range privateCertificate.CustomMetadata {
			convertedMap[k] = v
		}

		if err = d.Set("custom_metadata", flex.Flatten(convertedMap)); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting custom_metadata"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("description", privateCertificate.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting description"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("downloaded", privateCertificate.Downloaded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting downloaded"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("locks_total", flex.IntValue(privateCertificate.LocksTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locks_total"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("name", privateCertificate.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_group_id", privateCertificate.SecretGroupID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_group_id"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", privateCertificate.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state", flex.IntValue(privateCertificate.State)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("state_description", privateCertificate.StateDescription); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting state_description"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(privateCertificate.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("versions_total", flex.IntValue(privateCertificate.VersionsTotal)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting versions_total"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("signing_algorithm", privateCertificate.SigningAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting signing_algorithm"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("certificate_authority", privateCertificate.CertificateAuthority); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate_authority"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("certificate_template", privateCertificate.CertificateTemplate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate_template"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("common_name", privateCertificate.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(privateCertificate.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("issuer", privateCertificate.Issuer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuer"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("key_algorithm", privateCertificate.KeyAlgorithm); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_algorithm"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("next_rotation_date", DateTimeToRFC3339(privateCertificate.NextRotationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting next_rotation_date"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	rotation := []map[string]interface{}{}
	if privateCertificate.Rotation != nil {
		modelMap, err := dataSourceIbmSmPrivateCertificateRotationPolicyToMap(privateCertificate.Rotation)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		rotation = append(rotation, modelMap)
	}
	if err = d.Set("rotation", rotation); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting rotation"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("serial_number", privateCertificate.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	validity := []map[string]interface{}{}
	if privateCertificate.Validity != nil {
		modelMap, err := dataSourceIbmSmPrivateCertificateCertificateValidityToMap(privateCertificate.Validity)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
			return tfErr.GetDiag()
		}
		validity = append(validity, modelMap)
	}
	if err = d.Set("validity", validity); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting validity"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("revocation_time_seconds", flex.IntValue(privateCertificate.RevocationTimeSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting revocation_time_seconds"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("revocation_time_rfc3339", DateTimeToRFC3339(privateCertificate.RevocationTimeRfc3339)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting revocation_time_rfc3339"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("certificate", privateCertificate.Certificate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("private_key", privateCertificate.PrivateKey); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("issuing_ca", privateCertificate.IssuingCa); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuing_ca"), fmt.Sprintf("(Data) %s", PrivateCertSecretResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmSmPrivateCertificateRotationPolicyToMap(model secretsmanagerv2.RotationPolicyIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.CommonRotationPolicy); ok {
		return dataSourceIbmSmPrivateCertificateCommonRotationPolicyToMap(model.(*secretsmanagerv2.CommonRotationPolicy))
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
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.RotationPolicyIntf subtype encountered")
	}
}

func dataSourceIbmSmPrivateCertificateCommonRotationPolicyToMap(model *secretsmanagerv2.CommonRotationPolicy) (map[string]interface{}, error) {
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

func dataSourceIbmSmPrivateCertificateCertificateValidityToMap(model *secretsmanagerv2.CertificateValidity) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NotBefore != nil {
		modelMap["not_before"] = model.NotBefore.String()
	}
	if model.NotAfter != nil {
		modelMap["not_after"] = model.NotAfter.String()
	}
	return modelMap, nil
}
