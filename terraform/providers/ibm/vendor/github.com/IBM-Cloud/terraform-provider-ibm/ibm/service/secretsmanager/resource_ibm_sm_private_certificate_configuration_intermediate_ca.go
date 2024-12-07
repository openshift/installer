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

func ResourceIbmSmPrivateCertificateConfigurationIntermediateCA() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPrivateCertificateConfigurationIntermediateCACreate,
		ReadContext:   resourceIbmSmPrivateCertificateConfigurationIntermediateCARead,
		UpdateContext: resourceIbmSmPrivateCertificateConfigurationIntermediateCAUpdate,
		DeleteContext: resourceIbmSmPrivateCertificateConfigurationIntermediateCADelete,
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
			"secret_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The secret type. Supported types are arbitrary, certificates (imported, public, and private), IAM credentials, key-value, and user credentials.",
			},
			"max_ttl": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
			},
			"max_ttl_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
			},
			"crl_expiry": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The time until the certificate revocation list (CRL) expires.The value can be supplied as a string representation of a duration in hours, such as `48h`. The default is 72 hours. In the API response, this value is returned in seconds (integer).**Note:** The CRL is rotated automatically before it expires.",
			},
			"crl_expiry_seconds": &schema.Schema{
				Type:        schema.TypeInt,
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
				ForceNew:    true,
				Default:     "der",
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
			"exclude_cn_from_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				ForceNew:    true,
				Description: "Controls whether the common name is excluded from Subject Alternative Names (SANs).If the common name set to `true`, it is not included in DNS or Email SANs if they apply. This field can be useful if the common name is a human-readable identifier, instead of a hostname or an email address.",
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
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique serial number that was assigned to a certificate by the issuing certificate authority.",
			},
			"signing_method": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The signing method to use with this certificate authority to generate private certificates.You can choose between internal or externally signed options. For more information, see the [docs](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-intermediate-certificate-authorities).",
			},
			"issuer": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The distinguished name that identifies the entity that signed and issued the certificate.",
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
			// parameters for signing intermediate actions (internal)
			"ttl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the requested Time To Live (after which the certificate will be expired). The value can be provided provided as a string duration with time suffix (e.g. '24h') or the number of seconds as string (e.g. '86400').",
			},
			"max_path_length": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The maximum path length to encode in the generated certificate. `-1` means no limit.",
			},
			"permitted_dns_domains": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The allowed DNS domains or subdomains for the certificates that are to be signed and issued by this CA certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"use_csr_values": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Determines whether to use values from a certificate signing request (CSR) to complete a `private_cert_configuration_action_sign_csr` action.",
			},
		},
	}
}

func resourceIbmSmPrivateCertificateConfigurationIntermediateCACreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	configurationPrototypeModel, err := resourceIbmSmPrivateCertificateConfigurationIntermediateCAMapToConfigurationPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "create")
		return tfErr.GetDiag()
	}
	createConfigurationOptions.SetConfigurationPrototype(configurationPrototypeModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationWithContext failed: %s\n%s", err.Error(), response), PrivateCertConfigIntermediateCAResourceName, "create")
		return tfErr.GetDiag()
	}
	configuration := configurationIntf.(*secretsmanagerv2.PrivateCertificateConfigurationIntermediateCA)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	// signing the CSR
	if signingMethod, ok := d.GetOk("signing_method"); ok && signingMethod.(string) == "internal" {
		if _, ok := d.GetOk("issuer"); ok {
			createConfigurationActionOptions := &secretsmanagerv2.CreateConfigurationActionOptions{}

			createConfigurationActionOptions.SetName(d.Get("issuer").(string))
			configurationActionPrototypeModel, err := resourceIbmSmConfigurationActionPrivateCertificateSignIntermediateCAMapToConfigurationActionPrototype(d)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "create")
				return tfErr.GetDiag()
			}
			createConfigurationActionOptions.SetConfigActionPrototype(configurationActionPrototypeModel)

			_, responseAction, errAction := secretsManagerClient.CreateConfigurationActionWithContext(context, createConfigurationActionOptions)
			if errAction != nil {
				log.Printf("[DEBUG] CreateConfigurationActionWithContext failed %s\n%s", errAction, responseAction)
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationActionWithContext failed %s\n%s", errAction, responseAction), PrivateCertConfigIntermediateCAResourceName, "create")
				return tfErr.GetDiag()
			}
		} else {
			tfErr := flex.TerraformErrorf(nil, "`issuer` parameter is missing", PrivateCertConfigIntermediateCAResourceName, "create")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmPrivateCertificateConfigurationIntermediateCARead(context, d, meta)
}

func resourceIbmSmPrivateCertificateConfigurationIntermediateCARead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import an intermediate CA use the format `<region>/<instance_id>/<name>`", PrivateCertConfigIntermediateCAResourceName, "read")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	configuration := configurationIntf.(*secretsmanagerv2.PrivateCertificateConfigurationIntermediateCA)

	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", configuration.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("config_type", configuration.ConfigType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("max_ttl_seconds", flex.IntValue(configuration.MaxTtlSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting max_ttl_seconds"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if d.Get("max_ttl") == nil || d.Get("max_ttl") == "" {
		if err = d.Set("max_ttl", strconv.FormatInt(*configuration.MaxTtlSeconds, 10)+"s"); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting max_ttl"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("crl_expiry_seconds", flex.IntValue(configuration.CrlExpirySeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_expiry_seconds"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if d.Get("crl_expiry") == nil || d.Get("crl_expiry") == "" {
		if err = d.Set("crl_expiry", strconv.FormatInt(*configuration.CrlExpirySeconds, 10)+"s"); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_expiry"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("signing_method", configuration.SigningMethod); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting signing_method"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("issuer", configuration.Issuer); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuer"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("crl_disable", configuration.CrlDisable); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_disable"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("crl_distribution_points_encoded", configuration.CrlDistributionPointsEncoded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crl_distribution_points_encoded"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("issuing_certificates_urls_encoded", configuration.IssuingCertificatesUrlsEncoded); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting issuing_certificates_urls_encoded"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("common_name", configuration.CommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting common_name"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.AltNames != nil {
		if err = d.Set("alt_names", configuration.AltNames); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting alt_names"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("ip_sans", configuration.IpSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ip_sans"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("uri_sans", configuration.UriSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting uri_sans"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.OtherSans != nil {
		if err = d.Set("other_sans", configuration.OtherSans); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting other_sans"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("format", configuration.Format); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting format"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("private_key_format", configuration.PrivateKeyFormat); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting private_key_format"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("key_type", configuration.KeyType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_type"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("key_bits", flex.IntValue(configuration.KeyBits)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_bits"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("exclude_cn_from_sans", configuration.ExcludeCnFromSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting exclude_cn_from_sans"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.Ou != nil {
		if err = d.Set("ou", configuration.Ou); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ou"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Organization != nil {
		if err = d.Set("organization", configuration.Organization); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting organization"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Country != nil {
		if err = d.Set("country", configuration.Country); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting country"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Locality != nil {
		if err = d.Set("locality", configuration.Locality); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locality"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Province != nil {
		if err = d.Set("province", configuration.Province); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting province"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.StreetAddress != nil {
		if err = d.Set("street_address", configuration.StreetAddress); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting street_address"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.PostalCode != nil {
		if err = d.Set("postal_code", configuration.PostalCode); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting postal_code"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("serial_number", configuration.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("status", configuration.Status); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting status"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("expiration_date", DateTimeToRFC3339(configuration.ExpirationDate)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting expiration_date"), PrivateCertConfigIntermediateCAResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.Data != nil {
		dataMap, err := resourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateCADataToMap(configuration.Data)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
		if err = d.Set("data", []map[string]interface{}{dataMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting data"), PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.CryptoKey != nil {
		cryptoKeyMap, err := resourceIbmSmPrivateCertificateConfigurationCryptoKeyToMap(configuration.CryptoKey)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "read")
			return tfErr.GetDiag()
		}
		if len(cryptoKeyMap) > 0 {
			if err = d.Set("crypto_key", []map[string]interface{}{cryptoKeyMap}); err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting crypto_key"), PrivateCertConfigIntermediateCAResourceName, "read")
				return tfErr.GetDiag()
			}
		}
	}
	return nil
}

func resourceIbmSmPrivateCertificateConfigurationCryptoKeyToMap(model *secretsmanagerv2.PrivateCertificateCryptoKey) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Label != nil {
		modelMap["label"] = model.Label
	}
	if model.AllowGenerateKey != nil {
		modelMap["allow_generate_key"] = model.AllowGenerateKey
	}
	if model.Provider != nil {
		providerModelMap, err := resourceIbmSmPrivateCertificateConfigurationCryptoKeyProviderToMap(model.Provider)
		if err != nil {
			return modelMap, err
		}
		modelMap["provider"] = []map[string]interface{}{providerModelMap}
	}

	return modelMap, nil
}

func resourceIbmSmPrivateCertificateConfigurationCryptoKeyProviderToMap(providerModelIntf secretsmanagerv2.PrivateCertificateCryptoProviderIntf) (map[string]interface{}, error) {
	providerModelMap := make(map[string]interface{})
	providerModel := providerModelIntf.(*secretsmanagerv2.PrivateCertificateCryptoProviderHPCS)

	if providerModel.Type != nil {
		providerModelMap["type"] = providerModel.Type
	}
	if providerModel.InstanceCrn != nil {
		providerModelMap["instance_crn"] = providerModel.InstanceCrn
	}
	if providerModel.PinIamCredentialsSecretID != nil {
		providerModelMap["pin_iam_credentials_secret_id"] = providerModel.PinIamCredentialsSecretID
	}
	if providerModel.PrivateKeystoreID != nil {
		providerModelMap["private_keystore_id"] = providerModel.PrivateKeystoreID
	}

	return providerModelMap, nil
}

func resourceIbmSmPrivateCertificateConfigurationIntermediateCAUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateConfigurationOptions := &secretsmanagerv2.UpdateConfigurationOptions{}

	updateConfigurationOptions.SetName(configName)
	updateConfigurationOptions.SetXSmAcceptConfigurationType("private_cert_configuration_intermediate_ca")

	hasChange := false

	patchVals := &secretsmanagerv2.ConfigurationPatch{}

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
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigIntermediateCAResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmPrivateCertificateConfigurationIntermediateCARead(context, d, meta)
}

func resourceIbmSmPrivateCertificateConfigurationIntermediateCADelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigIntermediateCAResourceName, "delete")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigIntermediateCAResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmPrivateCertificateConfigurationIntermediateCAMapToConfigurationPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationPrototypeIntf, error) {
	model := &secretsmanagerv2.PrivateCertificateConfigurationIntermediateCAPrototype{}

	model.ConfigType = core.StringPtr("private_cert_configuration_intermediate_ca")

	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("issuer"); ok {
		model.Issuer = core.StringPtr(d.Get("issuer").(string))
	}
	if _, ok := d.GetOk("common_name"); ok {
		model.CommonName = core.StringPtr(d.Get("common_name").(string))
	}
	if _, ok := d.GetOk("signing_method"); ok {
		model.SigningMethod = core.StringPtr(d.Get("signing_method").(string))
	}
	if _, ok := d.GetOk("max_ttl"); ok {
		model.MaxTTL = core.StringPtr(d.Get("max_ttl").(string))
	}
	if _, ok := d.GetOk("issuer"); ok {
		model.Issuer = core.StringPtr(d.Get("issuer").(string))
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
	if _, ok := d.GetOk("alt_names"); ok {
		altNames := []string{}
		for _, altNamesItem := range d.Get("alt_names").([]interface{}) {
			altNames = append(altNames, altNamesItem.(string))
		}
		model.AltNames = altNames
	}
	if _, ok := d.GetOk("ip_sans"); ok {
		model.IpSans = core.StringPtr(d.Get("ip_sans").(string))
	}
	if _, ok := d.GetOk("uri_sans"); ok {
		model.UriSans = core.StringPtr(d.Get("uri_sans").(string))
	}
	if _, ok := d.GetOk("other_sans"); ok {
		otherSans := []string{}
		for _, otherSansItem := range d.Get("other_sans").([]interface{}) {
			otherSans = append(otherSans, otherSansItem.(string))
		}
		model.OtherSans = otherSans
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
	if _, ok := d.GetOk("exclude_cn_from_sans"); ok {
		model.ExcludeCnFromSans = core.BoolPtr(d.Get("exclude_cn_from_sans").(bool))
	}
	if _, ok := d.GetOk("ou"); ok {
		ou := []string{}
		for _, ouItem := range d.Get("ou").([]interface{}) {
			ou = append(ou, ouItem.(string))
		}
		model.Ou = ou
	}
	if _, ok := d.GetOk("organization"); ok {
		organization := []string{}
		for _, organizationItem := range d.Get("organization").([]interface{}) {
			organization = append(organization, organizationItem.(string))
		}
		model.Organization = organization
	}
	if _, ok := d.GetOk("country"); ok {
		country := []string{}
		for _, countryItem := range d.Get("country").([]interface{}) {
			country = append(country, countryItem.(string))
		}
		model.Country = country
	}
	if _, ok := d.GetOk("locality"); ok {
		locality := []string{}
		for _, localityItem := range d.Get("locality").([]interface{}) {
			locality = append(locality, localityItem.(string))
		}
		model.Locality = locality
	}
	if _, ok := d.GetOk("province"); ok {
		province := []string{}
		for _, provinceItem := range d.Get("province").([]interface{}) {
			province = append(province, provinceItem.(string))
		}
		model.Province = province
	}
	if _, ok := d.GetOk("street_address"); ok {
		streetAddress := []string{}
		for _, streetAddressItem := range d.Get("street_address").([]interface{}) {
			streetAddress = append(streetAddress, streetAddressItem.(string))
		}
		model.StreetAddress = streetAddress
	}
	if _, ok := d.GetOk("postal_code"); ok {
		postalCode := []string{}
		for _, postalCodeItem := range d.Get("postal_code").([]interface{}) {
			postalCode = append(postalCode, postalCodeItem.(string))
		}
		model.PostalCode = postalCode
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

func resourceIbmSmPrivateCertificateConfigurationMapToPrivateCertificateConfigurationCryptoKey(modelMap map[string]interface{}) (*secretsmanagerv2.PrivateCertificateCryptoKey, error) {
	model := &secretsmanagerv2.PrivateCertificateCryptoKey{}
	if modelMap["id"] != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["label"] != "" {
		model.Label = core.StringPtr(modelMap["label"].(string))
	}
	if modelMap["allow_generate_key"] != nil {
		model.AllowGenerateKey = core.BoolPtr(modelMap["allow_generate_key"].(bool))
	}
	if modelMap["provider"] != nil && len(modelMap["provider"].([]interface{})) > 0 {
		providerModel, err := resourceIbmSmPrivateCertificateConfigurationMapToPrivateCertificateConfigurationCryptoKeyProvider(modelMap["provider"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Provider = providerModel
	}

	return model, nil
}

func resourceIbmSmPrivateCertificateConfigurationMapToPrivateCertificateConfigurationCryptoKeyProvider(modelMapProvider map[string]interface{}) (secretsmanagerv2.PrivateCertificateCryptoProviderIntf, error) {
	modelProvider := &secretsmanagerv2.PrivateCertificateCryptoProviderHPCS{}
	if modelMapProvider["type"] != "" {
		modelProvider.Type = core.StringPtr(modelMapProvider["type"].(string))
	}
	if modelMapProvider["instance_crn"] != "" {
		modelProvider.InstanceCrn = core.StringPtr(modelMapProvider["instance_crn"].(string))
	}
	if modelMapProvider["pin_iam_credentials_secret_id"] != "" {
		modelProvider.PinIamCredentialsSecretID = core.StringPtr(modelMapProvider["pin_iam_credentials_secret_id"].(string))
	}
	if modelMapProvider["private_keystore_id"] != "" {
		modelProvider.PrivateKeystoreID = core.StringPtr(modelMapProvider["private_keystore_id"].(string))
	}

	return modelProvider, nil
}

func resourceIbmSmPrivateCertificateConfigurationIntermediateCAPrivateCertificateCADataToMap(modelIntf secretsmanagerv2.PrivateCertificateCADataIntf) (map[string]interface{}, error) {
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

func resourceIbmSmConfigurationActionPrivateCertificateSignIntermediateCAMapToConfigurationActionPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationActionPrototypeIntf, error) {
	model := &secretsmanagerv2.PrivateCertificateConfigurationActionSignIntermediatePrototype{}

	model.ActionType = core.StringPtr("private_cert_configuration_action_sign_intermediate")
	if _, ok := d.GetOk("name"); ok {
		model.IntermediateCertificateAuthority = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("common_name"); ok {
		model.CommonName = core.StringPtr(d.Get("common_name").(string))
	}
	if _, ok := d.GetOk("alt_names"); ok {
		altNames := []string{}
		for _, altNamesItem := range d.Get("alt_names").([]interface{}) {
			altNames = append(altNames, altNamesItem.(string))
		}
		model.AltNames = altNames
	}
	if _, ok := d.GetOk("ip_sans"); ok {
		model.IpSans = core.StringPtr(d.Get("ip_sans").(string))
	}
	if _, ok := d.GetOk("uri_sans"); ok {
		model.UriSans = core.StringPtr(d.Get("uri_sans").(string))
	}
	if _, ok := d.GetOk("other_sans"); ok {
		otherSans := []string{}
		for _, otherSansItem := range d.Get("other_sans").([]interface{}) {
			otherSans = append(otherSans, otherSansItem.(string))
		}
		model.OtherSans = otherSans
	}
	if _, ok := d.GetOk("ttl"); ok {
		model.TTL = core.StringPtr(d.Get("ttl").(string))
	}
	if _, ok := d.GetOk("max_path_length"); ok {
		model.MaxPathLength = core.Int64Ptr(int64(d.Get("max_path_length").(int)))
	}
	if _, ok := d.GetOk("exclude_cn_from_sans"); ok {
		model.ExcludeCnFromSans = core.BoolPtr(d.Get("exclude_cn_from_sans").(bool))
	}
	if _, ok := d.GetOk("permitted_dns_domains"); ok {
		permittedDnsDomains := []string{}
		for _, permittedDnsDomainsItem := range d.Get("permitted_dns_domains").([]interface{}) {
			permittedDnsDomains = append(permittedDnsDomains, permittedDnsDomainsItem.(string))
		}
		model.PermittedDnsDomains = permittedDnsDomains
	}
	if _, ok := d.GetOk("use_csr_values"); ok {
		model.UseCsrValues = core.BoolPtr(d.Get("use_csr_values").(bool))
	}
	if _, ok := d.GetOk("ou"); ok {
		ou := []string{}
		for _, ouItem := range d.Get("ou").([]interface{}) {
			ou = append(ou, ouItem.(string))
		}
		model.Ou = ou
	}
	if _, ok := d.GetOk("organization"); ok {
		organization := []string{}
		for _, organizationItem := range d.Get("organization").([]interface{}) {
			organization = append(organization, organizationItem.(string))
		}
		model.Organization = organization
	}
	if _, ok := d.GetOk("country"); ok {
		country := []string{}
		for _, countryItem := range d.Get("country").([]interface{}) {
			country = append(country, countryItem.(string))
		}
		model.Country = country
	}
	if _, ok := d.GetOk("locality"); ok {
		locality := []string{}
		for _, localityItem := range d.Get("locality").([]interface{}) {
			locality = append(locality, localityItem.(string))
		}
		model.Locality = locality
	}
	if _, ok := d.GetOk("province"); ok {
		province := []string{}
		for _, provinceItem := range d.Get("province").([]interface{}) {
			province = append(province, provinceItem.(string))
		}
		model.Province = province
	}
	if _, ok := d.GetOk("street_address"); ok {
		streetAddress := []string{}
		for _, streetAddressItem := range d.Get("street_address").([]interface{}) {
			streetAddress = append(streetAddress, streetAddressItem.(string))
		}
		model.StreetAddress = streetAddress
	}
	if _, ok := d.GetOk("postal_code"); ok {
		postalCode := []string{}
		for _, postalCodeItem := range d.Get("postal_code").([]interface{}) {
			postalCode = append(postalCode, postalCodeItem.(string))
		}
		model.PostalCode = postalCode
	}

	return model, nil
}
