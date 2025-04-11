// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package secretsmanager

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPrivateCertificateConfigurationTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPrivateCertificateConfigurationTemplateCreate,
		ReadContext:   resourceIbmSmPrivateCertificateConfigurationTemplateRead,
		UpdateContext: resourceIbmSmPrivateCertificateConfigurationTemplateUpdate,
		DeleteContext: resourceIbmSmPrivateCertificateConfigurationTemplateDelete,
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
				Optional:    true,
				Description: "The maximum time-to-live (TTL) for certificates that are created by this CA.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).Minimum value is one hour (`1h`). Maximum value is 100 years (`876000h`).",
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The requested time-to-live (TTL) for certificates that are created by this CA. This field's value cannot be longer than the `max_ttl` limit.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).",
			},
			"key_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of private key to generate.",
			},
			"key_bits": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The number of bits to use to generate the private key.Allowable values for RSA keys are: `2048` and `4096`. Allowable values for EC keys are: `224`, `256`, `384`, and `521`. The default for RSA keys is `2048`. The default for EC keys is `256`.",
			},
			"ou": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"organization": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"country": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The Country (C) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"locality": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"province": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"street_address": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The street address values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"postal_code": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The postal code values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Unused field.",
				Deprecated:  "This field is deprecated.",
			},
			"certificate_authority": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the intermediate certificate authority.",
			},
			"allowed_secret_groups": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Scopes the creation of private certificates to only the secret groups that you specify.This field can be supplied as a comma-delimited list of secret group IDs.",
			},
			"allow_localhost": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to allow `localhost` to be included as one of the requested common names.",
			},
			"allowed_domains": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and `allow_subdomains` options.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"allowed_domains_template": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to allow the domains that are supplied in the `allowed_domains` field to contain access control list (ACL) templates.",
			},
			"allow_bare_domains": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to allow clients to request private certificates that match the value of the actual domains on the final certificate.For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a certificate that contains the name `example.com` as one of the DNS values on the final certificate.**Important:** In some scenarios, allowing bare domains can be considered a security risk.",
			},
			"allow_subdomains": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to allow clients to request private certificates with common names (CN) that are subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.**Note:** This field is redundant if you use the `allow_any_name` option.",
			},
			"allow_glob_domains": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are specified in the `allowed_domains` field.If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.",
			},
			"allow_wildcard_certificates": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether the issuance of certificates with RFC 6125 wildcards in the CN field.When set to false, this field prevents wildcards from being issued even if they can be allowed by an option `allow_glob_domains`.",
			},
			"allow_any_name": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to allow clients to request a private certificate that matches any common name.",
			},
			"enforce_hostnames": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to enforce only valid host names for common names, DNS Subject Alternative Names, and the host section of email addresses.",
			},
			"allow_ip_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to allow clients to request a private certificate with IP Subject Alternative Names.",
			},
			"allowed_uri_sans": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The URI Subject Alternative Names to allow for private certificates.Values can contain glob patterns, for example `spiffe://hostname/_*`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"allowed_other_sans": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private certificates.The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any `other_sans` input.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"server_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for server use.",
			},
			"client_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for client use.",
			},
			"code_signing_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for code signing use.",
			},
			"email_protection_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for email protection use.",
			},
			"key_usage": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The allowed key usage constraint to define for private certificates.You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ext_key_usage": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The allowed extended key usage constraint on private certificates.You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage). Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ext_key_usage_oids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of extended key usage Object Identifiers (OIDs).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"use_csr_common_name": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the common name (CN) from a certificate signing request (CSR) instead of the CN that's included in the data of the certificate.Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include the `use_csr_sans` property.",
			},
			"use_csr_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the Subject Alternative Names(SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the certificate.Does not include the common name in the CSR. To use the common name, include the `use_csr_common_name` property.",
			},
			"require_cn": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to require a common name to create a private certificate.By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the `require_cn` option to `false`.",
			},
			"policy_identifiers": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of policy Object Identifiers (OIDs).",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"basic_constraints_valid_for_non_ca": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Determines whether to mark the Basic Constraints extension of an issued private certificate as valid for non-CA certificates.",
			},
			"not_before_duration": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The duration in seconds by which to backdate the `not_before` property of an issued private certificate.The value can be supplied as a string representation of a duration, such as `30s`. In the API response, this value is returned in seconds (integer).",
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
			"ttl_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The requested Time To Live, after which the certificate will be expired.",
			},
			"not_before_duration_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration in seconds by which to backdate the `not_before` property of an issued private certificate.",
			},
		},
	}
}

func resourceIbmSmPrivateCertificateConfigurationTemplateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigTemplateResourceName, "create")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationOptions := &secretsmanagerv2.CreateConfigurationOptions{}

	configurationPrototypeModel, err := resourceIbmSmPrivateCertificateConfigurationTemplateMapToConfigurationPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigTemplateResourceName, "create")
		return tfErr.GetDiag()
	}
	createConfigurationOptions.SetConfigurationPrototype(configurationPrototypeModel)

	configurationIntf, response, err := secretsManagerClient.CreateConfigurationWithContext(context, createConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationWithContext failed: %s\n%s", err.Error(), response), PrivateCertConfigTemplateResourceName, "create")
		return tfErr.GetDiag()
	}
	configuration := configurationIntf.(*secretsmanagerv2.PrivateCertificateConfigurationTemplate)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *configuration.Name))

	return resourceIbmSmPrivateCertificateConfigurationTemplateRead(context, d, meta)
}

func resourceIbmSmPrivateCertificateConfigurationTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	if len(id) != 3 {
		tfErr := flex.TerraformErrorf(nil, "Wrong format of resource ID. To import a certificate template use the format `<region>/<instance_id>/<name>`", PrivateCertConfigTemplateResourceName, "read")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	configuration := configurationIntf.(*secretsmanagerv2.PrivateCertificateConfigurationTemplate)

	if err = d.Set("instance_id", instanceId); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting instance_id"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("name", configuration.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting name"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("config_type", configuration.ConfigType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("secret_type", configuration.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("certificate_authority", configuration.CertificateAuthority); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate_authority"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("allowed_secret_groups", configuration.AllowedSecretGroups); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_secret_groups"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("max_ttl_seconds", flex.IntValue(configuration.MaxTtlSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting max_ttl_seconds"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("ttl_seconds", flex.IntValue(configuration.TtlSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ttl_seconds"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("allow_localhost", configuration.AllowLocalhost); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_localhost"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.AllowedDomains != nil {
		if err = d.Set("allowed_domains", configuration.AllowedDomains); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_domains"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("allowed_domains_template", configuration.AllowedDomainsTemplate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_domains_template"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("allow_bare_domains", configuration.AllowBareDomains); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_bare_domains"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("allow_subdomains", configuration.AllowSubdomains); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_subdomains"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("allow_glob_domains", configuration.AllowGlobDomains); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_glob_domains"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("allow_any_name", configuration.AllowAnyName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_any_name"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("enforce_hostnames", configuration.EnforceHostnames); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting enforce_hostnames"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("allow_ip_sans", configuration.AllowIpSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_ip_sans"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.AllowedUriSans != nil {
		if err = d.Set("allowed_uri_sans", configuration.AllowedUriSans); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_uri_sans"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.AllowedOtherSans != nil {
		if err = d.Set("allowed_other_sans", configuration.AllowedOtherSans); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_other_sans"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("server_flag", configuration.ServerFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting server_flag"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("client_flag", configuration.ClientFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting client_flag"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("code_signing_flag", configuration.CodeSigningFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting code_signing_flag"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("email_protection_flag", configuration.EmailProtectionFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting email_protection_flag"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("key_type", configuration.KeyType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_type"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("key_bits", flex.IntValue(configuration.KeyBits)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_bits"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.KeyUsage != nil {
		if err = d.Set("key_usage", configuration.KeyUsage); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_usage"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.ExtKeyUsage != nil {
		if err = d.Set("ext_key_usage", configuration.ExtKeyUsage); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ext_key_usage"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.ExtKeyUsageOids != nil {
		if err = d.Set("ext_key_usage_oids", configuration.ExtKeyUsageOids); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ext_key_usage_oids"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("use_csr_common_name", configuration.UseCsrCommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting use_csr_common_name"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("use_csr_sans", configuration.UseCsrSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting use_csr_sans"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.Ou != nil {
		if err = d.Set("ou", configuration.Ou); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ou"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Organization != nil {
		if err = d.Set("organization", configuration.Organization); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting organization"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Country != nil {
		if err = d.Set("country", configuration.Country); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting country"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Locality != nil {
		if err = d.Set("locality", configuration.Locality); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting locality"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.Province != nil {
		if err = d.Set("province", configuration.Province); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting province"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.StreetAddress != nil {
		if err = d.Set("street_address", configuration.StreetAddress); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting street_address"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if configuration.PostalCode != nil {
		if err = d.Set("postal_code", configuration.PostalCode); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting postal_code"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("require_cn", configuration.RequireCn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting require_cn"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if configuration.PolicyIdentifiers != nil {
		if err = d.Set("policy_identifiers", configuration.PolicyIdentifiers); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting policy_identifiers"), PrivateCertConfigTemplateResourceName, "read")
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("basic_constraints_valid_for_non_ca", configuration.BasicConstraintsValidForNonCa); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting basic_constraints_valid_for_non_ca"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}
	if err = d.Set("not_before_duration_seconds", flex.IntValue(configuration.NotBeforeDurationSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting not_before_duration_seconds"), PrivateCertConfigTemplateResourceName, "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmSmPrivateCertificateConfigurationTemplateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigTemplateResourceName, "update")
		return tfErr.GetDiag()
	}

	id := strings.Split(d.Id(), "/")
	region := id[0]
	instanceId := id[1]
	configName := id[2]
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	updateConfigurationOptions := &secretsmanagerv2.UpdateConfigurationOptions{}

	updateConfigurationOptions.SetName(configName)
	updateConfigurationOptions.SetXSmAcceptConfigurationType("private_cert_configuration_template")

	hasChange := false

	patchVals := &secretsmanagerv2.ConfigurationPatch{}

	if d.HasChange("max_ttl") {
		patchVals.MaxTTL = core.StringPtr(d.Get("max_ttl").(string))
		hasChange = true
	}

	if d.HasChange("ttl") {
		patchVals.TTL = core.StringPtr(d.Get("ttl").(string))
		hasChange = true
	}

	if d.HasChange("allowed_secret_groups") {
		patchVals.AllowedSecretGroups = core.StringPtr(d.Get("allowed_secret_groups").(string))
		hasChange = true
	}

	if d.HasChange("allow_localhost") {
		patchVals.AllowLocalhost = core.BoolPtr(d.Get("allow_localhost").(bool))
		hasChange = true
	}

	if d.HasChange("allowed_domains") {
		allowedDomains := []string{}
		for _, allowedDomainsItem := range d.Get("allowed_domains").([]interface{}) {
			allowedDomains = append(allowedDomains, allowedDomainsItem.(string))
		}
		patchVals.AllowedDomains = allowedDomains
		hasChange = true
	}

	if d.HasChange("allowed_domains_template") {
		patchVals.AllowedDomainsTemplate = core.BoolPtr(d.Get("allowed_domains_template").(bool))
		hasChange = true
	}

	if d.HasChange("allow_bare_domains") {
		patchVals.AllowBareDomains = core.BoolPtr(d.Get("allow_bare_domains").(bool))
		hasChange = true
	}

	if d.HasChange("allow_subdomains") {
		patchVals.AllowSubdomains = core.BoolPtr(d.Get("allow_subdomains").(bool))
		hasChange = true
	}

	if d.HasChange("allow_glob_domains") {
		patchVals.AllowGlobDomains = core.BoolPtr(d.Get("allow_glob_domains").(bool))
		hasChange = true
	}

	if d.HasChange("allow_any_name") {
		patchVals.AllowAnyName = core.BoolPtr(d.Get("allow_any_name").(bool))
		hasChange = true
	}

	if d.HasChange("enforce_hostnames") {
		patchVals.EnforceHostnames = core.BoolPtr(d.Get("enforce_hostnames").(bool))
		hasChange = true
	}

	if d.HasChange("allow_ip_sans") {
		patchVals.AllowIpSans = core.BoolPtr(d.Get("allow_ip_sans").(bool))
		hasChange = true
	}

	if d.HasChange("allowed_uri_sans") {
		allowedUriSans := []string{}
		for _, allowedUriSansItem := range d.Get("allowed_uri_sans").([]interface{}) {
			allowedUriSans = append(allowedUriSans, allowedUriSansItem.(string))
		}
		patchVals.AllowedUriSans = allowedUriSans
		hasChange = true
	}

	if d.HasChange("allowed_other_sans") {
		allowedOtherSans := []string{}
		for _, allowedOtherSansItem := range d.Get("allowed_other_sans").([]interface{}) {
			allowedOtherSans = append(allowedOtherSans, allowedOtherSansItem.(string))
		}
		patchVals.AllowedOtherSans = allowedOtherSans
		hasChange = true
	}

	if d.HasChange("server_flag") {
		patchVals.ServerFlag = core.BoolPtr(d.Get("server_flag").(bool))
		hasChange = true
	}

	if d.HasChange("client_flag") {
		patchVals.ClientFlag = core.BoolPtr(d.Get("client_flag").(bool))
		hasChange = true
	}

	if d.HasChange("code_signing_flag") {
		patchVals.CodeSigningFlag = core.BoolPtr(d.Get("code_signing_flag").(bool))
		hasChange = true
	}

	if d.HasChange("email_protection_flag") {
		patchVals.EmailProtectionFlag = core.BoolPtr(d.Get("email_protection_flag").(bool))
		hasChange = true
	}

	if d.HasChange("key_type") {
		patchVals.KeyType = core.StringPtr(d.Get("key_type").(string))
		hasChange = true
	}

	if d.HasChange("key_bits") {
		patchVals.KeyBits = core.Int64Ptr(int64(d.Get("key_bits").(int)))
		hasChange = true
	}

	if d.HasChange("key_usage") {
		keyUsage := []string{}
		for _, keyUsageItem := range d.Get("key_usage").([]interface{}) {
			keyUsage = append(keyUsage, keyUsageItem.(string))
		}
		patchVals.KeyUsage = keyUsage
		hasChange = true
	}

	if d.HasChange("ext_key_usage") {
		extKeyUsage := []string{}
		for _, extKeyUsageItem := range d.Get("ext_key_usage").([]interface{}) {
			extKeyUsage = append(extKeyUsage, extKeyUsageItem.(string))
		}
		patchVals.ExtKeyUsage = extKeyUsage
		hasChange = true
	}

	if d.HasChange("ext_key_usage_oids") {
		extKeyUsageOids := []string{}
		for _, extKeyUsageOidsItem := range d.Get("ext_key_usage_oids").([]interface{}) {
			extKeyUsageOids = append(extKeyUsageOids, extKeyUsageOidsItem.(string))
		}
		patchVals.ExtKeyUsageOids = extKeyUsageOids
		hasChange = true
	}

	if d.HasChange("use_csr_common_name") {
		patchVals.UseCsrCommonName = core.BoolPtr(d.Get("use_csr_common_name").(bool))
		hasChange = true
	}

	if d.HasChange("use_csr_sans") {
		patchVals.UseCsrSans = core.BoolPtr(d.Get("use_csr_sans").(bool))
		hasChange = true
	}

	if d.HasChange("ou") {
		ou := []string{}
		for _, ouItem := range d.Get("ou").([]interface{}) {
			ou = append(ou, ouItem.(string))
		}
		patchVals.Ou = ou
		hasChange = true
	}

	if d.HasChange("organization") {
		organization := []string{}
		for _, organizationItem := range d.Get("organization").([]interface{}) {
			organization = append(organization, organizationItem.(string))
		}
		patchVals.Organization = organization
		hasChange = true
	}

	if d.HasChange("country") {
		country := []string{}
		for _, countryItem := range d.Get("country").([]interface{}) {
			country = append(country, countryItem.(string))
		}
		patchVals.Country = country
		hasChange = true
	}

	if d.HasChange("locality") {
		locality := []string{}
		for _, localityItem := range d.Get("locality").([]interface{}) {
			locality = append(locality, localityItem.(string))
		}
		patchVals.Locality = locality
		hasChange = true
	}

	if d.HasChange("province") {
		province := []string{}
		for _, provinceItem := range d.Get("province").([]interface{}) {
			province = append(province, provinceItem.(string))
		}
		patchVals.Province = province
		hasChange = true
	}

	if d.HasChange("street_address") {
		streetAddress := []string{}
		for _, streetAddressItem := range d.Get("street_address").([]interface{}) {
			streetAddress = append(streetAddress, streetAddressItem.(string))
		}
		patchVals.StreetAddress = streetAddress
		hasChange = true
	}

	if d.HasChange("postal_code") {
		postalCode := []string{}
		for _, postalCodeItem := range d.Get("postal_code").([]interface{}) {
			postalCode = append(postalCode, postalCodeItem.(string))
		}
		patchVals.PostalCode = postalCode
		hasChange = true
	}

	if d.HasChange("not_before_duration") {
		patchVals.NotBeforeDuration = core.StringPtr(d.Get("not_before_duration").(string))
		hasChange = true
	}

	if d.HasChange("require_cn") {
		patchVals.RequireCn = core.BoolPtr(d.Get("require_cn").(bool))
		hasChange = true
	}

	if d.HasChange("policy_identifiers") {
		policyIdentifiers := []string{}
		for _, policyIdentifiersItem := range d.Get("policy_identifiers").([]interface{}) {
			policyIdentifiers = append(policyIdentifiers, policyIdentifiersItem.(string))
		}
		patchVals.PolicyIdentifiers = policyIdentifiers
		hasChange = true
	}

	if d.HasChange("basic_constraints_valid_for_non_ca") {
		patchVals.BasicConstraintsValidForNonCa = core.BoolPtr(d.Get("basic_constraints_valid_for_non_ca").(bool))
		hasChange = true
	}

	if hasChange {
		updateConfigurationOptions.ConfigurationPatch, _ = patchVals.AsPatch()
		_, response, err := secretsManagerClient.UpdateConfigurationWithContext(context, updateConfigurationOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigurationWithContext failed %s\n%s", err, response)
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigTemplateResourceName, "update")
			return tfErr.GetDiag()
		}
	}

	return resourceIbmSmPrivateCertificateConfigurationTemplateRead(context, d, meta)
}

func resourceIbmSmPrivateCertificateConfigurationTemplateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigTemplateResourceName, "delete")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigurationWithContext failed %s\n%s", err, response), PrivateCertConfigTemplateResourceName, "delete")
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func resourceIbmSmPrivateCertificateConfigurationTemplateMapToConfigurationPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationPrototypeIntf, error) {
	model := &secretsmanagerv2.PrivateCertificateConfigurationTemplatePrototype{}

	model.ConfigType = core.StringPtr("private_cert_configuration_template")

	if _, ok := d.GetOk("name"); ok {
		model.Name = core.StringPtr(d.Get("name").(string))
	}
	if _, ok := d.GetOk("certificate_authority"); ok {
		model.CertificateAuthority = core.StringPtr(d.Get("certificate_authority").(string))
	}
	if _, ok := d.GetOk("allowed_secret_groups"); ok {
		model.AllowedSecretGroups = core.StringPtr(d.Get("allowed_secret_groups").(string))
	}
	if _, ok := d.GetOk("max_ttl"); ok {
		model.MaxTTL = core.StringPtr(d.Get("max_ttl").(string))
	}
	if _, ok := d.GetOk("ttl"); ok {
		model.TTL = core.StringPtr(d.Get("ttl").(string))
	}
	if _, ok := d.GetOkExists("allow_localhost"); ok {
		model.AllowLocalhost = core.BoolPtr(d.Get("allow_localhost").(bool))
	}
	if _, ok := d.GetOk("allowed_domains"); ok {
		allowedDomains := []string{}
		for _, allowedDomainsItem := range d.Get("allowed_domains").([]interface{}) {
			allowedDomains = append(allowedDomains, allowedDomainsItem.(string))
		}
		model.AllowedDomains = allowedDomains
	}
	if _, ok := d.GetOk("allowed_domains_template"); ok {
		model.AllowedDomainsTemplate = core.BoolPtr(d.Get("allowed_domains_template").(bool))
	}
	if _, ok := d.GetOk("allow_bare_domains"); ok {
		model.AllowBareDomains = core.BoolPtr(d.Get("allow_bare_domains").(bool))
	}
	if _, ok := d.GetOk("allow_subdomains"); ok {
		model.AllowSubdomains = core.BoolPtr(d.Get("allow_subdomains").(bool))
	}
	if _, ok := d.GetOk("allow_glob_domains"); ok {
		model.AllowGlobDomains = core.BoolPtr(d.Get("allow_glob_domains").(bool))
	}
	if _, ok := d.GetOk("allow_wildcard_certificates"); ok {
		model.AllowWildcardCertificates = core.BoolPtr(d.Get("allow_wildcard_certificates").(bool))
	}
	if _, ok := d.GetOk("allow_any_name"); ok {
		model.AllowAnyName = core.BoolPtr(d.Get("allow_any_name").(bool))
	}
	if _, ok := d.GetOkExists("enforce_hostnames"); ok {
		model.EnforceHostnames = core.BoolPtr(d.Get("enforce_hostnames").(bool))
	}
	if _, ok := d.GetOkExists("allow_ip_sans"); ok {
		model.AllowIpSans = core.BoolPtr(d.Get("allow_ip_sans").(bool))
	}
	if _, ok := d.GetOk("allowed_uri_sans"); ok {
		allowedUriSans := []string{}
		for _, allowedUriSansItem := range d.Get("allowed_uri_sans").([]interface{}) {
			allowedUriSans = append(allowedUriSans, allowedUriSansItem.(string))
		}
		model.AllowedUriSans = allowedUriSans
	}
	if _, ok := d.GetOk("allowed_other_sans"); ok {
		allowedOtherSans := []string{}
		for _, allowedOtherSansItem := range d.Get("allowed_other_sans").([]interface{}) {
			allowedOtherSans = append(allowedOtherSans, allowedOtherSansItem.(string))
		}
		model.AllowedOtherSans = allowedOtherSans
	}
	if _, ok := d.GetOkExists("server_flag"); ok {
		model.ServerFlag = core.BoolPtr(d.Get("server_flag").(bool))
	}
	if _, ok := d.GetOkExists("client_flag"); ok {
		model.ClientFlag = core.BoolPtr(d.Get("client_flag").(bool))
	}
	if _, ok := d.GetOk("code_signing_flag"); ok {
		model.CodeSigningFlag = core.BoolPtr(d.Get("code_signing_flag").(bool))
	}
	if _, ok := d.GetOk("email_protection_flag"); ok {
		model.EmailProtectionFlag = core.BoolPtr(d.Get("email_protection_flag").(bool))
	}
	if _, ok := d.GetOk("key_type"); ok {
		model.KeyType = core.StringPtr(d.Get("key_type").(string))
	}
	if _, ok := d.GetOk("key_bits"); ok {
		model.KeyBits = core.Int64Ptr(int64(d.Get("key_bits").(int)))
	}
	if _, ok := d.GetOk("key_usage"); ok {
		keyUsage := []string{}
		for _, keyUsageItem := range d.Get("key_usage").([]interface{}) {
			keyUsage = append(keyUsage, keyUsageItem.(string))
		}
		model.KeyUsage = keyUsage
	}
	if _, ok := d.GetOk("ext_key_usage"); ok {
		extKeyUsage := []string{}
		for _, extKeyUsageItem := range d.Get("ext_key_usage").([]interface{}) {
			extKeyUsage = append(extKeyUsage, extKeyUsageItem.(string))
		}
		model.ExtKeyUsage = extKeyUsage
	}
	if _, ok := d.GetOk("ext_key_usage_oids"); ok {
		extKeyUsageOids := []string{}
		for _, extKeyUsageOidsItem := range d.Get("ext_key_usage_oids").([]interface{}) {
			extKeyUsageOids = append(extKeyUsageOids, extKeyUsageOidsItem.(string))
		}
		model.ExtKeyUsageOids = extKeyUsageOids
	}
	if _, ok := d.GetOkExists("use_csr_common_name"); ok {
		model.UseCsrCommonName = core.BoolPtr(d.Get("use_csr_common_name").(bool))
	}
	if _, ok := d.GetOkExists("use_csr_sans"); ok {
		model.UseCsrSans = core.BoolPtr(d.Get("use_csr_sans").(bool))
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
	if _, ok := d.GetOkExists("require_cn"); ok {
		model.RequireCn = core.BoolPtr(d.Get("require_cn").(bool))
	}
	if _, ok := d.GetOk("policy_identifiers"); ok {
		policyIdentifiers := []string{}
		for _, policyIdentifiersItem := range d.Get("policy_identifiers").([]interface{}) {
			policyIdentifiers = append(policyIdentifiers, policyIdentifiersItem.(string))
		}
		model.PolicyIdentifiers = policyIdentifiers
	}
	if _, ok := d.GetOk("basic_constraints_valid_for_non_ca"); ok {
		model.BasicConstraintsValidForNonCa = core.BoolPtr(d.Get("basic_constraints_valid_for_non_ca").(bool))
	}
	if _, ok := d.GetOk("not_before_duration"); ok {
		model.NotBeforeDuration = core.StringPtr(d.Get("not_before_duration").(string))
	}

	return model, nil
}
