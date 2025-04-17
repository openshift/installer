// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisCustomCertificates = "custom_certificates"
)

func DataSourceIBMCISCustomCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMCISCustomCertificatesRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_custom_certificates",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisCustomCertificates: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						cisCertificateUploadCustomCertID: {
							Type:     schema.TypeString,
							Computed: true,
						},
						cisCertificateUploadBundleMethod: {
							Type:        schema.TypeString,
							Description: "Certificate bundle method",
							Computed:    true,
						},
						cisCertificateUploadHosts: {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "hosts which the certificate uploaded to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						cisPageRulePriority: {
							Type:        schema.TypeInt,
							Description: "Certificate priority",
							Computed:    true,
						},
						cisCertificateUploadStatus: {
							Type:        schema.TypeString,
							Description: "certificate status",
							Computed:    true,
						},
						cisCertificateUploadIssuer: {
							Type:        schema.TypeString,
							Description: "certificate issuer",
							Computed:    true,
						},
						cisCertificateUploadSignature: {
							Type:        schema.TypeString,
							Description: "certificate signature",
							Computed:    true,
						},
						cisCertificateUploadUploadedOn: {
							Type:        schema.TypeString,
							Description: "certificate uploaded date",
							Computed:    true,
						},
						cisCertificateUploadModifiedOn: {
							Type:        schema.TypeString,
							Description: "certificate modified date",
							Computed:    true,
						},
						cisCertificateUploadExpiresOn: {
							Type:        schema.TypeString,
							Description: "certificate expires date",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMCISCustomCertificatesValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISCustomCertificatesValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_custom_certificates",
		Schema:       validateSchema}
	return &iBMCISCustomCertificatesValidator
}

func dataSourceIBMCISCustomCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewListCustomCertificatesOptions()
	result, resp, err := cisClient.ListCustomCertificates(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to list custom certificates: %v", resp)
	}
	certsList := make([]map[string]interface{}, 0)
	for _, r := range result.Result {
		cert := map[string]interface{}{}
		cert["id"] = flex.ConvertCisToTfThreeVar(*r.ID, zoneID, crn)
		cert[cisCertificateUploadCustomCertID] = *r.ID
		cert[cisCertificateUploadBundleMethod] = *r.BundleMethod
		cert[cisCertificateUploadHosts] = flex.FlattenStringList(r.Hosts)
		cert[cisCertificateUploadIssuer] = *r.Issuer
		cert[cisCertificateUploadSignature] = *r.Signature
		cert[cisCertificateUploadStatus] = *r.Status
		cert[cisCertificateUploadPriority] = *r.Priority
		cert[cisCertificateUploadUploadedOn] = *r.UploadedOn
		cert[cisCertificateUploadModifiedOn] = *r.ModifiedOn
		cert[cisCertificateUploadExpiresOn] = *r.ExpiresOn
		certsList = append(certsList, cert)
	}
	d.SetId(dataSourceIBMCISCustomCertificatesID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisCustomCertificates, certsList)
	return nil
}

func dataSourceIBMCISCustomCertificatesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
