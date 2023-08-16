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

func DataSourceIbmSmPrivateCertificateConfigurationRootCA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmPrivateCertificateConfigurationRootCARead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the configuration.",
			},
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration type.",
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
			"max_ttl_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum time-to-live (TTL) for certificates that are created by this CA in seconds.",
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
			"ttl_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The requested time-to-live (TTL) for certificates that are created by this CA. This field's value cannot be longer than the `max_ttl` limit.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).",
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
			"max_path_length": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum path length to encode in the generated certificate. `-1` means no limit.If the signing certificate has a maximum path length set, the path length is set to one less than that of the signing certificate. A limit of `0` means a literal path length of zero.",
			},
			"exclude_cn_from_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs).If the common name set to `true`, it is not included in DNS or Email SANs if they apply. This field can be useful if the common name is a human-readable identifier, instead of a hostname or an email address.",
			},
			"permitted_dns_domains": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Sensitive:   true,
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

func dataSourceIbmSmPrivateCertificateConfigurationRootCARead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		return diag.FromErr(err)
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(d.Get("name").(string))

	privateCertificateConfigurationRootCAIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigurationWithContext failed %s\n%s", err, response))
	}

	privateCertificateConfigurationRootCA := privateCertificateConfigurationRootCAIntf.(*secretsmanagerv2.PrivateCertificateConfigurationRootCA)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *getConfigurationOptions.Name))

	if err = d.Set("region", region); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region: %s", err))
	}

	if err = d.Set("config_type", privateCertificateConfigurationRootCA.ConfigType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config_type: %s", err))
	}

	if err = d.Set("secret_type", privateCertificateConfigurationRootCA.SecretType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_type: %s", err))
	}

	if err = d.Set("created_by", privateCertificateConfigurationRootCA.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}

	if err = d.Set("created_at", DateTimeToRFC3339(privateCertificateConfigurationRootCA.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(privateCertificateConfigurationRootCA.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	if err = d.Set("max_ttl_seconds", flex.IntValue(privateCertificateConfigurationRootCA.MaxTtlSeconds)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting max_ttl_seconds: %s", err))
	}

	if err = d.Set("crl_expiry_seconds", flex.IntValue(privateCertificateConfigurationRootCA.CrlExpirySeconds)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crl_expiry_seconds: %s", err))
	}

	if err = d.Set("crl_disable", privateCertificateConfigurationRootCA.CrlDisable); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crl_disable: %s", err))
	}

	if err = d.Set("crl_distribution_points_encoded", privateCertificateConfigurationRootCA.CrlDistributionPointsEncoded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crl_distribution_points_encoded: %s", err))
	}

	if err = d.Set("issuing_certificates_urls_encoded", privateCertificateConfigurationRootCA.IssuingCertificatesUrlsEncoded); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting issuing_certificates_urls_encoded: %s", err))
	}

	if err = d.Set("common_name", privateCertificateConfigurationRootCA.CommonName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting common_name: %s", err))
	}

	if privateCertificateConfigurationRootCA.AltNames != nil {
		if err = d.Set("alt_names", privateCertificateConfigurationRootCA.AltNames); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting alt_names: %s", err))
		}
	}

	if err = d.Set("ip_sans", privateCertificateConfigurationRootCA.IpSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ip_sans: %s", err))
	}

	if err = d.Set("uri_sans", privateCertificateConfigurationRootCA.UriSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting uri_sans: %s", err))
	}
	if privateCertificateConfigurationRootCA.OtherSans != nil {
		if err = d.Set("other_sans", privateCertificateConfigurationRootCA.OtherSans); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting other_sans: %s", err))
		}
	}
	if err = d.Set("ttl_seconds", flex.IntValue(privateCertificateConfigurationRootCA.TtlSeconds)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ttl_seconds: %s", err))
	}

	if err = d.Set("format", privateCertificateConfigurationRootCA.Format); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting format: %s", err))
	}

	if err = d.Set("private_key_format", privateCertificateConfigurationRootCA.PrivateKeyFormat); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting private_key_format: %s", err))
	}

	if err = d.Set("key_type", privateCertificateConfigurationRootCA.KeyType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_type: %s", err))
	}

	if err = d.Set("key_bits", flex.IntValue(privateCertificateConfigurationRootCA.KeyBits)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting key_bits: %s", err))
	}

	if err = d.Set("max_path_length", flex.IntValue(privateCertificateConfigurationRootCA.MaxPathLength)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting max_path_length: %s", err))
	}

	if err = d.Set("exclude_cn_from_sans", privateCertificateConfigurationRootCA.ExcludeCnFromSans); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting exclude_cn_from_sans: %s", err))
	}
	if privateCertificateConfigurationRootCA.PermittedDnsDomains != nil {
		if err = d.Set("permitted_dns_domains", privateCertificateConfigurationRootCA.PermittedDnsDomains); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting permitted_dns_domains: %s", err))
		}
	}
	if privateCertificateConfigurationRootCA.Ou != nil {
		if err = d.Set("ou", privateCertificateConfigurationRootCA.Ou); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting ou: %s", err))
		}
	}
	if privateCertificateConfigurationRootCA.Organization != nil {
		if err = d.Set("organization", privateCertificateConfigurationRootCA.Organization); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting organization: %s", err))
		}
	}
	if privateCertificateConfigurationRootCA.Country != nil {
		if err = d.Set("country", privateCertificateConfigurationRootCA.Country); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting country: %s", err))
		}
	}
	if privateCertificateConfigurationRootCA.Locality != nil {
		if err = d.Set("locality", privateCertificateConfigurationRootCA.Locality); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting locality: %s", err))
		}
	}
	if privateCertificateConfigurationRootCA.Province != nil {
		if err = d.Set("province", privateCertificateConfigurationRootCA.Province); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting province: %s", err))
		}
	}
	if privateCertificateConfigurationRootCA.StreetAddress != nil {
		if err = d.Set("street_address", privateCertificateConfigurationRootCA.StreetAddress); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting street_address: %s", err))
		}
	}
	if privateCertificateConfigurationRootCA.PostalCode != nil {
		if err = d.Set("postal_code", privateCertificateConfigurationRootCA.PostalCode); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting postal_code: %s", err))
		}
	}
	if err = d.Set("serial_number", privateCertificateConfigurationRootCA.SerialNumber); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting serial_number: %s", err))
	}

	if err = d.Set("status", privateCertificateConfigurationRootCA.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if err = d.Set("expiration_date", DateTimeToRFC3339(privateCertificateConfigurationRootCA.ExpirationDate)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiration_date: %s", err))
	}

	if privateCertificateConfigurationRootCA.Data != nil {
		dataMap, err := dataSourceIbmSmPrivateCertificateConfigurationRootCAPrivateCertificateCADataToMap(privateCertificateConfigurationRootCA.Data)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("data", []map[string]interface{}{dataMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting data: %s", err))
		}
	}

	return nil
}

func dataSourceIbmSmPrivateCertificateConfigurationRootCAPrivateCertificateCADataToMap(modelIntf secretsmanagerv2.PrivateCertificateCADataIntf) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	model := modelIntf.(*secretsmanagerv2.PrivateCertificateCAData)

	if model.Csr != nil {
		modelMap["csr"] = model.Csr
	}
	if model.PrivateKey != nil {
		modelMap["private_key"] = model.PrivateKey
	}
	if model.PrivateKeyType != nil {
		modelMap["private_key_type"] = model.PrivateKeyType
	}
	if model.Expiration != nil {
		modelMap["expiration"] = flex.IntValue(model.Expiration)
	}
	if model.Certificate != nil {
		modelMap["certificate"] = model.Certificate
	}
	if model.IssuingCa != nil {
		modelMap["issuing_ca"] = model.IssuingCa
	}
	if model.CaChain != nil {
		modelMap["ca_chain"] = model.CaChain
	}
	return modelMap, nil
}
