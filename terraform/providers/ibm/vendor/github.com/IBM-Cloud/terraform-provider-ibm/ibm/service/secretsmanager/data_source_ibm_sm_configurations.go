// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func DataSourceIbmSmConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmConfigurationsRead,

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
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources in a collection.",
			},
			"configurations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A collection of configuration metadata.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration type.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique name of your configuration.",
						},
						"secret_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
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
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date when a resource was recently modified. The date format follows RFC 3339.",
						},
						"lets_encrypt_environment": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration of the Let's Encrypt CA environment.",
						},
						"lets_encrypt_preferred_chain": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Prefer the chain with an issuer matching this Subject Common Name.",
						},
						"common_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Common Name (AKA CN) represents the server name that is protected by the SSL certificate.",
						},
						"crl_distribution_points_encoded": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are issued by this certificate authority.",
						},
						"expiration_date": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date a secret is expired. The date format follows RFC 3339.",
						},
						"key_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of private key to generate.",
						},
						"key_bits": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of bits to use to generate the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the certificate authority. The status of a root certificate authority is either `configured` or `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,`signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.",
						},
						"issuer": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
						},
						"signing_method": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The signing method to use with this certificate authority to generate private certificates.You can choose between internal or externally signed options. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).",
						},
						"certificate_authority": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the intermediate certificate authority.",
						},
						"crypto_key": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The data that is associated with a cryptographic key.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of a PKCS#11 key to use. If the key does not exist and generation is enabled, this ID is given to the generated key. If the key exists, and generation is disabled, then this ID is used to look up the key. This value or the crypto key label must be specified.",
									},
									"label": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The label of the key to use. If the key does not exist and generation is enabled, this field is the label that is given to the generated key. If the key exists, and generation is disabled, then this label is used to look up the key. This value or the crypto key ID must be specified.",
									},
									"allow_generate_key": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The indication of whether a new key is generated by the crypto provider if the given key name cannot be found.",
									},
									"provider": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The data that is associated with a cryptographic provider.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The type of cryptographic provider.",
												},
												"instance_crn": &schema.Schema{
													Description: "The HPCS instance CRN.",
													Computed:    true,
													Type:        schema.TypeString,
												},
												"pin_iam_credentials_secret_id": &schema.Schema{
													Description: "The secret Id of iam credentials with api key to access HPCS instance.",
													Computed:    true,
													Type:        schema.TypeString,
												},
												"private_keystore_id": &schema.Schema{
													Description: "The HPCS private key store space id.",
													Computed:    true,
													Type:        schema.TypeString,
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
		},
	}
}

func dataSourceIbmSmConfigurationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", ConfigurationsResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	listConfigurationsOptions := &secretsmanagerv2.ListConfigurationsOptions{}
	sort, ok := d.GetOk("sort")
	if ok {
		sortStr := sort.(string)
		listConfigurationsOptions.SetSort(sortStr)
	}
	search, ok := d.GetOk("search")
	if ok {
		searchStr := search.(string)
		listConfigurationsOptions.SetSearch(searchStr)
	}
	if _, ok := d.GetOk("secret_types"); ok {
		secretTypes := d.Get("secret_types").([]interface{})
		parsedTypes := make([]string, len(secretTypes))
		for i, v := range secretTypes {
			parsedTypes[i] = fmt.Sprint(v)
		}
		listConfigurationsOptions.SetSecretTypes(parsedTypes)
	}

	var pager *secretsmanagerv2.ConfigurationsPager
	pager, err = secretsManagerClient.NewConfigurationsPager(listConfigurationsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", ConfigurationsResourceName), "read")
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] ConfigurationsPager.GetAll() failed %s", err)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ConfigurationsPager.GetAll() %s", err), fmt.Sprintf("(Data) %s", ConfigurationsResourceName), "read")
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmSmConfigurationsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIbmSmConfigurationsConfigurationMetadataToMap(modelItem)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", ConfigurationsResourceName), "read")
			return tfErr.GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("configurations", mapSlice); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting configurations"), fmt.Sprintf("(Data) %s", ConfigurationsResourceName), "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("total_count", len(mapSlice)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting total_count"), fmt.Sprintf("(Data) %s", ConfigurationsResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}

// dataSourceIbmSmConfigurationsID returns a reasonable ID for the list.
func dataSourceIbmSmConfigurationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmSmConfigurationsConfigurationMetadataToMap(model secretsmanagerv2.ConfigurationMetadataIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.IAMCredentialsConfigurationMetadata); ok {
		return dataSourceIbmSmConfigurationsIAMCredentialsConfigurationMetadataToMap(model.(*secretsmanagerv2.IAMCredentialsConfigurationMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PublicCertificateConfigurationCALetsEncryptMetadata); ok {
		return dataSourceIbmSmConfigurationsPublicCertificateConfigurationCALetsEncryptMetadataToMap(model.(*secretsmanagerv2.PublicCertificateConfigurationCALetsEncryptMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServicesMetadata); ok {
		return dataSourceIbmSmConfigurationsPublicCertificateConfigurationDNSCloudInternetServicesMetadataToMap(model.(*secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServicesMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructureMetadata); ok {
		return dataSourceIbmSmConfigurationsPublicCertificateConfigurationDNSClassicInfrastructureMetadataToMap(model.(*secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructureMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PrivateCertificateConfigurationRootCAMetadata); ok {
		return dataSourceIbmSmConfigurationsPrivateCertificateConfigurationRootCAMetadataToMap(model.(*secretsmanagerv2.PrivateCertificateConfigurationRootCAMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PrivateCertificateConfigurationIntermediateCAMetadata); ok {
		return dataSourceIbmSmConfigurationsPrivateCertificateConfigurationIntermediateCAMetadataToMap(model.(*secretsmanagerv2.PrivateCertificateConfigurationIntermediateCAMetadata))
	} else if _, ok := model.(*secretsmanagerv2.PrivateCertificateConfigurationTemplateMetadata); ok {
		return dataSourceIbmSmConfigurationsPrivateCertificateConfigurationTemplateMetadataToMap(model.(*secretsmanagerv2.PrivateCertificateConfigurationTemplateMetadata))
	} else if _, ok := model.(*secretsmanagerv2.ConfigurationMetadata); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv2.ConfigurationMetadata)
		if model.ConfigType != nil {
			modelMap["config_type"] = *model.ConfigType
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.SecretType != nil {
			modelMap["secret_type"] = *model.SecretType
		}
		if model.CreatedBy != nil {
			modelMap["created_by"] = *model.CreatedBy
		}
		if model.CreatedAt != nil {
			modelMap["created_at"] = model.CreatedAt.String()
		}
		if model.UpdatedAt != nil {
			modelMap["updated_at"] = model.UpdatedAt.String()
		}
		if model.LetsEncryptEnvironment != nil {
			modelMap["lets_encrypt_environment"] = *model.LetsEncryptEnvironment
		}
		if model.LetsEncryptPreferredChain != nil {
			modelMap["lets_encrypt_preferred_chain"] = *model.LetsEncryptPreferredChain
		}
		if model.CommonName != nil {
			modelMap["common_name"] = *model.CommonName
		}
		if model.CrlDistributionPointsEncoded != nil {
			modelMap["crl_distribution_points_encoded"] = *model.CrlDistributionPointsEncoded
		}
		if model.ExpirationDate != nil {
			modelMap["expiration_date"] = model.ExpirationDate.String()
		}
		if model.KeyType != nil {
			modelMap["key_type"] = *model.KeyType
		}
		if model.KeyBits != nil {
			modelMap["key_bits"] = *model.KeyBits
		}
		if model.Status != nil {
			modelMap["status"] = *model.Status
		}
		if model.Issuer != nil {
			modelMap["issuer"] = *model.Issuer
		}
		if model.SigningMethod != nil {
			modelMap["signing_method"] = *model.SigningMethod
		}
		if model.CertificateAuthority != nil {
			modelMap["certificate_authority"] = *model.CertificateAuthority
		}

		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.ConfigurationMetadataIntf subtype encountered")
	}
}

func dataSourceIbmSmConfigurationsPrivateCertificateConfigurationIntermediateCAMetadataToMap(model *secretsmanagerv2.PrivateCertificateConfigurationIntermediateCAMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConfigType != nil {
		modelMap["config_type"] = *model.ConfigType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.CrlDistributionPointsEncoded != nil {
		modelMap["crl_distribution_points_encoded"] = *model.CrlDistributionPointsEncoded
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.KeyType != nil {
		modelMap["key_type"] = *model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = *model.KeyBits
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Issuer != nil {
		modelMap["issuer"] = *model.Issuer
	}
	if model.SigningMethod != nil {
		modelMap["signing_method"] = *model.SigningMethod
	}
	if model.CryptoKey != nil {
		cryptoModelMap, err := resourceIbmSmPrivateCertificateConfigurationCryptoKeyToMap(model.CryptoKey)
		if err != nil {
			return modelMap, err
		}
		modelMap["crypto_key"] = []map[string]interface{}{cryptoModelMap}
	}
	return modelMap, nil
}

func dataSourceIbmSmConfigurationsIAMCredentialsConfigurationMetadataToMap(model *secretsmanagerv2.IAMCredentialsConfigurationMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConfigType != nil {
		modelMap["config_type"] = *model.ConfigType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmConfigurationsPrivateCertificateConfigurationRootCAMetadataToMap(model *secretsmanagerv2.PrivateCertificateConfigurationRootCAMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConfigType != nil {
		modelMap["config_type"] = *model.ConfigType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.CommonName != nil {
		modelMap["common_name"] = *model.CommonName
	}
	if model.CrlDistributionPointsEncoded != nil {
		modelMap["crl_distribution_points_encoded"] = *model.CrlDistributionPointsEncoded
	}
	if model.ExpirationDate != nil {
		modelMap["expiration_date"] = model.ExpirationDate.String()
	}
	if model.KeyType != nil {
		modelMap["key_type"] = *model.KeyType
	}
	if model.KeyBits != nil {
		modelMap["key_bits"] = *model.KeyBits
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.CryptoKey != nil {
		cryptoModelMap, err := resourceIbmSmPrivateCertificateConfigurationCryptoKeyToMap(model.CryptoKey)
		if err != nil {
			return modelMap, err
		}
		modelMap["crypto_key"] = []map[string]interface{}{cryptoModelMap}
	}
	return modelMap, nil
}

func dataSourceIbmSmConfigurationsPublicCertificateConfigurationDNSClassicInfrastructureMetadataToMap(model *secretsmanagerv2.PublicCertificateConfigurationDNSClassicInfrastructureMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConfigType != nil {
		modelMap["config_type"] = *model.ConfigType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmConfigurationsPublicCertificateConfigurationCALetsEncryptMetadataToMap(model *secretsmanagerv2.PublicCertificateConfigurationCALetsEncryptMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConfigType != nil {
		modelMap["config_type"] = *model.ConfigType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.LetsEncryptEnvironment != nil {
		modelMap["lets_encrypt_environment"] = *model.LetsEncryptEnvironment
	}
	if model.LetsEncryptPreferredChain != nil {
		modelMap["lets_encrypt_preferred_chain"] = *model.LetsEncryptPreferredChain
	}
	return modelMap, nil
}

func dataSourceIbmSmConfigurationsPublicCertificateConfigurationDNSCloudInternetServicesMetadataToMap(model *secretsmanagerv2.PublicCertificateConfigurationDNSCloudInternetServicesMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConfigType != nil {
		modelMap["config_type"] = *model.ConfigType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSmConfigurationsPrivateCertificateConfigurationTemplateMetadataToMap(model *secretsmanagerv2.PrivateCertificateConfigurationTemplateMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConfigType != nil {
		modelMap["config_type"] = *model.ConfigType
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SecretType != nil {
		modelMap["secret_type"] = *model.SecretType
	}
	if model.CreatedBy != nil {
		modelMap["created_by"] = *model.CreatedBy
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.UpdatedAt != nil {
		modelMap["updated_at"] = model.UpdatedAt.String()
	}
	if model.CertificateAuthority != nil {
		modelMap["certificate_authority"] = *model.CertificateAuthority
	}
	return modelMap, nil
}
