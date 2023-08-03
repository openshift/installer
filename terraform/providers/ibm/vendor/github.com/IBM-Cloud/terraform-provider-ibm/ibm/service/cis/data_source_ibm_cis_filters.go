// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const cisFiltersList = "cis_filters_list"

func DataSourceIBMCISFilters() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISFiltersRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_filters",
					"cis_id"),
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
func DataSourceIBMCISFiltersValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	iBMCISFiltersValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_filters",
		Schema:       validateSchema}
	return &iBMCISFiltersValidator
}
func dataIBMCISFiltersRead(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(conns.ClientSession).CisFiltersSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisFiltersSession %s", err)
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))

	result, resp, err := cisClient.ListAllFilters(cisClient.NewListAllFiltersOptions(xAuthtoken, crn, zoneID))
	if err != nil || result == nil {
		return fmt.Errorf("[ERROR] Error Listing all filters %q: %s %s", d.Id(), err, resp)
	}

	filtersList := make([]map[string]interface{}, 0)

	for _, filtersObj := range result.Result {
		filtersOutput := map[string]interface{}{}
		filtersOutput[cisFilterID] = *filtersObj.ID
		filtersOutput[cisFilterDescription] = filtersObj.Description
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
