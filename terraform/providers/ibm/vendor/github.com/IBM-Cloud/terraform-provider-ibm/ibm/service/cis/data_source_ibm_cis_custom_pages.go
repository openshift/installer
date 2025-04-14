// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisCustomPages = "cis_custom_pages"
)

func DataSourceIBMCISCustomPages() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMCISCustomPagesRead,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_custom_pages",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisCustomPages: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisCustomPageIdentifier: {
							Type:        schema.TypeString,
							Description: "Custom page identifier",
							Computed:    true,
						},
						cisCustomPageURL: {
							Type:        schema.TypeString,
							Description: "Custom page url",
							Computed:    true,
						},
						cisCustomPageState: {
							Type:        schema.TypeString,
							Description: "Custom page state",
							Computed:    true,
						},
						cisCustomPageDesc: {
							Type:        schema.TypeString,
							Description: "Free text",
							Computed:    true,
						},
						cisCustomPageRequiredTokens: {
							Type:        schema.TypeList,
							Description: "Custom page state",
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						cisCustomPagePreviewTarget: {
							Type:        schema.TypeString,
							Description: "Custom page preview target",
							Computed:    true,
						},
						cisCustomPageCreatedOn: {
							Type:        schema.TypeString,
							Description: "Custom page created date",
							Computed:    true,
						},
						cisCustomPageModifiedOn: {
							Type:        schema.TypeString,
							Description: "Custom page modified date",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
func DataSourceIBMCISCustomPagesValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISCustomPagesValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_custom_pages",
		Schema:       validateSchema}
	return &iBMCISCustomPagesValidator
}
func dataSourceIBMCISCustomPagesRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisCustomPageClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID := d.Get(cisDomainID).(string)
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewListZoneCustomPagesOptions()

	result, response, err := cisClient.ListZoneCustomPages(opt)
	if err != nil {
		log.Printf("List custom pages failed: %v", response)
		return err
	}
	customPagesOutput := make([]map[string]interface{}, 0)
	for _, instance := range result.Result {
		customPage := make(map[string]interface{})
		customPage[cisCustomPageIdentifier] = *instance.ID
		customPage[cisCustomPageState] = *instance.State
		customPage[cisCustomPageDesc] = *instance.Description
		customPage[cisCustomPagePreviewTarget] = *instance.PreviewTarget
		customPage[cisCustomPageRequiredTokens] = flex.FlattenStringList(instance.RequiredTokens)
		if instance.CreatedOn != nil {
			customPage[cisCustomPageCreatedOn] = (*instance.CreatedOn).String()
		}
		if instance.ModifiedOn != nil {
			customPage[cisCustomPageModifiedOn] = (*instance.ModifiedOn).String()
		}
		if instance.URL != nil {
			customPage[cisCustomPageURL] = *instance.URL
		}

		customPagesOutput = append(customPagesOutput, customPage)
	}
	d.SetId(dataSourceIBMCISCustomPageID(d))
	d.Set(cisCustomPages, customPagesOutput)
	return nil
}

func dataSourceIBMCISCustomPageID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
