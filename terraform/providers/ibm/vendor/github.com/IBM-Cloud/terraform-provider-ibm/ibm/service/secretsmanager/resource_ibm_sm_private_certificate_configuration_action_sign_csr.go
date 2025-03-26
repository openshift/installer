package secretsmanager

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
)

func ResourceIbmSmPrivateCertificateConfigurationActionSignCsr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmSmPrivateCertificateConfigurationActionSignCsrCreateOrUpdate,
		ReadContext:   resourceIbmSmPrivateCertificateConfigurationActionSignCsrRead,
		UpdateContext: resourceIbmSmPrivateCertificateConfigurationActionSignCsrCreateOrUpdate,
		DeleteContext: resourceIbmSmPrivateCertificateConfigurationActionSignCsrDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name that uniquely identifies a configuration",
			},
			"csr": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if removeNewLineFromCertificate(oldValue) == removeNewLineFromCertificate(newValue) {
						return true
					}
					return false
				},
				ForceNew:    true,
				Description: "The certificate signing request.",
			},
			"common_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The Common Name (AKA CN) represents the server name that is protected by the SSL certificate.",
			},
			"alt_names": &schema.Schema{
				Type:             schema.TypeList,
				ForceNew:         true,
				Optional:         true,
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
				ForceNew:    true,
				Description: "The custom Object Identifier (OID) or UTF8-string Subject Alternative Names to define for the CA certificate.The alternative names must match the values that are specified in the `allowed_other_sans` field in the associated certificate template. The format is the same as OpenSSL: `<oid>:<type>:<value>` where the current valid type is `UTF8`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The requested time-to-live (TTL) for certificates that are created by this CA. This field's value cannot be longer than the `max_ttl` limit.The value can be supplied as a string representation of a duration in hours, for example '8760h'. In the API response, this value is returned in seconds (integer).",
			},
			"format": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The format of the returned data.",
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
			"ou": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The Organizational Unit (OU) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"organization": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The Organization (O) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"country": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The Country (C) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"locality": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The Locality (L) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"province": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The Province (ST) values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"street_address": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The street address values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"postal_code": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The postal code values to define in the subject field of the resulting certificate.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"serial_number": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The serial number to assign to the generated certificate. To assign a random serial number, you can omit this field.",
			},
			"data": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "The configuration action data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"expiration": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The certificate expiration time.",
						},
					},
				},
			},
		},
	}
}

func resourceIbmSmPrivateCertificateConfigurationActionSignCsrCreateOrUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	secretsManagerClient, err := meta.(conns.ClientSession).SecretsManagerV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigActionSignCsr, "create/update")
		return tfErr.GetDiag()
	}

	region := getRegion(secretsManagerClient, d)
	instanceId := d.Get("instance_id").(string)
	secretsManagerClient = getClientWithInstanceEndpoint(secretsManagerClient, instanceId, region, getEndpointType(secretsManagerClient, d))

	createConfigurationActionOptions := &secretsmanagerv2.CreateConfigurationActionOptions{}

	configurationActionPrototypeModel, err := resourceIbmSmPrivateCertificateConfigurationActionSignCsrMapToConfigurationActionPrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigActionSignCsr, "create/update")
		return tfErr.GetDiag()
	}
	createConfigurationActionOptions.SetConfigActionPrototype(configurationActionPrototypeModel)
	createConfigurationActionOptions.SetName(d.Get("name").(string))

	configurationActionIntf, response, err := secretsManagerClient.CreateConfigurationActionWithContext(context, createConfigurationActionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigurationActionWithContext failed %s\n%s", err, response)
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigurationActionWithContext failed: %s\n%s", err.Error(), response), PrivateCertConfigActionSignCsr, "create/update")
		return tfErr.GetDiag()
	}

	configurationAction := configurationActionIntf.(*secretsmanagerv2.PrivateCertificateConfigurationActionSignCSR)

	if configurationAction.Data != nil {
		dataMap, err := resourceIbmSmPrivateCertificateConfigurationActionSignCsrDataToMap(*configurationAction.Data)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, "", PrivateCertConfigActionSignCsr, "create/update")
			return tfErr.GetDiag()
		}
		if err = d.Set("data", []map[string]interface{}{dataMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting data: %s", err), PrivateCertConfigActionSignCsr, "create/update")
			return tfErr.GetDiag()
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/sign_csr", region, instanceId, d.Get("name").(string)))

	return nil
}

func resourceIbmSmPrivateCertificateConfigurationActionSignCsrRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIbmSmPrivateCertificateConfigurationActionSignCsrDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceIbmSmPrivateCertificateConfigurationActionSignCsrMapToConfigurationActionPrototype(d *schema.ResourceData) (secretsmanagerv2.ConfigurationActionPrototypeIntf, error) {
	model := &secretsmanagerv2.PrivateCertificateConfigurationActionSignCSRPrototype{
		ActionType: core.StringPtr("private_cert_configuration_action_sign_csr"),
	}
	if _, ok := d.GetOk("csr"); ok {
		model.Csr = core.StringPtr(formatCertificate(d.Get("csr").(string)))
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
	if _, ok := d.GetOk("use_csr_values"); ok {
		model.UseCsrValues = core.BoolPtr(d.Get("use_csr_values").(bool))
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
	if _, ok := d.GetOk("serial_number"); ok {
		model.SerialNumber = core.StringPtr(d.Get("serial_number").(string))
	}

	return model, nil
}

func resourceIbmSmPrivateCertificateConfigurationActionSignCsrDataToMap(model secretsmanagerv2.PrivateCertificateConfigurationCACertificate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})

	if model.Certificate != nil {
		modelMap["certificate"] = model.Certificate
	}
	if model.IssuingCa != nil {
		modelMap["issuing_ca"] = model.IssuingCa
	}
	if model.CaChain != nil {
		modelMap["ca_chain"] = model.CaChain
	}
	if model.Expiration != nil {
		modelMap["expiration"] = flex.IntValue(model.Expiration)
	}
	return modelMap, nil
}
