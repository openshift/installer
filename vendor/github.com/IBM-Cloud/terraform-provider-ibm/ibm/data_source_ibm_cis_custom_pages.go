// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisCustomPages = "cis_custom_pages"
)

func dataSourceIBMCISCustomPages() *schema.Resource {
	return &schema.Resource{
		Read:     dataSourceIBMCISCustomPagesRead,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:     schema.TypeString,
				Required: true,
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

func dataSourceIBMCISCustomPagesRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisCustomPageClientSession()
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
		customPage := make(map[string]interface{}, 0)
		customPage[cisCustomPageIdentifier] = *instance.ID
		customPage[cisCustomPageState] = *instance.State
		customPage[cisCustomPageDesc] = *instance.Description
		customPage[cisCustomPagePreviewTarget] = *instance.PreviewTarget
		customPage[cisCustomPageRequiredTokens] = flattenStringList(instance.RequiredTokens)
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
