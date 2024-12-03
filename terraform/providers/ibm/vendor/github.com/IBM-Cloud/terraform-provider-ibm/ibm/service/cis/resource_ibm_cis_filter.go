// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/networking-go-sdk/filtersv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCISFilters        = "ibm_cis_filter"
	cisFilterExpression  = "expression"
	cisFilterPaused      = "paused"
	cisFilterDescription = "description"
	cisFilterID          = "filter_id"
)

func ResourceIBMCISFilter() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISFilterCreate,
		Read:     ResourceIBMCISFilterRead,
		Update:   ResourceIBMCISFilterUpdate,
		Delete:   ResourceIBMCISFilterDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_filter",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisFilterPaused: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter Paused",
			},
			cisFilterID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Filter ID",
			},
			cisFilterExpression: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Filter Expression",
			},
			cisFilterDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter Description",
			},
		},
	}
}
func ResourceIBMCISFilterCreate(d *schema.ResourceData, meta interface{}) error {
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

	var newfilter filtersv1.FilterInput

	if p, ok := d.GetOkExists(cisFilterPaused); ok {
		paused := p.(bool)
		newfilter.Paused = &paused
	}
	if des, ok := d.GetOk(cisFilterDescription); ok {
		description := des.(string)
		newfilter.Description = &description
	}
	if e, ok := d.GetOk(cisFilterExpression); ok {
		expression := e.(string)
		newfilter.Expression = &expression
	}

	opt := cisClient.NewCreateFilterOptions(xAuthtoken, crn, zoneID)

	opt.SetFilterInput([]filtersv1.FilterInput{newfilter})

	result, resp, err := cisClient.CreateFilter(opt)
	if err != nil || result == nil {
		return fmt.Errorf("[ERROR] Error creating Filter for zone %q: %s %s", zoneID, err, resp)
	}
	d.SetId(flex.ConvertCisToTfThreeVar(*result.Result[0].ID, zoneID, crn))
	return ResourceIBMCISFilterRead(d, meta)

}
func ResourceIBMCISFilterRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(conns.ClientSession).CisFiltersSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisFiltersSession %s", err)
	}
	filterid, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	opt := cisClient.NewGetFilterOptions(xAuthtoken, crn, zoneID, filterid)

	result, response, err := cisClient.GetFilter(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Error GetFilter not found ")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("[ERROR] Error finding GetFilter %q: %s %s", d.Id(), err, response)
	}
	if result.Result != nil {
		d.Set(cisID, crn)
		d.Set(cisDomainID, zoneID)
		d.Set(cisFilterID, result.Result.ID)
		d.Set(cisFilterPaused, result.Result.Paused)
		d.Set(cisFilterDescription, result.Result.Description)
		d.Set(cisFilterExpression, result.Result.Expression)
	}
	return nil
}
func ResourceIBMCISFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while Getting IAM Access Token using BluemixSession %s", err)
	}
	xAuthtoken := sess.Config.IAMAccessToken

	cisClient, err := meta.(conns.ClientSession).CisFiltersSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisFiltersSession %s", err)
	}

	filterid, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}

	if d.HasChange(cisFilterExpression) ||
		d.HasChange(cisFilterPaused) ||
		d.HasChange(cisFilterDescription) {

		var updatefilter filtersv1.FilterUpdateInput
		updatefilter.ID = &filterid

		if p, ok := d.GetOkExists(cisFilterPaused); ok {
			paused := p.(bool)
			updatefilter.Paused = &paused
		}
		if des, ok := d.GetOk(cisFilterDescription); ok {
			description := des.(string)
			updatefilter.Description = &description
		}
		if e, ok := d.GetOk(cisFilterExpression); ok {
			expression := e.(string)
			updatefilter.Expression = &expression
		}

		opt := cisClient.NewUpdateFiltersOptions(xAuthtoken, crn, zoneID)

		opt.SetFilterUpdateInput([]filtersv1.FilterUpdateInput{updatefilter})

		result, resp, err := cisClient.UpdateFilters(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating Filter for zone %q: %s %s", zoneID, err, resp)
		}

		if *result.Result[0].ID == "" {
			return fmt.Errorf("[ERROR] Error failed to find id in Update response; resource was empty")
		}
	}
	return ResourceIBMCISFilterRead(d, meta)
}
func ResourceIBMCISFilterDelete(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return err
	}
	xAuthtoken := sess.Config.IAMAccessToken
	cisClient, err := meta.(conns.ClientSession).CisFiltersSession()
	if err != nil {
		return err
	}
	filterid, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return err
	}
	opt := cisClient.NewDeleteFiltersOptions(xAuthtoken, crn, zoneID, filterid)
	_, _, err = cisClient.DeleteFilters(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error deleting Filter: %s", err)
	}

	return nil
}
func ResourceIBMCISFilterValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisFilterDescription,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "Filter-creation"})

	ibmCISFiltersResourceValidator := validate.ResourceValidator{ResourceName: ibmCISFilters, Schema: validateSchema}
	return &ibmCISFiltersResourceValidator
}
