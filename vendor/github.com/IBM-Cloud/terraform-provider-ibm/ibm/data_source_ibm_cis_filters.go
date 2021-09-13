// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const cisFiltersList = "cis_filters_list"

func dataSourceIBMCISFilters() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISFiltersRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisFiltersList: {
				Type:        schema.TypeList,
				Description: "Collection of Filter detail",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisFilterPaused: {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Filter Paused",
						},
						cisFilterID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filter ID",
						},
						cisFilterExpression: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filter Expression",
						},
						cisFilterDescription: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Filter Description",
						},
					},
				},
			},
		},
	}
}
func dataIBMCISFiltersRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return fmt.Errorf("Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(ClientSession).CisFiltersSession()
	if err != nil {
		return fmt.Errorf("Error while getting the CisFiltersSession %s", err)
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := convertTftoCisTwoVar(d.Get(cisDomainID).(string))

	result, resp, err := cisClient.ListAllFilters(cisClient.NewListAllFiltersOptions(xAuthtoken, crn, zoneID))
	if err != nil || result == nil {
		return fmt.Errorf("Error Listing all filters %q: %s %s", d.Id(), err, resp)
	}

	filters := result.Result

	filtersList := make([]map[string]interface{}, 0)
	for _, filtersObj := range filters {
		filtersOutput := map[string]interface{}{}
		filtersOutput[cisFilterID] = *filtersObj.ID
		filtersOutput[cisFilterDescription] = *filtersObj.Description
		filtersOutput[cisFilterExpression] = *filtersObj.Expression
		filtersOutput[cisFilterPaused] = *filtersObj.Paused
		filtersList = append(filtersList, filtersOutput)
	}
	d.SetId(dataSourceCISFiltersCheckID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisFiltersList, filtersList)
	return nil
}

func dataSourceCISFiltersCheckID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
