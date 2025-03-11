// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPrivateCertificateConfigurationRootCA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPrivateCertificateConfigurationRootCACreate,
		ReadContext:   resourceIbmSmPrivateCertificateConfigurationRootCARead,
		UpdateContext: resourceIbmSmPrivateCertificateConfigurationRootCAUpdate,
		DeleteContext: resourceIbmSmPrivateCertificateConfigurationRootCADelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"config_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration type.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A human-readable unique name to assign to your configuration.To protect your privacy, do not use personal data, such as your name or location, as an name for your secret.",
			},
			"max_ttl": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
			},
			"crl_expiry": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The time until the certificate revocation list (CRL) expires.The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours. In the API response, this value is returned in seconds (integer).**Note:** The CRL is rotated automatically before it expires.",
			},
			"crl_disable": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Disables or enables certificate revocation list (CRL) building.If CRL building is disabled, a signed but zero-length CRL is returned when downloading the CRL. If CRL building is enabled, it will rebuild the CRL.",
			},
			"crl_distribution_points_encoded": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Determines whether to encode the certificate revocation list (CRL) distribution points in the certificates that are issued by this certificate authority.",
			},
			"issuing_certificates_urls_encoded": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Determines whether to encode the URL of the issuing certificate in the certificates that are issued by this certificate authority.",
			},
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Common Name (AKA CN) represents the server name that is protected by the SSL certificate.",
			},
			"alt_names": &schema.Schema{
				Type:             schema.TypeList,
				Optional:         true,
				ForceNew:         true,
				Description:      "With the Subject Alternative Name field, you can specify additional host names to be protected by a single SSL certificate.",
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: altNamesDiffSuppress,
			},
			"ip_sans": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The IP Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
			},
			"uri_sans": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The URI Subject Alternative Names to define for the CA certificate, in a comma-delimited list.",
			},
			"other_sans": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The requested time-to-live (TTL) for certificates that are created by this CA. This field's value cannot be longer than the `max_ttl` limit.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).",
			},
			"ttl_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The requested time-to-live (TTL) for certificates that are created by this CA. This field's value cannot be longer than the `max_ttl` limit.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).",
			},
			"format": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The format of the returned data.",
			},
			"private_key_format": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "der",
				ForceNew:    true,
				Description: "The format of the generated private key.",
			},
			"key_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The type of private key to generate.",
			},
			"key_bits": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The number of bits to use to generate the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.",
			},
			"max_path_length": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The maximum path length to encode in the generated certificate. `-1` means no limit.If the signing certificate has a maximum path length set, the path length is set to one less than that of the signing certificate. A limit of `0` means a literal path length of zero.",
			},
			"exclude_cn_from_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				ForceNew:    true,
				Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs).If the common name set to `true`, it is not included in DNS or Email SANs if they apply. This field can be useful if the common name is a human-readable identifier, instead of a hostname or an email address.",
			},
			"permitted_dns_domains": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ou": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"organization": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"country": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The Country (C) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"locality": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"province": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"street_address": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The street address values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"postal_code": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The postal code values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"crypto_key": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The data that is associated with a cryptographic key.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The ID of a PKCS#11 key to use. If the key does not exist and generation is enabled, this ID is given to the generated key. If the key exists, and generation is disabled, then this ID is used to look up the key. This value or the crypto key label must be specified.",
						},
						"label": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The label of the key to use. If the key does not exist and generation is enabled, this field is the label that is given to the generated key. If the key exists, and generation is disabled, then this label is used to look up the key. This value or the crypto key ID must be specified.",
						},
						"allow_generate_key": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The indication of whether a new key is generated by the crypto provider if the given key name cannot be found.",
						},
						"provider": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "The data that is associated with a cryptographic provider.",
							MaxItems:    1,
							ForceNew:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "The type of cryptographic provider.",
									},
									"instance_crn": &schema.Schema{
										Description: "The HPCS instance CRN.",
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Type:        schema.TypeString,
									},
									"pin_iam_credentials_secret_id": &schema.Schema{
										Description: "The secret Id of iam credentials with api key to access HPCS instance.",
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Type:        schema.TypeString,
									},
									"private_keystore_id": &schema.Schema{
										Description: "The HPCS private key store space id.",
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Type:        schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique serial number that was assigned to a certificate by the issuing certificate authority.",
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
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceIbmSmPrivateCertificateConfigurationRootCACreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigRootCAResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	configurationPrototypeModel, err := resourceIbmSmPrivateCertificateConfigurationRootCAMapToConfigurationPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigRootCAResourceName, "create")
		return tfErr.GetDiag()
	}
	createConfigurationOptions.SetConfigurationPrototype(configurationPrototypeModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationWithContext failed: %s\n%s", err.Error(), response), PrivateCertConfigRootCAResourceName, "create")
		return tfErr.GetDiag()
	}

	configuration := configurationIntf.(*secretsmanagerv2.PrivateCertificateConfigurationRootCA)
	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmPrivateCertificateConfigurationRootCARead(context, d, meta)
}

func resourceIbmSmPrivateCertificateConfigurationRootCARead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a root CA use the format `<region>/<instance_id>/<name>`", PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(configName)

	configurationIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}

	configuration := configurationIntf.(*secretsmanagerv2.PrivateCertificateConfigurationRootCA)

	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", configuration.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", configuration.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", DateTimeToRFC3339(configuration.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", DateTimeToRFC3339(configuration.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("max_ttl_seconds", flex.IntValue(configuration.MaxTtlSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting max_ttl_seconds"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if d.Get("max_ttl") == nil || d.Get("max_ttl") == "" {
		if err = d.Set("max_ttl", strconv.FormatInt(*configuration.MaxTtlSeconds, 10)+"s"); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting max_ttl"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("crl_expiry_seconds", flex.IntValue(configuration.CrlExpirySeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_expiry_seconds"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if d.Get("crl_expiry") == nil || d.Get("crl_expiry") == "" {
		if err = d.Set("crl_expiry", strconv.FormatInt(*configuration.CrlExpirySeconds, 10)+"s"); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_expiry"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("crl_disable", configuration.CrlDisable); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_disable"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("crl_distribution_points_encoded", configuration.CrlDistributionPointsEncoded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_distribution_points_encoded"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("issuing_certificates_urls_encoded", configuration.IssuingCertificatesUrlsEncoded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuing_certificates_urls_encoded"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("common_name", configuration.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.AltNames != nil {
		if err = d.Set("alt_names", configuration.AltNames); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting alt_names"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("ip_sans", configuration.IpSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ip_sans"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("uri_sans", configuration.UriSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting uri_sans"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.OtherSans != nil {
		if err = d.Set("other_sans", configuration.OtherSans); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting other_sans"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("ttl_seconds", flex.IntValue(configuration.TtlSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ttl_seconds"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("format", configuration.Format); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting format"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("private_key_format", configuration.PrivateKeyFormat); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key_format"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("key_type", configuration.KeyType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_type"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("key_bits", flex.IntValue(configuration.KeyBits)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_bits"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("max_path_length", flex.IntValue(configuration.MaxPathLength)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting max_path_length"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("exclude_cn_from_sans", configuration.ExcludeCnFromSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting exclude_cn_from_sans"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.PermittedDnsDomains != nil {
		if err = d.Set("permitted_dns_domains", configuration.PermittedDnsDomains); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting permitted_dns_domains"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Ou != nil {
		if err = d.Set("ou", configuration.Ou); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ou"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Organization != nil {
		if err = d.Set("organization", configuration.Organization); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting organization"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Country != nil {
		if err = d.Set("country", configuration.Country); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting country"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Locality != nil {
		if err = d.Set("locality", configuration.Locality); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locality"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Province != nil {
		if err = d.Set("province", configuration.Province); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting province"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.StreetAddress != nil {
		if err = d.Set("street_address", configuration.StreetAddress); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting street_address"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.PostalCode != nil {
		if err = d.Set("postal_code", configuration.PostalCode); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting postal_code"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("serial_number", configuration.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("status", configuration.Status); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting status"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("expiration_date", DateTimeToRFC3339(configuration.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), PrivateCertConfigRootCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.Data != nil {
		dataMap, err := resourceIbmSmPrivateCertificateConfigurationRootCAPrivateCertificateCADataToMap(configuration.Data)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
		if err = d.Set("data", []map[string]interface{}{dataMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting data"), PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.CryptoKey != nil {
		cryptoKeyMap, err := resourceIbmSmPrivateCertificateConfigurationCryptoKeyToMap(configuration.CryptoKey)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigRootCAResourceName, "read")
			return tfErr.GetDiag()
		}
		if len(cryptoKeyMap) > 0 {
			if err = d.Set("crypto_key", []map[string]interface{}{cryptoKeyMap}); err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crypto_key"), PrivateCertConfigRootCAResourceName, "read")
				return tfErr.GetDiag()
			}
		}
	}

	return nil
}

func resourceIbmSmPrivateCertificateConfigurationRootCAUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigRootCAResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateConfigurationOptions := &secretsmanagerv2.UpdateConfigurationOptions{}

	updateConfigurationOptions.SetName(configName)

	hasChange := false

	patchVals := &secretsmanagerv2.PrivateCertificateConfigurationRootCAPatch{}

	if d.HasChange("max_ttl") {
		patchVals.MaxTTL = core.StringPtr(d.Get("max_ttl").(string))
		hasChange = true
	}
	if d.HasChange("crl_expiry") {
		patchVals.CrlExpiry = core.StringPtr(d.Get("crl_expiry").(string))
		hasChange = true
	}
	if d.HasChange("crl_disable") {
		patchVals.CrlDisable = core.BoolPtr(d.Get("crl_disable").(bool))
		hasChange = true
	}
	if d.HasChange("crl_distribution_points_encoded") {
		patchVals.CrlDistributionPointsEncoded = core.BoolPtr(d.Get("crl_distribution_points_encoded").(bool))
		hasChange = true
	}
	if d.HasChange("issuing_certificates_urls_encoded") {
		patchVals.IssuingCertificatesUrlsEncoded = core.BoolPtr(d.Get("issuing_certificates_urls_encoded").(bool))
		hasChange = true
	}

	if hasChange {
		updateConfigurationOptions.ConfigurationPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateConfigurationWithContext(context, updateConfigurationOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigurationWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigRootCAResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmPrivateCertificateConfigurationRootCARead(context, d, meta)
}

func resourceIbmSmPrivateCertificateConfigurationRootCADelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigRootCAResourceName, "delete")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	deleteConfigurationOptions := &secretsmanagerv2.DeleteConfigurationOptions{}

	deleteConfigurationOptions.SetName(configName)

	response, err := secretsManagerClient.DeleteConfigurationWithContext(context, deleteConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigRootCAResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmPrivateCertificateConfigurationRootCAMapToConfigurationPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationPrototypeIntf, error) {
	model := &secretsmanagerv2.PrivateCertificateConfigurationRootCAPrototype{
		ConfigType: core.StringPtr("private_cert_configuration_root_ca"),
	}
	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("max_ttl"); ok {
		model.MaxTTL = core.StringPtr(d.Get("max_ttl").(string))
	}
	if _, ok := d.GetOk("crl_expiry"); ok {
		model.CrlExpiry = core.StringPtr(d.Get("crl_expiry").(string))
	}
	if _, ok := d.GetOk("crl_disable"); ok {
		model.CrlDisable = core.BoolPtr(d.Get("crl_disable").(bool))
	}
	if _, ok := d.GetOk("crl_distribution_points_encoded"); ok {
		model.CrlDistributionPointsEncoded = core.BoolPtr(d.Get("crl_distribution_points_encoded").(bool))
	}
	if _, ok := d.GetOk("issuing_certificates_urls_encoded"); ok {
		model.IssuingCertificatesUrlsEncoded = core.BoolPtr(d.Get("issuing_certificates_urls_encoded").(bool))
	}
	if _, ok := d.GetOk("common_name"); ok {
		model.CommonName = core.StringPtr(d.Get("common_name").(string))
	}
	if _, ok := d.GetOk("alt_names"); ok {
		altNames := d.Get("alt_names").([]interface{})
		altNamesParsed := make([]string, len(altNames))
		for i, v := range altNames {
			altNamesParsed[i] = fmt.Sprint(v)
		}
		model.AltNames = altNamesParsed
	}
	if _, ok := d.GetOk("ip_sans"); ok {
		model.IpSans = core.StringPtr(d.Get("ip_sans").(string))
	}
	if _, ok := d.GetOk("uri_sans"); ok {
		model.UriSans = core.StringPtr(d.Get("uri_sans").(string))
	}
	if _, ok := d.GetOk("other_sans"); ok {
		otherSans := d.Get("other_sans").([]interface{})
		otherSansParsed := make([]string, len(otherSans))
		for i, v := range otherSans {
			otherSansParsed[i] = fmt.Sprint(v)
		}
		model.OtherSans = otherSansParsed
	}
	if _, ok := d.GetOk("ttl"); ok {
		model.TTL = core.StringPtr(d.Get("ttl").(string))
	}
	if _, ok := d.GetOk("format"); ok {
		model.Format = core.StringPtr(d.Get("format").(string))
	}
	if _, ok := d.GetOk("private_key_format"); ok {
		model.PrivateKeyFormat = core.StringPtr(d.Get("private_key_format").(string))
	}
	if _, ok := d.GetOk("key_type"); ok {
		model.KeyType = core.StringPtr(d.Get("key_type").(string))
	}
	if _, ok := d.GetOk("key_bits"); ok {
		model.KeyBits = core.Int64Ptr(int64(d.Get("key_bits").(int)))
	}
	if _, ok := d.GetOk("max_path_length"); ok {
		model.MaxPathLength = core.Int64Ptr(int64(d.Get("max_path_length").(int)))
	}
	if _, ok := d.GetOk("exclude_cn_from_sans"); ok {
		model.ExcludeCnFromSans = core.BoolPtr(d.Get("exclude_cn_from_sans").(bool))
	}
	if _, ok := d.GetOk("permitted_dns_domains"); ok {
		permittedDnsDomains := d.Get("permitted_dns_domains").([]interface{})
		permittedDnsDomainsParsed := make([]string, len(permittedDnsDomains))
		for i, v := range permittedDnsDomains {
			permittedDnsDomainsParsed[i] = fmt.Sprint(v)
		}
		model.PermittedDnsDomains = permittedDnsDomainsParsed
	}
	if _, ok := d.GetOk("ou"); ok {
		ou := d.Get("ou").([]interface{})
		ouParsed := make([]string, len(ou))
		for i, v := range ou {
			ouParsed[i] = fmt.Sprint(v)
		}
		model.Ou = ouParsed
	}
	if _, ok := d.GetOk("organization"); ok {
		organization := d.Get("organization").([]interface{})
		organizationParsed := make([]string, len(organization))
		for i, v := range organization {
			organizationParsed[i] = fmt.Sprint(v)
		}
		model.Organization = organizationParsed
	}
	if _, ok := d.GetOk("country"); ok {
		country := d.Get("country").([]interface{})
		countryParsed := make([]string, len(country))
		for i, v := range country {
			countryParsed[i] = fmt.Sprint(v)
		}
		model.Country = countryParsed
	}
	if _, ok := d.GetOk("locality"); ok {
		locality := d.Get("locality").([]interface{})
		localityParsed := make([]string, len(locality))
		for i, v := range locality {
			localityParsed[i] = fmt.Sprint(v)
		}
		model.Locality = localityParsed
	}
	if _, ok := d.GetOk("province"); ok {
		province := d.Get("province").([]interface{})
		provinceParsed := make([]string, len(province))
		for i, v := range province {
			provinceParsed[i] = fmt.Sprint(v)
		}
		model.Province = provinceParsed
	}
	if _, ok := d.GetOk("street_address"); ok {
		streetAddress := d.Get("street_address").([]interface{})
		streetAddressParsed := make([]string, len(streetAddress))
		for i, v := range streetAddress {
			streetAddressParsed[i] = fmt.Sprint(v)
		}
		model.StreetAddress = streetAddressParsed
	}
	if _, ok := d.GetOk("postal_code"); ok {
		postalCode := d.Get("postal_code").([]interface{})
		postalCodeParsed := make([]string, len(postalCode))
		for i, v := range postalCode {
			postalCodeParsed[i] = fmt.Sprint(v)
		}
		model.PostalCode = postalCodeParsed
	}
	if _, ok := d.GetOk("crypto_key"); ok {
		if len(d.Get("crypto_key").([]interface{})) > 0 {
			CryptoKeyModel, err := resourceIbmSmPrivateCertificateConfigurationMapToPrivateCertificateConfigurationCryptoKey(d.Get("crypto_key").([]interface{})[0].(map[string]interface{}))
			if err != nil {
				return model, err
			}
			model.CryptoKey = CryptoKeyModel
		}
	}
	return model, nil
}

func resourceIbmSmPrivateCertificateConfigurationRootCAPrivateCertificateCADataToMap(modelIntf secretsmanagerv2.PrivateCertificateCADataIntf) (map[string]interface{}, error) {
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
