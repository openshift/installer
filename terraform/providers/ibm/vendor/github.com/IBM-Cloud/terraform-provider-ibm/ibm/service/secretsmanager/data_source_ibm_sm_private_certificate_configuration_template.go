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

func DataSourceIbmSmPrivateCertificateConfigurationTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSmPrivateCertificateConfigurationTemplateRead,

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
			"certificate_authority": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the intermediate certificate authority.",
			},
			"allowed_secret_groups": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scopes the creation of private certificates to only the secret groups that you specify.This field can be supplied as a comma-delimited list of secret group IDs.",
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
			"allow_localhost": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to allow `localhost` to be included as one of the requested common names.",
			},
			"allowed_domains": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The domains to define for the certificate template. This property is used along with the `allow_bare_domains` and `allow_subdomains` options.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_domains_template": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to allow the domains that are supplied in the `allowed_domains` field to contain access control list (ACL) templates.",
			},
			"allow_bare_domains": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to allow clients to request private certificates that match the value of the actual domains on the final certificate.For example, if you specify `example.com` in the `allowed_domains` field, you grant clients the ability to request a certificate that contains the name `example.com` as one of the DNS values on the final certificate.**Important:** In some scenarios, allowing bare domains can be considered a security risk.",
			},
			"allow_subdomains": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to allow clients to request private certificates with common names (CN) that are subdomains of the CNs that are allowed by the other certificate template options. This includes wildcard subdomains.For example, if `allowed_domains` has a value of `example.com` and `allow_subdomains`is set to `true`, then the following subdomains are allowed: `foo.example.com`, `bar.example.com`, `*.example.com`.**Note:** This field is redundant if you use the `allow_any_name` option.",
			},
			"allow_glob_domains": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to allow glob patterns, for example, `ftp*.example.com`, in the names that are specified in the `allowed_domains` field.If set to `true`, clients are allowed to request private certificates with names that match the glob patterns.",
			},
			"allow_any_name": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to allow clients to request a private certificate that matches any common name.",
			},
			"enforce_hostnames": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to enforce only valid host names for common names, DNS Subject Alternative Names, and the host section of email addresses.",
			},
			"allow_ip_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to allow clients to request a private certificate with IP Subject Alternative Names.",
			},
			"allowed_uri_sans": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The URI Subject Alternative Names to allow for private certificates.Values can contain glob patterns, for example `spiffe://hostname/_*`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_other_sans": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names (SANs) to allow for private certificates.The format for each element in the list is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`. To allow any value for an OID, use `*` as its value. Alternatively, specify a single `*` to allow any `other_sans` input.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"server_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for server use.",
			},
			"client_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for client use.",
			},
			"code_signing_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for code signing use.",
			},
			"email_protection_flag": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether private certificates are flagged for email protection use.",
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
			"key_usage": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The allowed key usage constraint to define for private certificates.You can find valid values in the [Go x509 package documentation](https://pkg.go.dev/crypto/x509#KeyUsage). Omit the `KeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ext_key_usage": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The allowed extended key usage constraint on private certificates.You can find valid values in the [Go x509 package documentation](https://golang.org/pkg/crypto/x509/#ExtKeyUsage). Omit the `ExtKeyUsage` part of the value. Values are not case-sensitive. To specify no key usage constraints, set this field to an empty list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ext_key_usage_oids": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of extended key usage Object Identifiers (OIDs).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"use_csr_common_name": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the common name (CN) from a certificate signing request (CSR) instead of the CN that's included in the data of the certificate.Does not include any requested Subject Alternative Names (SANs) in the CSR. To use the alternative names, include the `use_csr_sans` property.",
			},
			"use_csr_sans": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When used with the `private_cert_configuration_action_sign_csr` action, this field determines whether to use the Subject Alternative Names(SANs) from a certificate signing request (CSR) instead of the SANs that are included in the data of the certificate.Does not include the common name in the CSR. To use the common name, include the `use_csr_common_name` property.",
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
			"require_cn": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to require a common name to create a private certificate.By default, a common name is required to generate a certificate. To make the `common_name` field optional, set the `require_cn` option to `false`.",
			},
			"policy_identifiers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of policy Object Identifiers (OIDs).",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"basic_constraints_valid_for_non_ca": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Determines whether to mark the Basic Constraints extension of an issued private certificate as valid for non-CA certificates.",
			},
			"not_before_duration_seconds": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration in seconds by which to backdate the `not_before` property of an issued private certificate.",
			},
		},
	}
}

func dataSourceIbmSmPrivateCertificateConfigurationTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	getConfigurationOptions := &secretsmanagerv2.GetConfigurationOptions{}

	getConfigurationOptions.SetName(d.Get("name").(string))

	privateCertificateConfigurationTemplateIntf, response, err := secretsManagerClient.GetConfigurationWithContext(context, getConfigurationOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigurationWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigurationWithContext failed %s\n%s", err, response), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}
	privateCertificateConfigurationTemplate := privateCertificateConfigurationTemplateIntf.(*secretsmanagerv2.PrivateCertificateConfigurationTemplate)

	d.SetId(fmt.Sprintf("%s/%s/%s", region, instanceId, *getConfigurationOptions.Name))

	if err = d.Set("region", region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting region"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("config_type", privateCertificateConfigurationTemplate.ConfigType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting config_type"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("secret_type", privateCertificateConfigurationTemplate.SecretType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting secret_type"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", privateCertificateConfigurationTemplate.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_by"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", DateTimeToRFC3339(privateCertificateConfigurationTemplate.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting created_at"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", DateTimeToRFC3339(privateCertificateConfigurationTemplate.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting updated_at"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("certificate_authority", privateCertificateConfigurationTemplate.CertificateAuthority); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting certificate_authority"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allowed_secret_groups", privateCertificateConfigurationTemplate.AllowedSecretGroups); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_secret_groups"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("max_ttl_seconds", flex.IntValue(privateCertificateConfigurationTemplate.MaxTtlSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting max_ttl_seconds"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("ttl_seconds", flex.IntValue(privateCertificateConfigurationTemplate.TtlSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting ttl_seconds"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allow_localhost", privateCertificateConfigurationTemplate.AllowLocalhost); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_localhost"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allowed_domains_template", privateCertificateConfigurationTemplate.AllowedDomainsTemplate); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_domains_template"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allow_bare_domains", privateCertificateConfigurationTemplate.AllowBareDomains); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_bare_domains"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allow_subdomains", privateCertificateConfigurationTemplate.AllowSubdomains); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_subdomains"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allow_glob_domains", privateCertificateConfigurationTemplate.AllowGlobDomains); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_glob_domains"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allow_any_name", privateCertificateConfigurationTemplate.AllowAnyName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_any_name"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("enforce_hostnames", privateCertificateConfigurationTemplate.EnforceHostnames); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting enforce_hostnames"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("allow_ip_sans", privateCertificateConfigurationTemplate.AllowIpSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allow_ip_sans"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("server_flag", privateCertificateConfigurationTemplate.ServerFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting server_flag"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("client_flag", privateCertificateConfigurationTemplate.ClientFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting client_flag"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("code_signing_flag", privateCertificateConfigurationTemplate.CodeSigningFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting code_signing_flag"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("email_protection_flag", privateCertificateConfigurationTemplate.EmailProtectionFlag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting email_protection_flag"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("key_type", privateCertificateConfigurationTemplate.KeyType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_type"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("key_bits", flex.IntValue(privateCertificateConfigurationTemplate.KeyBits)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting key_bits"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("use_csr_common_name", privateCertificateConfigurationTemplate.UseCsrCommonName); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting use_csr_common_name"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("use_csr_sans", privateCertificateConfigurationTemplate.UseCsrSans); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting use_csr_sans"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("serial_number", privateCertificateConfigurationTemplate.SerialNumber); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting serial_number"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("require_cn", privateCertificateConfigurationTemplate.RequireCn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting require_cn"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("basic_constraints_valid_for_non_ca", privateCertificateConfigurationTemplate.BasicConstraintsValidForNonCa); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting basic_constraints_valid_for_non_ca"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("not_before_duration_seconds", flex.IntValue(privateCertificateConfigurationTemplate.NotBeforeDurationSeconds)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting not_before_duration_seconds"), fmt.Sprintf("(Data) %s", PrivateCertConfigTemplateResourceName), "read")
		return tfErr.GetDiag()
	}

	return nil
}
