// Copyright IBM Corp. 2017, 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisMtlsOutput                  = "mtls_certificates"
	cisMtlsCertID                  = "id"
	cisMtlsCertName                = "name"
	cisMtlsCertFingerprint         = "fingerprint"
	cisMtlsCertAssociatedHostnames = "associated_hostnames"
	cisMtlsCertCreatedAt           = "created_at"
	cisMtlsCertUpdatedAt           = "updated_at"
	cisMtlsCertExpiresOn           = "expires_on"
)

func DataSourceIBMCISMtls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataIBMCISMtlsRead,

		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_mtlss",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisMtlsOutput: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container for response information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisMtlsCertID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate ID",
						},
						cisMtlsCertName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate name",
						},
						cisMtlsCertFingerprint: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate Fingerprint",
						},
						cisMtlsCertAssociatedHostnames: {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Certificate Associated Hostnames",
						},
						cisMtlsCertCreatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate Created At",
						},
						cisMtlsCertUpdatedAt: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate Updated At",
						},
						cisMtlsCertExpiresOn: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate Expires On",
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMCISMtlsValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISMtlsValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_mtlss",
		Schema:       validateSchema}
	return &iBMCISMtlsValidator
}

func dataIBMCISMtlsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).CisMtlsSession()

	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error while getting the CisMtlsSession() %s %v", err, sess))
	}

	zoneID, crn, _ := flex.ConvertTftoCisTwoVar(d.Id())
	sess.Crn = core.StringPtr(crn)

	opt := sess.NewListAccessCertificatesOptions(zoneID)

	result, resp, err := sess.ListAccessCertificates(opt)
	if err != nil {
		log.Printf("[WARN] List all certificates failed: %v\n", resp)
		return diag.FromErr(err)
	}
	mtlsCertLists := make([]map[string]interface{}, 0)
	for _, certObj := range result.Result {
		mtlsCertList := map[string]interface{}{}
		mtlsCertList[cisMtlsCertID] = *certObj.ID
		mtlsCertList[cisMtlsCertName] = *certObj.Name
		mtlsCertList[cisMtlsCertFingerprint] = *certObj.Fingerprint
		mtlsCertList[cisMtlsCertAssociatedHostnames] = certObj.AssociatedHostnames
		mtlsCertList[cisMtlsCertCreatedAt] = *certObj.CreatedAt
		mtlsCertList[cisMtlsCertUpdatedAt] = *certObj.UpdatedAt
		mtlsCertList[cisMtlsCertExpiresOn] = *certObj.ExpiresOn

		mtlsCertLists = append(mtlsCertLists, mtlsCertList)

	}
	d.SetId(dataSourceCISMtlsCheckID(d))

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisMtlsOutput, mtlsCertLists)

	return nil
}
func dataSourceCISMtlsCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
