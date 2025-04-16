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

func DataSourceIbmSmPrivateCertificateConfigurationIntermediateCA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmPrivateCertificateConfigurationIntermediateCARead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the configuration.",
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Th configuration type.",
			},
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"max_ttl_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum time-to-live (TTL) for certificates that are created by this CA in seconds.",
			},
			"signing_method": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The signing method to use with this certificate authority to generate private certificates.You can choose between internal or externally signed options. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).",
			},
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
			},
			"crl_expiry_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The time until the certificate revocation list (CRL) expires, in seconds.",
			},
			"crl_disable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Disables or enables certificate revocation list (CRL) building.If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is enabled, it will rebuild the CRL.",
			},
			"crl_distribution_points_encoded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are issued by this certificate authority.",
			},
			"issuing_certificates_urls_encoded": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this certificate authority.",
			},
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Common Name (AKA CN) represents the server name that is protected by the SSL certificate.",
			},
			"alt_names": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "With the Subject Alternative Name field, you can specify additional host names to be protected by a single SSL certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_sans": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
			},
			"uri_sans": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
			},
			"other_sans": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"format": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The format of the returned data.",
			},
			"private_key_format": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The format of the generated private key.",
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
			"exclude_cn_from_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs).If the common name set to `true`, it is not included in DNS or Email SANs if they apply. This field can be useful if the common name is a human-readable identifier, instead of a hostname or an email address.",
			},
			"ou": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"organization": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"country": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Country (C) values to define in the subject field of the resulting certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"locality": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"province": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"street_address": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The street address values to define in the subject field of the resulting certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"postal_code": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The postal code values to define in the subject field of the resulting certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the certificate authority. The status of a root certificate authority is either `configured` or `expired`. For intermediate certificate authorities, possible statuses include `signing_required`,`signed_certificate_required`, `certificate_template_required`, `configured`, `expired` or `revoked`.",
			},
			"expiration_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date a secret is expired. The date format follows RFC 3339.",
			},
			"data": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The configuration data of your Private Certificate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"csr": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate signing request.",
						},
						"private_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "(Optional) The PEM-encoded private key to associate with the certificate.",
						},
						"private_key_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of private key to generate.",
						},
						"expiration": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The certificate expiration time.",
						},
						"certificate": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The PEM-encoded contents of your certificate.",
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
				},
			},
		},
	}
}

func dataSourceIbmSmPrivateCertificateConfigurationIntermediateCARead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(d.Get("name").(string))

	configurationIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigurationWithContext failed %s\n%s", err, response))
	}
	privateCertificateConfigurationIntermediateCA := configurationIntf.(*secretsmanagerv2.PrivateCertificateConfigurationIntermediateCA)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *getConfigurationOptions.Name))

	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if err = d.Set("config_type", privateCertificateConfigurationIntermediateCA.ConfigType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config_type: %s", err))
	}

	if err = d.Set("secret_type", privateCertificateConfigurationIntermediateCA.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}

	if err = d.Set("max_ttl_seconds", flex.IntValue(privateCertificateConfigurationIntermediateCA.MaxTtlSeconds)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting max_ttl_seconds: %s", err))
	}

	if err = d.Set("signing_method", privateCertificateConfigurationIntermediateCA.SigningMethod); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting signing_method: %s", err))
	}

	if err = d.Set("issuer", privateCertificateConfigurationIntermediateCA.Issuer); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuer: %s", err))
	}

	if err = d.Set("crl_expiry_seconds", flex.IntValue(privateCertificateConfigurationIntermediateCA.CrlExpirySeconds)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crl_expiry_seconds: %s", err))
	}

	if err = d.Set("crl_disable", privateCertificateConfigurationIntermediateCA.CrlDisable); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crl_disable: %s", err))
	}

	if err = d.Set("crl_distribution_points_encoded", privateCertificateConfigurationIntermediateCA.CrlDistributionPointsEncoded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crl_distribution_points_encoded: %s", err))
	}

	if err = d.Set("issuing_certificates_urls_encoded", privateCertificateConfigurationIntermediateCA.IssuingCertificatesUrlsEncoded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuing_certificates_urls_encoded: %s", err))
	}

	if err = d.Set("common_name", privateCertificateConfigurationIntermediateCA.CommonName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting common_name: %s", err))
	}

	if err = d.Set("ip_sans", privateCertificateConfigurationIntermediateCA.IpSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ip_sans: %s", err))
	}

	if err = d.Set("uri_sans", privateCertificateConfigurationIntermediateCA.UriSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting uri_sans: %s", err))
	}

	if err = d.Set("format", privateCertificateConfigurationIntermediateCA.Format); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting format: %s", err))
	}

	if err = d.Set("private_key_format", privateCertificateConfigurationIntermediateCA.PrivateKeyFormat); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_key_format: %s", err))
	}

	if err = d.Set("key_type", privateCertificateConfigurationIntermediateCA.KeyType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_type: %s", err))
	}

	if err = d.Set("key_bits", flex.IntValue(privateCertificateConfigurationIntermediateCA.KeyBits)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_bits: %s", err))
	}

	if err = d.Set("exclude_cn_from_sans", privateCertificateConfigurationIntermediateCA.ExcludeCnFromSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting exclude_cn_from_sans: %s", err))
	}

	if err = d.Set("serial_number", privateCertificateConfigurationIntermediateCA.SerialNumber); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting serial_number: %s", err))
	}

	if err = d.Set("status", privateCertificateConfigurationIntermediateCA.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(privateCertificateConfigurationIntermediateCA.ExpirationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiration_date: %s", err))
	}

	data := []map[string]interface{}{}
	if privateCertificateConfigurationIntermediateCA.Data != nil {
		modelMap, err := dataSourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateCADataToMap(privateCertificateConfigurationIntermediateCA.Data)
		if err != nil {
			return diag.FromErr(err)
		}
		data = append(data, modelMap)
	}
	if err = d.Set("data", data); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting data %s", err))
	}

	return nil
}

func dataSourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateCADataToMap(model secretsmanagerv2.PrivateCertificateCADataIntf) (map[string]interface{}, error) {
	if _, ok := model.(*secretsmanagerv2.PrivateCertificateConfigurationIntermediateCACSR); ok {
		return dataSourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateConfigurationIntermediateCACSRToMap(model.(*secretsmanagerv2.PrivateCertificateConfigurationIntermediateCACSR))
	} else if _, ok := model.(*secretsmanagerv2.PrivateCertificateConfigurationCACertificate); ok {
		return dataSourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateConfigurationCACertificateToMap(model.(*secretsmanagerv2.PrivateCertificateConfigurationCACertificate))
	} else if _, ok := model.(*secretsmanagerv2.PrivateCertificateCAData); ok {
		modelMap := make(map[string]interface{})
		model := model.(*secretsmanagerv2.PrivateCertificateCAData)
		if model.Csr != nil {
			modelMap["csr"] = *model.Csr
		}
		if model.PrivateKey != nil {
			modelMap["private_key"] = *model.PrivateKey
		}
		if model.PrivateKeyType != nil {
			modelMap["private_key_type"] = *model.PrivateKeyType
		}
		if model.Expiration != nil {
			modelMap["expiration"] = *model.Expiration
		}
		if model.Certificate != nil {
			modelMap["certificate"] = *model.Certificate
		}
		if model.IssuingCa != nil {
			modelMap["issuing_ca"] = *model.IssuingCa
		}
		if model.CaChain != nil {
			modelMap["ca_chain"] = model.CaChain
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized secretsmanagerv2.PrivateCertificateCADataIntf subtype encountered")
	}
}

func dataSourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateConfigurationIntermediateCACSRToMap(model *secretsmanagerv2.PrivateCertificateConfigurationIntermediateCACSR) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Csr != nil {
		modelMap["csr"] = *model.Csr
	}
	if model.PrivateKey != nil {
		modelMap["private_key"] = *model.PrivateKey
	}
	if model.PrivateKeyType != nil {
		modelMap["private_key_type"] = *model.PrivateKeyType
	}
	if model.Expiration != nil {
		modelMap["expiration"] = *model.Expiration
	}
	return modelMap, nil
}

func dataSourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateConfigurationCACertificateToMap(model *secretsmanagerv2.PrivateCertificateConfigurationCACertificate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Certificate != nil {
		modelMap["certificate"] = *model.Certificate
	}
	if model.IssuingCa != nil {
		modelMap["issuing_ca"] = *model.IssuingCa
	}
	if model.CaChain != nil {
		modelMap["ca_chain"] = model.CaChain
	}
	if model.Expiration != nil {
		modelMap["expiration"] = *model.Expiration
	}
	return modelMap, nil
}
